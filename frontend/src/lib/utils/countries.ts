import countries from '$lib/data/countries.json';

export type Country = { name: string; iso2: string; iso3?: string; emoji?: string };

export type CountryItem = {
	value: string;
	label: string;
	iso2: string;
};

export function getCountryItems(): CountryItem[] {
	return (countries as Country[])
		.filter((c) => c.iso3)
		.map((c) => ({
			value: c.iso3!.toUpperCase(),
			label: c.name,
			iso2: c.iso2.toUpperCase()
		}))
		.sort((a, b) => a.label.localeCompare(b.label));
}

export function filterCountries(items: CountryItem[], search: string): CountryItem[] {
	if (!search.trim()) return items.slice(0, 50);

	const lowerSearch = search.toLowerCase();
	return items
		.filter(
			(c) =>
				c.label.toLowerCase().includes(lowerSearch) || c.value.toLowerCase().includes(lowerSearch)
		)
		.slice(0, 50);
}

export function getCountryByIso3(iso3: string): CountryItem | null {
	const items = getCountryItems();
	return items.find((c) => c.value === iso3.toUpperCase()) ?? null;
}

export function getCountryIso2(iso3: string): string | null {
	const country = (countries as Country[]).find(
		(c) => c.iso3?.toUpperCase() === iso3?.toUpperCase()
	);
	return country?.iso2?.toUpperCase() ?? null;
}
