<!-- frontend/src/lib/components/SpreadsheetFindReplaceModal.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';
	
	export let show = false;
	export let data = [];
	
	const dispatch = createEventDispatcher();
	
	let findText = '';
	let replaceText = '';
	let caseSensitive = false;
	let wholeCell = false;
	let searchInFormulas = false;
	let matches = [];
	let currentMatchIndex = -1;
	let replacementCount = 0;
	let showResults = false;
	
	function handleFind() {
		if (!findText) {
			matches = [];
			currentMatchIndex = -1;
			return;
		}
		
		// Use the search utility
		import('../utils/spreadsheet-search.js').then(module => {
			const searchFn = searchInFormulas ? module.searchInFormulas : module.searchInSpreadsheet;
			matches = searchFn(data, findText, { caseSensitive, wholeCell });
			currentMatchIndex = matches.length > 0 ? 0 : -1;
			showResults = true;
			
			if (matches.length > 0) {
				dispatch('selectCell', matches[0]);
			}
		});
	}
	
	function handleFindNext() {
		if (matches.length === 0) {
			handleFind();
			return;
		}
		
		currentMatchIndex = (currentMatchIndex + 1) % matches.length;
		dispatch('selectCell', matches[currentMatchIndex]);
	}
	
	function handleFindPrevious() {
		if (matches.length === 0) {
			handleFind();
			return;
		}
		
		currentMatchIndex = currentMatchIndex <= 0 ? matches.length - 1 : currentMatchIndex - 1;
		dispatch('selectCell', matches[currentMatchIndex]);
	}
	
	function handleReplace() {
		if (currentMatchIndex < 0 || currentMatchIndex >= matches.length) {
			handleFind();
			return;
		}
		
		const match = matches[currentMatchIndex];
		dispatch('replace', {
			row: match.row,
			col: match.col,
			newValue: match.value.replace(
				new RegExp(escapeRegExp(findText), caseSensitive ? 'g' : 'gi'),
				replaceText
			)
		});
		
		// Remove this match and continue
		matches.splice(currentMatchIndex, 1);
		if (matches.length > 0) {
			currentMatchIndex = currentMatchIndex % matches.length;
			dispatch('selectCell', matches[currentMatchIndex]);
		} else {
			currentMatchIndex = -1;
		}
	}
	
	function handleReplaceAll() {
		import('../utils/spreadsheet-search.js').then(module => {
			const result = module.replaceAllInSpreadsheet(data, findText, replaceText, {
				caseSensitive,
				wholeCell
			});
			
			replacementCount = result.replacements;
			dispatch('replaceAll', result);
			showResults = true;
			matches = [];
			currentMatchIndex = -1;
		});
	}
	
	function escapeRegExp(string) {
		return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
	}
	
	function handleClose() {
		show = false;
		findText = '';
		replaceText = '';
		matches = [];
		currentMatchIndex = -1;
		showResults = false;
		dispatch('close');
	}
	
	function handleKeyDown(event) {
		if (event.key === 'Enter' && event.shiftKey) {
			event.preventDefault();
			handleFindPrevious();
		} else if (event.key === 'Enter') {
			event.preventDefault();
			handleFindNext();
		} else if (event.key === 'Escape') {
			handleClose();
		}
	}
</script>

{#if show}
	<div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
		<div class="bg-white rounded-lg shadow-xl max-w-lg w-full">
			<!-- Header -->
			<div class="flex items-center justify-between px-6 py-4 border-b border-gray-200">
				<h2 class="text-xl font-semibold text-gray-900">Find and Replace</h2>
				<button
					class="text-gray-400 hover:text-gray-600 transition-colors"
					on:click={handleClose}
				>
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
			
			<!-- Content -->
			<div class="p-6 space-y-4">
				<!-- Find -->
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Find</label>
					<input
						type="text"
						bind:value={findText}
						on:keydown={handleKeyDown}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						placeholder="Text to find..."
					/>
				</div>
				
				<!-- Replace -->
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Replace with</label>
					<input
						type="text"
						bind:value={replaceText}
						on:keydown={handleKeyDown}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						placeholder="Replacement text..."
					/>
				</div>
				
				<!-- Options -->
				<div class="flex flex-wrap gap-4">
					<label class="flex items-center space-x-2 cursor-pointer">
						<input
							type="checkbox"
							bind:checked={caseSensitive}
							class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
						/>
						<span class="text-sm text-gray-700">Case sensitive</span>
					</label>
					
					<label class="flex items-center space-x-2 cursor-pointer">
						<input
							type="checkbox"
							bind:checked={wholeCell}
							class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
						/>
						<span class="text-sm text-gray-700">Match whole cell</span>
					</label>
					
					<label class="flex items-center space-x-2 cursor-pointer">
						<input
							type="checkbox"
							bind:checked={searchInFormulas}
							class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
						/>
						<span class="text-sm text-gray-700">Search in formulas</span>
					</label>
				</div>
				
				<!-- Results -->
				{#if showResults}
					<div class="bg-gray-50 rounded-md p-3 text-sm">
						{#if matches.length > 0}
							<p class="text-gray-700">
								Match {currentMatchIndex + 1} of {matches.length}
								{#if matches[currentMatchIndex]}
									<span class="text-gray-500">
										({matches[currentMatchIndex].cellId})
									</span>
								{/if}
							</p>
						{:else if replacementCount > 0}
							<p class="text-green-600">
								✓ Replaced {replacementCount} occurrence{replacementCount !== 1 ? 's' : ''}
							</p>
						{:else}
							<p class="text-gray-500">No matches found</p>
						{/if}
					</div>
				{/if}
				
				<!-- Buttons -->
				<div class="flex flex-wrap gap-2 pt-2">
					<button
						class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
						on:click={handleFind}
					>
						Find
					</button>
					<button
						class="px-4 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
						on:click={handleFindNext}
						disabled={matches.length === 0}
					>
						Find Next
					</button>
					<button
						class="px-4 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
						on:click={handleFindPrevious}
						disabled={matches.length === 0}
					>
						Find Previous
					</button>
					<button
						class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
						on:click={handleReplace}
						disabled={matches.length === 0}
					>
						Replace
					</button>
					<button
						class="px-4 py-2 bg-green-700 text-white rounded hover:bg-green-800 transition-colors"
						on:click={handleReplaceAll}
					>
						Replace All
					</button>
					<div class="flex-1"></div>
					<button
						class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 transition-colors"
						on:click={handleClose}
					>
						Close
					</button>
				</div>
				
				<!-- Keyboard shortcuts -->
				<div class="text-xs text-gray-500 pt-2 border-t border-gray-200">
					<p>Shortcuts: Enter = Find Next, Shift+Enter = Find Previous, Escape = Close</p>
				</div>
			</div>
		</div>
	</div>
{/if}
