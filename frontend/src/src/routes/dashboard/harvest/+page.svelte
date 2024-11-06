<script lang="ts">
	import type { HoneyHarvest } from '$lib/types';

	let harvests: HoneyHarvest[] = [
		{
			id: '1',
			hiveId: '1',
			date: new Date('2024-03-10'),
			amount: 25.5,
			quality: 'Premium'
		}
		// Add more mock data
	];

	function addHarvest() {
		// Implementation
	}
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Honey Harvest Records</h1>
		<button
			class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
			on:click={addHarvest}
		>
			Record New Harvest
		</button>
	</div>

	<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg overflow-hidden">
		<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Date</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Hive ID</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Amount (kg)</th
					>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Quality</th
					>
				</tr>
			</thead>
			<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
				{#each harvests as harvest}
					<tr>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
							{new Date(harvest.date).toLocaleDateString()}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300"
							>{harvest.hiveId}</td
						>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300"
							>{harvest.amount}</td
						>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300"
							>{harvest.quality}</td
						>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
