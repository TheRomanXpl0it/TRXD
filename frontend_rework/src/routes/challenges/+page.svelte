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
  import {
    SpeedDial,
    SpeedDialTrigger,
    SpeedDialButton,
    Input,
    Card,
    Badge
  } from 'flowbite-svelte';
  import { Button } from '@/components/ui/button';
  import { Container, Droplet, X } from '@lucide/svelte';
  import { toast } from 'svelte-sonner';
  import * as Dialog from '$lib/components/ui/dialog/index.js';
  import SolveListSheet from '$lib/components/challenges/solvelist-sheet.svelte';
  import { Spinner } from "$lib/components/ui/spinner/index.js";

  // 🔁 Services (no fetch arg now)
  import { getChallenges } from '$lib/challenges';
  import { startInstance, stopInstance } from '$lib/instances';
  import { submitFlag } from '$lib/challenges'; // if submitFlag is in challenges.ts; otherwise adjust import

  // 🔐 Global auth store replaces { user } from parent layout
  import { user as authUser } from '$lib/stores/auth';
  import { onMount } from 'svelte';

  // ──────────────────────────────────────────────────────────────────────────────
  // Local state (Svelte 5 runes)
  // ──────────────────────────────────────────────────────────────────────────────
  let loading = $state(true);
  let error   = $state<string | null>(null);

  let challenges = $state<any[]>([]);
  let selected: any | null = $state(null);

  // seconds remaining per challenge id
  let countdowns: Record<string, number> = $state({});

  let openModal = $state(false);
  let solvesOpen = $state(false);
  let flag = $state('');
  let submittingFlag = $state(false);
  let flagError = $state(false);

  // ──────────────────────────────────────────────────────────────────────────────
  // Data loading (replaces +page.ts load)
  // ──────────────────────────────────────────────────────────────────────────────
  async function loadChallenges() {
    loading = true; error = null;
    try {
      challenges = await getChallenges();
      // seed countdowns from initial data
      const next: Record<string, number> = {};
      for (const c of challenges ?? []) {
        if (typeof (c as any).timeout === 'number' && (c as any).timeout > 0) {
          next[(c as any).id] = (c as any).timeout;
        }
      }
      countdowns = next;
    } catch (e: any) {
      error = e?.message ?? 'Failed to load challenges';
    } finally {
      loading = false;
    }
  }

  // run once when component mounts
  onMount(() => {
    loadChallenges();
    const timer = setInterval(() => {
      for (const id in countdowns) {
        if (countdowns[id] > 0) countdowns[id] = countdowns[id] - 1;
      }
    }, 1000);
    return () => clearInterval(timer);
  });


  $effect(() => {
    if (typeof window === 'undefined') return;
    const timer = setInterval(() => {
      for (const id in countdowns) {
        if (countdowns[id] > 0) countdowns[id] = countdowns[id] - 1;
      }
    }, 1000);
    return () => clearInterval(timer);
  });

  // ──────────────────────────────────────────────────────────────────────────────
  // Helpers
  // ──────────────────────────────────────────────────────────────────────────────
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
  const grouped = $derived(groupByCategory(challenges));

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
  function closeModal() { openModal = false; }
  $effect(() => { if (!openModal) selected = null; });

  function modifyChallenge(ch: any) {
    return () => {
      // navigate to modify page (implement your router navigation)
      // e.g., window.location.hash = `#/challenges/${ch.id}/edit`
      alert(1);
    };
  }

  function copyToClipboard(text: string) {
    if (typeof navigator === 'undefined') return;
    navigator.clipboard
      .writeText(text)
      .then(() => toast.success('Copied to clipboard!'))
      .catch(() => toast.error('Failed to copy to clipboard.'));
  }

  // ──────────────────────────────────────────────────────────────────────────────
  // Instance controls (no fetch arg)
  // ──────────────────────────────────────────────────────────────────────────────
  async function createInstance(ch: any) {
    try {
      const { host, port, timeout } = await startInstance(ch.id);
      ch.remote = host;
      ch.port = port;
      ch.timeout = timeout;
      if (typeof ch.timeout === 'number') {
        countdowns[ch.id] = Math.max(0, ch.timeout);
      }
      toast.success('Created instance!');
    } catch (err: any) {
      console.error(err);
      toast.error(`Failed to create instance: ${err?.message ?? err}`);
    }
  }

  async function destroyInstance(ch: any) {
    try {
      await stopInstance(ch.id);
      ch.remote = null;
      ch.port = null;
      ch.timeout = null;
      countdowns[ch.id] = 0;
      toast.success('Stopped instance!');
    } catch (err: any) {
      console.error(err);
      toast.error(`Failed to stop instance: ${err?.message ?? err}`);
    }
  }

  // ──────────────────────────────────────────────────────────────────────────────
  // Flag submission (no fetch arg)
  // ──────────────────────────────────────────────────────────────────────────────
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
      // mark solved locally
      selected.solved = true;
      const idx = challenges.findIndex((c: any) => c.id === selected!.id);
      if (idx !== -1) challenges[idx] = { ...challenges[idx], solved: true };
    } catch (e: any) {
      toast.error(e?.message ?? 'Flag submission failed');
    } finally {
      submittingFlag = false;
    }
  }
</script>

<p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Challenges</p>
<hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
<p class="mb-10 text-lg italic text-gray-500 dark:text-gray-400">
  "A man who loves to walk will walk more than a man who loves his destination"
</p>

{#if $authUser?.role === 'Admin'}
  <div class="fixed right-15 bottom-35 z-50">
    <SpeedDialTrigger color="green" />
    <SpeedDial>
      <SpeedDialButton name="Add a Challenge"><FlagSolid /></SpeedDialButton>
    </SpeedDial>
  </div>
{/if}

{#if loading}
  <div class="p-4">Loading challenges…</div>
{:else if error}
  <div class="p-4 text-red-600">{error}</div>
{:else}
  {#each grouped as [category, items]}
    <section class="mb-10">
      <div class="mb-3 flex items-center gap-3">
        <p class="text-2xl leading-tight font-bold text-gray-900 dark:text-white">{category}</p>
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
            <!-- First row, name with tags-->
            <div class="p-4">
              <p class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">{ch.name}</p>
              {#each ch.tags as tag}
                <Badge class="mr-1" color="gray">{tag}</Badge>
              {/each}
            </div>

            <!-- Bottom row, points, solved and dockerized-->
            <div class="mt-auto flex">
              {#if ch.solved}
                <Badge color="green" class="mr-auto">{ch.points}</Badge>
                <CheckCircleSolid class="mr-2 mb-2 text-green-500" />
              {:else}
                <Badge color="secondary" class="mb-1 ml-1">{ch.points}</Badge>
              {/if}
              {#if ch.instance}
                <Container class="mr-2 mb-2" />
                {#if countdowns[ch.id] > 0}
                  <Badge color="blue">{fmtTimeLeft(countdowns[ch.id])}</Badge>
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
        <BugSolid class="ml-auto h-6 w-6 text-gray-800" />
        {#if $authUser?.role === 'Admin'}
          <Button
            onclick={modifyChallenge(selected)}
            aria-label="Modify this challenge"
            variant="outline"
            size="sm"
            class="ml-auto hover:cursor-pointer"
          >
            <PenSolid />
          </Button>
          <Button
            onclick={modifyChallenge(selected)}
            variant="destructive"
            size="sm"
            class="hover:cursor-pointer mr-5"
          >
            <TrashBinSolid />
          </Button>
        {:else}
          <div class="ml-auto"></div>
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
            onclick={() => solvesOpen = true}
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

    <!-- Instance / remote -->
    <div class="mt-1 flex w-full flex-row items-center justify-center px-6">
      {#if selected?.instance}
        {#if countdowns[selected?.id] > 0}
          <Button
            size="sm"
            style="background-color:#779ecb;"
            disabled
            class="hover:cursor-pointer w-full mr-2"
          >
            <Container class="mr-1" />
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
      {#if selected?.remote}
        <Badge
          color="gray"
          class="cursor-pointer"
          onclick={() =>
            copyToClipboard(`${selected?.remote}${selected?.port ? `:${selected?.port}` : ''}`)
          }
        >
          <p class="text-lg">{selected?.remote}{selected?.port ? ` ${selected?.port}` : ''}</p>
        </Badge>
      {/if}
    </div>

    <!-- Submit flag -->
    <div class="mt-4 flex w-full items-center justify-between">
      <form class="mt-4 flex w-full items-center gap-2" class:justify-center={selected?.solved} onsubmit={onSubmitFlag}>
        {#if !selected?.solved}
            <Input
            class="ps-9 flex-1"
            placeholder="TRX{'...'}"
            bind:value={flag}
            color={flagError ? 'red' : 'gray'}
            oninput={() => (flagError = false)}
            >
            {#snippet left()}
                {#if flagError}
                <ExclamationCircleSolid class="h-5 w-5" />
                {:else}
                <FlagSolid class="h-5 w-5" />
                {/if}
            {/snippet}
            </Input>
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

<!-- All sheets that are imported -->
<SolveListSheet bind:open={solvesOpen} challenge={selected}/>
