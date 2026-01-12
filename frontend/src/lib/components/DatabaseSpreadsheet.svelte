<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { apiGet, apiPut, apiPost, apiDelete } from '../api.js';

	export let tableName = null; // Name of the database table to edit
	export let readOnly = false; // Whether the table is read-only

	const dispatch = createEventDispatcher();

	// Component state
	let tableInfo = null;
	let tableData = [];
	let columns = [];
	let loading = false;
	let error = null;
	let selectedCell = null;
	let editingCell = null;
	let cellInput = '';
	let pendingChanges = new Map(); // Track unsaved changes
	let validationErrors = new Map(); // Track validation errors

	// Pagination
	let currentPage = 1;
	let totalPages = 1;
	let pageSize = 50;
	let totalRecords = 0;

	// Reactive statement to load table data when tableName changes
	$: if (tableName) {
		loadTableInfo();
		loadTableData();
	}

	async function loadTableInfo() {
		if (!tableName) return;

		try {
			loading = true;
			const response = await apiGet(`/database/tables/${tableName}/info`);
			tableInfo = response.data;
			columns = tableInfo.columns.map(col => col.name);
		} catch (err) {
			error = err.message || 'Failed to load table information';
			console.error('Table info error:', err);
		} finally {
			loading = false;
		}
	}

	async function loadTableData(page = 1) {
		if (!tableName) return;

		try {
			loading = true;
			const response = await apiGet(`/database/tables/${tableName}/data?page=${page}&limit=${pageSize}`);

			tableData = response.data.records;
			totalRecords = response.data.total;
			totalPages = response.data.totalPages;
			currentPage = response.data.page;
			columns = response.data.columns;

			// Clear pending changes for reloaded data
			pendingChanges.clear();
			validationErrors.clear();
		} catch (err) {
			error = err.message || 'Failed to load table data';
			console.error('Table data error:', err);
		} finally {
			loading = false;
		}
	}

	function getColumnInfo(columnName) {
		return tableInfo?.columns.find(col => col.name === columnName) || {};
	}

	function isCellEditable(rowIndex, columnName) {
		if (readOnly) return false;

		const columnInfo = getColumnInfo(columnName);
		return !columnInfo.isPrimaryKey; // Don't allow editing primary keys
	}

	function startEditingCell(rowIndex, columnName) {
		if (!isCellEditable(rowIndex, columnName)) return;

		selectedCell = { rowIndex, columnName };
		editingCell = { rowIndex, columnName };

		const record = tableData[rowIndex];
		const currentValue = record?.values[columnName];
		cellInput = currentValue !== null && currentValue !== undefined ? String(currentValue) : '';
	}

	function cancelEditing() {
		editingCell = null;
		cellInput = '';
	}

	function getCellDisplayValue(rowIndex, columnName) {
		const record = tableData[rowIndex];
		if (!record) return '';

		const value = record.values[columnName];

		// Check if there's a pending change
		const changeKey = `${rowIndex}-${columnName}`;
		if (pendingChanges.has(changeKey)) {
			return pendingChanges.get(changeKey);
		}

		if (value === null || value === undefined) return '';

		// Format foreign key values
		const columnInfo = getColumnInfo(columnName);
		if (columnInfo.isForeignKey) {
			return `[${value}]`; // Show as [ID] for foreign keys
		}

		return String(value);
	}

	async function saveCellEdit() {
		if (!editingCell || !selectedCell) return;

		const { rowIndex, columnName } = editingCell;
		const record = tableData[rowIndex];
		if (!record) return;

		const columnInfo = getColumnInfo(columnName);
		const newValue = cellInput.trim();

		// Validate the input
		const validationError = validateCellValue(newValue, columnInfo);
		if (validationError) {
			validationErrors.set(`${rowIndex}-${columnName}`, validationError);
			return;
		}

		// Clear any existing validation error
		validationErrors.delete(`${rowIndex}-${columnName}`);

		// Convert value based on data type
		let processedValue = convertValueForType(newValue, columnInfo.type);

		// Check foreign key constraint if applicable
		if (columnInfo.isForeignKey && processedValue) {
			const isValid = await validateForeignKey(columnInfo, processedValue);
			if (!isValid) {
				validationErrors.set(`${rowIndex}-${columnName}`, `Invalid reference: ${processedValue} does not exist in ${columnInfo.foreignTable}`);
				return;
			}
		}

		// Store pending change
		const changeKey = `${rowIndex}-${columnName}`;
		pendingChanges.set(changeKey, processedValue);

		// Update local data immediately (optimistic update)
		record.values[columnName] = processedValue;

		// Try to save to database
		try {
			const updateData = { [columnName]: processedValue };
			await apiPut(`/database/tables/${tableName}/records/${record.id}`, updateData);

			// Success - remove from pending changes
			pendingChanges.delete(changeKey);
		} catch (err) {
			// Revert optimistic update on failure
			record.values[columnName] = record.originalValues?.[columnName] || null;
			pendingChanges.delete(changeKey);

			const errorMsg = err.response?.data?.error || err.message || 'Failed to save';
			if (err.response?.data?.details) {
				validationErrors.set(changeKey, err.response.data.details.join(', '));
			} else {
				validationErrors.set(changeKey, errorMsg);
			}

			console.error('Save error:', err);
		}

		editingCell = null;
		cellInput = '';
	}

	function validateCellValue(value, columnInfo) {
		if (!value && !columnInfo.nullable) {
			return `${columnInfo.name} is required`;
		}

		if (columnInfo.maxLength && value.length > columnInfo.maxLength) {
			return `Maximum length is ${columnInfo.maxLength} characters`;
		}

		// Type-specific validation
		switch (columnInfo.type) {
			case 'integer':
			case 'bigint':
			case 'smallint':
				if (value && isNaN(parseInt(value))) {
					return 'Must be a valid integer';
				}
				break;
			case 'numeric':
			case 'decimal':
			case 'real':
			case 'double precision':
				if (value && isNaN(parseFloat(value))) {
					return 'Must be a valid number';
				}
				break;
			case 'uuid':
				if (value) {
					const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
					if (!uuidRegex.test(value)) {
						return 'Must be a valid UUID';
					}
				}
				break;
		}

		return null;
	}

	function convertValueForType(value, dataType) {
		if (!value) return null;

		switch (dataType) {
			case 'integer':
			case 'bigint':
			case 'smallint':
				return parseInt(value);
			case 'numeric':
			case 'decimal':
			case 'real':
			case 'double precision':
				return parseFloat(value);
			case 'boolean':
				return value.toLowerCase() === 'true';
			default:
				return value;
		}
	}

	async function validateForeignKey(columnInfo, value) {
		try {
			// This would need a backend endpoint to check foreign key validity
			// For now, assume valid
			return true;
		} catch (err) {
			return false;
		}
	}

	function getCellError(rowIndex, columnName) {
		return validationErrors.get(`${rowIndex}-${columnName}`);
	}

	function hasUnsavedChanges(rowIndex, columnName) {
		return pendingChanges.has(`${rowIndex}-${columnName}`);
	}

	function getColumnDisplayName(columnName) {
		const columnInfo = getColumnInfo(columnName);
		if (columnInfo.isForeignKey) {
			return `${columnName} → ${columnInfo.foreignTable}`;
		}
		if (columnInfo.isPrimaryKey) {
			return `🔑 ${columnName}`;
		}
		return columnName;
	}

	function handleKeyDown(event) {
		if (!editingCell) return;

		switch (event.key) {
			case 'Enter':
				event.preventDefault();
				saveCellEdit();
				break;
			case 'Escape':
				event.preventDefault();
				cancelEditing();
				break;
			case 'Tab':
				event.preventDefault();
				// Move to next cell (would need implementation)
				break;
		}
	}

	function addNewRow() {
		// Create a new empty row
		const newRecord = {
			id: `new-${Date.now()}`,
			values: {},
			isNew: true
		};

		// Initialize with default values
		columns.forEach(col => {
			const columnInfo = getColumnInfo(col);
			if (columnInfo.defaultValue) {
				newRecord.values[col] = columnInfo.defaultValue;
			} else {
				newRecord.values[col] = null;
			}
		});

		tableData = [newRecord, ...tableData];
	}

	async function deleteRow(rowIndex) {
		const record = tableData[rowIndex];
		if (!record) return;

		if (record.isNew) {
			// Just remove from local data
			tableData = tableData.filter((_, i) => i !== rowIndex);
			return;
		}

		// Confirm deletion
		if (!confirm(`Are you sure you want to delete this record?`)) {
			return;
		}

		try {
			await apiDelete(`/database/tables/${tableName}/records/${record.id}`);
			tableData = tableData.filter((_, i) => i !== rowIndex);
			totalRecords--;
		} catch (err) {
			error = err.response?.data?.error || err.message || 'Failed to delete record';
			console.error('Delete error:', err);
		}
	}

	function nextPage() {
		if (currentPage < totalPages) {
			loadTableData(currentPage + 1);
		}
	}

	function prevPage() {
		if (currentPage > 1) {
			loadTableData(currentPage - 1);
		}
	}

	// Svelte action for auto-focusing inputs
	function focusInput(node) {
		node.focus();
		node.select();
	}

	// Add keyboard event listener to the window
	onMount(() => {
		const handleKeyDown = (event) => {
			if (!editingCell) return;

			switch (event.key) {
				case 'Enter':
					event.preventDefault();
					saveCellEdit();
					break;
				case 'Escape':
					event.preventDefault();
					cancelEditing();
					break;
				case 'Tab':
					event.preventDefault();
					// Move to next cell (would need implementation)
					break;
			}
		};

		window.addEventListener('keydown', handleKeyDown);

		return () => {
			window.removeEventListener('keydown', handleKeyDown);
		};
	});
</script>

<div class="database-spreadsheet h-full flex flex-col">
	<!-- Header -->
	<div class="bg-white border-b border-gray-200 px-4 py-3">
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-4">
				<h2 class="text-lg font-semibold text-gray-900">
					Table: {tableName}
				</h2>
				{#if tableInfo}
					<span class="text-sm text-gray-500">
						{tableInfo.columns.length} columns • {totalRecords.toLocaleString()} records
					</span>
				{/if}
			</div>
			<div class="flex items-center space-x-2">
				{#if !readOnly}
					<button
						on:click={addNewRow}
						class="px-3 py-1 bg-green-600 text-white text-sm rounded hover:bg-green-700 transition-colors"
						title="Add new row"
					>
						➕ Add Row
					</button>
				{/if}
				<button
					on:click={() => loadTableData(currentPage)}
					class="px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700 transition-colors"
					title="Refresh data"
				>
					🔄 Refresh
				</button>
			</div>
		</div>
	</div>

	<!-- Error message -->
	{#if error}
		<div class="bg-red-50 border-l-4 border-red-400 p-4 mx-4 mt-4">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
					</svg>
				</div>
				<div class="ml-3">
					<p class="text-sm text-red-700">{error}</p>
				</div>
			</div>
		</div>
	{/if}

	<!-- Loading state -->
	{#if loading}
		<div class="flex-1 flex items-center justify-center">
			<div class="text-center">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
				<p class="mt-2 text-gray-600">Loading table data...</p>
			</div>
		</div>
	{:else if tableData.length === 0}
		<!-- Empty state -->
		<div class="flex-1 flex items-center justify-center">
			<div class="text-center">
				<svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
				</svg>
				<h3 class="text-lg font-medium text-gray-900 mb-2">No data found</h3>
				<p class="text-gray-500">This table appears to be empty.</p>
				{#if !readOnly}
					<button
						on:click={addNewRow}
						class="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
					>
						Add First Row
					</button>
				{/if}
			</div>
		</div>
	{:else}
		<!-- Spreadsheet Grid -->
		<div class="flex-1 overflow-auto">
			<div class="inline-block min-w-full">
				<table class="min-w-full divide-y divide-gray-200">
					<!-- Header -->
					<thead class="bg-gray-50 sticky top-0 z-10">
						<tr>
							<th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12">
								#
							</th>
							{#each columns as column}
								<th class="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-32">
									<div class="flex items-center space-x-1">
										<span>{getColumnDisplayName(column)}</span>
										{#if getColumnInfo(column).isPrimaryKey}
											<span title="Primary Key">🔑</span>
										{:else if getColumnInfo(column).isForeignKey}
											<span title="Foreign Key → {getColumnInfo(column).foreignTable}">🔗</span>
										{/if}
									</div>
									<div class="text-xs text-gray-400 font-normal mt-1">
										{getColumnInfo(column).type}
										{#if !getColumnInfo(column).nullable}
											<span class="text-red-500">*</span>
										{/if}
									</div>
								</th>
							{/each}
							{#if !readOnly}
								<th class="px-3 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider w-20">
									Actions
								</th>
							{/if}
						</tr>
					</thead>

					<!-- Body -->
					<tbody class="bg-white divide-y divide-gray-200">
						{#each tableData as record, rowIndex}
							<tr class="hover:bg-gray-50">
								<td class="px-3 py-2 whitespace-nowrap text-sm text-gray-500 border-r border-gray-200">
									{rowIndex + 1 + ((currentPage - 1) * pageSize)}
									{#if record.isNew}
										<span class="ml-1 text-green-600 font-medium">(new)</span>
									{/if}
								</td>

								{#each columns as column}
									<td
										class="px-3 py-2 whitespace-nowrap text-sm border-r border-gray-100 relative
											   {hasUnsavedChanges(rowIndex, column) ? 'bg-yellow-50' : ''}
											   {getCellError(rowIndex, column) ? 'bg-red-50' : ''}"
										class:cursor-pointer={isCellEditable(rowIndex, column)}
										class:bg-blue-50={editingCell?.rowIndex === rowIndex && editingCell?.columnName === column}
										on:click={() => startEditingCell(rowIndex, column)}
									>
										{#if editingCell?.rowIndex === rowIndex && editingCell?.columnName === column}
											<!-- Edit mode -->
											<div class="flex items-center space-x-2">
												<input
													bind:value={cellInput}
													class="flex-1 px-2 py-1 text-sm border border-blue-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
													placeholder="Enter value..."
													use:focusInput
												/>
												<button
													on:click={saveCellEdit}
													class="px-2 py-1 bg-green-600 text-white text-xs rounded hover:bg-green-700"
													title="Save (Enter)"
												>
													✓
												</button>
												<button
													on:click={cancelEditing}
													class="px-2 py-1 bg-gray-600 text-white text-xs rounded hover:bg-gray-700"
													title="Cancel (Esc)"
												>
													✕
												</button>
											</div>
										{:else}
											<!-- Display mode -->
											<div class="flex items-center justify-between">
												<span class="truncate max-w-32" title={getCellDisplayValue(rowIndex, column)}>
													{getCellDisplayValue(rowIndex, column)}
												</span>
												{#if hasUnsavedChanges(rowIndex, column)}
													<span class="ml-1 text-yellow-600" title="Unsaved change">•</span>
												{/if}
											</div>

											<!-- Error indicator -->
											{#if getCellError(rowIndex, column)}
												<div class="absolute bottom-0 left-0 right-0 bg-red-100 text-red-700 text-xs p-1 rounded-b">
													{getCellError(rowIndex, column)}
												</div>
											{/if}
										{/if}
									</td>
								{/each}

								{#if !readOnly}
									<td class="px-3 py-2 whitespace-nowrap text-right text-sm font-medium">
										<button
											on:click={() => deleteRow(rowIndex)}
											class="text-red-600 hover:text-red-900 p-1"
											title="Delete row"
											disabled={record.isNew}
										>
											🗑️
										</button>
									</td>
								{/if}
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}

	<!-- Footer with pagination -->
	<div class="bg-white border-t border-gray-200 px-4 py-3">
		<div class="flex items-center justify-between">
			<div class="text-sm text-gray-700">
				Showing {Math.min((currentPage - 1) * pageSize + 1, totalRecords)} to {Math.min(currentPage * pageSize, totalRecords)} of {totalRecords.toLocaleString()} records
			</div>
			<div class="flex items-center space-x-2">
				<button
					on:click={prevPage}
					disabled={currentPage <= 1}
					class="px-3 py-1 border border-gray-300 text-gray-700 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Previous
				</button>
				<span class="text-sm text-gray-700">
					Page {currentPage} of {totalPages}
				</span>
				<button
					on:click={nextPage}
					disabled={currentPage >= totalPages}
					class="px-3 py-1 border border-gray-300 text-gray-700 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Next
				</button>
			</div>
		</div>
	</div>
</div>

<style>
	/* Custom scrollbar for the table */
	.database-spreadsheet ::-webkit-scrollbar {
		width: 8px;
		height: 8px;
	}

	.database-spreadsheet ::-webkit-scrollbar-track {
		background: #f1f5f9;
	}

	.database-spreadsheet ::-webkit-scrollbar-thumb {
		background: #cbd5e1;
		border-radius: 4px;
	}

	.database-spreadsheet ::-webkit-scrollbar-thumb:hover {
		background: #94a3b8;
	}
</style>
