<script lang="ts">
	import { Package } from '@lucide/svelte';

	let { icon: Icon = Package, title = 'No data', description = '', action = null } = $props();
</script>

<div class="flex flex-col items-center justify-center py-12 text-center">
	<div class="mb-4 rounded-full bg-gray-100 p-3 dark:bg-gray-800">
		<Icon class="h-6 w-6 text-gray-400" />
	</div>
	<h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">{title}</h3>
	{#if description}
		<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{description}</p>
	{/if}
	{#if action}
		<div class="mt-4">
			{@render action()}
		</div>
	{/if}
</div>
