<script lang="ts">
	import { Chart, Svg, Axis, Spline, Highlight, Tooltip, Points } from 'layerchart';
	import { scaleTime, scaleLinear } from 'd3-scale';
	// @ts-ignore - Ignore missing type declarations for d3-shape
	import { curveLinear } from 'd3-shape';

	// Props interface
	interface Props {
		data?: any[];
		timeMin?: number;
		timeMax?: number;
		teamNames?: Record<string, string>;
		userMode?: boolean;
	}

	let { data = [], timeMin, timeMax, teamNames, userMode = false }: Props = $props();

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

	// Dynamic vibrant colors
	const colors = [
		'#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6', '#06b6d4', '#f97316', '#84cc16', '#ec4899', '#6366f1'
	];

	const chartData = $derived.by(() => {
		const arr = Array.isArray(data) ? data : [];

		const ranked = [...arr]
			.map((e: any) => ({ ...e, total: totalPoints(e) }))
			.sort((a: any, b: any) => (b.total || 0) - (a.total || 0));

		let minTime = Date.now();
		let hasData = false;
		ranked.forEach((team: any) => {
			const solves = normalizeTeam(team).filter((s: any) => s.date !== null);
			if (solves.length > 0 && solves[0].date) {
				hasData = true;
				minTime = Math.min(minTime, solves[0].date.getTime() - 60000); // 1m buffer start
			}
		});

		const nowMs = Date.now();
		// Create 150 dense intervals for pixel-perfect smooth cursor tracking on lines
		const INTERVALS = 150;
		const timeStep = Math.max(1, (nowMs - minTime) / INTERVALS);
		const gridTimes: number[] = [];
		if (hasData) {
			for (let i = 0; i <= INTERVALS; i++) gridTimes.push(minTime + i * timeStep);
		}

		const series: Array<{
			name: string;
			data: Array<{ date: Date; value: number; name?: string; color?: string; }>;
			color: string;
		}> = [];

		ranked.forEach((team: any, index: number) => {
			const name = nameForTeam(team);
			const color = colors[index % colors.length];
			const solves = normalizeTeam(team)
				.filter((s: any) => s.date !== null)
				.sort((a: any, b: any) => a.date!.getTime() - b.date!.getTime());

			if (solves.length === 0) return;

			// Actual sparse solves
			const sparsePoints: Array<{ time: number; value: number }> = [];
			sparsePoints.push({ time: minTime, value: 0 }); 

			for (const s of solves) {
				sparsePoints.push({ time: s.date.getTime(), value: Number(s.points ?? 0) });
			}
			sparsePoints.push({ time: nowMs, value: sparsePoints[sparsePoints.length - 1].value });

			// Superimpose active solve moments onto the dense tracking grid
			const allT = [...new Set([...gridTimes, ...sparsePoints.map(sp => sp.time)])].sort((a, b) => a - b);
			const points: Array<{ date: Date; value: number; name: string; color: string; }> = [];

			for (const t of allT) {
				let val = 0;
				for (let i = 0; i < sparsePoints.length; i++) {
					if (sparsePoints[i].time === t) {
						val = sparsePoints[i].value;
						break;
					}
					// Linear interpolation
					if (i < sparsePoints.length - 1 && sparsePoints[i].time < t && sparsePoints[i + 1].time > t) {
						const p1 = sparsePoints[i];
						const p2 = sparsePoints[i + 1];
						const ratio = (t - p1.time) / (p2.time - p1.time);
						val = p1.value + ratio * (p2.value - p1.value);
						break;
					}
				}
				points.push({
					date: new Date(t),
					value: Math.round(val),
					name,
					color
				});
			}

			series.push({ name, data: points, color });
		});

		return series;
	});

	const textColor = $derived('hsl(var(--muted-foreground))');
	const gridColor = $derived('hsl(var(--border))');

	const timeRangeInfo = $derived.by(() => {
		if (chartData.length === 0) return { ticks: 6, format: (d: Date) => d.toLocaleDateString() };

		const allDates = chartData.flatMap((s) => s.data.map((p) => p.date.getTime()));
		const minTime = Math.min(...allDates);
		const maxTime = Math.max(...allDates);
		const rangeHours = (maxTime - minTime) / (1000 * 60 * 60);

		if (rangeHours <= 12) {
			return { ticks: Math.ceil(rangeHours / 2), format: (d: Date) => d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false }) };
		} else if (rangeHours <= 48) {
			return { ticks: Math.ceil(rangeHours / 5), format: (d: Date) => d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false }) };
		} else if (rangeHours <= 168) {
			return { ticks: Math.ceil(rangeHours / 12), format: (d: Date) => `${d.getDate()}/${d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false })}` };
		} else {
			return { ticks: Math.ceil(rangeHours / 24), format: (d: Date) => d.toLocaleDateString([], { month: 'short', day: 'numeric' }) };
		}
	});
</script>

<div class="w-full flex-col flex h-[480px]">
	<div class="flex-grow w-full">
		{#if chartData.length > 0}
			<Chart
				data={chartData.flatMap((s) => s.data)}
				x="date"
				xScale={scaleTime()}
				y="value"
				yScale={scaleLinear()}
				yDomain={[0, null]}
				padding={{ left: 24, bottom: 24, top: 12, right: 12 }}
				tooltip={{ mode: 'voronoi' }}
			>
				<Svg>
					<Axis
						placement="left"
						grid={{ style: `stroke: ${gridColor}; stroke-dasharray: 4` }}
						rule={{ style: `font-size: 11px; fill: ${textColor}; font-weight: bold;` }}
					/>
					<Axis
						placement="bottom"
						format={timeRangeInfo.format}
						ticks={timeRangeInfo.ticks}
						grid={{ style: `stroke: ${gridColor}; stroke-dasharray: 4` }}
						rule={{ style: `font-size: 11px; fill: ${textColor}; font-weight: bold;` }}
					/>
					{#each chartData as series}
						<Spline 
							data={series.data} 
							class="stroke-[3px]" 
							style={`stroke: ${series.color}; filter: drop-shadow(0 0 6px ${series.color}40);`} 
							curve={curveLinear} 
						/>
					{/each}
					<Highlight lines points={{ r: 3 }} />
				</Svg>
				<Tooltip.Root 
					class="bg-card/95 backdrop-blur-sm text-card-foreground shadow-xl border border-muted/30 rounded-lg p-3 min-w-[150px] z-50 text-sm"
				>
					{#snippet children({ data })}
						{@const active = Array.isArray(data) ? data[0] : data}
						{#if active && active.name}
							<div class="font-bold border-b border-muted/50 pb-1.5 mb-2 text-xs text-muted-foreground uppercase tracking-widest">
								{active.date ? new Date(active.date).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : ''}
							</div>
							<div class="flex items-center gap-3">
								<div class="w-3 h-3 rounded-full shadow-sm" style="background-color: {active.color}"></div>
								<div class="font-bold text-sm tracking-tight overflow-hidden text-ellipsis whitespace-nowrap max-w-[120px]">{active.name}</div>
								<div class="ml-auto font-black font-mono text-base">{active.value}</div>
							</div>
						{/if}
					{/snippet}
				</Tooltip.Root>
			</Chart>
		{:else}
			<div class="text-muted-foreground font-mono text-sm uppercase tracking-widest flex h-full items-center justify-center">
				No graph data visible
			</div>
		{/if}
	</div>

	<!-- Legend at bottom -->
	{#if chartData.length > 0}
		<div class="mt-6 flex flex-wrap justify-center gap-4 px-6 text-sm">
			{#each chartData as series}
				<div class="flex items-center gap-2 group cursor-pointer transition-opacity hover:opacity-100 opacity-80">
					<div
						class="h-4 w-4 rounded-full shadow-md"
						style="background-color: {series.color}; box-shadow: 0 0 10px {series.color}60;"
					></div>
					<span class="text-foreground tracking-tight font-bold">{series.name}</span>
				</div>
			{/each}
		</div>
	{/if}
</div>
