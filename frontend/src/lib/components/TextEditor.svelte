<script>
	import { onMount, createEventDispatcher } from 'svelte';
	import { jsPDF } from 'jspdf';
	import { createDocument, updateDocument, getDocuments, getDocument, getDocumentVersions, restoreDocumentVersion, deleteDocument } from '../api.js';

	export let editorMode = 'blocks'; // 'blocks', 'markdown', 'richtext'
	export let documentId = null; // For loading existing documents

	const dispatch = createEventDispatcher();

	let editorContent = '';
	let markdownContent = '';
	let blocks = [];
	let selectedBlockIndex = 0;
	let isComposing = false;

	// Document management
	let currentDocument = null;
	let documentTitle = 'Untitled Document';
	let isSaving = false;
	let isLoading = false;
	let saveStatus = '';
	let documentList = [];
	let showDocumentList = false;
	let showVersionHistory = false;
	let versionHistory = [];
	let hasUnsavedChanges = false;

	// Auto-save timer
	let autoSaveTimer = null;

	// Available block types
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

	// Math expression parser (simplified)
	function parseMathExpression(text) {
		// Convert natural language to LaTeX/math notation
		const expressions = {
			// Integrals
			'integral of ([^ ]+) dx': '\\int $1 \\, dx',
			'integral from ([^ ]+) to ([^ ]+) of ([^ ]+) d([^ ]+)': '\\int_{$1}^{$2} $3 \\, d$4',

			// Fractions
			'fraction ([^ ]+) over ([^ ]+)': '\\frac{$1}{$2}',

			// Square roots
			'square root of ([^ ]+)': '\\sqrt{$1}',

			// Summation
			'sum from ([^=]+)=([^ ]+) to ([^ ]+)': '\\sum_{$1=$2}^{$3}',
			'sum from ([^ ]+) to ([^ ]+) of ([^ ]+)': '\\sum_{$1}^{$2} $3',

			// Greek letters
			'alpha': '\\alpha',
			'beta': '\\beta',
			'gamma': '\\gamma',
			'delta': '\\delta',
			'pi': '\\pi',
			'sigma': '\\sigma',
			'omega': '\\omega',

			// Basic operators
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

	// Initialize with sample content
	onMount(() => {
		if (blocks.length === 0) {
			blocks = [
				{
					id: 1,
					type: 'heading1',
					content: 'Welcome to TPT Text Editor',
					properties: {}
				},
				{
					id: 2,
					type: 'text',
					content: 'This revolutionary text editor combines the best of Notion, Markdown, and Microsoft Word with groundbreaking features.',
					properties: {}
				},
				{
					id: 3,
					type: 'heading2',
					content: 'Natural Math Input (Better Than LaTeX!)',
					properties: {}
				},
				{
					id: 4,
					type: 'math',
					content: 'integral from 0 to infinity of e to the power of negative x squared dx',
					properties: {}
				},
				{
					id: 5,
					type: 'text',
					content: 'The equation above was created by simply typing "integral from 0 to infinity of e to the power of negative x squared dx" - no complex LaTeX syntax required!',
					properties: {}
				}
			];
		}
	});

	function addBlock(afterIndex = blocks.length - 1) {
		const newBlock = {
			id: Date.now() + Math.random(),
			type: 'text',
			content: '',
			properties: {}
		};

		blocks.splice(afterIndex + 1, 0, newBlock);
		blocks = [...blocks]; // Trigger reactivity
		selectedBlockIndex = afterIndex + 1;
	}

	function deleteBlock(blockIndex) {
		if (blocks.length > 1) {
			blocks.splice(blockIndex, 1);
			blocks = [...blocks];
			selectedBlockIndex = Math.max(0, blockIndex - 1);
		}
	}

	function updateBlockContent(blockIndex, content) {
		blocks[blockIndex].content = content;
		blocks = [...blocks];
	}

	function changeBlockType(blockIndex, newType) {
		blocks[blockIndex].type = newType;
		blocks = [...blocks];
	}

	function handleKeyDown(event, blockIndex) {
		const { key } = event;

		if (key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			addBlock(blockIndex);
		} else if (key === 'Backspace' && blocks[blockIndex].content === '') {
			event.preventDefault();
			deleteBlock(blockIndex);
		} else if (key === '/' && blocks[blockIndex].content === '') {
			// Show block type selector
			event.preventDefault();
			// This would trigger a block type selector UI
		}
	}

	function renderBlockContent(block) {
		if (block.type === 'math') {
			// Parse natural language to math notation
			return parseMathExpression(block.content);
		}
		return block.content;
	}

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

	// Document management functions
	async function saveDocument() {
		if (isSaving) return;

		isSaving = true;
		saveStatus = 'Saving...';

		try {
			const documentData = {
				title: documentTitle,
				content_type: 'text',
				content: {
					mode: editorMode,
					blocks: blocks,
					markdown: markdownContent
				}
			};

			let result;
			if (currentDocument) {
				result = await updateDocument(currentDocument.id, documentData);
			} else {
				result = await createDocument(documentData);
				currentDocument = result;
			}

			saveStatus = 'Saved';
			hasUnsavedChanges = false;

			setTimeout(() => {
				saveStatus = '';
			}, 2000);

		} catch (error) {
			saveStatus = 'Save failed';
			console.error('Save error:', error);
		} finally {
			isSaving = false;
		}
	}

	async function loadDocumentList() {
		try {
			const response = await getDocuments();
			documentList = response.documents.filter(doc => doc.content_type === 'text');
			showDocumentList = true;
		} catch (error) {
			console.error('Failed to load documents:', error);
		}
	}

	async function loadDocument(doc) {
		isLoading = true;
		try {
			const response = await getDocument(doc.id);
			currentDocument = doc;
			documentTitle = doc.title;

			// Load content based on editor mode
			if (response.content && typeof response.content === 'object') {
				if (response.content.blocks) {
					blocks = response.content.blocks;
				}
				if (response.content.markdown) {
					markdownContent = response.content.markdown;
				}
				if (response.content.mode) {
					editorMode = response.content.mode;
				}
			}

			showDocumentList = false;
			hasUnsavedChanges = false;
		} catch (error) {
			console.error('Failed to load document:', error);
		} finally {
			isLoading = false;
		}
	}

	async function loadVersionHistory() {
		if (!currentDocument) return;

		try {
			const response = await getDocumentVersions(currentDocument.id);
			versionHistory = response.versions;
			showVersionHistory = true;
		} catch (error) {
			console.error('Failed to load version history:', error);
		}
	}

	async function restoreVersion(version) {
		if (!currentDocument) return;

		try {
			await restoreDocumentVersion(currentDocument.id, version.version);
			await loadDocument(currentDocument); // Reload the document
			showVersionHistory = false;
		} catch (error) {
			console.error('Failed to restore version:', error);
		}
	}

	function createNewDocument() {
		currentDocument = null;
		documentTitle = 'Untitled Document';
		blocks = [{
			id: 1,
			type: 'heading1',
			content: 'New Document',
			properties: {}
		}];
		markdownContent = '';
		hasUnsavedChanges = false;
		showDocumentList = false;
		showVersionHistory = false;
	}

	function markAsChanged() {
		hasUnsavedChanges = true;
		if (autoSaveTimer) {
			clearTimeout(autoSaveTimer);
		}
		autoSaveTimer = setTimeout(() => {
			if (hasUnsavedChanges) {
				saveDocument();
			}
		}, 30000); // Auto-save after 30 seconds of inactivity
	}

	// Override existing functions to mark changes
	const originalUpdateBlockContent = updateBlockContent;
	updateBlockContent = (blockIndex, content) => {
		originalUpdateBlockContent(blockIndex, content);
		markAsChanged();
	};

	const originalAddBlock = addBlock;
	addBlock = (afterIndex) => {
		originalAddBlock(afterIndex);
		markAsChanged();
	};

	const originalDeleteBlock = deleteBlock;
	deleteBlock = (blockIndex) => {
		originalDeleteBlock(blockIndex);
		markAsChanged();
	};

	// Export to PDF
	function exportToPDF() {
		const doc = new jsPDF();

		// Set title
		doc.setFontSize(20);
		doc.text(documentTitle || 'Untitled Document', 20, 30);

		let yPosition = 50;

		// Export blocks
		blocks.forEach(block => {
			doc.setFontSize(getFontSize(block.type));

			// Handle different block types
			let content = block.content;
			if (block.type === 'math') {
				content = renderBlockContent(block);
			}

			// Split long text into lines
			const lines = doc.splitTextToSize(content, 170);

			// Add lines to PDF
			lines.forEach(line => {
				if (yPosition > 270) { // New page if near bottom
					doc.addPage();
					yPosition = 30;
				}
				doc.text(line, 20, yPosition);
				yPosition += getLineHeight(block.type);
			});

			// Add extra space after block
			yPosition += 10;
		});

		// Save the PDF
		doc.save(`${documentTitle || 'document'}.pdf`);
	}

	function getFontSize(blockType) {
		const sizes = {
			heading1: 18,
			heading2: 16,
			heading3: 14,
			text: 12,
			list: 12,
			quote: 12,
			code: 10,
			math: 12,
			table: 12,
			image: 12
		};
		return sizes[blockType] || 12;
	}

	function getLineHeight(blockType) {
		const heights = {
			heading1: 8,
			heading2: 7,
			heading3: 6,
			text: 6,
			list: 6,
			quote: 6,
			code: 5,
			math: 6,
			table: 6,
			image: 6
		};
		return heights[blockType] || 6;
	}

	// Svelte action for auto-resizing textarea
	function autoResize(node) {
		const resize = () => {
			node.style.height = 'auto';
			node.style.height = node.scrollHeight + 'px';
		};

		node.addEventListener('input', resize);
		node.addEventListener('focus', resize);

		// Initial resize
		setTimeout(resize, 0);

		return {
			destroy() {
				node.removeEventListener('input', resize);
				node.removeEventListener('focus', resize);
			}
		};
	}
</script>

<!-- Document Management Toolbar -->
<div class="bg-white border-b border-gray-200 px-8 py-4">
	<div class="max-w-4xl mx-auto flex items-center justify-between">
		<div class="flex items-center space-x-4">
			<input
				bind:value={documentTitle}
				placeholder="Document title..."
				class="text-xl font-semibold bg-transparent border-none outline-none focus:ring-2 focus:ring-blue-500 rounded px-2 py-1"
				on:input={markAsChanged}
			/>
			{#if hasUnsavedChanges}
				<span class="text-sm text-orange-600">•</span>
			{/if}
			<span class="text-sm text-gray-500">{saveStatus}</span>
		</div>

		<div class="flex items-center space-x-2">
			<button
				class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
				on:click={saveDocument}
				disabled={isSaving}
			>
				{isSaving ? 'Saving...' : 'Save'}
			</button>

			<button
				class="px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700"
				on:click={exportToPDF}
			>
				Export PDF
			</button>

			<button
				class="px-3 py-1 text-sm bg-gray-600 text-white rounded hover:bg-gray-700"
				on:click={createNewDocument}
			>
				New
			</button>

			<button
				class="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700"
				on:click={loadDocumentList}
			>
				Open
			</button>

			{#if currentDocument}
				<button
					class="px-3 py-1 text-sm bg-purple-600 text-white rounded hover:bg-purple-700"
					on:click={loadVersionHistory}
				>
					History
				</button>
			{/if}
		</div>
	</div>
</div>

<!-- Document List Modal -->
{#if showDocumentList}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-md max-h-96 overflow-y-auto">
			<h3 class="text-lg font-semibold mb-4">Open Document</h3>
			{#if documentList.length === 0}
				<p class="text-gray-500">No documents found.</p>
			{:else}
				<div class="space-y-2">
					{#each documentList as doc}
						<button
							class="w-full text-left p-3 border border-gray-200 rounded hover:bg-gray-50"
							on:click={() => loadDocument(doc)}
							disabled={isLoading}
						>
							<div class="font-medium">{doc.title}</div>
							<div class="text-sm text-gray-500">
								v{doc.version} • {new Date(doc.updated_at).toLocaleDateString()}
							</div>
						</button>
					{/each}
				</div>
			{/if}
			<div class="mt-4 flex justify-end">
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => showDocumentList = false}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Version History Modal -->
{#if showVersionHistory}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-lg max-h-96 overflow-y-auto">
			<h3 class="text-lg font-semibold mb-4">Version History</h3>
			<div class="space-y-2">
				{#each versionHistory as version}
					<div class="flex items-center justify-between p-3 border border-gray-200 rounded">
						<div>
							<div class="font-medium">Version {version.version}</div>
							<div class="text-sm text-gray-500">
								{new Date(version.created_at).toLocaleString()}
								{#if version.is_active}
									<span class="ml-2 text-green-600">(Current)</span>
								{/if}
							</div>
						</div>
						{#if !version.is_active}
							<button
								class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
								on:click={() => restoreVersion(version)}
							>
								Restore
							</button>
						{/if}
					</div>
				{/each}
			</div>
			<div class="mt-4 flex justify-end">
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => showVersionHistory = false}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<div class="h-full overflow-y-auto p-8 bg-white">
	<div class="max-w-4xl mx-auto">
		{#if editorMode === 'blocks'}
			<!-- Notion-Style Block Editor -->
			<div class="space-y-2">
				{#each blocks as block, index}
					<div
						class="group relative {selectedBlockIndex === index ? 'ring-2 ring-blue-500 rounded' : ''}"
						on:click={() => selectedBlockIndex = index}
					>
						<!-- Block Type Indicator -->
						<div class="absolute -left-8 top-1 opacity-0 group-hover:opacity-100 transition-opacity">
							<span class="text-xs text-gray-400">{blockTypes[block.type]?.icon || '📝'}</span>
						</div>

						<!-- Block Content -->
						{#if block.type === 'image'}
							<!-- Image Block -->
							<div class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center text-gray-500 hover:border-blue-400 transition-colors cursor-pointer">
								<svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
								</svg>
								<p>Click to add an image</p>
							</div>
						{:else if block.type === 'table'}
							<!-- Table Block -->
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
							<!-- Text/Math Blocks -->
							<div class={getBlockStyles(block.type)}>
								{#if block.type === 'list'}
									<span class="mr-2">•</span>
								{:else if block.type === 'quote'}
									<span class="mr-2 text-gray-400">"</span>
								{/if}

								{#if selectedBlockIndex === index}
									<!-- Editable Content -->
									<textarea
										bind:value={block.content}
										placeholder={blockTypes[block.type]?.placeholder || 'Start typing...'}
										class="w-full bg-transparent border-none outline-none resize-none {block.type === 'code' ? 'font-mono' : ''}"
										rows="1"
										on:input={(e) => {
											e.target.style.height = 'auto';
											e.target.style.height = e.target.scrollHeight + 'px';
											updateBlockContent(index, e.target.value);
										}}
										on:keydown={(e) => handleKeyDown(e, index)}
										on:focus={() => selectedBlockIndex = index}
										use:autoResize
									/>
								{:else}
									<!-- Display Content -->
									<div
										class="cursor-text"
										on:click={() => selectedBlockIndex = index}
									>
										{#if block.type === 'math'}
											<!-- Render Math Expression -->
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
								on:click={() => addBlock(index)}
							>
								+
							</button>
							{#if blocks.length > 1}
								<button
									class="w-6 h-6 bg-red-100 hover:bg-red-200 rounded text-xs text-red-600"
									title="Delete block"
									on:click={() => deleteBlock(index)}
								>
									×
								</button>
							{/if}
						</div>
					</div>
				{/each}
			</div>

		{:else if editorMode === 'markdown'}
			<!-- Markdown Editor -->
			<div class="grid grid-cols-2 gap-8 h-full">
				<!-- Editor Panel -->
				<div class="border-r border-gray-200 pr-4">
					<h3 class="text-sm font-medium text-gray-700 mb-2">Markdown Source</h3>
					<textarea
						bind:value={markdownContent}
						placeholder="# Start writing in Markdown...

## Features
- **Bold text**
- *Italic text*
- `Code snippets`
- [Links](https://example.com)

## Math Expressions
You can write natural math expressions:
- integral of x squared dx
- fraction a over b
- square root of 25"
						class="w-full h-96 p-3 border border-gray-300 rounded-md font-mono text-sm resize-none focus:ring-blue-500 focus:border-blue-500"
					></textarea>
				</div>

				<!-- Preview Panel -->
				<div class="pl-4">
					<h3 class="text-sm font-medium text-gray-700 mb-2">Live Preview</h3>
					<div class="w-full h-96 p-3 border border-gray-300 rounded-md bg-gray-50 overflow-y-auto prose prose-sm max-w-none">
						<!-- This would be a markdown renderer -->
						<div class="whitespace-pre-wrap font-mono text-sm">
							{markdownContent || '# Start writing in Markdown...\n\n## Features\n- **Bold text**\n- *Italic text*\n- `Code snippets`\n- [Links](https://example.com)\n\n## Math Expressions\nYou can write natural math expressions:\n- integral of x squared dx\n- fraction a over b\n- square root of 25'}
						</div>
					</div>
				</div>
			</div>

		{:else if editorMode === 'richtext'}
			<!-- Rich Text Editor (Traditional Word-like) -->
			<div class="border border-gray-300 rounded-lg overflow-hidden">
				<div class="bg-gray-50 px-4 py-2 border-b border-gray-200">
					<div class="flex items-center space-x-2 text-sm">
						<select class="border border-gray-300 rounded px-2 py-1">
							<option>Arial</option>
							<option>Times New Roman</option>
							<option>Inter</option>
						</select>
						<select class="border border-gray-300 rounded px-2 py-1">
							<option>12pt</option>
							<option>14pt</option>
							<option>16pt</option>
							<option>18pt</option>
						</select>
						<div class="w-px h-6 bg-gray-300"></div>
						<button class="p-1 hover:bg-white rounded" title="Bold"><strong>B</strong></button>
						<button class="p-1 hover:bg-white rounded" title="Italic"><em>I</em></button>
						<button class="p-1 hover:bg-white rounded" title="Underline"><u>U</u></button>
					</div>
				</div>
				<div
					contenteditable="true"
					class="min-h-96 p-4 focus:outline-none"
					placeholder="Start writing..."
				>
					<h1>Welcome to TPT Rich Text Editor</h1>
					<p>This is a traditional word processor interface with all the formatting options you expect.</p>
					<p>You can also insert <strong>natural math expressions</strong> like "integral from 0 to π of sin(x) dx" which will be automatically converted to proper mathematical notation.</p>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Auto-resize textarea */
	textarea {
		min-height: 1.5em;
	}

	.prose {
		color: #374151;
	}

	.prose strong {
		font-weight: 600;
	}

	.prose em {
		font-style: italic;
	}

	.prose code {
		background-color: #f3f4f6;
		padding: 0.125rem 0.25rem;
		border-radius: 0.25rem;
		font-size: 0.875em;
	}

	.prose h1 {
		font-size: 2rem;
		font-weight: 700;
		margin-bottom: 1rem;
	}

	.prose h2 {
		font-size: 1.5rem;
		font-weight: 600;
		margin-bottom: 0.75rem;
	}
</style>
