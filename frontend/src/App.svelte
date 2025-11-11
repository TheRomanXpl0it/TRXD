<script lang="ts">
	import './App.css';

	import { MoonIcon, SunIcon } from '@lucide/svelte';
	import favicon from '$lib/assets/favicon.ico';

	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppSidebar from '$lib/components/AppSidebar.svelte';

	import { toggleMode, ModeWatcher } from 'mode-watcher';
	import { Button } from '@/components/ui/button';
	import { Toaster } from '$lib/components/ui/sonner/index.js';

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
		<div class="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-gray-900 dark:border-gray-600 dark:border-t-gray-100"></div>
	</div>
{:else}
<Sidebar.Provider>
	<!-- pass the actual store value -->
	<AppSidebar user={$user} userMode={$userMode ?? false} />

	<!-- Mobile/Tablet header with menu trigger (fixed) -->
	<div class="fixed left-0 right-0 top-0 z-50 flex items-center justify-between border-b border-gray-200 bg-background p-3 dark:border-gray-700 lg:hidden overflow-visible">
		<Sidebar.Trigger size="icon" variant="ghost" class="cursor-pointer relative z-10">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<line x1="4" x2="20" y1="12" y2="12"></line>
				<line x1="4" x2="20" y1="6" y2="6"></line>
				<line x1="4" x2="20" y1="18" y2="18"></line>
			</svg>
		</Sidebar.Trigger>

		<span class="text-lg font-bold text-gray-900 dark:text-white">TRXD</span>

		<Button
			size="icon"
			onclick={toggleMode}
			variant="ghost"
			class="cursor-pointer"
			aria-label="Toggle theme"
		>
			<SunIcon
				class="h-5 w-5 rotate-0 scale-100 !transition-all dark:-rotate-90 dark:scale-0"
			/>
			<MoonIcon
				class="absolute h-5 w-5 rotate-90 scale-0 !transition-all dark:rotate-0 dark:scale-100"
			/>
		</Button>
	</div>

	<main class="flex min-h-screen w-full flex-col pt-[57px] lg:pt-0">
		<div class="router-content flex-1 mx-3 pb-24 lg:mx-6 lg:pb-16">
			<Router
				{routes}
				useHash={true}
				restoreScrollState={true}
			/>
		</div>

		<!-- Desktop theme toggle -->
		<Button
			size="icon"
			onclick={toggleMode}
			variant="outline"
			class="fixed bottom-6 right-6 z-50 hidden p-2 shadow-lg lg:flex"
			aria-label="Toggle theme"
		>
			<SunIcon
				class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 !transition-all dark:-rotate-90 dark:scale-0"
			/>
			<MoonIcon
				class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 !transition-all dark:rotate-0 dark:scale-100"
			/>
		</Button>
	</main>
</Sidebar.Provider>
{/if}
</QueryClientProvider>
