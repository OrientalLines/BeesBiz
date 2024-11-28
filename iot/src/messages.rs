use crate::types::{Sensor, SensorReading};
use actix::Message;
use serde::{Deserialize, Serialize};

#[derive(Message)]
#[rtype(result = "Result<(), ()>")]
pub struct CreateHive {
    pub hive_id: i32,
    pub sensors: Vec<Sensor>,
}

#[derive(Message)]
#[rtype(result = "Result<(), ()>")]
pub struct NewSensorData {
    pub sensor: Sensor,
}

#[derive(Message)]
#[rtype(result = "Result<(), ()>")]
pub struct SensorReadingData {
    pub reading: SensorReading,
}

#[derive(Message, Serialize, Deserialize)] // Added Serialize, Deserialize
#[rtype(result = "Result<(), ()>")]
pub struct DeleteSensor {
    pub hive_id: i32,
    pub sensor_id: i32,
}
