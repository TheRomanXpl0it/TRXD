<script lang="ts">
	import Label from '$lib/components/ui/label/label.svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import StatusBadge from '$lib/components/ui/status-badge.svelte';

	let { config, value = $bindable() } = $props<{
		config: {
			key: string;
			type: string;
			description?: string;
		};
		value: any;
	}>();
</script>

<div
	class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800"
>
	<div class="mb-3 flex items-start justify-between gap-3">
		<div class="min-w-0 flex-1">
			<Label class="mb-1 block font-semibold text-gray-900 dark:text-white">
				{config.key}
			</Label>
			{#if config.description}
				<p class="line-clamp-2 text-xs text-gray-600 dark:text-gray-400">{config.description}</p>
			{/if}
		</div>
		<StatusBadge variant="type">{config.type}</StatusBadge>
	</div>

	<div class="mt-3">
		{#if config.type === 'bool'}
			<div
				class="flex items-center gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 dark:border-gray-700 dark:bg-gray-900"
			>
				<Checkbox bind:checked={value} id={config.key} />
				<Label for={config.key} class="cursor-pointer text-sm font-medium">
					{value ? 'Enabled' : 'Disabled'}
				</Label>
			</div>
		{:else if config.type === 'int'}
			<Input type="number" class="w-full" bind:value placeholder="Enter number" />
		{:else}
			<Input type="text" class="w-full" bind:value placeholder="Enter value" />
		{/if}
	</div>
</div>
