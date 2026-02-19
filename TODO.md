# Spreadsheet Layout Improvement & Missing Functionality

## Overview
Consolidate the 3 horizontal menu bars into a cleaner 2-bar layout and add missing expected functionality.

## Implementation Steps

### Phase 1: Update Store (spreadsheet-store.js) ✅
- [x] Add state for multiple sheets management
- [x] Add state for status bar (selected cell count, sum, average)
- [x] Add state for zoom level
- [x] Add state for active sheet
- [x] Add canUndo/canRedo stores for UI state

### Phase 2: Create New Components ✅

#### 2.1 Create SpreadsheetStatusBar.svelte ✅
- [x] Display selected cell count, sum, average of selected cells
- [x] Add zoom controls (slider + percentage)
- [x] Add sheet navigation info
- [x] Display current cell reference

#### 2.2 Create SheetTabs.svelte ✅
- [x] Tab interface for multiple sheets
- [x] Add new sheet button
- [x] Sheet rename functionality
- [x] Sheet delete functionality
- [x] Active sheet highlighting

#### 2.3 Create QuickAccessToolbar.svelte ✅
- [x] Most common actions: Save, Undo, Redo
- [x] File operations: New, Open
- [x] Compact horizontal layout at top

### Phase 3: Enhance Existing Components ✅

#### 3.1 Update FormulaBar.svelte ✅
- [x] Add Name Box for direct cell navigation (e.g., type "A1" or "B2:D10")
- [x] Keep formula input field
- [x] Better styling integration

#### 3.2 Update SpreadsheetRibbon.svelte ✅
- [x] Remove canUndo/canRedo props, use stores directly
- [x] Keep existing tools organization

#### 3.3 Update SpreadsheetMenuBar.svelte
- [ ] Reduce to essential menus only (future enhancement)
- [ ] Remove redundant items now in ribbon (future enhancement)

### Phase 4: Update Main Spreadsheet.svelte ✅
- [x] Remove SpreadsheetToolbar import and usage
- [x] Add QuickAccessToolbar component
- [x] Add SheetTabs component
- [x] Add SpreadsheetStatusBar component
- [x] Update layout structure
- [x] Add status bar info calculation
- [x] Ensure all action handlers work with new components

### Phase 5: Update SpreadsheetGrid.svelte
- [ ] Add row/column resize handles (future enhancement)
- [x] Better cell reference handling for status bar

### Phase 6: Cleanup ✅
- [x] Remove SpreadsheetToolbar.svelte file
- [x] Update imports in Spreadsheet.svelte
- [x] Update SpreadsheetRibbon to use stores


## Files to Modify:
1. `frontend/src/lib/stores/spreadsheet-store.js` - Add new state
2. `frontend/src/lib/components/Spreadsheet.svelte` - Main layout
3. `frontend/src/lib/components/SpreadsheetRibbon.svelte` - Enhance
4. `frontend/src/lib/components/SpreadsheetMenuBar.svelte` - Reduce
5. `frontend/src/lib/components/FormulaBar.svelte` - Add Name Box
6. `frontend/src/lib/components/SpreadsheetGrid.svelte` - Add resize handles
7. Create `frontend/src/lib/components/QuickAccessToolbar.svelte` - New
8. Create `frontend/src/lib/components/SpreadsheetStatusBar.svelte` - New
9. Create `frontend/src/lib/components/SheetTabs.svelte` - New
10. Delete `frontend/src/lib/components/SpreadsheetToolbar.svelte` - Remove

## Expected Result:
- Cleaner 2-bar layout: Quick Access Toolbar + Ribbon
- Status bar at bottom with cell info and zoom
- Sheet tabs for multiple sheets
- Name box in formula bar for navigation
- Reduced code repetition
- Better user experience
