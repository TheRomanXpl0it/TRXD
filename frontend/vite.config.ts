import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import { resolve } from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

export default defineConfig({
	publicDir: 'static',
	plugins: [tailwindcss(), svelte()],
	optimizeDeps: {
		include: ['@iconify/svelte', '@iconify-json/circle-flags']
	},
	resolve: {
		alias: {
			$lib: resolve(dirname(fileURLToPath(import.meta.url)), 'src/lib'),
			'@': resolve(dirname(fileURLToPath(import.meta.url)), 'src/lib') // so "@/components/..." works
		}
	},
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:1337',
				changeOrigin: true,
				secure: false
			}
		}
	},
	ssr: {
		noExternal: ['bits-ui', 'vaul-svelte', 'svelte-sonner', 'svelte-motion', 'paneforge']
	}
});
