<script>
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';

	// Import extracted components
	import SpreadsheetMenuBar from './SpreadsheetMenuBar.svelte';
	import SpreadsheetRibbon from './SpreadsheetRibbon.svelte';
	import SpreadsheetToolbar from './SpreadsheetToolbar.svelte';
	import FormulaBar from './FormulaBar.svelte';
	import SpreadsheetGrid from './SpreadsheetGrid.svelte';
	import SpreadsheetModals from './SpreadsheetModals.svelte';

	// Import stores and utilities
	import {
		mode as modeStore,
		selectedTemplate as selectedTemplateStore,
		resetSpreadsheet
	} from '../stores/spreadsheet-store.js';

	export let mode = 'simple'; // 'simple' or 'advanced'
	export let selectedTemplate = null;

	// Event dispatcher for parent communication
	const dispatch = createEventDispatcher();

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
				dispatch('undo');
				break;
			case 'redo':
				dispatch('redo');
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
	});

	onDestroy(() => {
		// Cleanup if needed
	});
</script>



<div class="flex flex-col h-full">
	<!-- Menu Bar -->
	<SpreadsheetMenuBar on:action={handleAction} />

	<!-- Customizable Ribbon -->
	<SpreadsheetRibbon on:action={handleAction} />

	<!-- Professional Toolbar -->
	<SpreadsheetToolbar on:action={handleAction} />

	<!-- Formula bar -->
	<FormulaBar on:action={handleAction} />

	<!-- Spreadsheet grid -->
	<SpreadsheetGrid on:action={handleAction} />
</div>

<!-- Modals -->
<SpreadsheetModals on:action={handleAction} />
