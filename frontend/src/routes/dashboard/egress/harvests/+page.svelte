<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { HoneyHarvest } from '$lib/types';
	import { onMount } from 'svelte';
	import {
		getHoneyHarvests,
		deleteHoneyHarvest,
		updateHoneyHarvest,
		createHoneyHarvest
	} from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import HarvestEditModal from '$lib/components/modals/HarvestEditModal.svelte';

	const toastStore = getToastStore();
	let loading = true;
	let error = '';
	let harvests: HoneyHarvest[] = [];

	let showModal = false;
	let editingHarvest: HoneyHarvest | null = null;
	let searchQuery = '';
	let sortBy: 'harvest_date' | 'quantity' | 'quality_grade' = 'harvest_date';
	let filterQuality = 'all';
	let currentPage = 1;
	const itemsPerPage = 20;

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			harvests = await getHoneyHarvests();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load harvests';
			toastStore.trigger({
				message: '❌ ' + error,
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	async function handleSaveHarvest(harvest: HoneyHarvest) {
		try {
			if (editingHarvest?.harvest_id) {
				await updateHoneyHarvest({
					...harvest,
					harvest_id: editingHarvest.harvest_id
				});
			} else {
				const { harvest_id, ...newHarvest } = harvest;
				await createHoneyHarvest(newHarvest);
			}
			await loadData();
			showModal = false;
			editingHarvest = null;
			toastStore.trigger({
					message: '✨ Harvest saved successfully!',
					background: 'variant-filled-success'
			});
		} catch (e) {
			console.error('Save error:', e);
			toastStore.trigger({
				message: '❌ Failed to save harvest',
				background: 'variant-filled-error'
			});
		}
	}

	async function handleDelete(harvestId: number) {
		if (confirm('Are you sure you want to delete this harvest?')) {
			try {
				await deleteHoneyHarvest(harvestId);
				await loadData();
				toastStore.trigger({
					message: '✨ Harvest deleted successfully!',
					background: 'variant-filled-success'
				});
			} catch (e) {
				toastStore.trigger({
					message: '❌ Failed to delete harvest',
					background: 'variant-filled-error'
				});
			}
		}
	}

	// Filtering and pagination logic
	$: filteredHarvests = harvests
		.filter((h) => {
			const matchesSearch = searchQuery ? h.hive_id.toString().includes(searchQuery) : true;
			const matchesQuality = filterQuality === 'all' || h.quality_grade === filterQuality;
			return matchesSearch && matchesQuality;
		})
		.sort((a, b) => {
			if (sortBy === 'harvest_date') {
				return new Date(b.harvest_date || 0).getTime() - new Date(a.harvest_date || 0).getTime();
			}
			if (sortBy === 'quantity') {
				return b.quantity - a.quantity;
			}
			return a.quality_grade.localeCompare(b.quality_grade);
		});

	$: paginatedHarvests = filteredHarvests.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredHarvests.length / itemsPerPage);

	function handleEdit(harvest: HoneyHarvest) {
		editingHarvest = { ...harvest };
		showModal = true;
	}
</script>

<div class="container mx-auto p-4 space-y-6">
	<div class="flex justify-between items-center">
		<h1 class="text-2xl font-bold">Honey Harvests</h1>
		<button
			class="px-4 py-2 bg-amber-500 text-white rounded-lg hover:bg-amber-600"
			on:click={() => {
				editingHarvest = null;
				showModal = true;
			}}
		>
			New Harvest
		</button>
	</div>

	<!-- Filters -->
	<div class="flex gap-4 flex-wrap">
		<input
			type="text"
			placeholder="Search by hive ID..."
			bind:value={searchQuery}
			class="px-4 py-2 rounded-lg border"
		/>
		<select bind:value={filterQuality} class="px-4 py-2 rounded-lg border">
			<option value="all">All Qualities</option>
			<option value="A">Grade A</option>
			<option value="B">Grade B</option>
			<option value="C">Grade C</option>
			<option value="D">Grade D</option>
		</select>
		<select bind:value={sortBy} class="px-4 py-2 rounded-lg border">
			<option value="harvest_date">Sort by Date</option>
			<option value="quantity">Sort by Quantity</option>
			<option value="quality_grade">Sort by Quality</option>
		</select>
	</div>

	<!-- Harvests List -->
	{#if loading}
		<div class="text-center py-8">Loading...</div>
	{:else if error}
		<div class="text-center py-8 text-red-500">{error}</div>
	{:else if paginatedHarvests.length === 0}
		<div class="text-center py-8">No harvests found</div>
	{:else}
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each paginatedHarvests as harvest}
				<div class="p-4 border rounded-lg">
					<div class="flex justify-between items-start">
						<div>
							<div class="font-medium">Hive #{harvest.hive_id}</div>
							<div class="text-sm text-gray-600">
								{new Date(harvest.harvest_date || '').toLocaleDateString()}
							</div>
						</div>
						<div class="flex gap-2">
							<button
								class="text-amber-600 hover:text-amber-700"
								on:click={() => handleEdit(harvest)}
							>
								Edit
							</button>
							<button
								class="text-red-600 hover:text-red-700"
								on:click={() => handleDelete(harvest.harvest_id)}
							>
								Delete
							</button>
						</div>
					</div>
					<div class="mt-2">
						<div class="text-lg font-semibold">{harvest.quantity} kg</div>
						<div class="text-sm">Quality: Grade {harvest.quality_grade}</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	<!-- Pagination -->
	{#if totalPages > 1}
		<div class="flex justify-center gap-2 mt-4">
			<button
				class="px-4 py-2 rounded-lg border disabled:opacity-50"
				disabled={currentPage === 1}
				on:click={() => currentPage--}
			>
				Previous
			</button>
			<button
				class="px-4 py-2 rounded-lg border disabled:opacity-50"
				disabled={currentPage === totalPages}
				on:click={() => currentPage++}
			>
				Next
			</button>
		</div>
	{/if}

	<HarvestEditModal
		isOpen={showModal}
		harvest={editingHarvest}
		onClose={() => {
			showModal = false;
			editingHarvest = null;
		}}
		onSave={handleSaveHarvest}
	/>
</div>
