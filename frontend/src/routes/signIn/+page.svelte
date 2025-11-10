<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { toast } from 'svelte-sonner';
	import ErrorMessage from '$lib/components/ui/error-message.svelte';

	import { login } from '$lib/auth';
	import { link, push } from 'svelte-spa-router';
	import { user } from '@/stores/auth';
	import { loadUser } from '@/stores/auth';

	// --- State (Svelte 5 runes) ----------------------------------------------

	let email = $state('');
	let password = $state('');
	let remember = $state(false);

	let loading = $state(false);
	let errorMsg = $state<string | null>(null);

	// --- Helpers --------------------------------------------------------------

	function getErrorMessage(err: unknown): string {
		const defaultMessage = 'Login failed. Please try again.';

		if (!err || typeof err !== 'object') return defaultMessage;

		const maybeMessage = (err as { message?: string }).message;
		if (!maybeMessage) return defaultMessage;

		// Try to parse { "error": "..." } JSON in message
		try {
			const parsed = JSON.parse(maybeMessage);
			if (
				parsed &&
				typeof parsed === 'object' &&
				'error' in parsed &&
				typeof (parsed as any).error === 'string'
			) {
				return (parsed as any).error;
			}
		} catch {
			// Ignore JSON parse errors; fall back to raw message
		}

		return maybeMessage || defaultMessage;
	}

	function goToSignUp() {
		push('/signUp');
	}

	// --- Actions --------------------------------------------------------------

	async function onSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (loading) return;

		errorMsg = null;
		loading = true;

		try {
			const result = await login(email, password);

			if (result !== 'OK') {
				// If your login returns more info, you can improve this
				throw new Error('Login failed. Please try again.');
			}

			await loadUser();

			toast.success('Welcome back!');

			// Redirect based on team presence
			if ($user?.team_id) {
				push('/challenges');
			} else {
				push('/team');
			}
		} catch (err) {
			const message = getErrorMessage(err);
			errorMsg = message;
			toast.error(message);
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-svh flex items-center justify-center">
	<Card.Root class="mx-auto w-full max-w-sm flex flex-col">
		<Card.Header>
			<Card.Title>Welcome back hacker.</Card.Title>
			<Card.Description>Enter your email below to login to your account</Card.Description>
			<Card.Action>
				<Button
					variant="link"
					class="cursor-pointer"
					type="button"
					onclick={goToSignUp}
				>
					Sign Up
				</Button>
			</Card.Action>
		</Card.Header>

		<form onsubmit={onSubmit}>
			<Card.Content>
				<div class="flex flex-col gap-6">
					<div class="grid gap-2">
						<Label for="email">Email</Label>
						<Input
							id="email"
							name="email"
							type="email"
							placeholder="m@example.com"
							bind:value={email}
							required
						/>
					</div>

					<div class="grid gap-2">
						<div class="flex items-center">
							<Label for="password">Password</Label>
							<a
								use:link
								href="/forgot"
								class="ml-auto inline-block text-sm underline-offset-4 hover:underline"
							>
								Forgot your password?
							</a>
						</div>
						<Input
							id="password"
							name="password"
							type="password"
							placeholder="********"
							bind:value={password}
							required
						/>
					</div>

					<div class="flex select-none items-center gap-2 text-sm">
						<Checkbox id="remember" bind:checked={remember} />
						<Label for="remember">Remember me</Label>
					</div>

					<div class="min-h-5">
						{#if errorMsg}
							<ErrorMessage message={errorMsg} />
						{/if}
					</div>
				</div>
			</Card.Content>

			<Card.Footer class="flex-col gap-2">
				<Button type="submit" class="w-full cursor-pointer" disabled={loading}>
					{#if loading}
						Signing in...
					{:else}
						Sign in
					{/if}
				</Button>
			</Card.Footer>
		</form>
	</Card.Root>
</div>
