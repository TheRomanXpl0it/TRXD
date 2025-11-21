<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { toast } from 'svelte-sonner';
	import { getSolves } from '@/challenges';

	import * as Table from '$lib/components/ui/table/index.js';
	import { Droplet } from '@lucide/svelte';
	import { push } from 'svelte-spa-router';
	import { userMode } from '$lib/stores/auth';

	let { open = $bindable(false), challenge } = $props<{ open?: boolean; challenge: any }>();

	let loading = $state(false);
	let solves = $state<any[]>([]);

	// unified time getter
	const t = (s: any) => new Date(s?.createdAt ?? s?.time ?? s?.timestamp ?? Date.now()).getTime();

	function runSolvesEffect(isOpen: boolean, id?: string) {
		if (!isOpen || !id) return;
		const ac = new AbortController();
		loading = true;

		(async () => {
			try {
				const data = await getSolves(id);
				if (!ac.signal.aborted && open && challenge?.id === id) {
					// sort ascending by time (earliest first)
					solves = (data ?? []).slice().sort((a, b) => t(a) - t(b));
				}
			} catch (e: any) {
				if (e?.name !== 'AbortError')
					toast.error(e instanceof Error ? e.message : 'Failed to load solves');
			} finally {
				if (!ac.signal.aborted) loading = false;
			}
		})();

		return () => ac.abort();
	}

	$effect(() => runSolvesEffect(open, challenge?.id));

	function goItem(id: string | number | undefined, ev?: Event) {
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
		<Sheet.Header>
			<Sheet.Title>Solves</Sheet.Title>
			<Sheet.Description>Recent solvers for {challenge?.name}</Sheet.Description>
		</Sheet.Header>

		{#if loading}
			<p class="py-6 text-sm text-gray-500">Loading...</p>
		{:else if solves.length === 0}
			<p class="py-6 text-sm text-gray-500">No solves yet.</p>
		{:else}
			<Table.Root class="w-full">
				<Table.Header>
					<Table.Row>
						<Table.Head class="w-10"></Table.Head>
						<Table.Head>{$userMode ? 'Player' : 'Team'}</Table.Head>
						<Table.Head class="text-right whitespace-nowrap">Time</Table.Head>
					</Table.Row>
				</Table.Header>

				<Table.Body>
					{#each solves as s, i}
						<Table.Row>
							<Table.Cell class="py-2">
								{#if i === 0}
									<Droplet class="h-4 w-4 text-red-500" />
								{/if}
							</Table.Cell>

							<Table.Cell class="py-2">
								{#if s.id}
									{#key s.id}
										<a
											href={$userMode ? '/account/' + s.id : '/team/' + s.id}
											onclick={(e) => goItem(s.id, e)}
											class="cursor-pointer font-medium hover:underline"
										>
											{truncateName(s?.name ?? 'Anonymous')}
										</a>
									{/key}
								{:else}
									<span class="font-medium">{truncateName(s?.name ?? 'Anonymous')}</span>
								{/if}
							</Table.Cell>

							<Table.Cell class="py-2 text-right text-gray-600 dark:text-gray-400 text-xs sm:text-sm whitespace-nowrap">
								<span class="hidden sm:inline">{new Date(t(s)).toLocaleString()}</span>
								<span class="sm:hidden">{new Date(t(s)).toLocaleString(undefined, { 
									month: 'short', 
									day: 'numeric', 
									hour: '2-digit', 
									minute: '2-digit' 
								})}</span>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		{/if}
	</Sheet.Content>
</Sheet.Root>
