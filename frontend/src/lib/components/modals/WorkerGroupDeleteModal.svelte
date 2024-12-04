<script lang="ts">
	import Modal from './Modal.svelte';
	import type { WorkerGroup } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';

	export let isOpen: boolean;
	export let group: WorkerGroup | null;
	export let onClose: () => void;
	export let onConfirm: (groupId: number) => Promise<void>;

	const toastStore = getToastStore();
	let loading = false;

	async function handleDelete() {
		if (!group) return;

		loading = true;
		try {
			await onConfirm(group.group_id);
			toastStore.trigger({
				message: '✨ Worker group deleted successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to delete worker group',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}
</script>

<Modal {isOpen} title="Delete Worker Group" {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade class="space-y-6">
		<!-- Warning Icon -->
		<div class="flex justify-center" in:fly={{ y: 20, duration: 300, delay: 100 }}>
			<div class="w-16 h-16 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center">
				<svg class="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
					/>
				</svg>
			</div>
		</div>

		<!-- Warning Message -->
		<div class="text-center space-y-3" in:fly={{ y: 20, duration: 300, delay: 200 }}>
			<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Delete Worker Group</h3>
			<p class="text-gray-600 dark:text-gray-400">
				Are you sure you want to delete the worker group
				<span class="font-semibold text-red-500">"{group?.group_name}"</span>?
			</p>
			<p class="text-sm text-gray-500 dark:text-gray-500 bg-gray-50 dark:bg-gray-800/50 p-3 rounded-lg">
				This action cannot be undone. All data associated with this worker group will be permanently removed.
			</p>
		</div>

		<!-- Action Buttons -->
		<div class="flex justify-center gap-3 pt-4" in:fly={{ y: 20, duration: 300, delay: 300 }}>
			<button
				type="button"
				class="px-6 py-2.5 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700
					rounded-lg transition-all duration-300 transform hover:scale-105"
				on:click={onClose}
				disabled={loading}
			>
				Cancel
			</button>
			<button
				type="button"
				class="px-6 py-2.5 bg-gradient-to-r from-red-500 to-red-600 text-white rounded-lg
					hover:from-red-600 hover:to-red-700 transition-all duration-300 transform hover:scale-105
					disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
				on:click={handleDelete}
				disabled={loading}
			>
				{#if loading}
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
				{loading ? 'Deleting...' : 'Delete Worker Group'}
			</button>
		</div>
	</div>

	<Toast position="tr" />
</Modal>
