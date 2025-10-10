import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import path from 'path';

export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  resolve: {
    alias: {
      $lib: path.resolve(__dirname, 'src/lib'),
      '@': path.resolve(__dirname, 'src/lib'), // so "@/components/..." works
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:1337',
        changeOrigin: true,
        secure: false,
      }
    }
  },
  ssr: {
    noExternal: [
      'bits-ui',
      'vaul-svelte',
      'svelte-sonner',
      'svelte-motion',
      'paneforge',
    ]
  }
});
