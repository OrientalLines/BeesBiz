<script lang="ts">
	import Modal from './Modal.svelte';
	import type { Incident, Hive } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import { getHives } from '$lib/services/api';
	import { onMount } from 'svelte';
	import {
		AlertTriangle,
		Home,
		FileText,
		Tag,
		CheckCircle2,
		Bug,
		Cloud,
		Shield,
		AlertOctagon
	} from 'lucide-svelte';

	export let isOpen = false;
	export let incident: Incident | null = null;
	export let onClose = () => {};
	export let onSave = (updatedIncident: Incident) => {};

	const toastStore = getToastStore();
	let isLoading = false;

	// Add state for hives
	let hives: Hive[] = [];
	let filteredHives: Hive[] = [];
	let showHiveSuggestions = false;
	let hiveSearchTerm = '';

	onMount(async () => {
		try {
			hives = await getHives();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load hives',
				background: 'variant-filled-error'
			});
		}
	});

	// Update hive search when term changes
	$: {
		if (hiveSearchTerm) {
			filteredHives = hives
				.filter(
					(hive) =>
						hive.hive_id.toString().includes(hiveSearchTerm) ||
						hive.hive_type.toLowerCase().includes(hiveSearchTerm.toLowerCase())
				)
				.slice(0, 5); // Limit to 5 suggestions
			showHiveSuggestions = true;
		} else {
			filteredHives = [];
			showHiveSuggestions = false;
		}
	}

	function selectHive(hive: Hive) {
		formData.hive_id = hive.hive_id;
		hiveSearchTerm = `#${hive.hive_id} (${hive.hive_type})`;
		showHiveSuggestions = false;
	}

	const incidentTypes = [
		{ value: 'Disease', icon: Bug, color: 'text-red-500', bg: 'bg-red-50 dark:bg-red-900/10' },
		{
			value: 'Pest',
			icon: AlertOctagon,
			color: 'text-orange-500',
			bg: 'bg-orange-50 dark:bg-orange-900/10'
		},
		{ value: 'Weather', icon: Cloud, color: 'text-blue-500', bg: 'bg-blue-50 dark:bg-blue-900/10' },
		{
			value: 'Vandalism',
			icon: Shield,
			color: 'text-purple-500',
			bg: 'bg-purple-50 dark:bg-purple-900/10'
		},
		{
			value: 'Other',
			icon: AlertTriangle,
			color: 'text-gray-500',
			bg: 'bg-gray-50 dark:bg-gray-900/10'
		}
	];

	let formData: Partial<Incident> = {
		incident_id: incident?.incident_id || 0,
		hive_id: incident?.hive_id || 0,
		incident_date: incident?.incident_date || new Date(),
		description: incident?.description || '',
		severity: incident?.severity || 'low',
		actions_taken: incident?.actions_taken || ''
	};

	$: if (incident) {
		formData = { ...incident };
	}

	async function handleSubmit() {
		if (!formData) return;

		isLoading = true;
		try {
			await onSave(formData as Incident);
			toastStore.trigger({
				message: '✨ Incident report saved successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to save incident report',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	function handleDateChange(date: Date | null) {
		if (date) {
			formData.incident_date = date;
		}
	}

	// Add this helper function for status colors
	function getStatusClass(status: string): string {
		switch (status.toLowerCase()) {
			case 'active':
				return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300';
			case 'inactive':
				return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300';
			default:
				return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-300';
		}
	}
</script>

<Modal {isOpen} title={incident ? 'Edit Incident' : 'Report Incident'} {onClose}>
	<div
		in:fly={{ y: 50, duration: 400 }}
		out:fade
		class="max-h-[calc(100vh-10rem)] overflow-y-auto px-6 py-4"
	>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<!-- Alert Header -->
			<div
				class="bg-amber-50 dark:bg-amber-900/10 p-4 rounded-lg flex items-start gap-3"
				in:fly={{ y: 20, duration: 300, delay: 100 }}
			>
				<AlertTriangle class="w-6 h-6 text-amber-500 flex-shrink-0 mt-0.5" />
				<div>
					<h3 class="text-sm font-medium text-amber-800 dark:text-amber-200">Incident Report</h3>
					<p class="mt-1 text-sm text-amber-700 dark:text-amber-300">
						Please provide detailed information about the incident to help us track and resolve
						issues effectively.
					</p>
				</div>
			</div>

			<!-- Date Picker -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 200 }}>
				<div class="flex items-center gap-2">
					<Tag class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Incident Date
					</label>
				</div>
				<DatePicker
					startDate={formData.incident_date
						? formData.incident_date instanceof Date
							? formData.incident_date
							: new Date(formData.incident_date)
						: new Date()}
					onChange={handleDateChange}
				/>
			</div>

			<!-- Hive ID -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 300 }}>
				<div class="flex items-center gap-2">
					<Home class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Hive ID </label>
				</div>
				<div class="relative hive-search">
					<input
						type="text"
						bind:value={hiveSearchTerm}
						placeholder="Search by hive ID or type..."
						class="w-full px-4 py-2 rounded-lg border-2 border-gray-200 dark:border-gray-700
								bg-white dark:bg-gray-800 text-gray-900 dark:text-white
								focus:ring-2 focus:ring-amber-500 focus:border-transparent
								transition-all duration-300"
						on:focus={() => {
							if (hiveSearchTerm) showHiveSuggestions = true;
						}}
					/>

					{#if showHiveSuggestions && filteredHives.length > 0}
						<div
							class="absolute z-50 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-200
								dark:border-gray-700 rounded-lg shadow-lg overflow-hidden"
							transition:fade
						>
							{#each filteredHives as hive}
								<button
									type="button"
									class="w-full px-4 py-2 text-left hover:bg-gray-50 dark:hover:bg-gray-700
										flex items-center justify-between gap-2 transition-colors duration-200"
									on:click={() => selectHive(hive)}
								>
									<div>
										<span class="font-medium text-gray-900 dark:text-white">
											Hive #{hive.hive_id}
										</span>
										<span class="text-sm text-gray-500 dark:text-gray-400 ml-2">
											({hive.hive_type})
										</span>
									</div>
									<span
										class="text-xs px-2 py-1 rounded-full
										{getStatusClass(hive.current_status)}"
									>
										{hive.current_status}
									</span>
								</button>
							{/each}
						</div>
					{/if}

					{#if hiveSearchTerm && !showHiveSuggestions && formData.hive_id}
						<div class="mt-2 text-sm text-gray-500 dark:text-gray-400">
							Selected: Hive #{formData.hive_id} ({hives.find((h) => h.hive_id === formData.hive_id)
								?.hive_type})
						</div>
					{/if}
				</div>
			</div>

			<!-- Incident Type (Severity) -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 400 }}>
				<div class="flex items-center gap-2">
					<AlertTriangle class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Severity Level
					</label>
				</div>
				<div class="grid grid-cols-2 md:grid-cols-3 gap-3">
					{#each incidentTypes as type}
						<button
							type="button"
							class="flex items-center gap-2 p-3 border-2 rounded-lg transition-all duration-300
                                {formData.severity === type.value
								? `border-${type.color} ${type.bg}`
								: 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'}"
							on:click={() => (formData.severity = type.value)}
						>
							<svelte:component this={type.icon} class="w-4 h-4 {type.color}" />
							<span class="text-sm font-medium text-gray-900 dark:text-gray-100">
								{type.value}
							</span>
						</button>
					{/each}
				</div>
			</div>

			<!-- Description and Actions Taken -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 500 }}>
				<div class="flex items-center gap-2">
					<FileText class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Description </label>
				</div>
				<textarea
					bind:value={formData.description}
					rows="4"
					class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                        bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                        focus:ring-2 focus:ring-amber-500 focus:border-transparent
                        transition-all duration-300 resize-none"
					required
				/>
			</div>

			<!-- Actions Taken -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 600 }}>
				<div class="flex items-center gap-2">
					<CheckCircle2 class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Actions Taken
					</label>
				</div>
				<textarea
					bind:value={formData.actions_taken}
					rows="3"
					class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                        bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                        focus:ring-2 focus:ring-amber-500 focus:border-transparent
                        transition-all duration-300 resize-none"
				/>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 pt-6" in:fly={{ y: 20, duration: 300, delay: 700 }}>
				<button
					type="button"
					class="px-6 py-2.5 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700
                        rounded-lg transition-all duration-300 transform hover:scale-105"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-6 py-2.5 bg-gradient-to-r from-amber-500 to-amber-600 text-white rounded-lg
                        hover:from-amber-600 hover:to-amber-700 transition-all duration-300 transform hover:scale-105
                        disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
					disabled={isLoading}
				>
					{#if isLoading}
						<svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
								fill="none"
							/>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							/>
						</svg>
					{/if}
					{isLoading ? 'Saving...' : incident ? 'Save Changes' : 'Report Incident'}
				</button>
			</div>
		</form>
	</div>

	<Toast position="tr" />
</Modal>

<style lang="postcss">
	/* Стилизация скроллбара */
	::-webkit-scrollbar {
		width: 8px;
	}

	::-webkit-scrollbar-track {
		background: theme(colors.gray.100);
		border-radius: 4px;
	}

	::-webkit-scrollbar-thumb {
		background: theme(colors.amber.500);
		border-radius: 4px;
	}

	::-webkit-scrollbar-thumb:hover {
		background: theme(colors.amber.600);
	}

	/* Для Firefox */
	* {
		scrollbar-width: thin;
		scrollbar-color: theme(colors.amber.500) theme(colors.gray.100);
	}

	/* Add styles for the suggestions dropdown */
	.hive-search {
		position: relative;
	}
</style>
