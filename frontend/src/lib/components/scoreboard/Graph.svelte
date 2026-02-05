<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Chart, Svg, Axis, Spline, Highlight, Tooltip, Points } from 'layerchart';
	import { scaleTime, scaleLinear } from 'd3-scale';
	import { ChartLine } from '@lucide/svelte';

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

	// Detect dark mode
	let isDark = $state(false);

	$effect(() => {
		isDark = document.documentElement.classList.contains('dark');
	});

	// Parse ISO timestamps to Date
	function toDate(t: any): Date | null {
		if (t instanceof Date) return t;
		if (typeof t === 'number') {
			return new Date(t < 2_000_000_000 ? t * 1000 : t);
		}
		if (typeof t !== 'string') return null;
		const fixed = t.replace(/\.(\d{3})\d*(Z|[+\-]\d\d:\d\d)$/, '.$1$2');
		const d = new Date(fixed);
		return isNaN(d.getTime()) ? null : d;
	}

	function normalizeTeam(entry: any): any[] {
		if (!entry) return [];
		if (Array.isArray(entry.submissions)) {
			return entry.submissions.map((s: any) => ({
				points: Number(s?.score ?? 0),
				date: toDate(s?.timestamp),
				fb: !!s?.first_blood
			}));
		}
		if (Array.isArray(entry.solves)) {
			return entry.solves.map((s: any) => ({
				points: Number(s?.[1] ?? 0),
				date: toDate(s?.[2]),
				fb: !!s?.[3]
			}));
		}
		return [];
	}

	function totalPoints(entry: any): number {
		const solves = normalizeTeam(entry);
		if (solves.length === 0) return 0;
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

	const chartData = $derived.by(() => {
		const arr = Array.isArray(data) ? data : [];
		const n = Number(topN ?? 5) || 5;

		const ranked = [...arr]
			.map((e: any) => ({ ...e, total: totalPoints(e) }))
			.sort((a: any, b: any) => (b.total || 0) - (a.total || 0))
			.slice(0, n);

		const series: Array<{
			name: string;
			data: Array<{ date: Date; value: number; fb?: boolean }>;
			color: string;
		}> = [];

		ranked.forEach((team: any, index: number) => {
			const name = nameForTeam(team);
			const solves = normalizeTeam(team)
				.filter((s: any) => s.date !== null)
				.sort((a: any, b: any) => a.date!.getTime() - b.date!.getTime());

			const points: Array<{ date: Date; value: number; fb?: boolean }> = [];

			// Add starting point at 0
			if (solves.length > 0 && solves[0].date) {
				points.push({
					date: new Date(solves[0].date.getTime() - 60000),
					value: 0
				});
			}

			// Add cumulative points for each solve
			for (const s of solves) {
				if (s.date) {
					points.push({
						date: s.date,
						value: Number(s?.points ?? 0),
						fb: !!s?.fb
					});
				}
			}

			// Extend line to current time
			if (points.length > 0) {
				const lastPoint = points[points.length - 1];
				const now = new Date();
				if (lastPoint.date < now) {
					points.push({ date: now, value: lastPoint.value });
				}
			}

			if (points.length > 0) {
				series.push({
					name,
					data: points,
					color: colors[index % colors.length]
				});
			}
		});

		return series;
	});

	const textColor = $derived(
		isDark ? 'hsl(var(--muted-foreground))' : 'hsl(var(--muted-foreground))'
	);
	const gridColor = $derived(isDark ? 'hsl(var(--border))' : 'hsl(var(--border))');

	// Calculate time range and appropriate tick count
	const timeRangeInfo = $derived.by(() => {
		if (chartData.length === 0) return { ticks: 6, format: (d: Date) => d.toLocaleDateString() };

		const allDates = chartData.flatMap((s) => s.data.map((p) => p.date.getTime()));
		const minTime = Math.min(...allDates);
		const maxTime = Math.max(...allDates);
		const rangeMs = maxTime - minTime;
		const rangeHours = rangeMs / (1000 * 60 * 60);

		if (rangeHours <= 12) {
			// < 12 hours: show every 2 hours
			return {
				ticks: Math.ceil(rangeHours / 2),
				format: (d: Date) =>
					d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false })
			};
		} else if (rangeHours <= 48) {
			// 12-48 hours: show every 4-6 hours
			return {
				ticks: Math.ceil(rangeHours / 5),
				format: (d: Date) =>
					d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false })
			};
		} else if (rangeHours <= 168) {
			// 48h-7 days: show every 12 hours
			return {
				ticks: Math.ceil(rangeHours / 12),
				format: (d: Date) => {
					const day = d.getDate();
					const time = d.toLocaleTimeString([], {
						hour: '2-digit',
						minute: '2-digit',
						hour12: false
					});
					return `${day}/${time}`;
				}
			};
		} else {
			// > 7 days: show daily
			return {
				ticks: Math.ceil(rangeHours / 24),
				format: (d: Date) => d.toLocaleDateString([], { month: 'short', day: 'numeric' })
			};
		}
	});
</script>

<Card.Root class="w-full">
	<Card.Header>
		<div class="flex items-center gap-2">
			<ChartLine class="text-muted-foreground h-5 w-5" />
			<Card.Title>Score Progression</Card.Title>
		</div>
		<Card.Description>
			Top {topN}
			{userMode ? 'Players' : 'Teams'}
		</Card.Description>
	</Card.Header>

	<Card.Content class="pl-0">
		<div class="h-96 w-full md:p-4">
			{#if chartData.length > 0}
				<Chart
					data={chartData.flatMap((s) => s.data)}
					x="date"
					xScale={scaleTime()}
					y="value"
					yScale={scaleLinear()}
					yDomain={[0, null]}
					padding={{ left: 16, bottom: 24 }}
				>
					<Svg>
						<Axis
							placement="left"
							grid={{ style: `stroke: ${gridColor}` }}
							rule={{ style: `font-size: 11px; fill: ${textColor}` }}
						/>
						<Axis
							placement="bottom"
							format={timeRangeInfo.format}
							ticks={timeRangeInfo.ticks}
							grid={{ style: `stroke: ${gridColor}` }}
							rule={{ style: `font-size: 11px; fill: ${textColor}` }}
						/>
						{#each chartData as series}
							<Spline data={series.data} class="stroke-2" style={`stroke: ${series.color}`} />
							<Points data={series.data} r={4} fill={series.color} />
						{/each}
						<Highlight points lines />
					</Svg>
				</Chart>
			{:else}
				<div class="text-muted-foreground flex h-full items-center justify-center">
					No data available
				</div>
			{/if}
		</div>

		<!-- Legend at bottom -->
		{#if chartData.length > 0}
			<div class="mt-4 flex flex-wrap justify-center gap-3 px-6 text-sm">
				{#each chartData as series}
					<div class="flex items-center gap-1.5">
						<div
							class="h-3 w-3 rounded-full shadow-sm"
							style="background-color: {series.color}"
						></div>
						<span class="text-muted-foreground text-xs font-medium">{series.name}</span>
					</div>
				{/each}
			</div>
		{/if}
	</Card.Content>
</Card.Root>
