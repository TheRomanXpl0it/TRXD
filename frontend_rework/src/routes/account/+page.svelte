<script lang="ts">
  import { onMount } from 'svelte';
  import { user, authReady } from '$lib/stores/auth';
  import Spinner from '@/components/ui/spinner/spinner.svelte';
  import { Avatar } from 'flowbite-svelte';
  import { getTeam } from '$lib/team';
  import { getUserData } from '$lib/user';
  import { BugOutline } from 'flowbite-svelte-icons';
  import { Globe } from '@lucide/svelte';

  let team: any = null;
  let userVerboseData: any = null;
  let loading = false;
  let teamError: string | null = null;

  onMount(async () => {
    if (!$authReady) return;
    if (!$user?.team_id) return;           

    loading = true;
    try {
      team = await getTeam($user.team_id);
      userVerboseData = await getUserData($user.id)
    } catch (e: any) {
      teamError = e?.message ?? 'Failed to load team';
    } finally {
      loading = false;
      
    }
  });
</script>

<div>
{#if !$authReady}
  <div class="flex flex-row items-center align-center gap-2 text-gray-500">
    <Spinner /><p>Loading…</p>
  </div>
{:else if !$user}
  <p>You’re not signed in.</p>
{:else}
  <p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Account</p>
  <hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
  <p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
    "You didn't wake up today to be mediocre."
  </p>

  <div class="justify-center">
    {#if $user.profileImage}
      <Avatar src={$user.profileImage} class="h-24 w-24 mx-auto mb-4" />
    {:else}
      <Avatar class="h-24 w-24 mx-auto mb-4" >
        <BugOutline />
      </Avatar>
    {/if}
  </div>
  
  <div class="text-center">
    <h2 class="text-2xl font-semibold">{$user.name}</h2>
    <p class="text-gray-500">@{$user.name}</p>
      
  </div>
  
  <div class="flex flex-row justify-between mt-1">
    <div class="flex flex-row text-gray-500">
        {#if !userVerboseData?.nationality}
            <Globe class="mr-1"/>
            Unknown
        {:else}
            <img src={`https://flagcdn.com/16x12/${userVerboseData.nationality.toLowerCase()}.png`} alt={userVerboseData.nationality} class="h-4 w-6 me-2" />
            {userVerboseData.nationality}
        {/if}
    </div>
    <div class="flex flex-row text-gray-500">
        Joined on {new Date(userVerboseData?.joined_at).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })}
    
    </div>
  </div>
{/if}
</div>
