import { writable } from 'svelte/store';
import type { User } from '$lib/types';

interface AuthState {
	user: User | null;
	token: string | null;
}

function createAuthStore() {
	const storedAuth = typeof window !== 'undefined' ? localStorage.getItem('auth') : null;
	const initialValue: AuthState = storedAuth ? JSON.parse(storedAuth) : { user: null, token: null };

	const { subscribe, set, update } = writable<AuthState>(initialValue);

	return {
		subscribe,
		login: (user: User, token: string) => {
			const authState = { user, token };
			localStorage.setItem('auth', JSON.stringify(authState));
			set(authState);
		},
		logout: () => {
			localStorage.removeItem('auth');
			set({ user: null, token: null });
		},
		updateUser: (userData: Partial<User>) =>
			update((state) => {
				if (!state.user) return state;
				const updatedUser = { ...state.user, ...userData };
				const newState = { ...state, user: updatedUser };
				localStorage.setItem('auth', JSON.stringify(newState));
				return newState;
			})
	};
}

export const auth = createAuthStore();
