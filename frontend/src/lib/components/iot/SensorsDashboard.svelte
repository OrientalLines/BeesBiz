<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { Hive, Sensor, SensorReading } from '$lib/types';
	import { onMount } from 'svelte';
	import { getSensors, getSensorReadings, createSensor } from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import Modal from '../modals/Modal.svelte';

	export let selectedHive: Hive;
	export let onBack: () => void;
	export let onSelect: (sensor: Sensor) => void;

	const toastStore = getToastStore();
	let loading = true;
	let error = '';
	let sensors: Sensor[] = [];
	let readings: SensorReading[] = [];
	let showModal = false;

	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 9;
	let selectedTimeRange = '24h';

	// Form data for new sensor
	let formData = {
		hive_id: selectedHive.hive_id,
		sensor_type: 'temperature',
		initial_reading: 0
	};

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			const [sensorsResponse, readingsResponse] = await Promise.all([
				getSensors(),
				getSensorReadings()
			]);

			// Handle null/empty responses
			sensors = sensorsResponse || [];
			readings = readingsResponse || [];

			// Filter for current hive's sensors only
			sensors = sensors.filter((s) => s.hive_id === selectedHive.hive_id);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load sensors';
			toastStore.trigger({
				message: '❌ Failed to load sensors data',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	async function handleAddSensor() {
		try {
			// Convert the number to a string and then to Base64
			const value = formData.initial_reading.toString();
			const encoder = new TextEncoder();
			const bytes = encoder.encode(value);
			const base64 = btoa(String.fromCharCode(...bytes));

			const payload = {
				hive_id: formData.hive_id,
				sensor_type: formData.sensor_type,
				last_reading: base64,
				last_reading_time: new Date().toISOString()
			};

			const uint8Array = new Uint8Array(
				atob(payload.last_reading)
					.split('')
					.map((c) => c.charCodeAt(0))
			);
			const finalPayload = {
				...payload,
				last_reading: uint8Array
			};

			await createSensor(finalPayload);
			showModal = false;
			await loadData();
			toastStore.trigger({
				message: '✨ Sensor created successfully!',
				background: 'variant-filled-success'
			});
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to create sensor',
				background: 'variant-filled-error'
			});
		}
	}

	// Filter sensors based on search query
	$: filteredSensors = sensors.filter((sensor) =>
		sensor.sensor_type.toLowerCase().includes(searchQuery.toLowerCase())
	);

	// Get readings for the selected time range
	$: timeFilteredReadings = readings.filter((reading) => {
		if (!reading.timestamp) return false;
		const readingDate = new Date(reading.timestamp);
		const now = new Date();
		switch (selectedTimeRange) {
			case '1h':
				return now.getTime() - readingDate.getTime() <= 60 * 60 * 1000;
			case '24h':
				return now.getTime() - readingDate.getTime() <= 24 * 60 * 60 * 1000;
			case '7d':
				return now.getTime() - readingDate.getTime() <= 7 * 24 * 60 * 60 * 1000;
			case '30d':
				return now.getTime() - readingDate.getTime() <= 30 * 24 * 60 * 60 * 1000;
			default:
				return true;
		}
	});

	// Paginate filtered sensors
	$: paginatedSensors = filteredSensors.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredSensors.length / itemsPerPage);

	// Helper function to decode Base64 sensor reading
	function decodeSensorReading(reading: string | null): string {
		if (!reading) return '--';
		try {
			// Decode Base64 to bytes
			const binaryString = atob(reading);
			// Convert binary string to number
			const value = new TextDecoder().decode(
				new Uint8Array(binaryString.split('').map((c) => c.charCodeAt(0)))
			);
			return value;
		} catch (e) {
			console.error('Failed to decode sensor reading:', e);
			return '--';
		}
	}

	// Helper function to format sensor value with units
	function formatSensorValue(sensor: Sensor): string {
		const value = decodeSensorReading(sensor.last_reading);
		switch (sensor.sensor_type) {
			case 'temperature':
				return `${value}°C`;
			case 'humidity':
				return `${value}%`;
			case 'weight':
				return `${value}kg`;
			default:
				return value;
		}
	}
</script>

<div class="max-w mx-auto" in:fade>
	{#if loading}
		<div class="flex justify-center items-center h-64">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500" />
		</div>
	{:else if error}
		<div class="text-center text-red-500 p-4">
			{error}
		</div>
	{:else}
		<div class="flex justify-between items-center mb-8">
			<div>
				<button
					class="flex items-center text-amber-600 hover:text-amber-700 mb-4"
					on:click={onBack}
				>
					<Icon icon="mdi:arrow-left" class="w-5 h-5 mr-1" />
					Back
				</button>
				<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
					<Icon icon="mdi:chip" class="w-8 h-8 text-amber-500" />
					IoT Monitoring
				</h1>
				<p class="mt-2 text-gray-600 dark:text-gray-400">
					Monitor sensors for Hive {selectedHive.hive_id}
				</p>
			</div>
			<button
				class="bg-amber-500 text-white px-6 py-3 rounded-full
					hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
					flex items-center gap-2"
				on:click={() => (showModal = true)}
			>
				<Icon icon="mdi:plus" class="w-5 h-5" />
				Add Sensor
			</button>
		</div>

		{#if sensors.length === 0}
			<!-- Empty state -->
			<div class="text-center py-12 bg-white dark:bg-gray-800 rounded-xl shadow-lg">
				<Icon icon="mdi:sensor" class="w-16 h-16 text-amber-500 mx-auto mb-4" />
				<h2 class="text-2xl font-semibold text-gray-900 dark:text-white mb-2">No Sensors Found</h2>
				<p class="text-gray-600 dark:text-gray-400 mb-6">
					Get started by adding your first sensor to this hive.
				</p>
				<button
					class="bg-amber-500 text-white px-6 py-3 rounded-full
						hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
						inline-flex items-center gap-2"
					on:click={() => (showModal = true)}
				>
					<Icon icon="mdi:plus" class="w-5 h-5" />
					Add Your First Sensor
				</button>
			</div>
		{:else}
			<!-- Controls -->
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
				<div class="relative">
					<Icon
						icon="mdi:magnify"
						class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
					/>
					<input
						type="text"
						placeholder="Search sensors..."
						class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
							bg-white dark:bg-gray-800 text-gray-900 dark:text-white
							placeholder-gray-500 dark:placeholder-gray-400
							focus:ring-2 focus:ring-amber-500 focus:border-transparent"
						on:input={(e) => debouncedSearch(e.currentTarget.value)}
					/>
				</div>
				<select
					bind:value={selectedTimeRange}
					class="w-full px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
							bg-white dark:bg-gray-800 text-gray-900 dark:text-white
							placeholder-gray-500 dark:placeholder-gray-400"
				>
					<option value="1h">Last Hour</option>
					<option value="24h">Last 24 Hours</option>
					<option value="7d">Last 7 Days</option>
					<option value="30d">Last 30 Days</option>
				</select>
			</div>

			<!-- Sensors Grid -->
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each paginatedSensors as sensor}
					<button
						class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
							hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
							hover:border-amber-400"
						on:click={() => onSelect(sensor)}
					>
						<div class="flex items-center gap-3 mb-3">
							<Icon
								icon={sensor.sensor_type === 'temperature'
									? 'mdi:thermometer'
									: sensor.sensor_type === 'humidity'
										? 'mdi:water-percent'
										: 'mdi:scale'}
								class="w-6 h-6 text-amber-500"
							/>
							<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
								{sensor.sensor_type.charAt(0).toUpperCase() + sensor.sensor_type.slice(1)}
							</h2>
						</div>
						<div class="text-gray-600 dark:text-gray-400">
							<div class="flex items-center gap-2">
								<Icon icon="mdi:clock-outline" class="w-5 h-5" />
								<span>
									Last updated: {sensor.last_reading_time
										? new Date(sensor.last_reading_time).toLocaleString()
										: 'Never'}
								</span>
							</div>
							<div class="mt-2 text-2xl font-bold text-amber-600">
								{formatSensorValue(sensor)}
							</div>
						</div>
					</button>
				{/each}
			</div>

			<!-- Pagination -->
			{#if totalPages > 1}
				<div class="mt-8 flex justify-center gap-2">
					<button
						class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
							disabled:opacity-50 disabled:cursor-not-allowed"
						disabled={currentPage === 1}
						on:click={() => currentPage--}
					>
						Previous
					</button>

					<span class="px-4 py-2">
						Page {currentPage} of {totalPages}
					</span>

					<button
						class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
							disabled:opacity-50 disabled:cursor-not-allowed"
						disabled={currentPage === totalPages}
						on:click={() => currentPage++}
					>
						Next
					</button>
				</div>
			{/if}
		{/if}
	{/if}
</div>

<!-- Add Sensor Modal -->
<Modal
	isOpen={showModal}
	onClose={() => (showModal = false)}
	title={sensors.length === 0 ? 'Add Your First Sensor' : 'Add New Sensor'}
>
	<form on:submit|preventDefault={handleAddSensor} class="space-y-4">
		<div>
			<label for="sensor_type" class="block text-sm font-medium mb-1">Sensor Type</label>
			<select
				id="sensor_type"
				bind:value={formData.sensor_type}
				class="w-full px-4 py-2 rounded-lg border"
				required
			>
				<option value="temperature">Temperature</option>
				<option value="humidity">Humidity</option>
				<option value="weight">Weight</option>
			</select>
		</div>
		<div>
			<label for="initial_reading" class="block text-sm font-medium mb-1">Initial Reading</label>
			<input
				type="number"
				id="initial_reading"
				bind:value={formData.initial_reading}
				class="w-full px-4 py-2 rounded-lg border"
				required
				step="0.1"
				min="0"
			/>
		</div>
		<div class="flex justify-end gap-2">
			<button
				type="button"
				class="px-4 py-2 rounded-lg border"
				on:click={() => (showModal = false)}
			>
				Cancel
			</button>
			<button type="submit" class="px-4 py-2 rounded-lg bg-amber-500 text-white hover:bg-amber-600">
				{sensors.length === 0 ? 'Create First Sensor' : 'Create Sensor'}
			</button>
		</div>
	</form>
</Modal>
