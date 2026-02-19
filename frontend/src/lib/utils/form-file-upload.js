/**
 * Form File Upload Utility - Handle file uploads for forms
 */

/**
 * File upload configuration
 */
export const uploadConfig = {
	maxFileSize: 10 * 1024 * 1024, // 10MB default
	maxFiles: 5, // Maximum files per field
	acceptedTypes: {
		'file': '*',
		'image-upload': ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
		'video-upload': ['video/mp4', 'video/webm', 'video/ogg'],
		'audio-upload': ['audio/mpeg', 'audio/wav', 'audio/ogg', 'audio/webm']
	},
	chunkSize: 1024 * 1024, // 1MB chunks for large files
	uploadTimeout: 30000 // 30 seconds
};

/**
 * Validate file for upload
 * @param {File} file - File to validate
 * @param {Object} fieldConfig - Field configuration
 * @returns {Object} Validation result
 */
export function validateFile(file, fieldConfig = {}) {
	const props = fieldConfig.properties || {};
	const maxSize = (props.maxSize || 10) * 1024 * 1024; // Convert MB to bytes
	const acceptedTypes = props.accept ? props.accept.split(',').map(t => t.trim()) : ['*'];
	
	// Check file size
	if (file.size > maxSize) {
		return {
			valid: false,
			error: `File size exceeds maximum allowed (${props.maxSize || 10}MB)`
		};
	}
	
	// Check file type
	if (!acceptedTypes.includes('*')) {
		const isValidType = acceptedTypes.some(type => {
			// Handle wildcards like image/*
			if (type.endsWith('/*')) {
				const category = type.replace('/*', '');
				return file.type.startsWith(category + '/');
			}
			return file.type === type || file.name.endsWith(type.replace('.', ''));
		});
		
		if (!isValidType) {
			return {
				valid: false,
				error: `Invalid file type. Accepted: ${acceptedTypes.join(', ')}`
			};
		}
	}
	
	return { valid: true, error: null };
}

/**
 * Format file size for display
 * @param {number} bytes - File size in bytes
 * @returns {string} Formatted size
 */
export function formatFileSize(bytes) {
	if (bytes === 0) return '0 Bytes';
	
	const k = 1024;
	const sizes = ['Bytes', 'KB', 'MB', 'GB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	
	return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * Create file preview URL
 * @param {File} file 
 * @returns {Promise<string>} Preview URL
 */
export function createFilePreview(file) {
	return new Promise((resolve, reject) => {
		if (!file.type.startsWith('image/')) {
			resolve(null); // Only images get previews
			return;
		}
		
		const reader = new FileReader();
		reader.onload = (e) => resolve(e.target.result);
		reader.onerror = () => reject(new Error('Failed to create preview'));
		reader.readAsDataURL(file);
	});
}

/**
 * Upload file to server (simulated for demo)
 * @param {File} file - File to upload
 * @param {Function} onProgress - Progress callback
 * @returns {Promise<Object>} Upload result
 */
export async function uploadFile(file, onProgress = null) {
	// Simulate upload progress
	const totalChunks = Math.ceil(file.size / uploadConfig.chunkSize);
	
	for (let i = 0; i < totalChunks; i++) {
		// Simulate chunk upload
		await new Promise(resolve => setTimeout(resolve, 100));
		
		const progress = ((i + 1) / totalChunks) * 100;
		if (onProgress) {
			onProgress(progress);
		}
	}
	
	// Simulate server response
	const fileId = `file_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
	
	// Store file metadata in localStorage (demo mode)
	const fileData = {
		id: fileId,
		name: file.name,
		size: file.size,
		type: file.type,
		uploadedAt: new Date().toISOString(),
		// In production, this would be a server URL
		url: URL.createObjectURL(file)
	};
	
	const uploads = JSON.parse(localStorage.getItem('form_uploads') || '[]');
	uploads.push(fileData);
	localStorage.setItem('form_uploads', JSON.stringify(uploads));
	
	return {
		success: true,
		fileId: fileId,
		fileName: file.name,
		fileSize: file.size,
		fileType: file.type,
		url: fileData.url
	};
}

/**
 * Upload multiple files
 * @param {FileList} files - Files to upload
 * @param {Object} fieldConfig - Field configuration
 * @param {Function} onProgress - Progress callback for each file
 * @returns {Promise<Array>} Upload results
 */
export async function uploadMultipleFiles(files, fieldConfig = {}, onProgress = null) {
	const results = [];
	const maxFiles = fieldConfig.properties?.maxFiles || uploadConfig.maxFiles;
	
	// Check max files limit
	if (files.length > maxFiles) {
		return [{
			success: false,
			error: `Maximum ${maxFiles} files allowed`
		}];
	}
	
	for (let i = 0; i < files.length; i++) {
		const file = files[i];
		
		// Validate file
		const validation = validateFile(file, fieldConfig);
		if (!validation.valid) {
			results.push({
				success: false,
				fileName: file.name,
				error: validation.error
			});
			continue;
		}
		
		try {
			const fileProgress = (progress) => {
				if (onProgress) {
					onProgress(i, files.length, progress);
				}
			};
			
			const result = await uploadFile(file, fileProgress);
			results.push(result);
		} catch (error) {
			results.push({
				success: false,
				fileName: file.name,
				error: error.message
			});
		}
	}
	
	return results;
}

/**
 * Delete uploaded file
 * @param {string} fileId - File ID to delete
 * @returns {boolean} Success status
 */
export function deleteUploadedFile(fileId) {
	try {
		const uploads = JSON.parse(localStorage.getItem('form_uploads') || '[]');
		const file = uploads.find(u => u.id === fileId);
		
		if (file && file.url) {
			URL.revokeObjectURL(file.url);
		}
		
		const filtered = uploads.filter(u => u.id !== fileId);
		localStorage.setItem('form_uploads', JSON.stringify(filtered));
		
		return true;
	} catch {
		return false;
	}
}

/**
 * Get uploaded files
 * @returns {Array} Array of uploaded file metadata
 */
export function getUploadedFiles() {
	return JSON.parse(localStorage.getItem('form_uploads') || '[]');
}

/**
 * Create drag and drop handlers
 * @param {Function} onFilesDropped - Callback when files are dropped
 * @returns {Object} Drag and drop event handlers
 */
export function createDragAndDropHandlers(onFilesDropped) {
	return {
		onDragOver: (e) => {
			e.preventDefault();
			e.stopPropagation();
			e.currentTarget.classList.add('drag-over');
		},
		
		onDragLeave: (e) => {
			e.preventDefault();
			e.stopPropagation();
			e.currentTarget.classList.remove('drag-over');
		},
		
		onDrop: (e) => {
			e.preventDefault();
			e.stopPropagation();
			e.currentTarget.classList.remove('drag-over');
			
			const files = e.dataTransfer.files;
			if (files.length > 0) {
				onFilesDropped(files);
			}
		}
	};
}

/**
 * File type icons mapping
 */
export const fileTypeIcons = {
	'image/': '🖼️',
	'video/': '🎥',
	'audio/': '🎵',
	'application/pdf': '📄',
	'application/msword': '📝',
	'application/vnd.openxmlformats-officedocument.wordprocessingml.document': '📝',
	'application/vnd.ms-excel': '📊',
	'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': '📊',
	'text/': '📃',
	'application/zip': '📦',
	'default': '📎'
};

/**
 * Get icon for file type
 * @param {string} mimeType 
 * @returns {string} Icon emoji
 */
export function getFileIcon(mimeType) {
	for (const [type, icon] of Object.entries(fileTypeIcons)) {
		if (mimeType.startsWith(type.replace('/', '')) || mimeType === type) {
			return icon;
		}
	}
	return fileTypeIcons.default;
}

/**
 * Check if file is image
 * @param {string} mimeType 
 * @returns {boolean}
 */
export function isImageFile(mimeType) {
	return mimeType.startsWith('image/');
}

/**
 * Check if file is video
 * @param {string} mimeType 
 * @returns {boolean}
 */
export function isVideoFile(mimeType) {
	return mimeType.startsWith('video/');
}

/**
 * Check if file is audio
 * @param {string} mimeType 
 * @returns {boolean}
 */
export function isAudioFile(mimeType) {
	return mimeType.startsWith('audio/');
}

export default {
	uploadConfig,
	validateFile,
	formatFileSize,
	createFilePreview,
	uploadFile,
	uploadMultipleFiles,
	deleteUploadedFile,
	getUploadedFiles,
	createDragAndDropHandlers,
	fileTypeIcons,
	getFileIcon,
	isImageFile,
	isVideoFile,
	isAudioFile
};
