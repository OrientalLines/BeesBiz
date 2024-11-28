use crate::hive::Hive;
use crate::messages::{CreateHive, DeleteSensor, NewSensorData};
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

        // Declare queues with persistence and proper settings
        // TODO: Remove this hardcoded list
        for queue in &["sensor_queue", "hive_queue", "sensor_reading_queue"] {
            channel
                .queue_declare(
                    queue,
                    QueueDeclareOptions {
                        durable: true,
                        auto_delete: false,
                        exclusive: false,
                        passive: false,
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
                                exclusive: false,
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

                    // Ensure we properly cancel the consumer before starting regular consumers
                    debug!("Canceling initial load consumer...");
                    if let Err(e) = channel
                        .basic_cancel("initial_load_consumer", BasicCancelOptions::default())
                        .await
                    {
                        error!("Failed to cancel initial load consumer: {}", e);
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

impl Handler<CreateHive> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: CreateHive, _ctx: &mut Context<Self>) -> Self::Result {
        info!("Received CreateHive request for hive {}", msg.hive_id);

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
                "  - Sensor {}: {} (hive: {})",
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
        let sensor_delete_queue = "sensor_delete_queue".to_string();
        let rabbitmq_channel = self.rabbitmq_channel.clone();
        let apiary_addr = ctx.address();

        ctx.spawn(
            async move {
                debug!("Setting up consumers...");
                let (sensor_consumer, hive_consumer, delete_consumer, channel) = {
                    let channel = rabbitmq_channel.lock().await;

                    // Declare queues to ensure they exist
                    for queue in &[&sensor_queue, &hive_queue, &sensor_delete_queue] {
                        channel
                            .queue_declare(
                                queue,
                                QueueDeclareOptions {
                                    durable: true,
                                    auto_delete: false,
                                    exclusive: false,
                                    passive: false,
                                    ..Default::default()
                                },
                                FieldTable::default(),
                            )
                            .await
                            .expect("Failed to declare queue");
                    }

                    let sensor_consumer = channel
                        .basic_consume(
                            &sensor_queue,
                            "sensor_consumer",
                            BasicConsumeOptions {
                                no_ack: false,
                                exclusive: false,
                                ..Default::default()
                            },
                            FieldTable::default(),
                        )
                        .await
                        .expect("Failed to create sensor consumer");

                    let hive_consumer = channel
                        .basic_consume(
                            &hive_queue,
                            "hive_consumer",
                            BasicConsumeOptions {
                                no_ack: false,
                                exclusive: false,
                                ..Default::default()
                            },
                            FieldTable::default(),
                        )
                        .await
                        .expect("Failed to create hive consumer");

                    let delete_consumer = channel
                        .basic_consume(
                            &sensor_delete_queue,
                            "delete_consumer",
                            BasicConsumeOptions {
                                no_ack: false,
                                exclusive: false,
                                ..Default::default()
                            },
                            FieldTable::default(),
                        )
                        .await
                        .expect("Failed to create delete consumer");

                    (
                        sensor_consumer,
                        hive_consumer,
                        delete_consumer,
                        channel.clone(),
                    )
                }; // Lock is dropped here

                // Create separate channel clones for each consumer
                let sensor_channel = channel.clone();
                let hive_channel = channel.clone();
                let delete_channel = channel.clone();

                let apiary_addr_clone = apiary_addr.clone();
                let apiary_addr_clone2 = apiary_addr.clone();
                let apiary_addr_clone3 = apiary_addr.clone();

                let sensor_handle = tokio::spawn(async move {
                    let mut sensor_consumer = sensor_consumer;
                    while let Some(delivery_result) = sensor_consumer.next().await {
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
                                        if let Err(e) = sensor_channel
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

                                if let Err(e) =
                                    apiary_addr_clone.send(NewSensorData { sensor }).await
                                {
                                    error!("Failed to send NewSensorData message: {}", e);
                                }

                                if let Err(e) = sensor_channel
                                    .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                    .await
                                {
                                    error!("Failed to ack message: {}", e);
                                }
                            }
                            Err(e) => {
                                error!("Error in sensor_queue consumer: {}", e);
                                break;
                            }
                        }
                    }
                });

                let hive_handle = tokio::spawn(async move {
                    let mut hive_consumer = hive_consumer;
                    while let Some(delivery_result) = hive_consumer.next().await {
                        match delivery_result {
                            Ok(delivery) => {
                                info!("Received message from hive_queue");
                                let payload = &delivery.data;
                                let hive_request: Value = match serde_json::from_slice(payload) {
                                    Ok(val) => val,
                                    Err(e) => {
                                        error!("Failed to deserialize hive request: {}", e);
                                        if let Err(e) = hive_channel
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

                                let hive_id = hive_request["hive_id"].as_i64().unwrap_or(0) as i32;
                                let sensors = match serde_json::from_value::<Vec<Sensor>>(
                                    hive_request["sensors"].clone(),
                                ) {
                                    Ok(sensors) => sensors,
                                    Err(e) => {
                                        error!("Failed to deserialize sensors: {}", e);
                                        if let Err(e) = hive_channel
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

                                if let Err(e) = apiary_addr_clone2
                                    .send(CreateHive { hive_id, sensors })
                                    .await
                                {
                                    error!("Failed to send CreateHive message: {}", e);
                                }

                                if let Err(e) = hive_channel
                                    .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                    .await
                                {
                                    error!("Failed to ack message: {}", e);
                                }
                            }
                            Err(e) => {
                                error!("Error in hive_queue consumer: {}", e);
                                break;
                            }
                        }
                    }
                });

                let delete_handle = tokio::spawn(async move {
                    let mut delete_consumer = delete_consumer;
                    while let Some(delivery_result) = delete_consumer.next().await {
                        match delivery_result {
                            Ok(delivery) => {
                                info!("Received message from sensor_delete_queue");
                                let payload = &delivery.data;
                                match serde_json::from_slice::<DeleteSensor>(payload) {
                                    Ok(delete_request) => {
                                        info!(
                                            "Processing delete request for sensor {} in hive {}",
                                            delete_request.sensor_id, delete_request.hive_id
                                        );

                                        if let Err(e) =
                                            apiary_addr_clone3.send(delete_request).await
                                        {
                                            error!("Failed to send DeleteSensor message: {}", e);
                                            // Nack the message if we failed to process it
                                            if let Err(e) = delete_channel
                                                .basic_nack(
                                                    delivery.delivery_tag,
                                                    BasicNackOptions {
                                                        multiple: false,
                                                        requeue: true,
                                                    },
                                                )
                                                .await
                                            {
                                                error!("Failed to Nack message: {}", e);
                                            }
                                            continue;
                                        }

                                        // Acknowledge successful processing
                                        if let Err(e) = delete_channel
                                            .basic_ack(
                                                delivery.delivery_tag,
                                                BasicAckOptions::default(),
                                            )
                                            .await
                                        {
                                            error!("Failed to ack message: {}", e);
                                        }
                                    }
                                    Err(e) => {
                                        error!("Failed to deserialize delete request: {}", e);
                                        if let Err(e) = delete_channel
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
                                    }
                                }
                            }
                            Err(e) => {
                                error!("Error in delete_queue consumer: {}", e);
                                break;
                            }
                        }
                    }
                });

                tokio::try_join!(sensor_handle, hive_handle, delete_handle)
                    .expect("Consumer task failed");
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

impl Handler<DeleteSensor> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: DeleteSensor, _ctx: &mut Context<Self>) -> Self::Result {
        info!(
            "Processing delete request for sensor {} in hive {}",
            msg.sensor_id, msg.hive_id
        );

        // Find the hive
        if let Some(hive_addr) = self.hives.get(&msg.hive_id) {
            let hive_addr = hive_addr.clone();
            Box::pin(
                async move {
                    // Forward the delete request to the hive
                    hive_addr.send(msg).await.map_err(|e| {
                        error!("Failed to send DeleteSensor message to hive: {}", e);
                        ()
                    })?
                }
                .into_actor(self)
                .map(|result, _, _| result),
            )
        } else {
            error!("Hive {} not found", msg.hive_id);
            Box::pin(async move { Err(()) }.into_actor(self))
        }
    }
}
