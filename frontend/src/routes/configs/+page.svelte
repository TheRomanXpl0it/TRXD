<script lang="ts">
	// TODO: Refactor this file
	
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
	<div class="flex flex-col items-center justify-center py-12">
		<div class="mb-4 h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-500"></div>
		<p class="text-gray-600 dark:text-gray-400">Loading configs...</p>
	</div>
{:else if error}
	<div class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-600 dark:border-red-800 dark:bg-red-950/20">
		<p class="font-semibold">Error loading configs</p>
		<p class="text-sm">{error}</p>
	</div>
{:else}
	<div class="space-y-6">
		<!-- Config cards -->
		<div class="grid gap-4 md:grid-cols-2">
			{#each configs as c (c.key)}
				<div class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm transition-shadow hover:shadow-md dark:border-gray-700 dark:bg-gray-800">
					<div class="mb-3 flex items-start justify-between">
						<div class="min-w-0 flex-1">
							<Label class="mb-1 block text-base font-semibold text-gray-900 dark:text-white">
								{c.key}
							</Label>
							{#if c.description}
								<p class="text-xs text-gray-600 dark:text-gray-400">{c.description}</p>
							{/if}
						</div>
						<span class="ml-2 shrink-0 rounded bg-gray-100 px-2 py-1 text-xs font-medium text-gray-600 dark:bg-gray-700 dark:text-gray-300">
							{c.type}
						</span>
					</div>

					<div class="mt-3">
						{#if c.type === 'bool'}
							<div class="flex items-center gap-2">
								<Checkbox bind:checked={form[c.key]} id={c.key} />
								<Label for={c.key} class="cursor-pointer text-sm">
									{form[c.key] ? 'Enabled' : 'Disabled'}
								</Label>
							</div>
						{:else if c.type === 'int'}
							<Input
								type="number"
								class="w-full"
								bind:value={form[c.key]}
								placeholder="Enter number"
							/>
						{:else}
							<Input
								type="text"
								class="w-full"
								bind:value={form[c.key]}
								placeholder="Enter value"
							/>
						{/if}
					</div>
				</div>
			{/each}
		</div>

		<!-- Save section -->
		<div class="flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50 p-4 dark:border-gray-700 dark:bg-gray-800/50">
			<Button onclick={save} disabled={saving} class="cursor-pointer px-6">
				{#if saving}
					<div class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
					Saving...
				{:else}
					Save Changes
				{/if}
			</Button>
			{#if saveOk}
				<div class="flex items-center gap-2 text-green-600 dark:text-green-400">
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
					</svg>
					<span class="font-medium">Saved successfully!</span>
				</div>
			{/if}
			{#if saveError}
				<div class="flex items-center gap-2 text-red-600 dark:text-red-400">
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
					<span class="font-medium">{saveError}</span>
				</div>
			{/if}
		</div>
	</div>
{/if}
