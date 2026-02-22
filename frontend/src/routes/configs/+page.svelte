<script lang="ts">
	import { getConfigs, updateConfigs } from '@/config';
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { authState } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import LoadingState from '$lib/components/ui/loading-state.svelte';
	import ConfigCard from '$lib/components/configs/config-card.svelte';
	import { showError } from '$lib/utils/toast';
	import { Settings } from '@lucide/svelte';

	type ConfigType = 'bool' | 'int' | 'string' | (string & {});

	interface Config {
		key: string;
		type?: ConfigType | null;
		value?: string | number | boolean | null;
		description?: string;
		// Allow extra backend fields
		[key: string]: unknown;
	}

	type FormValue = string | boolean;

	let loading = $state(true);
	let error = $state<string | null>(null);

	let configs = $state<Config[]>([]);
	let form = $state<Record<string, FormValue>>({});

	let saving = $state(false);
	let saveError = $state<string | null>(null);
	let saveOk = $state(false);

	const isAdmin = $derived(authState.user?.role === 'Admin');

	const hasChanges = $derived(
		configs.some((config) => {
			const key = config.key;
			const current = form[key];
			const next = toConfigValue(config, current);
			const prev = String(config.value ?? '');
			return next !== prev;
		})
	);

	function normalizeType(type: Config['type']): ConfigType {
		if (type === 'bool' || type === 'int' || type === 'string') return type;
		return 'string';
	}

	function toFormValue(config: Config): FormValue {
		const t = normalizeType(config.type);
		const raw = config.value;

		if (t === 'bool') {
			return String(raw) === 'true';
		}

		return raw != null ? String(raw) : '';
	}

	function toConfigValue(config: Config, formValue: FormValue | undefined): string {
		const t = normalizeType(config.type);

		if (t === 'bool') {
			return formValue ? 'true' : 'false';
		}

		if (t === 'int') {
			const n = Number(formValue ?? 0);
			return Number.isFinite(n) ? String(n) : '0';
		}

		return String(formValue ?? '');
	}

	async function loadConfigs() {
		loading = true;
		error = null;

		try {
			const res = await getConfigs();
			const list = Array.isArray(res) ? (res as Config[]) : [];

			configs = list;

			const nextForm: Record<string, FormValue> = {};
			for (const c of list) {
				if (!c || typeof c !== 'object') continue;
				nextForm[c.key] = toFormValue(c);
			}
			form = nextForm;
		} catch (e: any) {
			error = e?.message ?? 'Failed to load configs';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (authState.user && !isAdmin) {
			goto('/404');
		}
	});

	onMount(loadConfigs);

	async function save() {
		if (saving || !hasChanges) return;

		saving = true;
		saveError = null;
		saveOk = false;

		try {
			// Build list of changed configs with new values
			const changes: Config[] = [];

			for (const config of configs) {
				const key = config.key;
				const value = toConfigValue(config, form[key]);
				const prev = String(config.value ?? '');

				if (value !== prev) {
					changes.push({ ...config, value });
				}
			}

			if (!changes.length) {
				saveOk = true;
				return;
			}

			const changesByKey = new Map(changes.map((c) => [c.key, c]));

			// Update all changed configs (in parallel)
			await Promise.all(changes.map((change) => updateConfigs(change)));

			// Sync local state with saved values
			configs = configs.map((config) => changesByKey.get(config.key) ?? config);

			saveOk = true;
		} catch (e: any) {
			saveError = e?.message ?? 'Failed to save configs';
		} finally {
			saving = false;
		}
	}
</script>

<div
	class="from-muted/20 to-background mb-6 mt-6 rounded-xl border-0 bg-gradient-to-br p-6 shadow-sm"
>
	<div class="flex items-center gap-4">
		<div
			class="bg-background flex h-16 w-16 shrink-0 items-center justify-center rounded-full shadow-sm"
		>
			<Settings class="text-muted-foreground h-8 w-8" />
		</div>
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Configuration</h1>
			<p class="text-muted-foreground mt-2 text-sm">Manage global instance settings.</p>
		</div>
	</div>
</div>

{#if loading}
	<LoadingState message="Loading configuration..." />
{:else if error}
	<div
		class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-600 dark:border-red-800 dark:bg-red-950/20"
	>
		<p class="font-semibold">Error loading configuration</p>
		<p class="text-sm">{error}</p>
	</div>
{:else}
	<div class="relative">
		<div class="space-y-6 pb-24">
			<!-- Config cards -->
			<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
				{#each configs as c (c.key)}
					<div
						class="bg-muted/20 hover:bg-background rounded-lg border-0 p-4 shadow-none transition-all hover:shadow-md"
					>
						<div class="mb-3 flex items-start justify-between gap-3">
							<div class="min-w-0 flex-1">
								<Label class="mb-1 block font-semibold text-gray-900 dark:text-white">
									{c.key}
								</Label>
								{#if c.description}
									<p class="line-clamp-2 text-xs text-gray-600 dark:text-gray-400">
										{c.description}
									</p>
								{/if}
							</div>
							<span
								class="shrink-0 rounded-md bg-gray-100 px-2 py-1 text-xs font-medium text-gray-600 dark:bg-gray-700 dark:text-gray-300"
							>
								{c.type}
							</span>
						</div>

						<div class="mt-3">
							{#if c.type === 'bool'}
								<div class="bg-background flex items-center gap-3 rounded-md border p-3">
									<Checkbox
										checked={form[c.key] === true}
										onCheckedChange={(v) => {
											form[c.key] = !!v;
										}}
										id={c.key}
									/>
									<Label for={c.key} class="cursor-pointer text-sm font-medium">
										{form[c.key] ? 'Enabled' : 'Disabled'}
									</Label>
								</div>
							{:else if c.type === 'int'}
								<Input
									type="number"
									class="bg-background w-full"
									bind:value={form[c.key]}
									placeholder="Enter number"
								/>
							{:else}
								<Input
									type="text"
									class="bg-background w-full"
									bind:value={form[c.key]}
									placeholder="Enter value"
								/>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		</div>

		<div class="bg-background sticky bottom-0 border-t p-4 shadow-lg dark:border-gray-700">
			<div class="flex flex-wrap items-center gap-3">
				<Button onclick={save} disabled={saving || !hasChanges} class="cursor-pointer">
					{#if saving}
						<div
							class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"
						></div>
						Saving...
					{:else}
						Save Changes
					{/if}
				</Button>

				{#if saveOk && !hasChanges}
					<div class="flex items-center gap-2 text-green-600 dark:text-green-400">
						<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							></path>
						</svg>
						<span class="font-medium">Changes saved</span>
					</div>
				{/if}

				{#if saveError}
					<div class="flex items-center gap-2 text-red-600 dark:text-red-400">
						<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							></path>
						</svg>
						<span class="font-medium">{saveError}</span>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
