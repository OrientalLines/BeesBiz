import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

function createThemeStore() {
	const { subscribe, set } = writable<Theme>('light');

	return {
		subscribe,
		toggle: () => {
			if (browser) {
				const html = document.documentElement;
				const currentTheme = html.classList.contains('dark') ? 'dark' : 'light';
				const newTheme = currentTheme === 'light' ? 'dark' : 'light';
				
				html.classList.remove(currentTheme);
				html.classList.add(newTheme);
				
				localStorage.setItem('theme', newTheme);
				set(newTheme);
			}
		},
		set: (value: Theme) => {
			if (browser) {
				const html = document.documentElement;
				html.classList.remove('light', 'dark');
				html.classList.add(value);
				localStorage.setItem('theme', value);
				set(value);
			}
		}
	};
}

export const theme = createThemeStore();
