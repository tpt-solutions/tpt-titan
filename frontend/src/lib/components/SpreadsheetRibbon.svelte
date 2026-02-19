<script>
	import {
		activeRibbonTab,
		ribbonTabs,
		showRibbonCustomizer,
		canUndo,
		canRedo
	} from '../stores/spreadsheet-store.js';

	// Default available tools

	const defaultAvailableTools = [
		// File operations
		{ id: 'save', name: 'Save', icon: '💾', category: 'file', shortcut: 'Ctrl+S' },
		{ id: 'open', name: 'Open', icon: '📂', category: 'file' },
		{ id: 'export', name: 'Export', icon: '📤', category: 'file' },

		// Editing
		{ id: 'undo', name: 'Undo', icon: '↶', category: 'edit', shortcut: 'Ctrl+Z' },
		{ id: 'redo', name: 'Redo', icon: '↷', category: 'edit', shortcut: 'Ctrl+Y' },
		{ id: 'copy', name: 'Copy', icon: '📋', category: 'edit', shortcut: 'Ctrl+C' },
		{ id: 'paste', name: 'Paste', icon: '📌', category: 'edit', shortcut: 'Ctrl+V' },
		{ id: 'cut', name: 'Cut', icon: '✂️', category: 'edit', shortcut: 'Ctrl+X' },

		// Formatting
		{ id: 'bold', name: 'Bold', icon: 'B', category: 'format' },
		{ id: 'italic', name: 'Italic', icon: 'I', category: 'format' },
		{ id: 'underline', name: 'Underline', icon: 'U', category: 'format' },
		{ id: 'align-left', name: 'Align Left', icon: '⬅️', category: 'format' },
		{ id: 'align-center', name: 'Align Center', icon: '⬌', category: 'format' },
		{ id: 'align-right', name: 'Align Right', icon: '➡️', category: 'format' },

		// Formulas
		{ id: 'sum', name: 'Sum', icon: '∑', category: 'formulas' },
		{ id: 'average', name: 'Average', icon: '📊', category: 'formulas' },
		{ id: 'count', name: 'Count', icon: '#', category: 'formulas' },
		{ id: 'min', name: 'Min', icon: '↓', category: 'formulas' },
		{ id: 'max', name: 'Max', icon: '↑', category: 'formulas' },

		// Insert
		{ id: 'insert-row', name: 'Insert Row', icon: '➕', category: 'insert' },
		{ id: 'insert-column', name: 'Insert Column', icon: '➕', category: 'insert' },
		{ id: 'insert-chart', name: 'Insert Chart', icon: '📈', category: 'insert' },
		{ id: 'insert-pivot', name: 'Insert Pivot', icon: '📋', category: 'insert' },

		// Data
		{ id: 'sort-asc', name: 'Sort Ascending', icon: '↑', category: 'data' },
		{ id: 'sort-desc', name: 'Sort Descending', icon: '↓', category: 'data' },
		{ id: 'filter', name: 'Filter', icon: '🔍', category: 'data' },
		{ id: 'data-cleanup', name: 'Data Cleanup', icon: '🧹', category: 'data' }
	];

	let availableTools = defaultAvailableTools;

	// Dispatch events to parent

	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();


	function dispatchAction(action, data = {}) {
		dispatch('action', { action, ...data });
	}

	// Handle tool click
	function handleToolClick(tool) {
		switch (tool.id) {
			case 'save':
				dispatchAction('save');
				break;
			case 'copy':
				dispatchAction('copy');
				break;
			case 'paste':
				dispatchAction('paste');
				break;
			case 'cut':
				dispatchAction('cut');
				break;
			case 'undo':
				dispatchAction('undo');
				break;
			case 'redo':
				dispatchAction('redo');
				break;
			case 'bold':
				dispatchAction('applyFormatting', { type: 'bold' });
				break;
			case 'italic':
				dispatchAction('applyFormatting', { type: 'italic' });
				break;
			case 'underline':
				dispatchAction('applyFormatting', { type: 'underline' });
				break;
			case 'align-left':
				dispatchAction('applyFormatting', { type: 'align', value: 'left' });
				break;
			case 'align-center':
				dispatchAction('applyFormatting', { type: 'align', value: 'center' });
				break;
			case 'align-right':
				dispatchAction('applyFormatting', { type: 'align', value: 'right' });
				break;
			case 'sum':
				dispatchAction('insertFormula', { formula: '=SUM()' });
				break;
			case 'average':
				dispatchAction('insertFormula', { formula: '=AVERAGE()' });
				break;
			case 'count':
				dispatchAction('insertFormula', { formula: '=COUNT()' });
				break;
			case 'min':
				dispatchAction('insertFormula', { formula: '=MIN()' });
				break;
			case 'max':
				dispatchAction('insertFormula', { formula: '=MAX()' });
				break;
			case 'insert-row':
				dispatchAction('insertRowBelow');
				break;
			case 'insert-column':
				dispatchAction('insertColumnRight');
				break;
			case 'insert-pivot':
				dispatchAction('insertPivotTable');
				break;
			case 'sort-asc':
				dispatchAction('sortColumn', { direction: 'asc' });
				break;
			case 'sort-desc':
				dispatchAction('sortColumn', { direction: 'desc' });
				break;
			case 'filter':
				dispatchAction('showFilterDialog');
				break;
			case 'data-cleanup':
				dispatchAction('showDataPrepTools');
				break;
			default:
				dispatchAction('toolAction', { toolId: tool.id, tool });
		}
	}

	// Get tool by ID
	function getToolById(toolId) {
		return availableTools.find(t => t.id === toolId);
	}
</script>

<!-- Customizable Ribbon -->
<div class="bg-white border-b border-gray-200">
	<!-- Ribbon Tabs -->
	<div class="flex border-b border-gray-200">
		{#each $ribbonTabs as tab, tabIndex}
			<button
				class="px-4 py-2 text-sm font-medium border-r border-gray-200 hover:bg-gray-50 transition-colors flex items-center space-x-2 {activeRibbonTab === tab.id ? 'bg-blue-50 text-blue-700 border-b-2 border-blue-500' : 'text-gray-700'}"
				on:click={() => activeRibbonTab.set(tab.id)}
			>
				<span>{tab.icon}</span>
				<span>{tab.name}</span>
			</button>
		{/each}
		<div class="flex-1"></div>
		<button
			class="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded transition-colors"
			on:click={() => showRibbonCustomizer.set(true)}
			title="Customize Ribbon"
		>
			⚙️
		</button>
	</div>

	<!-- Ribbon Content -->
	<div class="p-4 min-h-20">
		{#each $ribbonTabs as tab}
			{#if $activeRibbonTab === tab.id}
				<div class="flex flex-wrap gap-2">
					{#each tab.tools as tool, toolIndex}
						<div class="group relative">
						<button
							class="flex flex-col items-center justify-center w-12 h-12 p-2 text-sm border border-gray-300 rounded hover:bg-gray-50 hover:border-gray-400 transition-colors disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:bg-white"
							title="{tool.name}{tool.shortcut ? ` (${tool.shortcut})` : ''}"
							on:click={() => handleToolClick(tool)}
							disabled={(tool.id === 'undo' && !$canUndo) || (tool.id === 'redo' && !$canRedo)}

						>

								<span class="text-lg">{tool.icon}</span>
								<span class="text-xs mt-1 truncate w-full">{tool.name.split(' ')[0]}</span>
							</button>
							<!-- Remove button (shown on hover) -->
							<button
								class="absolute -top-1 -right-1 w-4 h-4 bg-red-500 text-white text-xs rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-600"
								on:click={() => dispatchAction('removeToolFromRibbon', { tabId: tab.id, toolIndex })}
								title="Remove from ribbon"
							>
								×
							</button>
						</div>
					{/each}
					<!-- Drop zone for new tools -->
					<div
						class="w-12 h-12 border-2 border-dashed border-gray-300 rounded flex items-center justify-center text-gray-400 hover:border-gray-500 hover:text-gray-600 transition-colors cursor-pointer"
						title="Drop tools here"
					>
						+
					</div>
				</div>
			{/if}
		{/each}
	</div>
</div>

<style>
	/* Ensure proper spacing and alignment */
	.group:hover .group-hover\:opacity-100 {
		opacity: 1;
	}
</style>
