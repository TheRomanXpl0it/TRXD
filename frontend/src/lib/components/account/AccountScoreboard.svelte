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

	const getSortArrow = (key: SortKey) => sortKey === key ? (sortDir === 'asc' ? ' ▲' : ' ▼') : '';

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

<div class="flex flex-col w-full">
	<div class="flex items-center gap-2">
		<ChartLine class="h-5 w-5 opacity-70" />
		<h3 class="text-xl font-semibold">Solves</h3>
	</div>

	<Table.Root class="w-full">
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-[40%] cursor-pointer" onclick={() => toggleSort('name')}>
					Challenge{getSortArrow('name')}
				</Table.Head>
				<Table.Head class="w-[20%] cursor-pointer" onclick={() => toggleSort('category')}>
					Category{getSortArrow('category')}
				</Table.Head>
				<Table.Head class="w-[15%] text-right cursor-pointer" onclick={() => toggleSort('points')}>
					Points{getSortArrow('points')}
				</Table.Head>
				<Table.Head class="w-[25%] cursor-pointer text-right sm:text-left" onclick={() => toggleSort('timestamp')}>
					Solved at{getSortArrow('timestamp')}
				</Table.Head>
			</Table.Row>
		</Table.Header>

		<Table.Body>
			{#if sortedSolves.length === 0}
				<Table.Row>
					<Table.Cell colspan={4} class="p-0">
						<EmptyState icon={Trophy} title="No solves yet" description="Solve challenges to see them here" />
					</Table.Cell>
				</Table.Row>
			{:else}
				{#each sortedSolves as solve (solve.id ?? solve.timestamp ?? solve.name)}
					<Table.Row>
						<Table.Cell class="font-medium">{solve.name ?? '-'}</Table.Cell>
						<Table.Cell>
							<StatusBadge variant="category">{solve.category ?? '-'}</StatusBadge>
						</Table.Cell>
						<Table.Cell class="text-right">{getPoints(solve)}</Table.Cell>
						<Table.Cell class="text-right sm:text-left">
							<div class="flex items-center justify-end sm:justify-start gap-2">
								<span>{formatDate(solve.timestamp)}</span>
								<span class="text-xs text-muted-foreground">({formatTimeSince(solve.timestamp)} ago)</span>
							</div>
						</Table.Cell>
					</Table.Row>
				{/each}
			{/if}
		</Table.Body>
	</Table.Root>
</div>
