<script lang="ts">
	import {
		Carousel,
		CarouselContent,
		CarouselItem,
		CarouselNext,
		CarouselPrevious
	} from '@/components/ui/carousel';
	import { Trophy, Award, Users, Clock, CalendarClock, Crown, Star } from '@lucide/svelte';
	import { push } from 'svelte-spa-router';
	import { Avatar } from 'flowbite-svelte';
	import { getUserData } from '$lib/user';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	let { team } = $props<{ team: any }>();

	// derived stats
	const badgesCount = $derived(team?.badges?.length ?? 0);
	const members = $derived(team?.members ?? []);
	const membersCount = $derived(members.length);
	const activeMembers = $derived(members.filter((m: any) => (m?.score ?? 0) > 0).length);
	const totalScore = $derived(team?.score ?? 0);
	const sortedMembers = $derived(
		[...(members as any[])].sort((a: any, b: any) => (b.score ?? 0) - (a.score ?? 0))
	);
	const topMember = $derived(sortedMembers[0]);
	const solves = $derived(team?.solves ?? []);
	const solvesCount = $derived(solves.length);

	const lastSolve: any = $derived(
		solves.length
			? [...solves].sort((a: any, b: any) => +new Date(b.timestamp) - +new Date(a.timestamp))[0]
			: null
	);

	const categories: any = $derived(() => {
		const map = new Map<string, number>();
		for (const s of solves) map.set(s.category, (map.get(s.category) ?? 0) + 1);
		const total = [...map.values()].reduce((a, b) => a + b, 0) || 1;
		return [...map.entries()]
			.sort((a, b) => b[1] - a[1])
			.map(([cat, count]) => ({ cat, count, pct: Math.round((count / total) * 100) }));
	});

	// member image enrichment for leaderboard
	const leaderboardMembers = $derived(sortedMembers.slice(0, 5));
	let memberImages: Record<string, string> = $state({});

	$effect(() => {
		const ids = leaderboardMembers.map((m: any) => m?.id).filter(Boolean);
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
				// ignore errors
			}
		})();
		return () => {
			cancelled = true;
		};
	});

	// utils
	const timeSince = (iso?: string) => {
		if (!iso) return '-';
		const sec = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 1000));
		const h = Math.floor(sec / 3600);
		const m = Math.floor((sec % 3600) / 60);
		const s = sec % 60;
		if (h > 0) return `${h}h ${m}m`;
		if (m > 0) return `${m}m ${s}s`;
		return `${s}s`;
	};
	const barWidth = (pct: number) => `width:${Math.max(3, Math.min(100, pct))}%`;

	// optional API handle
	let api: any = null;
	const setApi = (a: any) => (api = a);
</script>

<div class="mx-auto mt-5 w-full max-w-4xl">
	<Carousel
		on:init={(e) => setApi(e.detail)}
		opts={{ align: 'start', loop: true }}
		class="px-15 relative"
	>
		<CarouselContent>
			<!-- Slide 1: Overview -->
			<CarouselItem class="md:basis-3/3 h-full">
				<div class="w-full rounded-xl border p-6 dark:border-gray-700">
					<div class="flex items-center gap-3">
						<Crown class="h-6 w-6 shrink-0 opacity-80" />
						<h2 class="truncate text-2xl font-semibold">{team?.name ?? 'Team'}</h2>
						{#if team?.country}
							<span class="ml-2 shrink-0 rounded bg-gray-100 px-2 py-0.5 text-xs dark:bg-gray-800">
								{team.country}
							</span>
						{/if}
					</div>

					<div class="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-4">
						<div class="rounded-lg border p-4 dark:border-gray-700">
							<div class="text-muted-foreground flex items-center gap-2 text-sm">
								<Trophy class="h-4 w-4" /> Total points
							</div>
							<p class="mt-1 text-2xl font-semibold">{totalScore}</p>
						</div>

						<div class="rounded-lg border p-4 dark:border-gray-700">
							<div class="text-muted-foreground flex items-center gap-2 text-sm">
								<Users class="h-4 w-4" /> Members
							</div>
							<p class="mt-1 text-2xl font-semibold">{membersCount}</p>
							<p class="text-muted-foreground text-xs">{activeMembers} active</p>
						</div>

						<div class="rounded-lg border p-4 dark:border-gray-700">
							<div class="text-muted-foreground flex items-center gap-2 text-sm">
								<Award class="h-4 w-4" /> Badges
							</div>
							<p class="mt-1 text-2xl font-semibold">{badgesCount}</p>
						</div>

						<div class="rounded-lg border p-4 dark:border-gray-700">
							<div class="text-muted-foreground flex items-center gap-2 text-sm">
								<CalendarClock class="h-4 w-4" /> Solves
							</div>
							<p class="mt-1 text-2xl font-semibold">{solvesCount}</p>
						</div>
					</div>

					{#if topMember}
					<div class="bg-muted mt-6 rounded-lg p-4">
						<div class="text-muted-foreground flex items-center gap-2 text-sm">
							<Star class="h-4 w-4" /> Top member
						</div>
						<div class="mt-1 flex items-center justify-between">
							<button
								type="button"
								class="cursor-pointer text-lg font-medium hover:underline"
								onclick={() => push(`/account/${topMember.id}`)}>{topMember.name}</button
							>
								<p class="text-lg font-semibold">{topMember.score}pts</p>
							</div>
						</div>
					{/if}
				</div>
			</CarouselItem>

			<!-- Slide 2: Last activity -->
			<CarouselItem class="md:basis-3/3 h-full">
				<div class="w-full rounded-xl border p-6 dark:border-gray-700">
					<div class="flex items-center gap-2">
						<Clock class="h-5 w-5 opacity-80" />
						<h3 class="text-xl font-semibold">Last activity</h3>
					</div>

					{#if lastSolve}
						<div class="mt-4">
							<p class="text-muted-foreground text-sm">Most recent solve</p>
							<div class="mt-1 flex items-center justify-between">
								<div>
									<p class="text-lg font-medium">{lastSolve.name}</p>
									<span
										class="mt-1 inline-flex rounded border px-2 py-0.5 text-xs dark:border-gray-700"
									>
										{lastSolve.category}
									</span>
								</div>
								<div class="text-right">
									<p class="text-2xl font-semibold">{timeSince(lastSolve.timestamp)}</p>
									<p class="text-muted-foreground text-xs">ago</p>
								</div>
							</div>
							<p class="text-muted-foreground mt-2 text-xs">
								{new Date(lastSolve.timestamp).toLocaleString()}
							</p>
						</div>

						{#if categories.length}
							<div class="mt-6">
								<p class="text-muted-foreground text-sm">Categories hit recently</p>
								<div class="mt-2 flex flex-wrap gap-2">
									{#each categories.slice(0, 6) as c}
										<span
											class="inline-flex items-center rounded border px-2 py-0.5 text-xs dark:border-gray-700"
										>
											{c.cat} · {c.count}
										</span>
									{/each}
								</div>
							</div>
						{/if}
					{:else}
						<p class="text-muted-foreground mt-6">No solves yet.</p>
					{/if}
				</div>
			</CarouselItem>

			<!-- Slide 3: Members leaderboard -->
			<CarouselItem class="md:basis-3/3 h-full">
				<div class="w-full rounded-xl border p-6 dark:border-gray-700">
					<div class="flex items-center gap-2">
						<Users class="h-5 w-5 opacity-80" />
						<h3 class="text-xl font-semibold">Members</h3>
					</div>

					{#if sortedMembers.length}
						<div class="mt-4 space-y-3">
							{#each sortedMembers.slice(0, 5) as m, i}
								<div class="rounded-lg border p-3 dark:border-gray-700">
									<div class="flex items-center justify-between">
										<div class="flex items-center gap-3">
											<span
												class="bg-muted inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold"
											>
												{i + 1}
											</span>
										{#if memberImages[String(m.id)]}
											<Avatar src={memberImages[String(m.id)]} class="h-8 w-8" />
										{/if}
										<div>
											<button
												type="button"
												class="cursor-pointer font-medium leading-tight hover:underline"
												onclick={() => {
													push(`/account/${m.id}`);
												}}>{m.name}</button
											>
											<p class="text-muted-foreground text-xs">{m.role}</p>
										</div>
										</div>
										<p class="text-right text-lg font-semibold">{m.score}pts</p>
									</div>
								{#if topMember?.score}
									<div class="bg-muted mt-2 h-2 w-full rounded">
										<div
											class="bg-primary h-2 rounded"
											style={barWidth(Math.round((m.score / Math.max(1, topMember.score)) * 100))}
										></div>
									</div>
								{/if}
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-muted-foreground mt-6">No members.</p>
					{/if}
				</div>
			</CarouselItem>

			<!-- Slide 4: Category breakdown -->
			<CarouselItem class="md:basis-3/3 h-full">
				<div class="w-full rounded-xl border p-6 dark:border-gray-700">
					<div class="flex items-center gap-2">
						<CalendarClock class="h-5 w-5 opacity-80" />
						<h3 class="text-xl font-semibold">Category breakdown</h3>
					</div>

					{#if categories.length}
						<div class="mt-4 space-y-3">
							{#each categories as c}
								<div>
									<div class="flex items-center justify-between text-sm">
										<span class="font-medium">{c.cat}</span>
										<span class="text-muted-foreground">{c.count} ({c.pct}%)</span>
								</div>
								<div class="bg-muted mt-1 h-2 w-full rounded">
									<div class="h-2 rounded bg-emerald-500" style={barWidth(c.pct)}></div>
								</div>
							</div>
							{/each}
						</div>
					{:else}
						<p class="text-muted-foreground mt-6">No category stats yet.</p>
					{/if}
				</div>
			</CarouselItem>

			<!-- Slide 5: Badges -->
			<CarouselItem class="md:basis-3/3 h-full">
				<div class="w-full rounded-xl border p-6 dark:border-gray-700">
					<div class="flex items-center gap-2">
						<Award class="h-5 w-5 opacity-80" />
						<h3 class="text-xl font-semibold">Badges</h3>
					</div>

					{#if badgesCount === 0}
						<p class="text-muted-foreground mt-6">No badges yet.</p>
					{:else}
						<div class="mt-4 grid gap-3 sm:grid-cols-2">
							{#each team?.badges ?? [] as b}
								<Tooltip.Root>
									<Tooltip.Trigger asChild>
										<div class="rounded-lg border p-3 dark:border-gray-700">
											<div class="flex items-center gap-2">
												<Star class="h-4 w-4 text-yellow-500" />
												<p class="font-medium">{b.name}</p>
											</div>
											{#if b.description}
												<p class="text-muted-foreground mt-1 text-sm">{b.description}</p>
											{/if}
										</div>
									</Tooltip.Trigger>
									{#if b.description}
										<Tooltip.Content>
											<p>{b.description}</p>
										</Tooltip.Content>
									{/if}
								</Tooltip.Root>
							{/each}
						</div>
					{/if}
				</div>
			</CarouselItem>
		</CarouselContent>

		<CarouselPrevious class="absolute left-2 top-1/2 -translate-y-1/2" />
		<CarouselNext class="absolute right-2 top-1/2 -translate-y-1/2" />
	</Carousel>
</div>
