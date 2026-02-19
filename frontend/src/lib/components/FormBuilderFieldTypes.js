// frontend/src/lib/components/FormBuilderFieldTypes.js
// Centralised field-type catalogue used by FormBuilder and its sub-components.

export const fieldTypes = [
	// Basic Text Fields
	{ id: 'text',     name: 'Text Input',  icon: '📝', properties: { placeholder: '', required: false, maxLength: 100 } },
	{ id: 'textarea', name: 'Long Text',   icon: '📄', properties: { placeholder: '', required: false, rows: 3 } },
	{ id: 'password', name: 'Password',    icon: '🔒', properties: { placeholder: '', required: false } },

	// Numbers & Calculations
	{ id: 'number',      name: 'Number',           icon: '🔢', properties: { placeholder: '', required: false, min: null, max: null } },
	{ id: 'currency',    name: 'Currency',          icon: '💰', properties: { placeholder: '', required: false, currency: 'USD', min: null, max: null } },
	{ id: 'percentage',  name: 'Percentage',        icon: '📊', properties: { placeholder: '', required: false, min: 0, max: 100 } },
	{ id: 'calculation', name: 'Calculated Field',  icon: '🧮', properties: { formula: '', required: false, decimalPlaces: 2 } },

	// Contact Information
	{ id: 'email', name: 'Email',       icon: '📧', properties: { placeholder: '', required: false } },
	{ id: 'phone', name: 'Phone',       icon: '📱', properties: { placeholder: '', required: false } },
	{ id: 'url',   name: 'Website URL', icon: '🌐', properties: { placeholder: 'https://', required: false } },

	// Date & Time
	{ id: 'date',           name: 'Date',        icon: '📅',    properties: { required: false } },
	{ id: 'time',           name: 'Time',        icon: '🕐',    properties: { required: false } },
	{ id: 'datetime-local', name: 'Date & Time', icon: '📅🕐', properties: { required: false } },

	// Selection Fields
	{ id: 'select',   name: 'Dropdown',      icon: '▼',    properties: { options: ['Option 1', 'Option 2'], required: false } },
	{ id: 'radio',    name: 'Radio Buttons', icon: '⭕',   properties: { options: ['Option 1', 'Option 2'], required: false } },
	{ id: 'checkbox', name: 'Checkboxes',    icon: '☑️',  properties: { options: ['Option 1', 'Option 2'], required: false } },
	{ id: 'yesno',    name: 'Yes/No',        icon: '👍👎', properties: { required: false, defaultValue: null } },

	// Advanced Selection
	{ id: 'rating', name: 'Rating',        icon: '⭐', properties: { maxRating: 5, required: false } },
	{ id: 'scale',  name: 'Scale',         icon: '📏', properties: { min: 1, max: 10, minLabel: 'Poor', maxLabel: 'Excellent', required: false } },
	{ id: 'matrix', name: 'Rating Matrix', icon: '📋', properties: { rows: ['Feature 1', 'Feature 2'], columns: ['Poor', 'Good', 'Excellent'], required: false } },

	// Files & Media
	{ id: 'file',         name: 'File Upload',   icon: '📎', properties: { accept: '*', maxSize: 10, required: false } },
	{ id: 'image-upload', name: 'Image Upload',  icon: '🖼️', properties: { accept: 'image/*', maxSize: 5, maxFiles: 3, required: false } },
	{ id: 'video-upload', name: 'Video Upload',  icon: '🎥', properties: { accept: 'video/*', maxSize: 100, required: false } },
	{ id: 'audio-upload', name: 'Audio Upload',  icon: '🎵', properties: { accept: 'audio/*', maxSize: 50, required: false } },

	// Special Fields
	{ id: 'signature',   name: 'Digital Signature', icon: '✍️',   properties: { required: false, legalText: '' } },
	{ id: 'geolocation', name: 'GPS Location',      icon: '📍',   properties: { required: false, accuracy: 'high' } },
	{ id: 'qr-code',     name: 'QR Code Input',     icon: '📱📷', properties: { required: false, format: 'text' } },
	{ id: 'barcode',     name: 'Barcode Scanner',   icon: '📊📷', properties: { required: false, format: 'code128' } },

	// Advanced Fields
	{ id: 'address', name: 'Address',      icon: '🏠',  properties: { required: false, includeCountry: true, format: 'single' } },
	{ id: 'table',   name: 'Data Table',   icon: '📊',  properties: { columns: ['Column 1', 'Column 2'], rows: 3, required: false } },
	{ id: 'range',   name: 'Range Slider', icon: '🎚️', properties: { min: 0, max: 100, step: 1, required: false } },
	{ id: 'color',   name: 'Color Picker', icon: '🎨',  properties: { required: false, defaultColor: '#000000' } },

	// Utility Fields
	{ id: 'hidden',   name: 'Hidden Field',    icon: '👁️‍🗨️', properties: { value: '', required: false } },
	{ id: 'html',     name: 'HTML Content',    icon: '⚡',    properties: { content: '<p>Custom HTML content</p>', required: false } },
	{ id: 'divider',  name: 'Section Divider', icon: '➖',    properties: { title: 'Section Title', required: false } }
];

/** Returns a fresh field object for the given fieldType definition. */
export function createField(fieldType, order = 0) {
	return {
		id: Date.now() + Math.random(),
		type: fieldType.id,
		label: fieldType.name,
		properties: { ...fieldType.properties },
		order
	};
}
