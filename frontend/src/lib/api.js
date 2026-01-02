// API client utilities for TPT Titan frontend

const API_BASE_URL = import.meta.env.DEV ? 'http://localhost:8080/api/v1' : '/api/v1';

/**
 * Generic API request function
 * @param {string} endpoint - API endpoint (without base URL)
 * @param {Object} options - Fetch options
 * @returns {Promise<Object>} Response data
 */
export async function apiRequest(endpoint, options = {}) {
	const url = `${API_BASE_URL}${endpoint}`;

	const defaultOptions = {
		headers: {
			'Content-Type': 'application/json',
		},
	};

	// Add authorization header if token exists
	const token = getAuthToken();
	if (token) {
		defaultOptions.headers['Authorization'] = `Bearer ${token}`;
	}

	const response = await fetch(url, {
		...defaultOptions,
		...options,
		headers: {
			...defaultOptions.headers,
			...options.headers,
		},
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({ message: 'Network error' }));
		throw new Error(error.message || `HTTP ${response.status}`);
	}

	return response.json();
}

/**
 * GET request
 * @param {string} endpoint
 * @returns {Promise<Object>}
 */
export function apiGet(endpoint) {
	return apiRequest(endpoint);
}

/**
 * POST request
 * @param {string} endpoint
 * @param {Object} data
 * @returns {Promise<Object>}
 */
export function apiPost(endpoint, data) {
	return apiRequest(endpoint, {
		method: 'POST',
		body: JSON.stringify(data),
	});
}

/**
 * PUT request
 * @param {string} endpoint
 * @param {Object} data
 * @returns {Promise<Object>}
 */
export function apiPut(endpoint, data) {
	return apiRequest(endpoint, {
		method: 'PUT',
		body: JSON.stringify(data),
	});
}

/**
 * DELETE request
 * @param {string} endpoint
 * @returns {Promise<Object>}
 */
export function apiDelete(endpoint) {
	return apiRequest(endpoint, {
		method: 'DELETE',
	});
}

/**
 * Get stored authentication token
 * @returns {string|null}
 */
export function getAuthToken() {
	if (typeof window !== 'undefined') {
		return localStorage.getItem('auth_token');
	}
	return null;
}

/**
 * Set authentication token
 * @param {string} token
 */
export function setAuthToken(token) {
	if (typeof window !== 'undefined') {
		localStorage.setItem('auth_token', token);
	}
}

/**
 * Clear authentication token
 */
export function clearAuthToken() {
	if (typeof window !== 'undefined') {
		localStorage.removeItem('auth_token');
	}
}

/**
 * Health check
 * @returns {Promise<Object>}
 */
export function healthCheck() {
	return apiGet('/health');
}

/**
 * Document API functions
 */

/**
 * Get all documents for the authenticated user
 * @returns {Promise<Object>} Documents list
 */
export function getDocuments() {
	return apiGet('/documents');
}

/**
 * Get a specific document
 * @param {string} documentId
 * @returns {Promise<Object>} Document data
 */
export function getDocument(documentId) {
	return apiGet(`/documents/${documentId}`);
}

/**
 * Create a new document
 * @param {Object} documentData
 * @param {string} documentData.title
 * @param {string} documentData.content_type
 * @param {Object} documentData.content
 * @returns {Promise<Object>} Created document
 */
export function createDocument(documentData) {
	return apiPost('/documents', documentData);
}

/**
 * Update an existing document
 * @param {string} documentId
 * @param {Object} documentData
 * @returns {Promise<Object>} Updated document
 */
export function updateDocument(documentId, documentData) {
	return apiPut(`/documents/${documentId}`, documentData);
}

/**
 * Delete a document
 * @param {string} documentId
 * @returns {Promise<Object>} Deletion confirmation
 */
export function deleteDocument(documentId) {
	return apiDelete(`/documents/${documentId}`);
}

/**
 * Get version history for a document
 * @param {string} documentId
 * @returns {Promise<Object>} Version history
 */
export function getDocumentVersions(documentId) {
	return apiGet(`/documents/${documentId}/versions`);
}

/**
 * Restore a specific version of a document
 * @param {string} documentId
 * @param {number} version
 * @returns {Promise<Object>} Restoration confirmation
 */
export function restoreDocumentVersion(documentId, version) {
	return apiPost(`/documents/${documentId}/versions/${version}/restore`, {});
}

/**
 * AI API functions
 */

/**
 * Get all available AI models for the user
 * @returns {Promise<Object>} Available models
 */
export function getAIModels() {
	return apiGet('/ai/models');
}

/**
 * Create a custom AI model
 * @param {Object} modelData
 * @returns {Promise<Object>} Created model
 */
export function createAIModel(modelData) {
	return apiPost('/ai/models', modelData);
}

/**
 * Process an AI request
 * @param {Object} requestData
 * @param {string} requestData.task_id
 * @param {string} requestData.model_id
 * @param {string} requestData.input
 * @param {string} requestData.input_type
 * @returns {Promise<Object>} Request status
 */
export function processAIRequest(requestData) {
	return apiPost('/ai/requests', requestData);
}

/**
 * Get AI request status
 * @param {string} requestId
 * @returns {Promise<Object>} Request status and result
 */
export function getAIRequestStatus(requestId) {
	return apiGet(`/ai/requests/${requestId}`);
}

/**
 * List available Ollama models
 * @returns {Promise<Object>} Ollama models
 */
export function listOllamaModels() {
	return apiGet('/ai/ollama/models');
}

/**
 * Pull an Ollama model
 * @param {string} modelName
 * @returns {Promise<Object>} Pull status
 */
export function pullOllamaModel(modelName) {
	return apiPost(`/ai/ollama/models/${modelName}/pull`, {});
}

/**
 * Get AI usage statistics
 * @returns {Promise<Object>} Usage data
 */
export function getAIUsage() {
	return apiGet('/ai/usage');
}

/**
 * Detect system hardware capabilities
 * @returns {Promise<Object>} Hardware information
 */
export function detectHardware() {
	return apiGet('/ai/hardware');
}

/**
 * Get AI model recommendations based on hardware
 * @returns {Promise<Object>} Hardware and model recommendations
 */
export function getRecommendedModels() {
	return apiGet('/ai/recommendations');
}

/**
 * Setup recommended models and tasks for the user
 * @returns {Promise<Object>} Setup confirmation
 */
export function setupRecommendedModels() {
	return apiPost('/ai/setup', {});
}

/**
 * Check for available AI model upgrades
 * @returns {Promise<Object>} Upgrade check results with available options
 */
export function checkForAIUpgrades() {
	return apiPost('/ai/upgrades/check', {});
}

/**
 * Get AI upgrade check history
 * @returns {Promise<Object>} List of past upgrade checks
 */
export function getUpgradeHistory() {
	return apiGet('/ai/upgrades/history');
}

/**
 * Apply a selected AI model upgrade
 * @param {string} upgradeId - The ID of the upgrade check
 * @returns {Promise<Object>} Upgrade application status
 */
export function applyAIUpgrade(upgradeId) {
	return apiPost('/ai/upgrades/apply', { upgrade_id: upgradeId });
}

/**
 * Format API error for display
 * @param {Error} error
 * @returns {string}
 */
export function formatApiError(error) {
	if (error.message.includes('Network error')) {
		return 'Unable to connect to server. Please check your connection.';
	}
	return error.message;
}
