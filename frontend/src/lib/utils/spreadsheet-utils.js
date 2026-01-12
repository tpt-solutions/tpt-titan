// Spreadsheet utility functions

// Cell reference utilities
export function getCellId(row, col) {
	const cols = Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i));
	return `${cols[col]}${row + 1}`;
}

export function parseCellReference(ref) {
	const match = ref.match(/^([A-Z]+)(\d+)$/);
	if (!match) return null;

	const col = match[1].charCodeAt(0) - 65;
	const row = parseInt(match[2]) - 1;
	return { row, col };
}

// Range parsing and manipulation
export function parseRange(rangeStr) {
	if (!rangeStr.includes(':')) {
		return { start: parseCellReference(rangeStr), end: parseCellReference(rangeStr) };
	}

	const [startRef, endRef] = rangeStr.split(':');
	return {
		start: parseCellReference(startRef),
		end: parseCellReference(endRef)
	};
}

export function getCellsInRange(startRow, startCol, endRow, endCol) {
	const cells = [];
	const minRow = Math.min(startRow, endRow);
	const maxRow = Math.max(startRow, endRow);
	const minCol = Math.min(startCol, endCol);
	const maxCol = Math.max(startCol, endCol);

	for (let r = minRow; r <= maxRow; r++) {
		for (let c = minCol; c <= maxCol; c++) {
			cells.push({ row: r, col: c });
		}
	}
	return cells;
}

// Data type detection
export function detectDataType(value) {
	if (value === '' || value === null || value === undefined) {
		return 'empty';
	}

	const str = String(value);

	// Check if it's a formula
	if (str.startsWith('=')) {
		return 'formula';
	}

	// Check if it's a number
	if (!isNaN(str) && !isNaN(parseFloat(str))) {
		return Number(str) % 1 === 0 ? 'integer' : 'decimal';
	}

	// Check if it's a date (simple check)
	if (/^\d{1,2}[\/\-]\d{1,2}[\/\-]\d{2,4}$/.test(str)) {
		return 'date';
	}

	// Check if it's boolean
	if (str.toLowerCase() === 'true' || str.toLowerCase() === 'false') {
		return 'boolean';
	}

	return 'text';
}

// Cell formatting utilities
export function getCellStyle(format) {
	if (!format) return '';

	const styles = [];

	if (format.bold) styles.push('font-weight: bold');
	if (format.italic) styles.push('font-style: italic');
	if (format.underline) styles.push('text-decoration: underline');
	if (format.fontSize) styles.push(`font-size: ${format.fontSize}px`);
	if (format.color) styles.push(`color: ${format.color}`);
	if (format.backgroundColor) styles.push(`background-color: ${format.backgroundColor}`);
	if (format.align) styles.push(`text-align: ${format.align}`);

	// Borders
	if (format.borderTop) styles.push('border-top: 1px solid #ccc');
	if (format.borderBottom) styles.push('border-bottom: 1px solid #ccc');
	if (format.borderLeft) styles.push('border-left: 1px solid #ccc');
	if (format.borderRight) styles.push('border-right: 1px solid #ccc');
	if (format.borderAll) styles.push('border: 1px solid #ccc');

	return styles.join('; ');
}

// Clipboard utilities
export function copyToClipboard(data) {
	const text = data.map(row => row.join('\t')).join('\n');
	return navigator.clipboard.writeText(text);
}

export function parseClipboardData(clipboardText) {
	return clipboardText.split('\n').map(row => row.split('\t'));
}

// Formula utilities
export function extractFormulaDependencies(formula) {
	const dependencies = [];
	const cellRefRegex = /([A-Z]+\d+)/g;
	let match;

	while ((match = cellRefRegex.exec(formula)) !== null) {
		dependencies.push(match[1]);
	}

	return dependencies;
}

// Sorting utilities
export function sortData(data, columnIndex, direction = 'asc') {
	return [...data].sort((a, b) => {
		const aVal = a[columnIndex] || '';
		const bVal = b[columnIndex] || '';

		// Try numeric comparison first
		const aNum = parseFloat(aVal);
		const bNum = parseFloat(bVal);

		if (!isNaN(aNum) && !isNaN(bNum)) {
			return direction === 'asc' ? aNum - bNum : bNum - aNum;
		}

		// String comparison
		const comparison = aVal.toString().localeCompare(bVal.toString());
		return direction === 'asc' ? comparison : -comparison;
	});
}

// Filtering utilities
export function filterData(data, filters) {
	if (!filters || filters.size === 0) return data;

	return data.filter(row => {
		for (const [colIndex, allowedValues] of filters.entries()) {
			const cellValue = row[colIndex] || '';
			if (!allowedValues.has(cellValue)) {
				return false;
			}
		}
		return true;
	});
}

// Fill handle utilities
export function generateFillSequence(startValue, count, pattern = 'linear') {
	const result = [startValue];

	if (typeof startValue === 'number') {
		for (let i = 1; i < count; i++) {
			switch (pattern) {
				case 'linear':
					result.push(startValue + i);
					break;
				case 'copy':
					result.push(startValue);
					break;
				default:
					result.push(startValue + i);
			}
		}
	} else if (typeof startValue === 'string') {
		// Handle text patterns (simplified)
		for (let i = 1; i < count; i++) {
			result.push(startValue);
		}
	}

	return result;
}

// Validation utilities
export function validateFormula(formula) {
	if (!formula.startsWith('=')) return true;

	// Basic validation - check balanced parentheses
	let openCount = 0;
	for (const char of formula) {
		if (char === '(') openCount++;
		if (char === ')') openCount--;
		if (openCount < 0) return false;
	}
	return openCount === 0;
}

// Export utilities
export function dataToCSV(data, includeHeaders = false) {
	let csv = '';

	if (includeHeaders) {
		const headers = Array.from({ length: data[0]?.length || 0 }, (_, i) =>
			String.fromCharCode(65 + i)
		);
		csv += headers.join(',') + '\n';
	}

	for (const row of data) {
		const escapedRow = row.map(cell => {
			const str = String(cell || '');
			if (str.includes(',') || str.includes('"') || str.includes('\n')) {
				return '"' + str.replace(/"/g, '""') + '"';
			}
			return str;
		});
		csv += escapedRow.join(',') + '\n';
	}

	return csv;
}

// Template loading utilities
export function applyTemplate(data, template) {
	if (!template?.data) return data;

	const newData = Array(100).fill().map(() => Array(26).fill(''));

	template.data.forEach((row, rowIndex) => {
		if (Array.isArray(row) && rowIndex < 100) {
			row.forEach((cellValue, colIndex) => {
				if (colIndex < 26) {
					newData[rowIndex][colIndex] = cellValue || '';
				}
			});
		}
	});

	return newData;
}

// Keyboard navigation utilities
export const KEYBOARD_SHORTCUTS = {
	SAVE: ['ctrl+s', 'cmd+s'],
	UNDO: ['ctrl+z', 'cmd+z'],
	REDO: ['ctrl+y', 'cmd+y', 'ctrl+shift+z', 'cmd+shift+z'],
	COPY: ['ctrl+c', 'cmd+c'],
	PASTE: ['ctrl+v', 'cmd+v'],
	CUT: ['ctrl+x', 'cmd+x'],
	SELECT_ALL: ['ctrl+a', 'cmd+a'],
	FIND: ['ctrl+f', 'cmd+f'],
	BOLD: ['ctrl+b', 'cmd+b'],
	ITALIC: ['ctrl+i', 'cmd+i'],
};

export function normalizeKeyCombo(event) {
	const parts = [];
	if (event.ctrlKey || event.metaKey) parts.push('cmd');
	if (event.shiftKey) parts.push('shift');
	if (event.altKey) parts.push('alt');
	parts.push(event.key.toLowerCase());
	return parts.join('+');
}

// Auto-save utilities
export function shouldAutoSave(lastSave, intervalMinutes = 5) {
	const now = Date.now();
	return (now - lastSave) > (intervalMinutes * 60 * 1000);
}

// Performance utilities
export function debounce(func, wait) {
	let timeout;
	return function executedFunction(...args) {
		const later = () => {
			clearTimeout(timeout);
			func(...args);
		};
		clearTimeout(timeout);
		timeout = setTimeout(later, wait);
	};
}

export function throttle(func, limit) {
	let inThrottle;
	return function() {
		const args = arguments;
		const context = this;
		if (!inThrottle) {
			func.apply(context, args);
			inThrottle = true;
			setTimeout(() => inThrottle = false, limit);
		}
	};
}
