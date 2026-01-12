// Enhanced stores for TPT Titan
import { writable } from 'svelte/store';

// Import templates from separate file
import { templates as templateStore, templateCategories } from './templates.js';

// Export templates for backward compatibility
export { templateStore as templates, templateCategories };

// Chart system store
export const charts = writable([]);

// Plugin system store
export const plugins = writable([
	{
		id: 'currency-converter',
		name: 'Currency Converter',
		description: 'Convert between currencies with live rates',
		enabled: false,
		version: '1.0.0'
	},
	{
		id: 'data-validator',
		name: 'Data Validator',
		description: 'Advanced data validation and cleaning tools',
		enabled: false,
		version: '1.0.0'
	},
	{
		id: 'export-tools',
		name: 'Export Tools',
		description: 'Additional export formats and options',
		enabled: false,
		version: '1.0.0'
	}
]);

// User preferences store
export const userPreferences = writable({
	theme: 'light',
	defaultView: 'spreadsheet',
	autoSave: true,
	showFormulaBar: true,
	enableAnimations: true,
	touchMode: false,
	locale: 'en-US',
	dateFormat: 'MM/DD/YYYY',
	numberFormat: '1,234.56'
});

// Collaboration store (for small office)
export const collaboration = writable({
	enabled: false,
	users: [],
	currentUser: null,
	sharedDocuments: [],
	comments: [],
	versionHistory: []
});

// Contacts store
export const contacts = writable([]);

// Calendar stores
export const calendars = writable([]);
export const events = writable([]);
export const currentView = writable('month');
export const currentDate = writable(new Date());

// Email stores
export const emailSearchQuery = writable('');
export const emailAccounts = writable([]);
export const emails = writable([]);
export const selectedEmail = writable(null);
export const currentFolder = writable('inbox');
