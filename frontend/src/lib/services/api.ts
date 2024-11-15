import type { User } from '$lib/types';

interface LoginResponse {
	user: User;
	token: string;
}

/*
	beekeeper@test.com / password
	manager@test.com / password
	admin@test.com / password
 */
export async function login(email: string, password: string): Promise<LoginResponse> {
	// Mock API call
	const users = {
		'beekeeper@test.com': {
			id: '1',
			name: 'John Beekeeper',
			role: 'beekeeper',
			email: 'beekeeper@test.com'
		},
		'manager@test.com': {
			id: '2',
			name: 'Jane Manager',
			role: 'manager',
			email: 'manager@test.com'
		},
		'admin@test.com': {
			id: '3',
			name: 'Admin User',
			role: 'admin',
			email: 'admin@test.com'
		}
	} as const;

	// Simulate API delay
	await new Promise(resolve => setTimeout(resolve, 500));

	const user = users[email as keyof typeof users];
	if (user && password === 'password') {
		return {
			user,
			token: `mock-token-${user.role}-${Date.now()}`
		};
	}

	throw new Error('Invalid credentials');
}

export async function getBeeCommunities() {
	const response = await fetch('/api/communities');
	return response.json();
}

export async function getHoneyHarvests() {
	const response = await fetch('/api/harvests');
	return response.json();
}
