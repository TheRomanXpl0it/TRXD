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
		challenge = $bindable(),
		countdown = 0,
		isAdmin = false,
		onEdit,
		onDelete,
		onSolved,
		onCountdownUpdate,
		onOpenSolves,
		onInstanceChange
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
		onOpenSolves?: () => void;
		onInstanceChange?: (challenge?: any) => void;
	} = $props();

	function copyToClipboard(text: string) {
		if (typeof navigator === 'undefined') return;
		navigator.clipboard
			.writeText(text)
			.then(() => toast.success('Copied to clipboard!'))
			.catch(() => toast.error('Failed to copy to clipboard.'));
	}

	const connectionString = $derived.by(() => {
		const h = challenge?.instance
			? (challenge?.instance_host ?? challenge?.host ?? '')
			: (challenge?.host ?? '');
		const p = challenge?.instance ? challenge?.instance_port : challenge?.port;
		let str = p ? `${h}:${p}` : h;
		if (str && challenge?.conn_type === 'HTTP' && !str.startsWith('http')) {
			str = `http://${str}`;
		} else if (str && challenge?.conn_type === 'HTTPS' && !str.startsWith('http')) {
			str = `https://${str}`;
		}
		return str;
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content
		class="max-h-[95vh] max-w-[95vw] overflow-y-auto p-4 sm:max-w-[800px] sm:p-6"
		aria-describedby="challenge-description"
	>
		<Dialog.Header class="pb-4 sm:pb-6">
			<div class="min-w-0 flex-1">
				<Dialog.Title class="mb-3 break-words pr-8 text-2xl font-bold sm:text-3xl">
					{challenge?.name ?? 'Challenge'}
				</Dialog.Title>

				<!-- Tags & Metadata -->
				<div class="flex flex-wrap items-center gap-2">
					{#if challenge?.tags && challenge.tags.length > 0}
						<div role="list" aria-label="Challenge tags" class="contents">
							{#each challenge.tags as tag}
								<span
									class="inline-flex items-center rounded-full bg-black/5 px-2.5 py-0.5 text-xs font-medium dark:bg-white/10"
									role="listitem"
								>
									{tag}
								</span>
							{/each}
						</div>
					{/if}

					{#if challenge?.solves === 0}
						<span class="inline-flex items-center gap-1 text-xs font-medium">
							<Droplet class="h-4 w-4 text-red-500" aria-hidden="true" />
							<span class="opacity-70">0 solves</span>
						</span>
					{:else if challenge?.solves}
						<button
							type="button"
							onclick={onOpenSolves}
							class="cursor-pointer text-xs font-medium opacity-70 transition-opacity hover:underline hover:opacity-100 focus:underline focus:outline-none"
							aria-label="View {challenge.solves} solve{challenge.solves === 1 ? '' : 's'}"
						>
							{challenge.solves}
							{challenge.solves === 1 ? 'solve' : 'solves'}
						</button>
					{/if}
				</div>

				<!-- Authors -->
				{#if challenge?.authors && challenge.authors.length > 0}
					<div class="mt-2 flex items-center gap-1 text-xs font-medium opacity-70">
						<UserEditSolid class="h-4 w-4" aria-hidden="true" />
						<span>
							By {#each challenge.authors as author, i (author)}{author}{i <
								challenge.authors.length - 1
									? ', '
									: ''}{/each}
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
							class="cursor-pointer text-red-500 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950"
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
			<h3 id="description-heading" class="mb-2 text-sm font-semibold opacity-70">Description</h3>
			{#if challenge?.description}
				<Markdown content={challenge.description} class="text-base leading-relaxed" />
			{:else}
				<div class="text-base leading-relaxed opacity-60">No description available.</div>
			{/if}
		</section>

		<!-- Attachments -->
		{#if challenge?.attachments && challenge.attachments.length > 0}
			<section class="mb-6" aria-labelledby="attachments-heading">
				<h3 id="attachments-heading" class="mb-3 text-sm font-semibold opacity-70">Attachments</h3>
				<div class="flex flex-wrap gap-2">
					{#each challenge.attachments as attachment}
						<a
							href={`/attachments/${challenge.id || challenge.chall_id}/${attachment.replace(/^\/+/, '')}`}
							target="_blank"
							rel="noopener noreferrer"
							class="border-input bg-background hover:bg-accent hover:text-accent-foreground focus-visible:ring-ring inline-flex h-9 items-center justify-center gap-2 rounded-md border px-3 text-sm font-medium shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2"
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
				<h3 id="connection-heading" class="mb-3 text-sm font-semibold opacity-70">Connection</h3>
				<div class="flex items-center gap-3">
					<button
						type="button"
						onclick={() => copyToClipboard(connectionString)}
						class="focus:ring-primary inline-flex h-10 items-center gap-2 rounded-md bg-gray-100 px-4 font-mono text-sm font-medium transition-colors hover:bg-gray-200 focus:outline-none focus:ring-2 dark:bg-gray-800 dark:hover:bg-gray-700"
						aria-label="Copy connection string: {connectionString}"
					>
						<span>{connectionString}</span>
					</button>
					{#if connectionString.startsWith('http')}
						<a
							href={connectionString}
							target="_blank"
							rel="noopener noreferrer"
							class="text-primary text-sm font-semibold hover:underline"
						>
							Open
						</a>
					{/if}
				</div>
			</section>
		{/if}

		<!-- Instance Controls -->
		{#if challenge?.instance}
			<InstanceControls bind:challenge {countdown} {onCountdownUpdate} {onInstanceChange} />
		{/if}

		<!-- Submit Flag -->
		<FlagSubmission {challenge} {onSolved} />
	</Dialog.Content>
</Dialog.Root>
