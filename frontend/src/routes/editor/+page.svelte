<!-- frontend/src/routes/editor/+page.svelte -->
<script>
	import TextEditor from '$lib/components/TextEditor.svelte';

	// Accept framework-provided props to avoid warnings
	export const data = null;
	export const form = null;
	export const params = null;

	// Editor mode is owned here so the mode switcher controls the child via bind:
	let editorMode = 'blocks'; // 'blocks', 'markdown', 'richtext'
	let documentTitle = 'Untitled Document';
</script>

<svelte:head>
	<title>{documentTitle} - TPT Text Editor</title>
</svelte:head>

<div class="h-screen flex flex-col bg-white">
	<!-- TextEditor fills the remaining height and provides its own toolbar with integrated mode switcher -->
	<div class="flex-1 overflow-hidden flex flex-col">
		<TextEditor bind:editorMode />
	</div>
</div>
