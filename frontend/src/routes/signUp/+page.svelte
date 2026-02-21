<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { register } from '$lib/auth';
	import type { User } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { authState, loadUser } from '@/stores/auth';
	import { UserPlus } from '@lucide/svelte';

	let name = '';
	let email = '';
	let password = '';
	let confirm = '';

	let loading = false;
	let errorMsg: string | null = null;

	function validate(): string | null {
		if (!name.trim()) return 'Please enter your name.';
		if (!email.trim()) return 'Please enter your email.';
		if (password.length < 8) return 'Password must be at least 8 characters.';
		if (password !== confirm) return 'Passwords do not match.';
		return null;
	}

	async function onSubmit(e: Event) {
		e.preventDefault();
		errorMsg = validate();
		if (errorMsg) return;

		loading = true;
		try {
			const _user: User = await register(email, password, name);
			await loadUser();
			loading = false;
			toast.success('Welcome aboard!');
			// Check if user has a team, redirect accordingly
			if (authState.user?.team_id) {
				goto('/challenges');
			} else {
				goto('/team');
			}
		} catch (err: any) {
			// Extract error message from JSON response if present
			let message = 'Registration failed. Please try again.';
			if (err?.message) {
				try {
					const parsed = JSON.parse(err.message);
					message = parsed.error || message;
				} catch {
					message = err.message;
				}
			}
			errorMsg = message;
			toast.error(message);
			loading = false;
		}
	}
</script>

<div class="flex min-h-[80vh] items-center justify-center px-4 py-12">
	<Card.Root
		class="bg-card/50 mx-auto w-full max-w-md border-0 shadow-xl backdrop-blur-sm sm:max-w-[450px]"
	>
		<div class="p-8 pb-0">
			<Card.Header class="space-y-2 p-0 text-center">
				<div
					class="bg-primary/10 mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full"
				>
					<UserPlus class="text-primary h-6 w-6" />
				</div>
				<Card.Title class="text-2xl font-bold tracking-tight">Create an account</Card.Title>
				<Card.Description class="text-base">Join TRXD and start hacking</Card.Description>
			</Card.Header>
		</div>

		<form onsubmit={onSubmit} class="p-8 pt-6">
			<Card.Content class="space-y-6 p-0">
				<div class="space-y-4">
					<div class="space-y-2">
						<Label for="name" class="font-medium">Username</Label>
						<Input
							id="name"
							name="name"
							type="text"
							placeholder="Your username"
							bind:value={name}
							required
							class="bg-background/50"
						/>
					</div>

					<div class="space-y-2">
						<Label for="email" class="font-medium">Email</Label>
						<Input
							id="email"
							name="email"
							type="email"
							placeholder="name@email.com"
							bind:value={email}
							required
							class="bg-background/50"
						/>
					</div>

					<div class="space-y-2">
						<Label for="password" class="font-medium">Password</Label>
						<Input
							id="password"
							name="password"
							type="password"
							placeholder="********"
							minlength={8}
							bind:value={password}
							required
							class="bg-background/50"
						/>
						<p class="text-xs text-gray-500 dark:text-gray-400">At least 8 characters.</p>
					</div>

					<div class="space-y-2">
						<Label for="confirm" class="font-medium">Confirm password</Label>
						<Input
							id="confirm"
							name="confirm"
							type="password"
							placeholder="********"
							minlength={8}
							bind:value={confirm}
							required
							class="bg-background/50"
						/>
					</div>
				</div>

				{#if errorMsg}
					<div class="text-sm text-red-600 dark:text-red-400">{errorMsg}</div>
				{/if}

				<Button type="submit" class="w-full font-semibold shadow-sm" size="lg" disabled={loading}>
					{#if loading}
						<span class="inline-flex items-center gap-2"><Spinner /> Signing up...</span>
					{:else}
						Sign up
					{/if}
				</Button>
			</Card.Content>

			<Card.Footer class="text-muted-foreground mt-6 flex flex-col gap-4 p-0 text-center text-sm">
				<div class="flex w-full items-center gap-4">
					<span class="bg-border h-px flex-1"></span>
					<span class="text-muted-foreground text-xs uppercase">Or</span>
					<span class="bg-border h-px flex-1"></span>
				</div>
				<p>
					Already have an account?{' '}
					<Button
						variant="link"
						class="text-primary h-auto p-0 font-semibold"
						type="button"
						onclick={() => goto('/signIn')}
					>
						Sign in
					</Button>
				</p>
			</Card.Footer>
		</form>
	</Card.Root>
</div>
