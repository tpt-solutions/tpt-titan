<script>
	import {
		createSpreadsheet,
		exportSpreadsheetToExcel,
		importExcelToSpreadsheet,
		formatApiError
	} from '../api.js';

	import {
		workbookName,
		hasUnsavedChanges,
		isSaving,
		saveStatus,
		mode,
		showPivotTableBuilder,
		showDataPrepTools,
		showFindReplace,
		showFormulaHelp,
		spreadsheetData,
		cellFormats,
		currentWorkbook
	} from '../stores/spreadsheet-store.js';

	// File operations
	async function saveSpreadsheet() {
		if ($isSaving) return;

		isSaving.set(true);
		saveStatus.set('Saving...');

		try {
			const data = $spreadsheetData;
			const spreadsheetDataObj = {
				name: $workbookName,
				data: {},
				formulas: {}
			};

			// Convert data array to cell reference format
			for (let row = 0; row < data.length; row++) {
				for (let col = 0; col < data[row].length; col++) {
					const value = data[row][col];
					if (value !== '') {
						const cellRef = getCellId(row, col);
						spreadsheetDataObj.data[cellRef] = value;

						if (value.startsWith('=')) {
							spreadsheetDataObj.formulas[cellRef] = value;
						}
					}
				}
			}

			if ($currentWorkbook) {
				// Update existing spreadsheet
				await updateSpreadsheetCell($currentWorkbook.id, {
					cell_reference: 'A1', // Just a placeholder for now
					value: spreadsheetDataObj,
					formula: ''
				});
			} else {
				// Create new spreadsheet
				const result = await createSpreadsheet({ name: $workbookName });
				currentWorkbook.set(result);
			}

			saveStatus.set('Saved');
			hasUnsavedChanges.set(false);

			setTimeout(() => {
				saveStatus.set('');
			}, 2000);

		} catch (error) {
			saveStatus.set('Save failed');
			console.error('Save error:', error);
		} finally {
			isSaving.set(false);
		}
	}

	// Export to Excel using backend
	async function exportToExcel() {
		if (!$currentWorkbook) {
			alert('Please save the spreadsheet first before exporting to Excel.');
			return;
		}

		try {
			const excelBlob = await exportSpreadsheetToExcel($currentWorkbook.id);
			const url = URL.createObjectURL(excelBlob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `${$workbookName || 'spreadsheet'}.xlsx`;
			a.click();
			URL.revokeObjectURL(url);
		} catch (error) {
			console.error('Export error:', error);
			alert('Failed to export to Excel: ' + formatApiError(error));
		}
	}

	// Share spreadsheet (placeholder for now)
	function shareSpreadsheet() {
		alert('Sharing functionality will be implemented soon!');
	}

	// Import from file
	function importFromFile(event) {
		const file = event.target.files[0];
		if (!file) return;

		// This would be handled by the parent component
		// For now, just trigger a custom event
		dispatch('import', { file });
	}

	// Helper function (should be imported from utils)
	function getCellId(row, col) {
		const cols = Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i));
		return `${cols[col]}${row + 1}`;
	}

	// Dispatch events to parent
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	// Menu close handler
	function closeAllMenus() {
		dispatch('closeMenus');
	}
</script>

<!-- Professional Toolbar -->
<div class="bg-white border-b border-gray-200 px-4 py-2">
	<div class="flex items-center justify-between">
		<!-- Left side - File operations -->
		<div class="flex items-center space-x-2">
			<input
				bind:value={$workbookName}
				placeholder="Workbook name..."
				maxlength="400"
				class="text-lg font-semibold bg-transparent border-none outline-none focus:ring-2 focus:ring-blue-500 rounded px-2 py-1"
				on:input={() => hasUnsavedChanges.set(true)}
			/>
			{#if $hasUnsavedChanges}
				<span class="text-sm text-orange-600">•</span>
			{/if}
			<span class="text-sm text-gray-500">{$saveStatus}</span>
		</div>

		<!-- Right side - Tools and actions -->
		<div class="flex items-center space-x-2">
			<!-- Data Import/Export -->
			<div class="flex items-center space-x-1 mr-4">
				<label class="px-3 py-1 bg-green-600 text-white text-sm rounded hover:bg-green-700 cursor-pointer" title="Import Excel/CSV">
					📥 Import
					<input type="file" accept=".xlsx,.xls,.csv" class="hidden" on:change={importFromFile} />
				</label>
				<button class="px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700" on:click={exportToExcel} title="Export to Excel">
					📤 Excel
				</button>
			</div>

			<!-- Advanced Features (only show in advanced mode) -->
			{#if $mode === 'advanced'}
				<div class="flex items-center space-x-1 mr-4 border-l border-gray-300 pl-4">
					<button class="px-3 py-1 bg-purple-600 text-white text-sm rounded hover:bg-purple-700" title="Insert Chart">
						📈 Chart
					</button>
					<button class="px-3 py-1 bg-indigo-600 text-white text-sm rounded hover:bg-indigo-700" title="Data Analysis">
						📊 Analysis
					</button>
					<button
						class="px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700"
						on:click={() => showPivotTableBuilder.set(true)}
						title="Create Pivot Table"
					>
						📋 Pivot Table
					</button>
					<button
						class="px-3 py-1 bg-teal-600 text-white text-sm rounded hover:bg-teal-700"
						on:click={() => showDataPrepTools.set(true)}
						title="Data Preparation Tools"
					>
						🔧 Data Prep
					</button>
					<button
						class="px-3 py-1 bg-yellow-600 text-white text-sm rounded hover:bg-yellow-700"
						on:click={() => showFormulaHelp.set(true)}
						title="Formula Help"
					>
						🔢 Functions
					</button>
				</div>
			{/if}

			<!-- Save/Share -->
			<div class="flex items-center space-x-1">
				<button
					class="px-4 py-1 bg-green-600 text-white text-sm rounded hover:bg-green-700 disabled:opacity-50"
					disabled={$isSaving}
					on:click={saveSpreadsheet}
					title="Save Workbook"
				>
					{$isSaving ? '💾 Saving...' : '💾 Save'}
				</button>
				<button
					class="px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700"
					on:click={shareSpreadsheet}
					title="Share Workbook"
				>
					👥 Share
				</button>
			</div>
		</div>
	</div>
</div>

<style>
	input:focus {
		box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5);
	}
</style>
