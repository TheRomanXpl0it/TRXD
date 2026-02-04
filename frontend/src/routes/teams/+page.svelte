<script lang="ts">
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { getTeams } from '@/team';
	import { link, push } from 'svelte-spa-router';
	import ErrorMessage from '$lib/components/ui/error-message.svelte';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import countries from '$lib/data/countries.json';
	import { userMode } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import EmptyState from '$lib/components/ui/empty-state.svelte';
	import { Users, Globe } from '@lucide/svelte';
	import GeneratedAvatar from '$lib/components/ui/avatar/generated-avatar.svelte';

	// Redirect to accounts page if in user mode
	onMount(() => {
		if ($userMode) {
			push('/accounts');
		}
	});

	let perPage = $state(20);
	let currentPage = $state(1);

	const teamsQuery = createQuery(() => ({
		queryKey: ['teams', currentPage, perPage],
		queryFn: () => getTeams(currentPage, perPage),
		staleTime: 5 * 60 * 1000,
		placeholderData: (previousData) => previousData
	}));

	// Handle both PaginatedResponse and flat array (fallback)
	const rawData = $derived(teamsQuery.data);
	const isPaginated = $derived(!!rawData?.pagination);
	const teamsData = $derived(Array.isArray(rawData) ? rawData : (rawData?.data ?? []));
	const totalCount = $derived(isPaginated ? (rawData?.pagination?.total ?? 0) : teamsData.length);

	const loading = $derived(teamsQuery.isLoading);
	const error = $derived(teamsQuery.error?.message ?? null);

	// Helper to convert ISO3 to ISO2 for flag icons
	function getCountryIso2(iso3: string): string | null {
		const country = (countries as any[]).find((c) => c.iso3?.toUpperCase() === iso3?.toUpperCase());
		return country?.iso2?.toUpperCase() ?? null;
	}

	function truncateName(name: string, maxLength = 32): string {
		if (!name || name.length <= maxLength) return name;
		return name.slice(0, maxLength) + '...';
	}

	function getInitials(name: string): string {
		return name
			.split(' ')
			.map((n) => n[0])
			.slice(0, 2)
			.join('')
			.toUpperCase();
	}

	function getAvatarColor(name: string): string {
		const colors = [
			'bg-red-500',
			'bg-orange-500',
			'bg-amber-500',
			'bg-yellow-500',
			'bg-lime-500',
			'bg-green-500',
			'bg-emerald-500',
			'bg-teal-500',
			'bg-cyan-500',
			'bg-sky-500',
			'bg-blue-500',
			'bg-indigo-500',
			'bg-violet-500',
			'bg-purple-500',
			'bg-fuchsia-500',
			'bg-pink-500',
			'bg-rose-500'
		];
		let hash = 0;
		for (let i = 0; i < name.length; i++) {
			hash = name.charCodeAt(i) + ((hash << 5) - hash);
		}
		return colors[Math.abs(hash) % colors.length];
	}

	// Sort alphabetically by name
	const sorted = $derived(
		(Array.isArray(teamsData) ? [...teamsData] : []).sort((a, b) =>
			(a?.name || '').localeCompare(b?.name || '')
		)
	);
	const count = $derived(totalCount);

	// Track page changes
	$effect(() => {
		if (currentPage > 1) {
			setTimeout(() => {
				const paginationEl = document.getElementById('pagination-controls');
				if (paginationEl) {
					paginationEl.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
				}
			}, 0);
		}
	});
</script>

<div class="mx-auto max-w-5xl space-y-8 py-10">
	<div
		class="from-muted/20 to-background mb-6 mt-6 rounded-xl border-0 bg-gradient-to-br p-6 shadow-sm"
	>
		<div class="flex items-center gap-4">
			<div
				class="bg-background flex h-16 w-16 shrink-0 items-center justify-center rounded-full shadow-sm"
			>
				<Users class="text-muted-foreground h-8 w-8" />
			</div>
			<div>
				<h1 class="text-3xl font-bold tracking-tight">Teams</h1>
				<p class="text-muted-foreground mt-2 text-sm">View and manage competing teams.</p>
			</div>
		</div>
	</div>

	{#if error}
		<ErrorMessage title="Error loading teams" message={error} />
	{:else}
		<Card.Root class="overflow-hidden border-0 shadow-sm">
			<Card.Content class="p-0">
				<div class="relative w-full overflow-auto">
					<Table.Root>
						<Table.Header class="bg-transparent [&_tr]:border-b-0">
							<Table.Row class="hover:bg-transparent">
								<Table.Head
									class="text-muted-foreground/70 w-[400px] bg-transparent text-[10px] font-bold uppercase tracking-wider"
									>Team</Table.Head
								>
								<Table.Head
									class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider"
									>Country</Table.Head
								>
								<Table.Head
									class="text-muted-foreground/70 bg-transparent text-right text-[10px] font-bold uppercase tracking-wider"
									>Score</Table.Head
								>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#if loading && teamsData.length === 0}
								{#each Array(5) as _}
									<Table.Row>
										<Table.Cell>
											<Skeleton class="h-4 w-[150px]" />
										</Table.Cell>
										<Table.Cell><Skeleton class="h-5 w-[100px]" /></Table.Cell>
										<Table.Cell class="text-right"
											><Skeleton class="ml-auto h-4 w-[50px]" /></Table.Cell
										>
									</Table.Row>
								{/each}
							{:else}
								{@const pageRows = isPaginated
									? sorted
									: sorted.slice((currentPage - 1) * perPage, currentPage * perPage)}

								{#if pageRows.length === 0}
									<Table.Row>
										<Table.Cell colspan={3} class="h-[300px] text-center">
											<EmptyState
												icon={Users}
												title="No teams found"
												description="There are no teams to display at the moment."
											/>
										</Table.Cell>
									</Table.Row>
								{:else}
									{#each pageRows as team, i (team.id)}
										<Table.Row
											class="hover:bg-muted/50 cursor-pointer border-b-0 transition-colors"
											onclick={() => push(`/team/${team.id}`)}
										>
											<Table.Cell class="py-3">
												<div class="flex items-center gap-3">
													<div
														class="border-border h-8 w-8 shrink-0 overflow-hidden rounded-full border"
													>
														<GeneratedAvatar seed={team.name} class="h-full w-full" />
													</div>
													<span
														class="text-foreground decoration-primary/50 font-medium underline-offset-4 hover:underline"
													>
														{truncateName(team.name)}
													</span>
												</div>
											</Table.Cell>

											<Table.Cell>
												{#if team.country}
													{@const iso2 = getCountryIso2(team.country)}
													<div class="flex items-center gap-2">
														<div
															class="bg-background relative flex h-6 w-8 items-center justify-center overflow-hidden rounded border shadow-sm"
														>
															{#if iso2}
																<Icon
																	icon={`circle-flags:${iso2.toLowerCase()}`}
																	width="32"
																	height="32"
																/>
															{:else}
																<Globe class="text-muted-foreground h-4 w-4" />
															{/if}
														</div>
														<span class="text-foreground/80 text-sm">{team.country}</span>
													</div>
												{:else}
													<span class="text-muted-foreground text-sm italic">Not specified</span>
												{/if}
											</Table.Cell>

											<Table.Cell class="text-right">
												<div
													class="font-mono text-sm font-medium tabular-nums leading-none tracking-tight"
												>
													{team.score?.toLocaleString() ?? 0} pts
												</div>
											</Table.Cell>
										</Table.Row>
									{/each}
								{/if}
							{/if}
						</Table.Body>
					</Table.Root>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Pagination -->
		{#if count > perPage}
			<Pagination.Root {count} {perPage} bind:page={currentPage} siblingCount={1} class="mt-4">
				{#snippet children({ pages, currentPage: pageNum })}
					<div class="flex w-full justify-center overflow-x-auto py-4" id="pagination-controls">
						<Pagination.Content class="gap-4">
							<Pagination.Item class="mx-2">
								<Pagination.PrevButton class="h-9 w-9" />
							</Pagination.Item>

							{#each pages as page (page.key)}
								{#if page.type === 'ellipsis'}
									<Pagination.Item>
										<Pagination.Ellipsis />
									</Pagination.Item>
								{:else}
									<Pagination.Item>
										<Pagination.Link
											{page}
											isActive={pageNum === page.value}
											class="data-[selected]:bg-foreground data-[selected]:text-background h-9 w-9 transition-all data-[selected]:shadow-md"
										>
											{page.value}
										</Pagination.Link>
									</Pagination.Item>
								{/if}
							{/each}

							<Pagination.Item class="mx-2">
								<Pagination.NextButton class="h-9 w-9" />
							</Pagination.Item>
						</Pagination.Content>
					</div>
				{/snippet}
			</Pagination.Root>
		{/if}
	{/if}
</div>
