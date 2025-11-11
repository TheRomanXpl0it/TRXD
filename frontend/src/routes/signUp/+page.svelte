<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { register, type User } from '$lib/auth';
	import { toast } from 'svelte-sonner';
	import { link, push } from 'svelte-spa-router';
	import { user, loadUser } from '@/stores/auth';

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
			if ($user?.team_id) {
				push('/challenges');
			} else {
				push('/team');
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

<div class="flex min-h-full items-center justify-center py-12">  
    <Card.Root class="mx-auto w-full max-w-sm">
    	<Card.Header>
    		<Card.Title>Create your account</Card.Title>
    		<Card.Description>Join TRXD and start hacking.</Card.Description>
    		<Card.Action>
    			<div>
    				<Button variant="link" class="cursor-pointer" type="button" onclick={() => push('/signIn')}>
    					Sign in
    				</Button>
    			</div>
    		</Card.Action>
    	</Card.Header>
    
    	<!-- Wrap form so submit button works -->
    	<form onsubmit={onSubmit}>
    		<Card.Content>
    			<div class="flex flex-col gap-6">
    				<div class="grid gap-2">
    					<Label for="name">Username</Label>
    					<Input
    						id="name"
    						name="name"
    						type="text"
    						placeholder="Your username"
    						bind:value={name}
    						required
    					/>
    				</div>
    
    				<div class="grid gap-2">
    					<Label for="email">Email</Label>
    					<Input
    						id="email"
    						name="email"
    						type="email"
    						placeholder="name@email.com"
    						bind:value={email}
    						required
    					/>
    				</div>
    
    				<div class="grid gap-2">
    					<Label for="password">Password</Label>
    					<Input
    						id="password"
    						name="password"
    						type="password"
    						placeholder="********"
    						minlength={8}
    						bind:value={password}
    						required
    					/>
    					<p class="text-xs text-gray-500 dark:text-gray-400">At least 8 characters.</p>
    				</div>
    
    				<div class="grid gap-2">
    					<Label for="confirm">Confirm password</Label>
    					<Input
    						id="confirm"
    						name="confirm"
    						type="password"
    						placeholder="********"
    						minlength={8}
    						bind:value={confirm}
    						required
    					/>
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
    				{#if loading}
    					<span class="inline-flex items-center gap-2"><Spinner /> Signing up...</span>
    				{:else}
    					Sign up
    				{/if}
    			</Button>
    		</Card.Footer>
    	</form>
    </Card.Root>
</div>
