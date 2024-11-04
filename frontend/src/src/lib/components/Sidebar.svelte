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
		Moon
	} from 'lucide-svelte';
	import { slide } from 'svelte/transition';
	import { onMount } from 'svelte';

	$: currentPath = $page.url.pathname;
	$: isDark = $theme === 'dark';
	let isCollapsed = false;

	const menuItems = [
		{
			path: '/dashboard/hives',
			label: 'Hive Management',
			roles: ['beekeeper', 'admin'],
			icon: Kanban
		},
		{
			path: '/dashboard/harvest',
			label: 'Honey Harvest',
			roles: ['beekeeper', 'admin'],
			icon: Pickaxe
		},
		{
			path: '/dashboard/maintenance',
			label: 'Maintenance',
			roles: ['beekeeper', 'admin'],
			icon: PaintbrushVertical
		},
		{
			path: '/dashboard/incidents',
			label: 'Incidents',
			roles: ['beekeeper', 'admin'],
			icon: AlertTriangle
		},
		{
			path: '/dashboard/reports',
			label: 'Reports',
			roles: ['beekeeper', 'manager', 'admin'],
			icon: FileText
		},
		{
			path: '/dashboard/users',
			label: 'User Management',
			roles: ['manager', 'admin'],
			icon: Users
		},
		{ path: '/dashboard/settings', label: 'Settings', roles: ['admin'], icon: Settings }
	];

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
				<h2 class="font-semibold text-2xl bg-clip-text text-transparent bg-gradient-to-r from-amber-600 to-yellow-600">BeesBiz</h2>
			{:else}
				<div class="w-5" />
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
				{#if item.roles.includes($auth?.role || '')}
					<li>
						<a
							href={item.path}
							class="flex items-center px-4 py-3 rounded-lg transition-all duration-200 group relative {currentPath ===
							item.path
								? 'bg-amber-50 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300 font-medium'
								: 'text-gray-700 dark:text-gray-300 hover:bg-amber-50/50 dark:hover:bg-amber-900/20 hover:text-amber-900 dark:hover:text-amber-200'}"
						>
							<svelte:component
								this={item.icon}
								class="w-5 h-5 {currentPath === item.path
									? 'text-amber-600 dark:text-amber-400'
									: 'text-gray-400 dark:text-gray-500 group-hover:text-amber-600 dark:group-hover:text-amber-400'}"
							/>

							{#if !isCollapsed}
								<span class="ml-3" transition:slide|local={{ duration: 200 }}>
									{item.label}
								</span>
							{:else}
								<div
									class="absolute left-full ml-2 px-2 py-1 bg-amber-600 dark:bg-amber-700 text-white text-sm rounded-md opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap z-50"
								>
									{item.label}
								</div>
							{/if}

							{#if currentPath === item.path}
								<div
									class="absolute left-0 w-1 h-8 bg-gradient-to-b from-amber-500 to-yellow-500 rounded-r-full"
									transition:slide|local={{ duration: 200 }}
								/>
							{/if}
						</a>
					</li>
				{/if}
			{/each}
		</ul>
	</nav>

	<!-- User Profile Section -->
	{#if $auth}
		<div
			class="absolute bottom-0 left-0 right-0 p-4 border-t border-amber-100 dark:border-amber-900/30 bg-white dark:bg-gray-800"
		>
			<div class="flex items-center {isCollapsed ? 'justify-center' : 'gap-3'}">
				<div
					class="w-8 h-8 rounded-full bg-gradient-to-br from-amber-500 to-yellow-500 flex items-center justify-center"
				>
					<span class="text-white font-medium">{$auth.name[0].toUpperCase()}</span>
				</div>
				{#if !isCollapsed}
					<div class="flex-1" transition:slide|local={{ duration: 200 }}>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-200">{$auth.name}</p>
						<p class="text-xs text-amber-600 dark:text-amber-400 capitalize">{$auth.role}</p>
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
