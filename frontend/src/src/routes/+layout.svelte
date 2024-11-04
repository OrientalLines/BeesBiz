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

	// Page title mapping
	const getTitleFromPath = (path: string) => {
		// Special cases first
		if (path === '/') return 'Welcome';
		if (path === '/login') return 'Login';
		if (path === '/404') return 'Page Not Found';

		// Dashboard routes
		const dashboardRoutes: Record<string, string> = {
			'/dashboard/hives': 'Hive Management',
			'/dashboard/harvest': 'Honey Harvest',
			'/dashboard/maintenance': 'Maintenance',
			'/dashboard/incidents': 'Incidents',
			'/dashboard/reports': 'Reports',
			'/dashboard/users': 'User Management',
			'/dashboard/settings': 'Settings'
		};

		return dashboardRoutes[path] || 'Dashboard';
	};

	// Update title when route changes
	$: {
		if (browser) {
			const pageTitle = getTitleFromPath($page.url.pathname);
			document.title = `BeesBiz | ${pageTitle}`;
		}
	}

	// Subscribe to theme changes and update localStorage and HTML class
	$: if (browser && $theme) {
		localStorage.setItem('theme', $theme);
		document.documentElement.classList.remove('light', 'dark');
		document.documentElement.classList.add($theme);
	}

	onMount(() => {
		// Initialize theme from localStorage if available
		const savedTheme = localStorage.getItem('theme');
		if (savedTheme) {
			theme.set(savedTheme as 'light' | 'dark');
		} else {
			// Set default theme based on system preference
			const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
			theme.set(prefersDark ? 'dark' : 'light');
		}

		// Listen for system theme changes
		window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
			if (!localStorage.getItem('theme')) {
				theme.set(e.matches ? 'dark' : 'light');
			}
		});
	});
</script>

<Toast />

<div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
	{#if $auth}
		<div class="flex" in:fade={{ duration: 200 }}>
			<Sidebar />
			<main class="flex-1 p-8 dark:text-gray-100">
				<slot />
			</main>
		</div>
	{:else}
		<slot />
	{/if}
</div>
