<script lang="ts">
	import { getUserData } from '$lib/user';
	import GeneratedAvatar from '$lib/components/ui/avatar/generated-avatar.svelte';

	let { team } = $props<{ team: any }>();

	// local state
	let members = $state<any[]>([]);
	let filtered = $state<any[]>([]);

	// helpers

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
	// 2) recompute filtered whenever members change
	$effect(() => {
		const list = [...members];
		list.sort((a: any, b: any) => {
			if (a.score !== b.score) return b.score - a.score; // then score
			return a.name.localeCompare(b.name); // then name
		});
		filtered = list;
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
	<!-- Grid of members -->
	<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
		{#if filtered.length === 0}
			<div
				class="text-muted-foreground bg-muted/20 col-span-full rounded-lg border-0 p-6 text-center"
			>
				No members found.
			</div>
		{:else}
			{#each filtered as m (m.id ?? m.name)}
				<a
					href={`/account/${m.id}`}
					class="bg-muted/40 hover:bg-background group flex w-full cursor-pointer items-center gap-3 rounded-lg p-3 text-left transition-all hover:shadow-sm"
				>
					{#if memberImages[String(m.id)]}
						<div
							class="h-10 w-10 shrink-0 overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700"
						>
							<img
								src={memberImages[String(m.id)]}
								alt={m.name}
								class="h-full w-full rounded-full object-cover object-center"
							/>
						</div>
					{:else}
						<div class="border-border h-10 w-10 shrink-0 overflow-hidden rounded-full border">
							<GeneratedAvatar seed={m.name} class="h-full w-full" />
						</div>
					{/if}

					<div class="min-w-0 flex-1">
						<span class="block truncate text-sm font-medium">
							{m.name}
						</span>
						<p class="text-muted-foreground truncate text-xs">{m.role}</p>
					</div>

					<div class="ml-auto shrink-0 text-right">
						<p class="text-sm font-semibold">{prettyNum(m.score)} pts</p>
					</div>
				</a>
			{/each}
		{/if}
	</div>
</div>
