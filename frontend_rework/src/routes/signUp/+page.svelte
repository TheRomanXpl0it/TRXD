<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Checkbox } from "$lib/components/ui/checkbox/index.js";
  import * as Dialog from "$lib/components/ui/dialog/index.js";

  import Spinner from "$lib/components/ui/spinner/spinner.svelte";
  import { register, type User } from "$lib/auth";
  import { toast } from "svelte-sonner";
  import { link, push } from "svelte-spa-router";

  let name = "";
  let email = "";
  let password = "";
  let confirm = "";
  let agree = false;

  let loading = false;
  let errorMsg: string | null = null;

  // Terms modal
  let termsOpen = false;

  function validate(): string | null {
    if (!name.trim()) return "Please enter your name.";
    if (!email.trim()) return "Please enter your email.";
    if (password.length < 8) return "Password must be at least 8 characters.";
    if (password !== confirm) return "Passwords do not match.";
    if (!agree) return "You must agree to the Terms.";
    return null;
  }

  function getRedirect(): string {
    const q = new URLSearchParams(location.search);
    return q.get("redirect") || "/challenges";
  }

  async function onSubmit(e: Event) {
    e.preventDefault();
    errorMsg = validate();
    if (errorMsg) return;

    loading = true;
    try {
      const _user: User = await register(email,password,name);
      loading = false;
      toast.success("Welcome aboard!");
      await new Promise((r) => setTimeout(r, 500)) // wait a bit
      const q = new URLSearchParams(location.search);
      const dest = q.get("redirect") || "/challenges";
      window.location.replace(dest); // or: window.location.assign(dest)
    } catch (err: any) {
      errorMsg = err?.message ?? "Registration failed. Please try again.";
      toast.error(errorMsg as string);
    }
  }

  function acceptTerms() {
    agree = true;
    termsOpen = false;
  }
</script>

<Card.Root class="w-full max-w-sm mx-auto mt-50">
  <Card.Header>
    <Card.Title>Create your account</Card.Title>
    <Card.Description>Join TRXD and start hacking.</Card.Description>
    <Card.Action>
      <div>
        <Button variant="link" class="cursor-pointer" type="button" onclick={() => push("/signIn")}>
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
          <Input id="name" name="name" type="text" placeholder="Your username" bind:value={name} required />
        </div>

        <div class="grid gap-2">
          <Label for="email">Email</Label>
          <Input id="email" name="email" type="email" placeholder="name@email.com" bind:value={email} required />
        </div>

        <div class="grid gap-2">
          <Label for="password">Password</Label>
          <Input id="password" name="password" type="password" placeholder="********" minlength={8} bind:value={password} required />
          <p class="text-xs text-gray-500 dark:text-gray-400">At least 8 characters.</p>
        </div>

        <div class="grid gap-2">
          <Label for="confirm">Confirm password</Label>
          <Input id="confirm" name="confirm" type="password" placeholder="********" minlength={8} bind:value={confirm} required />
        </div>

        <div class="flex items-start gap-2 text-sm select-none mb-5">
            <Checkbox id="agree" class="mt-0.5" bind:checked={agree} />
          <Label for="agree" class="!mb-0">
            I agree to the&nbsp;
            <button type="button" class="underline underline-offset-4 hover:opacity-80 cursor-pointer" onclick={() => (termsOpen = true)}>
              Terms &amp; Conditions
            </button>
          </Label>
        </div>

        {#if errorMsg}
          <p class="text-red-600 text-sm">{errorMsg}</p>
        {/if}
      </div>
    </Card.Content>

    <Card.Footer class="flex-col gap-2">
      <Button type="submit" class="w-full cursor-pointer" disabled={loading}>
        {#if loading}
          <span class="inline-flex items-center gap-2"><Spinner /> Signing up…</span>
        {:else}
          Sign up
        {/if}
      </Button>
    </Card.Footer>
  </form>
</Card.Root>

<!-- Terms & Conditions Modal -->
<Dialog.Root bind:open={termsOpen}>
  <Dialog.Overlay />
  <Dialog.Content class="sm:max-w-[720px]">
    <Dialog.Header class="pb-2">
      <Dialog.Title>Terms &amp; Conditions</Dialog.Title>
      <Dialog.Description class="sr-only">
        Please read and accept the terms to continue.
      </Dialog.Description>
    </Dialog.Header>

    <div class="max-h-[60vh] overflow-y-auto pr-2 space-y-4 text-sm text-muted-foreground">
      <p><b>1. Introduction.</b> These Terms govern your use of the TRXD platform and services.</p>
      <p><b>2. Eligibility.</b> You must comply with applicable laws and be authorized to use this service.</p>
      <p><b>3. Account.</b> You’re responsible for maintaining the confidentiality of your credentials.</p>
      <p><b>4. Acceptable Use.</b> No abuse, unauthorized access attempts, or harmful activities.</p>
      <p><b>5. Content.</b> You’re responsible for the content you submit. Do not upload unlawful content.</p>
      <p><b>6. Privacy.</b> We process data as described in our Privacy Policy.</p>
      <p><b>7. Termination.</b> We may suspend or terminate accounts for violations.</p>
      <p><b>8. Disclaimer.</b> Service is provided “as is” without warranties.</p>
      <p><b>9. Limitation of Liability.</b> To the maximum extent permitted by law, our liability is limited.</p>
      <p><b>10. Changes.</b> We may update these Terms; continued use means acceptance.</p>
      <p><b>11. Contact.</b> For questions, contact support@theromanXpl0.it.</p>
      <p><b>12. Thanks Dario.</b> Always thank Dario.</p>
    </div>

    <div class="mt-4 flex justify-end gap-2">
      <Dialog.Close>
        <Button variant="outline" class="cursor-pointer">Close</Button>
      </Dialog.Close>
      <Button onclick={acceptTerms} class="cursor-pointer">I agree</Button>
    </div>
  </Dialog.Content>
</Dialog.Root>
