<script lang="ts">
	import './App.css';

	import favicon from '$lib/assets/favicon.ico';

	import Layout from '$lib/components/Layout.svelte';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import { ModeWatcher } from 'mode-watcher';
	import Router from 'svelte-spa-router';
	import routes from './routes';
	import { user, userMode, loadUser, authReady } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 30_000,
				gcTime: 5 * 60 * 1000,
				retry: 1,
				refetchOnWindowFocus: false
			}
		}
	});

	loadUser(false);

	onMount(() => {
		// Ensure user is loaded
		if (!$authReady) {
			loadUser(false);
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<QueryClientProvider client={queryClient}>
	<Toaster position="bottom-right" class="!justify-center md:!justify-end" />
	<ModeWatcher />

	{#if !$authReady}
		<!-- Loading state -->
		<div class="flex h-screen w-full items-center justify-center">
			<div
				class="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-gray-900 dark:border-gray-600 dark:border-t-gray-100"
			></div>
		</div>
	{:else}
		<Layout user={$user} userMode={$userMode ?? false}>
			<Router {routes} useHash={true} restoreScrollState={true} />
		</Layout>
	{/if}
</QueryClientProvider>
