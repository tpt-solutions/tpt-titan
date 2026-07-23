<!-- frontend/src/lib/components/FormBuilderCanvas.svelte -->
<!-- Centre canvas: live form preview with drag-to-reorder support. -->
<script>
// @ts-nocheck
	import { createEventDispatcher } from 'svelte';

	export let formData;
	export let selectedField = null;
	export let isDragging = false;
	export let draggedField = null;

	const dispatch = createEventDispatcher();

	function selectField(field) {
		dispatch('selectField', field);
	}

	function deleteField(fieldId) {
		dispatch('deleteField', fieldId);
	}

	function handleFieldDragStart(event, field) {
		dispatch('fieldDragStart', { event, field });
		event.target.style.opacity = '0.5';
	}

	function handleFieldDragEnd(event) {
		document.querySelectorAll('.form-field').forEach(el => {
			el.style.opacity = '1';
			el.style.borderColor = '';
		});
		dispatch('fieldDragEnd');
	}

	function handleFieldDragOver(event, dropIndex) {
		event.preventDefault();
		event.dataTransfer.dropEffect = 'move';
		const target = event.currentTarget;
		if (draggedField && draggedField.id !== formData.fields[dropIndex]?.id) {
			target.style.borderColor = '#3b82f6';
			target.style.backgroundColor = '#eff6ff';
		}
	}

	function handleFieldDrop(event, dropIndex) {
		event.preventDefault();
		dispatch('fieldDrop', { event, dropIndex });
		document.querySelectorAll('.form-field').forEach(el => {
			el.style.borderColor = '';
			el.style.backgroundColor = '';
			el.style.opacity = '1';
		});
	}

	function handleCanvasDrop(event) {
		event.preventDefault();
		dispatch('canvasDrop', { event, dropIndex: formData.fields.length });
	}
</script>

<div class="flex-1 bg-gray-50 overflow-y-auto">
	<div class="max-w-2xl mx-auto py-8">
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-8">
			<!-- Form header -->
			<div class="mb-8">
				<h2 class="text-2xl font-bold text-gray-900 mb-2">
					{formData.name || 'Untitled Form'}
				</h2>
				{#if formData.description}
					<p class="text-gray-600">{formData.description}</p>
				{/if}
			</div>

			<!-- Fields -->
			{#each formData.fields as field, index}
				<div
					class="form-field mb-6 p-4 border-2 border-dashed border-gray-300 rounded-lg cursor-move hover:border-blue-400 transition-colors {selectedField?.id === field.id ? 'border-blue-500 bg-blue-50' : ''} {isDragging && draggedField?.id === field.id ? 'opacity-50' : ''}"
					draggable="true"
					on:click={() => selectField(field)}
					on:dragstart={(e) => handleFieldDragStart(e, field)}
					on:dragend={handleFieldDragEnd}
					on:drop={(e) => handleFieldDrop(e, index)}
					on:dragover={(e) => handleFieldDragOver(e, index)}
					role="button"
					tabindex="0"
					aria-label="Form field: {field.label}"
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

					<!-- Field preview -->
					<div class="space-y-2">
						{#if field.type === 'text'}
							<input type="text" placeholder={field.properties.placeholder || 'Enter text...'} disabled class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50" />
						{:else if field.type === 'textarea'}
							<textarea placeholder={field.properties.placeholder || 'Enter text...'} rows={field.properties.rows || 3} disabled class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50"></textarea>
						{:else if field.type === 'number'}
							<input type="number" placeholder={field.properties.placeholder || 'Enter number...'} min={field.properties.min} max={field.properties.max} disabled class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50" />
						{:else if field.type === 'email'}
							<input type="email" placeholder={field.properties.placeholder || 'Enter email...'} disabled class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50" />
						{:else if field.type === 'date'}
							<input type="date" disabled class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50" />
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
									<label class="flex items-center"><input type="radio" disabled class="mr-2" /><span>{option}</span></label>
								{/each}
							</div>
						{:else if field.type === 'checkbox'}
							<div class="space-y-2">
								{#each field.properties.options || [] as option}
									<label class="flex items-center"><input type="checkbox" disabled class="mr-2" /><span>{option}</span></label>
								{/each}
							</div>
						{:else if field.type === 'signature'}
							<div class="border-2 border-dashed border-gray-300 rounded p-8 text-center text-gray-500">Signature Area</div>
						{:else}
							<div class="text-gray-500 italic">[{field.type} field preview]</div>
						{/if}

						{#if field.properties.required}
							<span class="text-red-500 text-sm">* Required</span>
						{/if}
					</div>
				</div>
			{/each}

			<!-- Empty drop zone -->
			{#if formData.fields.length === 0}
				<div
					class="mb-6 p-8 border-2 border-dashed border-gray-300 rounded-lg text-center text-gray-500 hover:border-blue-400 transition-colors"
					on:drop={handleCanvasDrop}
					on:dragover={(e) => e.preventDefault()}
					role="region"
					aria-label="Drop zone for form fields"
				>
					<svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
					</svg>
					<p>Drag field types here to start building your form</p>
				</div>
			{/if}

			<!-- Form actions bar -->
			<div class="flex justify-between items-center pt-6 border-t border-gray-200">
				<div class="text-sm text-gray-600">
					{formData.fields.length} field{formData.fields.length !== 1 ? 's' : ''}
				</div>
				<div class="flex space-x-3">
					<button
						on:click={() => dispatch('cancel')}
						class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
					>
						Cancel
					</button>
					<button
						on:click={() => dispatch('save')}
						class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					>
						Save Form
					</button>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.cursor-move { cursor: move; }
</style>
