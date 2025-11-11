<script lang="ts">
	let { title = "Error", message }: { title?: string; message: string } = $props();
</script>

<div class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-600 dark:border-red-800 dark:bg-red-950/20">
	<p class="font-semibold">{title}</p>
	<p class="text-sm">{message}</p>
</div>
