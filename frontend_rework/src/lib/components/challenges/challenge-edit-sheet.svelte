<script lang="ts">
  import * as Sheet from "$lib/components/ui/sheet/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Textarea } from "$lib/components/ui/textarea/index.js";
  import { Checkbox } from "$lib/components/ui/checkbox/index.js";
  import { Slider } from "$lib/components/ui/slider/index.js";
  import Label from "$lib/components/ui/label/label.svelte";
  import CategorySelect from "$lib/components/challenges/category-select.svelte";
  import { toast } from "svelte-sonner";
  import * as Accordion from "$lib/components/ui/accordion/index.js";
  import { updateChallenge, getChallenge } from "$lib/challenges";
  import { Check, Cpu, MemoryStick, Clock } from "@lucide/svelte";

  // --- CodeMirror imports (YAML + theming) ---
  import { EditorState, Compartment } from "@codemirror/state";
  import {
    EditorView,
    keymap,
    lineNumbers,
    highlightActiveLine,
    highlightActiveLineGutter
  } from "@codemirror/view";
  import {
    defaultKeymap,
    history,
    historyKeymap,
    indentWithTab
  } from "@codemirror/commands";
  import { yaml } from "@codemirror/lang-yaml";
  import { oneDark } from "@codemirror/theme-one-dark";
  import { syntaxHighlighting, defaultHighlightStyle } from "@codemirror/language";

  type Item = { value: string; label: string };

  // Props
  let {
    open = $bindable(false),
    challenge_user,
    categories = [] as Item[]
  } = $props<{ open?: boolean; challenge_user: any; categories?: Item[] }>();

  let challenge = $state<any>(null);

  // Form state
  let name           = $state("");
  let category       = $state("");
  let description    = $state("");
  let difficulty     = $state<"easy"|"medium"|"hard"|"insane">("easy");
  let type           = $state("Container");
  let hidden         = $state<boolean>(false);
  let maxPoints      = $state<number>(500);
  let dynamicScoring = $state(true);
  let host           = $state("");
  let portStr        = $state("");
  let attachments    = $state("");
  let authorsCsv     = $state("");
  let hashDomain     = $state(false);
  let imageName      = $state("");
  let composeFile    = $state("");
  let maxCPU         = $state("");
  let maxRam         = $state("");
  let lifetime       = $state("");
  let envs           = $state("");

  let saving = $state(false);

  // Options
  const typeOptions: Item[] = [
    { value: "Container", label: "Container" },
    { value: "Compose",   label: "Compose" },
    { value: "Normal",    label: "Normal" }
  ];

  const difficultyOptions: Item[] = [
    { value: "easy",   label: "Easy" },
    { value: "medium", label: "Medium" },
    { value: "hard",   label: "Hard" },
    { value: "insane", label: "Insane" }
  ];

  function toTitleCase(v: string) {
    const s = String(v || "").toLowerCase();
    return s.charAt(0).toUpperCase() + s.slice(1);
  }

  function normalizeDifficulty(input: any): "easy"|"medium"|"hard"|"insane" {
    const s = String(input ?? "").toLowerCase();
    if (s === "easy" || s === "medium" || s === "hard" || s === "insane") return s;
    return "easy";
  }

  function parseCsv(s: string): string[] {
    return s.split(",").map(x => x.trim()).filter(Boolean);
  }

  $effect(() => {
    if (!open) return;
    let cancelled = false;

    (async () => {
      try {
        const fetched = await getChallenge(challenge_user?.id);
        challenge = { ...(fetched ?? {}), ...(challenge_user ?? {}) };
        if (cancelled || !challenge) return;

        name        = String(challenge?.name ?? "");
        const rawCat = challenge?.category;
        category    = rawCat
          ? (typeof rawCat === "string" ? rawCat : (rawCat?.name ?? "")).toString()
          : "";

        description = String(challenge?.description ?? "");
        difficulty  = normalizeDifficulty(challenge?.difficulty);

        if (Array.isArray(challenge?.authors)) {
          authorsCsv = challenge.authors.join(", ");
        } else {
          authorsCsv = String(challenge?.authors ?? "");
        }

        type        = String(challenge?.type ?? "Container");
        hidden      = Boolean(challenge?.hidden ?? false);
        maxPoints   = Number.isFinite(+challenge?.max_points)
          ? Number(challenge.max_points)
          : Number.isFinite(+challenge?.points)
          ? Number(challenge.points)
          : 500;

        const st = String(challenge?.score_type ?? challenge?.scoreType ?? "STATIC").toUpperCase();
        dynamicScoring = st === "DYNAMIC";

        host        = String(challenge?.host ?? "");
        portStr     = challenge?.port != null ? String(challenge.port) : "";
        attachments = Array.isArray(challenge?.attachments)
          ? challenge.attachments.join(", ")
          : String(challenge?.attachments ?? "");

        hashDomain  = Boolean(challenge?.hash_domain ?? false);
        imageName   = String(challenge?.docker_config?.image ?? "");
        composeFile = String(challenge?.docker_config?.compose ?? "");
        maxCPU      = String(challenge?.docker_config?.max_cpu ?? "");
        maxRam      = String(challenge?.docker_config?.max_memory ?? "");
        lifetime    = String(challenge?.docker_config?.lifetime ?? "");
        envs        = String(challenge?.docker_config?.envs ?? "");
      } catch {
        /* noop */
      }
    })();

    return () => { cancelled = true; };
  });

  async function onSave(e: Event) {
    e.preventDefault();
    if (saving || !challenge?.id) return;

    if (!category.trim()) {
      toast.error("Please select a category.");
      return;
    }

    const portNum = portStr.trim() ? Number(portStr) : undefined;
    if (portStr.trim() && !Number.isFinite(portNum)) {
      toast.error("Port must be a number.");
      return;
    }

    const payload: any = {
      ChallID:     challenge.id,
      Name:        name.trim(),
      Category:    category.trim(),
      Description: description.trim(),
      Difficulty:  toTitleCase(difficulty),
      Authors:     parseCsv(authorsCsv),
      Type:        type,
      Hidden:      hidden,
      MaxPoints:   maxPoints,
      ScoreType:   dynamicScoring ? "DYNAMIC" : "STATIC",
      Host:        host.trim(),
      Port:        portNum,
      Attachments: parseCsv(attachments)
      // Add docker config fields if needed
    };

    saving = true;
    try {
      await updateChallenge(payload);
      toast.success("Challenge updated.");
      open = false;
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to update challenge.");
    } finally {
      saving = false;
    }
  }

  // ---------- CodeMirror setup for the Compose YAML field ----------
  let cmComposeHost: HTMLDivElement | undefined;
  let cmComposeView: EditorView | null = null;

  // Theme compartment + light theme that plays nice with your UI
  const themeComp = new Compartment();
  function lightTheme() {
    return EditorView.theme(
      {
        "&": { backgroundColor: "transparent" },
        ".cm-content": {
          fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
          fontSize: "0.9rem"
        },
        ".cm-gutters": {
          backgroundColor: "transparent",
          borderRight: "none",
          /* slate-500-ish for light mode line numbers */
          color: "rgb(100 116 139)"
        },
        "&.cm-editor": { borderRadius: "0.5rem" },
        "&.cm-focused": { outline: "none" }
      },
      { dark: false }
    );
  }

  // Create/destroy editor depending on section visibility
  $effect(() => {
    if (type === "Compose" && cmComposeHost && !cmComposeView) {
      cmComposeView = new EditorView({
        parent: cmComposeHost,
        state: EditorState.create({
          doc: composeFile || `version: '3'
services:
  web:
    image: TRX-Chall-1
    ports:
      - "31337:31337"
`,
          extensions: [
            lineNumbers(),
            highlightActiveLineGutter(),
            highlightActiveLine(),
            history(),
            keymap.of([indentWithTab, ...defaultKeymap, ...historyKeymap]),
            yaml(),
            /* crucial for colors when not using basicSetup */
            syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
            themeComp.of(lightTheme()), // start light; we swap below if needed
            EditorView.theme({ "&": { minHeight: "180px", borderRadius: "0.5rem" } }),
            EditorView.updateListener.of((u) => {
              if (u.docChanged) composeFile = u.state.doc.toString();
            })
          ]
        })
      });
    }

    if (type !== "Compose" && cmComposeView) {
      cmComposeView.destroy();
      cmComposeView = null;
    }
  });

  // Keep editor in sync with external composeFile changes
  $effect(() => {
    if (!cmComposeView) return;
    const current = cmComposeView.state.doc.toString();
    if (composeFile !== current) {
      cmComposeView.dispatch({
        changes: { from: 0, to: current.length, insert: composeFile }
      });
    }
  });

  // Auto-switch theme (Tailwind .dark or system)
  $effect(() => {
    if (!cmComposeView) return;

    const isDarkNow = () =>
      typeof document !== "undefined" &&
      (document.documentElement.classList.contains("dark") ||
        (typeof window !== "undefined" &&
          window.matchMedia?.("(prefers-color-scheme: dark)")?.matches));

    const apply = () => {
      cmComposeView?.dispatch({
        effects: themeComp.reconfigure(isDarkNow() ? oneDark : lightTheme())
      });
    };

    // initial apply
    apply();

    // listen to system changes
    const mq = window.matchMedia?.("(prefers-color-scheme: dark)");
    const mqHandler = () => apply();
    mq?.addEventListener?.("change", mqHandler);

    // listen to Tailwind's .dark class toggles
    const obs = new MutationObserver(apply);
    obs.observe(document.documentElement, { attributes: true, attributeFilter: ["class"] });

    return () => {
      mq?.removeEventListener?.("change", mqHandler);
      obs.disconnect();
    };
  });
</script>

<Sheet.Root bind:open>
  <Sheet.Content side="right" class="sm:max-w-[720px] px-5">
    <Sheet.Header>
      <Sheet.Title>Edit Challenge</Sheet.Title>
      <Sheet.Description>
        Modify settings for <b>{challenge?.name ?? "—"}</b>.
      </Sheet.Description>
    </Sheet.Header>

    <form class="space-y-5 mt-3" onsubmit={onSave}>
      <Accordion.Root type="single">
        <Accordion.Item value="Naming">
          <Accordion.Trigger class="text-xl font-semibold tracking-tight cursor-pointer">Naming and description</Accordion.Trigger>
          <Accordion.Content>
            <div>
              <Label for="ch-name" class="mb-1 block">Name</Label>
              <Input id="ch-name" bind:value={name} />
            </div>
            <div class="mt-3">
              <Label for="ch-desc" class="mb-1 block">Description</Label>
              <Textarea id="ch-desc" bind:value={description} class="min-h-28" />
            </div>
          </Accordion.Content>
        </Accordion.Item>

        <Accordion.Item value="Scoring">
          <Accordion.Trigger class="text-xl font-semibold tracking-tight cursor-pointer">Scoring and visibility</Accordion.Trigger>
          <Accordion.Content>
            <div class="flex flex-row justify-between mt-3">
              <div class="flex items-center gap-3">
                <Checkbox id="ch-hidden" bind:checked={hidden} />
                <Label for="ch-hidden">Hidden</Label>
              </div>
              <div class="flex items-center gap-3">
                <Checkbox id="ch-score" bind:checked={dynamicScoring} />
                <Label for="ch-score">Dynamic scoring</Label>
              </div>
            </div>
            <div class="mt-5">
              <Label for="ch-points" class="mb-1 block">Max Points: {maxPoints}</Label>
              <Slider id="ch-points" type="single" bind:value={maxPoints} min={0} max={1500} step={25} />
            </div>
          </Accordion.Content>
        </Accordion.Item>

        <Accordion.Item value="Difficulty">
          <Accordion.Trigger class="text-xl font-semibold tracking-tight cursor-pointer">Difficulty and categorization</Accordion.Trigger>
          <Accordion.Content>
            <div class="flex flex-row justify-between">
              <div>
                <Label for="ch-diff" class="mb-1 block">Difficulty</Label>
                <CategorySelect id="ch-diff" items={difficultyOptions} bind:value={difficulty} placeholder="Select difficulty…" />
              </div>
              <div>
                <Label for="ch-type" class="mb-1 block">Type</Label>
                <CategorySelect id="ch-type" items={typeOptions} bind:value={type} placeholder="Select type…" />
              </div>
              <div>
                <Label for="ch-cat" class="mb-1 block">Category</Label>
                <CategorySelect id="ch-cat" items={categories} bind:value={category} placeholder="Select a category…" />
              </div>
            </div>
          </Accordion.Content>
        </Accordion.Item>

        <Accordion.Item value="Authors">
          <Accordion.Trigger class="text-xl font-semibold tracking-tight cursor-pointer">Authors and tags</Accordion.Trigger>
          <Accordion.Content>
            <div class="flex flex-col">
              <div>
                <Label for="ch-auth" class="mb-1 block">Authors</Label>
                <Input id="ch-auth" bind:value={authorsCsv} placeholder="alice, bob" />
              </div>
              <div class="mt-3">
                <Label for="ch-att" class="mb-1 block">Attachments</Label>
                <Input id="ch-att" bind:value={attachments} placeholder="url1, url2" />
              </div>
            </div>
          </Accordion.Content>
        </Accordion.Item>

        {#if type === 'Normal'}
          <Accordion.Item value="Host">
            <Accordion.Trigger class="text-xl font-semibold tracking-tight cursor-pointer">Host and port</Accordion.Trigger>
            <Accordion.Content>
              <div class="flex flex-col">
                <div>
                  <Label for="ch-host" class="mb-1 block">Host</Label>
                  <Input id="ch-host" bind:value={host} placeholder="e.g. challenge.trxd.cc" />
                </div>
                <div>
                  <Label for="ch-port" class="mb-1 block">Port</Label>
                  <Input id="ch-port" bind:value={portStr} placeholder="e.g. 31337" />
                </div>
                <div class="mt-3 flex flex-row items-center">
                  <Checkbox id="ch-hashdomain" bind:checked={hashDomain} />
                  <Label for="ch-hashdomain" class="ml-2">Hash domain</Label>
                </div>
              </div>
            </Accordion.Content>
          </Accordion.Item>

        {:else if type === 'Compose'}
          <Accordion.Item value="Compose">
            <Accordion.Trigger class="text-xl font-semibold tracking-tight cursor-pointer">Compose Settings</Accordion.Trigger>
            <Accordion.Content>
              <div class="flex flex-col">
                <Label>Compose.yaml</Label>
                <div bind:this={cmComposeHost} class="mt-3 border rounded min-h-45" />
                <div class="mt-3 flex flex-row items-center">
                  <Checkbox id="com-hashdomain" bind:checked={hashDomain} />
                  <Label for="com-hashdomain" class="ml-1">Hash Domain</Label>
                </div>
              </div>
            </Accordion.Content>
          </Accordion.Item>

        {:else}
          <Accordion.Item value="Container">
            <Accordion.Trigger class="text-xl font-semibold tracking-tight cursor-pointer">Container Settings</Accordion.Trigger>
            <Accordion.Content>
              <div class="flex flex-col">
                <div>
                  <Label for="ho-host" class="mb-1 block">Container Image name</Label>
                  <Input id="ho-host" bind:value={imageName} placeholder="TRX-Chall-1" />
                </div>
                <div class="flex flex-row justify-between items-center">
                  <div class="flex flex-row mt-3 items-center">
                    <Checkbox id="ho-hashdomain" bind:checked={hashDomain} />
                    <Label for="ho-hashdomain" class="ml-2">Hash domain</Label>
                  </div>
                  <div class="mt-3">
                    <Label for="ho-port">Port</Label>
                    <Input id="ho-port" bind:value={portStr} placeholder="e.g. 31337" disabled={hashDomain} />
                  </div>
                </div>
              </div>
            </Accordion.Content>
          </Accordion.Item>
        {/if}

        {#if type === 'Container' || type === 'Compose'}
          <Accordion.Item value="Performance">
            <Accordion.Trigger class="text-xl font-semibold tracking-tight cursor-pointer">Performance Settings</Accordion.Trigger>
            <Accordion.Content>
              <div class="flex flex-col">
                <div class="flex flex-row justify-between">
                  <div class="mt-3 flex flex-row items-center">
                    <Cpu class="mr-2" />
                    <div>
                      <Label for="perf-cpu" class="mb-1 block">Max CPU</Label>
                      <Input id="perf-cpu" bind:value={maxCPU} placeholder="50%" />
                    </div>
                  </div>
                  <div class="mt-3 flex flex-row items-center">
                    <MemoryStick class="mr-2" />
                    <div>
                      <Label for="perf-ram" class="mb-1 block">Max RAM</Label>
                      <Input id="perf-ram" bind:value={maxRam} placeholder="1GB" />
                    </div>
                  </div>
                </div>
                <div class="flex flex-row justify-between">
                  <div class="mt-3 flex flex-row items-center">
                    <Clock class="mr-2" />
                    <div>
                      <Label for="perf-lifetime" class="mb-1 block">Lifetime</Label>
                      <Input id="perf-lifetime" bind:value={lifetime} placeholder="300" />
                    </div>
                  </div>
                </div>
              </div>
            </Accordion.Content>
          </Accordion.Item>
        {/if}
      </Accordion.Root>

      <div class="flex justify-end gap-2">
        <Sheet.Close>
          <Button type="button" variant="outline">Cancel</Button>
        </Sheet.Close>
        <Button type="submit" disabled={saving}>
          {#if saving}Saving…{:else}Save{/if}
        </Button>
      </div>
    </form>
  </Sheet.Content>
</Sheet.Root>
