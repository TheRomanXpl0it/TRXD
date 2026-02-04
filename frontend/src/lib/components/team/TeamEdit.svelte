<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import CountrySelect from '$lib/components/ui/country-select.svelte';
	import GeneratedAvatar from '$lib/components/ui/avatar/generated-avatar.svelte';
	import { updateTeam } from '$lib/team';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { showSuccess, showError } from '$lib/utils/toast';
	import { getCountryByIso3 } from '$lib/utils/countries';

	const queryClient = useQueryClient();

	let {
		open = $bindable(false),
		team,
		onupdated
	} = $props<{
		open?: boolean;
		team?: {
			id: number;
			name?: string;
			country?: string;
			tags?: string[];
		};
		onupdated?: (detail: { id: number }) => void;
	}>();

	let name = $state('');
	let countryCode = $state('');
	let tags = $state<string[]>([]);
	let newTag = $state('');
	let saving = $state(false);

	$effect(() => {
		if (team) {
			name = team.name ?? '';
			countryCode = team.country?.toUpperCase?.() ?? '';
			tags = team.tags ?? [];
		}
	});

	async function onSave(e: Event) {
		e.preventDefault();
		if (saving) return;

		const id = team?.id ?? 0;
		const n = name.trim();
		const c = countryCode.trim();

		if (!n && !c) {
			showError(null, 'Please fill at least one field.');
			return;
		}

		try {
			saving = true;
			await updateTeam(id, n, c, tags);
			open = false;
			queryClient.invalidateQueries({ queryKey: ['teams'] });
			onupdated?.({ id });
			showSuccess('Team updated.');
		} catch (err: any) {
			showError(err, 'Failed to update team.');
		} finally {
			saving = false;
		}
	}

	function addTag() {
		const t = newTag.trim();
		if (t && !tags.includes(t)) {
			tags = [...tags, t];
			newTag = '';
		}
	}

	function removeTag(tag: string) {
		tags = tags.filter((t) => t !== tag);
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="w-full px-5 sm:max-w-[640px]">
		<div
			class="from-muted/20 to-background mb-6 mt-4 rounded-xl border-0 bg-gradient-to-br p-6 shadow-sm"
		>
			<div class="flex items-center gap-4">
				<div
					class="bg-background border-background h-16 w-16 shrink-0 overflow-hidden rounded-full border-4 shadow-sm"
				>
					<GeneratedAvatar seed={name} class="h-full w-full" />
				</div>
				<div>
					<Sheet.Title class="text-xl font-bold">Edit Team</Sheet.Title>
					<Sheet.Description class="text-muted-foreground/80 mt-1">
						Update, modify or delete your team.
					</Sheet.Description>
				</div>
			</div>
		</div>

		<form class="mt-3 space-y-6" onsubmit={onSave}>
			<div class="space-y-4">
				<div>
					<Label for="pf-name" class="mb-1 block">Team name</Label>
					<Input id="pf-name" bind:value={name} placeholder={team?.name || 'Team name'} />
					{#if team?.name && team.name !== name}
						<p class="text-muted-foreground mt-1 text-sm">
							Current: {team.name}
						</p>
					{/if}
				</div>

				<div>
					<Label for="pf-country" class="mb-1 block">Country</Label>
					<CountrySelect id="pf-country" bind:value={countryCode} />
					{#if team?.country && team.country !== countryCode}
						{@const current = getCountryByIso3(team.country)}
						<p class="text-muted-foreground mt-1 text-sm">
							Current: {current?.name ?? team.country}
						</p>
					{/if}
				</div>
			</div>

			<div class="mt-8 flex justify-end gap-2">
				<Sheet.Close>
					<Button type="button" variant="outline">Cancel</Button>
				</Sheet.Close>
				<Button type="submit" disabled={saving}>
					{#if saving}Saving...{:else}Save{/if}
				</Button>
			</div>
		</form>
	</Sheet.Content>
</Sheet.Root>
