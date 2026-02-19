// templates.js - Template aggregator
// All template data lives in the individual category files below.
import { writable } from 'svelte/store';

import { financialTemplates } from './financial-templates.js';
import { operationsTemplates } from './operations-templates.js';
import { lifestyleTemplates } from './lifestyle-templates.js';
import { analyticsTemplates } from './analytics-templates.js';
import { businessTemplates } from './business-templates.js';
import { productivityTemplates } from './productivity-templates.js';
import { healthFitnessTemplates } from './health-fitness-templates.js';
import { educationTemplates } from './education-templates.js';
import { generalTemplates } from './general-templates.js';
import { additionalGeneralTemplates } from './additional-general-templates.js';

export const templateCategories = [
	{ id: 'all',            name: 'All Templates',   icon: '📋' },
	{ id: 'Financial',      name: 'Finance',          icon: '💰' },
	{ id: 'Operations',     name: 'Operations',       icon: '⚙️' },
	{ id: 'Business',       name: 'Business',         icon: '💼' },
	{ id: 'Productivity',   name: 'Productivity',     icon: '✅' },
	{ id: 'Health & Fitness', name: 'Health & Fitness', icon: '💪' },
	{ id: 'Lifestyle',      name: 'Lifestyle',        icon: '🏠' },
	{ id: 'Education',      name: 'Education',        icon: '📚' },
	{ id: 'Analytics',      name: 'Analytics',        icon: '📊' }
];

export const templates = writable([
	...financialTemplates,
	...operationsTemplates,
	...lifestyleTemplates,
	...analyticsTemplates,
	...businessTemplates,
	...productivityTemplates,
	...healthFitnessTemplates,
	...educationTemplates,
	...generalTemplates,
	...additionalGeneralTemplates
]);
