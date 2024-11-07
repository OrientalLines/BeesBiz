<script lang="ts">
	import { fade } from 'svelte/transition';
	import RegionSelect from '$lib/components/hives/RegionSelect.svelte';
	import ApiarySelect from '$lib/components/hives/ApiarySelect.svelte';
	import HivesList from '$lib/components/hives/HivesList.svelte';
	import SensorsDashboard from '$lib/components/iot/SensorsDashboard.svelte';
	import type { Region, Apiary, Hive, Sensor } from '$lib/types';

	let selectedRegion: Region | null = null;
	let selectedApiary: Apiary | null = null;
	let selectedHive: Hive | null = null;
	let selectedSensor: Sensor | null = null;

	function handleRegionSelect(region: Region) {
		selectedRegion = region;
		selectedApiary = null;
		selectedHive = null;
		selectedSensor = null;
	}

	function handleApiarySelect(apiary: Apiary) {
		selectedApiary = apiary;
		selectedHive = null;
		selectedSensor = null;
	}

	function handleHiveSelect(hive: Hive) {
		selectedHive = hive;
		selectedSensor = null;
	}

	async function handleSensorSelect(sensor: Sensor) {
		selectedSensor = sensor;
	}

	function handleBack() {
		if (selectedSensor) {
			selectedSensor = null;
		} else if (selectedHive) {
			selectedHive = null;
		} else if (selectedApiary) {
			selectedApiary = null;
		} else if (selectedRegion) {
			selectedRegion = null;
		}
	}
</script>

<div class="space-y-6">
	{#if selectedHive}
		<div in:fade>
			<SensorsDashboard {selectedHive} onBack={handleBack} onSelect={handleSensorSelect} />
		</div>
	{:else if selectedApiary}
		<div in:fade>
			<HivesList {selectedApiary} onSelect={handleHiveSelect} onBack={handleBack} />
		</div>
	{:else if selectedRegion}
		<div in:fade>
			<ApiarySelect {selectedRegion} onSelect={handleApiarySelect} onBack={handleBack} />
		</div>
	{:else}
		<div in:fade>
			<RegionSelect onSelect={handleRegionSelect} />
		</div>
	{/if}
</div>
