import { writable } from 'svelte/store';
import type { User } from '$lib/types';

interface AuthState {
	user: User | null;
	token: string | null;
}

function createAuthStore() {
	const { subscribe, set } = writable<AuthState>({
		user: null,
		token: null
	});

	// Initialize the store with stored data
	if (typeof window !== 'undefined') {
		const stored = localStorage.getItem('auth');
		if (stored) {
			try {
				const auth = JSON.parse(stored);
				set(auth);
			} catch (e) {
				console.error('Failed to parse stored auth', e);
			}
		}
	}

	return {
		subscribe,
		login: (user: User, token: string) => {
			const authState = { user, token };
			set(authState);
			localStorage.setItem('auth', JSON.stringify(authState));
		},
		logout: () => {
			set({ user: null, token: null });
			localStorage.removeItem('auth');
		},
		initialize: () => {
			const stored = localStorage.getItem('auth');
			if (stored) {
				try {
					const auth = JSON.parse(stored);
					set(auth);
				} catch (e) {
					console.error('Failed to parse stored auth', e);
					localStorage.removeItem('auth');
					set({ user: null, token: null });
				}
			}
		}
	};
}

export const auth = createAuthStore();
