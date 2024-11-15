<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { page } from '$app/stores';
	import type { SensorReading } from '$lib/types';

	const readingId = $page.params.readingId;

	// Mock data
	const reading: SensorReading = {
		readingId: parseInt(readingId),
		sensorId: 1,
		value: new Uint8Array([25]),
		timestamp: new Date(),
		sensorType: 'temperature'
	};

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
		<span class="font-medium">Back to Readings</span>
	</button>

	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon
					icon={reading.sensorType === 'temperature' ? 'mdi:thermometer' : 'mdi:gauge'}
					class="w-8 h-8 text-amber-500"
				/>
				Sensor Reading Details
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				Reading ID: {reading.readingId}
			</p>
		</div>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
		<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
			<h2 class="text-xl font-semibold mb-4">Reading Information</h2>
			<div class="space-y-4">
				<div>
					<label class="text-sm text-gray-500">Sensor Type</label>
					<p class="text-lg font-medium capitalize">{reading.sensorType}</p>
				</div>
				<div>
					<label class="text-sm text-gray-500">Value</label>
					<p class="text-2xl font-bold text-amber-600">
						{new TextDecoder().decode(reading.value)}°C
					</p>
				</div>
				<div>
					<label class="text-sm text-gray-500">Timestamp</label>
					<p class="text-lg font-medium">
						{new Date(reading.timestamp).toLocaleString()}
					</p>
				</div>
				<div>
					<label class="text-sm text-gray-500">Sensor ID</label>
					<p class="text-lg font-medium">{reading.sensorId}</p>
				</div>
			</div>
		</div>

		<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
			<h2 class="text-xl font-semibold mb-4">Historical Context</h2>
			<div class="h-[300px] flex items-center justify-center text-gray-500">
				Reading history chart placeholder
			</div>
		</div>
	</div>
</div>
