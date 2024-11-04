<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { fly, fade } from 'svelte/transition';

	let email = '';
	let password = '';
	let error = '';

	let mounted = false;
	onMount(() => {
		mounted = true;
		if ($auth) {
			goto('/dashboard/hives');
		}
	});

	async function handleLogin() {
		try {
			// Mock login - replace with actual API call
			if (email && password) {
				auth.login({
					id: '1',
					name: 'John Beekeeper',
					role: 'beekeeper',
					email
				});
				goto('/dashboard/hives');
			} else {
				error = 'Invalid credentials';
			}
		} catch (err) {
			error = 'Login failed';
		}
	}

	let isLogin = true;
	let name = '';
	let confirmPassword = '';
	
	async function handleRegister() {
		try {
			if (!name || !email || !password || !confirmPassword) {
				error = 'All fields are required';
				return;
			}
			
			if (password !== confirmPassword) {
				error = 'Passwords do not match';
				return;
			}

			// Mock registration - replace with actual API call
			auth.login({
				id: '1',
				name,
				role: 'beekeeper',
				email
			});
			goto('/dashboard/hives');
		} catch (err) {
			error = 'Registration failed';
		}
	}

	$: formTitle = isLogin ? 'Sign in to your account' : 'Create your account';
	$: handleSubmit = isLogin ? handleLogin : handleRegister;

	// Generate initial random positions for bees
	const bees = Array(5)
		.fill(null)
		.map(() => ({
			baseX: Math.random() * 100,
			baseY: Math.random() * 100
		}));
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-amber-50 via-blue-50 to-white relative overflow-hidden">
	<!-- Animated background bees -->
	{#if mounted}
		{#each bees as bee, i}
			<div
				class="absolute opacity-20 pointer-events-none"
				style="top: {bee.baseY}%; left: {bee.baseX}%;"
				in:fly={{ y: 100, duration: 1000, delay: i * 200 }}
			>
				<div class="text-4xl animate-float-{i + 1}">
					🐝
				</div>
			</div>
		{/each}
	{/if}

	<!-- Existing login form content -->
	{#if mounted}
		<div 
			class="max-w-md w-full space-y-8 p-8 bg-white/80 backdrop-blur-sm rounded-xl shadow-xl"
			in:fly={{ y: 50, duration: 1000 }}
		>
			<div in:fade={{ duration: 1000, delay: 300 }}>
				<h2 class="text-center text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-amber-600 to-yellow-600">
					{formTitle}
				</h2>
			</div>

			<form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
				{#if error}
					<div class="bg-red-100/80 backdrop-blur-sm text-red-700 p-3 rounded-lg">
						{error}
					</div>
				{/if}
				<div class="space-y-4">
					{#if !isLogin}
						<div>
							<label for="name" class="block text-sm font-medium text-gray-700">Full Name</label>
							<input
								id="name"
								type="text"
								bind:value={name}
								class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-amber-500 focus:ring-amber-500 bg-white/70 backdrop-blur-sm"
								required
							/>
						</div>
					{/if}
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
					{#if !isLogin}
						<div>
							<label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirm Password</label>
							<input
								id="confirmPassword"
								type="password"
								bind:value={confirmPassword}
								class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-amber-500 focus:ring-amber-500 bg-white/70 backdrop-blur-sm"
								required
							/>
						</div>
					{/if}
				</div>
				<button
					type="submit"
					class="w-full flex justify-center py-3 px-4 rounded-full text-white bg-gradient-to-r from-amber-500 to-yellow-500 hover:scale-105 transition-all duration-300 hover:shadow-lg"
				>
					{isLogin ? 'Sign in' : 'Create account'}
				</button>
				
				<div class="text-center">
					<button 
						type="button"
						class="text-amber-600 hover:text-amber-500 text-sm font-medium transition-colors"
						on:click={() => {
							isLogin = !isLogin;
							error = '';
						}}
					>
						{isLogin ? "Don't have an account? Sign up" : 'Already have an account? Sign in'}
					</button>
				</div>
			</form>
		</div>
	{/if}
</div>

<style>
	@keyframes float-1 {
		0%, 100% { transform: translate(0, 0) rotate(0deg); }
		50% { transform: translate(10px, -10px) rotate(5deg); }
	}
	@keyframes float-2 {
		0%, 100% { transform: translate(0, 0) rotate(0deg); }
		50% { transform: translate(-8px, -12px) rotate(-3deg); }
	}
	@keyframes float-3 {
		0%, 100% { transform: translate(0, 0) rotate(0deg); }
		50% { transform: translate(12px, -8px) rotate(8deg); }
	}
	@keyframes float-4 {
		0%, 100% { transform: translate(0, 0) rotate(0deg); }
		50% { transform: translate(-10px, -10px) rotate(-8deg); }
	}
	@keyframes float-5 {
		0%, 100% { transform: translate(0, 0) rotate(0deg); }
		50% { transform: translate(8px, -12px) rotate(3deg); }
	}

	.animate-float-1 { animation: float-1 8s infinite ease-in-out; }
	.animate-float-2 { animation: float-2 9s infinite ease-in-out; }
	.animate-float-3 { animation: float-3 10s infinite ease-in-out; }
	.animate-float-4 { animation: float-4 11s infinite ease-in-out; }
	.animate-float-5 { animation: float-5 12s infinite ease-in-out; }
</style>
