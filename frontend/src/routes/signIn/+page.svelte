<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { toast } from 'svelte-sonner';

	import { login } from '$lib/auth';
	import { link, push } from 'svelte-spa-router';
	import { user } from '@/stores/auth';
	import { loadUser } from '@/stores/auth';
	import { Lock } from '@lucide/svelte';

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

<div class="flex min-h-[80vh] items-center justify-center px-4 py-12">
	<Card.Root
		class="bg-card/50 mx-auto w-full max-w-md border-0 shadow-xl backdrop-blur-sm sm:max-w-[420px]"
	>
		<div class="p-8 pb-0">
			<Card.Header class="space-y-2 p-0 text-center">
				<div
					class="bg-primary/10 mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full"
				>
					<Lock class="text-primary h-6 w-6" />
				</div>
				<Card.Title class="text-2xl font-bold tracking-tight">Welcome back</Card.Title>
				<Card.Description class="text-base">Sign in to your account to continue</Card.Description>
			</Card.Header>
		</div>

		<form onsubmit={onSubmit} class="p-8 pt-6">
			<Card.Content class="space-y-6 p-0">
				<div class="space-y-4">
					<div class="space-y-2">
						<Label for="email" class="font-medium">Email</Label>
						<Input
							id="email"
							name="email"
							type="email"
							placeholder="name@example.com"
							bind:value={email}
							required
							class="bg-background/50"
						/>
					</div>

					<div class="space-y-2">
						<div class="flex items-center justify-between">
							<Label for="password" class="font-medium">Password</Label>
							<a
								use:link
								href="/forgot"
								class="text-muted-foreground hover:text-primary text-xs font-medium hover:underline"
							>
								Forgot password?
							</a>
						</div>
						<Input
							id="password"
							name="password"
							type="password"
							placeholder="••••••••"
							bind:value={password}
							required
							class="bg-background/50"
						/>
					</div>
				</div>

				<div class="flex items-center gap-2">
					<Checkbox id="remember" bind:checked={remember} />
					<Label for="remember" class="cursor-pointer text-sm font-normal">Remember me</Label>
				</div>

				{#if errorMsg}
					<div class="text-sm text-red-600 dark:text-red-400">{errorMsg}</div>
				{/if}

				<Button type="submit" class="w-full font-semibold shadow-sm" size="lg" disabled={loading}>
					{#if loading}
						Signing in...
					{:else}
						Sign in
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
					Don't have an account?{' '}
					<Button
						variant="link"
						class="text-primary h-auto p-0 font-semibold"
						type="button"
						onclick={goToSignUp}
					>
						Sign up
					</Button>
				</p>
			</Card.Footer>
		</form>
	</Card.Root>
</div>
