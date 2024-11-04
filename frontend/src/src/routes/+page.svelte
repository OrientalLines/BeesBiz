<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { fly, fade, scale } from 'svelte/transition';
	import { onMount } from 'svelte';

	let mounted = false;
	let mouseX = 0;
	let mouseY = 0;
	let honeyScore = 0;
	let lastClickTime = 0;

	// Add audio element
	let beeSound: HTMLAudioElement;

	// Track mouse position
	function handleMouseMove(event: MouseEvent) {
		mouseX = event.clientX;
		mouseY = event.clientY;
	}

	// Generate initial random positions for bees
	const bees = Array(7)
		.fill(null)
		.map(() => {
			// Define safe zones (top, bottom, or sides)
			const zone = Math.floor(Math.random() * 4); // 0: top, 1: right, 2: bottom, 3: left
			let baseX, baseY;

			switch (zone) {
				case 0: // top
					baseX = Math.random() * 100;
					baseY = Math.random() * 30;
					break;
				case 1: // right
					baseX = 70 + Math.random() * 30;
					baseY = Math.random() * 100;
					break;
				case 2: // bottom
					baseX = Math.random() * 100;
					baseY = 70 + Math.random() * 30;
					break;
				case 3: // left
					baseX = Math.random() * 30;
					baseY = Math.random() * 100;
					break;
			}

			return {
				baseX,
				baseY,
				speedMultiplier: 0.5 + Math.random() * 0.5,
				clickable: true
			};
		});

	// Calculate bee position based on mouse
	function getBeeStyle(bee: { baseX: number; baseY: number; speedMultiplier: number }) {
		if (!mounted) return '';

		// Calculate distance between mouse and bee
		const beeX = (bee.baseX / 100) * window.innerWidth;
		const beeY = (bee.baseY / 100) * window.innerHeight;

		// Vector from mouse to bee
		const dx = beeX - mouseX;
		const dy = beeY - mouseY;

		// Calculate distance
		const distance = Math.sqrt(dx * dx + dy * dy);

		// Repulsion force (stronger when closer)
		const repulsionRadius = 200; // Radius of effect
		const repulsion = Math.max(0, 1 - distance / repulsionRadius);

		// Calculate offset
		const offsetX = dx * repulsion * 0.5;
		const offsetY = dy * repulsion * 0.5;

		return `
            top: calc(${bee.baseY}% + ${offsetY}px);
            left: calc(${bee.baseX}% + ${offsetX}px);
            transition: all 0.3s ease-out;
        `;
	}
	onMount(() => {
		mounted = true;
		// Initialize audio with a short, pleasant sound
		beeSound = new Audio('/orb.mp3');
		beeSound.volume = 0.01; // Reduce volume to 2%
	});

	function handleGetStarted() {
		if ($auth) {
			goto('/dashboard/hives');
		} else {
			goto('/login');
		}
	}

	const features = [
		{
			title: 'Hive Management',
			description: 'Track and monitor all your hives in one place',
			icon: '🐝'
		},
		{
			title: 'Harvest Tracking',
			description: 'Record and analyze honey production data',
			icon: '🍯'
		},
		{
			title: 'Reports & Analytics',
			description: 'Generate insights from your beekeeping data',
			icon: '📊'
		}
	];

	// Modify click handler to include sound
	function handleBeeClick(index: number) {
		const now = Date.now();
		if (now - lastClickTime < 500) return;
		lastClickTime = now;
		
		if (bees[index].clickable) {
			honeyScore += 1;
			bees[index].clickable = false;
			
			// Play sound with a slight variation in pitch
			if (beeSound) {
				const clone = beeSound.cloneNode() as HTMLAudioElement;
				clone.playbackRate = 0.9 + Math.random() * 0.2; // Random pitch between 0.9 and 1.1
				clone.play().catch(() => {}); // Ignore errors if sound can't play
			}
			
			setTimeout(() => {
				bees[index].clickable = true;
			}, 2000);
		}
	}
</script>

<div
	class="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-amber-50 via-blue-50 to-white relative overflow-hidden select-none"
	on:mousemove={handleMouseMove}
>
	<!-- Animated background elements -->
	{#if mounted}
		{#each bees as bee, i}
			<div
				class="absolute opacity-20 pointer-events-auto cursor-pointer"
				style={getBeeStyle(bee)}
				in:fly={{ y: 100, duration: 1000, delay: i * 200 }}
				on:click={() => handleBeeClick(i)}
			>
				<div 
					class="text-4xl animate-float-{i + 1} transition-all duration-300"
					style={bee.clickable ? '' : 'opacity: 0.5; transform: scale(0.8);'}
				>
					🐝
				</div>
			</div>
		{/each}
	{/if}

	<div class="max-w-4xl mx-auto text-center px-4 relative z-10">
		{#if mounted}
			<div in:fly={{ y: -50, duration: 1000 }}>
				<h1
					class="text-6xl font-bold text-gray-900 mb-6 bg-clip-text text-transparent bg-gradient-to-r from-amber-600 to-yellow-600"
				>
					Beekeeping Management System
				</h1>
			</div>

			<div in:fade={{ duration: 1000, delay: 500 }}>
				<p class="text-xl text-gray-600 mb-12">
					Efficiently manage your apiaries, track honey harvests, and monitor hive health with our
					comprehensive beekeeping management solution.
				</p>
			</div>

			<div in:scale={{ duration: 800, delay: 800 }}>
				<button
					on:click={handleGetStarted}
					class="bg-gradient-to-r from-amber-500 to-yellow-500 text-white px-10 py-4 rounded-full text-lg font-semibold hover:scale-105 transition-all duration-300 hover:shadow-lg"
				>
					Get Started
				</button>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-16">
				{#each features as feature, i}
					<div
						in:fly={{ x: -50, duration: 1000, delay: 1000 + i * 200 }}
						class="p-8 bg-white/80 backdrop-blur-sm rounded-xl shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-2 select-none"
					>
						<div class="text-4xl mb-4 pointer-events-none">{feature.icon}</div>
						<h3 class="text-xl font-semibold text-gray-900 mb-3">{feature.title}</h3>
						<p class="text-gray-600">{feature.description}</p>
					</div>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Add score display -->
	{#if mounted}
		<div 
			class="fixed top-4 right-4 bg-amber-500/90 backdrop-blur-sm text-white px-4 py-2 rounded-full font-semibold z-20"
			in:fly={{ y: -20, duration: 800 }}
		>
			🍯 {honeyScore}
		</div>
	{/if}

	<!-- Add game instructions -->
	<div 
		class="fixed bottom-4 left-4 bg-white/80 backdrop-blur-sm px-4 py-2 rounded-lg shadow-lg text-gray-600 text-sm z-20"
		in:fly={{ x: -20, duration: 800, delay: 1500 }}
	>
		Click the bees to collect honey! 🍯
	</div>
</div>

<style>
	@keyframes float-1 {
		0%,
		100% {
			transform: translate(0, 0) rotate(0deg);
		}
		50% {
			transform: translate(10px, -10px) rotate(5deg);
		}
	}
	@keyframes float-2 {
		0%,
		100% {
			transform: translate(0, 0) rotate(0deg);
		}
		50% {
			transform: translate(-8px, -12px) rotate(-3deg);
		}
	}
	@keyframes float-3 {
		0%,
		100% {
			transform: translate(0, 0) rotate(0deg);
		}
		50% {
			transform: translate(12px, -8px) rotate(8deg);
		}
	}
	@keyframes float-4 {
		0%,
		100% {
			transform: translate(0, 0) rotate(0deg);
		}
		50% {
			transform: translate(-10px, -10px) rotate(-8deg);
		}
	}
	@keyframes float-5 {
		0%,
		100% {
			transform: translate(0, 0) rotate(0deg);
		}
		50% {
			transform: translate(8px, -12px) rotate(3deg);
		}
	}

	.animate-float-1 {
		animation: float-1 8s infinite ease-in-out;
	}
	.animate-float-2 {
		animation: float-2 9s infinite ease-in-out;
	}
	.animate-float-3 {
		animation: float-3 10s infinite ease-in-out;
	}
	.animate-float-4 {
		animation: float-4 11s infinite ease-in-out;
	}
	.animate-float-5 {
		animation: float-5 12s infinite ease-in-out;
	}
</style>
