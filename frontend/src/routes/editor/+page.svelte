<!-- frontend/src/routes/editor/+page.svelte -->
<script>
	import TextEditor from '$lib/components/TextEditor.svelte';

	// Accept framework-provided props to avoid warnings
	export const data = null;
	export const form = null;

	// Editor mode is owned here so the mode switcher controls the child via bind:
	let editorMode = 'blocks'; // 'blocks', 'markdown', 'richtext'
	let documentTitle = 'Untitled Document';
</script>

<svelte:head>
	<title>{documentTitle} - TPT Text Editor</title>
</svelte:head>

<div class="h-screen flex flex-col bg-white">
	<!-- Thin header: mode switcher only.
	     All document management (Save, Export PDF, New, Open, History, AI, Read Aloud)
	     is handled by TextEditor's own built-in toolbar. -->
	<header class="flex items-center px-6 py-2 border-b border-gray-200 bg-gray-50 shrink-0">
		<span class="text-sm font-medium text-gray-500 mr-3">Editor Mode:</span>
		<div class="flex items-center bg-white rounded-lg border border-gray-200 p-1">
			<button
				class="px-3 py-1 text-sm rounded-md transition-colors {editorMode === 'blocks' ? 'bg-blue-600 text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'}"
				on:click={() => editorMode = 'blocks'}
			>
				Blocks
			</button>
			<button
				class="px-3 py-1 text-sm rounded-md transition-colors {editorMode === 'markdown' ? 'bg-blue-600 text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'}"
				on:click={() => editorMode = 'markdown'}
			>
				Markdown
			</button>
			<button
				class="px-3 py-1 text-sm rounded-md transition-colors {editorMode === 'richtext' ? 'bg-blue-600 text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'}"
				on:click={() => editorMode = 'richtext'}
			>
				Rich Text
			</button>
		</div>
	</header>

	<!-- TextEditor fills the remaining height and provides its own toolbar -->
	<div class="flex-1 overflow-hidden flex flex-col">
		<TextEditor bind:editorMode />
	</div>
</div>
