<script>
	import { onMount, createEventDispatcher } from 'svelte';

	export let mode = 'simple'; // 'simple' or 'advanced'

	// Spreadsheet state
	let data = Array(100).fill().map(() => Array(26).fill('')); // 100 rows, 26 columns (A-Z)
	let selectedCell = { row: 0, col: 0 };
	let editingCell = null;
	let cellInput = '';
	let formulaBar = '';

	const dispatch = createEventDispatcher();

	// Column headers (A-Z)
	$: columns = Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i));

	// Row numbers (1-100)
	$: rows = Array.from({ length: 100 }, (_, i) => i + 1);

	function getCellId(row, col) {
		return `${columns[col]}${row + 1}`;
	}

	function selectCell(row, col) {
		selectedCell = { row, col };
		const value = data[row][col];
		formulaBar = value.startsWith('=') ? value : '';
		cellInput = value;
	}

	function startEditing(row, col) {
		editingCell = { row, col };
		cellInput = data[row][col];
	}

	function stopEditing(save = true) {
		if (editingCell && save) {
			data[editingCell.row][editingCell.col] = cellInput;
			formulaBar = cellInput.startsWith('=') ? cellInput : '';
		}
		editingCell = null;
	}

	function handleKeyDown(event) {
		const { row, col } = selectedCell;

		switch (event.key) {
			case 'Enter':
				if (editingCell) {
					stopEditing();
				} else {
					startEditing(row, col);
				}
				break;
			case 'Escape':
				stopEditing(false);
				break;
			case 'Tab':
				event.preventDefault();
				const nextCol = event.shiftKey ? Math.max(0, col - 1) : Math.min(25, col + 1);
				selectCell(row, nextCol);
				break;
			case 'ArrowUp':
				selectCell(Math.max(0, row - 1), col);
				break;
			case 'ArrowDown':
				selectCell(Math.min(99, row + 1), col);
				break;
			case 'ArrowLeft':
				selectCell(row, Math.max(0, col - 1));
				break;
			case 'ArrowRight':
				selectCell(row, Math.min(25, col + 1));
				break;
		}
	}

	// Evaluate basic formulas (for demo)
	function evaluateFormula(formula) {
		if (!formula.startsWith('=')) return formula;

		try {
			// Simple SUM function
			const sumMatch = formula.match(/^=SUM\((.+)\)$/i);
			if (sumMatch) {
				const range = sumMatch[1];
				const [start, end] = range.split(':');
				// For demo, just return a placeholder
				return `SUM(${range})`;
			}

			// Basic arithmetic
			return formula.substring(1); // Remove = and return as-is for demo
		} catch (error) {
			return '#ERROR!';
		}
	}

	// Get display value for cell
	function getDisplayValue(row, col) {
		const value = data[row][col];
		if (value.startsWith('=')) {
			return evaluateFormula(value);
		}
		return value;
	}

	// Svelte action for focusing input
	function focusInput(node) {
		node.focus();
	}

	// Export data to CSV
	function exportToCSV() {
		let csv = '';
		for (let row = 0; row < data.length; row++) {
			for (let col = 0; col < data[row].length; col++) {
				if (col > 0) csv += ',';
				let cellValue = data[row][col];
				// Escape quotes and wrap in quotes if necessary
				if (cellValue.includes(',') || cellValue.includes('"') || cellValue.includes('\n')) {
					cellValue = '"' + cellValue.replace(/"/g, '""') + '"';
				}
				csv += cellValue;
			}
			csv += '\n';
		}

		const blob = new Blob([csv], { type: 'text/csv' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = 'spreadsheet.csv';
		a.click();
		URL.revokeObjectURL(url);
	}

	// Import data from CSV
	function importFromCSV(event) {
		const file = event.target.files[0];
		if (!file) return;

		const reader = new FileReader();
		reader.onload = (e) => {
			const csv = e.target.result;
			const rows = csv.split('\n').filter(row => row.trim() !== '');
			const newData = [];

			for (let i = 0; i < Math.min(rows.length, 100); i++) {
				const cols = parseCSVRow(rows[i]);
				newData.push(cols.slice(0, 26));
			}

			// Pad to 100 rows and 26 columns
			while (newData.length < 100) {
				newData.push(Array(26).fill(''));
			}
			for (let row of newData) {
				while (row.length < 26) {
					row.push('');
				}
			}

			data = newData;
			selectedCell = { row: 0, col: 0 };
		};
		reader.readAsText(file);
	}

	// Parse a single CSV row, handling quotes
	function parseCSVRow(row) {
		const result = [];
		let current = '';
		let inQuotes = false;

		for (let i = 0; i < row.length; i++) {
			const char = row[i];
			if (char === '"') {
				if (inQuotes && row[i + 1] === '"') {
					current += '"';
					i++; // Skip next quote
				} else {
					inQuotes = !inQuotes;
				}
			} else if (char === ',' && !inQuotes) {
				result.push(current);
				current = '';
			} else {
				current += char;
			}
		}
		result.push(current);
		return result;
	}

	onMount(() => {
		// Initialize with some sample data
		data[0][0] = 'Product';
		data[0][1] = 'Price';
		data[0][2] = 'Quantity';
		data[0][3] = 'Total';

		data[1][0] = 'Widget A';
		data[1][1] = '10.50';
		data[1][2] = '5';
		data[1][3] = '=B2*C2';

		data[2][0] = 'Widget B';
		data[2][1] = '15.75';
		data[2][2] = '3';
		data[2][3] = '=B3*C3';
	});
</script>

<svelte:window on:keydown={handleKeyDown} />

<div class="flex flex-col h-full">
	<!-- Toolbar -->
	<div class="flex items-center px-4 py-2 border-b border-gray-200 bg-gray-50">
		<button
			class="px-3 py-1 bg-blue-500 text-white text-sm rounded hover:bg-blue-600 mr-2"
			on:click={exportToCSV}
		>
			Export CSV
		</button>
		<label class="px-3 py-1 bg-green-500 text-white text-sm rounded hover:bg-green-600 cursor-pointer mr-2">
			Import CSV
			<input
				type="file"
				accept=".csv"
				on:change={importFromCSV}
				class="hidden"
			/>
		</label>
	</div>

	<!-- Formula bar -->
	<div class="flex items-center px-4 py-2 border-b border-gray-200 bg-white">
		<span class="text-sm font-medium text-gray-600 mr-2 w-16">
			{selectedCell ? getCellId(selectedCell.row, selectedCell.col) : ''}
		</span>
		<input
			bind:value={formulaBar}
			placeholder="Enter formula or value"
			class="flex-1 px-3 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
			on:keydown={(e) => {
				if (e.key === 'Enter') {
					if (selectedCell) {
						data[selectedCell.row][selectedCell.col] = formulaBar;
						cellInput = formulaBar;
					}
				}
			}}
		/>
	</div>

	<!-- Spreadsheet grid -->
	<div class="flex-1 overflow-auto">
		<table class="border-collapse">
			<!-- Column headers -->
			<thead class="sticky top-0 z-10">
				<tr>
					<th class="w-12 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600">
						#
					</th>
					{#each columns as col, colIndex}
						<th class="w-24 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600">
							{col}
						</th>
					{/each}
				</tr>
			</thead>

			<!-- Data rows -->
			<tbody>
				{#each rows as rowNum, rowIndex}
					<tr class="hover:bg-gray-50">
						<!-- Row number -->
						<td class="w-12 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600 sticky left-0 z-10">
							{rowNum}
						</td>

						<!-- Data cells -->
						{#each columns as col, colIndex}
							<td
								class="w-24 h-8 border border-gray-300 p-0 relative cursor-cell {selectedCell?.row === rowIndex && selectedCell?.col === colIndex ? 'ring-2 ring-blue-500 ring-inset' : ''}"
								on:click={() => selectCell(rowIndex, colIndex)}
								on:dblclick={() => startEditing(rowIndex, colIndex)}
							>
								{#if editingCell?.row === rowIndex && editingCell?.col === colIndex}
									<input
										bind:value={cellInput}
										class="w-full h-full px-2 py-1 text-sm border-none outline-none bg-white"
										on:keydown={(e) => {
											if (e.key === 'Enter') stopEditing();
											if (e.key === 'Escape') stopEditing(false);
										}}
										on:blur={() => stopEditing()}
										use:focusInput
									/>
								{:else}
									<div class="w-full h-full px-2 py-1 text-sm overflow-hidden whitespace-nowrap">
										{getDisplayValue(rowIndex, colIndex)}
									</div>
								{/if}
							</td>
						{/each}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>

<!-- AI Formula Assistant (Simple Mode) -->
{#if mode === 'simple' && selectedCell}
	<div class="fixed bottom-4 right-4 bg-white border border-gray-200 rounded-lg shadow-lg p-4 max-w-sm">
		<h3 class="text-sm font-medium text-gray-900 mb-2">💡 AI Assistant</h3>
		<div class="text-xs text-gray-600 space-y-1">
			<p>Try typing:</p>
			<button class="block w-full text-left p-2 hover:bg-gray-50 rounded text-blue-600" on:click={() => formulaBar = '=SUM(A1:A10)'}>
				"=SUM(A1:A10)" - Add up a range
			</button>
			<button class="block w-full text-left p-2 hover:bg-gray-50 rounded text-blue-600" on:click={() => formulaBar = '=AVERAGE(B1:B10)'}>
				"=AVERAGE(B1:B10)" - Average values
			</button>
		</div>
	</div>
{/if}
