<script lang="ts">
  import * as Dialog from '$lib/components/ui/dialog/index.js';
  import Label from '@/components/ui/label/label.svelte';
  import { Input } from '$lib/components/ui/input/index.js';
  import { Textarea } from '$lib/components/ui/textarea/index.js';
  import { Button } from '@/components/ui/button';
  import { Spinner } from '$lib/components/ui/spinner/index.js';
  import MultiSelect from '$lib/components/challenges/category-select.svelte';
  import * as Tooltip from '$lib/components/ui/tooltip/index.js';
  import { Checkbox } from '$lib/components/ui/checkbox/index.js';
  import { toast } from 'svelte-sonner';
  import { createEventDispatcher } from 'svelte';
  import { createChallenge } from '$lib/challenges';

  const defaultChallengeTypes = [
    { value: 'Container', label: 'Container' },
    { value: 'Compose', label: 'Compose' },
    { value: 'Normal', label: 'Normal' }
  ];

  
  let {
    open = $bindable(false),
    challengeName = $bindable(""),
    challengeDescription = $bindable(""),
    category = $bindable(""),
    challengeType = $bindable(""),
    points = $bindable(500),
    dynamicScore = $bindable(true),
    categories = $bindable([]),
    challengeTypes = []
  } = $props<{
    open: boolean,
    challengeName: string
    challengeDescription: string;
    category: string;
    challengeType:string;
    points: number;
    dynamicScore: boolean;
    categories: Array<string>;
    challengeTypes: Array<{value:string,label:string}>
    
    
  }>();

  let createLoading = $state(false);
  const dispatch = createEventDispatcher();

  async function submitCreateChallenge(ev: SubmitEvent) {
    ev.preventDefault();
    if (createLoading) return;

    const trimmedName = challengeName.trim();
    if (!trimmedName) return toast.error('Please enter a challenge name.');
    if (!category) return toast.error('Please select a category.');
    if (!challengeType) return toast.error('Please select a challenge type.');
    if (typeof points !== 'number' || Number.isNaN(points) || points < 0) {
      return toast.error('Please choose a valid points value.');
    }

    createLoading = true;
    const scoretype = dynamicScore ? 'Dynamic' : 'Static';

    try {
      await createChallenge(
        trimmedName,
        category,
        challengeDescription.trim(),
        challengeType,
        points,
        scoretype
      );
      toast.success('Challenge created!');

      // reset local fields
      challengeName = '';
      challengeDescription = '';
      category = null;
      challengeType = 'Container';
      dynamicScore = false;
      points = 500;

      open = false;
      dispatch('created');
    } catch (err: any) {
      toast.error(err?.message ?? 'Failed to create challenge.');
    } finally {
      createLoading = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Overlay />
  <Dialog.Content class="sm:max-w-[720px]">
    <Dialog.Header class="pb-2">
      <Dialog.Title>Create Challenge</Dialog.Title>
      <Dialog.Description>
        Create the barebones skeleton of the challenge; edit later to upload files and set advanced options.
      </Dialog.Description>
    </Dialog.Header>

    <div class="mt-2 space-y-4">
      <form onsubmit={submitCreateChallenge}>
        <Label for="name" class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">
          Challenge Name*
        </Label>
        <Input id="name" type="text" bind:value={challengeName} required class="mb-4 w-full" />

        <Label for="description" class="mb-2 mt-4 block text-sm font-medium text-gray-900 dark:text-white">
          Description
        </Label>
        <Textarea id="description" bind:value={challengeDescription} class="mb-4 w-full" />

        <div class="flex flex-row flex-wrap items-start gap-6">
          <div class="flex min-w-56 flex-col">
            <Label for="category" class="mb-2 mt-4 block text-sm font-medium text-gray-900 dark:text-white">
              Category*
            </Label>
            <MultiSelect id="category" items={categories} bind:value={category} placeholder="Select a category…" />
          </div>

          <div class="flex min-w-56 flex-col">
            <Label for="type" class="mb-2 mt-4 block text-sm font-medium text-gray-900 dark:text-white">
              Type*
            </Label>
            <MultiSelect id="type" items={challengeTypes} bind:value={challengeType} placeholder="Select type..." />
          </div>

          <div class="flex min-w-40 flex-col">
            <Tooltip.Provider>
              <Tooltip.Root>
                <Tooltip.Trigger>
                  <Label for="scoretype" class="mb-2 mt-4 text-sm font-medium text-gray-900 dark:text-white">
                    Dynamic score*
                  </Label>
                  <div class="flex flex-row">
                    <Checkbox id="scoretype" class="mb-4 mt-2" bind:checked={dynamicScore} />
                  </div>
                </Tooltip.Trigger>
                <Tooltip.Content>
                  <p>Dynamic scoring decays challenge points over number of solves.</p>
                </Tooltip.Content>
              </Tooltip.Root>
            </Tooltip.Provider>
          </div>
        </div>

        <div class="mt-2">
          <Label for="points" class="mt-4 block text-sm font-medium text-gray-900 dark:text-white">
            Points
          </Label>
          <Input
            id="points"
            type="number"
            inputmode="numeric"
            min="0"
            max="1500"
            step="1"
            bind:value={points}
            oninput={(e) => {
              const v = (e.currentTarget as HTMLInputElement).valueAsNumber;
              const n = Number.isFinite(v) ? Math.max(0, Math.floor(v)) : 0;
              points = n;
              (e.currentTarget as HTMLInputElement).value = String(n);
            }}
          />
        </div>

        <div class="mt-6 flex justify-end gap-2">
          <Dialog.Close>
            <Button variant="outline" class="cursor-pointer" type="button">Cancel</Button>
          </Dialog.Close>
          <Button type="submit" class="cursor-pointer" disabled={createLoading}>
            {#if createLoading}<Spinner class="mr-2" />{/if}
            Create
          </Button>
        </div>
      </form>
    </div>
  </Dialog.Content>
</Dialog.Root>
