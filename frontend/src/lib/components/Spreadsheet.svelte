<script>
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';

	// Import extracted components
	import QuickAccessToolbar from './QuickAccessToolbar.svelte';
	import SpreadsheetMenuBar from './SpreadsheetMenuBar.svelte';
	import SpreadsheetRibbon from './SpreadsheetRibbon.svelte';
	import FormulaBar from './FormulaBar.svelte';
	import SpreadsheetGrid from './SpreadsheetGrid.svelte';
	import SheetTabs from './SheetTabs.svelte';
	import SpreadsheetStatusBar from './SpreadsheetStatusBar.svelte';
	import SpreadsheetModals from './SpreadsheetModals.svelte';


	// Import stores and utilities
	import {
		mode as modeStore,
		selectedTemplate as selectedTemplateStore,
		resetSpreadsheet,
		spreadsheetData,
		cellFormats,
		selectedCells,
		statusBarInfo,
		canUndo,
		canRedo
	} from '../stores/spreadsheet-store.js';


	// Import history management
	import SpreadsheetHistory, { createDebouncedPush, handleUndoRedoKeyboard } from '../utils/spreadsheet-history.js';


	export let mode = 'simple'; // 'simple' or 'advanced'
	export let selectedTemplate = null;

	$: modeStore.set(mode);


	// Event dispatcher for parent communication
	const dispatch = createEventDispatcher();

	// History management
	let spreadsheetHistory = new SpreadsheetHistory();
	let debouncedHistoryPush;

	// Update undo/redo availability in store
	function updateHistoryState() {
		canUndo.set(spreadsheetHistory.canUndo());
		canRedo.set(spreadsheetHistory.canRedo());
	}


	// Save current state to history
	function saveToHistory(actionType = 'cellChange') {
		const currentData = $spreadsheetData;
		const currentFormats = $cellFormats;
		
		debouncedHistoryPush({
			data: currentData,
			formats: currentFormats,
			actionType,
			timestamp: Date.now()
		});
		updateHistoryState();
	}

	// Undo handler
	function handleUndo() {
		const state = spreadsheetHistory.undo();
		if (state) {
			spreadsheetData.set(state.data);
			cellFormats.set(state.formats);
			updateHistoryState();
		}
	}

	// Redo handler
	function handleRedo() {
		const state = spreadsheetHistory.redo();
		if (state) {
			spreadsheetData.set(state.data);
			cellFormats.set(state.formats);
			updateHistoryState();
		}
	}

	// Update status bar when selection changes
	$: {
		const cells = Array.from($selectedCells);
		const data = $spreadsheetData;
		
		let sum = 0;
		let count = 0;
		let min = null;
		let max = null;
		let hasNumericValues = false;

		cells.forEach(cellId => {
			const match = cellId.match(/([A-Z]+)(\d+)/);
			if (match) {
				const col = match[1].charCodeAt(0) - 65;
				const row = parseInt(match[2]) - 1;
				const value = data[row]?.[col];
				const numValue = parseFloat(value);
				
				if (!isNaN(numValue) && value !== '') {
					sum += numValue;
					count++;
					hasNumericValues = true;
					if (min === null || numValue < min) min = numValue;
					if (max === null || numValue > max) max = numValue;
				}
			}
		});

		statusBarInfo.set({
			selectedCount: cells.length,
			sum: hasNumericValues ? sum : 0,
			average: hasNumericValues && count > 0 ? sum / count : 0,
			min: hasNumericValues ? min : null,
			max: hasNumericValues ? max : null
		});
	}


	// Global keyboard handler for undo/redo
	function handleGlobalKeyDown(event) {
		if (handleUndoRedoKeyboard(event, handleUndo, handleRedo)) {
			return;
		}
	}


	// Action handlers from child components - orchestrate actions between components
	function handleAction(event) {
		const { action, ...data } = event.detail;

		switch (action) {
			case 'save':
				dispatch('save');
				break;
			case 'exportExcel':
				dispatch('exportExcel');
				break;
			case 'exportCSV':
				dispatch('exportCSV');
				break;
			case 'import':
				dispatch('import', data);
				break;
			case 'share':
				dispatch('share');
				break;
			case 'applyFormatting':
				dispatch('applyFormatting', data);
				break;
			case 'applyNumberFormat':
				dispatch('applyNumberFormat', data);
				break;
			case 'applyBorder':
				dispatch('applyBorder', data);
				break;
			case 'insertRowAbove':
			case 'insertRowBelow':
			case 'insertColumnLeft':
			case 'insertColumnRight':
				dispatch('insert', { type: action, ...data });
				break;
			case 'deleteRow':
			case 'deleteColumn':
				dispatch('delete', { type: action, ...data });
				break;
			case 'clearContents':
				dispatch('clearContents', data);
				break;
			case 'sortColumn':
				dispatch('sortColumn', data);
				break;
			case 'showFilterDialog':
				dispatch('showFilterDialog', data);
				break;
			case 'showDataPrepTools':
				dispatch('showDataPrepTools', data);
				break;
			case 'showPivotTable':
				dispatch('showPivotTable', data);
				break;
			case 'insertFormula':
				dispatch('insertFormula', data);
				break;
			case 'copy':
				dispatch('copy');
				break;
			case 'paste':
				dispatch('paste');
				break;
			case 'cut':
				dispatch('cut');
				break;
			case 'undo':
				handleUndo();
				break;
			case 'redo':
				handleRedo();
				break;

			case 'selectAll':
				dispatch('selectAll');
				break;
			case 'newSpreadsheet':
				resetSpreadsheet();
				break;
			case 'openSpreadsheet':
				dispatch('openSpreadsheet');
				break;
			case 'saveAsSpreadsheet':
				dispatch('saveAsSpreadsheet');
				break;
			case 'printSpreadsheet':
				dispatch('printSpreadsheet');
				break;
			case 'zoomIn':
			case 'zoomOut':
			case 'resetZoom':
				dispatch('zoom', { action });
				break;
			case 'toggleGridlines':
			case 'toggleHeaders':
				dispatch('toggleView', { option: action.replace('toggle', '').toLowerCase() });
				break;
			case 'freezePanes':
				dispatch('freezePanes');
				break;
			case 'functionReference':
				dispatch('functionReference');
				break;
			case 'spellCheck':
			case 'autoCorrectOptions':
				dispatch('tools', { action });
				break;
			case 'removeToolFromRibbon':
				dispatch('customizeRibbon', { action: 'remove', ...data });
				break;
			default:
				// Forward unknown actions to parent
				dispatch(action, data);
		}
	}

	// Lifecycle
	onMount(() => {
		// Initialize any required setup
		console.log('Spreadsheet component mounted');
		
		// Initialize history with debounced push
		debouncedHistoryPush = createDebouncedPush(spreadsheetHistory, 300);
		
		// Add keyboard listener for undo/redo
		document.addEventListener('keydown', handleGlobalKeyDown);
		
		// Save initial state
		saveToHistory('initial');
	});

	onDestroy(() => {
		// Cleanup
		document.removeEventListener('keydown', handleGlobalKeyDown);
	});


	// Apply template when selectedTemplate changes
	$: if (selectedTemplate && selectedTemplate.data) {
		console.log('Applying template:', selectedTemplate.name);
		console.log('Template data:', selectedTemplate.data);

		// Apply template data to spreadsheet
		const templateData = selectedTemplate.data;
		const numRows = templateData.length;
		const numCols = Math.max(...templateData.map(row => row.length));

		// Create new data array with template data
		const newData = Array(Math.max(100, numRows)).fill().map(() => Array(Math.max(26, numCols)).fill(''));

		// Populate with template data
		templateData.forEach((row, rowIndex) => {
			row.forEach((cellValue, colIndex) => {
				if (cellValue !== undefined && cellValue !== null) {
					newData[rowIndex][colIndex] = String(cellValue);
				}
			});
		});

		// Update the spreadsheet data store
		spreadsheetData.set(newData);
		console.log('Template applied successfully, new data length:', newData.length);

		// Apply styles if available
		if (selectedTemplate.styles) {
			console.log('Applying template styles:', selectedTemplate.styles);
			const stylesMap = new Map();
			Object.entries(selectedTemplate.styles).forEach(([range, style]) => {
				stylesMap.set(range, style);
			});
			cellFormats.set(stylesMap);
		}
		
		// Save to history after template application
		setTimeout(() => saveToHistory('templateApply'), 0);
	}

</script>



<div class="flex flex-col h-full bg-white">
	{#if mode === 'advanced'}
		<!-- Quick Access Toolbar -->
		<QuickAccessToolbar on:action={handleAction} />

		<!-- Menu Bar -->
		<SpreadsheetMenuBar on:action={handleAction} />

		<!-- Customizable Ribbon -->
		<SpreadsheetRibbon on:action={handleAction} />
	{/if}

	<!-- Formula bar with Name Box -->
	<FormulaBar on:action={handleAction} />

	<!-- Spreadsheet grid -->
	<div class="flex-1 overflow-hidden relative">
		<SpreadsheetGrid on:action={handleAction} />
	</div>

	<!-- Sheet Tabs -->
	<SheetTabs on:action={handleAction} />

	<!-- Status Bar -->
	<SpreadsheetStatusBar on:action={handleAction} />
</div>

<!-- Modals -->
<SpreadsheetModals on:action={handleAction} />
