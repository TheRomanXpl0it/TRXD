<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import EmptyState from '$lib/components/ui/empty-state.svelte';
	import { Search } from '@lucide/svelte';
	import { getCountryItems, filterCountries, type CountryItem } from '$lib/utils/countries';

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
		id
	} = $props<{
		value?: string;
		placeholder?: string;
		id?: string;
	}>();

	const countryItems = getCountryItems();
	let countrySearch = $state('');

	const filteredCountries = $derived(filterCountries(countryItems, countrySearch));
	const selectedCountry = $derived(countryItems.find((c) => c.value === value));
</script>

<Select.Root
	type="single"
	bind:value
	items={filteredCountries}
	allowDeselect={true}
	onOpenChange={(open) => {
		if (!open) countrySearch = '';
	}}
>
	<Select.Trigger {id} class="w-full justify-between">
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
	</Select.Trigger>
	<Select.Content sideOffset={4} class="flex max-h-[300px] flex-col overflow-hidden">
		<div class="border-b px-2 py-2">
			<Input
				placeholder="Search countries..."
				bind:value={countrySearch}
				class="h-8 focus-visible:ring-0"
			/>
		</div>
		<div class="flex-1 overflow-y-auto p-1">
			{#each filteredCountries as item (item.value)}
				<Select.Item value={item.value} label={item.label} class="cursor-pointer">
					{#snippet children({ selected })}
						<div class="flex w-full items-center gap-2">
							<img
								src={getFlagUrl(item.iso2)}
								alt={item.label}
								class="h-3.5 w-5 rounded-[2px] object-cover shadow-sm"
							/>
							<span class="truncate">{item.label}</span>
							<span class="text-muted-foreground ml-auto text-xs opacity-50">{item.value}</span>
						</div>
					{/snippet}
				</Select.Item>
			{/each}
			{#if filteredCountries.length === 0}
				<div class="text-muted-foreground flex flex-col items-center justify-center py-6">
					<Search class="mb-2 h-8 w-8 opacity-20" />
					<p class="text-xs">No countries found</p>
				</div>
			{/if}
		</div>
	</Select.Content>
</Select.Root>
