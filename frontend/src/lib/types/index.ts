export type Time = string | Date;

export type Role = 'ADMIN' | 'WORKER' | 'MANAGER';

export interface LoginResponse {
	user: User;
	token: string;
}

export interface RegisterInput {
	email: string;
	password: string;
	full_name: string;
	username: string;
}

export interface User {
	user_id: number;
	username: string;
	full_name: string;
	role: Role;
	email: string;
	last_login?: Time;
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
	establishment_date: Time;
}

export interface Hive {
	hive_id: number;
	apiary_id: number;
	hive_type: string;
	installation_date: Time;
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
	harvest_date: Time;
	quantity: number;
	quality_grade: string;
}

export interface ObservationLog {
	log_id: number;
	hive_id: number;
	observation_date: Time;
	description: string;
	recommendations?: string;
}

export interface CreateObservationInput {
	hive_id: number;
	observation_date: string;
	description: string;
	recommendations?: string;
}

export interface MaintenancePlan {
	plan_id: number;
	apiary_id: number;
	planned_date: Date | null;
	work_type: string;
	assigned_to: number;
	status: 'Pending' | 'In Progress' | 'Completed';
	priority: 'High' | 'Medium' | 'Low';
	description: string;
}

export interface Incident {
	incident_id: number;
	hive_id: number;
	incident_date: Time;
	description: string;
	severity: string;
	actions_taken?: string;
}

export interface CreateIncidentInput {
	hive_id: number;
	incident_date: string;
	description: string;
	severity: string;
	actions_taken?: string;
}

export interface ProductionReport {
	report_id: number;
	apiary_id: number;
	start_date: Time;
	end_date: Time;
	total_honey_produced: number;
	total_expenses: number;
	curated_by: number;
}

export interface VeterinaryPassport {
	passport_id: number;
	bee_community_id: number;
	issue_date: Time;
	health_status: string;
	last_inspection_date: Time;
}

export interface VeterinaryRecord {
	record_id: number;
	passport_id: number;
	record_date: Time;
	description: string;
	treatment: string;
}

export interface WeatherData {
	weather_id: number;
	region_id: number;
	date: Time;
	temperature: number;
	humidity: number;
	wind_speed: number;
	precipitation: number;
}

export interface Sensor {
	sensor_id: number;
	hive_id: number;
	sensor_type: string;
	last_reading: string;
	last_reading_time: Time;
}

export interface SensorReading {
	reading_id: number;
	sensor_id: number;
	value: string;
	timestamp: Time;
}

export interface RegionFormData {
	name: string;
	climate_zone: string;
}

export interface WorkerGroup {
	group_id: number;
	manager_id: number;
	group_name: string;
	created_at: Time | null;
	updated_at: Time | null;
}

export interface WorkerGroupMember {
	id: number;
	group_id: number;
	worker_id: number;
	joined_at: Time | null;
}
