export type Time = string | Date;

export type Role = 'ADMIN' | 'WORKER' | 'MANAGER';

export interface User {
	user_id: number;
	username: string;
	full_name: string;
	role: Role;
	email: string;
	last_login: Time | null;
}

export interface Region {
	region_id: number;
	name: string;
	climate_zone: string;
}

export interface RegionApiary {
	id: number;
	apiary_id: number;
	region_id: number;
	region: Region;
	apiary: Apiary;
}

export interface AllowedRegion {
	id: number;
	user_id: number;
	region_id: number;
	region: Region;
	user: User;
}

export interface Apiary {
	apiary_id: number;
	location: string;
	manager_id: number;
	establishment_date: Time | null;
}

export interface Hive {
	hive_id: number;
	apiary_id: number;
	hive_type: string;
	installation_date: Time | null;
	current_status: string;
}

export interface BeeCommunity {
	community_id: number;
	hive_id: number;
	queen_age: number;
	population_estimate: number;
	health_status: string;
}

export interface HoneyHarvest {
	harvest_id: number;
	hive_id: number;
	harvest_date: Time | null;
	quantity: number;
	quality_grade: string;
}

export interface ObservationLog {
	log_id: number;
	hive_id: number;
	observation_date: Time | null;
	description: string;
	recommendations: string;
}

export interface MaintenancePlan {
	plan_id: number;
	apiary_id: number;
	planned_date: Time | null;
	work_type: string;
	assigned_to: number;
}

export interface Incident {
	incident_id: number;
	hive_id: number;
	incident_date: Time | null;
	description: string;
	severity: string;
	actions_taken: string;
}

export interface ProductionReport {
	report_id: number;
	apiary_id: number;
	start_date: Time | null;
	end_date: Time | null;
	total_honey: number;
	total_expenses: number;
}

export interface VeterinaryPassport {
	passport_id: number;
	bee_community_id: number;
	issue_date: Time | null;
	health_status: string;
	last_inspection_date: Time | null;
}

export interface VeterinaryRecord {
	record_id: number;
	passport_id: number;
	record_date: Time | null;
	description: string;
	treatment: string;
}

export interface WeatherData {
	weather_id: number;
	region_id: number;
	date: Time | null;
	temperature: number;
	humidity: number;
	wind_speed: number;
	precipitation: number;
}

export interface Sensor {
	sensor_id: number;
	hive_id: number;
	sensor_type: string;
	last_reading: Uint8Array;
	last_reading_time: Time | null;
}

export interface SensorReading {
	reading_id: number;
	sensor_id: number;
	value: Uint8Array;
	timestamp: Time | null;
}
