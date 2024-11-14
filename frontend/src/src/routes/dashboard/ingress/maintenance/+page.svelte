<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import type { MaintenanceTask } from '$lib/types';
	import MaintenanceTaskModal from '$lib/components/modals/MaintenanceTaskModal.svelte';

	let tasks: MaintenanceTask[] = [
		{
			id: '1',
			hiveId: '1',
			task: 'Hive Inspection',
			dueDate: new Date('2024-03-20'),
			status: 'Pending',
			priority: 'High'
		}
	];

	let showAddModal = false;
	let editingTask: MaintenanceTask | null = null;

	function handleSaveTask(task: MaintenanceTask) {
		if (editingTask) {
			const index = tasks.findIndex((t) => t.id === task.id);
			if (index !== -1) {
				tasks[index] = task;
			}
		} else {
			task.id = (Math.max(...tasks.map((t) => Number(t.id))) + 1).toString();
			tasks = [...tasks, task];
		}
		editingTask = null;
		showAddModal = false;
	}
	function addTask() {
		// Implementation
	}
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Maintenance Schedule</h1>
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

	<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg overflow-hidden">
		<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Task</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Hive ID</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Due Date</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Priority</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Status</th
					>
				</tr>
			</thead>
			<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
				{#each tasks as task}
					<tr>
						<td class="px-6 py-4 text-gray-900 dark:text-gray-300">{task.task}</td>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300"
							>{task.hiveId}</td
						>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
							{new Date(task.dueDate).toLocaleDateString()}
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
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
	{#if showAddModal}
		<MaintenanceTaskModal
			isOpen={true}
			task={null}
			onClose={() => (showAddModal = false)}
			onSave={handleSaveTask}
		/>
	{/if}

	{#if editingTask}
		<MaintenanceTaskModal
			isOpen={true}
			task={editingTask}
			onClose={() => (editingTask = null)}
			onSave={handleSaveTask}
		/>
	{/if}
</div>
