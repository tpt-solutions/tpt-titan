<script>
// @ts-nocheck
	import {
		spreadsheetData,
		selectedCell,
		editingCell,
		selectedCells,
		cellFormats,
		frozenRows,
		frozenCols,
		sortState,
		filterState,
		isDragging,
		dragStartCell,
		isDraggingFillHandle,
		fillStartCell,
		fillEndCell,
		fillPreviewCells,
		showFillHandle,
		showContextMenu,
		contextMenuPosition,
		contextMenuTarget,
		contextMenuItems,
		formulaBar,
		isFormulaRangeSelecting,
		formulaRangeStart,
		formulaRangeEnd
	} from '../stores/spreadsheet-store.js';

	import {
		getCellId,
		getCellStyle,
		detectDataType
	} from '../utils/spreadsheet-utils.js';

	// Props
	export const mode = 'simple'; // 'simple' or 'advanced'

	// Column headers (A-Z)
	$: columns = Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i));

	// Row numbers (1-100)
	$: rows = Array.from({ length: 100 }, (_, i) => i + 1);

	// Computed display values
	function getDisplayValue(row, col) {
		const value = $spreadsheetData[row]?.[col] || '';
		if (value.startsWith('=')) {
			// In a real implementation, this would evaluate the formula
			// For now, just return the formula or a placeholder
			return value; // Simplified - would need formula evaluation
		}
		return value;
	}

	// Format cell value based on formatting rules
	function formatCellValue(value, row, col) {
		const format = $cellFormats.get(getCellId(row, col));

		if (!format || !value) return value;

		// Handle currency formatting
		if (format.currency) {
			const numValue = parseFloat(value);
			if (!isNaN(numValue)) {
				return new Intl.NumberFormat('en-US', {
					style: 'currency',
					currency: format.currency === true ? 'USD' : format.currency,
					minimumFractionDigits: format.currencyDecimals !== undefined ? format.currencyDecimals : 2
				}).format(numValue);
			}
		}

		// Handle number formatting
		if (format.numberFormat) {
			const numValue = parseFloat(value);
			if (!isNaN(numValue)) {
				switch (format.numberFormat) {
					case 'percent':
						return new Intl.NumberFormat('en-US', {
							style: 'percent',
							minimumFractionDigits: 1,
							maximumFractionDigits: 1
						}).format(numValue / 100);
					case 'decimal':
						return new Intl.NumberFormat('en-US', {
							minimumFractionDigits: format.decimals || 2,
							maximumFractionDigits: format.decimals || 2
						}).format(numValue);
					case 'scientific':
						return numValue.toExponential(2);
					default:
						return value;
				}
			}
		}

		return value;
	}

	// Check if cell is in current selection
	function isCellSelected(row, col) {
		const cellId = getCellId(row, col);
		return $selectedCells.has(cellId);
	}

	// Check if cell is in formula range selection
	function isCellInFormulaRange(row, col) {
		if (!$isFormulaRangeSelecting || !$formulaRangeStart || !$formulaRangeEnd) return false;

		const startRow = Math.min($formulaRangeStart.row, $formulaRangeEnd.row);
		const endRow = Math.max($formulaRangeStart.row, $formulaRangeEnd.row);
		const startCol = Math.min($formulaRangeStart.col, $formulaRangeEnd.col);
		const endCol = Math.max($formulaRangeStart.col, $formulaRangeEnd.col);

		return row >= startRow && row <= endRow && col >= startCol && col <= endCol;
	}

	// Check if row passes all filters
	function rowPassesFilters(rowIndex) {
		for (const [colIndex, allowedValues] of $filterState.entries()) {
			const cellValue = $spreadsheetData[rowIndex]?.[colIndex] || '';
			if (!allowedValues.has(cellValue)) {
				return false;
			}
		}
		return true;
	}

	// Dispatch events to parent
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	function dispatchAction(action, data = {}) {
		dispatch('action', { action, ...data });
	}

	// Cell interaction handlers
	function handleCellClick(row, col, event) {
		event.preventDefault();

		// Handle formula range selection mode
		if ($isFormulaRangeSelecting) {
			dispatchAction('endFormulaRangeSelection', { row, col });
			return;
		}

		// Handle multiple cell selection with Ctrl/Cmd
		if (event.ctrlKey || event.metaKey) {
			dispatchAction('toggleCellSelection', { row, col });
			return;
		}

		// Handle range selection with Shift+click
		if (event.shiftKey) {
			dispatchAction('extendCellSelection', { row, col });
			return;
		}

		// Single click - just select the cell (don't edit)
		dispatchAction('selectCell', { row, col });
	}

	function handleCellDoubleClick(row, col, event) {
		event.preventDefault();
		dispatchAction('startEditingCell', { row, col });
	}

	function handleCellMouseDown(row, col, event) {
		// Start potential drag selection
		if (!$isFormulaRangeSelecting) {
			dispatchAction('startCellDrag', { row, col, clientX: event.clientX, clientY: event.clientY });
		} else if ($isFormulaRangeSelecting) {
			dispatchAction('startFormulaRangeSelection', { row, col });
		}
	}

	function handleCellMouseEnter(row, col, event) {
		// Handle drag selection
		if ($isDragging) {
			dispatchAction('updateCellDrag', { row, col });
		} else if ($isFormulaRangeSelecting) {
			dispatchAction('updateFormulaRangeSelection', { row, col });
		}
	}

	function handleCellMouseUp(row, col, event) {
		if ($isDragging) {
			dispatchAction('endCellDrag', { row, col });
		}
	}

	function handleCellContextMenu(row, col, event) {
		event.preventDefault();
		event.stopPropagation();

		dispatchAction('showContextMenu', {
			row,
			col,
			x: Math.min(event.clientX, window.innerWidth - 200),
			y: Math.min(event.clientY, window.innerHeight - 300)
		});
	}

	// Column header click for sorting
	function handleColumnHeaderClick(colIndex, event) {
		if (event.shiftKey) {
			// Could implement multi-column sorting
			return;
		}

		dispatchAction('toggleColumnSort', { colIndex });
	}

	// Fill handle interaction
	function handleFillHandleMouseDown(event) {
		event.preventDefault();
		event.stopPropagation();

		dispatchAction('startFillHandleDrag', {
			row: $selectedCell.row,
			col: $selectedCell.col
		});
	}

	// Svelte action for focusing input
	function focusInput(node) {
		node.focus();
	}
</script>

<!-- Spreadsheet grid with freeze panes support -->
<div class="flex-1 overflow-auto">
	<table class="border-collapse">
		<!-- Column headers -->
		<thead class="sticky top-0 z-20">
			<tr>
				<th class="w-12 h-8 bg-blue-100 border border-gray-300 text-center text-xs font-medium text-gray-600 sticky left-0 z-30 {$frozenCols >= 1 ? 'bg-blue-200 border-blue-400' : ''}">
					#
				</th>
				{#each columns as col, colIndex}
					<th
						class="w-24 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600 {$frozenCols > colIndex ? 'sticky left-12 bg-blue-100 border-blue-400 z-25' : ''} {$frozenCols > colIndex + 1 ? 'bg-blue-200' : ''} relative group cursor-pointer hover:bg-gray-200"
						on:click={(event) => handleColumnHeaderClick(colIndex, event)}
						title="Click to sort • Right-click for filter menu"
					>
						<div class="flex items-center justify-center px-1">
							<span>{col}</span>
							<!-- Sort indicator -->
							{#if $sortState.has(colIndex)}
								<span class="ml-1 text-blue-600">
									{$sortState.get(colIndex).direction === 'asc' ? '↑' : '↓'}
								</span>
							{/if}
						</div>
						<!-- Filter indicator -->
						{#if $filterState.has(colIndex)}
							<div class="absolute top-0 right-0 w-2 h-2 bg-blue-500 rounded-full"></div>
						{/if}
					</th>
				{/each}
			</tr>
		</thead>

		<!-- Data rows -->
		<tbody>
			{#each rows as rowNum, rowIndex}
				{#if rowPassesFilters(rowIndex)}
					<tr class="hover:bg-gray-50">
						<!-- Row number -->
						<td class="w-12 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600 sticky left-0 z-10">
							{rowNum}
						</td>

						<!-- Data cells -->
						{#each columns as col, colIndex}
							<td
								class="w-24 h-8 border border-gray-300 p-0 relative cursor-cell
									{$selectedCell?.row === rowIndex && $selectedCell?.col === colIndex && !$editingCell ? 'ring-2 ring-blue-500 ring-inset bg-blue-50' : ''}
									{$editingCell?.row === rowIndex && $editingCell?.col === colIndex ? 'ring-2 ring-green-500 ring-inset' : ''}
									{isCellSelected(rowIndex, colIndex) && !($selectedCell?.row === rowIndex && $selectedCell?.col === colIndex) ? 'bg-blue-100 ring-1 ring-blue-400' : ''}
									{isCellInFormulaRange(rowIndex, colIndex) ? 'bg-green-100 ring-2 ring-green-500' : ''}
									{$fillPreviewCells.has(getCellId(rowIndex, colIndex)) ? 'bg-yellow-200 ring-1 ring-yellow-400' : ''}
									hover:bg-gray-50 transition-colors"
								style={getCellStyle($cellFormats.get(getCellId(rowIndex, colIndex)))}
								on:click={(event) => handleCellClick(rowIndex, colIndex, event)}
								on:contextmenu={(event) => handleCellContextMenu(rowIndex, colIndex, event)}
								on:dblclick={(event) => handleCellDoubleClick(rowIndex, colIndex, event)}
								on:mousedown={(event) => handleCellMouseDown(rowIndex, colIndex, event)}
								on:mouseenter={(event) => handleCellMouseEnter(rowIndex, colIndex, event)}
								on:mouseup={(event) => handleCellMouseUp(rowIndex, colIndex, event)}
							>
								{#if $editingCell?.row === rowIndex && $editingCell?.col === colIndex}
									<input
										bind:value={$formulaBar}
										class="w-full h-full px-2 py-1 text-sm border-none outline-none bg-white"
										on:keydown={(event) => {
											if (event.key === 'Enter') {
												dispatchAction('stopEditingCell', { save: true });
											}
											if (event.key === 'Escape') {
												dispatchAction('stopEditingCell', { save: false });
											}
										}}
										on:blur={() => dispatchAction('stopEditingCell', { save: true })}
										use:focusInput
									/>
								{:else}
									<div class="w-full h-full px-2 py-1 text-sm overflow-hidden whitespace-nowrap select-none">
										{formatCellValue(getDisplayValue(rowIndex, colIndex), rowIndex, colIndex)}
									</div>

									<!-- Fill Handle -->
									{#if $selectedCell?.row === rowIndex && $selectedCell?.col === colIndex && !$editingCell && $showFillHandle}
										<div
											class="absolute bottom-0 right-0 w-2 h-2 bg-blue-500 border border-white cursor-crosshair z-10"
											on:mousedown={handleFillHandleMouseDown}
											title="Drag to auto-fill adjacent cells"
											role="button"
											tabindex="0"
										></div>
									{/if}
								{/if}
							</td>
						{/each}
					</tr>
				{/if}
			{/each}
		</tbody>
	</table>
</div>

<style>
	/* Ensure table cells don't have unwanted padding/margin */
	td {
		margin: 0;
		padding: 0;
	}

	/* Custom focus styles for inputs */
	input:focus {
		outline: none;
		box-shadow: inset 0 0 0 2px rgba(59, 130, 246, 0.5);
	}
</style>
