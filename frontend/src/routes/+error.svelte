<script lang="ts">
	import { fade } from 'svelte/transition';
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
</script>

<div
	class="min-h-[80vh] flex flex-col items-center justify-center p-8 text-center"
	in:fade={{ duration: 200 }}
>
	<div class="relative mb-8">
		<div class="text-8xl bee-flight">🐝</div>
		<div
			class="absolute -bottom-10 left-1/2 -translate-x-1/2 w-12 h-3 bg-black/20 dark:bg-white/20 rounded-full bee-shadow"
		></div>
	</div>

	<h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-4">
		{$page.status === 404 ? 'Oops! This hive seems empty' : 'Something went wrong'}
	</h1>

	<p class="text-lg text-gray-600 dark:text-gray-400 max-w-md mb-8">
		{$page.status === 404
			? "The page you're looking for has buzzed away or doesn't exist."
			: 'An unexpected error occurred while processing your request.'}
	</p>

	<button
		on:click={() => goto('/')}
		class="bg-amber-500 text-white px-6 py-3 rounded-full
        hover:bg-amber-600 transition-all shadow-lg hover:shadow-xl
        flex items-center gap-2"
	>
		<Icon icon="mdi:arrow-left" class="w-5 h-5" />
		Back to Home
	</button>
</div>

<style>
	.bee-flight {
		animation: fly 2s infinite ease-in-out;
	}

	@keyframes fly {
		0%,
		100% {
			transform: translateY(0) rotate(10deg);
		}
		50% {
			transform: translateY(-20px) rotate(-10deg);
		}
	}

	.bee-shadow {
		animation: shadow 2s infinite ease-in-out;
	}

	@keyframes shadow {
		0%,
		100% {
			transform: scale(1);
			opacity: 0.4;
		}
		50% {
			transform: scale(0.8);
			opacity: 0.2;
		}
	}
</style>
