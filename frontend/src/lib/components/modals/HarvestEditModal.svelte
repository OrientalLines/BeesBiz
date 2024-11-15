<script lang="ts">
	import Modal from './Modal.svelte';
	import type { HoneyHarvest } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import { Droplets, Home, Award, CalendarDays } from 'lucide-svelte';

	export let isOpen = false;
	export let harvest: HoneyHarvest | null = null;
	export let onClose = () => {};
	export let onSave = (updatedHarvest: HoneyHarvest) => {};

	const toastStore = getToastStore();
	let isLoading = false;
	let formData: Partial<HoneyHarvest> = {
		harvestDate: new Date(),
		hiveId: 0,
		quantity: 0,
		qualityGrade: 'A'
	};

	const qualityGrades = [
		{ value: 'A', label: 'Premium Grade A', description: 'Highest quality honey' },
		{ value: 'B', label: 'Grade B', description: 'Good quality honey' },
		{ value: 'C', label: 'Grade C', description: 'Standard quality honey' },
		{ value: 'D', label: 'Grade D', description: 'Processing grade honey' }
	];

	$: if (harvest) {
		formData = { ...harvest };
	}

	async function handleSubmit() {
		if (!formData) return;

		isLoading = true;
		try {
			await onSave(formData as HoneyHarvest);
			toastStore.trigger({
				message: '✨ Harvest record saved successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to save harvest record',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	function handleDateChange(date: Date | null) {
		if (date) {
			formData.harvestDate = date;
		}
	}
</script>

<Modal {isOpen} title={harvest ? 'Edit Harvest Record' : 'New Harvest Record'} {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<!-- Harvest Date -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 100 }}>
				<div class="flex items-center gap-2 mb-2">
					<CalendarDays class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Harvest Date </label>
				</div>
				<DatePicker startDate={formData.harvestDate} onChange={handleDateChange} />
			</div>

			<!-- Hive Selection -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 200 }}>
				<div class="flex items-center gap-2 mb-2">
					<Home class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Hive ID </label>
				</div>
				<div class="relative">
					<input
						type="number"
						bind:value={formData.hiveId}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
						required
					/>
				</div>
			</div>

			<!-- Quantity -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 300 }}>
				<div class="flex items-center gap-2 mb-2">
					<Droplets class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Quantity </label>
				</div>
				<div class="relative">
					<input
						type="number"
						step="0.1"
						bind:value={formData.quantity}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
						required
					/>
					<span class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500">kg</span>
				</div>
			</div>

			<!-- Quality Grade -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 400 }}>
				<div class="flex items-center gap-2 mb-2">
					<Award class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Quality Grade
					</label>
				</div>
				<div class="grid grid-cols-2 gap-3">
					{#each qualityGrades as grade}
						<button
							type="button"
							class="flex flex-col items-center p-3 border-2 rounded-lg transition-all duration-300
                                {formData.qualityGrade === grade.value
								? 'border-amber-500 bg-amber-50 dark:bg-amber-900/10'
								: 'border-gray-200 dark:border-gray-700 hover:border-amber-500/50'}"
							on:click={() => {
								formData.qualityGrade = grade.value as 'A' | 'B' | 'C' | 'D';
							}}
						>
							<span class="text-lg font-semibold text-amber-500">
								Grade {grade.value}
							</span>
							<span class="text-xs text-gray-500 dark:text-gray-400 text-center mt-1">
								{grade.description}
							</span>
						</button>
					{/each}
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 pt-6" in:fly={{ y: 20, duration: 300, delay: 500 }}>
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
					{isLoading ? 'Saving...' : harvest ? 'Save Changes' : 'Record Harvest'}
				</button>
			</div>
		</form>
	</div>

	<Toast position="tr" />
</Modal>

<style lang="postcss">
	/* Стилизация скроллбара для современного вида */
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
</style>
