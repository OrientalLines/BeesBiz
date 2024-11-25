<script lang="ts">
	import Modal from './Modal.svelte';
	import type { Region, User } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import Icon from '@iconify/svelte';
	import { getUserAllowedRegions } from '$lib/services/api';

	export let isOpen = false;
	export let user: User | null = null;
	export let regions: Region[] = [];
	export let onClose = () => {};
	export let onSave = (userId: number, regionIds: number[]) => {};

	const toastStore = getToastStore();
	let selectedRegions: number[] = [];
	let isLoading = false;

	// Load user's allowed regions when the user changes
	$: if (user) {
		loadUserRegions();
	}

	async function loadUserRegions() {
		if (!user) return;
		try {
			const allowedRegions = await getUserAllowedRegions(user.user_id);
			if (allowedRegions) {
				selectedRegions = allowedRegions.map(region => region.region_id);
			} else {
				selectedRegions = [];
			}
		} catch (error) {
			console.error('Error loading user regions:', error);
			selectedRegions = [];
		}
	}

	async function handleSubmit() {
		if (!user) return;

		isLoading = true;
		try {
			await onSave(user.user_id, selectedRegions);
			toastStore.trigger({
				message: '✨ Region access updated successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to update region access',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}
</script>

<Modal {isOpen} title="Edit Region Access" {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<div class="space-y-4">
				<h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
					<Icon icon="mdi:map-marker" class="w-5 h-5 text-amber-500" />
					Assign Regions
				</h3>

				<div class="space-y-2">
					{#each regions as region}
						<label
							class="flex items-center p-3 rounded-lg border border-gray-200
							dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800
							transition-colors duration-200"
						>
							<input
								type="checkbox"
								value={region.region_id}
								checked={selectedRegions.includes(region.region_id)}
								on:change={(e) => {
									if (e.currentTarget.checked) {
										selectedRegions = [...selectedRegions, region.region_id];
									} else {
										selectedRegions = selectedRegions.filter((id) => id !== region.region_id);
									}
								}}
								class="w-4 h-4 text-amber-500 border-gray-300 rounded
									focus:ring-amber-500 dark:focus:ring-offset-gray-800"
							/>
							<span class="ml-3 text-gray-900 dark:text-gray-300">{region.name}</span>
						</label>
					{/each}
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 pt-6">
				<button
					type="button"
					class="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100
						dark:hover:bg-gray-700 rounded-lg transition-colors duration-200"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-lg
						transition-colors duration-200 flex items-center gap-2"
					disabled={isLoading}
				>
					{#if isLoading}
						<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
					{/if}
					{isLoading ? 'Saving...' : 'Save Changes'}
				</button>
			</div>
		</form>
	</div>
</Modal>
