<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import * as monaco from 'monaco-editor';

	let {
		value = $bindable(''),
		language = 'yaml',
		onChange = undefined as ((value: string) => void) | undefined,
		class: className = ''
	} = $props();

	let editorContainer: HTMLDivElement;
	let editor: monaco.editor.IStandaloneCodeEditor | null = null;
	let isUpdatingFromProp = false;

	// Detect dark mode
	function isDarkMode(): boolean {
		if (typeof document === 'undefined') return false;
		return document.documentElement.classList.contains('dark');
	}

	onMount(() => {
		if (!editorContainer) return;

		editor = monaco.editor.create(editorContainer, {
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
		const observer = new MutationObserver(() => {
			if (editor) {
				monaco.editor.setTheme(isDarkMode() ? 'vs-dark' : 'vs');
			}
		});

		observer.observe(document.documentElement, {
			attributes: true,
			attributeFilter: ['class']
		});

		return () => {
			observer.disconnect();
		};
	});

	$effect(() => {
		if (!editor || !value) return;

		const currentValue = editor.getValue();
		if (currentValue !== value) {
			isUpdatingFromProp = true;
			editor.setValue(value);
			isUpdatingFromProp = false;
		}
	});

	onDestroy(() => {
		if (editor) {
			editor.dispose();
		}
	});
</script>

<div bind:this={editorContainer} class="{className} min-h-45 rounded border"></div>
