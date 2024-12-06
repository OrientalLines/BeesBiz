<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { ProductionReport } from '$lib/types';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import ReportEditModal from '$lib/components/modals/ReportEditModal.svelte';
	import ReportDetailsModal from '$lib/components/modals/ReportDetailsModal.svelte';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import { ReportPDFGenerator } from '$lib/pdf/reports';
	import {
		getProductionReports,
		createProductionReport,
		updateProductionReport,
		deleteProductionReport,
		getCuratedProductionReports
	} from '$lib/services/api';
	import { auth } from '$lib/stores/auth';
	import { onMount } from 'svelte';

	let showAddModal = false;
	let editingReport: ProductionReport | null = null;
	let viewingReport: ProductionReport | null = null;

	// State management
	let searchQuery = '';
	let currentAllPage = 1;
	let currentCuratedPage = 1;

	const itemsPerPage = 10;
	const toastStore = getToastStore();

	let dateRange = {
		start: null as Date | null,
		end: null as Date | null
	};

	// Debounced search function
	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentAllPage = 1;
		currentCuratedPage = 1;
	}, 300);

	// New state for reports
	let allReports: ProductionReport[] = [];
	let curatedReports: ProductionReport[] = [];

	let loading = {
		all: true,
		curated: true
	};
	let error = {
		all: null as string | null,
		curated: null as string | null
	};

	async function loadAllReports() {
		try {
			loading.all = true;
			error.all = null;
			allReports = await getProductionReports();
		} catch (e) {
			console.error('Error loading all reports:', e);
			error.all = 'Failed to load all reports';
			allReports = [];
		} finally {
			loading.all = false;
		}
	}

	async function loadCuratedReports() {
		try {
			loading.curated = true;
			error.curated = null;

			// If user is logged in, get their curated reports
			if ($auth?.user?.user_id) {
				const userId = Number($auth.user.user_id);
				if (isNaN(userId)) {
					throw new Error('Invalid user ID');
				}
				curatedReports = await getCuratedProductionReports(userId);
			} else {
				curatedReports = [];
			}
		} catch (e) {
			console.error('Error loading curated reports:', e);
			error.curated = 'Failed to load curated reports';
			curatedReports = [];
		} finally {
			loading.curated = false;
		}
	}

	onMount(() => {
		loadAllReports();
		loadCuratedReports();
	});

	// Modify existing filter logic to work with both report lists
	function createReportFilter(reports: ProductionReport[]) {
		return reports
			.filter((r) => {
				const matchesSearch = r.apiary_id.toString().includes(searchQuery);
				const matchesDate =
					(!dateRange.start || new Date(r.start_date) >= dateRange.start) &&
					(!dateRange.end || new Date(r.end_date) <= dateRange.end);
				return matchesSearch && matchesDate;
			})
			.sort((a, b) => new Date(b.end_date).getTime() - new Date(a.end_date).getTime());
	}

	$: filteredAllReports = createReportFilter(allReports);
	$: filteredCuratedReports = createReportFilter(curatedReports);

	$: paginatedAllReports = filteredAllReports.slice(
		(currentAllPage - 1) * itemsPerPage,
		currentAllPage * itemsPerPage
	);

	$: paginatedCuratedReports = filteredCuratedReports.slice(
		(currentCuratedPage - 1) * itemsPerPage,
		currentCuratedPage * itemsPerPage
	);

	$: totalAllPages = Math.ceil(filteredAllReports.length / itemsPerPage);
	$: totalCuratedPages = Math.ceil(filteredCuratedReports.length / itemsPerPage);

	// Add handler for date changes
	function handleDateRangeChange(start: Date | null, end?: Date | null) {
		dateRange = {
			start,
			end: end ?? null
		};
		currentAllPage = 1;
		currentCuratedPage = 1;
	}

	async function handleSaveReport(report: ProductionReport) {
		try {
			if (report.report_id) {
				// Update existing report
				const updatedReport = await updateProductionReport(report);
				allReports = allReports.map((r) =>
					r.report_id === updatedReport.report_id ? updatedReport : r
				);
				curatedReports = curatedReports.map((r) =>
					r.report_id === updatedReport.report_id ? updatedReport : r
				);
			} else {
				// Create new report
				const newReport = await createProductionReport(report);
				allReports = [...allReports, newReport];
				curatedReports = [...curatedReports, newReport];
			}

			showAddModal = false;
			editingReport = null;
			toastStore.trigger({
				message: '✨ Report saved successfully!',
				background: 'variant-filled-success'
			});
		} catch (e) {
			console.error('Save error:', e);
			toastStore.trigger({
				message: '❌ Failed to save report',
				background: 'variant-filled-error'
			});
		}
	}

	async function handleDeleteReport(id: number) {
		try {
			await deleteProductionReport(id);
			allReports = allReports.filter((r) => r.report_id !== id);
			curatedReports = curatedReports.filter((r) => r.report_id !== id);
			toastStore.trigger({
				message: '🗑️ Report deleted successfully',
				background: 'variant-filled-success'
			});
		} catch (e) {
			console.error('Delete error:', e);
			toastStore.trigger({
				message: '❌ Failed to delete report',
				background: 'variant-filled-error'
			});
		}
	}

	async function generatePDF(report: ProductionReport) {
		const pdfGenerator = new ReportPDFGenerator();
		pdfGenerator.generate(report);
	}

	// Calculate totals for summary cards
	$: totalHoneyProduction = filteredAllReports.reduce((sum, r) => sum + r.total_honey_produced, 0);
	$: totalExpenses = filteredAllReports.reduce((sum, r) => sum + r.total_expenses, 0);
	$: averageProduction = totalHoneyProduction / (filteredAllReports.length || 1);
</script>

<div class="max-w mx-auto" in:fade>
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:file-chart" class="w-8 h-8 text-amber-500" />
				Reports
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				View and analyze apiary performance reports
			</p>
		</div>
		<button
			class="bg-amber-500 text-white px-6 py-3 rounded-full
        hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
        flex items-center gap-2"
			on:click={() => (showAddModal = true)}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Generate Report
		</button>
	</div>

	<!-- Summary Cards -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
		<div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
			<div class="flex items-center gap-3">
				<Icon icon="mdi:honey" class="w-6 h-6 text-amber-500" />
				<h3 class="text-lg font-semibold">Total Honey</h3>
			</div>
			<p class="mt-2 text-3xl font-bold">{totalHoneyProduction} kg</p>
		</div>

		<div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
			<div class="flex items-center gap-3">
				<Icon icon="mdi:currency-usd" class="w-6 h-6 text-amber-500" />
				<h3 class="text-lg font-semibold">Total Expenses</h3>
			</div>
			<p class="mt-2 text-3xl font-bold">${totalExpenses}</p>
		</div>

		<div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
			<div class="flex items-center gap-3">
				<Icon icon="mdi:chart-line" class="w-6 h-6 text-amber-500" />
				<h3 class="text-lg font-semibold">Avg Production</h3>
			</div>
			<p class="mt-2 text-3xl font-bold">{averageProduction.toFixed(1)} kg</p>
		</div>
	</div>

	<!-- Controls -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
		<div class="relative">
			<Icon
				icon="mdi:magnify"
				class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
			/>
			<input
				type="text"
				placeholder="Search by apiary ID..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
					bg-white dark:bg-gray-800 text-gray-900 dark:text-white
					placeholder-gray-500 dark:placeholder-gray-400
					focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<!-- Replace the two date inputs with single DatePicker that spans 2 columns -->
		<div class="md:col-span-2">
			<DatePicker
				startDate={dateRange.start}
				endDate={dateRange.end}
				placeholder="Filter by report period"
				onChange={handleDateRangeChange}
			/>
		</div>
	</div>

	<!-- All Reports Section -->
	<div class="mb-8">
		<h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
			All Reports
		</h2>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Report Period
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Apiary ID
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Honey Production
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Expenses
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
						{#each paginatedAllReports as report}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									{new Date(report.start_date).toLocaleDateString()} -
									{new Date(report.end_date).toLocaleDateString()}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									{report.apiary_id}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									{report.total_honey_produced} kg
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									${report.total_expenses}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
									<button
										class="text-amber-600 hover:text-amber-700"
										on:click={() => (viewingReport = report)}
									>
										View Details
									</button>
									<button
										class="text-amber-600 hover:text-amber-700"
										on:click={() => generatePDF(report)}
									>
										Download PDF
									</button>
									<button
										class="text-red-600 hover:text-red-700"
										on:click={() => handleDeleteReport(report.report_id)}
									>
										Delete
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<!-- Pagination for All Reports -->
			{#if totalAllPages > 1}
				<div class="mt-4 flex justify-between items-center p-4">
					<div class="text-sm text-gray-700 dark:text-gray-300">
						Showing {(currentAllPage - 1) * itemsPerPage + 1} to
						{Math.min(currentAllPage * itemsPerPage, filteredAllReports.length)}
						of {filteredAllReports.length} entries
					</div>
					<div class="flex gap-2">
						<button
							class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700"
							disabled={currentAllPage === 1}
							on:click={() => currentAllPage--}
						>
							Previous
						</button>
						<button
							class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700"
							disabled={currentAllPage === totalAllPages}
							on:click={() => currentAllPage++}
						>
							Next
						</button>
					</div>
				</div>
			{/if}
		</div>
	</div>

	<!-- Curated Reports Section -->
	<div>
		<h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
			Curated Reports
		</h2>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Report Period
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Apiary ID
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Honey Production
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Expenses
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
						{#each paginatedCuratedReports as report}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									{new Date(report.start_date).toLocaleDateString()} -
									{new Date(report.end_date).toLocaleDateString()}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									{report.apiary_id}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									{report.total_honey_produced} kg
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									${report.total_expenses}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
									<button
										class="text-amber-600 hover:text-amber-700"
										on:click={() => (viewingReport = report)}
									>
										View Details
									</button>
									<button
										class="text-amber-600 hover:text-amber-700"
										on:click={() => generatePDF(report)}
									>
										Download PDF
									</button>
									<button
										class="text-red-600 hover:text-red-700"
										on:click={() => handleDeleteReport(report.report_id)}
									>
										Delete
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<!-- Pagination for Curated Reports -->
			{#if totalCuratedPages > 1}
				<div class="mt-4 flex justify-between items-center p-4">
					<div class="text-sm text-gray-700 dark:text-gray-300">
						Showing {(currentCuratedPage - 1) * itemsPerPage + 1} to
						{Math.min(currentCuratedPage * itemsPerPage, filteredCuratedReports.length)}
						of {filteredCuratedReports.length} entries
					</div>
					<div class="flex gap-2">
						<button
							class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700"
							disabled={currentCuratedPage === 1}
							on:click={() => currentCuratedPage--}
						>
							Previous
						</button>
						<button
							class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700"
							disabled={currentCuratedPage === totalCuratedPages}
							on:click={() => currentCuratedPage++}
						>
							Next
						</button>
					</div>
				</div>
			{/if}
		</div>
	</div>

	<!-- Modals remain the same -->
	{#if showAddModal}
		<ReportEditModal
			isOpen={true}
			report={null}
			onClose={() => {
				showAddModal = false;
			}}
			onSave={handleSaveReport}
		/>
	{/if}

	{#if editingReport}
		<ReportEditModal
			isOpen={true}
			report={editingReport}
			onClose={() => {
				editingReport = null;
			}}
			onSave={handleSaveReport}
		/>
	{/if}

	{#if viewingReport}
		<ReportDetailsModal
			isOpen={true}
			report={viewingReport}
			onClose={() => (viewingReport = null)}
		/>
	{/if}

	{#if loading.all || loading.curated}
		<div class="flex justify-center items-center h-64">
			<div class="spinner" />
		</div>
	{:else if error.all || error.curated}
		<div class="alert variant-filled-error">
			<span>{error.all || error.curated}</span>
			<button
				class="btn"
				on:click={() => {
					loadAllReports();
					loadCuratedReports();
				}}>Retry</button
			>
		</div>
	{/if}
</div>
