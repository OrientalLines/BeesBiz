export interface User {
	id: string;
	name: string;
	role: 'beekeeper' | 'manager' | 'admin';
	email: string;
}

export interface Region {
    region_id: string;
    name: string;
    climate_zone: 'tropical' | 'temperate' | 'continental' | 'mediterranean';
}

export interface Apiary {
    id: string;
    name: string;
    regionId: string;
    location: string;
    hives: number;
    honey: number;
}

export interface Hive {
    hive_id: string;
    apiaryId: string;
	hive_type: string;
	inspection_date: Date;
	current_status: 'active' | 'inactive';
}


export interface HoneyHarvest {
	id: string;
	hiveId: string;
	date: Date;
	amount: number;
	quality: string;
}

export interface Incident {
	id: string;
	hiveId: string;
	date: Date;
	type: string;
	description: string;
	status: 'open' | 'resolved';
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
