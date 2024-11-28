use actix::prelude::*;
use base64::{engine::general_purpose::STANDARD as BASE64, Engine};
use chrono::Utc;
use lapin::options::BasicPublishOptions;
use lapin::{Channel, Connection, ConnectionProperties};
use log::{debug, error, info};
use rand::Rng;
use std::time::Duration;

use crate::messages::{DeleteSensor, SensorReadingData};
use crate::types::SensorReading;

pub struct Hive {
    pub hive_id: i32,
    pub sensors: Vec<i32>,
    pub rabbitmq_url: String,
    pub channel: Option<Channel>,
}

impl Actor for Hive {
    type Context = Context<Self>;

    fn started(&mut self, ctx: &mut Self::Context) {
        info!(
            "Hive {} started with {} sensors",
            self.hive_id,
            self.sensors.len()
        );
        info!("Registered sensors: {:?}", self.sensors);

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
                info!("Hive {} initialized RabbitMQ channel", hive_id);

                // Set up sensor reading interval after channel is initialized
                ctx.run_interval(Duration::from_secs(10), move |act, _ctx| {
                    if let Some(channel) = &act.channel {
                        info!("Generating sensor readings for hive {}:", act.hive_id);
                        for sensor_id in &act.sensors {
                            let reading = SensorReading {
                                reading_id: None,
                                sensor_id: *sensor_id,
                                value: act.generate_mock_value(),
                                timestamp: Some(Utc::now().to_rfc3339()),
                            };

                            info!("Generated reading: {:?}", reading);

                            // Clone necessary data for async block
                            let channel = channel.clone();
                            let reading_clone = reading.clone();

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
                                                    Ok(_) => info!(
                                                        "Successfully published reading for sensor {} (value: {:?})",
                                                        reading_clone.sensor_id,
                                                        reading_clone.value
                                                    ),
                                                    Err(e) => error!("Failed to confirm publish: {}", e),
                                                }
                                            }
                                            Err(e) => error!("Failed to publish reading: {}", e),
                                        }
                                    }
                                    Err(e) => error!("Failed to serialize reading: {}", e),
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
    fn generate_mock_value(&self) -> String {
        let mut rng = rand::thread_rng();
        let values = vec![
            rng.gen_range(10..80),
            rng.gen_range(0..255),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
            rng.gen_range(0..100),
        ];

        BASE64.encode(&values)
    }
}

// Handler for SensorReadingData
impl Handler<SensorReadingData> for Hive {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: SensorReadingData, _ctx: &mut Context<Self>) -> Self::Result {
        info!(
            "Handling SensorReadingData for sensor {} in hive {}",
            msg.reading.sensor_id, self.hive_id
        );

        let channel = self.channel.clone();
        let reading = msg.reading;

        Box::pin(
            async move {
                debug!("Preparing to publish reading to RabbitMQ...");
                let payload = serde_json::to_vec(&reading).map_err(|e| {
                    error!("Failed to serialize reading: {}", e);
                })?;
                debug!("Acquiring channel lock...");
                let channel = channel.as_ref().ok_or(())?;

                // Add persistent message properties
                let properties = lapin::BasicProperties::default()
                    .with_delivery_mode(2) // Makes the message persistent
                    .with_content_type("application/json".into());

                debug!("Publishing message to sensor_reading_queue...");
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
                                info!(
                                    "Successfully published reading for sensor {} (value: {:?}) to RabbitMQ",
                                    reading.sensor_id,
                                    reading.value
                                );
                                Ok(())
                            }
                            Err(e) => {
                                error!("Failed to confirm publish: {}", e);
                                Err(())
                            }
                        }
                    }
                    Err(e) => {
                        error!("Failed to publish reading: {}", e);
                        Err(())
                    }
                }
            }
            .into_actor(self),
        )
    }
}

impl Handler<DeleteSensor> for Hive {
    type Result = Result<(), ()>;

    fn handle(&mut self, msg: DeleteSensor, _ctx: &mut Context<Self>) -> Self::Result {
        info!(
            "Removing sensor {} from hive {}",
            msg.sensor_id, self.hive_id
        );

        if self.hive_id != msg.hive_id {
            error!(
                "Hive ID mismatch: expected {}, got {}",
                self.hive_id, msg.hive_id
            );
            return Err(());
        }

        // Remove the sensor from the sensors list
        self.sensors.retain(|&sensor_id| sensor_id != msg.sensor_id);
        info!(
            "Hive {} now has {} sensors",
            self.hive_id,
            self.sensors.len()
        );
        Ok(())
    }
}
