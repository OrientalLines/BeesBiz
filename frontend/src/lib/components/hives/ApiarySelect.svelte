<script lang="ts">
	import type { Region, Apiary } from '$lib/types';
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import { onMount } from 'svelte';
	import { getApiaries, createApiary } from '$lib/services/api';
	import Modal from '$lib/components/modals/Modal.svelte';
	import { getToastStore } from '@skeletonlabs/skeleton';

	export let selectedRegion: Region;
	export let onSelect: (apiary: Apiary) => void;
	export let onBack: () => void;

	const toastStore = getToastStore();
	let showModal = false;

	// Form data matching the Apiary type exactly
	let formData: Omit<Apiary, 'apiary_id'> = {
		location: '',
		manager_id: 0,
		establishment_date: new Date().toISOString()
	};

	// State management
	let searchQuery = '';
	let sortBy: keyof Apiary = 'location';
	let currentPage = 1;
	const itemsPerPage = 12;

	// Debounced search function
	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	let apiaries: Apiary[] = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		await loadApiaries();
	});

	async function loadApiaries() {
		try {
			loading = true;
			apiaries = await getApiaries();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load apiaries';
		} finally {
			loading = false;
		}
	}

	// Filter and sort logic
	$: filteredApiaries = apiaries
		.filter((a) => a.location.toLowerCase().includes(searchQuery.toLowerCase()))
		.sort((a, b) => {
			const aValue = a[sortBy] ?? '';
			const bValue = b[sortBy] ?? '';
			return aValue < bValue ? -1 : aValue > bValue ? 1 : 0;
		});

	// Pagination logic
	$: paginatedApiaries = filteredApiaries.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredApiaries.length / itemsPerPage);

	async function handleCreateApiary() {
		try {
			// Ensure the date is in ISO format before sending
			const payload = {
				...formData,
				establishment_date: formData.establishment_date
					? new Date(formData.establishment_date).toISOString()
					: new Date().toISOString()
			};
			await createApiary(payload);
			showModal = false;
			await loadApiaries();
			toastStore.trigger({
				message: '✨ Apiary created successfully!',
				background: 'variant-filled-success'
			});
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to create apiary',
				background: 'variant-filled-error'
			});
		}
	}
</script>

<div class="max-w mx-auto" in:fade>
	<!-- Separate back button -->
	<button
		class="mb-6 flex items-center gap-2 text-amber-600 hover:text-amber-700 transition-colors"
		on:click={onBack}
	>
		<Icon icon="mdi:arrow-left" class="w-5 h-5" />
		<span class="font-medium">Back to Regions</span>
	</button>

	<!-- Header without back button -->
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:beehive-outline" class="w-8 h-8 text-amber-500" />
				{selectedRegion.name} Apiaries
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Select an apiary to view its hives</p>
		</div>

		<button
			on:click={() => (showModal = true)}
			class="bg-amber-500 text-white px-6 py-3 rounded-full
            hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
            flex items-center gap-2"
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			New Apiary
		</button>
	</div>

	<!-- Add search and filter controls below the header -->
	<div class="mb-6 flex flex-col sm:flex-row gap-4">
		<div class="relative flex-grow">
			<Icon
				icon="mdi:magnify"
				class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
			/>
			<input
				type="text"
				placeholder="Search apiaries..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                    bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                    placeholder-gray-500 dark:placeholder-gray-400
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>
		<select
			bind:value={sortBy}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                    bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                    focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="location">Sort by Location</option>
			<option value="manager_id">Sort by Manager</option>
			<option value="establishment_date">Sort by Date</option>
		</select>
	</div>

	<!-- Update the grid to use paginatedApiaries -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
		{#each paginatedApiaries as apiary}
			<button
				class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
            hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
            hover:border-amber-400 overflow-hidden"
				on:click={() => onSelect(apiary)}
			>
				<!-- Honeycomb background pattern -->
				<div class="absolute inset-0 opacity-5 group-hover:opacity-10 transition-opacity">
					<Icon icon="mdi:hexagon-outline" class="w-full h-full text-amber-500" />
				</div>

				<div class="relative">
					<div class="flex items-center gap-3 mb-3">
						<Icon icon="mdi:beehive-outline" class="w-6 h-6 text-amber-500" />
						<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
							{apiary.apiary_id}
						</h2>
					</div>

					<div class="flex items-center gap-2 text-gray-500 dark:text-gray-400">
						<Icon icon="mdi:map-marker" class="w-5 h-5" />
						<p>{apiary.location}</p>
					</div>

					<!-- Quick stats -->
					<div
						class="mt-4 grid grid-cols-2 gap-4 pt-4 border-t border-gray-100 dark:border-gray-700"
					>
						<div class="flex items-center gap-2">
							<Icon icon="mdi:calendar" class="w-5 h-5 text-amber-500" />
							<span class="text-sm">
								Established on{' '}
								{apiary.establishment_date
									? new Date(apiary.establishment_date).toLocaleDateString()
									: 'not set'}
							</span>
						</div>
						<div class="flex items-center gap-2">
							<Icon icon="mdi:account" class="w-5 h-5 text-amber-500" />
							<span class="text-sm"
								>Manager {apiary.manager_id ? `#${apiary.manager_id}` : 'not set'}
							</span>
						</div>
					</div>
				</div>
			</button>
		{/each}
	</div>

	<!-- Add pagination controls -->
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

	<!-- Add the Modal component for creating new apiary -->
	<Modal isOpen={showModal} title="Create New Apiary" onClose={() => (showModal = false)}>
		<form on:submit|preventDefault={handleCreateApiary} class="space-y-4">
			<div>
				<label for="location" class="block text-sm font-medium mb-1">Location</label>
				<input
					type="text"
					id="location"
					bind:value={formData.location}
					class="w-full px-4 py-2 rounded-lg border dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
					required
				/>
			</div>
			<div>
				<label for="establishment_date" class="block text-sm font-medium mb-1"
					>Establishment Date</label
				>
				<input
					type="datetime-local"
					id="establishment_date"
					bind:value={formData.establishment_date}
					class="w-full px-4 py-2 rounded-lg border dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
					required
				/>
			</div>
			<div class="flex justify-end gap-2">
				<button
					type="button"
					class="px-4 py-2 rounded-lg border"
					on:click={() => (showModal = false)}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 rounded-lg bg-amber-500 text-white hover:bg-amber-600"
				>
					Create
				</button>
			</div>
		</form>
	</Modal>
</div>
