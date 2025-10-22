<script lang="ts">
  import { onMount } from 'svelte';
  import './App.css';

  import { MoonIcon, SunIcon } from '@lucide/svelte';
  import favicon from '$lib/assets/favicon.ico';

  import * as Sidebar from '$lib/components/ui/sidebar/index.js';
  import AppSidebar from '$lib/components/app-sidebar.svelte';

  import { toggleMode, ModeWatcher } from 'mode-watcher';
  import { Button } from '@/components/ui/button';
  import { Toaster } from '$lib/components/ui/sonner/index.js';

  import Router from 'svelte-spa-router';
  import routes from './routes';

  import { user, authReady, loadUser } from '$lib/stores/auth';


  onMount(() => {
    loadUser(); // no const, don't shadow `user`
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<Toaster />
<ModeWatcher />

<Sidebar.Provider>
  <!-- pass the actual store value -->
  <AppSidebar user={$user} />

  <main class="flex w-full h-full">
    <Sidebar.Trigger size="lg" />

    <div class="ml-10 mr-15 flex w-full h-full flex-col">
      {#if $authReady}
        <Router {routes} useHash={false} />
      {:else}
        <div class="p-4">Loading…</div>
      {/if}
    </div>

    <Button
      size="icon"
      onclick={toggleMode}
      variant="outline"
      class="fixed bottom-4 right-4 p-2"
      aria-label="Toggle theme"
    >
      <SunIcon class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 !transition-all dark:-rotate-90 dark:scale-0" />
      <MoonIcon class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 !transition-all dark:rotate-0 dark:scale-100" />
    </Button>
  </main>
</Sidebar.Provider>
