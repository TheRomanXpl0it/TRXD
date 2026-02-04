<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';

	import createTeamImg from '$lib/assets/createTeam.webp?url';
	import joinTeamImg from '$lib/assets/joinTeam.webp?url';
	import { joinTeam, createTeam } from '@/team';
	import { toast } from 'svelte-sonner';

	let { onjoined, oncreated } = $props<{
		onjoined?: (detail: { joinName: string }) => void;
		oncreated?: (detail: { registerName: string }) => void;
	}>();

	// Modal state
	let joinOpen = $state(false);
	let joinName = $state('');
	let joinPassword = $state('');
	let joinLoading = $state(false);
	let joinError: string | null = $state(null);

	let registerOpen = $state(false);
	let registerName = $state('');
	let registerPassword = $state('');
	let confirmRegisterPassword = $state('');
	let registerLoading = $state(false);
	let registerError: string | null = $state(null);

	async function onJoinSubmit(e: Event) {
		e.preventDefault();
		joinError = null;

		if (!joinName.trim() || !joinPassword.trim()) {
			joinError = 'Please fill in both fields.';
			return;
		}

		joinLoading = true;
		try {
			const result = await joinTeam(joinName, joinPassword);
			// Success case - api function only returns on success (200)
			joinOpen = false;
			// reset fields after close
			const teamName = joinName;
			joinName = '';
			joinPassword = '';

			onjoined?.({ joinName: teamName });
			toast.success('Team Joined, welcome aboard!');
		} catch (err: any) {
			joinError = err?.message ?? 'Failed to join team.';
			toast.error(joinError ?? 'Error');
		} finally {
			joinLoading = false;
		}
	}

	async function onRegisterSubmit(e: Event) {
		e.preventDefault();
		registerError = null;

		if (!registerName.trim() || !registerPassword.trim() || !confirmRegisterPassword.trim()) {
			registerError = 'Please fill all fields.';
			return;
		}

		if (registerPassword !== confirmRegisterPassword) {
			registerError = 'Passwords do not match.';
			toast.error(registerError);
			return;
		}

		if (registerPassword.length < 8) {
			registerError = 'Password must be at least 8 characters.';
			toast.error(registerError);
			return;
		}

		registerLoading = true;
		try {
			const result = await createTeam(registerName, registerPassword);
			// Success case - api function only returns on success (200)
			registerOpen = false;
			// reset fields after close
			const teamName = registerName;
			registerName = '';
			registerPassword = '';

			oncreated?.({ registerName: teamName });
			toast.success('Team Created!');
		} catch (err: any) {
			registerError = err?.message ?? 'Failed to register team.';
			toast.error(registerError as string);
		} finally {
			registerLoading = false;
		}
	}
</script>

<div class="mb-12 mt-8 flex flex-col items-center justify-center gap-8 xl:flex-row xl:gap-6">
	<!-- Join card -->
	<Card.Root class="flex w-full flex-col p-6 transition-shadow hover:shadow-lg sm:w-[32rem] sm:p-8">
		<Card.Header class="flex-1 space-y-3 pb-4">
			<Card.Title class="Title">
				<h2 class="text-2xl font-bold tracking-tight">Join a Team</h2>
				<p class="text-muted-foreground mt-2 text-base">Collaborate with an existing team</p>
			</Card.Title>
			<Card.Content class="px-0 pt-4">
				<div class="bg-muted/30 border-border/50 mx-auto rounded-lg border p-8">
					<img
						src={joinTeamImg}
						alt="Join a Team"
						class="mx-auto h-auto max-h-56 w-full max-w-xs object-contain dark:invert"
					/>
				</div>
			</Card.Content>
		</Card.Header>
		<Button
			variant="outline"
			size="lg"
			class="w-full cursor-pointer text-base font-semibold"
			onclick={() => (joinOpen = true)}>Join Team</Button
		>
	</Card.Root>

	<!-- OR divider -->
	<div class="flex items-center justify-center xl:flex-col">
		<div class="bg-border h-px w-16 xl:h-16 xl:w-px"></div>
		<span class="text-muted-foreground px-4 text-sm font-medium xl:px-0 xl:py-4">OR</span>
		<div class="bg-border h-px w-16 xl:h-16 xl:w-px"></div>
	</div>

	<!-- Create card -->
	<Card.Root class="flex w-full flex-col p-6 transition-shadow hover:shadow-lg sm:w-[32rem] sm:p-8">
		<Card.Header class="flex-1 space-y-3 pb-4">
			<Card.Title class="Title">
				<h2 class="text-2xl font-bold tracking-tight">Create a Team</h2>
				<p class="text-muted-foreground mt-2 text-base">Start fresh with a new team</p>
			</Card.Title>
			<Card.Content class="px-0 pt-4">
				<div class="bg-muted/30 border-border/50 mx-auto rounded-lg border p-8">
					<img
						src={createTeamImg}
						alt="Create a Team"
						class="mx-auto h-auto max-h-56 w-full max-w-xs object-contain dark:invert"
					/>
				</div>
			</Card.Content>
		</Card.Header>
		<Button
			variant="default"
			size="lg"
			class="w-full cursor-pointer text-base font-semibold"
			onclick={() => (registerOpen = true)}
		>
			Create Team
		</Button>
	</Card.Root>
</div>

<!-- Join Team Modal -->
<Dialog.Root bind:open={joinOpen}>
	<Dialog.Overlay />
	<Dialog.Content class="sm:max-w-[480px]">
		<Dialog.Header>
			<Dialog.Title>Join a Team</Dialog.Title>
			<Dialog.Description>Enter the team name and password to join.</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={onJoinSubmit} class="mt-2 space-y-4">
			<div class="grid gap-2">
				<Label for="join-name">Team name</Label>
				<Input id="join-name" placeholder="e.g. ZeroDayCats" bind:value={joinName} required />
			</div>

			<div class="grid gap-2">
				<Label for="join-pass">Team password</Label>
				<Input
					id="join-pass"
					type="password"
					placeholder="••••••"
					bind:value={joinPassword}
					required
				/>
			</div>

			<div class="mt-4 flex justify-end gap-2">
				<Dialog.Close>
					<Button variant="outline" type="button" class="cursor-pointer">Cancel</Button>
				</Dialog.Close>
				<Button type="submit" disabled={joinLoading} class="cursor-pointer">
					{#if joinLoading}Joining...{:else}Join{/if}
				</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>

<!-- Register Team Modal -->
<Dialog.Root bind:open={registerOpen}>
	<Dialog.Overlay />
	<Dialog.Content class="sm:max-w-[480px]">
		<Dialog.Header>
			<Dialog.Title>Create a Team</Dialog.Title>
			<Dialog.Description>Enter the team name and password to register.</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={onRegisterSubmit} class="mt-2 space-y-4">
			<div class="grid gap-2">
				<Label for="reg-name">Team name</Label>
				<Input id="reg-name" placeholder="TRX" bind:value={registerName} required />
			</div>

			<div class="grid gap-2">
				<Label for="reg-pass">Team password</Label>
				<Input
					id="reg-pass"
					type="password"
					placeholder="••••••"
					bind:value={registerPassword}
					required
				/>
			</div>

			<div class="grid gap-2">
				<Label for="confirm-pass">Confirm password</Label>
				<Input
					id="confirm-pass"
					type="password"
					placeholder="••••••"
					bind:value={confirmRegisterPassword}
					required
				/>
			</div>

			<div class="mt-4 flex justify-end gap-2">
				<Dialog.Close>
					<Button variant="outline" type="button" class="cursor-pointer">Cancel</Button>
				</Dialog.Close>
				<Button type="submit" disabled={registerLoading} class="cursor-pointer">
					{#if registerLoading}Creating...{:else}Create{/if}
				</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
