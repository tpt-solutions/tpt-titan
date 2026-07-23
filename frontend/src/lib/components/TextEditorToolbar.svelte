<!-- frontend/src/lib/components/TextEditorToolbar.svelte -->
<script>
// @ts-nocheck
	import { createEventDispatcher } from 'svelte';

	export let documentTitle = 'Untitled Document';
	export let isSaving = false;
	export let saveStatus = '';
	export let hasUnsavedChanges = false;
	export let currentDocument = null;
	export let availableVoices = [];
	export let isReadingAloud = false;
	export let isGeneratingAI = false;
	export let isGeneratingSummary = false;
	export let canUndo = false;
	export let canRedo = false;
	export let editorMode = 'blocks';

	const dispatch = createEventDispatcher();

	let showExportMenu = false;
	let showFindReplace = false;

	function handleExport(format) {
		showExportMenu = false;
		dispatch('export', { format });
	}

	function setMode(mode) {
		editorMode = mode;
		dispatch('modeChange', mode);
	}

</script>



<div class="bg-white border-b border-gray-200">
	<!-- Top Row: Title and Document Info -->
	<div class="px-6 py-3 border-b border-gray-100">
		<div class="max-w-6xl mx-auto flex items-center justify-between">
			<div class="flex items-center space-x-3 flex-1 min-w-0">
				<input
					bind:value={documentTitle}
					placeholder="Document title..."
					class="text-xl font-semibold bg-transparent border-none outline-none focus:ring-2 focus:ring-blue-500 rounded px-2 py-1 w-full max-w-md"
					on:input={() => dispatch('titleChange', documentTitle)}
				/>
				<div class="flex items-center space-x-2 shrink-0">
					{#if hasUnsavedChanges}
						<span class="text-xs text-orange-600 font-medium">● Unsaved</span>
					{/if}
					{#if saveStatus}
						<span class="text-xs text-gray-400">{saveStatus}</span>
					{/if}
				</div>
			</div>
		</div>
	</div>

	<!-- Bottom Row: Organized Toolbar Sections -->
	<div class="px-6 py-2 bg-gray-50">
		<div class="max-w-6xl mx-auto flex items-center justify-between flex-wrap gap-2">
			<!-- Left: File Operations -->
			<div class="flex items-center space-x-1">
				<button
					class="px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 hover:border-gray-400 transition-colors flex items-center space-x-1"
					on:click={() => dispatch('newDocument')}
					title="Create new document"
				>
					<span>➕</span>
					<span>New</span>
				</button>
				<button
					class="px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 hover:border-gray-400 transition-colors flex items-center space-x-1"
					on:click={() => dispatch('openDocumentList')}
					title="Open existing document"
				>
					<span>📂</span>
					<span>Open</span>
				</button>
				<button
					class="px-3 py-1.5 text-sm font-medium text-white bg-blue-600 border border-blue-600 rounded-md hover:bg-blue-700 transition-colors flex items-center space-x-1 disabled:opacity-50 disabled:cursor-not-allowed"
					on:click={() => dispatch('save')}
					disabled={isSaving}
					title="Save document (Ctrl+S)"
				>
					<span>💾</span>
					<span>{isSaving ? 'Saving...' : 'Save'}</span>
				</button>
				
				<!-- Export Dropdown -->
				<div class="relative">
					<button
						class="px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 hover:border-gray-400 transition-colors flex items-center space-x-1"
						on:click={() => showExportMenu = !showExportMenu}
						title="Export document"
					>
						<span>📤</span>
						<span>Export</span>
						<span class="text-xs ml-1">▼</span>
					</button>
					{#if showExportMenu}
						<div class="absolute top-full left-0 mt-1 w-48 bg-white border border-gray-200 rounded-lg shadow-xl z-50 py-1">
							<div class="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-100">
								Export as
							</div>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-700 transition-colors flex items-center space-x-2"
								on:click={() => handleExport('pdf')}
							>
								<span>📄</span>
								<span>PDF Document</span>
							</button>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-700 transition-colors flex items-center space-x-2"
								on:click={() => handleExport('docx')}
							>
								<span>📝</span>
								<span>Word Document</span>
							</button>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-700 transition-colors flex items-center space-x-2"
								on:click={() => handleExport('html')}
							>
								<span>🌐</span>
								<span>HTML Page</span>
							</button>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-700 transition-colors flex items-center space-x-2"
								on:click={() => handleExport('md')}
							>
								<span>📋</span>
								<span>Markdown</span>
							</button>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-blue-50 hover:text-blue-700 transition-colors flex items-center space-x-2"
								on:click={() => handleExport('txt')}
							>
								<span>📃</span>
								<span>Plain Text</span>
							</button>
						</div>
					{/if}
				</div>
			</div>

			<!-- Center: Editor Mode Switcher -->
			<div class="flex items-center bg-white rounded-lg border border-gray-300 p-1 shadow-sm">
				<button
					class="px-3 py-1 text-sm font-medium rounded-md transition-all duration-200 {editorMode === 'blocks' ? 'bg-blue-600 text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'}"
					on:click={() => setMode('blocks')}
					title="Block-based editor (Notion-style)"
				>
					🧱 Blocks
				</button>
				<button
					class="px-3 py-1 text-sm font-medium rounded-md transition-all duration-200 {editorMode === 'markdown' ? 'bg-blue-600 text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'}"
					on:click={() => setMode('markdown')}
					title="Markdown editor"
				>
					📝 Markdown
				</button>
				<button
					class="px-3 py-1 text-sm font-medium rounded-md transition-all duration-200 {editorMode === 'richtext' ? 'bg-blue-600 text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'}"
					on:click={() => setMode('richtext')}
					title="Rich text editor (Word-style)"
				>
					✨ Rich Text
				</button>
			</div>

			<!-- Right: Tools & AI -->
			<div class="flex items-center space-x-1">
				<!-- Edit Actions -->
				<div class="flex items-center space-x-1 pr-2 border-r border-gray-300">
					<button
						class="p-1.5 text-gray-600 hover:bg-gray-200 rounded-md transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
						on:click={() => dispatch('undo')}
						disabled={!canUndo}
						title="Undo (Ctrl+Z)"
					>
						↩️
					</button>
					<button
						class="p-1.5 text-gray-600 hover:bg-gray-200 rounded-md transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
						on:click={() => dispatch('redo')}
						disabled={!canRedo}
						title="Redo (Ctrl+Y)"
					>
						↪️
					</button>
					<button
						class="p-1.5 text-gray-600 hover:bg-gray-200 rounded-md transition-colors"
						on:click={() => dispatch('openFindReplace')}
						title="Find & Replace (Ctrl+F)"
					>
						🔍
					</button>
				</div>

				<!-- AI Assistant Dropdown -->
				<div class="relative">
					<button
						class="px-3 py-1.5 text-sm font-medium text-indigo-700 bg-indigo-50 border border-indigo-200 rounded-md hover:bg-indigo-100 transition-colors flex items-center space-x-1"
						on:click={() => showFindReplace = !showFindReplace}
						title="AI Writing Assistant"
					>
						<span>🤖</span>
						<span>AI</span>
						<span class="text-xs ml-1">▼</span>
					</button>
					{#if showFindReplace}
						<div class="absolute top-full right-0 mt-1 w-56 bg-white border border-gray-200 rounded-lg shadow-xl z-50 py-1">
							<div class="px-3 py-2 text-xs font-semibold text-indigo-600 uppercase tracking-wider border-b border-gray-100">
								AI Assistant
							</div>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-700 transition-colors flex items-center space-x-2 disabled:opacity-50"
								on:click={() => { dispatch('aiSuggest'); showFindReplace = false; }}
								disabled={isGeneratingAI}
							>
								<span>✨</span>
								<span>Suggest Improvements</span>
							</button>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-700 transition-colors flex items-center space-x-2 disabled:opacity-50"
								on:click={() => { dispatch('aiContinue'); showFindReplace = false; }}
								disabled={isGeneratingAI}
							>
								<span>📝</span>
								<span>Continue Writing</span>
							</button>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-700 transition-colors flex items-center space-x-2 disabled:opacity-50"
								on:click={() => { dispatch('aiSummarize'); showFindReplace = false; }}
								disabled={isGeneratingSummary}
							>
								<span>📄</span>
								<span>Generate Summary</span>
							</button>
							<button
								class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-700 transition-colors flex items-center space-x-2"
								on:click={() => { dispatch('aiAnalyze'); showFindReplace = false; }}
							>
								<span>📊</span>
								<span>Analyze Document</span>
							</button>
						</div>
					{/if}
				</div>

				<!-- Tools -->
				<div class="flex items-center space-x-1 pl-2 border-l border-gray-300">
					<button
						class="p-1.5 text-gray-600 hover:bg-gray-200 rounded-md transition-colors"
						on:click={() => dispatch('openMathHelp')}
						title="Math Help & Examples"
					>
						🔢
					</button>
					{#if availableVoices.length > 0}
						<button
							class="p-1.5 text-gray-600 hover:bg-gray-200 rounded-md transition-colors disabled:opacity-50"
							on:click={() => dispatch('readAloud')}
							disabled={isReadingAloud}
							title="Read document aloud"
						>
							{isReadingAloud ? '🔊' : '📖'}
						</button>
					{/if}
					{#if currentDocument}
						<button
							class="p-1.5 text-gray-600 hover:bg-gray-200 rounded-md transition-colors"
							on:click={() => dispatch('openVersionHistory')}
							title="Version History"
						>
							🕐
						</button>
					{/if}
				</div>
			</div>
		</div>
	</div>
</div>
