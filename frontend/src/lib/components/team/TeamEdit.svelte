<script lang="ts">
	/* UI */
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';
	import { Avatar } from 'flowbite-svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import countries from '$lib/data/countries.json';

	import { updateTeam } from '$lib/team';
    import { createEventDispatcher } from 'svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
        
    const dispatch = createEventDispatcher();
	const queryClient = useQueryClient();

	let { open = $bindable(false), team } = $props<{
		open?: boolean;
		team?: { id: number; name?: string; country?: string; image?: string; bio?: string };
	}>();

	// Initialize form fields with current team data
	let name = $state('');
	let bio = $state('');
	let imageUrl = $state('');
	let countryCode = $state<string>('');
	let saving = $state(false);

	type Country = { name: string; iso2: string; iso3?: string; emoji?: string };
	const countryItems = (countries as Country[])
		.filter((c) => c.iso3) // Only include countries with iso3
		.map((c) => ({ value: c.iso3!.toUpperCase(), label: c.name, iso2: c.iso2.toUpperCase() }))
		.sort((a, b) => a.label.localeCompare(b.label));
	
	let countrySearch = $state('');
	const filteredCountries = $derived(
		countrySearch.trim() 
			? countryItems.filter(c => c.label.toLowerCase().includes(countrySearch.toLowerCase()) || c.value.toLowerCase().includes(countrySearch.toLowerCase())).slice(0, 50)
			: countryItems.slice(0, 50)
	);

	// Watch for team changes and update form fields
	$effect(() => {
		if (team) {
			name = team.name ?? '';
			bio = team.bio ?? '';
			imageUrl = team.image ?? '';
			countryCode = team.country?.toUpperCase?.() ?? '';
		}
	});

	// Reset form when sheet opens
	$effect(() => {
		if (open && team) {
			name = team.name ?? '';
			bio = team.bio ?? '';
			imageUrl = team.image ?? '';
			countryCode = team.country?.toUpperCase?.() ?? '';
		}
	});

	function isLikelyUrl(s: string): boolean {
		if (!s) return false;
		try {
			new URL(s);
			return true;
		} catch {
			return false;
		}
	}

	async function onSave(e: Event) {
		e.preventDefault();
		if (saving) return;

		const id = team?.id ?? 0;
		const n = name.trim();
		const b = bio.trim();
		const c = countryCode.trim();
		const i = imageUrl.trim();

		if (!n && !b && !c && !i) {
			toast.error('Please fill at least one field.');
			return;
		}

		if (i && !isLikelyUrl(i)) {
			toast.error('Image must be a valid URL.');
			return;
		}

		try {
			saving = true;
			await updateTeam(id, n, b, i, c);
			open = false;
			
			// Invalidate teams cache so the teams page updates
			queryClient.invalidateQueries({ queryKey: ['teams'] });
			
			dispatch('updated', {id: id})
			toast.success('Team updated.');
		} catch (err: any) {
			toast.error(err?.message ?? 'Failed to update team.');
		} finally {
			saving = false;
		}
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="w-full px-5 sm:max-w-[640px]">
		<Sheet.Header>
			<Sheet.Title>Update Team</Sheet.Title>
			<Sheet.Description>Change your team name, country, bio, and profile image.</Sheet.Description>
		</Sheet.Header>

		<form class="mt-3 space-y-6" onsubmit={onSave}>
			<Accordion.Root type="single" value="identity">
				<Accordion.Item value="identity">
					<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight">
						Team
					</Accordion.Trigger>
					<Accordion.Content>
						<div class="grid grid-cols-1 gap-4">
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

								<!-- Country Select -->
								<Select.Root
									type="single"
									bind:value={countryCode}
									items={filteredCountries}
									allowDeselect={true}
									onOpenChange={(open) => { if (!open) countrySearch = ''; }}
								>
									<Select.Trigger id="pf-country" class="w-full justify-between">
										{#if countryCode}
											{@const country = countryItems.find(c => c.value === countryCode)}
											<span class="flex items-center gap-2">
												{#if country?.iso2}
													<Icon
														icon={`circle-flags:${country.iso2.toLowerCase()}`}
														width="16"
														height="16"
													/>
												{/if}
												<span class="uppercase">{countryCode}</span>
											</span>
										{:else}
											<span class="text-muted-foreground">Select country</span>
										{/if}
									</Select.Trigger>
									<Select.Content sideOffset={4}>
										<div class="px-2 pb-2">
											<Input 
												placeholder="Search countries..." 
												bind:value={countrySearch}
												class="h-8"
											/>
										</div>
										{#each filteredCountries as item (item.value)}
											<Select.Item value={item.value} label={item.label}>
												{#snippet children({ selected })}
													<Icon
														icon={`circle-flags:${item.iso2.toLowerCase()}`}
														width="16"
														height="16"
													/>
													<span>{item.label} ({item.value})</span>
												{/snippet}
											</Select.Item>
										{/each}
										{#if filteredCountries.length === 0}
											<div class="px-2 py-6 text-center text-sm text-muted-foreground">
												No countries found
											</div>
										{/if}
									</Select.Content>
								</Select.Root>

								{#if team?.country && team.country.toUpperCase() !== countryCode}
									{@const teamCountry = countryItems.find(c => c.value === team.country.toUpperCase())}
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

							<div>
								<Label for="pf-image" class="mb-1 block">Image URL</Label>
								<Input
									id="pf-image"
									bind:value={imageUrl}
									placeholder={team?.image || 'https://.../avatar.png'}
								/>
								{#if team?.image && team.image !== imageUrl}
									<p class="text-muted-foreground mt-1 text-sm">Current image will be replaced</p>
								{/if}
								{#if imageUrl && isLikelyUrl(imageUrl)}
									<div class="mt-3">
										<Label class="mb-1 block">Preview</Label>
										<div>
											<Avatar src={imageUrl} class="h-24 w-24" />
										</div>
									</div>
								{:else if team?.image && isLikelyUrl(team.image)}
									<div class="mt-3">
										<Label class="mb-1 block">Current Image</Label>
										<div>
											<Avatar src={team.image} class="h-24 w-24" />
										</div>
									</div>
								{/if}
							</div>
						</div>
					</Accordion.Content>
				</Accordion.Item>
			</Accordion.Root>

			<div class="flex justify-end gap-2">
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
