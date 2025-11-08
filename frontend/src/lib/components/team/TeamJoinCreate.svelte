<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';

	// Images (transparent PNGs)
	import createTeamImg from '$lib/assets/createTeam.png?url';
	import joinTeamImg from '$lib/assets/joinTeam.png?url';
	import { joinTeam, createTeam } from '@/team';
	import { toast } from 'svelte-sonner';
	
	import { createEventDispatcher } from 'svelte';
        
    const dispatch = createEventDispatcher();

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
			joinName = '';
			joinPassword = '';

			dispatch('joined', {joinName: joinName})
			toast.success('Team Joined, welcome aboard!');
		} catch (err: any) {
			joinError = err?.message ?? 'Failed to join team.';
			toast.error(joinError??"Error");
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
			registerName = '';
			registerPassword = '';

			dispatch('created', {registerName: registerName})
			toast.success('Team Created!');
		} catch (err: any) {
			registerError = err?.message ?? 'Failed to register team.';
			toast.error(registerError as string);
		} finally {
			registerLoading = false;
		}
	}
</script>

<div class="mt-5 mb-10 flex flex-col items-stretch justify-center gap-6 xl:flex-row">
	<!-- Join card -->
	<Card.Root class="flex w-full flex-col p-4 sm:w-[28rem] sm:p-5">
		<Card.Header class="flex-1 pb-3">
			<Card.Title class="Title">
				<h3 class="text-lg font-bold">Join a Team</h3>
				<p class="text-sm italic text-gray-400">Join an already existing team</p>
			</Card.Title>
			<Card.Content class="px-0">
				<img
					src={joinTeamImg}
					alt="Join a Team"
					class="mx-auto mt-2 h-auto w-full max-h-48 max-w-xs rounded-md object-contain dark:invert sm:max-h-64"
				/>
			</Card.Content>
		</Card.Header>
		<Button variant="default" class="cursor-pointer" onclick={() => (joinOpen = true)}>Join</Button>
	</Card.Root>

	<!-- Create card -->
	<Card.Root class="flex w-full flex-col p-4 sm:w-[28rem] sm:p-5">
		<Card.Header class="flex-1 pb-3">
			<Card.Title class="Title">
				<h3 class="text-lg font-bold">Create a Team</h3>
				<p class="text-sm italic text-gray-400">Create a new team from scratch</p>
			</Card.Title>
			<Card.Content class="px-0">
				<img
					src={createTeamImg}
					alt="Create a Team"
					class="mx-auto mt-2 h-auto w-full max-h-48 max-w-xs rounded-md object-contain p-6 dark:invert sm:max-h-64 sm:p-8"
				/>
			</Card.Content>
		</Card.Header>
		<Button variant="default" class="cursor-pointer" onclick={() => (registerOpen = true)}>
			Create
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
				<Label for="team-name">Team name</Label>
				<Input id="team-name" placeholder="e.g. ZeroDayCats" bind:value={joinName} required />
			</div>

			<div class="grid gap-2">
				<Label for="team-pass">Team password</Label>
				<Input
					id="team-pass"
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
				<Label for="team-name">Team name</Label>
				<Input id="team-name" placeholder="TRX" bind:value={registerName} required />
			</div>

			<div class="grid gap-2">
				<Label for="team-pass">Team password</Label>
				<Input
					id="team-pass"
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
