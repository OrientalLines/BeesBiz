<script lang="ts">
	import Modal from './Modal.svelte';
	import type { ObservationLog, Hive } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import { getHives } from '$lib/services/api';
	import { onMount } from 'svelte';
	import { Home, ClipboardList, MessageSquare, Icon } from 'lucide-svelte';

	export let isOpen = false;
	export let observation: ObservationLog | null = null;
	export let onClose = () => {};
	export let onSave = (observation: ObservationLog) => {};

	const toastStore = getToastStore();
	let hives: Hive[] = [];
	let filteredHives: Hive[] = [];
	let showHiveSuggestions = false;
	let hiveSearchTerm = '';
	let selectedDate: Date = new Date();

	let formData: Partial<ObservationLog> = {
		hive_id: 0,
		description: '',
		recommendations: '',
		observation_date: new Date().toISOString()
	};

	onMount(async () => {
		try {
			hives = await getHives();
			if (observation) {
				formData = {
					...observation
				};
				selectedDate = new Date(observation.observation_date);
				const selectedHive = hives.find((h) => h.hive_id === observation.hive_id);
				if (selectedHive) {
					hiveSearchTerm = `#${selectedHive.hive_id} (${selectedHive.hive_type})`;
				}
			}
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load hives',
				background: 'variant-filled-error'
			});
		}
	});

	$: filteredHives = hives.filter(
		(hive) =>
			hive.hive_id.toString().includes(hiveSearchTerm) ||
			hive.hive_type.toLowerCase().includes(hiveSearchTerm.toLowerCase())
	);

	function handleSubmit() {
		if (formData) {
			onSave({
				...formData,
				log_id: observation?.log_id || 0,
				hive_id: formData.hive_id || 0,
				observation_date: selectedDate.toISOString(),
				description: formData.description || '',
				recommendations: formData.recommendations || ''
			} as ObservationLog);
		}
	}

	function handleDateChange(date: Date | null) {
		console.log('Date changed:', date);
		if (date) {
			selectedDate = date;
			formData.observation_date = date.toISOString();
		}
	}

	function handleHiveInput(event: Event) {
		const input = event.target as HTMLInputElement;
		const value = input.value;
		hiveSearchTerm = value;
		formData.hive_id = parseInt(value) || 0;
	}
</script>

<Modal {isOpen} title={observation ? 'Edit Observation' : 'New Observation'} {onClose}>
	<form class="space-y-4" on:submit|preventDefault={handleSubmit}>
		<div class="space-y-2">
			<div class="flex items-center gap-2">
				<Icon icon="mdi:calendar" class="w-4 h-4 text-amber-500" />
				<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
					Observation Date
				</label>
			</div>
			<DatePicker startDate={selectedDate} onChange={handleDateChange} singleDateMode={true} />
			<div class="text-xs text-gray-500">
				Selected date: {selectedDate.toLocaleString()}
			</div>
		</div>

		<div class="space-y-2">
			<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"> Hive </label>
			<div class="relative">
				<input
					type="text"
					value={hiveSearchTerm}
					on:input={handleHiveInput}
					class="w-full px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
						bg-white dark:bg-gray-800 text-gray-900 dark:text-white
						focus:ring-2 focus:ring-amber-500 focus:border-transparent"
					placeholder="Search by hive ID or type..."
					required
				/>
				{#if hiveSearchTerm && filteredHives.length > 0}
					<div
						class="absolute z-50 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-200
							dark:border-gray-700 rounded-lg shadow-lg max-h-48 overflow-y-auto"
					>
						{#each filteredHives as hive}
							<button
								type="button"
								class="w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700
									text-gray-900 dark:text-white"
								on:click={() => {
									formData.hive_id = hive.hive_id;
									hiveSearchTerm = `Hive #${hive.hive_id} (${hive.hive_type})`;
								}}
							>
								Hive #{hive.hive_id} ({hive.hive_type}) - {hive.current_status}
							</button>
						{/each}
					</div>
				{/if}
			</div>
			{#if formData?.hive_id && formData.hive_id > 0}
				<p class="text-sm text-gray-600 dark:text-gray-400">
					Selected: Hive #{formData.hive_id}
					{#if hives.find((h) => h.hive_id === formData.hive_id)}
						({hives.find((h) => h.hive_id === formData.hive_id)?.hive_type})
					{/if}
				</p>
			{/if}
		</div>

		<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 500 }}>
			<div class="flex items-center gap-2">
				<ClipboardList class="w-4 h-4 text-amber-500" />
				<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
					Observation Details
				</label>
			</div>
			<textarea
				bind:value={formData.description}
				rows="3"
				class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
					bg-white dark:bg-gray-800 text-gray-900 dark:text-white
					focus:ring-2 focus:ring-amber-500 focus:border-transparent
					transition-all duration-300 resize-none"
				required
				placeholder="Describe what you observed..."
			/>
		</div>

		<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 600 }}>
			<div class="flex items-center gap-2">
				<MessageSquare class="w-4 h-4 text-amber-500" />
				<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
					Recommendations
				</label>
			</div>
			<textarea
				bind:value={formData.recommendations}
				rows="2"
				class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
					bg-white dark:bg-gray-800 text-gray-900 dark:text-white
					focus:ring-2 focus:ring-amber-500 focus:border-transparent
					transition-all duration-300 resize-none"
				placeholder="Enter any recommendations for future actions..."
			/>
		</div>

		<div class="flex justify-end gap-3 mt-6">
			<button
				type="button"
				class="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg"
				on:click={onClose}
			>
				Cancel
			</button>
			<button
				type="submit"
				class="px-4 py-2 bg-amber-500 text-white rounded-lg hover:bg-amber-600"
				disabled={!formData.hive_id}
			>
				{observation ? 'Save Changes' : 'Create Observation'}
			</button>
		</div>
	</form>
</Modal>
