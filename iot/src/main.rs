use actix::prelude::*;
use std::env;

mod apiary;
mod config;
mod hive;
mod messages;
mod types;

use apiary::Apiary;
use config::Config;

#[actix::main]
async fn main() -> std::io::Result<()> {
    // Load configuration
    let config_path = env::args()
        .nth(1)
        .unwrap_or_else(|| "config.toml".to_string());
    let config = Config::from_file(&config_path).expect("Failed to load configuration");

    // Start Apiary Actor
    let apiary = Apiary::new(&config.apiary.rabbitmq_url)
        .await
        .expect("Failed to create Apiary");

    let apiary_addr = apiary.start();

    // Optionally, create initial hives based on configuration
    for hive_cfg in config.apiary.hives {
        let _ = apiary_addr
            .send(messages::CreateHive {
                hive_id: hive_cfg.name.parse().unwrap_or(0), // Assuming hive_id is derived from name
                sensors: hive_cfg
                    .sensors
                    .into_iter()
                    .map(|s| types::Sensor {
                        sensor_id: s.sensor_id,
                        hive_id: hive_cfg.name.parse().unwrap_or(0), // Assuming hive_id is derived from name
                        sensor_type: s.sensor_type,
                        last_reading: Vec::new(),
                        last_reading_time: None,
                    })
                    .collect(),
            })
            .await
            .unwrap_or_else(|e| {
                println!("Failed to send CreateHive message: {}", e);
                Ok(()) // Return Result<(), ()> as expected
            });
    }
    // Start the Actix system and keep it running
    actix_rt::System::new().block_on(async { Ok(()) })
}
