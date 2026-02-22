<script module lang="ts">
	declare const __GIT_HASH__: string;
</script>

<script lang="ts">
	import { authState } from '$lib/stores/auth';

	const gitHash = __GIT_HASH__;
	const isAdmin = $derived(authState.user?.role === 'Admin');
</script>

<svelte:head>
	<link rel="preload" as="image" href="/trx.svg" />
</svelte:head>

<!-- NOTE: Is this really useful? -->

<div class="flex flex-col items-center justify-center py-8">
	<!-- Hero Section -->
	<div class="mb-8 text-center">
		<img
			src="/trx.svg"
			alt="TRXD Logo"
			width="256"
			height="256"
			loading="eager"
			fetchpriority="high"
			class="mx-auto mb-8 h-auto w-48 sm:w-64"
		/>

		<h1 class="mb-4 text-4xl font-bold text-gray-900 sm:text-5xl dark:text-white">
			Welcome to TRXD
		</h1>

		<p class="mx-auto mb-6 max-w-2xl text-xl text-gray-600 dark:text-gray-300">
			A platform for hackers and cybersecurity enthusiasts
		</p>

		<p class="mx-auto mb-8 max-w-xl text-base italic text-gray-500 dark:text-gray-400">
			"ALZATI SUBITO"
		</p>

		<div class="flex flex-wrap justify-center gap-4">
			<a
				href="/challenges"
				class="bg-primary text-primary-foreground hover:bg-primary/90 rounded-lg px-6 py-3 font-semibold transition-colors"
			>
				View Challenges
			</a>
			<a
				href="/scoreboard"
				class="bg-secondary text-secondary-foreground hover:bg-secondary/80 rounded-lg px-6 py-3 font-semibold transition-colors"
			>
				Scoreboard
			</a>
		</div>
	</div>

	{#if isAdmin}
		<div class="mt-8 text-center text-xs text-gray-400 dark:text-gray-600">
			<span title={gitHash}>hash: {gitHash}</span>
		</div>
	{/if}
</div>
