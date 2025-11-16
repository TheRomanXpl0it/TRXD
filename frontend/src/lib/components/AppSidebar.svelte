<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Avatar } from 'flowbite-svelte';
	import { BugOutline } from 'flowbite-svelte-icons';
	import { link, push } from 'svelte-spa-router';
	import {
		HouseIcon,
		Joystick,
		ShieldHalf,
		Trophy,
		BookText,
		LogOut,
		LogIn,
		Settings,
		Users
	} from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button';
	import { getUserData } from '$lib/user';
	import { useSidebar } from '$lib/components/ui/sidebar/context.svelte.js';

	const sidebar = useSidebar();

	let { user = null, userMode = false } = $props<{
		user: {
			id?: number | string;
			name?: string;
			image?: string;
			profileImage?: string;
		} | null;
		userMode: boolean;
	}>();


	type NavItem = { title: string; url: string; icon: typeof HouseIcon };
	const baseItems: NavItem[] = [
		{ title: 'Home', url: '/', icon: HouseIcon },
		{ title: 'Challenges', url: '/challenges', icon: Joystick },
		{ title: 'Scoreboard', url: '/scoreboard', icon: Trophy },
		//{ title: "Writeups",   url: "/writeups",   icon: BookText },
	];

	// Accounts item shown only when logged in
	const accountsItem: NavItem = { title: 'Accounts', url: '/accounts', icon: Users };

	// Teams item shown only when userMode is false and logged in
	const teamsItem: NavItem = { title: 'Teams', url: '/teams', icon: ShieldHalf };

	// Team item shown only when userMode is false and logged in
	const teamItem: NavItem = { title: 'Team', url: '/team', icon: ShieldHalf };

	// Admin-only visibility
	const isAdmin = $derived((user as any)?.role === 'Admin');
	const configsItem: NavItem = { title: 'Configs', url: '/configs', icon: Settings };

	// Combine items based on userMode and role
	const allItems: NavItem[] = $derived(
		[...baseItems]
			.concat(user ? [accountsItem] : [])
			.concat(user && !userMode ? [teamsItem, teamItem] : [])
			.concat(isAdmin ? [configsItem] : [])
	);

	// Enrich user data to ensure profile image is available in `image`
	let enrichedUser = $state<any>(null);
	const displayImage = $derived(enrichedUser?.image ?? user?.image ?? user?.profileImage ?? null);

	$effect(() => {
		const id = user?.id
		if (!id) {
			enrichedUser = null;
			return;
		}
		(async () => {
			try {
				enrichedUser = await getUserData(id);
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
									<a {...props} href={item.url} use:link onclick={() => {
										// Close mobile menu on navigation
										if (sidebar.isMobile) {
											sidebar.setOpenMobile(false);
										}
									}}>
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
				onclick={() => {
					// Close mobile menu on navigation
					if (sidebar.isMobile) {
						sidebar.setOpenMobile(false);
					}
				}}
			>
				{#if displayImage}
					<Avatar src={displayImage} class="h-8 w-8 shrink-0" />
				{:else}
					<Avatar class="h-8 w-8 shrink-0">
						<BugOutline />
					</Avatar>
				{/if}

				<div class="min-w-0 flex-1">
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
