<script lang="ts">
	import { CheckCircle, Container } from '@lucide/svelte';

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
		class={`hover:ring-primary/20 focus-visible:ring-primary/50 flex w-full cursor-pointer items-center gap-4 rounded-lg p-3 text-left transition-colors hover:ring-2 focus-visible:outline-none focus-visible:ring-2 ${challenge.solved ? 'bg-green-50/80 dark:bg-green-950/30' : 'bg-gray-50 dark:bg-gray-800/50'} ${challenge.hidden ? 'bg-amber-50 ring-2 ring-inset ring-amber-400 dark:bg-amber-950/20 dark:ring-amber-600' : 'shadow-sm'}`}
		{onclick}
		aria-label="View details for {challenge.name}{challenge.solved ? ' (solved)' : ''}"
	>
		<!-- Points -->
		<div class="flex w-20 shrink-0 items-center gap-2">
			<span class="text-sm font-semibold" aria-label="{challenge.points} points">
				{challenge.points}
			</span>
			{#if challenge.solved}
				<CheckCircle class="h-4 w-4 text-green-500" aria-label="Solved" />
			{/if}
		</div>

		<!-- Challenge Name -->
		<div class="min-w-0 flex-1">
			<div class="flex items-center gap-2">
				<span class="truncate font-semibold">{challenge.name}</span>
				{#if challenge.instance}
					<Container class="h-4 w-4 shrink-0 opacity-50" aria-label="Instance-based challenge" />
				{/if}
			</div>
		</div>

		<!-- Tags -->
		<div class="flex shrink-0 flex-wrap gap-1.5" role="list" aria-label="Tags">
			{#each challenge.tags as tag}
				<span
					class="inline-flex items-center rounded-full bg-black/5 px-2.5 py-0.5 text-xs font-medium dark:bg-white/10"
					role="listitem"
				>
					{tag}
				</span>
			{/each}
		</div>

		{#if challenge.instance && countdown > 0}
			<span
				class="inline-flex shrink-0 items-center rounded-md bg-blue-100 px-2.5 py-0.5 text-xs font-semibold text-blue-800 transition-colors dark:border dark:border-blue-400 dark:bg-transparent dark:text-blue-400"
				aria-label="Instance expires in {fmtTimeLeft(countdown)}"
			>
				{fmtTimeLeft(countdown)}
			</span>
		{/if}
	</button>
{:else}
	<!-- Grid View -->
	<button
		type="button"
		class={`hover:ring-primary/20 focus-visible:ring-primary/50 group relative w-full cursor-pointer overflow-hidden rounded-lg p-5 text-left transition-all hover:ring-2 focus-visible:outline-none focus-visible:ring-2 ${challenge.solved ? 'bg-green-50/80 dark:bg-green-950/30' : 'bg-gray-50 dark:bg-gray-800/50'} shadow-md ${challenge.hidden ? 'bg-amber-50 ring-2 ring-inset ring-amber-400 dark:bg-amber-950/20 dark:ring-amber-600' : ''}`}
		{onclick}
		aria-label="View details for {challenge.name}, {challenge.points} points{challenge.solved
			? ', solved'
			: ''}"
	>
		<div class="mb-3 flex items-start justify-between">
			<p class="pr-2 text-lg font-semibold text-gray-900 dark:text-white">{challenge.name}</p>
		</div>

		<div class="mb-4 flex flex-wrap gap-1.5" role="list" aria-label="Tags">
			{#each challenge.tags as tag}
				<span
					class="inline-flex items-center rounded-full bg-black/5 px-2.5 py-0.5 text-xs font-medium dark:bg-white/10"
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
					<CheckCircle class="h-5 w-5 text-green-500" aria-label="Solved" />
				{/if}
			</div>

			{#if challenge.instance && countdown > 0}
				<span
					class="inline-flex shrink-0 items-center rounded-md bg-blue-100 px-2.5 py-0.5 text-xs font-semibold text-blue-800 transition-colors dark:border dark:border-blue-400 dark:bg-transparent dark:text-blue-400"
					aria-label="Instance expires in {fmtTimeLeft(countdown)}"
				>
					{fmtTimeLeft(countdown)}
				</span>
			{/if}
		</div>
	</button>
{/if}
