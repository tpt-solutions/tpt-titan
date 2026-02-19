/**
 * Editor Image Utility - Image upload and handling for text editor
 */

/**
 * Image upload configuration
 */
export const imageConfig = {
	maxFileSize: 5 * 1024 * 1024, // 5MB
	acceptedTypes: ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml'],
	maxWidth: 1200,
	maxHeight: 1200,
	quality: 0.9
};

/**
 * Validate image file
 * @param {File} file - Image file
 * @returns {Object} Validation result
 */
export function validateImage(file) {
	if (!file) {
		return { valid: false, error: 'No file selected' };
	}
	
	if (!imageConfig.acceptedTypes.includes(file.type)) {
		return { 
			valid: false, 
			error: `Invalid file type. Accepted: ${imageConfig.acceptedTypes.map(t => t.replace('image/', '.')).join(', ')}` 
		};
	}
	
	if (file.size > imageConfig.maxFileSize) {
		return { 
			valid: false, 
			error: `File too large. Max size: ${(imageConfig.maxFileSize / 1024 / 1024).toFixed(1)}MB` 
		};
	}
	
	return { valid: true, error: null };
}

/**
 * Resize image to fit within max dimensions
 * @param {HTMLImageElement} img - Image element
 * @returns {Object} New dimensions
 */
export function calculateResizeDimensions(img) {
	let { width, height } = img;
	
	if (width > imageConfig.maxWidth) {
		height = (height * imageConfig.maxWidth) / width;
		width = imageConfig.maxWidth;
	}
	
	if (height > imageConfig.maxHeight) {
		width = (width * imageConfig.maxHeight) / height;
		height = imageConfig.maxHeight;
	}
	
	return { width, height };
}

/**
 * Resize and compress image
 * @param {File} file - Image file
 * @returns {Promise<string>} Resized image as base64 data URL
 */
export function resizeImage(file) {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		
		reader.onload = (e) => {
			const img = new Image();
			
			img.onload = () => {
				const canvas = document.createElement('canvas');
				const ctx = canvas.getContext('2d');
				
				const { width, height } = calculateResizeDimensions(img);
				
				canvas.width = width;
				canvas.height = height;
				
				// Use better quality scaling
				ctx.imageSmoothingEnabled = true;
				ctx.imageSmoothingQuality = 'high';
				
				ctx.drawImage(img, 0, 0, width, height);
				
				// Convert to JPEG with compression
				const dataUrl = canvas.toDataURL('image/jpeg', imageConfig.quality);
				resolve(dataUrl);
			};
			
			img.onerror = () => reject(new Error('Failed to load image'));
			img.src = e.target.result;
		};
		
		reader.onerror = () => reject(new Error('Failed to read file'));
		reader.readAsDataURL(file);
	});
}

/**
 * Upload image to backend (simulated)
 * @param {string} dataUrl - Base64 image data
 * @param {string} filename - Original filename
 * @returns {Promise<Object>} Upload result with URL
 */
export async function uploadImage(dataUrl, filename) {
	// In a real implementation, this would upload to a server
	// For now, we store in localStorage or return the data URL
	
	try {
		// Simulate API call
		await new Promise(resolve => setTimeout(resolve, 500));
		
		// Generate a unique ID for the image
		const imageId = `img_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
		
		// Store in localStorage (for demo purposes)
		// In production, this would be stored on a server
		const imageData = {
			id: imageId,
			dataUrl: dataUrl,
			filename: filename,
			uploadedAt: new Date().toISOString()
		};
		
		// Store metadata
		const images = JSON.parse(localStorage.getItem('editor_images') || '[]');
		images.push(imageData);
		localStorage.setItem('editor_images', JSON.stringify(images));
		
		return {
			success: true,
			url: dataUrl, // In production, this would be a server URL
			id: imageId,
			filename: filename
		};
	} catch (error) {
		return {
			success: false,
			error: error.message
		};
	}
}

/**
 * Insert image into rich text editor
 * @param {HTMLElement} editorElement - Editor element
 * @param {string} imageUrl - Image URL or data URL
 * @param {Object} options - Image options (alignment, width, etc.)
 */
export function insertImageInEditor(editorElement, imageUrl, options = {}) {
	if (!editorElement) return;
	
	const { alignment = 'center', width = 'auto', caption = '' } = options;
	
	// Create image wrapper
	const wrapper = document.createElement('div');
	wrapper.className = `image-wrapper image-align-${alignment}`;
	wrapper.style.margin = '1em 0';
	wrapper.style.textAlign = alignment;
	wrapper.contentEditable = false;
	
	// Create image
	const img = document.createElement('img');
	img.src = imageUrl;
	img.style.maxWidth = width === 'auto' ? '100%' : width;
	img.style.height = 'auto';
	img.style.borderRadius = '4px';
	img.style.boxShadow = '0 2px 8px rgba(0,0,0,0.1)';
	
	wrapper.appendChild(img);
	
	// Add caption if provided
	if (caption) {
		const captionEl = document.createElement('div');
		captionEl.className = 'image-caption';
		captionEl.textContent = caption;
		captionEl.style.fontSize = '0.9em';
		captionEl.style.color = '#666';
		captionEl.style.marginTop = '0.5em';
		captionEl.style.fontStyle = 'italic';
		wrapper.appendChild(captionEl);
	}
	
	// Insert at cursor position or append
	const selection = window.getSelection();
	if (selection.rangeCount > 0) {
		const range = selection.getRangeAt(0);
		range.deleteContents();
		range.insertNode(wrapper);
		
		// Move cursor after image
		range.setStartAfter(wrapper);
		range.setEndAfter(wrapper);
		selection.removeAllRanges();
		selection.addRange(range);
	} else {
		editorElement.appendChild(wrapper);
	}
	
	// Add paragraph after image for continued typing
	const p = document.createElement('p');
	p.innerHTML = '<br>';
	editorElement.appendChild(p);
}

/**
 * Insert image in block editor
 * @param {Array} blocks - Blocks array
 * @param {number} selectedIndex - Selected block index
 * @param {string} imageUrl - Image URL
 * @param {string} caption - Image caption
 * @returns {Array} Updated blocks
 */
export function insertImageInBlocks(blocks, selectedIndex, imageUrl, caption = '') {
	const newBlock = {
		id: Date.now() + Math.random(),
		type: 'image',
		content: imageUrl,
		properties: {
			caption: caption,
			alignment: 'center'
		}
	};
	
	const newBlocks = [...blocks];
	newBlocks.splice(selectedIndex + 1, 0, newBlock);
	
	// Add empty text block after image
	newBlocks.splice(selectedIndex + 2, 0, {
		id: Date.now() + Math.random() + 1,
		type: 'text',
		content: '',
		properties: {}
	});
	
	return newBlocks;
}

/**
 * Insert image markdown
 * @param {string} markdown - Current markdown content
 * @param {string} imageUrl - Image URL
 * @param {string} altText - Alt text
 * @returns {string} Updated markdown
 */
export function insertImageInMarkdown(markdown, imageUrl, altText = 'Image') {
	const imageMarkdown = `\n\n![${altText}](${imageUrl})\n\n`;
	
	// For simplicity, append at end or insert at cursor position
	// In a real implementation, you'd track cursor position
	return markdown + imageMarkdown;
}

/**
 * Handle image file selection
 * @param {FileList} files - Selected files
 * @returns {Promise<Array>} Array of upload results
 */
export async function handleImageUpload(files) {
	const results = [];
	
	for (const file of Array.from(files)) {
		const validation = validateImage(file);
		
		if (!validation.valid) {
			results.push({ success: false, error: validation.error, filename: file.name });
			continue;
		}
		
		try {
			// Resize image
			const resizedDataUrl = await resizeImage(file);
			
			// Upload to storage
			const uploadResult = await uploadImage(resizedDataUrl, file.name);
			
			results.push(uploadResult);
		} catch (error) {
			results.push({ success: false, error: error.message, filename: file.name });
		}
	}
	
	return results;
}

/**
 * Get all stored images
 * @returns {Array} Array of image metadata
 */
export function getStoredImages() {
	return JSON.parse(localStorage.getItem('editor_images') || '[]');
}

/**
 * Delete stored image
 * @param {string} imageId - Image ID to delete
 */
export function deleteStoredImage(imageId) {
	const images = getStoredImages().filter(img => img.id !== imageId);
	localStorage.setItem('editor_images', JSON.stringify(images));
}

/**
 * Create image upload input element
 * @param {Function} onUpload - Callback when images are uploaded
 * @returns {HTMLInputElement} File input element
 */
export function createImageUploadInput(onUpload) {
	const input = document.createElement('input');
	input.type = 'file';
	input.accept = imageConfig.acceptedTypes.join(',');
	input.multiple = true;
	input.style.display = 'none';
	
	input.addEventListener('change', async (e) => {
		if (e.target.files && e.target.files.length > 0) {
			const results = await handleImageUpload(e.target.files);
			onUpload(results);
		}
		input.value = ''; // Reset for re-upload
	});
	
	return input;
}

/**
 * Image alignment options
 */
export const imageAlignments = [
	{ id: 'left', name: 'Left', icon: '⬅️' },
	{ id: 'center', name: 'Center', icon: '⬇️' },
	{ id: 'right', name: 'Right', icon: '➡️' }
];

/**
 * Image size options
 */
export const imageSizes = [
	{ id: 'small', name: 'Small', width: '200px' },
	{ id: 'medium', name: 'Medium', width: '400px' },
	{ id: 'large', name: 'Large', width: '600px' },
	{ id: 'full', name: 'Full Width', width: '100%' },
	{ id: 'auto', name: 'Auto', width: 'auto' }
];

export default {
	imageConfig,
	validateImage,
	resizeImage,
	uploadImage,
	insertImageInEditor,
	insertImageInBlocks,
	insertImageInMarkdown,
	handleImageUpload,
	getStoredImages,
	deleteStoredImage,
	createImageUploadInput,
	imageAlignments,
	imageSizes
};
