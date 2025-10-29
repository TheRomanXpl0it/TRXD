<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Chart from 'chart.js/auto';
	import 'chartjs-adapter-date-fns';

	// Chart.js/auto automatically registers all components, no need for manual registration

	// Props interface
	interface Props {
		data?: any[];
		topN?: number;
		timeMin?: number;
		timeMax?: number;
		teamNames?: Record<string, string>;
		userMode?: boolean;
	}

	let { data = [], topN = 5, timeMin, timeMax, teamNames, userMode = false }: Props = $props();

	let canvas: HTMLCanvasElement;
	let chart: Chart | null = null;

	// Parse ISO timestamps to epoch milliseconds
	function toEpochMs(t: any): number {
		if (typeof t === 'number') return t < 2_000_000_000 ? t * 1000 : t;
		if (typeof t !== 'string') return NaN;
		// Handle microseconds in ISO string
		const fixed = t.replace(/\.(\d{3})\d*(Z|[+\-]\d\d:\d\d)$/, '.$1$2');
		const n = Date.parse(fixed);
		return Number.isFinite(n) ? n : NaN;
	}

	function normalizeTeam(entry: any): any[] {
		if (!entry) return [];
		if (Array.isArray(entry.submissions)) {
			return entry.submissions.map((s: any) => ({
				points: Number(s?.score ?? 0),
				ts: toEpochMs(s?.timestamp),
				fb: !!s?.first_blood
			}));
		}
		if (Array.isArray(entry.solves)) {
			// [solve_id, points, timestamp, first_blood]
			return entry.solves.map((s: any) => ({
				points: Number(s?.[1] ?? 0),
				ts: toEpochMs(s?.[2]),
				fb: !!s?.[3]
			}));
		}
		return [];
	}

	function totalPoints(entry: any): number {
		const solves = normalizeTeam(entry);
		if (solves.length === 0) return 0;
		// Return the highest score (since backend provides total scores)
		return Math.max(...solves.map((s) => Number(s?.points ?? 0)));
	}

	function nameForTeam(team: any): string {
		if (teamNames && typeof teamNames === 'object') {
			const v = teamNames[team?.team_id];
			if (typeof v === 'string' && v.length) return v;
		}
		return String(team?.team_id ?? '');
	}

	// Color palette for different teams
	const colors = [
		'#3b82f6', // blue
		'#ef4444', // red
		'#10b981', // green
		'#f59e0b', // amber
		'#8b5cf6', // violet
		'#06b6d4', // cyan
		'#f97316', // orange
		'#84cc16', // lime
		'#ec4899', // pink
		'#6366f1' // indigo
	];

	function buildChartData() {
		const arr = Array.isArray(data) ? data : [];
		const n = Number(topN ?? 5) || 5;

		const ranked = [...arr]
			.map((e: any) => ({ ...e, total: totalPoints(e) }))
			.sort((a: any, b: any) => (b.total || 0) - (a.total || 0))
			.slice(0, n);

		//console.log('Ranked teams:', $state.snapshot(ranked));

		const datasets: any[] = [];

		ranked.forEach((team: any, index: number) => {
			const name = nameForTeam(team);
			const solves = normalizeTeam(team)
				.filter((s: any) => Number.isFinite(s?.ts))
				.sort((a: any, b: any) => a.ts - b.ts);

			//console.log(`Team ${name} solves:`, $state.snapshot(solves));

			const points: any[] = [];

			// Add starting point at 0
			if (solves.length > 0) {
				points.push({ x: new Date(solves[0].ts - 60000), y: 0 });
			}

			// Use raw score values (backend already provides total scores)
			// Add cumulative points for each solve
			for (const s of solves) {
				points.push({ x: new Date(s.ts), y: Number(s?.points ?? 0) });
			}

			// Extend line to current time to show score persistence
			if (points.length > 0) {
				const lastPoint = points[points.length - 1];
				const now = new Date();
				if (lastPoint.x < now) {
					points.push({ x: now, y: lastPoint.y });
				}
			}


			if (points.length === 0) {
				//console.log(`No points for team ${name}, skipping`);
				return;
			}

			const color = colors[index % colors.length];

			datasets.push({
				label: name,
				data: points,
				borderColor: color,
				backgroundColor: color + '20', // Add transparency
				borderWidth: 3,
				fill: false,
				tension: 0.1,
				pointRadius: 6,
				pointHoverRadius: 8,
				pointBackgroundColor: color,
				pointBorderColor: '#ffffff',
				pointBorderWidth: 2
			});
		});

		//console.log('Final datasets:', $state.snapshot(datasets));
		return { datasets };
	}

	function createChart() {
		if (!canvas) return;

		const chartData = buildChartData();
		//console.log('Chart data:', chartData);

		const config: any = {
			type: 'line' as const,
			data: chartData,
			options: {
				responsive: true,
				maintainAspectRatio: false,
				scales: {
					x: {
						type: 'time' as const,
						time: {
							unit: 'minute',
							displayFormats: {
								minute: 'MM/dd HH:mm',
								hour: 'MM/dd HH:mm',
								day: 'MM/dd',
								week: 'MM/dd',
								month: 'MMM yyyy'
							},
							tooltipFormat: 'yyyy-MM-dd HH:mm'
						},
						title: {
							display: true,
							text: 'Time',
							color: '#374151'
						},
						grid: {
							color: 'rgba(0, 0, 0, 0.1)'
						},
						ticks: {
							color: '#374151',
							maxTicksLimit: 8,
							source: 'data'
						}
					},
					y: {
						beginAtZero: true,
						title: {
							display: true,
							text: 'Points',
							color: '#374151'
						},
						grid: {
							color: 'rgba(0, 0, 0, 0.1)'
						},
						ticks: {
							color: '#374151'
						}
					}
				},
				plugins: {
					legend: {
						position: 'top' as const,
						labels: {
							usePointStyle: true,
							padding: 20,
							color: '#374151'
						}
					},
					tooltip: {
						mode: 'point' as const,
						intersect: true,
						backgroundColor: 'rgba(255, 255, 255, 0.95)',
						titleColor: '#374151',
						bodyColor: '#374151',
						borderColor: '#d1d5db',
						borderWidth: 1,

						callbacks: {
							title: function (context: any) {
								if (context.length > 0) {
									const date = new Date(context[0].parsed.x);
									return date.toLocaleString(undefined, {
										year: 'numeric',
										month: 'short',
										day: '2-digit',
										hour: '2-digit',
										minute: '2-digit'
									});
								}
								return '';
							},
							label: function (context: any) {
								return `${context.dataset.label}: ${context.parsed.y} pts`;
							}
						}
					}
				},
				interaction: {
					mode: 'point' as const,
					intersect: true
				},
				elements: {
					point: {
						hoverRadius: 10
					}
				}
			}
		};

		chart = new Chart(canvas, config as any);
	}

	function updateChart() {
		if (!chart) return;

		const chartData = buildChartData();
		chart.data = chartData;
		chart.update();
	}

	onMount(() => {
		createChart();
	});

	onDestroy(() => {
		if (chart) {
			chart.destroy();
			chart = null;
		}
	});

	// Watch for data changes
	$effect(() => {
		// Reactive dependencies
		void data;
		void topN;
		void teamNames;

		if (chart) {
			updateChart();
		}
	});
</script>

<div class="w-full p-3">
	<div class="mb-2 flex items-center justify-between">
		<h3 class="text-lg font-semibold">Top {topN} {userMode ? 'Players' : 'Teams'} — Score Chart</h3>
	</div>

	<div class="h-96 w-full p-4">
		<canvas bind:this={canvas}></canvas>
	</div>
</div>
