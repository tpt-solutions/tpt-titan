<!-- frontend/src/lib/components/FormPreview.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';
	import { validateField, validateForm, getDefaultValue, generateSampleValue } from '../utils/form-validation.js';

	export let formData = { fields: [], settings: {} };
	export let testMode = false;
	export let showValidation = true;

	const dispatch = createEventDispatcher();

	let formValues = {};
	let validationResults = {};
	let showErrors = false;
	let isSubmitting = false;
	let submitSuccess = false;

	// Initialize form values
	$: if (formData.fields) {
		formData.fields.forEach(field => {
			if (!(field.id in formValues)) {
				formValues[field.id] = getDefaultValue(field);
			}
		});
	}

	function handleInput(fieldId, value) {
		formValues[fieldId] = value;
		
		if (showValidation || showErrors) {
			const field = formData.fields.find(f => f.id === fieldId);
			if (field) {
				const result = validateField(value, field);
				validationResults[fieldId] = result;
			}
		}
		
		dispatch('change', { fieldId, value });
	}

	function handleSubmit() {
		showErrors = true;
		const validation = validateForm(formValues, formData.fields);
		validationResults = validation.results;
		
		if (validation.valid) {
			isSubmitting = true;
			dispatch('submit', { values: formValues });
			
			// Simulate submission
			setTimeout(() => {
				isSubmitting = false;
				submitSuccess = true;
				setTimeout(() => submitSuccess = false, 3000);
			}, 500);
		} else {
			dispatch('validationError', { results: validationResults });
		}
	}

	function fillSampleData() {
		formData.fields.forEach(field => {
			formValues[field.id] = generateSampleValue(field);
		});
		formValues = { ...formValues };
		showErrors = true;
		validateAllFields();
	}

	function validateAllFields() {
		const validation = validateForm(formValues, formData.fields);
		validationResults = validation.results;
	}

	function clearForm() {
		formData.fields.forEach(field => {
			formValues[field.id] = getDefaultValue(field);
		});
		formValues = { ...formValues };
		showErrors = false;
		validationResults = {};
	}

	function getFieldError(fieldId) {
		return showErrors && validationResults[fieldId]?.error;
	}

	function isFieldValid(fieldId) {
		return showErrors && validationResults[fieldId]?.valid;
	}

	function renderField(field) {
		const value = formValues[field.id];
		const error = getFieldError(field.id);
		const isValid = isFieldValid(field.id);
		const props = field.properties || {};
		
		const fieldClass = `w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 ${
			error ? 'border-red-500 focus:ring-red-200' : 
			isValid ? 'border-green-500 focus:ring-green-200' : 
			'border-gray-300 focus:ring-blue-200'
		}`;
		
		const labelClass = `block text-sm font-medium mb-1 ${
			props.required ? 'text-gray-900' : 'text-gray-700'
		}`;
		
		return { value, error, isValid, props, fieldClass, labelClass };
	}
</script>

<div class="bg-white rounded-lg shadow p-6">
	<!-- Header -->
	<div class="mb-6">
		<h2 class="text-2xl font-bold text-gray-900">{formData.name || 'Untitled Form'}</h2>
		{#if formData.description}
			<p class="text-gray-600 mt-2">{formData.description}</p>
		{/if}
	</div>

	<!-- Test Mode Controls -->
	{#if testMode}
		<div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
			<div class="flex items-center justify-between">
				<div>
					<h3 class="text-sm font-semibold text-blue-900">🧪 Test Mode</h3>
					<p class="text-sm text-blue-700 mt-1">
						Validation is {showValidation ? 'enabled' : 'disabled'}. 
						Errors will show after attempting to submit.
					</p>
				</div>
				<div class="flex space-x-2">
					<button
						class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
						on:click={fillSampleData}
					>
						Fill Sample Data
					</button>
					<button
						class="px-3 py-1 text-sm bg-gray-600 text-white rounded hover:bg-gray-700"
						on:click={clearForm}
					>
						Clear
					</button>
				</div>
			</div>
		</div>
	{/if}

	<!-- Form Fields -->
	<form on:submit|preventDefault={handleSubmit} class="space-y-4">
		{#each formData.fields as field (field.id)}
			{@const { value, error, isValid, props, fieldClass, labelClass } = renderField(field)}
			
			<div class="form-field">
				<label class={labelClass}>
					{field.label}
					{#if props.required}
						<span class="text-red-500 ml-1">*</span>
					{/if}
				</label>
				
				<!-- Text Input -->
				{#if field.type === 'text' || field.type === 'email' || field.type === 'url' || field.type === 'phone' || field.type === 'password'}
					<input
						type={field.type === 'password' ? 'password' : field.type === 'email' ? 'email' : field.type === 'url' ? 'url' : 'text'}
						class={fieldClass}
						placeholder={props.placeholder}
						value={value || ''}
						on:input={(e) => handleInput(field.id, e.target.value)}
					/>
				{/if}
				
				<!-- Textarea -->
				{#if field.type === 'textarea'}
					<textarea
						class={fieldClass}
						placeholder={props.placeholder}
						rows={props.rows || 3}
						value={value || ''}
						on:input={(e) => handleInput(field.id, e.target.value)}
					></textarea>
				{/if}
				
				<!-- Number, Currency, Percentage -->
				{#if field.type === 'number' || field.type === 'currency' || field.type === 'percentage'}
					<div class="relative">
						{#if field.type === 'currency'}
							<span class="absolute left-3 top-2 text-gray-500">
								{props.currency === 'USD' ? '$' : props.currency === 'EUR' ? '€' : props.currency === 'GBP' ? '£' : '$'}
							</span>
						{/if}
						<input
							type="number"
							class="{fieldClass} {field.type === 'currency' ? 'pl-8' : ''}"
							placeholder={props.placeholder}
							min={props.min}
							max={props.max}
							step={field.type === 'currency' || field.type === 'percentage' ? '0.01' : '1'}
							value={value || ''}
							on:input={(e) => handleInput(field.id, parseFloat(e.target.value))}
						/>
						{#if field.type === 'percentage'}
							<span class="absolute right-3 top-2 text-gray-500">%</span>
						{/if}
					</div>
				{/if}
				
				<!-- Select -->
				{#if field.type === 'select'}
					<select
						class={fieldClass}
						value={value || ''}
						on:change={(e) => handleInput(field.id, e.target.value)}
					>
						<option value="">Select an option...</option>
						{#each props.options || [] as option}
							<option value={option}>{option}</option>
						{/each}
					</select>
				{/if}
				
				<!-- Radio -->
				{#if field.type === 'radio'}
					<div class="space-y-2">
						{#each props.options || [] as option}
							<label class="flex items-center space-x-2 cursor-pointer">
								<input
									type="radio"
									name={field.id}
									value={option}
									checked={value === option}
									on:change={() => handleInput(field.id, option)}
									class="text-blue-600 focus:ring-blue-500"
								/>
								<span class="text-gray-700">{option}</span>
							</label>
						{/each}
					</div>
				{/if}
				
				<!-- Checkbox -->
				{#if field.type === 'checkbox'}
					<div class="space-y-2">
						{#each props.options || [] as option}
							<label class="flex items-center space-x-2 cursor-pointer">
								<input
									type="checkbox"
									value={option}
									checked={Array.isArray(value) && value.includes(option)}
									on:change={(e) => {
										const currentValues = Array.isArray(value) ? [...value] : [];
										if (e.target.checked) {
											handleInput(field.id, [...currentValues, option]);
										} else {
											handleInput(field.id, currentValues.filter(v => v !== option));
										}
									}}
									class="text-blue-600 focus:ring-blue-500 rounded"
								/>
								<span class="text-gray-700">{option}</span>
							</label>
						{/each}
					</div>
				{/if}
				
				<!-- Yes/No -->
				{#if field.type === 'yesno'}
					<div class="flex space-x-4">
						<label class="flex items-center space-x-2 cursor-pointer">
							<input
								type="radio"
								name={field.id}
								value={true}
								checked={value === true}
								on:change={() => handleInput(field.id, true)}
								class="text-blue-600 focus:ring-blue-500"
							/>
							<span class="text-gray-700">Yes 👍</span>
						</label>
						<label class="flex items-center space-x-2 cursor-pointer">
							<input
								type="radio"
								name={field.id}
								value={false}
								checked={value === false}
								on:change={() => handleInput(field.id, false)}
								class="text-blue-600 focus:ring-blue-500"
							/>
							<span class="text-gray-700">No 👎</span>
						</label>
					</div>
				{/if}
				
				<!-- Rating -->
				{#if field.type === 'rating'}
					<div class="flex space-x-1">
						{#each Array(props.maxRating || 5) as _, i}
							<button
								type="button"
								class="text-2xl focus:outline-none transition-colors"
								class:text-yellow-400={i < (value || 0)}
								class:text-gray-300={i >= (value || 0)}
								on:click={() => handleInput(field.id, i + 1)}
							>
								⭐
							</button>
						{/each}
					</div>
				{/if}
				
				<!-- Scale -->
				{#if field.type === 'scale'}
					<div class="space-y-2">
						<input
							type="range"
							min={props.min}
							max={props.max}
							value={value || props.min}
							class="w-full"
							on:input={(e) => handleInput(field.id, parseInt(e.target.value))}
						/>
						<div class="flex justify-between text-sm text-gray-500">
							<span>{props.minLabel || props.min}</span>
							<span class="font-semibold text-blue-600">{value || props.min}</span>
							<span>{props.maxLabel || props.max}</span>
						</div>
					</div>
				{/if}
				
				<!-- Date, Time, DateTime -->
				{#if field.type === 'date' || field.type === 'time' || field.type === 'datetime-local'}
					<input
						type={field.type}
						class={fieldClass}
						value={value || ''}
						on:input={(e) => handleInput(field.id, e.target.value)}
					/>
				{/if}
				
				<!-- Color -->
				{#if field.type === 'color'}
					<div class="flex items-center space-x-2">
						<input
							type="color"
							class="h-10 w-20 rounded border border-gray-300"
							value={value || '#000000'}
							on:input={(e) => handleInput(field.id, e.target.value)}
						/>
						<span class="text-gray-600 text-sm">{value || '#000000'}</span>
					</div>
				{/if}
				
				<!-- Range -->
				{#if field.type === 'range'}
					<div class="space-y-2">
						<input
							type="range"
							min={props.min}
							max={props.max}
							step={props.step || 1}
							value={value || props.min}
							class="w-full"
							on:input={(e) => handleInput(field.id, parseFloat(e.target.value))}
						/>
						<div class="text-center text-sm text-gray-600">
							Value: {value || props.min}
						</div>
					</div>
				{/if}
				
				<!-- File Upload -->
				{#if field.type === 'file' || field.type === 'image-upload' || field.type === 'video-upload' || field.type === 'audio-upload'}
					<div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 transition-colors">
						<input
							type="file"
							accept={props.accept}
							class="hidden"
							id="file-{field.id}"
							on:change={(e) => handleInput(field.id, e.target.files[0])}
						/>
						<label for="file-{field.id}" class="cursor-pointer">
							<div class="text-gray-500">
								{#if value && value.name}
									<div class="text-green-600 font-medium">
										📎 {value.name}
										<span class="text-sm text-gray-500">({(value.size / 1024).toFixed(1)} KB)</span>
									</div>
								{:else}
									<div class="text-4xl mb-2">📤</div>
									<div>Click to upload {field.type === 'image-upload' ? 'image' : field.type === 'video-upload' ? 'video' : field.type === 'audio-upload' ? 'audio' : 'file'}</div>
									<div class="text-sm text-gray-400 mt-1">
										Max size: {props.maxSize || 10}MB
										{#if props.accept && props.accept !== '*'}
											<br>Accepted: {props.accept}
										{/if}
									</div>
								{/if}
							</div>
						</label>
					</div>
				{/if}
				
				<!-- Signature -->
				{#if field.type === 'signature'}
					<div class="border-2 border-dashed border-gray-300 rounded-lg p-4">
						{#if value}
							<div class="bg-white p-4 rounded border">
								<img src={value} alt="Signature" class="max-h-24 mx-auto" />
								<button
									type="button"
									class="mt-2 text-sm text-red-600 hover:text-red-800"
									on:click={() => handleInput(field.id, null)}
								>
									Clear Signature
								</button>
							</div>
						{:else}
							<div class="text-center text-gray-500 py-8">
								<div class="text-4xl mb-2">✍️</div>
								<div>Signature pad would appear here</div>
								<div class="text-sm text-gray-400 mt-1">Click to sign</div>
							</div>
						{/if}
					</div>
				{/if}
				
				<!-- Calculated Field -->
				{#if field.type === 'calculation'}
					<div class="bg-gray-50 p-3 rounded border">
						<div class="text-sm text-gray-600">Formula: {props.formula}</div>
						<div class="text-lg font-semibold text-blue-600 mt-1">
							{value || '—'}
						</div>
					</div>
				{/if}
				
				<!-- HTML Content -->
				{#if field.type === 'html'}
					<div class="prose prose-sm max-w-none">
						{@html props.content || '<p>Custom HTML content</p>'}
					</div>
				{/if}
				
				<!-- Divider -->
				{#if field.type === 'divider'}
					<div class="border-t-2 border-gray-200 pt-4 mt-4">
						{#if props.title}
							<h3 class="text-lg font-semibold text-gray-800">{props.title}</h3>
						{/if}
					</div>
				{/if}
				
				<!-- Error Message -->
				{#if error}
					<p class="text-red-500 text-sm mt-1 flex items-center">
						<span class="mr-1">⚠️</span> {error}
					</p>
				{/if}
				
				<!-- Valid Indicator -->
				{#if isValid && testMode}
					<p class="text-green-500 text-sm mt-1 flex items-center">
						<span class="mr-1">✓</span> Valid
					</p>
				{/if}
			</div>
		{/each}
		
		<!-- Submit Button -->
		<div class="pt-4 border-t border-gray-200">
			{#if submitSuccess}
				<div class="bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded mb-4">
					✅ Form submitted successfully! (Test mode - no data was saved)
				</div>
			{/if}
			
			<div class="flex items-center justify-between">
				<button
					type="submit"
					class="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
					disabled={isSubmitting}
				>
					{isSubmitting ? 'Submitting...' : 'Submit'}
				</button>
				
				{#if testMode}
					<span class="text-sm text-gray-500">
						{Object.values(validationResults).filter(r => r.valid).length} / {formData.fields.length} fields valid
					</span>
				{/if}
			</div>
		</div>
	</form>
</div>

<style>
	.form-field {
		@apply mb-4;
	}
	
	/* Custom range slider styling */
	input[type="range"] {
		@apply h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer;
	}
	
	input[type="range"]::-webkit-slider-thumb {
		@apply appearance-none w-4 h-4 bg-blue-600 rounded-full cursor-pointer;
	}
	
	input[type="range"]::-moz-range-thumb {
		@apply w-4 h-4 bg-blue-600 rounded-full cursor-pointer border-0;
	}
</style>
