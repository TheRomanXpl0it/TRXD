<script lang="ts">
	import { getConfigs, updateConfigs } from '@/config';
	import { onMount } from 'svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Button } from '$lib/components/ui/button';
	import { user as authUser } from '$lib/stores/auth';
	import { push } from 'svelte-spa-router';

	let loading = $state(true);
	let error = $state<string | null>(null);
	let configs = $state<any[]>([]);
	const isAdmin = $derived(($authUser as any)?.role === 'Admin');

	// Editable form state keyed by config key
	let form = $state<Record<string, any>>({});

	// Save state
	let saving = $state(false);
	let saveError = $state<string | null>(null);
	let saveOk = $state(false);

	$effect(() => {
		if (!isAdmin) {
			push('/404');
		}
	});

	onMount(async () => {
		loading = true;
		error = null;
		try {
			const res = await getConfigs();
			configs = Array.isArray(res) ? res : [];
			// Initialize form values by type
			for (const c of configs) {
				if (!c || typeof c !== 'object') continue;
				const k = c.key as string;
				const t = (c.type as string) ?? 'string';
				const v = c.value;
				if (t === 'bool') {
					form[k] = String(v) === 'true';
				} else if (t === 'int') {
					form[k] = String(v ?? '');
				} else {
					form[k] = String(v ?? '');
				}
			}
		} catch (e: any) {
			error = e?.message ?? 'Failed to load configs';
		} finally {
			loading = false;
		}
	});

	async function save() {
		if (saving) return;
		saving = true;
		saveError = null;
		saveOk = false;
		try {
			// Build changes array by converting current form values back to strings
			const changes = [];
			for (const c of configs) {
				const t = (c.type as string) ?? 'string';
				const k = c.key as string;
				let value: string;
				if (t === 'bool') {
					value = form[k] ? 'true' : 'false';
				} else if (t === 'int') {
					const n = Number(form[k] ?? 0);
					value = String(Number.isFinite(n) ? n : 0);
				} else {
					value = String(form[k] ?? '');
				}
				const prev = String(c.value ?? '');
				if (value !== prev) {
					changes.push({ ...c, value });
				}
			}

			// Call updateConfigs once per changed config
			for (const change of changes) {
				await updateConfigs(change);
			}

			saveOk = true;
		} catch (e: any) {
			saveError = e?.message ?? 'Failed to save configs';
		} finally {
			saving = false;
		}
	}
</script>

<div>
	<p class="mt-5 text-3xl font-bold text-gray-800 dark:text-gray-100">Configs</p>
	<hr class="my-2 h-px border-0 bg-gray-200 dark:bg-gray-700" />
	<p class="mb-6 text-lg italic text-gray-500 dark:text-gray-400">
		"Tiny tweaks can lead to big changes"
	</p>
</div>

{#if loading}
	<div>Loading configs…</div>
{:else if error}
	<div class="text-red-600">{error}</div>
{:else}
	<div class="space-y-4">
		<!-- Editable list - tighter layout using shadcn components -->
		<div class="space-y-3">
			{#each configs as c (c.key)}
				<div class="flex items-center justify-between gap-3 rounded border p-3">
					<div class="min-w-0">
						<Label class="mb-0.5 block">{c.key}</Label>
						{#if c.description}
							<p class="text-muted-foreground text-xs">{c.description}</p>
						{/if}
					</div>
					<div class="flex items-center gap-3">
						<small class="text-gray-500">{c.type}</small>
						{#if c.type === 'bool'}
							<Checkbox bind:checked={form[c.key]} />
						{:else if c.type === 'int'}
							<Input type="number" class="w-32" bind:value={form[c.key]} />
						{:else}
							<Input type="text" class="w-64" bind:value={form[c.key]} />
						{/if}
					</div>
				</div>
			{/each}
		</div>

		<div class="flex items-center gap-3">
			<Button onclick={save} disabled={saving} class="cursor-pointer">
				{#if saving}Saving…{:else}Save{/if}
			</Button>
			{#if saveOk}
				<span class="text-green-600">Saved.</span>
			{/if}
			{#if saveError}
				<span class="text-red-600">{saveError}</span>
			{/if}
		</div>
	</div>
{/if}
