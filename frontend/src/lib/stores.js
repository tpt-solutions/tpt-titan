import { writable } from 'svelte/store';

// User authentication store
export const user = writable(null);

// Application theme store
export const theme = writable('light');

// Loading state store
export const isLoading = writable(false);

// Notification store
export const notifications = writable([]);

// Current workspace/project
export const currentWorkspace = writable(null);

// Application settings
export const settings = writable({
	language: 'en',
	timezone: 'UTC',
	dateFormat: 'MM/DD/YYYY',
});

// Initialize stores from localStorage if available
if (typeof window !== 'undefined') {
	// Load user from localStorage
	const savedUser = localStorage.getItem('user');
	if (savedUser) {
		try {
			user.set(JSON.parse(savedUser));
		} catch (e) {
			console.error('Failed to parse saved user data:', e);
		}
	}

	// Load theme preference
	const savedTheme = localStorage.getItem('theme') || 'light';
	theme.set(savedTheme);

	// Load settings
	const savedSettings = localStorage.getItem('settings');
	if (savedSettings) {
		try {
			settings.set(JSON.parse(savedSettings));
		} catch (e) {
			console.error('Failed to parse saved settings:', e);
		}
	}
}

// Subscribe to store changes and save to localStorage
user.subscribe(value => {
	if (typeof window !== 'undefined') {
		if (value) {
			localStorage.setItem('user', JSON.stringify(value));
		} else {
			localStorage.removeItem('user');
		}
	}
});

theme.subscribe(value => {
	if (typeof window !== 'undefined') {
		localStorage.setItem('theme', value);
		// Apply theme to document
		document.documentElement.classList.toggle('dark', value === 'dark');
	}
});

settings.subscribe(value => {
	if (typeof window !== 'undefined') {
		localStorage.setItem('settings', JSON.stringify(value));
	}
});

// Notification helper functions
export function addNotification(message, type = 'info', duration = 5000) {
	const id = Date.now() + Math.random();
	const notification = { id, message, type, duration };

	notifications.update(current => [...current, notification]);

	// Auto-remove notification
	if (duration > 0) {
		setTimeout(() => {
			removeNotification(id);
		}, duration);
	}

	return id;
}

export function removeNotification(id) {
	notifications.update(current => current.filter(n => n.id !== id));
}

export function clearNotifications() {
	notifications.set([]);
}

// Loading helper functions
export function startLoading() {
	isLoading.set(true);
}

export function stopLoading() {
	isLoading.set(false);
}

// User helper functions
export function setCurrentUser(userData) {
	user.set(userData);
}

export function clearCurrentUser() {
	user.set(null);
}

export function isAuthenticated() {
	let currentUser = null;
	user.subscribe(value => currentUser = value)();
	return currentUser !== null;
}
