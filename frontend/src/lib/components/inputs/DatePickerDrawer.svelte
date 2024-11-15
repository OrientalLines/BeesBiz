<script lang="ts">
	import { fly } from 'svelte/transition';
	import { Calendar } from 'lucide-svelte';
	import DatePicker from './DatePicker.svelte';

	export let isOpen = false;
	export let selectedDate: Date | null = null;
	export let onClose = () => {};
	export let onChange = (date: Date) => {};

	function handleDateSelect(date: Date) {
		onChange(date);
		onClose();
	}

	// Close on escape key
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			onClose();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isOpen}
	<!-- Backdrop -->
	<div
		class="fixed inset-0 bg-black/20 dark:bg-black/40 z-40"
		on:click={onClose}
		transition:fly={{ duration: 200, opacity: 0 }}
	/>

	<!-- Drawer -->
	<div
		class="fixed right-0 top-0 h-full w-full max-w-md bg-white dark:bg-gray-800 shadow-xl z-50 overflow-y-auto"
		transition:fly={{ x: 400, duration: 300 }}
	>
		<!-- Header -->
		<div
			class="p-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between"
		>
			<div class="flex items-center gap-2">
				<Calendar class="w-5 h-5 text-amber-500" />
				<h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Select Date</h2>
			</div>
			<button
				class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
				on:click={onClose}
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>
		</div>

		<!-- Calendar -->
		<div class="p-4">
			<DatePicker startDate={selectedDate} onChange={(start) => handleDateSelect(start as Date)} />
		</div>
	</div>
{/if}
