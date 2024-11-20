<script lang="ts">
	import Modal from './Modal.svelte';
	import type { MaintenanceTask } from '$lib/types';
	import { fade, fly, scale } from 'svelte/transition';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import DatePicker from '$lib/components/inputs/DatePicker.svelte';
	import {
		CalendarClock,
		Home,
		ClipboardList,
		AlertCircle,
		CheckCircle2,
		Clock,
		ArrowUp,
		ArrowRight,
		ArrowDown,
		Settings
	} from 'lucide-svelte';

	export let isOpen = false;
	export let task: MaintenanceTask | null = null;
	export let onClose = () => {};
	export let onSave = (updatedTask: MaintenanceTask) => {};

	const toastStore = getToastStore();
	let isLoading = false;

	const priorities = [
		{
			value: 'High',
			icon: ArrowUp,
			color: 'text-red-500',
			bg: 'bg-red-50 dark:bg-red-900/10',
			description: 'Requires immediate attention'
		},
		{
			value: 'Medium',
			icon: ArrowRight,
			color: 'text-amber-500',
			bg: 'bg-amber-50 dark:bg-amber-900/10',
			description: 'Should be addressed soon'
		},
		{
			value: 'Low',
			icon: ArrowDown,
			color: 'text-blue-500',
			bg: 'bg-blue-50 dark:bg-blue-900/10',
			description: 'Can be scheduled later'
		}
	];

	let formData: Partial<MaintenanceTask> = {
		task: '',
		hiveId: '',
		dueDate: new Date(),
		status: 'Pending',
		priority: 'Medium'
	};

	$: if (task) {
		formData = { ...task };
	}

	function getTimeRemaining(dueDate: Date): string {
		if (!dueDate) return '';

		const now = new Date();
		now.setHours(0, 0, 0, 0);

		const due = new Date(dueDate);
		due.setHours(0, 0, 0, 0);

		const diffTime = due.getTime() - now.getTime();
		const days = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

		if (days < 0) {
			return `Overdue by ${Math.abs(days)} day${Math.abs(days) !== 1 ? 's' : ''}`;
		}
		if (days === 0) {
			return 'Due today';
		}
		if (days === 1) {
			return 'Due tomorrow';
		}
		return `${days} day${days !== 1 ? 's' : ''} remaining`;
	}

	async function handleSubmit() {
		if (!formData) return;

		isLoading = true;
		try {
			await onSave(formData as MaintenanceTask);
			toastStore.trigger({
				message: '✨ Task saved successfully!',
				background: 'variant-filled-success'
			});
			onClose();
		} catch (error) {
			toastStore.trigger({
				message: '❌ Failed to save task',
				background: 'variant-filled-error'
			});
		} finally {
			isLoading = false;
		}
	}

	function handleDateChange(date: Date | null) {
		if (date) {
			formData.dueDate = date;
		}
	}
</script>

<Modal {isOpen} title={task ? 'Edit Task' : 'New Task'} {onClose}>
	<div in:fly={{ y: 50, duration: 400 }} out:fade class="max-w-2xl mx-auto">
		<!-- Header Banner -->
		<div
			class="bg-amber-50 dark:bg-amber-900/10 p-4 rounded-lg flex items-start gap-3 mb-6"
			in:fly={{ y: 20, duration: 300, delay: 100 }}
		>
			<Settings class="w-6 h-6 text-amber-500 flex-shrink-0 mt-0.5" />
			<div>
				<h3 class="text-sm font-medium text-amber-800 dark:text-amber-200">
					{task ? 'Edit Maintenance Task' : 'Create New Maintenance Task'}
				</h3>
				<p class="mt-1 text-sm text-amber-700 dark:text-amber-300">
					Schedule and track maintenance tasks to keep your hives in optimal condition.
				</p>
			</div>
		</div>

		<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<!-- Hive ID -->
				<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 300 }}>
					<div class="flex items-center gap-2">
						<Home class="w-4 h-4 text-amber-500" />
						<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Hive ID </label>
					</div>
					<input
						type="text"
						bind:value={formData.hiveId}
						class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                            bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                            focus:ring-2 focus:ring-amber-500 focus:border-transparent
                            transition-all duration-300"
						required
						placeholder="Enter hive identifier"
					/>
				</div>

				<!-- Due Date -->
				<div class="space-y-2 relative" in:fly={{ y: 20, duration: 300, delay: 400 }}>
					<div class="flex items-center justify-between">
						<div class="flex items-center gap-2">
							<CalendarClock class="w-4 h-4 text-amber-500" />
							<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Due Date </label>
						</div>
						{#if formData.dueDate}
							<span
								class="text-sm px-2 py-1 rounded-full
                                    {new Date(formData.dueDate) < new Date()
									? 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
									: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'}"
								in:scale={{ duration: 200 }}
							>
								<Clock class="w-3 h-3 inline-block mr-1" />
								{getTimeRemaining(formData.dueDate)}
							</span>
						{/if}
					</div>
					<div class="relative">
						<DatePicker startDate={formData.dueDate} onChange={handleDateChange} />
					</div>
				</div>
			</div>

			<!-- Task Description -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 200 }}>
				<div class="flex items-center gap-2">
					<ClipboardList class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Task Description
					</label>
				</div>
				<textarea
					bind:value={formData.task}
					rows="3"
					class="w-full px-4 py-3 rounded-lg border-2 border-gray-200 dark:border-gray-700
                        bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                        focus:ring-2 focus:ring-amber-500 focus:border-transparent
                        transition-all duration-300 resize-none"
					required
					placeholder="Describe what needs to be done..."
				/>
			</div>

			<!-- Priority Selection -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 500 }}>
				<div class="flex items-center gap-2">
					<AlertCircle class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300">
						Priority Level
					</label>
				</div>
				<div class="grid grid-cols-1 md:grid-cols-3 gap-3">
					{#each priorities as priority}
						<button
							type="button"
							class="flex flex-col items-center p-4 border-2 rounded-lg transition-all duration-300
                                {formData.priority === priority.value
								? `border-${priority.color} ${priority.bg}`
								: 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'}"
							on:click={() => (formData.priority = priority.value)}
						>
							<svelte:component this={priority.icon} class="w-5 h-5 {priority.color} mb-2" />
							<span class="text-sm font-medium text-gray-900 dark:text-gray-100">
								{priority.value}
							</span>
							<span class="text-xs text-gray-500 dark:text-gray-400 text-center mt-1">
								{priority.description}
							</span>
						</button>
					{/each}
				</div>
			</div>

			<!-- Status -->
			<div class="space-y-2" in:fly={{ y: 20, duration: 300, delay: 600 }}>
				<div class="flex items-center gap-2">
					<CheckCircle2 class="w-4 h-4 text-amber-500" />
					<label class="text-sm font-medium text-gray-700 dark:text-gray-300"> Status </label>
				</div>
				<div class="flex gap-4">
					{#each ['Pending', 'Completed'] as status}
						<label
							class="flex items-center gap-2 p-3 cursor-pointer rounded-lg
                                {formData.status === status
								? 'bg-amber-50 dark:bg-amber-900/10'
								: ''}"
						>
							<input
								type="radio"
								bind:group={formData.status}
								value={status}
								class="w-4 h-4 text-amber-500 focus:ring-amber-500 border-gray-300"
							/>
							<span class="text-sm text-gray-700 dark:text-gray-300">
								{status}
							</span>
						</label>
					{/each}
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex justify-end gap-3 pt-6" in:fly={{ y: 20, duration: 300, delay: 700 }}>
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
                        disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 min-w-[120px] justify-center"
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
					{isLoading ? 'Saving...' : task ? 'Save Changes' : 'Create Task'}
				</button>
			</div>
		</form>
	</div>

	<Toast position="tr" />
</Modal>
