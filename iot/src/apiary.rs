use crate::config::Config;
use crate::hive::Hive;
use crate::messages::{CreateHive, NewSensorData, SensorReadingData};
use crate::types::{Sensor, SensorReading};

use actix::prelude::*;
use futures_util::stream::StreamExt; // Import StreamExt for the `next` method
use lapin::{options::*, types::FieldTable, Channel, Connection, ConnectionProperties};
use serde_json::Value;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::Mutex; // Import the correct executor

pub struct Apiary {
    pub rabbitmq_url: String,
    pub rabbitmq_connection: Arc<Mutex<Connection>>,
    pub rabbitmq_channel: Arc<Mutex<Channel>>,
    pub hives: HashMap<i32, Addr<Hive>>, // Add this line
}

impl Apiary {
    pub async fn new(rabbitmq_url: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let conn = Connection::connect(rabbitmq_url, ConnectionProperties::default()).await?;

        let channel = conn.create_channel().await?;

        // Declare necessary queues with specific options
        for queue in &["sensor_queue", "hive_queue"] {
            channel
                .queue_declare(
                    queue,
                    QueueDeclareOptions {
                        durable: true,
                        ..Default::default()
                    },
                    FieldTable::default(),
                )
                .await?;
        }

        // Declare sensor_reading_queue separately with specific options
        channel
            .queue_declare(
                "sensor_reading_queue",
                QueueDeclareOptions {
                    durable: true,
                    auto_delete: false,
                    exclusive: false,
                    ..Default::default()
                },
                FieldTable::default(),
            )
            .await?;

        Ok(Apiary {
            rabbitmq_connection: Arc::new(Mutex::new(conn)),
            rabbitmq_channel: Arc::new(Mutex::new(channel)),
            hives: HashMap::new(),
            rabbitmq_url: rabbitmq_url.to_string(),
        })
    }
}

impl Actor for Apiary {
    type Context = Context<Self>;

    fn started(&mut self, ctx: &mut Context<Self>) {
        println!("Apiary started");
        self.start_consumers(ctx);
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
                println!("Attempting to lock the RabbitMQ channel...");
                let channel = rabbitmq_channel.lock().await;
                println!("Successfully locked the RabbitMQ channel.");

                println!("Starting to consume from sensor_queue...");
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
                        println!("Successfully created consumer for sensor_queue.");
                        consumer
                    }
                    Err(e) => {
                        println!("Failed to consume sensor_queue: {}", e);
                        return;
                    }
                };

                while let Some(delivery_result) = consumer.next().await {
                    match delivery_result {
                        Ok(delivery) => {
                            println!("Received a delivery from sensor_queue.");
                            let payload = &delivery.data;
                            let sensor: Sensor = match serde_json::from_slice(payload) {
                                Ok(sensor) => {
                                    println!("Successfully deserialized sensor data.");
                                    sensor
                                }
                                Err(e) => {
                                    println!("Failed to deserialize sensor: {}", e);
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
                                        println!("Failed to Nack message: {}", e);
                                    }
                                    continue;
                                }
                            };

                            // Send message to Apiary actor
                            println!("Sending NewSensorData message to Apiary actor...");
                            if let Err(e) = apiary_addr.send(NewSensorData { sensor }).await {
                                println!("Failed to send NewSensorData message: {}", e);
                            } else {
                                println!("NewSensorData message sent successfully.");
                            }

                            // Acknowledge the message
                            println!("Acknowledging the message...");
                            if let Err(e) = channel
                                .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                .await
                            {
                                println!("Failed to ack message: {}", e);
                            } else {
                                println!("Message acknowledged successfully.");
                            }
                        }
                        Err(e) => {
                            println!("Error in sensor_queue consumer: {}", e);
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
                println!("Starting hive_queue consumer...");
                let channel = rabbitmq_channel_hive.lock().await;
                println!("Acquired rabbitmq channel lock");
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
                        println!("Successfully created hive_queue consumer");
                        consumer
                    }
                    Err(e) => {
                        println!("Failed to consume hive_queue: {}", e);
                        return;
                    }
                };

                println!("Starting to consume messages from hive_queue");
                while let Some(delivery_result) = consumer.next().await {
                    match delivery_result {
                        Ok(delivery) => {
                            println!("Received message from hive_queue");
                            let payload = &delivery.data;
                            let hive_request: Value = match serde_json::from_slice(payload) {
                                Ok(val) => {
                                    println!("Successfully deserialized hive request");
                                    val
                                }
                                Err(e) => {
                                    println!("Failed to deserialize hive request: {}", e);
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
                                        println!("Failed to Nack message: {}", e);
                                    }
                                    continue;
                                }
                            };

                            println!("Extracting hive_id from request");
                            let hive_id = hive_request["hive_id"].as_i64().unwrap_or(0) as i32;
                            println!("Extracted hive_id: {}", hive_id);

                            println!("Deserializing sensors array");
                            let sensors = match serde_json::from_value::<Vec<Sensor>>(
                                hive_request["sensors"].clone(),
                            ) {
                                Ok(sensors) => {
                                    println!("Successfully deserialized {} sensors", sensors.len());
                                    sensors
                                }
                                Err(e) => {
                                    println!("Failed to deserialize sensors: {}", e);
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
                                        println!("Failed to Nack message: {}", e);
                                    }
                                    continue;
                                }
                            };

                            println!("Sending CreateHive message to Apiary actor");
                            // Send message to Apiary actor
                            if let Err(e) =
                                apiary_addr_hive.send(CreateHive { hive_id, sensors }).await
                            {
                                println!("Failed to send CreateHive message: {}", e);
                            } else {
                                println!("Successfully sent CreateHive message");
                            }

                            println!("Acknowledging message receipt");
                            // Acknowledge the message
                            if let Err(e) = channel
                                .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                .await
                            {
                                println!("Failed to ack message: {}", e);
                            } else {
                                println!("Successfully acknowledged message");
                            }
                        }
                        Err(e) => {
                            println!("Error in hive_queue consumer: {}", e);
                            println!("Breaking out of consumer loop");
                            break;
                        }
                    }
                }
                println!("Exiting hive_queue consumer");
            }
            .into_actor(self),
        );
    }
}

impl Handler<CreateHive> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: CreateHive, _ctx: &mut Context<Self>) -> Self::Result {
        println!("Handling CreateHive message for hive ID: {}", msg.hive_id);
        let rabbitmq_url = self.rabbitmq_url.clone();
        let hive_id = msg.hive_id;
        let sensors = msg.sensors.iter().map(|s| s.sensor_id).collect::<Vec<_>>();

        Box::pin(
            async move {
                let hive_addr = Hive {
                    hive_id,
                    sensors,
                    rabbitmq_url,
                    channel: None,
                }
                .start();

                Ok((hive_id, hive_addr))
            }
            .into_actor(self)
            .map(
                |result: Result<(i32, Addr<Hive>), ()>, actor, _ctx| match result {
                    Ok((hive_id, hive_addr)) => {
                        actor.hives.insert(hive_id, hive_addr);
                        println!("Hive {} inserted into Apiary's hives.", hive_id);
                        Ok(())
                    }
                    Err(_) => Err(()),
                },
            ),
        )
    }
}
// Handler for NewSensorData
impl Handler<NewSensorData> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: NewSensorData, _ctx: &mut Context<Self>) -> Self::Result {
        // Implement logic to handle new sensors
        println!("Handling NewSensorData message:");
        println!("  - Sensor ID: {}", msg.sensor.sensor_id);
        println!("  - Hive ID: {}", msg.sensor.hive_id);
        println!("  - Sensor Type: {}", msg.sensor.sensor_type);
        // For simplicity, acknowledge the sensor data
        Box::pin(async move { Ok(()) }.into_actor(self))
    }
}
