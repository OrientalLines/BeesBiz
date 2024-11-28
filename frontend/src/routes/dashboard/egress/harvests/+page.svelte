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

	// Add these helper functions
	function getQualityColor(grade: string) {
		const colors = {
			'A': 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
			'B': 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
			'C': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
			'D': 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
		};
		return colors[grade] || 'bg-gray-100 text-gray-800';
	}
</script>

<div class="container mx-auto p-4 space-y-8">
	<!-- Header Section -->
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold bg-gradient-to-r from-amber-500 to-amber-700 bg-clip-text text-transparent">
				Honey Harvests
			</h1>
			<p class="text-gray-600 dark:text-gray-400 mt-1">Track and manage your honey production</p>
		</div>
		<button
			class="px-6 py-2.5 bg-gradient-to-r from-amber-500 to-amber-600 text-white rounded-lg
				   hover:from-amber-600 hover:to-amber-700 transform hover:scale-105 transition-all duration-200
				   flex items-center gap-2 shadow-lg"
			on:click={() => {
				editingHarvest = null;
				showModal = true;
			}}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			New Harvest
		</button>
	</div>

	<!-- Filters Section -->
	<div class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-md">
		<div class="flex gap-4 flex-wrap">
			<div class="flex-1 min-w-[200px]">
				<div class="relative">
					<Icon icon="mdi:magnify" class="absolute left-3 top-3 text-gray-400 w-5 h-5" />
					<input
						type="text"
						placeholder="Search by hive ID..."
						bind:value={searchQuery}
						class="w-full pl-10 pr-4 py-2.5 rounded-lg border dark:border-gray-700
							   bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-white
							   focus:ring-2 focus:ring-amber-500 focus:border-transparent
							   transition-all duration-300"
					/>
				</div>
			</div>
			<select
				bind:value={filterQuality}
				class="px-4 py-2.5 rounded-lg border dark:border-gray-700
					   bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-white
					   focus:ring-2 focus:ring-amber-500 focus:border-transparent
					   transition-all duration-300"
			>
				<option value="all">All Qualities</option>
				<option value="A">Grade A</option>
				<option value="B">Grade B</option>
				<option value="C">Grade C</option>
				<option value="D">Grade D</option>
			</select>
			<select
				bind:value={sortBy}
				class="px-4 py-2.5 rounded-lg border dark:border-gray-700
					   bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-white
					   focus:ring-2 focus:ring-amber-500 focus:border-transparent
					   transition-all duration-300"
			>
				<option value="harvest_date">Sort by Date</option>
				<option value="quantity">Sort by Quantity</option>
				<option value="quality_grade">Sort by Quality</option>
			</select>
		</div>
	</div>

	<!-- Harvests Grid -->
	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500"></div>
		</div>
	{:else if error}
		<div class="text-center py-12 text-red-500 bg-red-50 dark:bg-red-900/20 rounded-lg">
			<Icon icon="mdi:alert-circle" class="w-12 h-12 mx-auto mb-2" />
			{error}
		</div>
	{:else if paginatedHarvests.length === 0}
		<div class="text-center py-12 text-gray-500 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
			<Icon icon="mdi:beehive-outline" class="w-16 h-16 mx-auto mb-2 opacity-50" />
			No harvests found
		</div>
	{:else}
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each paginatedHarvests as harvest}
				<div class="bg-white dark:bg-gray-800 rounded-xl shadow-md hover:shadow-lg transition-all duration-300
							border border-gray-100 dark:border-gray-700 overflow-hidden">
					<div class="p-5">
						<div class="flex justify-between items-start mb-4">
							<div>
								<div class="text-lg font-semibold flex items-center gap-2">
									<Icon icon="mdi:beehive" class="w-5 h-5 text-amber-500" />
									Hive #{harvest.hive_id}
								</div>
								<div class="text-sm text-gray-500 dark:text-gray-400">
									<Icon icon="mdi:calendar" class="w-4 h-4 inline mr-1" />
									{new Date(harvest.harvest_date || '').toLocaleDateString()}
								</div>
							</div>
							<div class="flex gap-2">
								<button
									class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
									on:click={() => handleEdit(harvest)}
								>
									<Icon icon="mdi:pencil" class="w-5 h-5 text-amber-600" />
								</button>
								<button
									class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
									on:click={() => handleDelete(harvest.harvest_id)}
								>
									<Icon icon="mdi:delete" class="w-5 h-5 text-red-600" />
								</button>
							</div>
						</div>
						<div class="flex justify-between items-center">
							<div class="text-2xl font-bold text-amber-600">
								{harvest.quantity} kg
							</div>
							<span class={`px-3 py-1 rounded-full text-sm font-medium ${getQualityColor(harvest.quality_grade)}`}>
								Grade {harvest.quality_grade}
							</span>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	<!-- Pagination -->
	{#if totalPages > 1}
		<div class="flex justify-center gap-3 mt-8">
			<button
				class="px-4 py-2 rounded-lg border dark:border-gray-700 disabled:opacity-50
					   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors
					   flex items-center gap-2"
				disabled={currentPage === 1}
				on:click={() => currentPage--}
			>
				<Icon icon="mdi:chevron-left" />
				Previous
			</button>
			<div class="flex items-center px-4 font-medium">
				Page {currentPage} of {totalPages}
			</div>
			<button
				class="px-4 py-2 rounded-lg border dark:border-gray-700 disabled:opacity-50
					   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors
					   flex items-center gap-2"
				disabled={currentPage === totalPages}
				on:click={() => currentPage++}
			>
				Next
				<Icon icon="mdi:chevron-right" />
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
