<script lang="ts">
	import { toast } from 'svelte-sonner';
	import SolveListSheet from '$lib/components/challenges/SolvelistSheet.svelte';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { push } from 'svelte-spa-router';
	import { getChallenges, deleteChallenge } from '$lib/challenges';
	import { getCategories } from '$lib/categories';
	import { user as authUser } from '$lib/stores/auth';
	import { onMount, untrack } from 'svelte';
	import { createQuery } from '@tanstack/svelte-query';
	
	import ChallengeFilters from '$lib/components/challenges/ChallengeFilters.svelte';
	import ChallengeCard from '$lib/components/challenges/ChallengeCard.svelte';
	import ChallengeModal from '$lib/components/challenges/ChallengeModal.svelte';
	import AdminControls from '$lib/components/challenges/AdminControls.svelte';

	import { config } from '$lib/env';

	// Lazy-loaded component handles
	type Cmp = typeof import('svelte').SvelteComponent;
	let CreateModalCmp: Cmp | null = $state(null);
	let DeleteDialogCmp: Cmp | null = $state(null);
	let ChallengeEditSheetCmp: Cmp | null = $state(null);

	// ** lazy load elements on first use **
	async function openCreate() {
		if (!CreateModalCmp) {
			const mod = await import('$lib/components/challenges/CreateChallengeModal.svelte');
			CreateModalCmp = mod.default;
		}
		createChallengeOpen = true;
	}

	async function requestDelete(ch: any) {
		toDelete = ch;
		if (!DeleteDialogCmp) {
			const mod = await import('$lib/components/challenges/DeleteChallengeDialog.svelte');
			DeleteDialogCmp = mod.default;
		}
		confirmDeleteOpen = true;
	}

	async function modifyChallenge(ch: any) {
     	if (!ChallengeEditSheetCmp) {
      		const mod = await import('$lib/components/challenges/ChallengeEditSheet.svelte');
      		ChallengeEditSheetCmp = mod.default;
     	}
     	editOpen = true;
	}

	// Local state
	let createChallengeOpen = $state(false);

	// Track selected challenge ID instead of the full object to avoid circular dependencies
	let selectedId = $state<number | null>(null);
	let countdowns: Record<string, number> = $state({});

	const challengesQuery = createQuery(() => ({
		queryKey: ['challenges'],
		queryFn: getChallenges,
		staleTime: 0 // Always fetch fresh data
	}));

	// Don't need to refetch categories
	const categoriesQuery = createQuery(() => ({
		queryKey: ['categories'],
		queryFn: getCategories,
		staleTime: 10 * 60 * 1000 // 10 minutes
	}));

	const challenges = $derived(challengesQuery.data ?? []);
	const loading = $derived(challengesQuery.isLoading || categoriesQuery.isLoading);
	const error = $derived(challengesQuery.error?.message ?? categoriesQuery.error?.message ?? null);

	// Derive the actual selected challenge from the ID - always fresh from challenges array
	const selected = $derived(
		selectedId ? challenges.find((c: any) => c.id === selectedId) ?? null : null
	);

	const categories = $derived(
		(categoriesQuery.data ?? [])
			.map((c: any) => ({
				value: c.name,
				label: c.name
			}))
			.sort((a: any, b: any) => a.label.localeCompare(b.label))
	);
	
	//comparators to sort challenges after filtering
	const comparators: Record<string, (a: Challenge, b: Challenge) => number> = {
        "points-min-to-max": (a, b) => (a.points ?? 0) - (b.points ?? 0),
        "points-max-to-min": (a, b) => (b.points ?? 0) - (a.points ?? 0),
        "solves-min-to-max": (a, b) => (a.solves ?? 0) - (b.solves ?? 0),
        "solves-max-to-min": (a, b) => (b.solves ?? 0) - (a.solves ?? 0),
        "alphabetical-a-to-z": (a, b) => (a.name ?? a.title ?? '').localeCompare(b.name ?? b.title ?? ''),
        "alphabetical-z-to-a": (a, b) => (b.name ?? b.title ?? '').localeCompare(a.name ?? a.title ?? '')
    };
	
	const sortedChallenges = $derived(
        [...filteredChallenges].sort(comparators[sortMethod] ?? comparators["alphabetical-a-to-z"])
    );

	let points: number = $state(500);
	let category: any = $state(null);
	let challengeType = $state('Container');
	let challengeName = $state('');
	let challengeDescription = $state('');
	let dynamicScore = $state(false);
	let createLoading = $state(false);

	let openModal = $state(false);
	let solvesOpen = $state(false);
	let editOpen = $state(false);
	let creatingInstance = $state<Record<number, boolean>>({});
	let destroyingInstance = $state<Record<number, boolean>>({});

	// Filters
	let filterCategories = $state<string[]>([]);
	let filterTags = $state<string[]>([]);
	let sortMethod = $state<string>('alphabetical-a-to-z');
	let search = $state('');
	let tagsOpen = $state(false);
	let categoriesOpen = $state(false);

	// Load compact view preference from localStorage
	let compactView = $state(false);

	onMount(() => {
		const saved = localStorage.getItem('challenges-compact-view');
		if (saved !== null) {
			compactView = saved === 'true';
		}
	});

	// Save to localStorage when compactView changes
	$effect(() => {
		if (typeof localStorage !== 'undefined') {
			localStorage.setItem('challenges-compact-view', String(compactView));
		}
	});

	// Fuzzy search helpers
	function norm(s: any) {
		return String(s ?? '')
			.trim()
			.toLowerCase();
	}
	
	function fuzzyScore(text: string, query: string) {
		const t = norm(text),
			q = norm(query);
		if (!q) return 1e9;
		if (t === q) return 1e6;
		if (t.startsWith(q)) return 5e5;
		if (t.includes(q)) return 3e5;
		let ti = 0,
			qi = 0,
			penalty = 0;
		while (ti < t.length && qi < q.length) {
			if (t[ti] === q[qi]) qi++;
			else penalty++;
			ti++;
		}
		return qi === q.length ? 1e5 - penalty : -Infinity;
	}

	const allTags = $derived(
		Array.from(
			new Set<string>(
				(challenges ?? []).flatMap((ch: any) => (ch?.tags ?? []).map((t: any) => String(t)))
			)
		).sort((a, b) => a.localeCompare(b))
	);

	// Optimize filtering with early returns and memoization
	const filteredChallenges = $derived.by(() => {
		const hasSearch = search.trim().length > 0;
		const hasCategoryFilter = filterCategories && filterCategories.length > 0;
		const hasTagsFilter = filterTags && filterTags.length > 0;
		
		// No filters at all - return all challenges
		if (!hasSearch && !hasCategoryFilter && !hasTagsFilter) {
			return challenges ?? [];
		}
		
		const searchQuery = hasSearch ? norm(search) : '';
		
		return (challenges ?? []).filter((c: any) => {
			// Category filter
			if (hasCategoryFilter) {
				const cat = norm(c?.category?.name ?? c?.category ?? '');
				if (!filterCategories.some((fc: string) => norm(fc) === cat)) {
					return false;
				}
			}
			
			// Tags filter
			if (hasTagsFilter) {
				const tags = (c?.tags ?? []).map((t: any) => String(t));
				if (!filterTags.every((t: string) => tags.includes(t))) {
					return false;
				}
			}
			
			// Search filter
			if (hasSearch) {
				const cat = c?.category?.name ?? c?.category ?? '';
				const tags = (c?.tags ?? []).map((t: any) => String(t));
				const name = c?.name ?? c?.title ?? '';
				
				// Quick check: exact name match
				if (norm(name).includes(searchQuery)) return true;
				if (norm(cat).includes(searchQuery)) return true;
				if (tags.some((t: string) => norm(t).includes(searchQuery))) return true;
				
				// Fallback to fuzzy if simple includes didn't work
				return (
					fuzzyScore(name, searchQuery) > -Infinity ||
					fuzzyScore(cat, searchQuery) > -Infinity ||
					tags.some((t: string) => fuzzyScore(t, searchQuery) > -Infinity)
				);
			}
			
			return true;
		});
	});

	const activeFiltersCount = $derived(
		(filterCategories?.length ?? 0) + (filterTags?.length ?? 0)
	);

	// delete confirmation modal state
	let confirmDeleteOpen = $state(false);
	let deleting = $state(false);
	let toDelete: any = $state(null);

	const challengeTypes = [
		{ value: 'Container', label: 'Container' },
		{ value: 'Compose', label: 'Compose' },
		{ value: 'Normal', label: 'Normal' }
	];

	// Update countdowns when challenges data changes
	$effect(() => {
		const next: Record<string, number> = {};
		for (const c of challenges) {
			if (typeof c?.timeout === 'number' && c.timeout > 0) next[c.id] = c.timeout;
		}
		countdowns = next;
	});

	// Handle errors
	$effect(() => {
		if (challengesQuery.error) {
			toast.error('You need to join a team first!');
			push('/team');
		}
	});

	onMount(() => {
		const timer = setInterval(() => {
			for (const id in countdowns) if (countdowns[id] > 0) countdowns[id] = countdowns[id] - 1;
		}, 1000);
		return () => clearInterval(timer);
	});


	function groupByCategory(list: Challenge[], cmp: (a: Challenge, b: Challenge) => number) {
        const map: Record<string, Challenge[]> = {};
        for (const c of list) {
            const label = (typeof c?.category === 'string' ? c.category : c?.category?.name) ?? 'Uncategorized';
            (map[label] ??= []).push(c);
        }
        return Object.entries(map)
            .sort(([a], [b]) => a.localeCompare(b))
            .map(([cat, items]) => [cat, items.sort(cmp)]) as [string, Challenge[]][];
    }
    
    const grouped = $derived.by(() => groupByCategory(sortedChallenges, comparators[sortMethod]));

	function openChallenge(ch: any) {
		selectedId = ch?.id ?? null;
		openModal = true;
	}
	
	function closeModal() {
		openModal = false;
		// Clear selection when closing
		setTimeout(() => {
			selectedId = null;
		}, 200);
	}

	function copyToClipboard(text: string) {
		if (typeof navigator === 'undefined') return;
		navigator.clipboard
			.writeText(text)
			.then(() => toast.success('Copied to clipboard!'))
			.catch(() => toast.error('Failed to copy to clipboard.'));
	}

	function updateCountdown(id: string | number, newCountdown: number) {
		countdowns[id] = newCountdown;
	}

	function handleChallengeSolved() {
		// Refetch challenges to get updated data
		challengesQuery.refetch();
	}

	async function confirmDelete() {
		if (!toDelete?.id) return;
		deleting = true;
		try {
			await deleteChallenge(toDelete.id);
			toast.success('Challenge deleted.');
			confirmDeleteOpen = false;
			openModal = false;
			toDelete = null;
			challengesQuery.refetch();
		} catch (err: any) {
			toast.error(err?.message ?? 'Failed to delete challenge.');
		} finally {
			deleting = false;
		}
	}
</script>

<p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Challenges</p>
<hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
<p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
	"A man who loves to walk will walk more than a man who loves his destination"
</p>

{#if $authUser?.role === 'Admin'}
	<AdminControls
		onopen-create={openCreate}
		oncategory-created={() => {
			challengesQuery.refetch();
			categoriesQuery.refetch();
		}}
	/>
{/if}

<ChallengeFilters
	bind:search
	bind:filterCategories
	bind:filterTags
	bind:sortMethod
	bind:compactView
	{categories}
	{allTags}
	{activeFiltersCount}
/>

{#if loading}
	<div class="flex flex-col items-center justify-center py-12">
		<Spinner class="mb-4 h-8 w-8" />
		<p class="text-gray-600 dark:text-gray-400">Loading challenges...</p>
	</div>
{:else if error}
	<div class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-600 dark:border-red-800 dark:bg-red-950/20">
		<p class="font-semibold">Error loading challenges</p>
		<p class="text-sm">{error}</p>
	</div>
{:else}
	{#each grouped as [category, items]}
		<section class={compactView ? 'mb-4' : 'mb-10'} aria-labelledby="category-{category.replace(/\s+/g, '-')}">
			<div class="{compactView ? 'mb-2' : 'mb-3'} flex items-center gap-3">
				<h2 id="category-{category.replace(/\s+/g, '-')}" class="text-2xl font-bold leading-tight text-gray-900 dark:text-white">
					{category}
				</h2>
				<span class="text-sm text-gray-500 dark:text-gray-400">
					{items.length} challenge{items.length === 1 ? '' : 's'}
				</span>
			</div>

			{#if compactView}
				<div class="space-y-2 px-0.5 py-0.5" role="list" aria-label="{category} challenges">
					{#each items as ch (ch.id)}
						<ChallengeCard
							challenge={ch}
							{compactView}
							countdown={countdowns[ch.id] ?? 0}
							onclick={() => openChallenge(ch)}
						/>
					{/each}
				</div>
			{:else}
				<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 px-0.5" role="list" aria-label="{category} challenges">
					{#each items as ch (ch.id)}
						<ChallengeCard
							challenge={ch}
							{compactView}
							countdown={countdowns[ch.id] ?? 0}
							onclick={() => openChallenge(ch)}
						/>
					{/each}
				</div>
			{/if}
		</section>
	{/each}
{/if}

<ChallengeModal
	bind:open={openModal}
	challenge={selected}
	countdown={selected?.id ? countdowns[selected.id] ?? 0 : 0}
	isAdmin={$authUser?.role === 'Admin'}
	onEdit={modifyChallenge}
	onDelete={requestDelete}
	onSolved={handleChallengeSolved}
	onCountdownUpdate={updateCountdown}
	onOpenSolves={() => (solvesOpen = true)}
/>

<SolveListSheet bind:open={solvesOpen} challenge={selected} />


{#if DeleteDialogCmp}
  <DeleteDialogCmp
    bind:open={confirmDeleteOpen}
    {toDelete}
    {deleting}
    on:confirm={confirmDelete}
  />
{/if}

{#if CreateModalCmp}
  <CreateModalCmp
    bind:open={createChallengeOpen}
    bind:challengeName
    bind:challengeDescription
    bind:category
    bind:challengeType
    bind:points
    bind:dynamicScore
    categories={categories}
    challengeTypes={challengeTypes}
    oncreated={() => challengesQuery.refetch()}
  />
{/if}

{#if ChallengeEditSheetCmp}
  <ChallengeEditSheetCmp
    bind:open={editOpen}
    challenge_user={selected}
    onupdated={() => challengesQuery.refetch()}
    all_tags={allTags}
  />
{/if}