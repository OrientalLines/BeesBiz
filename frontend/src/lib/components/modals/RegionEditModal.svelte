<script lang="ts">
	import type { Region } from '$lib/types';
	import Modal from './Modal.svelte';

	export let isOpen: boolean;
	export let region: Region | null;
	export let onClose: () => void;
	export let onSave: (region: Region) => void;

	const climateZones = ['tropical', 'temperate', 'continental', 'mediterranean'];

	let formData: Omit<Region, 'region_id'> = region
		? { name: region.name, climate_zone: region.climate_zone }
		: { name: '', climate_zone: 'temperate' };

	function handleSubmit() {
		onSave({
			...(region || {}),
			...formData
		} as Region);
	}
</script>

<Modal {isOpen} title={region ? 'Edit Region' : 'New Region'} {onClose}>
	<form on:submit|preventDefault={handleSubmit} class="space-y-4">
		<div>
			<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Region Name
			</label>
			<input
				type="text"
				bind:value={formData.name}
				class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600
                rounded-lg focus:ring-2 focus:ring-amber-500 dark:bg-gray-700"
				required
			/>
		</div>

		<div>
			<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Climate Zone
			</label>
			<select
				bind:value={formData.climate_zone}
				class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600
                rounded-lg focus:ring-2 focus:ring-amber-500 dark:bg-gray-700"
			>
				{#each climateZones as zone}
					<option value={zone}>{zone}</option>
				{/each}
			</select>
		</div>

		<div class="flex justify-end gap-3 pt-4">
			<button
				type="button"
				class="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100
                dark:hover:bg-gray-700 rounded-lg"
				on:click={onClose}
			>
				Cancel
			</button>
			<button
				type="submit"
				class="px-4 py-2 bg-amber-500 text-white rounded-lg hover:bg-amber-600
                transition-colors"
			>
				{region ? 'Update' : 'Create'}
			</button>
		</div>
	</form>
</Modal>
