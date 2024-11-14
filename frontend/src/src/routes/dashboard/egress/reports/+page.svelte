<script lang="ts">
	import { debounce } from 'lodash-es';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { Report } from '$lib/types';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import ReportEditModal from '$lib/components/modals/ReportEditModal.svelte';
	import ReportDetailsModal from '$lib/components/modals/ReportDetailsModal.svelte';

	let showAddModal = false;
	let editingReport: Report | null = null;
	let viewingReport: Report | null = null;

	// State management
	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 10;
	let dateRange = {
		start: null as Date | null,
		end: null as Date | null
	};

	// Debounced search function
	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	// Mock data - replace with API call
	let reports: Report[] = [
		{
			reportId: 1,
			apiaryId: 101,
			startDate: new Date('2024-01-01'),
			endDate: new Date('2024-03-31'),
			totalHoney: 450,
			totalExpenses: 2500
		}
	];

	// Add handler for date changes
	function handleDateRangeChange(start: Date | null, end: Date | null) {
		dateRange.start = start;
		dateRange.end = end;
	}

	function handleSaveReport(report: Report) {
		if (editingReport) {
			const index = reports.findIndex((r) => r.reportId === report.reportId);
			if (index !== -1) {
				reports[index] = report;
			}
		} else {
			report.reportId = Math.max(...reports.map((r) => r.reportId)) + 1;
			reports = [...reports, report];
		}
		editingReport = null;
		showAddModal = false;
	}

	// Filter logic
	$: filteredReports = reports
		.filter((r) => {
			const matchesSearch = r.apiaryId.toString().includes(searchQuery);
			const matchesDate =
				(!dateRange.start || new Date(r.startDate) >= dateRange.start) &&
				(!dateRange.end || new Date(r.endDate) <= dateRange.end);
			return matchesSearch && matchesDate;
		})
		.sort((a, b) => new Date(b.endDate).getTime() - new Date(a.endDate).getTime());

	$: paginatedReports = filteredReports.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredReports.length / itemsPerPage);

	// Calculate totals for summary cards
	$: totalHoneyProduction = filteredReports.reduce((sum, r) => sum + r.totalHoney, 0);
	$: totalExpenses = filteredReports.reduce((sum, r) => sum + r.totalExpenses, 0);
	$: averageProduction = totalHoneyProduction / filteredReports.length || 0;
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

	<!-- Reports Table -->
	<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
		<div class="overflow-x-auto">
			<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
				<thead class="bg-gray-50 dark:bg-gray-700">
					<tr>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Report Period
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Apiary ID
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Honey Production
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Expenses
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each paginatedReports as report}
						<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{new Date(report.startDate).toLocaleDateString()} -
								{new Date(report.endDate).toLocaleDateString()}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{report.apiaryId}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{report.totalHoney} kg
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								${report.totalExpenses}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
								<button
									class="text-amber-600 hover:text-amber-700"
									on:click={() => (viewingReport = report)}
								>
									View Details
								</button>
								<button class="text-amber-600 hover:text-amber-700"> Download PDF </button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>

	<!-- Pagination -->
	{#if totalPages > 1}
		<div class="mt-4 flex justify-between items-center">
			<div class="text-sm text-gray-700 dark:text-gray-300">
				Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(
					currentPage * itemsPerPage,
					filteredReports.length
				)} of {filteredReports.length} entries
			</div>
			<div class="flex gap-2">
				<button
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                        disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={currentPage === 1}
					on:click={() => currentPage--}
				>
					Previous
				</button>
				<button
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                        disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={currentPage === totalPages}
					on:click={() => currentPage++}
				>
					Next
				</button>
			</div>
		</div>
	{/if}
	{#if showAddModal}
		<ReportEditModal
			isOpen={true}
			report={null}
			onClose={() => (showAddModal = false)}
			onSave={handleSaveReport}
		/>
	{/if}

	{#if editingReport}
		<ReportEditModal
			isOpen={true}
			report={editingReport}
			onClose={() => (editingReport = null)}
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
</div>
