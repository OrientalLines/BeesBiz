import { writable } from 'svelte/store';
import type { User } from '$lib/types';

function createAuthStore() {
	const storedUser = typeof window !== 'undefined' ? localStorage.getItem('user') : null;
	const initialValue = storedUser ? JSON.parse(storedUser) : null;

	const { subscribe, set, update } = writable<User | null>(initialValue);

	return {
		subscribe,
		login: (user: User) => {
			localStorage.setItem('user', JSON.stringify(user));
			set(user);
		},
		logout: () => {
			localStorage.removeItem('user');
			set(null);
		},
		updateUser: (userData: Partial<User>) =>
			update((user) => {
				const updatedUser = user ? { ...user, ...userData } : null;
				if (updatedUser) {
					localStorage.setItem('user', JSON.stringify(updatedUser));
				}
				return updatedUser;
			})
	};
}

export const auth = createAuthStore();
