use crate::hive::Hive;
use crate::messages::{CreateHive, NewSensorData};
use crate::types::Sensor;

use actix::prelude::*;
use futures_util::stream::StreamExt;
use lapin::{options::*, types::FieldTable, Channel, Connection, ConnectionProperties};
use log::{debug, error, info, warn};
use serde_json::Value;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::Mutex;

pub struct Apiary {
    pub rabbitmq_url: String,
    pub rabbitmq_channel: Arc<Mutex<Channel>>,
    pub hives: HashMap<i32, Addr<Hive>>,
}

impl Apiary {
    pub async fn new(rabbitmq_url: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let conn = Connection::connect(rabbitmq_url, ConnectionProperties::default()).await?;
        let channel = conn.create_channel().await?;

        // Declare queues with persistence
        for queue in &["sensor_queue", "hive_queue", "sensor_reading_queue"] {
            channel
                .queue_declare(
                    queue,
                    QueueDeclareOptions {
                        durable: true,
                        auto_delete: false,
                        exclusive: false,
                        ..Default::default()
                    },
                    FieldTable::default(),
                )
                .await?;
        }

        Ok(Apiary {
            rabbitmq_channel: Arc::new(Mutex::new(channel)),
            hives: HashMap::new(),
            rabbitmq_url: rabbitmq_url.to_string(),
        })
    }
}

impl Actor for Apiary {
    type Context = Context<Self>;

    fn started(&mut self, _: &mut Context<Self>) {
        info!("Apiary started");
    }
}

#[derive(Message)]
#[rtype(result = "()")]
pub struct StartConsumers;

impl Handler<StartConsumers> for Apiary {
    type Result = ();

    fn handle(&mut self, _msg: StartConsumers, ctx: &mut Context<Self>) {
        info!("Starting regular consumers...");
        self.start_consumers(ctx);
    }
}

#[derive(Message)]
#[rtype(result = "Result<(), String>")]
pub struct LoadExistingHives;

impl Handler<LoadExistingHives> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), String>>;

    fn handle(&mut self, _msg: LoadExistingHives, _ctx: &mut Context<Self>) -> Self::Result {
        let channel = self.rabbitmq_channel.clone();

        Box::pin(
            async move {
                let mut stored_hives: Vec<(i32, Vec<Sensor>)> = Vec::new();
                let channel = channel.lock().await;

                info!("Starting initial load of hives from queue...");
                let queue = channel
                    .queue_declare(
                        "hive_queue",
                        QueueDeclareOptions {
                            passive: true,
                            ..Default::default()
                        },
                        FieldTable::default(),
                    )
                    .await
                    .map_err(|e| format!("Failed to check queue: {}", e))?;

                info!("Found {} messages in hive_queue", queue.message_count());

                if queue.message_count() > 0 {
                    let mut consumer = channel
                        .basic_consume(
                            "hive_queue",
                            "initial_load_consumer",
                            BasicConsumeOptions {
                                no_ack: false,
                                exclusive: true,
                                ..Default::default()
                            },
                            FieldTable::default(),
                        )
                        .await
                        .map_err(|e| format!("Failed to create consumer: {}", e))?;

                    let mut message_count = 0;
                    while let Some(delivery_result) = consumer.next().await {
                        match delivery_result {
                            Ok(delivery) => {
                                info!("Processing message {} from queue...", message_count + 1);
                                let payload = String::from_utf8_lossy(&delivery.data);
                                debug!("Message payload: {}", payload);

                                match serde_json::from_slice::<Value>(&delivery.data) {
                                    Ok(hive_request) => {
                                        let hive_id = hive_request["hive_id"].as_i64().unwrap_or(0) as i32;
                                        info!("Found hive_id: {} in queue", hive_id);

                                        match serde_json::from_value::<Vec<Sensor>>(hive_request["sensors"].clone()) {
                                            Ok(sensors) => {
                                                debug!("Successfully parsed {} sensors for hive {}", sensors.len(), hive_id);
                                                for sensor in &sensors {
                                                    info!("  - Sensor {}: {}", sensor.sensor_id, sensor.sensor_type);
                                                }
                                                stored_hives.push((hive_id, sensors));

                                                if let Err(e) = channel
                                                    .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                                    .await
                                                {
                                                    error!("Failed to acknowledge message: {}", e);
                                                } else {
                                                    info!("Successfully acknowledged message for hive {}", hive_id);
                                                }
                                            },
                                            Err(e) => {
                                                error!("Failed to parse sensors for hive {}: {}", hive_id, e);
                                                if let Err(e) = channel
                                                    .basic_reject(delivery.delivery_tag, BasicRejectOptions { requeue: true })
                                                    .await
                                                {
                                                    error!("Failed to reject message: {}", e);
                                                }
                                            }
                                        }
                                    },
                                    Err(e) => {
                                        error!("Failed to parse message as JSON: {}", e);
                                        if let Err(e) = channel
                                            .basic_reject(delivery.delivery_tag, BasicRejectOptions { requeue: true })
                                            .await
                                        {
                                            error!("Failed to reject message: {}", e);
                                        }
                                    }
                                }
                                message_count += 1;
                                if message_count >= queue.message_count() {
                                    break;
                                }
                            },
                            Err(e) => {
                                error!("Error receiving message: {}", e);
                                break;
                            }
                        }
                    }

                    // Cancel the consumer
                    if let Err(e) = channel
                        .basic_cancel("initial_load_consumer", BasicCancelOptions::default())
                        .await
                    {
                        error!("Failed to cancel consumer: {}", e);
                    }
                }

                info!("Finished loading {} hives from queue", stored_hives.len());
                Ok(stored_hives)
            }
            .into_actor(self)
            .map(|result, _, ctx| {
                match result {
                    Ok(hives) => {
                        info!("Creating {} hives from queue", hives.len());
                        for (hive_id, sensors) in hives {
                            info!("Creating hive {} from queue data", hive_id);
                            let create_hive = CreateHive { hive_id, sensors };
                            ctx.notify(create_hive);
                        }
                        Ok(())
                    },
                    Err(e) => Err(e),
                }
            }),
        )
    }
}

// Update CreateHive handler to be more verbose
impl Handler<CreateHive> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: CreateHive, _ctx: &mut Context<Self>) -> Self::Result {
        info!("Received CreateHive request for hive {}", msg.hive_id);

        // Check if hive already exists
        if self.hives.contains_key(&msg.hive_id) {
            warn!("Hive {} already exists, skipping creation", msg.hive_id);
            return Box::pin(async move { Ok(()) }.into_actor(self));
        }

        let rabbitmq_url = self.rabbitmq_url.clone();
        let hive_id = msg.hive_id;
        let sensors = msg.sensors;

        info!(
            "Creating new hive {} with {} sensors:",
            hive_id,
            sensors.len()
        );
        for sensor in &sensors {
            info!(
                "  - Sensor {}: {} (type: {})",
                sensor.sensor_id, sensor.sensor_type, sensor.hive_id
            );
        }

        let sensor_ids = sensors.iter().map(|s| s.sensor_id).collect::<Vec<_>>();

        Box::pin(
            async move {
                let hive = Hive {
                    hive_id,
                    sensors: sensor_ids,
                    rabbitmq_url,
                    channel: None,
                };
                Ok((hive_id, hive.start()))
            }
            .into_actor(self)
            .map(
                move |result: Result<(i32, Addr<Hive>), ()>, actor, _ctx| match result {
                    Ok((hive_id, addr)) => {
                        actor.hives.insert(hive_id, addr);
                        debug!("Hive {} successfully created and started", hive_id);
                        Ok(())
                    }
                    Err(_) => {
                        error!("Failed to create hive {}", hive_id);
                        Err(())
                    }
                },
            ),
        )
    }
}

impl Apiary {
    fn start_consumers(&self, ctx: &mut Context<Self>) {
        let sensor_queue = "sensor_queue".to_string();
        let hive_queue = "hive_queue".to_string();

        let rabbitmq_channel = self.rabbitmq_channel.clone();
        let apiary_addr = ctx.address();

        // Consumer for sensor_queue
        ctx.spawn(
            async move {
                debug!("Attempting to lock the RabbitMQ channel...");
                let channel = rabbitmq_channel.lock().await;
                debug!("Successfully locked the RabbitMQ channel.");

                debug!("Starting to consume from sensor_queue...");
                let mut consumer = match channel
                    .basic_consume(
                        &sensor_queue,
                        "sensor_consumer",
                        BasicConsumeOptions::default(),
                        FieldTable::default(),
                    )
                    .await
                {
                    Ok(consumer) => {
                        debug!("Successfully created consumer for sensor_queue.");
                        consumer
                    }
                    Err(e) => {
                        info!("Failed to consume sensor_queue: {}", e);
                        return;
                    }
                };

                while let Some(delivery_result) = consumer.next().await {
                    match delivery_result {
                        Ok(delivery) => {
                            info!("Received a delivery from sensor_queue.");
                            let payload = &delivery.data;
                            let sensor: Sensor = match serde_json::from_slice(payload) {
                                Ok(sensor) => {
                                    debug!("Successfully deserialized sensor data.");
                                    sensor
                                }
                                Err(e) => {
                                    error!("Failed to deserialize sensor: {}", e);
                                    // Nack the message without requeueing
                                    if let Err(e) = channel
                                        .basic_nack(
                                            delivery.delivery_tag,
                                            BasicNackOptions {
                                                multiple: false,
                                                requeue: false,
                                            },
                                        )
                                        .await
                                    {
                                        error!("Failed to Nack message: {}", e);
                                    }
                                    continue;
                                }
                            };

                            debug!("Sending NewSensorData message to Apiary actor...");
                            if let Err(e) = apiary_addr.send(NewSensorData { sensor }).await {
                                error!("Failed to send NewSensorData message: {}", e);
                            } else {
                                debug!("NewSensorData message sent successfully.");
                            }

                            debug!("Acknowledging the message...");
                            if let Err(e) = channel
                                .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                .await
                            {
                                error!("Failed to ack message: {}", e);
                            } else {
                                debug!("Message acknowledged successfully.");
                            }
                        }
                        Err(e) => {
                            error!("Error in sensor_queue consumer: {}", e);
                            break;
                        }
                    }
                }
            }
            .into_actor(self),
        );

        // Consumer for hive_queue
        let rabbitmq_channel_hive = self.rabbitmq_channel.clone();
        let apiary_addr_hive = ctx.address();

        ctx.spawn(
            async move {
                info!("Starting hive_queue consumer...");
                let channel = rabbitmq_channel_hive.lock().await;
                debug!("Acquired rabbitmq channel lock");
                let mut consumer = match channel
                    .basic_consume(
                        &hive_queue,
                        "hive_consumer",
                        BasicConsumeOptions::default(),
                        FieldTable::default(),
                    )
                    .await
                {
                    Ok(consumer) => {
                        info!("Successfully created hive_queue consumer");
                        consumer
                    }
                    Err(e) => {
                        error!("Failed to consume hive_queue: {}", e);
                        return;
                    }
                };

                debug!("Starting to consume messages from hive_queue");
                while let Some(delivery_result) = consumer.next().await {
                    match delivery_result {
                        Ok(delivery) => {
                            info!("Received message from hive_queue");
                            let payload = &delivery.data;
                            let hive_request: Value = match serde_json::from_slice(payload) {
                                Ok(val) => {
                                    debug!("Successfully deserialized hive request");
                                    val
                                }
                                Err(e) => {
                                    error!("Failed to deserialize hive request: {}", e);
                                    // Nack the message without requeueing
                                    if let Err(e) = channel
                                        .basic_nack(
                                            delivery.delivery_tag,
                                            BasicNackOptions {
                                                multiple: false,
                                                requeue: false,
                                            },
                                        )
                                        .await
                                    {
                                        error!("Failed to Nack message: {}", e);
                                    }
                                    continue;
                                }
                            };

                            debug!("Extracting hive_id from request");
                            let hive_id = hive_request["hive_id"].as_i64().unwrap_or(0) as i32;
                            debug!("Extracted hive_id: {}", hive_id);

                            debug!("Deserializing sensors array");
                            let sensors = match serde_json::from_value::<Vec<Sensor>>(
                                hive_request["sensors"].clone(),
                            ) {
                                Ok(sensors) => {
                                    debug!("Successfully deserialized {} sensors", sensors.len());
                                    sensors
                                }
                                Err(e) => {
                                    error!("Failed to deserialize sensors: {}", e);
                                    // Nack the message without requeueing
                                    if let Err(e) = channel
                                        .basic_nack(
                                            delivery.delivery_tag,
                                            BasicNackOptions {
                                                multiple: false,
                                                requeue: false,
                                            },
                                        )
                                        .await
                                    {
                                        error!("Failed to Nack message: {}", e);
                                    }
                                    continue;
                                }
                            };

                            info!("Sending CreateHive message to Apiary actor");
                            // Send message to Apiary actor
                            if let Err(e) =
                                apiary_addr_hive.send(CreateHive { hive_id, sensors }).await
                            {
                                error!("Failed to send CreateHive message: {}", e);
                            } else {
                                debug!("Successfully sent CreateHive message");
                            }

                            debug!("Acknowledging message receipt");
                            // Acknowledge the message
                            if let Err(e) = channel
                                .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                .await
                            {
                                error!("Failed to ack message: {}", e);
                            } else {
                                debug!("Successfully acknowledged message");
                            }
                        }
                        Err(e) => {
                            error!("Error in hive_queue consumer: {}", e);
                            error!("Breaking out of consumer loop");
                            break;
                        }
                    }
                }
                info!("Exiting hive_queue consumer");
            }
            .into_actor(self),
        );
    }
}

// Handler for NewSensorData
impl Handler<NewSensorData> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: NewSensorData, _ctx: &mut Context<Self>) -> Self::Result {
        // Implement logic to handle new sensors
        info!("Handling NewSensorData message:");
        info!("  - Sensor ID: {}", msg.sensor.sensor_id);
        info!("  - Hive ID: {}", msg.sensor.hive_id);
        info!("  - Sensor Type: {}", msg.sensor.sensor_type);
        // For simplicity, acknowledge the sensor data
        Box::pin(async move { Ok(()) }.into_actor(self))
    }
}
