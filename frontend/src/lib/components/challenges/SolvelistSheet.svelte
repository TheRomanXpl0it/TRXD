<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { toast } from 'svelte-sonner';
	import { getSolves } from '@/challenges';

	import * as Table from '$lib/components/ui/table/index.js';
	import { Droplet, Trophy } from '@lucide/svelte';
	import { push } from 'svelte-spa-router';
	import { userMode } from '$lib/stores/auth';
	import type { Solve } from '$lib/types';

	let { open = $bindable(false), challenge } = $props<{
		open?: boolean;
		challenge: { id: number; name: string };
	}>();

	let loading = $state(false);
	let solves = $state<Solve[]>([]);

	const t = (s: Solve) => new Date(s.timestamp).getTime();

	async function loadSolves(id: number) {
		if (!id) return;
		loading = true;
		try {
			const data = await getSolves(id);
			if (open && challenge?.id === id) {
				// sort ascending by time (earliest first)
				solves = (data ?? []).slice().sort((a, b) => t(a) - t(b));
			}
		} catch (e: any) {
			toast.error(e?.message ?? 'Failed to load solves');
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (open && challenge?.id) {
			loadSolves(challenge.id);
		}
	});

	function goItem(id: number | undefined, ev?: Event) {
		if (!id) return;
		// SPA navigation; also works if user middle-clicks (regular href below)
		if (ev) ev.preventDefault();
		$userMode ? push(`/account/${id}`) : push(`/team/${id}`);
	}

	function truncateName(name: string, maxLength = 32): string {
		if (!name || name.length <= maxLength) return name;
		return name.slice(0, maxLength) + '...';
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="left" class="w-full px-5 sm:max-w-[640px]">
		<div
			class="from-muted/20 to-background mb-6 mt-4 rounded-xl border-0 bg-gradient-to-br p-6 shadow-sm"
		>
			<div class="flex items-center gap-4">
				<div
					class="bg-background flex h-16 w-16 shrink-0 items-center justify-center rounded-full shadow-sm"
				>
					<Trophy class="text-muted-foreground h-8 w-8" />
				</div>
				<div>
					<Sheet.Title class="text-xl font-bold">Challenge Solves</Sheet.Title>
					<Sheet.Description class="text-muted-foreground/80 mt-1">
						Recent solvers for <b class="text-foreground">{challenge?.name}</b>
					</Sheet.Description>
				</div>
			</div>
		</div>

		{#if loading}
			<p class="py-6 text-sm text-gray-500">Loading...</p>
		{:else if solves.length === 0}
			<p class="py-6 text-sm text-gray-500">No solves yet.</p>
		{:else}
			<Table.Root class="w-full">
				<Table.Header>
					<Table.Row>
						<Table.Head class="w-[10%]">#</Table.Head>
						<Table.Head class="w-[50%]">{$userMode ? 'Player' : 'Team'}</Table.Head>
						<Table.Head class="w-[40%] text-right">Date</Table.Head>
					</Table.Row>
				</Table.Header>

				<Table.Body>
					{#each solves as s, i}
						<Table.Row>
							<Table.Cell class="font-medium">
								{#if i === 0}
									<Droplet class="h-4 w-4 text-red-500" />
								{:else}
									{i + 1}
								{/if}
							</Table.Cell>
							<Table.Cell class="py-2">
								<a
									href={$userMode ? '/account/' + s.id : '/team/' + s.id}
									onclick={(e) => goItem(s.id, e)}
									class="cursor-pointer font-medium hover:underline"
								>
									{truncateName(s.name)}
								</a>
							</Table.Cell>
							<Table.Cell
								class="whitespace-nowrap py-2 text-right text-xs text-gray-600 sm:text-sm dark:text-gray-400"
							>
								<span class="hidden sm:inline">{new Date(t(s)).toLocaleString()}</span>
								<span class="sm:hidden"
									>{new Date(t(s)).toLocaleString(undefined, {
										month: 'short',
										day: 'numeric',
										hour: '2-digit',
										minute: '2-digit'
									})}</span
								>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		{/if}
	</Sheet.Content>
</Sheet.Root>
