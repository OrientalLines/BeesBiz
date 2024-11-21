<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { Sensor, SensorReading } from '$lib/types';
	import { onMount } from 'svelte';
	import { getSensors, getSensorReadings } from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';

	const toastStore = getToastStore();
	let loading = true;
	let error = '';
	let readings: SensorReading[] = [];
	let sensors: Sensor[] = [];
	let sensorMap: Map<number, Sensor> = new Map();

	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 20;
	let selectedTimeRange = '24h';

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			const [sensorsData, readingsData] = await Promise.all([
				getSensors(),
				getSensorReadings()
			]);

			sensors = sensorsData || [];
			readings = readingsData || [];

			// Create a map of sensor_id to sensor for easy lookup
			sensorMap = new Map(sensors.map(sensor => [sensor.sensor_id, sensor]));
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load readings';
			toastStore.trigger({
				message: '❌ Failed to load sensor readings',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	function decodeSensorReading(reading: string | null): string {
		if (!reading) return '--';
		try {
			const binaryString = atob(reading);
			const value = new TextDecoder().decode(
				new Uint8Array(binaryString.split('').map((c) => c.charCodeAt(0)))
			);
			return value;
		} catch (e) {
			console.error('Failed to decode sensor reading:', e);
			return '--';
		}
	}

	function formatSensorValue(reading: SensorReading): string {
		const sensor = sensorMap.get(reading.sensor_id);
		if (!sensor) return decodeSensorReading(reading.value);

		const value = decodeSensorReading(reading.value);
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

	function getSensorType(reading: SensorReading): string {
		const sensor = sensorMap.get(reading.sensor_id);
		return sensor?.sensor_type || 'unknown';
	}

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	$: filteredReadings = readings
		.filter((r) => {
			const sensorType = getSensorType(r);
			return sensorType.toLowerCase().includes(searchQuery.toLowerCase());
		})
		.sort((a, b) => {
			const bTime = b.timestamp ? new Date(b.timestamp).getTime() : 0;
			const aTime = a.timestamp ? new Date(a.timestamp).getTime() : 0;
			return bTime - aTime;
		});

	$: paginatedReadings = filteredReadings.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredReadings.length / itemsPerPage);
</script>

<div class="max-w mx-auto space-y-6" in:fade>
	{#if loading}
		<div class="flex justify-center items-center h-64">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500" />
		</div>
	{:else if error}
		<div class="text-center text-red-500 p-4">
			{error}
		</div>
	{:else}
		<div class="flex justify-between items-center">
			<div>
				<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
					<Icon icon="mdi:chart-line" class="w-8 h-8 text-amber-500" />
					Sensor Readings
				</h1>
				<p class="mt-2 text-gray-600 dark:text-gray-400">View all sensor readings across hives</p>
			</div>
		</div>

		<!-- Search and Time Range -->
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<div class="relative">
				<Icon
					icon="mdi:magnify"
					class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
				/>
				<input
					type="text"
					placeholder="Search readings..."
					class="w-full pl-10 pr-4 py-2 rounded-lg border"
					on:input={(e) => debouncedSearch(e.currentTarget.value)}
				/>
			</div>
			<div class="flex justify-end">
				<select
					bind:value={selectedTimeRange}
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
						bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
				>
					<option value="1h">Last Hour</option>
					<option value="24h">Last 24 Hours</option>
					<option value="7d">Last 7 Days</option>
					<option value="30d">Last 30 Days</option>
				</select>
			</div>
		</div>

		<!-- Readings Table -->
		<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
			<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
				<thead class="bg-gray-50 dark:bg-gray-900">
					<tr>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>
							Timestamp
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>
							Sensor Type
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>
							Value
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>
							Sensor ID
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each paginatedReadings as reading}
						<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
								{reading.timestamp ? new Date(reading.timestamp).toLocaleString() : 'Never'}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white capitalize">
								{getSensorType(reading)}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
								{formatSensorValue(reading)}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
								<a
									href="/dashboard/iot/sensors/{reading.sensor_id}"
									class="text-amber-600 hover:text-amber-700"
								>
									{reading.sensor_id}
								</a>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
								<a
									href="/dashboard/iot/readings/{reading.reading_id}"
									class="text-amber-600 hover:text-amber-700"
								>
									<Icon icon="mdi:eye" class="w-5 h-5" />
								</a>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Empty State -->
		{#if readings.length === 0}
			<div class="text-center py-12">
				<Icon icon="mdi:chart-line" class="w-16 h-16 text-gray-400 mx-auto mb-4" />
				<h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Readings Found</h2>
				<p class="text-gray-600 dark:text-gray-400">
					There are no sensor readings available at the moment.
				</p>
			</div>
		{/if}

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
				<span class="px-4 py-2">Page {currentPage} of {totalPages}</span>
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
</div>
