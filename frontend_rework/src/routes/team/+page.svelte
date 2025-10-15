<script lang="ts">
  import { onMount } from 'svelte';
  import { user, authReady } from '$lib/stores/auth';
  import Spinner from '$lib/components/ui/spinner/spinner.svelte';
  import { Avatar } from 'flowbite-svelte';
  import { ShieldHalf } from '@lucide/svelte';
  import { Globe, Users } from '@lucide/svelte';

  import TeamStatsCarousel from '$lib/components/team/team-stats-carousel.svelte';
  import SolveListTable from '$lib/components/team/team-scoreboard.svelte';
  import TeamMembers from '$lib/components/team/team-memberlist.svelte';

  import TeamJoinCreate from '$lib/components/team/team-join-create.svelte';
  import { getTeam } from '$lib/team';

  let team: any = $state(null);
  let loading = $state(false);
  let teamError: string | null = $state(null);

  onMount(async () => {
    if (!$authReady) return;
    if (!$user?.team_id) return;

    loading = true;
    try {
      team = await getTeam($user.team_id);
    } catch (e: any) {
      teamError = e?.message ?? 'Failed to load team';
    } finally {
      loading = false;
    }
  });
</script>

<div>
  <p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Team</p>
  <hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
  <p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
    "None of us is as smart as all of us."
  </p>

  {#if loading}
    <div class="flex flex-row items-center gap-2 text-gray-500">
      <Spinner /> Loading...
    </div>
  {:else}
    {#if !team}
      <!-- Join/Create section extracted into its own component -->
      <TeamJoinCreate/>

    {:else if teamError}
      <p class="text-red-600 text-sm">{teamError}</p>

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
  {/if}
</div>
