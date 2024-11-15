<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { User } from '$lib/types';
	import UserEditModal from '$lib/components/modals/UserEditModal.svelte';
	import UserDeleteModal from '$lib/components/modals/UserDeleteModal.svelte';

	let searchQuery = '';
	let sortBy: 'name' | 'role' | 'lastActive' = 'name';
	let roleFilter = 'all';
	let editingUser: User | null = null;
	let deletingUser: User | null = null;
	let showAddModal = false;

	// Mock data - replace with API calls
	const users: User[] = [
		{
			id: '1',
			name: 'John Beekeeper',
			email: 'john@example.com',
			role: 'beekeeper',
			lastActive: new Date()
		}
		// Add more mock data as needed
	];

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
	}, 300);

	function handleEdit(user: User) {
		editingUser = { ...user };
	}

	function handleDelete(user: User) {
		deletingUser = { ...user };
	}

	function handleSaveUser(updatedUser: User) {
		// Implementation for saving user
		const index = users.findIndex((u) => u.id === updatedUser.id);
		if (index !== -1) {
			users[index] = updatedUser;
		}
		editingUser = null;
	}

	function handleDeleteUser() {
		if (deletingUser) {
			const index = users.findIndex((u) => u.id === deletingUser.id);
			if (index !== -1) {
				users.splice(index, 1);
			}
		}
		deletingUser = null;
	}

	function handleAddUser(newUser: User) {
		users.push(newUser);
		showAddModal = false;
	}

	$: filteredUsers = users
		.filter((user) => {
			const matchesSearch =
				user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				user.email.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesRole = roleFilter === 'all' || user.role === roleFilter;
			return matchesSearch && matchesRole;
		})
		.sort((a, b) => {
			if (sortBy === 'lastActive')
				return new Date(b.lastActive).getTime() - new Date(a.lastActive).getTime();
			return a[sortBy] > b[sortBy] ? 1 : -1;
		});
</script>

<div class="max-w mx-auto space-y-8" in:fade>
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:account-group" class="w-8 h-8 text-amber-500" />
				User Management
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Manage system users and their roles</p>
		</div>
		<button
			on:click={() => (showAddModal = true)}
			class="bg-amber-500 text-white px-6 py-3 rounded-full
            hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
            flex items-center gap-2"
		>
			<Icon icon="mdi:plus" class="w-5 h-5" />
			Add User
		</button>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
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
                    placeholder-gray-500 dark:placeholder-gray-400
                    focus:ring-2 focus:ring-amber-500 focus:border-transparent"
				on:input={(e) => debouncedSearch(e.currentTarget.value)}
			/>
		</div>

		<select
			bind:value={sortBy}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="name">Sort by Name</option>
			<option value="role">Sort by Role</option>
			<option value="lastActive">Sort by Last Active</option>
		</select>

		<select
			bind:value={roleFilter}
			class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                focus:ring-2 focus:ring-amber-500 focus:border-transparent"
		>
			<option value="all">All Roles</option>
			<option value="admin">Admins</option>
			<option value="manager">Managers</option>
			<option value="beekeeper">Beekeepers</option>
		</select>
	</div>

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
						Role
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Last Active
					</th>
					<th
						class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase"
					>
						Actions
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
				{#each filteredUsers as user}
					<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							{user.name}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							{user.email}
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="px-2 py-1 text-xs font-semibold rounded-full
                                {user.role === 'admin'
									? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300'
									: user.role === 'manager'
										? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
										: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'}"
							>
								{user.role}
							</span>
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							{new Date(user.lastActive).toLocaleString()}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							<button
								class="text-amber-600 hover:text-amber-700 mr-3"
								on:click={() => handleEdit(user)}
							>
								Edit
							</button>
							<button class="text-red-600 hover:text-red-700" on:click={() => handleDelete(user)}>
								Delete
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>

{#if editingUser}
	<UserEditModal
		isOpen={true}
		user={editingUser}
		onClose={() => (editingUser = null)}
		onSave={handleSaveUser}
	/>
{/if}

{#if deletingUser}
	<UserDeleteModal
		isOpen={true}
		user={deletingUser}
		onClose={() => (deletingUser = null)}
		onConfirm={handleDeleteUser}
	/>
{/if}

{#if showAddModal}
	<UserEditModal
		isOpen={true}
		user={null}
		onClose={() => (showAddModal = false)}
		onSave={handleAddUser}
	/>
{/if}
