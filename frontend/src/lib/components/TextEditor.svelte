<script>
	import { onMount, createEventDispatcher } from 'svelte';
	import { jsPDF } from 'jspdf';
	import { createDocument, updateDocument, getDocuments, getDocument, getDocumentVersions, restoreDocumentVersion, deleteDocument } from '../api.js';
	import SpeechService from '../services/speech.js';

	// AI Services
	import { aiService } from '../services/ai.js';

	export const editorMode = 'blocks'; // 'blocks', 'markdown', 'richtext'
	export const documentId = null; // For loading existing documents

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
	let showMathHelp = false;

	// Rich text editor state
	let editorElement = null;
	let selectedFontFamily = 'Arial';
	let selectedFontSize = '12pt';
	let isBold = false;
	let isItalic = false;
	let isUnderline = false;
	let textAlign = 'left';
	let showFindReplaceDialog = false;
	let findText = '';
	let replaceText = '';
	let findResults = [];
	let currentFindIndex = -1;

	// Speech functionality
	let isReadingAloud = false;
	let speechSettings = null;
	let availableVoices = [];
	let selectedVoiceModel = null;

	// AI functionality
	let showAIAssistant = false;
	let aiSuggestions = [];
	let isGeneratingAI = false;
	let selectedText = '';
	let aiContext = 'general';
	let documentSummary = '';
	let isGeneratingSummary = false;
	let textAnalysis = null;
	let showTextAnalysis = false;

	// Initialize speech functionality on mount
	onMount(async () => {
		try {
			// Load available voice models
			availableVoices = await SpeechService.getAvailableModels('tts');

			// Load user speech settings
			speechSettings = await SpeechService.getSpeechSettings();

			// Set default voice model
			if (availableVoices.length > 0) {
				selectedVoiceModel = availableVoices[0];
			}
		} catch (error) {
			console.error('Failed to initialize speech:', error);
		}
	});

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
				editorMode = response.content.mode || 'blocks';
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

	// Rich text editor functions
	function applyBold() {
		document.execCommand('bold');
		updateFormattingState();
	}

	function applyItalic() {
		document.execCommand('italic');
		updateFormattingState();
	}

	function applyUnderline() {
		document.execCommand('underline');
		updateFormattingState();
	}

	function applyFontFamily() {
		document.execCommand('fontName', false, selectedFontFamily);
	}

	function applyFontSize() {
		document.execCommand('fontSize', false, selectedFontSize.replace('pt', ''));
	}

	function applyAlignLeft() {
		document.execCommand('justifyLeft');
		textAlign = 'left';
	}

	function applyAlignCenter() {
		document.execCommand('justifyCenter');
		textAlign = 'center';
	}

	function applyAlignRight() {
		document.execCommand('justifyRight');
		textAlign = 'right';
	}

	function insertLink() {
		const url = prompt('Enter URL:');
		if (url) {
			document.execCommand('createLink', false, url);
		}
	}

	function showFindReplace() {
		showFindReplaceDialog = true;
		findText = '';
		replaceText = '';
		findResults = [];
		currentFindIndex = -1;
	}

	function performFind() {
		if (!editorElement || !findText) return;

		findResults = [];
		currentFindIndex = -1;

		const content = editorElement.innerHTML;
		const regex = new RegExp(findText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi');
		let match;

		while ((match = regex.exec(content)) !== null) {
			findResults.push({
				start: match.index,
				end: match.index + match[0].length,
				text: match[0]
			});
		}

		if (findResults.length > 0) {
			selectFindResult(0);
		}
	}

	function selectFindResult(index) {
		if (findResults.length === 0 || index < 0 || index >= findResults.length) return;

		currentFindIndex = index;
		const result = findResults[index];

		// Create a range and select the text
		const range = document.createRange();
		const textNode = editorElement.firstChild || editorElement;

		// This is a simplified implementation - a full implementation would need more complex text node traversal
		range.setStart(textNode, Math.min(result.start, textNode.textContent.length));
		range.setEnd(textNode, Math.min(result.end, textNode.textContent.length));

		const selection = window.getSelection();
		selection.removeAllRanges();
		selection.addRange(range);
	}

	function performReplace() {
		if (!editorElement || findResults.length === 0 || currentFindIndex === -1) return;

		document.execCommand('insertText', false, replaceText);
		findResults.splice(currentFindIndex, 1);

		if (findResults.length > 0) {
			selectFindResult(Math.min(currentFindIndex, findResults.length - 1));
		} else {
			currentFindIndex = -1;
		}
	}

	function replaceAll() {
		if (!editorElement || !findText || !replaceText) return;

		let content = editorElement.innerHTML;
		const regex = new RegExp(findText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi');
		content = content.replace(regex, replaceText);
		editorElement.innerHTML = content;

		findResults = [];
		currentFindIndex = -1;
		showFindReplaceDialog = false;
		markAsChanged();
	}

	function handleContentChange() {
		markAsChanged();
	}

	function updateFormattingState() {
		if (!editorElement) return;

		isBold = document.queryCommandState('bold');
		isItalic = document.queryCommandState('italic');
		isUnderline = document.queryCommandState('underline');

		// Update font family and size (simplified)
		selectedFontFamily = document.queryCommandValue('fontName') || 'Arial';
		const fontSizeValue = document.queryCommandValue('fontSize');
		selectedFontSize = fontSizeValue ? `${fontSizeValue}pt` : '12pt';

		// Update text alignment
		if (document.queryCommandState('justifyLeft')) textAlign = 'left';
		else if (document.queryCommandState('justifyCenter')) textAlign = 'center';
		else if (document.queryCommandState('justifyRight')) textAlign = 'right';
	}

	function handleRichTextKeyDown(event) {
		// Rich text specific shortcuts
		if (event.ctrlKey || event.metaKey) {
			switch (event.key.toLowerCase()) {
				case 'b':
					event.preventDefault();
					applyBold();
					break;
				case 'i':
					event.preventDefault();
					applyItalic();
					break;
				case 'u':
					event.preventDefault();
					applyUnderline();
					break;
				case 'l':
					event.preventDefault();
					applyAlignLeft();
					break;
				case 'e':
					event.preventDefault();
					applyAlignCenter();
					break;
				case 'r':
					event.preventDefault();
					applyAlignRight();
					break;
				case 'k':
					event.preventDefault();
					insertLink();
					break;
			}
		}
	}

	// Speech functionality
	async function readAloud() {
		if (isReadingAloud || !selectedVoiceModel) {
			return;
		}

		isReadingAloud = true;

		try {
			// Extract text content from blocks
			let fullText = '';

			if (editorMode === 'blocks') {
				fullText = blocks.map(block => block.content).join('\n\n');
			} else if (editorMode === 'markdown') {
				fullText = markdownContent;
			} else if (editorMode === 'richtext') {
				fullText = editorElement?.innerText || '';
			}

			if (!fullText.trim()) {
				alert('No text to read aloud');
				isReadingAloud = false;
				return;
			}

			// Split text into manageable chunks (for API limits)
			const chunks = splitTextIntoChunks(fullText, 4000); // 4K character chunks

			for (const chunk of chunks) {
				if (!isReadingAloud) break; // Allow stopping

				try {
					await SpeechService.textToSpeech(chunk, selectedVoiceModel.id, {
						voice: speechSettings?.defaultVoice || 'alloy',
						language: speechSettings?.defaultLanguage || 'en',
						speed: speechSettings?.ttsSpeed || 1.0,
						volume: speechSettings?.ttsVolume || 1.0
					});

					// Small pause between chunks
					await new Promise(resolve => setTimeout(resolve, 500));
				} catch (error) {
					console.error('TTS chunk error:', error);
					// Continue with next chunk instead of stopping
				}
			}

		} catch (error) {
			console.error('Read aloud error:', error);
			alert('Failed to read document aloud: ' + error.message);
		} finally {
			isReadingAloud = false;
		}
	}

	function splitTextIntoChunks(text, maxLength) {
		const chunks = [];
		let currentChunk = '';

		const sentences = text.split(/[.!?]+/).filter(s => s.trim());

		for (const sentence of sentences) {
			const trimmedSentence = sentence.trim();
			if (!trimmedSentence) continue;

			if (currentChunk.length + trimmedSentence.length + 1 <= maxLength) {
				currentChunk += (currentChunk ? '. ' : '') + trimmedSentence;
			} else {
				if (currentChunk) {
					chunks.push(currentChunk + '.');
				}
				currentChunk = trimmedSentence;
			}
		}

		if (currentChunk) {
			chunks.push(currentChunk + '.');
		}

		return chunks;
	}

	// AI Writing Assistant Functions
	async function getSelectedText() {
		let text = '';

		if (editorMode === 'blocks') {
			// Get text from selected block
			if (selectedBlockIndex >= 0 && selectedBlockIndex < blocks.length) {
				text = blocks[selectedBlockIndex].content;
			}
		} else if (editorMode === 'markdown') {
			text = markdownContent;
		} else if (editorMode === 'richtext') {
			// Get selected text from rich text editor
			const selection = window.getSelection();
			if (selection.rangeCount > 0) {
				text = selection.toString();
			}
			if (!text && editorElement) {
				text = editorElement.innerText;
			}
		}

		return text || '';
	}

	async function getAIWritingSuggestions() {
		const text = await getSelectedText();
		if (!text.trim()) {
			alert('Please select some text first');
			return;
		}

		isGeneratingAI = true;
		showAIAssistant = true;

		try {
			const suggestions = await aiService.getWritingSuggestions(text, aiContext);
			if (suggestions) {
				aiSuggestions = suggestions.suggestions || [];
			} else {
				aiSuggestions = [];
			}
		} catch (error) {
			console.error('Failed to get AI suggestions:', error);
			aiSuggestions = [];
			alert('Failed to get AI suggestions. Please try again.');
		} finally {
			isGeneratingAI = false;
		}
	}

	async function improveSelectedText(improvementType = 'general') {
		const text = await getSelectedText();
		if (!text.trim()) {
			alert('Please select some text first');
			return;
		}

		isGeneratingAI = true;

		try {
			const improvedText = await aiService.improveText(text, improvementType);
			if (improvedText) {
				// Replace the selected text with improved version
				replaceSelectedText(improvedText);
			}
		} catch (error) {
			console.error('Failed to improve text:', error);
			alert('Failed to improve text. Please try again.');
		} finally {
			isGeneratingAI = false;
		}
	}

	async function continueWriting() {
		const text = await getSelectedText();
		if (!text.trim()) {
			alert('Please select some text first');
			return;
		}

		isGeneratingAI = true;

		try {
			const continuation = await aiService.continueWriting(text, 'matching', 'medium');
			if (continuation) {
				// Insert continuation after selected text
				insertTextAfterSelection(continuation);
			}
		} catch (error) {
			console.error('Failed to continue writing:', error);
			alert('Failed to continue writing. Please try again.');
		} finally {
			isGeneratingAI = false;
		}
	}

	async function generateDocumentSummary() {
		let fullText = '';

		if (editorMode === 'blocks') {
			fullText = blocks.map(block => block.content).join('\n\n');
		} else if (editorMode === 'markdown') {
			fullText = markdownContent;
		} else if (editorMode === 'richtext') {
			fullText = editorElement?.innerText || '';
		}

		if (!fullText.trim()) {
			alert('No text to summarize');
			return;
		}

		isGeneratingSummary = true;

		try {
			const summary = await aiService.summarizeDocument(fullText, 'concise');
			if (summary) {
				documentSummary = summary;
			}
		} catch (error) {
			console.error('Failed to generate summary:', error);
			alert('Failed to generate document summary. Please try again.');
		} finally {
			isGeneratingSummary = false;
		}
	}

	async function analyzeDocument() {
		let fullText = '';

		if (editorMode === 'blocks') {
			fullText = blocks.map(block => block.content).join('\n\n');
		} else if (editorMode === 'markdown') {
			fullText = markdownContent;
		} else if (editorMode === 'richtext') {
			fullText = editorElement?.innerText || '';
		}

		if (!fullText.trim()) {
			alert('No text to analyze');
			return;
		}

		try {
			const analysis = await aiService.analyzeText(fullText);
			if (analysis) {
				textAnalysis = analysis;
				showTextAnalysis = true;
			}
		} catch (error) {
			console.error('Failed to analyze text:', error);
			alert('Failed to analyze document. Please try again.');
		}
	}

	function replaceSelectedText(newText) {
		if (editorMode === 'blocks') {
			if (selectedBlockIndex >= 0 && selectedBlockIndex < blocks.length) {
				blocks[selectedBlockIndex].content = newText;
				blocks = [...blocks]; // Trigger reactivity
				markAsChanged();
			}
		} else if (editorMode === 'markdown') {
			markdownContent = newText;
			markAsChanged();
		} else if (editorMode === 'richtext') {
			document.execCommand('insertText', false, newText);
		}
	}

	function insertTextAfterSelection(text) {
		if (editorMode === 'blocks') {
			if (selectedBlockIndex >= 0 && selectedBlockIndex < blocks.length) {
				blocks[selectedBlockIndex].content += ' ' + text;
				blocks = [...blocks]; // Trigger reactivity
				markAsChanged();
			}
		} else if (editorMode === 'markdown') {
			markdownContent += ' ' + text;
			markAsChanged();
		} else if (editorMode === 'richtext') {
			document.execCommand('insertText', false, ' ' + text);
		}
	}

	function applyAISuggestion(suggestion) {
		replaceSelectedText(suggestion);
		showAIAssistant = false;
		aiSuggestions = [];
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

			<button
				class="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700"
				on:click={() => showMathHelp = true}
				title="Math Help & Examples"
			>
				Math Help
			</button>

			{#if availableVoices.length > 0}
				<button
					class="px-3 py-1 text-sm bg-purple-600 text-white rounded hover:bg-purple-700 disabled:opacity-50"
					on:click={readAloud}
					disabled={isReadingAloud}
					title="Read document aloud"
				>
					{isReadingAloud ? '🔊 Reading...' : '📖 Read Aloud'}
				</button>
			{/if}

			<!-- AI Writing Assistant -->
			<div class="border-l border-gray-300 pl-4 ml-4 flex items-center space-x-2">
				<span class="text-xs text-gray-500 font-medium">AI</span>

				<!-- AI Features with Enhanced Tooltips -->
				<div class="relative">
					<button
						class="px-3 py-1 text-sm bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-50 relative group"
						on:click={getAIWritingSuggestions}
						disabled={isGeneratingAI}
						title="Get AI writing suggestions for grammar, style, and clarity improvements"
						aria-label="Get AI writing suggestions"
						type="button"
					>
						{isGeneratingAI ? '💭' : '✨ Suggest'}
						<!-- Enhanced Tooltip -->
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none">
							✨ AI Writing Assistant
							<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
						</div>
					</button>

					<button
						class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 relative group ml-1"
						on:click={continueWriting}
						disabled={isGeneratingAI}
						title="AI content continuation - select text to continue from that point"
						aria-label="Continue writing with AI"
						type="button"
					>
						📝 Continue
						<!-- Enhanced Tooltip -->
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none">
							📝 Continue Writing
							<div class="text-xs text-gray-300 mt-1">Select text to continue from that point</div>
							<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
						</div>
					</button>

					<button
						class="px-3 py-1 text-sm bg-yellow-600 text-white rounded hover:bg-yellow-700 disabled:opacity-50 relative group ml-1"
						on:click={generateDocumentSummary}
						disabled={isGeneratingSummary}
						title="Generate intelligent document summaries - works best with longer documents"
						aria-label="Generate document summary"
						type="button"
					>
						{isGeneratingSummary ? '📋' : '📄 Summary'}
						<!-- Enhanced Tooltip -->
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none">
							📄 Document Summary
							<div class="text-xs text-gray-300 mt-1">Best with 500+ words</div>
							<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
						</div>
					</button>

					<button
						class="px-3 py-1 text-sm bg-gray-600 text-white rounded hover:bg-gray-700 relative group ml-1"
						on:click={analyzeDocument}
						title="Analyze document for readability, sentiment, and key insights"
						aria-label="Analyze document"
						type="button"
					>
						📊 Analyze
						<!-- Enhanced Tooltip -->
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none">
							📊 Document Analysis
							<div class="text-xs text-gray-300 mt-1">Readability, sentiment & key phrases</div>
							<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
						</div>
					</button>
				</div>
			</div>
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
						<select bind:value={selectedFontFamily} on:change={applyFontFamily} class="border border-gray-300 rounded px-2 py-1">
							<option value="Arial">Arial</option>
							<option value="Times New Roman">Times New Roman</option>
							<option value="Inter">Inter</option>
							<option value="Georgia">Georgia</option>
							<option value="Verdana">Verdana</option>
						</select>
						<select bind:value={selectedFontSize} on:change={applyFontSize} class="border border-gray-300 rounded px-2 py-1">
							<option value="12pt">12pt</option>
							<option value="14pt">14pt</option>
							<option value="16pt">16pt</option>
							<option value="18pt">18pt</option>
							<option value="24pt">24pt</option>
							<option value="32pt">32pt</option>
						</select>
						<div class="w-px h-6 bg-gray-300"></div>
						<button on:click={applyBold} class="p-1 hover:bg-white rounded {isBold ? 'bg-blue-100' : ''}" title="Bold">
							<strong>B</strong>
						</button>
						<button on:click={applyItalic} class="p-1 hover:bg-white rounded {isItalic ? 'bg-blue-100' : ''}" title="Italic">
							<em>I</em>
						</button>
						<button on:click={applyUnderline} class="p-1 hover:bg-white rounded {isUnderline ? 'bg-blue-100' : ''}" title="Underline">
							<u>U</u>
						</button>
						<div class="w-px h-6 bg-gray-300"></div>
						<button on:click={applyAlignLeft} class="p-1 hover:bg-white rounded {textAlign === 'left' ? 'bg-blue-100' : ''}" title="Align Left">
							⬅️
						</button>
						<button on:click={applyAlignCenter} class="p-1 hover:bg-white rounded {textAlign === 'center' ? 'bg-blue-100' : ''}" title="Align Center">
							⬌
						</button>
						<button on:click={applyAlignRight} class="p-1 hover:bg-white rounded {textAlign === 'right' ? 'bg-blue-100' : ''}" title="Align Right">
							➡️
						</button>
						<div class="w-px h-6 bg-gray-300"></div>
						<button on:click={insertLink} class="p-1 hover:bg-white rounded" title="Insert Link">
							🔗
						</button>
						<button on:click={showFindReplace} class="p-1 hover:bg-white rounded" title="Find & Replace">
							🔍
						</button>
					</div>
				</div>
				<div
					bind:this={editorElement}
					contenteditable="true"
					class="min-h-96 p-4 focus:outline-none"
					placeholder="Start writing..."
					on:input={handleContentChange}
					on:keydown={handleRichTextKeyDown}
					on:mouseup={updateFormattingState}
					on:keyup={updateFormattingState}
				>
					<h1>Welcome to TPT Rich Text Editor</h1>
					<p>This is a traditional word processor interface with all the formatting options you expect.</p>
					<p>You can also insert <strong>natural math expressions</strong> like "integral from 0 to π of sin(x) dx" which will be automatically converted to proper mathematical notation.</p>
				</div>
			</div>
		{/if}
	</div>
</div>

<!-- Math Help Modal -->
{#if showMathHelp}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">Natural Math Input - Better Than LaTeX!</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showMathHelp = false}
				>
					×
				</button>
			</div>

			<div class="space-y-6">
				<!-- Introduction -->
				<div class="bg-blue-50 p-4 rounded-lg">
					<h4 class="text-lg font-medium text-blue-900 mb-2">✨ Revolutionary Math Input</h4>
					<p class="text-blue-800">
						Forget complex LaTeX syntax! TPT Titan understands natural language math expressions.
						Simply type what you would say out loud, and watch it transform into beautiful mathematical notation.
					</p>
				</div>

				<!-- Basic Examples -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">📚 Basic Expressions</h4>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div class="bg-gray-50 p-3 rounded">
							<div class="text-sm text-gray-600 mb-1">Type this:</div>
							<code class="bg-white px-2 py-1 rounded text-sm">fraction a over b</code>
							<div class="text-sm text-gray-600 mt-1">Becomes:</div>
							<div class="font-serif text-lg">a/b</div>
						</div>
						<div class="bg-gray-50 p-3 rounded">
							<div class="text-sm text-gray-600 mb-1">Type this:</div>
							<code class="bg-white px-2 py-1 rounded text-sm">square root of x</code>
							<div class="text-sm text-gray-600 mt-1">Becomes:</div>
							<div class="font-serif text-lg">√x</div>
						</div>
						<div class="bg-gray-50 p-3 rounded">
							<div class="text-sm text-gray-600 mb-1">Type this:</div>
							<code class="bg-white px-2 py-1 rounded text-sm">pi times r squared</code>
							<div class="text-sm text-gray-600 mt-1">Becomes:</div>
							<div class="font-serif text-lg">π × r²</div>
						</div>
						<div class="bg-gray-50 p-3 rounded">
							<div class="text-sm text-gray-600 mb-1">Type this:</div>
							<code class="bg-white px-2 py-1 rounded text-sm">alpha beta gamma</code>
							<div class="text-sm text-gray-600 mt-1">Becomes:</div>
							<div class="font-serif text-lg">αβγ</div>
						</div>
					</div>
				</div>

				<!-- Advanced Examples -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">🔬 Advanced Mathematics</h4>
					<div class="space-y-4">
						<div class="bg-purple-50 p-4 rounded">
							<div class="text-sm text-purple-700 mb-2">Integrals:</div>
							<code class="bg-white px-3 py-2 rounded text-sm block mb-2">integral from 0 to infinity of e to the power of negative x squared dx</code>
							<div class="text-sm text-purple-700 mb-1">Becomes:</div>
							<div class="font-serif text-xl">∫₀^∞ e^(-x²) dx</div>
						</div>

						<div class="bg-green-50 p-4 rounded">
							<div class="text-sm text-green-700 mb-2">Summations:</div>
							<code class="bg-white px-3 py-2 rounded text-sm block mb-2">sum from i equals 1 to n of x sub i</code>
							<div class="text-sm text-green-700 mb-1">Becomes:</div>
							<div class="font-serif text-xl">∑_{i=1}^n x_i</div>
						</div>

						<div class="bg-orange-50 p-4 rounded">
							<div class="text-sm text-orange-700 mb-2">Complex fractions:</div>
							<code class="bg-white px-3 py-2 rounded text-sm block mb-2">fraction x plus y over z minus w</code>
							<div class="text-sm text-orange-700 mb-1">Becomes:</div>
							<div class="font-serif text-xl">x+y/z-w</div>
						</div>
					</div>
				</div>

				<!-- Quick Reference -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">🚀 Quick Reference</h4>
					<div class="grid grid-cols-2 md:grid-cols-3 gap-3 text-sm">
						<div><strong>fractions:</strong> "fraction a over b"</div>
						<div><strong>roots:</strong> "square root of", "cube root of"</div>
						<div><strong>integrals:</strong> "integral of", "integral from a to b of"</div>
						<div><strong>summations:</strong> "sum from i=1 to n of"</div>
						<div><strong>greek:</strong> "alpha", "beta", "gamma", "pi", "sigma"</div>
						<div><strong>operators:</strong> "times", "divided by", "plus or minus"</div>
						<div><strong>superscripts:</strong> "squared", "cubed", "to the power of"</div>
						<div><strong>subscripts:</strong> "sub", "to the"</div>
						<div><strong>symbols:</strong> "therefore", "because", "infinity"</div>
					</div>
				</div>

				<!-- Try It Section -->
				<div class="bg-yellow-50 p-4 rounded">
					<h4 class="text-lg font-medium text-yellow-900 mb-2">🎯 Try It Now!</h4>
					<p class="text-yellow-800 mb-3">
						Create a new math block and try typing any of these expressions:
					</p>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm">
						<button class="text-left p-2 bg-white rounded hover:bg-gray-50" on:click={() => { showMathHelp = false; /* Would insert example */ }}>
							"fraction x squared plus y squared over z"
						</button>
						<button class="text-left p-2 bg-white rounded hover:bg-gray-50" on:click={() => { showMathHelp = false; /* Would insert example */ }}>
							"integral from negative infinity to infinity of e to the power of negative x squared over square root of pi dx"
						</button>
						<button class="text-left p-2 bg-white rounded hover:bg-gray-50" on:click={() => { showMathHelp = false; /* Would insert example */ }}>
							"sum from k equals 0 to infinity of fraction 1 over k factorial"
						</button>
						<button class="text-left p-2 bg-white rounded hover:bg-gray-50" on:click={() => { showMathHelp = false; /* Would insert example */ }}>
							"limit as x approaches 0 of fraction sin x over x equals 1"
						</button>
					</div>
				</div>

				<!-- Export Options -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">📤 Export Your Math</h4>
					<p class="text-gray-600 mb-3">
						Once you've created beautiful math expressions, export them to various formats:
					</p>
					<div class="grid grid-cols-2 md:grid-cols-4 gap-3">
						<div class="text-center p-3 bg-red-50 rounded">
							<div class="text-red-600 font-medium">PDF</div>
							<div class="text-xs text-red-600">Vector graphics</div>
						</div>
						<div class="text-center p-3 bg-blue-50 rounded">
							<div class="text-blue-600 font-medium">LaTeX</div>
							<div class="text-xs text-blue-600">Source code</div>
						</div>
						<div class="text-center p-3 bg-green-50 rounded">
							<div class="text-green-600 font-medium">MathML</div>
							<div class="text-xs text-green-600">Web standard</div>
						</div>
						<div class="text-center p-3 bg-purple-50 rounded">
							<div class="text-purple-600 font-medium">SVG</div>
							<div class="text-xs text-purple-600">Scalable images</div>
						</div>
					</div>
				</div>
			</div>

			<div class="mt-6 flex justify-end">
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
					on:click={() => showMathHelp = false}
				>
					Got it!
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Find & Replace Modal -->
{#if showFindReplaceDialog}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-md">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">Find & Replace</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showFindReplaceDialog = false}
				>
					×
				</button>
			</div>

			<div class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Find</label>
					<input
						type="text"
						bind:value={findText}
						placeholder="Enter text to find..."
						class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
						on:keydown={(e) => {
							if (e.key === 'Enter') performFind();
							if (e.key === 'Escape') showFindReplaceDialog = false;
						}}
					/>
				</div>

				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Replace with</label>
					<input
						type="text"
						bind:value={replaceText}
						placeholder="Enter replacement text..."
						class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
						on:keydown={(e) => {
							if (e.key === 'Enter') performReplace();
							if (e.key === 'Escape') showFindReplaceDialog = false;
						}}
					/>
				</div>

				{#if findResults.length > 0}
					<div class="text-sm text-gray-600">
						Found {findResults.length} occurrence{findResults.length !== 1 ? 's' : ''}
						{#if currentFindIndex >= 0}
							(currently at {currentFindIndex + 1})
						{/if}
					</div>
				{/if}
			</div>

			<div class="mt-6 flex justify-end space-x-3">
				<button
					class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50"
					on:click={() => showFindReplaceDialog = false}
				>
					Close
				</button>
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
					on:click={performFind}
				>
					Find
				</button>
				{#if findResults.length > 0}
					<button
						class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
						on:click={performReplace}
					>
						Replace
					</button>
					<button
						class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700"
						on:click={replaceAll}
					>
						Replace All
					</button>
				{/if}
			</div>
		</div>
	</div>
{/if}

<!-- AI Writing Assistant Modal -->
{#if showAIAssistant}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">✨ AI Writing Suggestions</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => { showAIAssistant = false; aiSuggestions = []; }}
				>
					×
				</button>
			</div>

			{#if aiSuggestions.length === 0}
				<div class="text-center py-8">
					<div class="text-4xl mb-4">💭</div>
					<p class="text-gray-600">Generating suggestions...</p>
				</div>
			{:else}
				<div class="space-y-4">
					{#each aiSuggestions as suggestion, index}
						<div class="border border-gray-200 rounded-lg p-4">
							<div class="flex items-start justify-between mb-2">
								<h4 class="font-medium text-gray-900">Suggestion {index + 1}</h4>
								<button
									class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
									on:click={() => applyAISuggestion(suggestion)}
								>
									Apply
								</button>
							</div>
							<p class="text-gray-700 whitespace-pre-wrap">{suggestion}</p>
						</div>
					{/each}
				</div>
			{/if}

			<div class="mt-6 flex justify-end">
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => { showAIAssistant = false; aiSuggestions = []; }}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Document Summary Modal -->
{#if documentSummary}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-3xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">📄 Document Summary</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => documentSummary = ''}
				>
					×
				</button>
			</div>

			<div class="prose prose-sm max-w-none">
				<div class="whitespace-pre-wrap text-gray-700">{documentSummary}</div>
			</div>

			<div class="mt-6 flex justify-end space-x-3">
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
					on:click={() => navigator.clipboard.writeText(documentSummary)}
				>
					Copy to Clipboard
				</button>
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => documentSummary = ''}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Text Analysis Modal -->
{#if showTextAnalysis}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">📊 Text Analysis</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => { showTextAnalysis = false; textAnalysis = null; }}
				>
					×
				</button>
			</div>

			{#if textAnalysis}
				<div class="grid grid-cols-2 gap-6">
					<div class="space-y-4">
						<h4 class="font-medium text-gray-900">Statistics</h4>
						<div class="space-y-2 text-sm">
							<div class="flex justify-between">
								<span class="text-gray-600">Words:</span>
								<span class="font-medium">{textAnalysis.word_count || 0}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Characters:</span>
								<span class="font-medium">{textAnalysis.char_count || 0}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Sentences:</span>
								<span class="font-medium">{textAnalysis.sentence_count || 0}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Paragraphs:</span>
								<span class="font-medium">{textAnalysis.paragraph_count || 0}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Reading Time:</span>
								<span class="font-medium">{textAnalysis.reading_time || 'N/A'}</span>
							</div>
						</div>
					</div>

					<div class="space-y-4">
						<h4 class="font-medium text-gray-900">Readability</h4>
						<div class="space-y-2 text-sm">
							<div class="flex justify-between">
								<span class="text-gray-600">Grade Level:</span>
								<span class="font-medium">{textAnalysis.grade_level || 'N/A'}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Readability:</span>
								<span class="font-medium">{textAnalysis.readability_score || 'N/A'}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Complexity:</span>
								<span class="font-medium">{textAnalysis.complexity || 'N/A'}</span>
							</div>
						</div>
					</div>
				</div>

				{#if textAnalysis.key_phrases && textAnalysis.key_phrases.length > 0}
					<div class="mt-6">
						<h4 class="font-medium text-gray-900 mb-3">Key Phrases</h4>
						<div class="flex flex-wrap gap-2">
							{#each textAnalysis.key_phrases as phrase}
								<span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
									{phrase}
								</span>
							{/each}
						</div>
					</div>
				{/if}

				{#if textAnalysis.sentiment}
					<div class="mt-6">
						<h4 class="font-medium text-gray-900 mb-3">Sentiment Analysis</h4>
						<div class="flex items-center space-x-4">
							<div class="flex-1 bg-gray-200 rounded-full h-2">
								<div
									class="bg-green-500 h-2 rounded-full"
									style="width: {textAnalysis.sentiment.positive || 0}%"
								></div>
							</div>
							<span class="text-sm text-gray-600">
								Positive: {textAnalysis.sentiment.positive || 0}%
							</span>
						</div>
						<div class="flex items-center space-x-4 mt-2">
							<div class="flex-1 bg-gray-200 rounded-full h-2">
								<div
									class="bg-red-500 h-2 rounded-full"
									style="width: {textAnalysis.sentiment.negative || 0}%"
								></div>
							</div>
							<span class="text-sm text-gray-600">
								Negative: {textAnalysis.sentiment.negative || 0}%
							</span>
						</div>
					</div>
				{/if}
			{:else}
				<div class="text-center py-8">
					<p class="text-gray-600">No analysis data available</p>
				</div>
			{/if}

			<div class="mt-6 flex justify-end">
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => { showTextAnalysis = false; textAnalysis = null; }}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	/* Auto-resize textarea */
	textarea {
		min-height: 1.5em;
	}

	.prose {
		color: #374151;
	}


</style>
