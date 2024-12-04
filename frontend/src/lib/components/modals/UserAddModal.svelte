<script lang="ts">
	import Modal from './Modal.svelte';
	import type { User, Region, Role } from '$lib/types';
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
	let formData: {
		username: string;
		full_name: string;
		email: string;
		password: string;
		role: Role;
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
					...formData
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
	<div in:fly={{ y: 50, duration: 400 }} out:fade class="p-4">
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<!-- User Information -->
			<div class="space-y-4">
				<!-- Username -->
				<div>
					<div class="flex items-center gap-2 mb-2">
						<div
							class="w-9 h-9 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center"
						>
							<Icon icon="mdi:account" class="w-5 h-5 text-amber-500" />
						</div>
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300" for="username">Username</label>
					</div>
					<input
						type="text"
						bind:value={formData.username}
						required
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-amber-500 transition-colors"
					/>
				</div>

				<!-- Full Name -->
				<div>
					<div class="flex items-center gap-2 mb-2">
						<div
							class="w-9 h-9 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center"
						>
							<Icon icon="mdi:card-account-details" class="w-5 h-5 text-amber-500" />
						</div>
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300" for="full_name">Full Name</label>
					</div>
					<input
						type="text"
						bind:value={formData.full_name}
						required
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-amber-500 transition-colors"
					/>
				</div>

				<!-- Email -->
				<div>
					<div class="flex items-center gap-2 mb-2">
						<div
							class="w-9 h-9 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center"
						>
							<Icon icon="mdi:email" class="w-5 h-5 text-amber-500" />
						</div>
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300" for="email">Email</label>
					</div>
					<input
						type="email"
						bind:value={formData.email}
						required
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-amber-500 transition-colors"
					/>
				</div>

				<!-- Password -->
				<div>
					<div class="flex items-center gap-2 mb-2">
						<div
							class="w-9 h-9 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center"
						>
							<Icon icon="mdi:lock" class="w-5 h-5 text-amber-500" />
						</div>
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300" for="password">Password</label>
					</div>
					<input
						type="password"
						bind:value={formData.password}
						required
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-amber-500 transition-colors"
					/>
				</div>

				<!-- Role -->
				<div>
					<div class="flex items-center gap-2 mb-2">
						<div
							class="w-9 h-9 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center"
						>
							<Icon icon="mdi:shield-account" class="w-5 h-5 text-amber-500" />
						</div>
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300" for="role">Role</label>
					</div>
					<select
						bind:value={formData.role}
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-amber-500 transition-colors appearance-none cursor-pointer"
					>
						<option value="ADMIN">Admin</option>
						<option value="MANAGER">Manager</option>
						<option value="WORKER">Worker</option>
					</select>
				</div>
			</div>

			<!-- Region Access -->
			<div class="space-y-4">
				<div class="flex items-center gap-2 mb-2">
					<div
						class="w-9 h-9 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center"
					>
						<Icon icon="mdi:map-marker" class="w-5 h-5 text-amber-500" />
					</div>
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300" for="region_access">Region Access</label>
				</div>

				<div class="space-y-2 max-h-48 overflow-y-auto">
					{#each regions as region}
						<label class="flex items-center p-3 rounded-lg border border-gray-300 dark:border-gray-600 
							bg-white dark:bg-gray-900 hover:bg-gray-50 dark:hover:bg-gray-800 
							transition-colors duration-200">
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
								class="w-4 h-4 text-amber-500 border-gray-300 rounded focus:ring-amber-500"
							/>
							<span class="ml-3 text-gray-700 dark:text-gray-300">{region.name}</span>
						</label>
					{/each}
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-center gap-4 pt-6 border-t border-gray-200 dark:border-gray-700">
				<button
					type="button"
					class="w-full px-6 py-3 text-sm font-medium text-gray-700 dark:text-gray-300 
					bg-white dark:bg-gray-900 hover:bg-gray-50 dark:hover:bg-gray-800 
					rounded-lg transition-colors border border-gray-300 dark:border-gray-600 min-w-[100px]"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="w-full px-6 py-3 text-sm font-medium text-white 
					bg-amber-500 hover:bg-amber-600 dark:bg-amber-600 dark:hover:bg-amber-700
					rounded-lg transition-colors flex items-center justify-center gap-2 
					shadow-lg hover:shadow-xl disabled:opacity-70 disabled:cursor-not-allowed min-w-[100px]"
					disabled={isLoading}
				>
					{#if isLoading}
						<Icon icon="mdi:loading" class="w-4 h-4 animate-spin" />
					{/if}
					<span>{isLoading ? 'Creating...' : 'Create User'}</span>
				</button>
			</div>
		</form>
	</div>
</Modal>

<style>
	/* Remove the previous style block as it's no longer needed */
</style>
