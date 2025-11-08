<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { toast } from 'svelte-sonner';
	import ErrorMessage from '$lib/components/ui/error-message.svelte';

	import { login, type User } from '$lib/auth';
	import { link, push } from 'svelte-spa-router';
	import { user, loadUser } from '@/stores/auth';
	let email = '';
	let password = '';
	let remember = false;

	let loading = false;
	let errorMsg: string | null = null;

	async function onSubmit(e: Event) {
		e.preventDefault();
		errorMsg = null;
		loading = true;
		try {
			const result = await login(email, password);
			if (result !== 'OK') {
				throw new Error('Login failed. Please try again.');
			}
			await loadUser();
			loading = false;
			toast.success('Welcome back!');
			// Check if user has a team, redirect accordingly
			if ($user?.team_id) {
				push('/challenges');
			} else {
				push('/team');
			}
		} catch (err: any) {
			// TODO: Refactor this
			// Extract error message from JSON response if present
			let message = 'Login failed. Please try again.';
			if (err?.message) {
				try {
					const parsed = JSON.parse(err.message);
					message = parsed.error || message;
				} catch {
					message = err.message;
				}
			}
			errorMsg = message;
			loading = false;
			toast.error(message);
		}
	}
</script>

<div class="min-h-svh flex items-center justify-center">
	<Card.Root class="mx-auto w-full max-w-sm flex flex-col">
		<Card.Header>
			<Card.Title>Welcome back hacker.</Card.Title>
			<Card.Description>Enter your email below to login to your account</Card.Description>
			<Card.Action>
				<Button variant="link" class="cursor-pointer" type="button" onclick={() => push('/signUp')}>
					Sign Up
				</Button>
			</Card.Action>
		</Card.Header>

		<form on:submit|preventDefault={onSubmit}>
			<Card.Content>
				<div class="flex flex-col gap-6">
					<div class="grid gap-2">
						<Label for="email">Email</Label>
						<Input id="email" name="email" type="email" placeholder="m@example.com" bind:value={email} required />
					</div>

					<div class="grid gap-2">
						<div class="flex items-center">
							<Label for="password">Password</Label>
							<a use:link href="/forgot" class="ml-auto inline-block text-sm underline-offset-4 hover:underline">
								Forgot your password?
							</a>
						</div>
						<Input id="password" name="password" type="password" placeholder="********" bind:value={password} required />
					</div>

					<div class="flex select-none items-center gap-2 text-sm">
						<Checkbox id="terms" bind:checked={remember} />
						<Label for="terms">Remember me</Label>
					</div>

					<div class="min-h-5">
						{#if errorMsg}
							<p class="text-sm text-red-600 dark:text-red-400">{errorMsg}</p>
						{/if}
					</div>
				</div>
			</Card.Content>

			<Card.Footer class="flex-col gap-2">
				<Button type="submit" class="w-full cursor-pointer" disabled={loading}>
					{#if loading}Signing in...{:else}Sign in{/if}
				</Button>
			</Card.Footer>
		</form>
	</Card.Root>
</div>

