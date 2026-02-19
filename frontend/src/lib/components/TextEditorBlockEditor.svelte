<!-- frontend/src/lib/components/TextEditorBlockEditor.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';

	export let blocks = [];
	export let selectedBlockIndex = 0;

	const dispatch = createEventDispatcher();

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

	function parseMathExpression(text) {
		const expressions = {
			'integral of ([^ ]+) dx': '\\int $1 \\, dx',
			'integral from ([^ ]+) to ([^ ]+) of ([^ ]+) d([^ ]+)': '\\int_{$1}^{$2} $3 \\, d$4',
			'fraction ([^ ]+) over ([^ ]+)': '\\frac{$1}{$2}',
			'square root of ([^ ]+)': '\\sqrt{$1}',
			'sum from ([^=]+)=([^ ]+) to ([^ ]+)': '\\sum_{$1=$2}^{$3}',
			'sum from ([^ ]+) to ([^ ]+) of ([^ ]+)': '\\sum_{$1}^{$2} $3',
			'alpha': '\\alpha',
			'beta': '\\beta',
			'gamma': '\\gamma',
			'delta': '\\delta',
			'pi': '\\pi',
			'sigma': '\\sigma',
			'omega': '\\omega',
			'plus or minus': '\\pm',
			'times': '\\times',
			'divided by': '\\div',
			'therefore': '\\therefore',
			'because': '\\because'
		};

		let result = text;
		for (const [pattern, replacement] of Object.entries(expressions)) {
			const regex = new RegExp(pattern, 'gi');
			result = result.replace(regex, replacement);
		}
		return result;
	}

	function renderBlockContent(block) {
		if (block.type === 'math') {
			return parseMathExpression(block.content);
		}
		return block.content;
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
							class="w-full bg-transparent border-none outline-none resize-none {block.type === 'code' ? 'font-mono' : ''}"
							rows="1"
							on:input={(e) => {
								e.target.style.height = 'auto';
								e.target.style.height = e.target.scrollHeight + 'px';
								dispatch('contentChange', { index, content: e.target.value });
							}}
							on:keydown={(e) => handleKeyDown(e, index)}
							on:focus={() => dispatch('selectBlock', index)}
							use:autoResize
						/>
					{:else}
						<div
							class="cursor-text"
							on:click={() => dispatch('selectBlock', index)}
						>
							{#if block.type === 'math'}
								<span class="font-serif">
									{renderBlockContent(block) || blockTypes[block.type]?.placeholder}
								</span>
								{#if block.content}
									<span class="ml-2 text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded">
										Math
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
</style>
