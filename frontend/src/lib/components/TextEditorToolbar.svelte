<!-- frontend/src/lib/components/TextEditorToolbar.svelte -->
<script>
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

	const dispatch = createEventDispatcher();

	let showExportMenu = false;

	function handleExport(format) {
		showExportMenu = false;
		dispatch('export', { format });
	}

</script>


<div class="bg-white border-b border-gray-200 px-8 py-4">
	<div class="max-w-4xl mx-auto flex items-center justify-between">
		<div class="flex items-center space-x-4">
			<input
				bind:value={documentTitle}
				placeholder="Document title..."
				class="text-xl font-semibold bg-transparent border-none outline-none focus:ring-2 focus:ring-blue-500 rounded px-2 py-1"
				on:input={() => dispatch('titleChange', documentTitle)}
			/>
			{#if hasUnsavedChanges}
				<span class="text-sm text-orange-600">•</span>
			{/if}
			<span class="text-sm text-gray-500">{saveStatus}</span>
		</div>

		<div class="flex items-center space-x-2">
			<!-- Undo/Redo -->
			<button
				class="px-2 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed"
				on:click={() => dispatch('undo')}
				disabled={!canUndo}
				title="Undo (Ctrl+Z)"
			>
				↩️ Undo
			</button>
			<button
				class="px-2 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed"
				on:click={() => dispatch('redo')}
				disabled={!canRedo}
				title="Redo (Ctrl+Y or Ctrl+Shift+Z)"
			>
				↪️ Redo
			</button>

			<div class="w-px h-6 bg-gray-300 mx-2"></div>

			<button
				class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
				on:click={() => dispatch('save')}
				disabled={isSaving}
			>
				{isSaving ? 'Saving...' : 'Save'}
			</button>

			<div class="relative">
				<button
					class="px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700 flex items-center"
					on:click={() => showExportMenu = !showExportMenu}
				>
					Export ▼
				</button>
				{#if showExportMenu}
					<div class="absolute top-full right-0 mt-1 w-48 bg-white border border-gray-200 rounded shadow-lg z-50">
						<button
							class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 first:rounded-t"
							on:click={() => handleExport('pdf')}
						>
							📄 Export as PDF
						</button>
						<button
							class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
							on:click={() => handleExport('docx')}
						>
							📝 Export as DOCX
						</button>
						<button
							class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
							on:click={() => handleExport('html')}
						>
							🌐 Export as HTML
						</button>
						<button
							class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
							on:click={() => handleExport('md')}
						>
							📋 Export as Markdown
						</button>
						<button
							class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 last:rounded-b"
							on:click={() => handleExport('txt')}
						>
							📃 Export as Text
						</button>
					</div>
				{/if}
			</div>


			<button
				class="px-3 py-1 text-sm bg-gray-600 text-white rounded hover:bg-gray-700"
				on:click={() => dispatch('newDocument')}
			>
				New
			</button>

			<button
				class="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700"
				on:click={() => dispatch('openDocumentList')}
			>
				Open
			</button>


			{#if currentDocument}
				<button
					class="px-3 py-1 text-sm bg-purple-600 text-white rounded hover:bg-purple-700"
					on:click={() => dispatch('openVersionHistory')}
				>
					History
				</button>
			{/if}

			<button
				class="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700"
				on:click={() => dispatch('openMathHelp')}
				title="Math Help & Examples"
			>
				Math Help
			</button>

			{#if availableVoices.length > 0}
				<button
					class="px-3 py-1 text-sm bg-purple-600 text-white rounded hover:bg-purple-700 disabled:opacity-50"
					on:click={() => dispatch('readAloud')}
					disabled={isReadingAloud}
					title="Read document aloud"
				>
					{isReadingAloud ? '🔊 Reading...' : '📖 Read Aloud'}
				</button>
			{/if}

			<!-- AI Writing Assistant -->
			<div class="border-l border-gray-300 pl-4 ml-4 flex items-center space-x-2">
				<span class="text-xs text-gray-500 font-medium">AI</span>

				<div class="relative">
					<button
						class="px-3 py-1 text-sm bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-50 relative group"
						on:click={() => dispatch('aiSuggest')}
						disabled={isGeneratingAI}
						title="Get AI writing suggestions for grammar, style, and clarity improvements"
						aria-label="Get AI writing suggestions"
						type="button"
					>
						{isGeneratingAI ? '💭' : '✨ Suggest'}
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none">
							✨ AI Writing Assistant
							<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
						</div>
					</button>

					<button
						class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 relative group ml-1"
						on:click={() => dispatch('aiContinue')}
						disabled={isGeneratingAI}
						title="AI content continuation - select text to continue from that point"
						aria-label="Continue writing with AI"
						type="button"
					>
						📝 Continue
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none">
							📝 Continue Writing
							<div class="text-xs text-gray-300 mt-1">Select text to continue from that point</div>
							<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
						</div>
					</button>

					<button
						class="px-3 py-1 text-sm bg-yellow-600 text-white rounded hover:bg-yellow-700 disabled:opacity-50 relative group ml-1"
						on:click={() => dispatch('aiSummarize')}
						disabled={isGeneratingSummary}
						title="Generate intelligent document summaries - works best with longer documents"
						aria-label="Generate document summary"
						type="button"
					>
						{isGeneratingSummary ? '📋' : '📄 Summary'}
						<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10 pointer-events-none">
							📄 Document Summary
							<div class="text-xs text-gray-300 mt-1">Best with 500+ words</div>
							<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
						</div>
					</button>

					<button
						class="px-3 py-1 text-sm bg-gray-600 text-white rounded hover:bg-gray-700 relative group ml-1"
						on:click={() => dispatch('aiAnalyze')}
						title="Analyze document for readability, sentiment, and key insights"
						aria-label="Analyze document"
						type="button"
					>
						📊 Analyze
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
