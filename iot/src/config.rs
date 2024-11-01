use log::LevelFilter;
use serde::Deserialize;
use std::fs;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub apiary: ApiaryConfig,
}

#[derive(Debug, Deserialize)]
pub struct ApiaryConfig {
    pub rabbitmq_url: String,
    pub log_level: String,
    pub hives: Vec<HiveConfig>,
}

impl ApiaryConfig {
    pub fn get_log_level(&self) -> LevelFilter {
        match self.log_level.to_lowercase().as_str() {
            "trace" => LevelFilter::Trace,
            "debug" => LevelFilter::Debug,
            "info" => LevelFilter::Info,
            "warn" => LevelFilter::Warn,
            "error" => LevelFilter::Error,
            _ => LevelFilter::Info,
        }
    }
}

#[derive(Debug, Deserialize)]
pub struct HiveConfig {
    pub id: i32,
    pub sensors: Vec<SensorConfig>,
}

#[derive(Debug, Deserialize)]
pub struct SensorConfig {
    pub sensor_id: i32,
    pub sensor_type: String,
}

impl Config {
    pub fn from_file(path: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let content = fs::read_to_string(path)?;
        let config = toml::from_str(&content)?;
        Ok(config)
    }
}
