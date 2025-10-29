<script lang="ts">
	import { params } from 'svelte-spa-router';
	import { user as authUser, authReady, userMode, loadUser } from '$lib/stores/auth';
	import { push } from 'svelte-spa-router';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Avatar } from 'flowbite-svelte';
	import { ShieldHalf, Globe, Users, Edit } from '@lucide/svelte';
	import Icon from '@iconify/svelte';

	import TeamStatsCarousel from '$lib/components/team/team-stats-carousel.svelte';
	import SolveListTable from '$lib/components/team/team-scoreboard.svelte';
	import TeamMembers from '$lib/components/team/team-memberlist.svelte';
	import TeamJoinCreate from '$lib/components/team/team-join-create.svelte';
	import TeamEdit from '$lib/components/team/team-edit.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { getTeam } from '$lib/team';

	let loading = $state(false);
	let teamError: string | null = $state(null);
	let team: any = $state(null);
	let teamEditOpen = $state(false);
	const isOwnTeam = $derived($authUser && team && String($authUser.team_id) === String(team.id));

	// race guard
	let lastKey: string | null = $state(null);
	let reqSeq = $state(0);

	function normalizeKey(x: unknown): string | null {
		const s = String(x ?? '').trim();
		return s ? s : null;
	}


	async function loadTeamByKey(key: string) {
		const mySeq = ++reqSeq;
		loading = true;
		teamError = null;
		team = null;

		try {
			const apiKey = /^\d+$/.test(key) ? Number(key) : (key as any);
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
		<div class="flex flex-row items-center gap-2 text-gray-500">
			<Spinner />
			<p>Loading…</p>
		</div>
	{:else if teamError}
		<p class="text-sm text-red-600">{teamError}</p>
	{:else if !team}
		<TeamJoinCreate on:created={() => loadUser()} on:joined={()=>loadUser()} />
	{:else}
		<!-- Team header -->
		<div class="flex flex-row">
			{#if team?.image || team?.profileImage}
				<Avatar src={team.image ?? team.profileImage} class="mx-auto mb-4 h-24 w-24" />
			{:else}
				<Avatar class="mx-auto mb-4 h-24 w-24">
					<ShieldHalf class="h-12 w-12 text-gray-400" />
				</Avatar>
			{/if}
		</div>

		<div class="flex flex-row justify-center text-center">
			<h2 class="text-2xl">{team.name}</h2>
		</div>

		<!-- Meta row -->
		<div class="mt-1 flex flex-row justify-between">
			<div class="flex flex-row items-center text-gray-500">
				{#if !team.country}
					<Globe class="mr-1" />
					Unknown
				{:else}
					<Icon
						icon={`circle-flags:${String(team.country).toLowerCase()}`}
						width="20"
						height="20"
						class="mr-2"
					/>
					{team.country}
				{/if}
			</div>
			<div class="flex flex-col items-end gap-2">
				{#if isOwnTeam}
					<Button
						variant="outline"
						size="sm"
						onclick={() => (teamEditOpen = true)}
						class="flex items-center gap-2"
					>
						<Edit class="h-4 w-4" />
						Edit Team
					</Button>
				{/if}
				<div class="flex flex-row text-gray-500">
					<Users class="mr-1" />
					{team.members?.length}
					{team.members?.length === 1 ? ' member' : ' members'}
				</div>
			</div>
		</div>

		<!-- Description -->
		<div class="mt-3 flex flex-row items-center justify-center">
			<p class="text-center text-lg italic text-gray-500 dark:text-gray-400">
				{team.bio}
			</p>
		</div>

		<!-- Stats -->
		<div class="mt-15 flex flex-row items-center">
			<TeamStatsCarousel {team} />
		</div>

		<!-- Members -->
		<div class="mt-15 flex flex-row items-center">
			<TeamMembers {team} />
		</div>

		<!-- Solves -->
		<div class="mt-15 flex flex-row items-center">
			<SolveListTable {team} />
		</div>
	{/if}
</div>
<TeamEdit bind:open={teamEditOpen} {team} on:updated={() => loadTeamByKey(team.id)} />
