<script lang="ts">
	import { renderMarkdown } from '$lib/utils/markdown';

	interface Props {
		content: string;
		class?: string;
	}

	let { content, class: className, ...rest }: Props = $props();

    // Memoize the HTML rendering to avoid re-rendering on every reactivity cycle
    let html = $derived.by(() => renderMarkdown(content));
</script>

<div class="markdown-content {className || ''}" {...rest}>
    {@html html}
</div>

<style>
	.markdown-content :global(h1) {
		font-size: 1.5rem;
		font-weight: bold;
		margin-top: 1.5rem;
		margin-bottom: 1rem;
	}

	.markdown-content :global(h2) {
		font-size: 1.25rem;
		font-weight: bold;
		margin-top: 1.25rem;
		margin-bottom: 0.75rem;
	}

	.markdown-content :global(h3) {
		font-size: 1.125rem;
		font-weight: 600;
		margin-top: 1rem;
		margin-bottom: 0.5rem;
	}

	.markdown-content :global(h4) {
		font-size: 1rem;
		font-weight: 600;
		margin-top: 0.75rem;
		margin-bottom: 0.5rem;
	}

	.markdown-content :global(p) {
		margin-top: 0.75rem;
		margin-bottom: 0.75rem;
	}

	.markdown-content :global(ul) {
		list-style-type: disc;
		list-style-position: inside;
		margin: 0.75rem 0;
		margin-left: 1rem;
	}

	.markdown-content :global(ol) {
		list-style-type: decimal;
		list-style-position: inside;
		margin: 0.75rem 0;
		margin-left: 1rem;
	}

	.markdown-content :global(li) {
		margin: 0.25rem 0;
	}

	.markdown-content :global(blockquote) {
		border-left: 4px solid;
		border-color: rgba(156, 163, 175, 1);
		padding-left: 1rem;
		padding-top: 0.5rem;
		padding-bottom: 0.5rem;
		margin: 0.75rem 0;
		font-style: italic;
	}

	:global(.dark) .markdown-content :global(blockquote) {
		border-color: rgba(75, 85, 99, 1);
	}

	.markdown-content :global(pre) {
		background-color: rgba(243, 244, 246, 1);
		border-radius: 0.375rem;
		padding: 1rem;
		margin: 0.75rem 0;
		overflow-x: auto;
	}

	:global(.dark) .markdown-content :global(pre) {
		background-color: rgba(31, 41, 55, 1);
	}

	.markdown-content :global(code) {
		background-color: rgba(243, 244, 246, 1);
		border-radius: 0.25rem;
		padding: 0.125rem 0.25rem;
		font-size: 0.875rem;
		font-family: ui-monospace, monospace;
	}

	:global(.dark) .markdown-content :global(code) {
		background-color: rgba(31, 41, 55, 1);
	}

	.markdown-content :global(pre code) {
		background-color: transparent;
		padding: 0;
	}

	.markdown-content :global(a) {
		color: rgb(37, 99, 235);
		text-decoration: underline;
	}

	.markdown-content :global(a:hover) {
		text-decoration: underline;
	}

	:global(.dark) .markdown-content :global(a) {
		color: rgb(96, 165, 250);
	}

	.markdown-content :global(img) {
		max-width: 100%;
		height: auto;
		border-radius: 0.375rem;
		margin: 0.75rem 0;
	}



	.markdown-content :global(table) {
		width: 100%;
		border-collapse: collapse;
		margin: 0.75rem 0;
	}

	.markdown-content :global(th) {
		border: 1px solid rgba(209, 213, 219, 1);
		background-color: rgba(243, 244, 246, 1);
		padding: 0.5rem 1rem;
		font-weight: 600;
		text-align: left;
	}

	:global(.dark) .markdown-content :global(th) {
		border-color: rgba(75, 85, 99, 1);
		background-color: rgba(31, 41, 55, 1);
	}

	.markdown-content :global(td) {
		border: 1px solid rgba(209, 213, 219, 1);
		padding: 0.5rem 1rem;
	}

	:global(.dark) .markdown-content :global(td) {
		border-color: rgba(75, 85, 99, 1);
	}

	.markdown-content :global(hr) {
		margin: 1.5rem 0;
		border-top: 1px solid rgba(209, 213, 219, 1);
	}

	:global(.dark) .markdown-content :global(hr) {
		border-color: rgba(75, 85, 99, 1);
	}

	.markdown-content :global(strong) {
		font-weight: bold;
	}

	.markdown-content :global(em) {
		font-style: italic;
	}
</style>
