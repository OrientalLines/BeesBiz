<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { theme } from '$lib/stores/theme';
	import { getToastStore } from '@skeletonlabs/skeleton';

	const toastStore = getToastStore();
	let notificationsEnabled = true;
	let emailNotifications = true;
	let autoBackup = true;
	let language = 'en';

	const languages = [
		{ code: 'en', name: 'English', flag: '🇬🇧' },
		{ code: 'es', name: 'Español', flag: '🇪🇸' },
		{ code: 'fr', name: 'Français', flag: '🇫🇷' },
		{ code: 'de', name: 'Deutsch', flag: '🇩🇪' }
	];

	async function saveSettings() {
		// Simulated save with loading state
		const button = document.querySelector('#saveButton');
		if (button) button.classList.add('loading');

		await new Promise((resolve) => setTimeout(resolve, 1000));

		toastStore.trigger({
			message: '✨ Settings saved successfully!',
			background: 'variant-filled-success'
		});

		if (button) button.classList.remove('loading');
	}
</script>

<div class="max-w-4xl mx-auto space-y-8 p-6" in:fade={{ duration: 300 }}>
	<!-- Header -->
	<div class="relative">
		<div
			class="relative bg-white/50 dark:bg-gray-800/50 backdrop-blur-xl rounded-xl p-8 border border-gray-200/30 dark:border-gray-700/30 shadow-xl"
		>
			<h1 class="text-4xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
				<Icon icon="mdi:cog" class="w-10 h-10 text-amber-500" />
				System Settings
			</h1>
			<p class="mt-2 text-gray-600 dark:text-gray-400 max-w-2xl">
				Customize your experience with advanced settings and preferences
			</p>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-6">
		<!-- Appearance -->
		<div
			class="group bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700
            shadow-lg hover:shadow-xl transition-all duration-300 relative overflow-hidden"
		>
			<div class="relative">
				<h2
					class="text-xl font-semibold mb-6 flex items-center gap-2 text-gray-900 dark:text-white"
				>
					<Icon icon="mdi:palette" class="w-6 h-6 text-amber-500" />
					Appearance
				</h2>
				<div class="space-y-6">
					<!-- Theme Toggle -->
					<div
						class="flex items-center justify-between group/item p-4 rounded-lg
                         dark:hover:bg-gray-700/50 transition-colors"
					>
						<div class="flex items-start gap-3">
							<div class="p-2 rounded-lg bg-amber-100 dark:bg-amber-900/30">
								<Icon icon="mdi:theme-light-dark" class="w-5 h-5 text-amber-500" />
							</div>
							<div>
								<h3 class="font-medium text-gray-900 dark:text-white flex items-center gap-2">
									Theme
									<span
										class="text-xs px-2 py-1 rounded-full bg-gray-100 dark:bg-gray-700
                                        text-gray-600 dark:text-gray-400"
									>
										{$theme === 'dark' ? 'Dark' : 'Light'}
									</span>
								</h3>
								<p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
									Choose your preferred visual theme
								</p>
							</div>
						</div>
						<button
							class="p-3 rounded-xl bg-gray-100 dark:bg-gray-700 hover:bg-amber-100
                                dark:hover:bg-amber-900/30 transition-all duration-300 group-hover/item:scale-105"
							on:click={() => theme.toggle()}
						>
							<Icon
								icon={$theme === 'dark' ? 'mdi:weather-sunny' : 'mdi:weather-night'}
								class="w-6 h-6 text-amber-500"
							/>
						</button>
					</div>

					<!-- Language Selector -->
					<div
						class="flex items-center justify-between group/item p-4 rounded-lg dark:hover:bg-gray-700/50 transition-colors"
					>
						<div class="flex items-start gap-3">
							<div class="p-2 rounded-lg bg-amber-100 dark:bg-amber-900/30">
								<Icon icon="mdi:translate" class="w-5 h-5 text-amber-500" />
							</div>
							<div>
								<h3 class="font-medium text-gray-900 dark:text-white">Language</h3>
								<p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
									Select your preferred language
								</p>
							</div>
						</div>
						<select
							bind:value={language}
							class="px-6 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                                bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                                focus:ring-2 focus:ring-amber-500 focus:border-transparent
                                appearance-none cursor-pointer"
						>
							{#each languages as lang}
								<option value={lang.code}>{lang.flag} {lang.name}</option>
							{/each}
						</select>
					</div>
				</div>
			</div>
		</div>

		<!-- Notifications Panel -->
		<div
			class="group bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700
            shadow-lg hover:shadow-xl transition-all duration-300 relative overflow-hidden"
		>
			<div class="relative">
				<h2
					class="text-xl font-semibold mb-6 flex items-center gap-2 text-gray-900 dark:text-white"
				>
					<Icon icon="mdi:bell" class="w-6 h-6 text-amber-500" />
					Notifications
				</h2>
				<div class="space-y-6">
					<!-- Notification Toggles -->
					<div class="space-y-6">
						{#each [{ title: 'Push Notifications', description: 'Receive instant alerts about important events', icon: 'mdi:bell-ring', bind: notificationsEnabled }, { title: 'Email Notifications', description: 'Get updates delivered to your inbox', icon: 'mdi:email', bind: emailNotifications }] as setting}
							<div
								class="flex items-center justify-between group/item p-4 rounded-lg
                                dark:hover:bg-gray-700/50 transition-colors"
							>
								<div class="flex items-start gap-3">
									<div class="p-2 rounded-lg bg-amber-100 dark:bg-amber-900/30">
										<Icon icon={setting.icon} class="w-5 h-5 text-amber-500" />
									</div>
									<div>
										<h3 class="font-medium text-gray-900 dark:text-white">{setting.title}</h3>
										<p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
											{setting.description}
										</p>
									</div>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" bind:checked={setting.bind} class="sr-only peer" />
									<div
										class="w-14 h-7 bg-gray-200 peer-focus:outline-none peer-focus:ring-4
                                        peer-focus:ring-amber-300 dark:peer-focus:ring-amber-800 rounded-full peer
                                        dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white
                                        after:content-[''] after:absolute after:top-0.5 after:left-[4px] after:bg-white
                                        after:border-gray-300 after:border after:rounded-full after:h-6 after:w-6
                                        after:transition-all dark:border-gray-600 peer-checked:bg-amber-500"
									></div>
								</label>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</div>

		<!-- Save Button -->
		<div class="flex justify-end pt-4">
			<button
				id="saveButton"
				on:click={saveSettings}
				class="relative bg-gradient-to-r from-amber-500 to-orange-500 text-white px-8 py-3 rounded-xl
                    hover:from-amber-600 hover:to-orange-600 transition-all shadow-lg hover:shadow-xl
                    flex items-center gap-3"
			>
				<Icon icon="mdi:content-save" class="w-5 h-5" />
				<span>Save Changes</span>
			</button>
		</div>
	</div>
</div>

<style lang="postcss">
	.loading-gradient {
		background: linear-gradient(
			270deg,
			theme(colors.amber.500 / 10%),
			theme(colors.amber.500 / 20%),
			theme(colors.amber.500 / 20%),
			theme(colors.amber.500 / 10%)
		);
		background-size: 400% 100%;
		animation: loading 2s ease-in-out infinite;
	}

	.loading .loading-gradient {
		opacity: 1;
	}

	@keyframes loading {
		0% {
			background-position: 200% 0;
		}
		100% {
			background-position: -200% 0;
		}
	}
</style>
