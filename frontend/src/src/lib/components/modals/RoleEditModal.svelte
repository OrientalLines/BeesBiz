<script lang="ts">
	import Modal from './Modal.svelte';
	import type { Role } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import { Shield, ShieldCheck, ShieldAlert } from 'lucide-svelte';

	export let isOpen = false;
	export let role: Role | null = null;
	export let onClose = () => {};
	export let onSave = (updatedRole: Role) => {};

	const toastStore = getToastStore();
	let formData: Partial<Role> = {};
	let isLoading = false;

	const availablePermissions = [
		{
			id: 'manage_users',
			label: 'Manage Users',
			description: 'Can create, edit and delete user accounts',
			icon: Shield
		},
		{
			id: 'manage_roles',
			label: 'Manage Roles',
			description: 'Can create and modify role permissions',
			icon: ShieldCheck
		},
		{
			id: 'manage_system',
			label: 'Manage System',
			description: 'Full system administration access',
			icon: ShieldAlert
		},
		{
			id: 'view_reports',
			label: 'View Reports',
			description: 'Access to view and export reports'
		},
		{
			id: 'manage_beekeepers',
			label: 'Manage Beekeepers',
			description: 'Can manage beekeeper profiles and assignments'
		},
		{
			id: 'manage_hives',
			label: 'Manage Hives',
			description: 'Can create and modify hive information'
		},
		{
			id: 'record_observations',
			label: 'Record Observations',
			description: 'Can add and edit hive observations'
		}
	];

	$: if (role) {
		formData = { ...role };
	}

	async function handleSubmit() {
		if (!formData || !formData.id) return;

		isLoading = true;
		try {
			await onSave(formData as Role);
			toastStore.trigger({
				message: '✨ Role updated successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to update role',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}
</script>

<Modal {isOpen} title="Edit Role" {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<!-- Role Info Section -->
			<div class="space-y-4" in:fly={{ y: 20, duration: 300, delay: 100 }}>
				<!-- Role Name Input -->
				<div class="relative">
					<input
						type="text"
						bind:value={formData.name}
						placeholder=" "
						class="peer w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
					/>
					<label
						class="absolute left-2 -top-2.5 bg-white dark:bg-gray-800 px-2 text-sm text-amber-500
                            transition-all duration-300"
					>
						Role Name
					</label>
				</div>

				<!-- Role Description Input -->
				<div class="relative">
					<textarea
						bind:value={formData.description}
						rows="3"
						placeholder=" "
						class="peer w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300 resize-none"
					/>
					<label
						class="absolute left-2 -top-2.5 bg-white dark:bg-gray-800 px-2 text-sm text-amber-500
                            transition-all duration-300"
					>
						Description
					</label>
				</div>
			</div>

			<!-- Permissions Section -->
			<div class="space-y-4" in:fly={{ y: 20, duration: 300, delay: 200 }}>
				<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 flex items-center gap-2">
					<Shield class="w-5 h-5 text-amber-500" />
					Permissions
				</h3>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					{#each availablePermissions as permission}
						<div
							class="relative flex items-start p-4 rounded-lg border-2 border-gray-200 dark:border-gray-700
                                hover:border-amber-500 dark:hover:border-amber-500 transition-all duration-300
                                {formData.permissions?.includes(permission.id)
								? 'border-amber-500 bg-amber-50 dark:bg-amber-900/10'
								: ''}"
						>
							<div class="flex items-center h-5">
								<input
									type="checkbox"
									checked={formData.permissions?.includes(permission.id)}
									on:change={(e) => {
										const checked = e.currentTarget.checked;
										formData.permissions = checked
											? [...(formData.permissions || []), permission.id]
											: (formData.permissions || []).filter((p) => p !== permission.id);
									}}
									class="h-4 w-4 rounded border-gray-300 text-amber-500 focus:ring-amber-500"
								/>
							</div>
							<div class="ml-3">
								<label class="font-medium text-gray-900 dark:text-gray-100">
									{permission.label}
								</label>
								<p class="text-sm text-gray-500 dark:text-gray-400">
									{permission.description}
								</p>
							</div>
							{#if permission.icon}
								<svelte:component
									this={permission.icon}
									class="absolute right-4 top-4 w-5 h-5 text-amber-500 opacity-50"
								/>
							{/if}
						</div>
					{/each}
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 pt-6" in:fly={{ y: 20, duration: 300, delay: 300 }}>
				<button
					type="button"
					class="px-6 py-2.5 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700
                        rounded-lg transition-all duration-300 transform hover:scale-105"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-6 py-2.5 bg-gradient-to-r from-amber-500 to-amber-600 text-white rounded-lg
                        hover:from-amber-600 hover:to-amber-700 transition-all duration-300 transform hover:scale-105
                        disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
					disabled={isLoading}
				>
					{#if isLoading}
						<svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
								fill="none"
							/>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							/>
						</svg>
					{/if}
					{isLoading ? 'Saving...' : 'Save Changes'}
				</button>
			</div>
		</form>
	</div>

	<Toast position="tr" />
</Modal>
