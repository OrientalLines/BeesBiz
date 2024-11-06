<script lang="ts">
import type { Region } from '$lib/types';
import Icon from '@iconify/svelte';
import { fade } from 'svelte/transition';

export let onSelect: (region: Region) => void;

// Climate zone icons mapping
const climateIcons = {
    'tropical': 'mdi:weather-sunny',
    'temperate': 'mdi:weather-partly-cloudy',
    'continental': 'mdi:pine-tree',
    'mediterranean': 'mdi:palm-tree',
} as const;

const regions: Region[] = [
    {
        region_id: 'reg1',
        name: 'Eastern Valley',
        climate_zone: 'temperate'
    },
    {
        region_id: 'reg2',
        name: 'Mountain Ridge',
        climate_zone: 'continental'
    },
    {
        region_id: 'reg3',
        name: 'Coastal Plains',
        climate_zone: 'mediterranean'
    },
    {
        region_id: 'reg4',
        name: 'Central Forest',
        climate_zone: 'temperate'
    }
];
</script>

<div class="max-w mx-auto" in:fade>
    <div class="flex justify-between items-center mb-8">
        <div>
            <h1 class="text-4xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-3">
                <Icon icon="mdi:map-marker-radius" class="w-8 h-8 text-amber-500" />
                Regions
            </h1>
            <p class="mt-2 text-gray-600 dark:text-gray-400">
                Select a region to view its apiaries
            </p>
        </div>
        <button class="bg-amber-500 text-white px-4 py-2 rounded-lg
            hover:bg-amber-600 transition-all flex items-center gap-2 text-sm">
            <Icon icon="mdi:plus" class="w-4 h-4" />
            New Region
        </button>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {#each regions as region}
        <button
            class="group flex flex-col p-5 bg-white dark:bg-gray-800 
            rounded-xl transition-all duration-300 hover:ring-2 hover:ring-amber-500
            border border-gray-100 dark:border-gray-700"
            on:click={() => onSelect(region)}
        >
            <div class="flex items-start justify-between mb-3">
                <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-lg bg-amber-100 dark:bg-amber-900/30 
                        flex items-center justify-center">
                        <Icon 
                            icon={climateIcons[region.climate_zone]} 
                            class="w-6 h-6 text-amber-600 dark:text-amber-400" 
                        />
                    </div>
                    <div class="text-left">
                        <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                            {region.name}
                        </h2>
                        <p class="text-sm text-gray-500 dark:text-gray-400 capitalize">
                            {region.climate_zone} Zone
                        </p>
                    </div>
                </div>
                <Icon 
                    icon="mdi:chevron-right" 
                    class="w-5 h-5 text-gray-400 group-hover:text-amber-500 
                    group-hover:translate-x-1 transition-all" 
                />
            </div>
        </button>
        {/each}
    </div>
</div>
