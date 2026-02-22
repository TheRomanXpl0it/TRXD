import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { resolve } from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const getGitHash = () => {
	return process.env.VITE_GIT_HASH ?? 'unknown';
};

export default defineConfig({
	define: {
		__GIT_HASH__: JSON.stringify(getGitHash())
	},
	publicDir: 'static',
	plugins: [tailwindcss(), sveltekit()],
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
		noExternal: ['bits-ui', 'vaul-svelte', 'svelte-sonner', 'svelte-motion', 'paneforge', '@lucide/svelte']
	},
	build: {
		rollupOptions: {
			output: {
				experimentalMinChunkSize: 10 * 1024 // 10KB
			}
		}
	}
});
