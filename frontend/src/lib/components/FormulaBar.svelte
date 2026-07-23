<script>
// @ts-nocheck
	import { createEventDispatcher } from 'svelte';
	import {
		selectedCell,
		formulaBar,
		showFormulaHelp,
		isFormulaRangeSelecting,
		formulaBeingEdited,
		selectedCells,
		selectCellRange,
		selectSingleCell
	} from '../stores/spreadsheet-store.js';

	import { getCellId, getCellFromId } from '../utils/spreadsheet-utils.js';

	const dispatch = createEventDispatcher();

	let nameBoxValue = '';

	// Update name box when selected cell changes
	$: if ($selectedCell) {
		const selectedCount = $selectedCells.size;
		if (selectedCount > 1) {
			// Show range if multiple cells selected
			const cells = Array.from($selectedCells).sort();
			nameBoxValue = `${cells[0]}:${cells[cells.length - 1]}`;
		} else {
			nameBoxValue = getCellId($selectedCell.row, $selectedCell.col);
		}
	}

	// Handle name box input for navigation
	function handleNameBoxKeyDown(event) {
		if (event.key === 'Enter') {
			navigateToCell(nameBoxValue);
		}
	}

	function navigateToCell(reference) {
		// Parse cell reference (e.g., "A1", "B2:D10")
		const rangeMatch = reference.match(/^([A-Z]+)(\d+):([A-Z]+)(\d+)$/);
		const singleMatch = reference.match(/^([A-Z]+)(\d+)$/);

		if (rangeMatch) {
			// Range selection
			const startCol = rangeMatch[1].charCodeAt(0) - 65;
			const startRow = parseInt(rangeMatch[2]) - 1;
			const endCol = rangeMatch[3].charCodeAt(0) - 65;
			const endRow = parseInt(rangeMatch[4]) - 1;

			if (isValidCell(startRow, startCol) && isValidCell(endRow, endCol)) {
				selectCellRange(startRow, startCol, endRow, endCol);
				selectedCell.set({ row: startRow, col: startCol });
				dispatch('action', { action: 'selectCell', row: startRow, col: startCol });
			}
		} else if (singleMatch) {
			// Single cell selection
			const col = singleMatch[1].charCodeAt(0) - 65;
			const row = parseInt(singleMatch[2]) - 1;

			if (isValidCell(row, col)) {
				selectSingleCell(row, col);
				dispatch('action', { action: 'selectCell', row, col });
			}
		}
	}

	function isValidCell(row, col) {
		return row >= 0 && row < 100 && col >= 0 && col < 26;
	}

	// Handle formula input to enable range selection for ANY function that needs ranges

	function handleFormulaInput(event) {
		const value = event.target.value;

		// Check if user is typing any function that expects a range parameter
		// Pattern: =FUNCTIONNAME( with no closing ) yet
		const functionMatch = value.toUpperCase().match(/^=([A-Z_]+)\([^)]*$/);

		if (functionMatch) {
			const functionName = functionMatch[1];

			// Common functions that expect ranges (can be extended)
			const rangeFunctions = [
				'SUM', 'AVERAGE', 'AVG', 'COUNT', 'MIN', 'MAX',
				'STDEV', 'STDEVP', 'VAR', 'VARP',
				'MEDIAN', 'MODE', 'RANK',
				'PRODUCT', 'GEOMEAN', 'HARMEAN',
				'SUBTOTAL', 'AGGREGATE'
			];

			if (rangeFunctions.includes(functionName)) {
				// Extract the partial formula up to the opening parenthesis
				const openParenIndex = value.indexOf('(');
				formulaBeingEdited.set(value.substring(0, openParenIndex + 1)); // "=FUNCTION("
				isFormulaRangeSelecting.set(true);
				console.log('Formula range selection enabled for:', functionName);
				return;
			}
		}

		// Reset formula range selection if formula is completed or changed
		isFormulaRangeSelecting.set(false);
		formulaBeingEdited.set('');
	}

	// Handle Enter key in formula bar
	function handleFormulaKeyDown(event) {
		if (event.key === 'Enter') {
			event.preventDefault();
			// Update the cell value and formula bar
			const { cellInput, spreadsheetData } = require('../stores/spreadsheet-store.js');
			if ($selectedCell) {
				spreadsheetData.update(data => {
					if (!data[$selectedCell.row]) data[$selectedCell.row] = [];
					data[$selectedCell.row][$selectedCell.col] = $formulaBar;
					return [...data];
				});
				cellInput.set($formulaBar);
			}
		} else if (event.key === 'Escape') {
			// Cancel range selection
			isFormulaRangeSelecting.set(false);
			formulaBeingEdited.set('');
		}
	}
</script>

<!-- Formula bar with Name Box -->
<div class="flex items-center px-4 py-2 border-b border-gray-200 bg-white">
	<!-- Name Box -->
	<div class="flex items-center mr-4">
		<span class="text-xs text-gray-500 mr-2">Name Box</span>
		<input
			type="text"
			bind:value={nameBoxValue}
			class="w-24 px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-medium"
			placeholder="A1"
			on:keydown={handleNameBoxKeyDown}
			title="Type cell reference (e.g., A1 or B2:D10) and press Enter"
		/>
	</div>

	<div class="w-px h-6 bg-gray-300 mx-2"></div>

	<!-- Formula Input -->
	<span class="text-sm font-medium text-gray-600 mr-2">fx</span>
	<div class="flex-1 relative">

		<input
			bind:value={$formulaBar}
			placeholder="Enter formula or value"
			class="w-full px-3 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
			on:input={handleFormulaInput}
			on:keydown={handleFormulaKeyDown}
		/>
		{#if $isFormulaRangeSelecting}
			<div class="absolute right-2 top-1/2 transform -translate-y-1/2 text-xs text-green-600 font-medium bg-green-100 px-2 py-1 rounded">
				Click and drag to select range
			</div>
		{/if}
	</div>
</div>

<style>
	input:focus {
		box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5);
	}
</style>
