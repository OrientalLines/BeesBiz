<script lang="ts">
	import Modal from './Modal.svelte';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import type { User, WorkerGroup } from '$lib/types';
	import { createWorkerGroup, addGroupMember } from '$lib/services/api';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';

	export let isOpen = false;
	export let freeUsers: User[] = [];
	export let onClose: () => void;
	export let onSave: (groupData: Omit<WorkerGroup, 'group_id'>) => Promise<void>;

	const toastStore = getToastStore();

	let groupName = '';
	let selectedUserIds: number[] = [];
	let loading = false;

	$: selectedUsers = freeUsers.filter((user) => selectedUserIds.includes(user.user_id));

	// Reset form when modal opens
	$: if (isOpen) {
		groupName = '';
		selectedUserIds = [];
	}

	async function handleSubmit() {
		if (!groupName.trim()) {
			toastStore.trigger({
				message: '⚠️ Please enter a group name',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});
			return;
		}

		if (selectedUserIds.length === 0) {
			toastStore.trigger({
				message: '⚠️ Please select at least one user',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});
			return;
		}

		try {
			loading = true;
			const authState = get(auth);
			if (!authState?.user?.user_id) {
				throw new Error('No manager ID available');
			}

			const groupData: Omit<WorkerGroup, 'group_id'> = {
				group_name: groupName,
				manager_id: authState.user.user_id,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString()
			};

			const createdGroup = await createWorkerGroup(groupData);

			// Add selected users to the group
			if (createdGroup.group_id) {
				await Promise.all(
					selectedUserIds.map((userId) => addGroupMember(createdGroup.group_id, userId))
				);
			}

			await onSave(groupData);

			toastStore.trigger({
				message: '✨ Worker group created successfully',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});

			handleClose();
		} catch (error) {
			console.error('Failed to create worker group:', error);
			toastStore.trigger({
				message: '❌ Failed to create worker group',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});
		} finally {
			loading = false;
		}
	}

	function handleClose() {
		groupName = '';
		selectedUserIds = [];
		onClose();
	}
</script>

<Modal {isOpen} title="Create Worker Group" {onClose}>
	<div class="p-4 h-[500px] flex flex-col justify-between">
		<form on:submit|preventDefault={handleSubmit} class="space-y-4 flex-1 flex flex-col">
			<div>
				<label for="groupName" class="block text-sm font-medium mb-1">Group Name</label>
				<input
					type="text"
					id="groupName"
					bind:value={groupName}
					class="w-full px-3 py-2 border rounded-lg dark:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-amber-500"
					placeholder="Enter group name"
				/>
			</div>

			<div class="flex-1 overflow-hidden">
				<label class="block text-sm font-medium mb-1" for="members">Select Members</label>
				<div
					class="border-2 border-amber-500/50 rounded-lg overflow-y-auto dark:bg-gray-800 h-[calc(100%-2rem)]"
				>
					<div class="p-2 space-y-2">
						{#each freeUsers as user}
							<label
								class="flex items-center p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md cursor-pointer transition-colors duration-150"
							>
								<input
									type="checkbox"
									bind:group={selectedUserIds}
									value={user.user_id}
									class="form-checkbox h-5 w-5 text-amber-500 rounded border-gray-300 focus:ring-amber-500 dark:border-gray-600 dark:bg-gray-700 dark:checked:bg-amber-500"
								/>
								<div class="ml-3">
									<span class="font-medium">{user.full_name}</span>
									<span class="text-sm text-gray-500 dark:text-gray-400 ml-1">({user.email})</span>
								</div>
							</label>
						{/each}
					</div>
				</div>
			</div>

			<div class="flex justify-end space-x-2 pt-4">
				<button
					type="button"
					class="px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-200 dark:bg-gray-700 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600"
					on:click={handleClose}
					disabled={loading}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 bg-amber-500 text-white rounded-lg hover:bg-amber-600 disabled:opacity-50"
					disabled={loading}
				>
					{loading ? 'Creating...' : 'Create Group'}
				</button>
			</div>
		</form>
	</div>
</Modal>
