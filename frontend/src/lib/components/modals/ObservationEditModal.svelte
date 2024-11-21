<script lang="ts">
	import Modal from './Modal.svelte';
	import type { ObservationLog, Hive } from '$lib/types';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import { ClipboardList, MessageSquare } from 'lucide-svelte';
	import { fly } from 'svelte/transition';
	import { getHives } from '$lib/services/api';
	import { onMount } from 'svelte';
	import { getToastStore } from '@skeletonlabs/skeleton';

	export let isOpen = false;
	export let observation: ObservationLog | null = null;
	export let onClose = () => {};
	export let onSave = (updatedObservation: ObservationLog) => {};

	const toastStore = getToastStore();
	let hives: Hive[] = [];
	let loading = false;
	let searchQuery = '';

	let formData: Partial<ObservationLog> = {
		observation_date: new Date(),
		hive_id: 0,
		description: '',
		recommendations: ''
	};

	onMount(async () => {
		await loadHives();
	});

	async function loadHives() {
		try {
			loading = true;
			hives = await getHives();
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to load hives',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	$: if (observation) {
		formData = {
			...observation,
			description: observation.description || '',
			recommendations: observation.recommendations || '',
			observation_date: observation.observation_date
				? new Date(observation.observation_date)
				: new Date()
		};
	}

	$: filteredHives = hives.filter(
		(hive) =>
			hive.hive_id.toString().includes(searchQuery) ||
			hive.hive_type.toLowerCase().includes(searchQuery.toLowerCase())
	);

	function handleSubmit() {
		if (formData) {
			const submissionData = {
				...formData,
				log_id: observation?.log_id,
				observation_date:
					formData.observation_date instanceof Date
						? formData.observation_date.toISOString()
						: formData.observation_date
							? new Date(formData.observation_date).toISOString()
							: new Date().toISOString(),
				description: formData.description || '',
				recommendations: formData.recommendations || ''
			};
			onSave(submissionData as ObservationLog);
			onClose();
		}
	}

	function handleDateChange(date: Date | null) {
		if (date) {
			formData.observation_date = date;
		}
	}

	function handleHiveInput(event: Event) {
		const input = event.target as HTMLInputElement;
		const value = input.value;
		searchQuery = value;
		formData.hive_id = parseInt(value) || 0;
	}
</script>

<Modal {isOpen} title={observation ? 'Edit Observation' : 'New Observation'} {onClose}>
	<form class="space-y-4" on:submit|preventDefault={handleSubmit}>
		<div>
			<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Observation Date
			</label>

			<DatePicker
				startDate={new Date(formData.observation_date ?? '')}
				onChange={handleDateChange}
			/>
		</div>

		<div class="space-y-2">
			<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"> Hive </label>
			<div class="relative">
				<input
					type="text"
					value={searchQuery}
					on:input={handleHiveInput}
					class="w-full px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
						bg-white dark:bg-gray-800 text-gray-900 dark:text-white
						focus:ring-2 focus:ring-amber-500 focus:border-transparent"
					placeholder="Search by hive ID or type..."
					required
				/>
				{#if searchQuery && filteredHives.length > 0}
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
									searchQuery = `Hive #${hive.hive_id} (${hive.hive_type})`;
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
