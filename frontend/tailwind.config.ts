import aspectRatio from '@tailwindcss/aspect-ratio';
import containerQueries from '@tailwindcss/container-queries';
import forms from '@tailwindcss/forms';
import typography from '@tailwindcss/typography';
import type { Config } from 'tailwindcss';
import { skeleton } from '@skeletonlabs/tw-plugin';
import { join } from 'path';

export default {
	content: ['./src/**/*.{html,js,svelte,ts}',
		join(require.resolve('@skeletonlabs/skeleton'), '../**/*.{html,js,svelte,ts}')],
	darkMode: 'class',
	theme: {
		extend: {
			fontFamily: {
				sans: ['Lexend', 'sans-serif']
			}
		}
	},

	plugins: [typography, forms, containerQueries, aspectRatio, skeleton({
		themes: { preset: ["skeleton"] }
	})]
} as Config;
