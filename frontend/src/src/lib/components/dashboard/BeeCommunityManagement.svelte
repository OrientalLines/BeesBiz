<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { BeeCommunity } from '$lib/types';

	// State management
	let searchQuery = '';
	let sortBy: 'hiveId' | 'queenAge' | 'populationEstimate' | 'healthStatus' = 'hiveId';
	let filterHealth = 'all';
	let currentPage = 1;
	const itemsPerPage = 20;

	// Debounced search function
	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	// Mock data - replace with API call
	const communities: BeeCommunity[] = [
		{
			communityId: 1,
			hiveId: 101,
			queenAge: 2,
			populationEstimate: 50000,
			healthStatus: 'healthy'
		}
		// Add more mock data as needed
	];

	// Filter and sort logic
	$: filteredCommunities = communities
		.filter(
			(c) =>
				(filterHealth === 'all' || c.healthStatus === filterHealth) &&
				(c.hiveId.toString().includes(searchQuery) ||
					c.healthStatus.toLowerCase().includes(searchQuery.toLowerCase()))
		)
		.sort((a, b) => {
			if (sortBy === 'populationEstimate') return b.populationEstimate - a.populationEstimate;
			if (sortBy === 'queenAge') return b.queenAge - a.queenAge;
			return a[sortBy] > b[sortBy] ? 1 : -1;
		});

	$: paginatedCommunities = filteredCommunities.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredCommunities.length / itemsPerPage);

	const healthStatusColors = {
		healthy: 'bg-green-100 text-green-800',
		weak: 'bg-yellow-100 text-yellow-800',
		critical: 'bg-red-100 text-red-800'
	};
</script>

<div class="max-w mx-auto" in:fade>
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:bee" class="w-8 h-8 text-amber-500" />
				Bee Communities
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				Monitor and manage bee populations across all hives
			</p>
		</div>
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
				placeholder="Search by hive ID or health status..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<select
			bind:value={sortBy}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
            focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="hive_id">Sort by Hive ID</option>
			<option value="queen_age">Sort by Queen Age</option>
			<option value="population">Sort by Population</option>
			<option value="health">Sort by Health Status</option>
		</select>

		<select
			bind:value={filterHealth}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
            focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="all">All Health Statuses</option>
			<option value="healthy">Healthy</option>
			<option value="weak">Weak</option>
			<option value="critical">Critical</option>
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
							Hive ID
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Queen Age (Years)
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Population
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Health Status
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each paginatedCommunities as community}
						<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{community.hiveId}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{community.queenAge}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{community.populationEstimate.toLocaleString()}
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span
									class="px-2 py-1 text-xs font-semibold rounded-full
                                    {healthStatusColors[community.healthStatus]}"
								>
									{community.healthStatus}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								<button class="text-amber-600 hover:text-amber-700"> Update </button>
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
					filteredCommunities.length
				)} of {filteredCommunities.length} entries
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
