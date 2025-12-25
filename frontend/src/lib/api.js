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
