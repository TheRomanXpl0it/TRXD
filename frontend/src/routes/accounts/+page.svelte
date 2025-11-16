<script lang="ts">
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import { getUsers } from '@/user';
	import { link, push } from 'svelte-spa-router';
	import ErrorMessage from '$lib/components/ui/error-message.svelte';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { Avatar } from 'flowbite-svelte';
	import { BugOutline } from 'flowbite-svelte-icons';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import countries from '$lib/data/countries.json';

	let perPage = $state(20);

	const usersQuery = createQuery(() => ({
		queryKey: ['users'],
		queryFn: getUsers,
		staleTime: 5 * 60 * 1000
	}));

	const usersData = $derived(usersQuery.data ?? []);
	const loading = $derived(usersQuery.isLoading);
	const error = $derived(usersQuery.error?.message ?? null);

	// Helper to convert ISO3 to ISO2 for flag icons
	function getCountryIso2(iso3: string): string | null {
		const country = (countries as any[]).find(c => c.iso3?.toUpperCase() === iso3?.toUpperCase());
		return country?.iso2?.toUpperCase() ?? null;
	}

	// Sort alphabetically by name
	const sorted = $derived(
		(Array.isArray(usersData) ? [...usersData] : []).sort(
			(a, b) => (a?.name || '').localeCompare(b?.name || '')
		)
	);
	const count = $derived(sorted.length);

	// Track page changes
	let currentPage = $state(1);
	$effect(() => {
		if (currentPage > 1) {
			setTimeout(() => {
				const paginationEl = document.getElementById('pagination-controls');
				if (paginationEl) {
					paginationEl.scrollIntoView({ behavior: 'instant', block: 'nearest' });
				}
			}, 0);
		}
	});
</script>

<p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Accounts</p>
<hr class="my-2 mb-10 h-px border-0 bg-gray-200 dark:bg-gray-700" />

{#if loading}
	<div class="flex flex-col items-center justify-center py-12">
		<Spinner class="mb-4 h-8 w-8" />
		<p class="text-gray-600 dark:text-gray-400">Loading accounts...</p>
	</div>
{:else if error}
	<ErrorMessage title="Error loading accounts" message={error} />
{:else}
	{@const startIndex = (currentPage - 1) * perPage}
	{@const pageRows = sorted.slice(startIndex, startIndex + perPage)}
	{@const totalPages = Math.max(1, Math.ceil(count / perPage))}
	{@const singlePage = totalPages <= 1}

	<!-- Table -->
	<div class="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head class="w-[4rem]">Photo</Table.Head>
					<Table.Head class="min-w-[10rem]">Account Name</Table.Head>
					<Table.Head class="w-[8rem]">Country</Table.Head>
				</Table.Row>
			</Table.Header>

			<Table.Body>
				{#if pageRows.length === 0}
					<Table.Row>
						<Table.Cell colspan={3} class="py-8 text-center text-gray-500">
							No accounts yet.
						</Table.Cell>
					</Table.Row>
				{:else}
					{#each pageRows as account (account.id)}
						<Table.Row class="h-16">
							<Table.Cell class="py-3">
								{#if account.image}
									<div class="h-12 w-12 overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700">
										<img
											src={account.image}
											alt={account.name}
											class="h-full w-full object-cover"
										/>
									</div>
								{:else}
									<div class="flex h-12 w-12 items-center justify-center rounded-full bg-gray-200 dark:bg-gray-700">
										<BugOutline class="h-6 w-6" />
									</div>
								{/if}
							</Table.Cell>

							<Table.Cell class="py-3">
								<a
									href={`#/account/${account.id}`}
									use:link
									onclick={(e) => {
										e.preventDefault();
										push(`/account/${account.id}`);
									}}
									class="text-primary cursor-pointer text-sm font-medium underline-offset-2 hover:underline sm:text-base"
									title={`View account ${account.name}`}
								>
									{account.name}
								</a>
							</Table.Cell>

							<Table.Cell class="py-3">
								{#if account.country}
									{@const iso2 = getCountryIso2(account.country)}
									<div class="flex items-center gap-2">
										{#if iso2}
											<Icon icon={`circle-flags:${iso2.toLowerCase()}`} width="20" height="20" />
										{/if}
										<span class="text-sm">{account.country}</span>
									</div>
								{:else}
									<span class="text-xs text-gray-500">-</span>
								{/if}
							</Table.Cell>
						</Table.Row>
					{/each}
				{/if}
			</Table.Body>
		</Table.Root>
	</div>

	<!-- Pagination -->
	<Pagination.Root {count} {perPage} bind:page={currentPage} siblingCount={1} class="mt-4">
		{#snippet children({ pages, currentPage: pageNum })}
			<div class="flex w-full justify-center overflow-x-auto" id="pagination-controls">
				<Pagination.Content
					class={`flex items-center justify-center gap-1 ${singlePage ? 'pointer-events-none opacity-50' : ''}`}
					aria-disabled={singlePage}
				>
					<Pagination.Item>
						<Pagination.PrevButton />
					</Pagination.Item>

					{#each pages as page (page.key)}
						{#if page.type === 'ellipsis'}
							<Pagination.Item>
								<Pagination.Ellipsis />
							</Pagination.Item>
						{:else}
							<Pagination.Item>
								<Pagination.Link {page} isActive={pageNum === page.value}>
									{page.value}
								</Pagination.Link>
							</Pagination.Item>
						{/if}
					{/each}

					<Pagination.Item>
						<Pagination.NextButton />
					</Pagination.Item>
				</Pagination.Content>
			</div>
		{/snippet}
	</Pagination.Root>
{/if}
