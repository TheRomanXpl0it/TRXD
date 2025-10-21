<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Checkbox } from "$lib/components/ui/checkbox/index.js";
  import { X } from "@lucide/svelte";
  import { createEventDispatcher } from "svelte";

  export type Item = { value: string; label: string };

  // Props (runes)
  let {
    items = [] as Item[],
    placeholder = "Select tags…",
    value = $bindable<string[]>([])
  } = $props<{ items?: Item[]; placeholder?: string; value?: string[] }>();

  const dispatch = createEventDispatcher<{ create: string }>();

  // Local state
  let open   = $state(false);
  let query  = $state("");
  let active = $state(0); // highlighted index in filtered list

  // Keep locally-created tags so they appear immediately
  let custom = $state<string[]>([]);

  // Normalized + available
  const normalized = $derived(items.map(i => ({ value: String(i.value), label: String(i.label) })));
  const available  = $derived([
    ...normalized,
    ...custom.map(v => ({ value: v, label: v }))
  ]);

  const q       = $derived(query.trim());
  const qLower  = $derived(q.toLowerCase());
  const filtered = $derived(
    q
      ? available.filter(i =>
          i.label.toLowerCase().includes(qLower) || i.value.toLowerCase().includes(qLower)
        )
      : available
  );

  const labelMap        = $derived(new Map(available.map(i => [i.value, i.label] as const)));
  const availableLower  = $derived(new Set(available.map(i => i.value.toLowerCase())));
  const selectedLower   = $derived(new Set(value.map(v => v.toLowerCase())));
  const canCreate       = $derived(q.length > 0 && !availableLower.has(qLower));

  const summary = $derived(
    value.length === 0
      ? placeholder
      : value.length === 1
      ? (labelMap.get(value[0]) ?? value[0])
      : `${value.length} tags selected`
  );

  function toggle(v: string) {
    const set = new Set(value);
    set.has(v) ? set.delete(v) : set.add(v);
    value = Array.from(set);
  }

  function remove(v: string) {
    value = value.filter(x => x !== v);
  }

  function createTagFromQuery() {
    if (!canCreate) return;
    custom = Array.from(new Set([...custom, q]));
    if (!selectedLower.has(qLower)) value = [...value, q];
    dispatch("create", q);
    query = "";
    active = 0;
  }

  // Adopt unknown selected tags so they render
  $effect(() => {
    // use availableLower so this effect re-runs when available changes
    const _dep = availableLower;
    const unknown = value.filter(v => !availableLower.has(v.toLowerCase()));
    if (unknown.length) custom = Array.from(new Set([...custom, ...unknown]));
  });

  // Reset highlight when the search query changes
  $effect(() => {
    const _q = q; // dependency
    active = 0;
  });

  // Clamp active when filtered length changes (no second arg!)
  $effect(() => {
    const len = filtered.length; // dependency
    if (len === 0) {
      active = 0;
    } else if (active >= len) {
      active = len - 1;
    } else if (active < 0) {
      active = 0;
    }
  });

  // Autofocus search when opened
  let searchEl: HTMLInputElement | null = null;
  $effect(() => {
    const _open = open; // dependency
    if (!_open) return;
    queueMicrotask(() => searchEl?.focus());
  });

  function onSearchKeydown(e: KeyboardEvent) {
    // Prevent outer form submit
    if (e.key === "Enter") e.preventDefault();

    if (e.key === "Escape") {
      open = false;
      return;
    }

    if (e.key === "Tab" && canCreate) {
      e.preventDefault();
      createTagFromQuery();
      return;
    }

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
      if (canCreate) {
        createTagFromQuery();
      } else if (filtered.length > 0) {
        const i = Math.max(0, Math.min(active, filtered.length - 1));
        toggle(filtered[i].value);
      }
    }
  }
</script>

<div class="flex flex-col gap-2">
  {#if value.length > 0}
    <div class="flex flex-wrap gap-2">
      {#each value as v (v)}
        <span class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-sm">
          {labelMap.get(v) ?? v}
          <button
            type="button"
            class="ml-1 opacity-70 hover:opacity-100"
            on:click={() => remove(v)}
            aria-label={`Remove ${labelMap.get(v) ?? v}`}
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
        <Input
          placeholder="Search or create…"
          bind:this={searchEl}
          bind:value={query}
          autocomplete="off"
          onkeydown={onSearchKeydown}
        />

        {#if canCreate}
          <button
            type="button"
            class="flex w-full items-center gap-2 px-2 py-1.5 rounded border border-dashed hover:bg-accent text-left"
            on:click={createTagFromQuery}
          >
            + Create tag “{q}”
          </button>
        {/if}

        <div class="max-h-60 overflow-auto pr-1">
          {#if filtered.length === 0 && !canCreate}
            <div class="text-sm text-muted-foreground px-2 py-1.5">No results</div>
          {:else}
            {#each filtered as it, i (it.value)}
              <button
                type="button"
                class={`flex w-full items-center gap-2 px-2 py-1.5 rounded hover:bg-accent text-left ${i === active ? 'bg-accent/60' : ''}`}
                on:click={() => toggle(it.value)}
                aria-selected={i === active}
              >
                <Checkbox checked={value.includes(it.value)} />
                <span class="truncate">{it.label}</span>
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
