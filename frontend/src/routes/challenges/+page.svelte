<script lang="ts">
	import {
		CheckCircleSolid,
		FlagSolid,
		BugSolid,
		PenSolid,
		TrashBinSolid,
		UserEditSolid,
		AwardSolid,
		ExclamationCircleSolid
	} from 'flowbite-svelte-icons';
	import { Card, Badge } from 'flowbite-svelte';
	import { Button } from '@/components/ui/button';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import { Container, Download, Droplet, X, Filter, Shapes } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import SolveListSheet from '$lib/components/challenges/solvelist-sheet.svelte';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { push } from 'svelte-spa-router';
	import { Input } from '$lib/components/ui/input/index.js';
	import { getChallenges, getCategories, deleteChallenge } from '$lib/challenges';
	import { startInstance, stopInstance } from '$lib/instances';
	import { submitFlag } from '$lib/challenges';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { user as authUser } from '$lib/stores/auth';
	import { onMount } from 'svelte';

	import { config } from '$lib/env';

	// Lazy-loaded component handles
	type Cmp = typeof import('svelte').SvelteComponent;
	let AdminControlsCmp: Cmp | null = $state(null);
	let CreateModalCmp: Cmp | null = $state(null);
	let DeleteDialogCmp: Cmp | null = $state(null);
	let ChallengeEditSheetCmp: Cmp | null = $state(null);

	// Load admin controls once for admins (small, cheap wrapper)
	$effect(async () => {
		if ($authUser?.role === 'Admin' && !AdminControlsCmp) {
			const mod = await import('$lib/components/challenges/admin-controls.svelte');
			AdminControlsCmp = mod.default;
		}
	});

	// Open “Create Challenge” (lazy-load modal on first use)
	async function openCreate() {
		if (!CreateModalCmp) {
			const mod = await import('$lib/components/challenges/create-challenge-modal.svelte');
			CreateModalCmp = mod.default;
		}
		createChallengeOpen = true;
	}

	// Open “Delete challenge” (lazy-load dialog on first use)
	async function requestDelete(ch: any) {
		toDelete = ch;
		if (!DeleteDialogCmp) {
			const mod = await import('$lib/components/challenges/delete-challenge-dialog.svelte');
			DeleteDialogCmp = mod.default;
		}
		confirmDeleteOpen = true;
	}

	// Open “Edit challenge” (lazy-load sheet on first use)
	function modifyChallenge(ch: any) {
		return async () => {
			if (!ChallengeEditSheetCmp) {
				const mod = await import('$lib/components/challenges/challenge-edit-sheet.svelte');
				ChallengeEditSheetCmp = mod.default;
			}
			editOpen = true;
		};
	}

	// ──────────────────────────────────────────────────────────────────────────────
	// Local state
	// ──────────────────────────────────────────────────────────────────────────────
	let loading = $state(true);
	let error = $state<string | null>(null);
	let createChallengeOpen = $state(false);

	let challenges = $state<any[]>([]);
	let categories = $state<any[]>([]);
	let selected: any | null = $state(null);
	let countdowns: Record<string, number> = $state({});

	let all_tags = $state<any[]>([]);

	let points: number = $state(500);
	let category: any = $state(null);
	let challengeType = $state('Container');
	let challengeName = $state('');
	let challengeDescription = $state('');
	let dynamicScore = $state(false);
	let createLoading = $state(false);

	// Admin controls moved to dedicated component (admin-controls.svelte)

	let openModal = $state(false);
	let solvesOpen = $state(false);
	let editOpen = $state(false);
	let flag = $state('');
	let submittingFlag = $state(false);
	let flagError = $state(false);

	// Filters
	let filterCategories = $state<string[]>([]);
	let filterTags = $state<string[]>([]);
	let search = $state('');
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

	const allTags = $derived(
		Array.from(
			new Set<string>(
				(challenges ?? []).flatMap((ch: any) => (ch?.tags ?? []).map((t: any) => String(t)))
			)
		).sort((a, b) => a.localeCompare(b))
	);

	const filteredChallenges = $derived(
		(challenges ?? [])
			.filter((c: any) => {
				if (!filterCategories || filterCategories.length === 0) return true;
				const cat = c?.category?.name ?? c?.category ?? '';
				return filterCategories.some((fc: string) => norm(cat) === norm(fc));
			})
			.filter((c: any) => {
				if (!filterTags || filterTags.length === 0) return true;
				const tags = (c?.tags ?? []).map((t: any) => String(t));
				return filterTags.every((t: string) => tags.includes(t));
			})
			.filter((c: any) => {
				const q = search.trim();
				if (!q) return true;
				const cat = c?.category?.name ?? c?.category ?? '';
				const tags = (c?.tags ?? []).map((t: any) => String(t));
				const name = c?.name ?? c?.title ?? '';
				return (
					fuzzyScore(name, q) > -Infinity ||
					fuzzyScore(cat, q) > -Infinity ||
					tags.some((t: string) => fuzzyScore(t, q) > -Infinity)
				);
			})
	);

	const activeFiltersCount = $derived(
		(filterCategories?.length ?? 0) + (filterTags?.length ?? 0) + (search.trim() ? 1 : 0)
	);

	// NEW: delete confirmation modal state
	let confirmDeleteOpen = $state(false);
	let deleting = $state(false);
	let toDelete: any = $state(null);

	const challengeTypes = [
		{ value: 'Container', label: 'Container' },
		{ value: 'Compose', label: 'Compose' },
		{ value: 'Normal', label: 'Normal' }
	];

	async function submitCreateChallenge(ev: SubmitEvent) {
		ev.preventDefault();
		if (createLoading) return;
		const trimmedName = challengeName.trim();
		if (!trimmedName) {
			toast.error('Please enter a challenge name.');
			return;
		}
		if (!category) {
			toast.error('Please select a category.');
			return;
		}
		if (!challengeType) {
			toast.error('Please select a challenge type.');
			return;
		}
		if (typeof points !== 'number' || Number.isNaN(points) || points < 0) {
			toast.error('Please choose a valid points value.');
			return;
		}
		createLoading = true;
		const scoretype = dynamicScore ? 'Dynamic' : 'Static';
		try {
			await createChallenge(
				trimmedName,
				category,
				challengeDescription.trim(),
				challengeType,
				points,
				scoretype
			);
			toast.success('Challenge created!');
			createChallengeOpen = false;
			challengeName = '';
			challengeDescription = '';
			category = null;
			challengeType = 'Container';
			dynamicScore = false;
			points = 500;
			await loadChallenges();
		} catch (err: any) {
			const msg = err?.message ?? 'Failed to create challenge.';
			toast.error(msg);
		} finally {
			createLoading = false;
		}
	}

	async function loadChallenges() {
		loading = true;
		error = null;
		try {
		    
		    const prevSelectedId = selected?.id; 
			challenges = await getChallenges();
			if (prevSelectedId != null) {
			//refresh the open modal
              const newer = challenges.find((c: any) => c.id === prevSelectedId);
              if (newer) selected = newer;
            }
			const next: Record<string, number> = {};
			for (const c of challenges ?? []) {
				if (typeof c?.timeout === 'number' && c.timeout > 0) next[c.id] = c.timeout;
				for (const t of c.tags ?? []) {
					if (!(t in all_tags)) all_tags.push(t);
				}
			}
			countdowns = next;
			const uniq = new Map<string, { value: string; label: string }>();
			for (const ch of challenges ?? []) {
				const rawCat = ch?.category;
				const list = Array.isArray(rawCat) ? rawCat : [rawCat];
				for (const item of list) {
					if (!item) continue;
					const label = typeof item === 'string' ? item : (item?.name ?? 'Uncategorized');
					const trimmed = String(label).trim();
					if (!trimmed) continue;
					const value = trimmed.toLowerCase();
					if (!uniq.has(value)) uniq.set(value, { value, label: trimmed });
				}
			}
			categories = Array.from(uniq.values()).sort((a, b) => a.label.localeCompare(b.label));
		} catch (e: any) {
			error = e?.message ?? 'Failed to load challenges';
			toast.error('You need to join a team first!');
			push('/team');
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadChallenges();
		const timer = setInterval(() => {
			for (const id in countdowns) if (countdowns[id] > 0) countdowns[id] = countdowns[id] - 1;
		}, 1000);
		return () => clearInterval(timer);
	});

	$effect(() => {
		if (typeof window === 'undefined') return;
		const timer = setInterval(() => {
			for (const id in countdowns) if (countdowns[id] > 0) countdowns[id] = countdowns[id] - 1;
		}, 1000);
		return () => clearInterval(timer);
	});

	function groupByCategory(list: any[]) {
		const map: Record<string, any[]> = {};
		for (const c of list ?? []) {
			const label = c?.category?.name ?? c?.category ?? 'Uncategorized';
			(map[label] ??= []).push(c);
		}
		return Object.entries(map)
			.sort(([a], [b]) => a.localeCompare(b))
			.map(([cat, items]) => [
				cat,
				items.sort((x, y) => String(x.title || '').localeCompare(String(y.title || '')))
			]) as [string, any[]][];
	}
	const grouped = $derived(groupByCategory(filteredChallenges));

	function fmtTimeLeft(total: number | undefined): string {
		if (!total || total < 0) total = 0;
		const h = Math.floor(total / 3600);
		const m = Math.floor((total % 3600) / 60);
		const s = Math.floor(total % 60);
		if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
		if (m > 0) return `${m}:${String(s).padStart(2, '0')}`;
		return `${s}`;
	}

	function openChallenge(ch: any) {
		selected = ch;
		openModal = true;
	}
	function closeModal() {
		openModal = false;
	}
	$effect(() => {
		if (!openModal) selected = null;
	});

	function copyToClipboard(text: string) {
		if (typeof navigator === 'undefined') return;
		navigator.clipboard
			.writeText(text)
			.then(() => toast.success('Copied to clipboard!'))
			.catch(() => toast.error('Failed to copy to clipboard.'));
	}

	async function createInstance(ch: any) {
		try {
			const { host, port, timeout } = await startInstance(ch.id);
			ch.host = host;
			ch.port = port;
			ch.timeout = timeout;
			if (typeof ch.timeout === 'number') countdowns[ch.id] = Math.max(0, ch.timeout);
			toast.success('Created instance!');
		} catch (err: any) {
			console.error(err);
			toast.error(`Failed to create instance: ${err?.message ?? err}`);
		}
	}

	async function destroyInstance(ch: any) {
		try {
			await stopInstance(ch.id);
			ch.host = null;
			ch.port = null;
			ch.timeout = null;
			countdowns[ch.id] = 0;
			toast.success('Stopped instance!');
		} catch (err: any) {
			console.error(err);
			toast.error(`Failed to stop instance: ${err?.message ?? err}`);
		}
	}

	async function onSubmitFlag(ev: SubmitEvent) {
		ev.preventDefault();
		if (!selected?.id) {
			toast.error('No challenge selected');
			return;
		}
		const value = flag.trim();
		if (!value) return;

		submittingFlag = true;
		try {
			const res = await submitFlag(selected.id, value);
			if ((res as any).status === 'Wrong') {
				flagError = true;
				toast.error('Incorrect flag');
				return;
			} else if ((res as any).first_blood) {
				toast.success('First blood! 🎉');
			} else {
				toast.success('Correct flag!');
			}
			flag = '';
			selected.solved = true;
			const idx = challenges.findIndex((c: any) => c.id === selected!.id);
			if (idx !== -1) challenges[idx] = { ...challenges[idx], solved: true };
		} catch (e: any) {
			toast.error(e?.message ?? 'Flag submission failed');
		} finally {
			submittingFlag = false;
		}
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
			await loadChallenges();
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

{#if $authUser?.role === 'Admin' && AdminControlsCmp}
  <AdminControlsCmp
    on:open-create={openCreate}
    on:category-created={() => loadChallenges()}
  />
{/if}


<!-- Filters -->
<div class="mb-4 flex flex-wrap items-center gap-3">
	<Popover.Root bind:open={categoriesOpen}>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button {...props} variant="outline" class="flex cursor-pointer items-center gap-2">
					<Shapes class="h-4 w-4" />
					Categories
					{#if filterCategories.length > 0}
						<span class="bg-primary text-primary-foreground ml-1 rounded px-2 py-0.5 text-xs"
							>{filterCategories.length}</span
						>
					{/if}
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-[260px] p-1">
			<Command.Root>
				<Command.Input
					placeholder="Search categories…"
					class="border-0 shadow-none ring-0 focus:outline-none focus:ring-0 focus-visible:outline-none focus-visible:ring-0"
				/>
				<Command.List>
					<Command.Empty>No categories.</Command.Empty>
					<Command.Group value="categories">
						{#each categories as c (c.value)}
							<Command.Item
								value={c.value}
								onSelect={() => {
									if (filterCategories.includes(c.value)) {
										filterCategories = filterCategories.filter((x) => x !== c.value);
									} else {
										filterCategories = [...filterCategories, c.value];
									}
								}}
							>
								<Checkbox checked={filterCategories.includes(c.value)} />
								<span class="ml-2">{c.label}</span>
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
		</Popover.Content>
	</Popover.Root>
	<Popover.Root bind:open={tagsOpen}>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button {...props} variant="outline" class="flex cursor-pointer items-center gap-2">
					<Filter class="h-4 w-4" />
					Tags
					{#if filterTags.length > 0}
						<span class="bg-primary text-primary-foreground ml-1 rounded px-2 py-0.5 text-xs"
							>{filterTags.length}</span
						>
					{/if}
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-[260px] p-1">
			<Command.Root>
				<Command.Input
					placeholder="Search tags…"
					class="border-0 shadow-none ring-0 focus:outline-none focus:ring-0 focus-visible:outline-none focus-visible:ring-0"
				/>
				<Command.List>
					<Command.Empty>No tags.</Command.Empty>
					<Command.Group value="tags">
						{#each allTags as t (t)}
							<Command.Item
								value={t}
								onSelect={() => {
									if (filterTags.includes(t)) {
										filterTags = filterTags.filter((x) => x !== t);
									} else {
										filterTags = [...filterTags, t];
									}
								}}
							>
								<Checkbox checked={filterTags.includes(t)} />
								<span class="ml-2">{t}</span>
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
		</Popover.Content>
	</Popover.Root>

	<div class="ml-auto flex items-center gap-2">
		<Input id="search" placeholder="Search challenges…" bind:value={search} class="w-[260px]" />
		{#if activeFiltersCount > 0}
			<Badge color="cyan">{activeFiltersCount}</Badge>
		{/if}
		<Button
			variant="ghost"
			size="sm"
			class="cursor-pointer"
			onclick={() => {
				filterCategories = [];
				filterTags = [];
				search = '';
			}}>Clear</Button
		>
	</div>
</div>

{#if loading}
	<div class="p-4">Loading challenges…</div>
{:else if error}
	<div class="p-4 text-red-600">{error}</div>
{:else}
	{#each grouped as [category, items]}
		<section class="mb-10">
			<div class="mb-3 flex items-center gap-3">
				<p class="text-2xl font-bold leading-tight text-gray-900 dark:text-white">{category}</p>
				<span class="text-sm text-gray-500 dark:text-gray-400">
					{items.length} challenge{items.length === 1 ? '' : 's'}
				</span>
			</div>

			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each items as ch}
					<Card
						class={`min-h-35 max-w-90 min-w-55 border-1 border-solid border-stone-900 transition-shadow hover:cursor-pointer hover:shadow-md dark:border-stone-300
            ${ch.hidden ? 'border-2 border-dashed !border-gray-300 dark:!border-gray-600' : ''}`}
						onclick={() => openChallenge(ch)}
					>
						<div class="p-4">
							<p class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">{ch.name}</p>
							{#each ch.tags as tag}
								<Badge class="mr-1" color="gray">{tag}</Badge>
							{/each}
						</div>

						<div class="mt-auto flex">
							{#if ch.solved}
								<Badge color="green" class="ml-1.5 mr-auto self-center">{ch.points}</Badge>
								<CheckCircleSolid class="mb-2 mr-2 self-center text-green-500" />
							{:else}
								<Badge color="secondary" class="mb-1.5 ml-1.5 self-center">{ch.points}</Badge>
							{/if}
							{#if ch.instance}
								<Container class="mb-2 mr-2 self-center {ch.solved ? '' : 'ml-auto'}" />
								{#if countdowns[ch.id] > 0}
									<Badge color="blue" class="self-center mr-1.5">{fmtTimeLeft(countdowns[ch.id])}</Badge>
								{/if}
							{/if}
						</div>
					</Card>
				{/each}
			</div>
		</section>
	{/each}
{/if}

<!-- One global dialog (not inside the loop) -->
<Dialog.Root bind:open={openModal}>
	<Dialog.Content class="sm:max-w-[720px]">
		<Dialog.Header class="pb-3">
			<div class="flex items-center gap-3">
				<Dialog.Title class="text-xl font-semibold text-gray-900 dark:text-white">
					{selected?.name}
				</Dialog.Title>
				<BugSolid class="ml-auto mr-auto h-6 w-6 text-gray-800" />
				{#if $authUser?.role === 'Admin'}
					<Button
						variant="ghost"
						size="icon"
						class="cursor-pointer"
						onclick={modifyChallenge(selected)}
					>
						<PenSolid class="h-5 w-5" />
					</Button>
					<Button
						variant="ghost"
						size="icon"
						class="cursor-pointer mr-5"
						onclick={() => requestDelete(selected)}
					>
						<TrashBinSolid class="h-5 w-5" />
					</Button>
				{/if}
			</div>
			<Dialog.Description class="sr-only">Challenge details</Dialog.Description>
		</Dialog.Header>

		<!-- Tags -->
		<div class="mb-4">
			{#each selected?.tags as tag}
				<Badge class="mr-1" color="cyan">{tag}</Badge>
			{/each}
		</div>

		<!-- Solves & authors -->
		<div class="flex flex-row">
			<span class="flex flex-row">
				{#if selected?.solves === 0}
					<Droplet class="mr-1 text-red-500" />
					<p>0 solves, be the first!</p>
				{:else}
					<Button
						onclick={() => (solvesOpen = true)}
						size="sm"
						class="hover:cursor-pointer"
						variant="outline"
					>
						<AwardSolid class="mr-1" />
						{#if selected?.solves === 1}
							<p>1 solve</p>
						{:else}
							<p>{selected?.solves} solves</p>
						{/if}
					</Button>
				{/if}
			</span>
			<span class="ml-auto flex flex-row">
				<UserEditSolid class="mr-1" />
				<span>
					{#each selected?.authors as author, i (author)}
						{author}{i < (selected?.authors?.length ?? 0) - 1 ? ', ' : ''}
					{/each}
				</span>
			</span>
		</div>

		<!-- Description -->
		<div class="mt-5 flex flex-row items-center">
			{selected?.description}
		</div>

		<!-- Attatchments  -->
		<div class="mt-3 flex flex-row items-center">
			{#each selected?.attachments as attatchment}
				<a
					href={config.getBackendUrl(attatchment)}
					target="_blank"
					rel="external"
					download
					class="mr-2 flex cursor-pointer items-center gap-1 rounded bg-gray-100 px-2 py-1 text-sm hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600"
				>
					<Download class="h-3 w-3" />
					{attatchment.split('/').pop()}
				</a>
			{/each}
		</div>

		<!-- Instance / remote -->
		<div class="mt-1 flex w-full flex-row items-center justify-center px-6">
			{#if selected?.instance}
				{#if countdowns[selected?.id] > 0}
					<Button
						size="sm"
						style="background-color:#779ecb;"
						disabled
						class="mr-2 w-full hover:cursor-pointer"
					>
						<Container class="mr-1 " />
						<span>Instance Running ({fmtTimeLeft(countdowns[selected?.id])})</span>
					</Button>
					<Button
						variant="destructive"
						size="sm"
						onclick={() => destroyInstance(selected)}
						class="hover:cursor-pointer"
					>
						<X />
					</Button>
				{:else}
					<Button
						style="background-color:#779ecb;"
						size="sm"
						onclick={() => createInstance(selected)}
						class="hover:cursor-pointer"
					>
						<Container class="mr-1" />
						<span>Start challenge instance</span>
					</Button>
				{/if}
			{/if}
		</div>

		<div class="mt-1 flex flex-row items-center justify-center">
			{#if selected?.host}
				<Badge
					color="gray"
					class="cursor-pointer"
					onclick={() =>
						copyToClipboard(`${selected?.host}${selected?.port ? `:${selected?.port}` : ''}`)}
				>
					<p class="text-lg">{selected?.host}{selected?.port ? ` ${selected?.port}` : ''}</p>
				</Badge>
			{/if}
		</div>

		<!-- Submit flag -->
		<div class="mt-4 flex w-full items-center justify-between">
			<form
				class="mt-4 flex w-full items-center gap-2"
				class:justify-center={selected?.solved}
				onsubmit={onSubmitFlag}
			>
				{#if !selected?.solved}
					<div class="relative flex-1">
						{#if flagError}
							<ExclamationCircleSolid
								class="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-red-500"
								aria-hidden="true"
							/>
						{:else}
							<FlagSolid
								class="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-500 dark:text-gray-400"
								aria-hidden="true"
							/>
						{/if}
						<Input
							class="pl-10"
							placeholder="TRX{'...'}"
							bind:value={flag}
							oninput={() => (flagError = false)}
							aria-invalid={flagError}
							data-error={flagError}
						/>
					</div>

					<Button
						type="submit"
						color="primary"
						class="h-full"
						disabled={submittingFlag || !flag.trim() || flagError}
					>
						{#if submittingFlag}
							<Spinner />
							Submitting...
						{:else}
							Submit
						{/if}
					</Button>
				{:else}
					<Badge color="green" class="flex items-center">Challenge solved</Badge>
				{/if}
			</form>
		</div>
	</Dialog.Content>
</Dialog.Root>

<!-- Solve list -->
<SolveListSheet bind:open={solvesOpen} challenge={selected} />


<!-- Delete Confirmation Modal -->
{#if DeleteDialogCmp}
  <DeleteDialogCmp
    bind:open={confirmDeleteOpen}
    {toDelete}
    {deleting}
    on:confirm={confirmDelete}
  />
{/if}

<!-- Create Challenge Modal -->
{#if CreateModalCmp}
  <CreateModalCmp
    bind:open={createChallengeOpen}
    bind:challengeName
    bind:challengeDescription
    bind:category
    bind:challengeType
    bind:points
    bind:dynamicScore
    bind:categories={categories}
    challengeTypes={challengeTypes}
    on:created={loadChallenges}
  />
{/if}

<!-- All sheets that are imported -->
{#if ChallengeEditSheetCmp}
  <ChallengeEditSheetCmp
    bind:open={editOpen}
    challenge_user={selected}
    on:updated={() => loadChallenges()}
    {all_tags}
  />
{/if}