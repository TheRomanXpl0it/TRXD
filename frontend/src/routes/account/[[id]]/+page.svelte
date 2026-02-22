<script lang="ts">
	import { page } from '$app/stores';
	import { authState } from '$lib/stores/auth';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { getTeam } from '$lib/team';
	import { getUserData } from '$lib/user';
	import Solvelist from '$lib/components/account/AccountScoreboard.svelte';
	import UserUpdate from '$lib/components/account/AccountEdit.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Edit, Shield, Trophy, Target, Users, Globe, Mail, Medal } from '@lucide/svelte';
	import ErrorMessage from '$lib/components/ui/error-message.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import GeneratedAvatar from '$lib/components/ui/avatar/generated-avatar.svelte';
	import RadarChart from '$lib/components/RadarChart.svelte';
	import CountryFlag from '$lib/components/ui/country-flag.svelte';
	import countries from '$lib/data/countries.json';

	let editSheetOpen = $state(false);
	let activeTab = $state<'overview' | 'solves'>('overview');

	// Derive currentUserId from params/user
	const currentUserId = $derived.by(() => {
		if (!authState.ready || !authState.user) return null;

		const routeKey = normalizeKey($page.params.id);
		if (routeKey) return validateId(routeKey);

		// Fallback to own profile
		const fallbackKey = authState.userMode ? authState.user.team_id : authState.user.id;

		return fallbackKey ? validateId(String(fallbackKey)) : null;
	});

	function normalizeKey(x: unknown): string | null {
		const s = String(x ?? '').trim();
		return s ? s : null;
	}

	function validateId(key: string): number | null {
		if (/^\d+$/.test(key)) {
			const num = Number(key);
			return num > 0 ? num : null;
		}
		return null;
	}

	const userQuery = createQuery(() => ({
		queryKey: ['user', currentUserId, authState.userMode],
		queryFn: () => getUserData(currentUserId!, authState.userMode),
		enabled: currentUserId !== null,
		staleTime: 10_000
	}));

	const userVerboseData = $derived(userQuery.data ?? null);

	// Derive currentTeamId from user data
	const currentTeamId = $derived.by(() => {
		if (authState.userMode || !userVerboseData) return null;
		return userVerboseData?.team_id ?? null;
	});

	const teamQuery = createQuery(() => ({
		queryKey: ['team', currentTeamId],
		queryFn: () => getTeam(currentTeamId!),
		enabled: currentTeamId !== null && !authState.userMode,
		staleTime: 10_000
	}));

	const team = $derived(teamQuery.data ?? null);

	const loading = $derived(
		(currentUserId !== null && !userQuery.isError && (userQuery.isFetching || !userVerboseData)) ||
			(currentTeamId !== null && !teamQuery.isError && (teamQuery.isFetching || !team))
	);
	const teamError = $derived(userQuery.error?.message ?? teamQuery.error?.message ?? null);

	const isOwnProfile = $derived.by(() => {
		if (!authState.user || !userVerboseData) return false;
		if (authState.userMode) {
			return String(authState.user.team_id) === String(userVerboseData.id);
		}
		return String(authState.user.id) === String(userVerboseData.id);
	});

	const solveCount = $derived.by(() => {
		if (authState.userMode) return userVerboseData?.solves?.length ?? 0;
		return (
			team?.solves?.filter((s) => String(s.user_id) === String(userVerboseData?.id)).length ?? 0
		);
	});

	const displayBadges = $derived.by(() => {
		if (authState.userMode && userVerboseData?.badges && userVerboseData.badges.length > 0)
			return userVerboseData.badges;
		if (!authState.userMode && team?.badges && team.badges.length > 0) return team.badges;
		return [];
	});
</script>

{#if !authState.ready || loading}
	<div class="mx-auto max-w-5xl space-y-8 py-10">
		<!-- Skeleton Header -->
		<div class="space-y-4">
			<div class="bg-muted h-8 w-48 animate-pulse rounded"></div>
			<Separator />
			<div class="flex flex-col items-center justify-center py-12">
				<Spinner class="mb-4 h-8 w-8" />
				<p class="text-muted-foreground">Loading profile...</p>
			</div>
		</div>
	</div>
{:else if !authState.user}
	<div class="mx-auto max-w-5xl py-10">
		<Card.Root>
			<Card.Content class="py-10 text-center">
				<p>You're not signed in.</p>
			</Card.Content>
		</Card.Root>
	</div>
{:else if !userVerboseData}
	<div class="mx-auto max-w-5xl py-10">
		<Card.Root>
			<Card.Content class="py-10 text-center">
				<p>Profile not found.</p>
			</Card.Content>
		</Card.Root>
	</div>
{:else}
	<div class="mx-auto max-w-5xl space-y-8 py-10">
		<div class="flex items-center justify-between">
			<h2 class="text-3xl font-bold tracking-tight">Profile</h2>
			{#if isOwnProfile}
				<Button variant="outline" size="sm" onclick={() => (editSheetOpen = true)} class="gap-2">
					<Edit class="h-4 w-4" />
					Edit Profile
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
						<GeneratedAvatar seed={userVerboseData?.name ?? 'user'} class="h-full w-full" />
						{#if userVerboseData?.country}
							<div
								class="absolute bottom-0 right-0 h-7 w-7 overflow-hidden rounded-full border-2 border-white shadow-md dark:border-gray-800"
							>
								<CountryFlag
									country={String(userVerboseData.country)}
									width={28}
									height={28}
									class="h-full w-full object-cover"
								/>
							</div>
						{/if}
						{#if userVerboseData?.country}
							{@const countryData = countries.find(
								(c) => c.iso3 === String(userVerboseData.country).toUpperCase()
							)}
							<span class="flex items-center gap-1.5">
								<Globe class="h-3.5 w-3.5" />
								{countryData?.name ?? userVerboseData.country}
							</span>
						{/if}
					</div>

					<div class="min-w-0 flex-1 space-y-1">
						<div class="flex items-center gap-2">
							<h1 class="truncate text-2xl font-bold sm:text-3xl">
								{userVerboseData?.name ?? 'Unknown User'}
							</h1>
							{#if userVerboseData?.role && userVerboseData.role !== 'User'}
								<span
									class="bg-primary/10 text-primary inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium uppercase tracking-wide"
								>
									{userVerboseData.role}
								</span>
							{/if}
						</div>

						<div class="text-muted-foreground flex flex-wrap items-center gap-x-4 gap-y-2 text-sm">
							{#if userVerboseData?.email}
								<span class="flex items-center gap-1.5">
									<Mail class="h-3.5 w-3.5" />
									{userVerboseData.email}
								</span>
							{/if}
							{#if userVerboseData?.joined_at}
								<span class="flex items-center gap-1.5">
									Joined {new Date(userVerboseData.joined_at).toLocaleDateString(undefined, {
										year: 'numeric',
										month: 'long'
									})}
								</span>
							{/if}
							{#if !authState.userMode && team}
								<span
									class="bg-muted/50 text-foreground/80 flex items-center gap-1.5 rounded-md px-2 py-0.5 font-medium"
								>
									<Users class="h-3.5 w-3.5" />
									{team.name}
								</span>
							{/if}
							{#if userVerboseData?.country}
								{@const countryData = countries.find(
									(c) => c.iso3 === String(userVerboseData.country).toUpperCase()
								)}
								<span class="flex items-center gap-1.5">
									<Globe class="h-3.5 w-3.5" />
									{countryData?.name ?? userVerboseData.country}
								</span>
							{/if}
						</div>
					</div>

					<!-- Key Stats (Right Side on Desktop) -->
					<div
						class="border-border mt-2 flex w-full gap-4 border-t pt-4 sm:mt-0 sm:w-auto sm:gap-8 sm:border-l sm:border-t-0 sm:pl-8 sm:pt-0"
					>
						<div class="text-center sm:text-right">
							<p class="text-muted-foreground text-xs font-semibold uppercase tracking-wider">
								Points
							</p>
							<p class="font-mono text-2xl font-semibold tabular-nums tracking-tight sm:text-3xl">
								{userVerboseData?.score?.toLocaleString() ?? 0}
							</p>
						</div>
						<div class="text-center sm:text-right">
							<p class="text-muted-foreground text-xs font-semibold uppercase tracking-wider">
								Solves
							</p>
							<p class="font-mono text-2xl font-semibold tabular-nums tracking-tight sm:text-3xl">
								{solveCount}
							</p>
						</div>
					</div>
				</div>

				{#if displayBadges.length > 0}
					<div class="border-border mt-6 flex flex-wrap items-center gap-3 border-t pt-4">
						<span class="text-muted-foreground text-xs font-semibold uppercase tracking-wider"
							>Badges</span
						>
						<div class="flex flex-wrap gap-2">
							{#each displayBadges as badge}
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
					class="ring-offset-background focus-visible:ring-ring inline-flex items-center justify-center whitespace-nowrap rounded-md px-6 py-1.5 text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 {activeTab ===
					'overview'
						? 'bg-background text-foreground shadow-sm'
						: 'hover:bg-background/50 hover:text-foreground'}"
					onclick={() => (activeTab = 'overview')}
				>
					Overview
				</button>
				<button
					class="ring-offset-background focus-visible:ring-ring inline-flex items-center justify-center whitespace-nowrap rounded-md px-6 py-1.5 text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 {activeTab ===
					'solves'
						? 'bg-background text-foreground shadow-sm'
						: 'hover:bg-background/50 hover:text-foreground'}"
					onclick={() => (activeTab = 'solves')}
				>
					Solves
				</button>
			</div>
		</div>

		<!-- Tab Content -->
		{#if activeTab === 'overview'}
			<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
				<!-- Stat Cards -->
				<Card.Root class="bg-card border-0 shadow-sm">
					<Card.Header class="pb-2">
						<Card.Title class="text-muted-foreground text-sm font-medium uppercase tracking-wider"
							>Total Score</Card.Title
						>
					</Card.Header>
					<Card.Content>
						<div class="font-mono text-2xl font-bold">
							{userVerboseData?.score?.toLocaleString() ?? 0} pts
						</div>
					</Card.Content>
				</Card.Root>

				<Card.Root class="bg-card border-0 shadow-sm">
					<Card.Header class="pb-2">
						<Card.Title class="text-muted-foreground text-sm font-medium uppercase tracking-wider"
							>Challenges Solved</Card.Title
						>
					</Card.Header>
					<Card.Content>
						<div class="font-mono text-2xl font-bold">{solveCount}</div>
						<p class="text-muted-foreground mt-1 text-xs">Across all categories</p>
					</Card.Content>
				</Card.Root>

				{#if !authState.userMode && team}
					<Card.Root class="bg-card border-0 shadow-sm">
						<Card.Header class="pb-2">
							<Card.Title class="text-muted-foreground text-sm font-medium uppercase tracking-wider"
								>Team Status</Card.Title
							>
						</Card.Header>
						<Card.Content>
							<div class="truncate text-lg font-bold">{team.name}</div>
							<p class="text-muted-foreground mt-1 text-xs">
								{team.members?.length ?? 0} members
							</p>
						</Card.Content>
					</Card.Root>
				{/if}

				<!-- Radar Chart -->
				<!-- Temporarily commented out
				<Card.Root class="bg-card border-0 shadow-sm md:col-span-2 lg:col-span-3">
					<Card.Header class="pb-2">
						<Card.Title class="text-muted-foreground text-sm font-medium uppercase tracking-wider"
							>Skill Radar</Card.Title
						>
					</Card.Header>
					<Card.Content>
						<RadarChart solves={userVerboseData?.solves} />
					</Card.Content>
				</Card.Root>
				-->
			</div>
		{:else if activeTab === 'solves'}
			<Card.Root class="overflow-hidden border-0 shadow-sm">
				<Card.Content class="p-0">
					{#if authState.userMode}
						<Solvelist
							solves={Array.isArray(userVerboseData?.solves) ? userVerboseData.solves : []}
						/>
					{:else if teamError}
						<ErrorMessage title="Error loading team data" message={teamError} />
					{:else if !team}
						<div class="text-muted-foreground p-8 text-center">
							This user has not joined a team yet.
						</div>
					{:else}
						<Solvelist
							solves={Array.isArray(team?.solves)
								? team.solves.filter((s: any) => String(s.user_id) === String(userVerboseData?.id))
								: []}
						/>
					{/if}
				</Card.Content>
			</Card.Root>
		{/if}
	</div>
{/if}

{#if userVerboseData}
	<UserUpdate
		bind:open={editSheetOpen}
		user={userVerboseData}
		onupdated={() => {
			userQuery.refetch();
			if (!authState.userMode) teamQuery.refetch();
		}}
	/>
{/if}
