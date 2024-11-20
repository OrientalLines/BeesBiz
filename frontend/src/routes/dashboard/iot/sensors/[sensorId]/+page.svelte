<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { page } from '$app/stores';
	import type { Sensor } from '$lib/types';

	// Get sensor ID from route params
	const sensorId = $page.params.sensorId;

	// Mock data - replace with API call
	const sensor: Sensor = {
		sensorId: parseInt(sensorId),
		hiveId: 1,
		sensorType: 'temperature',
		lastReading: new Uint8Array([25]),
		lastReadingTime: new Date()
	};

	let selectedTimeRange = '24h';

	function handleBack() {
		history.back();
	}
</script>

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
					icon={sensor.sensorType === 'temperature' ? 'mdi:thermometer' : 'mdi:gauge'}
					class="w-8 h-8 text-amber-500"
				/>
				{sensor.sensorType.charAt(0).toUpperCase() + sensor.sensorType.slice(1)} Sensor
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				Sensor ID: {sensor.sensorId}
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

	<!-- Sensor readings chart would go here -->
	<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
		<div class="h-[400px] flex items-center justify-center text-gray-500">
			Sensor readings chart placeholder
		</div>
	</div>

	<!-- Latest readings -->
	<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
		<h2 class="text-xl font-semibold mb-4">Latest Readings</h2>
		<div class="text-2xl font-bold text-amber-600">
			{new TextDecoder().decode(sensor.lastReading)}°C
		</div>
		<div class="text-gray-600 dark:text-gray-400 mt-2">
			Last updated: {new Date(sensor.lastReadingTime).toLocaleString()}
		</div>
	</div>
</div>
