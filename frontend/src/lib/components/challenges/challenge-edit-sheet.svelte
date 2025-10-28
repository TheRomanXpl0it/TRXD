<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Slider } from '$lib/components/ui/slider/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import CategorySelect from '$lib/components/challenges/category-select.svelte';
	import TagMultiSelect from '$lib/components/challenges/tag-multiselect.svelte';
	import { toast } from 'svelte-sonner';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import { updateChallengeMultipart, getChallenge } from '$lib/challenges';
	import { Check, Cpu, MemoryStick, Clock, X, Tags, UserPen, FileDown } from '@lucide/svelte';

	// --- CodeMirror imports (YAML + theming) ---
	import { EditorState, Compartment } from '@codemirror/state';
	import {
		EditorView,
		keymap,
		lineNumbers,
		highlightActiveLine,
		highlightActiveLineGutter
	} from '@codemirror/view';
	import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands';
	import { yaml } from '@codemirror/lang-yaml';
	import { oneDark } from '@codemirror/theme-one-dark';
	import { syntaxHighlighting, defaultHighlightStyle } from '@codemirror/language';

	type Item = { value: string; label: string };

	// Props
	let {
		open = $bindable(false),
		challenge_user,
		all_tags
	} = $props<{ open?: boolean; challenge_user: any; all_tags?: string[] }>();

	let challenge = $state<any>(null);

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

	// 👇 Deduped tags as **strings**
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

	// --- Attachments (existing + new uploads) ---
	let existingAttachments = $state<string[]>([]); // from backend
	let removedAttachmentNames = $state<Set<string>>(new Set()); // user-marked for deletion
	let newFiles = $state<File[]>([]); // newly selected files
	let fileInputEl: HTMLInputElement | null = null; // hidden <input type="file">

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
				console.log('Hash Domain:', challenge?.docker_config?.hash_domain);
				hashDomain = Boolean(challenge?.docker_config?.hash_domain ?? false);
				composeFile = String(challenge?.docker_config?.compose ?? '');
				maxCPU = String(challenge?.docker_config?.max_cpu ?? '');
				maxRam = String(challenge?.docker_config?.max_memory ?? '');
				lifetime = String(challenge?.docker_config?.lifetime ?? '');
				envs = String(challenge?.docker_config?.envs ?? '');

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

			// container/compose specifics
			image: type === 'Container' ? str(imageName) : undefined,
			compose: type === 'Compose' ? str(composeFile) : undefined,
			hash_domain: hashDomain,

			// performance / limits — only include if > 0
			max_points: toNum(maxPoints) && maxPoints > 0 ? maxPoints : undefined,
			lifetime: toNum(lifetime) && Number(lifetime) > 0 ? Number(lifetime) : undefined,
			max_memory: toNum(maxRam) && Number(maxRam) > 0 ? Number(maxRam) : undefined,

			// misc docker
			envs: str(envs),
			max_cpu: str(maxCPU)
		};

		// strip empty arrays
		if (Array.isArray(fields.authors) && fields.authors.length === 0) delete fields.authors;

		saving = true;

		// tags diffs
		const prev = Array.from(challenge_user.tags ?? []);
		const curr = Array.from(tags ?? []);
		const deletedTags = prev.filter((t: string) => !curr.includes(t));
		const newTags = curr.filter((t: string) => !prev.includes(t));

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

			for (const [k, v] of Object.entries(entries)) {
				if (v !== undefined && v !== null && v !== '') {
					fd.append(k, String(v));
				}
			}

			// 2) Existing attachments to KEEP (optional)
			for (const name of keptExisting) {
				fd.append('keep_attachments[]', name);
			}

			// 3) Existing attachments to DELETE
			if (removedAttachmentNames.size > 0) {
				fd.append('deleted_attachments', JSON.stringify(Array.from(removedAttachmentNames)));
			}

			// 4) New files to UPLOAD
			for (const f of newFiles) {
				fd.append('attachments[]', f, f.name);
			}

			// 5) Tags / Flags diffs
			if (deletedTags.length) {
				fd.append('deleted_tags', JSON.stringify(deletedTags));
			}
			if (newTags.length) {
				fd.append('new_tags', JSON.stringify(newTags));
			}
			if (deletedFlags.length) {
				fd.append('deleted_flags', JSON.stringify(deletedFlags));
			}
			if (newFlags.length) {
				fd.append('new_flags', JSON.stringify(newFlags));
			}

			// ---- Submit ----
			console.log('Multipartform:', fd);
			await updateChallengeMultipart(fd);

			open = false;

			setTimeout(() => {
				toast.success('Challenge updated.');
			}, 1000);
			// Reload the page to reflect the changes
			window.location.reload();
		} catch (err: any) {
			toast.error(err?.message ?? 'Failed to update challenge.');
		} finally {
			saving = false;
		}
	}

	// ---------- CodeMirror setup for the Compose YAML field ----------
	let cmComposeHost: HTMLDivElement | undefined;
	let cmComposeView: EditorView | null = null;

	// Theme compartment + light theme
	const themeComp = new Compartment();
	function lightTheme() {
		return EditorView.theme(
			{
				'&': { backgroundColor: 'transparent' },
				'.cm-content': {
					fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
					fontSize: '0.9rem'
				},
				'.cm-gutters': {
					backgroundColor: 'transparent',
					borderRight: 'none',
					color: 'rgb(100 116 139)'
				},
				'&.cm-editor': { borderRadius: '0.5rem' },
				'&.cm-focused': { outline: 'none' }
			},
			{ dark: false }
		);
	}

	$effect(() => {
		// Destroy editor when sheet closes or type changes away from Compose
		if ((!open || type !== 'Compose') && cmComposeView) {
			cmComposeView.destroy();
			cmComposeView = null;
			if (cmComposeHost) cmComposeHost.innerHTML = '';
		}

		// Create editor when sheet is open and type is Compose
		if (open && type === 'Compose' && cmComposeHost && !cmComposeView) {
			cmComposeView = new EditorView({
				parent: cmComposeHost,
				state: EditorState.create({
					doc: composeFile && composeFile.length ? composeFile : '\n',
					extensions: [
						lineNumbers(),
						highlightActiveLineGutter(),
						highlightActiveLine(),
						history(),
						keymap.of([indentWithTab, ...defaultKeymap, ...historyKeymap]),
						yaml(),
						syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
						themeComp.of(lightTheme()),
						EditorView.theme({ '&': { minHeight: '180px', borderRadius: '0.5rem' } }),
						EditorView.updateListener.of((u) => {
							if (u.docChanged) composeFile = u.state.doc.toString();
						})
					]
				})
			});
		}
	});

	$effect(() => {
		if (!cmComposeView) return;
		const current = cmComposeView.state.doc.toString();
		const desired = composeFile && composeFile.length ? composeFile : '\n';
		if (desired !== current) {
			cmComposeView.dispatch({
				changes: { from: 0, to: current.length, insert: desired }
			});
		}
	});

	$effect(() => {
		if (!cmComposeView) return;

		const isDarkNow = () =>
			typeof document !== 'undefined' &&
			(document.documentElement.classList.contains('dark') ||
				(typeof window !== 'undefined' &&
					window.matchMedia?.('(prefers-color-scheme: dark)')?.matches));

		const apply = () => {
			cmComposeView?.dispatch({
				effects: themeComp.reconfigure(isDarkNow() ? oneDark : lightTheme())
			});
		};

		apply();

		const mq = window.matchMedia?.('(prefers-color-scheme: dark)');
		const mqHandler = () => apply();
		mq?.addEventListener?.('change', mqHandler);

		const obs = new MutationObserver(apply);
		obs.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });

		return () => {
			mq?.removeEventListener?.('change', mqHandler);
			obs.disconnect();
		};
	});
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="px-5 sm:max-w-[720px]">
		<Sheet.Header>
			<Sheet.Title>Edit Challenge</Sheet.Title>
			<Sheet.Description>
				Modify settings for <b>{challenge?.name ?? '—'}</b>.
			</Sheet.Description>
		</Sheet.Header>

		<form class="mt-3 space-y-5" onsubmit={onSave}>
			<Accordion.Root type="single">
				<Accordion.Item value="Naming">
					<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight"
						>Naming and description</Accordion.Trigger
					>
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
					<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight"
						>Scoring and visibility</Accordion.Trigger
					>
					<Accordion.Content>
						<div class="mt-3 flex flex-row justify-between">
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
							<Input
								id="ch-points"
								type="number"
								inputmode="numeric"
								min="0"
								max="1500"
								step="1"
								bind:value={maxPoints}
								oninput={(e) => {
									const v = (e.currentTarget as HTMLInputElement).valueAsNumber;
									const n = Number.isFinite(v) ? Math.max(0, Math.floor(v)) : 0;
									maxPoints = n;
									e.currentTarget.value = String(n);
								}}
							/>
						</div>
						<div class="mt-5">
							<Label for="ch-flags" class="mb-1 block">Flags</Label>
							<div class="flex flex-col gap-2">
								{#each flags as flag, index (index)}
									<div class="flex items-center gap-2">
										<Input bind:value={flags[index].flag} placeholder="Flag value" class="flex-1" />
										<Checkbox id={'flag-' + index} bind:checked={flags[index].regex} />
										<Label for={'flag-' + index}>Regex</Label>

										<Button
											type="button"
											variant="destructive"
											size="icon"
											onclick={() => (flags = flags.filter((_: any, i: any) => i !== index))}
										>
											<X class="h-4 w-4" />
										</Button>
									</div>
								{/each}
								<Button
									type="button"
									variant="outline"
									size="sm"
									onclick={() => (flags = [...flags, { flag: '', regex: false }])}
								>
									Add Flag
								</Button>
							</div>
						</div>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="Difficulty">
					<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight"
						>Difficulty and categorization</Accordion.Trigger
					>
					<Accordion.Content>
						<div class="flex flex-row justify-between">
							<div>
								<Label for="ch-diff" class="mb-1 block">Difficulty</Label>
								<CategorySelect
									id="ch-diff"
									items={difficultyOptions}
									bind:value={difficulty}
									placeholder="Select difficulty…"
									searchPlaceholder="Search difficulty"
								/>
							</div>
							<div>
								<Label for="ch-type" class="mb-1 block">Type</Label>
								<CategorySelect
									id="ch-type"
									items={typeOptions}
									bind:value={type}
									placeholder="Select type…"
									searchPlaceholder="Search type"
								/>
							</div>
							<div>
								<Label for="ch-cat" class="mb-1 block">Category</Label>
								<Input id="ch-cat" bind:value={category} placeholder="Enter category" />
							</div>
						</div>
						<div class="mt-3">
							<Label class="mb-1 block">Tags</Label>
							<TagMultiSelect
								all_tags={uniqAllTags}
								bind:value={tags}
								on:create={(e) => {
									const newTag = String(e.detail ?? '').trim();
									if (newTag && !tags.includes(newTag)) tags = [...tags, newTag];
								}}
							/>
						</div>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="Authors">
					<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight"
						>Authors and attachments</Accordion.Trigger
					>
					<Accordion.Content>
						<div class="flex flex-col">
							<div class="flex flex-row items-center">
								<UserPen class="mr-2" />
								<div>
									<Label for="ch-auth" class="mb-1 block">Authors</Label>
									<Input id="ch-auth" bind:value={authorsCsv} placeholder="alice, bob" />
								</div>
							</div>

							<!-- Attachments: shadcn-svelte upload UI -->
							<div class="mt-3 flex flex-row">
								<FileDown class="mr-2" />
								<div class="flex-1">
									<Label class="mb-1 block">Attachments</Label>

									<input
										bind:this={fileInputEl}
										type="file"
										multiple
										class="hidden"
										onchange={(e) => addFiles((e.currentTarget as HTMLInputElement).files)}
									/>

									<div class="flex flex-wrap items-center gap-2">
										<Button
											type="button"
											variant="outline"
											size="sm"
											onclick={() => fileInputEl?.click()}
										>
											Select files…
										</Button>
										<span class="text-muted-foreground text-sm">
											{newFiles.length} new file{newFiles.length === 1 ? '' : 's'} selected
										</span>
									</div>

									<div
										class="text-muted-foreground mt-3 rounded border border-dashed p-4 text-sm"
										ondragover={(e) => e.preventDefault()}
										ondrop={(e) => {
											e.preventDefault();
											const dt = e.dataTransfer;
											if (dt?.files?.length) addFiles(dt.files);
										}}
									>
										Drag & drop files here, or use “Select files…”
									</div>

									<!-- Existing attachments -->
									{#if keptExisting.length > 0 || removedAttachmentNames.size > 0}
										<div class="mt-4">
											<Label class="mb-2 block">Existing files</Label>
											<div class="flex flex-col gap-2">
												{#each existingAttachments as a (a)}
													<div class="flex items-center justify-between rounded border p-2">
														<div class="flex items-center gap-2">
															<Tags class="h-4 w-4" />
															<span class="max-w-[380px] truncate">{a}</span>
														</div>
														<div class="flex items-center gap-2">
															<Checkbox
																id={`rm-${a}`}
																checked={removedAttachmentNames.has(a)}
																onchange={() => toggleRemoveExisting(a)}
															/>
															<Label for={`rm-${a}`} class="text-sm">Remove</Label>
														</div>
													</div>
												{/each}
											</div>
										</div>
									{/if}

									<!-- Newly selected files -->
									{#if newFiles.length > 0}
										<div class="mt-4">
											<Label class="mb-2 block">New uploads</Label>
											<div class="flex flex-col gap-2">
												{#each newFiles as f, i (f.name + '::' + f.size)}
													<div class="flex items-center justify-between rounded border p-2">
														<div class="flex items-center gap-2">
															<Tags class="h-4 w-4" />
															<span class="max-w-[320px] truncate">{f.name}</span>
															<span class="text-muted-foreground text-xs"
																>({formatBytes(f.size)})</span
															>
														</div>
														<Button
															type="button"
															variant="destructive"
															size="icon"
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
						</div>
					</Accordion.Content>
				</Accordion.Item>

				{#if type === 'Normal'}
					<Accordion.Item value="Host">
						<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight"
							>Host and port</Accordion.Trigger
						>
						<Accordion.Content>
							<div class="flex flex-col">
								<div>
									<Label for="ch-host" class="mb-1 block">Host</Label>
									<Input id="ch-host" bind:value={host} placeholder="e.g. challenge.trxd.cc" />
								</div>
								<div class="mt-3 flex flex-row items-center justify-between">
									<div>
										<Label for="ch-port" class="mb-1 block" aria-disabled={hashDomain}>Port</Label>
										<Input
											id="ch-port"
											bind:value={portStr}
											placeholder="e.g. 31337"
											disabled={hashDomain}
										/>
									</div>
								</div>
							</div>
						</Accordion.Content>
					</Accordion.Item>
				{:else if type === 'Compose'}
					<Accordion.Item value="Compose">
						<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight"
							>Compose Settings</Accordion.Trigger
						>
						<Accordion.Content>
							<div class="flex flex-col">
								<Label>Compose.yaml</Label>
								<div bind:this={cmComposeHost} class="min-h-45 mt-3 rounded border"></div>
								<div class="mt-3 flex flex-row items-center">
									<Checkbox id="com-hashdomain" bind:checked={hashDomain} />
									<Label for="com-hashdomain" class="ml-1">Hash Domain</Label>
								</div>
								<div class="mt-3">
									<Label for="com-envs" class="mb-1 block">Envs</Label>
									<Input id="com-envs" bind:value={envs} placeholder="LD_PRELOAD=..." />
								</div>
							</div>
						</Accordion.Content>
					</Accordion.Item>
				{:else}
					<Accordion.Item value="Container">
						<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight"
							>Container Settings</Accordion.Trigger
						>
						<Accordion.Content>
							<div class="flex flex-col">
								<div>
									<Label for="ho-host" class="mb-1 block">Container Image name</Label>
									<Input id="ho-host" bind:value={imageName} placeholder="TRX-Chall-1" />
								</div>
								<div class="flex flex-row items-center justify-between">
									<div class="mt-3 flex flex-row items-center">
										<Checkbox id="ho-hashdomain" bind:checked={hashDomain} />
										<Label for="ho-hashdomain" class="ml-2">Hash domain</Label>
									</div>
									<div class="mt-3">
										<Label for="ho-port">Port</Label>
										<Input
											id="ho-port"
											bind:value={portStr}
											placeholder="e.g. 31337"
											disabled={hashDomain}
										/>
									</div>
								</div>
								<div class="mt-3">
									<Label for="con-envs" class="mb-1 block">Envs</Label>
									<Input id="con-envs" bind:value={envs} placeholder="LD_PRELOAD=..." />
								</div>
							</div>
						</Accordion.Content>
					</Accordion.Item>
				{/if}

				{#if type === 'Container' || type === 'Compose'}
					<Accordion.Item value="Performance">
						<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight"
							>Performance Settings</Accordion.Trigger
						>
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
