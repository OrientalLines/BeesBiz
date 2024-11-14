export interface User {
	id: string;
	name: string;
	role: 'beekeeper' | 'manager' | 'admin';
	email: string;
}

export interface Role {
	id: string;
	name: 'beekeeper' | 'manager' | 'admin';
	description: string;
	permissions: string[];
	userCount: number;
}

export interface AuditLog {
	id: string;
	userId: string;
	userName: string;
	action: string;
	resource: string;
	timestamp: Date;
	details: string;
}

export interface Region {
	region_id: number;
	name: string;
	climate_zone: 'tropical' | 'temperate' | 'continental' | 'mediterranean';
}

export interface Apiary {
	id: number;
	name: string;
	regionId: number;
	location: string;
	hives: number;
	honey: number;
}

export interface Hive {
	hiveId: number;
	apiaryId: number;
	hiveType: string;
	inspectionDate: Date;
	currentStatus: 'active' | 'inactive';
}

export interface HoneyHarvest {
	harvestId: number;
	hiveId: number;
	harvestDate: Date;
	quantity: number;
	qualityGrade: 'A' | 'B' | 'C' | 'D';
}

export interface Incident {
	id: string;
	hiveId: string;
	date: Date;
	type: string;
	description: string;
	status: 'open' | 'resolved';
}

export interface MaintenanceTask {
	id: string;
	hiveId: string;
	task: string;
	dueDate: Date;
	status: string;
	priority: string;
}

export type BeeCommunity = {
	communityId: number;
	hiveId: number;
	queenAge: number;
	populationEstimate: number;
	healthStatus: 'healthy' | 'weak' | 'critical';
};


export type Report = {
	reportId: number;
	apiaryId: number;
	startDate: Date;
	endDate: Date;
	totalHoney: number;
	totalExpenses: number;
};

export type ObservationLog = {
	logId: number;
	hiveId: number;
	observationDate: Date;
	description: string;
	recommendations: string;
};

export type Sensor = {
	sensorId: number;
	hiveId: number;
	sensorType: string;
	lastReading: Uint8Array;
	lastReadingTime: Date;
};

export type SensorReading = {
	readingId: number;
	sensorId: number;
	value: Uint8Array;
	timestamp: Date;
	sensorType: string;
};
