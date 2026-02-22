<script lang="ts">
	import Label from '$lib/components/ui/label/label.svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';

	let {
		id,
		label,
		value = $bindable(''),
		type = 'text',
		placeholder = '',
		disabled = false,
		hint,
		multiline = false,
		rows = 3
	} = $props<{
		id: string;
		label: string;
		value?: string;
		type?: string;
		placeholder?: string;
		disabled?: boolean;
		hint?: string;
		multiline?: boolean;
		rows?: number;
	}>();
</script>

<div class="space-y-2">
	<Label for={id} class="text-sm font-medium">{label}</Label>
	{#if multiline}
		<Textarea {id} bind:value {placeholder} {disabled} {rows} />
	{:else}
		<Input {id} {type} bind:value {placeholder} {disabled} />
	{/if}
	{#if hint}
		<p class="text-muted-foreground text-xs">{hint}</p>
	{/if}
</div>
