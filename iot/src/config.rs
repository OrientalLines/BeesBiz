use serde::Deserialize;
use std::fs;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub apiary: ApiaryConfig,
}

#[derive(Debug, Deserialize)]
pub struct ApiaryConfig {
    pub name: String,
    pub rabbitmq_url: String,
    pub hives: Vec<HiveConfig>,
}

#[derive(Debug, Deserialize)]
pub struct HiveConfig {
    pub name: String,
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
