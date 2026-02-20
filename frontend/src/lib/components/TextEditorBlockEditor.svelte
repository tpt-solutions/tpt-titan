<!-- frontend/src/lib/components/TextEditorBlockEditor.svelte -->
<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import mathService from '$lib/services/math.js';

	export let blocks = [];
	export let selectedBlockIndex = 0;

	const dispatch = createEventDispatcher();

	// Track math conversion state
	let mathConversions = {};
	let convertingMath = {};


	const blockTypes = {
		text: { icon: '📝', placeholder: 'Type something...' },
		heading1: { icon: 'H1', placeholder: 'Heading 1' },
		heading2: { icon: 'H2', placeholder: 'Heading 2' },
		heading3: { icon: 'H3', placeholder: 'Heading 3' },
		list: { icon: '•', placeholder: 'List item' },
		quote: { icon: '"', placeholder: 'Quote' },
		code: { icon: '</>', placeholder: 'Code block' },
		math: { icon: '∫', placeholder: 'Math expression (e.g., integral of x squared dx)' },
		table: { icon: '📊', placeholder: 'Create a table' },
		image: { icon: '🖼️', placeholder: 'Add an image' }
	};

	function getBlockStyles(blockType) {
		const styles = {
			heading1: 'text-3xl font-bold mb-4',
			heading2: 'text-2xl font-semibold mb-3',
			heading3: 'text-xl font-medium mb-2',
			text: 'text-base mb-2',
			list: 'text-base mb-1 ml-4',
			quote: 'text-base mb-2 pl-4 border-l-4 border-gray-300 italic',
			code: 'text-sm font-mono bg-gray-100 p-3 rounded mb-2',
			math: 'text-base mb-2 font-serif bg-purple-50 p-2 rounded border border-purple-200',
			table: 'text-base mb-2',
			image: 'text-base mb-2'
		};
		return styles[blockType] || styles.text;
	}

	async function parseMathExpression(text) {
		if (!text || text.trim() === '') return '';
		
		// Use the math service for superior natural language processing
		try {
			const result = await mathService.parseNaturalMath(text);
			return result.latex || text;
		} catch (error) {
			console.error('Math parsing error:', error);
			return text;
		}
	}

	async function convertMathBlock(blockIndex) {
		const block = blocks[blockIndex];
		if (!block || block.type !== 'math' || !block.content) return;
		
		const content = block.content.trim();
		if (!content) return;
		
		// Check cache first
		if (mathConversions[content]) {
			blocks[blockIndex].mathData = mathConversions[content];
			blocks = [...blocks];
			return;
		}
		
		convertingMath[blockIndex] = true;
		
		try {
			const result = await mathService.parseNaturalMath(content);
			mathConversions[content] = result;
			blocks[blockIndex].mathData = result;
			blocks = [...blocks];
		} catch (error) {
			console.error('Math conversion error:', error);
		} finally {
			convertingMath[blockIndex] = false;
		}
	}

	function renderBlockContent(block) {
		if (block.type === 'math' && block.mathData) {
			return block.mathData.latex || block.content;
		}
		return block.content;
	}

	function renderMathHTML(block) {
		if (block.type === 'math' && block.mathData) {
			return mathService.renderToHTML(block.mathData.latex);
		}
		return block.content;
	}

	// Convert math when block loses focus
	function handleMathBlur(blockIndex) {
		convertMathBlock(blockIndex);
	}


	function handleKeyDown(event, blockIndex) {
		const { key } = event;

		if (key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			dispatch('addBlock', blockIndex);
		} else if (key === 'Backspace' && blocks[blockIndex].content === '') {
			event.preventDefault();
			dispatch('deleteBlock', blockIndex);
		}
	}

	function autoResize(node) {
		const resize = () => {
			node.style.height = 'auto';
			node.style.height = node.scrollHeight + 'px';
		};

		node.addEventListener('input', resize);
		node.addEventListener('focus', resize);
		setTimeout(resize, 0);

		return {
			destroy() {
				node.removeEventListener('input', resize);
				node.removeEventListener('focus', resize);
			}
		};
	}
</script>

<div class="space-y-2">
	{#each blocks as block, index}
		<div
			class="group relative {selectedBlockIndex === index ? 'ring-2 ring-blue-500 rounded' : ''}"
			on:click={() => dispatch('selectBlock', index)}
		>
			<!-- Block Type Indicator -->
			<div class="absolute -left-8 top-1 opacity-0 group-hover:opacity-100 transition-opacity">
				<span class="text-xs text-gray-400">{blockTypes[block.type]?.icon || '📝'}</span>
			</div>

			<!-- Block Content -->
			{#if block.type === 'image'}
				<div class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center text-gray-500 hover:border-blue-400 transition-colors cursor-pointer">
					<svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
					</svg>
					<p>Click to add an image</p>
				</div>
			{:else if block.type === 'table'}
				<div class="border border-gray-300 rounded overflow-hidden">
					<table class="w-full">
						<thead class="bg-gray-50">
							<tr>
								<th class="border border-gray-300 p-2 text-left">Column 1</th>
								<th class="border border-gray-300 p-2 text-left">Column 2</th>
								<th class="border border-gray-300 p-2 text-left">Column 3</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td class="border border-gray-300 p-2">Data 1</td>
								<td class="border border-gray-300 p-2">Data 2</td>
								<td class="border border-gray-300 p-2">Data 3</td>
							</tr>
						</tbody>
					</table>
				</div>
			{:else}
				<div class={getBlockStyles(block.type)}>
					{#if block.type === 'list'}
						<span class="mr-2">•</span>
					{:else if block.type === 'quote'}
						<span class="mr-2 text-gray-400">"</span>
					{/if}

					{#if selectedBlockIndex === index}
						<textarea
							bind:value={block.content}
							placeholder={blockTypes[block.type]?.placeholder || 'Start typing...'}
							class="w-full bg-transparent border-none outline-none resize-none {block.type === 'code' ? 'font-mono' : ''} {block.type === 'math' ? 'font-serif' : ''}"
							rows="1"
							on:input={(e) => {
								e.target.style.height = 'auto';
								e.target.style.height = e.target.scrollHeight + 'px';
								dispatch('contentChange', { index, content: e.target.value });
								// Clear math data when editing
								if (block.type === 'math') {
									block.mathData = null;
								}
							}}
							on:keydown={(e) => handleKeyDown(e, index)}
							on:focus={() => dispatch('selectBlock', index)}
							on:blur={() => {
								if (block.type === 'math') {
									handleMathBlur(index);
								}
							}}
							use:autoResize
						/>

					{:else}
						<div
							class="cursor-text"
							on:click={() => dispatch('selectBlock', index)}
						>
					{#if block.type === 'math'}
						<div class="math-display font-serif text-lg">
							{#if convertingMath[index]}
								<span class="text-gray-400 italic">Converting...</span>
							{:else if block.mathData}
								<span class="math-rendered" title="LaTeX: {block.mathData.latex}">
									{@html renderMathHTML(block)}
								</span>
							{:else}
								<span class="text-gray-600">{block.content || blockTypes[block.type]?.placeholder}</span>
							{/if}
						</div>
						{#if block.mathData}
							<span class="ml-2 text-xs text-green-600 bg-green-50 px-2 py-1 rounded border border-green-200">
								✓ Math
							</span>
						{:else if block.content}
							<span class="ml-2 text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded">
								Press Enter to convert
							</span>
						{/if}

							{:else}
								{block.content || blockTypes[block.type]?.placeholder}
							{/if}
						</div>
					{/if}
				</div>
			{/if}

			<!-- Block Actions (visible on hover) -->
			<div class="absolute -right-8 top-1 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col space-y-1">
				<button
					class="w-6 h-6 bg-gray-100 hover:bg-gray-200 rounded text-xs"
					title="Add block below"
					on:click|stopPropagation={() => dispatch('addBlock', index)}
				>
					+
				</button>
				{#if blocks.length > 1}
					<button
						class="w-6 h-6 bg-red-100 hover:bg-red-200 rounded text-xs text-red-600"
						title="Delete block"
						on:click|stopPropagation={() => dispatch('deleteBlock', index)}
					>
						×
					</button>
				{/if}
			</div>
		</div>
	{/each}
</div>

<style>
	textarea {
		min-height: 1.5em;
	}

	.math-display {
		min-height: 1.5em;
		padding: 0.25rem 0;
	}

	.math-rendered :global(.fraction) {
		display: inline-flex;
		flex-direction: column;
		vertical-align: middle;
		text-align: center;
		margin: 0 0.2em;
	}

	.math-rendered :global(.fraction .num) {
		border-bottom: 1px solid currentColor;
		padding: 0 0.2em;
	}

	.math-rendered :global(.fraction .den) {
		padding: 0 0.2em;
	}

	.math-rendered :global(sup) {
		font-size: 0.75em;
		vertical-align: super;
		line-height: 0;
	}

	.math-rendered :global(sub) {
		font-size: 0.75em;
		vertical-align: sub;
		line-height: 0;
	}
</style>
