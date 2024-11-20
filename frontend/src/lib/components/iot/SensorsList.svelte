<script lang="ts">
	import type { Hive, Sensor } from '$lib/types';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import { debounce } from 'lodash-es';
	import { onMount } from 'svelte';
	import { getSensorsByHive, createSensor } from '$lib/services/api';

	export let selectedHive: Hive;
	export let onBack: () => void;
	export let onSelect: (apiary: Sensor) => void;

	// State management following pattern from BeeCommunityManagement.svelte
	let searchQuery = '';
	let sortBy: 'sensorType' | 'lastReadingTime' = 'lastReadingTime';
	let currentPage = 1;
	const itemsPerPage = 12;

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	// Mock data following pattern from HivesList.svelte
	let sensors: Sensor[] = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		try {
			sensors = await getSensorsByHive(selectedHive.hiveId);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load sensors';
		} finally {
			loading = false;
		}
	});

	async function handleCreateSensor() {
		// Implement create sensor modal/form
	}

	$: filteredSensors = sensors
		.filter((s) => s.hiveId === selectedHive.hiveId)
		.filter((s) => s.sensorType.toLowerCase().includes(searchQuery.toLowerCase()))
		.sort((a, b) => {
			if (sortBy === 'lastReadingTime') {
				return new Date(b.lastReadingTime).getTime() - new Date(a.lastReadingTime).getTime();
			}
			return a.sensorType.localeCompare(b.sensorType);
		});

	$: paginatedSensors = filteredSensors.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredSensors.length / itemsPerPage);
</script>

<!-- Layout pattern following ApiarySelect.svelte -->
<div class="max-w mx-auto" in:fade>
	<button
		class="mb-6 flex items-center gap-2 text-amber-600 hover:text-amber-700 transition-colors"
		on:click={onBack}
	>
		<Icon icon="mdi:arrow-left" class="w-5 h-5" />
		<span class="font-medium">Back to Hive</span>
	</button>

	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:chip" class="w-8 h-8 text-amber-500" />
				Hive Sensors
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				Monitor IoT sensors for {selectedHive.hiveType} Hive
			</p>
		</div>

		<button
			class="bg-amber-500 text-white px-6 py-3 rounded-full
            hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
            flex items-center gap-2"
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Add Sensor
		</button>
	</div>

	<!-- Grid layout following ApiarySelect.svelte -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
		{#each paginatedSensors as sensor}
			<button
				class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
                hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
                hover:border-amber-400 overflow-hidden"
				on:click={() => onSelect(sensor)}
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
			</button>
		{/each}
	</div>

	<!-- Pagination controls following ApiarySelect.svelte -->
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
</div>
