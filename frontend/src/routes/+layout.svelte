<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Toast from '$lib/components/Toast.svelte';
	import { auth } from '$lib/stores/auth';
	import { theme } from '$lib/stores/theme';
	import { fade } from 'svelte/transition';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { navigating } from '$app/stores';
	import { initializeStores } from '@skeletonlabs/skeleton';
	import { goto } from '$app/navigation';
	import process from 'process';

	const getTitleFromPath = (path: string) => {
		if (path === '/') return 'Welcome';
		if (path === '/login') return 'Login';
		if (path === '/404') return 'Page Not Found';

		const dashboardRoutes: Record<string, string> = {
			'/dashboard/hives': 'Hive Management',
			'/dashboard/harvest': 'Honey Harvest',
			'/dashboard/maintenance': 'Maintenance',
			'/dashboard/incidents': 'Incidents',
			'/dashboard/reports': 'Reports',
			'/dashboard/users': 'User Management',
			'/dashboard/users/roles': 'Role Management',
			'/dashboard/users/audit': 'Audit Log',
			'/dashboard/settings': 'Settings'
		};

		const matchingRoute = Object.entries(dashboardRoutes).find(([route]) => path.startsWith(route));

		return matchingRoute ? matchingRoute[1] : 'Dashboard';
	};

	$: {
		if (browser) {
			const pageTitle = getTitleFromPath($page.url.pathname);
			document.title = `BeesBiz | ${pageTitle}`;
		}
	}

	$: if (browser && $theme) {
		localStorage.setItem('theme', $theme);
		document.documentElement.classList.remove('light', 'dark');
		document.documentElement.classList.add($theme);
	}

	let pageLoaded = false;

	onMount(() => {
		if (browser) {
			auth.initialize();
			theme.initialize();
			// Check for stored token
			const storedAuth = localStorage.getItem('auth');
			if (storedAuth) {
				try {
					const authData = JSON.parse(storedAuth);
					auth.login(authData.user, authData.token);
				} catch (error) {
					console.error('Failed to parse stored auth data:', error);
					localStorage.removeItem('auth');
				}
			}
		}
		pageLoaded = true;
	});

	$: if (browser && $navigating) {
		const targetPath = $navigating.to?.url.pathname || '';
		if (isDashboardRoute(targetPath) && !hasRouteAccess(targetPath)) {
			goto('/dashboard/hives');
		}
	}

	initializeStores();

	function isDashboardRoute(path: string): boolean {
		return path.startsWith('/dashboard');
	}

	function hasRouteAccess(path: string): boolean {
		if (!$auth?.user) return false;

		const userRole = $auth.user.role.toUpperCase();

		const routeAccess: Record<string, string[]> = {
			'/dashboard/users': ['MANAGER', 'ADMIN'],
			'/dashboard/users/roles': ['ADMIN'],
			'/dashboard/users/audit': ['ADMIN'],
			'/dashboard/settings': ['ADMIN']
		};

		const matchingRoute = Object.entries(routeAccess).find(([route]) => path.startsWith(route));

		if (!matchingRoute) return true;
		return matchingRoute[1].includes(userRole);
	}

	if (browser) {
		window.process = process;
	}

	$: currentPath = $page.url.pathname;
</script>

<Toast />

<div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
	{#if !pageLoaded}
		<div class="h-screen flex flex-col items-center justify-center gap-8">
			<div class="relative">
				<div class="text-6xl bee-flight">🐝</div>
				<div
					class="absolute -bottom-8 left-1/2 -translate-x-1/2 w-8 h-2 bg-black/20 dark:bg-white/20 rounded-full bee-shadow"
				></div>
			</div>
			<p
				class="text-transparent bg-clip-text bg-gradient-to-r from-amber-500 via-yellow-400 to-amber-500 text-lg font-medium animate-pulse tracking-wider"
			>
				Buzzing into action...
			</p>
		</div>
	{:else if isDashboardRoute(currentPath)}
		{#if $auth?.user}
			<div class="flex" in:fade={{ duration: 200 }}>
				<Sidebar />
				<main class="flex-1 p-8 dark:text-gray-100 relative">
					{#if $navigating}
						<div
							class="absolute inset-0 bg-gray-50/50 dark:bg-gray-900/50 flex flex-col items-center justify-center gap-4 z-50"
						>
							<div class="relative">
								<div class="text-4xl bee-flight">🐝</div>
								<div
									class="absolute -bottom-6 left-1/2 -translate-x-1/2 w-6 h-1.5 bg-black/20 dark:bg-white/20 rounded-full bee-shadow"
								></div>
							</div>
						</div>
					{/if}
					<slot />
				</main>
			</div>
		{:else if browser}
			{goto('/login')}
		{/if}
	{:else}
		<slot />
	{/if}
</div>

<style>
	.bee-flight {
		animation: fly 2s infinite ease-in-out;
	}

	@keyframes fly {
		0%,
		100% {
			transform: translateY(0) rotate(10deg);
		}

		50% {
			transform: translateY(-20px) rotate(-10deg);
		}
	}

	.bee-shadow {
		animation: shadow 2s infinite ease-in-out;
	}

	@keyframes shadow {
		0%,
		100% {
			transform: scale(1);
			opacity: 0.4;
		}

		50% {
			transform: scale(0.8);
			opacity: 0.2;
		}
	}
</style>
