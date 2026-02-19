<!-- frontend/src/lib/components/FormBuilderFieldPanel.svelte -->
<!-- Right sidebar: draggable field-type catalogue + quick-add buttons. -->
<script>
	import { createEventDispatcher } from 'svelte';
	import { fieldTypes } from './FormBuilderFieldTypes.js';

	export let isDragging = false;
	export let draggedFieldType = null;

	const dispatch = createEventDispatcher();

	function handleDragStart(event, fieldType) {
		dispatch('dragStart', { event, fieldType });
	}

	function handleDragEnd() {
		dispatch('dragEnd');
	}

	function quickAdd(fieldType) {
		dispatch('quickAdd', fieldType);
	}
</script>

<div class="w-64 bg-white border-l border-gray-200 p-4">
	<h3 class="text-lg font-semibold text-gray-900 mb-4">Field Types</h3>

	<div class="space-y-2">
		{#each fieldTypes as fieldType}
			<div
				class="p-3 border border-gray-200 rounded-lg cursor-move hover:bg-gray-50 transition-colors {isDragging && draggedFieldType?.id === fieldType.id ? 'opacity-50' : ''}"
				draggable="true"
				on:dragstart={(e) => handleDragStart(e, fieldType)}
				on:dragend={handleDragEnd}
				role="button"
				tabindex="0"
				aria-label="Drag {fieldType.name} field to add to form"
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

	<!-- Quick add -->
	<div class="mt-6 pt-4 border-t border-gray-200">
		<h4 class="text-sm font-medium text-gray-900 mb-2">Quick Add</h4>
		<div class="space-y-1">
			{#each fieldTypes.slice(0, 3) as fieldType}
				<button
					on:click={() => quickAdd(fieldType)}
					class="w-full text-left p-2 text-sm hover:bg-gray-50 rounded transition-colors"
				>
					+ {fieldType.name}
				</button>
			{/each}
		</div>
	</div>
</div>
