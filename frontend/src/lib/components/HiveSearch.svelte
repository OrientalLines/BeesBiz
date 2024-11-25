<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { searchHives } from '$lib/services/api';
	import Icon from '@iconify/svelte';
	import { clickOutside } from '$lib/actions/clickOutside';
	import type { Hive } from '$lib/types';
	import Fuse from 'fuse.js';

	export let selectedHiveId: string = '';
	export let placeholder = 'Search for a hive...';

	const dispatch = createEventDispatcher<{
		select: { hiveId: string };
	}>();

	let searchTerm = '';
	let results: Hive[] = [];
	let loading = false;
	let showDropdown = false;
	let searchInput: HTMLInputElement;

	let debounceTimer: NodeJS.Timeout;

	// Configure Fuse.js options
	const fuseOptions = {
		keys: ['hive_id'],
		threshold: 0.3,
		includeScore: true,
		shouldSort: true
	};

	let allHives: Hive[] = [];

	// Load all hives when component mounts
	async function loadAllHives() {
		try {
			loading = true;
			allHives = await searchHives(''); // Assuming empty string returns all hives
		} catch (error) {
			console.error('Failed to load hives:', error);
		} finally {
			loading = false;
		}
	}

	async function handleSearch() {
		clearTimeout(debounceTimer);

		if (!searchTerm.trim()) {
			results = allHives;
			return;
		}

		debounceTimer = setTimeout(() => {
			const fuse = new Fuse(allHives, fuseOptions);
			const fuzzyResults = fuse.search(searchTerm);
			results = fuzzyResults
				.sort((a, b) => (a.score || 0) - (b.score || 0))
				.map((result) => result.item);
		}, 300);
	}

	// Show all hives when focusing on empty input
	function handleFocus() {
		if (!searchTerm) {
			results = allHives;
		}
		showDropdown = true;
	}

	// Load hives when component mounts
	onMount(() => {
		loadAllHives();
	});

	function handleSelect(hive: Hive) {
		selectedHiveId = hive.hive_id.toString();
		// Use optional chaining and type assertion for hive_name
		searchTerm = `Hive ${hive.hive_id}${(hive as any).hive_name ? ` - ${(hive as any).hive_name}` : ''}`;
		showDropdown = false;
		dispatch('select', { hiveId: hive.hive_id.toString() });
	}

	function handleClickOutside() {
		showDropdown = false;
	}

	function updateDropdownOffset() {
		if (searchInput) {
			const rect = searchInput.getBoundingClientRect();
			document.documentElement.style.setProperty('--dropdown-offset', `${rect.bottom + 16}px`);
		}
	}

	$: if (showDropdown) {
		updateDropdownOffset();
	}
</script>

<div class="relative w-full" use:clickOutside={handleClickOutside}>
	<div class="relative">
		<input
			bind:this={searchInput}
			type="text"
			bind:value={searchTerm}
			on:input={handleSearch}
			on:focus={handleFocus}
			{placeholder}
			class="w-full px-4 py-2 text-sm border rounded-md
                border-gray-300 dark:border-gray-600
                focus:outline-none focus:ring-2 focus:ring-amber-500
                dark:bg-gray-700 dark:text-white
                pr-10"
		/>
		{#if loading}
			<div class="absolute right-3 top-2.5">
				<div class="animate-spin rounded-full h-5 w-5 border-b-2 border-amber-500" />
			</div>
		{:else if selectedHiveId}
			<button
				class="absolute right-3 top-2.5 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
				on:click={() => {
					selectedHiveId = '';
					searchTerm = '';
					dispatch('select', { hiveId: '' });
				}}
			>
				<Icon icon="mdi:close" class="w-5 h-5" />
			</button>
		{/if}
	</div>

	{#if showDropdown && (results.length > 0 || loading)}
		<div
			class="absolute top-full mt-1 z-50 w-full bg-white dark:bg-gray-800
                rounded-md shadow-lg border border-gray-200 dark:border-gray-700
                overflow-y-auto"
			style="max-height: calc(100vh - var(--dropdown-offset, 200px));"
		>
			{#if loading}
				<div class="p-4 text-center text-gray-500">Loading hives...</div>
			{:else}
				{#each results as hive}
					<button
						on:click={() => handleSelect(hive)}
						class="w-full px-4 py-2 text-left text-sm hover:bg-gray-100
                            dark:hover:bg-gray-700 flex items-center justify-between
                            transition-colors duration-150"
					>
						<div class="flex flex-col">
							<span class="font-medium">Hive {hive.hive_id}</span>
							{#if (hive as any).hive_name || (hive as any).location}
								<span class="text-xs text-gray-500 dark:text-gray-400">
									{[(hive as any).hive_name, (hive as any).location].filter(Boolean).join(' - ')}
								</span>
							{/if}
						</div>
						<Icon icon="mdi:arrow-right" class="w-4 h-4 text-amber-500 flex-shrink-0 ml-2" />
					</button>
				{/each}
			{/if}
		</div>
	{:else if showDropdown && searchTerm && !loading}
		<div
			class="absolute top-full mt-1 z-50 w-full bg-white dark:bg-gray-800
			rounded-md shadow-lg border border-gray-200 dark:border-gray-700
			p-4"
		>
			<p class="text-sm text-gray-500 dark:text-gray-400 text-center">No hives found</p>
		</div>
	{/if}
</div>
