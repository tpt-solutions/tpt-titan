// @ts-nocheck
/**
 * Spreadsheet Search Utilities - Find and replace functionality for spreadsheets
 */

/**
 * Search for text in spreadsheet data
 * @param {Array} data - 2D array of cell data
 * @param {string} searchText - Text to search for
 * @param {Object} options - Search options
 * @param {boolean} options.caseSensitive - Case sensitive search
 * @param {boolean} options.wholeCell - Match whole cell only
 * @returns {Array} Array of matches with {row, col, value}
 */
export function searchInSpreadsheet(data, searchText, options = {}) {
	const { caseSensitive = false, wholeCell = false } = options;
	const matches = [];
	
	if (!searchText) return matches;
	
	const searchRegex = new RegExp(
		wholeCell ? `^${escapeRegExp(searchText)}$` : escapeRegExp(searchText),
		caseSensitive ? 'g' : 'gi'
	);
	
	for (let row = 0; row < data.length; row++) {
		for (let col = 0; col < data[row].length; col++) {
			const cellValue = String(data[row][col] || '');
			if (searchRegex.test(cellValue)) {
				matches.push({
					row,
					col,
					value: cellValue,
					cellId: getCellId(row, col)
				});
			}
		}
	}
	
	return matches;
}

/**
 * Replace text in spreadsheet data
 * @param {Array} data - 2D array of cell data
 * @param {string} searchText - Text to search for
 * @param {string} replaceText - Text to replace with
 * @param {Object} options - Replace options
 * @param {boolean} options.caseSensitive - Case sensitive search
 * @param {boolean} options.wholeCell - Match whole cell only
 * @returns {Object} {newData, replacements} - New data array and count of replacements
 */
export function replaceInSpreadsheet(data, searchText, replaceText, options = {}) {
	const { caseSensitive = false, wholeCell = false } = options;
	let replacements = 0;
	
	if (!searchText) return { newData: data, replacements: 0 };
	
	const newData = data.map(row => [...row]);
	const searchRegex = new RegExp(
		wholeCell ? `^${escapeRegExp(searchText)}$` : escapeRegExp(searchText),
		caseSensitive ? 'g' : 'gi'
	);
	
	for (let row = 0; row < newData.length; row++) {
		for (let col = 0; col < newData[row].length; col++) {
			const cellValue = String(newData[row][col] || '');
			if (searchRegex.test(cellValue)) {
				newData[row][col] = cellValue.replace(searchRegex, replaceText);
				replacements++;
			}
		}
	}
	
	return { newData, replacements };
}

/**
 * Replace all occurrences in spreadsheet
 * @param {Array} data - 2D array of cell data
 * @param {string} searchText - Text to search for
 * @param {string} replaceText - Text to replace with
 * @param {Object} options - Replace options
 * @returns {Object} {newData, replacements}
 */
export function replaceAllInSpreadsheet(data, searchText, replaceText, options = {}) {
	return replaceInSpreadsheet(data, searchText, replaceText, options);
}

/**
 * Find next match from current position
 * @param {Array} matches - Array of match objects
 * @param {number} currentRow - Current row position
 * @param {number} currentCol - Current column position
 * @returns {Object|null} Next match or null
 */
export function findNextMatch(matches, currentRow, currentCol) {
	if (matches.length === 0) return null;
	
	// Find first match after current position
	for (const match of matches) {
		if (match.row > currentRow || (match.row === currentRow && match.col > currentCol)) {
			return match;
		}
	}
	
	// Wrap around to first match
	return matches[0];
}

/**
 * Find previous match from current position
 * @param {Array} matches - Array of match objects
 * @param {number} currentRow - Current row position
 * @param {number} currentCol - Current column position
 * @returns {Object|null} Previous match or null
 */
export function findPreviousMatch(matches, currentRow, currentCol) {
	if (matches.length === 0) return null;
	
	// Find last match before current position
	for (let i = matches.length - 1; i >= 0; i--) {
		const match = matches[i];
		if (match.row < currentRow || (match.row === currentRow && match.col < currentCol)) {
			return match;
		}
	}
	
	// Wrap around to last match
	return matches[matches.length - 1];
}

/**
 * Get cell ID from row and column (e.g., "A1", "B2")
 * @param {number} row - Row index (0-based)
 * @param {number} col - Column index (0-based)
 * @returns {string} Cell ID
 */
export function getCellId(row, col) {
	const colLetter = columnIndexToLetter(col);
	return `${colLetter}${row + 1}`;
}

/**
 * Convert column index to letter (0 = A, 1 = B, etc.)
 * @param {number} index - Column index
 * @returns {string} Column letter
 */
export function columnIndexToLetter(index) {
	let result = '';
	let num = index;
	
	do {
		result = String.fromCharCode(65 + (num % 26)) + result;
		num = Math.floor(num / 26) - 1;
	} while (num >= 0);
	
	return result;
}

/**
 * Escape special regex characters
 * @param {string} string - String to escape
 * @returns {string} Escaped string
 */
function escapeRegExp(string) {
	return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/**
 * Search in formulas only
 * @param {Array} data - 2D array of cell data
 * @param {string} searchText - Text to search for
 * @param {Object} options - Search options
 * @returns {Array} Array of formula matches
 */
export function searchInFormulas(data, searchText, options = {}) {
	const { caseSensitive = false } = options;
	const matches = [];
	
	if (!searchText) return matches;
	
	const searchRegex = new RegExp(
		escapeRegExp(searchText),
		caseSensitive ? 'g' : 'gi'
	);
	
	for (let row = 0; row < data.length; row++) {
		for (let col = 0; col < data[row].length; col++) {
			const cellValue = String(data[row][col] || '');
			// Check if cell starts with = (formula)
			if (cellValue.startsWith('=') && searchRegex.test(cellValue)) {
				matches.push({
					row,
					col,
					value: cellValue,
					cellId: getCellId(row, col),
					isFormula: true
				});
			}
		}
	}
	
	return matches;
}

export default {
	searchInSpreadsheet,
	replaceInSpreadsheet,
	replaceAllInSpreadsheet,
	findNextMatch,
	findPreviousMatch,
	getCellId,
	columnIndexToLetter,
	searchInFormulas
};
