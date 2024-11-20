<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { theme } from '$lib/stores/theme';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		AlertTriangle,
		Users,
		Settings,
		ChevronRight,
		Kanban,
		Pickaxe,
		PaintbrushVertical,
		FileText,
		LogOut,
		Sun,
		Moon,
		Bug,
		ClipboardList
	} from 'lucide-svelte';
	import { slide } from 'svelte/transition';
	import { onMount } from 'svelte';

	$: currentPath = $page.url.pathname;
	$: isDark = $theme === 'dark';
	let isCollapsed = false;

	const menuItems = [
		{
			label: 'Hive Management',
			icon: Kanban,
			roles: ['WORKER', 'ADMIN'],
			items: [
				{
					path: '/dashboard/hives',
					label: 'Hives Overview',
					roles: ['WORKER', 'ADMIN']
				},
				{
					path: '/dashboard/communities',
					label: 'Bee Communities',
					roles: ['WORKER', 'ADMIN']
				},
				{
					path: '/dashboard/iot',
					label: 'IoT Monitoring',
					roles: ['WORKER', 'ADMIN']
				}
			]
		},
		{
			label: 'Data Input',
			icon: ClipboardList,
			roles: ['WORKER', 'ADMIN'],
			items: [
				{
					path: '/dashboard/ingress/observations',
					label: 'Observations',
					roles: ['WORKER', 'ADMIN']
				},
				{
					path: '/dashboard/ingress/maintenance',
					label: 'Maintenance Plans',
					roles: ['WORKER', 'ADMIN']
				},
				{
					path: '/dashboard/ingress/incidents',
					label: 'Incidents',
					roles: ['WORKER', 'ADMIN']
				}
			]
		},
		{
			label: 'Data Output',
			icon: FileText,
			roles: ['WORKER', 'ADMIN'],
			items: [
				{
					path: '/dashboard/egress/harvests',
					label: 'Honey Harvest',
					roles: ['WORKER', 'ADMIN']
				},
				{
					path: '/dashboard/egress/reports',
					label: 'Reports',
					roles: ['WORKER', 'manager', 'ADMIN']
				}
			]
		},
		{
			label: 'User Management',
			icon: Users,
			roles: ['manager', 'ADMIN'],
			items: [
				{
					path: '/dashboard/users',
					label: 'Users Overview',
					roles: ['manager', 'ADMIN']
				},
				{
					path: '/dashboard/users/roles',
					label: 'Role Management',
					roles: ['ADMIN']
				},
				{
					path: '/dashboard/users/audit',
					label: 'Audit Log',
					roles: ['ADMIN']
				}
			]
		},
		{
			label: 'Settings',
			icon: Settings,
			roles: ['ADMIN'],
			items: [
				{
					path: '/dashboard/settings',
					label: 'General Settings',
					roles: ['ADMIN']
				}
				// {
				// 	path: '/dashboard/settings/notifications',
				// 	label: 'Notifications',
				// 	roles: ['ADMIN']
				// },
				// {
				// 	path: '/dashboard/settings/system',
				// 	label: 'System',
				// 	roles: ['ADMIN']
				// }
			]
		}
	];

	let expandedSection: string | null = null;

	function toggleSection(label: string) {
		expandedSection = expandedSection === label ? null : label;
	}
	function toggleTheme() {
		theme.toggle();
	}

	async function handleLogout() {
		await auth.logout();
		await goto('/');
	}
</script>

<aside
	class="bg-white dark:bg-gray-800 border-r border-amber-100 dark:border-amber-900/30 min-h-screen sticky top-0 shadow-sm transition-all duration-300 {isCollapsed
		? 'w-20'
		: 'w-64'}"
>
	<!-- Logo/Header Section -->
	<div class="p-4 border-b border-amber-100 dark:border-amber-900/30">
		<div class="flex items-center justify-between">
			{#if !isCollapsed}
				<h2
					class="font-semibold text-2xl bg-clip-text text-transparent bg-gradient-to-r from-amber-600 to-yellow-600"
				>
					BeesBiz
				</h2>
			{:else}
				<div class="w-5"></div>
			{/if}
			<button
				on:click={() => (isCollapsed = !isCollapsed)}
				class="p-2 rounded-lg hover:bg-amber-50 dark:hover:bg-amber-900/30 transition-colors duration-200"
			>
				<ChevronRight
					class="w-5 h-5 text-amber-600 dark:text-amber-400 transition-transform duration-300 {isCollapsed
						? 'rotate-180'
						: ''}"
				/>
			</button>
		</div>
	</div>

	<!-- Navigation Section -->
	<nav class="p-4">
		<ul class="space-y-2">
			{#each menuItems as item}
				{#if item.roles.includes($auth?.user?.role || '')}
					<li class="relative">
						{#if item.items}
							<!-- Section with subitems -->
							<button
								on:click={() => toggleSection(item.label)}
								class="w-full flex items-center justify-between px-4 py-3 rounded-lg transition-all duration-200
									text-gray-700 dark:text-gray-300 hover:bg-amber-50/50 dark:hover:bg-amber-900/20
									hover:text-amber-900 dark:hover:text-amber-200"
							>
								<div class="flex items-center">
									<svelte:component
										this={item.icon}
										class="w-5 h-5 text-gray-400 dark:text-gray-500"
									/>
									{#if !isCollapsed}
										<span class="ml-3">{item.label}</span>
									{/if}
								</div>
								{#if !isCollapsed}
									<ChevronRight
										class="w-4 h-4 transition-transform duration-200
										{expandedSection === item.label ? 'rotate-90' : ''}"
									/>
								{/if}
							</button>

							{#if expandedSection === item.label && !isCollapsed}
								<div class="relative">
									<ul
										class="mt-2 ml-4 space-y-2 overflow-hidden"
										transition:slide|local={{ duration: 200 }}
									>
										{#each item.items as subItem}
											{#if subItem.roles.includes($auth?.user?.role || '')}
												<li>
													<a
														href={subItem.path}
														class="flex items-center px-4 py-2 rounded-lg transition-all duration-200
															{currentPath === subItem.path
															? 'bg-amber-50 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300'
															: 'text-gray-600 dark:text-gray-400 hover:bg-amber-50/50 dark:hover:bg-amber-900/20'}"
													>
														<span class="ml-2">{subItem.label}</span>
													</a>
												</li>
											{/if}
										{/each}
									</ul>
								</div>
							{/if}
						{:else}
							<!-- Single item code remains unchanged -->
						{/if}
					</li>
				{/if}
			{/each}
		</ul>
	</nav>

	<!-- User Profile Section -->
	{#if $auth?.user}
		<div
			class="absolute bottom-0 left-0 right-0 p-4 border-t border-amber-100 dark:border-amber-900/30 bg-white dark:bg-gray-800"
		>
			<div class="flex items-center {isCollapsed ? 'justify-center' : 'gap-3'}">
				<div
					class="w-8 h-8 rounded-full bg-gradient-to-br from-amber-500 to-yellow-500 flex items-center justify-center"
				>
					<span class="text-white font-medium">
						{$auth.user.full_name ? $auth.user.full_name[0].toUpperCase() : '?'}
					</span>
				</div>
				{#if !isCollapsed}
					<div class="flex-1 min-w-0" transition:slide|local={{ duration: 200 }}>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-200 truncate">
							{$auth.user.full_name || 'User'}
						</p>
						<p class="text-xs text-amber-600 dark:text-amber-400 capitalize">
							{$auth.user.role || 'Guest'}
						</p>
					</div>
					<div class="flex items-center gap-2">
						<button
							on:click={toggleTheme}
							class="p-2 rounded-lg hover:bg-amber-50 dark:hover:bg-amber-900/30 transition-colors duration-200 text-amber-600 dark:text-amber-400"
						>
							{#if isDark}
								<Sun class="w-5 h-5" />
							{:else}
								<Moon class="w-5 h-5" />
							{/if}
						</button>
						<button
							on:click={handleLogout}
							class="p-2 rounded-lg hover:bg-amber-50 dark:hover:bg-amber-900/30 transition-colors duration-200 text-amber-600 dark:text-amber-400 hover:text-amber-800 dark:hover:text-amber-300"
						>
							<LogOut class="w-5 h-5" />
						</button>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</aside>
