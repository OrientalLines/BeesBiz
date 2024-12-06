<script lang="ts">
	import Modal from './Modal.svelte';
	import type { ProductionReport, User } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import {
		FileSpreadsheet,
		TrendingUp,
		DollarSign,
		Droplets,
		Users as UsersIcon
	} from 'lucide-svelte';
	import { getApiaries, getUsers } from '$lib/services/api';
	import type { Apiary } from '$lib/types';
	import { Home } from 'lucide-svelte';
	import { onMount } from 'svelte';

	export let isOpen = false;
	export let report: ProductionReport | null = null;
	export let onClose = () => {};
	export let onSave = (updatedReport: ProductionReport) => {};

	const toastStore = getToastStore();
	let isLoading = false;
	let formData: Partial<ProductionReport & { curated_by?: number }> = {
		report_id: undefined,
		apiary_id: 0,
		curated_by: undefined,
		start_date: new Date(),
		end_date: new Date(),
		total_honey_produced: 0,
		total_expenses: 0
	};

	let apiaries: Apiary[] = [];
	let filteredApiaries: Apiary[] = [];
	let apiarySearchQuery = '';
	let showApiaryDropdown = false;
	let selectedApiary: Apiary | null = null;

	// User selection state
	let users: User[] = [];
	let selectedCurator: User | null = null;

	// Load apiaries and users on mount
	onMount(async () => {
		await Promise.all([loadApiaries(), loadUsers()]);

		// Populate form if editing an existing report
		if (report) {
			formData = {
				...report,
				start_date: report.start_date ? new Date(report.start_date) : new Date(),
				end_date: report.end_date ? new Date(report.end_date) : new Date()
			};

			// Find and set selected apiary
			const matchingApiary = apiaries.find((a) => a.apiary_id === report.apiary_id);
			if (matchingApiary) {
				selectedApiary = matchingApiary;
				apiarySearchQuery = `Apiary #${matchingApiary.apiary_id} (${matchingApiary.location})`;
			}

			// Find and set selected curator
			if (report.curated_by) {
				const matchingCurator = users.find((u) => u.user_id === report.curated_by);
				if (matchingCurator) {
					selectedCurator = matchingCurator;
				}
			}
		}
	});

	async function loadApiaries() {
		try {
			isLoading = true;
			apiaries = await getApiaries();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load apiaries',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	async function loadUsers() {
		try {
			isLoading = true;
			users = await getUsers();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load users',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	function handleDateRangeChange(start: Date | null, end: Date | null) {
		if (start) formData.start_date = start;
		if (end) formData.end_date = end;
	}

	async function handleSubmit() {
		// Validate required fields
		if (!formData.apiary_id || !formData.curated_by) {
			toastStore.trigger({
				message: '❌ Please select an apiary and a curator',
				background: 'variant-filled-error'
			});
			return;
		}

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
		formData.total_honey_produced && formData.total_expenses
			? (
					((formData.total_honey_produced * 25 - formData.total_expenses) / (formData.total_honey_produced * 25)) *
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
						on:focus={() => (showApiaryDropdown = true)}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
							bg-white dark:bg-gray-800 text-gray-900 dark:text-white
							focus:ring-2 focus:ring-amber-500 focus:border-transparent
							transition-all duration-300"
						required
					/>

					{#if showApiaryDropdown && apiaries.length > 0}
						<div
							class="absolute z-50 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-200
							dark:border-gray-700 rounded-lg shadow-lg max-h-48 overflow-y-auto dropdown-scroll"
						>
							{#each apiaries.filter((a) => a.apiary_id
										.toString()
										.includes(apiarySearchQuery) || a.location
										.toLowerCase()
										.includes(apiarySearchQuery.toLowerCase())) as apiary}
								<button
									type="button"
									class="w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700
										text-gray-900 dark:text-white"
									on:click={() => {
										selectedApiary = apiary;
										formData.apiary_id = apiary.apiary_id;
										apiarySearchQuery = `Apiary #${apiary.apiary_id} (${apiary.location})`;
										showApiaryDropdown = false;
									}}
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

			<!-- Curator Selection -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 350 }}>
				<div class="flex items-center gap-2 mb-2">
					<UsersIcon class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Report Curator
					</label>
				</div>
				<div class="border-2 border-amber-500/50 rounded-lg overflow-y-auto dark:bg-gray-800">
					<div class="p-2 space-y-2">
						{#each users.filter((u) => !selectedCurator || u.user_id === selectedCurator.user_id || u.full_name
									.toLowerCase()
									.includes(apiarySearchQuery.toLowerCase()) || u.email
									.toLowerCase()
									.includes(apiarySearchQuery.toLowerCase())) as user}
							<label
								class="flex items-center p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md cursor-pointer transition-colors duration-150"
								class:bg-amber-100={selectedCurator?.user_id === user.user_id}
							>
								<input
									type="radio"
									name="curator"
									value={user.user_id}
									checked={selectedCurator?.user_id === user.user_id}
									on:change={() => {
										selectedCurator = user;
										formData.curated_by = user.user_id;
									}}
									class="form-radio h-5 w-5 text-amber-500 rounded border-gray-300 focus:ring-amber-500 dark:border-gray-600 dark:bg-gray-700 dark:checked:bg-amber-500"
								/>
								<div class="ml-3">
									<span class="font-medium">{user.full_name}</span>
									<span class="text-sm text-gray-500 dark:text-gray-400 ml-1">({user.email})</span>
								</div>
							</label>
						{/each}
					</div>
				</div>
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
						bind:value={formData.total_honey_produced}
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
			{#if formData.total_honey_produced && formData.total_expenses && (formData.total_honey_produced > 0 || formData.total_expenses > 0)}
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
</Modal>
