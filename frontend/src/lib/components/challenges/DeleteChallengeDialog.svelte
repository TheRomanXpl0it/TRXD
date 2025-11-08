<script lang="ts">
  import * as Dialog from '$lib/components/ui/dialog/index.js';
  import { Button } from '@/components/ui/button';
  import { Input } from '$lib/components/ui/input/index.js';
  import Label from '$lib/components/ui/label/label.svelte';
  import { Spinner } from '$lib/components/ui/spinner/index.js';
  import { createEventDispatcher } from 'svelte';



  let { open = $bindable(false), toDelete = $bindable(null), deleting = $bindable(false) } = $props<{ open?: boolean; toDelete: any, deleting: boolean }>();

  const dispatch = createEventDispatcher();

  let confirmationText = $state('');

  const expectedText = $derived(`Yes, I want to delete '${toDelete?.name ?? ''}'`);
  const canDelete = $derived(confirmationText === expectedText);

  // Reset confirmation text when dialog opens/closes
  $effect(() => {
    if (!open) {
      confirmationText = '';
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Overlay />
  <Dialog.Content class="sm:max-w-[520px]">
    <Dialog.Header class="pb-2">
      <Dialog.Title>Delete challenge?</Dialog.Title>
      <Dialog.Description>
        You're about to permanently delete <b>{toDelete?.name ?? 'this challenge'}</b>. This action cannot be undone.
      </Dialog.Description>
    </Dialog.Header>

    <div class="mt-4 space-y-3">
      <div class="space-y-2 text-sm text-gray-600 dark:text-gray-300">
        <p>All related data (like attachments and configuration) may be removed.</p>
        <p>To confirm, type the following text:</p>
      </div>

      <div class="rounded-md bg-gray-100 dark:bg-gray-800 px-3 py-2">
        <code class="text-sm font-mono break-all">{expectedText}</code>
      </div>

      <div>
        <Label for="confirm-delete" class="mb-2 block text-sm font-medium">
          Confirmation
        </Label>
        <Input
          id="confirm-delete"
          type="text"
          bind:value={confirmationText}
          placeholder="Type here to confirm..."
          disabled={deleting}
          class="w-full"
        />
      </div>
    </div>

    <div class="mt-6 flex justify-end gap-2">
      <Dialog.Close>
        <Button variant="outline" class="cursor-pointer" type="button" disabled={deleting}>
          Cancel
        </Button>
      </Dialog.Close>
      <Button variant="destructive" class="cursor-pointer" disabled={deleting || !canDelete} onclick={() => dispatch('confirm')}>
        {#if deleting}
          <Spinner class="mr-2" /> Deleting...
        {:else}
          Delete
        {/if}
      </Button>
    </div>
  </Dialog.Content>
</Dialog.Root>
