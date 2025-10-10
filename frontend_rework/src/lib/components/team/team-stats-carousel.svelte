<script lang="ts">
  import {
    Carousel,
    CarouselContent,
    CarouselItem,
    CarouselNext,
    CarouselPrevious,
  } from "@/components/ui/carousel";
  import { Trophy, Award, Users, Clock, CalendarClock, Crown, Star } from "@lucide/svelte";

  let { team } = $props<{ team: any }>();

  // ── derived stats (all `any`) ────────────────────────────────────────────────
  const badgesCount   = $derived(team?.badges?.length ?? 0);
  const members       = $derived(team?.members ?? []);
  const membersCount  = $derived(members.length);
  const activeMembers = $derived(members.filter((m: any) => (m?.score ?? 0) > 0).length);
  const totalScore    = $derived(team?.score ?? 0);
  const sortedMembers = $derived([...(members as any[])].sort((a: any, b: any) => (b.score ?? 0) - (a.score ?? 0)));
  const topMember     = $derived(sortedMembers[0]);
  const solves        = $derived(team?.solves ?? []);
  const solvesCount   = $derived(solves.length);

  const lastSolve:any = $derived(
    solves.length
      ? [...solves].sort((a: any, b: any) => +new Date(b.timestamp) - +new Date(a.timestamp))[0]
      : null
  );

  const categories:any = $derived(() => {
    const map = new Map<string, number>();
    for (const s of solves) map.set(s.category, (map.get(s.category) ?? 0) + 1);
    const total = [...map.values()].reduce((a, b) => a + b, 0) || 1;
    return [...map.entries()]
      .sort((a, b) => b[1] - a[1])
      .map(([cat, count]) => ({ cat, count, pct: Math.round((count / total) * 100) }));
  });

  // utils
  const timeSince = (iso?: string) => {
    if (!iso) return "—";
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

<div class="mx-auto w-full max-w-4xl mt-5">
  <Carousel
    on:init={(e) => setApi(e.detail)}
    opts={{ align: "start", loop: true }}
    class="relative px-15"
  >
    <CarouselContent>
      <!-- Slide 1: Overview -->
      <CarouselItem class="md:basis-3/3 h-full">
        <div class="w-full rounded-xl border p-6 dark:border-gray-700">
          <div class="flex items-center gap-3">
            <Crown class="h-6 w-6 opacity-80" />
            <h2 class="text-2xl font-semibold">{team?.name ?? "Team"}</h2>
            {#if team?.country}
              <span class="ml-2 rounded bg-gray-100 px-2 py-0.5 text-xs dark:bg-gray-800">
                {team.country}
              </span>
            {/if}
          </div>

          <div class="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-4">
            <div class="rounded-lg border p-4 dark:border-gray-700">
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <Trophy class="h-4 w-4" /> Total points
              </div>
              <p class="mt-1 text-2xl font-semibold">{totalScore}</p>
            </div>

            <div class="rounded-lg border p-4 dark:border-gray-700">
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <Users class="h-4 w-4" /> Members
              </div>
              <p class="mt-1 text-2xl font-semibold">{membersCount}</p>
              <p class="text-xs text-muted-foreground">{activeMembers} active</p>
            </div>

            <div class="rounded-lg border p-4 dark:border-gray-700">
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <Award class="h-4 w-4" /> Badges
              </div>
              <p class="mt-1 text-2xl font-semibold">{badgesCount}</p>
            </div>

            <div class="rounded-lg border p-4 dark:border-gray-700">
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CalendarClock class="h-4 w-4" /> Solves
              </div>
              <p class="mt-1 text-2xl font-semibold">{solvesCount}</p>
            </div>
          </div>

          {#if topMember}
            <div class="mt-6 rounded-lg bg-muted p-4">
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <Star class="h-4 w-4" /> Top member
              </div>
              <div class="mt-1 flex items-center justify-between">
                <p class="text-lg font-medium">{topMember.name}</p>
                <p class="text-lg font-semibold">{topMember.score}</p>
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
              <p class="text-sm text-muted-foreground">Most recent solve</p>
              <div class="mt-1 flex items-center justify-between">
                <div>
                  <p class="text-lg font-medium">{lastSolve.name}</p>
                  <span class="mt-1 inline-flex rounded border px-2 py-0.5 text-xs dark:border-gray-700">
                    {lastSolve.category}
                  </span>
                </div>
                <div class="text-right">
                  <p class="text-2xl font-semibold">{timeSince(lastSolve.timestamp)}</p>
                  <p class="text-xs text-muted-foreground">ago</p>
                </div>
              </div>
              <p class="mt-2 text-xs text-muted-foreground">
                {new Date(lastSolve.timestamp).toLocaleString()}
              </p>
            </div>

            {#if categories.length}
              <div class="mt-6">
                <p class="text-sm text-muted-foreground">Categories hit recently</p>
                <div class="mt-2 flex flex-wrap gap-2">
                  {#each categories.slice(0, 6) as c}
                    <span class="inline-flex items-center rounded border px-2 py-0.5 text-xs dark:border-gray-700">
                      {c.cat} · {c.count}
                    </span>
                  {/each}
                </div>
              </div>
            {/if}
          {:else}
            <p class="mt-6 text-muted-foreground">No solves yet.</p>
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
                      <span class="inline-flex h-6 w-6 items-center justify-center rounded-full bg-muted text-xs font-semibold">
                        {i + 1}
                      </span>
                      <div>
                        <p class="font-medium leading-tight">{m.name}</p>
                        <p class="text-xs text-muted-foreground">{m.role}</p>
                      </div>
                    </div>
                    <p class="text-right text-lg font-semibold">{m.score}pts</p>
                  </div>
                  {#if topMember?.score}
                    <div class="mt-2 h-2 w-full rounded bg-muted">
                      <div
                        class="h-2 rounded bg-primary"
                        style={barWidth(Math.round((m.score / Math.max(1, topMember.score)) * 100))}
                      />
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {:else}
            <p class="mt-6 text-muted-foreground">No members.</p>
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
                  <div class="mt-1 h-2 w-full rounded bg-muted">
                    <div class="h-2 rounded bg-emerald-500" style={barWidth(c.pct)} />
                  </div>
                </div>
              {/each}
            </div>
          {:else}
            <p class="mt-6 text-muted-foreground">No category stats yet.</p>
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
            <p class="mt-6 text-muted-foreground">No badges yet.</p>
          {:else}
            <div class="mt-4 grid gap-3 sm:grid-cols-2">
              {#each team?.badges ?? [] as b}
                <div class="rounded-lg border p-3 dark:border-gray-700">
                  <div class="flex items-center gap-2">
                    <Star class="h-4 w-4 text-yellow-500" />
                    <p class="font-medium">{b.name}</p>
                  </div>
                  {#if b.description}
                    <p class="mt-1 text-sm text-muted-foreground">{b.description}</p>
                  {/if}
                </div>
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