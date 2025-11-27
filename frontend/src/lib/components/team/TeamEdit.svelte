<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import { Avatar } from 'flowbite-svelte';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import CountrySelect from '$lib/components/ui/country-select.svelte';
	import Icon from '@iconify/svelte';
	import { updateTeam } from '$lib/team';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { isValidUrl } from '$lib/utils/validation';
	import { showSuccess, showError } from '$lib/utils/toast';
	import { getCountryByIso3 } from '$lib/utils/countries';
        
	const queryClient = useQueryClient();

	let { 
		open = $bindable(false), 
		team,
		onupdated
	} = $props<{
		open?: boolean;
		team?: { id: number; name?: string; country?: string; image?: string; bio?: string };
		onupdated?: (detail: { id: number }) => void;
	}>();

	let name = $state('');
	let bio = $state('');
	let imageUrl = $state('');
	let countryCode = $state('');
	let saving = $state(false);

	$effect(() => {
		if (team) {
			name = team.name ?? '';
			bio = team.bio ?? '';
			imageUrl = team.image ?? '';
			countryCode = team.country?.toUpperCase?.() ?? '';
		}
	});

	async function onSave(e: Event) {
		e.preventDefault();
		if (saving) return;

		const id = team?.id ?? 0;
		const n = name.trim();
		const b = bio.trim();
		const c = countryCode.trim();
		const i = imageUrl.trim();

		if (!n && !b && !c && !i) {
			showError(null, 'Please fill at least one field.');
			return;
		}

		if (i && !isValidUrl(i)) {
			showError(null, 'Image must be a valid URL.');
			return;
		}

		try {
			saving = true;
			await updateTeam(id, n, b, i, c);
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
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="w-full px-5 sm:max-w-[640px]">
		<form class="mt-3 space-y-6" onsubmit={onSave}>
			<div class="space-y-4">
				<h2 class="text-xl font-semibold tracking-tight">Team</h2>
				
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
					<Label for="pf-bio" class="mb-1 block">Bio</Label>
					<Textarea
						id="pf-bio"
						bind:value={bio}
						rows={5}
						placeholder={team?.bio || 'Tell us about your team'}
					/>
					{#if team?.bio && team.bio !== bio}
						<p class="text-muted-foreground mt-1 text-sm">Current bio will be replaced</p>
					{/if}
				</div>

				<div>
					<Label for="pf-country" class="mb-1 block">Country</Label>
					<CountrySelect id="pf-country" bind:value={countryCode} />
					{#if team?.country && team.country.toUpperCase() !== countryCode}
						{@const teamCountry = getCountryByIso3(team.country)}
						<div class="text-muted-foreground mt-1 flex items-center gap-2 text-sm">
							{#if teamCountry?.iso2}
								<Icon
									icon={`circle-flags:${teamCountry.iso2.toLowerCase()}`}
									width="16"
									height="16"
								/>
							{/if}
							<span>Current: {team.country.toUpperCase()}</span>
						</div>
					{/if}
				</div>

				<div class="flex items-start gap-4">
					<!-- Preview/Current Image -->
					<div class="shrink-0">
						{#if imageUrl && isValidUrl(imageUrl)}
							<Avatar src={imageUrl} class="h-24 w-24" />
						{:else if team?.image && isValidUrl(team.image)}
							<Avatar src={team.image} class="h-24 w-24" />
						{:else}
							<div class="flex h-24 w-24 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800">
								<span class="text-xs text-gray-400">No image</span>
							</div>
						{/if}
					</div>
					
					<!-- Input Field -->
					<div class="flex-1 min-w-0">
						<Label for="pf-image" class="mb-1 block">Image URL</Label>
						<Input
							id="pf-image"
							bind:value={imageUrl}
							placeholder={team?.image || 'https://.../avatar.png'}
						/>
						{#if team?.image && team.image !== imageUrl}
							<p class="text-muted-foreground mt-1 text-sm">Current image will be replaced</p>
						{/if}
					</div>
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
