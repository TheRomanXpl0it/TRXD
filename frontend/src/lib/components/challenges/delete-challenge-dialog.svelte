<script lang="ts">
  import * as Dialog from '$lib/components/ui/dialog/index.js';
  import { Button } from '@/components/ui/button';
  import { Spinner } from '$lib/components/ui/spinner/index.js';
  import { createEventDispatcher } from 'svelte';


  
  let { open = $bindable(false), toDelete = $bindable(null), deleting = $bindable(false) } = $props<{ open?: boolean; toDelete: any, deleting: boolean }>();

  const dispatch = createEventDispatcher();
</script>

<Dialog.Root bind:open>
  <Dialog.Overlay />
  <Dialog.Content class="sm:max-w-[520px]">
    <Dialog.Header class="pb-2">
      <Dialog.Title>Delete challenge?</Dialog.Title>
      <Dialog.Description>
        You’re about to permanently delete <b>{toDelete?.name ?? 'this challenge'}</b>. This action cannot be undone.
      </Dialog.Description>
    </Dialog.Header>

    <div class="mt-4 space-y-2 text-sm text-gray-600 dark:text-gray-300">
      <p>All related data (like attachments and configuration) may be removed.</p>
      <p>Please confirm to proceed.</p>
    </div>

    <div class="mt-6 flex justify-end gap-2">
      <Dialog.Close>
        <Button variant="outline" class="cursor-pointer" type="button" disabled={deleting}>
          Cancel
        </Button>
      </Dialog.Close>
      <Button variant="destructive" class="cursor-pointer" disabled={deleting} onclick={() => dispatch('confirm')}>
        {#if deleting}
          <Spinner class="mr-2" /> Deleting…
        {:else}
          Delete
        {/if}
      </Button>
    </div>
  </Dialog.Content>
</Dialog.Root>
