import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 5173,
		host: true,
		// Optimize HMR performance
		hmr: {
			port: 24678
		},
		// Enable filesystem watching optimizations
		watch: {
			usePolling: false,
			interval: 100
		}
	},
	// Optimize dependency pre-bundling for faster cold starts
	optimizeDeps: {
		include: [
			'@sveltejs/kit',
			'svelte',
			'svelte/store',
			'svelte/transition',
			'svelte/animate'
		],
		// Exclude large libraries that don't need pre-bundling
		exclude: []
	},
	// Enable route-level code splitting for lazy loading
	build: {
		target: 'esnext',
		rollupOptions: {
			output: {
				manualChunks: (id) => {
					// Core SvelteKit chunks
					if (id.includes('@sveltejs/kit')) {
						return 'svelte-kit';
					}
					// Spreadsheet-specific chunks
					if (id.includes('Spreadsheet') || id.includes('spreadsheet-store')) {
						return 'spreadsheet';
					}
					// Forms-specific chunks
					if (id.includes('Form')) {
						return 'forms';
					}
					// Calendar-specific chunks
					if (id.includes('Calendar')) {
						return 'calendar';
					}
					// Email-specific chunks
					if (id.includes('Email')) {
						return 'email';
					}
					// Task-specific chunks
					if (id.includes('Task')) {
						return 'tasks';
					}
				}
			}
		}
	},
	// Enable caching for better performance
	cacheDir: '.vite-cache'
});
