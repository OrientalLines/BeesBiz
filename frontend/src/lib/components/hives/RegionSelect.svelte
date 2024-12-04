<script lang="ts">
	import { onMount } from 'svelte';
	import { getRegions, createRegion, updateRegion, deleteRegion } from '$lib/services/api';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import type { Region } from '$lib/types';
	import RegionEditModal from '../modals/RegionEditModal.svelte';
	import RegionDeleteModal from '../modals/RegionDeleteModal.svelte';
	import RegionAddModal from '../modals/RegionAddModal.svelte';

	export let onSelect: (region: Region) => void;
	const toastStore = getToastStore();

	let regions: Region[] | null = null;
	let loading = true;
	let error = '';
	let showAddModal = false;
	let editingRegion: Region | null = null;
	let deletingRegion: Region | null = null;

	const climateIcons = {
		tropical: 'mdi:weather-sunny',
		temperate: 'mdi:weather-partly-cloudy',
		continental: 'mdi:pine-tree',
		mediterranean: 'mdi:palm-tree'
	} as const;

	onMount(async () => {
		await loadRegions();
	});

	async function loadRegions() {
		try {
			loading = true;
			regions = await getRegions();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load regions';
			toastStore.trigger({
				message: '❌ Failed to load regions',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	async function handleSaveRegion(region: Region) {
		try {
			if (editingRegion) {
				await updateRegion(region);
				toastStore.trigger({
					message: '✨ Region updated successfully!',
					background: 'variant-filled-success'
				});
			} else {
				await createRegion(region);
				toastStore.trigger({
					message: '✨ Region created successfully!',
					background: 'variant-filled-success'
				});
			}
			editingRegion = null;
			showAddModal = false;
			await loadRegions();
		} catch (e) {
			toastStore.trigger({
				message: `❌ Failed to ${editingRegion ? 'update' : 'create'} region`,
				background: 'variant-filled-error'
			});
		}
	}

	async function handleDelete(regionId: number) {
		if (!deletingRegion) return;

		deletingRegion = null;
		await deleteRegion(regionId);
		toastStore.trigger({
			message: '✨ Region deleted successfully!',
			background: 'variant-filled-success'
		});
		await loadRegions();
	}
</script>

<div class="max-w mx-auto" in:fade>
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-4xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-3">
				<Icon icon="mdi:map-marker-radius" class="w-8 h-8 text-amber-500" />
				Regions
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Select a region to view its apiaries</p>
		</div>
		<button
			class="px-4 py-2 bg-amber-500 text-white rounded-lg hover:bg-amber-600
            transition-colors flex items-center gap-2"
			on:click={() => (showAddModal = true)}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Add Region
		</button>
	</div>

	{#if loading}
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each Array(6) as _}
				<div class="animate-pulse bg-gray-200 dark:bg-gray-700 rounded-xl h-32" />
			{/each}
		</div>
	{:else if error}
		<div class="text-red-500 text-center p-4">{error}</div>
	{:else if !regions || regions.length === 0}
		<div class="text-center p-8 bg-gray-50 dark:bg-gray-800/50 rounded-xl border-2 border-dashed border-gray-200 dark:border-gray-700">
			<Icon icon="mdi:map-marker-off" class="w-16 h-16 mx-auto text-gray-400 dark:text-gray-600 mb-4" />
			<h3 class="text-xl font-semibold text-gray-700 dark:text-gray-300 mb-2">No Regions Yet</h3>
			<p class="text-gray-500 dark:text-gray-400 mb-4">Start by adding your first region to manage your apiaries</p>
			<button
				class="px-4 py-2 bg-amber-500 text-white rounded-lg hover:bg-amber-600 transition-colors flex items-center gap-2 mx-auto"
				on:click={() => (showAddModal = true)}
			>
				<Icon icon="mdi:plus" class="w-5 h-5" />
				Add First Region
			</button>
		</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each regions as region}
				<div class="group relative">
					<button
						class="w-full flex flex-col p-5 bg-white dark:bg-gray-800
                        rounded-xl transition-all duration-300 hover:ring-2 hover:ring-amber-500
                        border border-gray-100 dark:border-gray-700"
						on:click={() => onSelect(region)}
					>
						<div class="flex items-start justify-between mb-3">
							<div class="flex items-center gap-3">
								<div
									class="w-10 h-10 rounded-lg bg-amber-100 dark:bg-amber-900/30
                                    flex items-center justify-center"
								>
									<Icon
										icon={climateIcons[region.climate_zone as keyof typeof climateIcons] || 'mdi:map-marker'}
										class="w-6 h-6 text-amber-600 dark:text-amber-400"
									/>
								</div>
								<div class="text-left">
									<h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
										{region.name}
									</h2>
									<p class="text-sm text-gray-500 dark:text-gray-400 capitalize">
										{region.climate_zone} Zone
									</p>
								</div>
							</div>
							<Icon
								icon="mdi:chevron-right"
								class="w-5 h-5 text-gray-400 group-hover:text-amber-500
                                group-hover:translate-x-1 transition-all"
							/>
						</div>
					</button>
					<div
						class="absolute top-2 right-2 opacity-0 group-hover:opacity-100
                        transition-opacity flex gap-2"
					>
						<button
							class="p-1 rounded-lg bg-amber-100 dark:bg-amber-900/30
                            text-amber-600 dark:text-amber-400 hover:bg-amber-200
                            dark:hover:bg-amber-900/50"
							on:click|stopPropagation={() => (editingRegion = region)}
						>
							<Icon icon="mdi:pencil" class="w-4 h-4" />
						</button>
						<button
							class="p-1 rounded-lg bg-red-100 dark:bg-red-900/30
                            text-red-600 dark:text-red-400 hover:bg-red-200
                            dark:hover:bg-red-900/50"
							on:click|stopPropagation={() => (deletingRegion = region)}
						>
							<Icon icon="mdi:trash" class="w-4 h-4" />
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	{#if showAddModal}
		<RegionAddModal
			isOpen={true}
			onClose={() => (showAddModal = false)}
		/>
	{/if}

	{#if editingRegion}
		<RegionEditModal
			isOpen={true}
			onClose={() => (editingRegion = null)}
		/>
	{/if}

	<RegionDeleteModal
		isOpen={!!deletingRegion}
		region={deletingRegion || null}
		onClose={() => (deletingRegion = null)}
		onConfirm={handleDelete}
	/>
</div>
