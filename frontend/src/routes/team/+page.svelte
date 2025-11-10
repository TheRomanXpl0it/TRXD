<script lang="ts">
	import { params } from 'svelte-spa-router';
	import { user as authUser, authReady, userMode, loadUser } from '$lib/stores/auth';
	import { push } from 'svelte-spa-router';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Avatar } from 'flowbite-svelte';
	import { ShieldHalf, Globe, Users, Edit, Award } from '@lucide/svelte';
	import Icon from '@iconify/svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	import TeamStatsCarousel from '$lib/components/team/TeamStatsCarousel.svelte';
	import SolveListTable from '$lib/components/team/TeamScoreboard.svelte';
	import TeamMembers from '$lib/components/team/TeamMemberlist.svelte';
	import TeamJoinCreate from '$lib/components/team/TeamJoinCreate.svelte';
	import TeamEdit from '$lib/components/team/TeamEdit.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { getTeam } from '$lib/team';

	let loading = $state(false);
	let teamError: string | null = $state(null);
	let team: any = $state(null);
	let teamEditOpen = $state(false);
	let activeTab = $state<'overview' | 'members' | 'solves'>('overview');
	const isOwnTeam = $derived($authUser && team && String($authUser.team_id) === String(team.id));

	// race guard
	let lastKey: string | null = $state(null);
	let reqSeq = $state(0);

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

	async function loadTeamByKey(key: string) {
		const mySeq = ++reqSeq;
		loading = true;
		teamError = null;

		try {
			const apiKey = validateId(key);
			if (apiKey === null) {
				throw new Error('Invalid team ID format');
			}
			const t = await getTeam(apiKey);
			if (mySeq !== reqSeq) return;
			team = t ?? null;
		} catch (e: any) {
			if (mySeq !== reqSeq) return;
			teamError = e?.message ?? 'Failed to load team';
			team = null;
		} finally {
			if (mySeq === reqSeq) loading = false;
		}
	}

	// React to: auth ready, URL param changes
	$effect(() => {
		if (!$authReady) return;

		// Redirect to 404 if userMode is enabled
		if ($userMode) {
			push('/not-found');
			return;
		}

		const softRouteKey = normalizeKey($params?.id);
		const fallbackKey = normalizeKey($authUser?.team_id);

		const effectiveKey = softRouteKey ?? fallbackKey;

		if (!effectiveKey) {
			lastKey = null;
			team = null;
			teamError = null;
			loading = false;
			return;
		}

		if (effectiveKey === lastKey) return;
		lastKey = effectiveKey;
		void loadTeamByKey(effectiveKey);
	});
</script>

<div>
	<p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Team</p>
	<hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
	<p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
		"None of us is as smart as all of us."
	</p>

	{#if !$authReady || loading}
		<div class="flex flex-col items-center justify-center py-12">
			<Spinner class="mb-4 h-8 w-8" />
			<p class="text-gray-600 dark:text-gray-400">Loading...</p>
		</div>
	{:else if teamError}
		<div class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-600 dark:border-red-800 dark:bg-red-950/20">
			<p class="font-semibold">Error loading team</p>
			<p class="text-sm">{teamError}</p>
		</div>
	{:else if !team}
		<TeamJoinCreate oncreated={() => loadUser()} onjoined={() => loadUser()} />
	{:else}
		<!-- Header -->
		<div class="mb-8 flex items-start justify-between pb-6">
			<div class="flex min-w-0 flex-1 items-center gap-4">
				{#if team?.image || team?.profileImage}
					<Avatar src={team.image ?? team.profileImage} class="h-16 w-16" />
				{:else}
					<Avatar class="h-16 w-16">
						<ShieldHalf class="h-8 w-8" />
					</Avatar>
				{/if}
				<div class="min-w-0 flex-1">
					<h2 class="truncate text-2xl font-bold">{team.name}</h2>
					<div class="mt-1 flex items-center gap-3 text-sm">
						{#if team.country}
							<div class="flex items-center gap-1">
								<Icon icon={`circle-flags:${String(team.country).toLowerCase()}`} width="14" height="14" />
								<span>{team.country}</span>
							</div>
						{/if}
						<div class="flex items-center gap-1">
							<Users class="h-3.5 w-3.5" />
							<span>{team.members?.length} {team.members?.length === 1 ? 'member' : 'members'}</span>
						</div>
					</div>
				</div>
			</div>
			{#if isOwnTeam}
				<Button variant="outline" size="sm" onclick={() => (teamEditOpen = true)}>
					<Edit class="h-4 w-4" />
				</Button>
			{/if}
		</div>

		<!-- Tabs -->
		<div class="mb-8 border-b border-black dark:border-white">
			<div class="flex">
				<button
					class="cursor-pointer border-b-2 px-4 py-3 text-sm font-medium text-black transition-colors dark:text-white {activeTab === 'overview' ? 'border-black dark:border-white' : 'border-transparent hover:bg-black/5 dark:hover:bg-white/5'}"
					onclick={() => activeTab = 'overview'}
				>
					Overview
				</button>
				<button
					class="cursor-pointer border-b-2 px-4 py-3 text-sm font-medium text-black transition-colors dark:text-white {activeTab === 'members' ? 'border-black dark:border-white' : 'border-transparent hover:bg-black/5 dark:hover:bg-white/5'}"
					onclick={() => activeTab = 'members'}
				>
					Members
				</button>
				<button
					class="cursor-pointer border-b-2 px-4 py-3 text-sm font-medium text-black transition-colors dark:text-white {activeTab === 'solves' ? 'border-black dark:border-white' : 'border-transparent hover:bg-black/5 dark:hover:bg-white/5'}"
					onclick={() => activeTab = 'solves'}
				>
					Solves
				</button>
			</div>
		</div>

		<!-- Tab Content -->
		{#if activeTab === 'overview'}
			<div class="space-y-8">
				{#if team.bio}
					<div>
						<h3 class="mb-2 text-lg font-semibold">About</h3>
						<p>{team.bio}</p>
					</div>
				{/if}

				<!-- Stats Grid -->
				<div>
					<h3 class="mb-4 text-lg font-semibold">Statistics</h3>
					<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
						<div class="rounded bg-gray-50 dark:bg-gray-800/50 p-4 shadow-sm">
							<div class="text-sm">Total Points</div>
							<p class="mt-1 text-2xl font-semibold">{team?.score ?? 0}</p>
						</div>
						<div class="rounded bg-gray-50 dark:bg-gray-800/50 p-4 shadow-sm">
							<div class="text-sm">Members</div>
							<p class="mt-1 text-2xl font-semibold">{team?.members?.length ?? 0}</p>
							<p class="text-xs opacity-70">{team?.members?.filter((m: any) => (m?.score ?? 0) > 0).length ?? 0} active</p>
						</div>
						<div class="rounded bg-gray-50 dark:bg-gray-800/50 p-4 shadow-sm">
							<div class="text-sm">Solves</div>
							<p class="mt-1 text-2xl font-semibold">{team?.solves?.length ?? 0}</p>
						</div>
						<div class="rounded bg-gray-50 dark:bg-gray-800/50 p-4 shadow-sm">
							<div class="text-sm">Badges</div>
							<p class="mt-1 text-2xl font-semibold">{team?.badges?.length ?? 0}</p>
						</div>
					</div>
				</div>

				<!-- Latest Activity. This should probably be removed -->
				{#if team?.solves && team.solves.length > 0}
					{@const lastSolve = [...team.solves].sort((a, b) => +new Date(b.timestamp) - +new Date(a.timestamp))[0]}
					{@const timeSince = (iso) => {
						if (!iso) return '-';
						const sec = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 1000));
						const h = Math.floor(sec / 3600);
						const m = Math.floor((sec % 3600) / 60);
						const s = sec % 60;
						if (h > 0) return `${h}h ${m}m`;
						if (m > 0) return `${m}m ${s}s`;
						return `${s}s`;
					}}
					<div>
						<h3 class="mb-4 text-lg font-semibold">Latest Activity</h3>
						<div class="rounded bg-gray-50 dark:bg-gray-800/50 p-4 shadow-sm">
							<div class="text-sm">Most recent solve</div>
							<div class="mt-2 flex items-center justify-between">
								<div>
									<p class="text-lg font-medium">{lastSolve.name}</p>
									<span class="mt-1 inline-flex items-center rounded-full bg-black/5 dark:bg-white/10 px-2.5 py-0.5 text-xs font-medium">
										{lastSolve.category}
									</span>
								</div>
								<div class="text-right">
									<p class="text-2xl font-semibold">{timeSince(lastSolve.timestamp)}</p>
									<p class="text-xs opacity-70">ago</p>
								</div>
							</div>
							<p class="mt-2 text-xs opacity-70">
								{new Date(lastSolve.timestamp).toLocaleString()}
							</p>
						</div>
					</div>
				{/if}

				<!-- Category Breakdown. Improve this in future updates -->
				{#if team?.solves && team.solves.length > 0}
					{@const categories = (() => {
						const map = new Map();
						for (const s of team.solves) map.set(s.category, (map.get(s.category) ?? 0) + 1);
						const total = [...map.values()].reduce((a, b) => a + b, 0) || 1;
						return [...map.entries()]
							.sort((a, b) => b[1] - a[1])
							.map(([cat, count]) => ({ cat, count, pct: Math.round((count / total) * 100) }));
					})()}
					<div>
						<h3 class="mb-4 text-lg font-semibold">Category Breakdown</h3>
						<div class="space-y-3">
							{#each categories as c}
								<div>
									<div class="flex items-center justify-between text-sm">
										<span class="font-medium">{c.cat}</span>
										<span class="opacity-70">{c.count} ({c.pct}%)</span>
									</div>
									<div class="mt-1 h-2 w-full rounded bg-gray-100 dark:bg-gray-700">
										<div class="h-full rounded bg-gray-900 dark:bg-gray-200" style="width:{Math.max(3, Math.min(100, c.pct))}%"></div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}

				<!-- Badges Section -->
				{#if team?.badges && team.badges.length > 0}
					<div>
						<h3 class="mb-4 text-lg font-semibold">Badges</h3>
						<div class="flex flex-wrap gap-3">
							{#each team.badges as badge}
								<Tooltip.Root>
									<Tooltip.Trigger>
										<div class="flex items-center gap-2 rounded bg-gray-50 dark:bg-gray-800/50 px-3 py-2 shadow-sm">
											{#if badge.icon}
												<span class="text-xl">{badge.icon}</span>
											{/if}
											<span class="text-sm font-medium">{badge.name}</span>
										</div>
									</Tooltip.Trigger>
									{#if badge.description}
										<Tooltip.Content>
											<p>{badge.description}</p>
										</Tooltip.Content>
									{/if}
								</Tooltip.Root>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		{:else if activeTab === 'members'}
			<TeamMembers {team} />
		{:else if activeTab === 'solves'}
			<SolveListTable {team} />
		{/if}
	{/if}
</div>
<TeamEdit bind:open={teamEditOpen} {team} onupdated={() => loadTeamByKey(team.id)} />
