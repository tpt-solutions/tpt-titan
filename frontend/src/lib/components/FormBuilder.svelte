<!-- frontend/src/lib/components/FormBuilder.svelte -->
<!-- Orchestrator: wires together all FormBuilder sub-components. -->
<script>
// @ts-nocheck
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import { createField } from './FormBuilderFieldTypes.js';
	import FormBuilderProperties    from './FormBuilderProperties.svelte';
	import FormBuilderCanvas        from './FormBuilderCanvas.svelte';
	import FormBuilderFieldPanel    from './FormBuilderFieldPanel.svelte';
	import FormBuilderTemplatesModal   from './FormBuilderTemplatesModal.svelte';
	import FormBuilderResponsesModal   from './FormBuilderResponsesModal.svelte';
	import FormBuilderConditionalModal from './FormBuilderConditionalModal.svelte';
	import FormHistory, { createDebouncedPush, handleUndoRedoKeyboard } from '../utils/form-history.js';


	/** @type {any} */
	export let form = null; // Existing form to edit, or null for new

	const dispatch = createEventDispatcher();

	// ── History management ────────────────────────────────────────────────────
	let formHistory = new FormHistory();
	/** @type {any} */
	let debouncedHistoryPush;
	let canUndo = false;
	let canRedo = false;

	// Update undo/redo availability
	function updateHistoryState() {
		canUndo = formHistory.canUndo();
		canRedo = formHistory.canRedo();
	}

	// Save current state to history
	function saveToHistory(actionType = 'fieldChange', description = 'Field updated') {
		debouncedHistoryPush({
			fields: formData.fields,
			settings: formData.settings,
			actionType,
			description
		});
		updateHistoryState();
	}

	// Undo handler
	function handleUndo() {
		const state = /** @type {any} */ (formHistory.undo());
		if (state) {
			formData.fields = state.fields;
			formData.settings = state.settings;
			selectedField = null;
			updateHistoryState();
		}
	}

	// Redo handler
	function handleRedo() {
		const state = /** @type {any} */ (formHistory.redo());
		if (state) {
			formData.fields = state.fields;
			formData.settings = state.settings;
			selectedField = null;
			updateHistoryState();
		}
	}

	// Global keyboard handler for undo/redo
	/** @param {any} event */
	function handleGlobalKeyDown(event) {
		if (handleUndoRedoKeyboard(event, handleUndo, handleRedo)) {
			return;
		}
	}

	// ── Form data ──────────────────────────────────────────────────────────────

	let formData = {
		name:        form?.name        || '',
		description: form?.description || '',
		fields:      form?.fields      || [],
		settings:    form?.settings    || {
			allowMultipleResponses: false,
			requireLogin:           false,
			emailNotifications:     true,
			autoSave:               true
		}
	};

	// ── Selection & drag state ─────────────────────────────────────────────────
	/** @type {any} */
	let selectedField    = null;
	/** @type {any} */
	let draggedFieldType = null; // dragging a new field type from the sidebar
	/** @type {any} */
	let draggedField     = null; // reordering an existing field on the canvas
	let isDragging       = false;

	// ── Modal visibility ───────────────────────────────────────────────────────
	let showTemplates  = false;
	let showConditional = false;
	let showResponses  = false;

	// ── Response data (loaded lazily) ──────────────────────────────────────────
	/** @type {any[]} */
	let formResponses = [];

	// ── Field CRUD ─────────────────────────────────────────────────────────────
	/** @param {any} fieldType */
	function addField(fieldType) {
		const f = createField(fieldType, formData.fields.length);
		formData.fields = [...formData.fields, f];
		selectedField = f;
		saveToHistory('addField', `Added ${f.label} field`);
	}

	/** @param {any} fieldId */
	function deleteField(fieldId) {
		const field = formData.fields.find(f => f.id === fieldId);
		formData.fields = formData.fields.filter(f => f.id !== fieldId);
		if (selectedField?.id === fieldId) selectedField = null;
		saveToHistory('deleteField', `Deleted ${field?.label || 'field'}`);
	}

	/**
	 * @param {any} fieldId
	 * @param {any} label
	 */
	function updateFieldLabel(fieldId, label) {
		formData.fields = formData.fields.map(f => f.id === fieldId ? { ...f, label } : f);
		if (selectedField?.id === fieldId) selectedField.label = label;
		saveToHistory('updateLabel', `Updated label to "${label}"`);
	}

	/**
	 * @param {any} fieldId
	 * @param {any} property
	 * @param {any} value
	 */
	function updateFieldProperty(fieldId, property, value) {
		formData.fields = formData.fields.map(f =>
			f.id === fieldId ? { ...f, properties: { ...f.properties, [property]: value } } : f
		);
		if (selectedField?.id === fieldId) selectedField.properties[property] = value;
		saveToHistory('updateProperty', `Updated ${property}`);
	}

	/**
	 * @param {any} property
	 * @param {any} value
	 */
	function updateSetting(property, value) {
		formData.settings = { ...formData.settings, [property]: value };
		saveToHistory('updateSetting', `Updated ${property}`);
	}

	// ── Calculated Fields ──────────────────────────────────────────────────────
	function executeCalculatedFields() {
		const calculatedFields = formData.fields.filter(f => f.type === 'calculation');
		if (calculatedFields.length === 0) return;

		/** @type {Record<string, any>} */
		const fieldValues = {};
		formData.fields.forEach(f => {
			fieldValues[f.label] = f.properties.value || 0;
		});

		calculatedFields.forEach(field => {
			try {
				let formula = field.properties.formula || '';
				// Replace field references {FieldName} with values
				formula = formula.replace(/\{([^}]+)\}/g, (match, fieldName) => {
					const value = fieldValues[fieldName.trim()];
					return value !== undefined ? value : 0;
				});

				// Basic math functions
				formula = formula.replace(/SUM\(([^)]+)\)/g, (match, args) => {
					return args.split(',').reduce((sum, val) => sum + (parseFloat(val) || 0), 0);
				});
				formula = formula.replace(/AVG\(([^)]+)\)/g, (match, args) => {
					const vals = args.split(',');
					return vals.reduce((sum, val) => sum + (parseFloat(val) || 0), 0) / vals.length;
				});
				formula = formula.replace(/MIN\(([^)]+)\)/g, (match, args) => {
					return Math.min(...args.split(',').map(v => parseFloat(v) || 0));
				});
				formula = formula.replace(/MAX\(([^)]+)\)/g, (match, args) => {
					return Math.max(...args.split(',').map(v => parseFloat(v) || 0));
				});

				// Evaluate safely
				const result = Function('"use strict"; return (' + formula + ')')();
				field.properties.value = parseFloat(result.toFixed(field.properties.decimalPlaces || 2));
			} catch (e) {
				console.error('Formula error:', e);
				field.properties.value = null;
			}
		});
	}

	// ── Reorder helpers ────────────────────────────────────────────────────────

	/**
	 * @param {any} fromIndex
	 * @param {any} toIndex
	 */
	function reorderFields(fromIndex, toIndex) {
		if (fromIndex === toIndex) return;
		const fields = [...formData.fields];
		const [moved] = fields.splice(fromIndex, 1);
		fields.splice(toIndex, 0, moved);
		fields.forEach((f, i) => { f.order = i; });
		formData.fields = fields;
		saveToHistory('reorderFields', 'Reordered fields');
	}


	// ── Canvas drag events ─────────────────────────────────────────────────────
	/** @param {{ detail: { field: any } }} e */
	function onFieldDragStart({ detail: { field } }) {
		draggedField = field;
		isDragging   = true;
	}

	function onFieldDragEnd() {
		draggedField = null;
		isDragging   = false;
	}

	/** @param {{ detail: { event: any, dropIndex: any } }} e */
	function onFieldDrop({ detail: { event, dropIndex } }) {
		const draggedId  = event.dataTransfer.getData('text/plain');
		const fromIndex  = formData.fields.findIndex(f => f.id.toString() === draggedId);
		if (fromIndex !== -1) reorderFields(fromIndex, dropIndex);
		draggedField = null;
		isDragging   = false;
	}

	/** @param {{ detail: { event: any, dropIndex: any } }} e */
	function onCanvasDrop({ detail: { event, dropIndex } }) {
		event.preventDefault();
		if (draggedFieldType) {
			const f = createField(draggedFieldType, dropIndex);
			formData.fields.splice(dropIndex, 0, f);
			formData.fields.forEach((field, i) => { field.order = i; });
			formData.fields = [...formData.fields];
			selectedField    = f;
		}
		draggedFieldType = null;
		isDragging       = false;
	}

	// ── Field panel drag events ────────────────────────────────────────────────
	/** @param {{ detail: { fieldType: any } }} e */
	function onPanelDragStart({ detail: { fieldType } }) {
		draggedFieldType = fieldType;
		isDragging       = true;
	}

	function onPanelDragEnd() {
		draggedFieldType = null;
		isDragging       = false;
	}

	// ── Templates ──────────────────────────────────────────────────────────────
	/** @param {{ detail: any }} e */
	function applyTemplate({ detail: template }) {
		formData.fields = template.fields.map((f, i) => ({
			id:         Date.now() + Math.random() + i,
			type:       f.type,
			label:      f.label,
			properties: { ...f.properties },
			order:      i
		}));
		showTemplates = false;
		saveToHistory('applyTemplate', `Applied ${template.name} template`);
	}


	// ── Responses ──────────────────────────────────────────────────────────────
	async function loadResponses() {
		formResponses = [
			{ id: 1, submitted_at: new Date().toISOString(),              responses: { name: 'John Doe',   email: 'john@example.com', message: 'Great form!' } },
			{ id: 2, submitted_at: new Date(Date.now() - 3600000).toISOString(), responses: { name: 'Jane Smith', email: 'jane@example.com', message: 'Very useful form.' } }
		];
		showResponses = true;
	}

	// ── Save / Cancel ──────────────────────────────────────────────────────────
	function saveForm() {
		if (!formData.name.trim()) { alert('Please enter a form name'); return; }
		if (formData.fields.length === 0) { alert('Please add at least one field to your form'); return; }
		dispatch('save', formData);
	}

	// ── Lifecycle ─────────────────────────────────────────────────────────────
	onMount(() => {
		// Initialize history with debounced push
		debouncedHistoryPush = createDebouncedPush(formHistory, 500);
		
		// Add keyboard listener for undo/redo
		document.addEventListener('keydown', handleGlobalKeyDown);
		
		// Save initial state
		setTimeout(() => saveToHistory('initial', 'Initial state'), 0);
	});

	onDestroy(() => {
		// Cleanup
		document.removeEventListener('keydown', handleGlobalKeyDown);
	});
</script>


<div class="h-full flex flex-col">
	<!-- Toolbar with Undo/Redo -->
	<div class="flex items-center space-x-2 px-4 py-2 border-b border-gray-200 bg-gray-50">
		<button
			class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed"
			on:click={handleUndo}
			disabled={!canUndo}
			title="Undo (Ctrl+Z)"
		>
			↶ Undo
		</button>
		<button
			class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed"
			on:click={handleRedo}
			disabled={!canRedo}
			title="Redo (Ctrl+Y)"
		>
			↷ Redo
		</button>
		<div class="flex-1"></div>
		<span class="text-sm text-gray-500">
			{formData.fields.length} field{formData.fields.length !== 1 ? 's' : ''}
		</span>
	</div>

	<div class="h-full flex flex-1">
		<!-- Left: settings + field properties -->
	<FormBuilderProperties
		{formData}
		{selectedField}
		on:showTemplates={() => showTemplates = true}
		on:showConditional={() => showConditional = true}
		on:showResponses={loadResponses}
		on:updateLabel={({ detail }) => updateFieldLabel(detail.fieldId, detail.label)}
		on:updateProperty={({ detail }) => updateFieldProperty(detail.fieldId, detail.property, detail.value)}
		on:updateSetting={({ detail }) => updateSetting(detail.property, detail.value)}
	/>



	<!-- Centre: canvas -->
	<FormBuilderCanvas
		{formData}
		{selectedField}
		{isDragging}
		{draggedField}
		on:selectField={({ detail }) => selectedField = detail}
		on:deleteField={({ detail }) => deleteField(detail)}
		on:fieldDragStart={onFieldDragStart}
		on:fieldDragEnd={onFieldDragEnd}
		on:fieldDrop={onFieldDrop}
		on:canvasDrop={onCanvasDrop}
		on:save={saveForm}
		on:cancel={() => dispatch('cancel')}
	/>

		<!-- Right: draggable field-type catalogue -->
		<FormBuilderFieldPanel
			{isDragging}
			{draggedFieldType}
			on:dragStart={onPanelDragStart}
			on:dragEnd={onPanelDragEnd}
			on:quickAdd={({ detail }) => addField(detail)}
		/>
	</div>
</div>


<!-- Modals (rendered at root level so they overlay everything) -->
{#if showTemplates}
	<FormBuilderTemplatesModal
		on:apply={applyTemplate}
		on:close={() => showTemplates = false}
	/>
{/if}

{#if showResponses}
	<FormBuilderResponsesModal
		{formData}
		{formResponses}
		on:close={() => showResponses = false}
	/>
{/if}

{#if showConditional}
	<FormBuilderConditionalModal
		{selectedField}
		allFields={formData.fields}
		on:updated={() => { formData.fields = [...formData.fields]; }}
		on:close={() => showConditional = false}
	/>
{/if}
