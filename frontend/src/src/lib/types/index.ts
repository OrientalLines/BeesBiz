export interface User {
	id: string;
	name: string;
	role: 'beekeeper' | 'manager' | 'admin';
	email: string;
}

export interface Hive {
	id: string;
	name: string;
	location: string;
	status: 'active' | 'inactive' | 'maintenance';
	lastInspection: Date;
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
