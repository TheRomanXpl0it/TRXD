<script lang="ts">
	import { getSubmissions, deleteSubmission } from '$lib/challenges';
	import { authState } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Trash2, MessageSquareShare, RefreshCw, CheckCircle, XCircle, Droplet, Eye, EyeOff } from '@lucide/svelte';
	import { showSuccess, showError } from '$lib/utils/toast';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Pagination from '$lib/components/ui/pagination';
	import * as Card from '$lib/components/ui/card/index.js';

	import * as Table from '$lib/components/ui/table';

	let submissions = $state<any[]>([]);
	let paginationInfo = $state<any>({ page: 1, pages: 1, total: 0 });
	let currentPage = $state(1);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let deleting = $state<Record<string, boolean>>({});
	let expandedFlags = $state<Record<string, boolean>>({});

	const limit = 50;
	
	const count = $derived(paginationInfo.total);
	const perPage = $derived(limit);

	const isAdmin = $derived(authState.user?.role === 'Admin');

	async function loadSubmissions(page = 1) {
		if (!isAdmin) return;
		loading = true;
		error = null;
		try {
			const res = await getSubmissions(page, limit);
			submissions = Array.isArray(res?.submissions) ? res.submissions : [];
			paginationInfo = {
				page: page,
				pages: Math.ceil((res?.total || 0) / limit),
				total: res?.total || 0,
			};
			currentPage = page;
		} catch (err: any) {
			error = err?.message ?? 'Failed to load submissions';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (authState.user && !isAdmin) {
			goto('/404');
		} else if (authState.user && isAdmin) {
			loadSubmissions(currentPage);
		}
	});

	async function remove(id: number | string) {
		if (!confirm('Are you sure you want to delete this submission?')) return;
		if (deleting[id]) return;
		
		deleting[id] = true;
		try {
			await deleteSubmission(id);
			showSuccess('Submission deleted successfully.');
			loadSubmissions(currentPage); // Refresh current page
		} catch (err: any) {
			showError(err, 'Failed to delete submission.');
		} finally {
			deleting[id] = false;
		}
	}
	
	$effect(() => {
		if (currentPage !== paginationInfo.page && !loading) {
			loadSubmissions(currentPage);
		}
	});
</script>

<div class="mx-auto max-w-7xl space-y-8 px-6 py-10">
	<div
		class="from-muted/20 to-background mb-6 mt-6 rounded-xl border-0 bg-gradient-to-br p-6 shadow-sm"
	>
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-4">
				<div
					class="bg-background flex h-16 w-16 shrink-0 items-center justify-center rounded-full shadow-sm"
				>
					<MessageSquareShare class="text-muted-foreground h-8 w-8" />
				</div>
				<div>
					<h1 class="text-3xl font-bold tracking-tight">Submissions</h1>
					<p class="text-muted-foreground mt-2 text-sm">View and manage challenge flag submissions</p>
				</div>
			</div>
			<Button variant="outline" size="sm" onclick={() => loadSubmissions(currentPage)} disabled={loading} class="cursor-pointer">
				<RefreshCw class="mr-2 h-4 w-4 {loading ? 'animate-spin' : ''}" />
				Refresh
			</Button>
		</div>
	</div>

	{#if loading && submissions.length === 0}
		<div class="flex flex-col items-center justify-center py-12">
			<Spinner class="mb-4 h-8 w-8" />
			<p class="text-gray-600 dark:text-gray-400">Loading submissions...</p>
		</div>
	{:else if error && submissions.length === 0}
		<div
			class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-600 dark:border-red-800 dark:bg-red-950/20"
		>
			<p class="font-semibold">Error loading submissions</p>
			<p class="text-sm">{error}</p>
		</div>
	{:else if submissions.length === 0}
		<div class="rounded-lg border bg-card text-card-foreground shadow-sm p-8 text-center text-muted-foreground">
			No submissions found.
		</div>
	{:else}
		<Card.Root class="overflow-hidden border-0 shadow-sm">
			<Card.Content class="p-0">
				<div class="relative mx-4 overflow-auto sm:mx-6">
					<Table.Root>
						<Table.Header class="bg-transparent [&_tr]:border-b-0">
							<Table.Row class="hover:bg-transparent">
								<Table.Head class="text-muted-foreground/70 w-[60px] bg-transparent text-[10px] font-bold uppercase tracking-wider">ID</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">User</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Team</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Challenge</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Flag</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Status</Table.Head>
								<Table.Head class="text-muted-foreground/70 bg-transparent text-[10px] font-bold uppercase tracking-wider">Timestamp</Table.Head>
								<Table.Head class="text-muted-foreground/70 w-[80px] bg-transparent text-right text-[10px] font-bold uppercase tracking-wider">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each submissions as sub (sub.id)}
								<Table.Row class="hover:bg-muted/50 cursor-pointer border-b-0 transition-colors">
									<Table.Cell class="font-mono text-xs">{sub.id}</Table.Cell>
									<Table.Cell class="max-w-[120px]">
										<a href={`/account/${sub.user_id}`} class="hover:underline font-medium text-primary block truncate" title={sub.user_name}>
											{sub.user_name}
										</a>
									</Table.Cell>
									<Table.Cell class="max-w-[120px]">
										<a href={`/team/${sub.team_id}`} class="hover:underline text-muted-foreground block truncate" title={sub.team_name}>
											{sub.team_name}
										</a>
									</Table.Cell>
									<Table.Cell class="font-medium max-w-[120px] truncate" title={sub.chall_name}>
										{sub.chall_name}
									</Table.Cell>
									<Table.Cell class="max-w-[200px] sm:max-w-[300px]">
										<div class="flex items-center gap-2">
											<code class="px-1.5 py-0.5 rounded bg-muted/50 text-xs text-muted-foreground w-full {expandedFlags[sub.id] ? 'break-all' : 'truncate block'}" title={sub.flag}>
												{sub.flag}
											</code>
											<Button variant="ghost" size="icon" class="h-6 w-6 shrink-0" onclick={(e) => { e.stopPropagation(); expandedFlags[sub.id] = !expandedFlags[sub.id]; }}>
												{#if expandedFlags[sub.id]}
													<EyeOff class="h-3 w-3" />
												{:else}
													<Eye class="h-3 w-3" />
												{/if}
											</Button>
										</div>
									</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-1.5 whitespace-nowrap">
											{#if sub.status === 'Correct'}
												<span class="inline-flex items-center gap-1 text-green-600 dark:text-green-500 font-medium text-sm">
													<CheckCircle class="w-4 h-4" /> Correct
												</span>
											{:else}
												<span class="inline-flex items-center gap-1 text-red-600 dark:text-red-500 font-medium text-sm">
													<XCircle class="w-4 h-4" /> Incorrect
												</span>
											{/if}
											{#if sub.first_blood}
												<span title="First Blood!" class="inline-flex items-center shrink-0">
													<Droplet class="w-4 h-4 text-red-500 fill-red-500 drop-shadow-sm" />
												</span>
											{/if}
										</div>
									</Table.Cell>
									<Table.Cell class="whitespace-nowrap text-muted-foreground text-sm">
										{new Date(sub.timestamp).toLocaleString()}
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button
											variant="outline"
											size="icon"
											class="text-red-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950/50"
											onclick={() => remove(sub.id)}
											disabled={deleting[sub.id]}
											title="Delete Submission"
										>
											{#if deleting[sub.id]}
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
		
		{#if paginationInfo.pages > 1}
			<div class="mt-6 flex justify-center">
				<Pagination.Root {count} {perPage} bind:page={currentPage} siblingCount={1} class="mt-4">
					{#snippet children({ pages, currentPage })}
						<Pagination.Content class="gap-4">
							<Pagination.Item class="mx-2">
								<Pagination.PrevButton class="h-9 w-9" />
							</Pagination.Item>
							{#each pages as page (page.key)}
								{#if page.type === 'ellipsis'}
									<Pagination.Item>
										<Pagination.Ellipsis />
									</Pagination.Item>
								{:else}
									<Pagination.Item>
										<Pagination.Link
											{page}
											isActive={currentPage === page.value}
											class="h-9 w-9"
										>
											{page.value}
										</Pagination.Link>
									</Pagination.Item>
								{/if}
							{/each}
							<Pagination.Item class="mx-2">
								<Pagination.NextButton class="h-9 w-9" />
							</Pagination.Item>
						</Pagination.Content>
					{/snippet}
				</Pagination.Root>
			</div>
		{/if}
	{/if}
</div>
