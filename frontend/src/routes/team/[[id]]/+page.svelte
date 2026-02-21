<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Globe, Users, Edit, Award, LayoutGrid, List, Mail, Medal } from '@lucide/svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	import SolveListTable from '$lib/components/team/TeamScoreboard.svelte';
	import TeamMembers from '$lib/components/team/TeamMemberlist.svelte';
	import TeamJoinCreate from '$lib/components/team/TeamJoinCreate.svelte';
	import TeamEdit from '$lib/components/team/TeamEdit.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { getTeam } from '$lib/team';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import GeneratedAvatar from '$lib/components/ui/avatar/generated-avatar.svelte';
	import RadarChart from '$lib/components/RadarChart.svelte';
	import CountryFlag from '$lib/components/ui/country-flag.svelte';
	import countries from '$lib/data/countries.json';
	import { createQuery } from '@tanstack/svelte-query';
	import { authState, loadUser } from '$lib/stores/auth';

	let teamEditOpen = $state(false);
	let activeTab = $state<'overview' | 'solves'>('overview');

	const currentTeamId = $derived.by(() => {
		if (!authState.ready) return null;
		const routeKey = normalizeKey($page.params.id);
		const fallbackKey = normalizeKey(authState.user?.team_id);
		const effectiveKey = routeKey ?? fallbackKey;
		return effectiveKey ? validateId(effectiveKey) : null;
	});

	const teamQuery = createQuery(() => ({
		queryKey: ['team', currentTeamId],
		queryFn: () => getTeam(currentTeamId!),
		enabled: currentTeamId !== null,
		staleTime: 10_000
	}));

	const team = $derived(teamQuery.data ?? null);
	const loading = $derived(teamQuery.isLoading && currentTeamId !== null);
	const teamError = $derived(teamQuery.error?.message ?? null);

	const isOwnTeam = $derived(
		authState.user && team && String(authState.user.team_id) === String(team.id)
	);

	function normalizeKey(x: unknown): string | null {
		const s = String(x ?? '').trim();
		return s ? s : null;
	}

	function validateId(key: string): number | null {
		// Only allow positive integers
		if (/^\d+$/.test(key)) {
			const num = Number(key);
			return num > 0 ? num : null;
		}
		return null;
	}

	$effect(() => {
		if (authState.ready && authState.userMode) {
			goto('/accounts');
		}
	});
</script>

{#if !authState.ready || loading}
	<div class="mx-auto max-w-6xl space-y-8 py-10">
		<div class="space-y-4">
			<div class="bg-muted h-8 w-48 animate-pulse rounded"></div>
			<Separator />
			<div class="flex flex-col items-center justify-center py-12">
				<Spinner class="mb-4 h-8 w-8" />
				<p class="text-muted-foreground">Loading team...</p>
			</div>
		</div>
	</div>
{:else if teamError}
	<div class="mx-auto max-w-6xl py-10">
		<div
			class="border-destructive/50 bg-destructive/10 text-destructive dark:border-destructive dark:bg-destructive/20 rounded-lg border p-6"
		>
			<p class="text-lg font-semibold">Error loading team</p>
			<p>{teamError}</p>
		</div>
	</div>
{:else if !team}
	<div class="mx-auto max-w-4xl py-10">
		<TeamJoinCreate oncreated={() => loadUser()} onjoined={() => loadUser()} />
	</div>
{:else}
	<div class="mx-auto max-w-6xl space-y-8 py-10">
		<!-- Header -->
		<div class="flex items-center justify-between">
			<h2 class="text-3xl font-bold tracking-tight">Team Profile</h2>
			{#if isOwnTeam}
				<Button variant="outline" size="sm" onclick={() => (teamEditOpen = true)} class="gap-2">
					<Edit class="h-4 w-4" />
					Edit Team
				</Button>
			{/if}
		</div>

		<!-- Hero Card -->
		<Card.Root
			class="from-muted/20 to-background overflow-hidden border-0 bg-gradient-to-br shadow-sm"
		>
			<Card.Content class="p-6 sm:p-8">
				<div class="flex flex-col items-start gap-6 sm:flex-row sm:items-center">
					<!-- Avatar with Country Badge -->
					<div
						class="border-background bg-background relative h-20 w-20 shrink-0 overflow-hidden rounded-full border-4 shadow-md"
					>
						<GeneratedAvatar seed={team.name} class="h-full w-full" />
						{#if team.country}
							<div
								class="absolute bottom-0 right-0 h-7 w-7 overflow-hidden rounded-full border-2 border-white shadow-md dark:border-gray-800"
							>
								<CountryFlag
									country={String(team.country)}
									width={28}
									height={28}
									class="h-full w-full object-cover"
								/>
							</div>
						{/if}
					</div>

					<div class="min-w-0 flex-1 space-y-1">
						<h1 class="truncate text-3xl font-bold tracking-tight">{team.name}</h1>
						<div class="text-muted-foreground mt-1 flex flex-wrap items-center gap-4 text-sm">
							{#if team.email}
								<span class="flex items-center gap-1.5">
									<Mail class="h-4 w-4" />
									{team.email}
								</span>
							{/if}
							<span class="flex items-center gap-1.5">
								<Users class="h-4 w-4" />
								{team.members?.length ?? 0} members
							</span>
							{#if team.country}
								{@const countryData = countries.find(
									(c) => c.iso3 === String(team.country).toUpperCase()
								)}
								<span class="flex items-center gap-1.5">
									<Globe class="h-4 w-4" />
									{countryData?.name ?? team.country}
								</span>
							{/if}
						</div>
					</div>

					<!-- Key Stats -->
					<div
						class="border-border mt-2 flex w-full justify-between gap-8 border-t pt-4 sm:mt-0 sm:w-auto sm:justify-end sm:border-l sm:border-t-0 sm:pl-8 sm:pt-0"
					>
						<div class="text-center sm:text-right">
							<p class="text-muted-foreground text-[10px] font-bold uppercase tracking-wider">
								Total Score
							</p>
							<p class="font-mono text-3xl font-bold tabular-nums tracking-tight">
								{team.score?.toLocaleString() ?? 0}
							</p>
						</div>
						<div class="text-center sm:text-right">
							<p class="text-muted-foreground text-[10px] font-bold uppercase tracking-wider">
								Solves
							</p>
							<p class="font-mono text-3xl font-bold tabular-nums tracking-tight">
								{team.solves?.length ?? 0}
							</p>
						</div>
					</div>
				</div>

				{#if Array.isArray(team.badges) && team.badges.length > 0}
					<div class="border-border mt-6 flex flex-wrap items-center gap-3 border-t pt-4">
						<span class="text-muted-foreground text-xs font-semibold uppercase tracking-wider"
							>Badges</span
						>
						<div class="flex flex-wrap gap-2">
							{#each team.badges as badge}
								<Tooltip.Root>
									<Tooltip.Trigger>
										<span
											class="bg-primary/10 text-primary inline-flex cursor-help items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
										>
											<Medal class="mr-1 h-3 w-3" />
											{badge.name}
										</span>
									</Tooltip.Trigger>
									<Tooltip.Content>
										<p>{badge.description}</p>
									</Tooltip.Content>
								</Tooltip.Root>
							{/each}
						</div>
					</div>
				{/if}
			</Card.Content>
		</Card.Root>

		<!-- Tabs Control -->
		<div class="flex justify-center">
			<div
				class="bg-muted text-muted-foreground inline-flex h-10 items-center justify-center gap-1 rounded-lg p-1"
			>
				<button
					class="ring-offset-background focus-visible:ring-ring inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md px-6 py-1.5 text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 {activeTab ===
					'overview'
						? 'bg-background text-foreground shadow-sm'
						: 'hover:bg-background/50 hover:text-foreground'}"
					onclick={() => (activeTab = 'overview')}
				>
					<LayoutGrid class="h-4 w-4" />
					Overview
				</button>
				<button
					class="ring-offset-background focus-visible:ring-ring inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md px-6 py-1.5 text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 {activeTab ===
					'solves'
						? 'bg-background text-foreground shadow-sm'
						: 'hover:bg-background/50 hover:text-foreground'}"
					onclick={() => (activeTab = 'solves')}
				>
					<List class="h-4 w-4" />
					Solves
				</button>
			</div>
		</div>

		{#if activeTab === 'overview'}
			<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
				<!-- Row 1: Left: Category List (4 cols), Right: Stats (3 cols) -->

				<!-- Category Breakdown (List) -->
				<Card.Root class="bg-card border-0 shadow-sm lg:col-span-4">
					<Card.Header class="pb-2">
						<Card.Title class="text-muted-foreground text-sm font-medium uppercase tracking-wider"
							>Category Breakdown</Card.Title
						>
					</Card.Header>
					<Card.Content>
						{#if team?.solves && team.solves.length > 0}
							{@const categories = (() => {
								const map = new Map();
								for (const s of team.solves) map.set(s.category, (map.get(s.category) ?? 0) + 1);
								const total = [...map.values()].reduce((a, b) => a + b, 0) || 1;
								return [...map.entries()]
									.sort((a, b) => b[1] - a[1])
									.map(([cat, count]) => ({ cat, count, pct: Math.round((count / total) * 100) }));
							})()}
							<div class="grid gap-4 sm:grid-cols-2">
								{#each categories as c}
									<div class="space-y-1">
										<div class="flex justify-between text-xs font-medium">
											<span>{c.cat}</span>
											<span class="text-muted-foreground">{c.count} ({c.pct}%)</span>
										</div>
										<div class="bg-muted h-1.5 w-full overflow-hidden rounded-full">
											<div class="bg-primary h-full" style="width: {c.pct}%"></div>
										</div>
									</div>
								{/each}
							</div>
						{:else}
							<p class="text-muted-foreground text-sm">No solves yet.</p>
						{/if}
					</Card.Content>
				</Card.Root>

				<!-- Stats Card -->
				<Card.Root class="bg-card border-0 shadow-sm lg:col-span-3">
					<Card.Header class="pb-2">
						<Card.Title class="text-muted-foreground text-sm font-medium uppercase tracking-wider"
							>Team Status</Card.Title
						>
					</Card.Header>
					<Card.Content>
						<div class="grid grid-cols-2 gap-4">
							<div>
								<p class="text-muted-foreground text-xs uppercase">Score</p>
								<p class="font-mono text-xl font-bold">{team.score}</p>
							</div>
							<div>
								<p class="text-muted-foreground text-xs uppercase">Solves</p>
								<p class="font-mono text-xl font-bold">{team.solves?.length ?? 0}</p>
							</div>
						</div>
					</Card.Content>
				</Card.Root>

				<!-- Row 2: Radar Chart (Full Width) -->
				<!-- Temporarily commented out
				<Card.Root class="bg-card border-0 shadow-sm md:col-span-2 lg:col-span-7">
					<Card.Header class="pb-2">
						<Card.Title class="text-muted-foreground text-sm font-medium uppercase tracking-wider"
							>Skill Radar</Card.Title
						>
					</Card.Header>
					<Card.Content>
						<RadarChart solves={team.solves} />
					</Card.Content>
				</Card.Root>
				-->

				<!-- Row 3: Members (Full Width) -->
				<Card.Root class="bg-card border-0 shadow-sm md:col-span-2 lg:col-span-7">
					<Card.Header class="pb-2">
						<Card.Title class="text-muted-foreground text-sm font-medium uppercase tracking-wider"
							>Team Members</Card.Title
						>
					</Card.Header>
					<Card.Content>
						{#if team.members && team.members.length > 0}
							<div class="grid gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
								{#each team.members as member}
									<a
										href={`/account/${member.id}`}
										class="hover:bg-muted/50 flex items-center gap-3 rounded-lg p-2 transition-colors"
									>
										<GeneratedAvatar seed={member.name} size={40} />
										<div class="min-w-0">
											<div class="truncate text-sm font-medium">{member.name}</div>
											<div class="text-muted-foreground text-xs">{member.role}</div>
										</div>
									</a>
								{/each}
							</div>
						{:else}
							<p class="text-muted-foreground text-sm">No members found.</p>
						{/if}
					</Card.Content>
				</Card.Root>
			</div>
		{:else if activeTab === 'solves'}
			<Card.Root class="overflow-hidden border-0 shadow-sm">
				<Card.Content class="p-0">
					<SolveListTable {team} />
				</Card.Content>
			</Card.Root>
		{/if}
	</div>
{/if}

{#if team}
	<TeamEdit bind:open={teamEditOpen} {team} onupdated={() => teamQuery.refetch()} />
{/if}
