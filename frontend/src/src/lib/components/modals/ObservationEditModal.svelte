<script lang="ts">
	import Modal from './Modal.svelte';
	import type { ObservationLog } from '$lib/types';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import { ClipboardList, MessageSquare } from 'lucide-svelte';
	import { fly } from 'svelte/transition';

	export let isOpen = false;
	export let observation: ObservationLog | null = null;
	export let onClose = () => {};
	export let onSave = (updatedObservation: ObservationLog) => {};

	let formData: Partial<ObservationLog> = {
		observationDate: new Date(),
		hiveId: 0,
		description: '',
		recommendations: ''
	};

	$: if (observation) {
		formData = { ...observation };
	}

	function handleSubmit() {
		if (formData) {
			onSave(formData as ObservationLog);
			onClose();
		}
	}

	function handleDateChange(date: Date | null) {
		if (date) {
			formData.observationDate = date;
		}
	}
</script>

<Modal {isOpen} title={observation ? 'Edit Observation' : 'New Observation'} {onClose}>
	<form class="space-y-4" on:submit|preventDefault={handleSubmit}>
		<div>
			<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Observation Date
			</label>
			<DatePicker startDate={formData.observationDate} onChange={handleDateChange} />
		</div>

		<div>
			<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Hive ID
			</label>
			<input
				type="number"
				bind:value={formData.hiveId}
				class="w-full px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                    bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                    focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				required
			/>
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
			<button type="submit" class="px-4 py-2 bg-amber-500 text-white rounded-lg hover:bg-amber-600">
				{observation ? 'Save Changes' : 'Create Observation'}
			</button>
		</div>
	</form>
</Modal>
