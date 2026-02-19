<!-- frontend/src/lib/components/TextEditorMarkdownEditor.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';

	export let markdownContent = '';

	const dispatch = createEventDispatcher();

	const defaultPreview = `# Start writing in Markdown...

## Features
- **Bold text**
- *Italic text*
- \`Code snippets\`
- [Links](https://example.com)

## Math Expressions
You can write natural math expressions:
- integral of x squared dx
- fraction a over b
- square root of 25`;
</script>

<div class="grid grid-cols-2 gap-8 h-full">
	<!-- Editor Panel -->
	<div class="border-r border-gray-200 pr-4">
		<h3 class="text-sm font-medium text-gray-700 mb-2">Markdown Source</h3>
		<textarea
			bind:value={markdownContent}
			placeholder={defaultPreview}
			class="w-full h-96 p-3 border border-gray-300 rounded-md font-mono text-sm resize-none focus:ring-blue-500 focus:border-blue-500"
			on:input={() => dispatch('change', markdownContent)}
		></textarea>
	</div>

	<!-- Preview Panel -->
	<div class="pl-4">
		<h3 class="text-sm font-medium text-gray-700 mb-2">Live Preview</h3>
		<div class="w-full h-96 p-3 border border-gray-300 rounded-md bg-gray-50 overflow-y-auto prose prose-sm max-w-none">
			<div class="whitespace-pre-wrap font-mono text-sm">
				{markdownContent || defaultPreview}
			</div>
		</div>
	</div>
</div>

<style>
	.prose {
		color: #374151;
	}
</style>
