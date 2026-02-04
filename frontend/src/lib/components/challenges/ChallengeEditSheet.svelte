<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Slider } from '$lib/components/ui/slider/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import CategorySelect from '$lib/components/challenges/CategorySelect.svelte';
	import TagMultiSelect from '$lib/components/challenges/TagMultiselect.svelte';
	import { toast } from 'svelte-sonner';
	import {
		updateChallenge,
		getChallenge,
		uploadAttachments,
		deleteAttachments
	} from '$lib/challenges';
	import { Cpu, MemoryStick, Clock, X, Tags as TagsIcon } from '@lucide/svelte';
	import { createFlags, deleteFlags } from '$lib/flags';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';

	type Item = { value: string; label: string };

	// Props
	let {
		open = $bindable(false),
		challenge_user,
		all_tags,
		onupdated
	} = $props<{
		open?: boolean;
		challenge_user: any;
		all_tags?: string[];
		onupdated?: (detail: { id: number }) => void;
	}>();

	let challenge = $state<any>(null);
	let activeTab = $state<'meta' | 'settings' | 'flags' | 'deployment'>('meta');

	// Form state
	let name = $state('');
	let category = $state('');
	let description = $state('');
	let difficulty = $state<'easy' | 'medium' | 'hard' | 'insane'>('easy');
	let type = $state('Container');
	let hidden = $state<boolean>(false);
	let maxPoints = $state<number>(500);
	let dynamicScoring = $state(true);
	let host = $state('');
	let portStr = $state('');
	let connType = $state('NONE');
	let authorsCsv = $state('');
	let hashDomain = $state(false);
	let imageName = $state('');
	let composeFile = $state('');
	let maxCPU = $state('');
	let maxRam = $state('');
	let lifetime = $state('');
	let envs = $state('');

	let flags_og = $state<any[]>([]);
	let flags: any = $state<any[]>([]);
	let flagsRegex: any = $state<any[]>([]);

	const uniqAllTags = $derived(
		Array.from(
			new Set(
				(Array.isArray(all_tags) ? all_tags : []).map((t) => String(t ?? '').trim()).filter(Boolean)
			)
		)
	);

	// Current selected tags, as **strings**, deduped
	let tags = $state<string[]>(
		Array.from(new Set((Array.isArray(all_tags) ? all_tags : []).map(String)))
	);

	let existingAttachments = $state<string[]>([]); // from backend
	let removedAttachmentNames = $state<Set<string>>(new Set()); // user-marked for deletion
	let newFiles = $state<File[]>([]); // newly selected files
	let fileInputEl = $state<HTMLInputElement | null>(null); // hidden <input type="file">

	function addFiles(files: FileList | null) {
		if (!files || files.length === 0) return;
		const incoming = Array.from(files);
		// de-dupe by name + size (simple heuristic)
		const dedup = new Map(newFiles.map((f) => [f.name + '::' + f.size, f]));
		for (const f of incoming) {
			const k = f.name + '::' + f.size;
			if (!dedup.has(k)) dedup.set(k, f);
		}
		newFiles = Array.from(dedup.values());
		// reset input so same file can be re-added if removed
		if (fileInputEl) fileInputEl.value = '';
	}

	function removeNewFile(index: number) {
		newFiles = newFiles.filter((_, i) => i !== index);
	}

	function toggleRemoveExisting(name: string) {
		const s = new Set(removedAttachmentNames);
		if (s.has(name)) s.delete(name);
		else s.add(name);
		removedAttachmentNames = s;
	}

	const keptExisting = $derived(existingAttachments.filter((n) => !removedAttachmentNames.has(n)));

	function formatBytes(n: number) {
		if (!Number.isFinite(n)) return '';
		const units = ['B', 'KB', 'MB', 'GB', 'TB'];
		let i = 0;
		while (n >= 1024 && i < units.length - 1) {
			n /= 1024;
			i++;
		}
		return `${n.toFixed(n < 10 && i > 0 ? 1 : 0)} ${units[i]}`;
	}

	let saving = $state(false);

	// Options
	const typeOptions: Item[] = [
		{ value: 'Container', label: 'Container' },
		{ value: 'Compose', label: 'Compose' },
		{ value: 'Normal', label: 'Normal' }
	];

	const difficultyOptions: Item[] = [
		{ value: 'easy', label: 'Easy' },
		{ value: 'medium', label: 'Medium' },
		{ value: 'hard', label: 'Hard' },
		{ value: 'insane', label: 'Insane' }
	];

	function toTitleCase(v: string) {
		const s = String(v || '').toLowerCase();
		return s.charAt(0).toUpperCase() + s.slice(1);
	}

	function normalizeDifficulty(input: any): 'easy' | 'medium' | 'hard' | 'insane' {
		const s = String(input ?? '').toLowerCase();
		if (s === 'easy' || s === 'medium' || s === 'hard' || s === 'insane') return s;
		return 'easy';
	}

	function parseCsv(s: string): string[] {
		return s
			.split(',')
			.map((x) => x.trim())
			.filter(Boolean);
	}

	$effect(() => {
		if (!open) return;
		let cancelled = false;

		(async () => {
			try {
				const fetched = await getChallenge(challenge_user?.id);
				challenge = { ...(fetched ?? {}), ...(challenge_user ?? {}) };
				if (cancelled || !challenge) return;

				flags_og = Array.isArray(challenge?.flags) ? challenge.flags : [];

				name = String(challenge?.name ?? '');
				const rawCat = challenge?.category;
				category = rawCat
					? (typeof rawCat === 'string' ? rawCat : (rawCat?.name ?? '')).toString()
					: '';

				description = String(challenge?.description ?? '');
				difficulty = normalizeDifficulty(challenge?.difficulty);

				if (Array.isArray(challenge?.authors)) {
					authorsCsv = challenge.authors.join(', ');
				} else {
					authorsCsv = String(challenge?.authors ?? '');
				}

				type = String(challenge?.type ?? 'Container');
				hidden = Boolean(challenge?.hidden ?? false);
				maxPoints = Number.isFinite(+challenge?.max_points)
					? Number(challenge.max_points)
					: Number.isFinite(+challenge?.points)
						? Number(challenge.points)
						: 500;

				const st = String(challenge?.score_type ?? challenge?.scoreType ?? 'Static');
				dynamicScoring = st === 'Dynamic';

				host = String(challenge?.host ?? '');
				portStr = challenge?.port != null ? String(challenge.port) : '';
				connType = String(challenge?.conn_type ?? 'NONE');

				existingAttachments = Array.isArray(challenge?.attachments)
					? challenge.attachments.map((a: any) => String(a)).filter(Boolean)
					: [];
				removedAttachmentNames = new Set();
				newFiles = [];

				// 👇 Keep tags as strings; normalize + dedupe
				tags = Array.from(
					new Set(
						(Array.isArray(challenge?.tags) ? challenge.tags : [])
							.map((t: any) => String(t ?? '').trim())
							.filter(Boolean)
					)
				);
				imageName = String(challenge?.docker_config?.image ?? '');
				hashDomain = Boolean(challenge?.docker_config?.hash_domain ?? false);
				composeFile = String(challenge?.docker_config?.compose ?? '');
				maxCPU = String(challenge?.docker_config?.max_cpu ?? '');
				maxRam = String(challenge?.docker_config?.max_memory ?? '');
				lifetime = String(challenge?.docker_config?.lifetime ?? '');
				envs = (() => {
					const e = challenge?.docker_config?.envs;
					if (!e) return '';
					try {
						const obj = JSON.parse(String(e));
						return Object.entries(obj)
							.map(([k, v]) => `${k}=${v}`)
							.join('\n');
					} catch {
						return String(e);
					}
				})();

				// ✅ Properly map flags array (no for..in on arrays)
				flags = Array.isArray(challenge?.flags)
					? challenge.flags.map((f: any) => ({
							flag: String(f?.flag ?? ''),
							regex: Boolean(f?.regex)
						}))
					: [];
			} catch {
				/* noop */
			}
		})();

		return () => {
			cancelled = true;
		};
	});

	async function onSave(e: Event) {
		e.preventDefault();
		if (saving || !challenge_user?.id) return;

		if (!category.trim()) {
			toast.error('Please select a category.');
			return;
		}

		const portNum = portStr.trim() ? Number(portStr) : undefined;
		if (portStr.trim() && !Number.isFinite(portNum)) {
			toast.error('Port must be a number.');
			return;
		}

		// helpers
		const toInt = (x: any) => {
			const n = parseInt(String(x), 10);
			return !isNaN(n) ? n : undefined;
		};
		const toNum = (x: any) => {
			const n = Number(x);
			return Number.isFinite(n) ? n : undefined;
		};
		const str = (x: any) => {
			const s = String(x ?? '').trim();
			return s || undefined;
		};

		// build the fields exactly as backend expects (snake_case)
		const fields: any = {
			chall_id: challenge_user.id,
			name: name.trim(),
			category: category.trim(),
			description: str(description),
			difficulty: toTitleCase(difficulty), // "Easy" | "Medium" | ...
			authors: (authorsCsv || '')
				.split(',')
				.map((a) => a.trim())
				.filter(Boolean),
			type,
			hidden,
			score_type: dynamicScoring ? 'Dynamic' : 'Static',
			host: str(host),
			port: portNum,
			conn_type: connType,

			// container/compose specifics
			image: type === 'Container' ? str(imageName) : undefined,
			compose: type === 'Compose' ? str(composeFile) : undefined,
			hash_domain: hashDomain,

			// performance / limits
			max_points: toInt(maxPoints) ?? 500,
			lifetime: toInt(lifetime) ?? 0,
			max_memory: toInt(maxRam) ?? 0,

			// misc docker
			envs: (() => {
				if (!envs.trim()) return undefined;
				try {
					JSON.parse(envs);
					return envs.trim();
				} catch {
					const obj: Record<string, string> = {};
					envs.split(/[\n,;]/).forEach((line) => {
						const parts = line.split('=');
						if (parts.length >= 2) {
							const key = parts[0].trim();
							const val = parts.slice(1).join('=').trim();
							if (key) obj[key] = val;
						}
					});
					return Object.keys(obj).length > 0 ? JSON.stringify(obj) : undefined;
				}
			})(),
			tags,
			max_cpu: toNum(maxCPU) && toNum(maxCPU) > 0 ? String(toNum(maxCPU)) : ''
		};

		// strip empty arrays
		if (Array.isArray(fields.authors) && fields.authors.length === 0) delete fields.authors;

		saving = true;

		// tags diffs
		// tags diffs (removed separate logic, now handled in updateChallenge)

		// flags diffs
		const prevFlags = Array.from(flags_og ?? []);
		const currFlags = Array.from(flags ?? []);
		const deletedFlags = prevFlags.filter(
			(pf: any) => !currFlags.some((cf: any) => cf.flag === pf.flag && !!cf.regex === !!pf.regex)
		);
		const newFlags = currFlags.filter(
			(cf: any) => !prevFlags.some((pf: any) => cf.flag === pf.flag && !!cf.regex === !!pf.regex)
		);

		try {
			// ---- Build FormData (multipart) ----
			const fd = new FormData();

			// Append chall_id first, explicitly (Svelte 5 value guard)
			fd.append('chall_id', String(challenge_user?.id ?? ''));

			// 1) Append scalar + array fields (EXCLUDING chall_id)
			const entries: Record<string, any> = {
				name: fields.name,
				category: fields.category,
				description: fields.description,
				difficulty: fields.difficulty,
				type: fields.type,
				hidden: fields.hidden ? 'true' : 'false',
				score_type: fields.score_type,
				host: fields.host,
				port: fields.port,
				image: fields.image,
				compose: fields.compose,
				hash_domain: fields.hash_domain ? 'true' : 'false',
				max_points: fields.max_points,
				lifetime: fields.lifetime,
				max_memory: fields.max_memory,
				envs: fields.envs,
				max_cpu: fields.max_cpu
			};

			// Send authors as indexed array for Go form decoder
			if (Array.isArray(fields.authors)) {
				fields.authors.forEach((author, index) => {
					fd.append(`authors[${index}]`, author);
				});
			}

			if (Array.isArray(fields.tags)) {
				fields.tags.forEach((tag: string, index: number) => {
					fd.append(`tags[${index}]`, tag);
				});
			}

			for (const [k, v] of Object.entries(entries)) {
				if (v !== undefined && v !== null && v !== '') {
					fd.append(k, String(v));
				}
			}

			// 4) Update Metadata (JSON)
			await updateChallenge(fields);

			// 5) Upload New Attachments (if any)
			if (newFiles.length > 0) {
				const fd = new FormData();
				fd.append('chall_id', String(challenge_user?.id));
				for (const f of newFiles) {
					// Backend iterates multipartForm.File, so Key 'attachments' is standard
					fd.append('attachments', f, f.name);
				}
				await uploadAttachments(fd);
			}

			// 6) Delete Attachments (if any)
			if (removedAttachmentNames.size > 0) {
				const namesToDelete = Array.from(removedAttachmentNames).map(
					(n) => n.split('/').pop() || n
				);
				await deleteAttachments(challenge_user.id, namesToDelete);
			}

			await createFlags(newFlags, challenge_user?.id);
			await deleteFlags(deletedFlags, challenge_user?.id);

			onupdated?.({ id: challenge_user.id });
			open = false;
			toast.success('Challenge updated.');
		} catch (err: any) {
			toast.error(err?.message ?? 'Failed to update challenge.');
		} finally {
			saving = false;
		}
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="w-full px-5 sm:max-w-[720px]">
		<div
			class="from-muted/20 to-background mb-6 mt-4 rounded-xl border-0 bg-gradient-to-br p-6 shadow-sm"
		>
			<div class="flex items-center gap-4">
				<div
					class="bg-background flex h-16 w-16 shrink-0 items-center justify-center rounded-full shadow-sm"
				>
					<Cpu class="text-muted-foreground h-8 w-8" />
				</div>
				<div>
					<Sheet.Title class="text-xl font-bold">Edit Challenge</Sheet.Title>
					<Sheet.Description class="text-muted-foreground/80 mt-1">
						Modify settings for <b class="text-foreground">{challenge?.name ?? '-'}</b>
					</Sheet.Description>
				</div>
			</div>
		</div>

		<div class="flex flex-col overflow-hidden">
			<form class="mt-3 flex flex-col overflow-hidden" onsubmit={onSave}>
				<!-- Tabs -->
				<div class="bg-background sticky top-0 z-10 border-b border-gray-200 dark:border-gray-700">
					<div class="flex">
						<button
							type="button"
							class="cursor-pointer border-b-2 px-4 py-3 text-sm font-medium transition-colors focus:outline-none {activeTab ===
							'meta'
								? 'border-black text-black dark:border-white dark:text-white'
								: 'border-transparent text-gray-500 hover:bg-black/5 dark:text-gray-400 dark:hover:bg-white/5'}"
							onclick={() => (activeTab = 'meta')}
						>
							Meta
						</button>
						<button
							type="button"
							class="cursor-pointer border-b-2 px-4 py-3 text-sm font-medium transition-colors focus:outline-none {activeTab ===
							'settings'
								? 'border-black text-black dark:border-white dark:text-white'
								: 'border-transparent text-gray-500 hover:bg-black/5 dark:text-gray-400 dark:hover:bg-white/5'}"
							onclick={() => (activeTab = 'settings')}
						>
							Settings
						</button>
						<button
							type="button"
							class="cursor-pointer border-b-2 px-4 py-3 text-sm font-medium transition-colors focus:outline-none {activeTab ===
							'flags'
								? 'border-black text-black dark:border-white dark:text-white'
								: 'border-transparent text-gray-500 hover:bg-black/5 dark:text-gray-400 dark:hover:bg-white/5'}"
							onclick={() => (activeTab = 'flags')}
						>
							Flags
						</button>
						<button
							type="button"
							class="cursor-pointer border-b-2 px-4 py-3 text-sm font-medium transition-colors focus:outline-none {activeTab ===
							'deployment'
								? 'border-black text-black dark:border-white dark:text-white'
								: 'border-transparent text-gray-500 hover:bg-black/5 dark:text-gray-400 dark:hover:bg-white/5'}"
							onclick={() => (activeTab = 'deployment')}
						>
							Deployment
						</button>
					</div>
				</div>

				<!-- Tab Content Container -->
				<div class="flex-1 overflow-y-auto">
					<!-- Meta Tab -->
					{#if activeTab === 'meta'}
						<div class="mt-6 space-y-6 pb-4">
							<!-- Basic Information -->
							<div class="grid gap-6">
								<div class="bg-muted/20 rounded-xl border-0 p-5">
									<h4
										class="text-muted-foreground mb-4 text-sm font-semibold uppercase tracking-wider"
									>
										General
									</h4>
									<div class="space-y-4">
										<div>
											<Label for="ch-name" class="mb-2 block text-sm font-semibold"
												>Challenge Name</Label
											>
											<Input
												id="ch-name"
												bind:value={name}
												placeholder="Enter challenge name"
												class="bg-background"
											/>
										</div>
										<div>
											<Label for="ch-desc" class="mb-2 block text-sm font-semibold"
												>Description</Label
											>
											<Textarea
												id="ch-desc"
												bind:value={description}
												class="bg-background min-h-32"
												placeholder="Describe the challenge..."
											/>
										</div>
									</div>
								</div>

								<div class="bg-muted/20 rounded-xl border-0 p-5">
									<h4
										class="text-muted-foreground mb-4 text-sm font-semibold uppercase tracking-wider"
									>
										Credits
									</h4>
									<div>
										<Label for="ch-auth" class="mb-2 block text-sm font-semibold">Authors</Label>
										<Input
											id="ch-auth"
											bind:value={authorsCsv}
											placeholder="alice, bob"
											class="bg-background"
										/>
										<p class="text-muted-foreground mt-1 text-xs">
											Comma separated list of authors.
										</p>
									</div>
								</div>
							</div>

							<!-- Attachments Section -->
							<div class="border-t pt-6">
								<h4 class="mb-4 text-sm font-semibold">Attachments</h4>

								<input
									bind:this={fileInputEl}
									type="file"
									multiple
									class="hidden"
									onchange={(e) => addFiles((e.currentTarget as HTMLInputElement).files)}
								/>

								<div class="space-y-4">
									<!-- Upload Area -->
									<div>
										<div
											role="button"
											tabindex="0"
											class="text-muted-foreground hover:bg-muted/50 rounded-lg border-2 border-dashed p-6 text-center text-sm transition-colors hover:border-gray-400"
											ondragover={(e) => e.preventDefault()}
											ondrop={(e) => {
												e.preventDefault();
												const dt = e.dataTransfer;
												if (dt?.files?.length) addFiles(dt.files);
											}}
										>
											<p class="mb-2">Drag & drop files here</p>
											<p class="text-xs">or</p>
											<Button
												type="button"
												variant="outline"
												size="sm"
												class="mt-3 focus:outline-none"
												onclick={() => fileInputEl?.click()}
											>
												Browse files
											</Button>
											{#if newFiles.length > 0}
												<p class="text-muted-foreground mt-3 text-xs">
													{newFiles.length} new file{newFiles.length === 1 ? '' : 's'} selected
												</p>
											{/if}
										</div>
									</div>
								</div>

								<!-- Existing attachments -->
								{#if keptExisting.length > 0}
									<div>
										<h5 class="mb-3 text-sm font-medium">Existing Files</h5>
										<div class="flex flex-col gap-2">
											{#each keptExisting as a (a)}
												<div
													class="bg-muted/30 hover:bg-muted/50 flex items-center justify-between rounded-lg border p-3 transition-colors"
												>
													<div class="flex min-w-0 items-center gap-3">
														<TagsIcon class="text-muted-foreground h-4 w-4 shrink-0" />
														<span class="truncate text-sm">{a.split('/').pop() || a}</span>
													</div>
													<Button
														type="button"
														variant="ghost"
														size="icon"
														class="h-8 w-8 shrink-0 focus:outline-none"
														aria-label={`Remove ${a}`}
														title={`Remove ${a}`}
														onclick={() => toggleRemoveExisting(a)}
													>
														<X class="h-4 w-4" />
													</Button>
												</div>
											{/each}
										</div>
									</div>
								{/if}

								<!-- Newly selected files -->
								{#if newFiles.length > 0}
									<div>
										<h5 class="mb-3 text-sm font-medium">New Uploads</h5>
										<div class="flex flex-col gap-2">
											{#each newFiles as f, i (f.name + '::' + f.size)}
												<div
													class="bg-muted/30 hover:bg-muted/50 flex items-center justify-between rounded-lg border p-3 transition-colors"
												>
													<div class="flex min-w-0 items-center gap-3">
														<TagsIcon class="text-muted-foreground h-4 w-4 shrink-0" />
														<span class="truncate text-sm">{f.name}</span>
														<span class="text-muted-foreground shrink-0 text-xs"
															>({formatBytes(f.size)})</span
														>
													</div>
													<Button
														type="button"
														variant="ghost"
														size="icon"
														class="h-8 w-8 shrink-0 focus:outline-none"
														aria-label={`Remove ${f.name}`}
														title={`Remove ${f.name}`}
														onclick={() => removeNewFile(i)}
													>
														<X class="h-4 w-4" />
													</Button>
												</div>
											{/each}
										</div>
									</div>
								{/if}
							</div>
						</div>
					{/if}

					<!-- Settings Tab -->
					{#if activeTab === 'settings'}
						<div class="mt-6 space-y-6 pb-4">
							<div class="bg-muted/20 rounded-xl border-0 p-5">
								<h4
									class="text-muted-foreground mb-4 text-sm font-semibold uppercase tracking-wider"
								>
									Organization
								</h4>
								<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
									<div class="col-span-1 sm:col-span-2">
										<Label for="ch-cat" class="mb-2 block text-sm font-semibold">Category</Label>
										<Input
											id="ch-cat"
											bind:value={category}
											placeholder="Enter category"
											class="bg-background"
										/>
									</div>

									<div class="col-span-1 sm:col-span-2">
										<Label class="mb-2 block text-sm font-semibold">Tags</Label>
										<TagMultiSelect
											all_tags={uniqAllTags}
											bind:value={tags}
											oncreate={(newTag) => {
												const tag = String(newTag ?? '').trim();
												if (tag && !tags.includes(tag)) tags = [...tags, tag];
											}}
										/>
									</div>
								</div>
							</div>

							<!-- Difficulty & Type -->
							<div class="bg-muted/20 rounded-xl border-0 p-5">
								<h4
									class="text-muted-foreground mb-4 text-sm font-semibold uppercase tracking-wider"
								>
									Classification
								</h4>
								<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
									<div>
										<Label for="ch-diff" class="mb-2 block text-sm font-medium">Difficulty</Label>
										<CategorySelect
											id="ch-diff"
											items={difficultyOptions}
											bind:value={difficulty}
											placeholder="Select difficulty..."
											searchPlaceholder="Search difficulty"
										/>
									</div>
									<div>
										<Label for="ch-type" class="mb-2 block text-sm font-medium">Type</Label>
										<CategorySelect
											id="ch-type"
											items={typeOptions}
											bind:value={type}
											placeholder="Select type..."
											searchPlaceholder="Search type"
										/>
									</div>
								</div>
							</div>

							<!-- Visibility & Scoring -->
							<div class="border-t pt-6">
								<h4 class="mb-4 text-sm font-semibold">Options</h4>
								<div class="flex flex-col gap-6 sm:flex-row">
									<div class="flex items-center gap-3">
										<Checkbox id="ch-hidden" bind:checked={hidden} />
										<Label for="ch-hidden" class="cursor-pointer">Hidden</Label>
									</div>
									<div class="flex items-center gap-3">
										<Checkbox id="ch-score" bind:checked={dynamicScoring} />
										<Label for="ch-score" class="cursor-pointer">Dynamic scoring</Label>
									</div>
								</div>
							</div>

							<!-- Points -->
							<div class="border-t pt-6">
								<div>
									<Label for="ch-points" class="mb-2 block text-sm font-semibold">Max Points</Label>
									<div class="flex items-center gap-3">
										<Input
											id="ch-points"
											type="number"
											inputmode="numeric"
											min="0"
											max="1500"
											step="1"
											bind:value={maxPoints}
											class="max-w-[200px]"
											oninput={(e) => {
												const v = (e.currentTarget as HTMLInputElement).valueAsNumber;
												const n = Number.isFinite(v) ? Math.max(0, Math.floor(v)) : 0;
												maxPoints = n;
												e.currentTarget.value = String(n);
											}}
										/>
										<span class="text-muted-foreground text-sm font-medium">points</span>
									</div>
								</div>
							</div>
						</div>
					{/if}

					<!-- Flags Tab -->
					{#if activeTab === 'flags'}
						<div class="mt-6 space-y-6 pb-4">
							<div>
								<div class="mb-4 flex items-center justify-between">
									<h4 class="text-sm font-semibold">Challenge Flags</h4>
									<Button
										type="button"
										variant="outline"
										size="sm"
										class="focus:outline-none"
										onclick={() => (flags = [...flags, { flag: '', regex: false }])}
									>
										Add Flag
									</Button>
								</div>

								<div class="flex flex-col gap-3">
									{#each flags as flag, index (index)}
										<div class="bg-muted/30 flex items-center gap-3 rounded-lg border p-3">
											<Input
												bind:value={flags[index].flag}
												placeholder="Enter flag value"
												class="flex-1"
											/>
											<div class="flex shrink-0 items-center gap-2">
												<Checkbox id={'flag-' + index} bind:checked={flags[index].regex} />
												<Label for={'flag-' + index} class="cursor-pointer whitespace-nowrap"
													>Regex</Label
												>
											</div>
											<Button
												type="button"
												variant="ghost"
												size="icon"
												class="h-9 w-9 shrink-0 focus:outline-none"
												onclick={() => (flags = flags.filter((_: any, i: any) => i !== index))}
												aria-label="Remove flag"
											>
												<X class="h-4 w-4" />
											</Button>
										</div>
									{/each}

									{#if flags.length === 0}
										<div
											class="text-muted-foreground rounded-lg border-2 border-dashed p-8 text-center text-sm"
										>
											<p>No flags added yet</p>
											<p class="mt-1 text-xs">Click "Add Flag" to create one</p>
										</div>
									{/if}
								</div>
							</div>
						</div>
					{/if}

					<!-- Deployment Tab -->
					{#if activeTab === 'deployment'}
						<div class="mt-6 space-y-6 pb-4">
							{#if type === 'Normal'}
								<!-- Normal Challenge Deployment -->
								<div class="bg-muted/20 rounded-xl border-0 p-5">
									<h4
										class="text-muted-foreground mb-4 text-sm font-semibold uppercase tracking-wider"
									>
										Connection Details
									</h4>
									<div class="grid gap-4 sm:grid-cols-2">
										<div class="col-span-2">
											<Label for="ch-host" class="mb-2 block text-sm font-medium">Host</Label>
											<Input
												id="ch-host"
												bind:value={host}
												placeholder="e.g. challenge.trxd.cc"
												class="bg-background"
											/>
										</div>
										<div>
											<Label for="ch-port" class="mb-2 block text-sm font-medium">Port</Label>
											<Input
												id="ch-port"
												bind:value={portStr}
												placeholder="e.g. 31337"
												class="bg-background"
											/>
										</div>
										<div>
											<Label for="ch-conn-type" class="mb-2 block text-sm font-medium"
												>Connection Type</Label
											>
											<select
												id="ch-conn-type"
												bind:value={connType}
												class="border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border px-3 py-2 text-sm file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
											>
												<option value="NONE">None</option>
												<option value="TCP">TCP</option>
												<option value="TCP_TLS">TCP + TLS</option>
												<option value="HTTP">HTTP</option>
												<option value="HTTPS">HTTPS</option>
											</select>
										</div>
									</div>
								</div>
							{:else if type === 'Compose'}
								<!-- Compose Challenge Deployment -->
								<div class="space-y-4">
									<div>
										<Label class="mb-2 block text-sm font-semibold">Compose Configuration</Label>
										{#if open && activeTab === 'deployment'}
											<MonacoEditor bind:value={composeFile} language="yaml" class="mt-3" />
										{/if}
									</div>
									<div class="flex items-center gap-3">
										<Checkbox id="com-hashdomain" bind:checked={hashDomain} />
										<Label for="com-hashdomain" class="cursor-pointer">Hash Domain</Label>
									</div>
								</div>

								<div class="space-y-4 border-t pt-6">
									<h4 class="mb-4 text-sm font-semibold">Connection Details</h4>
									<div>
										<Label for="ch-host" class="mb-2 block text-sm font-medium">Host</Label>
										<Input id="ch-host" bind:value={host} placeholder="e.g. challenge.trxd.cc" />
									</div>
									<div>
										<Label for="com-envs" class="mb-2 block text-sm font-medium"
											>Environment Variables</Label
										>
										<Textarea
											id="com-envs"
											bind:value={envs}
											placeholder="FLAG1=flag&#10;FLAG2=flag"
										/>
									</div>
									<div>
										<Label for="com-conn-type" class="mb-2 block text-sm font-medium"
											>Connection Type</Label
										>
										<select
											id="com-conn-type"
											bind:value={connType}
											class="border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border px-3 py-2 text-sm file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
										>
											<option value="NONE">None</option>
											<option value="TCP">TCP</option>
											<option value="TCP_TLS">TCP + TLS</option>
											<option value="HTTP">HTTP</option>
											<option value="HTTPS">HTTPS</option>
										</select>
									</div>
								</div>
							{:else}
								<!-- Container Challenge Deployment -->
								<div class="space-y-4">
									<h4 class="mb-4 text-sm font-semibold">Container Configuration</h4>
									<div>
										<Label for="ho-host" class="mb-2 block text-sm font-medium"
											>Container Image Name</Label
										>
										<Input id="ho-host" bind:value={imageName} placeholder="TRX-Chall-1" />
									</div>
									<div class="flex items-center gap-3">
										<Checkbox id="ho-hashdomain" bind:checked={hashDomain} />
										<Label for="ho-hashdomain" class="cursor-pointer">Hash domain</Label>
									</div>
								</div>

								<div class="space-y-4 border-t pt-6">
									<h4 class="mb-4 text-sm font-semibold">Connection Details</h4>
									<div>
										<Label for="ch-host" class="mb-2 block text-sm font-medium">Host</Label>
										<Input id="ch-host" bind:value={host} placeholder="e.g. challenge.trxd.cc" />
									</div>
									<div>
										<Label for="ho-port" class="mb-2 block text-sm font-medium">Port</Label>
										<Input id="ho-port" bind:value={portStr} placeholder="e.g. 31337" />
									</div>
									<div>
										<Label for="ho-conn-type" class="mb-2 block text-sm font-medium"
											>Connection Type</Label
										>
										<select
											id="ho-conn-type"
											bind:value={connType}
											class="border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border px-3 py-2 text-sm file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
										>
											<option value="NONE">None</option>
											<option value="TCP">TCP</option>
											<option value="TCP_TLS">TCP + TLS</option>
											<option value="HTTP">HTTP</option>
											<option value="HTTPS">HTTPS</option>
										</select>
									</div>
									<div>
										<Label for="con-envs" class="mb-2 block text-sm font-medium"
											>Environment Variables</Label
										>
										<Textarea
											id="con-envs"
											bind:value={envs}
											placeholder="FLAG1=flag&#10;FLAG2=flag"
										/>
									</div>
								</div>
							{/if}

							<!-- Performance settings for Container and Compose -->
							{#if type === 'Container' || type === 'Compose'}
								<div class="border-t pt-6">
									<h4 class="mb-4 text-sm font-semibold">Performance Limits</h4>
									<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
										<div class="flex items-start gap-3">
											<Cpu class="text-muted-foreground mt-2 h-5 w-5 shrink-0" />
											<div class="flex-1">
												<Label for="perf-cpu" class="mb-2 block text-sm font-medium"
													>Max CPU (Cores)</Label
												>
												<Input id="perf-cpu" bind:value={maxCPU} placeholder="e.g. 1" />
											</div>
										</div>
										<div class="flex items-start gap-3">
											<MemoryStick class="text-muted-foreground mt-2 h-5 w-5 shrink-0" />
											<div class="flex-1">
												<Label for="perf-ram" class="mb-2 block text-sm font-medium"
													>Max RAM (MB)</Label
												>
												<Input id="perf-ram" bind:value={maxRam} placeholder="e.g. 512" />
											</div>
										</div>
									</div>
									<div class="mt-4 flex items-start gap-3">
										<Clock class="text-muted-foreground mt-2 h-5 w-5 shrink-0" />
										<div class="flex-1">
											<Label for="perf-lifetime" class="mb-2 block text-sm font-medium"
												>Lifetime (seconds)</Label
											>
											<Input
												id="perf-lifetime"
												bind:value={lifetime}
												placeholder="300"
												class="max-w-[200px]"
											/>
										</div>
									</div>
								</div>
							{/if}
						</div>
					{/if}
				</div>
				<!-- End Tab Content Container -->

				<div class="mb-4 mt-6 flex justify-end gap-2">
					<Sheet.Close>
						<Button type="button" variant="outline" class="focus:outline-none">Cancel</Button>
					</Sheet.Close>
					<Button type="submit" disabled={saving} class="focus:outline-none">
						{#if saving}Saving...{:else}Save{/if}
					</Button>
				</div>
			</form>
		</div>
	</Sheet.Content>
</Sheet.Root>
