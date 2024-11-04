import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
	id: string;
	message: string;
	type: ToastType;
	duration?: number;
}

function createToastStore() {
	const { subscribe, update } = writable<Toast[]>([]);

	function addToast(message: string, type: ToastType = 'info', duration = 3000) {
		const id = Math.random().toString(36).substring(2);
		update((toasts) => [...toasts, { id, message, type, duration }]);

		if (duration) {
			setTimeout(() => {
				removeToast(id);
			}, duration);
		}
	}

	function removeToast(id: string) {
		update((toasts) => toasts.filter((t) => t.id !== id));
	}

	return {
		subscribe,
		success: (msg: string) => addToast(msg, 'success'),
		error: (msg: string) => addToast(msg, 'error'),
		info: (msg: string) => addToast(msg, 'info'),
		warning: (msg: string) => addToast(msg, 'warning'),
		remove: removeToast
	};
}

export const toasts = createToastStore();
