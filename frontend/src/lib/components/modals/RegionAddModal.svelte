<script lang="ts">
	import Modal from './Modal.svelte';
	import type { RegionFormData } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';

	export let isOpen = false;
	export let onClose = () => {};
	export let onSave = (regionData: RegionFormData) => {};

	const toastStore = getToastStore();
	let isLoading = false;
	let formData: RegionFormData = {
		name: '',
		climate_zone: ''
	};

	async function handleSubmit() {
		if (!formData.name || !formData.climate_zone) {
			toastStore.trigger({
				message: '❌ Please fill in all required fields',
				background: 'variant-filled-error'
			});
			return;
		}

		isLoading = true;
		try {
			await onSave(formData);
			toastStore.trigger({
				message: '✨ Region created successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to create region',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}
</script>

<Modal {isOpen} title="Add New Region" {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<div class="space-y-4">
				<div class="text-center space-y-3" in:fly={{ y: 20, duration: 300, delay: 100 }}>
					<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Create New Region</h3>
					<p class="text-gray-600 dark:text-gray-400">Enter the details for the new region</p>
				</div>

				<!-- Region Form Fields -->
				<div class="space-y-4 mt-6">
					<div class="space-y-2">
						<label for="name" class="text-sm font-medium text-gray-700 dark:text-gray-300">
							Region Name *
						</label>
						<input
							type="text"
							id="name"
							bind:value={formData.name}
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-amber-500
								focus:border-amber-500 dark:bg-gray-700 dark:border-gray-600
								dark:text-white"
							placeholder="Enter region name"
							required
						/>
					</div>

					<div class="space-y-2">
						<label for="climate_zone" class="text-sm font-medium text-gray-700 dark:text-gray-300">
							Climate Zone *
						</label>
						<input
							type="text"
							id="climate_zone"
							bind:value={formData.climate_zone}
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-amber-500
								focus:border-amber-500 dark:bg-gray-700 dark:border-gray-600
								dark:text-white"
							placeholder="Enter climate zone"
							required
						/>
					</div>
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 pt-6 border-t border-gray-200 dark:border-gray-700">
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
						transition-colors duration-200 flex items-center gap-2 disabled:opacity-50"
					disabled={isLoading}
				>
					{#if isLoading}
						<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
					{/if}
					{isLoading ? 'Creating...' : 'Create Region'}
				</button>
			</div>
		</form>
	</div>
</Modal>
