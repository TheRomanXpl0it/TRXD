<script lang="ts">
	import { Button } from '@/components/ui/button';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { X, Filter, Shapes, LayoutGrid, List } from '@lucide/svelte';

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

	function clearFilters() {
		filterCategories = [];
		filterTags = [];
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
					aria-label="Filter by categories"
				>
					<Shapes class="h-4 w-4" aria-hidden="true" />
					<span class="hidden sm:inline">Categories</span>
					{#if filterCategories.length > 0}
						<span
							class="bg-primary text-primary-foreground rounded px-1.5 py-0.5 text-xs"
							aria-label="{filterCategories.length} categories selected"
						>
							{filterCategories.length}
						</span>
					{/if}
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-[260px] p-1">
			<Command.Root>
				<Command.Input
					placeholder="Search categories..."
					class="border-0 shadow-none ring-0 focus:outline-none focus:ring-0 focus-visible:outline-none focus-visible:ring-0"
				/>
				<Command.List>
					<Command.Empty>No categories.</Command.Empty>
					<Command.Group value="categories">
						{#each categories as c (c.value)}
							<Command.Item
								value={c.value}
								onSelect={() => {
									if (filterCategories.includes(c.value)) {
										filterCategories = filterCategories.filter((x) => x !== c.value);
									} else {
										filterCategories = [...filterCategories, c.value];
									}
								}}
							>
								<Checkbox
									checked={filterCategories.includes(c.value)}
									aria-label="Filter by {c.label}"
								/>
								<span class="ml-2">{c.label}</span>
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
		</Popover.Content>
	</Popover.Root>

	<Popover.Root bind:open={tagsOpen}>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					variant="outline"
					class="flex cursor-pointer items-center gap-1 shrink-0"
					aria-label="Filter by tags"
				>
					<Filter class="h-4 w-4" aria-hidden="true" />
					<span class="hidden sm:inline">Tags</span>
					{#if filterTags.length > 0}
						<span
							class="bg-primary text-primary-foreground rounded px-1.5 py-0.5 text-xs"
							aria-label="{filterTags.length} tags selected"
						>
							{filterTags.length}
						</span>
					{/if}
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-[260px] p-1">
			<Command.Root>
				<Command.Input
					placeholder="Search tags..."
					class="border-0 shadow-none ring-0 focus:outline-none focus:ring-0 focus-visible:outline-none focus-visible:ring-0"
				/>
				<Command.List>
					<Command.Empty>No tags.</Command.Empty>
					<Command.Group value="tags">
						{#each allTags as t (t)}
							<Command.Item
								value={t}
								onSelect={() => {
									if (filterTags.includes(t)) {
										filterTags = filterTags.filter((x) => x !== t);
									} else {
										filterTags = [...filterTags, t];
									}
								}}
							>
								<Checkbox checked={filterTags.includes(t)} aria-label="Filter by {t}" />
								<span class="ml-2">{t}</span>
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
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
