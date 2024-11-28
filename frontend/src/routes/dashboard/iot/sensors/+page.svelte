<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { Sensor } from '$lib/types';
	import { onMount } from 'svelte';
	import { getSensors } from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';

	const toastStore = getToastStore();
	let loading = true;
	let error = '';
	let sensors: Sensor[] = [];

	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 12;
	let selectedTimeRange = '24h';

	onMount(async () => {
		await loadSensors();
	});

	async function loadSensors() {
		try {
			loading = true;
			sensors = (await getSensors()) || [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load sensors';
			toastStore.trigger({
				message: '❌ Failed to load sensors',
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

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	$: filteredSensors = sensors
		.filter((s) => s.sensor_type.toLowerCase().includes(searchQuery.toLowerCase()))
		.sort((a, b) => {
			const bTime = b.last_reading_time ? new Date(b.last_reading_time).getTime() : 0;
			const aTime = a.last_reading_time ? new Date(a.last_reading_time).getTime() : 0;
			return bTime - aTime;
		});

	$: paginatedSensors = filteredSensors.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredSensors.length / itemsPerPage);
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
					<Icon icon="mdi:chip" class="w-8 h-8 text-amber-500" />
					All Sensors
				</h1>
				<p class="mt-2 text-gray-600 dark:text-gray-400">Monitor all IoT sensors across hives</p>
			</div>
		</div>

		<!-- Search -->
		<div class="relative">
			<Icon
				icon="mdi:magnify"
				class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
			/>
			<input
				type="text"
				placeholder="Search sensors..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each paginatedSensors as sensor}
				<a
					href="/dashboard/iot/sensors/{sensor.sensor_id}"
					class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
						hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
						hover:border-amber-400 overflow-hidden"
				>
					<div class="relative">
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

						<div class="mt-4 space-y-2">
							<div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
								<Icon icon="mdi:clock-outline" class="w-5 h-5" />
								<span>
									Last updated: {sensor.last_reading_time
										? new Date(sensor.last_reading_time).toLocaleString()
										: 'Never'}
								</span>
							</div>
							<div class="text-2xl font-bold text-amber-600">
								{formatSensorValue(sensor)}
							</div>
						</div>
					</div>
				</a>
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
