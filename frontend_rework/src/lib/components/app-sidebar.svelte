<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { Avatar } from "flowbite-svelte";
  import { BugOutline } from "flowbite-svelte-icons";
  import { link } from "svelte-spa-router";
  import { HouseIcon, Joystick, ShieldHalf, Trophy, BookText } from "@lucide/svelte";

  export let user: { name?: string; profileImage?: string } | null = null;

  const items = [
    { title: "Home",       url: "/",           icon: HouseIcon },
    { title: "Scoreboard", url: "/scoreboard", icon: Trophy },
    { title: "Team",       url: "/team",       icon: ShieldHalf },
    { title: "Challenges", url: "/challenges", icon: Joystick },
    { title: "Writeups",   url: "/writeups",   icon: BookText },
  ];
</script>

<Sidebar.Root>
  <Sidebar.Content>
    <Sidebar.Group>
      <Sidebar.GroupLabel>TRXD</Sidebar.GroupLabel>
      <Sidebar.GroupContent>
        <Sidebar.Menu>
          {#each items as item (item.title)}
            <Sidebar.MenuItem>
              <Sidebar.MenuButton class="cursor-pointer">
                {#snippet child({ props })}
                  <!-- Spread first, then override href and add use:link for SPA navigation -->
                  <a {...props} href={item.url} use:link>
                    <item.icon />
                    <span>{item.title}</span>
                  </a>
                {/snippet}
              </Sidebar.MenuButton>
            </Sidebar.MenuItem>
          {/each}
        </Sidebar.Menu>
      </Sidebar.GroupContent>
    </Sidebar.Group>
  </Sidebar.Content>

  <Sidebar.Footer>
    {#if user}
      <a
        href="/account"
        use:link
        class="flex cursor-pointer items-center gap-3 rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
      >
        {#if user.profileImage}
          <Avatar src={user.profileImage} class="h-8 w-8" />
        {:else}
          <Avatar class="h-8 w-8">
            <BugOutline />
          </Avatar>
        {/if}
        <div class="min-w-0">
          <p class="truncate text-sm font-medium text-gray-700 dark:text-gray-100">{user.name}</p>
          <p class="truncate text-xs text-gray-500 dark:text-gray-400">@{user.name}</p>
        </div>
      </a>
    {:else}
      <a
        href="/signIn"
        use:link
        class="block rounded-lg p-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
      >
        Sign in
      </a>
    {/if}
  </Sidebar.Footer>
</Sidebar.Root>
