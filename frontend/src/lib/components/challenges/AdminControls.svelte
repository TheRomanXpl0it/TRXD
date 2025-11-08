<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button } from '@/components/ui/button';
  import * as Popover from '$lib/components/ui/popover/index.js';
  import * as Command from '$lib/components/ui/command/index.js';
  import Label from '@/components/ui/label/label.svelte';
  import { Input } from '$lib/components/ui/input/index.js';
  import { Spinner } from '$lib/components/ui/spinner/index.js';
  import { Shapes, NotebookPenIcon } from '@lucide/svelte';
  import { toast } from 'svelte-sonner';
  import { createCategory as createCategoryApi } from '$lib/categories';

  const dispatch = createEventDispatcher();

  // local state
  let catPopoverOpen = $state(false);
  let creatingCat = $state(false);
  let newCategoryName = $state('');
  let newCategoryIcon = $state('');

  // categories are managed by parent; we just reload after creation
  // parent listens to `category-created` and refreshes
  async function submitCreateCategory(ev?: SubmitEvent) {
    ev?.preventDefault();
    const name = newCategoryName.trim();
    const icon = newCategoryIcon.trim();
    if (!name || !icon) {
      toast.error('Please enter a category name and an icon.');
      return;
    }
    creatingCat = true;
    try {
      await createCategoryApi(name, icon);
      toast.success('Category created!');
      newCategoryName = '';
      newCategoryIcon = '';
      catPopoverOpen = false;
      dispatch('category-created');
    } catch (e: any) {
      toast.error(e?.message ?? 'Failed to create category.');
    } finally {
      creatingCat = false;
    }
  }
</script>

<div class="mb-6 flex flex-wrap items-center gap-2">
    <Button variant="outline" onclick={() => dispatch('open-create')} class="cursor-pointer">
        <NotebookPenIcon class="mr-2 h-4 w-4" />
        Create Challenge
    </Button>

    <Popover.Root bind:open={catPopoverOpen}>
        <Popover.Trigger>
        {#snippet child({ props })}
            <Button {...props} variant="outline" class="flex cursor-pointer items-center gap-2">
            <Shapes class="h-4 w-4" />
            New Category
            </Button>
        {/snippet}
        </Popover.Trigger>
        <Popover.Content class="w-[320px] p-3">
        <form class="space-y-3" onsubmit={submitCreateCategory}>
            <div>
            <Label for="cat-name" class="mb-1 block text-sm">Category name</Label>
            <Input id="cat-name" placeholder="e.g. Forensics" bind:value={newCategoryName} />
            </div>
            <div>
            <Label for="cat-icon" class="mb-1 block text-sm">Icon (lucide component name)</Label>
            <Input id="cat-icon" placeholder="e.g. Bug, Shield, Lock" bind:value={newCategoryIcon} />
            <p class="mt-1 text-xs text-gray-500">Use any lucide-svelte icon component name.</p>
            </div>
            <div class="flex justify-end gap-2">
            <Button type="button" variant="outline" class="cursor-pointer" onclick={() => (catPopoverOpen = false)}>
                Cancel
            </Button>
            <Button type="submit" class="cursor-pointer" disabled={creatingCat || !newCategoryName.trim() || !newCategoryIcon.trim()}>
                {#if creatingCat}<Spinner class="mr-2" />{/if}
                Create
            </Button>
            </div>
        </form>
        </Popover.Content>
    </Popover.Root>
</div>
