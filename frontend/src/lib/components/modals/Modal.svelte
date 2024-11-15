<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import Icon from '@iconify/svelte';

	export let isOpen = false;
	export let title = '';
	export let onClose = () => {};
</script>

{#if isOpen}
	<div
		class="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 flex items-center justify-center p-4"
		in:fade={{ duration: 200 }}
		on:click={onClose}
	>
		<div
			class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-lg w-full border border-gray-100 dark:border-gray-700"
			in:scale={{ duration: 200, start: 0.95, opacity: 0 }}
			on:click|stopPropagation
		>
			<div
				class="flex justify-between items-center px-6 py-4 border-b border-gray-100 dark:border-gray-700"
			>
				<h2 class="text-xl font-semibold text-gray-900 dark:text-white tracking-tight">{title}</h2>
				<button
					class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200
                    transition-colors duration-200 rounded-lg p-1 hover:bg-gray-100 dark:hover:bg-gray-700"
					on:click={onClose}
				>
					<Icon icon="mdi:close" class="w-5 h-5" />
				</button>
			</div>
			<div class="px-6 py-5">
				<slot />
			</div>
		</div>
	</div>
{/if}
