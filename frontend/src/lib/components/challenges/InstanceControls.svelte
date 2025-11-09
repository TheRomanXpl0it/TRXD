<script lang="ts">
	import { Container, X } from '@lucide/svelte';
	import { Button } from '@/components/ui/button';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { toast } from 'svelte-sonner';
	import { startInstance, stopInstance } from '$lib/instances';

	let {
		challenge,
		countdown = 0,
		onCountdownUpdate
	}: {
		challenge: any;
		countdown?: number;
		onCountdownUpdate?: (id: string | number, newCountdown: number) => void;
	} = $props();

	let creatingInstance = $state(false);
	let destroyingInstance = $state(false);

	function fmtTimeLeft(total: number): string {
		if (!total || total < 0) total = 0;
		const h = Math.floor(total / 3600);
		const m = Math.floor((total % 3600) / 60);
		const s = Math.floor(total % 60);
		if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
		if (m > 0) return `${m}:${String(s).padStart(2, '0')}`;
		return `${s}`;
	}

	function copyToClipboard(text: string) {
		if (typeof navigator === 'undefined') return;
		navigator.clipboard
			.writeText(text)
			.then(() => toast.success('Copied to clipboard!'))
			.catch(() => toast.error('Failed to copy to clipboard.'));
	}

	async function createInstance() {
		if (creatingInstance || !challenge?.id) return;
		creatingInstance = true;
		try {
			const { host, port, timeout } = await startInstance(challenge.id);
			// Force reactivity by creating a new object reference
			challenge = { ...challenge, host, port, timeout };
			if (typeof timeout === 'number' && onCountdownUpdate) {
				onCountdownUpdate(challenge.id, Math.max(0, timeout));
			}
			toast.success('Instance created!');
		} catch (err: any) {
			console.error(err);
			toast.error(`Failed to create instance: ${err?.message ?? err}`);
		} finally {
			creatingInstance = false;
		}
	}

	async function destroyInstance() {
		if (destroyingInstance || !challenge?.id) return;
		destroyingInstance = true;
		try {
			await stopInstance(challenge.id);
			// Force reactivity by creating a new object reference
			challenge = { ...challenge, host: null, port: null, timeout: null };
			if (onCountdownUpdate) {
				onCountdownUpdate(challenge.id, 0);
			}
			toast.success('Instance stopped!');
		} catch (err: any) {
			console.error(err);
			toast.error(`Failed to stop instance: ${err?.message ?? err}`);
		} finally {
			destroyingInstance = false;
		}
	}

	const connectionString = $derived(
		`${challenge?.host ?? ''}${challenge?.port ? `:${challenge.port}` : ''}`
	);
</script>

<div class="mb-6">
	<h3 class="text-sm font-semibold mb-3 opacity-70">Instance</h3>
	<div class="flex w-full flex-row items-center gap-2 sm:gap-3">
		{#if countdown > 0}
			<button
				type="button"
				onclick={() => copyToClipboard(connectionString)}
				class="cursor-pointer flex-1 h-11 bg-green-600 text-white rounded-md px-3 sm:px-4 flex items-center justify-between gap-2 sm:gap-3 font-semibold text-sm transition-colors hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2"
				aria-label="Copy instance connection: {connectionString}"
			>
				<div class="flex items-center gap-2">
					<Container class="h-5 w-5 shrink-0" aria-hidden="true" />
					<span class="hidden sm:inline">Running</span>
				</div>
				<span class="font-mono text-xs sm:text-sm truncate">{connectionString}</span>
				<span class="text-xs sm:text-sm shrink-0" aria-label="Expires in {fmtTimeLeft(countdown)}">
					{fmtTimeLeft(countdown)}
				</span>
			</button>
			<Button
				variant="destructive"
				onclick={destroyInstance}
				disabled={destroyingInstance}
				class="cursor-pointer h-11 px-4 font-semibold shrink-0"
				aria-label="Stop instance"
			>
				{#if destroyingInstance}
					<Spinner class="h-4 w-4" />
					<span class="sr-only">Stopping instance...</span>
				{:else}
					<X class="h-5 w-5" aria-hidden="true" />
					<span class="sr-only">Stop</span>
				{/if}
			</Button>
		{:else}
			<Button
				onclick={createInstance}
				disabled={creatingInstance}
				class="flex-1 h-11 bg-blue-600 hover:bg-blue-700 text-white font-semibold cursor-pointer"
				aria-label="Start challenge instance"
			>
				{#if creatingInstance}
					<Spinner class="mr-2 h-5 w-5" />
					Starting...
				{:else}
					<Container class="mr-2 h-5 w-5" aria-hidden="true" />
					Start Instance
				{/if}
			</Button>
		{/if}
	</div>
</div>
