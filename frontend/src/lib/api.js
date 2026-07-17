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
 * Spreadsheet API functions
 */

/**
 * Create a new spreadsheet
 * @param {Object} spreadsheetData
 * @param {string} spreadsheetData.name
 * @returns {Promise<Object>} Created spreadsheet
 */
export function createSpreadsheet(spreadsheetData) {
	return apiPost('/spreadsheets', spreadsheetData);
}

/**
 * Get a specific spreadsheet
 * @param {string} spreadsheetId
 * @returns {Promise<Object>} Spreadsheet data
 */
export function getSpreadsheet(spreadsheetId) {
	return apiGet(`/spreadsheets/${spreadsheetId}`);
}

/**
 * Update a spreadsheet cell
 * @param {string} spreadsheetId
 * @param {Object} cellData
 * @param {string} cellData.cell_reference
 * @param {any} cellData.value
 * @param {string} cellData.formula
 * @returns {Promise<Object>} Update confirmation
 */
export function updateSpreadsheetCell(spreadsheetId, cellData) {
	return apiPut(`/spreadsheets/${spreadsheetId}/cells`, cellData);
}

/**
 * Update multiple spreadsheet cells (batch)
 * @param {string} spreadsheetId
 * @param {Object} batchData
 * @param {Array} batchData.updates
 * @param {number} batchData.version
 * @returns {Promise<Object>} Batch update confirmation
 */
export function updateSpreadsheetCells(spreadsheetId, batchData) {
	return apiPut(`/spreadsheets/${spreadsheetId}/batch`, batchData);
}

/**
 * Evaluate a spreadsheet formula
 * @param {Object} formulaData
 * @param {string} formulaData.formula
 * @param {Object} formulaData.cell_context
 * @returns {Promise<Object>} Formula evaluation result
 */
export function evaluateFormula(formulaData) {
	return apiPost('/spreadsheets/evaluate', formulaData);
}

/**
 * Get available spreadsheet functions
 * @returns {Promise<Object>} List of available functions
 */
export function getAvailableFunctions() {
	return apiGet('/spreadsheets/functions');
}

/**
 * Validate a spreadsheet formula
 * @param {Object} validationData
 * @param {string} validationData.formula
 * @returns {Promise<Object>} Validation result
 */
export function validateFormula(validationData) {
	return apiPost('/spreadsheets/validate', validationData);
}

/**
 * Get chart suggestions for spreadsheet data
 * @param {Object} chartData
 * @param {Object} chartData.data
 * @param {string} chartData.range
 * @param {Object} chartData.data_types
 * @returns {Promise<Object>} Chart suggestions
 */
export function getChartSuggestions(chartData) {
	return apiPost('/spreadsheets/charts/suggestions', chartData);
}

/**
 * Create a chart from spreadsheet data
 * @param {Object} chartData
 * @returns {Promise<Object>} Created chart
 */
export function createChart(chartData) {
	return apiPost('/spreadsheets/charts', chartData);
}

/**
 * Get charts for a spreadsheet
 * @param {string} spreadsheetId
 * @returns {Promise<Object>} Charts list
 */
export function getSpreadsheetCharts(spreadsheetId) {
	return apiGet(`/spreadsheets/${spreadsheetId}/charts`);
}

/**
 * Export spreadsheet to Excel
 * @param {string} spreadsheetId
 * @param {Object} exportOptions
 * @returns {Promise<Blob>} Excel file
 */
export function exportSpreadsheetToExcel(spreadsheetId, exportOptions = {}) {
	const url = `${API_BASE_URL}/spreadsheets/${spreadsheetId}/export/excel`;

	const defaultOptions = {
		headers: {
			'Content-Type': 'application/json',
		},
	};

	const token = getAuthToken();
	if (token) {
		defaultOptions.headers['Authorization'] = `Bearer ${token}`;
	}

	return fetch(url, {
		method: 'POST',
		body: JSON.stringify(exportOptions),
		...defaultOptions,
		headers: defaultOptions.headers,
	}).then(response => {
		if (!response.ok) {
			throw new Error(`Export failed: ${response.status}`);
		}
		return response.blob();
	});
}

/**
 * Import Excel file to spreadsheet
 * @param {FormData} formData - Form data with file
 * @returns {Promise<Object>} Import result
 */
export function importExcelToSpreadsheet(formData) {
	const url = `${API_BASE_URL}/spreadsheets/import/excel`;

	const defaultOptions = {
		headers: {},
	};

	const token = getAuthToken();
	if (token) {
		defaultOptions.headers['Authorization'] = `Bearer ${token}`;
	}

	return fetch(url, {
		method: 'POST',
		body: formData,
		...defaultOptions,
		headers: defaultOptions.headers,
	}).then(response => {
		if (!response.ok) {
			const error = response.json().catch(() => ({ message: 'Import failed' }));
			throw new Error(error.message || `HTTP ${response.status}`);
		}
		return response.json();
	});
}

/**
 * Get spreadsheet version info
 * @param {string} spreadsheetId
 * @returns {Promise<Object>} Version information
 */
export function getSpreadsheetVersion(spreadsheetId) {
	return apiGet(`/spreadsheets/${spreadsheetId}/version`);
}

/**
 * Get spreadsheet changes since version
 * @param {string} spreadsheetId
 * @param {number} sinceVersion
 * @returns {Promise<Object>} Changes list
 */
export function getSpreadsheetChanges(spreadsheetId, sinceVersion) {
	return apiGet(`/spreadsheets/${spreadsheetId}/changes?since_version=${sinceVersion}`);
}

/**
 * Lock spreadsheet cells for editing
 * @param {string} spreadsheetId
 * @param {Object} lockData
 * @returns {Promise<Object>} Lock confirmation
 */
export function lockSpreadsheetCells(spreadsheetId, lockData) {
	return apiPost(`/spreadsheets/${spreadsheetId}/lock`, lockData);
}

/**
 * Unlock spreadsheet cells
 * @param {string} spreadsheetId
 * @param {Object} unlockData
 * @returns {Promise<Object>} Unlock confirmation
 */
export function unlockSpreadsheetCells(spreadsheetId, unlockData) {
	return apiPost(`/spreadsheets/${spreadsheetId}/unlock`, unlockData);
}

/**
 * Form API functions
 */

/**
 * Get all forms for the authenticated user
 * @returns {Promise<Object>} Forms list
 */
export function getForms() {
	return apiGet('/forms');
}

/**
 * Get a specific form
 * @param {string} formId
 * @returns {Promise<Object>} Form data
 */
export function getForm(formId) {
	return apiGet(`/forms/${formId}`);
}

/**
 * Create a new form
 * @param {Object} formData
 * @param {string} formData.name
 * @param {string} formData.description
 * @param {Array} formData.fields
 * @returns {Promise<Object>} Created form
 */
export function createForm(formData) {
	return apiPost('/forms', formData);
}

/**
 * Update an existing form
 * @param {string} formId
 * @param {Object} formData
 * @returns {Promise<Object>} Updated form
 */
export function updateForm(formId, formData) {
	return apiPut(`/forms/${formId}`, formData);
}

/**
 * Delete a form
 * @param {string} formId
 * @returns {Promise<Object>} Deletion confirmation
 */
export function deleteForm(formId) {
	return apiDelete(`/forms/${formId}`);
}

/**
 * Get form responses
 * @param {string} formId
 * @returns {Promise<Object>} Form responses
 */
export function getFormResponses(formId) {
	return apiGet(`/forms/${formId}/responses`);
}

/**
 * Submit a form response
 * @param {string} formId
 * @param {Object} responseData
 * @returns {Promise<Object>} Submission confirmation
 */
export function submitFormResponse(formId, responseData) {
	return apiPost(`/forms/${formId}/responses`, responseData);
}

/**
 * Advanced form modules (templates, relationships, reports, query builder,
 * email distribution, workflow) reachable via the formsAdvancedGroup routes.
 */

// ── Templates ───────────────────────────────────────────────────────
export function getFormTemplates(params = {}) {
	const qs = new URLSearchParams(params).toString();
	return apiGet(`/forms/templates${qs ? `?${qs}` : ''}`);
}

export function createFormTemplate(templateData) {
	return apiPost('/forms/templates', templateData);
}

export function getFormTemplateCategories() {
	return apiGet('/forms/templates/categories');
}

export function useFormTemplate(templateId, data) {
	return apiPost(`/forms/templates/${templateId}/use`, data);
}

// ── Relationships ───────────────────────────────────────────────────
export function getFormRelationships(formId) {
	return apiGet(`/forms/${formId}/relationships`);
}

export function createRelationship(formId, relData) {
	return apiPost(`/forms/${formId}/relationships`, relData);
}

export function createLookupField(formId, fieldData) {
	return apiPost(`/forms/${formId}/lookup-fields`, fieldData);
}

export function getFormHierarchy(formId) {
	return apiGet(`/forms/${formId}/hierarchy`);
}

export function getRelatedData(formId, recordId) {
	return apiGet(`/forms/${formId}/related-data`);
}

// ── Reports & dashboards ────────────────────────────────────────────
export function getFormReports(formId) {
	return apiGet(`/forms/${formId}/reports`);
}

export function createReport(formId, reportData) {
	return apiPost(`/forms/${formId}/reports`, reportData);
}

export function executeReport(formId, reportId) {
	return apiPost(`/forms/${formId}/reports/${reportId}/execute`, {});
}

export function generateAdHocReport(formId, config) {
	return apiPost(`/forms/${formId}/reports/ad-hoc`, config);
}

export function exportReport(formId, reportId, format = 'csv') {
	return apiGet(`/forms/${formId}/reports/${reportId}/export?format=${format}`);
}

export function getDashboard(formId, dashboardId) {
	return apiGet(`/forms/${formId}/dashboards/${dashboardId}`);
}

export function createDashboard(formId, dashboardData) {
	return apiPost(`/forms/${formId}/dashboards`, dashboardData);
}

// ── Query builder ───────────────────────────────────────────────────
export function getAvailableTables(formId) {
	return apiGet(`/forms/${formId}/query/tables`);
}

export function buildSQL(formId, elements) {
	return apiPost(`/forms/${formId}/query/build`, elements);
}

export function executeVisualQuery(formId, elements) {
	return apiPost(`/forms/${formId}/query/execute`, elements);
}

export function validateVisualQuery(formId, elements) {
	return apiPost(`/forms/${formId}/query/validate`, elements);
}

export function getQuerySuggestions(formId, elements) {
	return apiPost(`/forms/${formId}/query/suggestions`, elements);
}

export function saveQueryTemplate(formId, data) {
	return apiPost(`/forms/${formId}/query/templates`, data);
}

export function getQueryTemplates(formId) {
	return apiGet(`/forms/${formId}/query/templates`);
}

// ── Email distribution ──────────────────────────────────────────────
export function getEmailDistributions(formId) {
	return apiGet(`/forms/${formId}/email-distributions`);
}

export function createEmailDistribution(formId, data) {
	return apiPost(`/forms/${formId}/email-distributions`, data);
}

export function sendFormResponseEmail(formId, responseId, data) {
	return apiPost(`/forms/${formId}/email-distributions/${responseId}/send`, data);
}

// ── Workflow ────────────────────────────────────────────────────────
export function createFormWorkflow(formId, data) {
	return apiPost(`/forms/${formId}/workflows`, data);
}

export function startWorkflow(formId, workflowId, recordId) {
	return apiPost(`/forms/${formId}/workflows/start`, { workflow_id: workflowId, record_id: recordId });
}

export function processApproval(formId, approvalId, data) {
	return apiPost(`/forms/${formId}/workflows/${approvalId}/approve`, data);
}

export function getPendingApprovals(formId) {
	return apiGet(`/forms/${formId}/workflows/approvals`);
}

export function createNotificationTemplate(formId, data) {
	return apiPost(`/forms/${formId}/workflows/notification-templates`, data);
}

/**
 * Document AI processing/analysis endpoints.
 */

export function processDocument(documentId, analysisType) {
	return apiPost(`/documents/${documentId}/process`, { analysis_type: analysisType });
}

export function getDocumentAnalysis(documentId) {
	return apiGet(`/documents/${documentId}/analysis`);
}

export function getDocumentAnalyses(documentId) {
	return apiGet(`/documents/${documentId}/analyses`);
}

export function getDocumentProcessingStatus(documentId) {
	return apiGet(`/documents/${documentId}/status`);
}

export function uploadDocumentWithAI(data) {
	return apiPost('/documents/upload', data);
}

/**
 * AI Settings API functions
 */

/**
 * Get user's AI settings
 * @returns {Promise<Object>} AI settings
 */
export function getAISettings() {
	return apiGet('/settings/ai');
}

/**
 * Update user's AI settings
 * @param {Object} settingsData
 * @returns {Promise<Object>} Update confirmation
 */
export function updateAISettings(settingsData) {
	return apiPut('/settings/ai', settingsData);
}

/**
 * Get user's speech settings
 * @returns {Promise<Object>} Speech settings
 */
export function getSpeechSettings() {
	return apiGet('/settings/speech');
}

/**
 * Update user's speech settings
 * @param {Object} settingsData
 * @returns {Promise<Object>} Update confirmation
 */
export function updateSpeechSettings(settingsData) {
	return apiPut('/settings/speech', settingsData);
}

/**
 * Get AI usage statistics
 * @param {Object} options
 * @param {string} options.period - "7d", "30d", "90d"
 * @param {string} options.provider - Optional provider filter
 * @param {string} options.service - Optional service filter
 * @returns {Promise<Object>} Usage statistics
 */
export function getAIUsageStats(options = {}) {
	const params = new URLSearchParams();
	if (options.period) params.append('period', options.period);
	if (options.provider) params.append('provider', options.provider);
	if (options.service) params.append('service', options.service);

	const queryString = params.toString();
	const endpoint = queryString ? `/settings/ai/usage?${queryString}` : '/settings/ai/usage';

	return apiGet(endpoint);
}

/**
 * Test an API key for a provider
 * @param {string} provider - "openai", "elevenlabs", "replicate", etc.
 * @param {string} apiKey
 * @returns {Promise<Object>} Test result
 */
export function testAPIKey(provider, apiKey) {
	return apiPost('/settings/ai/test-key', { provider, api_key: apiKey });
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

/**
 * Download a binary response as a file
 * @param {string} endpoint
 * @param {Object} body
 * @returns {Promise<Blob>}
 */
export async function apiDownload(endpoint, body = {}) {
	const url = `${API_BASE_URL}${endpoint}`;
	const token = getAuthToken();
	const response = await fetch(url, {
		method: body && Object.keys(body).length ? 'POST' : 'GET',
		headers: {
			'Content-Type': 'application/json',
			...(token ? { Authorization: `Bearer ${token}` } : {}),
		},
		...(body && Object.keys(body).length ? { body: JSON.stringify(body) } : {}),
	});
	if (!response.ok) {
		const error = await response.json().catch(() => ({ message: 'Download failed' }));
		throw new Error(error.message || `HTTP ${response.status}`);
	}
	return response.blob();
}

/**
 * Admin API functions
 */

export function getAdminStats() {
	return apiGet('/admin/stats');
}

export function getAdminUsers(params = {}) {
	const q = new URLSearchParams();
	if (params.page) q.append('page', params.page);
	if (params.limit) q.append('limit', params.limit);
	if (params.search) q.append('search', params.search);
	if (params.status) q.append('status', params.status);
	const qs = q.toString();
	return apiGet(`/admin/users${qs ? `?${qs}` : ''}`);
}

export function updateAdminUserStatus(id, isActive) {
	return apiPut(`/admin/users/${id}/status`, { is_active: isActive });
}

export function deleteAdminUser(id) {
	return apiDelete(`/admin/users/${id}`);
}

export function getAdminLogs(params = {}) {
	const q = new URLSearchParams();
	if (params.limit) q.append('limit', params.limit);
	if (params.level) q.append('level', params.level);
	if (params.event_type) q.append('event_type', params.event_type);
	const qs = q.toString();
	return apiGet(`/admin/logs${qs ? `?${qs}` : ''}`);
}

export function getAdminDatabaseStats() {
	return apiGet('/admin/database/stats');
}

export function runAdminDatabaseMaintenance() {
	return apiPost('/admin/database/maintenance', {});
}

export function getAdminSecurityAlerts(limit) {
	return apiGet(`/admin/security/alerts${limit ? `?limit=${limit}` : ''}`);
}

export function resolveAdminSecurityAlert(id) {
	return apiPost(`/admin/security/alerts/${id}/resolve`, {});
}

export function getAdminSettings() {
	return apiGet('/admin/settings');
}

export function updateAdminSettings(settings) {
	return apiPut('/admin/settings', settings);
}

/**
 * Speech (TTS/STT) API functions
 */

export function getSpeechModels(type) {
	return apiGet(`/speech/models${type ? `?type=${type}` : ''}`);
}

export function createSpeechModel(data) {
	return apiPost('/speech/models', data);
}

export function textToSpeech(data) {
	return apiPost('/speech/tts', data);
}

export function speechToText(formData) {
	const url = `${API_BASE_URL}/speech/stt`;
	const token = getAuthToken();
	return fetch(url, {
		method: 'POST',
		headers: token ? { Authorization: `Bearer ${token}` } : {},
		body: formData,
	}).then(async response => {
		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Speech recognition failed' }));
			throw new Error(error.message || `HTTP ${response.status}`);
		}
		return response.json();
	});
}

export function getSpeechRequestStatus(requestId) {
	return apiGet(`/speech/requests/${requestId}`);
}

export function getSpeechHistory(params = {}) {
	const q = new URLSearchParams();
	if (params.type) q.append('type', params.type);
	if (params.status) q.append('status', params.status);
	if (params.limit) q.append('limit', params.limit);
	const qs = q.toString();
	return apiGet(`/speech/history${qs ? `?${qs}` : ''}`);
}

/**
 * Voice notes & annotations API functions
 */

export function getVoiceNotes(params = {}) {
	const q = new URLSearchParams();
	if (params.limit) q.append('limit', params.limit);
	if (params.offset) q.append('offset', params.offset);
	if (params.favorites) q.append('favorites', params.favorites);
	if (params.tag) q.append('tag', params.tag);
	const qs = q.toString();
	return apiGet(`/voice/notes${qs ? `?${qs}` : ''}`);
}

export function getVoiceNote(id) {
	return apiGet(`/voice/notes/${id}`);
}

export function createVoiceNote(data) {
	return apiPost('/voice/notes', data);
}

export function updateVoiceNote(id, data) {
	return apiPut(`/voice/notes/${id}`, data);
}

export function deleteVoiceNote(id) {
	return apiDelete(`/voice/notes/${id}`);
}

export function getVoiceAnnotations(params = {}) {
	const q = new URLSearchParams();
	if (params.content_type) q.append('content_type', params.content_type);
	if (params.content_id) q.append('content_id', params.content_id);
	const qs = q.toString();
	return apiGet(`/voice/annotations${qs ? `?${qs}` : ''}`);
}

export function getVoiceAnnotation(id) {
	return apiGet(`/voice/annotations/${id}`);
}

export function createVoiceAnnotation(data) {
	return apiPost('/voice/annotations', data);
}

export function deleteVoiceAnnotation(id) {
	return apiDelete(`/voice/annotations/${id}`);
}

/**
 * Math API functions
 */

export function validateExpression(expression) {
	return apiPost('/math/validate', { expression });
}

export function optimizeExpression(expression) {
	return apiPost('/math/optimize', { expression });
}

export function convertExpression(expression, fromFormat, toFormat) {
	return apiPost('/math/convert', { expression, from_format: fromFormat, to_format: toFormat });
}

export function getMathFunctions(category) {
	return apiGet(`/math/functions${category ? `?category=${category}` : ''}`);
}

export function getMathSymbols(category) {
	return apiGet(`/math/symbols${category ? `?category=${category}` : ''}`);
}

export function getMathConstants() {
	return apiGet('/math/constants');
}

export function getMathTheorems(category) {
	return apiGet(`/math/theorems${category ? `?category=${category}` : ''}`);
}

export function recognizeHandwriting(strokes, width, height) {
	return apiPost('/math/recognize', { strokes, width, height });
}

export function recognizeEquationFromImage(formData) {
	const url = `${API_BASE_URL}/math/recognize-image`;
	const token = getAuthToken();
	return fetch(url, {
		method: 'POST',
		headers: token ? { Authorization: `Bearer ${token}` } : {},
		body: formData,
	}).then(async response => {
		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Recognition failed' }));
			throw new Error(error.message || `HTTP ${response.status}`);
		}
		return response.json();
	});
}

export function getEquationTemplates(params = {}) {
	const q = new URLSearchParams();
	if (params.category) q.append('category', params.category);
	if (params.search) q.append('search', params.search);
	const qs = q.toString();
	return apiGet(`/math/templates${qs ? `?${qs}` : ''}`);
}

export function saveEquationTemplate(data) {
	return apiPost('/math/templates', data);
}

export function searchEquations(query) {
	return apiGet(`/math/templates/search?q=${encodeURIComponent(query)}`);
}

export function getEquationTemplateCategories() {
	return apiGet('/math/templates/categories');
}

export function saveMathCanvas(data) {
	return apiPost('/math/canvas', data);
}

export function getMathCanvases() {
	return apiGet('/math/canvas');
}

export function generateEquationImage(data) {
	return apiPost('/math/canvas/generate-image', data);
}

export function exportEquation(expression, format) {
	return apiPost('/math/export', { expression, format });
}

export function batchExportEquations(expressions, format) {
	return apiPost('/math/export/batch', { expressions, format });
}

/**
 * Document export API functions
 */

export function getDocumentExportFormats() {
	return apiGet('/documents/export/formats');
}

export function getDOCXTemplates(category) {
	return apiGet(`/documents/export/docx/templates${category ? `?category=${category}` : ''}`);
}

export function getDOCXFeatures() {
	return apiGet('/documents/export/docx/features');
}

export function exportDocument(documentId, body) {
	return apiDownload(`/documents/${documentId}/export`, body);
}

export function convertDocument(body) {
	return apiDownload('/documents/export/convert', body);
}

export function batchExportDocuments(body) {
	return apiPost('/documents/export/batch', body);
}

export function getDocumentStatistics(content) {
	return apiPost('/documents/export/statistics', { content });
}

export function validateDocumentContent(content, format) {
	return apiPost('/documents/export/validate', { content, format });
}

/**
 * Workflow API functions
 */

export function getWorkflows(params = {}) {
	const q = new URLSearchParams();
	if (params.category) q.append('category', params.category);
	if (params.active) q.append('active', params.active);
	const qs = q.toString();
	return apiGet(`/workflows${qs ? `?${qs}` : ''}`);
}

export function createWorkflow(data) {
	return apiPost('/workflows', data);
}

export function getWorkflow(id) {
	return apiGet(`/workflows/${id}`);
}

export function updateWorkflow(id, data) {
	return apiPut(`/workflows/${id}`, data);
}

export function deleteWorkflow(id) {
	return apiDelete(`/workflows/${id}`);
}

export function executeWorkflow(id, triggerData = {}) {
	return apiPost(`/workflows/${id}/execute`, { trigger_data: triggerData });
}

export function dryRunWorkflow(id, triggerData = {}) {
	return apiPost(`/workflows/${id}/execute`, { trigger_data: triggerData, dry_run: true });
}

export function getWorkflowExecutions(id, limit) {
	return apiGet(`/workflows/${id}/executions${limit ? `?limit=${limit}` : ''}`);
}

export function getWorkflowExecution(executionId) {
	return apiGet(`/executions/${executionId}`);
}

export function updateWorkflowNodes(id, nodes) {
	return apiPut(`/workflows/${id}/nodes`, { nodes });
}

export function updateWorkflowConnections(id, connections) {
	return apiPut(`/workflows/${id}/connections`, { connections });
}

export function getWorkflowTemplates(category) {
	return apiGet(`/workflow-templates${category ? `?category=${category}` : ''}`);
}

export function createWorkflowFromTemplate(templateId) {
	return apiPost(`/workflow-templates/${templateId}/create`, {});
}

export function getWorkflowInsights() {
	return apiGet('/ai/workflows/insights');
}

export function analyzeWorkflowUsage() {
	return apiGet('/ai/workflows/usage-analysis');
}

export function getWorkflowPredictions() {
	return apiGet('/ai/workflows/predictions');
}

export function optimizeWorkflow(id) {
	return apiGet(`/ai/workflows/${id}/optimization`);
}

export function getIntegrationConnectors() {
	return apiGet('/connectors');
}

/**
 * MCP (Model Context Protocol) servers — bridge to external systems.
 */
export function getMCPServers() {
	return apiGet('/mcp/servers');
}

export function createMCPServer(payload) {
	return apiPost('/mcp/servers', payload);
}

export function deleteMCPServer(id) {
	return apiDelete(`/mcp/servers/${id}`);
}

export function testMCPServer(payload) {
	return apiPost('/mcp/servers/test', payload);
}

export function getMCPConnectors() {
	return apiGet('/mcp/connectors');
}

/**
 * Monitoring API functions
 */

export function getMonitoringMetrics() {
	return apiGet('/monitoring/metrics');
}

export function getMonitoringHealth() {
	return apiGet('/monitoring/health');
}

export function getMonitoringPerformance() {
	return apiGet('/monitoring/performance');
}

export function getMonitoringAlerts(limit) {
	return apiGet(`/monitoring/alerts${limit ? `?limit=${limit}` : ''}`);
}
