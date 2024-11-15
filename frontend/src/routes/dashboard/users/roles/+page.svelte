<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import type { Role } from '$lib/types';
	import RoleEditModal from '$lib/components/modals/RoleEditModal.svelte';

	let editingRole: Role | null = null;

	// Mock data - replace with API calls
	const roles: Role[] = [
		{
			id: '1',
			name: 'admin',
			description: 'Full system access and control',
			permissions: ['manage_users', 'manage_roles', 'manage_system'],
			userCount: 2
		},
		{
			id: '2',
			name: 'manager',
			description: 'Manage beekeepers and view reports',
			permissions: ['view_reports', 'manage_beekeepers'],
			userCount: 5
		},
		{
			id: '3',
			name: 'beekeeper',
			description: 'Manage hives and record observations',
			permissions: ['manage_hives', 'record_observations'],
			userCount: 15
		}
	];

	function editRole(roleId: string) {
		const role = roles.find((r) => r.id === roleId);
		if (role) {
			editingRole = { ...role };
		}
	}

	function handleSaveRole(updatedRole: Role) {
		const index = roles.findIndex((r) => r.id === updatedRole.id);
		if (index !== -1) {
			roles[index] = updatedRole;
		}
		editingRole = null;
	}
</script>

<div class="max-w mx-auto space-y-8" in:fade>
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:shield-account" class="w-8 h-8 text-amber-500" />
				Role Management
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">Manage system roles and permissions</p>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-6">
		{#each roles as role}
			<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
				<div class="flex items-start justify-between">
					<div>
						<h2 class="text-xl font-semibold text-gray-900 dark:text-white capitalize">
							{role.name}
						</h2>
						<p class="mt-1 text-gray-600 dark:text-gray-400">{role.description}</p>
					</div>
					<button on:click={() => editRole(role.id)} class="text-amber-600 hover:text-amber-700">
						Edit Role
					</button>
				</div>

				<div class="mt-4">
					<h3 class="text-sm font-medium text-gray-900 dark:text-white mb-2">Permissions</h3>
					<div class="flex flex-wrap gap-2">
						{#each role.permissions as permission}
							<span
								class="px-2 py-1 text-xs font-semibold rounded-full
                                bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300"
							>
								{permission}
							</span>
						{/each}
					</div>
				</div>

				<div class="mt-4 flex items-center text-sm text-gray-600 dark:text-gray-400">
					<Icon icon="mdi:account-group" class="w-5 h-5 mr-1" />
					{role.userCount} users assigned
				</div>
			</div>
		{/each}
	</div>
</div>

{#if editingRole}
	<RoleEditModal
		isOpen={true}
		role={editingRole}
		onClose={() => (editingRole = null)}
		onSave={handleSaveRole}
	/>
{/if}
