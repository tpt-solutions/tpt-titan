/**
 * Form Analytics Utility - Track and analyze form submissions
 */

/**
 * Analytics event types
 */
export const analyticsEvents = {
	VIEW: 'form_view',
	START: 'form_start',
	SUBMIT: 'form_submit',
	ABANDON: 'form_abandon',
	ERROR: 'form_error',
	FIELD_FOCUS: 'field_focus',
	FIELD_BLUR: 'field_blur',
	FIELD_CHANGE: 'field_change'
};

/**
 * Initialize analytics for a form
 * @param {string} formId - Form ID
 * @returns {Object} Analytics session
 */
export function initAnalytics(formId) {
	const sessionId = generateSessionId();
	const startTime = Date.now();
	
	const session = {
		formId,
		sessionId,
		startTime,
		events: [],
		fieldInteractions: {},
		completionStatus: 'started'
	};
	
	// Store in localStorage
	const sessions = getAnalyticsSessions();
	sessions.push(session);
	saveAnalyticsSessions(sessions);
	
	return session;
}

/**
 * Generate unique session ID
 * @returns {string} Session ID
 */
function generateSessionId() {
	return 'sess_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
}

/**
 * Track analytics event
 * @param {string} formId - Form ID
 * @param {string} eventType - Event type
 * @param {Object} data - Event data
 */
export function trackEvent(formId, eventType, data = {}) {
	const sessions = getAnalyticsSessions();
	const session = sessions.find(s => s.formId === formId && s.completionStatus === 'started');
	
	if (!session) return;
	
	const event = {
		type: eventType,
		timestamp: Date.now(),
		data
	};
	
	session.events.push(event);
	
	// Update completion status
	if (eventType === analyticsEvents.SUBMIT) {
		session.completionStatus = 'completed';
		session.completionTime = Date.now();
	} else if (eventType === analyticsEvents.ABANDON) {
		session.completionStatus = 'abandoned';
	}
	
	saveAnalyticsSessions(sessions);
}

/**
 * Track field interaction
 * @param {string} formId - Form ID
 * @param {string} fieldId - Field ID
 * @param {string} interactionType - Interaction type
 * @param {Object} data - Additional data
 */
export function trackFieldInteraction(formId, fieldId, interactionType, data = {}) {
	const sessions = getAnalyticsSessions();
	const session = sessions.find(s => s.formId === formId && s.completionStatus === 'started');
	
	if (!session) return;
	
	if (!session.fieldInteractions[fieldId]) {
		session.fieldInteractions[fieldId] = {
			focusCount: 0,
			blurCount: 0,
			changeCount: 0,
			timeSpent: 0,
			lastFocusTime: null
		};
	}
	
	const interaction = session.fieldInteractions[fieldId];
	
	switch (interactionType) {
		case 'focus':
			interaction.focusCount++;
			interaction.lastFocusTime = Date.now();
			break;
		case 'blur':
			interaction.blurCount++;
			if (interaction.lastFocusTime) {
				interaction.timeSpent += Date.now() - interaction.lastFocusTime;
			}
			break;
		case 'change':
			interaction.changeCount++;
			break;
	}
	
	saveAnalyticsSessions(sessions);
}

/**
 * Get all analytics sessions
 * @returns {Array} Analytics sessions
 */
export function getAnalyticsSessions() {
	return JSON.parse(localStorage.getItem('form_analytics_sessions') || '[]');
}

/**
 * Save analytics sessions
 * @param {Array} sessions 
 */
function saveAnalyticsSessions(sessions) {
	localStorage.setItem('form_analytics_sessions', JSON.stringify(sessions));
}

/**
 * Get form analytics summary
 * @param {string} formId - Form ID
 * @returns {Object} Analytics summary
 */
export function getFormAnalytics(formId) {
	const sessions = getAnalyticsSessions().filter(s => s.formId === formId);
	
	if (sessions.length === 0) {
		return null;
	}
	
	const totalViews = sessions.length;
	const completedSessions = sessions.filter(s => s.completionStatus === 'completed');
	const abandonedSessions = sessions.filter(s => s.completionStatus === 'abandoned');
	const startedSessions = sessions.filter(s => s.completionStatus === 'started');
	
	const totalSubmissions = completedSessions.length;
	const conversionRate = totalViews > 0 ? (totalSubmissions / totalViews * 100).toFixed(2) : 0;
	const abandonmentRate = totalViews > 0 ? (abandonedSessions.length / totalViews * 100).toFixed(2) : 0;
	
	// Calculate average completion time
	const completionTimes = completedSessions
		.filter(s => s.completionTime && s.startTime)
		.map(s => s.completionTime - s.startTime);
	
	const avgCompletionTime = completionTimes.length > 0
		? completionTimes.reduce((a, b) => a + b, 0) / completionTimes.length
		: 0;
	
	// Field analytics
	const fieldStats = {};
	sessions.forEach(session => {
		Object.entries(session.fieldInteractions).forEach(([fieldId, interactions]) => {
			if (!fieldStats[fieldId]) {
				fieldStats[fieldId] = {
					totalFocuses: 0,
					totalChanges: 0,
					totalTimeSpent: 0,
					interactions: 0
				};
			}
			
			fieldStats[fieldId].totalFocuses += interactions.focusCount;
			fieldStats[fieldId].totalChanges += interactions.changeCount;
			fieldStats[fieldId].totalTimeSpent += interactions.timeSpent;
			fieldStats[fieldId].interactions++;
		});
	});
	
	// Calculate averages per field
	Object.keys(fieldStats).forEach(fieldId => {
		const stats = fieldStats[fieldId];
		stats.avgFocuses = (stats.totalFocuses / stats.interactions).toFixed(2);
		stats.avgChanges = (stats.totalChanges / stats.interactions).toFixed(2);
		stats.avgTimeSpent = (stats.totalTimeSpent / stats.interactions / 1000).toFixed(2); // in seconds
	});
	
	// Time series data (submissions per day)
	const timeSeries = {};
	sessions.forEach(session => {
		const date = new Date(session.startTime).toISOString().split('T')[0];
		if (!timeSeries[date]) {
			timeSeries[date] = { views: 0, submissions: 0 };
		}
		timeSeries[date].views++;
		if (session.completionStatus === 'completed') {
			timeSeries[date].submissions++;
		}
	});
	
	return {
		totalViews,
		totalSubmissions,
		conversionRate,
		abandonmentRate,
		activeSessions: startedSessions.length,
		avgCompletionTime: formatDuration(avgCompletionTime),
		avgCompletionTimeMs: avgCompletionTime,
		fieldStats,
		timeSeries,
		sessions: sessions.map(s => ({
			sessionId: s.sessionId,
			startTime: s.startTime,
			completionStatus: s.completionStatus,
			completionTime: s.completionTime,
			duration: s.completionTime ? s.completionTime - s.startTime : null
		}))
	};
}

/**
 * Format duration in milliseconds to readable string
 * @param {number} ms - Milliseconds
 * @returns {string} Formatted duration
 */
function formatDuration(ms) {
	if (ms === 0) return '0s';
	
	const seconds = Math.floor(ms / 1000);
	const minutes = Math.floor(seconds / 60);
	const hours = Math.floor(minutes / 60);
	
	if (hours > 0) {
		return `${hours}h ${minutes % 60}m`;
	} else if (minutes > 0) {
		return `${minutes}m ${seconds % 60}s`;
	} else {
		return `${seconds}s`;
	}
}

/**
 * Get field drop-off analysis
 * @param {string} formId - Form ID
 * @returns {Array} Field drop-off data
 */
export function getFieldDropOffAnalysis(formId) {
	const sessions = getAnalyticsSessions().filter(s => s.formId === formId);
	const analytics = getFormAnalytics(formId);
	
	if (!analytics) return [];
	
	const fieldIds = Object.keys(analytics.fieldStats);
	const totalSessions = sessions.length;
	
	return fieldIds.map(fieldId => {
		const stats = analytics.fieldStats[fieldId];
		const interactedSessions = sessions.filter(s => s.fieldInteractions[fieldId]).length;
		const dropOffRate = totalSessions > 0 
			? ((totalSessions - interactedSessions) / totalSessions * 100).toFixed(2)
			: 0;
		
		return {
			fieldId,
			fieldName: `Field ${fieldId}`, // Would be replaced with actual field name
			interactions: interactedSessions,
			dropOffRate,
			avgTimeSpent: stats.avgTimeSpent,
			avgChanges: stats.avgChanges
		};
	}).sort((a, b) => parseFloat(b.dropOffRate) - parseFloat(a.dropOffRate));
}

/**
 * Export analytics to CSV
 * @param {string} formId - Form ID
 * @returns {string} CSV content
 */
export function exportAnalyticsToCSV(formId) {
	const analytics = getFormAnalytics(formId);
	if (!analytics) return '';
	
	let csv = 'Date,Views,Submissions,Conversion Rate\n';
	
	Object.entries(analytics.timeSeries).forEach(([date, data]) => {
		const rate = data.views > 0 ? (data.submissions / data.views * 100).toFixed(2) : 0;
		csv += `${date},${data.views},${data.submissions},${rate}%\n`;
	});
	
	return csv;
}

/**
 * Get comparison data (this week vs last week)
 * @param {string} formId - Form ID
 * @returns {Object} Comparison data
 */
export function getWeeklyComparison(formId) {
	const sessions = getAnalyticsSessions().filter(s => s.formId === formId);
	const now = new Date();
	
	const thisWeekStart = new Date(now);
	thisWeekStart.setDate(now.getDate() - 7);
	
	const lastWeekStart = new Date(thisWeekStart);
	lastWeekStart.setDate(thisWeekStart.getDate() - 7);
	
	const thisWeek = sessions.filter(s => s.startTime >= thisWeekStart.getTime());
	const lastWeek = sessions.filter(s => s.startTime >= lastWeekStart.getTime() && s.startTime < thisWeekStart.getTime());
	
	const thisWeekSubmissions = thisWeek.filter(s => s.completionStatus === 'completed').length;
	const lastWeekSubmissions = lastWeek.filter(s => s.completionStatus === 'completed').length;
	
	const change = lastWeekSubmissions > 0
		? ((thisWeekSubmissions - lastWeekSubmissions) / lastWeekSubmissions * 100).toFixed(2)
		: 0;
	
	return {
		thisWeek: {
			views: thisWeek.length,
			submissions: thisWeekSubmissions
		},
		lastWeek: {
			views: lastWeek.length,
			submissions: lastWeekSubmissions
		},
		change: {
			percentage: change,
			direction: change >= 0 ? 'up' : 'down'
		}
	};
}

/**
 * Clear analytics data
 * @param {string} formId - Form ID (optional, clears all if not provided)
 */
export function clearAnalytics(formId = null) {
	if (formId) {
		const sessions = getAnalyticsSessions().filter(s => s.formId !== formId);
		saveAnalyticsSessions(sessions);
	} else {
		localStorage.removeItem('form_analytics_sessions');
	}
}

/**
 * Get top performing forms
 * @returns {Array} Top forms by conversion rate
 */
export function getTopPerformingForms() {
	const sessions = getAnalyticsSessions();
	const formStats = {};
	
	sessions.forEach(session => {
		if (!formStats[session.formId]) {
			formStats[session.formId] = {
				views: 0,
				submissions: 0
			};
		}
		
		formStats[session.formId].views++;
		if (session.completionStatus === 'completed') {
			formStats[session.formId].submissions++;
		}
	});
	
	return Object.entries(formStats)
		.map(([formId, stats]) => ({
			formId,
			views: stats.views,
			submissions: stats.submissions,
			conversionRate: stats.views > 0 ? (stats.submissions / stats.views * 100).toFixed(2) : 0
		}))
		.sort((a, b) => parseFloat(b.conversionRate) - parseFloat(a.conversionRate));
}

export default {
	analyticsEvents,
	initAnalytics,
	trackEvent,
	trackFieldInteraction,
	getFormAnalytics,
	getFieldDropOffAnalysis,
	exportAnalyticsToCSV,
	getWeeklyComparison,
	clearAnalytics,
	getTopPerformingForms
};
