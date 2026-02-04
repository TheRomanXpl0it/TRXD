<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import CountrySelect from '$lib/components/ui/country-select.svelte';
	import { updateUser } from '$lib/user';
	import { showSuccess, showError } from '$lib/utils/toast';
	import { getCountryByIso3 } from '$lib/utils/countries';
	import Icon from '@iconify/svelte';
	import GeneratedAvatar from '$lib/components/ui/avatar/generated-avatar.svelte';

	let {
		open = $bindable(false),
		user,
		onupdated
	} = $props<{
		open?: boolean;
		user?: { id: number; name?: string; country?: string };
		onupdated?: (detail: { id: number }) => void;
	}>();

	let name = $state('');
	let countryCode = $state('');
	let saving = $state(false);

	$effect(() => {
		if (user) {
			name = user.name ?? '';
			countryCode = user.country?.toUpperCase?.() ?? '';
		}
	});

	async function onSave(e: Event) {
		e.preventDefault();
		if (saving) return;

		const trimmedName = name.trim();
		const trimmedCountry = countryCode.trim();

		if (!trimmedName && !trimmedCountry) {
			showError(null, 'Please fill at least one field.');
			return;
		}

		try {
			saving = true;
			await updateUser(user?.id ?? 0, trimmedName, trimmedCountry);
			open = false;
			onupdated?.({ id: user?.id ?? 0 });
			showSuccess('Profile updated.');
		} catch (err: any) {
			showError(err, 'Failed to update profile.');
		} finally {
			saving = false;
		}
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
					<Sheet.Title class="text-xl font-bold">Edit Profile</Sheet.Title>
					<Sheet.Description class="text-muted-foreground/80 mt-1">
						Update your personal details.
					</Sheet.Description>
				</div>
			</div>
		</div>

		<form class="mt-3 space-y-6" onsubmit={onSave}>
			<div class="space-y-4">
				<div>
					<Label for="pf-name" class="mb-1 block">Display name</Label>
					<Input id="pf-name" bind:value={name} placeholder={user?.name || 'Your name'} />
					{#if user?.name && user.name !== name}
						<p class="text-muted-foreground mt-1 text-sm">
							Current: {user.name}
						</p>
					{/if}
				</div>

				<div>
					<Label for="pf-country" class="mb-1 block">Nationality</Label>
					<CountrySelect id="pf-country" bind:value={countryCode} />
					{#if user?.country && user.country.toUpperCase() !== countryCode}
						{@const userCountry = getCountryByIso3(user.country)}
						<div class="text-muted-foreground mt-1 flex items-center gap-2 text-sm">
							{#if userCountry?.iso2}
								<Icon
									icon={`circle-flags:${userCountry.iso2.toLowerCase()}`}
									width="16"
									height="16"
								/>
							{/if}
							<span>Current: {user.country.toUpperCase()}</span>
						</div>
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
