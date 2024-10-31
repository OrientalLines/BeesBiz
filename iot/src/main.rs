// iot/src/main.rs

use actix::prelude::*;
use lapin::{options::QueueDeclareOptions, types::FieldTable, Connection};
use std::env;

mod apiary;
mod config;
mod hive;
mod messages;
mod types;

use apiary::Apiary;
use config::Config;

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    let args: Vec<String> = env::args().collect();
    let config_path = args.get(1).expect("Config file path required");
    let config = Config::from_file(config_path).expect("Failed to load configuration");

    // Create RabbitMQ connection
    let conn = Connection::connect(&config.apiary.rabbitmq_url, Default::default())
        .await
        .expect("Failed to connect to RabbitMQ");

    // Create and start the Apiary actor
    let apiary_addr = Apiary::new(&config.apiary.rabbitmq_url)
        .await
        .expect("Failed to create Apiary")
        .start();

    // Create initial hives based on configuration
    for hive_cfg in config.apiary.hives {
        let _ = apiary_addr
            .send(messages::CreateHive {
                hive_id: hive_cfg.id, // Use numeric ID from config
                sensors: hive_cfg
                    .sensors
                    .into_iter()
                    .map(|s| types::Sensor {
                        sensor_id: s.sensor_id,
                        hive_id: hive_cfg.id, // Use numeric ID
                        sensor_type: s.sensor_type,
                        last_reading: String::new(),
                        last_reading_time: None,
                    })
                    .collect(),
            })
            .await
            .unwrap_or_else(|e| {
                println!("Failed to send CreateHive message: {}", e);
                Ok(())
            });
    }

    println!("System running. Press Ctrl-C to exit.");

    // Keep the system running
    tokio::signal::ctrl_c().await?;
    println!("Shutting down...");

    // After starting the Apiary
    println!("Checking RabbitMQ queues...");
    if let Ok(channel) = conn.create_channel().await {
        if let Ok(queue) = channel
            .queue_declare(
                "sensor_reading_queue",
                QueueDeclareOptions::default(),
                FieldTable::default(),
            )
            .await
        {
            println!(
                "sensor_reading_queue status: {} messages, {} consumers",
                queue.message_count(),
                queue.consumer_count()
            );
        }
    }

    Ok(())
}
