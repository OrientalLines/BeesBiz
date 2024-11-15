<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { HoneyHarvest } from '$lib/types';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import HarvestEditModal from '$lib/components/modals/HarvestEditModal.svelte';

	let showAddModal = false;
	let editingHarvest: HoneyHarvest | null = null;
	// State management
	let searchQuery = '';
	let sortBy: 'harvestDate' | 'quantity' | 'qualityGrade' = 'harvestDate';
	let filterQuality = 'all';
	let currentPage = 1;
	const itemsPerPage = 20;
	let dateRange = {
		start: null as Date | null,
		end: null as Date | null
	};

	// Debounced search function
	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	// Mock data - replace with API call
	let harvests: HoneyHarvest[] = [
		{
			harvestId: 1,
			hiveId: 101,
			harvestDate: new Date('2024-03-15'),
			quantity: 25.5,
			qualityGrade: 'A'
		}
		// Add more mock data as needed
	];

	function handleSaveHarvest(harvest: HoneyHarvest) {
		if (editingHarvest) {
			const index = harvests.findIndex((h) => h.harvestId === harvest.harvestId);
			if (index !== -1) {
				harvests[index] = harvest;
			}
		} else {
			harvest.harvestId = Math.max(...harvests.map((h) => h.harvestId)) + 1;
			harvests = [...harvests, harvest];
		}
		editingHarvest = null;
		showAddModal = false;
	}

	// Filter and sort logic
	$: filteredHarvests = harvests
		.filter((h) => {
			const matchesSearch =
				h.hiveId.toString().includes(searchQuery) ||
				h.qualityGrade.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesQuality = filterQuality === 'all' || h.qualityGrade === filterQuality;
			const matchesDate =
				(!dateRange.start || new Date(h.harvestDate) >= dateRange.start) &&
				(!dateRange.end || new Date(h.harvestDate) <= dateRange.end);
			return matchesSearch && matchesQuality && matchesDate;
		})
		.sort((a, b) => {
			// TODO: Fix This comparison appears to be unintentional because the types '"harvestDate"' and '"quantity"' have no overlap.
			// if (sortBy === 'quantity') return b.quantity - a.quantity;
			if (sortBy === 'harvestDate')
				return new Date(b.harvestDate).getTime() - new Date(a.harvestDate).getTime();
			if (sortBy === 'qualityGrade') return a.qualityGrade.localeCompare(b.qualityGrade);
			return 0; // Default case
		});

	$: paginatedHarvests = filteredHarvests.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredHarvests.length / itemsPerPage);

	// Calculate total harvest quantity
	$: totalHarvest = filteredHarvests.reduce((sum, h) => sum + h.quantity, 0);

	const qualityColors = {
		A: 'bg-green-100 text-green-800',
		B: 'bg-blue-100 text-blue-800',
		C: 'bg-yellow-100 text-yellow-800',
		D: 'bg-red-100 text-red-800'
	};

	function handleDateRangeChange(start: Date | null, end: Date | null) {
		dateRange.start = start;
		dateRange.end = end;
	}
</script>

<div class="max-w mx-auto" in:fade>
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:honey" class="w-8 h-8 text-amber-500" />
				Honey Harvests
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				Track and manage honey production across all hives
			</p>
		</div>
		<button
			class="bg-amber-500 text-white px-6 py-3 rounded-full
			hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
			flex items-center gap-2"
			on:click={() => (showAddModal = true)}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			New Harvest
		</button>
	</div>

	<!-- Stats Cards -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
		<div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
			<div class="flex items-center gap-3">
				<Icon icon="mdi:chart-line" class="w-6 h-6 text-amber-500" />
				<h3 class="text-lg font-semibold">Total Harvest</h3>
			</div>
			<p class="mt-2 text-3xl font-bold">{totalHarvest.toFixed(1)} kg</p>
		</div>
		<!-- Add more stat cards as needed -->
	</div>

	<!-- Controls -->
	<div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
		<div class="relative">
			<Icon
				icon="mdi:magnify"
				class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
			/>
			<input
				type="text"
				placeholder="Search harvests..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                    bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                    placeholder-gray-500 dark:placeholder-gray-400
                    focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<div class="md:col-span-2">
			<DatePicker
				startDate={dateRange.start}
				endDate={dateRange.end}
				placeholder="Filter by date range"
				onChange={handleDateRangeChange}
			/>
		</div>

		<select
			bind:value={filterQuality}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                    bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                    placeholder-gray-500 dark:placeholder-gray-400
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="all">All Qualities</option>
			<option value="A">Grade A</option>
			<option value="B">Grade B</option>
			<option value="C">Grade C</option>
			<option value="D">Grade D</option>
		</select>
	</div>

	<!-- Data Grid -->
	<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
		<div class="overflow-x-auto">
			<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
				<thead class="bg-gray-50 dark:bg-gray-700">
					<tr>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Harvest Date
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Hive ID
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Quantity (kg)
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Quality Grade
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each paginatedHarvests as harvest}
						<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{new Date(harvest.harvestDate).toLocaleDateString()}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{harvest.hiveId}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{harvest.quantity.toFixed(1)}
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span
									class="px-2 py-1 text-xs font-semibold rounded-full
                                        {qualityColors[harvest.qualityGrade]}"
								>
									Grade {harvest.qualityGrade}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								<button
									class="text-amber-600 hover:text-amber-700"
									on:click={() => (editingHarvest = harvest)}
								>
									Edit
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>

	<!-- Pagination -->
	{#if totalPages > 1}
		<div class="mt-4 flex justify-between items-center">
			<div class="text-sm text-gray-700 dark:text-gray-300">
				Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(
					currentPage * itemsPerPage,
					filteredHarvests.length
				)} of {filteredHarvests.length} entries
			</div>
			<div class="flex gap-2">
				<button
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                        disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={currentPage === 1}
					on:click={() => currentPage--}
				>
					Previous
				</button>
				<button
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                        disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={currentPage === totalPages}
					on:click={() => currentPage++}
				>
					Next
				</button>
			</div>
		</div>
	{/if}

	{#if showAddModal}
		<HarvestEditModal
			isOpen={true}
			harvest={null}
			onClose={() => (showAddModal = false)}
			onSave={handleSaveHarvest}
		/>
	{/if}

	{#if editingHarvest}
		<HarvestEditModal
			isOpen={true}
			harvest={editingHarvest}
			onClose={() => (editingHarvest = null)}
			onSave={handleSaveHarvest}
		/>
	{/if}
</div>
