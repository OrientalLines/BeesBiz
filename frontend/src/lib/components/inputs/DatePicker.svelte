<script lang="ts">
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import { clickOutside } from '$lib/actions/clickOutside';

	export let startDate: Date | null = null;
	export let endDate: Date | null = null;
	export let placeholder = 'Select date range';
	export let onChange: (start: Date | null, end: Date | null) => void;

	let isOpen = false;
	let inputRef: HTMLDivElement;
	let currentMonth = new Date();
	let selectingEnd = false;

	const presets = [
		{ label: '🐝 Last 24 hours', value: '24h' },
		{ label: '🍯 Last 7 days', value: '7d' },
		{ label: '🐝 Last 30 days', value: '30d' },
		{ label: '🍯 Custom range', value: 'custom' }
	];

	const DAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
	const MONTHS = [
		'January',
		'February',
		'March',
		'April',
		'May',
		'June',
		'July',
		'August',
		'September',
		'October',
		'November',
		'December'
	];

	function getDaysInMonth(date: Date): Date[] {
		const year = date.getFullYear();
		const month = date.getMonth();
		const firstDay = new Date(year, month, 1);
		const lastDay = new Date(year, month + 1, 0);

		const days: Date[] = [];

		// Add previous month's days
		const firstDayOfWeek = firstDay.getDay();
		for (let i = firstDayOfWeek; i > 0; i--) {
			days.push(new Date(year, month, -i + 1));
		}

		// Add current month's days
		for (let i = 1; i <= lastDay.getDate(); i++) {
			days.push(new Date(year, month, i));
		}

		// Add next month's days
		const remainingDays = 42 - days.length; // 6 rows * 7 days = 42
		for (let i = 1; i <= remainingDays; i++) {
			days.push(new Date(year, month + 1, i));
		}

		return days;
	}

	$: calendarDays = getDaysInMonth(currentMonth);

	function isDateInRange(date: Date): boolean {
		if (!startDate || !endDate) return false;
		return date >= startDate && date <= endDate;
	}
	function isDateSelected(date: Date): boolean {
		return !!(
			(startDate && date.toDateString() === startDate.toDateString()) ||
			(endDate && date.toDateString() === endDate.toDateString())
		);
	}

	function handleDateClick(date: Date) {
		if (!selectingEnd && !startDate) {
			startDate = date;
			selectingEnd = true;
		} else if (selectingEnd) {
			if (date < startDate!) {
				endDate = startDate;
				startDate = date;
			} else {
				endDate = date;
			}
			selectingEnd = false;
			onChange(startDate, endDate);
		} else {
			startDate = date;
			endDate = null;
			selectingEnd = true;
		}
	}

	function formatDate(date: Date | null): string {
		if (!date) return '';
		return date.toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function getDisplayText(): string {
		if (!startDate && !endDate) return placeholder;
		if (startDate && !endDate) return formatDate(startDate);
		return `${formatDate(startDate)} - ${formatDate(endDate)}`;
	}

	function handlePresetClick(preset: string) {
		const now = new Date();
		let start: Date;

		switch (preset) {
			case '24h':
				start = new Date(now.getTime() - 24 * 60 * 60 * 1000);
				break;
			case '7d':
				start = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
				break;
			case '30d':
				start = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
				break;
			default:
				return;
		}

		startDate = start;
		endDate = now;
		onChange(startDate, endDate);
		isOpen = false;
	}

	function handleClear() {
		startDate = null;
		endDate = null;
		onChange(null, null);
		isOpen = false;
	}
</script>

<div class="relative" use:clickOutside={() => (isOpen = false)}>
	<div
		bind:this={inputRef}
		class="w-full px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700
        bg-white dark:bg-gray-800 text-gray-900 dark:text-white cursor-pointer
        hover:border-amber-400 transition-colors flex items-center justify-between"
		on:click={() => (isOpen = !isOpen)}
	>
		<div class="flex items-center gap-2">
			<Icon icon="mdi:calendar" class="w-5 h-5 text-amber-500" />
			<span class={!startDate && !endDate ? 'text-gray-500 dark:text-gray-400' : ''}>
				{getDisplayText()}
			</span>
		</div>
		<Icon icon={isOpen ? 'mdi:chevron-up' : 'mdi:chevron-down'} class="w-5 h-5 text-gray-400" />
	</div>

	{#if isOpen}
		<div
			class="absolute top-full left-0 mt-2 w-[340px] bg-white dark:bg-gray-800 rounded-xl
            shadow-lg border border-gray-200 dark:border-gray-700 p-4 z-50"
			in:fade={{ duration: 100 }}
		>
			<!-- Presets -->
			<div class="space-y-2 mb-4">
				{#each presets as preset}
					<button
						class="w-full px-4 py-2 text-left rounded-lg hover:bg-amber-50
                        dark:hover:bg-amber-900/20 transition-colors"
						on:click={() => handlePresetClick(preset.value)}
					>
						{preset.label}
					</button>
				{/each}
			</div>

			<!-- Custom Calendar -->
			<div class="border-t border-gray-200 dark:border-gray-700 pt-4">
				<!-- Calendar Header -->
				<div class="flex justify-between items-center mb-4">
					<button
						class="p-1 hover:bg-amber-50 dark:hover:bg-amber-900/20 rounded-full"
						on:click={() =>
							(currentMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1))}
					>
						<Icon icon="mdi:chevron-left" class="w-5 h-5" />
					</button>
					<h3 class="font-medium">
						{MONTHS[currentMonth.getMonth()]}
						{currentMonth.getFullYear()}
					</h3>
					<button
						class="p-1 hover:bg-amber-50 dark:hover:bg-amber-900/20 rounded-full"
						on:click={() =>
							(currentMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1))}
					>
						<Icon icon="mdi:chevron-right" class="w-5 h-5" />
					</button>
				</div>

				<!-- Calendar Grid -->
				<div class="grid grid-cols-7 gap-1">
					{#each DAYS as day}
						<div class="text-center text-xs font-medium text-gray-600 dark:text-gray-400 py-1">
							{day}
						</div>
					{/each}

					{#each calendarDays as date}
						<button
							class="aspect-square p-1 relative rounded-full
								{date.getMonth() === currentMonth.getMonth()
								? 'text-gray-800 dark:text-white'
								: 'text-gray-400 dark:text-gray-600'}
								{isDateSelected(date)
								? 'bg-amber-500 text-white dark:text-white hover:bg-amber-600 dark:hover:bg-amber-600'
								: isDateInRange(date)
									? 'bg-amber-100 text-gray-800 dark:bg-amber-900/20 dark:text-white hover:bg-amber-200 dark:hover:bg-amber-900/30'
									: 'hover:bg-amber-100/50 dark:hover:bg-amber-900/10'}
								transition-colors duration-200"
							on:click={() => handleDateClick(date)}
						>
							<span class="text-sm">
								{date.getDate()}
							</span>
							{#if isDateSelected(date)}
								<span class="absolute -top-1 -right-1 text-xs">🐝</span>
							{/if}
						</button>
					{/each}
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="mt-4 flex justify-end gap-2">
				<button
					class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400
                    hover:text-gray-900 dark:hover:text-white transition-colors"
					on:click={handleClear}
				>
					Clear
				</button>
				<button
					class="px-4 py-2 text-sm bg-amber-500 text-white rounded-lg
                    hover:bg-amber-600 transition-colors flex items-center gap-1"
					on:click={() => (isOpen = false)}
				>
					<span>Apply</span>
					<span>🍯</span>
				</button>
			</div>
		</div>
	{/if}
</div>
