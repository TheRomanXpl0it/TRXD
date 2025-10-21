<script lang="ts">
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Pagination from "$lib/components/ui/pagination/index.js";
  import { Medal } from "@lucide/svelte";
  import { getScoreboard } from "@/scoreboard";
  import { onMount } from "svelte";
  import { link, push } from "svelte-spa-router";

  let scoreboardData: any[] = $state([]);
  let perPage = $state(10);

  async function loadScoreboard() {
    const result = await getScoreboard();
    scoreboardData = Array.isArray(result) ? result : [];
  }
  onMount(loadScoreboard);

  // sort by score desc
  const sorted = $derived(
    (Array.isArray(scoreboardData) ? [...scoreboardData] : []).sort(
      (a, b) => (Number(b?.score) || 0) - (Number(a?.score) || 0)
    )
  );
  const count = $derived(sorted.length);

  function medalClass(rank: number) {
    if (rank === 1) return "text-yellow-500";
    if (rank === 2) return "text-gray-300";
    if (rank === 3) return "text-amber-700";
    return "";
  }

  // If using hash routing (default), use "#/team/ID". For history mode, drop the "#".
  const hrefForTeam = (id: number | string) => `#/team/${id}`;
  const pushTeam = (id: number | string) => push(`/team/${id}`);
</script>

<!-- Left-aligned heading/quote (unchanged) -->
<div>
  <p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Scoreboard</p>
  <hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
  <p class="mb-6 text-lg italic text-gray-500 dark:text-gray-400">
    "True success is measured by effort and inner resolve, not just the final score."
  </p>
</div>

<div class="mx-auto max-w-4xl w-full px-4">
  <Pagination.Root {count} {perPage} class="flex flex-col w-full">
    {#snippet children({ pages, currentPage })}
      {@const startIndex = (currentPage - 1) * perPage}
      {@const pageRows   = sorted.slice(startIndex, startIndex + perPage)}
      {@const totalPages = Math.max(1, Math.ceil(count / perPage))}
      {@const singlePage = totalPages <= 1}

      <!-- Table -->
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 w-full">
        <Table.Root class="w-full">
          <Table.Header>
            <Table.Row>
              <Table.Head class="w-20">Rank</Table.Head>
              <Table.Head>Team</Table.Head>
              <Table.Head class="text-right">Score</Table.Head>
              <Table.Head class="w-40">Badges</Table.Head>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            {#if pageRows.length === 0}
              <Table.Row>
                <Table.Cell colspan={4} class="py-8 text-center text-gray-500">
                  No teams yet.
                </Table.Cell>
              </Table.Row>
            {:else}
              {#each pageRows as row, i (row.id)}
                {@const rank = startIndex + i + 1}
                <Table.Row>
                  <Table.Cell class="font-medium">
                    <div class="flex items-center gap-2">
                      <span>#{rank}</span>
                      {#if rank <= 3}
                        <Medal class={`h-4 w-4 ${medalClass(rank)}`} aria-label="Medal" />
                      {/if}
                    </div>
                  </Table.Cell>

                  <Table.Cell>
                    <!-- svelte-spa-router in-page navigation -->
                    <a
                      href={hrefForTeam(row.id)}
                      use:link
                      on:click={(e) => { e.preventDefault(); pushTeam(row.id); }}
                      class="text-primary underline-offset-2 hover:underline cursor-pointer"
                      title={`View team ${row.name}`}
                    >
                      {row.name}
                    </a>
                  </Table.Cell>

                  <Table.Cell class="text-right tabular-nums">{row.score}</Table.Cell>

                  <Table.Cell>
                    {#if Array.isArray(row.badges) && row.badges.length}
                      <div class="flex flex-wrap gap-1">
                        {#each row.badges as b, bi (bi)}
                          <span class="rounded-full border px-2 py-0.5 text-xs">{b.name}</span>
                        {/each}
                      </div>
                    {:else}
                      <span class="text-xs text-gray-500">—</span>
                    {/if}
                  </Table.Cell>
                </Table.Row>
              {/each}
            {/if}
          </Table.Body>
        </Table.Root>
      </div>

      <!-- Pagination directly UNDER the table, centered -->
      <div class="mt-4 flex w-full justify-center">
        <Pagination.Content
          class={`flex items-center justify-center gap-1 ${singlePage ? "opacity-50 pointer-events-none" : ""}`}
          aria-disabled={singlePage}
        >
          <Pagination.Item>
            <Pagination.PrevButton />
          </Pagination.Item>

          {#each pages as page (page.key)}
            {#if page.type === "ellipsis"}
              <Pagination.Item>
                <Pagination.Ellipsis />
              </Pagination.Item>
            {:else}
              <Pagination.Item>
                <Pagination.Link {page} isActive={currentPage === page.value}>
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

      <!-- Optional: range text under pagination, also centered -->
      <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
        Showing <span class="font-medium">{Math.min(count, startIndex + 1)}</span>–
        <span class="font-medium">{Math.min(count, startIndex + perPage)}</span>
        of <span class="font-medium">{count}</span>
      </p>
    {/snippet}
  </Pagination.Root>
</div>
