// @ts-nocheck
/**
 * Spreadsheet History Manager - Undo/Redo functionality for spreadsheets
 * Tracks cell changes, row/column operations, and formatting changes
 */

export class SpreadsheetHistory {
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

		// Remove any future history if we're not at the end
		if (this.currentIndex < this.history.length - 1) {
			this.history = this.history.slice(0, this.currentIndex + 1);
		}

		// Add new state
		this.history.push({
			...state,
			timestamp: Date.now()
		});

		// Limit history size
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

/**
 * Create a debounced history push function
 * @param {SpreadsheetHistory} history - The history instance
 * @param {number} delay - Debounce delay in ms
 * @returns {Function}
 */
export function createDebouncedPush(history, delay = 300) {
	let timeoutId = null;
	let lastState = null;

	return (state) => {
		if (timeoutId) {
			clearTimeout(timeoutId);
		}

		// Check if state has actually changed
		const stateString = JSON.stringify(state.data);
		if (lastState === stateString) return;

		lastState = stateString;

		timeoutId = setTimeout(() => {
			history.push(state);
			timeoutId = null;
		}, delay);
	};
}

/**
 * Keyboard shortcut handler for undo/redo
 * @param {KeyboardEvent} event 
 * @param {Function} onUndo 
 * @param {Function} onRedo 
 * @returns {boolean} - Whether the event was handled
 */
export function handleUndoRedoKeyboard(event, onUndo, onRedo) {
	const isCtrlOrCmd = event.ctrlKey || event.metaKey;
	
	if (isCtrlOrCmd && event.key.toLowerCase() === 'z') {
		event.preventDefault();
		if (event.shiftKey) {
			onRedo();
		} else {
			onUndo();
		}
		return true;
	}

	if (isCtrlOrCmd && event.key.toLowerCase() === 'y') {
		event.preventDefault();
		onRedo();
		return true;
	}

	return false;
}

export default SpreadsheetHistory;
