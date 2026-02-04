<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppSidebar from '$lib/components/AppSidebar.svelte';
	import { Button } from '@/components/ui/button';
	import { MoonIcon, SunIcon } from '@lucide/svelte';
	import { toggleMode } from 'mode-watcher';

	interface Props {
		user: any;
		userMode: boolean;
		children?: import('svelte').Snippet;
	}

	let { user, userMode, children }: Props = $props();
</script>

<Sidebar.Provider>
	<!-- pass the actual store value -->
	<AppSidebar {user} {userMode} />

	<!-- Mobile/Tablet header with menu trigger (fixed) -->
	<div
		class="bg-background fixed left-0 right-0 top-0 z-50 flex items-center justify-between overflow-visible border-b border-gray-200 p-3 lg:hidden dark:border-gray-700"
	>
		<Sidebar.Trigger size="icon" variant="ghost" class="relative z-10 cursor-pointer">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="24"
				height="24"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<line x1="4" x2="20" y1="12" y2="12"></line>
				<line x1="4" x2="20" y1="6" y2="6"></line>
				<line x1="4" x2="20" y1="18" y2="18"></line>
			</svg>
		</Sidebar.Trigger>

		<span class="text-lg font-bold text-gray-900 dark:text-white">TRXD</span>

		<Button
			size="icon"
			onclick={toggleMode}
			variant="ghost"
			class="cursor-pointer"
			aria-label="Toggle theme"
		>
			<SunIcon class="h-5 w-5 rotate-0 scale-100 !transition-all dark:-rotate-90 dark:scale-0" />
			<MoonIcon
				class="absolute h-5 w-5 rotate-90 scale-0 !transition-all dark:rotate-0 dark:scale-100"
			/>
		</Button>
	</div>

	<main class="flex min-h-screen w-full flex-col pt-[57px] lg:pt-0">
		<div class="router-content mx-3 flex-1 pb-24 lg:mx-6 lg:pb-16">
			{@render children?.()}
		</div>

		<!-- Desktop theme toggle -->
		<Button
			size="icon"
			onclick={toggleMode}
			variant="outline"
			class="fixed bottom-6 right-6 z-50 hidden p-2 shadow-lg lg:flex"
			aria-label="Toggle theme"
		>
			<SunIcon
				class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 !transition-all dark:-rotate-90 dark:scale-0"
			/>
			<MoonIcon
				class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 !transition-all dark:rotate-0 dark:scale-100"
			/>
		</Button>
	</main>
</Sidebar.Provider>
