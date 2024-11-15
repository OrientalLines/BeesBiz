<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { Sensor, SensorReading } from '$lib/types';

	let searchQuery = '';
	let selectedTimeRange = '24h';

	// Mock data - replace with API calls
	const sensors: Sensor[] = [
		{
			sensorId: 1,
			hiveId: 1,
			sensorType: 'temperature',
			lastReading: new Uint8Array([25]),
			lastReadingTime: new Date()
		}
	];

	const recentReadings: SensorReading[] = [
		{
			readingId: 1,
			sensorId: 1,
			value: new Uint8Array([25]),
			timestamp: new Date(),
			sensorType: 'temperature'
		}
	];

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
	}, 300);

	$: filteredSensors = sensors.filter((s) =>
		s.sensorType.toLowerCase().includes(searchQuery.toLowerCase())
	);

	$: filteredReadings = recentReadings
		.filter((r) => r.sensorType.toLowerCase().includes(searchQuery.toLowerCase()))
		.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
</script>

<div class="max-w mx-auto space-y-8" in:fade>
	<!-- Header -->
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:chip" class="w-8 h-8 text-amber-500" />
				IoT Dashboard
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Monitor all IoT sensors and readings</p>
		</div>
	</div>

	<!-- Controls -->
	<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
		<div class="relative">
			<Icon
				icon="mdi:magnify"
				class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
			/>
			<input
				type="text"
				placeholder="Search sensors and readings..."
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
                    placeholder-gray-500 dark:placeholder-gray-400
            focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="1h">Last Hour</option>
			<option value="24h">Last 24 Hours</option>
			<option value="7d">Last 7 Days</option>
			<option value="30d">Last 30 Days</option>
		</select>
	</div>

	<!-- Quick Stats -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
			<div class="flex items-center gap-3 mb-2">
				<Icon icon="mdi:chip" class="w-6 h-6 text-amber-500" />
				<h2 class="text-lg font-semibold">Active Sensors</h2>
			</div>
			<p class="text-3xl font-bold text-amber-600">{sensors.length}</p>
		</div>
		<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
			<div class="flex items-center gap-3 mb-2">
				<Icon icon="mdi:chart-line" class="w-6 h-6 text-amber-500" />
				<h2 class="text-lg font-semibold">Total Readings</h2>
			</div>
			<p class="text-3xl font-bold text-amber-600">{recentReadings.length}</p>
		</div>
		<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
			<div class="flex items-center gap-3 mb-2">
				<Icon icon="mdi:clock-outline" class="w-6 h-6 text-amber-500" />
				<h2 class="text-lg font-semibold">Last Update</h2>
			</div>
			<p class="text-xl font-medium text-amber-600">
				{new Date(
					Math.max(...sensors.map((s) => new Date(s.lastReadingTime).getTime()))
				).toLocaleString()}
			</p>
		</div>
	</div>

	<!-- Active Sensors Section -->
	<section class="space-y-4">
		<div class="flex justify-between items-center">
			<h2 class="text-2xl font-bold text-gray-900 dark:text-white">Active Sensors</h2>
			<a
				href="/dashboard/iot/sensors"
				class="text-amber-600 hover:text-amber-700 flex items-center gap-1"
			>
				View All
				<Icon icon="mdi:arrow-right" class="w-5 h-5" />
			</a>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each filteredSensors.slice(0, 6) as sensor}
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
						<h3 class="text-xl font-semibold text-gray-900 dark:text-white">
							{sensor.sensorType.charAt(0).toUpperCase() + sensor.sensorType.slice(1)}
						</h3>
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
	</section>

	<!-- Recent Readings Section -->
	<section class="space-y-4">
		<div class="flex justify-between items-center">
			<h2 class="text-2xl font-bold text-gray-900 dark:text-white">Recent Readings</h2>
			<a
				href="/dashboard/iot/readings"
				class="text-amber-600 hover:text-amber-700 flex items-center gap-1"
			>
				View All
				<Icon icon="mdi:arrow-right" class="w-5 h-5" />
			</a>
		</div>
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
					{#each filteredReadings.slice(0, 5) as reading}
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
	</section>
</div>
