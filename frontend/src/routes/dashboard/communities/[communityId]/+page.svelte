<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { getToastStore } from '@skeletonlabs/skeleton';
	import {
		getBeeCommunities,
		getHives,
		updateBeeCommunity,
		deleteBeeCommunity
	} from '$lib/services/api';
	import type { BeeCommunity, Hive } from '$lib/types';
	import Modal from '$lib/components/modals/Modal.svelte';
	import { goto } from '$app/navigation';

	const communityId = Number($page.params.communityId) || 0;

	const toastStore = getToastStore();

	let loading = true;
	let error = '';
	let community: BeeCommunity | null = null;
	let hive: Hive | null = null;
	let showEditModal = false;
	let showDeleteModal = false;

	let formData = {
		community_id: 0,
		hive_id: 0,
		queen_age: 0,
		population_estimate: 0,
		health_status: 'healthy'
	};

	onMount(async () => {
		if (communityId === 0) {
			error = 'Invalid community ID';
			return;
		}
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			const [communities, hives] = await Promise.all([getBeeCommunities(), getHives()]);

			const foundCommunity = communities.find((c) => Number(c.community_id) === communityId);
			community = foundCommunity || null;
			if (community) {
				const foundHive = hives.find((h) => Number(h.hive_id) === Number(community?.hive_id));
				hive = foundHive || null;
				formData = { ...community };
			}

			if (!community) {
				throw new Error(`Community with ID ${communityId} not found`);
			}
			if (!hive) {
				throw new Error(`Hive with ID ${community?.hive_id} not found`);
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load bee community';
			toastStore.trigger({
				message: '❌ ' + error,
				background: 'variant-filled-error'
			});
		} finally {
			loading = false;
		}
	}

	function getHealthStatusColor(status: string): string {
		switch (status.toLowerCase()) {
			case 'healthy':
				return 'text-green-600 bg-green-100';
			case 'weak':
				return 'text-yellow-600 bg-yellow-100';
			case 'critical':
				return 'text-red-600 bg-red-100';
			default:
				return 'text-gray-600 bg-gray-100';
		}
	}

	async function handleUpdateCommunity(event: SubmitEvent) {
		event.preventDefault();
		try {
			await updateBeeCommunity(formData);
			showEditModal = false;
			await loadData();
			toastStore.trigger({
				message: '✨ Bee community updated successfully!',
				background: 'variant-filled-success'
			});
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to update bee community',
				background: 'variant-filled-error'
			});
		}
	}

	async function handleDelete() {
		try {
			await deleteBeeCommunity(communityId);
			toastStore.trigger({
				message: '✨ Bee community deleted successfully!',
				background: 'variant-filled-success'
			});
			goto('/dashboard/communities');
		} catch (e) {
			toastStore.trigger({
				message: '❌ Failed to delete bee community',
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
	{:else if community && hive}
		<!-- Header -->
		<div class="flex justify-between items-start">
			<div>
				<div class="flex items-center gap-3">
					<a
						href="/dashboard/communities"
						class="text-gray-500 hover:text-amber-500 transition-colors"
					>
						<Icon icon="mdi:arrow-left" class="w-6 h-6" />
					</a>
					<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Bee Community Details</h1>
				</div>
				<p class="mt-2 text-gray-600 dark:text-gray-400">
					Community ID: {community.community_id}
				</p>
			</div>
			<div class="flex gap-2">
				<button
					class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
                        hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors
                        flex items-center gap-2"
					on:click={() => (showEditModal = true)}
				>
					<Icon icon="mdi:pencil" class="w-5 h-5" />
					Edit
				</button>
				<button
					class="px-4 py-2 rounded-lg border border-red-200 text-red-600
                        hover:bg-red-50 transition-colors flex items-center gap-2"
					on:click={() => (showDeleteModal = true)}
				>
					<Icon icon="mdi:delete" class="w-5 h-5" />
					Delete
				</button>
			</div>
		</div>

		<!-- Community Info -->
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<!-- Community Details -->
			<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
				<h2 class="text-xl font-semibold mb-4 flex items-center gap-2">
					<Icon icon="mdi:bee" class="w-6 h-6 text-amber-500" />
					Community Status
				</h2>
				<div class="space-y-4">
					<div>
						<span class="text-gray-500 dark:text-gray-400">Health Status</span>
						<div
							class={`mt-1 px-3 py-1 rounded-full inline-block ${getHealthStatusColor(community.health_status)}`}
						>
							{community.health_status}
						</div>
					</div>
					<div>
						<span class="text-gray-500 dark:text-gray-400">Queen Age</span>
						<div class="text-lg font-semibold">
							{community.queen_age} months
						</div>
					</div>
					<div>
						<span class="text-gray-500 dark:text-gray-400">Population Estimate</span>
						<div class="text-lg font-semibold">
							{community.population_estimate.toLocaleString()} bees
						</div>
					</div>
				</div>
			</div>

			<!-- Hive Details -->
			<div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
				<h2 class="text-xl font-semibold mb-4 flex items-center gap-2">
					<Icon icon="mdi:home" class="w-6 h-6 text-amber-500" />
					Hive Information
				</h2>
				<div class="space-y-4">
					<div>
						<span class="text-gray-500 dark:text-gray-400">Hive Type</span>
						<div class="text-lg font-semibold">
							{hive.hive_type}
						</div>
					</div>
					<div>
						<span class="text-gray-500 dark:text-gray-400">Current Status</span>
						<div class="text-lg font-semibold">
							{hive.current_status}
						</div>
					</div>
					<div>
						<span class="text-gray-500 dark:text-gray-400">Installation Date</span>
						<div class="text-lg font-semibold">
							{hive.installation_date
								? new Date(hive.installation_date).toLocaleDateString()
								: 'Not recorded'}
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Edit Modal -->
		<Modal
			isOpen={showEditModal}
			title="Edit Bee Community"
			onClose={() => (showEditModal = false)}
		>
			<form on:submit|preventDefault={handleUpdateCommunity} class="space-y-4">
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
						on:click={() => (showEditModal = false)}
					>
						Cancel
					</button>
					<button
						type="submit"
						class="px-4 py-2 rounded-lg bg-amber-500 text-white hover:bg-amber-600"
					>
						Update
					</button>
				</div>
			</form>
		</Modal>

		<!-- Delete Confirmation Modal -->
		<Modal
			isOpen={showDeleteModal}
			title="Delete Bee Community"
			onClose={() => (showDeleteModal = false)}
		>
			<div class="space-y-4">
				<p class="text-gray-600 dark:text-gray-400">
					Are you sure you want to delete this bee community? This action cannot be undone.
				</p>
				<div class="flex justify-end gap-2">
					<button class="px-4 py-2 rounded-lg border" on:click={() => (showDeleteModal = false)}>
						Cancel
					</button>
					<button
						class="px-4 py-2 rounded-lg bg-red-500 text-white hover:bg-red-600"
						on:click={handleDelete}
					>
						Delete
					</button>
				</div>
			</div>
		</Modal>
	{/if}
</div>
