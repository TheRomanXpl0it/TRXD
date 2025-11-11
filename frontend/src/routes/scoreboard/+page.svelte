<script lang="ts">
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import ScoreHistory from '$lib/components/scoreboard/Graph.svelte';
	import { Medal, Rows3, Grid2x2 } from '@lucide/svelte';
	import { getGraphData, getScoreboard } from '@/scoreboard';
	import { onMount } from 'svelte';
	import { link, push } from 'svelte-spa-router';
	import { userMode } from '$lib/stores/auth';
	import ErrorMessage from '$lib/components/ui/error-message.svelte';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Button } from '@/components/ui/button';
	import { Avatar } from 'flowbite-svelte';
	import { BugOutline } from 'flowbite-svelte-icons';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { createQuery } from '@tanstack/svelte-query';

	let perPage = $state(10);
	let compactMode = $state(false);

	const scoreboardQuery = createQuery(() => ({
		queryKey: ['scoreboard'],
		queryFn: getScoreboard,
		staleTime: 30_000
	}));

	const graphQuery = createQuery(() => ({
		queryKey: ['scoreboard-graph'],
		queryFn: getGraphData,
		staleTime: 30_000
	}));

	// Derived values from queries
	const scoreboardData = $derived(scoreboardQuery.data ?? []);
	const graphData = $derived(graphQuery.data ?? []);
	const loading = $derived(scoreboardQuery.isLoading || graphQuery.isLoading);
	const error = $derived(scoreboardQuery.error?.message ?? graphQuery.error?.message ?? null);

	const teamNames = $derived(
		scoreboardData.reduce(
			(acc, team) => {
				acc[team.id] = team.name;
				return acc;
			},
			{} as Record<string, string>
		)
	);

	// Load compact mode from localStorage
	onMount(() => {
		const saved = localStorage.getItem('scoreboard-compact-mode');
		if (saved !== null) {
			compactMode = saved === 'true';
		}
	});

	$effect(() => {
		localStorage.setItem('scoreboard-compact-mode', String(compactMode));
	});

	// sort by score desc
	const sorted = $derived(
		(Array.isArray(scoreboardData) ? [...scoreboardData] : []).sort(
			(a, b) => (Number(b?.score) || 0) - (Number(a?.score) || 0)
		)
	);
	const count = $derived(sorted.length);

	function medalClass(rank: number) {
		if (rank === 1) return 'text-yellow-500';
		if (rank === 2) return 'text-gray-300';
		if (rank === 3) return 'text-amber-700';
		return '';
	}

	// Conditional navigation based on userMode
	const hrefForItem = (id: number | string) => ($userMode ? `#/account/${id}` : `#/team/${id}`);
	const pushItem = (id: number | string) =>
		$userMode ? push(`/account/${id}`) : push(`/team/${id}`);

	// TODO: Fix this, this does not work
	let currentPage = $state(1);
	$effect(() => {
		if (currentPage > 1) {
			setTimeout(() => {
				const paginationEl = document.getElementById('pagination-controls');
				if (paginationEl) {
					paginationEl.scrollIntoView({ behavior: 'instant', block: 'nearest' });
				}
			}, 0);
		}
	});
</script>

<div class="w-full overflow-hidden">
	<div class="mb-6">
		<div class="mt-5 flex flex-wrap items-center justify-between gap-3">
			<p class="text-2xl font-bold text-gray-800 dark:text-gray-100 sm:text-3xl">Scoreboard</p>
			<div class="flex gap-2">
				<Button
					variant={!compactMode ? 'default' : 'outline'}
					size="sm"
					onclick={() => compactMode = false}
					class="cursor-pointer gap-1.5"
				>
					<Grid2x2 class="h-4 w-4" />
					<span class="hidden sm:inline">Full</span>
				</Button>
				<Button
					variant={compactMode ? 'default' : 'outline'}
					size="sm"
					onclick={() => compactMode = true}
					class="cursor-pointer gap-1.5"
				>
					<Rows3 class="h-4 w-4" />
					<span class="hidden sm:inline">Compact</span>
				</Button>
			</div>
		</div>
		<hr class="my-2 mb-10 h-px border-0 bg-gray-200 dark:bg-gray-700" />
	</div>

	{#if loading}
		<div class="flex flex-col items-center justify-center py-12">
			<Spinner class="mb-4 h-8 w-8" />
			<p class="text-gray-600 dark:text-gray-400">Loading scoreboard...</p>
		</div>
	{:else if error}
		<ErrorMessage title="Error loading scoreboard" message={error} />
	{:else}
		<div class="mb-6 max-w-full">
			<ScoreHistory data={graphData} topN={10} {teamNames} userMode={$userMode} />
		</div>

	{@const startIndex = (currentPage - 1) * perPage}
	{@const pageRows = sorted.slice(startIndex, startIndex + perPage)}
	{@const totalPages = Math.max(1, Math.ceil(count / perPage))}
	{@const singlePage = totalPages <= 1}

	<!-- Table -->
	<div class="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head class="w-[4rem] sm:w-[5rem]">Rank</Table.Head>
					<Table.Head class="min-w-[7rem]">{$userMode ? 'Player' : 'Team'}</Table.Head>
					<Table.Head class="w-[5rem] text-right {compactMode ? 'pr-6' : ''}">Score</Table.Head>
					{#if !compactMode}
						<Table.Head class="w-[10rem] min-w-[10rem]">Badges</Table.Head>
					{/if}
				</Table.Row>
			</Table.Header>

			<Table.Body>
				{#if pageRows.length === 0}
					<Table.Row>
						<Table.Cell colspan={compactMode ? 3 : 4} class="py-8 text-center text-gray-500">
							{$userMode ? 'No players yet.' : 'No teams yet.'}
						</Table.Cell>
					</Table.Row>
				{:else}
					{#each pageRows as row, i (row.id)}
						{@const rank = startIndex + i + 1}
						<Table.Row>
							<Table.Cell class="font-medium">
								<div class="flex items-center gap-1.5">
									<span class="text-sm">#{rank}</span>
									{#if rank <= 3}
										<Medal class={`h-4 w-4 ${medalClass(rank)}`} aria-label="Medal" />
									{/if}
								</div>
							</Table.Cell>

							<Table.Cell class="max-w-32 sm:max-w-xs">
								<div class="flex items-center gap-2 min-w-0">
									{#if !compactMode}
										{#if row.image}
											<div class="h-8 w-8 shrink-0 overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700">
												<img src={row.image} alt={row.name} class="h-full w-full object-cover" />
											</div>
										{:else}
											<Avatar class="h-8 w-8 shrink-0">
												<BugOutline class="h-4 w-4" />
											</Avatar>
										{/if}
									{/if}
									<!-- svelte-spa-router in-page navigation -->
									<a
										href={hrefForItem(row.id)}
										use:link
										onclick={(e) => {
											e.preventDefault();
											pushItem(row.id);
										}}
										class="text-primary cursor-pointer text-sm underline-offset-2 hover:underline sm:text-base truncate"
										title={`View ${$userMode ? 'player' : 'team'} ${row.name}`}
									>
										{row.name}
									</a>
								</div>
							</Table.Cell>

							<Table.Cell class="text-right tabular-nums text-sm text-gray-600 dark:text-gray-400 {compactMode ? 'pr-6' : ''}">{row.score}</Table.Cell>

							{#if !compactMode}
								<Table.Cell>
									{#if Array.isArray(row.badges) && row.badges.length}
										<div class="flex flex-wrap gap-1">
											{#each row.badges as b, bi (bi)}
												<Tooltip.Root>
													<Tooltip.Trigger>
														<span class="rounded-full border px-2 py-0.5 text-xs">{b.name}</span>
													</Tooltip.Trigger>
													{#if b.description}
														<Tooltip.Content>
															<p>{b.description}</p>
														</Tooltip.Content>
													{/if}
												</Tooltip.Root>
											{/each}
										</div>
									{:else}
										<span class="text-xs text-gray-500">-</span>
									{/if}
								</Table.Cell>
							{/if}
						</Table.Row>
					{/each}
				{/if}
			</Table.Body>
		</Table.Root>
	</div>

	<!-- Pagination directly under the table -->
	<Pagination.Root {count} {perPage} bind:page={currentPage} siblingCount={0} class="mt-4">
		{#snippet children({ pages, currentPage: pageNum })}
			<div class="flex w-full justify-center overflow-x-auto" id="pagination-controls">
				<Pagination.Content
					class={`flex items-center justify-center gap-1 ${singlePage ? 'pointer-events-none opacity-50' : ''}`}
					aria-disabled={singlePage}
				>
					<Pagination.Item>
						<Pagination.PrevButton />
					</Pagination.Item>

					{#each pages as page (page.key)}
						{#if page.type === 'ellipsis'}
							<Pagination.Item>
								<Pagination.Ellipsis />
							</Pagination.Item>
						{:else}
							<Pagination.Item>
								<Pagination.Link {page} isActive={pageNum === page.value}>
									{page.value}
								</Pagination.Link>
							</Pagination.Item>
						{/if}
					{/each}

					<Pagination.Item>
						<Pagination.NextButton />
					</Pagination.Item>
				</Pagination.Content>
			</div>
		{/snippet}
	</Pagination.Root>
{/if}
</div>
