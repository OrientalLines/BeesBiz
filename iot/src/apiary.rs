use crate::config::Config;
use crate::hive::Hive;
use crate::messages::{CreateHive, NewSensorData, SensorReadingData};
use crate::types::{Sensor, SensorReading};

use actix::prelude::*;
use futures_util::stream::StreamExt; // Import StreamExt for the `next` method
use lapin::{options::*, types::FieldTable, Channel, Connection, ConnectionProperties};
use serde_json::Value;
use std::sync::Arc;
use tokio::sync::Mutex; // Import the correct executor

pub struct Apiary {
    pub rabbitmq_connection: Arc<Mutex<Connection>>,
    pub rabbitmq_channel: Arc<Mutex<Channel>>,
}

impl Apiary {
    pub async fn new(rabbitmq_url: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let conn = Connection::connect(rabbitmq_url, ConnectionProperties::default()).await?;

        let channel = conn.create_channel().await?;

        // Declare necessary queues
        for queue in &["sensor_queue", "sensor_reading_queue", "hive_queue"] {
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

        Ok(Apiary {
            rabbitmq_connection: Arc::new(Mutex::new(conn)),
            rabbitmq_channel: Arc::new(Mutex::new(channel)),
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
                let channel = rabbitmq_channel.lock().await;
                let mut consumer = match channel
                    .basic_consume(
                        &sensor_queue,
                        "sensor_consumer",
                        BasicConsumeOptions::default(),
                        FieldTable::default(),
                    )
                    .await
                {
                    Ok(consumer) => consumer,
                    Err(e) => {
                        println!("Failed to consume sensor_queue: {}", e);
                        return;
                    }
                };

                while let Some(delivery_result) = consumer.next().await {
                    match delivery_result {
                        Ok(delivery) => {
                            let payload = &delivery.data;
                            let sensor: Sensor = match serde_json::from_slice(payload) {
                                Ok(sensor) => sensor,
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
                            if let Err(e) = apiary_addr.send(NewSensorData { sensor }).await {
                                println!("Failed to send NewSensorData message: {}", e);
                            }

                            // Acknowledge the message
                            if let Err(e) = channel
                                .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                .await
                            {
                                println!("Failed to ack message: {}", e);
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
                let channel = rabbitmq_channel_hive.lock().await;
                let mut consumer = match channel
                    .basic_consume(
                        &hive_queue,
                        "hive_consumer",
                        BasicConsumeOptions::default(),
                        FieldTable::default(),
                    )
                    .await
                {
                    Ok(consumer) => consumer,
                    Err(e) => {
                        println!("Failed to consume hive_queue: {}", e);
                        return;
                    }
                };

                while let Some(delivery_result) = consumer.next().await {
                    match delivery_result {
                        Ok(delivery) => {
                            let payload = &delivery.data;
                            let hive_request: Value = match serde_json::from_slice(payload) {
                                Ok(val) => val,
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

                            let hive_id = hive_request["hive_id"].as_i64().unwrap_or(0) as i32;
                            let sensors = match serde_json::from_value::<Vec<Sensor>>(
                                hive_request["sensors"].clone(),
                            ) {
                                Ok(sensors) => sensors,
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

                            // Send message to Apiary actor
                            if let Err(e) =
                                apiary_addr_hive.send(CreateHive { hive_id, sensors }).await
                            {
                                println!("Failed to send CreateHive message: {}", e);
                            }

                            // Acknowledge the message
                            if let Err(e) = channel
                                .basic_ack(delivery.delivery_tag, BasicAckOptions::default())
                                .await
                            {
                                println!("Failed to ack message: {}", e);
                            }
                        }
                        Err(e) => {
                            println!("Error in hive_queue consumer: {}", e);
                            break;
                        }
                    }
                }
            }
            .into_actor(self),
        );
    }
}

// Handler for CreateHive
impl Handler<CreateHive> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: CreateHive, _ctx: &mut Context<Self>) -> Self::Result {
        let rabbitmq_channel = self.rabbitmq_channel.clone();

        Box::pin(
            async move {
                let channel = rabbitmq_channel.lock().await.clone();
                // Start Hive actor
                let _hive_actor = Hive {
                    hive_id: msg.hive_id,
                    sensors: msg.sensors.iter().map(|s| s.sensor_id).collect(),
                    rabbitmq_channel: Arc::new(Mutex::new(channel)),
                }
                .start();

                println!("Created Hive actor with ID {}", msg.hive_id);
                Ok(())
            }
            .into_actor(self),
        )
    }
}

// Handler for NewSensorData
impl Handler<NewSensorData> for Apiary {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: NewSensorData, _ctx: &mut Context<Self>) -> Self::Result {
        // Implement logic to handle new sensors
        println!("Received new sensor data: {:?}", msg.sensor);
        // For simplicity, acknowledge the sensor data
        Box::pin(async move { Ok(()) }.into_actor(self))
    }
}
