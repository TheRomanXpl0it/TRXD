<script lang="ts">
  import * as Sheet from "$lib/components/ui/sheet/index.js";
  import { toast } from "svelte-sonner";
  import { getSolves } from "@/challenges";

  let { open = $bindable(false), challenge } =
    $props<{ open?: boolean; challenge: any }>();

  let loading = $state(false);
  let solves = $state<any[]>([]);

  function runSolvesEffect(isOpen: boolean, id?: string) {
    if (!isOpen || !id) return;
    const ac = new AbortController();
    loading = true;

    (async () => {
      try {
        const data = await getSolves(id)
        if (!ac.signal.aborted && open && challenge?.id === id) {
          solves = data;
        }
      } catch (e: any) {
        if (e?.name !== 'AbortError') toast.error(e instanceof Error ? e.message : 'Failed to load solves');
      } finally {
        if (!ac.signal.aborted) loading = false;
      }
    })();
  
    return () => ac.abort();
  }
  
  $effect(() => runSolvesEffect(open, challenge?.id));
</script>

<Sheet.Root bind:open>
  <Sheet.Content side="left" class="sm:max-w-[640px]">
    <Sheet.Header>
      <Sheet.Title>Solves</Sheet.Title>
      <Sheet.Description>Recent solvers for {challenge?.name}</Sheet.Description>
    </Sheet.Header>

    {#if loading}
      <p class="py-6 text-sm text-gray-500">Loading…</p>
    {:else if solves.length === 0}
      <p class="py-6 text-sm text-gray-500">No solves yet.</p>
    {:else}
      <ul class="divide-y dark:divide-gray-800">
        {#each solves as s}
          <li class="py-2 text-sm">
            <span class="font-medium">{s.user?.name ?? s.user ?? 'Anonymous'}</span>
            <span class="ml-2 text-gray-500">
              {new Date(s.createdAt ?? s.time ?? Date.now()).toLocaleString()}
            </span>
            {#if s.points}
              <span class="ml-2 text-gray-500">{s.points} pts</span>
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </Sheet.Content>
</Sheet.Root>