<script lang="ts">
	import { onMount } from 'svelte';
	import { getSensorById, getSensorReadings } from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import Icon from '@iconify/svelte';
	import { page } from '$app/stores';
	import type { Sensor, SensorReading } from '$lib/types';
	import { fade } from 'svelte/transition';

	const sensorId = $page.params.sensorId;
	const toastStore = getToastStore();
	let loading = true;
	let error = '';
	let sensor: Sensor | null = null;
	let readings: SensorReading[] = [];
	let selectedTimeRange = '24h';

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

	onMount(async () => {
		try {
			loading = true;
			[sensor, readings] = await Promise.all([
				getSensorById(parseInt(sensorId)),
				getSensorReadings()
			]);

			if (!sensor) throw new Error('Sensor not found');

			// Filter readings for this sensor
			readings = readings?.filter((r) => r.sensor_id === parseInt(sensorId)) || [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load sensor';
			toastStore.trigger({
				message: '❌ Failed to load sensor details',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	});

	function handleBack() {
		history.back();
	}
</script>

{#if loading}
	<div class="flex justify-center items-center h-64">
		<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500" />
	</div>
{:else if error}
	<div class="text-center text-red-500 p-4">
		{error}
	</div>
{:else if sensor}
	<div class="max-w mx-auto space-y-6" in:fade>
		<button
			class="flex items-center gap-2 text-amber-600 hover:text-amber-700 transition-colors"
			on:click={handleBack}
		>
			<Icon icon="mdi:arrow-left" class="w-5 h-5" />
			<span class="font-medium">Back to Sensors</span>
		</button>

		<div class="flex justify-between items-center">
			<div>
				<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
					<Icon
						icon={sensor.sensor_type === 'temperature'
							? 'mdi:thermometer'
							: sensor.sensor_type === 'humidity'
								? 'mdi:water-percent'
								: 'mdi:scale'}
						class="w-8 h-8 text-amber-500"
					/>
					{sensor.sensor_type.charAt(0).toUpperCase() + sensor.sensor_type.slice(1)} Sensor
				</h1>
				<p class="mt-2 text-gray-600 dark:text-gray-400">
					Sensor ID: {sensor.sensor_id}
				</p>
			</div>

			<select
				bind:value={selectedTimeRange}
				class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
					bg-white dark:bg-gray-800 text-gray-900 dark:text-white
					placeholder-gray-500 dark:placeholder-gray-400"
			>
				<option value="1h">Last Hour</option>
				<option value="24h">Last 24 Hours</option>
				<option value="7d">Last 7 Days</option>
				<option value="30d">Last 30 Days</option>
			</select>
		</div>

		<!-- Latest readings -->
		<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
			<h2 class="text-xl font-semibold mb-4">Latest Reading</h2>
			<div class="text-2xl font-bold text-amber-600">
				{formatSensorValue(sensor)}
			</div>
			<div class="text-gray-600 dark:text-gray-400 mt-2">
				Last updated: {sensor.last_reading_time
					? new Date(sensor.last_reading_time).toLocaleString()
					: 'Never'}
			</div>
		</div>

		<!-- Readings History -->
		<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
			<h2 class="text-xl font-semibold mb-4">Reading History</h2>
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead>
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
								Timestamp
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
								Value
							</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
						{#each readings as reading}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
									{reading.timestamp ? new Date(reading.timestamp).toLocaleString() : 'Never'}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
									{decodeSensorReading(reading.value)}
									{sensor.sensor_type === 'temperature'
										? '°C'
										: sensor.sensor_type === 'humidity'
											? '%'
											: 'kg'}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
{/if}
