<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { SensorReading } from '$lib/types';

	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 20;
	let selectedTimeRange = '24h';

	// Mock data
	const readings: SensorReading[] = [
		{
			readingId: 1,
			sensorId: 1,
			value: new Uint8Array([25]),
			timestamp: new Date(),
			sensorType: 'temperature'
		}
	];

	$: filteredReadings = readings
		.filter((r) => r.sensorType.toLowerCase().includes(searchQuery.toLowerCase()))
		.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());

	$: paginatedReadings = filteredReadings.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredReadings.length / itemsPerPage);
</script>

<div class="max-w mx-auto space-y-6" in:fade>
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:chart-line" class="w-8 h-8 text-amber-500" />
				Sensor Readings
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">View all sensor readings across hives</p>
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
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
				{#each paginatedReadings as reading}
					<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
							{new Date(reading.timestamp).toLocaleString()}
						</td>
						<td
							class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white capitalize"
						>
							{reading.sensorType}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
							{new TextDecoder().decode(reading.value)}°C
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
							{reading.sensorId}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
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
