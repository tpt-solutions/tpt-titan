/**
 * Editor Formatting Utility - Rich text formatting enhancements
 */

/**
 * Formatting commands for rich text editor
 */
export const formattingCommands = {
	// Text style
	bold: { command: 'bold', icon: '𝐁', name: 'Bold', shortcut: 'Ctrl+B' },
	italic: { command: 'italic', icon: '𝐼', name: 'Italic', shortcut: 'Ctrl+I' },
	underline: { command: 'underline', icon: '𝑈', name: 'Underline', shortcut: 'Ctrl+U' },
	strikethrough: { command: 'strikeThrough', icon: 'S̶', name: 'Strikethrough', shortcut: '' },
	
	// Font
	fontColor: { command: 'foreColor', icon: '🎨', name: 'Text Color', shortcut: '' },
	backgroundColor: { command: 'hiliteColor', icon: '🖍️', name: 'Background Color', shortcut: '' },
	fontSize: { command: 'fontSize', icon: '𝐀', name: 'Font Size', shortcut: '' },
	fontName: { command: 'fontName', icon: '𝐀', name: 'Font Family', shortcut: '' },
	
	// Alignment
	alignLeft: { command: 'justifyLeft', icon: '⬅️', name: 'Align Left', shortcut: 'Ctrl+L' },
	alignCenter: { command: 'justifyCenter', icon: '↔️', name: 'Align Center', shortcut: 'Ctrl+E' },
	alignRight: { command: 'justifyRight', icon: '➡️', name: 'Align Right', shortcut: 'Ctrl+R' },
	alignJustify: { command: 'justifyFull', icon: '⬌', name: 'Justify', shortcut: 'Ctrl+J' },
	
	// Lists
	bulletList: { command: 'insertUnorderedList', icon: '•', name: 'Bullet List', shortcut: '' },
	numberList: { command: 'insertOrderedList', icon: '1.', name: 'Numbered List', shortcut: '' },
	
	// Indentation
	indent: { command: 'indent', icon: '→', name: 'Increase Indent', shortcut: 'Tab' },
	outdent: { command: 'outdent', icon: '←', name: 'Decrease Indent', shortcut: 'Shift+Tab' },
	
	// Script
	superscript: { command: 'superscript', icon: 'x²', name: 'Superscript', shortcut: '' },
	subscript: { command: 'subscript', icon: 'x₂', name: 'Subscript', shortcut: '' },
	
	// Headings
	heading1: { command: 'formatBlock', value: 'H1', icon: 'H1', name: 'Heading 1', shortcut: '' },
	heading2: { command: 'formatBlock', value: 'H2', icon: 'H2', name: 'Heading 2', shortcut: '' },
	heading3: { command: 'formatBlock', value: 'H3', icon: 'H3', name: 'Heading 3', shortcut: '' },
	paragraph: { command: 'formatBlock', value: 'P', icon: '¶', name: 'Paragraph', shortcut: '' },
	
	// Other
	removeFormat: { command: 'removeFormat', icon: '✕', name: 'Clear Formatting', shortcut: '' },
	selectAll: { command: 'selectAll', icon: '☐', name: 'Select All', shortcut: 'Ctrl+A' },
	undo: { command: 'undo', icon: '↩️', name: 'Undo', shortcut: 'Ctrl+Z' },
	redo: { command: 'redo', icon: '↪️', name: 'Redo', shortcut: 'Ctrl+Y' }
};

/**
 * Execute formatting command
 * @param {string} command - Command name
 * @param {string} value - Command value (optional)
 * @returns {boolean} Success status
 */
export function execCommand(command, value = null) {
	try {
		return document.execCommand(command, false, value);
	} catch (error) {
		console.error('Formatting command failed:', error);
		return false;
	}
}

/**
 * Check if command is active (for toolbar state)
 * @param {string} command - Command name
 * @returns {boolean} Active state
 */
export function queryCommandState(command) {
	try {
		return document.queryCommandState(command);
	} catch {
		return false;
	}
}

/**
 * Check if command is supported
 * @param {string} command - Command name
 * @returns {boolean} Supported state
 */
export function queryCommandSupported(command) {
	try {
		return document.queryCommandSupported(command);
	} catch {
		return false;
	}
}

/**
 * Get current font color
 * @returns {string} Color value
 */
export function getFontColor() {
	try {
		return document.queryCommandValue('foreColor');
	} catch {
		return '#000000';
	}
}

/**
 * Get current background color
 * @returns {string} Color value
 */
export function getBackgroundColor() {
	try {
		return document.queryCommandValue('hiliteColor');
	} catch {
		return 'transparent';
	}
}

/**
 * Color palette for text/background colors
 */
export const colorPalette = [
	// Grayscale
	'#000000', '#333333', '#666666', '#999999', '#CCCCCC', '#FFFFFF',
	// Reds
	'#FF0000', '#FF4444', '#FF6666', '#CC0000', '#990000',
	// Oranges
	'#FF8800', '#FFAA44', '#FFCC88', '#CC6600', '#994400',
	// Yellows
	'#FFFF00', '#FFFF44', '#FFFF88', '#CCCC00', '#999900',
	// Greens
	'#00FF00', '#44FF44', '#88FF88', '#00CC00', '#009900',
	// Blues
	'#0000FF', '#4444FF', '#8888FF', '#0000CC', '#000099',
	// Purples
	'#8800FF', '#AA44FF', '#CC88FF', '#6600CC', '#440099',
	// Pinks
	'#FF00FF', '#FF44FF', '#FF88FF', '#CC00CC', '#990099'
];

/**
 * Font size options
 */
export const fontSizes = [
	{ value: '1', label: 'Small', size: '12px' },
	{ value: '2', label: 'Normal', size: '14px' },
	{ value: '3', label: 'Medium', size: '16px' },
	{ value: '4', label: 'Large', size: '18px' },
	{ value: '5', label: 'X-Large', size: '24px' },
	{ value: '6', label: 'XX-Large', size: '32px' },
	{ value: '7', label: 'Huge', size: '48px' }
];

/**
 * Font family options
 */
export const fontFamilies = [
	{ value: 'Arial, sans-serif', label: 'Arial' },
	{ value: 'Times New Roman, serif', label: 'Times New Roman' },
	{ value: 'Courier New, monospace', label: 'Courier New' },
	{ value: 'Georgia, serif', label: 'Georgia' },
	{ value: 'Verdana, sans-serif', label: 'Verdana' },
	{ value: 'Helvetica, sans-serif', label: 'Helvetica' },
	{ value: 'Tahoma, sans-serif', label: 'Tahoma' },
	{ value: 'Trebuchet MS, sans-serif', label: 'Trebuchet MS' },
	{ value: 'Impact, sans-serif', label: 'Impact' },
	{ value: 'Comic Sans MS, cursive', label: 'Comic Sans MS' }
];

/**
 * Apply font color
 * @param {string} color - Color value
 */
export function applyFontColor(color) {
	execCommand('foreColor', color);
}

/**
 * Apply background color
 * @param {string} color - Color value
 */
export function applyBackgroundColor(color) {
	execCommand('hiliteColor', color);
}

/**
 * Apply font size
 * @param {string} size - Size value (1-7)
 */
export function applyFontSize(size) {
	execCommand('fontSize', size);
}

/**
 * Apply font family
 * @param {string} fontFamily - Font family value
 */
export function applyFontFamily(fontFamily) {
	execCommand('fontName', fontFamily);
}

/**
 * Toggle bold formatting
 */
export function toggleBold() {
	execCommand('bold');
}

/**
 * Toggle italic formatting
 */
export function toggleItalic() {
	execCommand('italic');
}

/**
 * Toggle underline formatting
 */
export function toggleUnderline() {
	execCommand('underline');
}

/**
 * Toggle strikethrough formatting
 */
export function toggleStrikethrough() {
	execCommand('strikeThrough');
}

/**
 * Toggle superscript
 */
export function toggleSuperscript() {
	execCommand('superscript');
}

/**
 * Toggle subscript
 */
export function toggleSubscript() {
	execCommand('subscript');
}

/**
 * Apply heading format
 * @param {string} level - Heading level (H1, H2, H3, P)
 */
export function applyHeading(level) {
	execCommand('formatBlock', level);
}

/**
 * Insert bullet list
 */
export function insertBulletList() {
	execCommand('insertUnorderedList');
}

/**
 * Insert numbered list
 */
export function insertNumberedList() {
	execCommand('insertOrderedList');
}

/**
 * Increase indent
 */
export function increaseIndent() {
	execCommand('indent');
}

/**
 * Decrease indent
 */
export function decreaseIndent() {
	execCommand('outdent');
}

/**
 * Align text left
 */
export function alignLeft() {
	execCommand('justifyLeft');
}

/**
 * Align text center
 */
export function alignCenter() {
	execCommand('justifyCenter');
}

/**
 * Align text right
 */
export function alignRight() {
	execCommand('justifyRight');
}

/**
 * Justify text
 */
export function alignJustify() {
	execCommand('justifyFull');
}

/**
 * Remove all formatting
 */
export function clearFormatting() {
	execCommand('removeFormat');
}

/**
 * Get word count from HTML content
 * @param {string} html - HTML content
 * @returns {Object} Word and character counts
 */
export function getWordCount(html) {
	// Strip HTML tags
	const text = html.replace(/<[^>]*>/g, ' ');
	
	// Count words
	const words = text.trim().split(/\s+/).filter(word => word.length > 0);
	const wordCount = words.length;
	
	// Count characters
	const charCount = text.replace(/\s/g, '').length;
	const charCountWithSpaces = text.length;
	
	return {
		words: wordCount,
		characters: charCount,
		charactersWithSpaces: charCountWithSpaces
	};
}

/**
 * Print editor content
 * @param {HTMLElement} element - Editor element
 */
export function printContent(element) {
	if (!element) return;
	
	const printWindow = window.open('', '_blank');
	if (!printWindow) {
		alert('Please allow popups to print');
		return;
	}
	
	printWindow.document.write(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Print Document</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; padding: 20px; }
				@media print { body { padding: 0; } }
			</style>
		</head>
		<body>
			${element.innerHTML}
		</body>
		</html>
	`);
	
	printWindow.document.close();
	printWindow.focus();
	
	setTimeout(() => {
		printWindow.print();
		printWindow.close();
	}, 250);
}

/**
 * Enable spell check
 * @param {HTMLElement} element - Editor element
 * @param {boolean} enabled - Enable/disable
 */
export function setSpellCheck(element, enabled) {
	if (element) {
		element.spellcheck = enabled;
	}
}

/**
 * Format keyboard shortcuts handler
 * @param {KeyboardEvent} event 
 * @returns {boolean} Whether shortcut was handled
 */
export function handleFormattingShortcuts(event) {
	if (!event.ctrlKey && !event.metaKey) return false;
	
	const key = event.key.toLowerCase();
	
	switch (key) {
		case 'b':
			event.preventDefault();
			toggleBold();
			return true;
		case 'i':
			event.preventDefault();
			toggleItalic();
			return true;
		case 'u':
			event.preventDefault();
			toggleUnderline();
			return true;
		case 'l':
			event.preventDefault();
			alignLeft();
			return true;
		case 'e':
			event.preventDefault();
			alignCenter();
			return true;
		case 'r':
			event.preventDefault();
			alignRight();
			return true;
		case 'j':
			event.preventDefault();
			alignJustify();
			return true;
		default:
			return false;
	}
}

export default {
	formattingCommands,
	execCommand,
	queryCommandState,
	queryCommandSupported,
	colorPalette,
	fontSizes,
	fontFamilies,
	applyFontColor,
	applyBackgroundColor,
	applyFontSize,
	applyFontFamily,
	toggleBold,
	toggleItalic,
	toggleUnderline,
	toggleStrikethrough,
	toggleSuperscript,
	toggleSubscript,
	applyHeading,
	insertBulletList,
	insertNumberedList,
	increaseIndent,
	decreaseIndent,
	alignLeft,
	alignCenter,
	alignRight,
	alignJustify,
	clearFormatting,
	getWordCount,
	printContent,
	setSpellCheck,
	handleFormattingShortcuts
};
