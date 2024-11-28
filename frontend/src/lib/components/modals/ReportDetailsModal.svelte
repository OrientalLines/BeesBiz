<script lang="ts">
	import Modal from './Modal.svelte';
	import type { ProductionReport } from '$lib/types';
	import Icon from '@iconify/svelte';
	import { tweened } from 'svelte/motion';
	import { cubicOut } from 'svelte/easing';
	import { onMount } from 'svelte';

	export let isOpen = false;
	export let report: ProductionReport;
	export let onClose = () => {};

	// Animated number stores
	const honeyAmount = tweened(0, { duration: 1500, easing: cubicOut });
	const expenses = tweened(0, { duration: 1500, easing: cubicOut });
	const profitMarginAnim = tweened(0, { duration: 1500, easing: cubicOut });
	const revenuePerHiveAnim = tweened(0, { duration: 1500, easing: cubicOut });
	const efficiencyAnim = tweened(0, { duration: 1500, easing: cubicOut });

	// Start animations when modal opens
	$: if (isOpen) {
		honeyAmount.set(report.total_honey);
		expenses.set(report.total_expenses);
		profitMarginAnim.set(profitMargin);
		revenuePerHiveAnim.set(revenuePerHive);
		efficiencyAnim.set(productionEfficiency);
	}

	// Calculate some dummy metrics
	$: profitMargin =
		((report.total_honey * 50 - report.total_expenses) / (report.total_honey * 50)) * 100;
	$: revenuePerHive = (report.total_honey * 50) / 5; // Assuming 5 hives for demo
	$: productionEfficiency = report.total_honey / report.total_expenses;
</script>

<Modal {isOpen} title="Production Report Details" {onClose}>
	<div class="space-y-8 max-h-[80vh] overflow-y-auto pr-2">
		<!-- Period and Basic Info -->
		<div
			class="bg-gray-50 dark:bg-gray-800 p-6 rounded-xl border border-gray-100 dark:border-gray-700"
		>
			<h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
				Production Report Overview
			</h3>
			<div class="grid grid-cols-2 gap-6">
				<div>
					<p class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Period</p>
					<p class="text-gray-900 dark:text-white text-lg">
						{new Date(report.start_date).toLocaleDateString()} -
						{new Date(report.end_date).toLocaleDateString()}
					</p>
				</div>
				<div>
					<p class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Apiary ID</p>
					<p class="text-gray-900 dark:text-white text-lg">{report.apiary_id}</p>
				</div>
			</div>
		</div>

		<!-- Key Metrics -->
		<div class="grid grid-cols-2 gap-6">
			<div
				class="bg-amber-50 dark:bg-amber-900/20 p-6 rounded-xl border border-amber-100 dark:border-amber-800/30"
			>
				<div class="flex items-center gap-3">
					<div class="p-2 bg-amber-100 dark:bg-amber-800/30 rounded-lg">
						<Icon icon="mdi:honey" class="w-6 h-6 text-amber-600 dark:text-amber-400" />
					</div>
					<h4 class="font-medium text-gray-900 dark:text-white">Production</h4>
				</div>
				<p class="text-3xl font-bold text-gray-900 dark:text-white mt-3">
					{$honeyAmount.toFixed(1)} kg
					<span class="text-sm font-normal text-gray-500 dark:text-gray-400 ml-1">of honey</span>
				</p>
			</div>
			<div
				class="bg-amber-50 dark:bg-amber-900/20 p-6 rounded-xl border border-amber-100 dark:border-amber-800/30"
			>
				<div class="flex items-center gap-3">
					<div class="p-2 bg-amber-100 dark:bg-amber-800/30 rounded-lg">
						<Icon icon="mdi:currency-usd" class="w-6 h-6 text-amber-600 dark:text-amber-400" />
					</div>
					<h4 class="font-medium text-gray-900 dark:text-white">Expenses</h4>
				</div>
				<p class="text-3xl font-bold text-gray-900 dark:text-white mt-3">
					${$expenses.toFixed(2)}
					<span class="text-sm font-normal text-gray-500 dark:text-gray-400 ml-1">total</span>
				</p>
			</div>
		</div>

		<!-- Performance Metrics -->
		<div class="space-y-4">
			<h3 class="text-xl font-semibold text-gray-900 dark:text-white">Performance Metrics</h3>
			<div class="grid grid-cols-1 gap-4">
				{#each [{ label: 'Profit Margin', value: `${$profitMarginAnim.toFixed(1)}%`, icon: 'mdi:chart-line', trend: profitMargin > 15 ? 'up' : 'down', color: profitMargin > 15 ? 'text-green-500' : 'text-red-500' }, { label: 'Revenue per Hive', value: `$${$revenuePerHiveAnim.toFixed(2)}`, icon: 'mdi:beehive', trend: revenuePerHive > 1000 ? 'up' : 'down', color: revenuePerHive > 1000 ? 'text-green-500' : 'text-red-500' }, { label: 'Production Efficiency', value: `${$efficiencyAnim.toFixed(2)} kg/$`, icon: 'mdi:trending-up', trend: productionEfficiency > 2 ? 'up' : 'down', color: productionEfficiency > 2 ? 'text-green-500' : 'text-red-500' }] as metric}
					<div class="border dark:border-gray-700 p-6 rounded-xl relative overflow-hidden">
						<div class="flex items-center justify-between mb-2">
							<div class="flex items-center gap-3">
								<Icon icon={metric.icon} class="w-5 h-5 text-amber-500" />
								<p class="text-sm font-medium text-gray-500 dark:text-gray-400">{metric.label}</p>
							</div>
							<Icon
								icon={metric.trend === 'up' ? 'mdi:trending-up' : 'mdi:trending-down'}
								class="w-5 h-5 {metric.color}"
							/>
						</div>
						<p
							class="text-2xl font-semibold text-gray-900 dark:text-white flex items-baseline gap-2"
						>
							{metric.value}
							<span class="text-sm {metric.color}">
								{metric.trend === 'up' ? '+' : '-'}12%
							</span>
						</p>
						<!-- Add subtle gradient background -->
						<div
							class="absolute inset-0 bg-gradient-to-br from-transparent to-amber-50/30 dark:to-amber-900/10"
						/>
					</div>
				{/each}
			</div>
		</div>
	</div>
	<div class="flex justify-end pt-4 mt-4 border-t border-gray-200 dark:border-gray-700">
		<button
			class="px-6 py-2.5 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg
				   transition-colors duration-200 font-medium"
			on:click={onClose}
		>
			Close
		</button>
	</div>
</Modal>
