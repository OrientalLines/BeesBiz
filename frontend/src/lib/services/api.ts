import { auth } from '$lib/stores/auth';
import type { Region, Apiary, Hive, BeeCommunity, HoneyHarvest, ProductionReport, Sensor, SensorReading, WeatherData, ObservationLog } from '$lib/types';
import { get } from 'svelte/store';
import type { User } from 'lucide-svelte';

// Types
export type Role = 'ADMIN' | 'WORKER' | 'MANAGER';

// API Base URL
const API_BASE_URL = 'http://localhost:4040';

// Auth Types & Functions
interface LoginResponse {
	user: User;
	token: string;
}

interface RegisterInput {
	email: string;
	password: string;
	full_name: string;
	username: string;
}

export async function login(email: string, password: string): Promise<LoginResponse> {
	const response = await fetch(`${API_BASE_URL}/auth/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email_or_username: email, password })
	});
	handleResponse(response);
	return response.json();
}

export async function register(input: RegisterInput): Promise<{ message: string }> {
	const response = await fetch(`${API_BASE_URL}/auth/register`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(input)
	});
	handleResponse(response);
	return response.json();
}

// Region Operations
export async function getRegions(): Promise<Region[]> {
	return getResource<Region[]>('region');
}

export async function getRegionById(id: number): Promise<Region> {
	return getResource<Region>(`region/${id}`);
}

export async function createRegion(data: Omit<Region, 'region_id'>): Promise<Region> {
	return createResource<Region>('region', data);
}

export async function updateRegion(data: Region): Promise<Region> {
	return updateResource<Region>('region', data);
}

export async function deleteRegion(id: number): Promise<void> {
	return deleteResource('region', id);
}

// Apiary Operations
export async function getApiaries(): Promise<Apiary[]> {
	return getResource<Apiary[]>('apiary');
}

export async function getApiaryById(id: number): Promise<Apiary> {
	return getResource<Apiary>(`apiary/${id}`);
}

export async function createApiary(data: Omit<Apiary, 'apiary_id'>): Promise<Apiary> {
	return createResource<Apiary>('apiary', data);
}

export async function updateApiary(data: Apiary): Promise<Apiary> {
	return updateResource<Apiary>('apiary', data);
}

export async function deleteApiary(id: number): Promise<void> {
	return deleteResource('apiary', id);
}

// Hive Operations
export async function getHives(): Promise<Hive[]> {
	return getResource<Hive[]>('hive');
}

export async function getHivesByApiaryId(apiaryId: number): Promise<Hive[]> {
	return getResource<Hive[]>(`hive/${apiaryId}/hives`);
}

export async function createHive(data: Omit<Hive, 'hive_id'>): Promise<Hive> {
	return createResource<Hive>('hive', data);
}

export async function updateHive(data: Hive): Promise<Hive> {
	return updateResource<Hive>('hive', data);
}

export async function deleteHive(id: number): Promise<void> {
	return deleteResource('hive', id);
}

// BeeCommunity Operations
export async function getBeeCommunities(): Promise<BeeCommunity[]> {
	return getResource<BeeCommunity[]>('bee-community');
}

export async function getBeeCommunityById(id: number): Promise<BeeCommunity> {
	const communities = await getResource<BeeCommunity[]>('bee-community');
	const community = communities.find(c => c.community_id === id);
	if (!community) {
		throw new Error('Bee community not found');
	}
	return community;
}

export async function createBeeCommunity(data: Omit<BeeCommunity, 'community_id'>): Promise<BeeCommunity> {
	return createResource<BeeCommunity>('bee-community', data);
}

export async function updateBeeCommunity(data: BeeCommunity): Promise<BeeCommunity> {
	return updateResource<BeeCommunity>('bee-community', data);
}

export async function deleteBeeCommunity(id: number): Promise<void> {
	return deleteResource('bee-community', id);
}

// HoneyHarvest Operations
export async function getHoneyHarvests(): Promise<HoneyHarvest[]> {
	return getResource<HoneyHarvest[]>('honey-harvest');
}

export async function getHoneyHarvestById(id: number): Promise<HoneyHarvest> {
	return getResource<HoneyHarvest>(`honey-harvest/${id}`);
}

export async function createHoneyHarvest(data: Omit<HoneyHarvest, 'harvest_id'>): Promise<HoneyHarvest> {
	return createResource<HoneyHarvest>('honey-harvest', data);
}

export async function updateHoneyHarvest(data: HoneyHarvest): Promise<HoneyHarvest> {
	return updateResource<HoneyHarvest>('honey-harvest', data);
}

export async function deleteHoneyHarvest(id: number): Promise<void> {
	return deleteResource('honey-harvest', id);
}

// ProductionReport Operations
export async function getProductionReports(): Promise<ProductionReport[]> {
	return getResource<ProductionReport[]>('production-report');
}

export async function getProductionReportById(id: number): Promise<ProductionReport> {
	return getResource<ProductionReport>(`production-report/${id}`);
}

export async function createProductionReport(data: Omit<ProductionReport, 'report_id'>): Promise<ProductionReport> {
	return createResource<ProductionReport>('production-report', data);
}

export async function updateProductionReport(data: ProductionReport): Promise<ProductionReport> {
	return updateResource<ProductionReport>('production-report', data);
}

export async function deleteProductionReport(id: number): Promise<void> {
	return deleteResource('production-report', id);
}

// Sensor Operations
export async function getSensors(): Promise<Sensor[]> {
	return getResource<Sensor[]>('sensor');
}

export async function getSensorById(id: number): Promise<Sensor> {
	return getResource<Sensor>(`sensor/${id}`);
}

export async function createSensor(data: {
	hive_id: number;
	sensor_type: string;
	last_reading: Uint8Array;
	last_reading_time: string;
}): Promise<Sensor> {
	const response = await fetch(`${API_BASE_URL}/api/sensor`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			...data,
			last_reading: Array.from(data.last_reading) // Convert Uint8Array to regular array for JSON
		})
	});
	handleResponse(response);
	return response.json();
}

export async function updateSensor(data: Sensor): Promise<Sensor> {
	return updateResource<Sensor>('sensor', data);
}

export async function deleteSensor(id: number): Promise<void> {
	return deleteResource('sensor', id);
}

// SensorReading Operations
export async function getSensorReadings(): Promise<SensorReading[]> {
	const response = await fetch(`${API_BASE_URL}/api/sensor-reading`, {
		headers: getAuthHeaders()
	});
	handleResponse(response);
	return response.json();
}

export async function getSensorReadingById(readingId: number): Promise<SensorReading | null> {
	if (!readingId || isNaN(readingId)) {
		throw new Error('Invalid reading ID');
	}

	try {
		const response = await fetch(`${API_BASE_URL}/api/sensor-reading/${readingId}`, {
			headers: getAuthHeaders()
		});

		if (!response.ok) {
			if (response.status === 404) {
				return null;
			}
			throw new Error(`Failed to fetch reading: ${response.statusText}`);
		}

		return response.json();
	} catch (error) {
		console.error('Error fetching sensor reading:', error);
		throw error;
	}
}

export async function createSensorReading(data: {
	sensor_id: number;
	value: string; // Base64 encoded string
	timestamp: string;
}): Promise<SensorReading> {
	const response = await fetch(`${API_BASE_URL}/api/sensor-reading`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

export async function updateSensorReading(data: {
	reading_id: number;
	sensor_id: number;
	value: string;
	timestamp: string;
}): Promise<SensorReading> {
	if (!data.reading_id || isNaN(data.reading_id)) {
		throw new Error('Invalid reading ID');
	}

	const response = await fetch(`${API_BASE_URL}/api/sensor-reading/${data.reading_id}`, {
		method: 'PUT',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

export async function deleteSensorReading(readingId: number): Promise<void> {
	if (!readingId || isNaN(readingId)) {
		throw new Error('Invalid reading ID');
	}

	const response = await fetch(`${API_BASE_URL}/api/sensor-reading/${readingId}`, {
		method: 'DELETE',
		headers: getAuthHeaders()
	});
	handleResponse(response);
}

// WeatherData Operations
export async function getWeatherData(): Promise<WeatherData[]> {
	return getResource<WeatherData[]>('weather-data');
}

export async function getWeatherDataById(id: number): Promise<WeatherData> {
	return getResource<WeatherData>(`weather-data/${id}`);
}

export async function createWeatherData(data: Omit<WeatherData, 'weather_id'>): Promise<WeatherData> {
	return createResource<WeatherData>('weather-data', data);
}

export async function updateWeatherData(data: WeatherData): Promise<WeatherData> {
	return updateResource<WeatherData>('weather-data', data);
}

export async function deleteWeatherData(id: number): Promise<void> {
	return deleteResource('weather-data', id);
}

// ObservationLog Operations (using Incident endpoints)
export async function getObservationLogs(): Promise<ObservationLog[]> {
	return getResource<ObservationLog[]>('incident');
}

export async function getObservationLogById(id: number): Promise<ObservationLog> {
	return getResource<ObservationLog>(`incident/${id}`);
}

export async function createObservationLog(data: Omit<ObservationLog, 'log_id'>): Promise<ObservationLog> {
	const payload = {
		hive_id: data.hive_id,
		incident_date: data.observation_date,
		description: data.description,
		severity: 'observation', // To mark it as an observation
		actions_taken: data.recommendations || ''
	};
	return createResource<ObservationLog>('incident', payload);
}

export async function updateObservationLog(data: ObservationLog): Promise<ObservationLog> {
	const payload = {
		incident_id: data.log_id,
		hive_id: data.hive_id,
		incident_date: data.observation_date,
		description: data.description,
		severity: 'observation',
		actions_taken: data.recommendations || ''
	};
	return updateResource<ObservationLog>('incident', payload);
}

export async function deleteObservationLog(id: number): Promise<void> {
	return deleteResource('incident', id);
}

// Generic CRUD functions
async function getResource<T>(endpoint: string): Promise<T> {
	const response = await fetch(`${API_BASE_URL}/api/${endpoint}`, {
		headers: getAuthHeaders()
	});
	handleResponse(response);
	return response.json();
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
async function createResource<T>(endpoint: string, data: any): Promise<T> {
	const response = await fetch(`${API_BASE_URL}/api/${endpoint}`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
async function updateResource<T>(endpoint: string, data: any): Promise<T> {
	const response = await fetch(`${API_BASE_URL}/api/${endpoint}`, {
		method: 'PUT',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

async function deleteResource(endpoint: string, id: number): Promise<void> {
	const response = await fetch(`${API_BASE_URL}/api/${endpoint}/${id}`, {
		method: 'DELETE',
		headers: getAuthHeaders()
	});
	handleResponse(response);
}

// Helper functions
function getAuthHeaders(): HeadersInit {
	const authState = get(auth);
	return {
		'Content-Type': 'application/json',
		'Authorization': `Bearer ${authState?.token || ''}`
	};
}

function handleResponse(response: Response) {
	if (!response.ok) {
		if (response.status === 401) {
			throw new Error('Unauthorized. Please log in again.');
		}
		throw new Error(`Request failed with status ${response.status}`);
	}
}
