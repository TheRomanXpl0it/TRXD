<script lang="ts">
  import { params } from 'svelte-spa-router';
  import { user as authUser, authReady } from '$lib/stores/auth';
  import Spinner from '$lib/components/ui/spinner/spinner.svelte';
  import { Avatar } from 'flowbite-svelte';
  import { ShieldHalf, Globe, Users } from '@lucide/svelte';

  import TeamStatsCarousel from '$lib/components/team/team-stats-carousel.svelte';
  import SolveListTable from '$lib/components/team/team-scoreboard.svelte';
  import TeamMembers from '$lib/components/team/team-memberlist.svelte';
  import TeamJoinCreate from '$lib/components/team/team-join-create.svelte';
  import { getTeam } from '$lib/team';

  let loading = $state(false);
  let teamError: string | null = $state(null);
  let team: any = $state(null);

  // race guard
  let lastKey: string | null = $state(null);
  let reqSeq = $state(0);

  function normalizeKey(x: unknown): string | null {
    const s = String(x ?? '').trim();
    return s ? s : null;
  }

  // Strongest source of truth: URL path /team/:id
  function routeIdFromLocation(): string | null {
    if (typeof window === 'undefined') return null;
    const m = window.location.pathname.match(/^\/team\/([^/]+)\/?$/);
    return m?.[1] ? m[1] : null;
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

    // Prefer explicit :id from the URL path; then params store; finally fallback to user's team_id
    const hardRouteKey = normalizeKey(routeIdFromLocation());
    const softRouteKey = normalizeKey($params?.id);
    const fallbackKey  = normalizeKey($authUser?.team_id);

    const effectiveKey = hardRouteKey ?? softRouteKey ?? fallbackKey;

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
      <Spinner /> <p>Loading…</p>
    </div>
  {:else if teamError}
    <p class="text-red-600 text-sm">{teamError}</p>

  {:else if !team}
    <TeamJoinCreate/>

  {:else}
    <!-- Team header -->
    <div class="flex flex-row">
      {#if team.profileImage}
        <Avatar src={team.profileImage} class="h-24 w-24 mx-auto mb-4" />
      {:else}
        <Avatar class="h-24 w-24 mx-auto mb-4">
          <ShieldHalf />
        </Avatar>
      {/if}
    </div>

    <div class="flex flex-row justify-center text-center">
      <h2 class="text-2xl">{team.name}</h2>
    </div>

    <!-- Meta row -->
    <div class="flex flex-row justify-between mt-1">
      <div class="flex flex-row text-gray-500">
        {#if !team.country}
          <Globe class="mr-1" />
          Unknown
        {:else}
          <ShieldHalf class="mr-1" />
          {team.country}
        {/if}
      </div>
      <div class="flex flex-row text-gray-500">
        <Users class="mr-1" />
        {team.members?.length}
        {team.members?.length === 1 ? ' member' : ' members'}
      </div>
    </div>

    <!-- Stats -->
    <div class="flex flex-row items-center mt-15">
      <TeamStatsCarousel {team} />
    </div>

    <!-- Members -->
    <div class="flex flex-row items-center mt-15">
      <TeamMembers {team} />
    </div>

    <!-- Solves -->
    <div class="flex flex-row items-center mt-15">
      <SolveListTable team={team} />
    </div>
  {/if}
</div>
