// @ts-nocheck
/**
 * Editor Search Utility - Find and replace functionality for all editor modes
 */

/**
 * Search result type
 * @typedef {Object} SearchResult
 * @property {number} index - Index of the match
 * @property {string} text - Matched text
 * @property {number} length - Length of matched text
 * @property {Object} [block] - Block reference (for block editor)
 * @property {number} [blockIndex] - Block index (for block editor)
 * @property {number} [startOffset] - Start offset in content
 * @property {number} [endOffset] - End offset in content
 */

/**
 * Search in block editor content
 * @param {Array} blocks - Array of block objects
 * @param {string} searchText - Text to search for
 * @param {boolean} caseSensitive - Whether search is case sensitive
 * @returns {SearchResult[]} Array of search results
 */
export function searchInBlocks(blocks, searchText, caseSensitive = false) {
	if (!searchText || !blocks || blocks.length === 0) return [];
	
	const results = [];
	const flags = caseSensitive ? 'g' : 'gi';
	const regex = new RegExp(escapeRegExp(searchText), flags);
	
	blocks.forEach((block, blockIndex) => {
		if (!block.content) return;
		
		const content = block.content;
		let match;
		
		while ((match = regex.exec(content)) !== null) {
			results.push({
				index: results.length,
				text: match[0],
				length: match[0].length,
				block: block,
				blockIndex: blockIndex,
				startOffset: match.index,
				endOffset: match.index + match[0].length
			});
			
			// Prevent infinite loop on zero-length matches
			if (match[0].length === 0) break;
		}
	});
	
	return results;
}

/**
 * Search in markdown text
 * @param {string} markdown - Markdown content
 * @param {string} searchText - Text to search for
 * @param {boolean} caseSensitive - Whether search is case sensitive
 * @returns {SearchResult[]} Array of search results
 */
export function searchInMarkdown(markdown, searchText, caseSensitive = false) {
	if (!searchText || !markdown) return [];
	
	const results = [];
	const flags = caseSensitive ? 'g' : 'gi';
	const regex = new RegExp(escapeRegExp(searchText), flags);
	
	let match;
	while ((match = regex.exec(markdown)) !== null) {
		results.push({
			index: results.length,
			text: match[0],
			length: match[0].length,
			startOffset: match.index,
			endOffset: match.index + match[0].length
		});
		
		// Prevent infinite loop on zero-length matches
		if (match[0].length === 0) break;
	}
	
	return results;
}

/**
 * Search in rich text HTML
 * @param {HTMLElement} element - Rich text editor element
 * @param {string} searchText - Text to search for
 * @param {boolean} caseSensitive - Whether search is case sensitive
 * @returns {SearchResult[]} Array of search results with DOM ranges
 */
export function searchInRichText(element, searchText, caseSensitive = false) {
	if (!searchText || !element) return [];
	
	const results = [];
	const walker = document.createTreeWalker(
		element,
		NodeFilter.SHOW_TEXT,
		null,
		false
	);
	
	const flags = caseSensitive ? 'g' : 'gi';
	const regex = new RegExp(escapeRegExp(searchText), flags);
	
	let node;
	let globalOffset = 0;
	const textNodes = [];
	
	// Collect all text nodes with their offsets
	while ((node = walker.nextNode()) !== null) {
		textNodes.push({
			node: node,
			startOffset: globalOffset,
			endOffset: globalOffset + node.textContent.length,
			text: node.textContent
		});
		globalOffset += node.textContent.length;
	}
	
	// Search in concatenated text
	const fullText = textNodes.map(tn => tn.text).join('');
	let match;
	
	while ((match = regex.exec(fullText)) !== null) {
		const matchStart = match.index;
		const matchEnd = match.index + match[0].length;
		
		// Find which text node(s) contain this match
		const startNodeInfo = textNodes.find(tn => 
			matchStart >= tn.startOffset && matchStart < tn.endOffset
		);
		const endNodeInfo = textNodes.find(tn => 
			matchEnd > tn.startOffset && matchEnd <= tn.endOffset
		);
		
		if (startNodeInfo && endNodeInfo) {
			results.push({
				index: results.length,
				text: match[0],
				length: match[0].length,
				startNode: startNodeInfo.node,
				endNode: endNodeInfo.node,
				startOffsetInNode: matchStart - startNodeInfo.startOffset,
				endOffsetInNode: matchEnd - endNodeInfo.startOffset,
				range: createRange(startNodeInfo.node, matchStart - startNodeInfo.startOffset,
				                  endNodeInfo.node, matchEnd - endNodeInfo.startOffset)
			});
		}
		
		// Prevent infinite loop on zero-length matches
		if (match[0].length === 0) break;
	}
	
	return results;
}

/**
 * Create a DOM range
 */
function createRange(startNode, startOffset, endNode, endOffset) {
	const range = document.createRange();
	range.setStart(startNode, startOffset);
	range.setEnd(endNode, endOffset);
	return range;
}

/**
 * Replace text in block editor
 * @param {Array} blocks - Blocks array
 * @param {SearchResult} result - Search result to replace
 * @param {string} replacement - Replacement text
 * @returns {Array} Updated blocks array
 */
export function replaceInBlocks(blocks, result, replacement) {
	if (!result || result.blockIndex === undefined) return blocks;
	
	const newBlocks = [...blocks];
	const block = { ...newBlocks[result.blockIndex] };
	const content = block.content;
	
	block.content = content.substring(0, result.startOffset) + 
	                replacement + 
	                content.substring(result.endOffset);
	
	newBlocks[result.blockIndex] = block;
	return newBlocks;
}

/**
 * Replace text in markdown
 * @param {string} markdown - Markdown content
 * @param {SearchResult} result - Search result to replace
 * @param {string} replacement - Replacement text
 * @returns {string} Updated markdown
 */
export function replaceInMarkdown(markdown, result, replacement) {
	if (!result || result.startOffset === undefined) return markdown;
	
	return markdown.substring(0, result.startOffset) + 
	       replacement + 
	       markdown.substring(result.endOffset);
}

/**
 * Replace all occurrences in text
 * @param {string} text - Original text
 * @param {string} searchText - Text to search for
 * @param {string} replacement - Replacement text
 * @param {boolean} caseSensitive - Whether search is case sensitive
 * @returns {string} Updated text
 */
export function replaceAll(text, searchText, replacement, caseSensitive = false) {
	if (!searchText) return text;
	
	const flags = caseSensitive ? 'g' : 'gi';
	const regex = new RegExp(escapeRegExp(searchText), flags);
	return text.replace(regex, replacement);
}

/**
 * Replace all in blocks
 * @param {Array} blocks - Blocks array
 * @param {string} searchText - Text to search for
 * @param {string} replacement - Replacement text
 * @param {boolean} caseSensitive - Whether search is case sensitive
 * @returns {Array} Updated blocks array
 */
export function replaceAllInBlocks(blocks, searchText, replacement, caseSensitive = false) {
	if (!searchText || !blocks) return blocks;
	
	const flags = caseSensitive ? 'g' : 'gi';
	const regex = new RegExp(escapeRegExp(searchText), flags);
	
	return blocks.map(block => {
		if (!block.content) return block;
		return {
			...block,
			content: block.content.replace(regex, replacement)
		};
	});
}

/**
 * Escape special regex characters
 * @param {string} string 
 * @returns {string}
 */
function escapeRegExp(string) {
	return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/**
 * Highlight search result in DOM
 * @param {SearchResult} result 
 * @param {string} highlightClass - CSS class for highlighting
 */
export function highlightResult(result, highlightClass = 'search-highlight') {
	if (result.range) {
		const selection = window.getSelection();
		selection.removeAllRanges();
		selection.addRange(result.range);
		
		// Scroll into view
		const rect = result.range.getBoundingClientRect();
		if (rect) {
			result.range.startContainer.parentElement?.scrollIntoView({
				behavior: 'smooth',
				block: 'center'
			});
		}
	}
}

/**
 * Clear all highlights
 */
export function clearHighlights() {
	const selection = window.getSelection();
	selection.removeAllRanges();
}

export default {
	searchInBlocks,
	searchInMarkdown,
	searchInRichText,
	replaceInBlocks,
	replaceInMarkdown,
	replaceAll,
	replaceAllInBlocks,
	highlightResult,
	clearHighlights
};
