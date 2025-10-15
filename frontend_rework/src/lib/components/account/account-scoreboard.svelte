<script lang="ts">
  import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow
  } from "@/components/ui/table";
  import { ChartLine } from "@lucide/svelte";

  // Props: only solves[]
  let { solves } = $props<{ solves: any[] | undefined }>();

  type SortKey = 'name' | 'category' | 'points' | 'timestamp';
  let sortKey = $state<SortKey>('timestamp');
  let sortDir = $state<'asc' | 'desc'>('desc');
  let rows: any[] = $state([]);          // what we actually render
  let totalPoints = $state(0);

  // Helpers
  const getPoints = (s: any) => Number(s?.points ?? s?.score ?? 0);

  const fmtDate = (iso?: string) => {
    if (!iso) return '—';
    const d = new Date(iso);
    return Number.isNaN(+d) ? '—' : d.toLocaleString();
  };

  const timeSince = (iso?: string) => {
    if (!iso) return '—';
    const sec = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 1000));
    const h = Math.floor(sec / 3600);
    const m = Math.floor((sec % 3600) / 60);
    const s = sec % 60;
    if (h > 0) return `${h}h ${m}m`;
    if (m > 0) return `${m}m ${s}s`;
    return `${s}s`;
  };

  function toggleSort(key: SortKey) {
    if (sortKey === key) sortDir = sortDir === 'asc' ? 'desc' : 'asc';
    else { sortKey = key; sortDir = key === 'timestamp' ? 'desc' : 'asc'; }
  }
  const arrow = (key: SortKey | string) =>
    (sortKey === key ? (sortDir === 'asc' ? ' ▲' : ' ▼') : '');

  // Compute rows reactively (robust against undefined / late props)
  $effect(() => {
    const src = Array.isArray(solves) ? solves : [];
    const arr = [...src];

    arr.sort((a: any, b: any) => {
      let av: any, bv: any;
      switch (sortKey) {
        case 'name':
          av = a?.name ?? ''; bv = b?.name ?? ''; break;
        case 'category':
          av = a?.category ?? ''; bv = b?.category ?? ''; break;
        case 'points':
          av = getPoints(a); bv = getPoints(b); break;
        case 'timestamp':
          av = new Date(a?.timestamp ?? 0).getTime();
          bv = new Date(b?.timestamp ?? 0).getTime(); break;
      }
      if (av < bv) return sortDir === 'asc' ? -1 : 1;
      if (av > bv) return sortDir === 'asc' ?  1 : -1;
      return 0;
    });

    rows = arr;
    totalPoints = arr.reduce((acc: number, s: any) => acc + getPoints(s), 0);
  });
</script>

<div class="flex flex-col w-full">
  <div class="flex items-center gap-2">
    <ChartLine class="h-5 w-5 opacity-70" />
    <h3 class="text-xl font-semibold">Solves</h3>
  </div>

  <Table class="w-full">
    <TableCaption class="text-sm">
      {rows.length} solve{rows.length === 1 ? '' : 's'}
      {#if totalPoints > 0} · {totalPoints} pts total{/if}
    </TableCaption>

    <TableHeader>
      <TableRow>
        <TableHead class="w-[40%] cursor-pointer" on:click={() => toggleSort('name')}>
          Challenge {arrow('name')}
        </TableHead>
        <TableHead class="w-[20%] cursor-pointer" on:click={() => toggleSort('category')}>
          Category {arrow('category')}
        </TableHead>
        <TableHead class="w-[15%] text-right cursor-pointer" on:click={() => toggleSort('points')}>
          Points {arrow('points')}
        </TableHead>
        <TableHead class="w-[25%] cursor-pointer text-right sm:text-left" on:click={() => toggleSort('timestamp')}>
          Solved at {arrow('timestamp')}
        </TableHead>
      </TableRow>
    </TableHeader>

    <TableBody>
      {#if rows.length === 0}
        <TableRow>
          <TableCell colspan={4} class="py-10 text-center text-muted-foreground">
            No solves yet.
          </TableCell>
        </TableRow>
      {:else}
        {#each rows as s (s.id ?? s.timestamp ?? s.name ?? Math.random())}
          <TableRow>
            <TableCell class="font-medium">{s.name ?? '—'}</TableCell>
            <TableCell>
              <span class="inline-flex items-center rounded border px-2 py-0.5 text-xs dark:border-gray-700">
                {s.category ?? '—'}
              </span>
            </TableCell>
            <TableCell class="text-right">{getPoints(s)}</TableCell>
            <TableCell class="text-right sm:text-left">
              <div class="flex items-center justify-end sm:justify-start gap-2">
                <span>{fmtDate(s.timestamp)}</span>
                <span class="text-xs text-muted-foreground">({timeSince(s.timestamp)} ago)</span>
              </div>
            </TableCell>
          </TableRow>
        {/each}
      {/if}
    </TableBody>
  </Table>
</div>
