<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import type { Incident } from '$lib/types';
	import IncidentEditModal from '$lib/components/modals/IncidentEditModal.svelte';
	import {
		getIncidents,
		createIncident,
		updateIncidentStatus,
		updateIncident
	} from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import { onMount } from 'svelte';

	const toastStore = getToastStore();
	let incidents: Incident[] = [];
	let loading = true;
	let showAddModal = false;
	let editingIncident: Incident | null = null;

	onMount(async () => {
		await loadIncidents();
	});

	async function loadIncidents() {
		try {
			loading = true;
			incidents = await getIncidents();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load incidents',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	async function handleSaveIncident(incident: Incident) {
		try {
			if (editingIncident) {
				await updateIncident(incident);
				const index = incidents.findIndex((i) => i.incident_id === incident.incident_id);
				if (index !== -1) {
					incidents[index] = incident;
					incidents = [...incidents];
				}
			} else {
				const newIncident = await createIncident({
					hive_id: incident.hive_id,
					incident_date: new Date(incident.incident_date).toISOString(),
					description: incident.description,
					severity: incident.severity,
					actions_taken: incident.actions_taken
				});
				incidents = [...incidents, newIncident];
			}
			toastStore.trigger({
				message: `✨ Incident ${editingIncident ? 'updated' : 'created'} successfully`,
				background: 'variant-filled-success'
			});
		} catch (error) {
			toastStore.trigger({
				message: `❌ Failed to ${editingIncident ? 'update' : 'create'} incident`,
				background: 'variant-filled-error'
			});
		} finally {
			editingIncident = null;
			showAddModal = false;
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString(undefined, {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	function getSeverityClass(severity: string): string {
		switch (severity.toLowerCase()) {
			case 'critical':
				return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300';
			case 'high':
				return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-300';
			case 'medium':
				return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300';
			case 'low':
				return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300';
			default:
				return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-300';
		}
	}
</script>

<div class="h-full flex flex-col space-y-6 overflow-hidden" in:fade>
	<!-- Header - Fixed at top -->
	<div class="flex-shrink-0 flex justify-between items-center">
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

	<!-- Scrollable Content Area -->
	<div class="flex-1 overflow-y-auto min-h-0 h-full">
		{#if loading}
			<div class="flex justify-center py-8">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500" />
			</div>
		{:else if incidents.length === 0}
			<!-- Empty State -->
			<div class="bg-white dark:bg-gray-800 rounded-lg p-8 text-center space-y-4">
				<div class="flex justify-center">
					<Icon
						icon="mdi:clipboard-text-outline"
						class="w-16 h-16 text-gray-400 dark:text-gray-600"
					/>
				</div>
				<h3 class="text-lg font-medium text-gray-900 dark:text-white">No Incidents Reported</h3>
				<p class="text-gray-500 dark:text-gray-400 max-w-sm mx-auto">
					There are currently no incident reports. Click the button above to report a new incident.
				</p>
				<button
					class="mt-4 inline-flex items-center gap-2 px-6 py-3
							bg-amber-500 text-white rounded-full
							hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl"
					on:click={() => (showAddModal = true)}
				>
					<Icon icon="mdi:plus" class="w-5 h-5" />
					Report First Incident
				</button>
			</div>
		{:else}
			<!-- Grid Layout -->
			<div class="grid gap-4">
				<!-- Header -->
				<div
					class="grid grid-cols-[1fr_1fr_1fr_2fr_1fr] gap-6 p-6 bg-white dark:bg-gray-800 rounded-lg items-center font-medium"
				>
					<div class="text-sm text-gray-500 dark:text-gray-400">Date</div>
					<div class="text-sm text-gray-500 dark:text-gray-400">Hive ID</div>
					<div class="text-sm text-gray-500 dark:text-gray-400">Severity</div>
					<div class="text-sm text-gray-500 dark:text-gray-400">Description</div>
					<div class="text-sm text-gray-500 dark:text-gray-400">Actions</div>
				</div>

				<!-- Incidents -->
				{#each incidents as incident}
					<div
						class="grid grid-cols-[1fr_1fr_1fr_2fr_1fr] gap-6 p-6 bg-white dark:bg-gray-800 rounded-lg items-center"
					>
						<div class="text-gray-900 dark:text-gray-300">
							{new Date(incident.incident_date).toLocaleDateString()}
						</div>
						<div class="text-gray-900 dark:text-gray-300">
							{incident.hive_id}
						</div>
						<div>
							<span
								class="px-2 py-1 text-xs font-semibold rounded-full
									{getSeverityClass(incident.severity)}"
							>
								{incident.severity}
							</span>
						</div>
						<div class="text-gray-900 dark:text-gray-300">
							{incident.description}
						</div>
						<div class="text-gray-900 dark:text-gray-300 text-sm">
							{incident.actions_taken || 'No actions recorded'}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Modals -->
	{#if showAddModal || editingIncident}
		<IncidentEditModal
			isOpen={true}
			incident={editingIncident}
			onClose={() => {
				showAddModal = false;
				editingIncident = null;
			}}
			onSave={handleSaveIncident}
		/>
	{/if}
</div>
