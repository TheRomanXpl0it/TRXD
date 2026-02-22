<script lang="ts">
	import { toast } from 'svelte-sonner';
	import SolveListSheet from '$lib/components/challenges/SolvelistSheet.svelte';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { goto } from '$app/navigation';
	import { getChallenges, deleteChallenge } from '$lib/challenges';
	import { getCategories } from '$lib/categories';
	import { authState } from '$lib/stores/auth';
	import { onMount, untrack } from 'svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	import ChallengeFilters from '$lib/components/challenges/ChallengeFilters.svelte';
	import ChallengeCard from '$lib/components/challenges/ChallengeCard.svelte';
	import ChallengeModal from '$lib/components/challenges/ChallengeModal.svelte';
	import AdminControls from '$lib/components/challenges/AdminControls.svelte';
	import { Flag } from '@lucide/svelte';
	import type { Challenge } from '$lib/types';

	import { config } from '$lib/env';

	// Lazy-loaded admin component handles
	let CreateModalCmp: any | null = $state(null);
	let DeleteDialogCmp: any | null = $state(null);
	let ChallengeEditSheetCmp: any | null = $state(null);

	const isAdmin = $derived(authState.user?.role === 'Admin');

	// Preload admin components when page loads if user is admin
	onMount(() => {
		if (isAdmin) {
			import('$lib/components/challenges/CreateChallengeModal.svelte').then((mod) => {
				CreateModalCmp = mod.default;
			});
			import('$lib/components/challenges/DeleteChallengeDialog.svelte').then((mod) => {
				DeleteDialogCmp = mod.default;
			});
			import('$lib/components/challenges/ChallengeEditSheet.svelte').then((mod) => {
				ChallengeEditSheetCmp = mod.default;
			});
		}
	});

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
	const queryClient = useQueryClient();

	const challengesQuery = createQuery(() => ({
		queryKey: ['challenges'],
		queryFn: getChallenges,
		staleTime: 5 * 60 * 1000 // 5 minutes
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
		selectedId ? (challenges.find((c: any) => c.id === selectedId) ?? null) : null
	);

	const categories = $derived(
		(categoriesQuery.data ?? [])
			.map((c: any) => ({
				value: typeof c === 'string' ? c : c.name,
				label: typeof c === 'string' ? c : c.name
			}))
			.sort((a: any, b: any) => (a.label || '').localeCompare(b.label || ''))
	);

	// Filters
	let filterCategories = $state<string[]>([]);
	let filterTags = $state<string[]>([]);
	let search = $state('');
	let debouncedSearch = $state('');
	let tagsOpen = $state(false);
	let categoriesOpen = $state(false);

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

	// Optimize filtering with early returns and memoization
	const filteredChallenges = $derived.by(() => {
		const hasSearch = debouncedSearch.trim().length > 0;
		const hasCategoryFilter = filterCategories && filterCategories.length > 0;
		const hasTagsFilter = filterTags && filterTags.length > 0;

		// No filters at all - return all challenges
		if (!hasSearch && !hasCategoryFilter && !hasTagsFilter) {
			return challenges ?? [];
		}

		const searchQuery = hasSearch ? norm(debouncedSearch) : '';

		return (challenges ?? []).filter((c: any) => {
			if (hasCategoryFilter) {
				const cat = norm(c.category);
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

			if (hasSearch) {
				const cat = c.category;
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

	// Sort challenges by points (lowest first), stable sort for equal points
	const sortedChallenges = $derived(
		filteredChallenges
			.map((c: any, i: number) => ({ ...c, _index: i }))
			.sort((a: any, b: any) => {
				const pointsDiff = (a.points ?? 0) - (b.points ?? 0);
				return pointsDiff !== 0 ? pointsDiff : a._index - b._index;
			})
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

	const allTags = $derived(
		Array.from(
			new Set<string>(
				(challenges ?? []).flatMap((ch: any) => (ch?.tags ?? []).map((t: any) => String(t)))
			)
		).sort((a, b) => a.localeCompare(b))
	);

	// Debounce search — cleanup on unmount to avoid stale state updates
	$effect(() => {
		const q = search;
		const t = setTimeout(() => {
			debouncedSearch = q;
		}, 300);
		return () => clearTimeout(t);
	});

	const activeFiltersCount = $derived((filterCategories?.length ?? 0) + (filterTags?.length ?? 0));

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
			goto('/team');
		}
	});

	// Self-stopping countdown timer — only runs while at least one instance is active.
	let countdownTimer: ReturnType<typeof setInterval> | undefined;

	function startCountdownTimer() {
		if (countdownTimer) return;
		countdownTimer = setInterval(() => {
			let hasActive = false;
			for (const id in countdowns) {
				if (countdowns[id] > 0) {
					countdowns[id] = countdowns[id] - 1;
					hasActive = true;
				}
			}
			if (!hasActive) {
				clearInterval(countdownTimer);
				countdownTimer = undefined;
			}
		}, 1000);
	}

	// Start timer whenever a non-zero countdown appears; stop when all expire.
	$effect(() => {
		if (Object.values(countdowns).some((v) => v > 0)) {
			startCountdownTimer();
		}
		return () => {
			clearInterval(countdownTimer);
			countdownTimer = undefined;
		};
	});

	function groupByCategory(list: Challenge[]) {
		const map: Record<string, Challenge[]> = {};
		for (const c of list) {
			const label = c.category ?? 'Uncategorized';
			(map[label] ??= []).push(c);
		}
		return Object.entries(map)
			.sort(([a], [b]) => a.localeCompare(b))
			.map(([cat, items]) => [cat, items]) as [string, Challenge[]][];
	}

	const grouped = $derived.by(() => groupByCategory(sortedChallenges));

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

<div
	class="from-muted/20 to-background mb-6 mt-6 rounded-xl border-0 bg-gradient-to-br p-6 shadow-sm"
>
	<div class="flex items-center gap-4">
		<div
			class="bg-background flex h-16 w-16 shrink-0 items-center justify-center rounded-full shadow-sm"
		>
			<Flag class="text-muted-foreground h-8 w-8" />
		</div>
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Challenges</h1>
			<p class="text-muted-foreground mt-2 text-sm">
				"A man who loves to walk will walk more than a man who loves his destination"
			</p>
		</div>
	</div>
</div>

{#if isAdmin}
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
	<div
		class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-600 dark:border-red-800 dark:bg-red-950/20"
	>
		<p class="font-semibold">Error loading challenges</p>
		<p class="text-sm">{error}</p>
	</div>
{:else}
	{#each grouped as [category, items]}
		<section
			class={compactView ? 'mb-4' : 'mb-10'}
			aria-labelledby="category-{category.replace(/\s+/g, '-')}"
		>
			<div class="{compactView ? 'mb-2' : 'mb-3'} flex items-center gap-3">
				<h2
					id="category-{category.replace(/\s+/g, '-')}"
					class="text-2xl font-bold leading-tight text-gray-900 dark:text-white"
				>
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
				<div
					class="grid gap-4 px-0.5 sm:grid-cols-2 lg:grid-cols-3"
					role="list"
					aria-label="{category} challenges"
				>
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
	countdown={selected?.id ? (countdowns[selected.id] ?? 0) : 0}
	{isAdmin}
	onEdit={modifyChallenge}
	onDelete={requestDelete}
	onSolved={handleChallengeSolved}
	onCountdownUpdate={updateCountdown}
	onOpenSolves={() => (solvesOpen = true)}
	onInstanceChange={(updatedChallenge) => {
		if (updatedChallenge) {
			queryClient.setQueryData(['challenges'], (old: any[]) => {
				return old?.map((c) => (c.id === updatedChallenge.id ? updatedChallenge : c)) ?? [];
			});
		} else {
			challengesQuery.refetch();
		}
	}}
/>

{#if selected}
	<SolveListSheet bind:open={solvesOpen} challenge={selected} />
{/if}

{#if DeleteDialogCmp}
	<DeleteDialogCmp bind:open={confirmDeleteOpen} {toDelete} {deleting} onconfirm={confirmDelete} />
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
		{categories}
		{challengeTypes}
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
