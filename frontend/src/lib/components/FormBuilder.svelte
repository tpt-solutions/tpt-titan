<script>
	import { createEventDispatcher } from 'svelte';

	export let form = null; // Existing form to edit, or null for new

	const dispatch = createEventDispatcher();

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

<style>
	/* Custom drag styles */
	.cursor-move {
		cursor: move;
	}
</style>
