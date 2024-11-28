<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { MaintenancePlan } from '$lib/types';
	import MaintenanceTaskModal from '$lib/components/modals/MaintenanceTaskModal.svelte';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import {
		getMaintenanceTasks,
		createMaintenanceTask,
		updateMaintenanceTask
	} from '$lib/services/api';

	let showAddModal = false;
	let editingTask: MaintenancePlan | null = null;
	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 10;
	const toastStore = getToastStore();

	let dateRange = {
		start: null as Date | null,
		end: null as Date | null
	};

	let tasks: MaintenancePlan[] = [];

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	function handleDateRangeChange(start: Date | null, end: Date | null) {
		dateRange = {
			start,
			end
		};
		currentPage = 1;
	}

	async function handleSaveTask(task: MaintenancePlan) {
		try {
			if (editingTask) {
				const index = tasks.findIndex((t) => t.plan_id === task.plan_id);
				if (index !== -1) {
					tasks[index] = task;
					tasks = [...tasks];
				}
			} else {
				task.plan_id = Math.max(...tasks.map((t) => t.plan_id)) + 1;
				tasks = [...tasks, task];
			}
			editingTask = null;
			showAddModal = false;
			toastStore.trigger({
				message: '✨ Task saved successfully!',
				background: 'variant-filled-success'
			});
		} catch (e) {
			console.error('Save error:', e);
			toastStore.trigger({
				message: '❌ Failed to save task',
				background: 'variant-filled-error'
			});
		}
	}

	$: filteredTasks = tasks
		.filter((t) => {
			const matchesSearch =
				t.work_type.toLowerCase().includes(searchQuery.toLowerCase()) ||
				t.apiary_id.toString().includes(searchQuery);
			const matchesDate =
				(!dateRange.start || new Date(t.planned_date!) >= dateRange.start) &&
				(!dateRange.end || new Date(t.planned_date!) <= dateRange.end);
			return matchesSearch && matchesDate;
		})
		.sort((a, b) => new Date(a.planned_date!).getTime() - new Date(b.planned_date!).getTime());

	$: paginatedTasks = filteredTasks.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);
	$: totalPages = Math.ceil(filteredTasks.length / itemsPerPage);

	// Summary calculations
	$: pendingTasks = filteredTasks.filter((t) => t.status === 'Pending').length;
	$: highPriorityTasks = filteredTasks.filter((t) => t.priority === 'High').length;
	$: completedTasks = filteredTasks.filter((t) => t.status === 'Completed').length;

	function handleSearchInput(event: Event) {
		const target = event.target as HTMLInputElement;
		debouncedSearch(target.value);
	}
</script>

<div class="max-w mx-auto" in:fade>
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:calendar-check" class="w-8 h-8 text-amber-500" />
				Maintenance Schedule
			</h1>
			<p class="mt-2 text-sm text-gray-500 dark:text-gray-300">
				Manage and track your maintenance tasks efficiently.
			</p>
		</div>
		<div class="flex items-center gap-4">
			<DatePicker
				startDate={dateRange.start}
				endDate={dateRange.end}
				onChange={(start, end) => handleDateRangeChange(start, end ?? null)}
			/>
			<input
				type="text"
				placeholder="Search tasks..."
				class="w-64 px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-amber-500"
				bind:value={searchQuery}
				on:input={handleSearchInput}
			/>
			<button
				class="bg-amber-500 text-white px-6 py-3 rounded-full
                hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
                flex items-center gap-2"
				on:click={() => (showAddModal = true)}
			>
				<Icon icon="mdi:plus" class="w-5 h-5" />
				Add Task
			</button>
		</div>
	</div>

	<!-- Summary Cards -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
		<div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
			<h3 class="text-lg font-semibold text-gray-900 dark:text-white">Pending Tasks</h3>
			<p class="text-3xl font-bold text-amber-500">{pendingTasks}</p>
		</div>
		<div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
			<h3 class="text-lg font-semibold text-gray-900 dark:text-white">High Priority</h3>
			<p class="text-3xl font-bold text-red-500">{highPriorityTasks}</p>
		</div>
		<div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
			<h3 class="text-lg font-semibold text-gray-900 dark:text-white">Completed</h3>
			<p class="text-3xl font-bold text-green-500">{completedTasks}</p>
		</div>
	</div>

	<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg overflow-hidden">
		<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Task
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Apiary ID
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Due Date
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Priority
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Status
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Actions
					</th>
				</tr>
			</thead>
			<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
				{#each paginatedTasks as task}
					<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
						<td class="px-6 py-4 text-gray-900 dark:text-gray-300">{task.work_type}</td>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300"
							>{task.apiary_id}</td
						>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
							{new Date(task.planned_date!).toLocaleDateString()}
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                                {task.priority === 'High'
									? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
									: task.priority === 'Medium'
										? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
										: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'}"
							>
								{task.priority}
							</span>
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                                {task.status === 'Pending'
									? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'
									: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'}"
							>
								{task.status}
							</span>
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
							<button
								class="text-amber-600 hover:text-amber-700"
								on:click={() => (editingTask = task)}
							>
								Edit
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<!-- Pagination -->
	{#if totalPages > 1}
		<div class="mt-4 flex justify-between items-center">
			<div class="text-sm text-gray-700 dark:text-gray-300">
				Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(
					currentPage * itemsPerPage,
					filteredTasks.length
				)} of {filteredTasks.length} entries
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

	{#if showAddModal || editingTask}
		<MaintenanceTaskModal
			isOpen={true}
			task={editingTask}
			onClose={() => {
				showAddModal = false;
				editingTask = null;
			}}
			onSave={handleSaveTask}
		/>
	{/if}
</div>
