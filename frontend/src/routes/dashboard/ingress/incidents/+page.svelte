<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import type { Incident } from '$lib/types';
	import IncidentEditModal from '$lib/components/modals/IncidentEditModal.svelte';

	let incidents: Incident[] = [
		{
			id: '1',
			hiveId: '1',
			date: new Date('2024-03-12'),
			type: 'Disease',
			description: 'Varroa mites detected',
			status: 'open'
		}
	];

	let showAddModal = false;
	let editingIncident: Incident | null = null;

	function handleSaveIncident(incident: Incident) {
		if (editingIncident) {
			const index = incidents.findIndex((i) => i.id === incident.id);
			if (index !== -1) {
				incidents[index] = incident;
			}
		} else {
			incident.id = (Math.max(...incidents.map((i) => Number(i.id))) + 1).toString();
			incidents = [...incidents, incident];
		}
		editingIncident = null;
		showAddModal = false;
	}
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
			<Icon icon="mdi:alert" class="w-8 h-8 text-amber-500" />
			Incident Reports
		</h1>
		<button
			class="bg-amber-500 text-white px-6 py-3 rounded-full
            hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
            flex items-center gap-2"
			on:click={() => (showAddModal = true)}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Report Incident
		</button>
	</div>

	<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg overflow-hidden">
		<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Date</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Hive ID</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Type</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Description</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Status</th
					>
				</tr>
			</thead>
			<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
				{#each incidents as incident}
					<tr>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
							{new Date(incident.date).toLocaleDateString()}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300"
							>{incident.hiveId}</td
						>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300"
							>{incident.type}</td
						>
						<td class="px-6 py-4 text-gray-900 dark:text-gray-300">{incident.description}</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                                {incident.status === 'open'
									? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
									: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'}"
							>
								{incident.status}
							</span>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
	{#if showAddModal}
		<IncidentEditModal
			isOpen={true}
			incident={null}
			onClose={() => (showAddModal = false)}
			onSave={handleSaveIncident}
		/>
	{/if}

	{#if editingIncident}
		<IncidentEditModal
			isOpen={true}
			incident={editingIncident}
			onClose={() => (editingIncident = null)}
			onSave={handleSaveIncident}
		/>
	{/if}
</div>
