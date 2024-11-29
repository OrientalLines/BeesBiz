<script lang="ts">
	import Modal from './Modal.svelte';
	import type { User } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import Icon from '@iconify/svelte';

	export let isOpen = false;
	export let user: User | null = null;
	export let onClose = () => {};
	export let onSave = (updatedUser: User) => {};

	const toastStore = getToastStore();
	let formData: Partial<User> = {};
	let isLoading = false;
	let avatarPreview: string | null = null;

	$: if (user) {
		formData = {
			user_id: user.user_id,
			username: user.username,
			full_name: user.full_name,
			email: user.email,
			role: user.role
		};
		avatarPreview = `https://api.dicebear.com/7.x/avataaars/svg?seed=${user.email}`;
	}

	async function handleSubmit() {
		if (!formData || !formData.username) return;

		isLoading = true;
		try {
			await onSave(formData as User);
			toastStore.trigger({
				message: '✨ User updated successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to update user',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}
</script>

<Modal {isOpen} title="Edit User Profile" {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade class="p-4">
		<form on:submit|preventDefault={handleSubmit} class="space-y-6">
			<!-- Form Fields -->
			<div class="space-y-4">
				<!-- Username -->
				<div>
					<div class="flex items-center gap-2 mb-2">
						<div
							class="w-9 h-9 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center"
						>
							<Icon icon="mdi:account" class="w-5 h-5 text-amber-500" />
						</div>
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300">Username</label>
					</div>
					<input
						type="text"
						bind:value={formData.username}
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-800 border border-amber-200/70
                        dark:border-amber-800/30 rounded-lg focus:ring-2 focus:ring-amber-500/50
                        focus:border-amber-500 transition-colors"
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
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300">Full Name</label>
					</div>
					<input
						type="text"
						bind:value={formData.full_name}
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-800 border border-amber-200/70
                        dark:border-amber-800/30 rounded-lg focus:ring-2 focus:ring-amber-500/50
                        focus:border-amber-500 transition-colors"
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
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300">Email</label>
					</div>
					<input
						type="email"
						bind:value={formData.email}
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-800 border border-amber-200/70
                        dark:border-amber-800/30 rounded-lg focus:ring-2 focus:ring-amber-500/50
                        focus:border-amber-500 transition-colors"
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
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300">Role</label>
					</div>
					<select
						bind:value={formData.role}
						class="w-full px-4 py-2.5 bg-white dark:bg-gray-800 border border-amber-200/70
                        dark:border-amber-800/30 rounded-lg focus:ring-2 focus:ring-amber-500/50
                        focus:border-amber-500 transition-colors appearance-none cursor-pointer"
					>
						<option value="ADMIN">Admin</option>
						<option value="MANAGER">Manager</option>
						<option value="WORKER">Worker</option>
					</select>
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-center gap-4 pt-6 border-t border-gray-200 dark:border-gray-700">
				<button
					type="button"
					class="w-full px-6 py-3 text-sm font-medium text-gray-700 dark:text-gray-300
                    bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700/50
                    rounded-lg transition-colors border border-gray-200 dark:border-gray-700
                    min-w-[100px]"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="w-full px-6 py-3 text-sm font-medium text-white bg-gradient-to-r
                    from-amber-500 to-orange-500 hover:from-amber-600 hover:to-orange-600
                    rounded-lg transition-colors flex items-center justify-center gap-2
                    shadow-lg hover:shadow-xl disabled:opacity-70 disabled:cursor-not-allowed
                    min-w-[100px]"
					disabled={isLoading}
				>
					{#if isLoading}
						<Icon icon="mdi:loading" class="w-4 h-4 animate-spin" />
					{/if}
					<span>{isLoading ? 'Saving...' : 'Save Changes'}</span>
				</button>
			</div>
		</form>
	</div>
</Modal>
