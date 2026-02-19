/**
 * Form Conditional Logic Utility - Show/hide/enable fields based on conditions
 */

/**
 * Condition operators
 */
export const operators = {
	equals: { name: 'Equals', symbol: '==' },
	notEquals: { name: 'Not Equals', symbol: '!=' },
	contains: { name: 'Contains', symbol: 'contains' },
	notContains: { name: 'Does Not Contain', symbol: '!contains' },
	greaterThan: { name: 'Greater Than', symbol: '>' },
	lessThan: { name: 'Less Than', symbol: '<' },
	greaterThanOrEqual: { name: 'Greater Than or Equal', symbol: '>=' },
	lessThanOrEqual: { name: 'Less Than or Equal', symbol: '<=' },
	isEmpty: { name: 'Is Empty', symbol: 'isEmpty' },
	isNotEmpty: { name: 'Is Not Empty', symbol: 'isNotEmpty' },
	startsWith: { name: 'Starts With', symbol: 'startsWith' },
	endsWith: { name: 'Ends With', symbol: 'endsWith' },
	matches: { name: 'Matches Regex', symbol: 'matches' }
};

/**
 * Action types
 */
export const actions = {
	show: { name: 'Show Field', description: 'Show the field when condition is true' },
	hide: { name: 'Hide Field', description: 'Hide the field when condition is true' },
	enable: { name: 'Enable Field', description: 'Enable the field when condition is true' },
	disable: { name: 'Disable Field', description: 'Disable the field when condition is true' },
	require: { name: 'Make Required', description: 'Make field required when condition is true' },
	optional: { name: 'Make Optional', description: 'Make field optional when condition is true' },
	setValue: { name: 'Set Value', description: 'Set a specific value when condition is true' },
	clearValue: { name: 'Clear Value', description: 'Clear the field value when condition is true' }
};

/**
 * Evaluate a single condition
 * @param {any} fieldValue - The field value to check
 * @param {string} operator - The operator to use
 * @param {any} value - The value to compare against
 * @returns {boolean} Whether the condition is met
 */
export function evaluateCondition(fieldValue, operator, value) {
	// Handle empty checks
	if (operator === 'isEmpty') {
		return fieldValue === undefined || fieldValue === null || fieldValue === '' || 
		       (Array.isArray(fieldValue) && fieldValue.length === 0);
	}
	
	if (operator === 'isNotEmpty') {
		return fieldValue !== undefined && fieldValue !== null && fieldValue !== '' &&
		       (!Array.isArray(fieldValue) || fieldValue.length > 0);
	}
	
	// Normalize values for comparison
	const normalizedFieldValue = normalizeValue(fieldValue);
	const normalizedCompareValue = normalizeValue(value);
	
	switch (operator) {
		case 'equals':
			return normalizedFieldValue === normalizedCompareValue;
			
		case 'notEquals':
			return normalizedFieldValue !== normalizedCompareValue;
			
		case 'contains':
			if (Array.isArray(normalizedFieldValue)) {
				return normalizedFieldValue.includes(normalizedCompareValue);
			}
			return String(normalizedFieldValue).includes(String(normalizedCompareValue));
			
		case 'notContains':
			if (Array.isArray(normalizedFieldValue)) {
				return !normalizedFieldValue.includes(normalizedCompareValue);
			}
			return !String(normalizedFieldValue).includes(String(normalizedCompareValue));
			
		case 'greaterThan':
			return parseFloat(normalizedFieldValue) > parseFloat(normalizedCompareValue);
			
		case 'lessThan':
			return parseFloat(normalizedFieldValue) < parseFloat(normalizedCompareValue);
			
		case 'greaterThanOrEqual':
			return parseFloat(normalizedFieldValue) >= parseFloat(normalizedCompareValue);
			
		case 'lessThanOrEqual':
			return parseFloat(normalizedFieldValue) <= parseFloat(normalizedCompareValue);
			
		case 'startsWith':
			return String(normalizedFieldValue).startsWith(String(normalizedCompareValue));
			
		case 'endsWith':
			return String(normalizedFieldValue).endsWith(String(normalizedCompareValue));
			
		case 'matches':
			try {
				const regex = new RegExp(normalizedCompareValue);
				return regex.test(String(normalizedFieldValue));
			} catch {
				return false;
			}
			
		default:
			return false;
	}
}

/**
 * Normalize value for comparison
 * @param {any} value 
 * @returns {any}
 */
function normalizeValue(value) {
	if (value === undefined || value === null) return '';
	if (typeof value === 'boolean') return value;
	if (typeof value === 'number') return value;
	return String(value).trim().toLowerCase();
}

/**
 * Evaluate a rule group (AND/OR logic)
 * @param {Array} conditions - Array of conditions
 * @param {string} logic - 'and' or 'or'
 * @param {Object} formValues - Current form values
 * @returns {boolean} Whether the rule group is satisfied
 */
export function evaluateRuleGroup(conditions, logic, formValues) {
	if (!conditions || conditions.length === 0) return true;
	
	const results = conditions.map(condition => {
		const fieldValue = formValues[condition.fieldId];
		return evaluateCondition(fieldValue, condition.operator, condition.value);
	});
	
	if (logic === 'or') {
		return results.some(r => r);
	}
	
	return results.every(r => r);
}

/**
 * Evaluate all rules for a field
 * @param {Array} rules - Array of conditional rules
 * @param {Object} formValues - Current form values
 * @returns {Object} Result with visibility, enabled, required states
 */
export function evaluateFieldRules(rules, formValues) {
	const result = {
		visible: true,
		enabled: true,
		required: null, // null means use default
		value: null
	};
	
	if (!rules || rules.length === 0) return result;
	
	for (const rule of rules) {
		const conditionMet = evaluateRuleGroup(rule.conditions, rule.logic, formValues);
		
		if (conditionMet) {
			switch (rule.action) {
				case 'show':
					result.visible = true;
					break;
				case 'hide':
					result.visible = false;
					break;
				case 'enable':
					result.enabled = true;
					break;
				case 'disable':
					result.enabled = false;
					break;
				case 'require':
					result.required = true;
					break;
				case 'optional':
					result.required = false;
					break;
				case 'setValue':
					result.value = rule.actionValue;
					break;
				case 'clearValue':
					result.value = null;
					break;
			}
		}
	}
	
	return result;
}

/**
 * Get visible fields based on conditional logic
 * @param {Array} fields - All form fields
 * @param {Object} formValues - Current form values
 * @returns {Array} Visible fields
 */
export function getVisibleFields(fields, formValues) {
	return fields.filter(field => {
		if (!field.conditionalRules || field.conditionalRules.length === 0) {
			return true;
		}
		const result = evaluateFieldRules(field.conditionalRules, formValues);
		return result.visible;
	});
}

/**
 * Get field state (visible, enabled, required)
 * @param {Object} field - Field definition
 * @param {Object} formValues - Current form values
 * @returns {Object} Field state
 */
export function getFieldState(field, formValues) {
	if (!field.conditionalRules || field.conditionalRules.length === 0) {
		return {
			visible: true,
			enabled: true,
			required: field.properties?.required || false
		};
	}
	
	const result = evaluateFieldRules(field.conditionalRules, formValues);
	return {
		visible: result.visible,
		enabled: result.enabled,
		required: result.required !== null ? result.required : (field.properties?.required || false)
	};
}

/**
 * Create a new condition
 * @param {string} fieldId - Field to check
 * @param {string} operator - Operator to use
 * @param {any} value - Value to compare
 * @returns {Object} Condition object
 */
export function createCondition(fieldId, operator, value) {
	return {
		id: Date.now() + Math.random(),
		fieldId,
		operator,
		value
	};
}

/**
 * Create a new rule
 * @param {string} action - Action to take
 * @param {Array} conditions - Conditions to evaluate
 * @param {string} logic - 'and' or 'or'
 * @returns {Object} Rule object
 */
export function createRule(action, conditions = [], logic = 'and') {
	return {
		id: Date.now() + Math.random(),
		action,
		conditions,
		logic,
		actionValue: null // For setValue action
	};
}

/**
 * Get available fields for conditions (exclude current field)
 * @param {Array} allFields - All form fields
 * @param {string} currentFieldId - Current field ID to exclude
 * @returns {Array} Available fields
 */
export function getAvailableConditionFields(allFields, currentFieldId) {
	return allFields.filter(f => f.id !== currentFieldId && isFieldValidForCondition(f));
}

/**
 * Check if field type supports conditional logic
 * @param {Object} field 
 * @returns {boolean}
 */
function isFieldValidForCondition(field) {
	const validTypes = [
		'text', 'textarea', 'email', 'phone', 'url', 'number', 'currency', 
		'percentage', 'select', 'radio', 'checkbox', 'yesno', 'rating', 'scale',
		'date', 'time', 'datetime-local', 'range', 'color'
	];
	return validTypes.includes(field.type);
}

/**
 * Get operators for a field type
 * @param {string} fieldType 
 * @returns {Array} Available operators
 */
export function getOperatorsForFieldType(fieldType) {
	const commonOperators = ['equals', 'notEquals', 'isEmpty', 'isNotEmpty'];
	const textOperators = ['contains', 'notContains', 'startsWith', 'endsWith', 'matches'];
	const numberOperators = ['greaterThan', 'lessThan', 'greaterThanOrEqual', 'lessThanOrEqual'];
	
	switch (fieldType) {
		case 'text':
		case 'textarea':
		case 'email':
		case 'phone':
		case 'url':
			return [...commonOperators, ...textOperators];
			
		case 'number':
		case 'currency':
		case 'percentage':
		case 'rating':
		case 'scale':
		case 'range':
			return [...commonOperators, ...numberOperators];
			
		case 'select':
		case 'radio':
		case 'checkbox':
		case 'yesno':
		case 'color':
			return commonOperators;
			
		case 'date':
		case 'time':
		case 'datetime-local':
			return [...commonOperators, ...numberOperators];
			
		default:
			return commonOperators;
	}
}

/**
 * Apply conditional logic to form values
 * @param {Array} fields - Form fields
 * @param {Object} formValues - Current values
 * @returns {Object} Updated values and field states
 */
export function applyConditionalLogic(fields, formValues) {
	const newValues = { ...formValues };
	const fieldStates = {};
	
	fields.forEach(field => {
		const state = getFieldState(field, formValues);
		fieldStates[field.id] = state;
		
		// Apply value changes from conditional rules
		if (state.value !== null) {
			newValues[field.id] = state.value;
		}
	});
	
	return { values: newValues, fieldStates };
}

/**
 * Serialize conditional rules for backend
 * @param {Array} rules 
 * @returns {string} JSON string
 */
export function serializeRules(rules) {
	return JSON.stringify(rules);
}

/**
 * Deserialize conditional rules from backend
 * @param {string} rulesJson 
 * @returns {Array} Rules array
 */
export function deserializeRules(rulesJson) {
	try {
		return JSON.parse(rulesJson) || [];
	} catch {
		return [];
	}
}

export default {
	operators,
	actions,
	evaluateCondition,
	evaluateRuleGroup,
	evaluateFieldRules,
	getVisibleFields,
	getFieldState,
	createCondition,
	createRule,
	getAvailableConditionFields,
	getOperatorsForFieldType,
	applyConditionalLogic,
	serializeRules,
	deserializeRules
};
