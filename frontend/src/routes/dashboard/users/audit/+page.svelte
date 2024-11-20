<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { AuditLog } from '$lib/types';

	let searchQuery = '';
	let filterAction = 'all';
	let startDate: string | null = null;
	let endDate: string | null = null;

	// Mock data - replace with API calls
	const auditLogs: AuditLog[] = [
		{
			id: '1',
			userId: '1',
			userName: 'John Admin',
			action: 'user_create',
			resource: 'users',
			timestamp: new Date(),
			details: 'Created new user account for Jane Beekeeper'
		}
		// Add more mock data as needed
	];

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
	}, 300);

	$: filteredLogs = auditLogs
		.filter((log) => {
			const matchesSearch =
				log.userName.toLowerCase().includes(searchQuery.toLowerCase()) ||
				log.action.toLowerCase().includes(searchQuery.toLowerCase()) ||
				log.details.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesAction = filterAction === 'all' || log.action === filterAction;
			const matchesDateRange =
				(!startDate || new Date(log.timestamp) >= new Date(startDate)) &&
				(!endDate || new Date(log.timestamp) <= new Date(endDate));
			return matchesSearch && matchesAction && matchesDateRange;
		})
		.sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime());
</script>

<div class="max-w mx-auto space-y-8" in:fade>
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:history" class="w-8 h-8 text-amber-500" />
				Audit Log
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Track system activities and changes</p>
		</div>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
		<div class="relative">
			<Icon
				icon="mdi:magnify"
				class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
			/>
			<input
				type="text"
				placeholder="Search audit logs..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                    bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                    focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<select
			bind:value={filterAction}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="all">All Actions</option>
			<option value="user_create">User Creation</option>
			<option value="user_update">User Update</option>
			<option value="user_delete">User Deletion</option>
			<option value="role_change">Role Change</option>
		</select>

		<input
			type="date"
			bind:value={startDate}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		/>

		<input
			type="date"
			bind:value={endDate}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		/>
	</div>

	<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
		<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Timestamp
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						User
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Action
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Details
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
				{#each filteredLogs as log}
					<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							{new Date(log.timestamp).toLocaleString()}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							{log.userName}
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="px-2 py-1 text-xs font-semibold rounded-full
                                bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300"
							>
								{log.action}
							</span>
						</td>
						<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300">
							{log.details}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
