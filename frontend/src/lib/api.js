// API client utilities for TPT Titan frontend

const API_BASE_URL = import.meta.env.DEV ? 'http://localhost:8080/api/v1' : '/api/v1';

/**
 * Generic API request function
 * @param {string} endpoint - API endpoint (without base URL)
 * @param {any} options - Fetch options
 * @returns {Promise<any>} Response data
 */
export async function apiRequest(endpoint, options = {}) {
	const url = `${API_BASE_URL}${endpoint}`;

	/** @type {Record<string, string>} */
	const defaultHeaders = {
		'Content-Type': 'application/json',
	};

	// Add authorization header if token exists
	const token = getAuthToken();
	if (token) {
		defaultHeaders['Authorization'] = `Bearer ${token}`;
	}

	const response = await fetch(url, {
		...options,
		headers: {
			...defaultHeaders,
			...(options.headers || {}),
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
 * @returns {Promise<any>}
 */
export function apiGet(endpoint) {
	return apiRequest(endpoint);
}

/**
 * POST request
 * @param {string} endpoint
 * @param {any} data
 * @returns {Promise<any>}
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
 * @param {any} data
 * @returns {Promise<any>}
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
 * @returns {Promise<any>}
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
 * @returns {Promise<any>}
 */
export function healthCheck() {
	return apiGet('/health');
}

/**
 * Document API functions
 */

/**
 * Get all documents for the authenticated user
 * @returns {Promise<any>} Documents list
 */
export function getDocuments() {
	return apiGet('/documents');
}

/**
 * Get a specific document
 * @param {any} documentId
 * @returns {Promise<any>} Document data
 */
export function getDocument(documentId) {
	return apiGet(`/documents/${documentId}`);
}

/**
 * Create a new document
 * @param {any} documentData
 * @returns {Promise<any>} Created document
 */
export function createDocument(documentData) {
	return apiPost('/documents', documentData);
}

/**
 * Update an existing document
 * @param {any} documentId
 * @param {any} documentData
 * @returns {Promise<any>} Updated document
 */
export function updateDocument(documentId, documentData) {
	return apiPut(`/documents/${documentId}`, documentData);
}

/**
 * Delete a document
 * @param {any} documentId
 * @returns {Promise<any>} Deletion confirmation
 */
export function deleteDocument(documentId) {
	return apiDelete(`/documents/${documentId}`);
}

/**
 * Get version history for a document
 * @param {any} documentId
 * @returns {Promise<any>} Version history
 */
export function getDocumentVersions(documentId) {
	return apiGet(`/documents/${documentId}/versions`);
}

/**
 * Restore a specific version of a document
 * @param {any} documentId
 * @param {any} version
 * @returns {Promise<any>} Restoration confirmation
 */
export function restoreDocumentVersion(documentId, version) {
	return apiPost(`/documents/${documentId}/versions/${version}/restore`, {});
}

/**
 * AI API functions
 */

/**
 * Get all available AI models for the user
 * @returns {Promise<any>} Available models
 */
export function getAIModels() {
	return apiGet('/ai/models');
}

/**
 * Create a custom AI model
 * @param {any} modelData
 * @returns {Promise<any>} Created model
 */
export function createAIModel(modelData) {
	return apiPost('/ai/models', modelData);
}

/**
 * Process an AI request
 * @param {any} requestData
 * @returns {Promise<any>} Request status
 */
export function processAIRequest(requestData) {
	return apiPost('/ai/requests', requestData);
}

/**
 * Get AI request status
 * @param {any} requestId
 * @returns {Promise<any>} Request status and result
 */
export function getAIRequestStatus(requestId) {
	return apiGet(`/ai/requests/${requestId}`);
}

/**
 * List available Ollama models
 * @returns {Promise<any>} Ollama models
 */
export function listOllamaModels() {
	return apiGet('/ai/ollama/models');
}

/**
 * Pull an Ollama model
 * @param {any} modelName
 * @returns {Promise<any>} Pull status
 */
export function pullOllamaModel(modelName) {
	return apiPost(`/ai/ollama/models/${modelName}/pull`, {});
}

/**
 * Get AI usage statistics
 * @returns {Promise<any>} Usage data
 */
export function getAIUsage() {
	return apiGet('/ai/usage');
}

/**
 * Detect system hardware capabilities
 * @returns {Promise<any>} Hardware information
 */
export function detectHardware() {
	return apiGet('/ai/hardware');
}

/**
 * Get AI model recommendations based on hardware
 * @returns {Promise<any>} Hardware and model recommendations
 */
export function getRecommendedModels() {
	return apiGet('/ai/recommendations');
}

/**
 * Setup recommended models and tasks for the user
 * @returns {Promise<any>} Setup confirmation
 */
export function setupRecommendedModels() {
	return apiPost('/ai/setup', {});
}

/**
 * Check for available AI model upgrades
 * @returns {Promise<any>} Upgrade check results with available options
 */
export function checkForAIUpgrades() {
	return apiPost('/ai/upgrades/check', {});
}

/**
 * Get AI upgrade check history
 * @returns {Promise<any>} List of past upgrade checks
 */
export function getUpgradeHistory() {
	return apiGet('/ai/upgrades/history');
}

/**
 * Apply a selected AI model upgrade
 * @param {any} upgradeId - The ID of the upgrade check
 * @returns {Promise<any>} Upgrade application status
 */
export function applyAIUpgrade(upgradeId) {
	return apiPost('/ai/upgrades/apply', { upgrade_id: upgradeId });
}

/**
 * Spreadsheet API functions
 */

/**
 * Create a new spreadsheet
 * @param {any} spreadsheetData
 * @returns {Promise<any>} Created spreadsheet
 */
export function createSpreadsheet(spreadsheetData) {
	return apiPost('/spreadsheets', spreadsheetData);
}

/**
 * Get a specific spreadsheet
 * @param {any} spreadsheetId
 * @returns {Promise<any>} Spreadsheet data
 */
export function getSpreadsheet(spreadsheetId) {
	return apiGet(`/spreadsheets/${spreadsheetId}`);
}

/**
 * Update a spreadsheet cell
 * @param {any} spreadsheetId
 * @param {any} cellData
 * @returns {Promise<any>} Update confirmation
 */
export function updateSpreadsheetCell(spreadsheetId, cellData) {
	return apiPut(`/spreadsheets/${spreadsheetId}/cells`, cellData);
}

/**
 * Update multiple spreadsheet cells (batch)
 * @param {any} spreadsheetId
 * @param {any} batchData
 * @returns {Promise<any>} Batch update confirmation
 */
export function updateSpreadsheetCells(spreadsheetId, batchData) {
	return apiPut(`/spreadsheets/${spreadsheetId}/batch`, batchData);
}

/**
 * Evaluate a spreadsheet formula
 * @param {any} formulaData
 * @returns {Promise<any>} Formula evaluation result
 */
export function evaluateFormula(formulaData) {
	return apiPost('/spreadsheets/evaluate', formulaData);
}

/**
 * Get available spreadsheet functions
 * @returns {Promise<any>} List of available functions
 */
export function getAvailableFunctions() {
	return apiGet('/spreadsheets/functions');
}

/**
 * Validate a spreadsheet formula
 * @param {any} validationData
 * @returns {Promise<any>} Validation result
 */
export function validateFormula(validationData) {
	return apiPost('/spreadsheets/validate', validationData);
}

/**
 * Get chart suggestions for spreadsheet data
 * @param {any} chartData
 * @returns {Promise<any>} Chart suggestions
 */
export function getChartSuggestions(chartData) {
	return apiPost('/spreadsheets/charts/suggestions', chartData);
}

/**
 * Create a chart from spreadsheet data
 * @param {any} chartData
 * @returns {Promise<any>} Created chart
 */
export function createChart(chartData) {
	return apiPost('/spreadsheets/charts', chartData);
}

/**
 * Get charts for a spreadsheet
 * @param {any} spreadsheetId
 * @returns {Promise<any>} Charts list
 */
export function getSpreadsheetCharts(spreadsheetId) {
	return apiGet(`/spreadsheets/${spreadsheetId}/charts`);
}

/**
 * Export spreadsheet to Excel
 * @param {any} spreadsheetId
 * @param {any} exportOptions
 * @returns {Promise<Blob>} Excel file
 */
export function exportSpreadsheetToExcel(spreadsheetId, exportOptions = {}) {
	const url = `${API_BASE_URL}/spreadsheets/${spreadsheetId}/export/excel`;

	/** @type {Record<string, string>} */
	const exportHeaders = {
		'Content-Type': 'application/json',
	};

	const token = getAuthToken();
	if (token) {
		exportHeaders['Authorization'] = `Bearer ${token}`;
	}

	return fetch(url, {
		method: 'POST',
		body: JSON.stringify(exportOptions),
		headers: exportHeaders,
	}).then(response => {
		if (!response.ok) {
			throw new Error(`Export failed: ${response.status}`);
		}
		return response.blob();
	});
}

/**
 * Import Excel file to spreadsheet
 * @param {any} formData - Form data with file
 * @returns {Promise<any>} Import result
 */
export function importExcelToSpreadsheet(formData) {
	const url = `${API_BASE_URL}/spreadsheets/import/excel`;

	/** @type {Record<string, string>} */
	const importHeaders = {};

	const token = getAuthToken();
	if (token) {
		importHeaders['Authorization'] = `Bearer ${token}`;
	}

	return fetch(url, {
		method: 'POST',
		body: formData,
		headers: importHeaders,
	}).then(async response => {
		if (!response.ok) {
			const error = /** @type {any} */ (await response.json().catch(() => ({ message: 'Import failed' })));
			throw new Error(error.message || `HTTP ${response.status}`);
		}
		return response.json();
	});
}

/**
 * Get spreadsheet version info
 * @param {any} spreadsheetId
 * @returns {Promise<any>} Version information
 */
export function getSpreadsheetVersion(spreadsheetId) {
	return apiGet(`/spreadsheets/${spreadsheetId}/version`);
}

/**
 * Get spreadsheet changes since version
 * @param {any} spreadsheetId
 * @param {any} sinceVersion
 * @returns {Promise<any>} Changes list
 */
export function getSpreadsheetChanges(spreadsheetId, sinceVersion) {
	return apiGet(`/spreadsheets/${spreadsheetId}/changes?since_version=${sinceVersion}`);
}

/**
 * Lock spreadsheet cells for editing
 * @param {any} spreadsheetId
 * @param {any} lockData
 * @returns {Promise<any>} Lock confirmation
 */
export function lockSpreadsheetCells(spreadsheetId, lockData) {
	return apiPost(`/spreadsheets/${spreadsheetId}/lock`, lockData);
}

/**
 * Unlock spreadsheet cells
 * @param {any} spreadsheetId
 * @param {any} unlockData
 * @returns {Promise<any>} Unlock confirmation
 */
export function unlockSpreadsheetCells(spreadsheetId, unlockData) {
	return apiPost(`/spreadsheets/${spreadsheetId}/unlock`, unlockData);
}

/**
 * Form API functions
 */

/**
 * Get all forms for the authenticated user
 * @returns {Promise<any>} Forms list
 */
export function getForms() {
	return apiGet('/forms');
}

/**
 * Get a specific form
 * @param {any} formId
 * @returns {Promise<any>} Form data
 */
export function getForm(formId) {
	return apiGet(`/forms/${formId}`);
}

/**
 * Create a new form
 * @param {any} formData
 * @returns {Promise<any>} Created form
 */
export function createForm(formData) {
	return apiPost('/forms', formData);
}

/**
 * Update an existing form
 * @param {any} formId
 * @param {any} formData
 * @returns {Promise<any>} Updated form
 */
export function updateForm(formId, formData) {
	return apiPut(`/forms/${formId}`, formData);
}

/**
 * Delete a form
 * @param {any} formId
 * @returns {Promise<any>} Deletion confirmation
 */
export function deleteForm(formId) {
	return apiDelete(`/forms/${formId}`);
}

/**
 * Get form responses
 * @param {any} formId
 * @returns {Promise<any>} Form responses
 */
export function getFormResponses(formId) {
	return apiGet(`/forms/${formId}/responses`);
}

/**
 * Submit a form response
 * @param {any} formId
 * @param {any} responseData
 * @returns {Promise<any>} Submission confirmation
 */
export function submitFormResponse(formId, responseData) {
	return apiPost(`/forms/${formId}/responses`, responseData);
}

/**
 * Advanced form modules (templates, relationships, reports, query builder,
 * email distribution, workflow) reachable via the formsAdvancedGroup routes.
 */

// ── Templates ───────────────────────────────────────────────────────
/**
 * @param {any} params
 * @returns {Promise<any>}
 */
export function getFormTemplates(params = {}) {
	const qs = new URLSearchParams(params).toString();
	return apiGet(`/forms/templates${qs ? `?${qs}` : ''}`);
}

/**
 * @param {any} templateData
 * @returns {Promise<any>}
 */
export function createFormTemplate(templateData) {
	return apiPost('/forms/templates', templateData);
}

/**
 * @returns {Promise<any>}
 */
export function getFormTemplateCategories() {
	return apiGet('/forms/templates/categories');
}

/**
 * @param {any} templateId
 * @param {any} data
 * @returns {Promise<any>}
 */
export function useFormTemplate(templateId, data) {
	return apiPost(`/forms/templates/${templateId}/use`, data);
}

// ── Relationships ───────────────────────────────────────────────────
/**
 * @param {any} formId
 * @returns {Promise<any>}
 */
export function getFormRelationships(formId) {
	return apiGet(`/forms/${formId}/relationships`);
}

/**
 * @param {any} formId
 * @param {any} relData
 * @returns {Promise<any>}
 */
export function createRelationship(formId, relData) {
	return apiPost(`/forms/${formId}/relationships`, relData);
}

/**
 * @param {any} formId
 * @param {any} fieldData
 * @returns {Promise<any>}
 */
export function createLookupField(formId, fieldData) {
	return apiPost(`/forms/${formId}/lookup-fields`, fieldData);
}

/**
 * @param {any} formId
 * @returns {Promise<any>}
 */
export function getFormHierarchy(formId) {
	return apiGet(`/forms/${formId}/hierarchy`);
}

/**
 * @param {any} formId
 * @param {any} recordId
 * @returns {Promise<any>}
 */
export function getRelatedData(formId, recordId) {
	return apiGet(`/forms/${formId}/related-data`);
}

// ── Reports & dashboards ────────────────────────────────────────────
/**
 * @param {any} formId
 * @returns {Promise<any>}
 */
export function getFormReports(formId) {
	return apiGet(`/forms/${formId}/reports`);
}

/**
 * @param {any} formId
 * @param {any} reportData
 * @returns {Promise<any>}
 */
export function createReport(formId, reportData) {
	return apiPost(`/forms/${formId}/reports`, reportData);
}

/**
 * @param {any} formId
 * @param {any} reportId
 * @returns {Promise<any>}
 */
export function executeReport(formId, reportId) {
	return apiPost(`/forms/${formId}/reports/${reportId}/execute`, {});
}

/**
 * @param {any} formId
 * @param {any} config
 * @returns {Promise<any>}
 */
export function generateAdHocReport(formId, config) {
	return apiPost(`/forms/${formId}/reports/ad-hoc`, config);
}

/**
 * @param {any} formId
 * @param {any} reportId
 * @param {string} format
 * @returns {Promise<any>}
 */
export function exportReport(formId, reportId, format = 'csv') {
	return apiGet(`/forms/${formId}/reports/${reportId}/export?format=${format}`);
}

/**
 * @param {any} formId
 * @param {any} dashboardId
 * @returns {Promise<any>}
 */
export function getDashboard(formId, dashboardId) {
	return apiGet(`/forms/${formId}/dashboards/${dashboardId}`);
}

/**
 * @param {any} formId
 * @param {any} dashboardData
 * @returns {Promise<any>}
 */
export function createDashboard(formId, dashboardData) {
	return apiPost(`/forms/${formId}/dashboards`, dashboardData);
}

// ── Query builder ───────────────────────────────────────────────────
/**
 * @param {any} formId
 * @returns {Promise<any>}
 */
export function getAvailableTables(formId) {
	return apiGet(`/forms/${formId}/query/tables`);
}

/**
 * @param {any} formId
 * @param {any} elements
 * @returns {Promise<any>}
 */
export function buildSQL(formId, elements) {
	return apiPost(`/forms/${formId}/query/build`, elements);
}

/**
 * @param {any} formId
 * @param {any} elements
 * @returns {Promise<any>}
 */
export function executeVisualQuery(formId, elements) {
	return apiPost(`/forms/${formId}/query/execute`, elements);
}

/**
 * @param {any} formId
 * @param {any} elements
 * @returns {Promise<any>}
 */
export function validateVisualQuery(formId, elements) {
	return apiPost(`/forms/${formId}/query/validate`, elements);
}

/**
 * @param {any} formId
 * @param {any} elements
 * @returns {Promise<any>}
 */
export function getQuerySuggestions(formId, elements) {
	return apiPost(`/forms/${formId}/query/suggestions`, elements);
}

/**
 * @param {any} formId
 * @param {any} data
 * @returns {Promise<any>}
 */
export function saveQueryTemplate(formId, data) {
	return apiPost(`/forms/${formId}/query/templates`, data);
}

/**
 * @param {any} formId
 * @returns {Promise<any>}
 */
export function getQueryTemplates(formId) {
	return apiGet(`/forms/${formId}/query/templates`);
}

// ── Email distribution ──────────────────────────────────────────────
/**
 * @param {any} formId
 * @returns {Promise<any>}
 */
export function getEmailDistributions(formId) {
	return apiGet(`/forms/${formId}/email-distributions`);
}

/**
 * @param {any} formId
 * @param {any} data
 * @returns {Promise<any>}
 */
export function createEmailDistribution(formId, data) {
	return apiPost(`/forms/${formId}/email-distributions`, data);
}

/**
 * @param {any} formId
 * @param {any} responseId
 * @param {any} data
 * @returns {Promise<any>}
 */
export function sendFormResponseEmail(formId, responseId, data) {
	return apiPost(`/forms/${formId}/email-distributions/${responseId}/send`, data);
}

// ── Workflow ────────────────────────────────────────────────────────
/**
 * @param {any} formId
 * @param {any} data
 * @returns {Promise<any>}
 */
export function createFormWorkflow(formId, data) {
	return apiPost(`/forms/${formId}/workflows`, data);
}

/**
 * @param {any} formId
 * @param {any} workflowId
 * @param {any} recordId
 * @returns {Promise<any>}
 */
export function startWorkflow(formId, workflowId, recordId) {
	return apiPost(`/forms/${formId}/workflows/start`, { workflow_id: workflowId, record_id: recordId });
}

/**
 * @param {any} formId
 * @param {any} approvalId
 * @param {any} data
 * @returns {Promise<any>}
 */
export function processApproval(formId, approvalId, data) {
	return apiPost(`/forms/${formId}/workflows/${approvalId}/approve`, data);
}

/**
 * @param {any} formId
 * @returns {Promise<any>}
 */
export function getPendingApprovals(formId) {
	return apiGet(`/forms/${formId}/workflows/approvals`);
}

/**
 * @param {any} formId
 * @param {any} data
 * @returns {Promise<any>}
 */
export function createNotificationTemplate(formId, data) {
	return apiPost(`/forms/${formId}/workflows/notification-templates`, data);
}

/**
 * Document AI processing/analysis endpoints.
 */

/**
 * @param {any} documentId
 * @param {any} analysisType
 * @returns {Promise<any>}
 */
export function processDocument(documentId, analysisType) {
	return apiPost(`/documents/${documentId}/process`, { analysis_type: analysisType });
}

/**
 * @param {any} documentId
 * @returns {Promise<any>}
 */
export function getDocumentAnalysis(documentId) {
	return apiGet(`/documents/${documentId}/analysis`);
}

/**
 * @param {any} documentId
 * @returns {Promise<any>}
 */
export function getDocumentAnalyses(documentId) {
	return apiGet(`/documents/${documentId}/analyses`);
}

/**
 * @param {any} documentId
 * @returns {Promise<any>}
 */
export function getDocumentProcessingStatus(documentId) {
	return apiGet(`/documents/${documentId}/status`);
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function uploadDocumentWithAI(data) {
	return apiPost('/documents/upload', data);
}

/**
 * AI Settings API functions
 */

/**
 * Get user's AI settings
 * @returns {Promise<any>} AI settings
 */
export function getAISettings() {
	return apiGet('/settings/ai');
}

/**
 * Update user's AI settings
 * @param {any} settingsData
 * @returns {Promise<any>} Update confirmation
 */
export function updateAISettings(settingsData) {
	return apiPut('/settings/ai', settingsData);
}

/**
 * Get user's speech settings
 * @returns {Promise<any>} Speech settings
 */
export function getSpeechSettings() {
	return apiGet('/settings/speech');
}

/**
 * Update user's speech settings
 * @param {any} settingsData
 * @returns {Promise<any>} Update confirmation
 */
export function updateSpeechSettings(settingsData) {
	return apiPut('/settings/speech', settingsData);
}

/**
 * Get AI usage statistics
 * @param {any} options
 * @returns {Promise<any>} Usage statistics
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
 * @param {any} provider - "openai", "elevenlabs", "replicate", etc.
 * @param {any} apiKey
 * @returns {Promise<any>} Test result
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
 * @param {any} body
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

/**
 * @returns {Promise<any>}
 */
export function getAdminStats() {
	return apiGet('/admin/stats');
}

/**
 * @param {any} params
 * @returns {Promise<any>}
 */
export function getAdminUsers(params = {}) {
	const q = new URLSearchParams();
	if (params.page) q.append('page', params.page);
	if (params.limit) q.append('limit', params.limit);
	if (params.search) q.append('search', params.search);
	if (params.status) q.append('status', params.status);
	const qs = q.toString();
	return apiGet(`/admin/users${qs ? `?${qs}` : ''}`);
}

/**
 * @param {any} id
 * @param {any} isActive
 * @returns {Promise<any>}
 */
export function updateAdminUserStatus(id, isActive) {
	return apiPut(`/admin/users/${id}/status`, { is_active: isActive });
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function deleteAdminUser(id) {
	return apiDelete(`/admin/users/${id}`);
}

/**
 * @param {any} params
 * @returns {Promise<any>}
 */
export function getAdminLogs(params = {}) {
	const q = new URLSearchParams();
	if (params.limit) q.append('limit', params.limit);
	if (params.level) q.append('level', params.level);
	if (params.event_type) q.append('event_type', params.event_type);
	const qs = q.toString();
	return apiGet(`/admin/logs${qs ? `?${qs}` : ''}`);
}

/**
 * @returns {Promise<any>}
 */
export function getAdminDatabaseStats() {
	return apiGet('/admin/database/stats');
}

/**
 * @returns {Promise<any>}
 */
export function runAdminDatabaseMaintenance() {
	return apiPost('/admin/database/maintenance', {});
}

/**
 * @param {any} limit
 * @returns {Promise<any>}
 */
export function getAdminSecurityAlerts(limit) {
	return apiGet(`/admin/security/alerts${limit ? `?limit=${limit}` : ''}`);
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function resolveAdminSecurityAlert(id) {
	return apiPost(`/admin/security/alerts/${id}/resolve`, {});
}

/**
 * @returns {Promise<any>}
 */
export function getAdminSettings() {
	return apiGet('/admin/settings');
}

/**
 * @param {any} settings
 * @returns {Promise<any>}
 */
export function updateAdminSettings(settings) {
	return apiPut('/admin/settings', settings);
}

/**
 * Speech (TTS/STT) API functions
 */

/**
 * @param {any} type
 * @returns {Promise<any>}
 */
export function getSpeechModels(type) {
	return apiGet(`/speech/models${type ? `?type=${type}` : ''}`);
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function createSpeechModel(data) {
	return apiPost('/speech/models', data);
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function textToSpeech(data) {
	return apiPost('/speech/tts', data);
}

/**
 * @param {any} formData
 * @returns {Promise<any>}
 */
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

/**
 * @param {any} requestId
 * @returns {Promise<any>}
 */
export function getSpeechRequestStatus(requestId) {
	return apiGet(`/speech/requests/${requestId}`);
}

/**
 * @param {any} params
 * @returns {Promise<any>}
 */
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

/**
 * @param {any} params
 * @returns {Promise<any>}
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

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function getVoiceNote(id) {
	return apiGet(`/voice/notes/${id}`);
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function createVoiceNote(data) {
	return apiPost('/voice/notes', data);
}

/**
 * @param {any} id
 * @param {any} data
 * @returns {Promise<any>}
 */
export function updateVoiceNote(id, data) {
	return apiPut(`/voice/notes/${id}`, data);
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function deleteVoiceNote(id) {
	return apiDelete(`/voice/notes/${id}`);
}

/**
 * @param {any} params
 * @returns {Promise<any>}
 */
export function getVoiceAnnotations(params = {}) {
	const q = new URLSearchParams();
	if (params.content_type) q.append('content_type', params.content_type);
	if (params.content_id) q.append('content_id', params.content_id);
	const qs = q.toString();
	return apiGet(`/voice/annotations${qs ? `?${qs}` : ''}`);
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function getVoiceAnnotation(id) {
	return apiGet(`/voice/annotations/${id}`);
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function createVoiceAnnotation(data) {
	return apiPost('/voice/annotations', data);
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function deleteVoiceAnnotation(id) {
	return apiDelete(`/voice/annotations/${id}`);
}

/**
 * Math API functions
 */

/**
 * @param {any} expression
 * @returns {Promise<any>}
 */
export function validateExpression(expression) {
	return apiPost('/math/validate', { expression });
}

/**
 * @param {any} expression
 * @returns {Promise<any>}
 */
export function optimizeExpression(expression) {
	return apiPost('/math/optimize', { expression });
}

/**
 * @param {any} expression
 * @param {any} fromFormat
 * @param {any} toFormat
 * @returns {Promise<any>}
 */
export function convertExpression(expression, fromFormat, toFormat) {
	return apiPost('/math/convert', { expression, from_format: fromFormat, to_format: toFormat });
}

/**
 * @param {any} category
 * @returns {Promise<any>}
 */
export function getMathFunctions(category) {
	return apiGet(`/math/functions${category ? `?category=${category}` : ''}`);
}

/**
 * @param {any} category
 * @returns {Promise<any>}
 */
export function getMathSymbols(category) {
	return apiGet(`/math/symbols${category ? `?category=${category}` : ''}`);
}

/**
 * @returns {Promise<any>}
 */
export function getMathConstants() {
	return apiGet('/math/constants');
}

/**
 * @param {any} category
 * @returns {Promise<any>}
 */
export function getMathTheorems(category) {
	return apiGet(`/math/theorems${category ? `?category=${category}` : ''}`);
}

/**
 * @param {any} strokes
 * @param {any} width
 * @param {any} height
 * @returns {Promise<any>}
 */
export function recognizeHandwriting(strokes, width, height) {
	return apiPost('/math/recognize', { strokes, width, height });
}

/**
 * @param {any} formData
 * @returns {Promise<any>}
 */
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

/**
 * @param {any} params
 * @returns {Promise<any>}
 */
export function getEquationTemplates(params = {}) {
	const q = new URLSearchParams();
	if (params.category) q.append('category', params.category);
	if (params.search) q.append('search', params.search);
	const qs = q.toString();
	return apiGet(`/math/templates${qs ? `?${qs}` : ''}`);
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function saveEquationTemplate(data) {
	return apiPost('/math/templates', data);
}

/**
 * @param {any} query
 * @returns {Promise<any>}
 */
export function searchEquations(query) {
	return apiGet(`/math/templates/search?q=${encodeURIComponent(query)}`);
}

/**
 * @returns {Promise<any>}
 */
export function getEquationTemplateCategories() {
	return apiGet('/math/templates/categories');
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function saveMathCanvas(data) {
	return apiPost('/math/canvas', data);
}

/**
 * @returns {Promise<any>}
 */
export function getMathCanvases() {
	return apiGet('/math/canvas');
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function generateEquationImage(data) {
	return apiPost('/math/canvas/generate-image', data);
}

/**
 * @param {any} expression
 * @param {any} format
 * @returns {Promise<any>}
 */
export function exportEquation(expression, format) {
	return apiPost('/math/export', { expression, format });
}

/**
 * @param {any} expressions
 * @param {any} format
 * @returns {Promise<any>}
 */
export function batchExportEquations(expressions, format) {
	return apiPost('/math/export/batch', { expressions, format });
}

/**
 * Document export API functions
 */

/**
 * @returns {Promise<any>}
 */
export function getDocumentExportFormats() {
	return apiGet('/documents/export/formats');
}

/**
 * @param {any} category
 * @returns {Promise<any>}
 */
export function getDOCXTemplates(category) {
	return apiGet(`/documents/export/docx/templates${category ? `?category=${category}` : ''}`);
}

/**
 * @returns {Promise<any>}
 */
export function getDOCXFeatures() {
	return apiGet('/documents/export/docx/features');
}

/**
 * @param {any} documentId
 * @param {any} body
 * @returns {Promise<Blob>}
 */
export function exportDocument(documentId, body) {
	return apiDownload(`/documents/${documentId}/export`, body);
}

/**
 * @param {any} body
 * @returns {Promise<Blob>}
 */
export function convertDocument(body) {
	return apiDownload('/documents/export/convert', body);
}

/**
 * @param {any} body
 * @returns {Promise<any>}
 */
export function batchExportDocuments(body) {
	return apiPost('/documents/export/batch', body);
}

/**
 * @param {any} content
 * @returns {Promise<any>}
 */
export function getDocumentStatistics(content) {
	return apiPost('/documents/export/statistics', { content });
}

/**
 * @param {any} content
 * @param {any} format
 * @returns {Promise<any>}
 */
export function validateDocumentContent(content, format) {
	return apiPost('/documents/export/validate', { content, format });
}

/**
 * Workflow API functions
 */

/**
 * @param {any} params
 * @returns {Promise<any>}
 */
export function getWorkflows(params = {}) {
	const q = new URLSearchParams();
	if (params.category) q.append('category', params.category);
	if (params.active) q.append('active', params.active);
	const qs = q.toString();
	return apiGet(`/workflows${qs ? `?${qs}` : ''}`);
}

/**
 * @param {any} data
 * @returns {Promise<any>}
 */
export function createWorkflow(data) {
	return apiPost('/workflows', data);
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function getWorkflow(id) {
	return apiGet(`/workflows/${id}`);
}

/**
 * @param {any} id
 * @param {any} data
 * @returns {Promise<any>}
 */
export function updateWorkflow(id, data) {
	return apiPut(`/workflows/${id}`, data);
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function deleteWorkflow(id) {
	return apiDelete(`/workflows/${id}`);
}

/**
 * @param {any} id
 * @param {any} triggerData
 * @returns {Promise<any>}
 */
export function executeWorkflow(id, triggerData = {}) {
	return apiPost(`/workflows/${id}/execute`, { trigger_data: triggerData });
}

/**
 * @param {any} id
 * @param {any} triggerData
 * @returns {Promise<any>}
 */
export function dryRunWorkflow(id, triggerData = {}) {
	return apiPost(`/workflows/${id}/execute`, { trigger_data: triggerData, dry_run: true });
}

/**
 * @param {any} id
 * @param {any} limit
 * @returns {Promise<any>}
 */
export function getWorkflowExecutions(id, limit) {
	return apiGet(`/workflows/${id}/executions${limit ? `?limit=${limit}` : ''}`);
}

/**
 * @param {any} executionId
 * @returns {Promise<any>}
 */
export function getWorkflowExecution(executionId) {
	return apiGet(`/executions/${executionId}`);
}

/**
 * @param {any} id
 * @param {any} nodes
 * @returns {Promise<any>}
 */
export function updateWorkflowNodes(id, nodes) {
	return apiPut(`/workflows/${id}/nodes`, { nodes });
}

/**
 * @param {any} id
 * @param {any} connections
 * @returns {Promise<any>}
 */
export function updateWorkflowConnections(id, connections) {
	return apiPut(`/workflows/${id}/connections`, { connections });
}

/**
 * @param {any} category
 * @returns {Promise<any>}
 */
export function getWorkflowTemplates(category) {
	return apiGet(`/workflow-templates${category ? `?category=${category}` : ''}`);
}

/**
 * @param {any} templateId
 * @returns {Promise<any>}
 */
export function createWorkflowFromTemplate(templateId) {
	return apiPost(`/workflow-templates/${templateId}/create`, {});
}

/**
 * @returns {Promise<any>}
 */
export function getWorkflowInsights() {
	return apiGet('/ai/workflows/insights');
}

/**
 * @returns {Promise<any>}
 */
export function analyzeWorkflowUsage() {
	return apiGet('/ai/workflows/usage-analysis');
}

/**
 * @returns {Promise<any>}
 */
export function getWorkflowPredictions() {
	return apiGet('/ai/workflows/predictions');
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function optimizeWorkflow(id) {
	return apiGet(`/ai/workflows/${id}/optimization`);
}

// AI-authored workflow generation: turns a natural-language description into a
// draft workflow canvas (reviewed by the user before it is persisted).
/**
 * @param {any} prompt
 * @returns {Promise<any>}
 */
export function generateWorkflowFromPrompt(prompt) {
	return apiPost('/ai/workflows/generate', { prompt });
}

/**
 * @returns {Promise<any>}
 */
export function getIntegrationConnectors() {
	return apiGet('/connectors');
}

/**
 * @returns {Promise<any>}
 */
export function getNamedConnectorTemplates() {
	return apiGet('/connectors/templates');
}

// Webhook delivery log (inbound/outbound call audit trail).
/**
 * @returns {Promise<any>}
 */
export function getWebhookDeliveryLogs() {
	return apiGet('/admin/webhook-logs');
}

// Admin outbound domain allowlist.
/**
 * @returns {Promise<any>}
 */
export function getOutboundDomainAllowlist() {
	return apiGet('/admin/outbound-domains');
}

/**
 * @param {any} domains
 * @returns {Promise<any>}
 */
export function updateOutboundDomainAllowlist(domains) {
	return apiPut('/admin/outbound-domains', { domains });
}

/**
 * MCP (Model Context Protocol) servers — bridge to external systems.
 */
/**
 * @returns {Promise<any>}
 */
export function getMCPServers() {
	return apiGet('/mcp/servers');
}

/**
 * @param {any} payload
 * @returns {Promise<any>}
 */
export function createMCPServer(payload) {
	return apiPost('/mcp/servers', payload);
}

/**
 * @param {any} id
 * @returns {Promise<any>}
 */
export function deleteMCPServer(id) {
	return apiDelete(`/mcp/servers/${id}`);
}

/**
 * @param {any} payload
 * @returns {Promise<any>}
 */
export function testMCPServer(payload) {
	return apiPost('/mcp/servers/test', payload);
}

/**
 * @returns {Promise<any>}
 */
export function getMCPConnectors() {
	return apiGet('/mcp/connectors');
}

/**
 * Monitoring API functions
 */

/**
 * @returns {Promise<any>}
 */
export function getMonitoringMetrics() {
	return apiGet('/monitoring/metrics');
}

/**
 * @returns {Promise<any>}
 */
export function getMonitoringHealth() {
	return apiGet('/monitoring/health');
}

/**
 * @returns {Promise<any>}
 */
export function getMonitoringPerformance() {
	return apiGet('/monitoring/performance');
}

/**
 * @param {any} limit
 * @returns {Promise<any>}
 */
export function getMonitoringAlerts(limit) {
	return apiGet(`/monitoring/alerts${limit ? `?limit=${limit}` : ''}`);
}
