<script lang="ts">
	import Modal from './Modal.svelte';
	import type { User, Region } from '$lib/types';
	import type { Role } from '$lib/services/api';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import Icon from '@iconify/svelte';

	export let isOpen = false;
	export let regions: Region[] = [];
	export let onClose = () => {};
	export let onSave = (userData: Omit<User, 'user_id'>, selectedRegions: number[]) => {};

	const toastStore = getToastStore();
	let isLoading = false;
	let selectedRegions: number[] = [];

	// Form data
	let formData = {
		username: '',
		full_name: '',
		email: '',
		password: '',
		role: 'WORKER' as Role
	};

	// Reset form when modal opens
	$: if (isOpen) {
		formData = {
			username: '',
			full_name: '',
			email: '',
			password: '',
			role: 'WORKER'
		};
		selectedRegions = [];
	}

	async function handleSubmit() {
		if (!formData.username || !formData.email || !formData.password) {
			toastStore.trigger({
				message: '❌ Please fill in all required fields',
				background: 'variant-filled-error'
			});
			return;
		}

		isLoading = true;
		try {
			await onSave(
				{
					...formData,
					last_login: new Date()
				},
				selectedRegions
			);
			toastStore.trigger({
				message: '✨ User created successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to create user',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}
</script>

<Modal {isOpen} title="Add New User" {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<!-- User Information -->
			<div class="space-y-4">
				<h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
					<Icon icon="mdi:user" class="w-5 h-5 text-amber-500" />
					User Information
				</h3>

				<div class="grid grid-cols-1 gap-4">
					<input
						type="text"
						bind:value={formData.username}
						placeholder="Username *"
						class="input"
						required
					/>
					<input
						type="text"
						bind:value={formData.full_name}
						placeholder="Full Name *"
						class="input"
						required
					/>
					<input
						type="email"
						bind:value={formData.email}
						placeholder="Email *"
						class="input"
						required
					/>
					<input
						type="password"
						bind:value={formData.password}
						placeholder="Password *"
						class="input"
						required
					/>
					<select bind:value={formData.role} class="select">
						<option value="ADMIN">Admin</option>
						<option value="MANAGER">Manager</option>
						<option value="WORKER">Worker</option>
					</select>
				</div>
			</div>

			<!-- Region Access -->
			<div class="space-y-4">
				<h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
					<Icon icon="mdi:map-marker" class="w-5 h-5 text-amber-500" />
					Region Access
				</h3>

				<div class="space-y-2 max-h-48 overflow-y-auto">
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
					{isLoading ? 'Creating...' : 'Create User'}
				</button>
			</div>
		</form>
	</div>
</Modal>

<style lang="postcss">
	.input {
		@apply w-full px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
            focus:ring-2 focus:ring-amber-500 focus:border-transparent;
	}

	.select {
		@apply w-full px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
            focus:ring-2 focus:ring-amber-500 focus:border-transparent;
	}
</style>
