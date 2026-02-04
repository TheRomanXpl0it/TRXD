<script lang="ts">
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';

	let { message = 'Loading...', class: className = '' } = $props<{
		message?: string;
		class?: string;
	}>();
</script>

<div class="flex flex-col items-center justify-center py-12 {className}">
	<Spinner class="mb-4 h-8 w-8" />
	<p class="text-gray-600 dark:text-gray-400">{message}</p>
</div>
