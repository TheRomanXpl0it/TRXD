<script lang="ts">
	import '../App.css';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import { ModeWatcher } from 'mode-watcher';
	import Layout from '$lib/components/Layout.svelte';
	import { authState, loadUser } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { addCollection } from '@iconify/svelte';
	import circleFlagsData from '@iconify-json/circle-flags/icons.json';

	addCollection(circleFlagsData);

	let { children } = $props();

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
		if (!authState.ready) {
			loadUser(false);
		}
	});
</script>

<QueryClientProvider client={queryClient}>
	<Toaster position="bottom-right" class="!justify-center md:!justify-end" />
	<ModeWatcher />

	{#if !authState.ready}
		<!-- Loading state -->
		<div class="flex h-screen w-full items-center justify-center">
			<div
				class="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-gray-900 dark:border-gray-600 dark:border-t-gray-100"
			></div>
		</div>
	{:else}
		<Layout user={authState.user} userMode={authState.userMode ?? false}>
			{@render children()}
		</Layout>
	{/if}
</QueryClientProvider>
