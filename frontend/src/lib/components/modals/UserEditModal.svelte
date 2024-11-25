<script lang="ts">
	import Modal from './Modal.svelte';
	import type { User } from '$lib/types';
	import type { Role } from '$lib/services/api';
	import { fade, fly } from 'svelte/transition';
	import { Avatar } from '@skeletonlabs/skeleton';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';

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
	<div in:fly={{ y: 50, duration: 400 }} out:fade>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<!-- Form Fields with Modern Styling -->
			<div class="space-y-4">
				<div class="relative">
					<input type="text" bind:value={formData.username} placeholder="Username" class="input" />
				</div>
				<div class="relative">
					<input
						type="text"
						bind:value={formData.full_name}
						placeholder="Full Name"
						class="input"
					/>
				</div>
				<div class="relative">
					<input type="email" bind:value={formData.email} placeholder="Email" class="input" />
				</div>
				<div class="relative">
					<select bind:value={formData.role} class="select">
						<option value="ADMIN">Admin</option>
						<option value="MANAGER">Manager</option>
						<option value="WORKER">Worker</option>
					</select>
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 mt-8" in:fly={{ y: 20, duration: 300, delay: 600 }}>
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
