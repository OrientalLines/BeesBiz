use actix::prelude::*;
use lapin::{
    options::BasicPublishOptions, publisher_confirm::Confirmation, BasicProperties, Channel,
};
use serde_json::json;
use std::sync::Arc;
use tokio::sync::Mutex;

use crate::messages::SensorReadingData;
use crate::types::SensorReading;

pub struct Hive {
    pub hive_id: i32,
    pub sensors: Vec<i32>, // List of sensor IDs
    pub rabbitmq_channel: Arc<Mutex<Channel>>,
}

impl Actor for Hive {
    type Context = Context<Self>;

    fn started(&mut self, _ctx: &mut Self::Context) {
        println!("Hive {} started", self.hive_id);
    }
}

// Handler for SensorReadingData
impl Handler<SensorReadingData> for Hive {
    type Result = ResponseActFuture<Self, Result<(), ()>>;

    fn handle(&mut self, msg: SensorReadingData, _ctx: &mut Context<Self>) -> Self::Result {
        let channel = self.rabbitmq_channel.clone();
        let reading = msg.reading;

        Box::pin(
            async move {
                let payload = match serde_json::to_vec(&reading) {
                    Ok(data) => data,
                    Err(e) => {
                        println!("Failed to serialize reading: {}", e);
                        return Err(());
                    }
                };

                let confirm = channel
                    .lock()
                    .await
                    .basic_publish(
                        "",
                        "sensor_reading_queue",
                        BasicPublishOptions::default(),
                        &payload,
                        BasicProperties::default(),
                    )
                    .await
                    .map_err(|e| {
                        println!("Failed to publish reading: {}", e);
                        ()
                    })?
                    .await
                    .map_err(|e| {
                        println!("Failed to confirm publish: {}", e);
                        ()
                    })?;

                match confirm {
                    Confirmation::NotRequested => Ok(()),
                    Confirmation::Ack(_) => Ok(()),
                    Confirmation::Nack(_) => {
                        println!("Publish Nacked");
                        Err(())
                    }
                }
            }
            .into_actor(self),
        )
    }
}
