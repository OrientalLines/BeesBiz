<script lang="ts">
import {
    fade
} from 'svelte/transition';
import RegionSelect from '$lib/components/hives/RegionSelect.svelte';
import ApiarySelect from '$lib/components/hives/ApiarySelect.svelte';
import HivesList from '$lib/components/hives/HivesList.svelte';
import type {
    Region,
    Apiary
} from '$lib/types';

let selectedRegion: Region | null = null;
let selectedApiary: Apiary | null = null;

function handleRegionSelect(region: Region) {
    selectedRegion = region;
    selectedApiary = null;
}

function handleApiarySelect(apiary: Apiary) {
    selectedApiary = apiary;
}

function handleBack() {
    if (selectedApiary) {
        selectedApiary = null;
    } else if (selectedRegion) {
        selectedRegion = null;
    }
}
</script>

<div class="space-y-6">
    {#if selectedApiary}
    <div in:fade>
        <HivesList
            {selectedApiary}
            onBack={handleBack}
            />
            </div>
            {:else if selectedRegion}
            <div in:fade>
                <ApiarySelect
                    {selectedRegion}
                    onSelect={handleApiarySelect}
                    onBack={handleBack}
                    />
                    </div>
                    {:else}
                    <div in:fade>
                        <RegionSelect
                            onSelect={handleRegionSelect}
                            />
                            </div>
                            {/if}
                            </div>
