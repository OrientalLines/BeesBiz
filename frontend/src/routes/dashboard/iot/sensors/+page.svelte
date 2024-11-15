<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { Sensor } from '$lib/types';

	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 12;
	let selectedTimeRange = '24h';

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	// Mock data - replace with API call
	const sensors: Sensor[] = [
		{
			sensorId: 1,
			hiveId: 1,
			sensorType: 'temperature',
			lastReading: new Uint8Array([25]),
			lastReadingTime: new Date()
		}
	];

	$: filteredSensors = sensors
		.filter((s) => s.sensorType.toLowerCase().includes(searchQuery.toLowerCase()))
		.sort((a, b) => new Date(b.lastReadingTime).getTime() - new Date(a.lastReadingTime).getTime());

	$: paginatedSensors = filteredSensors.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredSensors.length / itemsPerPage);
</script>

<div class="max-w mx-auto space-y-6" in:fade>
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:chip" class="w-8 h-8 text-amber-500" />
				All Sensors
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Monitor all IoT sensors across hives</p>
		</div>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
		{#each paginatedSensors as sensor}
			<a
				href="/dashboard/iot/sensors/{sensor.sensorId}"
				class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
                hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
                hover:border-amber-400 overflow-hidden"
			>
				<div class="relative">
					<div class="flex items-center gap-3 mb-3">
						<Icon
							icon={sensor.sensorType === 'temperature' ? 'mdi:thermometer' : 'mdi:gauge'}
							class="w-6 h-6 text-amber-500"
						/>
						<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
							{sensor.sensorType.charAt(0).toUpperCase() + sensor.sensorType.slice(1)}
						</h2>
					</div>

					<div class="mt-4 space-y-2">
						<div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
							<Icon icon="mdi:clock-outline" class="w-5 h-5" />
							<span>Last updated: {new Date(sensor.lastReadingTime).toLocaleString()}</span>
						</div>
						<div class="text-2xl font-bold text-amber-600">
							{new TextDecoder().decode(sensor.lastReading)}°C
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
</div>
