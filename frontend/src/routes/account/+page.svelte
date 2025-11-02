<script lang="ts">
	import { onMount } from 'svelte';
	import { params } from 'svelte-spa-router';
	import { user, authReady, userMode } from '$lib/stores/auth';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Avatar } from 'flowbite-svelte';
	import { BugOutline } from 'flowbite-svelte-icons';
	import { Globe } from '@lucide/svelte';
	import Icon from '@iconify/svelte';
	import { getTeam } from '$lib/team';
	import { getUserData } from '$lib/user';
	import Solvelist from '$lib/components/account/account-scoreboard.svelte';
	import UserUpdate from '$lib/components/account/account-edit.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Edit } from '@lucide/svelte';

	let loading = $state(false);
	let teamError: string | null = $state(null);
	let editSheetOpen = $state(false);

	let userVerboseData: any = $state(null);
	let team: any = $state(null);

	// to avoid races
	let reqSeq = $state(0);

	function normalizeKey(x: unknown): string | null {
		const s = String(x ?? '').trim();
		return s ? s : null;
	}

	async function loadUserAndTeamByKey(key: string, inUserMode: boolean) {
		const mySeq = ++reqSeq;

		loading = true;
		teamError = null;
		userVerboseData = null;
		team = null;

		try {
			const apiKey = /^\d+$/.test(key) ? Number(key) : (key as any);
			const userData = await getUserData(apiKey);
			if (mySeq !== reqSeq) return;
			userVerboseData = userData ?? null;

			// if we're in userMode, stops here
			if (inUserMode) {
				team = null;
				return;
			}

			if (userVerboseData?.team_id != null) {
				const t = await getTeam(userVerboseData.team_id);
				if (mySeq !== reqSeq) return;
				team = t ?? null;
			} else {
				team = null;
			}
		} catch (e: any) {
			if (mySeq !== reqSeq) return;
			teamError = e?.message ?? 'Failed to load team';
			userVerboseData = null;
			team = null;
		} finally {
			if (mySeq === reqSeq) loading = false;
		}
	}

	// manual wiring so scoreboard updates don't re-run the fetch
	onMount(() => {
		let currentKey: string | null = null;
		let currentMode = false;
		let currentReady = false;

		const unsubParams = params.subscribe(($p) => {
			const routeKey = normalizeKey($p?.id);
			maybeReload(routeKey, currentMode, currentReady);
		});

		const unsubAuth = authReady.subscribe(($ready) => {
			currentReady = $ready;
			// re-evaluate with last known route/user
			const routeKey = normalizeKey(get(params)?.id); // small helper below
			maybeReload(routeKey, currentMode, currentReady);
		});

		const unsubUser = user.subscribe(($u) => {
			// if no route param, we fallback to user.id
			const routeKey = normalizeKey(get(params)?.id);
			const fallbackKey = normalizeKey($u?.id);
			const effectiveKey = routeKey ?? fallbackKey;
			maybeReload(effectiveKey, currentMode, get(authReady));
		});

		const unsubMode = userMode.subscribe(($m) => {
			currentMode = !!$m;
			const routeKey = normalizeKey(get(params)?.id);
			const fallbackKey = normalizeKey(get(user)?.id);
			const effectiveKey = routeKey ?? fallbackKey;
			maybeReload(effectiveKey, currentMode, get(authReady));
		});

		// helper to re-load only when the identity or mode really changed
		function maybeReload(nextKey: string | null, nextMode: boolean, ready: boolean) {
			if (!ready) return;
			// fallback to current user if no route key
			if (!nextKey) {
				const fallbackKey = normalizeKey(get(user)?.id);
				nextKey = fallbackKey;
			}
			if (!nextKey) return;

			// if key OR mode changed → load
			if (nextKey !== currentKey || nextMode !== currentMode) {
				currentKey = nextKey;
				currentMode = nextMode;
				void loadUserAndTeamByKey(nextKey, nextMode);
			}
		}

		// small util to get latest value of a store
		function get<T>(store: any): T {
			let v: T;
			store.subscribe((x: T) => (v = x))();
			return v!;
		}

		return () => {
			unsubParams();
			unsubAuth();
			unsubUser();
			unsubMode();
		};
	});

	// derived: is it my profile?
	const isOwnProfile = $derived(
		$user && userVerboseData && String($user.id) === String(userVerboseData.id)
	);
</script>

{#if !$authReady || loading}
	<div class="flex flex-row items-center gap-2 text-gray-500">
		<Spinner />
		<p>Loading…</p>
	</div>
{:else if !$user && !$params?.id}
	<p>You’re not signed in.</p>
{:else}
	<p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Account</p>
	<hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
	<p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
		"You didn't wake up today to be mediocre."
	</p>

	<div class="justify-center">
		{#if userVerboseData?.image}
			<Avatar src={userVerboseData.image} class="mx-auto mb-4 h-24 w-24" />
		{:else}
			<Avatar class="mx-auto mb-4 h-24 w-24">
				<BugOutline class="h-12 w-12 text-gray-400" />
			</Avatar>
		{/if}
	</div>

	<div class="text-center">
		<h2 class="text-2xl font-semibold">{userVerboseData?.name ?? '—'}</h2>
		<p class="text-gray-500">@{userVerboseData?.name ?? '—'}</p>
	</div>

	<div class="mt-1 flex flex-row justify-between">
		<div class="flex flex-row items-center text-gray-500">
			{#if !userVerboseData?.country}
				<Globe class="mr-1" /> Unknown
			{:else}
				<Icon
					icon={`circle-flags:${String(userVerboseData.country).toLowerCase()}`}
					width="20"
					height="20"
					class="mr-2"
				/>
				{userVerboseData.country}
			{/if}
		</div>
		<div class="flex flex-col items-end gap-2">
			{#if isOwnProfile}
				<Button
					variant="outline"
					size="sm"
					onclick={() => (editSheetOpen = true)}
					class="flex items-center gap-2"
				>
					<Edit class="h-4 w-4" />
					Edit Profile
				</Button>
			{/if}
			<div class="flex flex-row text-gray-500">
				Joined on {userVerboseData?.joined_at
					? new Date(userVerboseData.joined_at).toLocaleDateString(undefined, {
							year: 'numeric',
							month: 'long',
							day: 'numeric'
						})
					: '—'}
			</div>
		</div>
	</div>

	{#if $userMode}
		<!-- solves are in user payload -->
		<div class="mt-6">
			<Solvelist solves={Array.isArray(userVerboseData?.solves) ? userVerboseData.solves : []} />
		</div>
	{:else}
		{#if teamError}
			<p class="mt-4 text-sm text-red-500">{teamError}</p>
		{/if}
		{#if !team}
			<p class="mt-6 text-gray-500">This user has not joined a team yet.</p>
		{:else}
			<div class="mt-6">
				<Solvelist
					solves={Array.isArray(team?.solves)
						? team.solves.filter((s: any) => String(s.user_id) === String(userVerboseData?.id))
						: []}
				/>
			</div>
		{/if}
	{/if}
{/if}

<UserUpdate
	bind:open={editSheetOpen}
	user={userVerboseData}
	on:updated={() => {
		const effectiveId = userVerboseData?.id;
		if (effectiveId) {
			// reload with CURRENT mode, but do it once, not via $effect
			void loadUserAndTeamByKey(String(effectiveId), !!$userMode);
		}
	}}
/>
