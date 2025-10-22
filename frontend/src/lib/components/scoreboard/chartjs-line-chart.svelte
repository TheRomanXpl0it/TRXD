<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import {
		Chart,
		LineElement,
		PointElement,
		LinearScale,
		TimeScale,
		Title,
		Tooltip,
		Legend,
		CategoryScale
	} from 'chart.js';
	import 'chartjs-adapter-date-fns';

	// Register Chart.js components
	Chart.register(
		LineElement,
		PointElement,
		LinearScale,
		TimeScale,
		Title,
		Tooltip,
		Legend,
		CategoryScale
	);

	// Props interface
	interface Props {
		data?: any[];
		topN?: number;
		timeMin?: number;
		timeMax?: number;
		teamNames?: Record<string, string>;
	}

	let { data = [], topN = 5, timeMin, timeMax, teamNames }: Props = $props();

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
		let sum = 0;
		for (const s of normalizeTeam(entry)) {
			sum += Number(s?.points ?? 0);
		}
		return sum;
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
		'#6366f1'  // indigo
	];

	function buildChartData() {
		const arr = Array.isArray(data) ? data : [];
		const n = Number(topN ?? 5) || 5;

		const ranked = [...arr]
			.map((e: any) => ({ ...e, total: totalPoints(e) }))
			.sort((a: any, b: any) => (b.total || 0) - (a.total || 0))
			.slice(0, n);

		const datasets: any[] = [];

		ranked.forEach((team: any, index: number) => {
			const name = nameForTeam(team);
			const solves = normalizeTeam(team)
				.filter((s: any) => Number.isFinite(s?.ts))
				.sort((a: any, b: any) => a.ts - b.ts);

			let cum = 0;
			const points: any[] = [];

			// Add starting point
			if (solves.length > 0) {
				points.push({ x: solves[0].ts, y: 0 });
			}

			// Add cumulative points
			for (const s of solves) {
				cum += Number(s?.points ?? 0);
				points.push({ x: s.ts, y: cum });
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
				pointRadius: 4,
				pointHoverRadius: 6,
				pointBackgroundColor: color,
				pointBorderColor: '#ffffff',
				pointBorderWidth: 2
			});
		});

		return { datasets };
	}

	function createChart() {
		if (!canvas) return;

		const chartData = buildChartData();

		const config = {
			type: 'line' as const,
			data: chartData,
			options: {
				responsive: true,
				maintainAspectRatio: false,
				scales: {
					x: {
						type: 'time' as const,
						time: {
							displayFormats: {
								hour: 'MMM dd HH:mm',
								day: 'MMM dd',
								week: 'MMM dd',
								month: 'MMM yyyy'
							},
							tooltipFormat: 'MMM dd, yyyy HH:mm'
						},
						title: {
							display: true,
							text: 'Time'
						},
						grid: {
							color: 'rgba(0, 0, 0, 0.1)'
						}
					},
					y: {
						beginAtZero: true,
						title: {
							display: true,
							text: 'Points'
						},
						grid: {
							color: 'rgba(0, 0, 0, 0.1)'
						}
					}
				},
				plugins: {
					legend: {
						position: 'top' as const,
						labels: {
							usePointStyle: true,
							padding: 20
						}
					},
					tooltip: {
						mode: 'index' as const,
						intersect: false,
						callbacks: {
							title: function(context: any) {
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
							label: function(context: any) {
								return `${context.dataset.label}: ${context.parsed.y} pts`;
							}
						}
					}
				},
				interaction: {
					mode: 'index' as const,
					intersect: false
				},
				elements: {
					point: {
						hoverRadius: 8
					}
				}
			}
		};

		chart = new Chart(canvas, config);
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

<div class="w-full h-96">
	<canvas bind:this={canvas}></canvas>
</div>
