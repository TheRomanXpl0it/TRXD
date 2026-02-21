<script lang="ts">
	import { Container, X } from '@lucide/svelte';
	import { Button } from '@/components/ui/button';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { toast } from 'svelte-sonner';
	import { startInstance, stopInstance } from '$lib/instances';

	let {
		challenge = $bindable(),
		countdown = 0,
		onCountdownUpdate,
		onInstanceChange
	}: {
		challenge: any;
		countdown?: number;
		onCountdownUpdate?: (id: string | number, newCountdown: number) => void;
		onCountdownUpdate?: (id: string | number, newCountdown: number) => void;
		onInstanceChange?: (challenge?: any) => void;
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
			// Update instance-specific fields
			challenge = { ...challenge, instance_host: host, instance_port: port, timeout };
			if (typeof timeout === 'number' && onCountdownUpdate) {
				onCountdownUpdate(challenge.id, Math.max(0, timeout));
			}
			toast.success('Instance created!');
			// Trigger update in parent without refetching
			onInstanceChange?.(challenge);
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
			// Clear instance-specific fields
			challenge = { ...challenge, instance_host: null, instance_port: null, timeout: null };
			if (onCountdownUpdate) {
				onCountdownUpdate(challenge.id, 0);
			}
			toast.success('Instance stopped!');
			// Trigger update in parent
			onInstanceChange?.(challenge);
		} catch (err: any) {
			console.error(err);
			toast.error(`Failed to stop instance: ${err?.message ?? err}`);
		} finally {
			destroyingInstance = false;
		}
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

<div class="mb-6">
	<h3 class="mb-3 text-sm font-semibold opacity-70">Instance</h3>
	<div class="flex w-full flex-row items-center gap-2 sm:gap-3">
		{#if countdown > 0}
			<button
				type="button"
				onclick={() => copyToClipboard(connectionString)}
				class="flex h-11 flex-1 cursor-pointer items-center justify-between gap-2 rounded-md bg-green-600 px-3 text-sm font-semibold text-white transition-colors hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 sm:gap-3 sm:px-4"
				aria-label="Copy instance connection: {connectionString}"
			>
				<div class="flex items-center gap-2">
					<Container class="h-5 w-5 shrink-0" aria-hidden="true" />
					<span class="hidden sm:inline">Running</span>
				</div>
				<span class="truncate font-mono text-xs sm:text-sm">{connectionString}</span>
				<span class="shrink-0 text-xs sm:text-sm" aria-label="Expires in {fmtTimeLeft(countdown)}">
					{fmtTimeLeft(countdown)}
				</span>
			</button>
			{#if connectionString.startsWith('http')}
				<a
					href={connectionString}
					target="_blank"
					rel="noopener noreferrer"
					class="bg-secondary/80 hover:bg-secondary flex h-11 shrink-0 items-center justify-center rounded-md px-4 text-sm font-semibold transition-colors focus:outline-none"
				>
					Open
				</a>
			{/if}
			<Button
				variant="destructive"
				onclick={destroyInstance}
				disabled={destroyingInstance}
				class="h-11 shrink-0 cursor-pointer px-4 font-semibold"
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
				class="h-11 flex-1 cursor-pointer bg-blue-600 font-semibold text-white hover:bg-blue-700"
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
