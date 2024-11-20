<script lang="ts">
	import type { Apiary, Hive } from '$lib/types';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import { debounce } from 'lodash-es';

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

	// Mock data
	const hives: Hive[] = [
		{
			hiveId: 1,
			apiaryId: selectedApiary.id,
			hiveType: 'Langstroth',
			currentStatus: 'active',
			inspectionDate: new Date()
		}
	];

	$: filteredHives = hives
		.filter((h) => h.apiaryId === selectedApiary.id)
		.filter(
			(h) =>
				h.hiveId.toString().includes(searchQuery.toLowerCase()) ||
				h.currentStatus.toLowerCase().includes(searchQuery.toLowerCase())
		)
		.sort((a, b) => {
			// FIXME
			// if (sortBy === 'currentStatus') {
			// 	return a.currentStatus.localeCompare(b.currentStatus);
			// }
			if (sortBy === 'hiveId') {
				return a.hiveId - b.hiveId;
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
				Manage hives in {selectedApiary.name}
			</p>
		</div>

		<button
			class="bg-amber-500 text-white px-6 py-3 rounded-full
            hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
            flex items-center gap-2"
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Add Hive
		</button>
	</div>

	<!-- Grid layout -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
		{#each paginatedHives as hive}
			<button
				class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
                hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
                hover:border-amber-400 overflow-hidden"
				on:click={() => onSelect(hive)}
			>
				<div class="relative">
					<div class="flex items-center gap-3 mb-3">
						<Icon icon="mdi:beehive-outline" class="w-6 h-6 text-amber-500" />
						<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
							Hive #{hive.hiveId}
						</h2>
					</div>

					<div class="mt-4 space-y-2">
						<div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
							<Icon icon="mdi:clock-outline" class="w-5 h-5" />
							<span>Last inspection: {new Date(hive.inspectionDate).toLocaleDateString()}</span>
						</div>
						<div class="flex items-center gap-2">
							<Icon
								icon={hive.currentStatus === 'active' ? 'mdi:check-circle' : 'mdi:alert-circle'}
								class="w-5 h-5 {hive.currentStatus === 'active'
									? 'text-green-500'
									: 'text-amber-500'}"
							/>
							<span class="capitalize">{hive.currentStatus}</span>
						</div>
					</div>
				</div>
			</button>
		{/each}
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
</div>
