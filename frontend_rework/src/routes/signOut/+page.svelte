<script lang="ts">
      import Spinner from '$lib/components/ui/spinner/spinner.svelte';
      import { logout } from '@/auth';
      import { user } from '$lib/stores/auth';
      
      
      let loading = $state(true);
      $effect(() => {
        (async () => {
          try {
            await logout();
          } catch (e) {
            console.error('Logout failed', e);
          } finally {
            loading = false;
            $user = null;
            const q = new URLSearchParams(location.search);
            const dest = q.get("redirect") || "/";
            window.location.replace(dest);
            
          }
        })();
      });
        
</script>


{#if loading}
  <div class="flex flex-row items-center gap-2 text-gray-500">
    <Spinner /><p>Logging out…</p>
  </div>
{:else}
  <p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">You have been logged out.</p>
  <hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
  <p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
    <a href="/signIn" class="text-primary-700 dark:text-primary-500 hover:underline">Click here to log in again.</a>
  </p>
{/if}

