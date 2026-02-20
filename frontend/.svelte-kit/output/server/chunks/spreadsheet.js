import { w as writable, h as derived, c as create_ssr_component, s as subscribe, f as createEventDispatcher, b as escape, d as each, a as add_attribute, o as onDestroy, v as validate_component } from "./calendar.js";
import { g as getCellId, a as getCellStyle, F as FormulaBar } from "./forms.js";
const spreadsheetData = writable(Array(100).fill().map(() => Array(26).fill("")));
const selectedCell = writable({ row: 0, col: 0 });
const editingCell = writable(null);
const selectedCells = writable(/* @__PURE__ */ new Set());
const isDragging = writable(false);
const cellFormats = writable(/* @__PURE__ */ new Map());
const workbookName = writable("Untitled Spreadsheet");
const hasUnsavedChanges = writable(false);
const formulaBar = writable("");
const showFormulaHelp = writable(false);
const showPivotTableBuilder = writable(false);
const showDataPrepTools = writable(false);
const showFormatDialog = writable(false);
const showRibbonCustomizer = writable(false);
const showFileMenu = writable(false);
const showEditMenu = writable(false);
const showFormatMenu = writable(false);
const showViewMenu = writable(false);
const showInsertMenu = writable(false);
const showToolsMenu = writable(false);
const activeRibbonTab = writable("home");
const ribbonTabs = writable([
  { id: "home", name: "Home", icon: "🏠", tools: [] },
  { id: "insert", name: "Insert", icon: "➕", tools: [] },
  { id: "formulas", name: "Formulas", icon: "🔢", tools: [] },
  { id: "data", name: "Data", icon: "📊", tools: [] },
  { id: "review", name: "Review", icon: "👁️", tools: [] }
]);
const frozenCols = writable(0);
const sortState = writable(/* @__PURE__ */ new Map());
const filterState = writable(/* @__PURE__ */ new Map());
const fillPreviewCells = writable(/* @__PURE__ */ new Set());
const showFillHandle = writable(false);
const showContextMenu = writable(false);
const contextMenuPosition = writable({ x: 0, y: 0 });
const contextMenuItems = writable([]);
const sheets = writable([
  { id: "sheet1", name: "Sheet1", data: Array(100).fill().map(() => Array(26).fill("")) }
]);
const activeSheetId = writable("sheet1");
const sheetCounter = writable(1);
const statusBarInfo = writable({
  selectedCount: 0,
  sum: 0,
  average: 0,
  min: null,
  max: null
});
const zoomLevel = writable(100);
const canUndo = writable(false);
const canRedo = writable(false);
const isFormulaRangeSelecting = writable(false);
const formulaRangeStart = writable(null);
const formulaRangeEnd = writable(null);
derived([], () => Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i)));
derived([], () => Array.from({ length: 100 }, (_, i) => i + 1));
const QuickAccessToolbar = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $hasUnsavedChanges, $$unsubscribe_hasUnsavedChanges;
  let $canUndo, $$unsubscribe_canUndo;
  let $canRedo, $$unsubscribe_canRedo;
  let $workbookName, $$unsubscribe_workbookName;
  $$unsubscribe_hasUnsavedChanges = subscribe(hasUnsavedChanges, (value) => $hasUnsavedChanges = value);
  $$unsubscribe_canUndo = subscribe(canUndo, (value) => $canUndo = value);
  $$unsubscribe_canRedo = subscribe(canRedo, (value) => $canRedo = value);
  $$unsubscribe_workbookName = subscribe(workbookName, (value) => $workbookName = value);
  createEventDispatcher();
  $$unsubscribe_hasUnsavedChanges();
  $$unsubscribe_canUndo();
  $$unsubscribe_canRedo();
  $$unsubscribe_workbookName();
  return ` <div class="bg-gray-900 text-white px-3 py-1.5 flex items-center space-x-1"> <div class="flex items-center space-x-2 mr-4" data-svelte-h="svelte-1rx1vhw"><span class="text-lg">📊</span> <span class="font-semibold text-sm">TPT Titan</span></div>  <div class="flex items-center space-x-1"> <button class="p-1.5 rounded hover:bg-gray-700 transition-colors" title="New Spreadsheet (Ctrl+N)" data-svelte-h="svelte-19tlspf">📄</button>  <button class="p-1.5 rounded hover:bg-gray-700 transition-colors" title="Open Spreadsheet (Ctrl+O)" data-svelte-h="svelte-1f4damo">📂</button>  <button class="p-1.5 rounded hover:bg-gray-700 transition-colors relative" title="Save (Ctrl+S)">💾
			${$hasUnsavedChanges ? `<span class="absolute -top-0.5 -right-0.5 w-2 h-2 bg-orange-500 rounded-full"></span>` : ``}</button> <div class="w-px h-5 bg-gray-600 mx-2"></div>  <button class="p-1.5 rounded hover:bg-gray-700 transition-colors disabled:opacity-30 disabled:cursor-not-allowed" ${!$canUndo ? "disabled" : ""} title="Undo (Ctrl+Z)">↶</button>  <button class="p-1.5 rounded hover:bg-gray-700 transition-colors disabled:opacity-30 disabled:cursor-not-allowed" ${!$canRedo ? "disabled" : ""} title="Redo (Ctrl+Y)">↷</button></div>  <div class="flex-1 text-center"><span class="text-sm text-gray-300 truncate max-w-md inline-block">${escape($workbookName)} ${$hasUnsavedChanges ? `<span class="text-orange-400" data-svelte-h="svelte-1hbje4n">●</span>` : ``}</span></div>  <div class="flex items-center space-x-1"><button class="p-1.5 rounded hover:bg-gray-700 transition-colors" title="Share" data-svelte-h="svelte-1gq6w0m">👥</button></div></div>`;
});
const SpreadsheetMenuBar = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $showFileMenu, $$unsubscribe_showFileMenu;
  let $showEditMenu, $$unsubscribe_showEditMenu;
  let $showFormatMenu, $$unsubscribe_showFormatMenu;
  let $showViewMenu, $$unsubscribe_showViewMenu;
  let $showInsertMenu, $$unsubscribe_showInsertMenu;
  let $showToolsMenu, $$unsubscribe_showToolsMenu;
  $$unsubscribe_showFileMenu = subscribe(showFileMenu, (value) => $showFileMenu = value);
  $$unsubscribe_showEditMenu = subscribe(showEditMenu, (value) => $showEditMenu = value);
  $$unsubscribe_showFormatMenu = subscribe(showFormatMenu, (value) => $showFormatMenu = value);
  $$unsubscribe_showViewMenu = subscribe(showViewMenu, (value) => $showViewMenu = value);
  $$unsubscribe_showInsertMenu = subscribe(showInsertMenu, (value) => $showInsertMenu = value);
  $$unsubscribe_showToolsMenu = subscribe(showToolsMenu, (value) => $showToolsMenu = value);
  createEventDispatcher();
  $$unsubscribe_showFileMenu();
  $$unsubscribe_showEditMenu();
  $$unsubscribe_showFormatMenu();
  $$unsubscribe_showViewMenu();
  $$unsubscribe_showInsertMenu();
  $$unsubscribe_showToolsMenu();
  return ` <div class="bg-gray-800 text-white px-4 py-1 flex items-center space-x-6 text-sm"> <div class="relative"><button class="${"hover:bg-gray-700 px-3 py-1 rounded transition-colors " + escape($showFileMenu ? "bg-gray-700" : "", true)}">File</button> ${$showFileMenu ? `<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50"><button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1072cf2">📄 New Spreadsheet</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1ayjopy">📂 Open...</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center justify-between" data-svelte-h="svelte-b6jo40"><span>💾 Save</span> <span class="text-xs text-gray-500">Ctrl+S</span></button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center justify-between" data-svelte-h="svelte-x2v39k"><span>💾 Save As...</span> <span class="text-xs text-gray-500">Ctrl+Shift+S</span></button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-8ecd84">📊 Export to CSV</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1wc31wk">📈 Export to Excel</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-k0bc2z">🖨️ Print...</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-ht6634">👥 Share...</button></div>` : ``}</div>  <div class="relative"><button class="${"hover:bg-gray-700 px-3 py-1 rounded transition-colors " + escape($showEditMenu ? "bg-gray-700" : "", true)}">Edit</button> ${$showEditMenu ? `<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50"><button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1o8w9vf">↶ Undo</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-av8l50">↷ Redo</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1vl4qff">📋 Copy</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-rrrtge">✂️ Cut</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-8ch926">📌 Paste</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-xaks1f">🔍 Find &amp; Replace...</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-11alknb">📋 Select All</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1gnm6sc">🧹 Clear Contents</button></div>` : ``}</div>  <div class="relative"><button class="${"hover:bg-gray-700 px-3 py-1 rounded transition-colors " + escape($showFormatMenu ? "bg-gray-700" : "", true)}">Format</button> ${$showFormatMenu ? `<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50"><button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-nkw5yx">🎨 Format Cells...</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-13s7nri"><strong>B</strong> Bold</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1ean2fv"><em>I</em> Italic</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-14t9dc5"><u>U</u> Underline</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1dyrf7">⬅️ Align Left</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-3g16kk">⬌ Center</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-zu65u7">➡️ Align Right</button></div>` : ``}</div>  <div class="relative"><button class="${"hover:bg-gray-700 px-3 py-1 rounded transition-colors " + escape($showViewMenu ? "bg-gray-700" : "", true)}">View</button> ${$showViewMenu ? `<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50"><button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1q1uigh">🔍 Zoom In</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-tbnqtu">🔎 Zoom Out</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1r1eexo">📏 100% Zoom</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1vhlgir">📊 Show Gridlines</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1r5grke">📋 Show Headers</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-dz74f">🧊 Freeze Panes</button></div>` : ``}</div>  <div class="relative"><button class="${"hover:bg-gray-700 px-3 py-1 rounded transition-colors " + escape($showInsertMenu ? "bg-gray-700" : "", true)}">Insert</button> ${$showInsertMenu ? `<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50"><button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1upm52i">➕ Insert Row Above</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1j2yjku">➕ Insert Row Below</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-x93hg2">➕ Insert Column Left</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-ii2z9e">➕ Insert Column Right</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1mufa7b">📋 Insert Pivot Table</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1v9wx94">📈 Insert Chart</button></div>` : ``}</div>  <div class="relative"><button class="${"hover:bg-gray-700 px-3 py-1 rounded transition-colors " + escape($showToolsMenu ? "bg-gray-700" : "", true)}">Tools</button> ${$showToolsMenu ? `<div class="absolute top-full left-0 mt-1 bg-white border border-gray-300 rounded shadow-lg py-1 min-w-48 z-50"><button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1exrujz">🔧 Data Preparation</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-cvk7kk">📊 Data Analysis</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-5dzlgz">🔢 Formula Help</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-1n67q3b">📋 Function Reference</button> <div class="border-t border-gray-200 my-1"></div> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-p56uz">🔍 Spelling Check</button> <button class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" data-svelte-h="svelte-qcx1a5">📝 AutoCorrect Options</button></div>` : ``}</div></div>`;
});
const css$1 = {
  code: ".group.svelte-1rf0shy:hover .group-hover\\:opacity-100.svelte-1rf0shy{opacity:1}",
  map: `{"version":3,"file":"SpreadsheetRibbon.svelte","sources":["SpreadsheetRibbon.svelte"],"sourcesContent":["<script>\\r\\n\\timport {\\r\\n\\t\\tactiveRibbonTab,\\r\\n\\t\\tribbonTabs,\\r\\n\\t\\tshowRibbonCustomizer,\\r\\n\\t\\tcanUndo,\\r\\n\\t\\tcanRedo\\r\\n\\t} from '../stores/spreadsheet-store.js';\\r\\n\\r\\n\\t// Default available tools\\r\\n\\r\\n\\tconst defaultAvailableTools = [\\r\\n\\t\\t// File operations\\r\\n\\t\\t{ id: 'save', name: 'Save', icon: '💾', category: 'file', shortcut: 'Ctrl+S' },\\r\\n\\t\\t{ id: 'open', name: 'Open', icon: '📂', category: 'file' },\\r\\n\\t\\t{ id: 'export', name: 'Export', icon: '📤', category: 'file' },\\r\\n\\r\\n\\t\\t// Editing\\r\\n\\t\\t{ id: 'undo', name: 'Undo', icon: '↶', category: 'edit', shortcut: 'Ctrl+Z' },\\r\\n\\t\\t{ id: 'redo', name: 'Redo', icon: '↷', category: 'edit', shortcut: 'Ctrl+Y' },\\r\\n\\t\\t{ id: 'copy', name: 'Copy', icon: '📋', category: 'edit', shortcut: 'Ctrl+C' },\\r\\n\\t\\t{ id: 'paste', name: 'Paste', icon: '📌', category: 'edit', shortcut: 'Ctrl+V' },\\r\\n\\t\\t{ id: 'cut', name: 'Cut', icon: '✂️', category: 'edit', shortcut: 'Ctrl+X' },\\r\\n\\r\\n\\t\\t// Formatting\\r\\n\\t\\t{ id: 'bold', name: 'Bold', icon: 'B', category: 'format' },\\r\\n\\t\\t{ id: 'italic', name: 'Italic', icon: 'I', category: 'format' },\\r\\n\\t\\t{ id: 'underline', name: 'Underline', icon: 'U', category: 'format' },\\r\\n\\t\\t{ id: 'align-left', name: 'Align Left', icon: '⬅️', category: 'format' },\\r\\n\\t\\t{ id: 'align-center', name: 'Align Center', icon: '⬌', category: 'format' },\\r\\n\\t\\t{ id: 'align-right', name: 'Align Right', icon: '➡️', category: 'format' },\\r\\n\\r\\n\\t\\t// Formulas\\r\\n\\t\\t{ id: 'sum', name: 'Sum', icon: '∑', category: 'formulas' },\\r\\n\\t\\t{ id: 'average', name: 'Average', icon: '📊', category: 'formulas' },\\r\\n\\t\\t{ id: 'count', name: 'Count', icon: '#', category: 'formulas' },\\r\\n\\t\\t{ id: 'min', name: 'Min', icon: '↓', category: 'formulas' },\\r\\n\\t\\t{ id: 'max', name: 'Max', icon: '↑', category: 'formulas' },\\r\\n\\r\\n\\t\\t// Insert\\r\\n\\t\\t{ id: 'insert-row', name: 'Insert Row', icon: '➕', category: 'insert' },\\r\\n\\t\\t{ id: 'insert-column', name: 'Insert Column', icon: '➕', category: 'insert' },\\r\\n\\t\\t{ id: 'insert-chart', name: 'Insert Chart', icon: '📈', category: 'insert' },\\r\\n\\t\\t{ id: 'insert-pivot', name: 'Insert Pivot', icon: '📋', category: 'insert' },\\r\\n\\r\\n\\t\\t// Data\\r\\n\\t\\t{ id: 'sort-asc', name: 'Sort Ascending', icon: '↑', category: 'data' },\\r\\n\\t\\t{ id: 'sort-desc', name: 'Sort Descending', icon: '↓', category: 'data' },\\r\\n\\t\\t{ id: 'filter', name: 'Filter', icon: '🔍', category: 'data' },\\r\\n\\t\\t{ id: 'data-cleanup', name: 'Data Cleanup', icon: '🧹', category: 'data' }\\r\\n\\t];\\r\\n\\r\\n\\tlet availableTools = defaultAvailableTools;\\r\\n\\r\\n\\t// Dispatch events to parent\\r\\n\\r\\n\\timport { createEventDispatcher } from 'svelte';\\r\\n\\tconst dispatch = createEventDispatcher();\\r\\n\\r\\n\\r\\n\\tfunction dispatchAction(action, data = {}) {\\r\\n\\t\\tdispatch('action', { action, ...data });\\r\\n\\t}\\r\\n\\r\\n\\t// Handle tool click\\r\\n\\tfunction handleToolClick(tool) {\\r\\n\\t\\tswitch (tool.id) {\\r\\n\\t\\t\\tcase 'save':\\r\\n\\t\\t\\t\\tdispatchAction('save');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'copy':\\r\\n\\t\\t\\t\\tdispatchAction('copy');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'paste':\\r\\n\\t\\t\\t\\tdispatchAction('paste');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'cut':\\r\\n\\t\\t\\t\\tdispatchAction('cut');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'undo':\\r\\n\\t\\t\\t\\tdispatchAction('undo');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'redo':\\r\\n\\t\\t\\t\\tdispatchAction('redo');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'bold':\\r\\n\\t\\t\\t\\tdispatchAction('applyFormatting', { type: 'bold' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'italic':\\r\\n\\t\\t\\t\\tdispatchAction('applyFormatting', { type: 'italic' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'underline':\\r\\n\\t\\t\\t\\tdispatchAction('applyFormatting', { type: 'underline' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'align-left':\\r\\n\\t\\t\\t\\tdispatchAction('applyFormatting', { type: 'align', value: 'left' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'align-center':\\r\\n\\t\\t\\t\\tdispatchAction('applyFormatting', { type: 'align', value: 'center' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'align-right':\\r\\n\\t\\t\\t\\tdispatchAction('applyFormatting', { type: 'align', value: 'right' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'sum':\\r\\n\\t\\t\\t\\tdispatchAction('insertFormula', { formula: '=SUM()' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'average':\\r\\n\\t\\t\\t\\tdispatchAction('insertFormula', { formula: '=AVERAGE()' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'count':\\r\\n\\t\\t\\t\\tdispatchAction('insertFormula', { formula: '=COUNT()' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'min':\\r\\n\\t\\t\\t\\tdispatchAction('insertFormula', { formula: '=MIN()' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'max':\\r\\n\\t\\t\\t\\tdispatchAction('insertFormula', { formula: '=MAX()' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'insert-row':\\r\\n\\t\\t\\t\\tdispatchAction('insertRowBelow');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'insert-column':\\r\\n\\t\\t\\t\\tdispatchAction('insertColumnRight');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'insert-pivot':\\r\\n\\t\\t\\t\\tdispatchAction('insertPivotTable');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'sort-asc':\\r\\n\\t\\t\\t\\tdispatchAction('sortColumn', { direction: 'asc' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'sort-desc':\\r\\n\\t\\t\\t\\tdispatchAction('sortColumn', { direction: 'desc' });\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'filter':\\r\\n\\t\\t\\t\\tdispatchAction('showFilterDialog');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tcase 'data-cleanup':\\r\\n\\t\\t\\t\\tdispatchAction('showDataPrepTools');\\r\\n\\t\\t\\t\\tbreak;\\r\\n\\t\\t\\tdefault:\\r\\n\\t\\t\\t\\tdispatchAction('toolAction', { toolId: tool.id, tool });\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\t// Get tool by ID\\r\\n\\tfunction getToolById(toolId) {\\r\\n\\t\\treturn availableTools.find(t => t.id === toolId);\\r\\n\\t}\\r\\n<\/script>\\r\\n\\r\\n<!-- Customizable Ribbon -->\\r\\n<div class=\\"bg-white border-b border-gray-200\\">\\r\\n\\t<!-- Ribbon Tabs -->\\r\\n\\t<div class=\\"flex border-b border-gray-200\\">\\r\\n\\t\\t{#each $ribbonTabs as tab, tabIndex}\\r\\n\\t\\t\\t<button\\r\\n\\t\\t\\t\\tclass=\\"px-4 py-2 text-sm font-medium border-r border-gray-200 hover:bg-gray-50 transition-colors flex items-center space-x-2 {activeRibbonTab === tab.id ? 'bg-blue-50 text-blue-700 border-b-2 border-blue-500' : 'text-gray-700'}\\"\\r\\n\\t\\t\\t\\ton:click={() => activeRibbonTab.set(tab.id)}\\r\\n\\t\\t\\t>\\r\\n\\t\\t\\t\\t<span>{tab.icon}</span>\\r\\n\\t\\t\\t\\t<span>{tab.name}</span>\\r\\n\\t\\t\\t</button>\\r\\n\\t\\t{/each}\\r\\n\\t\\t<div class=\\"flex-1\\"></div>\\r\\n\\t\\t<button\\r\\n\\t\\t\\tclass=\\"px-3 py-2 text-gray-600 hover:bg-gray-100 rounded transition-colors\\"\\r\\n\\t\\t\\ton:click={() => showRibbonCustomizer.set(true)}\\r\\n\\t\\t\\ttitle=\\"Customize Ribbon\\"\\r\\n\\t\\t>\\r\\n\\t\\t\\t⚙️\\r\\n\\t\\t</button>\\r\\n\\t</div>\\r\\n\\r\\n\\t<!-- Ribbon Content -->\\r\\n\\t<div class=\\"p-4 min-h-20\\">\\r\\n\\t\\t{#each $ribbonTabs as tab}\\r\\n\\t\\t\\t{#if $activeRibbonTab === tab.id}\\r\\n\\t\\t\\t\\t<div class=\\"flex flex-wrap gap-2\\">\\r\\n\\t\\t\\t\\t\\t{#each tab.tools as tool, toolIndex}\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"group relative\\">\\r\\n\\t\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\t\\tclass=\\"flex flex-col items-center justify-center w-12 h-12 p-2 text-sm border border-gray-300 rounded hover:bg-gray-50 hover:border-gray-400 transition-colors disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:bg-white\\"\\r\\n\\t\\t\\t\\t\\t\\t\\ttitle=\\"{tool.name}{tool.shortcut ? \` (\${tool.shortcut})\` : ''}\\"\\r\\n\\t\\t\\t\\t\\t\\t\\ton:click={() => handleToolClick(tool)}\\r\\n\\t\\t\\t\\t\\t\\t\\tdisabled={(tool.id === 'undo' && !$canUndo) || (tool.id === 'redo' && !$canRedo)}\\r\\n\\r\\n\\t\\t\\t\\t\\t\\t>\\r\\n\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-lg\\">{tool.icon}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"text-xs mt-1 truncate w-full\\">{tool.name.split(' ')[0]}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t\\t<!-- Remove button (shown on hover) -->\\r\\n\\t\\t\\t\\t\\t\\t\\t<button\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"absolute -top-1 -right-1 w-4 h-4 bg-red-500 text-white text-xs rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-600\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:click={() => dispatchAction('removeToolFromRibbon', { tabId: tab.id, toolIndex })}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ttitle=\\"Remove from ribbon\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t×\\r\\n\\t\\t\\t\\t\\t\\t\\t</button>\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t{/each}\\r\\n\\t\\t\\t\\t\\t<!-- Drop zone for new tools -->\\r\\n\\t\\t\\t\\t\\t<div\\r\\n\\t\\t\\t\\t\\t\\tclass=\\"w-12 h-12 border-2 border-dashed border-gray-300 rounded flex items-center justify-center text-gray-400 hover:border-gray-500 hover:text-gray-600 transition-colors cursor-pointer\\"\\r\\n\\t\\t\\t\\t\\t\\ttitle=\\"Drop tools here\\"\\r\\n\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t+\\r\\n\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t{/if}\\r\\n\\t\\t{/each}\\r\\n\\t</div>\\r\\n</div>\\r\\n\\r\\n<style>\\r\\n\\t/* Ensure proper spacing and alignment */\\r\\n\\t.group:hover .group-hover\\\\:opacity-100 {\\r\\n\\t\\topacity: 1;\\r\\n\\t}\\r\\n</style>\\r\\n"],"names":[],"mappings":"AAwNC,qBAAM,MAAM,CAAC,wCAA0B,CACtC,OAAO,CAAE,CACV"}`
};
const SpreadsheetRibbon = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $ribbonTabs, $$unsubscribe_ribbonTabs;
  let $activeRibbonTab, $$unsubscribe_activeRibbonTab;
  let $canUndo, $$unsubscribe_canUndo;
  let $canRedo, $$unsubscribe_canRedo;
  $$unsubscribe_ribbonTabs = subscribe(ribbonTabs, (value) => $ribbonTabs = value);
  $$unsubscribe_activeRibbonTab = subscribe(activeRibbonTab, (value) => $activeRibbonTab = value);
  $$unsubscribe_canUndo = subscribe(canUndo, (value) => $canUndo = value);
  $$unsubscribe_canRedo = subscribe(canRedo, (value) => $canRedo = value);
  createEventDispatcher();
  $$result.css.add(css$1);
  $$unsubscribe_ribbonTabs();
  $$unsubscribe_activeRibbonTab();
  $$unsubscribe_canUndo();
  $$unsubscribe_canRedo();
  return ` <div class="bg-white border-b border-gray-200"> <div class="flex border-b border-gray-200">${each($ribbonTabs, (tab, tabIndex) => {
    return `<button class="${"px-4 py-2 text-sm font-medium border-r border-gray-200 hover:bg-gray-50 transition-colors flex items-center space-x-2 " + escape(
      activeRibbonTab === tab.id ? "bg-blue-50 text-blue-700 border-b-2 border-blue-500" : "text-gray-700",
      true
    )}"><span>${escape(tab.icon)}</span> <span>${escape(tab.name)}</span> </button>`;
  })} <div class="flex-1"></div> <button class="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded transition-colors" title="Customize Ribbon" data-svelte-h="svelte-1lppj72">⚙️</button></div>  <div class="p-4 min-h-20">${each($ribbonTabs, (tab) => {
    return `${$activeRibbonTab === tab.id ? `<div class="flex flex-wrap gap-2">${each(tab.tools, (tool, toolIndex) => {
      return `<div class="group relative svelte-1rf0shy"><button class="flex flex-col items-center justify-center w-12 h-12 p-2 text-sm border border-gray-300 rounded hover:bg-gray-50 hover:border-gray-400 transition-colors disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:bg-white" title="${escape(tool.name, true) + escape(tool.shortcut ? ` (${tool.shortcut})` : "", true)}" ${tool.id === "undo" && !$canUndo || tool.id === "redo" && !$canRedo ? "disabled" : ""}><span class="text-lg">${escape(tool.icon)}</span> <span class="text-xs mt-1 truncate w-full">${escape(tool.name.split(" ")[0])}</span></button>  <button class="absolute -top-1 -right-1 w-4 h-4 bg-red-500 text-white text-xs rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-600 svelte-1rf0shy" title="Remove from ribbon" data-svelte-h="svelte-1wv2qzp">×</button> </div>`;
    })}  <div class="w-12 h-12 border-2 border-dashed border-gray-300 rounded flex items-center justify-center text-gray-400 hover:border-gray-500 hover:text-gray-600 transition-colors cursor-pointer" title="Drop tools here" data-svelte-h="svelte-tsyjhf">+</div> </div>` : ``}`;
  })}</div> </div>`;
});
const css = {
  code: "td.svelte-1aqqtc9{margin:0;padding:0}input.svelte-1aqqtc9:focus{outline:none;box-shadow:inset 0 0 0 2px rgba(59, 130, 246, 0.5)}",
  map: `{"version":3,"file":"SpreadsheetGrid.svelte","sources":["SpreadsheetGrid.svelte"],"sourcesContent":["<script>\\r\\n\\timport {\\r\\n\\t\\tspreadsheetData,\\r\\n\\t\\tselectedCell,\\r\\n\\t\\teditingCell,\\r\\n\\t\\tselectedCells,\\r\\n\\t\\tcellFormats,\\r\\n\\t\\tfrozenRows,\\r\\n\\t\\tfrozenCols,\\r\\n\\t\\tsortState,\\r\\n\\t\\tfilterState,\\r\\n\\t\\tisDragging,\\r\\n\\t\\tdragStartCell,\\r\\n\\t\\tisDraggingFillHandle,\\r\\n\\t\\tfillStartCell,\\r\\n\\t\\tfillEndCell,\\r\\n\\t\\tfillPreviewCells,\\r\\n\\t\\tshowFillHandle,\\r\\n\\t\\tshowContextMenu,\\r\\n\\t\\tcontextMenuPosition,\\r\\n\\t\\tcontextMenuTarget,\\r\\n\\t\\tcontextMenuItems,\\r\\n\\t\\tformulaBar,\\r\\n\\t\\tisFormulaRangeSelecting,\\r\\n\\t\\tformulaRangeStart,\\r\\n\\t\\tformulaRangeEnd\\r\\n\\t} from '../stores/spreadsheet-store.js';\\r\\n\\r\\n\\timport {\\r\\n\\t\\tgetCellId,\\r\\n\\t\\tgetCellStyle,\\r\\n\\t\\tdetectDataType\\r\\n\\t} from '../utils/spreadsheet-utils.js';\\r\\n\\r\\n\\t// Props\\r\\n\\texport const mode = 'simple'; // 'simple' or 'advanced'\\r\\n\\r\\n\\t// Column headers (A-Z)\\r\\n\\t$: columns = Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i));\\r\\n\\r\\n\\t// Row numbers (1-100)\\r\\n\\t$: rows = Array.from({ length: 100 }, (_, i) => i + 1);\\r\\n\\r\\n\\t// Computed display values\\r\\n\\tfunction getDisplayValue(row, col) {\\r\\n\\t\\tconst value = $spreadsheetData[row]?.[col] || '';\\r\\n\\t\\tif (value.startsWith('=')) {\\r\\n\\t\\t\\t// In a real implementation, this would evaluate the formula\\r\\n\\t\\t\\t// For now, just return the formula or a placeholder\\r\\n\\t\\t\\treturn value; // Simplified - would need formula evaluation\\r\\n\\t\\t}\\r\\n\\t\\treturn value;\\r\\n\\t}\\r\\n\\r\\n\\t// Format cell value based on formatting rules\\r\\n\\tfunction formatCellValue(value, row, col) {\\r\\n\\t\\tconst format = $cellFormats.get(getCellId(row, col));\\r\\n\\r\\n\\t\\tif (!format || !value) return value;\\r\\n\\r\\n\\t\\t// Handle currency formatting\\r\\n\\t\\tif (format.currency) {\\r\\n\\t\\t\\tconst numValue = parseFloat(value);\\r\\n\\t\\t\\tif (!isNaN(numValue)) {\\r\\n\\t\\t\\t\\treturn new Intl.NumberFormat('en-US', {\\r\\n\\t\\t\\t\\t\\tstyle: 'currency',\\r\\n\\t\\t\\t\\t\\tcurrency: format.currency === true ? 'USD' : format.currency,\\r\\n\\t\\t\\t\\t\\tminimumFractionDigits: format.currencyDecimals !== undefined ? format.currencyDecimals : 2\\r\\n\\t\\t\\t\\t}).format(numValue);\\r\\n\\t\\t\\t}\\r\\n\\t\\t}\\r\\n\\r\\n\\t\\t// Handle number formatting\\r\\n\\t\\tif (format.numberFormat) {\\r\\n\\t\\t\\tconst numValue = parseFloat(value);\\r\\n\\t\\t\\tif (!isNaN(numValue)) {\\r\\n\\t\\t\\t\\tswitch (format.numberFormat) {\\r\\n\\t\\t\\t\\t\\tcase 'percent':\\r\\n\\t\\t\\t\\t\\t\\treturn new Intl.NumberFormat('en-US', {\\r\\n\\t\\t\\t\\t\\t\\t\\tstyle: 'percent',\\r\\n\\t\\t\\t\\t\\t\\t\\tminimumFractionDigits: 1,\\r\\n\\t\\t\\t\\t\\t\\t\\tmaximumFractionDigits: 1\\r\\n\\t\\t\\t\\t\\t\\t}).format(numValue / 100);\\r\\n\\t\\t\\t\\t\\tcase 'decimal':\\r\\n\\t\\t\\t\\t\\t\\treturn new Intl.NumberFormat('en-US', {\\r\\n\\t\\t\\t\\t\\t\\t\\tminimumFractionDigits: format.decimals || 2,\\r\\n\\t\\t\\t\\t\\t\\t\\tmaximumFractionDigits: format.decimals || 2\\r\\n\\t\\t\\t\\t\\t\\t}).format(numValue);\\r\\n\\t\\t\\t\\t\\tcase 'scientific':\\r\\n\\t\\t\\t\\t\\t\\treturn numValue.toExponential(2);\\r\\n\\t\\t\\t\\t\\tdefault:\\r\\n\\t\\t\\t\\t\\t\\treturn value;\\r\\n\\t\\t\\t\\t}\\r\\n\\t\\t\\t}\\r\\n\\t\\t}\\r\\n\\r\\n\\t\\treturn value;\\r\\n\\t}\\r\\n\\r\\n\\t// Check if cell is in current selection\\r\\n\\tfunction isCellSelected(row, col) {\\r\\n\\t\\tconst cellId = getCellId(row, col);\\r\\n\\t\\treturn $selectedCells.has(cellId);\\r\\n\\t}\\r\\n\\r\\n\\t// Check if cell is in formula range selection\\r\\n\\tfunction isCellInFormulaRange(row, col) {\\r\\n\\t\\tif (!$isFormulaRangeSelecting || !$formulaRangeStart || !$formulaRangeEnd) return false;\\r\\n\\r\\n\\t\\tconst startRow = Math.min($formulaRangeStart.row, $formulaRangeEnd.row);\\r\\n\\t\\tconst endRow = Math.max($formulaRangeStart.row, $formulaRangeEnd.row);\\r\\n\\t\\tconst startCol = Math.min($formulaRangeStart.col, $formulaRangeEnd.col);\\r\\n\\t\\tconst endCol = Math.max($formulaRangeStart.col, $formulaRangeEnd.col);\\r\\n\\r\\n\\t\\treturn row >= startRow && row <= endRow && col >= startCol && col <= endCol;\\r\\n\\t}\\r\\n\\r\\n\\t// Check if row passes all filters\\r\\n\\tfunction rowPassesFilters(rowIndex) {\\r\\n\\t\\tfor (const [colIndex, allowedValues] of $filterState.entries()) {\\r\\n\\t\\t\\tconst cellValue = $spreadsheetData[rowIndex]?.[colIndex] || '';\\r\\n\\t\\t\\tif (!allowedValues.has(cellValue)) {\\r\\n\\t\\t\\t\\treturn false;\\r\\n\\t\\t\\t}\\r\\n\\t\\t}\\r\\n\\t\\treturn true;\\r\\n\\t}\\r\\n\\r\\n\\t// Dispatch events to parent\\r\\n\\timport { createEventDispatcher } from 'svelte';\\r\\n\\tconst dispatch = createEventDispatcher();\\r\\n\\r\\n\\tfunction dispatchAction(action, data = {}) {\\r\\n\\t\\tdispatch('action', { action, ...data });\\r\\n\\t}\\r\\n\\r\\n\\t// Cell interaction handlers\\r\\n\\tfunction handleCellClick(row, col, event) {\\r\\n\\t\\tevent.preventDefault();\\r\\n\\r\\n\\t\\t// Handle formula range selection mode\\r\\n\\t\\tif ($isFormulaRangeSelecting) {\\r\\n\\t\\t\\tdispatchAction('endFormulaRangeSelection', { row, col });\\r\\n\\t\\t\\treturn;\\r\\n\\t\\t}\\r\\n\\r\\n\\t\\t// Handle multiple cell selection with Ctrl/Cmd\\r\\n\\t\\tif (event.ctrlKey || event.metaKey) {\\r\\n\\t\\t\\tdispatchAction('toggleCellSelection', { row, col });\\r\\n\\t\\t\\treturn;\\r\\n\\t\\t}\\r\\n\\r\\n\\t\\t// Handle range selection with Shift+click\\r\\n\\t\\tif (event.shiftKey) {\\r\\n\\t\\t\\tdispatchAction('extendCellSelection', { row, col });\\r\\n\\t\\t\\treturn;\\r\\n\\t\\t}\\r\\n\\r\\n\\t\\t// Single click - just select the cell (don't edit)\\r\\n\\t\\tdispatchAction('selectCell', { row, col });\\r\\n\\t}\\r\\n\\r\\n\\tfunction handleCellDoubleClick(row, col, event) {\\r\\n\\t\\tevent.preventDefault();\\r\\n\\t\\tdispatchAction('startEditingCell', { row, col });\\r\\n\\t}\\r\\n\\r\\n\\tfunction handleCellMouseDown(row, col, event) {\\r\\n\\t\\t// Start potential drag selection\\r\\n\\t\\tif (!$isFormulaRangeSelecting) {\\r\\n\\t\\t\\tdispatchAction('startCellDrag', { row, col, clientX: event.clientX, clientY: event.clientY });\\r\\n\\t\\t} else if ($isFormulaRangeSelecting) {\\r\\n\\t\\t\\tdispatchAction('startFormulaRangeSelection', { row, col });\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\tfunction handleCellMouseEnter(row, col, event) {\\r\\n\\t\\t// Handle drag selection\\r\\n\\t\\tif ($isDragging) {\\r\\n\\t\\t\\tdispatchAction('updateCellDrag', { row, col });\\r\\n\\t\\t} else if ($isFormulaRangeSelecting) {\\r\\n\\t\\t\\tdispatchAction('updateFormulaRangeSelection', { row, col });\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\tfunction handleCellMouseUp(row, col, event) {\\r\\n\\t\\tif ($isDragging) {\\r\\n\\t\\t\\tdispatchAction('endCellDrag', { row, col });\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\tfunction handleCellContextMenu(row, col, event) {\\r\\n\\t\\tevent.preventDefault();\\r\\n\\t\\tevent.stopPropagation();\\r\\n\\r\\n\\t\\tdispatchAction('showContextMenu', {\\r\\n\\t\\t\\trow,\\r\\n\\t\\t\\tcol,\\r\\n\\t\\t\\tx: Math.min(event.clientX, window.innerWidth - 200),\\r\\n\\t\\t\\ty: Math.min(event.clientY, window.innerHeight - 300)\\r\\n\\t\\t});\\r\\n\\t}\\r\\n\\r\\n\\t// Column header click for sorting\\r\\n\\tfunction handleColumnHeaderClick(colIndex, event) {\\r\\n\\t\\tif (event.shiftKey) {\\r\\n\\t\\t\\t// Could implement multi-column sorting\\r\\n\\t\\t\\treturn;\\r\\n\\t\\t}\\r\\n\\r\\n\\t\\tdispatchAction('toggleColumnSort', { colIndex });\\r\\n\\t}\\r\\n\\r\\n\\t// Fill handle interaction\\r\\n\\tfunction handleFillHandleMouseDown(event) {\\r\\n\\t\\tevent.preventDefault();\\r\\n\\t\\tevent.stopPropagation();\\r\\n\\r\\n\\t\\tdispatchAction('startFillHandleDrag', {\\r\\n\\t\\t\\trow: $selectedCell.row,\\r\\n\\t\\t\\tcol: $selectedCell.col\\r\\n\\t\\t});\\r\\n\\t}\\r\\n\\r\\n\\t// Svelte action for focusing input\\r\\n\\tfunction focusInput(node) {\\r\\n\\t\\tnode.focus();\\r\\n\\t}\\r\\n<\/script>\\r\\n\\r\\n<!-- Spreadsheet grid with freeze panes support -->\\r\\n<div class=\\"flex-1 overflow-auto\\">\\r\\n\\t<table class=\\"border-collapse\\">\\r\\n\\t\\t<!-- Column headers -->\\r\\n\\t\\t<thead class=\\"sticky top-0 z-20\\">\\r\\n\\t\\t\\t<tr>\\r\\n\\t\\t\\t\\t<th class=\\"w-12 h-8 bg-blue-100 border border-gray-300 text-center text-xs font-medium text-gray-600 sticky left-0 z-30 {$frozenCols >= 1 ? 'bg-blue-200 border-blue-400' : ''}\\">\\r\\n\\t\\t\\t\\t\\t#\\r\\n\\t\\t\\t\\t</th>\\r\\n\\t\\t\\t\\t{#each columns as col, colIndex}\\r\\n\\t\\t\\t\\t\\t<th\\r\\n\\t\\t\\t\\t\\t\\tclass=\\"w-24 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600 {$frozenCols > colIndex ? 'sticky left-12 bg-blue-100 border-blue-400 z-25' : ''} {$frozenCols > colIndex + 1 ? 'bg-blue-200' : ''} relative group cursor-pointer hover:bg-gray-200\\"\\r\\n\\t\\t\\t\\t\\t\\ton:click={(event) => handleColumnHeaderClick(colIndex, event)}\\r\\n\\t\\t\\t\\t\\t\\ttitle=\\"Click to sort • Right-click for filter menu\\"\\r\\n\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t<div class=\\"flex items-center justify-center px-1\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t<span>{col}</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t<!-- Sort indicator -->\\r\\n\\t\\t\\t\\t\\t\\t\\t{#if $sortState.has(colIndex)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t<span class=\\"ml-1 text-blue-600\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{$sortState.get(colIndex).direction === 'asc' ? '↑' : '↓'}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t</span>\\r\\n\\t\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t\\t</div>\\r\\n\\t\\t\\t\\t\\t\\t<!-- Filter indicator -->\\r\\n\\t\\t\\t\\t\\t\\t{#if $filterState.has(colIndex)}\\r\\n\\t\\t\\t\\t\\t\\t\\t<div class=\\"absolute top-0 right-0 w-2 h-2 bg-blue-500 rounded-full\\"></div>\\r\\n\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t</th>\\r\\n\\t\\t\\t\\t{/each}\\r\\n\\t\\t\\t</tr>\\r\\n\\t\\t</thead>\\r\\n\\r\\n\\t\\t<!-- Data rows -->\\r\\n\\t\\t<tbody>\\r\\n\\t\\t\\t{#each rows as rowNum, rowIndex}\\r\\n\\t\\t\\t\\t{#if rowPassesFilters(rowIndex)}\\r\\n\\t\\t\\t\\t\\t<tr class=\\"hover:bg-gray-50\\">\\r\\n\\t\\t\\t\\t\\t\\t<!-- Row number -->\\r\\n\\t\\t\\t\\t\\t\\t<td class=\\"w-12 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600 sticky left-0 z-10\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t{rowNum}\\r\\n\\t\\t\\t\\t\\t\\t</td>\\r\\n\\r\\n\\t\\t\\t\\t\\t\\t<!-- Data cells -->\\r\\n\\t\\t\\t\\t\\t\\t{#each columns as col, colIndex}\\r\\n\\t\\t\\t\\t\\t\\t\\t<td\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"w-24 h-8 border border-gray-300 p-0 relative cursor-cell\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{$selectedCell?.row === rowIndex && $selectedCell?.col === colIndex && !$editingCell ? 'ring-2 ring-blue-500 ring-inset bg-blue-50' : ''}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{$editingCell?.row === rowIndex && $editingCell?.col === colIndex ? 'ring-2 ring-green-500 ring-inset' : ''}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{isCellSelected(rowIndex, colIndex) && !($selectedCell?.row === rowIndex && $selectedCell?.col === colIndex) ? 'bg-blue-100 ring-1 ring-blue-400' : ''}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{isCellInFormulaRange(rowIndex, colIndex) ? 'bg-green-100 ring-2 ring-green-500' : ''}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{$fillPreviewCells.has(getCellId(rowIndex, colIndex)) ? 'bg-yellow-200 ring-1 ring-yellow-400' : ''}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\thover:bg-gray-50 transition-colors\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\tstyle={getCellStyle($cellFormats.get(getCellId(rowIndex, colIndex)))}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:click={(event) => handleCellClick(rowIndex, colIndex, event)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:contextmenu={(event) => handleCellContextMenu(rowIndex, colIndex, event)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:dblclick={(event) => handleCellDoubleClick(rowIndex, colIndex, event)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:mousedown={(event) => handleCellMouseDown(rowIndex, colIndex, event)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:mouseenter={(event) => handleCellMouseEnter(rowIndex, colIndex, event)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\ton:mouseup={(event) => handleCellMouseUp(rowIndex, colIndex, event)}\\r\\n\\t\\t\\t\\t\\t\\t\\t>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t{#if $editingCell?.row === rowIndex && $editingCell?.col === colIndex}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<input\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\tbind:value={$formulaBar}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"w-full h-full px-2 py-1 text-sm border-none outline-none bg-white\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\ton:keydown={(event) => {\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\tif (event.key === 'Enter') {\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\tdispatchAction('stopEditingCell', { save: true });\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\tif (event.key === 'Escape') {\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\tdispatchAction('stopEditingCell', { save: false });\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t}}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\ton:blur={() => dispatchAction('stopEditingCell', { save: true })}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\tuse:focusInput\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t/>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t{:else}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<div class=\\"w-full h-full px-2 py-1 text-sm overflow-hidden whitespace-nowrap select-none\\">\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t{formatCellValue(getDisplayValue(rowIndex, colIndex), rowIndex, colIndex)}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t</div>\\r\\n\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t<!-- Fill Handle -->\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{#if $selectedCell?.row === rowIndex && $selectedCell?.col === colIndex && !$editingCell && $showFillHandle}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t<div\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\tclass=\\"absolute bottom-0 right-0 w-2 h-2 bg-blue-500 border border-white cursor-crosshair z-10\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\ton:mousedown={handleFillHandleMouseDown}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\ttitle=\\"Drag to auto-fill adjacent cells\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\trole=\\"button\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t\\ttabindex=\\"0\\"\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t\\t></div>\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t\\t\\t\\t\\t</td>\\r\\n\\t\\t\\t\\t\\t\\t{/each}\\r\\n\\t\\t\\t\\t\\t</tr>\\r\\n\\t\\t\\t\\t{/if}\\r\\n\\t\\t\\t{/each}\\r\\n\\t\\t</tbody>\\r\\n\\t</table>\\r\\n</div>\\r\\n\\r\\n<style>\\r\\n\\t/* Ensure table cells don't have unwanted padding/margin */\\r\\n\\ttd {\\r\\n\\t\\tmargin: 0;\\r\\n\\t\\tpadding: 0;\\r\\n\\t}\\r\\n\\r\\n\\t/* Custom focus styles for inputs */\\r\\n\\tinput:focus {\\r\\n\\t\\toutline: none;\\r\\n\\t\\tbox-shadow: inset 0 0 0 2px rgba(59, 130, 246, 0.5);\\r\\n\\t}\\r\\n</style>\\r\\n"],"names":[],"mappings":"AA6UC,iBAAG,CACF,MAAM,CAAE,CAAC,CACT,OAAO,CAAE,CACV,CAGA,oBAAK,MAAO,CACX,OAAO,CAAE,IAAI,CACb,UAAU,CAAE,KAAK,CAAC,CAAC,CAAC,CAAC,CAAC,CAAC,CAAC,GAAG,CAAC,KAAK,EAAE,CAAC,CAAC,GAAG,CAAC,CAAC,GAAG,CAAC,CAAC,GAAG,CACnD"}`
};
const SpreadsheetGrid = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let columns;
  let rows;
  let $selectedCell, $$unsubscribe_selectedCell;
  let $$unsubscribe_isDragging;
  let $isFormulaRangeSelecting, $$unsubscribe_isFormulaRangeSelecting;
  let $spreadsheetData, $$unsubscribe_spreadsheetData;
  let $filterState, $$unsubscribe_filterState;
  let $formulaRangeEnd, $$unsubscribe_formulaRangeEnd;
  let $formulaRangeStart, $$unsubscribe_formulaRangeStart;
  let $selectedCells, $$unsubscribe_selectedCells;
  let $cellFormats, $$unsubscribe_cellFormats;
  let $frozenCols, $$unsubscribe_frozenCols;
  let $sortState, $$unsubscribe_sortState;
  let $editingCell, $$unsubscribe_editingCell;
  let $fillPreviewCells, $$unsubscribe_fillPreviewCells;
  let $formulaBar, $$unsubscribe_formulaBar;
  let $showFillHandle, $$unsubscribe_showFillHandle;
  $$unsubscribe_selectedCell = subscribe(selectedCell, (value) => $selectedCell = value);
  $$unsubscribe_isDragging = subscribe(isDragging, (value) => value);
  $$unsubscribe_isFormulaRangeSelecting = subscribe(isFormulaRangeSelecting, (value) => $isFormulaRangeSelecting = value);
  $$unsubscribe_spreadsheetData = subscribe(spreadsheetData, (value) => $spreadsheetData = value);
  $$unsubscribe_filterState = subscribe(filterState, (value) => $filterState = value);
  $$unsubscribe_formulaRangeEnd = subscribe(formulaRangeEnd, (value) => $formulaRangeEnd = value);
  $$unsubscribe_formulaRangeStart = subscribe(formulaRangeStart, (value) => $formulaRangeStart = value);
  $$unsubscribe_selectedCells = subscribe(selectedCells, (value) => $selectedCells = value);
  $$unsubscribe_cellFormats = subscribe(cellFormats, (value) => $cellFormats = value);
  $$unsubscribe_frozenCols = subscribe(frozenCols, (value) => $frozenCols = value);
  $$unsubscribe_sortState = subscribe(sortState, (value) => $sortState = value);
  $$unsubscribe_editingCell = subscribe(editingCell, (value) => $editingCell = value);
  $$unsubscribe_fillPreviewCells = subscribe(fillPreviewCells, (value) => $fillPreviewCells = value);
  $$unsubscribe_formulaBar = subscribe(formulaBar, (value) => $formulaBar = value);
  $$unsubscribe_showFillHandle = subscribe(showFillHandle, (value) => $showFillHandle = value);
  const mode = "simple";
  function getDisplayValue(row, col) {
    const value = $spreadsheetData[row]?.[col] || "";
    if (value.startsWith("=")) {
      return value;
    }
    return value;
  }
  function formatCellValue(value, row, col) {
    const format = $cellFormats.get(getCellId(row, col));
    if (!format || !value) return value;
    if (format.currency) {
      const numValue = parseFloat(value);
      if (!isNaN(numValue)) {
        return new Intl.NumberFormat(
          "en-US",
          {
            style: "currency",
            currency: format.currency === true ? "USD" : format.currency,
            minimumFractionDigits: format.currencyDecimals !== void 0 ? format.currencyDecimals : 2
          }
        ).format(numValue);
      }
    }
    if (format.numberFormat) {
      const numValue = parseFloat(value);
      if (!isNaN(numValue)) {
        switch (format.numberFormat) {
          case "percent":
            return new Intl.NumberFormat(
              "en-US",
              {
                style: "percent",
                minimumFractionDigits: 1,
                maximumFractionDigits: 1
              }
            ).format(numValue / 100);
          case "decimal":
            return new Intl.NumberFormat(
              "en-US",
              {
                minimumFractionDigits: format.decimals || 2,
                maximumFractionDigits: format.decimals || 2
              }
            ).format(numValue);
          case "scientific":
            return numValue.toExponential(2);
          default:
            return value;
        }
      }
    }
    return value;
  }
  function isCellSelected(row, col) {
    const cellId = getCellId(row, col);
    return $selectedCells.has(cellId);
  }
  function isCellInFormulaRange(row, col) {
    if (!$isFormulaRangeSelecting || !$formulaRangeStart || !$formulaRangeEnd) return false;
    const startRow = Math.min($formulaRangeStart.row, $formulaRangeEnd.row);
    const endRow = Math.max($formulaRangeStart.row, $formulaRangeEnd.row);
    const startCol = Math.min($formulaRangeStart.col, $formulaRangeEnd.col);
    const endCol = Math.max($formulaRangeStart.col, $formulaRangeEnd.col);
    return row >= startRow && row <= endRow && col >= startCol && col <= endCol;
  }
  function rowPassesFilters(rowIndex) {
    for (const [colIndex, allowedValues] of $filterState.entries()) {
      const cellValue = $spreadsheetData[rowIndex]?.[colIndex] || "";
      if (!allowedValues.has(cellValue)) {
        return false;
      }
    }
    return true;
  }
  createEventDispatcher();
  if ($$props.mode === void 0 && $$bindings.mode && mode !== void 0) $$bindings.mode(mode);
  $$result.css.add(css);
  columns = Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i));
  rows = Array.from({ length: 100 }, (_, i) => i + 1);
  $$unsubscribe_selectedCell();
  $$unsubscribe_isDragging();
  $$unsubscribe_isFormulaRangeSelecting();
  $$unsubscribe_spreadsheetData();
  $$unsubscribe_filterState();
  $$unsubscribe_formulaRangeEnd();
  $$unsubscribe_formulaRangeStart();
  $$unsubscribe_selectedCells();
  $$unsubscribe_cellFormats();
  $$unsubscribe_frozenCols();
  $$unsubscribe_sortState();
  $$unsubscribe_editingCell();
  $$unsubscribe_fillPreviewCells();
  $$unsubscribe_formulaBar();
  $$unsubscribe_showFillHandle();
  return ` <div class="flex-1 overflow-auto"><table class="border-collapse"> <thead class="sticky top-0 z-20"><tr><th class="${"w-12 h-8 bg-blue-100 border border-gray-300 text-center text-xs font-medium text-gray-600 sticky left-0 z-30 " + escape($frozenCols >= 1 ? "bg-blue-200 border-blue-400" : "", true)}">#</th> ${each(columns, (col, colIndex) => {
    return `<th class="${"w-24 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600 " + escape(
      $frozenCols > colIndex ? "sticky left-12 bg-blue-100 border-blue-400 z-25" : "",
      true
    ) + " " + escape($frozenCols > colIndex + 1 ? "bg-blue-200" : "", true) + " relative group cursor-pointer hover:bg-gray-200"}" title="Click to sort • Right-click for filter menu"><div class="flex items-center justify-center px-1"><span>${escape(col)}</span>  ${$sortState.has(colIndex) ? `<span class="ml-1 text-blue-600">${escape($sortState.get(colIndex).direction === "asc" ? "↑" : "↓")} </span>` : ``}</div>  ${$filterState.has(colIndex) ? `<div class="absolute top-0 right-0 w-2 h-2 bg-blue-500 rounded-full"></div>` : ``} </th>`;
  })}</tr></thead>  <tbody>${each(rows, (rowNum, rowIndex) => {
    return `${rowPassesFilters(rowIndex) ? `<tr class="hover:bg-gray-50"> <td class="w-12 h-8 bg-gray-100 border border-gray-300 text-center text-xs font-medium text-gray-600 sticky left-0 z-10 svelte-1aqqtc9">${escape(rowNum)}</td>  ${each(columns, (col, colIndex) => {
      return `<td class="${"w-24 h-8 border border-gray-300 p-0 relative cursor-cell " + escape(
        $selectedCell?.row === rowIndex && $selectedCell?.col === colIndex && !$editingCell ? "ring-2 ring-blue-500 ring-inset bg-blue-50" : "",
        true
      ) + " " + escape(
        $editingCell?.row === rowIndex && $editingCell?.col === colIndex ? "ring-2 ring-green-500 ring-inset" : "",
        true
      ) + " " + escape(
        isCellSelected(rowIndex, colIndex) && !($selectedCell?.row === rowIndex && $selectedCell?.col === colIndex) ? "bg-blue-100 ring-1 ring-blue-400" : "",
        true
      ) + " " + escape(
        isCellInFormulaRange(rowIndex, colIndex) ? "bg-green-100 ring-2 ring-green-500" : "",
        true
      ) + " " + escape(
        $fillPreviewCells.has(getCellId(rowIndex, colIndex)) ? "bg-yellow-200 ring-1 ring-yellow-400" : "",
        true
      ) + " hover:bg-gray-50 transition-colors svelte-1aqqtc9"}"${add_attribute("style", getCellStyle($cellFormats.get(getCellId(rowIndex, colIndex))), 0)}>${$editingCell?.row === rowIndex && $editingCell?.col === colIndex ? `<input class="w-full h-full px-2 py-1 text-sm border-none outline-none bg-white svelte-1aqqtc9"${add_attribute("value", $formulaBar, 0)}>` : `<div class="w-full h-full px-2 py-1 text-sm overflow-hidden whitespace-nowrap select-none">${escape(formatCellValue(getDisplayValue(rowIndex, colIndex), rowIndex, colIndex))}</div>  ${$selectedCell?.row === rowIndex && $selectedCell?.col === colIndex && !$editingCell && $showFillHandle ? `<div class="absolute bottom-0 right-0 w-2 h-2 bg-blue-500 border border-white cursor-crosshair z-10" title="Drag to auto-fill adjacent cells" role="button" tabindex="0"></div>` : ``}`} </td>`;
    })} </tr>` : ``}`;
  })}</tbody></table> </div>`;
});
const SheetTabs = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $sheets, $$unsubscribe_sheets;
  let $activeSheetId, $$unsubscribe_activeSheetId;
  let $$unsubscribe_sheetCounter;
  $$unsubscribe_sheets = subscribe(sheets, (value) => $sheets = value);
  $$unsubscribe_activeSheetId = subscribe(activeSheetId, (value) => $activeSheetId = value);
  $$unsubscribe_sheetCounter = subscribe(sheetCounter, (value) => value);
  createEventDispatcher();
  let editingSheetId = null;
  let editingSheetName = "";
  $$unsubscribe_sheets();
  $$unsubscribe_activeSheetId();
  $$unsubscribe_sheetCounter();
  return ` <div class="bg-gray-100 border-t border-gray-300 flex items-center"> <div class="flex items-center overflow-x-auto">${each($sheets, (sheet) => {
    return `<div class="${"group relative flex items-center min-w-24 max-w-40 px-3 py-2 border-r border-gray-300 cursor-pointer transition-colors " + escape(
      $activeSheetId === sheet.id ? "bg-white border-t-2 border-t-blue-500 text-gray-900 font-medium" : "bg-gray-100 hover:bg-gray-200 text-gray-600",
      true
    )}">${editingSheetId === sheet.id ? `<input type="text" class="w-full px-1 py-0.5 text-sm border border-blue-500 rounded focus:outline-none" autofocus${add_attribute("value", editingSheetName, 0)}>` : `<span class="truncate text-sm flex-1">${escape(sheet.name)}</span>  <div class="flex items-center space-x-1 ml-2 opacity-0 group-hover:opacity-100 transition-opacity"><button class="p-0.5 rounded hover:bg-gray-300 text-gray-500" title="Rename Sheet" data-svelte-h="svelte-oo27c0">✏️</button> ${$sheets.length > 1 ? `<button class="p-0.5 rounded hover:bg-red-100 text-red-500" title="Delete Sheet" data-svelte-h="svelte-ueq4y3">×
							</button>` : ``} </div>`} </div>`;
  })}</div>  <button class="flex items-center justify-center w-8 h-8 mx-2 rounded hover:bg-gray-200 text-gray-600 transition-colors" title="Add New Sheet" data-svelte-h="svelte-15scmxs">+</button></div>`;
});
function formatNumber(num) {
  if (num === null || num === void 0) return "-";
  if (Number.isInteger(num)) return num.toString();
  return num.toFixed(2);
}
const SpreadsheetStatusBar = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $selectedCell, $$unsubscribe_selectedCell;
  let $statusBarInfo, $$unsubscribe_statusBarInfo;
  let $sheets, $$unsubscribe_sheets;
  let $activeSheetId, $$unsubscribe_activeSheetId;
  let $zoomLevel, $$unsubscribe_zoomLevel;
  $$unsubscribe_selectedCell = subscribe(selectedCell, (value) => $selectedCell = value);
  $$unsubscribe_statusBarInfo = subscribe(statusBarInfo, (value) => $statusBarInfo = value);
  $$unsubscribe_sheets = subscribe(sheets, (value) => $sheets = value);
  $$unsubscribe_activeSheetId = subscribe(activeSheetId, (value) => $activeSheetId = value);
  $$unsubscribe_zoomLevel = subscribe(zoomLevel, (value) => $zoomLevel = value);
  createEventDispatcher();
  $$unsubscribe_selectedCell();
  $$unsubscribe_statusBarInfo();
  $$unsubscribe_sheets();
  $$unsubscribe_activeSheetId();
  $$unsubscribe_zoomLevel();
  return ` <div class="bg-gray-100 border-t border-gray-300 px-3 py-1.5 flex items-center justify-between text-sm"> <div class="flex items-center space-x-4">${$selectedCell ? `<span class="text-gray-600 font-medium">${escape(getCellId($selectedCell.row, $selectedCell.col))}</span>` : ``} ${$statusBarInfo.selectedCount > 0 ? `<div class="flex items-center space-x-3 text-gray-600"><span>Count: <strong class="text-gray-800">${escape($statusBarInfo.selectedCount)}</strong></span> <span>Sum: <strong class="text-gray-800">${escape(formatNumber($statusBarInfo.sum))}</strong></span> <span>Avg: <strong class="text-gray-800">${escape(formatNumber($statusBarInfo.average))}</strong></span> ${$statusBarInfo.min !== null ? `<span>Min: <strong class="text-gray-800">${escape(formatNumber($statusBarInfo.min))}</strong></span>` : ``} ${$statusBarInfo.max !== null ? `<span>Max: <strong class="text-gray-800">${escape(formatNumber($statusBarInfo.max))}</strong></span>` : ``}</div>` : ``}</div>  <div class="flex items-center space-x-2"><span class="text-gray-500 text-xs">${escape($sheets.find((s) => s.id === $activeSheetId)?.name || "Sheet1")}</span></div>  <div class="flex items-center space-x-2"><button class="p-1 rounded hover:bg-gray-200 transition-colors text-gray-600" title="Zoom Out" data-svelte-h="svelte-m7l0fl">−</button> <div class="flex items-center space-x-1"><input type="number"${add_attribute("value", $zoomLevel, 0)} min="50" max="200" class="w-14 px-1 py-0.5 text-center text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"> <span class="text-gray-600" data-svelte-h="svelte-19z6mdy">%</span></div> <button class="p-1 rounded hover:bg-gray-200 transition-colors text-gray-600" title="Zoom In" data-svelte-h="svelte-5832v6">+</button> <button class="px-2 py-1 text-xs rounded hover:bg-gray-200 transition-colors text-gray-600 ml-2" title="Reset Zoom" data-svelte-h="svelte-19tne6a">100%</button></div></div>`;
});
const SpreadsheetModals = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $showContextMenu, $$unsubscribe_showContextMenu;
  let $showFormulaHelp, $$unsubscribe_showFormulaHelp;
  let $showPivotTableBuilder, $$unsubscribe_showPivotTableBuilder;
  let $showDataPrepTools, $$unsubscribe_showDataPrepTools;
  let $showFormatDialog, $$unsubscribe_showFormatDialog;
  let $contextMenuPosition, $$unsubscribe_contextMenuPosition;
  let $contextMenuItems, $$unsubscribe_contextMenuItems;
  let $showRibbonCustomizer, $$unsubscribe_showRibbonCustomizer;
  $$unsubscribe_showContextMenu = subscribe(showContextMenu, (value) => $showContextMenu = value);
  $$unsubscribe_showFormulaHelp = subscribe(showFormulaHelp, (value) => $showFormulaHelp = value);
  $$unsubscribe_showPivotTableBuilder = subscribe(showPivotTableBuilder, (value) => $showPivotTableBuilder = value);
  $$unsubscribe_showDataPrepTools = subscribe(showDataPrepTools, (value) => $showDataPrepTools = value);
  $$unsubscribe_showFormatDialog = subscribe(showFormatDialog, (value) => $showFormatDialog = value);
  $$unsubscribe_contextMenuPosition = subscribe(contextMenuPosition, (value) => $contextMenuPosition = value);
  $$unsubscribe_contextMenuItems = subscribe(contextMenuItems, (value) => $contextMenuItems = value);
  $$unsubscribe_showRibbonCustomizer = subscribe(showRibbonCustomizer, (value) => $showRibbonCustomizer = value);
  createEventDispatcher();
  $$unsubscribe_showContextMenu();
  $$unsubscribe_showFormulaHelp();
  $$unsubscribe_showPivotTableBuilder();
  $$unsubscribe_showDataPrepTools();
  $$unsubscribe_showFormatDialog();
  $$unsubscribe_contextMenuPosition();
  $$unsubscribe_contextMenuItems();
  $$unsubscribe_showRibbonCustomizer();
  return ` ${$showFormulaHelp ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-96 overflow-y-auto"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-1tpkm3p">🔢 Advanced Spreadsheet Functions</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-e6ffdl">×</button></div> <div class="space-y-6" data-svelte-h="svelte-1dzkg2m"> <div class="bg-blue-50 p-4 rounded-lg"><h4 class="text-lg font-medium text-blue-900 mb-2">⚡ Powerful Mathematical Engine</h4> <p class="text-blue-800">TPT Titan includes a comprehensive mathematical function library with 50+ functions,
						advanced Excel compatibility, real-time collaboration, and automatic chart generation.</p></div>  <div class="grid grid-cols-1 md:grid-cols-2 gap-6"> <div><h4 class="text-lg font-medium text-gray-900 mb-3">🧮 Mathematical Functions</h4> <div class="space-y-2 text-sm"><div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">SUM(range)</code> <span class="text-gray-600">Add numbers</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">AVERAGE(range)</code> <span class="text-gray-600">Calculate mean</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">SIN(angle)</code> <span class="text-gray-600">Sine function</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">COS(angle)</code> <span class="text-gray-600">Cosine function</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">SQRT(number)</code> <span class="text-gray-600">Square root</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">POWER(base,exp)</code> <span class="text-gray-600">Power function</span></div></div></div>  <div><h4 class="text-lg font-medium text-gray-900 mb-3">📊 Statistical Functions</h4> <div class="space-y-2 text-sm"><div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">MIN(range)</code> <span class="text-gray-600">Minimum value</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">MAX(range)</code> <span class="text-gray-600">Maximum value</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">COUNT(range)</code> <span class="text-gray-600">Count numbers</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">STDEV(range)</code> <span class="text-gray-600">Standard deviation</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">MEDIAN(range)</code> <span class="text-gray-600">Median value</span></div> <div class="flex justify-between"><code class="bg-gray-100 px-2 py-1 rounded">MODE(range)</code> <span class="text-gray-600">Most frequent value</span></div></div></div></div>  <div><h4 class="text-lg font-medium text-gray-900 mb-3">🚀 Advanced Features</h4> <div class="grid grid-cols-1 md:grid-cols-3 gap-4"><div class="bg-purple-50 p-4 rounded"><h5 class="font-medium text-purple-900 mb-2">📈 Charts &amp; Visualization</h5> <p class="text-sm text-purple-700">Automatic chart suggestions based on your data patterns. Support for bar, line, pie, scatter, and area charts.</p></div> <div class="bg-green-50 p-4 rounded"><h5 class="font-medium text-green-900 mb-2">🔄 Real-time Collaboration</h5> <p class="text-sm text-green-700">Work with others simultaneously. See changes in real-time with conflict resolution and version control.</p></div> <div class="bg-blue-50 p-4 rounded"><h5 class="font-medium text-blue-900 mb-2">📊 Excel Compatibility</h5> <p class="text-sm text-blue-700">Full .xlsx import/export with formulas, styles, multiple sheets, and formatting preservation.</p></div></div></div>  <div><h4 class="text-lg font-medium text-gray-900 mb-3">💡 Quick Examples</h4> <div class="bg-gray-50 p-4 rounded"><div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm"><div><p class="font-medium mb-2">Basic Calculations:</p> <code class="block bg-white p-2 rounded mb-1">=SUM(A1:A10)</code> <code class="block bg-white p-2 rounded mb-1">=AVERAGE(B1:B20)</code> <code class="block bg-white p-2 rounded mb-1">=A1*B1 + C1</code></div> <div><p class="font-medium mb-2">Advanced Functions:</p> <code class="block bg-white p-2 rounded mb-1">=SIN(RADIANS(30))</code> <code class="block bg-white p-2 rounded mb-1">=SQRT(A1^2 + B1^2)</code> <code class="block bg-white p-2 rounded mb-1">=IF(A1&gt;100, &quot;High&quot;, &quot;Low&quot;)</code></div></div></div></div>  <div><h4 class="text-lg font-medium text-gray-900 mb-3">📤 Export &amp; Share</h4> <p class="text-gray-600 mb-3">Export your spreadsheets to various formats and share with collaborators:</p> <div class="grid grid-cols-2 md:grid-cols-4 gap-3"><div class="text-center p-3 bg-green-50 rounded"><div class="text-green-600 font-medium">Excel</div> <div class="text-xs text-green-600">.xlsx format</div></div> <div class="text-center p-3 bg-blue-50 rounded"><div class="text-blue-600 font-medium">CSV</div> <div class="text-xs text-blue-600">Data export</div></div> <div class="text-center p-3 bg-purple-50 rounded"><div class="text-purple-600 font-medium">PDF</div> <div class="text-xs text-purple-600">Printable format</div></div> <div class="text-center p-3 bg-orange-50 rounded"><div class="text-orange-600 font-medium">Share</div> <div class="text-xs text-orange-600">Collaborate</div></div></div></div></div> <div class="mt-6 flex justify-end"><button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-qu9n1s">Got it!</button></div></div></div>` : ``}  ${$showPivotTableBuilder ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-6xl max-h-[90vh] overflow-y-auto"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-b1uhui">📋 TPT Titan Pivot Table Builder</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-17g0rab">×</button></div> <div class="grid grid-cols-1 lg:grid-cols-3 gap-6"> <div class="lg:col-span-1"><h4 class="text-lg font-medium text-gray-900 mb-4" data-svelte-h="svelte-17mscfa">📊 Data Source</h4> <div class="space-y-4" data-svelte-h="svelte-ck2q6t"><div><label for="pivot-range" class="block text-sm font-medium text-gray-700 mb-2">Select Range</label> <input id="pivot-range" type="text" placeholder="Click and drag on spreadsheet to select range" value="A1:D10" class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500" readonly> <p class="text-xs text-gray-500 mt-1">💡 Open the pivot table builder and select cells visually</p></div> <div><label for="pivot-headers" class="block text-sm font-medium text-gray-700 mb-2">Data has headers</label> <input id="pivot-headers" type="checkbox" checked class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"></div></div> <h4 class="text-lg font-medium text-gray-900 mb-4 mt-6" data-svelte-h="svelte-1ez9v51">🎯 Pivot Configuration</h4> <div class="space-y-4"><div><label for="pivot-rows" class="block text-sm font-medium text-gray-700 mb-2" data-svelte-h="svelte-12vfsjl">Rows</label> <select id="pivot-rows" class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"><option value="Select field..." data-svelte-h="svelte-7ywo6k">Select field...</option><option value="Product" data-svelte-h="svelte-12nhss3">Product</option><option value="Category" data-svelte-h="svelte-oa7kxc">Category</option><option value="Region" data-svelte-h="svelte-inu57w">Region</option></select></div> <div><label for="pivot-columns" class="block text-sm font-medium text-gray-700 mb-2" data-svelte-h="svelte-1hb1vzn">Columns</label> <select id="pivot-columns" class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"><option value="Select field..." data-svelte-h="svelte-7ywo6k">Select field...</option><option value="Month" data-svelte-h="svelte-1qv0ks4">Month</option><option value="Quarter" data-svelte-h="svelte-xeyhs2">Quarter</option><option value="Year" data-svelte-h="svelte-2cuc65">Year</option></select></div> <div><label for="pivot-values" class="block text-sm font-medium text-gray-700 mb-2" data-svelte-h="svelte-32oywn">Values</label> <select id="pivot-values" class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"><option value="Select field..." data-svelte-h="svelte-7ywo6k">Select field...</option><option value="Sales" data-svelte-h="svelte-18jwroq">Sales</option><option value="Revenue" data-svelte-h="svelte-1k22hwu">Revenue</option><option value="Quantity" data-svelte-h="svelte-1hdxa1t">Quantity</option></select></div></div></div>  <div class="lg:col-span-2" data-svelte-h="svelte-17keho3"><h4 class="text-lg font-medium text-gray-900 mb-4">👁️ Preview</h4> <div class="bg-gray-50 p-4 rounded-lg min-h-64"><table class="w-full text-sm"><thead><tr class="bg-gray-200"><th class="p-2 text-left">Product</th> <th class="p-2 text-left">Q1</th> <th class="p-2 text-left">Q2</th> <th class="p-2 text-left">Q3</th> <th class="p-2 text-left">Q4</th> <th class="p-2 text-left">Total</th></tr></thead> <tbody><tr><td class="p-2 border">Widget A</td> <td class="p-2 border">12,500</td> <td class="p-2 border">15,200</td> <td class="p-2 border">18,900</td> <td class="p-2 border">14,300</td> <td class="p-2 border font-bold">60,900</td></tr> <tr><td class="p-2 border">Widget B</td> <td class="p-2 border">8,900</td> <td class="p-2 border">11,200</td> <td class="p-2 border">9,800</td> <td class="p-2 border">12,400</td> <td class="p-2 border font-bold">42,300</td></tr> <tr class="bg-gray-200 font-bold"><td class="p-2 border">Total</td> <td class="p-2 border">21,400</td> <td class="p-2 border">26,400</td> <td class="p-2 border">28,700</td> <td class="p-2 border">26,700</td> <td class="p-2 border">103,200</td></tr></tbody></table></div> <div class="mt-4 flex justify-end space-x-3"><button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700">📊 Create Chart</button> <button class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700">✅ Insert Pivot Table</button></div></div></div>  <div class="mt-8" data-svelte-h="svelte-1w62v9d"><h4 class="text-lg font-medium text-gray-900 mb-4">🚀 Advanced Pivot Features</h4> <div class="grid grid-cols-1 md:grid-cols-3 gap-4"><div class="bg-blue-50 p-4 rounded"><h5 class="font-medium text-blue-900 mb-2">🎲 Calculated Fields</h5> <p class="text-sm text-blue-700">Add custom calculations like profit margins, growth rates, or percentage changes.</p></div> <div class="bg-green-50 p-4 rounded"><h5 class="font-medium text-green-900 mb-2">🔍 Filters &amp; Slicers</h5> <p class="text-sm text-green-700">Interactive filters to drill down into your data and focus on specific segments.</p></div> <div class="bg-purple-50 p-4 rounded"><h5 class="font-medium text-purple-900 mb-2">📈 Multiple Aggregations</h5> <p class="text-sm text-purple-700">Sum, average, count, min, max, and custom aggregations for comprehensive analysis.</p></div></div></div></div></div>` : ``}  ${$showDataPrepTools ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-5xl max-h-[90vh] overflow-y-auto"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-lzcjhf">🔧 Data Preparation Tools</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-4z3hdy">×</button></div> <div class="grid grid-cols-1 lg:grid-cols-2 gap-6" data-svelte-h="svelte-1ra2bmr"> <div><h4 class="text-lg font-medium text-gray-900 mb-4">⚡ Quick Fixes</h4> <div class="space-y-3"><button class="w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 text-left transition-colors"><div class="flex items-center"><span class="text-2xl mr-3">🧹</span> <div><h5 class="font-medium">Remove Empty Rows/Columns</h5> <p class="text-sm text-gray-600">Clean up your data by removing blank entries</p></div></div></button> <button class="w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 text-left transition-colors"><div class="flex items-center"><span class="text-2xl mr-3">🔄</span> <div><h5 class="font-medium">Convert to Table Format</h5> <p class="text-sm text-gray-600">Transform your data into proper table structure</p></div></div></button> <button class="w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 text-left transition-colors"><div class="flex items-center"><span class="text-2xl mr-3">📋</span> <div><h5 class="font-medium">Normalize Data Types</h5> <p class="text-sm text-gray-600">Ensure consistent data types across columns</p></div></div></button> <button class="w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 text-left transition-colors"><div class="flex items-center"><span class="text-2xl mr-3">🔀</span> <div><h5 class="font-medium">Unpivot Data</h5> <p class="text-sm text-gray-600">Convert wide tables to long format for analysis</p></div></div></button></div></div>  <div><h4 class="text-lg font-medium text-gray-900 mb-4">🛠️ Advanced Tools</h4> <div class="space-y-3"><div class="p-4 border border-gray-200 rounded-lg"><h5 class="font-medium mb-2">📊 Data Profiling</h5> <p class="text-sm text-gray-600 mb-3">Analyze your data structure, identify patterns, and get insights about data quality.</p> <div class="grid grid-cols-3 gap-2 text-xs"><div class="text-center p-2 bg-blue-50 rounded"><div class="font-medium">15</div> <div class="text-gray-600">Columns</div></div> <div class="text-center p-2 bg-green-50 rounded"><div class="font-medium">1,247</div> <div class="text-gray-600">Rows</div></div> <div class="text-center p-2 bg-yellow-50 rounded"><div class="font-medium">98.5%</div> <div class="text-gray-600">Complete</div></div></div></div> <div class="p-4 border border-gray-200 rounded-lg"><h5 class="font-medium mb-2">🔄 Data Transformation</h5> <p class="text-sm text-gray-600 mb-3">Apply advanced transformations to prepare data for analysis.</p> <div class="space-y-2"><label class="flex items-center text-sm"><input type="checkbox" class="mr-2">
									Convert text to proper case</label> <label class="flex items-center text-sm"><input type="checkbox" class="mr-2">
									Remove duplicate rows</label> <label class="flex items-center text-sm"><input type="checkbox" class="mr-2">
									Standardize date formats</label></div></div> <div class="p-4 border border-gray-200 rounded-lg"><h5 class="font-medium mb-2">🎯 Pivot-Ready Conversion</h5> <p class="text-sm text-gray-600 mb-3">One-click conversion to pivot table friendly format.</p> <div class="space-y-2"><button class="w-full px-3 py-2 bg-blue-600 text-white text-sm rounded hover:bg-blue-700">Convert to Long Format</button> <button class="w-full px-3 py-2 bg-green-600 text-white text-sm rounded hover:bg-green-700">Create Fact-Dimension Tables</button> <button class="w-full px-3 py-2 bg-purple-600 text-white text-sm rounded hover:bg-purple-700">Generate Pivot Template</button></div></div></div></div></div>  <div class="mt-8" data-svelte-h="svelte-1meqwrv"><h4 class="text-lg font-medium text-gray-900 mb-4">👁️ Data Preview</h4> <div class="bg-gray-50 p-4 rounded-lg overflow-x-auto"><table class="w-full text-sm"><thead><tr class="bg-gray-200"><th class="p-2 text-left">Product</th> <th class="p-2 text-left">Category</th> <th class="p-2 text-left">Jan</th> <th class="p-2 text-left">Feb</th> <th class="p-2 text-left">Mar</th> <th class="p-2 text-left">Total</th></tr></thead> <tbody><tr><td class="p-2 border">Widget A</td> <td class="p-2 border">Electronics</td> <td class="p-2 border">1,200</td> <td class="p-2 border">1,500</td> <td class="p-2 border">1,800</td> <td class="p-2 border font-bold">4,500</td></tr> <tr><td class="p-2 border">Widget B</td> <td class="p-2 border">Electronics</td> <td class="p-2 border">900</td> <td class="p-2 border">1,100</td> <td class="p-2 border">1,200</td> <td class="p-2 border font-bold">3,200</td></tr></tbody></table></div></div> <div class="mt-6 flex justify-end space-x-3"><button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-9kew67">Cancel</button> <button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-1h1nlqj">Apply Changes</button></div></div></div>` : ``}  ${$showFormatDialog ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg p-6 w-full max-w-lg"><div class="flex items-center justify-between mb-6"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-12gquwn">🎨 Format Cells</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-1vc1wvn">×</button></div> <div class="space-y-6"> <div><h4 class="text-lg font-medium text-gray-900 mb-3" data-svelte-h="svelte-1jdi5dr">📊 Number</h4> <div class="grid grid-cols-2 gap-3"><button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-left" data-svelte-h="svelte-1h0r1up"><div class="font-medium">General</div> <div class="text-sm text-gray-500">No specific formatting</div></button> <button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-left" data-svelte-h="svelte-1akbs3i"><div class="font-medium">Currency</div> <div class="text-sm text-gray-500">$1,234.56</div></button> <button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-left" data-svelte-h="svelte-krd6kf"><div class="font-medium">Percentage</div> <div class="text-sm text-gray-500">123.46%</div></button> <button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-left" data-svelte-h="svelte-gqcn0q"><div class="font-medium">Number</div> <div class="text-sm text-gray-500">1,234.56</div></button></div></div>  <div><h4 class="text-lg font-medium text-gray-900 mb-3" data-svelte-h="svelte-1gb42nn">📐 Alignment</h4> <div class="grid grid-cols-3 gap-3"><button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center" data-svelte-h="svelte-gnaufw"><div class="text-2xl mb-2">⬅️</div> <div class="font-medium">Left</div></button> <button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center" data-svelte-h="svelte-ms7bbu"><div class="text-2xl mb-2">⬌</div> <div class="font-medium">Center</div></button> <button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center" data-svelte-h="svelte-zh8bge"><div class="text-2xl mb-2">➡️</div> <div class="font-medium">Right</div></button></div></div>  <div><h4 class="text-lg font-medium text-gray-900 mb-3" data-svelte-h="svelte-yfdfw">✏️ Font</h4> <div class="grid grid-cols-3 gap-3"><button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center" data-svelte-h="svelte-17076gd"><div class="text-2xl mb-2"><strong>B</strong></div> <div class="font-medium">Bold</div></button> <button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center" data-svelte-h="svelte-1p1wtga"><div class="text-2xl mb-2"><em>I</em></div> <div class="font-medium">Italic</div></button> <button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center" data-svelte-h="svelte-nrp6d8"><div class="text-2xl mb-2"><u>U</u></div> <div class="font-medium">Underline</div></button></div></div>  <div><h4 class="text-lg font-medium text-gray-900 mb-3" data-svelte-h="svelte-puozcj">🔲 Borders</h4> <div class="grid grid-cols-2 gap-3"><button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center" data-svelte-h="svelte-qzehuh"><div class="text-2xl mb-2">🔲</div> <div class="font-medium">All Borders</div></button> <button class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center" data-svelte-h="svelte-oy313o"><div class="text-2xl mb-2">▭</div> <div class="font-medium">Outline</div></button></div></div></div> <div class="mt-6 flex justify-end space-x-3"><button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-cg1v3s">Cancel</button> <button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-1qto3k0">Apply</button></div></div></div>` : ``}  ${$showContextMenu ? `<div class="fixed z-50" style="${"left: " + escape($contextMenuPosition.x, true) + "px; top: " + escape($contextMenuPosition.y, true) + "px;"}"><div class="bg-white border border-gray-300 rounded-lg shadow-lg py-1 min-w-48">${each($contextMenuItems, (item) => {
    return `${item.type === "divider" ? `<div class="border-t border-gray-200 my-1"></div>` : `<button class="${"w-full text-left px-4 py-2 text-sm hover:bg-gray-100 flex items-center justify-between " + escape(item.disabled ? "opacity-50 cursor-not-allowed" : "", true)}" ${item.disabled ? "disabled" : ""}><span class="flex items-center"><span class="mr-3 text-lg">${escape(item.icon)}</span> <span>${escape(item.label)}</span></span> ${item.shortcut ? `<span class="text-xs text-gray-500 ml-4">${escape(item.shortcut)}</span>` : ``} </button>`}`;
  })}</div></div>` : ``}  ${$showRibbonCustomizer ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"><div class="bg-white rounded-lg w-full max-w-6xl max-h-[90vh] overflow-hidden"><div class="flex items-center justify-between p-6 border-b border-gray-200"><h3 class="text-xl font-semibold text-gray-900" data-svelte-h="svelte-7axkfu">🎨 Customize Ribbon</h3> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-10elzsr">×</button></div> <div class="p-6"><div class="bg-blue-50 p-4 rounded-lg mb-4" data-svelte-h="svelte-13z73i4"><h5 class="font-medium text-blue-900 mb-2">📖 How to Customize</h5> <ul class="text-sm text-blue-800 space-y-1"><li>• <strong>Drag &amp; Drop:</strong> Drag tools from the palette to add them to the ribbon</li> <li>• <strong>Remove Tools:</strong> Hover over ribbon tools and click the × to remove them</li> <li>• <strong>Switch Tabs:</strong> Click on different ribbon tabs to customize each one</li> <li>• <strong>Reset:</strong> Use &quot;Reset to Default&quot; to restore the original ribbon layout</li> <li>• <strong>Save:</strong> Your customizations are automatically saved to your browser</li></ul></div> <div class="flex justify-end space-x-3"><button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" data-svelte-h="svelte-wxi2b8">Close</button> <button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-hhdl8p">Save &amp; Close</button></div></div></div></div>` : ``}`;
});
class SpreadsheetHistory {
  constructor(maxHistory = 50) {
    this.history = [];
    this.currentIndex = -1;
    this.maxHistory = maxHistory;
    this.isUndoing = false;
  }
  /**
   * Save a state to history
   * @param {Object} state - The state to save
   * @param {Array} state.data - 2D array of cell data
   * @param {Map} state.formats - Cell formatting map
   * @param {Object} state.metadata - Additional metadata (selected cell, etc.)
   * @param {string} state.actionType - Type of action ('cellChange', 'rowInsert', 'colDelete', etc.)
   */
  push(state) {
    if (this.isUndoing) return;
    if (this.currentIndex < this.history.length - 1) {
      this.history = this.history.slice(0, this.currentIndex + 1);
    }
    this.history.push({
      ...state,
      timestamp: Date.now()
    });
    if (this.history.length > this.maxHistory) {
      this.history.shift();
    } else {
      this.currentIndex++;
    }
  }
  /**
   * Check if undo is available
   * @returns {boolean}
   */
  canUndo() {
    return this.currentIndex > 0;
  }
  /**
   * Check if redo is available
   * @returns {boolean}
   */
  canRedo() {
    return this.currentIndex < this.history.length - 1;
  }
  /**
   * Undo - go back one state
   * @returns {Object|null} The previous state or null if can't undo
   */
  undo() {
    if (!this.canUndo()) return null;
    this.isUndoing = true;
    this.currentIndex--;
    const state = this.history[this.currentIndex];
    setTimeout(() => {
      this.isUndoing = false;
    }, 0);
    return state;
  }
  /**
   * Redo - go forward one state
   * @returns {Object|null} The next state or null if can't redo
   */
  redo() {
    if (!this.canRedo()) return null;
    this.isUndoing = true;
    this.currentIndex++;
    const state = this.history[this.currentIndex];
    setTimeout(() => {
      this.isUndoing = false;
    }, 0);
    return state;
  }
  /**
   * Get current state without changing position
   * @returns {Object|null}
   */
  getCurrentState() {
    if (this.currentIndex >= 0 && this.currentIndex < this.history.length) {
      return this.history[this.currentIndex];
    }
    return null;
  }
  /**
   * Clear all history
   */
  clear() {
    this.history = [];
    this.currentIndex = -1;
    this.isUndoing = false;
  }
  /**
   * Get history statistics
   * @returns {Object}
   */
  getStats() {
    return {
      totalStates: this.history.length,
      currentIndex: this.currentIndex,
      canUndo: this.canUndo(),
      canRedo: this.canRedo()
    };
  }
}
function handleUndoRedoKeyboard(event, onUndo, onRedo) {
  const isCtrlOrCmd = event.ctrlKey || event.metaKey;
  if (isCtrlOrCmd && event.key.toLowerCase() === "z") {
    event.preventDefault();
    if (event.shiftKey) {
      onRedo();
    } else {
      onUndo();
    }
    return true;
  }
  if (isCtrlOrCmd && event.key.toLowerCase() === "y") {
    event.preventDefault();
    onRedo();
    return true;
  }
  return false;
}
const Spreadsheet = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $spreadsheetData, $$unsubscribe_spreadsheetData;
  let $selectedCells, $$unsubscribe_selectedCells;
  let $cellFormats, $$unsubscribe_cellFormats;
  $$unsubscribe_spreadsheetData = subscribe(spreadsheetData, (value) => $spreadsheetData = value);
  $$unsubscribe_selectedCells = subscribe(selectedCells, (value) => $selectedCells = value);
  $$unsubscribe_cellFormats = subscribe(cellFormats, (value) => $cellFormats = value);
  let { mode = "simple" } = $$props;
  let { selectedTemplate = null } = $$props;
  createEventDispatcher();
  let spreadsheetHistory = new SpreadsheetHistory();
  let debouncedHistoryPush;
  function updateHistoryState() {
    canUndo.set(spreadsheetHistory.canUndo());
    canRedo.set(spreadsheetHistory.canRedo());
  }
  function saveToHistory(actionType = "cellChange") {
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
  function handleUndo() {
    const state = spreadsheetHistory.undo();
    if (state) {
      spreadsheetData.set(state.data);
      cellFormats.set(state.formats);
      updateHistoryState();
    }
  }
  function handleRedo() {
    const state = spreadsheetHistory.redo();
    if (state) {
      spreadsheetData.set(state.data);
      cellFormats.set(state.formats);
      updateHistoryState();
    }
  }
  function handleGlobalKeyDown(event) {
    if (handleUndoRedoKeyboard(event, handleUndo, handleRedo)) {
      return;
    }
  }
  onDestroy(() => {
    document.removeEventListener("keydown", handleGlobalKeyDown);
  });
  if ($$props.mode === void 0 && $$bindings.mode && mode !== void 0) $$bindings.mode(mode);
  if ($$props.selectedTemplate === void 0 && $$bindings.selectedTemplate && selectedTemplate !== void 0) $$bindings.selectedTemplate(selectedTemplate);
  {
    {
      const cells = Array.from($selectedCells);
      const data = $spreadsheetData;
      let sum = 0;
      let count = 0;
      let min = null;
      let max = null;
      let hasNumericValues = false;
      cells.forEach((cellId) => {
        const match = cellId.match(/([A-Z]+)(\d+)/);
        if (match) {
          const col = match[1].charCodeAt(0) - 65;
          const row = parseInt(match[2]) - 1;
          const value = data[row]?.[col];
          const numValue = parseFloat(value);
          if (!isNaN(numValue) && value !== "") {
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
  }
  {
    if (selectedTemplate && selectedTemplate.data) {
      console.log("Applying template:", selectedTemplate.name);
      console.log("Template data:", selectedTemplate.data);
      const templateData = selectedTemplate.data;
      const numRows = templateData.length;
      const numCols = Math.max(...templateData.map((row) => row.length));
      const newData = Array(Math.max(100, numRows)).fill().map(() => Array(Math.max(26, numCols)).fill(""));
      templateData.forEach((row, rowIndex) => {
        row.forEach((cellValue, colIndex) => {
          if (cellValue !== void 0 && cellValue !== null) {
            newData[rowIndex][colIndex] = String(cellValue);
          }
        });
      });
      spreadsheetData.set(newData);
      console.log("Template applied successfully, new data length:", newData.length);
      if (selectedTemplate.styles) {
        console.log("Applying template styles:", selectedTemplate.styles);
        const stylesMap = /* @__PURE__ */ new Map();
        Object.entries(selectedTemplate.styles).forEach(([range, style]) => {
          stylesMap.set(range, style);
        });
        cellFormats.set(stylesMap);
      }
      setTimeout(() => saveToHistory("templateApply"), 0);
    }
  }
  $$unsubscribe_spreadsheetData();
  $$unsubscribe_selectedCells();
  $$unsubscribe_cellFormats();
  return `<div class="flex flex-col h-full bg-white"> ${validate_component(QuickAccessToolbar, "QuickAccessToolbar").$$render($$result, {}, {}, {})}  ${validate_component(SpreadsheetMenuBar, "SpreadsheetMenuBar").$$render($$result, {}, {}, {})}  ${validate_component(SpreadsheetRibbon, "SpreadsheetRibbon").$$render($$result, {}, {}, {})}  ${validate_component(FormulaBar, "FormulaBar").$$render($$result, {}, {}, {})}  <div class="flex-1 overflow-hidden relative">${validate_component(SpreadsheetGrid, "SpreadsheetGrid").$$render($$result, {}, {}, {})}</div>  ${validate_component(SheetTabs, "SheetTabs").$$render($$result, {}, {}, {})}  ${validate_component(SpreadsheetStatusBar, "SpreadsheetStatusBar").$$render($$result, {}, {}, {})}</div>  ${validate_component(SpreadsheetModals, "SpreadsheetModals").$$render($$result, {}, {}, {})}`;
});
export {
  Spreadsheet as S,
  selectedCells as a,
  formulaBar as f,
  isFormulaRangeSelecting as i,
  selectedCell as s
};
