<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { Hive, Sensor } from '$lib/types';

	export let selectedHive: Hive;
	export let onBack: () => void;
	export let onSelect: (sensor: Sensor) => void;

	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 9;
	let selectedTimeRange = '24h';

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	// Mock data - replace with API call
	const sensors: Sensor[] = [
		{
			sensorId: 1,
			hiveId: selectedHive.hiveId,
			sensorType: 'temperature',
			lastReading: new Uint8Array([25]),
			lastReadingTime: new Date()
		}
	];

	$: filteredSensors = sensors.filter((sensor) =>
		sensor.sensorType.toLowerCase().includes(searchQuery.toLowerCase())
	);

	$: paginatedSensors = filteredSensors.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);
</script>

<div class="max-w mx-auto" in:fade>
	<div class="flex justify-between items-center mb-8">
		<div>
			<button class="flex items-center text-amber-600 hover:text-amber-700 mb-4" on:click={onBack}>
				<Icon icon="mdi:arrow-left" class="w-5 h-5 mr-1" />
				Back
			</button>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:chip" class="w-8 h-8 text-amber-500" />
				IoT Monitoring
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				Monitor sensors for Hive {selectedHive.hiveId}
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
			<a
				href="/dashboard/iot/sensors/{sensor.sensorId}"
				class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
                hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
                hover:border-amber-400"
			>
				<div class="flex items-center gap-3 mb-3">
					<Icon
						icon={sensor.sensorType === 'temperature' ? 'mdi:thermometer' : 'mdi:gauge'}
						class="w-6 h-6 text-amber-500"
					/>
					<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
						{sensor.sensorType.charAt(0).toUpperCase() + sensor.sensorType.slice(1)}
					</h2>
				</div>
				<div class="text-gray-600 dark:text-gray-400">
					<div class="flex items-center gap-2">
						<Icon icon="mdi:clock-outline" class="w-5 h-5" />
						<span>Last updated: {new Date(sensor.lastReadingTime).toLocaleString()}</span>
					</div>
					<div class="mt-2 text-2xl font-bold text-amber-600">
						{new TextDecoder().decode(sensor.lastReading)}°C
					</div>
				</div>
			</a>
		{/each}
	</div>
</div>
