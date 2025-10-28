<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Avatar } from 'flowbite-svelte';
	import { BugOutline } from 'flowbite-svelte-icons';
	import { link, push } from 'svelte-spa-router';
	import { HouseIcon, Joystick, ShieldHalf, Trophy, BookText, LogOut, LogIn } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button';
	import { getUserData } from '$lib/user';

	let { user = null, userMode = false } = $props<{
		user: {
			id?: number | string;
			name?: string;
			image?: string;
			profileImage?: string;
		} | null;
		userMode: boolean;
	}>();

	console.log(user);

	const baseItems = [
		{ title: 'Home', url: '/', icon: HouseIcon },
		{ title: 'Scoreboard', url: '/scoreboard', icon: Trophy },
		{ title: 'Challenges', url: '/challenges', icon: Joystick }
		//{ title: "Writeups",   url: "/writeups",   icon: BookText },
	];

	// Team item shown only when userMode is false
	const teamItem = { title: 'Team', url: '/team', icon: ShieldHalf };

	// Combine items based on userMode
	const allItems = $derived(userMode ? baseItems : [...baseItems, teamItem]);

	// Enrich user data to ensure profile image is available in `image`
	let enrichedUser = $state<any>(null);
	const displayImage = $derived(enrichedUser?.image ?? user?.image ?? user?.profileImage ?? null);

	$effect(() => {
		const id = (user as any)?.id;
		if (!id) {
			enrichedUser = null;
			return;
		}
		const apiKey = /^\d+$/.test(String(id)) ? Number(id) : id;
		(async () => {
			try {
				enrichedUser = await getUserData(apiKey);
			} catch {
				// leave enrichedUser as null on failure
			}
		})();
	});
</script>

<Sidebar.Root>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupLabel>TRXD</Sidebar.GroupLabel>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					{#each allItems as item (item.title)}
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
				{#if displayImage}
					<Avatar src={displayImage} class="h-8 w-8" />
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
					class="pointer-events-none ml-auto cursor-pointer opacity-0 transition-opacity
                 duration-150 focus:pointer-events-auto
                 focus:opacity-100 group-hover:pointer-events-auto group-hover:opacity-100"
					variant="outline"
					title="Log out"
					aria-label="Log out"
					onclick={() => {
						push('/signOut');
					}}
				>
					<LogOut />
				</Button>
			</a>
		{:else}
			<a
				href="/signIn"
				use:link
				class="flex flex-row rounded-lg p-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
			>
				<LogIn class="mr-3" />
				Sign in
			</a>
		{/if}
	</Sidebar.Footer>
</Sidebar.Root>
