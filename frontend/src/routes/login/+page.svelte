<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { fly, fade } from 'svelte/transition';
	import { login, register } from '$lib/services/api';

	let activeTab: 'login' | 'register' = 'login';
	let email = '';
	let password = '';
	let username = '';
	let fullName = '';
	let error = '';

	let mounted = false;
	onMount(() => {
		mounted = true;
		if ($auth?.user) {
			goto('/dashboard/hives');
		}
	});

	async function handleLogin() {
		try {
			error = '';
			if (!email || !password) {
				error = 'Please fill in all fields';
				return;
			}

			const response = await login(email, password);
			// @ts-ignore fuck typescript
			auth.login(response.user, response.token);
			goto('/dashboard/hives');
		} catch (err) {
			error = 'Invalid credentials';
		}
	}

	async function handleRegister() {
		try {
			error = '';
			if (!email || !password || !username || !fullName) {
				error = 'Please fill in all fields';
				return;
			}

			await register({ email, password, username, fullName });
			// Switch to login tab after successful registration
			activeTab = 'login';
			error = '';
			// Clear form
			email = '';
			password = '';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Registration failed';
		}
	}

	const bees = Array(5)
		.fill(null)
		.map(() => ({
			baseX: Math.random() * 100,
			baseY: Math.random() * 100
		}));
</script>

<div
	class="min-h-screen flex items-center justify-center bg-gradient-to-br from-amber-50 via-blue-50 to-white relative overflow-hidden"
>
	<!-- Animated background bees -->
	{#if mounted}
		{#each bees as bee, i}
			<div
				class="absolute opacity-20 pointer-events-none"
				style="top: {bee.baseY}%; left: {bee.baseX}%;"
				in:fly={{ y: 100, duration: 1000, delay: i * 200 }}
			>
				<div class="text-4xl animate-float-{i + 1}">🐝</div>
			</div>
		{/each}
	{/if}

	<!-- Updated form content -->
	{#if mounted}
		<div
			class="max-w-md w-full space-y-8 p-8 bg-white/80 backdrop-blur-sm rounded-xl shadow-xl"
			in:fly={{ y: 50, duration: 1000 }}
		>
			<div in:fade={{ duration: 1000, delay: 300 }}>
				<div class="flex space-x-4 mb-8">
					<button
						class="flex-1 py-2 border-b-2 transition-all duration-300 {activeTab === 'login'
							? 'border-amber-500 text-amber-600'
							: 'border-transparent text-gray-500 hover:text-amber-600'}"
						on:click={() => (activeTab = 'login')}
					>
						Sign In
					</button>
					<button
						class="flex-1 py-2 border-b-2 transition-all duration-300 {activeTab === 'register'
							? 'border-amber-500 text-amber-600'
							: 'border-transparent text-gray-500 hover:text-amber-600'}"
						on:click={() => (activeTab = 'register')}
					>
						Register
					</button>
				</div>
			</div>

			{#if error}
				<div class="bg-red-100/80 backdrop-blur-sm text-red-700 p-3 rounded-lg">
					{error}
				</div>
			{/if}

			{#if activeTab === 'login'}
				<form class="space-y-6" on:submit|preventDefault={handleLogin}>
					<div class="space-y-4">
						<div>
							<label for="email" class="block text-sm font-medium text-gray-700">Email</label>
							<input
								id="email"
								type="email"
								bind:value={email}
								class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-amber-500 focus:ring-amber-500 bg-white/70 backdrop-blur-sm"
								required
							/>
						</div>
						<div>
							<label for="password" class="block text-sm font-medium text-gray-700">Password</label>
							<input
								id="password"
								type="password"
								bind:value={password}
								class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-amber-500 focus:ring-amber-500 bg-white/70 backdrop-blur-sm"
								required
							/>
						</div>
					</div>
					<button
						type="submit"
						class="w-full flex justify-center py-3 px-4 rounded-full text-white bg-gradient-to-r from-amber-500 to-yellow-500 hover:scale-105 transition-all duration-300 hover:shadow-lg"
					>
						Sign in
					</button>
				</form>
			{:else}
				<form class="space-y-6" on:submit|preventDefault={handleRegister}>
					<div class="space-y-4">
						<div>
							<label for="reg-email" class="block text-sm font-medium text-gray-700">Email</label>
							<input
								id="reg-email"
								type="email"
								bind:value={email}
								class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-amber-500 focus:ring-amber-500 bg-white/70 backdrop-blur-sm"
								required
							/>
						</div>
						<div>
							<label for="username" class="block text-sm font-medium text-gray-700">Username</label>
							<input
								id="username"
								type="text"
								bind:value={username}
								class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-amber-500 focus:ring-amber-500 bg-white/70 backdrop-blur-sm"
								required
							/>
						</div>
						<div>
							<label for="fullName" class="block text-sm font-medium text-gray-700">Full Name</label
							>
							<input
								id="fullName"
								type="text"
								bind:value={fullName}
								class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-amber-500 focus:ring-amber-500 bg-white/70 backdrop-blur-sm"
								required
							/>
						</div>
						<div>
							<label for="reg-password" class="block text-sm font-medium text-gray-700"
								>Password</label
							>
							<input
								id="reg-password"
								type="password"
								bind:value={password}
								class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-amber-500 focus:ring-amber-500 bg-white/70 backdrop-blur-sm"
								required
							/>
						</div>
					</div>
					<button
						type="submit"
						class="w-full flex justify-center py-3 px-4 rounded-full text-white bg-gradient-to-r from-amber-500 to-yellow-500 hover:scale-105 transition-all duration-300 hover:shadow-lg"
					>
						Register
					</button>
				</form>
			{/if}
		</div>
	{/if}
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
