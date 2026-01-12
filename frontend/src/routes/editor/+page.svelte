<script>
	import TextEditor from '$lib/components/TextEditor.svelte';
	import { onMount } from 'svelte';

	// Accept framework-provided props to avoid warnings
	export let params = null;
	export let data = null;
	export let form = null;

	let documentTitle = 'Untitled Document';
	let isAutoSaving = false;
	let lastSaved = null;

	// Document settings
	let editorMode = 'blocks'; // 'blocks', 'markdown', 'richtext'
	let showMathPanel = false;

	// AI features
	let showAIPanel = false;
	let aiPrompt = '';

	onMount(() => {
		// Initialize editor
		console.log('Text Editor loaded');

		// Simulate auto-save
		setInterval(() => {
			if (!isAutoSaving) {
				isAutoSaving = true;
				setTimeout(() => {
					isAutoSaving = false;
					lastSaved = new Date();
				}, 1000);
			}
		}, 5000);
	});

	function handleSave() {
		// TODO: Implement save functionality
		console.log('Saving document...');
		lastSaved = new Date();
	}

	function handleExport() {
		// TODO: Implement export functionality
		console.log('Exporting document...');
	}

	function insertMathExpression(expression) {
		// TODO: Insert math expression into editor
		console.log('Inserting math:', expression);
		showMathPanel = false;
	}

	function runAIPrompt() {
		// TODO: Implement AI assistance
		console.log('Running AI prompt:', aiPrompt);
		showAIPanel = false;
		aiPrompt = '';
	}
</script>

<svelte:head>
	<title>{documentTitle} - TPT Text Editor</title>
</svelte:head>

<div class="h-screen flex flex-col bg-white">
	<!-- Header -->
	<header class="flex items-center justify-between px-6 py-3 border-b border-gray-200 bg-white">
		<div class="flex items-center space-x-4">
			<h1 class="text-xl font-semibold text-gray-900">{documentTitle}</h1>
			<div class="flex items-center space-x-2">
				{#if isAutoSaving}
					<span class="text-sm text-orange-600 flex items-center">
						<svg class="w-4 h-4 mr-1 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
						</svg>
						Saving...
					</span>
				{:else if lastSaved}
					<span class="text-sm text-green-600">
						Saved {lastSaved.toLocaleTimeString()}
					</span>
				{/if}
			</div>
		</div>

		<div class="flex items-center space-x-2">
			<!-- Editor Mode Toggle -->
			<div class="flex items-center bg-gray-100 rounded-lg p-1">
				<button
					class="px-3 py-1 text-sm rounded-md transition-colors {editorMode === 'blocks' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
					on:click={() => editorMode = 'blocks'}
				>
					Blocks
				</button>
				<button
					class="px-3 py-1 text-sm rounded-md transition-colors {editorMode === 'markdown' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
					on:click={() => editorMode = 'markdown'}
				>
					Markdown
				</button>
				<button
					class="px-3 py-1 text-sm rounded-md transition-colors {editorMode === 'richtext' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
					on:click={() => editorMode = 'richtext'}
				>
					Rich Text
				</button>
			</div>

			<!-- Action Buttons -->
			<button
				class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
				on:click={handleSave}
			>
				Save
			</button>
			<button
				class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
				on:click={handleExport}
			>
				Export
			</button>
		</div>
	</header>

	<!-- Toolbar -->
	<div class="flex items-center space-x-2 px-6 py-3 border-b border-gray-200 bg-gray-50">
		<!-- Formatting Tools -->
		<div class="flex items-center space-x-1">
			<button class="p-2 hover:bg-white rounded border border-transparent hover:border-gray-300 transition-colors" title="Bold">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 4h8a4 4 0 014 4v8a4 4 0 01-4 4H6a4 4 0 01-4-4V8a4 4 0 014-4z"></path>
				</svg>
			</button>
			<button class="p-2 hover:bg-white rounded border border-transparent hover:border-gray-300 transition-colors" title="Italic">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
				</svg>
			</button>
			<button class="p-2 hover:bg-white rounded border border-transparent hover:border-gray-300 transition-colors" title="Underline">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"></path>
				</svg>
			</button>
		</div>

		<div class="w-px h-6 bg-gray-300"></div>

		<!-- Special Features -->
		<button
			class="p-2 hover:bg-white rounded border border-transparent hover:border-gray-300 transition-colors text-purple-600"
			title="Insert Math"
			on:click={() => showMathPanel = !showMathPanel}
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
			</svg>
		</button>

		<button
			class="p-2 hover:bg-white rounded border border-transparent hover:border-gray-300 transition-colors text-green-600"
			title="AI Assistant"
			on:click={() => showAIPanel = !showAIPanel}
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"></path>
			</svg>
		</button>

		<button class="p-2 hover:bg-white rounded border border-transparent hover:border-gray-300 transition-colors" title="Insert Table">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
			</svg>
		</button>

		<button class="p-2 hover:bg-white rounded border border-transparent hover:border-gray-300 transition-colors" title="Insert Image">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
			</svg>
		</button>
	</div>

	<!-- Main Editor Area -->
	<div class="flex-1 overflow-hidden relative">
		<TextEditor {editorMode} />

		<!-- Math Panel -->
		{#if showMathPanel}
			<div class="absolute top-4 right-4 bg-white border border-gray-200 rounded-lg shadow-lg p-4 w-80 z-50">
				<h3 class="text-lg font-semibold text-gray-900 mb-3">Natural Math</h3>
				<div class="space-y-2">
					<p class="text-sm text-gray-600 mb-3">Type natural language or use templates:</p>
					<button
						class="w-full text-left p-2 hover:bg-gray-50 rounded text-sm"
						on:click={() => insertMathExpression('integral')}
					>
						∫ integral of x² dx
					</button>
					<button
						class="w-full text-left p-2 hover:bg-gray-50 rounded text-sm"
						on:click={() => insertMathExpression('fraction')}
					>
						½ fraction a over b
					</button>
					<button
						class="w-full text-left p-2 hover:bg-gray-50 rounded text-sm"
						on:click={() => insertMathExpression('summation')}
					>
						∑ sum from i=1 to n
					</button>
					<button
						class="w-full text-left p-2 hover:bg-gray-50 rounded text-sm"
						on:click={() => insertMathExpression('sqrt')}
					>
						√ square root of x
					</button>
				</div>
				<button
					class="w-full mt-3 px-3 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 transition-colors"
					on:click={() => showMathPanel = false}
				>
					Close
				</button>
			</div>
		{/if}

		<!-- AI Panel -->
		{#if showAIPanel}
			<div class="absolute top-4 right-4 bg-white border border-gray-200 rounded-lg shadow-lg p-4 w-80 z-50">
				<h3 class="text-lg font-semibold text-gray-900 mb-3">AI Writing Assistant</h3>
				<textarea
					bind:value={aiPrompt}
					placeholder="Describe what you want to write..."
					rows="4"
					class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 resize-none"
				></textarea>
				<div class="flex space-x-2 mt-3">
					<button
						class="flex-1 px-3 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
						on:click={runAIPrompt}
					>
						Generate
					</button>
					<button
						class="px-3 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400 transition-colors"
						on:click={() => showAIPanel = false}
					>
						Cancel
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
