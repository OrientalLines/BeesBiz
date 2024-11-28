<script lang="ts">
	import { onMount } from 'svelte';
	import { getHivesByApiaryId, createHive, deleteHive } from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import { debounce } from 'lodash-es';
	import Modal from '../modals/Modal.svelte';
	import type { Apiary, Hive } from '$lib/types';

	export let selectedApiary: Apiary;
	export let onBack: () => void;
	export let onSelect: (hive: Hive) => void;

	// State management
	let searchQuery = '';
	let sortBy: 'hiveId' | 'status' = 'hiveId';
	let currentPage = 1;
	const itemsPerPage = 12;

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	const toastStore = getToastStore();
	let showModal = false;
	let loading = false;
	let error = '';

	// Form data for new hive
	let formData: Omit<Hive, 'hive_id'> = {
		apiary_id: selectedApiary.apiary_id,
		hive_type: '',
		current_status: 'active',
		installation_date: new Date().toISOString()
	};

	let hives: Hive[] = [];

	let deletingHive: Hive | null = null;

	onMount(async () => {
		await loadHives();
	});

	async function loadHives() {
		try {
			loading = true;
			const response = await getHivesByApiaryId(selectedApiary.apiary_id);
			hives = response || [];
		} catch (error) {
			console.error('Failed to load hives:', error);
			toastStore.trigger({
				message: '❌ Failed to load hives',
				background: 'variant-filled-error'
			});
			hives = [];
		} finally {
			loading = false;
		}
	}

	async function handleCreateHive() {
		try {
			const payload = {
				...formData,
				installation_date: new Date(formData.installation_date).toISOString()
			};
			await createHive(payload);
			showModal = false;
			await loadHives();
			toastStore.trigger({
				message: '✨ Hive created successfully!',
				background: 'variant-filled-success'
			});
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to create hive',
				background: 'variant-filled-error'
			});
		}
	}

	$: filteredHives = hives
		.filter((h) => h.apiary_id === selectedApiary.apiary_id)
		.filter(
			(h) =>
				h.hive_id.toString().includes(searchQuery.toLowerCase()) ||
				h.current_status.toLowerCase().includes(searchQuery.toLowerCase())
		)
		.sort((a, b) => {
			// FIXME
			// if (sortBy === 'currentStatus') {
			// 	return a.currentStatus.localeCompare(b.currentStatus);
			// }
			if (sortBy === 'hiveId') {
				return a.hive_id - b.hive_id;
			}
			return 0;
		});

	$: paginatedHives = filteredHives.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredHives.length / itemsPerPage);
</script>

<div class="max-w mx-auto" in:fade>
	<button
		class="mb-6 flex items-center gap-2 text-amber-600 hover:text-amber-700 transition-colors"
		on:click={onBack}
	>
		<Icon icon="mdi:arrow-left" class="w-5 h-5" />
		<span class="font-medium">Back to Apiary</span>
	</button>

	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:beehive-outline" class="w-8 h-8 text-amber-500" />
				Hives
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				Manage hives in {selectedApiary.location}
			</p>
		</div>

		<button
			class="bg-amber-500 text-white px-6 py-3 rounded-full
            hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
            flex items-center gap-2"
			on:click={() => (showModal = true)}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Add Hive
		</button>
	</div>

	<!-- Grid layout -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
		{#if loading}
			<div class="flex justify-center items-center h-64">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500" />
			</div>
		{:else if error}
			<div class="text-center text-red-500 p-4">
				{error}
			</div>
		{:else if paginatedHives.length === 0}
			<div class="col-span-full flex flex-col items-center justify-center py-16 px-4">
				<div class="bg-amber-50 dark:bg-amber-900/20 rounded-full p-4 mb-4">
					<Icon icon="mdi:beehive-outline" class="w-12 h-12 text-amber-500" />
				</div>
				<h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">No Hives Found</h3>
				<p class="text-gray-600 dark:text-gray-400 text-center max-w-md mb-6">
					{searchQuery
						? 'No hives match your search criteria. Try adjusting your search.'
						: "This apiary doesn't have any hives yet. Add your first hive to get started!"}
				</p>
				<button
					class="bg-amber-500 text-white px-6 py-3 rounded-full
					hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
					flex items-center gap-2"
					on:click={() => (showModal = true)}
				>
					<Icon icon="mdi:plus" class="w-5 h-5" />
					Add First Hive
				</button>
			</div>
		{:else}
			{#each paginatedHives as hive}
				<div
					class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
					hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
					hover:border-amber-400 overflow-hidden"
				>
					<button
						class="absolute top-2 right-2 p-2 text-gray-400 hover:text-red-500
							opacity-0 group-hover:opacity-100 transition-all duration-300"
						on:click|stopPropagation={(e) => {
							e.stopPropagation();
							deletingHive = hive;
						}}
					>
						<Icon icon="mdi:delete" class="w-5 h-5" />
					</button>

					<button class="w-full text-left" on:click={() => onSelect(hive)}>
						<div class="relative">
							<div class="flex items-center gap-3 mb-3">
								<Icon icon="mdi:beehive-outline" class="w-6 h-6 text-amber-500" />
								<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
									Hive #{hive.hive_id}
								</h2>
							</div>

							<div class="mt-4 space-y-2">
								<div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
									<Icon icon="mdi:clock-outline" class="w-5 h-5" />
									<span
										>Last inspection: {hive.installation_date
											? new Date(hive.installation_date).toLocaleDateString()
											: 'Never'}</span
									>
								</div>
								<div class="flex items-center gap-2">
									<Icon
										icon={hive.current_status === 'active'
											? 'mdi:check-circle'
											: 'mdi:alert-circle'}
										class="w-5 h-5 {hive.current_status === 'active'
											? 'text-green-500'
											: 'text-amber-500'}"
									/>
									<span class="capitalize">{hive.current_status}</span>
								</div>
							</div>
						</div>
					</button>
				</div>
			{/each}
		{/if}
	</div>

	<!-- Pagination controls -->
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

	<!-- Add Modal for creating new hive -->
	<Modal isOpen={showModal} title="Create New Hive" onClose={() => (showModal = false)}>
		<form on:submit|preventDefault={handleCreateHive} class="space-y-4">
			<div>
				<label for="hive_type" class="block text-sm font-medium mb-1">Hive Type</label>
				<select
					id="hive_type"
					bind:value={formData.hive_type}
					class="w-full px-4 py-2 rounded-lg border dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
					required
				>
					<option value="Langstroth">Langstroth</option>
					<option value="Top Bar">Top Bar</option>
					<option value="Warre">Warre</option>
				</select>
			</div>
			<div>
				<label for="current_status" class="block text-sm font-medium mb-1">Status</label>
				<select
					id="current_status"
					bind:value={formData.current_status}
					class="w-full px-4 py-2 rounded-lg border dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
					required
				>
					<option value="active">Active</option>
					<option value="inactive">Inactive</option>
					<option value="maintenance">Maintenance</option>
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
				<button
					type="submit"
					class="px-4 py-2 rounded-lg bg-amber-500 text-white hover:bg-amber-600"
				>
					Create
				</button>
			</div>
		</form>
	</Modal>

	{#if deletingHive}
		<Modal isOpen={true} title="Delete Hive" onClose={() => (deletingHive = null)}>
			<div class="space-y-4" in:fade={{ duration: 300 }}>
				<div class="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg flex items-start gap-3">
					<div class="p-2 bg-red-100 dark:bg-red-900/30 rounded-full">
						<Icon icon="mdi:alert" class="w-6 h-6 text-red-600 dark:text-red-400" />
					</div>
					<div>
						<h3 class="text-lg font-semibold text-red-800 dark:text-red-200">Confirm Deletion</h3>
						<p class="mt-1 text-red-700 dark:text-red-300">
							Are you sure you want to delete Hive #{deletingHive.hive_id}? This action cannot be
							undone.
						</p>
						<div class="mt-2 text-sm text-red-600 dark:text-red-400">
							<p>Hive details:</p>
							<ul class="list-disc list-inside mt-1">
								<li>Type: {deletingHive.hive_type}</li>
								<li>Status: {deletingHive.current_status}</li>
								<li>
									Installation Date: {new Date(deletingHive.installation_date).toLocaleDateString()}
								</li>
							</ul>
						</div>
					</div>
				</div>

				<div class="flex justify-end gap-3">
					<button
						type="button"
						class="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100
							dark:hover:bg-gray-700 rounded-lg transition-colors"
						on:click={() => (deletingHive = null)}
					>
						Cancel
					</button>
					<button
						type="button"
						class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg
							transition-colors flex items-center gap-2"
						on:click={async () => {
							if (deletingHive) {
								try {
									await deleteHive(deletingHive.hive_id);
									await loadHives();
									toastStore.trigger({
										message: '✨ Hive deleted successfully!',
										background: 'variant-filled-success'
									});
								} catch (e) {
									toastStore.trigger({
										message: '❌ Failed to delete hive',
										background: 'variant-filled-error'
									});
								} finally {
									deletingHive = null;
								}
							}
						}}
					>
						<Icon icon="mdi:delete" class="w-5 h-5" />
						Delete Hive
					</button>
				</div>
			</div>
		</Modal>
	{/if}
</div>
