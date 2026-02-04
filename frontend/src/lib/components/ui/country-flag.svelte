<script lang="ts">
	// Load all flag SVGs using Vite's glob import
	const flags = import.meta.glob('/node_modules/country-flag-icons/3x2/*.svg', {
		query: '?url',
		eager: true
	}) as Record<string, { default: string }>;

	const getFlagUrl = (iso: string) => {
		// Handle both ISO2 and ISO3 codes
		const code = iso?.toUpperCase() || '';
		const key = `/node_modules/country-flag-icons/3x2/${code}.svg`;
		return flags[key]?.default || '';
	};

	let {
		country,
		class: className = '',
		width = 32,
		height = 32
	} = $props<{
		country: string;
		class?: string;
		width?: number | string;
		height?: number | string;
	}>();

	const flagUrl = $derived(getFlagUrl(country));
</script>

{#if flagUrl}
	<img src={flagUrl} alt={`${country} flag`} class={className} {width} {height} />
{:else}
	<!-- Fallback for missing flags -->
	<div class="bg-muted {className}" style="width: {width}px; height: {height}px;"></div>
{/if}
