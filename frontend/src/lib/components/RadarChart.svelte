<script lang="ts">
	import { onMount } from 'svelte';

	const browser = typeof window !== 'undefined';

	let { solves = [] } = $props<{ solves?: any[] }>();

	let chart: any;
	let chartEl: HTMLDivElement;

	// Handle null prop explicit (default value only handles undefined)
	const safeSolves = $derived(Array.isArray(solves) ? solves : []);

	const categories = $derived.by(() => {
		const cats = new Set<string>();
		// Process solves to get all categories
		safeSolves.forEach((s: any) => {
			if (s.category) cats.add(s.category);
		});
		return Array.from(cats).sort();
	});

	const seriesData = $derived.by(() => {
		if (!categories.length) return [0, 0, 0, 0, 0]; // dummy

		// Sum points per category
		const sums: Record<string, number> = {};
		categories.forEach((c) => (sums[c] = 0));

		safeSolves.forEach((s: any) => {
			if (s.category) {
				sums[s.category] += Number(s.points || 0);
			}
		});

		return categories.map((c) => sums[c]);
	});

	$effect(() => {
		if (!browser || !chartEl) return;

		// Dynamic import to avoid SSR issues
		import('apexcharts').then((mod) => {
			const ApexCharts = mod.default;

			const options = {
				series: [
					{
						name: 'Points',
						data: seriesData.length ? seriesData : [0]
					}
				],
				chart: {
					height: 250,
					type: 'radar',
					toolbar: { show: false },
					parentHeightOffset: 0,
					background: 'transparent'
				},
				labels: categories.length ? categories : ['None'],
				title: {
					text: undefined
				},
				stroke: {
					width: 2,
					colors: ['#3b82f6'] // blue-500
				},
				fill: {
					opacity: 0.2,
					colors: ['#3b82f6'] // blue-500
				},
				markers: {
					size: 4,
					colors: ['#3b82f6'],
					strokeColors: '#fff',
					strokeWidth: 2
				},
				yaxis: {
					show: false
				},
				xaxis: {
					labels: {
						style: {
							colors: [], // inherit
							fontSize: '12px',
							fontFamily: 'inherit',
							fontWeight: 600
						}
					}
				},
				theme: {
					mode: document.documentElement.classList.contains('dark') ? 'dark' : 'light',
					palette: 'palette1'
				},
				plotOptions: {
					radar: {
						polygons: {
							strokeColors: 'var(--border)', // subtle grid
							connectorColors: 'var(--border)'
						}
					}
				}
			};

			// Update theme dynamically
			const observer = new MutationObserver(() => {
				const isDark = document.documentElement.classList.contains('dark');
				if (chart) {
					chart.updateOptions({
						theme: { mode: isDark ? 'dark' : 'light' },
						plotOptions: {
							radar: {
								polygons: {
									strokeColors: isDark ? '#374151' : '#e5e7eb',
									connectorColors: isDark ? '#374151' : '#e5e7eb'
								}
							}
						}
					});
				}
			});
			observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });

			// Initial chart
			if (chart) chart.destroy();
			// Only render if we have categories (even if dummy data generated for derived, we check explicit length or solves)
			if (categories.length > 0) {
				chart = new ApexCharts(chartEl, options);
				chart.render();
			}

			return () => {
				observer.disconnect();
				if (chart) chart.destroy();
			};
		});
	});
</script>

<div class="flex h-[250px] w-full items-center justify-center p-2" bind:this={chartEl}>
	{#if safeSolves.length === 0}
		<p class="text-muted-foreground w-full text-center text-sm">No solve data available</p>
	{/if}
</div>
