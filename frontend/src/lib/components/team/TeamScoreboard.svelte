<script lang="ts">
	import {
		Table,
		TableBody,
		TableCaption,
		TableCell,
		TableHead,
		TableHeader,
		TableRow
	} from '@/components/ui/table';
	import { ChartLine } from '@lucide/svelte';

	// Prop (team contains solves[] and members[])
	let { team } = $props<{ team: any }>();
	let solves = $state(team?.solves);

	// Map of memberId -> memberName for quick lookup
	let memberNameById: Record<string, string> = $state({});

	// State
	type SortKey = 'name' | 'category' | 'points' | 'timestamp' | 'solver';
	let sortKey = $state<SortKey>('timestamp');
	let sortDir = $state<'asc' | 'desc'>('desc');
	let rows: any[] = $state([]); // what we actually render
	let totalPoints = $state(0);

	// Helpers
	const getPoints = (s: any) => Number(s?.points ?? s?.value ?? s?.score ?? 0);

	const fmtDate = (iso?: string) => {
		if (!iso) return '-';
		const d = new Date(iso);
		return Number.isNaN(+d) ? '-' : d.toLocaleString();
	};

	const timeSince = (iso?: string) => {
		if (!iso) return '-';
		const sec = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 1000));
		const h = Math.floor(sec / 3600);
		const m = Math.floor((sec % 3600) / 60);
		const s = sec % 60;
		if (h > 0) return `${h}h ${m}m`;
		if (m > 0) return `${m}m ${s}s`;
		return `${s}s`;
	};

	const truncateName = (name: string, maxLength = 32): string => {
		if (!name || name.length <= maxLength) return name;
		return name.slice(0, maxLength) + '...';
	};

	function toggleSort(key: SortKey) {
		if (sortKey === key) sortDir = sortDir === 'asc' ? 'desc' : 'asc';
		else {
			sortKey = key;
			sortDir = key === 'timestamp' ? 'desc' : 'asc';
		}
	}
	const arrow = (key: SortKey | string) =>
		sortKey === key ? (sortDir === 'asc' ? ' ▲' : ' ▼') : '';

	// Build member lookup map whenever team.members changes
	$effect(() => {
		const src = Array.isArray(team?.members) ? team.members : [];
		const map: Record<string, string> = {};
		for (const m of src) {
			const id = String(m?.id ?? m?.user_id ?? '');
			if (!id) continue;
			map[id] = String(m?.name ?? m?.username ?? m?.displayName ?? 'Unknown');
		}
		memberNameById = map;
	});

	const solverName = (uid: any) => {
		const key = uid == null ? '' : String(uid);
		return memberNameById[key] ?? '-';
	};

	// Compute rows reactively (robust against undefined / late props)
	$effect(() => {
		const src = Array.isArray(solves) ? solves : [];
		const arr = [...src];

		arr.sort((a: any, b: any) => {
			let av: any, bv: any;
			switch (sortKey) {
				case 'name':
					av = a?.name ?? '';
					bv = b?.name ?? '';
					break;
				case 'category':
					av = a?.category ?? '';
					bv = b?.category ?? '';
					break;
				case 'points':
					av = getPoints(a);
					bv = getPoints(b);
					break;
				case 'timestamp':
					av = new Date(a?.timestamp ?? 0).getTime();
					bv = new Date(b?.timestamp ?? 0).getTime();
					break;
				case 'solver':
					av = solverName(a?.user_id) ?? '';
					bv = solverName(b?.user_id) ?? '';
					break;
			}
			if (av < bv) return sortDir === 'asc' ? -1 : 1;
			if (av > bv) return sortDir === 'asc' ? 1 : -1;
			return 0;
		});

		rows = arr;
		totalPoints = arr.reduce((acc: number, s: any) => acc + getPoints(s), 0);
	});
</script>

<div class="flex w-full flex-col">
	<div class="mb-4 flex items-center gap-2 pl-1">
		<ChartLine class="h-5 w-5 opacity-70" />
		<h3 class="text-xl font-semibold">Solves</h3>
	</div>
	<Table class="w-full">
		<TableCaption class="text-sm">
			{#if totalPoints > 0}{totalPoints} pts total{/if}
		</TableCaption>

		<TableHeader class="bg-transparent [&_tr]:border-b-0">
			<TableRow class="hover:bg-transparent">
				<TableHead
					class="text-muted-foreground/70 w-[28%] cursor-pointer text-[10px] font-bold uppercase tracking-wider"
					onclick={() => toggleSort('name')}
				>
					Challenge {arrow('name')}
				</TableHead>
				<TableHead
					class="text-muted-foreground/70 w-[16%] cursor-pointer text-[10px] font-bold uppercase tracking-wider"
					onclick={() => toggleSort('category')}
				>
					Category {arrow('category')}
				</TableHead>
				<TableHead
					class="text-muted-foreground/70 w-[12%] cursor-pointer text-right text-[10px] font-bold uppercase tracking-wider"
					onclick={() => toggleSort('points')}
				>
					Points {arrow('points')}
				</TableHead>
				<TableHead
					class="text-muted-foreground/70 w-[18%] cursor-pointer text-[10px] font-bold uppercase tracking-wider"
					onclick={() => toggleSort('solver')}
				>
					Solved by {arrow('solver')}
				</TableHead>
				<TableHead
					class="text-muted-foreground/70 w-[16%] cursor-pointer text-[10px] font-bold uppercase tracking-wider"
					onclick={() => toggleSort('timestamp')}
				>
					Solved at {arrow('timestamp')}
				</TableHead>
				<TableHead
					class="text-muted-foreground/70 w-[10%] text-right text-[10px] font-bold uppercase tracking-wider"
				>
					Ago
				</TableHead>
			</TableRow>
		</TableHeader>

		<TableBody>
			{#if rows.length === 0}
				<TableRow class="border-b-0">
					<TableCell colspan={6} class="text-muted-foreground py-10 text-center">
						No solves yet.
					</TableCell>
				</TableRow>
			{:else}
				{#each rows as s (s.id ?? s.timestamp ?? s.name ?? Math.random())}
					<TableRow class="border-b-0 transition-colors">
						<TableCell class="font-medium">{truncateName(s.name ?? '-')}</TableCell>
						<TableCell>
							<span class="text-muted-foreground text-xs font-medium">
								{s.category ?? '-'}
							</span>
						</TableCell>
						<TableCell class="text-right">
							<div class="font-mono text-sm font-medium tabular-nums leading-none tracking-tight">
								{getPoints(s)}
							</div>
						</TableCell>
						<TableCell>{truncateName(solverName(s.user_id))}</TableCell>
						<TableCell>{fmtDate(s.timestamp)}</TableCell>
						<TableCell class="text-right">{timeSince(s.timestamp)}</TableCell>
					</TableRow>
				{/each}
			{/if}
		</TableBody>
	</Table>
</div>
