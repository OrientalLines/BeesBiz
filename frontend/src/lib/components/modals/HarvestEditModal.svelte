<script lang="ts">
	import Modal from './Modal.svelte';
	import type { HoneyHarvest } from '$lib/types';
	import { fade, fly } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import { Droplets, Home, Award, CalendarDays } from 'lucide-svelte';
	import { getHives } from '$lib/services/api';
	import { onMount } from 'svelte';
	import { debounce } from 'lodash-es';
	import type { Hive } from '$lib/types';

	export let isOpen = false;
	export let harvest: HoneyHarvest | null = null;
	export let onClose = () => {};
	export let onSave = (updatedHarvest: HoneyHarvest) => {};

	const toastStore = getToastStore();
	let isLoading = false;
	let formData: Partial<HoneyHarvest> = {
		harvest_id: undefined,
		harvest_date: new Date(),
		hive_id: 0,
		quantity: 0,
		quality_grade: 'A'
	};

	const qualityGrades = [
		{ value: 'A', label: 'Premium Grade A', description: 'Highest quality honey' },
		{ value: 'B', label: 'Grade B', description: 'Good quality honey' },
		{ value: 'C', label: 'Grade C', description: 'Standard quality honey' },
		{ value: 'D', label: 'Grade D', description: 'Processing grade honey' }
	];

	let hives: Hive[] = [];
	let filteredHives: Hive[] = [];
	let hiveSearchQuery = '';
	let showHiveDropdown = false;
	let selectedHive: Hive | null = null;

	onMount(async () => {
		await loadHives();
		if (harvest) {
			formData = {
				...harvest,
				harvest_date: harvest.harvest_date ? new Date(harvest.harvest_date) : new Date()
			};
			const matchingHive = hives.find((h) => h.hive_id === harvest.hive_id);
			if (matchingHive) {
				selectedHive = matchingHive;
				hiveSearchQuery = `Hive #${matchingHive.hive_id} (${matchingHive.hive_type})`;
			}
		}
	});

	async function loadHives() {
		try {
			isLoading = true;
			hives = await getHives();
			filterHives();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to load hives',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	const filterHives = debounce(() => {
		if (!hiveSearchQuery.trim()) {
			filteredHives = hives;
		} else {
			filteredHives = hives.filter(
				(hive) =>
					hive.hive_id.toString().includes(hiveSearchQuery) ||
					(hive.hive_type && hive.hive_type.toLowerCase().includes(hiveSearchQuery.toLowerCase()))
			);
		}
	}, 300);

	$: if (hiveSearchQuery !== undefined) {
		filterHives();
	}

	function handleHiveInput(event: Event) {
		const input = event.target as HTMLInputElement;
		hiveSearchQuery = input.value;
		showHiveDropdown = true;
		selectedHive = null;
		formData.hive_id = 0;
	}

	function selectHive(hive: Hive) {
		selectedHive = hive;
		formData.hive_id = hive.hive_id;
		hiveSearchQuery = `Hive #${hive.hive_id} (${hive.hive_type})`;
		showHiveDropdown = false;
	}

	async function handleSubmit() {
		if (!formData.hive_id) return;

		isLoading = true;
		try {
			const submitData: HoneyHarvest = {
				...formData,
				harvest_id: harvest?.harvest_id // Do not set to 0; undefined indicates creation
			} as HoneyHarvest;

			await onSave(submitData);
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to save harvest record',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	function handleDateChange(date: Date | null) {
		if (date) {
			formData.harvest_date = date;
		}
	}

	$: if (harvest && isOpen) {
		formData = {
			...harvest,
			harvest_date: harvest.harvest_date ? new Date(harvest.harvest_date) : new Date()
		};
		const matchingHive = hives.find((h) => h.hive_id === harvest.hive_id);
		if (matchingHive) {
			selectedHive = matchingHive;
			hiveSearchQuery = `Hive #${matchingHive.hive_id} (${matchingHive.hive_type})`;
		}
	} else if (!isOpen) {
		formData = {
			harvest_id: undefined,
			harvest_date: new Date(),
			hive_id: 0,
			quantity: 0,
			quality_grade: 'A'
		};
		selectedHive = null;
		hiveSearchQuery = '';
	}
</script>

<Modal {isOpen} title={harvest ? 'Edit Harvest Record' : 'New Harvest Record'} {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade>
		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<!-- Harvest Date -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 100 }}>
				<div class="flex items-center gap-2 mb-2">
					<CalendarDays class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Harvest Date </label>
				</div>
				<DatePicker startDate={new Date(formData.harvest_date ?? '')} onChange={handleDateChange} />
			</div>

			<!-- Hive Selection -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 200 }}>
				<div class="flex items-center gap-2 mb-2">
					<Home class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Hive Selection
					</label>
				</div>
				<div class="relative">
					<input
						type="text"
						placeholder="Search by hive ID or type..."
						bind:value={hiveSearchQuery}
						on:input={handleHiveInput}
						on:focus={() => (showHiveDropdown = true)}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
							bg-white dark:bg-gray-800 text-gray-900 dark:text-white
							focus:ring-2 focus:ring-amber-500 focus:border-transparent
							transition-all duration-300"
						required
					/>

					{#if showHiveDropdown && filteredHives.length > 0}
						<div
							class="absolute z-50 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-200
							dark:border-gray-700 rounded-lg shadow-lg max-h-48 overflow-y-auto dropdown-scroll"
						>
							{#each filteredHives as hive}
								<button
									type="button"
									class="w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700
										text-gray-900 dark:text-white"
									on:click={() => selectHive(hive)}
								>
									Hive #{hive.hive_id} ({hive.hive_type}) - {hive.current_status}
								</button>
							{/each}
						</div>
					{/if}
				</div>
				{#if selectedHive}
					<p class="text-sm text-gray-600 dark:text-gray-400">
						Selected: Hive #{selectedHive.hive_id} ({selectedHive.hive_type}) - {selectedHive.current_status}
					</p>
				{/if}
			</div>

			<!-- Quantity -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 300 }}>
				<div class="flex items-center gap-2 mb-2">
					<Droplets class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Quantity </label>
				</div>
				<div class="relative">
					<input
						type="number"
						step="0.1"
						bind:value={formData.quantity}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
							bg-white dark:bg-gray-800 text-gray-900 dark:text-white
							focus:ring-2 focus:ring-amber-500 focus:border-transparent
							transition-all duration-300"
						required
					/>
					<span class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500">kg</span>
				</div>
			</div>

			<!-- Quality Grade -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 400 }}>
				<div class="flex items-center gap-2 mb-2">
					<Award class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Quality Grade
					</label>
				</div>
				<div class="grid grid-cols-2 gap-3">
					{#each qualityGrades as grade}
						<button
							type="button"
							class="flex flex-col items-center p-3 border-2 rounded-lg transition-all duration-300
								{formData.quality_grade === grade.value
								? 'border-amber-500 bg-amber-50 dark:bg-amber-900/10'
								: 'border-gray-200 dark:border-gray-700 hover:border-amber-500/50'}"
							on:click={() => {
								formData.quality_grade = grade.value as 'A' | 'B' | 'C' | 'D';
							}}
						>
							<span class="text-lg font-semibold text-amber-500">
								Grade {grade.value}
							</span>
							<span class="text-xs text-gray-500 dark:text-gray-400 text-center mt-1">
								{grade.description}
							</span>
						</button>
					{/each}
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 pt-6" in:fly={{ y: 20, duration: 300, delay: 500 }}>
				<button
					type="button"
					class="px-6 py-2.5 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700
						rounded-lg transition-all duration-300 transform hover:scale-105"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-6 py-2.5 bg-gradient-to-r from-amber-500 to-amber-600 text-white rounded-lg
						hover:from-amber-600 hover:to-amber-700 transition-all duration-300 transform hover:scale-105
						disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
					disabled={isLoading}
				>
					{#if isLoading}
						<svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
								fill="none"
							/>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							/>
						</svg>
					{/if}
					{isLoading ? 'Saving...' : harvest ? 'Save Changes' : 'Record Harvest'}
				</button>
			</div>
		</form>
	</div>

	<Toast position="tr" />
</Modal>

<style lang="postcss">
	/* Стилизация скроллбара для современного вида */
	::-webkit-scrollbar {
		width: 8px;
	}

	::-webkit-scrollbar-track {
		background: theme(colors.gray.100);
		border-radius: 4px;
	}

	::-webkit-scrollbar-thumb {
		background: theme(colors.amber.500);
		border-radius: 4px;
	}

	::-webkit-scrollbar-thumb:hover {
		background: theme(colors.amber.600);
	}

	/* Для Firefox */
	* {
		scrollbar-width: thin;
		scrollbar-color: theme(colors.amber.500) theme(colors.gray.100);
	}

	/* Add styles for the dropdown */
	.dropdown-scroll {
		scrollbar-width: thin;
		scrollbar-color: theme(colors.amber.500) theme(colors.gray.100);
	}
</style>
