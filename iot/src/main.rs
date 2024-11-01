use actix::prelude::*;
use apiary::{Apiary, LoadExistingHives, StartConsumers};
use config::Config;
use lapin::{options::QueueDeclareOptions, types::FieldTable, Connection};
use log::{error, info};
use std::env;

mod apiary;
mod config;
mod hive;
mod logging;
mod messages;
mod types;

#[actix_rt::main]
async fn main() -> std::io::Result<()> {
    let args: Vec<String> = env::args().collect();
    let config_path = args.get(1).expect("Config file path required");
    let config = Config::from_file(config_path).expect("Failed to load configuration");

    logging::init_logger(config.apiary.get_log_level());
    logging::print_banner();

    info!("Connecting to RabbitMQ...");
    let conn = Connection::connect(&config.apiary.rabbitmq_url, Default::default())
        .await
        .map_err(|e| {
            error!("Failed to connect to RabbitMQ: {}", e);
            e
        })
        .expect("Failed to connect to RabbitMQ");

    let apiary_addr = Apiary::new(&config.apiary.rabbitmq_url)
        .await
        .expect("Failed to create Apiary")
        .start();

    info!("Starting regular consumers...");
    apiary_addr
        .send(StartConsumers)
        .await
        .expect("Failed to start consumers");

    info!("Loading existing hives from queue...");
    if let Err(e) = apiary_addr.send(LoadExistingHives).await {
        error!("Failed to load existing hives: {}", e);
    }

    info!("Loading hives from configuration...");
    for hive_cfg in config.apiary.hives {
        let create_hive = messages::CreateHive {
            hive_id: hive_cfg.id,
            sensors: hive_cfg
                .sensors
                .into_iter()
                .map(|s| types::Sensor {
                    sensor_id: s.sensor_id,
                    hive_id: hive_cfg.id,
                    sensor_type: s.sensor_type,
                    last_reading: String::new(),
                    last_reading_time: None,
                })
                .collect(),
        };

        if let Err(e) = apiary_addr.send(create_hive).await {
            error!("Failed to create hive from config: {}", e);
        }
    }

    info!("System running. Press Ctrl-C to exit.");

    tokio::signal::ctrl_c().await?;
    info!("Shutting down...");

    info!("Checking RabbitMQ queues...");
    if let Ok(channel) = conn.create_channel().await {
        if let Ok(queue) = channel
            .queue_declare(
                "sensor_reading_queue",
                QueueDeclareOptions::default(),
                FieldTable::default(),
            )
            .await
        {
            info!(
                "sensor_reading_queue status: {} messages, {} consumers",
                queue.message_count(),
                queue.consumer_count()
            );
        }
    }

    Ok(())
}
