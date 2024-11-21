<script lang="ts">
	import { onMount } from 'svelte';
	import { getSensorReadingById } from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import Icon from '@iconify/svelte';
	import { page } from '$app/stores';
	import type { SensorReading } from '$lib/types';
	import { fade } from 'svelte/transition';

	const readingId = $page.params.readingId;

	const toastStore = getToastStore();
	let loading = true;
	let error = '';
	let reading: SensorReading | null = null;

	onMount(async () => {
		try {
			loading = true;
			reading = await getSensorReadingById(parseInt(readingId));
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load reading';
			toastStore.trigger({
				message: '❌ Failed to load reading details',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	});

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
{:else if reading}
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
					<Icon icon={'mdi:gauge'} class="w-8 h-8 text-amber-500" />
					Sensor Reading Details
				</h1>
				<p class="mt-2 text-gray-600 dark:text-gray-400">
					Reading ID: {reading.reading_id}
				</p>
			</div>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
				<h2 class="text-xl font-semibold mb-4">Reading Information</h2>
				<div class="space-y-4">
					<div>
						<label class="text-sm text-gray-500">Value</label>
						<p class="text-2xl font-bold text-amber-600">
							{decodeSensorReading(reading.value)}°C
						</p>
					</div>
					<div>
						<label class="text-sm text-gray-500">Timestamp</label>
						<p class="text-lg font-medium">
							{reading.timestamp ? new Date(reading.timestamp).toLocaleString() : ''}
						</p>
					</div>
					<div>
						<label class="text-sm text-gray-500">Sensor ID</label>
						<p class="text-lg font-medium">{reading.sensor_id}</p>
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
{/if}
