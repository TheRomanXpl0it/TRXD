<script lang="ts">
	import { getInstances, adminStopInstance } from '$lib/instances';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authState } from '$lib/stores/auth';
	import { Button } from '$lib/components/ui/button';
	import { Trash2, ServerIcon, RefreshCw } from '@lucide/svelte';
	import { showSuccess, showError } from '$lib/utils/toast';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Card from '$lib/components/ui/card/index.js';

	import * as Table from '$lib/components/ui/table';

	let instances = $state<any[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let stopping = $state<Record<string, boolean>>({});

	const isAdmin = $derived(authState.user?.role === 'Admin');

	async function loadInstances() {
		if (!isAdmin) return;
		loading = true;
		error = null;
		try {
			instances = await getInstances();
		} catch (err: any) {
			error = err?.message ?? 'Failed to load instances';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (authState.user && !isAdmin) {
			goto('/404');
		} else if (authState.user && isAdmin) {
			loadInstances();
		}
	});

	async function stop(teamId: number, challId: number) {
		if (!confirm('Are you sure you want to stop this instance?')) return;
		const key = `${teamId}-${challId}`;
		if (stopping[key]) return;
		
		stopping[key] = true;
		try {
			await adminStopInstance(teamId, challId);
			showSuccess('Instance stopped successfully.');
			loadInstances(); // Refresh the list
		} catch (err: any) {
			showError(err, 'Failed to stop instance.');
		} finally {
			stopping[key] = false;
		}
	}

	function formatConn(inst: any) {
		const h = inst.host || '';
		const p = inst.port || '';
		const str = p ? `${h}:${p}` : h;
		if (inst.conn_type === 'TCP') {
			return p ? `nc ${h} ${p}` : `nc ${h}`;
		}
		if (inst.conn_type === 'HTTP' && !str.startsWith('http')) return `http://${str}`;
		if (inst.conn_type === 'HTTPS' && !str.startsWith('http')) return `https://${str}`;
		return str;
	}
</script>

<div class="mx-auto max-w-5xl space-y-8 px-6 py-10">
	<div
		class="from-muted/20 to-background mb-6 mt-6 rounded-xl border-0 bg-gradient-to-br p-6 shadow-sm"
	>
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-4">
				<div
					class="bg-background flex h-16 w-16 shrink-0 items-center justify-center rounded-full shadow-sm"
				>
					<ServerIcon class="text-muted-foreground h-8 w-8" />
				</div>
				<div>
					<h1 class="text-3xl font-bold tracking-tight">Instances</h1>
					<p class="text-muted-foreground mt-2 text-sm">Manage active challenge instances</p>
				</div>
			</div>
			<Button variant="outline" size="sm" onclick={loadInstances} disabled={loading} class="cursor-pointer">
				<RefreshCw class="mr-2 h-4 w-4 {loading ? 'animate-spin' : ''}" />
				Refresh
			</Button>
		</div>
	</div>

	{#if loading && instances.length === 0}
		<div class="flex flex-col items-center justify-center py-12">
			<Spinner class="mb-4 h-8 w-8" />
			<p class="text-gray-600 dark:text-gray-400">Loading instances...</p>
		</div>
	{:else if error && instances.length === 0}
		<div
			class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-600 dark:border-red-800 dark:bg-red-950/20"
		>
			<p class="font-semibold">Error loading instances</p>
			<p class="text-sm">{error}</p>
		</div>
	{:else if instances.length === 0}
		<div class="rounded-lg border bg-card text-card-foreground shadow-sm p-8 text-center text-muted-foreground">
			No active instances found.
		</div>
	{:else}
		<Card.Root class="overflow-hidden border-0 shadow-sm">
			<Card.Content class="p-0">
				<div class="relative mx-4 overflow-auto sm:mx-6">
					<Table.Root>
						<Table.Header class="bg-transparent [&_tr]:border-b-0">
							<Table.Row class="hover:bg-transparent">
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Team (ID)</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Challenge (ID)</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Connection</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Expires</Table.Head>
								<Table.Head class="text-muted-foreground/70 w-[100px] bg-transparent text-right text-[10px] font-bold uppercase tracking-wider">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each instances as inst (inst.docker_id)}
								<Table.Row class="hover:bg-muted/50 transition-colors border-b-0">
									<Table.Cell class="font-medium max-w-[150px] truncate" title={inst.team_name}>
										{inst.team_name} <span class="text-muted-foreground text-xs">({inst.team_id})</span>
									</Table.Cell>
									<Table.Cell class="max-w-[150px] truncate" title={inst.chall_name}>
										{inst.chall_name} <span class="text-muted-foreground text-xs">({inst.chall_id})</span>
									</Table.Cell>
									<Table.Cell>
										<span class="font-mono text-xs">{formatConn(inst)}</span>
									</Table.Cell>
									<Table.Cell>
										{#if inst.expires_at}
											{new Date(inst.expires_at).toLocaleString()}
										{:else}
											<span class="text-muted-foreground">Never</span>
										{/if}
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button
											variant="outline"
											size="icon"
											class="text-red-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950/50"
											onclick={() => stop(inst.team_id, inst.chall_id)}
											disabled={stopping[`${inst.team_id}-${inst.chall_id}`]}
											title="Stop Instance"
										>
											{#if stopping[`${inst.team_id}-${inst.chall_id}`]}
												<Spinner class="h-4 w-4" />
											{:else}
												<Trash2 class="h-4 w-4" />
											{/if}
										</Button>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</div>
			</Card.Content>
		</Card.Root>
	{/if}
</div>
