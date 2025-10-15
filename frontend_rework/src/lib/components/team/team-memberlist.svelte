<script lang="ts">
  import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card";
  import { Input } from "@/components/ui/input";
  import { Users } from "@lucide/svelte";
  import { push } from "svelte-spa-router";

  let { team } = $props<{ team: any }>();

  // local state
  let q = $state("");
  let members = $state<any[]>([]);
  let filtered = $state<any[]>([]);

  // helpers
  const norm = (s: any) => String(s ?? "").trim().toLowerCase();

  // tiny fuzzy: exact / prefix / substring / subsequence
  function fuzzyScore(text: string, query: string) {
    const t = norm(text);
    const qn = norm(query);
    if (!qn) return 1e9;              // no query => keep everything (float to top)
    if (t === qn) return 1e6;         // exact
    if (t.startsWith(qn)) return 5e5; // prefix
    if (t.includes(qn)) return 3e5;   // substring
    // subsequence
    let ti = 0, qi = 0, penalty = 0;
    while (ti < t.length && qi < qn.length) {
      if (t[ti] === qn[qi]) qi++;
      else penalty++;
      ti++;
    }
    return qi === qn.length ? 1e5 - penalty : -Infinity;
  }

  const initials = (name: string) =>
    String(name ?? "")
      .split(/\s+/)
      .filter(Boolean)
      .slice(0, 2)
      .map((s) => s[0]?.toUpperCase() ?? "")
      .join("");

  const prettyNum = (n: number) => new Intl.NumberFormat().format(Number(n ?? 0));

  // 1) derive members when team changes
  $effect(() => {
    const raw = Array.isArray(team?.members) ? team.members : [];
    members = raw.map((m: any) => ({
      id: m?.id,
      name: m?.name ?? "—",
      role: m?.role ?? "Member",
      score: Number(m?.score ?? 0),
    }));
  });

  // 2) recompute filtered whenever members or q change
  $effect(() => {
    const list = [...members];
    list.sort((a: any, b: any) => {
      const fa = fuzzyScore(a.name, q);
      const fb = fuzzyScore(b.name, q);
      if (fa !== fb) return fb - fa;           // better fuzzy first
      if (a.score !== b.score) return b.score - a.score; // then score
      return a.name.localeCompare(b.name);     // then name
    });
    filtered = list.filter((m) => fuzzyScore(m.name, q) > -Infinity);
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
      <Input
        placeholder="Search members…"
        bind:value={q}
      />
    </div>
  </div>

  <p class="mb-3 text-sm text-muted-foreground">
    Showing {filtered.length} of {members.length}
  </p>

  <!-- Grid of members -->
  <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
    {#if filtered.length === 0}
      <div class="col-span-full rounded-lg border p-6 text-center text-muted-foreground dark:border-gray-700">
        No members match “{q}”.
      </div>
    {:else}
      {#each filtered as m (m.id ?? m.name)}
        <HoverCard>
          <HoverCardTrigger>
            <button
              type="button"
              class="group flex w-full items-center gap-3 rounded-lg border p-3 text-left transition-colors hover:bg-muted dark:border-gray-700 cursor-pointer"
            >
              <div class="flex h-10 w-10 items-center justify-center rounded-full bg-muted font-semibold">
                {initials(m.name)}
              </div>

              <div class="min-w-0">
                <p class="truncate text-sm font-medium hover:underline cursor-pointer" onclick={() => {push(`/account/${m.id}`)}}>{m.name}</p>
                <p class="truncate text-xs text-muted-foreground">{m.role}</p>
              </div>

              <div class="ml-auto text-right">
                <p class="text-sm font-semibold">{prettyNum(m.score)} pts</p>
              </div>
            </button>
          </HoverCardTrigger>

          <HoverCardContent class="w-80">
            <div class="flex items-start gap-3">
              <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted font-semibold">
                {initials(m.name)}
              </div>
              <div class="min-w-0">
                <p class="truncate text-base font-semibold leading-tight">{m.name}</p>
                <p class="text-xs text-muted-foreground">{m.role}</p>

                <div class="mt-3 grid grid-cols-2 gap-2">
                  <div class="rounded-md border p-2 text-center text-xs dark:border-gray-700">
                    <div class="text-muted-foreground">Score</div>
                    <div class="mt-1 text-sm font-semibold">{prettyNum(m.score)}</div>
                  </div>
                  <div class="rounded-md border p-2 text-center text-xs dark:border-gray-700">
                    <div class="text-muted-foreground">Team total</div>
                    <div class="mt-1 text-sm font-semibold">{prettyNum(team?.score ?? 0)}</div>
                  </div>
                </div>

                {#if team?.name}
                  <p class="mt-2 truncate text-xs text-muted-foreground">
                    Team: <span class="font-medium">{team.name}</span>
                  </p>
                {/if}
              </div>
            </div>
          </HoverCardContent>
        </HoverCard>
      {/each}
    {/if}
  </div>
</div>
