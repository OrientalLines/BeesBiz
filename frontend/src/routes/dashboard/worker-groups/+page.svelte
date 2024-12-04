<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { User, WorkerGroup } from '$lib/types';
	import { getToastStore, Toast } from '@skeletonlabs/skeleton';
	import { onMount } from 'svelte';
	import {
		deleteWorkerGroup,
		getFreeUsers,
		getWorkerGroupById,
		getWorkerGroups
	} from '$lib/services/api';
	import WorkerGroupAddModal from '$lib/components/modals/WorkerGroupAddModal.svelte';
	import WorkerGroupEditModal from '$lib/components/modals/WorkerGroupEditModal.svelte';
	import WorkerGroupDeleteModal from '$lib/components/modals/WorkerGroupDeleteModal.svelte';

	const toastStore = getToastStore();

	let workerGroups: (WorkerGroup & { members: User[] })[] = [];
	let freeUsers: User[] = [];
	let loading = true;
	let searchQuery = '';
	let addModalOpen = false;
	let editModalOpen = false;
	let deleteModalOpen = false;
	let selectedGroup: (WorkerGroup & { members: User[] }) | null = null;

	onMount(async () => {
		await Promise.all([loadWorkerGroups(), loadFreeUsers()]);
	});

	async function loadWorkerGroups() {
		try {
			loading = true;
			const groups = await getWorkerGroups();
			workerGroups = await Promise.all(
				groups.map(async (group) => ({
					...group,
					members: await getWorkerGroupById(group.group_id).then((data) => data.members)
				}))
			);
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load worker groups',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});
		} finally {
			loading = false;
		}
	}

	async function loadFreeUsers() {
		try {
			freeUsers = await getFreeUsers();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load available workers',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});
		}
	}

	async function handleDeleteGroup(groupId: number) {
		try {
			await deleteWorkerGroup(groupId);
			await Promise.all([loadWorkerGroups(), loadFreeUsers()]);
			toastStore.trigger({
				message: '✨ Worker group deleted successfully',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to delete worker group',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium'
			});
		}
	}

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
	}, 300);

	$: filteredGroups = workerGroups.filter((group) =>
		group.group_name.toLowerCase().includes(searchQuery.toLowerCase())
	);
	function handleEdit(group: WorkerGroup & { members: User[] }) {
		selectedGroup = group;
		editModalOpen = true;
	}
	function handleDelete(group: WorkerGroup & { members: User[] }) {
		selectedGroup = group;
		deleteModalOpen = true;
	}

	async function handleAddGroup(groupData: Omit<WorkerGroup, 'group_id'>) {
		await loadWorkerGroups();
		await loadFreeUsers();
	}
</script>

<div class="max-w mx-auto space-y-8" in:fade>
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:account-group" class="w-8 h-8 text-amber-500" />
				Worker Groups
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Manage worker groups and their members</p>
		</div>
		<button
			class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-lg
				transition-colors duration-200 flex items-center gap-2"
			on:click={() => (addModalOpen = true)}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Create Group
		</button>
	</div>

	<div class="relative">
		<Icon
			icon="mdi:magnify"
			class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
		/>
		<input
			type="text"
			placeholder="Search groups..."
			class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
				bg-white dark:bg-gray-800 text-gray-900 dark:text-white
				focus:ring-2 focus:ring-amber-500 focus:border-transparent"
			on:input={(e) => debouncedSearch(e.currentTarget.value)}
		/>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-8">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500" />
		</div>
	{:else}
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
			<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
				<thead class="bg-gray-50 dark:bg-gray-700">
					<tr>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Group Name
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Members Count
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Updated At
						</th>
						<th
							class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each filteredGroups as group}
						<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{group.group_name}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{group?.members?.length ?? 0}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{group.updated_at ? new Date(group.updated_at).toLocaleString() : ''}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
								<button
									class="text-amber-600 hover:text-amber-900 dark:hover:text-amber-400 mr-3"
									on:click={() => handleEdit(group)}
								>
									<Icon icon="mdi:pencil" class="w-5 h-5" />
								</button>
								<button
									class="text-red-600 hover:text-red-900 dark:hover:text-red-400"
									on:click={() => handleDelete(group)}
								>
									<Icon icon="mdi:delete" class="w-5 h-5" />
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<WorkerGroupAddModal
	isOpen={addModalOpen}
	{freeUsers}
	onClose={() => (addModalOpen = false)}
	onSave={async (groupData) => {
		await handleAddGroup(groupData);
		addModalOpen = false;
	}}
/>

<WorkerGroupEditModal
	isOpen={editModalOpen}
	group={selectedGroup}
	{freeUsers}
	onClose={() => {
		editModalOpen = false;
		selectedGroup = null;
	}}
	onSave={async () => {
		await loadWorkerGroups();
		await loadFreeUsers();
		editModalOpen = false;
	}}
/>

<WorkerGroupDeleteModal
	isOpen={deleteModalOpen}
	group={selectedGroup}
	onClose={() => {
		deleteModalOpen = false;
		selectedGroup = null;
	}}
	onConfirm={async (groupId) => {
		await handleDeleteGroup(groupId);
		deleteModalOpen = false;
	}}
/>

<Toast position="br" />
