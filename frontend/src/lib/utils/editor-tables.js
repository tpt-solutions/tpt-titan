// @ts-nocheck
/**
 * Editor Table Utility - Table creation and editing for text editor
 */

/**
 * Create a new table structure
 * @param {number} rows - Number of rows
 * @param {number} columns - Number of columns
 * @param {Array} headers - Optional header row
 * @returns {Object} Table structure
 */
export function createTable(rows = 3, columns = 3, headers = null) {
	const table = {
		rows: [],
		hasHeader: !!headers
	};
	
	// Create header row if provided
	if (headers && headers.length > 0) {
		table.rows.push({
			type: 'header',
			cells: headers.map(h => ({ content: h, align: 'left' }))
		});
	}
	
	// Create data rows
	for (let i = 0; i < rows; i++) {
		table.rows.push({
			type: 'data',
			cells: Array(columns).fill(null).map(() => ({
				content: '',
				align: 'left'
			}))
		});
	}
	
	return table;
}

/**
 * Insert row at position
 * @param {Object} table - Table structure
 * @param {number} index - Row index to insert at
 * @returns {Object} Updated table
 */
export function insertRow(table, index) {
	const columns = table.rows[0]?.cells.length || 3;
	const newRow = {
		type: 'data',
		cells: Array(columns).fill(null).map(() => ({
			content: '',
			align: 'left'
		}))
	};
	
	const newRows = [...table.rows];
	newRows.splice(index, 0, newRow);
	
	return { ...table, rows: newRows };
}

/**
 * Delete row at index
 * @param {Object} table - Table structure
 * @param {number} index - Row index to delete
 * @returns {Object} Updated table
 */
export function deleteRow(table, index) {
	if (table.rows.length <= 1) return table; // Keep at least one row
	
	const newRows = [...table.rows];
	newRows.splice(index, 1);
	
	return { ...table, rows: newRows };
}

/**
 * Insert column at position
 * @param {Object} table - Table structure
 * @param {number} index - Column index to insert at
 * @returns {Object} Updated table
 */
export function insertColumn(table, index) {
	const newRows = table.rows.map(row => {
		const newCells = [...row.cells];
		newCells.splice(index, 0, {
			content: '',
			align: 'left'
		});
		return { ...row, cells: newCells };
	});
	
	return { ...table, rows: newRows };
}

/**
 * Delete column at index
 * @param {Object} table - Table structure
 * @param {number} index - Column index to delete
 * @returns {Object} Updated table
 */
export function deleteColumn(table, index) {
	if (table.rows[0]?.cells.length <= 1) return table; // Keep at least one column
	
	const newRows = table.rows.map(row => {
		const newCells = [...row.cells];
		newCells.splice(index, 1);
		return { ...row, cells: newCells };
	});
	
	return { ...table, rows: newRows };
}

/**
 * Update cell content
 * @param {Object} table - Table structure
 * @param {number} rowIndex - Row index
 * @param {number} colIndex - Column index
 * @param {string} content - New content
 * @returns {Object} Updated table
 */
export function updateCell(table, rowIndex, colIndex, content) {
	const newRows = [...table.rows];
	newRows[rowIndex] = {
		...newRows[rowIndex],
		cells: [...newRows[rowIndex].cells]
	};
	newRows[rowIndex].cells[colIndex] = {
		...newRows[rowIndex].cells[colIndex],
		content
	};
	
	return { ...table, rows: newRows };
}

/**
 * Set cell alignment
 * @param {Object} table - Table structure
 * @param {number} rowIndex - Row index
 * @param {number} colIndex - Column index
 * @param {string} align - Alignment (left, center, right)
 * @returns {Object} Updated table
 */
export function setCellAlignment(table, rowIndex, colIndex, align) {
	const newRows = [...table.rows];
	newRows[rowIndex] = {
		...newRows[rowIndex],
		cells: [...newRows[rowIndex].cells]
	};
	newRows[rowIndex].cells[colIndex] = {
		...newRows[rowIndex].cells[colIndex],
		align
	};
	
	return { ...table, rows: newRows };
}

/**
 * Merge cells (simplified - only supports merging within a row)
 * @param {Object} table - Table structure
 * @param {number} rowIndex - Row index
 * @param {number} startCol - Start column
 * @param {number} endCol - End column
 * @returns {Object} Updated table
 */
export function mergeCells(table, rowIndex, startCol, endCol) {
	const newRows = [...table.rows];
	const row = newRows[rowIndex];
	
	// Mark cells as merged
	for (let i = startCol; i <= endCol; i++) {
		row.cells[i] = {
			...row.cells[i],
			colspan: i === startCol ? (endCol - startCol + 1) : 0,
			merged: i !== startCol
		};
	}
	
	return { ...table, rows: newRows };
}

/**
 * Convert table to HTML
 * @param {Object} table - Table structure
 * @returns {string} HTML string
 */
export function tableToHTML(table) {
	let html = '<table style="width: 100%; border-collapse: collapse; margin: 1em 0;">\n';
	
	table.rows.forEach((row, rowIndex) => {
		html += '  <tr>\n';
		
		row.cells.forEach((cell, colIndex) => {
			if (cell.merged) return; // Skip merged cells
			
			const tag = row.type === 'header' ? 'th' : 'td';
			const align = cell.align || 'left';
			const colspan = cell.colspan > 1 ? ` colspan="${cell.colspan}"` : '';
			
			const style = `
				border: 1px solid #ddd;
				padding: 8px;
				text-align: ${align};
				background: ${row.type === 'header' ? '#f5f5f5' : 'white'};
				font-weight: ${row.type === 'header' ? 'bold' : 'normal'};
			`.replace(/\s+/g, ' ').trim();
			
			html += `    <${tag}${colspan} style="${style}">${cell.content || '&nbsp;'}</${tag}>\n`;
		});
		
		html += '  </tr>\n';
	});
	
	html += '</table>';
	return html;
}

/**
 * Convert table to Markdown
 * @param {Object} table - Table structure
 * @returns {string} Markdown string
 */
export function tableToMarkdown(table) {
	if (table.rows.length === 0) return '';
	
	let md = '';
	
	table.rows.forEach((row, rowIndex) => {
		md += '|';
		
		row.cells.forEach(cell => {
			if (cell.merged) return;
			const content = (cell.content || '').replace(/\|/g, '\\|');
			md += ` ${content} |`;
		});
		
		md += '\n';
		
		// Add separator after header row
		if (rowIndex === 0 && row.type === 'header') {
			md += '|';
			row.cells.forEach(cell => {
				if (cell.merged) return;
				const colspan = cell.colspan || 1;
				const separator = ' --- '.repeat(colspan).trim();
				md += ` ${separator} |`;
			});
			md += '\n';
		}
	});
	
	return md;
}

/**
 * Insert table in rich text editor
 * @param {HTMLElement} editorElement - Editor element
 * @param {Object} table - Table structure
 */
export function insertTableInEditor(editorElement, table) {
	if (!editorElement) return;
	
	const html = tableToHTML(table);
	const wrapper = document.createElement('div');
	wrapper.innerHTML = html;
	
	// Insert at cursor position
	const selection = window.getSelection();
	if (selection.rangeCount > 0) {
		const range = selection.getRangeAt(0);
		range.deleteContents();
		range.insertNode(wrapper.firstChild);
		
		// Add space after
		const p = document.createElement('p');
		p.innerHTML = '<br>';
		range.setStartAfter(wrapper.firstChild);
		range.setEndAfter(wrapper.firstChild);
		range.insertNode(p);
	} else {
		editorElement.appendChild(wrapper.firstChild);
	}
}

/**
 * Insert table in block editor
 * @param {Array} blocks - Blocks array
 * @param {number} selectedIndex - Selected block index
 * @param {Object} table - Table structure
 * @returns {Array} Updated blocks
 */
export function insertTableInBlocks(blocks, selectedIndex, table) {
	const newBlock = {
		id: Date.now() + Math.random(),
		type: 'table',
		content: '',
		properties: {
			rows: table.rows.map(r => r.cells.map(c => c.content)),
			columns: table.rows[0]?.cells.map((_, i) => `Column ${i + 1}`) || [],
			hasHeader: table.hasHeader
		}
	};
	
	const newBlocks = [...blocks];
	newBlocks.splice(selectedIndex + 1, 0, newBlock);
	
	// Add empty text block after
	newBlocks.splice(selectedIndex + 2, 0, {
		id: Date.now() + Math.random() + 1,
		type: 'text',
		content: '',
		properties: {}
	});
	
	return newBlocks;
}

/**
 * Parse HTML table to structure
 * @param {HTMLTableElement} tableElement - Table DOM element
 * @returns {Object} Table structure
 */
export function parseHTMLTable(tableElement) {
	const rows = [];
	const trs = tableElement.querySelectorAll('tr');
	
	trs.forEach((tr, index) => {
		const cells = [];
		const tds = tr.querySelectorAll('td, th');
		
		tds.forEach(td => {
			cells.push({
				content: td.textContent || '',
				align: td.style.textAlign || 'left',
				colspan: parseInt(td.getAttribute('colspan')) || 1
			});
		});
		
		rows.push({
			type: tr.querySelector('th') ? 'header' : 'data',
			cells
		});
	});
	
	return { rows, hasHeader: rows[0]?.type === 'header' };
}

/**
 * Table styling presets
 */
export const tableStyles = {
	default: {
		border: '1px solid #ddd',
		headerBg: '#f5f5f5',
		cellPadding: '8px'
	},
	minimal: {
		border: 'none',
		headerBg: 'transparent',
		cellPadding: '4px',
		headerBorderBottom: '2px solid #333'
	},
	striped: {
		border: '1px solid #ddd',
		headerBg: '#333',
		headerColor: '#fff',
		cellPadding: '8px',
		stripedBg: '#f9f9f9'
	},
	bordered: {
		border: '2px solid #333',
		headerBg: '#f5f5f5',
		cellPadding: '10px',
		cellBorder: '1px solid #333'
	}
};

/**
 * Apply table style
 * @param {Object} table - Table structure
 * @param {string} styleName - Style name
 * @returns {Object} Table with style
 */
export function applyTableStyle(table, styleName) {
	const style = tableStyles[styleName] || tableStyles.default;
	return { ...table, style };
}

export default {
	createTable,
	insertRow,
	deleteRow,
	insertColumn,
	deleteColumn,
	updateCell,
	setCellAlignment,
	mergeCells,
	tableToHTML,
	tableToMarkdown,
	insertTableInEditor,
	insertTableInBlocks,
	parseHTMLTable,
	tableStyles,
	applyTableStyle
};
