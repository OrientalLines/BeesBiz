import { auth } from '$lib/stores/auth';
import type { Region, Apiary, Hive, BeeCommunity, HoneyHarvest, ProductionReport, Sensor, SensorReading, WeatherData } from '$lib/types';
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
	fullName: string;
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

export async function getBeeCommunitiesByHiveId(hiveId: number): Promise<BeeCommunity[]> {
	return getResource<BeeCommunity[]>(`bee-community/${hiveId}/bee-communities`);
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

export async function createSensor(data: Omit<Sensor, 'sensor_id'>): Promise<Sensor> {
	return createResource<Sensor>('sensor', data);
}

export async function updateSensor(data: Sensor): Promise<Sensor> {
	return updateResource<Sensor>('sensor', data);
}

export async function deleteSensor(id: number): Promise<void> {
	return deleteResource('sensor', id);
}

// SensorReading Operations
export async function getSensorReadings(): Promise<SensorReading[]> {
	return getResource<SensorReading[]>('sensor-reading');
}

export async function getSensorReadingById(id: number): Promise<SensorReading> {
	return getResource<SensorReading>(`sensor-reading/${id}`);
}

export async function createSensorReading(data: Omit<SensorReading, 'reading_id'>): Promise<SensorReading> {
	return createResource<SensorReading>('sensor-reading', data);
}

export async function updateSensorReading(data: SensorReading): Promise<SensorReading> {
	return updateResource<SensorReading>('sensor-reading', data);
}

export async function deleteSensorReading(id: number): Promise<void> {
	return deleteResource('sensor-reading', id);
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
