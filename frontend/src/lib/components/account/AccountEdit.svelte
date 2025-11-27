<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import { Avatar } from 'flowbite-svelte';
	import CountrySelect from '$lib/components/ui/country-select.svelte';
	import { updateUser } from '$lib/user';
	import { isValidUrl } from '$lib/utils/validation';
	import { showSuccess, showError } from '$lib/utils/toast';
	import { getCountryByIso3 } from '$lib/utils/countries';
	import Icon from '@iconify/svelte';

	let {
		open = $bindable(false),
		user,
		onupdated
	} = $props<{
		open?: boolean;
		user?: { id: number; name?: string; country?: string; image?: string };
		onupdated?: (detail: { id: number }) => void;
	}>();

	let name = $state('');
	let imageUrl = $state('');
	let countryCode = $state('');
	let saving = $state(false);

	$effect(() => {
		if (user) {
			name = user.name ?? '';
			imageUrl = user.image ?? '';
			countryCode = user.country?.toUpperCase?.() ?? '';
		}
	});

	async function onSave(e: Event) {
		e.preventDefault();
		if (saving) return;

		const trimmedName = name.trim();
		const trimmedCountry = countryCode.trim();
		const trimmedImage = imageUrl.trim();

		if (!trimmedName && !trimmedCountry && !trimmedImage) {
			showError(null, 'Please fill at least one field.');
			return;
		}

		if (trimmedImage && !isValidUrl(trimmedImage)) {
			showError(null, 'Image must be a valid URL.');
			return;
		}

		try {
			saving = true;
			await updateUser(user?.id ?? 0, trimmedName, trimmedCountry, trimmedImage);
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
		<form class="mt-3 space-y-6" onsubmit={onSave}>
			<div class="space-y-4">
				<h2 class="text-xl font-semibold tracking-tight">Identity</h2>
				
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

				<div class="flex items-start gap-4">
					<div class="shrink-0">
						{#if imageUrl && isValidUrl(imageUrl)}
							<Avatar src={imageUrl} class="h-24 w-24" />
						{:else if user?.image && isValidUrl(user.image)}
							<Avatar src={user.image} class="h-24 w-24" />
						{:else}
							<div class="flex h-24 w-24 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800">
								<span class="text-xs text-gray-400">No image</span>
							</div>
						{/if}
					</div>

					<div class="flex-1 min-w-0">
						<Label for="pf-image" class="mb-1 block">Image URL</Label>
						<Input
							id="pf-image"
							bind:value={imageUrl}
							placeholder={user?.image || 'https://.../avatar.png'}
						/>
						{#if user?.image && user.image !== imageUrl}
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
