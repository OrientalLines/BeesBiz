import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

function createThemeStore() {
	const initialTheme = browser ? (localStorage.getItem('theme') as Theme) || 'light' : 'light';

	if (browser) {
		document.documentElement.classList.remove('light', 'dark');
		document.documentElement.classList.add(initialTheme);
	}

	const { subscribe, set } = writable<Theme>(initialTheme);

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
		},
		initialize: () => {
			if (browser) {
				const savedTheme = localStorage.getItem('theme') as Theme;
				if (savedTheme) {
					set(savedTheme);
				} else {
					const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
					set(prefersDark ? 'dark' : 'light');
				}
			}
		}
	};
}

export const theme = createThemeStore();
