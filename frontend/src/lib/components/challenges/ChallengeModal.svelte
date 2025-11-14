<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '@/components/ui/button';
	import { Download, Droplet } from '@lucide/svelte';
	import { PenSolid, TrashBinSolid, UserEditSolid } from 'flowbite-svelte-icons';
	import { toast } from 'svelte-sonner';
	import InstanceControls from './InstanceControls.svelte';
	import FlagSubmission from './FlagSubmission.svelte';
	import Markdown from '$lib/components/Markdown.svelte';
	import { config } from '$lib/env';

	let {
		open = $bindable(false),
		challenge,
		countdown = 0,
		isAdmin = false,
		onEdit,
		onDelete,
		onSolved,
		onCountdownUpdate,
		onOpenSolves
	}: {
		open: boolean;
		challenge: any;
		countdown?: number;
		isAdmin?: boolean;
		onEdit?: (challenge: any) => void;
		onDelete?: (challenge: any) => void;
		onSolved?: () => void;
		onCountdownUpdate?: (id: string | number, newCountdown: number) => void;
		onOpenSolves?: () => void;
	} = $props();

	function copyToClipboard(text: string) {
		if (typeof navigator === 'undefined') return;
		navigator.clipboard
			.writeText(text)
			.then(() => toast.success('Copied to clipboard!'))
			.catch(() => toast.error('Failed to copy to clipboard.'));
	}

	const connectionString = $derived(
		`${challenge?.host ?? ''}${challenge?.port ? `:${challenge.port}` : ''}`
	);
</script>

<Dialog.Root bind:open>
	<Dialog.Content
		class="max-w-[95vw] sm:max-w-[800px] max-h-[95vh] overflow-y-auto p-4 sm:p-6"
		aria-describedby="challenge-description"
	>
		<Dialog.Header class="pb-4 sm:pb-6">
			<div class="min-w-0 flex-1">
				<Dialog.Title class="text-2xl sm:text-3xl font-bold break-words mb-3 pr-8">
					{challenge?.name ?? 'Challenge'}
				</Dialog.Title>

				<!-- Tags & Metadata -->
				<div class="flex flex-wrap items-center gap-2">
					{#if challenge?.tags && challenge.tags.length > 0}
						<div role="list" aria-label="Challenge tags" class="contents">
							{#each challenge.tags as tag}
								<span
									class="inline-flex items-center rounded-full bg-black/5 dark:bg-white/10 px-2.5 py-0.5 text-xs font-medium"
									role="listitem"
								>
									{tag}
								</span>
							{/each}
						</div>
					{/if}

					{#if challenge?.difficulty}
						<span
							class="inline-flex items-center rounded-full bg-black/5 dark:bg-white/10 px-2.5 py-0.5 text-xs font-medium"
							aria-label="Difficulty: {challenge.difficulty}"
						>
							{challenge.difficulty}
						</span>
					{/if}

					{#if challenge?.solves === 0}
						<span class="inline-flex items-center gap-1 text-xs font-medium">
							<Droplet class="h-3 w-3 text-red-500" aria-hidden="true" />
							<span class="opacity-70">0 solves</span>
						</span>
					{:else if challenge?.solves}
						<button
							type="button"
							onclick={onOpenSolves}
							class="text-xs font-medium opacity-70 hover:opacity-100 hover:underline focus:outline-none focus:underline"
							aria-label="View {challenge.solves} solve{challenge.solves === 1 ? '' : 's'}"
						>
							{challenge.solves} {challenge.solves === 1 ? 'solve' : 'solves'}
						</button>
					{/if}
				</div>

				<!-- Authors -->
				{#if challenge?.authors && challenge.authors.length > 0}
					<div class="mt-2 flex items-center gap-1 text-xs font-medium opacity-70">
						<UserEditSolid class="h-3 w-3" aria-hidden="true" />
						<span>
							By {#each challenge.authors as author, i (author)}{author}{i <
								challenge.authors.length - 1 ? ', ' : ''}{/each}
						</span>
					</div>
				{/if}

				<!-- Admin Controls -->
				{#if isAdmin}
					<div class="mt-3 flex items-center gap-2" role="group" aria-label="Admin actions">
						<Button
							variant="outline"
							size="sm"
							class="cursor-pointer"
							onclick={() => onEdit?.(challenge)}
							aria-label="Edit challenge"
						>
							<PenSolid class="h-3.5 w-3.5" aria-hidden="true" />
							<span>Edit</span>
						</Button>
						<Button
							variant="outline"
							size="sm"
							class="text-red-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950 cursor-pointer"
							onclick={() => onDelete?.(challenge)}
							aria-label="Delete challenge"
						>
							<TrashBinSolid class="h-3.5 w-3.5" aria-hidden="true" />
							<span>Delete</span>
						</Button>
					</div>
				{/if}
			</div>
			<Dialog.Description id="challenge-description" class="sr-only">
				Challenge details and submission form
			</Dialog.Description>
		</Dialog.Header>

		<!-- Description -->
		<section class="mb-6" aria-labelledby="description-heading">
			<h3 id="description-heading" class="text-sm font-semibold mb-2 opacity-70">Description</h3>
			{#if challenge?.description}
				<Markdown content={challenge.description} class="text-base leading-relaxed" />
			{:else}
				<div class="text-base leading-relaxed opacity-60">No description available.</div>
			{/if}
		</section>

		<!-- Attachments -->
		{#if challenge?.attachments && challenge.attachments.length > 0}
			<section class="mb-6" aria-labelledby="attachments-heading">
				<h3 id="attachments-heading" class="text-sm font-semibold mb-3 opacity-70">Attachments</h3>
				<div class="flex flex-wrap gap-2">
					{#each challenge.attachments as attachment}
						<a
							href={attachment}
							target="_blank"
							rel="noopener noreferrer"
							class="inline-flex h-9 items-center justify-center gap-2 rounded-md border border-input bg-background px-3 text-sm font-medium shadow-sm transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
							aria-label="Download {attachment.split('/').pop()}"
						>
							<Download class="h-4 w-4" aria-hidden="true" />
							<span>{attachment.split('/').pop()}</span>
						</a>
					{/each}
				</div>
			</section>
		{/if}

		<!-- Connection Info (only for non-instance challenges) -->
		{#if challenge?.host && !challenge.instance}
			<section class="mb-6" aria-labelledby="connection-heading">
				<h3 id="connection-heading" class="text-sm font-semibold mb-3 opacity-70">Connection</h3>
				<button
					type="button"
					onclick={() => copyToClipboard(connectionString)}
					class="inline-flex h-10 items-center gap-2 rounded-md bg-gray-100 dark:bg-gray-800 px-4 font-mono text-sm font-medium transition-colors hover:bg-gray-200 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary"
					aria-label="Copy connection string: {connectionString}"
				>
					<span>{connectionString}</span>
				</button>
			</section>
		{/if}

		<!-- Instance Controls -->
		{#if challenge?.instance}
			<InstanceControls {challenge} {countdown} {onCountdownUpdate} />
		{/if}

		<!-- Submit Flag -->
		<FlagSubmission {challenge} {onSolved} />
	</Dialog.Content>
</Dialog.Root>
