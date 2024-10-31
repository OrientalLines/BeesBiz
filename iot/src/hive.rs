// iot/src/hive.rs

use actix::prelude::*;
use base64::{engine::general_purpose::STANDARD as BASE64, Engine};
use chrono::Utc;
use lapin::options::BasicPublishOptions;
use lapin::{BasicProperties, Channel, Connection, ConnectionProperties};
use rand::Rng;
use serde_json::json;
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::Mutex;

use crate::messages::SensorReadingData;
use crate::types::SensorReading;

pub struct Hive {
    pub hive_id: i32,
    pub sensors: Vec<i32>,        // List of sensor IDs
    pub rabbitmq_url: String,     // Store URL instead of channel
    pub channel: Option<Channel>, // Make it Option<Channel>
}

impl Hive {
    // Initialize the channel when the actor starts
    async fn initialize_channel(&mut self) -> Result<(), Box<dyn std::error::Error>> {
        let conn = Connection::connect(&self.rabbitmq_url, ConnectionProperties::default()).await?;
        self.channel = Some(conn.create_channel().await?);
        Ok(())
    }
}

impl Actor for Hive {
    type Context = Context<Self>;

    fn started(&mut self, ctx: &mut Self::Context) {
        println!(
            "Hive {} started with {} sensors",
            self.hive_id,
            self.sensors.len()
        );
        println!("Registered sensors: {:?}", self.sensors);

        // Initialize channel
        let hive_id = self.hive_id;
        let rabbitmq_url = self.rabbitmq_url.clone();

        ctx.wait(
            async move {
                let conn = Connection::connect(&rabbitmq_url, ConnectionProperties::default())
                    .await
                    .expect("Failed to connect to RabbitMQ");
                conn.create_channel()
                    .await
                    .expect("Failed to create channel")
            }
            .into_actor(self)
            .map(move |channel, act, ctx| {
                act.channel = Some(channel);
                println!("Hive {} initialized RabbitMQ channel", hive_id);

                // Set up sensor reading interval after channel is initialized
                ctx.run_interval(Duration::from_secs(5), move |act, _ctx| {
                    if let Some(channel) = &act.channel {
                        println!("\nGenerating sensor readings for hive {}:", act.hive_id);
                        for sensor_id in &act.sensors {
                            let reading = SensorReading {
                                reading_id: None,
                                sensor_id: *sensor_id,
                                value: act.generate_mock_value(*sensor_id),
                                timestamp: Some(Utc::now().to_rfc3339()),
                            };

                            println!("  - Generated reading: {:?}", reading);

                            // Clone necessary data for async block
                            let channel = channel.clone();
                            let reading_clone = reading.clone();

                            // Spawn a future to publish the message
                            actix::spawn(async move {
                                match serde_json::to_vec(&reading_clone) {
                                    Ok(payload) => {
                                        let props = lapin::BasicProperties::default()
                                            .with_delivery_mode(2) // persistent
                                            .with_content_type("application/json".into())
                                            .with_priority(0);

                                        match channel
                                            .basic_publish(
                                                "",
                                                "sensor_reading_queue",
                                                BasicPublishOptions {
                                                    mandatory: true,
                                                    ..Default::default()
                                                },
                                                &payload,
                                                props,
                                            )
                                            .await
                                        {
                                            Ok(confirm) => {
                                                match confirm.await {
                                                    Ok(_) => println!(
                                                        "✅ Successfully published reading for sensor {} (value: {:?})",
                                                        reading_clone.sensor_id,
                                                        reading_clone.value
                                                    ),
                                                    Err(e) => println!("❌ Failed to confirm publish: {}", e),
                                                }
                                            }
                                            Err(e) => println!("❌ Failed to publish reading: {}", e),
                                        }
                                    }
                                    Err(e) => println!("Failed to serialize reading: {}", e),
                                }
                            });
                        }
                    }
                });
            }),
        );
    }
}

impl Hive {
    /// Generates a mock sensor value based on sensor type
    fn generate_mock_value(&self, sensor_id: i32) -> String {
        // Generate random values as before
        let mut rng = rand::thread_rng();
        let values = vec![rng.gen_range(20..30), rng.gen_range(0..100)];

        // Convert to base64
        BASE64.encode(&values)
    }
}

// Handler for SensorReadingData
impl Handler<SensorReadingData> for Hive {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: SensorReadingData, _ctx: &mut Context<Self>) -> Self::Result {
        println!(
            "Handling SensorReadingData for sensor {} in hive {}",
            msg.reading.sensor_id, self.hive_id
        );

        let channel = self.channel.clone();
        let reading = msg.reading;

        Box::pin(
            async move {
                println!("Preparing to publish reading to RabbitMQ...");
                let payload = serde_json::to_vec(&reading).map_err(|e| {
                    println!("Failed to serialize reading: {}", e);
                    ()
                })?;
                println!("Acquiring channel lock...");
                let channel = channel.as_ref().ok_or(())?;

                // Add persistent message properties
                let properties = lapin::BasicProperties::default()
                    .with_delivery_mode(2) // Makes the message persistent
                    .with_content_type("application/json".into());

                println!("Publishing message to sensor_reading_queue...");
                match channel
                    .basic_publish(
                        "",
                        "sensor_reading_queue",
                        lapin::options::BasicPublishOptions::default(),
                        &payload,
                        properties,
                    )
                    .await
                {
                    Ok(confirm) => {
                        match confirm.await {
                            Ok(_) => {
                                println!(
                                    "✅ Successfully published reading for sensor {} (value: {:?}) to RabbitMQ",
                                    reading.sensor_id,
                                    reading.value
                                );
                                Ok(())
                            }
                            Err(e) => {
                                println!("❌ Failed to confirm publish: {}", e);
                                Err(())
                            }
                        }
                    }
                    Err(e) => {
                        println!("❌ Failed to publish reading: {}", e);
                        Err(())
                    }
                }
            }
            .into_actor(self),
        )
    }
}
