<script lang="ts">
	import type { Hive, Sensor } from '$lib/types';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import { debounce } from 'lodash-es';
	import { onMount } from 'svelte';
	import { getSensors, createSensor } from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import Modal from '../modals/Modal.svelte';

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
	let showModal = false;

	const toastStore = getToastStore();

	// Form data for new sensor
	let formData = {
		hive_id: selectedHive.hive_id,
		sensor_type: 'temperature'
	};

	onMount(async () => {
		await loadSensors();
	});

	async function loadSensors() {
		try {
			loading = true;
			sensors = await getSensors();
			// Filter for current hive's sensors only
			sensors = sensors.filter((s) => s.hive_id === selectedHive.hive_id);
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

	async function handleCreateSensor() {
		try {
			await createSensor({
				hive_id: formData.hive_id,
				sensor_type: formData.sensor_type,
				last_reading: new Uint8Array([0]), // Initial reading of 0
				last_reading_time: new Date().toISOString() // Current timestamp
			});
			showModal = false;
			await loadSensors();
			toastStore.trigger({
				message: '✨ Sensor created successfully!',
				background: 'variant-filled-success'
			});
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to create sensor',
				background: 'variant-filled-error'
			});
		}
	}

	function handleAddSensor() {
		showModal = true;
	}

	$: filteredSensors = sensors
		.filter((s) => s.hive_id === selectedHive.hive_id)
		.filter((s) => s.sensor_type.toLowerCase().includes(searchQuery.toLowerCase()))
		.sort((a, b) => {
			if (sortBy === 'lastReadingTime') {
				// Handle null case by treating null as oldest (smallest timestamp)
				if (!a.last_reading_time) return 1;
				if (!b.last_reading_time) return -1;
				if (!a.last_reading_time && !b.last_reading_time) return 0;
				return new Date(b.last_reading_time).getTime() - new Date(a.last_reading_time).getTime();
			}
			return a.sensor_type.localeCompare(b.sensor_type);
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
				Monitor IoT sensors for {selectedHive.hive_type} Hive
			</p>
		</div>

		<button
			class="bg-amber-500 text-white px-6 py-3 rounded-full
            hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
            flex items-center gap-2"
			on:click={handleAddSensor}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Add Sensor
		</button>
	</div>

	<!-- Add loading state -->
	{#if loading}
		<div class="flex justify-center items-center h-64">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500" />
		</div>
	{:else if error}
		<div class="text-center text-red-500 p-4">
			{error}
		</div>
	{:else}
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
								icon={sensor.sensor_type === 'temperature' ? 'mdi:thermometer' : 'mdi:gauge'}
								class="w-6 h-6 text-amber-500"
							/>
							<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
								{sensor.sensor_type.charAt(0).toUpperCase() + sensor.sensor_type.slice(1)}
							</h2>
						</div>

						<div class="mt-4 space-y-2">
							<div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
								<Icon icon="mdi:clock-outline" class="w-5 h-5" />
								<span
									>Last updated: {sensor.last_reading_time
										? new Date(sensor.last_reading_time).toLocaleString()
										: 'Never'}</span
								>
							</div>
							<div class="text-2xl font-bold text-amber-600">
								{sensor.last_reading ? decodeSensorReading(sensor.last_reading) : '--'}°C
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
	{/if}
</div>

<!-- Add Modal component -->
<Modal isOpen={showModal} title="Add New Sensor" onClose={() => (showModal = false)}>
	<form on:submit|preventDefault={handleCreateSensor} class="space-y-4">
		<div>
			<label for="sensor_type" class="block text-sm font-medium mb-1">Sensor Type</label>
			<select
				id="sensor_type"
				bind:value={formData.sensor_type}
				class="w-full px-4 py-2 rounded-lg border"
				required
			>
				<option value="temperature">Temperature</option>
				<option value="humidity">Humidity</option>
				<option value="weight">Weight</option>
			</select>
		</div>
		<div class="flex justify-end gap-2">
			<button
				type="button"
				class="px-4 py-2 rounded-lg border"
				on:click={() => (showModal = false)}
			>
				Cancel
			</button>
			<button type="submit" class="px-4 py-2 rounded-lg bg-amber-500 text-white hover:bg-amber-600">
				Create
			</button>
		</div>
	</form>
</Modal>
