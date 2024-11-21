<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { ObservationLog } from '$lib/types';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import ObservationEditModal from '$lib/components/modals/ObservationEditModal.svelte';
	import {
		getObservationLogs,
		createObservationLog,
		updateObservationLog,
		deleteObservationLog
	} from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import { onMount } from 'svelte';

	const toastStore = getToastStore();
	let loading = true;
	let error = '';
	let logs: ObservationLog[] = [];
	let editingObservation: ObservationLog | null = null;
	let showAddModal = false;
	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 20;
	let dateRange = {
		start: null as Date | null,
		end: null as Date | null
	};

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			logs = await getObservationLogs();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load observation logs';
			toastStore.trigger({
				message: '❌ Failed to load observation logs',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	function handleDateRangeChange(start: Date | null, end: Date | null) {
		dateRange.start = start;
		dateRange.end = end;
		currentPage = 1;
	}

	function handleEdit(observation: ObservationLog) {
		editingObservation = { ...observation };
	}

	async function handleSaveObservation(updatedObservation: ObservationLog) {
		try {
			if (editingObservation) {
				await updateObservationLog(updatedObservation);
				toastStore.trigger({
					message: '✨ Observation updated successfully!',
					background: 'variant-filled-success'
				});
			} else {
				await createObservationLog(updatedObservation);
				toastStore.trigger({
					message: '✨ Observation created successfully!',
					background: 'variant-filled-success'
				});
			}
			await loadData();
			editingObservation = null;
			showAddModal = false;
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to save observation',
				background: 'variant-filled-error'
			});
		}
	}

	async function handleDelete(logId: number) {
		if (confirm('Are you sure you want to delete this observation?')) {
			try {
				await deleteObservationLog(logId);
				await loadData();
				toastStore.trigger({
					message: '✨ Observation deleted successfully!',
					background: 'variant-filled-success'
				});
			} catch (e) {
				toastStore.trigger({
					message: '❌ Failed to delete observation',
					background: 'variant-filled-error'
				});
			}
		}
	}

	$: filteredLogs = logs
		.filter((log) => {
			const matchesSearch =
				log.hive_id.toString().includes(searchQuery) ||
				log.description.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesDate =
				(!dateRange.start || new Date(log.observation_date) >= dateRange.start) &&
				(!dateRange.end || new Date(log.observation_date) <= dateRange.end);
			return matchesSearch && matchesDate;
		})
		.sort((a, b) => {
			const dateA = a.observation_date ? new Date(a.observation_date).getTime() : 0;
			const dateB = b.observation_date ? new Date(b.observation_date).getTime() : 0;
			return dateB - dateA;
		});

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
			on:click={() => (showAddModal = true)}
			class="bg-amber-500 text-white px-6 py-3 rounded-full
                hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
                flex items-center gap-2"
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			New Observation
		</button>
	</div>

	<!-- Controls -->
	<div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
		<div class="relative">
			<Icon
				icon="mdi:magnify"
				class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
			/>
			<input
				type="text"
				placeholder="Search observations..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
					bg-white dark:bg-gray-800 text-gray-900 dark:text-white
					placeholder-gray-500 dark:placeholder-gray-400
					focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<DatePicker
			startDate={dateRange.start}
			endDate={dateRange.end}
			placeholder="Filter by observation date"
			onChange={handleDateRangeChange}
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
								{log.observation_date ? new Date(log.observation_date).toLocaleDateString() : 'N/A'}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{log.hive_id}
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300">
								{log.description}
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300">
								{log.recommendations}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
								<button
									class="text-amber-600 hover:text-amber-700"
									on:click={() => handleEdit(log)}
								>
									<Icon icon="mdi:pencil" class="w-5 h-5" />
								</button>
								<button
									class="text-red-600 hover:text-red-700"
									on:click={() => handleDelete(log.log_id)}
								>
									<Icon icon="mdi:delete" class="w-5 h-5" />
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

{#if editingObservation}
	<ObservationEditModal
		isOpen={true}
		observation={editingObservation}
		onClose={() => (editingObservation = null)}
		onSave={handleSaveObservation}
	/>
{/if}

{#if showAddModal}
	<ObservationEditModal
		isOpen={true}
		observation={null}
		onClose={() => (showAddModal = false)}
		onSave={handleSaveObservation}
	/>
{/if}
