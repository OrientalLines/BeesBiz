<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { ObservationLog } from '$lib/types';

	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 20;
	let dateRange = {
		start: '',
		end: ''
	};

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	// Mock data - replace with API call
	const logs: ObservationLog[] = [
		{
			logId: 1,
			hiveId: 101,
			observationDate: new Date(),
			description: 'Healthy brood pattern observed',
			recommendations: 'Continue regular monitoring'
		}
	];

	$: filteredLogs = logs
		.filter((log) => {
			const matchesSearch =
				log.hiveId.toString().includes(searchQuery) ||
				log.description.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesDate =
				(!dateRange.start || new Date(log.observationDate) >= new Date(dateRange.start)) &&
				(!dateRange.end || new Date(log.observationDate) <= new Date(dateRange.end));
			return matchesSearch && matchesDate;
		})
		.sort((a, b) => new Date(b.observationDate).getTime() - new Date(a.observationDate).getTime());

	$: paginatedLogs = filteredLogs.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);
	$: totalPages = Math.ceil(filteredLogs.length / itemsPerPage);
</script>

<div class="max-w mx-auto" in:fade>
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:notebook" class="w-8 h-8 text-amber-500" />
				Observation Logs
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Track and manage hive observations</p>
		</div>
		<button
			class="bg-amber-500 text-white px-6 py-3 rounded-full
                hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
                flex items-center gap-2"
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			New Observation
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
				placeholder="Search observations..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                    focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<input
			type="date"
			bind:value={dateRange.start}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		/>

		<input
			type="date"
			bind:value={dateRange.end}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		/>
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
							Date
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Hive ID
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Description
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Recommendations
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each paginatedLogs as log}
						<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{new Date(log.observationDate).toLocaleDateString()}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{log.hiveId}
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300">
								{log.description}
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300">
								{log.recommendations}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
								<button class="text-amber-600 hover:text-amber-700"> Edit </button>
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
					filteredLogs.length
				)} of {filteredLogs.length} entries
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
</div>
