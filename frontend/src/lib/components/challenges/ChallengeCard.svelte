<script lang="ts">
	import { CheckCircleSolid } from 'flowbite-svelte-icons';
	import { Badge } from 'flowbite-svelte';
	import { Container } from '@lucide/svelte';

	let {
		challenge,
		compactView = false,
		countdown = 0,
		onclick
	}: {
		challenge: any;
		compactView?: boolean;
		countdown?: number;
		onclick: () => void;
	} = $props();

	function fmtTimeLeft(total: number): string {
		if (!total || total < 0) total = 0;
		const h = Math.floor(total / 3600);
		const m = Math.floor((total % 3600) / 60);
		const s = Math.floor(total % 60);
		if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
		if (m > 0) return `${m}:${String(s).padStart(2, '0')}`;
		return `${s}`;
	}
</script>

{#if compactView}
	<!-- Compact View -->
	<button
		type="button"
		class={`flex w-full items-center gap-4 rounded-lg p-3 text-left transition-colors hover:ring-2 hover:ring-primary/20 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50 cursor-pointer ${challenge.solved ? 'bg-green-50/80 dark:bg-green-950/30' : 'bg-gray-50 dark:bg-gray-800/50'} ${challenge.hidden ? 'bg-amber-50 dark:bg-amber-950/20 ring-2 ring-amber-400 dark:ring-amber-600 ring-inset' : 'shadow-sm'}`}
		{onclick}
		aria-label="View details for {challenge.name}{challenge.solved ? ' (solved)' : ''}"
	>
		<!-- Points -->
		<div class="flex items-center gap-2 shrink-0 w-20">
			<span class="font-semibold text-sm" aria-label="{challenge.points} points">
				{challenge.points}
			</span>
			{#if challenge.solved}
				<CheckCircleSolid class="h-4 w-4 text-green-500" aria-label="Solved" />
			{/if}
		</div>

		<!-- Challenge Name -->
		<div class="flex-1 min-w-0">
			<div class="flex items-center gap-2">
				<span class="font-semibold truncate">{challenge.name}</span>
				{#if challenge.instance}
					<Container class="h-4 w-4 shrink-0 opacity-50" aria-label="Instance-based challenge" />
				{/if}
			</div>
		</div>

		<!-- Tags -->
		<div class="flex flex-wrap gap-1.5 shrink-0" role="list" aria-label="Tags">
			{#each challenge.tags as tag}
				<span
					class="inline-flex items-center rounded-full bg-black/5 dark:bg-white/10 px-2.5 py-0.5 text-xs font-medium"
					role="listitem"
				>
					{tag}
				</span>
			{/each}
		</div>

		<!-- Instance Timer -->
		{#if challenge.instance && countdown > 0}
			<Badge color="blue" class="text-xs shrink-0" aria-label="Instance expires in {fmtTimeLeft(countdown)}">
				{fmtTimeLeft(countdown)}
			</Badge>
		{/if}
	</button>
{:else}
	<!-- Grid View -->
	<button
		type="button"
		class={`group relative w-full overflow-hidden rounded-lg p-5 text-left transition-all hover:ring-2 hover:ring-primary/20 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50 cursor-pointer ${challenge.solved ? 'bg-green-50/80 dark:bg-green-950/30' : 'bg-gray-50 dark:bg-gray-800/50'} shadow-md ${challenge.hidden ? 'bg-amber-50 dark:bg-amber-950/20 ring-2 ring-amber-400 dark:ring-amber-600 ring-inset' : ''}`}
		{onclick}
		aria-label="View details for {challenge.name}, {challenge.points} points{challenge.solved ? ', solved' : ''}"
	>
		<div class="mb-3 flex items-start justify-between">
			<p class="text-lg font-semibold text-gray-900 dark:text-white pr-2">{challenge.name}</p>
			{#if challenge.instance}
				<Container class="h-5 w-5 shrink-0 text-gray-500 dark:text-gray-400" aria-label="Instance-based challenge" />
			{/if}
		</div>

		<div class="mb-4 flex flex-wrap gap-1.5" role="list" aria-label="Tags">
			{#each challenge.tags as tag}
				<span
					class="inline-flex items-center rounded-full bg-black/5 dark:bg-white/10 px-2.5 py-0.5 text-xs font-medium"
					role="listitem"
				>
					{tag}
				</span>
			{/each}
		</div>

		<div class="mt-auto flex items-center justify-between">
			<div class="flex items-center gap-2">
				<span
					class={`text-sm font-semibold ${challenge.solved ? 'text-green-600 dark:text-green-400' : ''}`}
					aria-label="{challenge.points} points"
				>
					{challenge.points} pts
				</span>
				{#if challenge.solved}
					<CheckCircleSolid class="h-5 w-5 text-green-500" aria-label="Solved" />
				{/if}
			</div>

			{#if challenge.instance && countdown > 0}
				<Badge color="blue" class="text-xs" aria-label="Instance expires in {fmtTimeLeft(countdown)}">
					{fmtTimeLeft(countdown)}
				</Badge>
			{/if}
		</div>
	</button>
{/if}
