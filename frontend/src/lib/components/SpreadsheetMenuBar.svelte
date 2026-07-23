<script>
// @ts-nocheck
	import {
		showFileMenu,
		showEditMenu,
		showFormatMenu,
		showViewMenu,
		showInsertMenu,
		showToolsMenu,
		showFreezePaneDialog,
		showFindReplace,
		showFormulaHelp
	} from '../stores/spreadsheet-store.js';

	import { dataToCSV } from '../utils/spreadsheet-utils.js';

	// Menu toggle functions
	function toggleFileMenu() {
		showFileMenu.update(v => !v);
		showEditMenu.set(false);
		showFormatMenu.set(false);
		showViewMenu.set(false);
		showInsertMenu.set(false);
		showToolsMenu.set(false);
	}

	function toggleEditMenu() {
		showEditMenu.update(v => !v);
		showFileMenu.set(false);
		showFormatMenu.set(false);
		showViewMenu.set(false);
		showInsertMenu.set(false);
		showToolsMenu.set(false);
	}

	function toggleFormatMenu() {
		showFormatMenu.update(v => !v);
		showFileMenu.set(false);
		showEditMenu.set(false);
		showViewMenu.set(false);
		showInsertMenu.set(false);
		showToolsMenu.set(false);
	}

	function toggleViewMenu() {
		showViewMenu.update(v => !v);
		showFileMenu.set(false);
		showEditMenu.set(false);
		showFormatMenu.set(false);
		showInsertMenu.set(false);
		showToolsMenu.set(false);
	}

	function toggleInsertMenu() {
		showInsertMenu.update(v => !v);
		showFileMenu.set(false);
		showEditMenu.set(false);
		showFormatMenu.set(false);
		showViewMenu.set(false);
		showToolsMenu.set(false);
	}

	function toggleToolsMenu() {
		showToolsMenu.update(v => !v);
		showFileMenu.set(false);
		showEditMenu.set(false);
		showFormatMenu.set(false);
		showViewMenu.set(false);
		showInsertMenu.set(false);
	}

	// Menu action handlers (these will be passed from parent component)
	let handlers = {};

	// Dispatch events to parent component
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	function dispatchAction(action, data = {}) {
		dispatch('action', { action, ...data });
	}

	// Export to CSV
	function exportToCSV() {
		dispatchAction('exportCSV');
		showFileMenu.set(false);
	}

	// File menu actions
	function newSpreadsheet() {
		dispatchAction('newSpreadsheet');
		showFileMenu.set(false);
	}

	function openSpreadsheet() {
		dispatchAction('openSpreadsheet');
		showFileMenu.set(false);
	}

	function saveAsSpreadsheet() {
		dispatchAction('saveAsSpreadsheet');
		showFileMenu.set(false);
	}

	function printSpreadsheet() {
		dispatchAction('printSpreadsheet');
		showFileMenu.set(false);
	}

	// Edit menu actions
	function undo() {
		dispatchAction('undo');
		showEditMenu.set(false);
	}

	function redo() {
		dispatchAction('redo');
		showEditMenu.set(false);
	}

	function cutSelection() {
		dispatchAction('cut');
		showEditMenu.set(false);
	}

	function copySelection() {
		dispatchAction('copy');
		showEditMenu.set(false);
	}

	function pasteClipboard() {
		dispatchAction('paste');
		showEditMenu.set(false);
	}

	function findReplace() {
		showFindReplace.set(true);
		showEditMenu.set(false);
	}

	function selectAll() {
		dispatchAction('selectAll');
		showEditMenu.set(false);
	}

	function clearContents() {
		dispatchAction('clearContents');
		showEditMenu.set(false);
	}

	// Format menu actions
	function formatCells() {
		dispatchAction('formatCells');
		showFormatMenu.set(false);
	}

	function applyBold() {
		dispatchAction('applyFormatting', { type: 'bold' });
		showFormatMenu.set(false);
	}

	function applyItalic() {
		dispatchAction('applyFormatting', { type: 'italic' });
		showFormatMenu.set(false);
	}

	function applyUnderline() {
		dispatchAction('applyFormatting', { type: 'underline' });
		showFormatMenu.set(false);
	}

	function alignLeft() {
		dispatchAction('applyFormatting', { type: 'align', value: 'left' });
		showFormatMenu.set(false);
	}

	function alignCenter() {
		dispatchAction('applyFormatting', { type: 'align', value: 'center' });
		showFormatMenu.set(false);
	}

	function alignRight() {
		dispatchAction('applyFormatting', { type: 'align', value: 'right' });
		showFormatMenu.set(false);
	}

	// View menu actions
	function zoomIn() {
		dispatchAction('zoomIn');
		showViewMenu.set(false);
	}

	function zoomOut() {
		dispatchAction('zoomOut');
		showViewMenu.set(false);
	}

	function resetZoom() {
		dispatchAction('resetZoom');
		showViewMenu.set(false);
	}

	function toggleGridlines() {
		dispatchAction('toggleGridlines');
		showViewMenu.set(false);
	}

	function toggleHeaders() {
		dispatchAction('toggleHeaders');
		showViewMenu.set(false);
	}

	function freezePanes() {
		showFreezePaneDialog.set(true);
		showViewMenu.set(false);
	}

	// Insert menu actions
	function insertRowAbove() {
		dispatchAction('insertRowAbove');
		showInsertMenu.set(false);
	}

	function insertRowBelow() {
		dispatchAction('insertRowBelow');
		showInsertMenu.set(false);
	}

	function insertColumnLeft() {
		dispatchAction('insertColumnLeft');
		showInsertMenu.set(false);
	}

	function insertColumnRight() {
		dispatchAction('insertColumnRight');
		showInsertMenu.set(false);
	}

	function insertChart() {
		dispatchAction('insertChart');
		showInsertMenu.set(false);
	}

	function insertPivotTable() {
		dispatchAction('insertPivotTable');
		showInsertMenu.set(false);
	}

	// Tools menu actions
	function dataPreparation() {
		dispatchAction('dataPreparation');
		showToolsMenu.set(false);
	}

	function dataAnalysis() {
		dispatchAction('dataAnalysis');
		showToolsMenu.set(false);
	}

	function showFormulaHelpDialog() {
		showFormulaHelp.set(true);
		showToolsMenu.set(false);
	}

	function functionReference() {
		dispatchAction('functionReference');
		showToolsMenu.set(false);
	}

	function spellCheck() {
		dispatchAction('spellCheck');
		showToolsMenu.set(false);
	}

	function autoCorrectOptions() {
		dispatchAction('autoCorrectOptions');
		showToolsMenu.set(false);
	}
</script>

<!-- Menu Bar -->
<div class="bg-gray-800 text-white px-4 py-1 flex items-center space-x-6 text-sm">
	<!-- File Menu -->
	<div class="relative">
		<button
			class="hover:bg-gray-700 px-3 py-1 rounded transition-colors {$showFileMenu ? 'bg-gray-700' : ''}"
			on:click={toggleFileMenu}
		>
			File
		</button>
		{#if $showFileMenu}
			<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50">
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={newSpreadsheet}>
					📄 New Spreadsheet
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={openSpreadsheet}>
					📂 Open...
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center justify-between" on:click={() => dispatchAction('save')}>
					<span>💾 Save</span>
					<span class="text-xs text-gray-500">Ctrl+S</span>
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center justify-between" on:click={saveAsSpreadsheet}>
					<span>💾 Save As...</span>
					<span class="text-xs text-gray-500">Ctrl+Shift+S</span>
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={exportToCSV}>
					📊 Export to CSV
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={() => dispatchAction('exportExcel')}>
					📈 Export to Excel
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={printSpreadsheet}>
					🖨️ Print...
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={() => dispatchAction('share')}>
					👥 Share...
				</button>
			</div>
		{/if}
	</div>

	<!-- Edit Menu -->
	<div class="relative">
		<button
			class="hover:bg-gray-700 px-3 py-1 rounded transition-colors {$showEditMenu ? 'bg-gray-700' : ''}"
			on:click={toggleEditMenu}
		>
			Edit
		</button>
		{#if $showEditMenu}
			<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50">
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={undo}>
					↶ Undo
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={redo}>
					↷ Redo
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={copySelection}>
					📋 Copy
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={cutSelection}>
					✂️ Cut
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={pasteClipboard}>
					📌 Paste
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={findReplace}>
					🔍 Find & Replace...
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={selectAll}>
					📋 Select All
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={clearContents}>
					🧹 Clear Contents
				</button>
			</div>
		{/if}
	</div>

	<!-- Format Menu -->
	<div class="relative">
		<button
			class="hover:bg-gray-700 px-3 py-1 rounded transition-colors {$showFormatMenu ? 'bg-gray-700' : ''}"
			on:click={toggleFormatMenu}
		>
			Format
		</button>
		{#if $showFormatMenu}
			<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50">
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={formatCells}>
					🎨 Format Cells...
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={applyBold}>
					<strong>B</strong> Bold
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={applyItalic}>
					<em>I</em> Italic
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={applyUnderline}>
					<u>U</u> Underline
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={alignLeft}>
					⬅️ Align Left
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={alignCenter}>
					⬌ Center
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={alignRight}>
					➡️ Align Right
				</button>
			</div>
		{/if}
	</div>

	<!-- View Menu -->
	<div class="relative">
		<button
			class="hover:bg-gray-700 px-3 py-1 rounded transition-colors {$showViewMenu ? 'bg-gray-700' : ''}"
			on:click={toggleViewMenu}
		>
			View
		</button>
		{#if $showViewMenu}
			<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50">
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={zoomIn}>
					🔍 Zoom In
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={zoomOut}>
					🔎 Zoom Out
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={resetZoom}>
					📏 100% Zoom
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={toggleGridlines}>
					📊 Show Gridlines
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={toggleHeaders}>
					📋 Show Headers
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={freezePanes}>
					🧊 Freeze Panes
				</button>
			</div>
		{/if}
	</div>

	<!-- Insert Menu -->
	<div class="relative">
		<button
			class="hover:bg-gray-700 px-3 py-1 rounded transition-colors {$showInsertMenu ? 'bg-gray-700' : ''}"
			on:click={toggleInsertMenu}
		>
			Insert
		</button>
		{#if $showInsertMenu}
			<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50">
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={insertRowAbove}>
					➕ Insert Row Above
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={insertRowBelow}>
					➕ Insert Row Below
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={insertColumnLeft}>
					➕ Insert Column Left
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={insertColumnRight}>
					➕ Insert Column Right
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={insertPivotTable}>
					📋 Insert Pivot Table
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={insertChart}>
					📈 Insert Chart
				</button>
			</div>
		{/if}
	</div>

	<!-- Tools Menu -->
	<div class="relative">
		<button
			class="hover:bg-gray-700 px-3 py-1 rounded transition-colors {$showToolsMenu ? 'bg-gray-700' : ''}"
			on:click={toggleToolsMenu}
		>
			Tools
		</button>
		{#if $showToolsMenu}
			<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50">
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={dataPreparation}>
					🔧 Data Preparation
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={dataAnalysis}>
					📊 Data Analysis
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={showFormulaHelpDialog}>
					🔢 Formula Help
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={functionReference}>
					📋 Function Reference
				</button>
				<div class="border-t border-gray-200 my-1"></div>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={spellCheck}>
					🔍 Spelling Check
				</button>
				<button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" on:click={autoCorrectOptions}>
					📝 AutoCorrect Options
				</button>
			</div>
		{/if}
	</div>
</div>
