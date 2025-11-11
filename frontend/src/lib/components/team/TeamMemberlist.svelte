<script lang="ts">
	import { Input } from '@/components/ui/input';
	import { Users } from '@lucide/svelte';
	import { push } from 'svelte-spa-router';
	import { Avatar } from 'flowbite-svelte';
	import { getUserData } from '$lib/user';

	let { team } = $props<{ team: any }>();

	// local state
	let q = $state('');
	let members = $state<any[]>([]);
	let filtered = $state<any[]>([]);

	// helpers
	const norm = (s: any) =>
		String(s ?? '')
			.trim()
			.toLowerCase();

	// tiny fuzzy: exact / prefix / substring / subsequence
	function fuzzyScore(text: string, query: string) {
		const t = norm(text);
		const qn = norm(query);
		if (!qn) return 1e9; // no query => keep everything (float to top)
		if (t === qn) return 1e6; // exact
		if (t.startsWith(qn)) return 5e5; // prefix
		if (t.includes(qn)) return 3e5; // substring
		// subsequence
		let ti = 0,
			qi = 0,
			penalty = 0;
		while (ti < t.length && qi < qn.length) {
			if (t[ti] === qn[qi]) qi++;
			else penalty++;
			ti++;
		}
		return qi === qn.length ? 1e5 - penalty : -Infinity;
	}

	const initials = (name: string) =>
		String(name ?? '')
			.split(/\s+/)
			.filter(Boolean)
			.slice(0, 2)
			.map((s) => s[0]?.toUpperCase() ?? '')
			.join('');

	const prettyNum = (n: number) => new Intl.NumberFormat().format(Number(n ?? 0));

	// 1) derive members when team changes
	$effect(() => {
		const raw = Array.isArray(team?.members) ? team.members : [];
		members = raw.map((m: any) => ({
			id: m?.id,
			name: m?.name ?? '-',
			role: m?.role ?? 'Member',
			score: Number(m?.score ?? 0)
		}));
	});

	// 2) recompute filtered whenever members or q change
	$effect(() => {
		const list = [...members];
		list.sort((a: any, b: any) => {
			const fa = fuzzyScore(a.name, q);
			const fb = fuzzyScore(b.name, q);
			if (fa !== fb) return fb - fa; // better fuzzy first
			if (a.score !== b.score) return b.score - a.score; // then score
			return a.name.localeCompare(b.name); // then name
		});
		filtered = list.filter((m) => fuzzyScore(m.name, q) > -Infinity);
	});
	// 3) fetch member images (enrichment)
	let memberImages = $state<Record<string, string>>({});

	$effect(() => {
		const ids = members.map((m: any) => m?.id).filter(Boolean);
		if (ids.length === 0) {
			memberImages = {};
			return;
		}
		let cancelled = false;
		(async () => {
			try {
				const results = await Promise.all(
					ids.map(async (id: any) => {
						try {
							const data = await getUserData(/^\d+$/.test(String(id)) ? Number(id) : id);
							const img = data?.image ?? data?.profileImage ?? null;
							return [String(id), img] as const;
						} catch {
							return [String(id), null] as const;
						}
					})
				);
				if (cancelled) return;
				const map: Record<string, string> = {};
				for (const [id, img] of results) {
					if (img) map[id] = img;
				}
				memberImages = map;
			} catch {
				// ignore enrichment errors
			}
		})();
		return () => {
			cancelled = true;
		};
	});
</script>

<div class="w-full">
	<!-- Header / search -->
	<div class="mb-3 flex items-center gap-3">
		<div class="flex items-center gap-2">
			<Users class="h-5 w-5 opacity-70" />
			<h3 class="text-xl font-semibold">Members</h3>
		</div>

		<div class="ml-auto w-full max-w-xs">
			<!-- bind:value ensures q updates properly -->
			<Input placeholder="Search members..." bind:value={q} />
		</div>
	</div>

	<p class="text-muted-foreground mb-3 text-sm">
		Showing {filtered.length} of {members.length}
	</p>

	<!-- Grid of members -->
	<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
		{#if filtered.length === 0}
			<div
				class="text-muted-foreground col-span-full rounded-lg border p-6 text-center dark:border-gray-700"
			>
				No members match “{q}”.
			</div>
		{:else}
			{#each filtered as m (m.id ?? m.name)}
				<button
					type="button"
					class="hover:bg-muted group flex w-full cursor-pointer items-center gap-3 rounded-lg border p-3 text-left transition-colors dark:border-gray-700"
				>
					{#if memberImages[String(m.id)]}
						<Avatar src={memberImages[String(m.id)]} class="h-10 w-10 shrink-0" />
					{:else}
						<div
							class="bg-muted flex h-10 w-10 shrink-0 items-center justify-center rounded-full font-semibold"
						>
							{initials(m.name)}
						</div>
					{/if}

					<div class="min-w-0 flex-1">
						<span
							class="block cursor-pointer truncate text-sm font-medium hover:underline"
							onclick={() => {
								push(`/account/${m.id}`);
							}}
							onkeydown={(e) => {
								if (e.key === 'Enter' || e.key === ' ') {
									e.preventDefault();
									push(`/account/${m.id}`);
								}
							}}
							role="link"
							tabindex="0"
						>
							{m.name}
						</span>
						<p class="text-muted-foreground truncate text-xs">{m.role}</p>
					</div>

					<div class="ml-auto shrink-0 text-right">
						<p class="text-sm font-semibold">{prettyNum(m.score)} pts</p>
					</div>
				</button>
			{/each}
		{/if}
	</div>
</div>
