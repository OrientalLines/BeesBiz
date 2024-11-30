<script lang="ts">
	import { fade } from 'svelte/transition';
	import RegionSelect from '$lib/components/hives/RegionSelect.svelte';
	import ApiarySelect from '$lib/components/hives/ApiarySelect.svelte';
	import HivesList from '$lib/components/hives/HivesList.svelte';
	import SensorsDashboard from '$lib/components/iot/SensorsDashboard.svelte';
	import type { Region, Apiary, Hive, Sensor } from '$lib/types';

	let selectedRegion: Region | undefined = undefined;
	let selectedApiary: Apiary | undefined = undefined;
	let selectedHive: Hive | undefined = undefined;
	let selectedSensor: Sensor | undefined = undefined;

	function handleRegionSelect(region: Region) {
		selectedRegion = region;
		selectedApiary = undefined;
		selectedHive = undefined;
		selectedSensor = undefined;
	}

	function handleApiarySelect(apiary: Apiary) {
		selectedApiary = apiary;
		selectedHive = undefined;
		selectedSensor = undefined;
	}

	function handleHiveSelect(hive: Hive) {
		selectedHive = hive;
		selectedSensor = undefined;
	}

	async function handleSensorSelect(sensor: Sensor) {
		selectedSensor = sensor;
	}

	function handleBack() {
		if (selectedSensor) {
			selectedSensor = undefined;
		} else if (selectedHive) {
			selectedHive = undefined;
		} else if (selectedApiary) {
			selectedApiary = undefined;
		} else if (selectedRegion) {
			selectedRegion = undefined;
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
