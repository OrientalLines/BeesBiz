use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Debug, Serialize, Deserialize)]
pub struct Sensor {
    pub sensor_id: i32,
    pub hive_id: i32,
    pub sensor_type: String,
    pub last_reading: Vec<u8>,
    pub last_reading_time: Option<String>, // Using String for simplicity
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SensorReading {
    pub reading_id: Option<i32>,
    pub sensor_id: i32,
    pub value: Vec<u8>,
    pub timestamp: Option<String>, // Using String for simplicity
}
