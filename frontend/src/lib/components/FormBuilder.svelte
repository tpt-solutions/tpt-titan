<script>
	import { createEventDispatcher } from 'svelte';

	export let form = null; // Existing form to edit, or null for new

	const dispatch = createEventDispatcher();

	// Form response management
	let showResponseManager = false;
	let formResponses = [];
	let selectedResponse = null;
	let showFormTemplates = false;
	let showConditionalLogic = false;

	let formData = {
		name: form?.name || '',
		description: form?.description || '',
		fields: form?.fields || [],
		settings: form?.settings || {
			allowMultipleResponses: false,
			requireLogin: false,
			emailNotifications: true,
			autoSave: true
		}
	};

	let selectedField = null;
	let draggedField = null;
	let isDragging = false;

	// Available field types
	const fieldTypes = [
		{ id: 'text', name: 'Text Input', icon: '📝', properties: { placeholder: '', required: false, maxLength: 100 } },
		{ id: 'textarea', name: 'Long Text', icon: '📄', properties: { placeholder: '', required: false, rows: 3 } },
		{ id: 'number', name: 'Number', icon: '🔢', properties: { placeholder: '', required: false, min: null, max: null } },
		{ id: 'email', name: 'Email', icon: '📧', properties: { placeholder: '', required: false } },
		{ id: 'phone', name: 'Phone', icon: '📱', properties: { placeholder: '', required: false } },
		{ id: 'date', name: 'Date', icon: '📅', properties: { required: false } },
		{ id: 'select', name: 'Dropdown', icon: '▼', properties: { options: ['Option 1', 'Option 2'], required: false } },
		{ id: 'radio', name: 'Radio Buttons', icon: '⭕', properties: { options: ['Option 1', 'Option 2'], required: false } },
		{ id: 'checkbox', name: 'Checkboxes', icon: '☑️', properties: { options: ['Option 1', 'Option 2'], required: false } },
		{ id: 'file', name: 'File Upload', icon: '📎', properties: { accept: '*', maxSize: 10, required: false } },
		{ id: 'signature', name: 'Digital Signature', icon: '✍️', properties: { required: false, legalText: '' } },
		{ id: 'rating', name: 'Rating', icon: '⭐', properties: { maxRating: 5, required: false } }
	];

	function addField(fieldType) {
		const newField = {
			id: Date.now() + Math.random(),
			type: fieldType.id,
			label: fieldType.name,
			properties: { ...fieldType.properties },
			order: formData.fields.length
		};
		formData.fields = [...formData.fields, newField];
		selectedField = newField;
	}

	function deleteField(fieldId) {
		formData.fields = formData.fields.filter(f => f.id !== fieldId);
		if (selectedField?.id === fieldId) {
			selectedField = null;
		}
	}

	function selectField(field) {
		selectedField = field;
	}

	function updateFieldProperty(fieldId, property, value) {
		formData.fields = formData.fields.map(field =>
			field.id === fieldId
				? { ...field, properties: { ...field.properties, [property]: value } }
				: field
		);
		if (selectedField?.id === fieldId) {
			selectedField.properties[property] = value;
		}
	}

	function updateFieldLabel(fieldId, label) {
		formData.fields = formData.fields.map(field =>
			field.id === fieldId ? { ...field, label } : field
		);
		if (selectedField?.id === fieldId) {
			selectedField.label = label;
		}
	}

	function moveField(fromIndex, toIndex) {
		const fields = [...formData.fields];
		const [movedField] = fields.splice(fromIndex, 1);
		fields.splice(toIndex, 0, movedField);

		// Update order property
		fields.forEach((field, index) => {
			field.order = index;
		});

		formData.fields = fields;
	}

	function saveForm() {
		// Basic validation
		if (!formData.name.trim()) {
			alert('Please enter a form name');
			return;
		}

		if (formData.fields.length === 0) {
			alert('Please add at least one field to your form');
			return;
		}

		dispatch('save', formData);
	}

	function cancel() {
		dispatch('cancel');
	}

	// Response management functions
	async function loadFormResponses() {
		// This would connect to the backend to load responses
		// For now, showing mock data
		formResponses = [
			{
				id: 1,
				submitted_at: new Date().toISOString(),
				responses: {
					name: 'John Doe',
					email: 'john@example.com',
					message: 'Great form!'
				}
			},
			{
				id: 2,
				submitted_at: new Date(Date.now() - 3600000).toISOString(),
				responses: {
					name: 'Jane Smith',
					email: 'jane@example.com',
					message: 'Very useful form.'
				}
			}
		];
		showResponseManager = true;
	}

	function exportResponses() {
		let csv = 'ID,Submitted At';
		// Add field headers
		formData.fields.forEach(field => {
			csv += ',' + field.label;
		});
		csv += '\n';

		// Add response data
		formResponses.forEach(response => {
			csv += response.id + ',' + response.submitted_at;
			formData.fields.forEach(field => {
				const value = response.responses[field.label.toLowerCase()] || '';
				csv += ',"' + value.replace(/"/g, '""') + '"';
			});
			csv += '\n';
		});

		const blob = new Blob([csv], { type: 'text/csv' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `${formData.name || 'form'}_responses.csv`;
		a.click();
		URL.revokeObjectURL(url);
	}

	// Form templates
	const formTemplates = [
		{
			id: 'contact',
			name: 'Contact Form',
			description: 'Basic contact information collection',
			icon: '📞',
			fields: [
				{ type: 'text', label: 'Full Name', properties: { required: true } },
				{ type: 'email', label: 'Email Address', properties: { required: true } },
				{ type: 'phone', label: 'Phone Number', properties: { required: false } },
				{ type: 'textarea', label: 'Message', properties: { required: true, rows: 4 } }
			]
		},
		{
			id: 'survey',
			name: 'Customer Survey',
			description: 'Gather feedback from customers',
			icon: '📊',
			fields: [
				{ type: 'text', label: 'Name', properties: { required: false } },
				{ type: 'email', label: 'Email', properties: { required: false } },
				{ type: 'radio', label: 'How satisfied are you?', properties: { options: ['Very Satisfied', 'Satisfied', 'Neutral', 'Dissatisfied', 'Very Dissatisfied'], required: true } },
				{ type: 'textarea', label: 'Comments', properties: { required: false, rows: 3 } }
			]
		},
		{
			id: 'registration',
			name: 'Event Registration',
			description: 'Register attendees for events',
			icon: '🎟️',
			fields: [
				{ type: 'text', label: 'Full Name', properties: { required: true } },
				{ type: 'email', label: 'Email Address', properties: { required: true } },
				{ type: 'phone', label: 'Phone Number', properties: { required: true } },
				{ type: 'select', label: 'Ticket Type', properties: { options: ['General Admission', 'VIP', 'Student'], required: true } },
				{ type: 'number', label: 'Number of Tickets', properties: { required: true, min: 1, max: 10 } },
				{ type: 'textarea', label: 'Special Requests', properties: { required: false, rows: 2 } }
			]
		},
		{
			id: 'feedback',
			name: 'Product Feedback',
			description: 'Collect feedback on products or services',
			icon: '💬',
			fields: [
				{ type: 'text', label: 'Product Name', properties: { required: true } },
				{ type: 'rating', label: 'Overall Rating', properties: { required: true, maxRating: 5 } },
				{ type: 'radio', label: 'Would you recommend?', properties: { options: ['Yes', 'No', 'Maybe'], required: true } },
				{ type: 'textarea', label: 'What did you like?', properties: { required: false, rows: 2 } },
				{ type: 'textarea', label: 'What could be improved?', properties: { required: false, rows: 2 } }
			]
		}
	];

	function applyTemplate(template) {
		formData.fields = template.fields.map((field, index) => ({
			id: Date.now() + Math.random() + index,
			type: field.type,
			label: field.label,
			properties: { ...field.properties },
			order: index
		}));
		showFormTemplates = false;
	}

	// Conditional logic
	function addConditionalLogic() {
		if (!selectedField) return;

		// Add conditional logic to the selected field
		if (!selectedField.conditionalLogic) {
			selectedField.conditionalLogic = [];
		}

		selectedField.conditionalLogic.push({
			field: '',
			operator: 'equals',
			value: '',
			action: 'show'
		});

		showConditionalLogic = true;
	}

	function updateConditionalRule(ruleIndex, property, value) {
		if (selectedField?.conditionalLogic?.[ruleIndex]) {
			selectedField.conditionalLogic[ruleIndex][property] = value;
			formData.fields = [...formData.fields]; // Trigger reactivity
		}
	}

	function removeConditionalRule(ruleIndex) {
		if (selectedField?.conditionalLogic) {
			selectedField.conditionalLogic.splice(ruleIndex, 1);
			formData.fields = [...formData.fields];
		}
	}

	// Drag and drop handlers
	function handleDragStart(event, fieldType) {
		draggedField = fieldType;
		isDragging = true;
	}

	function handleDragEnd() {
		draggedField = null;
		isDragging = false;
	}

	function handleDrop(event, dropIndex) {
		event.preventDefault();
		if (draggedField) {
			const newField = {
				id: Date.now() + Math.random(),
				type: draggedField.id,
				label: draggedField.name,
				properties: { ...draggedField.properties },
				order: dropIndex
			};

			// Insert at the drop position
			formData.fields.splice(dropIndex, 0, newField);

			// Update order for all fields
			formData.fields.forEach((field, index) => {
				field.order = index;
			});

			formData.fields = [...formData.fields];
			selectedField = newField;
		}
		isDragging = false;
		draggedField = null;
	}
</script>

<div class="h-full flex">
	<!-- Left Sidebar - Form Settings -->
	<div class="w-80 bg-white border-r border-gray-200 flex flex-col">
		<div class="p-4 border-b border-gray-200">
			<h3 class="text-lg font-semibold text-gray-900 mb-3">Form Settings</h3>
			<div class="space-y-3">
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Form Name</label>
					<input
						bind:value={formData.name}
						placeholder="Enter form name"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
					<textarea
						bind:value={formData.description}
						placeholder="Describe your form"
						rows="3"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
					></textarea>
				</div>
			</div>
		</div>

		<!-- Field Properties Panel -->
		{#if selectedField}
			<div class="flex-1 p-4 overflow-y-auto">
				<h3 class="text-lg font-semibold text-gray-900 mb-3">Field Properties</h3>
				<div class="space-y-4">
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-1">Label</label>
						<input
							bind:value={selectedField.label}
							on:input={(e) => updateFieldLabel(selectedField.id, e.target.value)}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>

					<!-- Type-specific properties -->
					{#if selectedField.type === 'text' || selectedField.type === 'textarea' || selectedField.type === 'email'}
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Placeholder</label>
							<input
								bind:value={selectedField.properties.placeholder}
								on:input={(e) => updateFieldProperty(selectedField.id, 'placeholder', e.target.value)}
								class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
					{/if}

					{#if selectedField.type === 'number'}
						<div class="grid grid-cols-2 gap-2">
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-1">Min</label>
								<input
									type="number"
									bind:value={selectedField.properties.min}
									on:input={(e) => updateFieldProperty(selectedField.id, 'min', e.target.value ? parseFloat(e.target.value) : null)}
									class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-1">Max</label>
								<input
									type="number"
									bind:value={selectedField.properties.max}
									on:input={(e) => updateFieldProperty(selectedField.id, 'max', e.target.value ? parseFloat(e.target.value) : null)}
									class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
						</div>
					{/if}

					{#if ['select', 'radio', 'checkbox'].includes(selectedField.type)}
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Options</label>
							<textarea
								bind:value={selectedField.properties.options}
								on:input={(e) => updateFieldProperty(selectedField.id, 'options', e.target.value.split('\n').filter(opt => opt.trim()))}
								placeholder="One option per line"
								rows="4"
								class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							></textarea>
						</div>
					{/if}

					{#if selectedField.type === 'file'}
						<div class="grid grid-cols-2 gap-2">
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-1">Max Size (MB)</label>
								<input
									type="number"
									bind:value={selectedField.properties.maxSize}
									on:input={(e) => updateFieldProperty(selectedField.id, 'maxSize', parseInt(e.target.value))}
									class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-1">File Types</label>
								<input
									bind:value={selectedField.properties.accept}
									on:input={(e) => updateFieldProperty(selectedField.id, 'accept', e.target.value)}
									placeholder="e.g., .pdf,.doc"
									class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
						</div>
					{/if}

					<!-- Common properties -->
					<div class="flex items-center">
						<input
							type="checkbox"
							id="required-{selectedField.id}"
							bind:checked={selectedField.properties.required}
							on:change={(e) => updateFieldProperty(selectedField.id, 'required', e.target.checked)}
							class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
						/>
						<label for="required-{selectedField.id}" class="ml-2 text-sm text-gray-700">
							Required field
						</label>
					</div>
				</div>
			</div>
		{:else}
			<div class="flex-1 p-4 flex items-center justify-center text-gray-500">
				<div class="text-center">
					<svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
					</svg>
					<p>Select a field to edit its properties</p>
				</div>
			</div>
		{/if}
	</div>

	<!-- Main Canvas - Form Preview -->
	<div class="flex-1 bg-gray-50 overflow-y-auto">
		<div class="max-w-2xl mx-auto py-8">
			<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-8">
				<!-- Form Header -->
				<div class="mb-8">
					<h2 class="text-2xl font-bold text-gray-900 mb-2">
						{formData.name || 'Untitled Form'}
					</h2>
					{#if formData.description}
						<p class="text-gray-600">{formData.description}</p>
					{/if}
				</div>

				<!-- Form Fields -->
				{#each formData.fields as field, index}
					<div
						class="mb-6 p-4 border-2 border-dashed border-gray-300 rounded-lg cursor-move hover:border-blue-400 transition-colors {selectedField?.id === field.id ? 'border-blue-500 bg-blue-50' : ''}"
						on:click={() => selectField(field)}
						on:drop={(e) => handleDrop(e, index)}
						on:dragover={(e) => e.preventDefault()}
					>
						<div class="flex items-center justify-between mb-2">
							<span class="text-sm font-medium text-gray-700">{field.label}</span>
							<button
								class="text-red-500 hover:text-red-700 p-1"
								on:click={(e) => { e.stopPropagation(); deleteField(field.id); }}
								title="Delete field"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
								</svg>
							</button>
						</div>

						<!-- Field Preview -->
						<div class="space-y-2">
							{#if field.type === 'text'}
								<input
									type="text"
									placeholder={field.properties.placeholder || 'Enter text...'}
									disabled
									class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50"
								/>
							{:else if field.type === 'textarea'}
								<textarea
									placeholder={field.properties.placeholder || 'Enter text...'}
									rows={field.properties.rows || 3}
									disabled
									class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50"
								></textarea>
							{:else if field.type === 'number'}
								<input
									type="number"
									placeholder={field.properties.placeholder || 'Enter number...'}
									min={field.properties.min}
									max={field.properties.max}
									disabled
									class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50"
								/>
							{:else if field.type === 'email'}
								<input
									type="email"
									placeholder={field.properties.placeholder || 'Enter email...'}
									disabled
									class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50"
								/>
							{:else if field.type === 'date'}
								<input
									type="date"
									disabled
									class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50"
								/>
							{:else if field.type === 'select'}
								<select disabled class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50">
									<option>Select an option...</option>
									{#each field.properties.options || [] as option}
										<option>{option}</option>
									{/each}
								</select>
							{:else if field.type === 'radio'}
								<div class="space-y-2">
									{#each field.properties.options || [] as option}
										<label class="flex items-center">
											<input type="radio" disabled class="mr-2" />
											<span>{option}</span>
										</label>
									{/each}
								</div>
							{:else if field.type === 'checkbox'}
								<div class="space-y-2">
									{#each field.properties.options || [] as option}
										<label class="flex items-center">
											<input type="checkbox" disabled class="mr-2" />
											<span>{option}</span>
										</label>
									{/each}
								</div>
							{:else if field.type === 'signature'}
								<div class="border-2 border-dashed border-gray-300 rounded p-8 text-center text-gray-500">
									Signature Area
								</div>
							{:else}
								<div class="text-gray-500 italic">[{field.type} field preview]</div>
							{/if}

							{#if field.properties.required}
								<span class="text-red-500 text-sm">* Required</span>
							{/if}
						</div>
					</div>
				{/each}

				<!-- Drop zone for new fields -->
				{#if formData.fields.length === 0}
					<div
						class="mb-6 p-8 border-2 border-dashed border-gray-300 rounded-lg text-center text-gray-500 hover:border-blue-400 transition-colors"
						on:drop={(e) => handleDrop(e, 0)}
						on:dragover={(e) => e.preventDefault()}
					>
						<svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
						</svg>
						<p>Drag field types here to start building your form</p>
					</div>
				{/if}

				<!-- Form Actions -->
				<div class="flex justify-between items-center pt-6 border-t border-gray-200">
					<div class="text-sm text-gray-600">
						{formData.fields.length} field{formData.fields.length !== 1 ? 's' : ''}
					</div>
					<div class="flex space-x-3">
						<button
							on:click={cancel}
							class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
						>
							Cancel
						</button>
						<button
							on:click={saveForm}
							class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
						>
							Save Form
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Right Sidebar - Field Types -->
	<div class="w-64 bg-white border-l border-gray-200 p-4">
		<h3 class="text-lg font-semibold text-gray-900 mb-4">Field Types</h3>
		<div class="space-y-2">
			{#each fieldTypes as fieldType}
				<div
					class="p-3 border border-gray-200 rounded-lg cursor-move hover:bg-gray-50 transition-colors {isDragging && draggedField?.id === fieldType.id ? 'opacity-50' : ''}"
					draggable="true"
					on:dragstart={(e) => handleDragStart(e, fieldType)}
					on:dragend={handleDragEnd}
				>
					<div class="flex items-center space-x-3">
						<span class="text-lg">{fieldType.icon}</span>
						<div>
							<div class="font-medium text-sm">{fieldType.name}</div>
							<div class="text-xs text-gray-500">Drag to add</div>
						</div>
					</div>
				</div>
			{/each}
		</div>

		<!-- Quick Add Buttons -->
		<div class="mt-6 pt-4 border-t border-gray-200">
			<h4 class="text-sm font-medium text-gray-900 mb-2">Quick Add</h4>
			<div class="space-y-1">
				{#each fieldTypes.slice(0, 3) as fieldType}
					<button
						on:click={() => addField(fieldType)}
						class="w-full text-left p-2 text-sm hover:bg-gray-50 rounded transition-colors"
					>
						+ {fieldType.name}
					</button>
				{/each}
			</div>
		</div>
	</div>
</div>

<!-- Form Templates Modal -->
{#if showFormTemplates}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<div>
					<h2 class="text-2xl font-bold text-gray-900">Choose a Form Template</h2>
					<p class="text-gray-600 mt-1">Start with a pre-built template and customize it for your needs</p>
				</div>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showFormTemplates = false}
				>
					×
				</button>
			</div>

			<div class="flex-1 overflow-y-auto p-6">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					{#each formTemplates as template}
						<div
							class="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-lg transition-all cursor-pointer group"
							on:click={() => applyTemplate(template)}
						>
							<div class="flex items-center mb-4">
								<div class="text-3xl mr-3">{template.icon}</div>
								<div>
									<h3 class="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">
										{template.name}
									</h3>
									<p class="text-sm text-gray-600">{template.description}</p>
								</div>
							</div>

							<div class="space-y-2">
								{#each template.fields.slice(0, 4) as field}
									<div class="flex items-center text-sm text-gray-600">
										<span class="w-4 h-4 mr-2 bg-gray-200 rounded flex items-center justify-center text-xs">
											{field.type === 'text' ? '📝' : field.type === 'email' ? '📧' : field.type === 'radio' ? '⭕' : '📄'}
										</span>
										<span>{field.label}</span>
										{#if field.properties.required}
											<span class="text-red-500 ml-1">*</span>
										{/if}
									</div>
								{/each}
								{#if template.fields.length > 4}
									<div class="text-sm text-gray-500">
										+{template.fields.length - 4} more fields
									</div>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			</div>

			<div class="border-t border-gray-200 p-6 bg-gray-50">
				<div class="flex items-center justify-between text-sm text-gray-600">
					<span>{formTemplates.length} templates available</span>
					<span>💡 Templates can be fully customized after selection</span>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Response Manager Modal -->
{#if showResponseManager}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-6xl max-h-[90vh] overflow-hidden flex flex-col">
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<div>
					<h2 class="text-2xl font-bold text-gray-900">Form Responses</h2>
					<p class="text-gray-600 mt-1">{formResponses.length} responses received</p>
				</div>
				<div class="flex items-center space-x-3">
					<button
						class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
						on:click={exportResponses}
					>
						📊 Export CSV
					</button>
					<button
						class="text-gray-400 hover:text-gray-600 text-2xl"
						on:click={() => showResponseManager = false}
					>
						×
					</button>
				</div>
			</div>

			<div class="flex-1 overflow-y-auto">
				{#if formResponses.length === 0}
					<div class="text-center py-12">
						<div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
							<svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
							</svg>
						</div>
						<h3 class="text-lg font-medium text-gray-900 mb-2">No responses yet</h3>
						<p class="text-gray-500">Responses will appear here once people submit your form</p>
					</div>
				{:else}
					<div class="p-6">
						<!-- Summary Stats -->
						<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
							<div class="bg-blue-50 p-4 rounded-lg">
								<div class="text-2xl font-bold text-blue-600">{formResponses.length}</div>
								<div class="text-sm text-blue-600">Total Responses</div>
							</div>
							<div class="bg-green-50 p-4 rounded-lg">
								<div class="text-2xl font-bold text-green-600">
									{Math.round((formResponses.length / (formResponses.length + 5)) * 100)}%
								</div>
								<div class="text-sm text-green-600">Completion Rate</div>
							</div>
							<div class="bg-purple-50 p-4 rounded-lg">
								<div class="text-2xl font-bold text-purple-600">
									{new Date(Math.max(...formResponses.map(r => new Date(r.submitted_at)))).toLocaleDateString()}
								</div>
								<div class="text-sm text-purple-600">Latest Response</div>
							</div>
						</div>

						<!-- Responses Table -->
						<div class="bg-white border border-gray-200 rounded-lg overflow-hidden">
							<div class="overflow-x-auto">
								<table class="w-full">
									<thead class="bg-gray-50">
										<tr>
											<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
											<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Submitted</th>
											{#each formData.fields as field}
												<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{field.label}</th>
											{/each}
										</tr>
									</thead>
									<tbody class="bg-white divide-y divide-gray-200">
										{#each formResponses as response}
											<tr class="hover:bg-gray-50">
												<td class="px-4 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
													{response.id}
												</td>
												<td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500">
													{new Date(response.submitted_at).toLocaleString()}
												</td>
												{#each formData.fields as field}
													<td class="px-4 py-4 whitespace-nowrap text-sm text-gray-900">
														{response.responses[field.label.toLowerCase()] || '-'}
													</td>
												{/each}
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}

<!-- Conditional Logic Modal -->
{#if showConditionalLogic}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-2xl max-h-[90vh] overflow-hidden flex flex-col">
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<div>
					<h2 class="text-2xl font-bold text-gray-900">Conditional Logic</h2>
					<p class="text-gray-600 mt-1">Show/hide fields based on responses to other fields</p>
				</div>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showConditionalLogic = false}
				>
					×
				</button>
			</div>

			<div class="flex-1 overflow-y-auto p-6">
				{#if selectedField}
					<div class="space-y-6">
						<div class="bg-blue-50 p-4 rounded-lg">
							<h3 class="font-medium text-blue-900 mb-2">Field: {selectedField.label}</h3>
							<p class="text-sm text-blue-700">
								Set up rules to show or hide this field based on responses to other questions.
							</p>
						</div>

						{#if selectedField.conditionalLogic && selectedField.conditionalLogic.length > 0}
							<div class="space-y-4">
								{#each selectedField.conditionalLogic as rule, index}
									<div class="border border-gray-200 rounded-lg p-4">
										<div class="flex items-center justify-between mb-3">
											<h4 class="font-medium text-gray-900">Rule {index + 1}</h4>
											<button
												class="text-red-500 hover:text-red-700"
												on:click={() => removeConditionalRule(index)}
											>
												Remove
											</button>
										</div>

										<div class="grid grid-cols-1 md:grid-cols-4 gap-3">
											<div>
												<label class="block text-sm font-medium text-gray-700 mb-1">If field</label>
												<select
													bind:value={rule.field}
													on:change={(e) => updateConditionalRule(index, 'field', e.target.value)}
													class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
												>
													<option value="">Select field...</option>
													{#each formData.fields.filter(f => f.id !== selectedField.id) as field}
														<option value={field.id}>{field.label}</option>
													{/each}
												</select>
											</div>

											<div>
												<label class="block text-sm font-medium text-gray-700 mb-1">Operator</label>
												<select
													bind:value={rule.operator}
													on:change={(e) => updateConditionalRule(index, 'operator', e.target.value)}
													class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
												>
													<option value="equals">Equals</option>
													<option value="not_equals">Not equals</option>
													<option value="contains">Contains</option>
													<option value="not_contains">Does not contain</option>
													<option value="greater_than">Greater than</option>
													<option value="less_than">Less than</option>
												</select>
											</div>

											<div>
												<label class="block text-sm font-medium text-gray-700 mb-1">Value</label>
												<input
													type="text"
													bind:value={rule.value}
													on:input={(e) => updateConditionalRule(index, 'value', e.target.value)}
													placeholder="Enter value..."
													class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
												/>
											</div>

											<div>
												<label class="block text-sm font-medium text-gray-700 mb-1">Then</label>
												<select
													bind:value={rule.action}
													on:change={(e) => updateConditionalRule(index, 'action', e.target.value)}
													class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
												>
													<option value="show">Show field</option>
													<option value="hide">Hide field</option>
												</select>
											</div>
										</div>
									</div>
								{/each}
							</div>
						{:else}
							<div class="text-center py-8 text-gray-500">
								<p>No conditional rules set up yet.</p>
								<p class="text-sm mt-1">Add a rule below to show/hide this field based on other responses.</p>
							</div>
						{/if}

						<div class="flex justify-center">
							<button
								class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
								on:click={addConditionalLogic}
							>
								➕ Add Rule
							</button>
						</div>
					</div>
				{/if}
			</div>

			<div class="border-t border-gray-200 p-6 bg-gray-50">
				<div class="flex justify-end space-x-3">
					<button
						class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50"
						on:click={() => showConditionalLogic = false}
					>
						Cancel
					</button>
					<button
						class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
						on:click={() => showConditionalLogic = false}
					>
						Save Logic
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Enhanced Field Properties Panel -->
{#if selectedField}
	<!-- Add buttons to the field properties panel -->
	<div class="border-t border-gray-200 p-4">
		<div class="flex space-x-2">
			<button
				class="flex-1 px-3 py-2 bg-purple-600 text-white text-sm rounded hover:bg-purple-700 transition-colors"
				on:click={() => showConditionalLogic = true}
			>
				🎯 Add Logic
			</button>
			<button
				class="flex-1 px-3 py-2 bg-green-600 text-white text-sm rounded hover:bg-green-700"
				on:click={loadFormResponses}
			>
				📊 View Responses
			</button>
		</div>
	</div>
{/if}

<!-- Add template button to main form actions -->
<div class="flex justify-between items-center pt-6 border-t border-gray-200">
	<div class="flex space-x-3">
		<button
			class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
			on:click={() => showFormTemplates = true}
		>
			📋 Use Template
		</button>
		<div class="text-sm text-gray-600 flex items-center">
			{formData.fields.length} field{formData.fields.length !== 1 ? 's' : ''}
		</div>
	</div>
	<div class="flex space-x-3">
		<button
			on:click={cancel}
			class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
		>
			Cancel
		</button>
		<button
			on:click={saveForm}
			class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
		>
			Save Form
		</button>
	</div>
</div>

<style>
	/* Custom drag styles */
	.cursor-move {
		cursor: move;
	}
</style>
