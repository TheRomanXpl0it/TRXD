<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import { tick } from 'svelte';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { cn } from '$lib/utils.js';
	import VirtualList from '@sveltejs/svelte-virtual-list';

	type Item = { value: string; label: string };

	let {
		items = [] as Item[],
		value = $bindable<string>(''),
		placeholder = 'Select a category...',
		searchPlaceholder = 'Search category...',
		groupLabel = 'categories',
		className = '',
		widthClass = 'w-[220px]',
		id = ''
	} = $props<{
		items?: Item[];
		value?: string;
		placeholder?: string;
		searchPlaceholder?: string;
		groupLabel?: string;
		className?: string;
		widthClass?: string;
		id?: string;
	}>();

	let open = $state(false);
	let triggerRef = $state<HTMLButtonElement>(null!);
	let searchValue = $state('');

	const selectedLabel = $derived(items.find((i: any) => i.value === value)?.label);

	const filteredItems = $derived(
		searchValue.trim()
			? items.filter(
					(item: Item) =>
						(item.label || '').toLowerCase().includes(searchValue.toLowerCase()) ||
						(item.value || '').toLowerCase().includes(searchValue.toLowerCase())
				)
			: items
	);

	function closeAndFocusTrigger() {
		open = false;
		searchValue = '';
		tick().then(() => triggerRef?.focus());
	}

	function selectItem(item: Item) {
		value = item.value;
		closeAndFocusTrigger();
	}
</script>

<Popover.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (!isOpen) searchValue = '';
	}}
>
	<Popover.Trigger bind:ref={triggerRef}>
		{#snippet child({ props })}
			<Button
				{...props}
				variant="outline"
				role="combobox"
				aria-expanded={open}
				class={cn(widthClass, 'cursor-pointer justify-between', className)}
			>
				{selectedLabel || placeholder}
				<ChevronsUpDownIcon class="opacity-50" />
			</Button>
		{/snippet}
	</Popover.Trigger>

	<Popover.Content class={cn(widthClass, 'p-0')} {id}>
		<div class="p-2">
			<Input placeholder={searchPlaceholder} bind:value={searchValue} class="h-9" />
		</div>

		{#if filteredItems.length === 0}
			<div class="text-muted-foreground border-t py-6 text-center text-sm">No results.</div>
		{:else if filteredItems.length > 20}
			<div class="h-[300px] border-t" data-command-group data-value={groupLabel}>
				<VirtualList items={filteredItems} let:item height="300px">
					<button
						type="button"
						role="option"
						aria-selected={value === item.value}
						class="hover:bg-accent hover:text-accent-foreground relative flex w-full cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
						onclick={() => selectItem(item)}
					>
						<CheckIcon class={cn('h-4 w-4', value !== item.value && 'text-transparent')} />
						{item.label}
					</button>
				</VirtualList>
			</div>
		{:else}
			<div class="border-t" data-command-group data-value={groupLabel}>
				{#each filteredItems as item (item.value)}
					<button
						type="button"
						role="option"
						aria-selected={value === item.value}
						class="hover:bg-accent hover:text-accent-foreground relative flex w-full cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
						onclick={() => selectItem(item)}
					>
						<CheckIcon class={cn('h-4 w-4', value !== item.value && 'text-transparent')} />
						{item.label}
					</button>
				{/each}
			</div>
		{/if}
	</Popover.Content>
</Popover.Root>
