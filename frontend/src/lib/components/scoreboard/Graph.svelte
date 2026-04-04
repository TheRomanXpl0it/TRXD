<script lang="ts">
	import { Chart, Svg, Axis, Spline, Highlight, Tooltip as ChartTooltip } from 'layerchart';
	import { scaleTime, scaleLinear } from 'd3-scale';
	// @ts-ignore - Ignore missing type declarations for d3-shape
	import { curveStepAfter } from 'd3-shape';

	// Props interface
	interface Props {
		data?: any[];
		timeMin?: number;
		timeMax?: number;
		teamNames?: Record<string, string>;
		userMode?: boolean;
		compact?: boolean;
		height?: string;
	}

	let { data = [], timeMin, timeMax, teamNames, userMode = false, compact = false, height = '' }: Props = $props();

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

	// Competitive colors for top 3
	const top3Colors = ['#fbbf24', '#94a3b8', '#cd7f32'];
	// Dynamic vibrant colors for others
	const colors = [
		'#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6', '#06b6d4', '#f97316', '#84cc16', '#ec4899', '#6366f1'
	];

	let highlightedTeam = $state<string | number | null>(null);

	const chartData = $derived.by(() => {
		const arr = Array.isArray(data) ? data : [];

		const ranked = [...arr]
			.map((e: any) => {
				const solves = normalizeTeam(e);
				const total = solves.length > 0 ? Math.max(...solves.map((s) => Number(s?.points ?? 0))) : 0;
				const lastSolve = solves.length > 0 ? Math.max(...solves.map((s) => s.date?.getTime() || 0)) : 0;
				return { ...e, total, lastSolve };
			})
			.sort((a: any, b: any) => 
				(b.total - a.total) || 
				(a.lastSolve - b.lastSolve) || 
				((a.id || a.team_id || 0) - (b.id || b.team_id || 0))
			);

		let minTime = Date.now();
		let hasData = false;
		ranked.forEach((team: any) => {
			const solves = normalizeTeam(team).filter((s: any) => s.date !== null);
			if (solves.length > 0 && solves[0].date) {
				hasData = true;
				minTime = Math.min(minTime, solves[0].date.getTime() - 60000);
			}
		});

		const nowMs = Date.now();

		const series: Array<{
			id: string | number;
			name: string;
			data: Array<{ date: Date; value: number; name?: string; color?: string; }>;
			color: string;
		}> = [];

		ranked.forEach((team: any, index: number) => {
			const name = nameForTeam(team);
			const color = index < 3 ? top3Colors[index] : colors[(index - 3) % colors.length];
			const solves = normalizeTeam(team)
				.filter((s: any) => s.date !== null)
				.sort((a: any, b: any) => a.date!.getTime() - b.date!.getTime());

			if (solves.length === 0) return;

			// Use only actual solve data points — step function via curveStepAfter
			const points: Array<{ date: Date; value: number; name: string; color: string; }> = [];

			// Start at 0
			points.push({ date: new Date(minTime), value: 0, name, color });

			for (const s of solves) {
				points.push({
					date: s.date,
					value: Number(s.points ?? 0),
					name,
					color
				});
			}

			// Extend to now
			points.push({
				date: new Date(nowMs),
				value: points[points.length - 1].value,
				name,
				color
			});

			series.push({ id: team.team_id || team.id, name, data: points, color });
		});

		return series;
	});

	const textColor = $derived('hsl(var(--muted-foreground))');
	const gridColor = $derived('hsl(var(--border))');
</script>

<div class="w-full flex-col flex" style={height ? `height: ${height}` : `height: ${compact ? '280px' : '380px'}`}>
	<div class="flex-grow w-full">
		{#if chartData.length > 0}
			<Chart
				data={chartData.flatMap((s) => s.data)}
				x="date"
				xScale={scaleTime()}
				y="value"
				yScale={scaleLinear()}
				yDomain={[0, null]}
				padding={{ left: compact ? 16 : 24, bottom: 8, top: compact ? 8 : 12, right: compact ? 8 : 12 }}
				tooltip={{ mode: 'voronoi' }}
			>
				<Svg>
					<Axis
						placement="left"
						grid={{ style: `stroke: ${gridColor}; stroke-dasharray: 4` }}
						rule={{ style: `font-size: ${compact ? '9px' : '11px'}; fill: ${textColor}; font-weight: bold;` }}
					/>
					{#each chartData as series}
						<Spline 
							data={series.data} 
							class="stroke-[3px] transition-all duration-300" 
							style={`stroke: ${series.color}; opacity: ${highlightedTeam === null || highlightedTeam === series.id ? 1 : 0.1};`} 
							curve={curveStepAfter} 
						/>
					{/each}
					<Highlight lines points={{ r: 6 }} />
				</Svg>
				<ChartTooltip.Root 
					class="bg-card/95 backdrop-blur-sm text-card-foreground shadow-xl border border-muted/30 rounded-lg p-3 min-w-[150px] z-50 text-sm"
				>
					{#snippet children({ data })}
						{@const active = Array.isArray(data) ? data[0] : data}
						{#if active && active.name}
							<div class="font-bold border-b border-muted/50 pb-1.5 mb-2 text-xs text-muted-foreground uppercase tracking-widest">
								{active.date ? new Date(active.date).toLocaleString([], { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit', hour12: false }) : ''}
							</div>
							<div class="flex items-center gap-3">
								<div class="w-3 h-3 rounded-full shadow-sm" style="background-color: {active.color}"></div>
								<div class="font-bold text-sm tracking-tight overflow-hidden text-ellipsis whitespace-nowrap max-w-[120px]">{active.name}</div>
								<div class="ml-auto font-black font-mono text-base">{active.value}</div>
							</div>
						{/if}
					{/snippet}
				</ChartTooltip.Root>
			</Chart>
		{:else}
			<div class="text-muted-foreground font-mono text-sm uppercase tracking-widest flex h-full items-center justify-center">
				No graph data visible
			</div>
		{/if}
	</div>

	<!-- Legend at bottom -->
	{#if chartData.length > 0 && !compact}
		<div class="mt-6 flex flex-wrap justify-center gap-4 px-6 text-sm">
			{#each chartData as series}
				{@const lastScore = series.data.length > 0 ? series.data[series.data.length - 1].value : 0}
				<button
					type="button"
					class="flex items-center gap-2 group cursor-pointer transition-all duration-200 rounded-md px-2 py-1 hover:bg-muted/50 focus:outline-none focus:ring-2 focus:ring-primary/30 {highlightedTeam !== null && highlightedTeam !== series.id ? 'opacity-30 grayscale-[50%]' : 'opacity-100'}"
					title="{series.name}: {lastScore} pts"
					onclick={() => {
						if (highlightedTeam === series.id) {
							highlightedTeam = null;
						} else {
							highlightedTeam = series.id;
						}
					}}
				>
					<div
						class="h-4 w-4 rounded-full shadow-md transition-transform duration-200 {highlightedTeam === series.id ? 'scale-125 border-2 border-white' : ''}"
						style="background-color: {series.color};"
					></div>
					<span class="text-foreground tracking-tight font-bold">{series.name}</span>
				</button>
			{/each}
		</div>
	{/if}
</div>
