<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type monaco from 'monaco-editor';
	type Monaco = typeof import('monaco-editor');
	type Editor = monaco.editor.IStandaloneCodeEditor;

	let {
		value = $bindable(''),
		language = 'yaml',
		onChange = undefined as ((value: string) => void) | undefined,
		class: className = ''
	} = $props();

	let editorContainer: HTMLDivElement;
	let editor: Editor | null = null;
	let isUpdatingFromProp = false;
	let monacoLib: Monaco | null = null;

	// Detect dark mode
	function isDarkMode(): boolean {
		if (typeof document === 'undefined') return false;
		return document.documentElement.classList.contains('dark');
	}

	let observer: MutationObserver | null = null;

	onMount(async () => {
		if (!editorContainer) return;

		monacoLib = await import('monaco-editor');

		editor = monacoLib.editor.create(editorContainer, {
			value: value || '',
			language,
			theme: isDarkMode() ? 'vs-dark' : 'vs',
			automaticLayout: true,
			minimap: { enabled: false },
			scrollBeyondLastLine: false,
			lineNumbers: 'on',
			renderWhitespace: 'selection',
			tabSize: 2,
			insertSpaces: true,
			wordWrap: 'on'
		});

		editor.onDidChangeModelContent(() => {
			if (!editor || isUpdatingFromProp) return;

			const newValue = editor.getValue();
			value = newValue;

			if (onChange) {
				onChange(newValue);
			}
		});

		// Watch for theme changes
		observer = new MutationObserver(() => {
			if (editor) {
				monacoLib?.editor.setTheme(isDarkMode() ? 'vs-dark' : 'vs');
			}
		});

		observer.observe(document.documentElement, {
			attributes: true,
			attributeFilter: ['class']
		});
	});

	$effect(() => {
		if (!editor || !monacoLib) return;
		if (!value) return;

		const currentValue = editor.getValue();
		if (currentValue !== value) {
			isUpdatingFromProp = true;
			editor.setValue(value);
			isUpdatingFromProp = false;
		}
	});

	onDestroy(() => {
		observer?.disconnect();
		if (editor) {
			editor.dispose();
		}
	});
</script>

<div bind:this={editorContainer} class="{className} min-h-45 rounded border"></div>
