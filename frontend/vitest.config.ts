import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { resolve } from 'path';

export default defineConfig({
	plugins: [svelte({ hot: !process.env.VITEST })],
	test: {
		environment: 'jsdom',
		globals: true,
		setupFiles: ['./src/test/setup.ts'],
		include: ['src/**/*.{test,spec}.{js,ts}'],
		server: {
			deps: {
				inline: ['@lucide/svelte']
			}
		}
	},
	resolve: {
		alias: {
			$lib: resolve('./src/lib'),
			'@': resolve('./src/lib'),
			'$app/navigation': resolve('./src/test/mocks/app-navigation.ts'),
			'$app/stores': resolve('./src/test/mocks/app-stores.ts')
		},
		conditions: ['browser']
	}
});
