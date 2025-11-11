<script lang="ts">
	import { Button } from '@/components/ui/button';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { X, Filter, Shapes, LayoutGrid, List } from '@lucide/svelte';
	import VirtualList from '@sveltejs/svelte-virtual-list';

	let {
		search = $bindable(''),
		filterCategories = $bindable([]),
		filterTags = $bindable([]),
		categories = [],
		allTags = [],
		compactView = $bindable(false),
		activeFiltersCount = 0
	}: {
		search: string;
		filterCategories: string[];
		filterTags: string[];
		categories: Array<{ value: string; label: string }>;
		allTags: string[];
		compactView: boolean;
		activeFiltersCount: number;
	} = $props();

	let tagsOpen = $state(false);
	let categoriesOpen = $state(false);
	let categorySearch = $state('');
	let tagSearch = $state('');

	const filteredCategories = $derived(
		categorySearch
			? categories.filter(c => 
				c.label.toLowerCase().includes(categorySearch.toLowerCase())
			)
			: categories
	);

	const filteredTags = $derived(
		tagSearch
			? allTags.filter(t => 
				t.toLowerCase().includes(tagSearch.toLowerCase())
			)
			: allTags
	);

	function clearFilters() {
		filterCategories = [];
		filterTags = [];
	}
	
	function toggleCategory(value: string) {
		const idx = filterCategories.indexOf(value);
		if (idx > -1) {
			filterCategories.splice(idx, 1);
			filterCategories = filterCategories;
		} else {
			filterCategories.push(value);
			filterCategories = filterCategories;
		}
	}
	
	function toggleTag(value: string) {
		const idx = filterTags.indexOf(value);
		if (idx > -1) {
			filterTags.splice(idx, 1);
			filterTags = filterTags;
		} else {
			filterTags.push(value);
			filterTags = filterTags;
		}
	}
</script>

<div class="mb-6 flex items-center gap-2">
	<div class="relative flex-1 min-w-0">
		<label for="search-challenges" class="sr-only">Search challenges</label>
		<Input
			id="search-challenges"
			placeholder="Search challenges by name, category, or tag"
			bind:value={search}
			class="pr-8"
			aria-label="Search challenges"
		/>
		{#if search}
			<button
				type="button"
				onclick={() => (search = '')}
				class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
				aria-label="Clear search"
			>
				<X class="h-4 w-4" />
			</button>
		{/if}
	</div>

	<Popover.Root bind:open={categoriesOpen}>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					variant="outline"
					class="flex cursor-pointer items-center gap-1 shrink-0"
					aria-label={filterCategories.length > 0
						? `${filterCategories.length} categories selected`
						: "Filter by categories"}
				>
					<Shapes class="h-4 w-4" aria-hidden="true" />
					<span class="hidden sm:inline">Categories</span>
					{#if filterCategories.length > 0}
						<span class="ml-1 flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1 text-xs text-primary-foreground">
							{filterCategories.length}
						</span>
					{/if}
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-[260px] p-1">
			<div class="px-2 py-1.5">
				<Input
					bind:value={categorySearch}
					placeholder="Search categories..."
					class="h-8"
				/>
			</div>
			{#if filteredCategories.length === 0}
				<div class="py-6 text-center text-sm text-muted-foreground">
					No categories found.
				</div>
			{:else if filteredCategories.length > 20}
				<div class="px-1" style="height: 300px;">
					<VirtualList items={filteredCategories} let:item>
						<button
							type="button"
							class="relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
							onclick={() => toggleCategory(item.value)}
						>
							<Checkbox
								checked={filterCategories.includes(item.value)}
								aria-label="Filter by {item.label}"
							/>
							<span class="ml-2">{item.label}</span>
						</button>
					</VirtualList>
				</div>
			{:else}
				<div class="px-1 max-h-[300px] overflow-y-auto">
					{#each filteredCategories as item (item.value)}
						<button
							type="button"
							class="relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
							onclick={() => toggleCategory(item.value)}
						>
							<Checkbox
								checked={filterCategories.includes(item.value)}
								aria-label="Filter by {item.label}"
							/>
							<span class="ml-2">{item.label}</span>
						</button>
					{/each}
				</div>
			{/if}
		</Popover.Content>
	</Popover.Root>

	<Popover.Root bind:open={tagsOpen}>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					variant="outline"
					class="flex cursor-pointer items-center gap-1 shrink-0"
					aria-label={filterTags.length > 0
						? `${filterTags.length} tags selected`
						: "Filter by tags"}
				>
					<Filter class="h-4 w-4" aria-hidden="true" />
					<span class="hidden sm:inline">Tags</span>
					{#if filterTags.length > 0}
						<span class="ml-1 flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1 text-xs text-primary-foreground">
							{filterTags.length}
						</span>
					{/if}
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-[260px] p-1">
			<div class="px-2 py-1.5">
				<Input
					bind:value={tagSearch}
					placeholder="Search tags..."
					class="h-8"
				/>
			</div>
			{#if filteredTags.length === 0}
				<div class="py-6 text-center text-sm text-muted-foreground">
					No tags found.
				</div>
			{:else if filteredTags.length > 20}
				<div class="px-1" style="height: 300px;">
					<VirtualList items={filteredTags} let:item>
						<button
							type="button"
							class="relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
							onclick={() => toggleTag(item)}
						>
							<Checkbox checked={filterTags.includes(item)} aria-label="Filter by {item}" />
							<span class="ml-2">{item}</span>
						</button>
					</VirtualList>
				</div>
			{:else}
				<div class="px-1 max-h-[300px] overflow-y-auto">
					{#each filteredTags as item}
						<button
							type="button"
							class="relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
							onclick={() => toggleTag(item)}
						>
							<Checkbox checked={filterTags.includes(item)} aria-label="Filter by {item}" />
							<span class="ml-2">{item}</span>
						</button>
					{/each}
				</div>
			{/if}
		</Popover.Content>
	</Popover.Root>

	{#if activeFiltersCount > 0}
		<Button
			variant="ghost"
			size="icon"
			class="cursor-pointer shrink-0"
			onclick={clearFilters}
			aria-label="Clear all filters ({activeFiltersCount} active)"
		>
			<X class="h-4 w-4" aria-hidden="true" />
		</Button>
	{/if}

	<div class="ml-auto flex items-center gap-1" role="group" aria-label="View mode">
		<Button
			variant={compactView ? 'ghost' : 'secondary'}
			size="icon"
			class="cursor-pointer shrink-0"
			onclick={() => (compactView = false)}
			aria-label="Grid view"
			aria-pressed={!compactView}
		>
			<LayoutGrid class="h-4 w-4" aria-hidden="true" />
		</Button>
		<Button
			variant={compactView ? 'secondary' : 'ghost'}
			size="icon"
			class="cursor-pointer shrink-0"
			onclick={() => (compactView = true)}
			aria-label="Compact view"
			aria-pressed={compactView}
		>
			<List class="h-4 w-4" aria-hidden="true" />
		</Button>
	</div>
</div>
