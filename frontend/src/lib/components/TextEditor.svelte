<!-- frontend/src/lib/components/TextEditor.svelte -->
<script>
	import { onMount } from 'svelte';
	import { jsPDF } from 'jspdf';
	import { createDocument, updateDocument, getDocuments, getDocument, getDocumentVersions, restoreDocumentVersion } from '../api.js';
	import SpeechService from '../services/speech.js';
	import { aiService } from '../services/ai.js';
	import EditorHistory, { createDebouncedPush, handleUndoRedoKeyboard } from '../utils/editor-history.js';
	import { 
		searchInBlocks, 
		searchInMarkdown, 
		searchInRichText,
		replaceInBlocks,
		replaceInMarkdown,
		replaceAllInBlocks,
		highlightResult,
		clearHighlights 
	} from '../utils/editor-search.js';
	import { exportDocument } from '../utils/editor-export.js';




	import TextEditorToolbar from './TextEditorToolbar.svelte';

	import TextEditorBlockEditor from './TextEditorBlockEditor.svelte';
	import TextEditorMarkdownEditor from './TextEditorMarkdownEditor.svelte';
	import TextEditorRichText from './TextEditorRichText.svelte';
	import TextEditorModals from './TextEditorModals.svelte';

	// ── Props ──────────────────────────────────────────────────────
	export let editorMode = 'blocks'; // 'blocks' | 'markdown' | 'richtext'
	export let documentId = null;

	// ── Document state ─────────────────────────────────────────────
	let currentDocument = null;
	let documentTitle = 'Untitled Document';
	let isSaving = false;
	let isLoading = false;
	let saveStatus = '';
	let hasUnsavedChanges = false;
	let autoSaveTimer = null;

	// ── Editor content ─────────────────────────────────────────────
	let blocks = [];
	let markdownContent = '';
	let selectedBlockIndex = 0;

	// ── Rich-text DOM ref ──────────────────────────────────────────
	let richTextEditorElement = null;

	// ── Modal visibility ───────────────────────────────────────────
	let showDocumentList = false;
	let showVersionHistory = false;
	let showMathHelp = false;
	let showFindReplaceDialog = false;
	let showAIAssistant = false;

	// ── Modal data ─────────────────────────────────────────────────
	let documentList = [];
	let versionHistory = [];
	let aiSuggestions = [];
	let documentSummary = '';
	let showTextAnalysis = false;
	let textAnalysis = null;

	// ── Find & replace ─────────────────────────────────────────────
	let findText = '';
	let replaceText = '';
	let findResults = [];
	let currentFindIndex = -1;
	let isCaseSensitive = false;


	// ── Speech ─────────────────────────────────────────────────────
	let isReadingAloud = false;
	let speechSettings = null;
	let availableVoices = [];
	let selectedVoiceModel = null;

	// ── AI ─────────────────────────────────────────────────────────
	let isGeneratingAI = false;
	let isGeneratingSummary = false;
	let aiContext = 'general';

	// ── History (Undo/Redo) ────────────────────────────────────────
	let editorHistory = new EditorHistory(100);
	let canUndo = false;
	let canRedo = false;
	let debouncedHistoryPush;

	// ── Lifecycle ──────────────────────────────────────────────────
	onMount(async () => {
		// Initialize history
		debouncedHistoryPush = createDebouncedPush(editorHistory, 500);
		saveToHistory();

		// Speech init
		try {
			availableVoices = await SpeechService.getAvailableModels('tts');
			speechSettings = await SpeechService.getSpeechSettings();
			if (availableVoices.length > 0) selectedVoiceModel = availableVoices[0];
		} catch (error) {
			console.error('Failed to initialize speech:', error);
		}

		// Default blocks
		if (blocks.length === 0) {
			blocks = [
				{ id: 1, type: 'heading1', content: 'Welcome to TPT Text Editor', properties: {} },
				{ id: 2, type: 'text', content: 'This revolutionary text editor combines the best of Notion, Markdown, and Microsoft Word with groundbreaking features.', properties: {} },
				{ id: 3, type: 'heading2', content: 'Natural Math Input (Better Than LaTeX!)', properties: {} },
				{ id: 4, type: 'math', content: 'integral from 0 to infinity of e to the power of negative x squared dx', properties: {} },
				{ id: 5, type: 'text', content: 'The equation above was created by simply typing it in plain English - no complex LaTeX syntax required!', properties: {} }
			];
		}
	});

	// ── History helpers ───────────────────────────────────────────
	function saveToHistory() {
		let content;
		if (editorMode === 'blocks') content = JSON.parse(JSON.stringify(blocks));
		else if (editorMode === 'markdown') content = markdownContent;
		else if (editorMode === 'richtext') content = richTextEditorElement?.innerHTML || '';
		
		debouncedHistoryPush({
			type: editorMode,
			content: content,
			metadata: { selectedBlockIndex }
		});
		updateHistoryState();
	}

	function updateHistoryState() {
		canUndo = editorHistory.canUndo();
		canRedo = editorHistory.canRedo();
	}

	function undo() {
		const state = editorHistory.undo();
		if (state) restoreFromHistory(state);
		updateHistoryState();
	}

	function redo() {
		const state = editorHistory.redo();
		if (state) restoreFromHistory(state);
		updateHistoryState();
	}

	function restoreFromHistory(state) {
		if (state.type === 'blocks') {
			blocks = JSON.parse(JSON.stringify(state.content));
			selectedBlockIndex = state.metadata?.selectedBlockIndex || 0;
		} else if (state.type === 'markdown') {
			markdownContent = state.content;
		} else if (state.type === 'richtext' && richTextEditorElement) {
			richTextEditorElement.innerHTML = state.content;
		}
		hasUnsavedChanges = true;
	}

	// ── Auto-save helper ───────────────────────────────────────────
	function markAsChanged() {
		hasUnsavedChanges = true;
		saveToHistory();
		if (autoSaveTimer) clearTimeout(autoSaveTimer);
		autoSaveTimer = setTimeout(() => {
			if (hasUnsavedChanges) saveDocument();
		}, 30000);
	}


	// ── Document CRUD ──────────────────────────────────────────────
	async function saveDocument() {
		if (isSaving) return;
		isSaving = true;
		saveStatus = 'Saving...';
		try {
			const documentData = {
				title: documentTitle,
				content_type: 'text',
				content: { mode: editorMode, blocks, markdown: markdownContent }
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
			setTimeout(() => { saveStatus = ''; }, 2000);
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
			if (response.content && typeof response.content === 'object') {
				if (response.content.blocks) blocks = response.content.blocks;
				if (response.content.markdown) markdownContent = response.content.markdown;
				if (response.content.mode) editorMode = response.content.mode || 'blocks';
			}
			showDocumentList = false;
			hasUnsavedChanges = false;
			// Reset history after loading
			editorHistory.clear();
			saveToHistory();
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
			await loadDocument(currentDocument);
			showVersionHistory = false;
		} catch (error) {
			console.error('Failed to restore version:', error);
		}
	}

	function createNewDocument() {
		currentDocument = null;
		documentTitle = 'Untitled Document';
		blocks = [{ id: 1, type: 'heading1', content: 'New Document', properties: {} }];
		markdownContent = '';
		hasUnsavedChanges = false;
		showDocumentList = false;
		showVersionHistory = false;
		// Reset history
		editorHistory.clear();
		saveToHistory();
	}


	// ── PDF Export ─────────────────────────────────────────────────
	function getFontSize(blockType) {
		const sizes = { heading1: 18, heading2: 16, heading3: 14, text: 12, list: 12, quote: 12, code: 10, math: 12, table: 12, image: 12 };
		return sizes[blockType] || 12;
	}

	function getLineHeight(blockType) {
		const heights = { heading1: 8, heading2: 7, heading3: 6, text: 6, list: 6, quote: 6, code: 5, math: 6, table: 6, image: 6 };
		return heights[blockType] || 6;
	}

	function exportToPDF() {
		const doc = new jsPDF();
		doc.setFontSize(20);
		doc.text(documentTitle || 'Untitled Document', 20, 30);

		let yPosition = 50;
		blocks.forEach(block => {
			doc.setFontSize(getFontSize(block.type));
			const content = block.content;
			const lines = doc.splitTextToSize(content, 170);
			lines.forEach(line => {
				if (yPosition > 270) { doc.addPage(); yPosition = 30; }
				doc.text(line, 20, yPosition);
				yPosition += getLineHeight(block.type);
			});
			yPosition += 10;
		});

		doc.save(`${documentTitle || 'document'}.pdf`);
	}

	// ── Block editor helpers ───────────────────────────────────────
	function handleAddBlock(event) {
		const afterIndex = event.detail;
		const newBlock = { id: Date.now() + Math.random(), type: 'text', content: '', properties: {} };
		blocks.splice(afterIndex + 1, 0, newBlock);
		blocks = [...blocks];
		selectedBlockIndex = afterIndex + 1;
		markAsChanged();
	}

	function handleDeleteBlock(event) {
		const blockIndex = event.detail;
		if (blocks.length > 1) {
			blocks.splice(blockIndex, 1);
			blocks = [...blocks];
			selectedBlockIndex = Math.max(0, blockIndex - 1);
			markAsChanged();
		}
	}

	function handleBlockContentChange(event) {
		const { index, content } = event.detail;
		blocks[index].content = content;
		blocks = [...blocks];
		markAsChanged();
	}

	// ── Keyboard shortcuts ─────────────────────────────────────────
	function handleGlobalKeyDown(event) {
		// Handle undo/redo
		if (handleUndoRedoKeyboard(event, undo, redo)) {
			return;
		}
	}


	// ── Speech ─────────────────────────────────────────────────────
	function getFullText() {
		if (editorMode === 'blocks') return blocks.map(b => b.content).join('\n\n');
		if (editorMode === 'markdown') return markdownContent;
		if (editorMode === 'richtext') return richTextEditorElement?.innerText || '';
		return '';
	}

	function splitTextIntoChunks(text, maxLength) {
		const chunks = [];
		let currentChunk = '';
		const sentences = text.split(/[.!?]+/).filter(s => s.trim());
		for (const sentence of sentences) {
			const trimmed = sentence.trim();
			if (!trimmed) continue;
			if (currentChunk.length + trimmed.length + 1 <= maxLength) {
				currentChunk += (currentChunk ? '. ' : '') + trimmed;
			} else {
				if (currentChunk) chunks.push(currentChunk + '.');
				currentChunk = trimmed;
			}
		}
		if (currentChunk) chunks.push(currentChunk + '.');
		return chunks;
	}

	async function readAloud() {
		if (isReadingAloud || !selectedVoiceModel) return;
		isReadingAloud = true;
		try {
			const fullText = getFullText();
			if (!fullText.trim()) { alert('No text to read aloud'); isReadingAloud = false; return; }
			const chunks = splitTextIntoChunks(fullText, 4000);
			for (const chunk of chunks) {
				if (!isReadingAloud) break;
				try {
					await SpeechService.textToSpeech(chunk, selectedVoiceModel.id, {
						voice: speechSettings?.defaultVoice || 'alloy',
						language: speechSettings?.defaultLanguage || 'en',
						speed: speechSettings?.ttsSpeed || 1.0,
						volume: speechSettings?.ttsVolume || 1.0
					});
					await new Promise(resolve => setTimeout(resolve, 500));
				} catch (error) {
					console.error('TTS chunk error:', error);
				}
			}
		} catch (error) {
			console.error('Read aloud error:', error);
			alert('Failed to read document aloud: ' + error.message);
		} finally {
			isReadingAloud = false;
		}
	}

	// ── AI helpers ─────────────────────────────────────────────────
	function getSelectedText() {
		if (editorMode === 'blocks') {
			return (selectedBlockIndex >= 0 && selectedBlockIndex < blocks.length)
				? blocks[selectedBlockIndex].content
				: '';
		}
		if (editorMode === 'markdown') return markdownContent;
		if (editorMode === 'richtext') {
			const sel = window.getSelection();
			if (sel && sel.rangeCount > 0) {
				const text = sel.toString();
				if (text) return text;
			}
			return richTextEditorElement?.innerText || '';
		}
		return '';
	}

	function replaceSelectedText(newText) {
		if (editorMode === 'blocks') {
			if (selectedBlockIndex >= 0 && selectedBlockIndex < blocks.length) {
				blocks[selectedBlockIndex].content = newText;
				blocks = [...blocks];
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
				blocks = [...blocks];
				markAsChanged();
			}
		} else if (editorMode === 'markdown') {
			markdownContent += ' ' + text;
			markAsChanged();
		} else if (editorMode === 'richtext') {
			document.execCommand('insertText', false, ' ' + text);
		}
	}

	async function getAIWritingSuggestions() {
		const text = getSelectedText();
		if (!text.trim()) { alert('Please select some text first'); return; }
		isGeneratingAI = true;
		showAIAssistant = true;
		try {
			const suggestions = await aiService.getWritingSuggestions(text, aiContext);
			aiSuggestions = suggestions?.suggestions || [];
		} catch (error) {
			console.error('Failed to get AI suggestions:', error);
			aiSuggestions = [];
			alert('Failed to get AI suggestions. Please try again.');
		} finally {
			isGeneratingAI = false;
		}
	}

	async function continueWriting() {
		const text = getSelectedText();
		if (!text.trim()) { alert('Please select some text first'); return; }
		isGeneratingAI = true;
		try {
			const continuation = await aiService.continueWriting(text, 'matching', 'medium');
			if (continuation) insertTextAfterSelection(continuation);
		} catch (error) {
			console.error('Failed to continue writing:', error);
			alert('Failed to continue writing. Please try again.');
		} finally {
			isGeneratingAI = false;
		}
	}

	async function generateDocumentSummary() {
		const fullText = getFullText();
		if (!fullText.trim()) { alert('No text to summarize'); return; }
		isGeneratingSummary = true;
		try {
			const summary = await aiService.summarizeDocument(fullText, 'concise');
			if (summary) documentSummary = summary;
		} catch (error) {
			console.error('Failed to generate summary:', error);
			alert('Failed to generate document summary. Please try again.');
		} finally {
			isGeneratingSummary = false;
		}
	}

	async function analyzeDocument() {
		const fullText = getFullText();
		if (!fullText.trim()) { alert('No text to analyze'); return; }
		try {
			const analysis = await aiService.analyzeText(fullText);
			if (analysis) { textAnalysis = analysis; showTextAnalysis = true; }
		} catch (error) {
			console.error('Failed to analyze text:', error);
			alert('Failed to analyze document. Please try again.');
		}
	}

	function applyAISuggestion(suggestion) {
		replaceSelectedText(suggestion);
		showAIAssistant = false;
		aiSuggestions = [];
	}

	// ── Find & Replace (all editor modes) ────────────────────────────
	function performFind() {
		if (!findText) return;
		findResults = [];
		currentFindIndex = -1;
		clearHighlights();

		if (editorMode === 'blocks') {
			findResults = searchInBlocks(blocks, findText, isCaseSensitive);
		} else if (editorMode === 'markdown') {
			findResults = searchInMarkdown(markdownContent, findText, isCaseSensitive);
		} else if (editorMode === 'richtext' && richTextEditorElement) {
			findResults = searchInRichText(richTextEditorElement, findText, isCaseSensitive);
		}

		if (findResults.length > 0) {
			currentFindIndex = 0;
			selectFindResult(0);
		}
	}

	function selectFindResult(index) {
		if (findResults.length === 0 || index < 0 || index >= findResults.length) return;
		currentFindIndex = index;
		const result = findResults[index];

		if (editorMode === 'blocks') {
			// Select the block containing the result
			selectedBlockIndex = result.blockIndex;
			// Scroll to the block
			setTimeout(() => {
				const blockElement = document.querySelector(`[data-block-index="${result.blockIndex}"]`);
				blockElement?.scrollIntoView({ behavior: 'smooth', block: 'center' });
			}, 0);
		} else if (editorMode === 'markdown') {
			// For markdown, we just know the position
			// The editor would need to handle cursor positioning
		} else if (editorMode === 'richtext') {
			highlightResult(result);
		}
	}

	function performReplace() {
		if (findResults.length === 0 || currentFindIndex === -1) return;
		const result = findResults[currentFindIndex];

		if (editorMode === 'blocks') {
			blocks = replaceInBlocks(blocks, result, replaceText);
			markAsChanged();
			// Re-search to update indices
			performFind();
		} else if (editorMode === 'markdown') {
			markdownContent = replaceInMarkdown(markdownContent, result, replaceText);
			markAsChanged();
			// Re-search to update indices
			performFind();
		} else if (editorMode === 'richtext') {
			document.execCommand('insertText', false, replaceText);
			markAsChanged();
			// Re-search in rich text
			setTimeout(() => performFind(), 0);
		}
	}

	function replaceAll() {
		if (!findText || !replaceText) return;

		if (editorMode === 'blocks') {
			blocks = replaceAllInBlocks(blocks, findText, replaceText, isCaseSensitive);
			markAsChanged();
		} else if (editorMode === 'markdown') {
			const flags = isCaseSensitive ? 'g' : 'gi';
			const regex = new RegExp(findText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), flags);
			markdownContent = markdownContent.replace(regex, replaceText);
			markAsChanged();
		} else if (editorMode === 'richtext' && richTextEditorElement) {
			const flags = isCaseSensitive ? 'g' : 'gi';
			const regex = new RegExp(findText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), flags);
			richTextEditorElement.innerHTML = richTextEditorElement.innerHTML.replace(regex, replaceText);
			markAsChanged();
		}

		findResults = [];
		currentFindIndex = -1;
		showFindReplaceDialog = false;
	}

	// ── Export Handler ───────────────────────────────────────────────
	function handleExport(event) {
		const format = event.detail?.format || 'pdf';
		const content = {
			title: documentTitle,
			mode: editorMode,
			blocks: blocks,
			markdown: markdownContent,
			html: richTextEditorElement?.innerHTML || ''
		};
		
		try {
			exportDocument(content, format);
		} catch (error) {
			console.error('Export failed:', error);
			alert('Export failed: ' + error.message);
		}
	}

	function findNext() {

		if (findResults.length === 0) {
			performFind();
			return;
		}
		currentFindIndex = (currentFindIndex + 1) % findResults.length;
		selectFindResult(currentFindIndex);
	}

	function findPrevious() {
		if (findResults.length === 0) {
			performFind();
			return;
		}
		currentFindIndex = currentFindIndex <= 0 ? findResults.length - 1 : currentFindIndex - 1;
		selectFindResult(currentFindIndex);
	}

</script>

<!-- ── Toolbar ──────────────────────────────────────────────────── -->
<TextEditorToolbar
	{documentTitle}
	{isSaving}
	{saveStatus}
	{hasUnsavedChanges}
	{currentDocument}
	{availableVoices}
	{isReadingAloud}
	{isGeneratingAI}
	{isGeneratingSummary}
	{canUndo}
	{canRedo}
	on:titleChange={(e) => { documentTitle = e.detail; markAsChanged(); }}
	on:save={saveDocument}
	on:export={handleExport}

	on:newDocument={createNewDocument}
	on:openDocumentList={loadDocumentList}
	on:openVersionHistory={loadVersionHistory}
	on:openMathHelp={() => showMathHelp = true}
	on:readAloud={readAloud}
	on:aiSuggest={getAIWritingSuggestions}
	on:aiContinue={continueWriting}
	on:aiSummarize={generateDocumentSummary}
	on:aiAnalyze={analyzeDocument}
	on:undo={undo}
	on:redo={redo}
/>

<!-- ── Global Keyboard Handler ───────────────────────────────────── -->
<svelte:window on:keydown={handleGlobalKeyDown} />


<!-- ── Editor canvas ────────────────────────────────────────────── -->
<div class="h-full overflow-y-auto p-8 bg-white">
	<div class="max-w-4xl mx-auto">
		{#if editorMode === 'blocks'}
			<TextEditorBlockEditor
				{blocks}
				{selectedBlockIndex}
				on:addBlock={handleAddBlock}
				on:deleteBlock={handleDeleteBlock}
				on:contentChange={handleBlockContentChange}
				on:selectBlock={(e) => selectedBlockIndex = e.detail}
			/>

		{:else if editorMode === 'markdown'}
			<TextEditorMarkdownEditor
				{markdownContent}
				on:change={(e) => { markdownContent = e.detail; markAsChanged(); }}
			/>

		{:else if editorMode === 'richtext'}
			<TextEditorRichText
				bind:editorElement={richTextEditorElement}
				on:change={markAsChanged}
	on:openFindReplace={() => {
					showFindReplaceDialog = true;
					findText = '';
					replaceText = '';
					findResults = [];
					currentFindIndex = -1;
					isCaseSensitive = false;
				}}

			/>
		{/if}
	</div>
</div>

<!-- ── All modals ────────────────────────────────────────────────── -->
<TextEditorModals
	{showDocumentList}
	{documentList}
	{isLoading}
	{showVersionHistory}
	{versionHistory}
	{showMathHelp}
	{showFindReplaceDialog}
	{showAIAssistant}
	{aiSuggestions}
	{documentSummary}
	{showTextAnalysis}
	{textAnalysis}
	bind:findText
	bind:replaceText
	{findResults}
	{currentFindIndex}
	{isCaseSensitive}
	on:findNext={findNext}
	on:findPrevious={findPrevious}
	on:toggleCaseSensitive={() => isCaseSensitive = !isCaseSensitive}
	on:loadDocument={(e) => loadDocument(e.detail)}

	on:closeDocumentList={() => showDocumentList = false}
	on:restoreVersion={(e) => restoreVersion(e.detail)}
	on:closeVersionHistory={() => showVersionHistory = false}
	on:closeMathHelp={() => showMathHelp = false}
	on:performFind={performFind}
	on:performReplace={performReplace}
	on:replaceAll={replaceAll}
	on:closeFindReplace={() => showFindReplaceDialog = false}
	on:applyAISuggestion={(e) => applyAISuggestion(e.detail)}
	on:closeAIAssistant={() => { showAIAssistant = false; aiSuggestions = []; }}
	on:closeSummary={() => documentSummary = ''}
	on:closeTextAnalysis={() => { showTextAnalysis = false; textAnalysis = null; }}
/>
