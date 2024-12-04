<script lang="ts">
	import { Modal } from '@skeletonlabs/skeleton';
	import type { WorkerGroup } from '$lib/types';

	export let isOpen: boolean;
	export let group: WorkerGroup | null;
	export let onClose: () => void;
	export let onConfirm: (groupId: number) => Promise<void>;

	let loading = false;

	async function handleDelete() {
		if (!group) return;
		
		try {
			loading = true;
			await onConfirm(group.group_id);
			handleClose();
		} finally {
			loading = false;
		}
	}

	function handleClose() {
		onClose();
	}
</script>

<Modal bind:open={isOpen} on:close={handleClose}>
	<div class="p-4 w-full max-w-md">
		<h3 class="text-2xl font-bold mb-4">Delete Worker Group</h3>
		<p class="mb-4">
			Are you sure you want to delete the worker group "{group?.name}"? This action cannot be undone.
		</p>
		<div class="flex justify-end space-x-2">
			<button
				type="button"
				class="px-4 py-2 text-gray-700 bg-gray-200 rounded-lg hover:bg-gray-300"
				on:click={handleClose}
				disabled={loading}
			>
				Cancel
			</button>
			<button
				type="button"
				class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 disabled:opacity-50"
				on:click={handleDelete}
				disabled={loading}
			>
				{loading ? 'Deleting...' : 'Delete'}
			</button>
		</div>
	</div>
</Modal> 