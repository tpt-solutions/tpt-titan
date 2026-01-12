<script>
	import {
		selectedCell,
		formulaBar,
		showFormulaHelp,
		isFormulaRangeSelecting,
		formulaBeingEdited
	} from '../stores/spreadsheet-store.js';

	import { getCellId } from '../utils/spreadsheet-utils.js';

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

<!-- Formula bar -->
<div class="flex items-center px-4 py-2 border-b border-gray-200 bg-white">
	<span class="text-sm font-medium text-gray-600 mr-2 w-16">
		{$selectedCell ? getCellId($selectedCell.row, $selectedCell.col) : ''}
	</span>
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
