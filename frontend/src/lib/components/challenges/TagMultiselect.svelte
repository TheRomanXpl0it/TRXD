<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Checkbox } from "$lib/components/ui/checkbox/index.js";
  import { X } from "@lucide/svelte";

  let {
    all_tags = [] as string[],
    placeholder = "Select tags...",
    value = $bindable<string[]>([]),
    oncreate
  } = $props<{ 
    all_tags?: string[]; 
    placeholder?: string; 
    value?: string[];
    oncreate?: (tag: string) => void;
  }>();

  let open   = $state(false);
  let query  = $state("");
  let active = $state(0);
  let custom = $state<string[]>([]); // newly created tags (raw strings)

  const available = $derived([...(all_tags ?? []), ...(custom ?? [])]);

  // filter (raw, case-sensitive); empty query -> show everything
  const filtered = $derived(query ? available.filter(v => String(v).includes(query)) : available);

  const canCreate = $derived(query.length > 0 && !available.includes(query));

  const summary = $derived(
    Array.isArray(value) && value.length > 0 ? `${value.length} selected` : placeholder
  );

  function toggle(v: string) {
    const arr = Array.isArray(value) ? [...value] : [];
    const idx = arr.indexOf(v);
    if (idx !== -1) arr.splice(idx, 1);
    else arr.push(v);
    value = arr;
  }

  function remove(v: string) {
    value = (Array.isArray(value) ? value : []).filter(x => x !== v);
  }

  function createTagFromQuery() {
    if (!canCreate) return;
    const tag = query;
    custom = [...custom, tag];
    if (!value.includes(tag)) value = [...value, tag];
    oncreate?.(tag);
    query = "";
    active = 0;
  }

  // focus search when opened
  let searchEl: HTMLInputElement | null = null;
  $effect(() => {
    if (open) queueMicrotask(() => searchEl?.focus());
  });

  // reset highlight on query change
  $effect(() => {
    const _ = query; // just to trigger
    active = 0;
  });

  // keyboard navigation + create on Enter
  function onSearchKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") { open = false; return; }
    if (e.key === "ArrowDown") {
      e.preventDefault();
      if (filtered.length) active = (active + 1) % filtered.length;
      return;
    }
    if (e.key === "ArrowUp") {
      e.preventDefault();
      if (filtered.length) active = (active - 1 + filtered.length) % filtered.length;
      return;
    }
    if (e.key === "Enter") {
      e.preventDefault();
      if (canCreate) createTagFromQuery();
      else if (filtered.length) {
        const idx = Math.max(0, Math.min(active, filtered.length - 1));
        toggle(filtered[idx]);
      }
    }
  }
</script>

<div class="flex flex-col gap-2">
  {#if Array.isArray(value) && value.length > 0}
    <div class="flex flex-wrap gap-2">
      {#each value as v (v)}
        <span class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-sm">
          {v}
          <button
            type="button"
            class="ml-1 opacity-70 hover:opacity-100"
            onclick={() => remove(v)}
            aria-label={`Remove ${v}`}
            title="Remove"
          >
            <X class="h-4 w-4" />
          </button>
        </span>
      {/each}
    </div>
  {/if}

  <Popover.Root bind:open>
    <Popover.Trigger>
      <Button type="button" variant="outline" class="justify-between w-full">
        <span class="truncate">{summary}</span>
      </Button>
    </Popover.Trigger>

    <Popover.Content class="w-80 p-3">
      <div class="flex flex-col gap-3">
        <input
          placeholder="Search or create..."
          bind:this={searchEl}
          bind:value={query}
          autocomplete="off"
          onkeydown={onSearchKeydown}
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm
                 placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2
                 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
        />

        {#if canCreate}
          <button
            type="button"
            class="flex w-full items-center gap-2 px-2 py-1.5 rounded border border-dashed hover:bg-accent text-left"
            onclick={createTagFromQuery}
          >
            + Create tag “{query}”
          </button>
        {/if}

        <div class="max-h-60 overflow-auto pr-1">
          {#if filtered.length === 0 && !canCreate}
            <div class="text-sm text-muted-foreground px-2 py-1.5">No results</div>
          {:else}
            {#each filtered as v, i (i)}
              <button
                type="button"
                class={`flex w-full items-center gap-2 px-2 py-1.5 rounded hover:bg-accent text-left ${i === active ? 'bg-accent/60' : ''}`}
                onclick={() => toggle(v)}
                aria-pressed={(value ?? []).includes(v)}
              >
                <Checkbox checked={(value ?? []).includes(v)} />
                <span class="truncate">{v}</span>
              </button>
            {/each}
          {/if}
        </div>

        <div class="flex justify-end">
          <Button type="button" size="sm" onclick={() => (open = false)}>Done</Button>
        </div>
      </div>
    </Popover.Content>
  </Popover.Root>
</div>