<script lang="ts">
import type {
    Apiary,
    Hive
} from '$lib/types';

export let selectedApiary: Apiary;
export let onBack: () => void;

// Mock data - replace with actual API calls
const hives: Hive[] = [{
    hive_id: '1',
    apiaryId: '1',
    hive_type: 'Langstroth',
    inspection_date: new Date('2024-03-15'),
    current_status: 'active'
}, ];

$: filteredHives = hives.filter(h => h.apiaryId === selectedApiary.id);

function addHive() {
    // Implementation
}

function updateHive(id: string) {
    // Implementation
}
</script>

<div class="flex justify-between items-center mb-6">
    <div class="flex items-center gap-4">
        <button
            class="text-gray-600 dark:text-gray-400 hover:text-primary"
            on:click={onBack}
            >
            ← Back to Apiaries
        </button>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Hives</h1>
    </div>
    <button
        class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary/90"
        on:click={addHive}
        >
        Add New Hive
    </button>
</div>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow-md">
    <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-700">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">ID</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">Type</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">Status</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">Last Inspection</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">Actions</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                {#each filteredHives as hive}
                <tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
                    <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
                        {hive.hive_id}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
                        {hive.hive_type}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                        <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                            {hive.current_status === 'active'
                            ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
                            : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'}">
                            {hive.current_status}
                        </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-300">
                        {new Date(hive.inspection_date).toLocaleDateString()}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap space-x-2">
                        <button
                            class="text-primary hover:text-primary/80"
                            on:click={() => updateHive(hive.hive_id)}
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
