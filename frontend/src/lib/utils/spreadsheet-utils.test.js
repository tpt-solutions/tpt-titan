// frontend/src/lib/utils/spreadsheet-utils.test.js
// Run with:  cd frontend && npx vitest run src/lib/utils/spreadsheet-utils.test.js

import { describe, it, expect } from 'vitest';
import {
	getCellId,
	parseCellReference,
	parseRange,
	getCellsInRange,
	detectDataType,
	getCellStyle,
	extractFormulaDependencies,
	sortData,
	filterData,
	generateFillSequence,
	validateFormula,
	dataToCSV,
	shouldAutoSave,
	normalizeKeyCombo,
} from './spreadsheet-utils.js';

// ─── getCellId ──────────────────────────────────────────────────────────────
describe('getCellId', () => {
	it('converts first column correctly', () => {
		expect(getCellId(0, 0)).toBe('A1');
	});

	it('converts row correctly', () => {
		expect(getCellId(4, 0)).toBe('A5');
	});

	it('converts last single-letter column (Z)', () => {
		expect(getCellId(0, 25)).toBe('Z1');
	});

	it('converts first double-letter column (AA)', () => {
		expect(getCellId(0, 26)).toBe('AA1');
	});

	it('converts AB correctly', () => {
		expect(getCellId(0, 27)).toBe('AB1');
	});

	it('converts AZ correctly (col 51)', () => {
		expect(getCellId(0, 51)).toBe('AZ1');
	});

	it('converts BA correctly (col 52)', () => {
		expect(getCellId(0, 52)).toBe('BA1');
	});

	it('combines row and column', () => {
		expect(getCellId(9, 3)).toBe('D10');
	});
});

// ─── parseCellReference ──────────────────────────────────────────────────────
describe('parseCellReference', () => {
	it('parses A1', () => {
		expect(parseCellReference('A1')).toEqual({ row: 0, col: 0 });
	});

	it('parses Z10', () => {
		expect(parseCellReference('Z10')).toEqual({ row: 9, col: 25 });
	});

	it('returns null for invalid ref', () => {
		expect(parseCellReference('notaref')).toBeNull();
	});

	it('returns null for lowercase ref', () => {
		// lowercase not supported per the regex
		expect(parseCellReference('a1')).toBeNull();
	});
});

// ─── parseRange ──────────────────────────────────────────────────────────────
describe('parseRange', () => {
	it('parses a single cell range', () => {
		const r = parseRange('B3');
		expect(r.start).toEqual({ row: 2, col: 1 });
		expect(r.end).toEqual({ row: 2, col: 1 });
	});

	it('parses a multi-cell range', () => {
		const r = parseRange('A1:C3');
		expect(r.start).toEqual({ row: 0, col: 0 });
		expect(r.end).toEqual({ row: 2, col: 2 });
	});
});

// ─── getCellsInRange ─────────────────────────────────────────────────────────
describe('getCellsInRange', () => {
	it('returns a 2×2 set of cells', () => {
		const cells = getCellsInRange(0, 0, 1, 1);
		expect(cells).toHaveLength(4);
		expect(cells).toContainEqual({ row: 0, col: 0 });
		expect(cells).toContainEqual({ row: 1, col: 1 });
	});

	it('handles inverted ranges (end < start)', () => {
		const cells = getCellsInRange(2, 2, 0, 0);
		expect(cells).toHaveLength(9); // 3×3
	});

	it('returns a single cell for equal start/end', () => {
		expect(getCellsInRange(3, 3, 3, 3)).toHaveLength(1);
	});
});

// ─── detectDataType ──────────────────────────────────────────────────────────
describe('detectDataType', () => {
	it('detects empty string', () => {
		expect(detectDataType('')).toBe('empty');
	});

	it('detects null as empty', () => {
		expect(detectDataType(null)).toBe('empty');
	});

	it('detects formula', () => {
		expect(detectDataType('=SUM(A1:A5)')).toBe('formula');
	});

	it('detects integer', () => {
		expect(detectDataType('42')).toBe('integer');
	});

	it('detects decimal', () => {
		expect(detectDataType('3.14')).toBe('decimal');
	});

	it('detects date', () => {
		expect(detectDataType('12/25/2024')).toBe('date');
	});

	it('detects boolean true', () => {
		expect(detectDataType('true')).toBe('boolean');
	});

	it('detects boolean false', () => {
		expect(detectDataType('FALSE')).toBe('boolean');
	});

	it('detects plain text', () => {
		expect(detectDataType('Hello World')).toBe('text');
	});
});

// ─── getCellStyle ────────────────────────────────────────────────────────────
describe('getCellStyle', () => {
	it('returns empty string for no format', () => {
		expect(getCellStyle(null)).toBe('');
		expect(getCellStyle(undefined)).toBe('');
	});

	it('applies bold', () => {
		expect(getCellStyle({ bold: true })).toContain('font-weight: bold');
	});

	it('applies multiple styles', () => {
		const style = getCellStyle({ italic: true, color: '#ff0000' });
		expect(style).toContain('font-style: italic');
		expect(style).toContain('color: #ff0000');
	});

	it('applies alignment', () => {
		expect(getCellStyle({ align: 'center' })).toContain('text-align: center');
	});
});

// ─── extractFormulaDependencies ───────────────────────────────────────────────
describe('extractFormulaDependencies', () => {
	it('extracts simple cell ref', () => {
		const deps = extractFormulaDependencies('=A1+B2');
		expect(deps).toContain('A1');
		expect(deps).toContain('B2');
	});

	it('does NOT treat function names as cell refs', () => {
		const deps = extractFormulaDependencies('=SUM(A1:A5)');
		expect(deps).not.toContain('SUM');
	});

	it('does NOT treat MAX as a cell ref', () => {
		const deps = extractFormulaDependencies('=MAX(B1,B2)');
		expect(deps).not.toContain('MAX');
		expect(deps).toContain('B1');
		expect(deps).toContain('B2');
	});

	it('deduplicates repeated refs', () => {
		const deps = extractFormulaDependencies('=A1+A1*2');
		const a1Count = deps.filter((d) => d === 'A1').length;
		expect(a1Count).toBe(1);
	});

	it('handles empty/no-formula strings', () => {
		expect(extractFormulaDependencies('hello')).toHaveLength(0);
	});
});

// ─── sortData ────────────────────────────────────────────────────────────────
describe('sortData', () => {
	const rows = [
		['banana', 3],
		['apple', 1],
		['cherry', 2],
	];

	it('sorts strings ascending by col 0', () => {
		const sorted = sortData(rows, 0, 'asc');
		expect(sorted[0][0]).toBe('apple');
		expect(sorted[2][0]).toBe('cherry');
	});

	it('sorts strings descending by col 0', () => {
		const sorted = sortData(rows, 0, 'desc');
		expect(sorted[0][0]).toBe('cherry');
	});

	it('sorts numbers ascending by col 1', () => {
		const sorted = sortData(rows, 1, 'asc');
		expect(sorted[0][1]).toBe(1);
		expect(sorted[2][1]).toBe(3);
	});

	it('does not mutate the original array', () => {
		sortData(rows, 0, 'asc');
		expect(rows[0][0]).toBe('banana'); // unchanged
	});
});

// ─── filterData ──────────────────────────────────────────────────────────────
describe('filterData', () => {
	const rows = [
		['apple', 'red'],
		['banana', 'yellow'],
		['cherry', 'red'],
	];

	it('returns all rows when no filters', () => {
		expect(filterData(rows, new Map())).toHaveLength(3);
	});

	it('filters by column 1', () => {
		const filters = new Map([[1, new Set(['red'])]]);
		const result = filterData(rows, filters);
		expect(result).toHaveLength(2);
		expect(result.every((r) => r[1] === 'red')).toBe(true);
	});
});

// ─── generateFillSequence ────────────────────────────────────────────────────
describe('generateFillSequence', () => {
	it('generates a linear numeric sequence', () => {
		expect(generateFillSequence(1, 5, 'linear')).toEqual([1, 2, 3, 4, 5]);
	});

	it('generates a copy sequence', () => {
		expect(generateFillSequence(7, 3, 'copy')).toEqual([7, 7, 7]);
	});

	it('copies text for text values', () => {
		expect(generateFillSequence('hello', 3)).toEqual(['hello', 'hello', 'hello']);
	});
});

// ─── validateFormula ─────────────────────────────────────────────────────────
describe('validateFormula', () => {
	it('returns true for non-formula string', () => {
		expect(validateFormula('hello')).toBe(true);
	});

	it('returns true for balanced formula', () => {
		expect(validateFormula('=SUM(A1:A5)')).toBe(true);
	});

	it('returns false for unbalanced open paren', () => {
		expect(validateFormula('=SUM(A1')).toBe(false);
	});

	it('returns false for unbalanced close paren', () => {
		expect(validateFormula('=A1)')).toBe(false);
	});

	it('returns true for nested balanced parens', () => {
		expect(validateFormula('=IF(SUM(A1:A3)>0,"pos","neg")')).toBe(true);
	});
});

// ─── dataToCSV ───────────────────────────────────────────────────────────────
describe('dataToCSV', () => {
	it('converts simple data to CSV', () => {
		const csv = dataToCSV([['a', 'b'], ['c', 'd']]);
		expect(csv).toContain('a,b');
		expect(csv).toContain('c,d');
	});

	it('escapes commas inside values', () => {
		const csv = dataToCSV([['hello, world']]);
		expect(csv).toContain('"hello, world"');
	});

	it('escapes double quotes inside values', () => {
		const csv = dataToCSV([['say "hi"']]);
		expect(csv).toContain('"say ""hi"""');
	});
});

// ─── shouldAutoSave ──────────────────────────────────────────────────────────
describe('shouldAutoSave', () => {
	it('returns false when last save was recent', () => {
		expect(shouldAutoSave(Date.now() - 1000, 5)).toBe(false);
	});

	it('returns true when enough time has passed', () => {
		const sixMinutesAgo = Date.now() - 6 * 60 * 1000;
		expect(shouldAutoSave(sixMinutesAgo, 5)).toBe(true);
	});
});

// ─── normalizeKeyCombo ───────────────────────────────────────────────────────
describe('normalizeKeyCombo', () => {
	it('builds ctrl+s from event', () => {
		const event = { ctrlKey: true, metaKey: false, shiftKey: false, altKey: false, key: 's' };
		expect(normalizeKeyCombo(event)).toBe('cmd+s');
	});

	it('builds ctrl+shift+z from event', () => {
		const event = { ctrlKey: true, metaKey: false, shiftKey: true, altKey: false, key: 'z' };
		expect(normalizeKeyCombo(event)).toBe('cmd+shift+z');
	});

	it('handles plain letter keys', () => {
		const event = { ctrlKey: false, metaKey: false, shiftKey: false, altKey: false, key: 'a' };
		expect(normalizeKeyCombo(event)).toBe('a');
	});
});
