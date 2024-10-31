use crate::types::{Sensor, SensorReading};
use actix::Message;

// Message to create a new Hive
#[derive(Message)]
#[rtype(result = "Result<(), ()>")]
pub struct CreateHive {
    pub hive_id: i32,
    pub sensors: Vec<Sensor>,
}

// Message to handle new sensor data
#[derive(Message)]
#[rtype(result = "Result<(), ()>")]
pub struct NewSensorData {
    pub sensor: Sensor,
}

// Message to handle sensor readings
#[derive(Message)]
#[rtype(result = "Result<(), ()>")]
pub struct SensorReadingData {
    pub reading: SensorReading,
}
