<script lang="ts">
	/* UI */
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';
	import { Avatar } from 'flowbite-svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import countries from '$lib/data/countries.json';
	/* Your provided function */
	import { updateUser } from '$lib/user';

	let { open = $bindable(false), user } = $props<{
		open?: boolean;
		user?: { id: number; name?: string; country?: string; image?: string };
	}>();

	// Initialize form fields with current user data
	let name = $state('');
	let imageUrl = $state('');
	let countryCode = $state<string>('');
	let saving = $state(false);
	let confirmTeamPropagation = $state(true);

	type Country = { name: string; iso2: string; iso3?: string; emoji?: string };
	const countryItems = (countries as Country[])
		.map((c) => ({ value: c.iso2.toUpperCase(), label: c.name }))
		.sort((a, b) => a.label.localeCompare(b.label));

	// Watch for user changes and update form fields
	$effect(() => {
		if (user) {
			name = user.name ?? '';
			imageUrl = user.image ?? '';
			countryCode = user.country?.toUpperCase?.() ?? '';
		}
	});

	// Reset form when sheet opens
	$effect(() => {
		if (open && user) {
			name = user.name ?? '';
			imageUrl = user.image ?? '';
			countryCode = user.country?.toUpperCase?.() ?? '';
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

		const id = user?.id ?? 0;
		const n = name.trim();
		const c = countryCode.trim();
		const i = imageUrl.trim();

		if (!n && !c && !i) {
			toast.error('Please fill at least one field.');
			return;
		}

		if (i && !isLikelyUrl(i)) {
			toast.error('Image must be a valid URL.');
			return;
		}

		try {
			saving = true;
			await updateUser(id, n, c, i);
			open = false;
			toast.success('Profile updated.');
			// Reload page to show changes immediately
			setTimeout(() => {
				window.location.reload();
			}, 1000);
		} catch (err: any) {
			toast.error(err?.message ?? 'Failed to update profile.');
		} finally {
			saving = false;
		}
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="px-5 sm:max-w-[640px]">
		<Sheet.Header>
			<Sheet.Title>Update Profile</Sheet.Title>
			<Sheet.Description
				>Change your display name, nationality, and profile image.</Sheet.Description
			>
		</Sheet.Header>

		<form class="mt-3 space-y-6" onsubmit={onSave}>
			<Accordion.Root type="single" value="identity">
				<Accordion.Item value="identity">
					<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight">
						Identity
					</Accordion.Trigger>
					<Accordion.Content>
						<div class="grid grid-cols-1 gap-4">
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

								<!-- Country Select -->
								<Select.Root
									type="single"
									bind:value={countryCode}
									items={countryItems}
									allowDeselect={true}
								>
									<Select.Trigger id="pf-country" class="w-full justify-between">
										{#if countryCode}
											<span class="flex items-center gap-2">
												<Icon
													icon={`circle-flags:${countryCode.toLowerCase()}`}
													width="16"
													height="16"
												/>
												<span class="uppercase">{countryCode}</span>
											</span>
										{:else}
											<span class="text-muted-foreground">Select country</span>
										{/if}
									</Select.Trigger>
									<Select.Content sideOffset={4}>
										{#each countryItems as item (item.value)}
											<Select.Item value={item.value} label={item.label}>
												{#snippet children({ selected })}
													<Icon
														icon={`circle-flags:${item.value.toLowerCase()}`}
														width="16"
														height="16"
													/>
													<span>{item.label}</span>
												{/snippet}
											</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>

								{#if user?.country && user.country.toUpperCase() !== countryCode}
									<div class="text-muted-foreground mt-1 flex items-center gap-2 text-sm">
										<Icon
											icon={`circle-flags:${user.country.toLowerCase()}`}
											width="16"
											height="16"
										/>
										<span>Current: {user.country.toUpperCase()}</span>
									</div>
								{/if}
							</div>

							<div>
								<Label for="pf-image" class="mb-1 block">Image URL</Label>
								<Input
									id="pf-image"
									bind:value={imageUrl}
									placeholder={user?.image || 'https://…/avatar.png'}
								/>
								{#if user?.image && user.image !== imageUrl}
									<p class="text-muted-foreground mt-1 text-sm">Current image will be replaced</p>
								{/if}
								{#if imageUrl && isLikelyUrl(imageUrl)}
									<div class="mt-3">
										<Label class="mb-1 block">Preview</Label>
										<div class="rounded-lg border p-2">
											<Avatar src={imageUrl} class="h-24 w-24" />
										</div>
									</div>
								{:else if user?.image && isLikelyUrl(user.image)}
									<div class="mt-3">
										<Label class="mb-1 block">Current Image</Label>
										<div class="rounded-lg border p-2">
											<Avatar src={user.image} class="h-24 w-24" />
										</div>
									</div>
								{/if}
							</div>
						</div>
					</Accordion.Content>
				</Accordion.Item>

				<Accordion.Item value="advanced">
					<Accordion.Trigger class="cursor-pointer text-xl font-semibold tracking-tight">
						Advanced
					</Accordion.Trigger>
					<Accordion.Content>
						<div class="flex items-center gap-2">
							<Checkbox id="pf-propagate" bind:checked={confirmTeamPropagation} />
							<Label for="pf-propagate">Also update team display (if enabled by admins)</Label>
						</div>
					</Accordion.Content>
				</Accordion.Item>
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
