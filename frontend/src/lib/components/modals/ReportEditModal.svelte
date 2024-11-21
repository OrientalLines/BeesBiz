<script lang="ts">
	import Modal from './Modal.svelte';
	import type { ProductionReport } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import { FileSpreadsheet, TrendingUp, DollarSign, Droplets } from 'lucide-svelte';
	import { getApiaries } from '$lib/services/api';
	import type { Apiary } from '$lib/types';
	import { Home } from 'lucide-svelte';
	import { debounce } from 'lodash-es';
	import { onMount } from 'svelte';

	export let isOpen = false;
	export let report: ProductionReport | null = null;
	export let onClose = () => {};
	export let onSave = (updatedReport: ProductionReport) => {};

	const toastStore = getToastStore();
	let isLoading = false;
	let formData: Partial<ProductionReport> = {
		report_id: undefined,
		apiary_id: 0,
		start_date: new Date(),
		end_date: new Date(),
		total_honey: 0,
		total_expenses: 0
	};

	let apiaries: Apiary[] = [];
	let filteredApiaries: Apiary[] = [];
	let apiarySearchQuery = '';
	let showApiaryDropdown = false;
	let selectedApiary: Apiary | null = null;

	// Load apiaries on mount
	onMount(async () => {
		await loadApiaries();
		if (report) {
			formData = {
				...report,
				start_date: report.start_date ? new Date(report.start_date) : new Date(),
				end_date: report.end_date ? new Date(report.end_date) : new Date()
			};
			const matchingApiary = apiaries.find((a) => a.apiary_id === report.apiary_id);
			if (matchingApiary) {
				selectedApiary = matchingApiary;
				apiarySearchQuery = `Apiary #${matchingApiary.apiary_id} (${matchingApiary.location})`;
			}
		}
	});

	async function loadApiaries() {
		try {
			isLoading = true;
			apiaries = await getApiaries();
			filterApiaries();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load apiaries',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	const filterApiaries = debounce(() => {
		if (!apiarySearchQuery.trim()) {
			filteredApiaries = apiaries;
		} else {
			filteredApiaries = apiaries.filter(
				(apiary) =>
					apiary.apiary_id.toString().includes(apiarySearchQuery) ||
					(apiary.location &&
						apiary.location.toLowerCase().includes(apiarySearchQuery.toLowerCase()))
			);
		}
	}, 300);

	function handleApiaryInput(event: Event) {
		const input = event.target as HTMLInputElement;
		apiarySearchQuery = input.value;
		showApiaryDropdown = true;
		selectedApiary = null;
		formData.apiary_id = 0;
	}

	function selectApiary(apiary: Apiary) {
		selectedApiary = apiary;
		formData.apiary_id = apiary.apiary_id;
		apiarySearchQuery = `Apiary #${apiary.apiary_id} (${apiary.location})`;
		showApiaryDropdown = false;
	}

	// Add reactive statement for apiary search
	$: if (apiarySearchQuery !== undefined) {
		filterApiaries();
	}

	function handleDateRangeChange(start: Date | null, end: Date | null) {
		if (start) formData.start_date = start;
		if (end) formData.end_date = end;
	}

	async function handleSubmit() {
		if (!formData.apiary_id) return;

		isLoading = true;
		try {
			const submitData = {
				...formData,
				report_id: report?.report_id // Keep reportId undefined for new reports
			} as ProductionReport;

			await onSave(submitData);
			toastStore.trigger({
				message: '✨ Report saved successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to save report',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	// Calculate report metrics
	$: profitMargin =
		formData.total_honey && formData.total_expenses
			? (
					((formData.total_honey * 25 - formData.total_expenses) / (formData.total_honey * 25)) *
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
					startDate={formData.start_date ? new Date(formData.start_date) : undefined}
					endDate={formData.end_date ? new Date(formData.end_date) : undefined}
					onChange={handleDateRangeChange}
				/>
			</div>

			<!-- Apiary Selection -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 300 }}>
				<div class="flex items-center gap-2 mb-2">
					<Home class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Apiary Selection
					</label>
				</div>
				<div class="relative">
					<input
						type="text"
						bind:value={apiarySearchQuery}
						placeholder="Search by apiary ID or location..."
						on:input={handleApiaryInput}
						on:focus={() => (showApiaryDropdown = true)}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
							bg-white dark:bg-gray-800 text-gray-900 dark:text-white
							focus:ring-2 focus:ring-amber-500 focus:border-transparent
							transition-all duration-300"
						required
					/>

					{#if showApiaryDropdown && filteredApiaries.length > 0}
						<div
							class="absolute z-50 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-200
							dark:border-gray-700 rounded-lg shadow-lg max-h-48 overflow-y-auto dropdown-scroll"
						>
							{#each filteredApiaries as apiary}
								<button
									type="button"
									class="w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700
										text-gray-900 dark:text-white"
									on:click={() => selectApiary(apiary)}
								>
									Apiary #{apiary.apiary_id} ({apiary.location})
								</button>
							{/each}
						</div>
					{/if}
				</div>
				{#if selectedApiary}
					<p class="text-sm text-gray-600 dark:text-gray-400">
						Selected: Apiary #{selectedApiary.apiary_id} ({selectedApiary.location})
					</p>
				{/if}
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
						bind:value={formData.total_honey}
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
						bind:value={formData.total_expenses}
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
			{#if formData.total_honey && formData.total_expenses && (formData.total_honey > 0 || formData.total_expenses > 0)}
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
