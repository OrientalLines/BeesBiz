<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { debounce } from 'lodash-es';
	import type { BeeCommunity, Hive } from '$lib/types';
	import { onMount } from 'svelte';
	import { createBeeCommunity, getBeeCommunities, getHives } from '$lib/services/api';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import Modal from '$lib/components/modals/Modal.svelte';

	const toastStore = getToastStore();
	let loading = true;
	let error = '';
	let communities: BeeCommunity[] = [];
	let hives: Hive[] = [];
	let hiveMap: Map<number, Hive> = new Map();

	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 12;
	let selectedHealthStatus = 'all';

	// Add modal state and form data
	let showModal = false;
	let formData = {
		hive_id: 0,
		queen_age: 0,
		population_estimate: 0,
		health_status: 'healthy'
	};

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			const [communitiesData, hivesData] = await Promise.all([getBeeCommunities(), getHives()]);

			communities = communitiesData || [];
			hives = hivesData || [];

			// Create a map of hive_id to hive for easy lookup
			hiveMap = new Map(hives.map((hive) => [hive.hive_id, hive]));
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load bee communities';
			toastStore.trigger({
				message: '❌ Failed to load bee communities',
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	function getHiveName(hiveId: number): string {
		const hive = hiveMap.get(hiveId);
		return hive ? `${hive.hive_type} Hive #${hive.hive_id}` : `Hive ${hiveId}`;
	}

	function getHealthStatusColor(status: string): string {
		switch (status.toLowerCase()) {
			case 'healthy':
				return 'text-green-600';
			case 'weak':
				return 'text-yellow-600';
			case 'critical':
				return 'text-red-600';
			default:
				return 'text-gray-600';
		}
	}

	const debouncedSearch = debounce((value: string) => {
		searchQuery = value;
		currentPage = 1;
	}, 300);

	$: filteredCommunities = communities
		.filter((c) => {
			const matchesSearch = getHiveName(c.hive_id)
				.toLowerCase()
				.includes(searchQuery.toLowerCase());
			const matchesStatus =
				selectedHealthStatus === 'all' ||
				c.health_status.toLowerCase() === selectedHealthStatus.toLowerCase();
			return matchesSearch && matchesStatus;
		})
		.sort((a, b) => b.community_id - a.community_id);

	$: paginatedCommunities = filteredCommunities.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredCommunities.length / itemsPerPage);

	async function handleCreateCommunity(event: SubmitEvent) {
		event.preventDefault();
		try {
			await createBeeCommunity(formData);
			showModal = false;
			await loadData();
			toastStore.trigger({
				message: '✨ Bee community created successfully!',
				background: 'variant-filled-success'
			});
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to create bee community',
				background: 'variant-filled-error'
			});
		}
	}
</script>

<div class="max-w mx-auto space-y-6" in:fade>
	{#if loading}
		<div class="flex justify-center items-center h-64">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500" />
		</div>
	{:else if error}
		<div class="text-center text-red-500 p-4">
			{error}
		</div>
	{:else}
		<div class="flex justify-between items-center">
			<div>
				<h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
					<Icon icon="mdi:bee" class="w-8 h-8 text-amber-500" />
					Bee Communities
				</h1>
				<p class="mt-2 text-gray-600 dark:text-gray-400">
					Monitor and manage bee populations across hives
				</p>
			</div>
			<button
				class="px-4 py-2 bg-amber-500 text-white rounded-lg hover:bg-amber-600 transition-colors
                    flex items-center gap-2"
				on:click={() => (showModal = true)}
			>
				<Icon icon="mdi:plus" class="w-5 h-5" />
				Add Community
			</button>
		</div>

		<!-- Search and Filter -->
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<div class="relative">
				<Icon
					icon="mdi:magnify"
					class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
				/>
				<input
					type="text"
					placeholder="Search communities..."
					class="w-full pl-10 pr-4 py-2 rounded-lg border"
					on:input={(e) => debouncedSearch(e.currentTarget.value)}
				/>
			</div>
			<div class="flex justify-end">
				<select
					bind:value={selectedHealthStatus}
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                        bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
				>
					<option value="all">All Statuses</option>
					<option value="healthy">Healthy</option>
					<option value="weak">Weak</option>
					<option value="critical">Critical</option>
				</select>
			</div>
		</div>

		<!-- Communities Grid -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each paginatedCommunities as community}
				<a
					href="/dashboard/communities/{community.community_id}"
					class="group relative p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg
                        hover:shadow-xl transition-all duration-300 text-left border-2 border-transparent
                        hover:border-amber-400"
				>
					<div class="flex items-center gap-3 mb-3">
						<Icon icon="mdi:bee" class="w-6 h-6 text-amber-500" />
						<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
							{getHiveName(community.hive_id)}
						</h2>
					</div>

					<div class="space-y-2">
						<div class="flex items-center gap-2">
							<Icon icon="mdi:heart-pulse" class="w-5 h-5" />
							<span class={getHealthStatusColor(community.health_status)}>
								Status: {community.health_status}
							</span>
						</div>
						<div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
							<Icon icon="mdi:crown" class="w-5 h-5" />
							<span>Queen Age: {community.queen_age} months</span>
						</div>
						<div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
							<Icon icon="mdi:bee-flower" class="w-5 h-5" />
							<span>Population: {community.population_estimate}</span>
						</div>
					</div>
				</a>
			{/each}
		</div>

		<!-- Empty State -->
		{#if communities.length === 0}
			<div class="text-center py-12">
				<Icon icon="mdi:bee" class="w-16 h-16 text-gray-400 mx-auto mb-4" />
				<h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
					No Bee Communities Found
				</h2>
				<p class="text-gray-600 dark:text-gray-400">
					Get started by adding your first bee community.
				</p>
			</div>
		{/if}

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="mt-8 flex justify-center gap-2">
				<button
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                        disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={currentPage === 1}
					on:click={() => currentPage--}
				>
					Previous
				</button>
				<span class="px-4 py-2">Page {currentPage} of {totalPages}</span>
				<button
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                        disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={currentPage === totalPages}
					on:click={() => currentPage++}
				>
					Next
				</button>
			</div>
		{/if}
	{/if}
</div>

<!-- Add Modal Component -->
<Modal isOpen={showModal} title="Create New Bee Community" onClose={() => (showModal = false)}>
	<form on:submit|preventDefault={handleCreateCommunity} class="space-y-4">
		<div>
			<label for="hive_id" class="block text-sm font-medium mb-1">Hive</label>
			<select
				id="hive_id"
				bind:value={formData.hive_id}
				class="w-full px-4 py-2 rounded-lg border"
				required
			>
				<option value={0} disabled>Select a hive...</option>
				{#each hives.filter((h) => !communities.some((c) => c.hive_id === h.hive_id)) as hive}
					<option value={hive.hive_id}>
						{hive.hive_type} Hive #{hive.hive_id} - {hive.current_status}
					</option>
				{/each}
			</select>
		</div>

		<div>
			<label for="queen_age" class="block text-sm font-medium mb-1">Queen Age (months)</label>
			<input
				type="number"
				id="queen_age"
				bind:value={formData.queen_age}
				min="0"
				max="60"
				class="w-full px-4 py-2 rounded-lg border"
				required
			/>
		</div>

		<div>
			<label for="population_estimate" class="block text-sm font-medium mb-1">
				Population Estimate
			</label>
			<input
				type="number"
				id="population_estimate"
				bind:value={formData.population_estimate}
				min="0"
				step="100"
				class="w-full px-4 py-2 rounded-lg border"
				required
			/>
		</div>

		<div>
			<label for="health_status" class="block text-sm font-medium mb-1">Health Status</label>
			<select
				id="health_status"
				bind:value={formData.health_status}
				class="w-full px-4 py-2 rounded-lg border"
				required
			>
				<option value="healthy">Healthy</option>
				<option value="weak">Weak</option>
				<option value="critical">Critical</option>
			</select>
		</div>

		<div class="flex justify-end gap-2">
			<button
				type="button"
				class="px-4 py-2 rounded-lg border"
				on:click={() => (showModal = false)}
			>
				Cancel
			</button>
			<button
				type="submit"
				class="px-4 py-2 rounded-lg bg-amber-500 text-white hover:bg-amber-600"
				disabled={formData.hive_id === 0}
			>
				Create
			</button>
		</div>
	</form>
</Modal>
