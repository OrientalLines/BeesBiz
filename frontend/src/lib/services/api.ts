import { auth } from '$lib/stores/auth';
import type { Region, Apiary, Hive, BeeCommunity, HoneyHarvest, ProductionReport, Sensor, SensorReading, WeatherData, ObservationLog, User, Incident, CreateIncidentInput, CreateObservationInput, Role, RegisterInput, LoginResponse } from '$lib/types';
import { get } from 'svelte/store';

// API Base URL
const API_BASE_URL = 'http://localhost:4040';

const PROXY_API_BASE_URL = 'http://localhost:4041';

// Helper functions
function getAuthHeaders(): HeadersInit {
	const authState = get(auth);
	if (!authState?.token) {
		console.warn('No auth token available');
		throw new Error('Authentication required');
	}

	return {
		'Content-Type': 'application/json',
		'Authorization': `Bearer ${authState.token}`
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

// Generic CRUD functions
async function getResource<T>(endpoint: string): Promise<T> {
	try {
		const response = await fetch(`${API_BASE_URL}/api/${endpoint}`, {
			headers: getAuthHeaders()
		});
		handleResponse(response);
		return response.json();
	} catch (error) {
		console.error(`Error fetching ${endpoint}:`, error);
		throw error;
	}
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

// User Management
export async function getUsers(): Promise<User[]> {
	return getResource<User[]>('user');
}

export async function getUserById(id: number): Promise<User> {
	return getResource<User>(`user/${id}`);
}

export async function createUser(data: Omit<User, 'user_id'>): Promise<User> {
	const response = await fetch(`${API_BASE_URL}/api/user`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

export async function updateUser(data: User): Promise<User> {
	return updateResource<User>('user', data);
}

export async function deleteUser(id: number): Promise<void> {
	return deleteResource('user', id);
}

// User Role Management
export async function updateUserRole(userId: number, role: Role): Promise<User> {
	const response = await fetch(`${API_BASE_URL}/api/user/role`, {
		method: 'PUT',
		headers: getAuthHeaders(),
		body: JSON.stringify({ user_id: userId, role })
	});
	handleResponse(response);
	return response.json();
}

// Region Access Management
export async function getUserAllowedRegions(userId: number): Promise<Region[]> {
	return getResource<Region[]>(`user/${userId}/allowed-regions`);
}

export async function updateUserAllowedRegions(userId: number, regionIds: number[]): Promise<void> {
	// First, get current allowed regions to determine which ones to remove
	const currentRegions = await getUserAllowedRegions(userId);
	const currentRegionIds = currentRegions.map(r => r.region_id);

	// Remove regions that are no longer selected
	const regionsToRemove = currentRegionIds.filter(id => !regionIds.includes(id));
	for (const regionId of regionsToRemove) {
		const response = await fetch(`${API_BASE_URL}/api/allowed-region/${regionId}`, {
			method: 'DELETE',
			headers: getAuthHeaders()
		});
		handleResponse(response);
	}

	// Add new regions
	const regionsToAdd = regionIds.filter(id => !currentRegionIds.includes(id));
	for (const regionId of regionsToAdd) {
		const payload = {
			user_id: userId,
			region_id: regionId
		};

		const response = await fetch(`${API_BASE_URL}/api/allowed-region`, {
			method: 'POST',
			headers: getAuthHeaders(),
			body: JSON.stringify(payload)
		});
		handleResponse(response);
	}
}

// Replace the old updateRegionAccess function
export async function updateRegionAccess(userId: number, regionIds: number[]): Promise<void> {
	return updateUserAllowedRegions(userId, regionIds);
}

// Add a function to check if a user has access to a specific region
export async function checkRegionAccess(userId: number, regionId: number): Promise<boolean> {
	try {
		const allowedRegions = await getUserAllowedRegions(userId);
		return allowedRegions.some(region => region.region_id === regionId);
	} catch (error) {
		console.error('Error checking region access:', error);
		return false;
	}
}

// Add gRPC-specific functions
export async function getTotalHoneyHarvested(hiveId: number, startDate: string, endDate: string): Promise<number> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/getTotalHoneyHarvested`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			hive_id: hiveId,
			start_date: startDate,
			end_date: endDate
		})
	});
	handleResponse(response);
	const data = await response.json();
	return data.total_honey;
}

export async function addObservation(
	hiveId: number,
	observationDate: string,
	description: string,
	recommendations: string
): Promise<void> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/addObservation`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			hive_id: hiveId,
			observation_date: observationDate,
			description,
			recommendations
		})
	});
	handleResponse(response);
}

export async function getCommunityHealthStatus(communityId: number): Promise<string> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/getCommunityHealthStatus`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			community_id: communityId
		})
	});
	handleResponse(response);
	const data = await response.json();
	return data.health_status;
}

export async function updateHiveStatus(hiveId: number, newStatus: string): Promise<void> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/updateHiveStatus`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			hive_id: hiveId,
			new_status: newStatus
		})
	});
	handleResponse(response);
}

export async function getAverageTemperature(regionId: number, days: number): Promise<number> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/getAverageTemperature`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			region_id: regionId,
			days
		})
	});
	handleResponse(response);
	const data = await response.json();
	return data.avg_temperature;
}

export async function assignMaintenancePlan(planId: number, userId: number): Promise<void> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/assignMaintenancePlan`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			plan_id: planId,
			user_id: userId
		})
	});
	handleResponse(response);
}

export async function registerIncident(
	hiveId: number,
	incidentDate: string,
	description: string,
	severity: string
): Promise<void> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/registerIncident`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			hive_id: hiveId,
			incident_date: incidentDate,
			description,
			severity
		})
	});
	handleResponse(response);
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function getLatestSensorReading(hiveId: number, sensorType: string): Promise<any> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/getLatestSensorReading`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			hive_id: hiveId,
			sensor_type: sensorType
		})
	});
	handleResponse(response);
	return response.json();
}

export async function createProductionReport(
	apiaryId: number,
	startDate: string,
	endDate: string
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
): Promise<any> {
	const response = await fetch(`${PROXY_API_BASE_URL}/api/grpc/createProductionReport`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify({
			apiary_id: apiaryId,
			start_date: startDate,
			end_date: endDate
		})
	});
	handleResponse(response);
	return response.json();
}

export async function searchHives(query: string): Promise<Hive[]> {
	const hives = await getHives();
	const searchTerm = query.toLowerCase();

	return hives.filter(hive =>
		hive.hive_id.toString().includes(searchTerm) ||
		hive.hive_type.toLowerCase().includes(searchTerm) ||
		hive.current_status.toLowerCase().includes(searchTerm)
	);
}


// Update the incident functions
export async function getIncidents(): Promise<Incident[]> {
	return getResource<Incident[]>('incident');
}

export async function createIncident(data: CreateIncidentInput): Promise<Incident> {
	const response = await fetch(`${API_BASE_URL}/api/incident`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

export async function updateIncidentStatus(incidentId: string, status: string): Promise<void> {
	const response = await fetch(`${API_BASE_URL}/api/incident/${incidentId}/status`, {
		method: 'PUT',
		headers: getAuthHeaders(),
		body: JSON.stringify({ severity: status })
	});
	handleResponse(response);
}

export async function updateIncident(data: Incident): Promise<Incident> {
	const response = await fetch(`${API_BASE_URL}/api/incident`, {
		method: 'PUT',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

export async function deleteIncident(id: number): Promise<void> {
	const response = await fetch(`${API_BASE_URL}/api/incident/${id}`, {
		method: 'DELETE',
		headers: getAuthHeaders()
	});
	handleResponse(response);
}

export async function getObservationLogs(): Promise<ObservationLog[]> {
	return getResource<ObservationLog[]>('observation');
}

export async function getObservationLogById(id: number): Promise<ObservationLog> {
	return getResource<ObservationLog>(`observation/${id}`);
}

export async function createObservationLog(data: CreateObservationInput): Promise<ObservationLog> {
	const response = await fetch(`${API_BASE_URL}/api/observation`, {
		method: 'POST',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

export async function updateObservationLog(data: ObservationLog): Promise<ObservationLog> {
	const response = await fetch(`${API_BASE_URL}/api/observation`, {
		method: 'PUT',
		headers: getAuthHeaders(),
		body: JSON.stringify(data)
	});
	handleResponse(response);
	return response.json();
}

export async function deleteObservationLog(id: number): Promise<void> {
	const response = await fetch(`${API_BASE_URL}/api/observation/${id}`, {
		method: 'DELETE',
		headers: getAuthHeaders()
	});
	handleResponse(response);
}
