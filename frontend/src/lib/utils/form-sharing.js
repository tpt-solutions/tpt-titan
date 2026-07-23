// @ts-nocheck
/**
 * Form Sharing Utility - Generate public links, embed codes, QR codes
 */

/**
 * Generate a unique public form ID
 * @returns {string} Unique ID
 */
export function generateFormId() {
	return 'form_' + Date.now().toString(36) + '_' + Math.random().toString(36).substr(2, 9);
}

/**
 * Generate public form URL
 * @param {string} formId - Form ID
 * @param {Object} options - URL options
 * @returns {string} Public URL
 */
export function generatePublicUrl(formId, options = {}) {
	const { customDomain = false, subdomain = null } = options;
	
	// In production, this would use the actual domain
	const baseUrl = customDomain && subdomain 
		? `https://${subdomain}.tpt-titan.com`
		: window.location.origin;
	
	return `${baseUrl}/forms/public/${formId}`;
}

/**
 * Generate embed code for external sites
 * @param {string} formId - Form ID
 * @param {Object} options - Embed options
 * @returns {string} HTML embed code
 */
export function generateEmbedCode(formId, options = {}) {
	const {
		width = '100%',
		height = '600px',
		theme = 'light',
		hideHeader = false,
		hideFooter = false
	} = options;
	
	const url = generatePublicUrl(formId, options);
	const params = new URLSearchParams({
		theme,
		hideHeader: hideHeader ? '1' : '0',
		hideFooter: hideFooter ? '1' : '0',
		embed: '1'
	});
	
	const embedUrl = `${url}?${params.toString()}`;
	
	return `<iframe 
	src="${embedUrl}" 
	width="${width}" 
	height="${height}" 
	frameborder="0" 
	style="border: none; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"
	allow="camera; microphone; geolocation"
	title="TPT Titan Form"
></iframe>`;
}

/**
 * Generate QR code URL (using external service)
 * @param {string} formUrl - Form URL
 * @param {Object} options - QR code options
 * @returns {string} QR code image URL
 */
export function generateQRCodeUrl(formUrl, options = {}) {
	const {
		size = 300,
		color = '000000',
		bgColor = 'FFFFFF',
		errorCorrection = 'H'
	} = options;
	
	// Using QRServer API (free, no key required for basic usage)
	const params = new URLSearchParams({
		data: formUrl,
		size: `${size}x${size}`,
		color: color,
		bgcolor: bgColor,
		qzone: '1',
		format: 'png',
		eclevel: errorCorrection
	});
	
	return `https://api.qrserver.com/v1/create-qr-code/?${params.toString()}`;
}

/**
 * Generate shareable links for different platforms
 * @param {string} formUrl - Form URL
 * @param {Object} formData - Form metadata
 * @returns {Object} Platform-specific share URLs
 */
export function generateSocialShareLinks(formUrl, formData = {}) {
	const { name = 'Form', description = '' } = formData;
	const text = encodeURIComponent(`Check out this form: ${name}`);
	const textWithDesc = encodeURIComponent(`${description || name}`);
	
	return {
		email: `mailto:?subject=${encodeURIComponent(name)}&body=${encodeURIComponent(formUrl)}`,
		facebook: `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(formUrl)}`,
		twitter: `https://twitter.com/intent/tweet?url=${encodeURIComponent(formUrl)}&text=${text}`,
		linkedin: `https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(formUrl)}`,
		whatsapp: `https://wa.me/?text=${text}%20${encodeURIComponent(formUrl)}`,
		telegram: `https://t.me/share/url?url=${encodeURIComponent(formUrl)}&text=${text}`,
		reddit: `https://reddit.com/submit?url=${encodeURIComponent(formUrl)}&title=${text}`,
		copy: formUrl
	};
}

/**
 * Form visibility settings
 */
export const visibilityOptions = {
	public: {
		id: 'public',
		name: 'Public',
		description: 'Anyone with the link can view and submit',
		icon: '🌐'
	},
	private: {
		id: 'private',
		name: 'Private',
		description: 'Only you can access',
		icon: '🔒'
	},
	restricted: {
		id: 'restricted',
		name: 'Restricted',
		description: 'Only specific people can access',
		icon: '👥'
	},
	password: {
		id: 'password',
		name: 'Password Protected',
		description: 'Requires password to access',
		icon: '🔐'
	}
};

/**
 * Generate form access settings
 * @param {Object} settings - Access configuration
 * @returns {Object} Normalized settings
 */
export function generateAccessSettings(settings = {}) {
	return {
		visibility: settings.visibility || 'public',
		password: settings.password || null,
		allowedEmails: settings.allowedEmails || [],
		requireLogin: settings.requireLogin || false,
		allowMultipleSubmissions: settings.allowMultipleSubmissions !== false,
		expiryDate: settings.expiryDate || null,
		startDate: settings.startDate || null,
		maxSubmissions: settings.maxSubmissions || null,
		captcha: settings.captcha || false
	};
}

/**
 * Create form sharing configuration
 * @param {string} formId - Form ID
 * @param {Object} config - Sharing configuration
 * @returns {Object} Complete sharing config
 */
export function createSharingConfig(formId, config = {}) {
	const accessSettings = generateAccessSettings(config);
	
	return {
		formId,
		publicUrl: generatePublicUrl(formId, config),
		embedCode: generateEmbedCode(formId, config),
		qrCodeUrl: generateQRCodeUrl(generatePublicUrl(formId), config),
		socialLinks: generateSocialShareLinks(generatePublicUrl(formId), config),
		accessSettings,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
		analyticsEnabled: config.analyticsEnabled !== false,
		notificationSettings: {
			emailOnSubmit: config.emailOnSubmit || false,
			webhookUrl: config.webhookUrl || null
		}
	};
}

/**
 * Validate form access
 * @param {Object} accessSettings - Access configuration
 * @param {Object} context - Access context (user, password, etc.)
 * @returns {Object} Validation result
 */
export function validateAccess(accessSettings, context = {}) {
	const { user = null, password = null, submissionCount = 0 } = context;
	
	// Check visibility
	if (accessSettings.visibility === 'private') {
		return { valid: false, error: 'This form is private' };
	}
	
	// Check password
	if (accessSettings.visibility === 'password' && accessSettings.password) {
		if (password !== accessSettings.password) {
			return { valid: false, error: 'Invalid password' };
		}
	}
	
	// Check allowed emails
	if (accessSettings.visibility === 'restricted' && accessSettings.allowedEmails.length > 0) {
		if (!user || !accessSettings.allowedEmails.includes(user.email)) {
			return { valid: false, error: 'You do not have access to this form' };
		}
	}
	
	// Check login requirement
	if (accessSettings.requireLogin && !user) {
		return { valid: false, error: 'Please log in to access this form' };
	}
	
	// Check dates
	const now = new Date();
	if (accessSettings.startDate && new Date(accessSettings.startDate) > now) {
		return { valid: false, error: 'This form is not yet open' };
	}
	if (accessSettings.expiryDate && new Date(accessSettings.expiryDate) < now) {
		return { valid: false, error: 'This form has closed' };
	}
	
	// Check max submissions
	if (accessSettings.maxSubmissions && submissionCount >= accessSettings.maxSubmissions) {
		return { valid: false, error: 'This form has reached its submission limit' };
	}
	
	return { valid: true, error: null };
}

/**
 * Generate form short link (using external service or custom)
 * @param {string} formId - Form ID
 * @returns {Promise<string>} Short URL
 */
export async function generateShortLink(formId) {
	// In production, this would use a URL shortener API
	// For demo, we'll create a simple hash-based short code
	const shortCode = btoa(formId).replace(/[^a-zA-Z0-9]/g, '').substr(0, 8);
	return `${window.location.origin}/f/${shortCode}`;
}

/**
 * Track form view (analytics)
 * @param {string} formId - Form ID
 * @param {Object} metadata - View metadata
 */
export function trackFormView(formId, metadata = {}) {
	const views = JSON.parse(localStorage.getItem('form_analytics_views') || '[]');
	views.push({
		formId,
		timestamp: new Date().toISOString(),
		userAgent: navigator.userAgent,
		referrer: document.referrer,
		...metadata
	});
	localStorage.setItem('form_analytics_views', JSON.stringify(views));
}

/**
 * Track form submission (analytics)
 * @param {string} formId - Form ID
 * @param {Object} metadata - Submission metadata
 */
export function trackFormSubmission(formId, metadata = {}) {
	const submissions = JSON.parse(localStorage.getItem('form_analytics_submissions') || '[]');
	submissions.push({
		formId,
		timestamp: new Date().toISOString(),
		...metadata
	});
	localStorage.setItem('form_analytics_submissions', JSON.stringify(submissions));
}

/**
 * Get form analytics
 * @param {string} formId - Form ID
 * @returns {Object} Analytics data
 */
export function getFormAnalytics(formId) {
	const views = JSON.parse(localStorage.getItem('form_analytics_views') || '[]')
		.filter(v => v.formId === formId);
	const submissions = JSON.parse(localStorage.getItem('form_analytics_submissions') || '[]')
		.filter(s => s.formId === formId);
	
	const uniqueViews = new Set(views.map(v => v.sessionId || v.timestamp)).size;
	
	return {
		totalViews: views.length,
		uniqueViews,
		totalSubmissions: submissions.length,
		conversionRate: views.length > 0 ? (submissions.length / views.length * 100).toFixed(2) : 0,
		views,
		submissions
	};
}

/**
 * Export sharing configuration
 * @param {Object} sharingConfig 
 * @returns {string} JSON string
 */
export function exportSharingConfig(sharingConfig) {
	return JSON.stringify(sharingConfig, null, 2);
}

/**
 * Import sharing configuration
 * @param {string} jsonString 
 * @returns {Object} Sharing config
 */
export function importSharingConfig(jsonString) {
	try {
		return JSON.parse(jsonString);
	} catch {
		return null;
	}
}

export default {
	generateFormId,
	generatePublicUrl,
	generateEmbedCode,
	generateQRCodeUrl,
	generateSocialShareLinks,
	generateAccessSettings,
	createSharingConfig,
	validateAccess,
	generateShortLink,
	trackFormView,
	trackFormSubmission,
	getFormAnalytics,
	exportSharingConfig,
	importSharingConfig,
	visibilityOptions
};
