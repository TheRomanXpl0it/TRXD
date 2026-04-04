<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';

	interface CategoryTotal {
		category: string;
		count: number;
	}

	interface Props {
		solves?: any[];
		totalChallenges?: CategoryTotal[];
	}

	let { solves = [], totalChallenges = [] }: Props = $props();

	let chartEl = $state<HTMLDivElement>();
	let chart = $state<any>();

	// Process data for the radar chart
	const processedData = $derived.by(() => {
		if (!totalChallenges || totalChallenges.length === 0) return { labels: [], series: [] };

		const solveCounts: Record<string, number> = {};
		solves.forEach((s) => {
			if (s.category) {
				solveCounts[s.category] = (solveCounts[s.category] || 0) + 1;
			}
		});

		const labels = totalChallenges.map((tc) => tc.category);
		const seriesData = totalChallenges.map((tc) => {
			const count = solveCounts[tc.category] || 0;
			const total = tc.count || 1;
			return Math.round((count / total) * 100);
		});

		return { labels, series: seriesData };
	});

	$effect(() => {
		if (!browser || !chartEl || processedData.labels.length === 0) return;

		import('apexcharts').then((mod) => {
			const ApexCharts = mod.default;

			const isDark = document.documentElement.classList.contains('dark');
			const primaryColor = 'hsl(217.2 91.2% 59.8%)'; // Standard blue, adjust if needed

			const options = {
				series: [
					{
						name: 'Completion',
						data: processedData.series
					}
				],
				chart: {
					height: 350,
					type: 'radar',
					toolbar: { show: false },
					background: 'transparent',
					dropShadow: {
						enabled: true,
						blur: 8,
						left: 1,
						top: 1,
						opacity: 0.2
					}
				},
				colors: [primaryColor],
				plotOptions: {
					radar: {
						size: 110,
						polygons: {
							strokeColors: isDark ? '#334155' : '#e2e8f0',
							strokeWidth: '1px',
							connectorColors: isDark ? '#334155' : '#e2e8f0',
							fill: {
								colors: isDark ? ['#1e293b', '#0f172a'] : ['#f8fafc', '#fff']
							}
						}
					}
				},
				stroke: {
					width: 3,
					curve: 'smooth'
				},
				fill: {
					opacity: 0.4
				},
				markers: {
					size: 5,
					colors: ['#fff'],
					strokeColors: primaryColor,
					strokeWidth: 3
				},
				labels: processedData.labels,
				xaxis: {
					labels: {
						show: true,
						style: {
							colors: isDark ? '#94a3b8' : '#64748b',
							fontSize: '11px',
							fontWeight: 800,
							fontFamily: 'inherit'
						}
					}
				},
				yaxis: {
					show: false,
					min: 0,
					max: 100,
					tickAmount: 5
				},
				tooltip: {
					theme: isDark ? 'dark' : 'light',
					y: {
						formatter: (val: number) => val + '%'
					}
				}
			};

			if (chart) chart.destroy();
			chart = new ApexCharts(chartEl, options);
			chart.render();
		});

		return () => {
			if (chart) chart.destroy();
		};
	});
</script>

<div class="flex min-h-[350px] w-full items-center justify-center p-4">
	{#if processedData.labels.length > 0}
		<div bind:this={chartEl} class="w-full"></div>
	{:else}
		<div class="text-muted-foreground font-mono text-xs uppercase tracking-widest opacity-50">
			The competition hasn't started yet
		</div>
	{/if}
</div>
