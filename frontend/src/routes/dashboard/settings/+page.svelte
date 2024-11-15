<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { theme } from '$lib/stores/theme';

	let notificationsEnabled = true;
	let emailNotifications = true;
	let autoBackup = true;
	let language = 'en';

	const languages = [
		{ code: 'en', name: 'English' },
		{ code: 'es', name: 'Español' },
		{ code: 'fr', name: 'Français' },
		{ code: 'de', name: 'Deutsch' }
	];

	function saveSettings() {
		// Implementation for saving settings
	}
</script>

<div class="max-w mx-auto space-y-8" in:fade>
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:cog" class="w-8 h-8 text-amber-500" />
				System Settings
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400">
				Configure system preferences and notifications
			</p>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-8">
		<!-- Appearance -->
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
			<h2 class="text-xl font-semibold mb-4 flex items-center gap-2">
				<Icon icon="mdi:palette" class="w-6 h-6 text-amber-500" />
				Appearance
			</h2>
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<div>
						<h3 class="font-medium text-gray-900 dark:text-white">Theme</h3>
						<p class="text-sm text-gray-600 dark:text-gray-400">Choose your preferred theme</p>
					</div>
					<button
						class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
						on:click={() => theme.toggle()}
					>
						<Icon
							icon={$theme === 'dark' ? 'mdi:weather-sunny' : 'mdi:weather-night'}
							class="w-6 h-6 text-amber-500"
						/>
					</button>
				</div>

				<div class="flex items-center justify-between">
					<div>
						<h3 class="font-medium text-gray-900 dark:text-white">Language</h3>
						<p class="text-sm text-gray-600 dark:text-gray-400">Select your preferred language</p>
					</div>
					<select
						bind:value={language}
						class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
							bg-white dark:bg-gray-800 text-gray-900 dark:text-white
							focus:ring-2 focus:ring-amber-500 focus:border-transparent"
					>
						{#each languages as lang}
							<option value={lang.code}>{lang.name}</option>
						{/each}
					</select>
				</div>
			</div>
		</div>

		<!-- Notifications -->
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
			<h2 class="text-xl font-semibold mb-4 flex items-center gap-2">
				<Icon icon="mdi:bell" class="w-6 h-6 text-amber-500" />
				Notifications
			</h2>
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<div>
						<h3 class="font-medium text-gray-900 dark:text-white">Enable Notifications</h3>
						<p class="text-sm text-gray-600 dark:text-gray-400">
							Receive alerts about important events
						</p>
					</div>
					<label class="relative inline-flex items-center cursor-pointer">
						<input type="checkbox" bind:checked={notificationsEnabled} class="sr-only peer" />
						<div
							class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4
							peer-focus:ring-amber-300 dark:peer-focus:ring-amber-800 rounded-full peer
							dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white
							after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white
							after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5
							after:transition-all dark:border-gray-600 peer-checked:bg-amber-500"
						></div>
					</label>
				</div>

				<div class="flex items-center justify-between">
					<div>
						<h3 class="font-medium text-gray-900 dark:text-white">Email Notifications</h3>
						<p class="text-sm text-gray-600 dark:text-gray-400">Receive notifications via email</p>
					</div>
					<label class="relative inline-flex items-center cursor-pointer">
						<input type="checkbox" bind:checked={emailNotifications} class="sr-only peer" />
						<div
							class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4
							peer-focus:ring-amber-300 dark:peer-focus:ring-amber-800 rounded-full peer
							dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white
							after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white
							after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5
							after:transition-all dark:border-gray-600 peer-checked:bg-amber-500"
						></div>
					</label>
				</div>
			</div>
		</div>

		<!-- System -->
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
			<h2 class="text-xl font-semibold mb-4 flex items-center gap-2">
				<Icon icon="mdi:database" class="w-6 h-6 text-amber-500" />
				System
			</h2>
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<div>
						<h3 class="font-medium text-gray-900 dark:text-white">Automatic Backup</h3>
						<p class="text-sm text-gray-600 dark:text-gray-400">Enable daily automatic backups</p>
					</div>
					<label class="relative inline-flex items-center cursor-pointer">
						<input type="checkbox" bind:checked={autoBackup} class="sr-only peer" />
						<div
							class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4
							peer-focus:ring-amber-300 dark:peer-focus:ring-amber-800 rounded-full peer
							dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white
							after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white
							after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5
							after:transition-all dark:border-gray-600 peer-checked:bg-amber-500"
						></div>
					</label>
				</div>
			</div>
		</div>
	</div>

	<div class="flex justify-end">
		<button
			on:click={saveSettings}
			class="bg-amber-500 text-white px-6 py-3 rounded-full
				hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
				flex items-center gap-2"
		>
			<Icon icon="mdi:content-save" class="w-5 h-5" />
			Save Settings
		</button>
	</div>
</div>
