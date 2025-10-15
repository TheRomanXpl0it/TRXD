<script lang="ts">
  import CheckIcon from "@lucide/svelte/icons/check";
  import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
  import { tick } from "svelte";
  import * as Command from "$lib/components/ui/command/index.js";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { cn } from "$lib/utils.js";

  type Item = { value: string; label: string };

  let {
    items = [] as Item[],
    // 👇 make this bindable so parent can use `bind:value`
    value = $bindable<string>(""),
    placeholder = "Select a category…",
    searchPlaceholder = "Search category…",
    groupLabel = "categories",
    className = "",
    widthClass = "w-[220px]"
  } = $props<{
    items?: Item[];
    value?: string;                // now bindable thanks to $bindable above
    placeholder?: string;
    searchPlaceholder?: string;
    groupLabel?: string;
    className?: string;
    widthClass?: string;
  }>();

  let open = $state(false);
  let triggerRef = $state<HTMLButtonElement>(null!);

  const selectedLabel = $derived(items.find((i) => i.value === value)?.label);

  function closeAndFocusTrigger() {
    open = false;
    tick().then(() => triggerRef?.focus());
  }
</script>

<Popover.Root bind:open>
  <Popover.Trigger bind:ref={triggerRef}>
    {#snippet child({ props })}
      <Button
        {...props}
        variant="outline"
        role="combobox"
        aria-expanded={open}
        class={cn(widthClass, "justify-between", className)}
      >
        {selectedLabel || placeholder}
        <ChevronsUpDownIcon class="opacity-50" />
      </Button>
    {/snippet}
  </Popover.Trigger>

  <Popover.Content class={cn(widthClass, "p-1")}>
    <Command.Root>
      <Command.Input placeholder={searchPlaceholder} class="border-0 shadow-none ring-0 focus:ring-0 focus:outline-none focus-visible:ring-0 focus-visible:outline-none" />
      <Command.List>
        <Command.Empty>No results.</Command.Empty>
        <Command.Group value={groupLabel}>
          {#each items as item (item.value)}
            <Command.Item
              value={item.value}
              onSelect={() => {
                value = item.value;       // updates parent via binding
                closeAndFocusTrigger();
              }}
            >
              <CheckIcon class={cn(value !== item.value && "text-transparent")} />
              {item.label}
            </Command.Item>
          {/each}
        </Command.Group>
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
