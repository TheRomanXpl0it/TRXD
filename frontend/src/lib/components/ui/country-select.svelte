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
	<Select.Content sideOffset={4} class="max-h-[300px]">
		<div class="px-2 pb-2">
			<Input placeholder="Search countries..." bind:value={countrySearch} class="h-8" />
		</div>
		{#each filteredCountries as item (item.value)}
			<Select.Item value={item.value} label={item.label}>
				{#snippet children({ selected })}
					<img src={getFlagUrl(item.iso2)} alt={item.label} class="h-4 w-6 object-cover" />
					<span>{item.label} ({item.value})</span>
				{/snippet}
			</Select.Item>
		{/each}
		{#if filteredCountries.length === 0}
			<EmptyState icon={Search} title="No countries found" />
		{/if}
	</Select.Content>
</Select.Root>
