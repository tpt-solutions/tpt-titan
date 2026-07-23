// @ts-nocheck
/**
 * Editor Export Utility - Export documents to various formats
 */

/**
 * Convert blocks to plain text
 * @param {Array} blocks - Array of block objects
 * @returns {string} Plain text content
 */
export function blocksToText(blocks) {
	if (!blocks || blocks.length === 0) return '';
	
	return blocks.map(block => {
		const content = block.content || '';
		switch (block.type) {
			case 'heading1':
				return `# ${content}\n\n`;
			case 'heading2':
				return `## ${content}\n\n`;
			case 'heading3':
				return `### ${content}\n\n`;
			case 'list':
				return `- ${content}\n`;
			case 'quote':
				return `> ${content}\n\n`;
			case 'code':
				return `\`\`\`\n${content}\n\`\`\`\n\n`;
			case 'math':
				return `[MATH: ${content}]\n\n`;
			case 'divider':
				return '---\n\n';
			default:
				return `${content}\n\n`;
		}
	}).join('').trim();
}

/**
 * Convert blocks to HTML
 * @param {Array} blocks - Array of block objects
 * @returns {string} HTML content
 */
export function blocksToHTML(blocks) {
	if (!blocks || blocks.length === 0) return '';
	
	const htmlParts = blocks.map(block => {
		const content = escapeHtml(block.content || '');
		switch (block.type) {
			case 'heading1':
				return `<h1>${content}</h1>`;
			case 'heading2':
				return `<h2>${content}</h2>`;
			case 'heading3':
				return `<h3>${content}</h3>`;
			case 'list':
				return `<li>${content}</li>`;
			case 'quote':
				return `<blockquote>${content}</blockquote>`;
			case 'code':
				return `<pre><code>${content}</code></pre>`;
			case 'math':
				return `<div class="math">${content}</div>`;
			case 'divider':
				return '<hr>';
			case 'image':
				return `<img src="${content}" alt="Image">`;
			case 'table':
				return renderTableHTML(block);
			default:
				return `<p>${content}</p>`;
		}
	});
	
	// Wrap lists in ul/ol tags
	let html = '';
	let inList = false;
	
	for (const part of htmlParts) {
		if (part.startsWith('<li>')) {
			if (!inList) {
				html += '<ul>';
				inList = true;
			}
			html += part;
		} else {
			if (inList) {
				html += '</ul>';
				inList = false;
			}
			html += part;
		}
	}
	
	if (inList) {
		html += '</ul>';
	}
	
	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Exported Document</title>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; max-width: 800px; margin: 0 auto; padding: 20px; }
		h1, h2, h3 { color: #333; }
		blockquote { border-left: 3px solid #ccc; margin: 0; padding-left: 20px; color: #666; }
		pre { background: #f4f4f4; padding: 15px; border-radius: 5px; overflow-x: auto; }
		code { font-family: 'Courier New', monospace; }
		.math { font-style: italic; background: #f9f9f9; padding: 10px; border-radius: 5px; }
		img { max-width: 100%; height: auto; }
		table { border-collapse: collapse; width: 100%; margin: 15px 0; }
		td, th { border: 1px solid #ddd; padding: 8px; }
		th { background: #f4f4f4; }
	</style>
</head>
<body>
${html}
</body>
</html>`;
}

/**
 * Render table block as HTML
 */
function renderTableHTML(block) {
	if (!block.properties || !block.properties.rows) return '<table></table>';
	
	const rows = block.properties.rows || [];
	const columns = block.properties.columns || [];
	
	let html = '<table><thead><tr>';
	columns.forEach(col => {
		html += `<th>${escapeHtml(col)}</th>`;
	});
	html += '</tr></thead><tbody>';
	
	rows.forEach(row => {
		html += '<tr>';
		columns.forEach((col, i) => {
			html += `<td>${escapeHtml(row[i] || '')}</td>`;
		});
		html += '</tr>';
	});
	
	html += '</tbody></table>';
	return html;
}

/**
 * Convert blocks to Markdown
 * @param {Array} blocks - Array of block objects
 * @returns {string} Markdown content
 */
export function blocksToMarkdown(blocks) {
	if (!blocks || blocks.length === 0) return '';
	
	return blocks.map(block => {
		const content = block.content || '';
		switch (block.type) {
			case 'heading1':
				return `# ${content}\n\n`;
			case 'heading2':
				return `## ${content}\n\n`;
			case 'heading3':
				return `### ${content}\n\n`;
			case 'list':
				return `- ${content}\n`;
			case 'quote':
				return `> ${content}\n\n`;
			case 'code':
				return `\`\`\`\n${content}\n\`\`\`\n\n`;
			case 'math':
				return `$$\n${content}\n$$\n\n`;
			case 'divider':
				return '---\n\n';
			case 'image':
				return `![Image](${content})\n\n`;
			case 'table':
				return renderTableMarkdown(block);
			default:
				return `${content}\n\n`;
		}
	}).join('').trim();
}

/**
 * Render table block as Markdown
 */
function renderTableMarkdown(block) {
	if (!block.properties || !block.properties.rows) return '';
	
	const rows = block.properties.rows || [];
	const columns = block.properties.columns || [];
	
	if (columns.length === 0) return '';
	
	let md = '| ' + columns.join(' | ') + ' |\n';
	md += '| ' + columns.map(() => '---').join(' | ') + ' |\n';
	
	rows.forEach(row => {
		md += '| ' + columns.map((col, i) => row[i] || '').join(' | ') + ' |\n';
	});
	
	return md + '\n';
}

/**
 * Convert markdown to DOCX-compatible HTML
 * @param {string} markdown - Markdown content
 * @returns {string} HTML for DOCX conversion
 */
export function markdownToDocxHTML(markdown) {
	if (!markdown) return '';
	
	// Simple markdown to HTML conversion
	let html = markdown
		.replace(/^### (.*$)/gim, '<h3>$1</h3>')
		.replace(/^## (.*$)/gim, '<h2>$1</h2>')
		.replace(/^# (.*$)/gim, '<h1>$1</h1>')
		.replace(/\*\*(.*)\*\*/gim, '<strong>$1</strong>')
		.replace(/\*(.*)\*/gim, '<em>$1</em>')
		.replace(/`([^`]+)`/gim, '<code>$1</code>')
		.replace(/^\> (.*$)/gim, '<blockquote>$1</blockquote>')
		.replace(/^\- (.*$)/gim, '<li>$1</li>')
		.replace(/^\d+\. (.*$)/gim, '<li>$1</li>')
		.replace(/\n/gim, '<br>');
	
	return html;
}

/**
 * Create a full HTML document for DOCX export
 * @param {string} title - Document title
 * @param {string} contentHTML - HTML content
 * @returns {string} Full HTML document
 */
export function createDocxHTML(title, contentHTML) {
	return `<!DOCTYPE html>
<html xmlns:o='urn:schemas-microsoft-com:office:office' xmlns:w='urn:schemas-microsoft-com:office:word' xmlns='http://www.w3.org/TR/REC-html40'>
<head>
	<meta charset="utf-8">
	<title>${escapeHtml(title)}</title>
	<!--[if gte mso 9]>
	<xml>
		<w:WordDocument>
			<w:View>Print</w:View>
			<w:Zoom>100</w:Zoom>
		</w:WordDocument>
	</xml>
	<![endif]-->
	<style>
		body { font-family: 'Times New Roman', serif; font-size: 12pt; line-height: 1.5; }
		h1 { font-size: 24pt; font-weight: bold; margin-top: 24pt; margin-bottom: 12pt; }
		h2 { font-size: 18pt; font-weight: bold; margin-top: 18pt; margin-bottom: 9pt; }
		h3 { font-size: 14pt; font-weight: bold; margin-top: 14pt; margin-bottom: 7pt; }
		p { margin-top: 6pt; margin-bottom: 6pt; }
		blockquote { margin-left: 36pt; margin-right: 36pt; font-style: italic; }
		pre { font-family: 'Courier New', monospace; background: #f4f4f4; padding: 10pt; }
		code { font-family: 'Courier New', monospace; }
		table { border-collapse: collapse; width: 100%; }
		td, th { border: 1px solid black; padding: 5pt; }
	</style>
</head>
<body>
	<h1>${escapeHtml(title)}</h1>
	${contentHTML}
</body>
</html>`;
}

/**
 * Download content as a file
 * @param {string} content - File content
 * @param {string} filename - File name
 * @param {string} mimeType - MIME type
 */
export function downloadFile(content, filename, mimeType) {
	const blob = new Blob([content], { type: mimeType });
	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = filename;
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
	URL.revokeObjectURL(url);
}

/**
 * Export blocks to TXT file
 * @param {Array} blocks - Array of block objects
 * @param {string} filename - File name (without extension)
 */
export function exportToTXT(blocks, filename = 'document') {
	const content = blocksToText(blocks);
	downloadFile(content, `${filename}.txt`, 'text/plain');
}

/**
 * Export blocks to HTML file
 * @param {Array} blocks - Array of block objects
 * @param {string} filename - File name (without extension)
 */
export function exportToHTML(blocks, filename = 'document') {
	const content = blocksToHTML(blocks);
	downloadFile(content, `${filename}.html`, 'text/html');
}

/**
 * Export blocks to Markdown file
 * @param {Array} blocks - Array of block objects
 * @param {string} filename - File name (without extension)
 */
export function exportToMarkdown(blocks, filename = 'document') {
	const content = blocksToMarkdown(blocks);
	downloadFile(content, `${filename}.md`, 'text/markdown');
}

/**
 * Export to DOCX-compatible HTML (can be opened in Word)
 * @param {Array} blocks - Array of block objects
 * @param {string} title - Document title
 * @param {string} filename - File name (without extension)
 */
export function exportToDOCX(blocks, title = 'Document', filename = 'document') {
	const markdown = blocksToMarkdown(blocks);
	const contentHTML = markdownToDocxHTML(markdown);
	const fullHTML = createDocxHTML(title, contentHTML);
	downloadFile(fullHTML, `${filename}.docx.html`, 'application/msword');
}

/**
 * Export based on editor mode
 * @param {string} mode - Editor mode ('blocks', 'markdown', 'richtext')
 * @param {Array|string} content - Content (blocks array or markdown string)
 * @param {string} title - Document title
 * @param {string} format - Export format ('txt', 'html', 'md', 'docx')
 * @param {string} filename - File name (without extension)
 */
export function exportDocument(mode, content, title, format, filename = 'document') {
	let blocks;
	
	if (mode === 'blocks') {
		blocks = content;
	} else if (mode === 'markdown') {
		// Convert markdown to blocks for consistent export
		blocks = markdownToBlocks(content);
	} else if (mode === 'richtext') {
		// Convert HTML to blocks
		blocks = htmlToBlocks(content);
	}
	
	switch (format) {
		case 'txt':
			exportToTXT(blocks, filename);
			break;
		case 'html':
			exportToHTML(blocks, filename);
			break;
		case 'md':
			exportToMarkdown(blocks, filename);
			break;
		case 'docx':
			exportToDOCX(blocks, title, filename);
			break;
		default:
			exportToTXT(blocks, filename);
	}
}

/**
 * Simple markdown to blocks converter
 * @param {string} markdown - Markdown content
 * @returns {Array} Blocks array
 */
function markdownToBlocks(markdown) {
	if (!markdown) return [];
	
	const lines = markdown.split('\n');
	const blocks = [];
	let id = 1;
	
	for (const line of lines) {
		const trimmed = line.trim();
		if (!trimmed) continue;
		
		if (trimmed.startsWith('# ')) {
			blocks.push({ id: id++, type: 'heading1', content: trimmed.substring(2), properties: {} });
		} else if (trimmed.startsWith('## ')) {
			blocks.push({ id: id++, type: 'heading2', content: trimmed.substring(3), properties: {} });
		} else if (trimmed.startsWith('### ')) {
			blocks.push({ id: id++, type: 'heading3', content: trimmed.substring(4), properties: {} });
		} else if (trimmed.startsWith('- ') || trimmed.startsWith('* ')) {
			blocks.push({ id: id++, type: 'list', content: trimmed.substring(2), properties: {} });
		} else if (trimmed.startsWith('> ')) {
			blocks.push({ id: id++, type: 'quote', content: trimmed.substring(2), properties: {} });
		} else if (trimmed === '---' || trimmed === '***') {
			blocks.push({ id: id++, type: 'divider', content: '', properties: {} });
		} else {
			blocks.push({ id: id++, type: 'text', content: trimmed, properties: {} });
		}
	}
	
	return blocks;
}

/**
 * Simple HTML to blocks converter
 * @param {string} html - HTML content
 * @returns {Array} Blocks array
 */
function htmlToBlocks(html) {
	if (!html) return [];
	
	// Create a temporary div to parse HTML
	const div = document.createElement('div');
	div.innerHTML = html;
	
	const blocks = [];
	let id = 1;
	
	const elements = div.querySelectorAll('h1, h2, h3, p, blockquote, li, pre, hr, div');
	
	for (const el of elements) {
		const text = el.textContent || '';
		
		switch (el.tagName.toLowerCase()) {
			case 'h1':
				blocks.push({ id: id++, type: 'heading1', content: text, properties: {} });
				break;
			case 'h2':
				blocks.push({ id: id++, type: 'heading2', content: text, properties: {} });
				break;
			case 'h3':
				blocks.push({ id: id++, type: 'heading3', content: text, properties: {} });
				break;
			case 'blockquote':
				blocks.push({ id: id++, type: 'quote', content: text, properties: {} });
				break;
			case 'li':
				blocks.push({ id: id++, type: 'list', content: text, properties: {} });
				break;
			case 'pre':
				blocks.push({ id: id++, type: 'code', content: text, properties: {} });
				break;
			case 'hr':
				blocks.push({ id: id++, type: 'divider', content: '', properties: {} });
				break;
			default:
				if (text.trim()) {
					blocks.push({ id: id++, type: 'text', content: text, properties: {} });
				}
		}
	}
	
	return blocks;
}

/**
 * Escape HTML special characters
 * @param {string} text 
 * @returns {string}
 */
function escapeHtml(text) {
	if (!text) return '';
	return text
		.replace(/&/g, '&amp;')
		.replace(/</g, '<')
		.replace(/>/g, '>')
		.replace(/"/g, '"')
		.replace(/'/g, '&#039;');
}

export default {
	blocksToText,
	blocksToHTML,
	blocksToMarkdown,
	markdownToDocxHTML,
	createDocxHTML,
	downloadFile,
	exportToTXT,
	exportToHTML,
	exportToMarkdown,
	exportToDOCX,
	exportDocument
};
