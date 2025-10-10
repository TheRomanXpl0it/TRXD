<script lang="ts">
  import { Card, Label, Input, Checkbox } from "flowbite-svelte";
  import { Button } from "@/components/ui/button"
  import { login, type User } from "$lib/auth";

  let email = "";
  let password = "";
  let remember = false;

  let loading = false;
  let errorMsg: string | null = null;

  async function onSubmit(e: Event) {
    e.preventDefault();
    errorMsg = null;
    loading = true;
    try {
      const user: User = await login(email, password);

      if (remember) {
        localStorage.setItem("remember_me", "1");
      } else {
        localStorage.removeItem("remember_me");
      }

      location.assign("/challenges"); // change to your desired route
    } catch (err) {
      errorMsg = err instanceof Error ? err.message : "Login failed. Please try again.";
    } finally {
      loading = false;
    }
  }
</script>

<div class="flex flex-col w-full h-full items-center justify-center mt-50">
  <Card class="p-4 sm:p-6 md:p-8">
    <form class="flex flex-col space-y-6" on:submit|preventDefault={onSubmit}>
      <h3 class="text-xl font-medium text-gray-900 dark:text-white">Welcome back hacker.</h3>

      <Label class="space-y-2">
        <span>Email</span>
        <Input type="email" name="email" bind:value={email} placeholder="name@email.com" required />
      </Label>

      <Label class="space-y-2">
        <span>Your password</span>
        <Input type="password" name="password" bind:value={password} placeholder="•••••" required />
      </Label>

      <div class="flex items-start gap-2">
        <Checkbox bind:checked={remember}>Remember me</Checkbox>
        <a href="/" class="text-primary-700 dark:text-primary-500 ms-auto text-sm hover:underline">Lost password?</a>
      </div>

      {#if errorMsg}
        <p class="text-red-600 text-sm">{errorMsg}</p>
      {/if}

      <Button type="submit" class="w-full" disabled={loading}>
        {#if loading}Signing in…{/if}
        {#if !loading}Login to your account{/if}
      </Button>

      <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
        Not registered? <a href="/" class="text-primary-700 dark:text-primary-500 hover:underline">Create account</a>
      </div>
    </form>
  </Card>
</div>
