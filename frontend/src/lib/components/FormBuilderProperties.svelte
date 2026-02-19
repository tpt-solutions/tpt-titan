<!-- frontend/src/lib/components/FormBuilderProperties.svelte -->
<!-- Left-sidebar panel: form settings at top, selected-field properties below. -->
<script>
	import { createEventDispatcher } from 'svelte';

	export let formData;
	export let selectedField = null;

	const dispatch = createEventDispatcher();

	function updateFieldProperty(fieldId, property, value) {
		dispatch('updateProperty', { fieldId, property, value });
	}

	function updateFieldLabel(fieldId, label) {
		dispatch('updateLabel', { fieldId, label });
	}
</script>

<div class="w-80 bg-white border-r border-gray-200 flex flex-col">
	<!-- Form name / description -->
	<div class="p-4 border-b border-gray-200">
		<h3 class="text-lg font-semibold text-gray-900 mb-3">Form Settings</h3>
		<div class="space-y-3">
			<div>
				<label for="form-name" class="block text-sm font-medium text-gray-700 mb-1">Form Name</label>
				<input
					id="form-name"
					bind:value={formData.name}
					placeholder="Enter form name"
					class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
				/>
			</div>
			<div>
				<label for="form-description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
				<textarea
					id="form-description"
					bind:value={formData.description}
					placeholder="Describe your form"
					rows="3"
					class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
				></textarea>
			</div>
				<button
					class="w-full px-3 py-2 bg-indigo-600 text-white text-sm rounded-lg hover:bg-indigo-700 transition-colors"
					on:click={() => dispatch('showTemplates')}
				>
					📋 Use Template
				</button>
			</div>

			<!-- Form Settings -->
			<div class="mt-4 pt-4 border-t border-gray-200">
				<h4 class="text-sm font-medium text-gray-900 mb-3">Form Options</h4>
				
				<!-- Multi-page toggle -->
				<div class="flex items-center mb-3">
					<input
						type="checkbox"
						id="multi-page"
						bind:checked={formData.settings.isMultiPage}
						on:change={(e) => dispatch('updateSetting', { property: 'isMultiPage', value: e.target.checked })}
						class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
					/>
					<label for="multi-page" class="ml-2 text-sm text-gray-700">Multi-page form</label>
				</div>

				<!-- CAPTCHA toggle -->
				<div class="flex items-center mb-3">
					<input
						type="checkbox"
						id="captcha"
						bind:checked={formData.settings.captcha}
						on:change={(e) => dispatch('updateSetting', { property: 'captcha', value: e.target.checked })}
						class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
					/>
					<label for="captcha" class="ml-2 text-sm text-gray-700">Enable CAPTCHA protection</label>
				</div>

				<!-- Allow multiple submissions -->
				<div class="flex items-center mb-3">
					<input
						type="checkbox"
						id="multiple-submissions"
						bind:checked={formData.settings.allowMultipleSubmissions}
						on:change={(e) => dispatch('updateSetting', { property: 'allowMultipleSubmissions', value: e.target.checked })}
						class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
					/>
					<label for="multiple-submissions" class="ml-2 text-sm text-gray-700">Allow multiple submissions</label>
				</div>

				<!-- Max submissions -->
				{#if !formData.settings.allowMultipleSubmissions}
					<div class="mb-3">
						<label for="max-submissions" class="block text-sm font-medium text-gray-700 mb-1">Max Submissions</label>
						<input
							id="max-submissions"
							type="number"
							bind:value={formData.settings.maxSubmissions}
							on:input={(e) => dispatch('updateSetting', { property: 'maxSubmissions', value: parseInt(e.target.value) || null })}
							placeholder="Unlimited"
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
				{/if}
			</div>

			<!-- Theme Settings -->
			<div class="mt-4 pt-4 border-t border-gray-200">
				<h4 class="text-sm font-medium text-gray-900 mb-3">Theme & Styling</h4>
				
				<div class="space-y-3">
					<div>
						<label for="theme-color" class="block text-sm font-medium text-gray-700 mb-1">Primary Color</label>
						<div class="flex items-center space-x-2">
							<input
								type="color"
								id="theme-color"
								bind:value={formData.settings.themeColor}
								on:input={(e) => dispatch('updateSetting', { property: 'themeColor', value: e.target.value })}
								class="h-8 w-16 border border-gray-300 rounded cursor-pointer"
							/>
							<span class="text-sm text-gray-600">{formData.settings.themeColor || '#4F46E5'}</span>
						</div>
					</div>

					<div>
						<label for="theme-style" class="block text-sm font-medium text-gray-700 mb-1">Form Style</label>
						<select
							id="theme-style"
							bind:value={formData.settings.themeStyle}
							on:change={(e) => dispatch('updateSetting', { property: 'themeStyle', value: e.target.value })}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						>
							<option value="default">Default</option>
							<option value="minimal">Minimal</option>
							<option value="card">Card Style</option>
							<option value="bordered">Bordered</option>
							<option value="rounded">Rounded</option>
						</select>
					</div>

					<div class="flex items-center">
						<input
							type="checkbox"
							id="show-progress"
							bind:checked={formData.settings.showProgressBar}
							on:change={(e) => dispatch('updateSetting', { property: 'showProgressBar', value: e.target.checked })}
							class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
						/>
						<label for="show-progress" class="ml-2 text-sm text-gray-700">Show progress bar</label>
					</div>
				</div>
			</div>
		</div>


	<!-- Field properties -->
	{#if selectedField}
		<div class="flex-1 p-4 overflow-y-auto">
			<h3 class="text-lg font-semibold text-gray-900 mb-3">Field Properties</h3>
			<div class="space-y-4">
				<!-- Label -->
				<div>
					<label for="field-label-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">Label</label>
					<input
						id="field-label-{selectedField.id}"
						bind:value={selectedField.label}
						on:input={(e) => updateFieldLabel(selectedField.id, e.target.value)}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>

				<!-- Placeholder (text / textarea / email) -->
				{#if ['text', 'textarea', 'email'].includes(selectedField.type)}
					<div>
						<label for="field-placeholder-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">Placeholder</label>
						<input
							id="field-placeholder-{selectedField.id}"
							bind:value={selectedField.properties.placeholder}
							on:input={(e) => updateFieldProperty(selectedField.id, 'placeholder', e.target.value)}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
				{/if}

				<!-- Min / Max (number) -->
				{#if selectedField.type === 'number'}
					<div class="grid grid-cols-2 gap-2">
						<div>
							<label for="field-min-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">Min</label>
							<input
								id="field-min-{selectedField.id}"
								type="number"
								bind:value={selectedField.properties.min}
								on:input={(e) => updateFieldProperty(selectedField.id, 'min', e.target.value ? parseFloat(e.target.value) : null)}
								class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
						<div>
							<label for="field-max-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">Max</label>
							<input
								id="field-max-{selectedField.id}"
								type="number"
								bind:value={selectedField.properties.max}
								on:input={(e) => updateFieldProperty(selectedField.id, 'max', e.target.value ? parseFloat(e.target.value) : null)}
								class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
					</div>
				{/if}

				<!-- Options (select / radio / checkbox) -->
				{#if ['select', 'radio', 'checkbox'].includes(selectedField.type)}
					<div>
						<label for="field-options-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">Options</label>
						<textarea
							id="field-options-{selectedField.id}"
							bind:value={selectedField.properties.options}
							on:input={(e) => updateFieldProperty(selectedField.id, 'options', e.target.value.split('\n').filter(opt => opt.trim()))}
							placeholder="One option per line"
							rows="4"
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>
				{/if}

				<!-- File upload settings -->
				{#if selectedField.type === 'file'}
					<div class="grid grid-cols-2 gap-2">
						<div>
							<label for="field-maxsize-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">Max Size (MB)</label>
							<input
								id="field-maxsize-{selectedField.id}"
								type="number"
								bind:value={selectedField.properties.maxSize}
								on:input={(e) => updateFieldProperty(selectedField.id, 'maxSize', parseInt(e.target.value))}
								class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
						<div>
							<label for="field-accept-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">File Types</label>
							<input
								id="field-accept-{selectedField.id}"
								bind:value={selectedField.properties.accept}
								on:input={(e) => updateFieldProperty(selectedField.id, 'accept', e.target.value)}
								placeholder="e.g., .pdf,.doc"
								class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
					</div>
				{/if}

				<!-- Required toggle -->
				<div class="flex items-center">
					<input
						type="checkbox"
						id="required-{selectedField.id}"
						bind:checked={selectedField.properties.required}
						on:change={(e) => updateFieldProperty(selectedField.id, 'required', e.target.checked)}
						class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
					/>
					<label for="required-{selectedField.id}" class="ml-2 text-sm text-gray-700">Required field</label>
				</div>

				<!-- Calculated Field Formula -->
				{#if selectedField.type === 'calculation'}
					<div class="mt-4 pt-4 border-t border-gray-200">
						<label for="field-formula-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">
							Formula
							<span class="text-xs text-gray-500 ml-1">(Use field labels in braces, e.g., {`{Price}`} * {`{Quantity}`})</span>
						</label>
						<textarea
							id="field-formula-{selectedField.id}"
							bind:value={selectedField.properties.formula}
							on:input={(e) => updateFieldProperty(selectedField.id, 'formula', e.target.value)}
							placeholder="e.g., {`{Field1}`} + {`{Field2}`}"
							rows="3"
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
						></textarea>
						
						<div class="mt-2">
							<label for="decimal-places-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">Decimal Places</label>
							<input
								id="decimal-places-{selectedField.id}"
								type="number"
								min="0"
								max="10"
								bind:value={selectedField.properties.decimalPlaces}
								on:input={(e) => updateFieldProperty(selectedField.id, 'decimalPlaces', parseInt(e.target.value) || 2)}
								class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>

						<div class="mt-2 p-2 bg-blue-50 rounded text-xs text-blue-700">
							<strong>Available operators:</strong> +, -, *, /, %, ^, (), SUM(), AVG(), MIN(), MAX()
						</div>
					</div>
				{/if}

				<!-- Page Assignment (for multi-page forms) -->
				{#if formData.settings?.isMultiPage}
					<div class="mt-4 pt-4 border-t border-gray-200">
						<label for="field-page-{selectedField.id}" class="block text-sm font-medium text-gray-700 mb-1">Page Number</label>
						<input
							id="field-page-{selectedField.id}"
							type="number"
							min="1"
							bind:value={selectedField.properties.page}
							on:input={(e) => updateFieldProperty(selectedField.id, 'page', parseInt(e.target.value) || 1)}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
				{/if}


				<!-- Action buttons -->
				<div class="mt-4 pt-4 border-t border-gray-200 flex space-x-2">
					<button
						class="flex-1 px-3 py-2 bg-purple-600 text-white text-sm rounded hover:bg-purple-700 transition-colors"
						on:click={() => dispatch('showConditional')}
					>
						🎯 Add Logic
					</button>
					<button
						class="flex-1 px-3 py-2 bg-green-600 text-white text-sm rounded hover:bg-green-700 transition-colors"
						on:click={() => dispatch('showResponses')}
					>
						📊 Responses
					</button>
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
