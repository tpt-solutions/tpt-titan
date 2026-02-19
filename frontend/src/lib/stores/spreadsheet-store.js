import { writable, derived } from 'svelte/store';

// Core spreadsheet state
export const spreadsheetData = writable(Array(100).fill().map(() => Array(26).fill('')));
export const selectedCell = writable({ row: 0, col: 0 });
export const editingCell = writable(null);
export const cellInput = writable('');

// Selection and interaction state
export const selectedCells = writable(new Set());
export const selectionStart = writable(null);
export const selectionEnd = writable(null);
export const isDragging = writable(false);
export const dragStartCell = writable(null);

// Cell formatting state
export const cellFormats = writable(new Map());

// Workbook and file state
export const currentWorkbook = writable(null);
export const workbookName = writable('Untitled Spreadsheet');
export const hasUnsavedChanges = writable(false);
export const isSaving = writable(false);
export const saveStatus = writable('');

// Formula and calculation state
export const formulaBar = writable('');
export const showFormulaHelp = writable(false);

// UI visibility state
export const showPivotTableBuilder = writable(false);
export const showDataPrepTools = writable(false);
export const showFindReplace = writable(false);
export const showFormatDialog = writable(false);
export const showRibbonCustomizer = writable(false);
export const showFreezePaneDialog = writable(false);

// Menu and toolbar state
export const showFileMenu = writable(false);
export const showEditMenu = writable(false);
export const showFormatMenu = writable(false);
export const showViewMenu = writable(false);
export const showInsertMenu = writable(false);
export const showToolsMenu = writable(false);
export const menuPosition = writable({ x: 0, y: 0 });

// Ribbon state
export const activeRibbonTab = writable('home');
export const ribbonTabs = writable([
	{ id: 'home', name: 'Home', icon: '🏠', tools: [] },
	{ id: 'insert', name: 'Insert', icon: '➕', tools: [] },
	{ id: 'formulas', name: 'Formulas', icon: '🔢', tools: [] },
	{ id: 'data', name: 'Data', icon: '📊', tools: [] },
	{ id: 'review', name: 'Review', icon: '👁️', tools: [] }
]);

// Freeze panes state
export const frozenRows = writable(0);
export const frozenCols = writable(0);

// Sorting and filtering state
export const sortState = writable(new Map());
export const filterState = writable(new Map());
export const showSortDialog = writable(false);
export const showFilterDialog = writable(false);

// Fill handle state
export const isDraggingFillHandle = writable(false);
export const fillStartCell = writable(null);
export const fillEndCell = writable(null);
export const fillPreviewCells = writable(new Set());
export const showFillHandle = writable(false);

// Context menu state
export const showContextMenu = writable(false);
export const contextMenuPosition = writable({ x: 0, y: 0 });
export const contextMenuTarget = writable(null);
export const contextMenuItems = writable([]);

// Clipboard and undo/redo
export const clipboardData = writable(null);
export const undoStack = writable([]);
export const redoStack = writable([]);
export const maxUndoSteps = writable(50);

// Template and mode state
export const mode = writable('simple'); // 'simple' or 'advanced'
export const selectedTemplate = writable(null);

// Sheet management state
export const sheets = writable([
	{ id: 'sheet1', name: 'Sheet1', data: Array(100).fill().map(() => Array(26).fill('')) }
]);
export const activeSheetId = writable('sheet1');
export const sheetCounter = writable(1); // For generating new sheet names

// Status bar state
export const statusBarInfo = writable({
	selectedCount: 0,
	sum: 0,
	average: 0,
	min: null,
	max: null
});

// Zoom state
export const zoomLevel = writable(100); // Percentage

// Undo/Redo availability state (for UI)
export const canUndo = writable(false);
export const canRedo = writable(false);

// Find & Replace state

export const findText = writable('');
export const replaceText = writable('');
export const findResults = writable([]);
export const currentFindIndex = writable(-1);

// Range selection state
export const selectedRange = writable('');
export const isSelectingRange = writable(false);
export const rangeStart = writable(null);
export const rangeEnd = writable(null);
export const isFormulaRangeSelecting = writable(false);
export const formulaRangeStart = writable(null);
export const formulaRangeEnd = writable(null);
export const formulaBeingEdited = writable('');
export const isInEditMode = writable(false);
export const pendingInput = writable('');

// Auto-save state
export const lastAutoSave = writable(Date.now());
export const autoSaveInterval = writable(null);
export const dailyBackupInterval = writable(null);

// Derived stores for computed values
export const columns = derived([], () => Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i)));
export const rows = derived([], () => Array.from({ length: 100 }, (_, i) => i + 1));

// Helper functions
export function getCellId(row, col) {
	const cols = Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i));
	return `${cols[col]}${row + 1}`;
}

export function getCellFromId(cellId) {
	const match = cellId.match(/([A-Z]+)(\d+)/);
	if (!match) return null;

	const col = match[1].charCodeAt(0) - 65;
	const row = parseInt(match[2]) - 1;
	return { row, col };
}

// Cell content getters/setters
export function getCellValue(row, col) {
	let data;
	spreadsheetData.subscribe(value => data = value)();
	return data?.[row]?.[col] || '';
}

export function setCellValue(row, col, value) {
	spreadsheetData.update(data => {
		if (!data[row]) data[row] = [];
		data[row][col] = value;
		return [...data];
	});
	hasUnsavedChanges.set(true);
}

export function getCellFormat(row, col) {
	let formats;
	cellFormats.subscribe(value => formats = value)();
	const cellId = getCellId(row, col);
	return formats.get(cellId) || {};
}

export function setCellFormat(row, col, format) {
	cellFormats.update(formats => {
		const cellId = getCellId(row, col);
		formats.set(cellId, format);
		return new Map(formats);
	});
	hasUnsavedChanges.set(true);
}

// Selection helpers
export function selectSingleCell(row, col) {
	selectedCell.set({ row, col });
	selectedCells.set(new Set([getCellId(row, col)]));
	const value = getCellValue(row, col);
	cellInput.set(value);
	formulaBar.set(value.startsWith('=') ? value : '');
}

export function selectCellRange(startRow, startCol, endRow, endCol) {
	selectedCells.update(cells => {
		cells.clear();
		const minRow = Math.min(startRow, endRow);
		const maxRow = Math.max(startRow, endRow);
		const minCol = Math.min(startCol, endCol);
		const maxCol = Math.max(startCol, endCol);

		for (let r = minRow; r <= maxRow; r++) {
			for (let c = minCol; c <= maxCol; c++) {
				cells.add(getCellId(r, c));
			}
		}
		return new Set(cells);
	});
}

export function clearSelection() {
	selectedCells.set(new Set());
}

// Undo/Redo helpers
export function addToUndoStack(action, data) {
	undoStack.update(stack => {
		stack.push({ action, data, timestamp: Date.now() });
		if (stack.length > 50) { // maxUndoSteps
			stack.shift();
		}
		return stack;
	});
	redoStack.set([]);
	hasUnsavedChanges.set(true);
}

// Reset all state (for new spreadsheet)
export function resetSpreadsheet() {
	spreadsheetData.set(Array(100).fill().map(() => Array(26).fill('')));
	selectedCell.set({ row: 0, col: 0 });
	editingCell.set(null);
	cellInput.set('');
	selectedCells.set(new Set());
	cellFormats.set(new Map());
	currentWorkbook.set(null);
	workbookName.set('Untitled Spreadsheet');
	hasUnsavedChanges.set(false);
	isSaving.set(false);
	saveStatus.set('');
	formulaBar.set('');
	sortState.set(new Map());
	filterState.set(new Map());
	undoStack.set([]);
	redoStack.set([]);
}
