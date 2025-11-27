<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import EmptyState from '$lib/components/ui/empty-state.svelte';
	import Icon from '@iconify/svelte';
	import { Search } from '@lucide/svelte';
	import { getCountryItems, filterCountries, type CountryItem } from '$lib/utils/countries';

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
				<Icon icon={`circle-flags:${selectedCountry.iso2.toLowerCase()}`} width="16" height="16" />
				<span class="uppercase">{value}</span>
			</span>
		{:else}
			<span class="text-muted-foreground">{placeholder}</span>
		{/if}
	</Select.Trigger>
	<Select.Content sideOffset={4}>
		<div class="px-2 pb-2">
			<Input placeholder="Search countries..." bind:value={countrySearch} class="h-8" />
		</div>
		{#each filteredCountries as item (item.value)}
			<Select.Item value={item.value} label={item.label}>
				{#snippet children({ selected })}
					<Icon icon={`circle-flags:${item.iso2.toLowerCase()}`} width="16" height="16" />
					<span>{item.label} ({item.value})</span>
				{/snippet}
			</Select.Item>
		{/each}
		{#if filteredCountries.length === 0}
			<EmptyState icon={Search} title="No countries found" />
		{/if}
	</Select.Content>
</Select.Root>
