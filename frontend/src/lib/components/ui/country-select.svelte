<script lang="ts">
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { cn } from '$lib/utils.js';
	import VirtualList from '@sveltejs/svelte-virtual-list';
	import CheckIcon from '@lucide/svelte/icons/check';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import { Search } from '@lucide/svelte';
	import { getCountryItems, filterCountries, type CountryItem } from '$lib/utils/countries';
	import { tick } from 'svelte';

	// Load all flag SVGs using Vite's glob import
	const flags = import.meta.glob('/node_modules/country-flag-icons/3x2/*.svg', {
		query: '?url',
		eager: true
	}) as Record<string, { default: string }>;

	const getFlagUrl = (iso2: string) => {
		const key = `/node_modules/country-flag-icons/3x2/${iso2.toUpperCase()}.svg`;
		return flags[key]?.default || '';
	};

	let {
		value = $bindable(''),
		placeholder = 'Select country',
		id,
		className = ''
	} = $props<{
		value?: string;
		placeholder?: string;
		id?: string;
		className?: string;
	}>();

	const countryItems = getCountryItems();
	let countrySearch = $state('');
	let open = $state(false);
	let triggerRef = $state<HTMLButtonElement>(null!);

	const filteredCountries = $derived.by(() => {
		const items = filterCountries(countryItems, countrySearch);
		// Pre-calculate URLs to avoid re-computation during virtual scrolling
		return items.map((c) => ({
			...c,
			flagUrl: getFlagUrl(c.iso2)
		}));
	});
	const selectedCountry = $derived(countryItems.find((c) => c.value === value));

	function selectItem(item: CountryItem) {
		value = item.value;
		open = false;
		countrySearch = '';
		tick().then(() => triggerRef?.focus());
	}
</script>

<Popover.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (!isOpen) countrySearch = '';
	}}
>
	<Popover.Trigger bind:ref={triggerRef}>
		{#snippet child({ props })}
			<Button
				{...props}
				variant="outline"
				role="combobox"
				aria-expanded={open}
				{id}
				class={cn('w-full justify-between font-normal', className)}
			>
				{#if selectedCountry}
					<span class="flex items-center gap-2">
						<img
							src={getFlagUrl(selectedCountry.iso2)}
							alt={selectedCountry.label}
							class="h-4 w-6 object-cover"
						/>
						<span class="uppercase">{value}</span>
					</span>
				{:else}
					<span class="text-muted-foreground">{placeholder}</span>
				{/if}
				<ChevronsUpDownIcon class="ml-2 h-4 w-4 shrink-0 opacity-50" />
			</Button>
		{/snippet}
	</Popover.Trigger>

	<Popover.Content class="w-[300px] p-0" align="start">
		<div class="flex flex-col overflow-hidden">
			<div class="border-b px-3 py-2">
				<div class="flex items-center gap-2">
					<Search class="h-4 w-4 opacity-50" />
					<Input
						placeholder="Search countries..."
						bind:value={countrySearch}
						class="h-8 border-0 p-0 focus-visible:ring-0"
					/>
				</div>
			</div>

			{#if filteredCountries.length === 0}
				<div class="text-muted-foreground flex flex-col items-center justify-center py-6">
					<Search class="mb-2 h-8 w-8 opacity-20" />
					<p class="text-xs">No countries found</p>
				</div>
			{:else if filteredCountries.length > 20}
				<div class="h-[300px] p-1">
					<VirtualList items={filteredCountries} let:item height="300px">
						<button
							type="button"
							role="option"
							aria-selected={value === item.value}
							class="hover:bg-accent hover:text-accent-foreground relative flex w-full cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none"
							onclick={() => selectItem(item)}
						>
							<div class="flex w-full items-center gap-2">
								<div class="flex h-4 w-4 items-center justify-center">
									{#if value === item.value}
										<CheckIcon class="h-4 w-4" />
									{/if}
								</div>
								<img
									src={item.flagUrl}
									alt={item.label}
									loading="lazy"
									class="h-3.5 w-5 rounded-[2px] object-cover shadow-sm"
								/>
								<span class="truncate">{item.label}</span>
								<span class="text-muted-foreground ml-auto text-xs opacity-50">{item.value}</span>
							</div>
						</button>
					</VirtualList>
				</div>
			{:else}
				<div class="max-h-[300px] overflow-y-auto p-1">
					{#each filteredCountries as item (item.value)}
						<button
							type="button"
							role="option"
							aria-selected={value === item.value}
							class="hover:bg-accent hover:text-accent-foreground relative flex w-full cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none"
							onclick={() => selectItem(item)}
						>
							<div class="flex w-full items-center gap-2">
								<div class="flex h-4 w-4 items-center justify-center">
									{#if value === item.value}
										<CheckIcon class="h-4 w-4" />
									{/if}
								</div>
								<img
									src={item.flagUrl}
									alt={item.label}
									loading="lazy"
									class="h-3.5 w-5 rounded-[2px] object-cover shadow-sm"
								/>
								<span class="truncate">{item.label}</span>
								<span class="text-muted-foreground ml-auto text-xs opacity-50">{item.value}</span>
							</div>
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</Popover.Content>
</Popover.Root>
