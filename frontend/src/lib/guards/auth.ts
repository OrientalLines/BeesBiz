import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { auth } from '$lib/stores/auth';
import { get } from 'svelte/store';

const PUBLIC_ROUTES = ['/login'];

export function requireAuth(path: string): boolean {
	if (!browser) return true;
	if (PUBLIC_ROUTES.includes(path)) return true;

	const authState = get(auth);

	if (!authState.user || !authState.token) {
		if (path !== '/login') {
			goto('/login');
		}
		return false;
	}

	const adminOnlyPaths = ['/dashboard/users', '/dashboard/settings'];
	const managerPaths = ['/dashboard/reports'];

	if (adminOnlyPaths.some((p) => path.startsWith(p)) && authState.user.role !== 'ADMIN') {
		goto('/dashboard/hives');
		return false;
	}

	if (
		managerPaths.some((p) => path.startsWith(p)) &&
		authState.user.role !== 'ADMIN' &&
		authState.user.role !== 'MANAGER'
	) {
		goto('/dashboard/hives');
		return false;
	}

	return true;
}
