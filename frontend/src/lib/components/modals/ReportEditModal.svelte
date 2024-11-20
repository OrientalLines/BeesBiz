<script lang="ts">
	import Modal from './Modal.svelte';
	import type { Report } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import { FileSpreadsheet, TrendingUp, DollarSign, Droplets } from 'lucide-svelte';

	export let isOpen = false;
	export let report: Report | null = null;
	export let onClose = () => {};
	export let onSave = (updatedReport: Report) => {};

	const toastStore = getToastStore();
	let isLoading = false;
	let formData: Partial<Report> = {
		apiaryId: 0,
		startDate: new Date(),
		endDate: new Date(),
		totalHoney: 0,
		totalExpenses: 0
	};

	$: if (report) {
		formData = { ...report };
	}

	function handleDateRangeChange(start: Date | null, end: Date | null) {
		if (start) formData.startDate = start;
		if (end) formData.endDate = end;
	}

	async function handleSubmit() {
		if (!formData) return;

		isLoading = true;
		try {
			await onSave(formData as Report);
			toastStore.trigger({
				message: '✨ Report updated successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to update report',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	// Calculate report metrics
	$: profitMargin =
		formData.totalHoney && formData.totalExpenses
			? (
					((formData.totalHoney * 25 - formData.totalExpenses) / (formData.totalHoney * 25)) *
					100
				).toFixed(1)
			: 0;
</script>

<Modal {isOpen} title={report ? 'Edit Report' : 'Generate Report'} {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<!-- Report Header -->
			<div class="flex items-center gap-3 mb-6" in:fly={{ y: 20, duration: 300, delay: 100 }}>
				<FileSpreadsheet class="w-6 h-6 text-amber-500" />
				<h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100">Report Details</h2>
			</div>

			<!-- Date Range Picker -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 200 }}>
				<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
					Report Period
				</label>
				<DatePicker
					startDate={formData.startDate}
					endDate={formData.endDate}
					onChange={handleDateRangeChange}
				/>
			</div>

			<!-- Apiary Selection -->
			<div class="relative" in:fly={{ y: 20, duration: 300, delay: 300 }}>
				<input
					type="number"
					bind:value={formData.apiaryId}
					placeholder=" "
					class="peer w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                        bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                        focus:ring-2 focus:ring-amber-500 focus:border-transparent
                        transition-all duration-300"
					required
				/>
				<label
					class="absolute left-2 -top-2.5 bg-white dark:bg-gray-800 px-2 text-sm text-amber-500
                        transition-all duration-300"
				>
					Apiary ID
				</label>
			</div>

			<!-- Metrics Grid -->
			<div
				class="grid grid-cols-1 md:grid-cols-2 gap-4"
				in:fly={{ y: 20, duration: 300, delay: 400 }}
			>
				<!-- Honey Production -->
				<div class="relative space-y-2">
					<div class="flex items-center gap-2 mb-2">
						<Droplets class="w-4 h-4 text-amber-500" />
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
							Honey Production
						</label>
					</div>
					<input
						type="number"
						step="0.1"
						bind:value={formData.totalHoney}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
						required
					/>
					<span class="absolute right-3 top-[2.8rem] text-gray-500">kg</span>
				</div>

				<!-- Expenses -->
				<div class="relative space-y-2">
					<div class="flex items-center gap-2 mb-2">
						<DollarSign class="w-4 h-4 text-amber-500" />
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
							Total Expenses
						</label>
					</div>
					<input
						type="number"
						step="0.01"
						bind:value={formData.totalExpenses}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
						required
					/>
					<span class="absolute right-3 top-[2.8rem] text-gray-500">$</span>
				</div>
			</div>

			<!-- Profit Margin Indicator -->
			{#if formData.totalHoney && formData.totalExpenses && (formData.totalHoney > 0 || formData.totalExpenses > 0)}
				<div
					class="bg-gray-50 dark:bg-gray-800/50 p-4 rounded-lg space-y-2"
					in:fly={{ y: 20, duration: 300, delay: 500 }}
				>
					<div class="flex items-center gap-2">
						<TrendingUp class="w-4 h-4 text-amber-500" />
						<span class="text-sm font-medium text-gray-700 dark:text-gray-300">
							Estimated Profit Margin
						</span>
					</div>
					<div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
						<div
							class="bg-gradient-to-r from-amber-500 to-amber-600 h-2.5 rounded-full transition-all duration-500"
							style="width: {Math.max(0, Math.min(100, Number(profitMargin)))}%"
						/>
					</div>
					<span class="text-sm text-gray-600 dark:text-gray-400">
						{profitMargin}% profit margin
					</span>
				</div>
			{/if}

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 pt-6" in:fly={{ y: 20, duration: 300, delay: 600 }}>
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
					{isLoading ? 'Saving...' : report ? 'Save Changes' : 'Generate Report'}
				</button>
			</div>
		</form>
	</div>

	<Toast position="tr" />
</Modal>
