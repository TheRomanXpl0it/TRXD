<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { Avatar } from "flowbite-svelte";
  import { BugOutline } from "flowbite-svelte-icons";
  import { link, push } from "svelte-spa-router";
  import { HouseIcon, Joystick, ShieldHalf, Trophy, BookText, LogOut, LogIn } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button";

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
        class="group flex cursor-pointer items-center gap-3 rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
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

        <!-- Hidden until parent <a> is hovered (or button focused) -->
        <Button
          class="ml-auto opacity-0 pointer-events-none transition-opacity duration-150
                 group-hover:opacity-100 group-hover:pointer-events-auto
                 focus:opacity-100 focus:pointer-events-auto cursor-pointer"
          variant="outline"
          title="Log out"
          aria-label="Log out"
          onclick={()=>{push('/signOut')}}
        >
          <LogOut />
        </Button>
      </a>
    {:else}
      <a
        href="/signIn"
        use:link
        class="rounded-lg p-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700 flex flex-row"
      >
          <LogIn class="mr-3"/>
        Sign in
      </a>
    {/if}
  </Sidebar.Footer>
</Sidebar.Root>
