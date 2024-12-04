<script lang="ts">
	import Modal from './Modal.svelte';

	import { getToastStore } from '@skeletonlabs/skeleton';
	import type { User, WorkerGroup } from '$lib/types';
	import {
		getWorkerGroupById,
		updateWorkerGroup,
		addGroupMember,
		removeGroupMember
	} from '$lib/services/api';

	export let isOpen: boolean;
	export let group: WorkerGroup | null;
	export let freeUsers: User[] = [];
	export let onClose: () => void;
	export let onSave: () => Promise<void>;

	const toastStore = getToastStore();

	let groupName = '';
	let selectedUserIds: number[] = [];
	let loading = false;
	let loadingMembers = false;

	$: if (group) {
		groupName = group.group_name;
		loadingMembers = true;
		getWorkerGroupById(group.group_id).then((data) => {
			selectedUserIds = data.members.map((member) => member.user_id);
			loadingMembers = false;
		});
	}

	let availableUsers: User[] = [];
	$: if (group) {
		loadingMembers = true;
		getWorkerGroupById(group.group_id).then((data) => {
			availableUsers = data.members;
			if (freeUsers.length > 0) {
				availableUsers = [...freeUsers, ...data.members];
			}
			loadingMembers = false;
		});
	} else {
		availableUsers = freeUsers;
	}

	async function handleSubmit() {
		if (!group) return;

		if (!groupName.trim()) {
			toastStore.trigger({
				message: '⚠️ Please enter a group name',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});
			return;
		}

		try {
			loading = true;

			// Get current group members
			const currentGroup = await getWorkerGroupById(group.group_id);
			const currentMemberIds = currentGroup.members.map((member) => member.user_id);

			// Update group name
			await updateWorkerGroup(group.group_id, {
				group_name: groupName,
				manager_id: group.manager_id,
				created_at: group.created_at,
				updated_at: new Date().toISOString()
			});

			// Determine which members to add and remove
			const membersToAdd = selectedUserIds.filter((id) => !currentMemberIds.includes(id));
			const membersToRemove = currentMemberIds.filter((id) => !selectedUserIds.includes(id));

			// Process member changes
			await Promise.all([
				...membersToAdd.map((userId) => addGroupMember(group.group_id, userId)),
				...membersToRemove.map((userId) => removeGroupMember(group.group_id, userId))
			]);

			toastStore.trigger({
				message: '✨ Worker group updated successfully',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});

			await onSave();
			handleClose();
		} catch (error) {
			console.error('Failed to update worker group:', error);
			toastStore.trigger({
				message: '❌ Failed to update worker group',
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

<Modal {isOpen} onClose={handleClose} title="Edit Worker Group">
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
				<label class="block text-sm font-medium mb-1" for="members">Members</label>
				<div
					class="border-2 border-amber-500/50 rounded-lg overflow-y-auto dark:bg-gray-800 h-[calc(100%-2rem)]"
				>
					<div class="p-2 space-y-2">
						{#if loadingMembers}
							<div class="flex justify-center py-4">
								<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500"></div>
							</div>
						{:else}
							{#each availableUsers as user}
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
										<span class="text-sm text-gray-500 dark:text-gray-400 ml-1">({user.email})</span
										>
									</div>
								</label>
							{/each}
						{/if}
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
					{loading ? 'Saving...' : 'Save Changes'}
				</button>
			</div>
		</form>
	</div>
</Modal>
