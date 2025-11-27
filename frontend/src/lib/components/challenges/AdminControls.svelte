<script lang="ts">
	import { Button } from '@/components/ui/button';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Label from '@/components/ui/label/label.svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { Shapes, Plus } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import { createCategory } from '$lib/categories';

	let {
		'onopen-create': onOpenCreate,
		'oncategory-created': onCategoryCreated
	} = $props<{
		'onopen-create'?: () => void;
		'oncategory-created'?: () => void;
	}>();

	let categoryPopoverOpen = $state(false);
	let creating = $state(false);
	let categoryName = $state('');

	const isFormValid = $derived(categoryName.trim().length > 0);

	async function handleCreateCategory(e?: SubmitEvent) {
		e?.preventDefault();

		const name = categoryName.trim();

		if (!name) {
			toast.error('Category name is required.');
			return;
		}

		creating = true;

		try {
			await createCategory(name);
			toast.success(`Category "${name}" created successfully.`);
			categoryName = '';
			categoryPopoverOpen = false;
			onCategoryCreated?.();
		} catch (err: any) {
			toast.error(err?.message ?? 'Failed to create category.');
		} finally {
			creating = false;
		}
	}
</script>

<div class="mb-6 flex flex-wrap items-center gap-3">
	<Button variant="default" onclick={() => onOpenCreate?.()} class="cursor-pointer">
		<Plus class="mr-2 h-4 w-4" />
		Create Challenge
	</Button>

	<Popover.Root bind:open={categoryPopoverOpen}>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button {...props} variant="outline" class="cursor-pointer">
					<Shapes class="mr-2 h-4 w-4" />
					New Category
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-[320px] p-4">
			<div class="mb-4">
				<h4 class="font-semibold text-base">Create Category</h4>
				<p class="text-sm text-muted-foreground mt-1">Add a new challenge category</p>
			</div>

			<form class="space-y-4" onsubmit={handleCreateCategory}>
				<div class="space-y-2">
					<Label for="cat-name" class="text-sm font-medium">Category Name</Label>
					<Input
						id="cat-name"
						placeholder="e.g., Forensics, Web, Crypto"
						bind:value={categoryName}
						disabled={creating}
						autofocus
					/>
				</div>

				<div class="flex justify-end gap-2 pt-2">
					<Button
						type="button"
						variant="outline"
						class="cursor-pointer"
						onclick={() => (categoryPopoverOpen = false)}
						disabled={creating}
					>
						Cancel
					</Button>
					<Button type="submit" class="cursor-pointer" disabled={creating || !isFormValid}>
						{#if creating}
							<Spinner class="mr-2 h-4 w-4" />
							Creating...
						{:else}
							Create
						{/if}
					</Button>
				</div>
			</form>
		</Popover.Content>
	</Popover.Root>
</div>
