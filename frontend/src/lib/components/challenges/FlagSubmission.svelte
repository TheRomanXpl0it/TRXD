<script lang="ts">
	import { CheckCircleSolid, FlagSolid, ExclamationCircleSolid } from 'flowbite-svelte-icons';
	import { Button } from '@/components/ui/button';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { toast } from 'svelte-sonner';
	import { submitFlag } from '$lib/challenges';
	import { onMount } from 'svelte';

	let {
		challenge,
		onSolved
	}: {
		challenge: any;
		onSolved?: () => void;
	} = $props();

	let flag = $state('');
	let submittingFlag = $state(false);
	let flagError = $state(false);
	let flagInputElement = $state<HTMLInputElement>();

	async function onSubmitFlag(ev: SubmitEvent) {
		ev.preventDefault();
		if (!challenge?.id) {
			toast.error('No challenge selected');
			return;
		}
		const value = flag.trim();
		if (!value) return;

		submittingFlag = true;
		flagError = false;
		try {
			const res = await submitFlag(challenge.id, value);
			if ((res as any).status === 'Wrong') {
				flagError = true;
				toast.error('Incorrect flag');
				return;
			} else if ((res as any).first_blood) {
				toast.success('First blood!');
			} else {
				toast.success('Correct flag!');
			}
			flag = '';
			challenge.solved = true;
			if (onSolved) onSolved();
		} catch (e: any) {
			flagError = true;
			toast.error(e?.message ?? 'Flag submission failed');
		} finally {
			submittingFlag = false;
		}
	}
</script>

<div class="border-t border-black/10 pt-6 dark:border-white/10">
	<form
		class="flex w-full items-center gap-2 sm:gap-3"
		class:justify-center={challenge?.solved}
		onsubmit={onSubmitFlag}
	>
		{#if !challenge?.solved}
			<div class="relative flex-1">
				{#if flagError}
					<ExclamationCircleSolid
						class="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-red-500"
						aria-hidden="true"
					/>
				{:else}
					<FlagSolid
						class="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 opacity-50"
						aria-hidden="true"
					/>
				{/if}
				<Input
					bind:this={flagInputElement}
					class={`h-11 pl-10 text-sm sm:text-base ${flagError ? 'border-red-500 focus-visible:ring-red-500' : ''}`}
					placeholder="TRX{'{'}...{'}'}"
					bind:value={flag}
					oninput={() => (flagError = false)}
					aria-invalid={flagError}
					aria-label="Enter flag"
					aria-describedby={flagError ? 'flag-error' : undefined}
				/>
				{#if flagError}
					<p id="flag-error" class="sr-only">Incorrect flag entered</p>
				{/if}
			</div>

			<Button
				type="submit"
				color="primary"
				class="h-11 shrink-0 px-6 sm:px-8"
				disabled={submittingFlag || !flag.trim() || flagError}
				aria-label="Submit flag"
			>
				{#if submittingFlag}
					<Spinner class="mr-2" aria-hidden="true" />
					<span class="hidden sm:inline">Submitting...</span>
					<span class="inline sm:hidden">...</span>
					<span class="sr-only">Submitting flag</span>
				{:else}
					Submit
				{/if}
			</Button>
		{:else}
			<div
				class="flex items-center gap-2 rounded-lg border border-green-500 bg-green-50 px-6 py-3 dark:bg-green-950/30"
				role="status"
				aria-live="polite"
			>
				<CheckCircleSolid class="h-5 w-5 text-green-500" aria-hidden="true" />
				<span class="font-semibold text-green-700 dark:text-green-400">Challenge solved!</span>
			</div>
		{/if}
	</form>
</div>
