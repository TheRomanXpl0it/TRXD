<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { toast } from 'svelte-sonner';

	import { login, type User } from '$lib/auth';
	import { link, push } from 'svelte-spa-router';
	import { user } from '@/stores/auth';

	$: if ($user) {
		toast.success('already logged in!');
		push('/challenges');
	}
	let email = '';
	let password = '';
	let remember = false;

	let loading = false;
	let errorMsg: string | null = null;

	function getRedirect(): string {
		const q = new URLSearchParams(location.search);
		return q.get('redirect') || '/challenges';
	}

	async function onSubmit(e: Event) {
		e.preventDefault();
		errorMsg = null;
		loading = true;
		try {
			const result = await login(email, password);
			if (result !== 'OK') {
				throw new Error('Login failed. Please try again.');
			}
			loading = false;
			toast.success('Welcome back!');
			await new Promise((r) => setTimeout(r, 500)); // wait a bit
			const q = new URLSearchParams(location.search);
			const dest = q.get('redirect') || '/challenges';
			window.location.replace(dest); // or: window.location.assign(dest)
		} catch (err: any) {
			errorMsg = err?.message ?? 'Login failed. Please try again.';
			loading = false;
			toast.error(errorMsg as string);
		}
	}
</script>

<Card.Root class="mt-50 mx-auto flex h-full w-full max-w-sm flex-col">
	<Card.Header>
		<Card.Title>Welcome back hacker.</Card.Title>
		<Card.Description>Enter your email below to login to your account</Card.Description>
		<Card.Action>
			<Button variant="link" class="cursor-pointer" type="button" onclick={() => push('/signup')}>
				Sign Up
			</Button>
		</Card.Action>
	</Card.Header>

	<!-- Wrap Content + Footer in one form so the submit button works -->
	<form on:submit|preventDefault={onSubmit}>
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

				<!-- Remember me (native checkbox to avoid extra deps) -->
				<div class="mb-5 flex select-none items-center gap-2 text-sm">
					<Checkbox id="terms" bind:checked={remember} />
					<Label for="terms">Remember me</Label>
				</div>
			</div>
		</Card.Content>

		<Card.Footer class="flex-col gap-2">
			<Button type="submit" class="w-full cursor-pointer" disabled={loading}>
				{#if loading}Signing in…{:else}Sign in{/if}
			</Button>
		</Card.Footer>
	</form>
</Card.Root>
