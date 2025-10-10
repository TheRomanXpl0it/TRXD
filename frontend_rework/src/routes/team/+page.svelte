<script lang="ts">
  import { getTeam } from "@/team";
  import { onMount } from "svelte";
  import { user, authReady } from '$lib/stores/auth';
  import Spinner from '@/components/ui/spinner/spinner.svelte';
  import { Avatar } from "flowbite-svelte";
  import { ShieldHalf } from "@lucide/svelte";
  import { Globe, Users } from '@lucide/svelte';
  import TeamStatsCarousel from "@/components/team/team-stats-carousel.svelte";
  import SolveListTable from '$lib/components/team/team-scoreboard.svelte';
  import TeamMembers from '$lib/components/team/team-memberlist.svelte';
  
  let team: any = null;
  let loading = false;
  let teamError: string | null = null;

  onMount(async () => {
    if (!$authReady) return;
    if (!$user?.team_id) return;           

    loading = true;
    try {
      team = await getTeam($user.team_id);
      console.log(team);
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
        <div class="flex flex-row items-center align-center gap-2 text-gray-500">
            <Spinner/>
            Loading...
        </div>
    {:else}
        {#if !team}
            You're not in a team.
        {:else if teamError}
            <p class="text-red-600 text-sm">{teamError}</p>
        {:else}
            <!-- First two rows, team propic and name underneat-->
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
                <h2 class="text-2xl">
                    {team.name}
                </h2>
            </div>
            
        <!-- Third row, team nationality and member count -->
            <div class="flex flex-row justify-between mt-1">
                <div class="flex flex-row text-gray-500">
                    {#if !team.nationality}
                        <Globe class="mr-1"/>
                        Unknown
                    {:else}
                        <ShieldHalf class="mr-1"/>
                        {team.nationality}
                    {/if}
                </div>
                <div class="flex flex-row text-gray-500">
                    <Users class="mr-1"/>
                    {team.members?.length} {team.members?.length === 1 ? 'member' : 'members'}
                </div>
            </div>
            
        <!-- Fourth row, carousels-->
        <div class="flex flex-row items-center mt-15">
            <TeamStatsCarousel {team} />
        </div>
        
        <!-- Fifth row, TeamMembers-->
        <div class="flex flex-row items-center mt-15">
            <TeamMembers {team} />
        </div>
        
        <!-- Sixth row, solvelist-->
        <div class="flex flex-row items-center mt-15">
            <SolveListTable solves={team.solves} />
        </div>
        
                    
        {/if}
    {/if}
</div>