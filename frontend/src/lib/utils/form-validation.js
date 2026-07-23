// @ts-nocheck
/**
 * Form Validation Utility - Validate form fields and show errors
 */

/**
 * Validation result type
 * @typedef {Object} ValidationResult
 * @property {boolean} valid - Whether the field is valid
 * @property {string} error - Error message if invalid
 */

/**
 * Validate a single field value
 * @param {any} value - Field value
 * @param {Object} field - Field definition with properties
 * @returns {ValidationResult} Validation result
 */
export function validateField(value, field) {
	const props = field.properties || {};
	
	// Required check
	if (props.required && (value === undefined || value === null || value === '')) {
		return { valid: false, error: 'This field is required' };
	}
	
	// Skip other validations if value is empty and not required
	if (!value && !props.required) {
		return { valid: true, error: null };
	}
	
	// Type-specific validations
	switch (field.type) {
		case 'email':
			if (!isValidEmail(value)) {
				return { valid: false, error: 'Please enter a valid email address' };
			}
			break;
			
		case 'url':
			if (!isValidURL(value)) {
				return { valid: false, error: 'Please enter a valid URL' };
			}
			break;
			
		case 'phone':
			if (!isValidPhone(value)) {
				return { valid: false, error: 'Please enter a valid phone number' };
			}
			break;
			
		case 'number':
		case 'currency':
		case 'percentage':
			const numValue = parseFloat(value);
			if (isNaN(numValue)) {
				return { valid: false, error: 'Please enter a valid number' };
			}
			if (props.min !== null && props.min !== undefined && numValue < props.min) {
				return { valid: false, error: `Value must be at least ${props.min}` };
			}
			if (props.max !== null && props.max !== undefined && numValue > props.max) {
				return { valid: false, error: `Value must be at most ${props.max}` };
			}
			break;
			
		case 'text':
		case 'textarea':
		case 'password':
			const strValue = String(value);
			if (props.maxLength && strValue.length > props.maxLength) {
				return { valid: false, error: `Maximum ${props.maxLength} characters allowed` };
			}
			if (props.minLength && strValue.length < props.minLength) {
				return { valid: false, error: `Minimum ${props.minLength} characters required` };
			}
			if (props.pattern && !new RegExp(props.pattern).test(strValue)) {
				return { valid: false, error: props.patternMessage || 'Invalid format' };
			}
			break;
			
		case 'file':
		case 'image-upload':
		case 'video-upload':
		case 'audio-upload':
			if (value && props.maxSize) {
				const fileSizeMB = value.size / (1024 * 1024);
				if (fileSizeMB > props.maxSize) {
					return { valid: false, error: `File size must be less than ${props.maxSize}MB` };
				}
			}
			if (value && props.accept && props.accept !== '*') {
				const acceptedTypes = props.accept.split(',').map(t => t.trim());
				const fileType = value.type;
				const isAccepted = acceptedTypes.some(type => {
					if (type.includes('*')) {
						return fileType.startsWith(type.replace('/*', ''));
					}
					return fileType === type;
				});
				if (!isAccepted) {
					return { valid: false, error: `File type must be: ${props.accept}` };
				}
			}
			break;
			
		case 'select':
		case 'radio':
			const options = props.options || [];
			if (!options.includes(value)) {
				return { valid: false, error: 'Please select a valid option' };
			}
			break;
			
		case 'checkbox':
			const checkboxOptions = props.options || [];
			if (Array.isArray(value)) {
				const invalid = value.some(v => !checkboxOptions.includes(v));
				if (invalid) {
					return { valid: false, error: 'Invalid selection' };
				}
			}
			break;
			
		case 'rating':
			const rating = parseInt(value);
			const maxRating = props.maxRating || 5;
			if (isNaN(rating) || rating < 1 || rating > maxRating) {
				return { valid: false, error: `Please rate between 1 and ${maxRating}` };
			}
			break;
			
		case 'scale':
			const scaleValue = parseInt(value);
			if (isNaN(scaleValue) || scaleValue < props.min || scaleValue > props.max) {
				return { valid: false, error: `Value must be between ${props.min} and ${props.max}` };
			}
			break;
			
		case 'date':
			if (!isValidDate(value)) {
				return { valid: false, error: 'Please enter a valid date' };
			}
			break;
			
		case 'time':
			if (!isValidTime(value)) {
				return { valid: false, error: 'Please enter a valid time' };
			}
			break;
			
		case 'datetime-local':
			if (!isValidDateTime(value)) {
				return { valid: false, error: 'Please enter a valid date and time' };
			}
			break;
			
		case 'signature':
			if (props.required && (!value || value === '')) {
				return { valid: false, error: 'Signature is required' };
			}
			break;
			
		case 'calculation':
			// Calculated fields are always valid (they're computed)
			return { valid: true, error: null };
	}
	
	return { valid: true, error: null };
}

/**
 * Validate all form fields
 * @param {Object} values - Form values keyed by field ID
 * @param {Array} fields - Array of field definitions
 * @returns {Object} Validation results keyed by field ID
 */
export function validateForm(values, fields) {
	const results = {};
	let isValid = true;
	
	for (const field of fields) {
		const value = values[field.id];
		const result = validateField(value, field);
		results[field.id] = result;
		if (!result.valid) {
			isValid = false;
		}
	}
	
	return { valid: isValid, results };
}

/**
 * Get default value for a field type
 * @param {Object} field - Field definition
 * @returns {any} Default value
 */
export function getDefaultValue(field) {
	const props = field.properties || {};
	
	switch (field.type) {
		case 'checkbox':
			return [];
		case 'yesno':
			return props.defaultValue || null;
		case 'rating':
			return 0;
		case 'scale':
			return props.min || 1;
		case 'range':
			return props.min || 0;
		case 'color':
			return props.defaultColor || '#000000';
		case 'number':
		case 'currency':
		case 'percentage':
			return props.min || '';
		case 'calculation':
			return '';
		case 'hidden':
			return props.value || '';
		case 'table':
			return props.rows ? props.rows.map(() => props.columns.map(() => '')) : [];
		default:
			return '';
	}
}

/**
 * Check if email is valid
 * @param {string} email 
 * @returns {boolean}
 */
function isValidEmail(email) {
	const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	return re.test(String(email).toLowerCase());
}

/**
 * Check if URL is valid
 * @param {string} url 
 * @returns {boolean}
 */
function isValidURL(url) {
	try {
		new URL(url);
		return true;
	} catch {
		return false;
	}
}

/**
 * Check if phone number is valid
 * @param {string} phone 
 * @returns {boolean}
 */
function isValidPhone(phone) {
	// Basic phone validation - allows various formats
	const cleaned = String(phone).replace(/\D/g, '');
	return cleaned.length >= 10 && cleaned.length <= 15;
}

/**
 * Check if date is valid
 * @param {string} date 
 * @returns {boolean}
 */
function isValidDate(date) {
	const d = new Date(date);
	return !isNaN(d.getTime());
}

/**
 * Check if time is valid
 * @param {string} time 
 * @returns {boolean}
 */
function isValidTime(time) {
	const re = /^([01]?[0-9]|2[0-3]):[0-5][0-9]$/;
	return re.test(time);
}

/**
 * Check if datetime is valid
 * @param {string} datetime 
 * @returns {boolean}
 */
function isValidDateTime(datetime) {
	const d = new Date(datetime);
	return !isNaN(d.getTime());
}

/**
 * Format field value for display
 * @param {any} value 
 * @param {Object} field 
 * @returns {string}
 */
export function formatFieldValue(value, field) {
	if (value === undefined || value === null || value === '') {
		return '-';
	}
	
	const props = field.properties || {};
	
	switch (field.type) {
		case 'currency':
			const currency = props.currency || 'USD';
			return new Intl.NumberFormat('en-US', {
				style: 'currency',
				currency: currency
			}).format(value);
			
		case 'percentage':
			return `${value}%`;
			
		case 'date':
			return new Date(value).toLocaleDateString();
			
		case 'datetime-local':
			return new Date(value).toLocaleString();
			
		case 'checkbox':
			if (Array.isArray(value)) {
				return value.join(', ');
			}
			return String(value);
			
		case 'rating':
			return '⭐'.repeat(parseInt(value));
			
		case 'yesno':
			return value === true ? 'Yes' : value === false ? 'No' : '-';
			
		case 'file':
		case 'image-upload':
		case 'video-upload':
		case 'audio-upload':
			if (value && value.name) {
				return value.name;
			}
			return '-';
			
		case 'signature':
			return value ? '[Signed]' : '[Not signed]';
			
		case 'geolocation':
			if (value && typeof value === 'object') {
				return `${value.lat?.toFixed(4)}, ${value.lng?.toFixed(4)}`;
			}
			return String(value);
			
		default:
			return String(value);
	}
}

/**
 * Generate sample data for a field (for testing)
 * @param {Object} field 
 * @returns {any} Sample value
 */
export function generateSampleValue(field) {
	const props = field.properties || {};
	
	switch (field.type) {
		case 'text':
			return 'Sample text';
		case 'textarea':
			return 'This is a sample long text response for testing purposes.';
		case 'email':
			return 'test@example.com';
		case 'phone':
			return '+1 (555) 123-4567';
		case 'url':
			return 'https://example.com';
		case 'number':
			return Math.floor(Math.random() * 100);
		case 'currency':
			return (Math.random() * 1000).toFixed(2);
		case 'percentage':
			return Math.floor(Math.random() * 100);
		case 'select':
		case 'radio':
			const options = props.options || ['Option 1', 'Option 2'];
			return options[0];
		case 'checkbox':
			const cbOptions = props.options || ['Option 1', 'Option 2'];
			return [cbOptions[0]];
		case 'yesno':
			return true;
		case 'rating':
			return Math.floor(Math.random() * (props.maxRating || 5)) + 1;
		case 'scale':
			return Math.floor(Math.random() * (props.max - props.min + 1)) + props.min;
		case 'date':
			return new Date().toISOString().split('T')[0];
		case 'time':
			return '14:30';
		case 'datetime-local':
			return new Date().toISOString().slice(0, 16);
		case 'color':
			return '#3b82f6';
		case 'range':
			return Math.floor((props.min + props.max) / 2);
		case 'file':
			return { name: 'sample.pdf', size: 1024000, type: 'application/pdf' };
		case 'image-upload':
			return { name: 'sample.jpg', size: 2048000, type: 'image/jpeg' };
		case 'signature':
			return 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==';
		case 'calculation':
			return '42';
		case 'hidden':
			return props.value || 'hidden-value';
		default:
			return 'Sample value';
	}
}

export default {
	validateField,
	validateForm,
	getDefaultValue,
	formatFieldValue,
	generateSampleValue
};
