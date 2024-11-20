<script lang="ts">
	import type { Hive } from '$lib/types';

	let hives: Hive[] = [
		{
			hiveId: 1,
			apiaryId: 1,
			hiveType: 'Langstroth',
			inspectionDate: new Date('2024-03-15'),
			currentStatus: 'active'
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
		<h1 class="text-2xl font-bold text-gray-900">Hive Management</h1>
		<button
			class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
			on:click={addHive}
		>
			Add New Hive
		</button>
	</div>

	<div class="bg-white shadow-md rounded-lg overflow-hidden">
		<table class="min-w-full divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>Type</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>Apiary</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>Status</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>Last Inspection</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
						>Actions</th
					>
				</tr>
			</thead>
			<tbody class="bg-white divide-y divide-gray-200">
				{#each hives as hive}
					<tr>
						<td class="px-6 py-4 whitespace-nowrap">{hive.hiveType}</td>
						<td class="px-6 py-4 whitespace-nowrap">{hive.apiaryId}</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                                {hive.currentStatus === 'active'
									? 'bg-green-100 text-green-800'
									: hive.currentStatus === 'inactive'
										? 'bg-red-100 text-red-800'
										: 'bg-yellow-100 text-yellow-800'}"
							>
								{hive.currentStatus}
							</span>
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							{new Date(hive.inspectionDate).toLocaleDateString()}
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<button
								class="text-blue-600 hover:text-blue-900"
								on:click={() => updateHive(hive.hiveId.toString())}
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
