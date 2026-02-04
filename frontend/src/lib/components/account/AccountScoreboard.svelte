<script lang="ts">
	import * as Table from '@/components/ui/table';
	import EmptyState from '$lib/components/ui/empty-state.svelte';
	import StatusBadge from '$lib/components/ui/status-badge.svelte';
	import { ChartLine, Trophy } from '@lucide/svelte';
	import { formatDate, formatTimeSince, formatNumber } from '$lib/utils/formatting';

	let { solves } = $props<{ solves: any[] | undefined }>();

	type SortKey = 'name' | 'category' | 'points' | 'timestamp';

	let sortKey = $state<SortKey>('timestamp');
	let sortDir = $state<'asc' | 'desc'>('desc');

	const getPoints = (s: any) => formatNumber(s?.points ?? s?.score);

	function toggleSort(key: SortKey) {
		if (sortKey === key) {
			sortDir = sortDir === 'asc' ? 'desc' : 'asc';
		} else {
			sortKey = key;
			sortDir = key === 'timestamp' ? 'desc' : 'asc';
		}
	}

	const getSortArrow = (key: SortKey) => (sortKey === key ? (sortDir === 'asc' ? ' ▲' : ' ▼') : '');

	const sortedSolves = $derived.by(() => {
		const source = Array.isArray(solves) ? [...solves] : [];

		return source.sort((a: any, b: any) => {
			let valueA: any, valueB: any;

			switch (sortKey) {
				case 'name':
					valueA = a?.name ?? '';
					valueB = b?.name ?? '';
					break;
				case 'category':
					valueA = a?.category ?? '';
					valueB = b?.category ?? '';
					break;
				case 'points':
					valueA = getPoints(a);
					valueB = getPoints(b);
					break;
				case 'timestamp':
					valueA = new Date(a?.timestamp ?? 0).getTime();
					valueB = new Date(b?.timestamp ?? 0).getTime();
					break;
			}

			if (valueA < valueB) return sortDir === 'asc' ? -1 : 1;
			if (valueA > valueB) return sortDir === 'asc' ? 1 : -1;
			return 0;
		});
	});

	const totalPoints = $derived(sortedSolves.reduce((acc, s) => acc + getPoints(s), 0));
</script>

<div class="flex w-full flex-col">
	<div class="flex items-center gap-2 px-6 py-4">
		<ChartLine class="h-5 w-5 opacity-70" />
		<h3 class="text-xl font-semibold">Solves</h3>
	</div>

	<Table.Root class="w-full">
		<Table.Header class="bg-transparent [&_tr]:border-b-0">
			<Table.Row class="hover:bg-transparent">
				<Table.Head
					class="text-muted-foreground/70 w-[40%] cursor-pointer text-[10px] font-bold uppercase tracking-wider"
					onclick={() => toggleSort('name')}
				>
					Challenge{getSortArrow('name')}
				</Table.Head>
				<Table.Head
					class="text-muted-foreground/70 w-[20%] cursor-pointer text-[10px] font-bold uppercase tracking-wider"
					onclick={() => toggleSort('category')}
				>
					Category{getSortArrow('category')}
				</Table.Head>
				<Table.Head
					class="text-muted-foreground/70 w-[15%] cursor-pointer text-right text-[10px] font-bold uppercase tracking-wider"
					onclick={() => toggleSort('points')}
				>
					Points{getSortArrow('points')}
				</Table.Head>
				<Table.Head
					class="text-muted-foreground/70 w-[25%] cursor-pointer text-right text-[10px] font-bold uppercase tracking-wider sm:text-left"
					onclick={() => toggleSort('timestamp')}
				>
					Solved at{getSortArrow('timestamp')}
				</Table.Head>
			</Table.Row>
		</Table.Header>

		<Table.Body>
			{#if sortedSolves.length === 0}
				<Table.Row class="border-b-0">
					<Table.Cell colspan={4} class="p-0">
						<EmptyState
							icon={Trophy}
							title="No solves yet"
							description="Solve challenges to see them here"
						/>
					</Table.Cell>
				</Table.Row>
			{:else}
				{#each sortedSolves as solve (solve.id ?? solve.timestamp ?? solve.name)}
					<Table.Row class="border-b-0 transition-colors">
						<Table.Cell class="font-medium">{solve.name ?? '-'}</Table.Cell>
						<Table.Cell>
							<StatusBadge variant="category">{solve.category ?? '-'}</StatusBadge>
						</Table.Cell>
						<Table.Cell class="text-right">
							<div class="font-mono text-sm font-medium tabular-nums leading-none tracking-tight">
								{getPoints(solve)}
							</div>
						</Table.Cell>
						<Table.Cell class="text-right sm:text-left">
							<div class="flex items-center justify-end gap-2 sm:justify-start">
								<span>{formatDate(solve.timestamp)}</span>
								<span class="text-muted-foreground text-xs"
									>({formatTimeSince(solve.timestamp)} ago)</span
								>
							</div>
						</Table.Cell>
					</Table.Row>
				{/each}
			{/if}
		</Table.Body>
	</Table.Root>
</div>
