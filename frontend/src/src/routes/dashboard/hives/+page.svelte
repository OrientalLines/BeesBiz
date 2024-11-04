<script lang="ts">
	import type { Hive } from '$lib/types';

	let hives: Hive[] = [
		{
			id: '1',
			name: 'Hive A1',
			location: 'Field 1',
			status: 'active',
			lastInspection: new Date('2024-03-15')
		}
		// Add more mock data
	];

	function addHive() {
		// Implementation
	}

	function updateHive(id: string) {
		// Implementation
	}
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Hive Management</h1>
		<button
			class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
			on:click={addHive}
		>
			Add New Hive
		</button>
	</div>

	<div class="bg-white dark:bg-gray-800 shadow-md rounded-lg overflow-hidden">
		<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Name</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Location</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Status</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Last Inspection</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
						>Actions</th
					>
				</tr>
			</thead>
			<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
				{#each hives as hive}
					<tr>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">{hive.name}</td>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">{hive.location}</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                                {hive.status === 'active'
									? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
									: hive.status === 'inactive'
										? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
										: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300'}"
							>
								{hive.status}
							</span>
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
							{new Date(hive.lastInspection).toLocaleDateString()}
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<button
								class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
								on:click={() => updateHive(hive.id)}
							>
								Edit
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
