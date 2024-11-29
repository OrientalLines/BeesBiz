<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { User, Region } from '$lib/types';
	import {
		getUsers,
		deleteUser,
		updateUser,
		updateRegionAccess,
		getRegions,
		createUser
	} from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import { onMount } from 'svelte';
	import UserEditModal from '$lib/components/modals/UserEditModal.svelte';
	import UserDeleteModal from '$lib/components/modals/UserDeleteModal.svelte';
	import RegionEditModal from '$lib/components/modals/RegionEditModal.svelte';
	import UserAddModal from '$lib/components/modals/UserAddModal.svelte';
	import { Toast } from '@skeletonlabs/skeleton';

	const toastStore = getToastStore();

	let users: User[] = [];
	let loading = true;
	let searchQuery = '';
	let roleFilter = 'all';
	let editModalOpen = false;
	let deleteModalOpen = false;
	let regionModalOpen = false;
	let selectedUser: User | null = null;
	let regions: Region[] = [];
	let addModalOpen = false;

	onMount(async () => {
		await Promise.all([loadUsers(), loadRegions()]);
	});

	async function loadUsers() {
		try {
			loading = true;
			users = await getUsers();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load users',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium '
			});
		} finally {
			loading = false;
		}
	}

	async function loadRegions() {
		try {
			regions = await getRegions();
			console.log('Loaded regions:', regions);
			if (!regions || regions.length === 0) {
				console.log('Triggering toast for no regions');
				toastStore.trigger({
					message: '⚠️ No regions available',
					background: 'bg-amber-500 rounded-lg',
					classes: 'text-white font-medium '
				});
			}
		} catch (error) {
			console.error('Error loading regions:', error);
			toastStore.trigger({
				message: '❌ Failed to load regions',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium '
			});
		}
	}

	async function handleDeleteUser(userId: number) {
		try {
			await deleteUser(userId);
			await loadUsers();
			toastStore.trigger({
				message: '✨ User deleted successfully',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium '
			});
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to delete user',
				background: 'bg-amber-500 rounded-lg',
				classes: 'text-white font-medium '
			});
		}
	}

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
	}, 300);

	$: filteredUsers = users.filter((user) => {
		const matchesSearch =
			user.full_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			user.email.toLowerCase().includes(searchQuery.toLowerCase());
		const matchesRole = roleFilter === 'all' || user.role === roleFilter;
		return matchesSearch && matchesRole;
	});

	function handleEdit(user: User) {
		selectedUser = user;
		editModalOpen = true;
	}

	function handleDelete(user: User) {
		selectedUser = user;
		deleteModalOpen = true;
	}

	function handleRegionEdit(user: User) {
		selectedUser = user;
		regionModalOpen = true;
	}

	async function handleUserUpdate(updatedUser: User) {
		try {
			await updateUser(updatedUser);
			await loadUsers();
		} catch (error) {
			throw error;
		}
	}

	async function handleRegionUpdate(userId: number, regionIds: number[]) {
		try {
			await updateRegionAccess(userId, regionIds);
		} catch (error) {
			throw error;
		}
	}

	async function handleCreateUser(userData: Omit<User, 'user_id'>, selectedRegions: number[]) {
		try {
			const newUser = await createUser(userData);
			if (selectedRegions.length > 0) {
				await updateRegionAccess(newUser.user_id, selectedRegions);
			}
			await loadUsers();
		} catch (error) {
			throw error;
		}
	}
</script>

<div class="max-w mx-auto space-y-8" in:fade>
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:users" class="w-8 h-8 text-amber-500" />
				Users Overview
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Manage system users and their roles</p>
		</div>
		<button
			class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-lg
				transition-colors duration-200 flex items-center gap-2"
			on:click={() => (addModalOpen = true)}
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Add User
		</button>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
		<div class="relative">
			<Icon
				icon="mdi:magnify"
				class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
			/>
			<input
				type="text"
				placeholder="Search users..."
				class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
					bg-white dark:bg-gray-800 text-gray-900 dark:text-white
					focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<select
			bind:value={roleFilter}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
				bg-white dark:bg-gray-800 text-gray-900 dark:text-white
				focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="all">All Roles</option>
			<option value="ADMIN">Admin</option>
			<option value="MANAGER">Manager</option>
			<option value="WORKER">Worker</option>
		</select>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-8">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500"></div>
		</div>
	{:else}
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
			<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
				<thead class="bg-gray-50 dark:bg-gray-700">
					<tr>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Name
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Email
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							settingsRole
						</th>
						<th
							class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Last Login
						</th>
						<th
							class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each filteredUsers as user}
						<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{user.full_name}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{user.email}
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span
									class="px-2 py-1 text-xs font-semibold rounded-full
									{user.role === 'ADMIN'
										? 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'
										: user.role === 'MANAGER'
											? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300'
											: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'}"
								>
									{user.role}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
								{new Date(user.last_login).toLocaleString()}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
								<button
									class="text-amber-600 hover:text-amber-900 dark:hover:text-amber-400 mr-3"
									on:click={() => handleEdit(user)}
								>
									<Icon icon="mdi:pencil" class="w-5 h-5" />
								</button>
								<button
									class="text-amber-600 hover:text-amber-900 dark:hover:text-amber-400 mr-3"
									on:click={() => handleRegionEdit(user)}
								>
									<Icon icon="mdi:map-marker" class="w-5 h-5" />
								</button>
								<button
									class="text-red-600 hover:text-red-900 dark:hover:text-red-400"
									on:click={() => handleDelete(user)}
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

<UserEditModal
	isOpen={editModalOpen}
	user={selectedUser}
	onClose={() => {
		editModalOpen = false;
		selectedUser = null;
	}}
	onSave={handleUserUpdate}
/>

<UserDeleteModal
	isOpen={deleteModalOpen}
	user={selectedUser}
	onClose={() => {
		deleteModalOpen = false;
		selectedUser = null;
	}}
	onConfirm={handleDeleteUser}
/>

<RegionEditModal
	isOpen={regionModalOpen}
	user={selectedUser}
	{regions}
	onClose={() => {
		regionModalOpen = false;
		selectedUser = null;
	}}
	onSave={handleRegionUpdate}
/>

<UserAddModal
	isOpen={addModalOpen}
	{regions}
	onClose={() => (addModalOpen = false)}
	onSave={handleCreateUser}
/>

<Toast position="br" />
