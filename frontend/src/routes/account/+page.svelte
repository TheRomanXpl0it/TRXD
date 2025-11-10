<script lang="ts">
	import { params } from 'svelte-spa-router';
	import { user, authReady, userMode } from '$lib/stores/auth';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Avatar } from 'flowbite-svelte';
	import { BugOutline } from 'flowbite-svelte-icons';
	import { Globe } from '@lucide/svelte';
	import Icon from '@iconify/svelte';
	import { getTeam } from '$lib/team';
	import { getUserData } from '$lib/user';
	import Solvelist from '$lib/components/account/AccountScoreboard.svelte';
	import UserUpdate from '$lib/components/account/AccountEdit.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Edit } from '@lucide/svelte';
	import ErrorMessage from '$lib/components/ui/error-message.svelte';
	import { createQuery } from '@tanstack/svelte-query';

	let editSheetOpen = $state(false);
	let activeTab = $state<'overview' | 'solves'>('overview');

	// Derive currentUserId from params/user - no need for $effect!
	const currentUserId = $derived.by(() => {
		if (!$authReady) return null;
		
		const routeKey = normalizeKey($params?.id);
		const fallbackKey = normalizeKey($user?.id);
		const effectiveKey = routeKey ?? fallbackKey;
		
		if (!effectiveKey) return null;
		
		return validateId(effectiveKey);
	});

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

		// Reject anything else (to avoid future issues in the future)
		return null;
	}

	const userQuery = createQuery(() => ({
		queryKey: ['user', currentUserId],
		queryFn: () => getUserData(currentUserId!),
		enabled: currentUserId !== null,
		staleTime: 10_000,
		refetchInterval: 10_000
	}));

	const userVerboseData = $derived(userQuery.data ?? null);

	// Derive currentTeamId from user data
	const currentTeamId = $derived.by(() => {
		if ($userMode || !userVerboseData) return null;
		return userVerboseData?.team_id ?? null;
	});

	const teamQuery = createQuery(() => ({
		queryKey: ['team', currentTeamId],
		queryFn: () => getTeam(currentTeamId!),
		enabled: currentTeamId !== null && !$userMode,
		staleTime: 10_000,
		refetchInterval: 10_000
	}));

	const team = $derived(teamQuery.data ?? null);
	const loading = $derived(userQuery.isLoading || teamQuery.isLoading);
	const teamError = $derived(userQuery.error?.message ?? teamQuery.error?.message ?? null);

	// derived: is it my profile?
	const isOwnProfile = $derived(
		$user && userVerboseData && String($user.id) === String(userVerboseData.id)
	);
</script>

{#if !$authReady || loading}
	<div>
		<p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Account</p>
		<hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
		<p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
			"You didn't wake up today to be mediocre."
		</p>

		<div class="flex flex-col items-center justify-center py-12">
			<Spinner class="mb-4 h-8 w-8" />
			<p class="text-gray-600 dark:text-gray-400">Loading...</p>
		</div>
	</div>
{:else if !$user && !$params?.id}
	<p>You're not signed in.</p>
{:else}
	<div>
		<p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Account</p>
		<hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
		<p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
			"You didn't wake up today to be mediocre."
		</p>

		<!-- Header -->
		<div class="mb-8 flex items-start justify-between pb-6">
			<div class="flex items-center gap-4 min-w-0 flex-1">
				{#if userVerboseData?.image}
					<Avatar src={userVerboseData.image} class="h-16 w-16 shrink-0" />
				{:else}
					<Avatar class="h-16 w-16 shrink-0">
						<BugOutline class="h-8 w-8" />
					</Avatar>
				{/if}
				<div class="min-w-0 flex-1">
					<h2 class="text-2xl font-bold truncate">{userVerboseData?.name ?? '-'}</h2>
					<div class="mt-1 flex items-center gap-3 text-sm">
						{#if userVerboseData?.country}
							<div class="flex items-center gap-1">
								<Icon icon={`circle-flags:${String(userVerboseData.country).toLowerCase()}`} width="14" height="14" />
								<span>{userVerboseData.country}</span>
							</div>
						{:else}
							<div class="flex items-center gap-1">
								<Globe class="h-3.5 w-3.5" />
								<span>Unknown</span>
							</div>
						{/if}
						{#if userVerboseData?.joined_at}
							<span>•</span>
							<span>Joined {new Date(userVerboseData.joined_at).toLocaleDateString(undefined, {
								year: 'numeric',
								month: 'short'
							})}</span>
						{/if}
					</div>
				</div>
			</div>
			{#if isOwnProfile}
				<Button variant="outline" size="sm" onclick={() => (editSheetOpen = true)} class="shrink-0">
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
				{#if userVerboseData?.bio}
					<div>
						<h3 class="mb-2 text-lg font-semibold">About</h3>
						<p>{userVerboseData.bio}</p>
					</div>
				{/if}

				<!-- Stats Grid -->
				<div>
					<h3 class="mb-4 text-lg font-semibold">Statistics</h3>
					<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						<div class="rounded bg-gray-50 dark:bg-gray-800/50 p-4 shadow-sm">
							<div class="text-sm">Total Points</div>
							<p class="mt-1 text-2xl font-semibold">{userVerboseData?.score ?? 0}</p>
						</div>
						<div class="rounded bg-gray-50 dark:bg-gray-800/50 p-4 shadow-sm">
							<div class="text-sm">Solves</div>
							<p class="mt-1 text-2xl font-semibold">
								{#if $userMode}
									{userVerboseData?.solves?.length ?? 0}
								{:else}
									{team?.solves?.filter((s) => String(s.user_id) === String(userVerboseData?.id)).length ?? 0}
								{/if}
							</p>
						</div>
						{#if !$userMode}
							<div class="rounded bg-gray-50 dark:bg-gray-800/50 p-4 shadow-sm">
								<div class="text-sm">Team</div>
								<p class="mt-1 truncate text-lg font-semibold">{team?.name ?? 'No team'}</p>
							</div>
						{/if}
					</div>
				</div>
			</div>
		{:else if activeTab === 'solves'}
			{#if $userMode}
				<Solvelist solves={Array.isArray(userVerboseData?.solves) ? userVerboseData.solves : []} />
			{:else}
				{#if teamError}
					<ErrorMessage title="Error loading team data" message={teamError} />
				{:else if !team}
					<p class="text-gray-500">This user has not joined a team yet.</p>
				{:else}
					<Solvelist
						solves={Array.isArray(team?.solves)
							? team.solves.filter((s: any) => String(s.user_id) === String(userVerboseData?.id))
							: []}
					/>
				{/if}
			{/if}
		{/if}
	</div>
{/if}

<UserUpdate
	bind:open={editSheetOpen}
	user={userVerboseData}
	onupdated={() => {
		// Invalidate queries to refresh data
		userQuery.refetch();
		if (!$userMode) {
			teamQuery.refetch();
		}
	}}
/>
