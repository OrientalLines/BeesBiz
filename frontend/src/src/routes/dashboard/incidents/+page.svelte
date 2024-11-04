<script lang="ts">
    import type { Incident } from '$lib/types';

    let incidents: Incident[] = [
        {
            id: '1',
            hiveId: '1',
            date: new Date('2024-03-12'),
            type: 'Disease',
            description: 'Varroa mites detected',
            status: 'open'
        }
        // Add more mock data as needed
    ];

    function addIncident() {
        // Implementation
    }
</script>

<div class="space-y-6">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Incident Reports</h1>
        <button
            class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
            on:click={addIncident}
        >
            Report Incident
        </button>
    </div>

    <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-700">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Date</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Hive ID</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Type</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Description</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
                </tr>
            </thead>
            <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                {#each incidents as incident}
                    <tr>
                        <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
                            {new Date(incident.date).toLocaleDateString()}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">{incident.hiveId}</td>
                        <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">{incident.type}</td>
                        <td class="px-6 py-4 text-gray-900 dark:text-gray-300">{incident.description}</td>
                        <td class="px-6 py-4 whitespace-nowrap">
                            <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                                {incident.status === 'open' ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300' : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'}">
                                {incident.status}
                            </span>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</div>
