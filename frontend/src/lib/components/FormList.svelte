<script>
	import { createEventDispatcher } from 'svelte';

	export let forms = [];

	const dispatch = createEventDispatcher();

	let draggedItem = null;
	let draggedIndex = -1;

	function handleDragStart(event, form, index) {
		draggedItem = form;
		draggedIndex = index;
		event.dataTransfer.effectAllowed = 'move';
		event.dataTransfer.setData('text/plain', form.id.toString());

		// Add visual feedback
		event.target.style.opacity = '0.5';
	}

	function handleDragOver(event) {
		event.preventDefault();
		event.dataTransfer.dropEffect = 'move';
	}

	function handleDragEnter(event, index) {
		// Visual feedback for drop target
		const card = event.target.closest('.form-card');
		if (card) card.style.borderColor = '#3b82f6';
	}

	function handleDragLeave(event) {
		// Remove visual feedback
		const card = event.target.closest('.form-card');
		if (card) card.style.borderColor = '#e5e7eb';
	}

	function handleDrop(event, dropIndex) {
		event.preventDefault();

		const draggedId = parseInt(event.dataTransfer.getData('text/plain'));
		const draggedForm = forms.find(f => f.id === draggedId);

		if (draggedForm && draggedIndex !== dropIndex) {
			// Remove dragged item from current position
			const newForms = forms.filter(f => f.id !== draggedId);

			// Insert at new position
			newForms.splice(dropIndex, 0, draggedForm);

			forms = newForms;
			dispatch('reorder', { forms });
		}

		// Reset visual feedback
		document.querySelectorAll('.form-card').forEach(card => {
			card.style.borderColor = '#e5e7eb';
			if (card.style.opacity === '0.5') {
				card.style.opacity = '1';
			}
		});

		draggedItem = null;
		draggedIndex = -1;
	}

	function handleDragEnd(event) {
		// Reset visual feedback
		document.querySelectorAll('.form-card').forEach(card => {
			card.style.borderColor = '#e5e7eb';
			card.style.opacity = '1';
		});

		draggedItem = null;
		draggedIndex = -1;
	}

	function handleEdit(form) {
		dispatch('edit', form);
	}

	function handleDelete(form) {
		if (confirm(`Are you sure you want to delete "${form.name}"? This action cannot be undone.`)) {
			dispatch('delete', form.id);
		}
	}

	function handleViewResponses(form) {
		dispatch('viewResponses', form);
	}

	function handleDatabaseRelations() {
		dispatch('openDatabaseRelations');
	}

	function handleAdvancedReports() {
		dispatch('openAdvancedReports');
	}

	function handleWorkflowDesigner() {
		dispatch('openWorkflowDesigner');
	}

	function handleAdvanced(form) {
		dispatch('openAdvanced', form);
	}

	function getStatusColor(status) {
		switch (status) {
			case 'active': return 'bg-green-100 text-green-800';
			case 'draft': return 'bg-yellow-100 text-yellow-800';
			case 'archived': return 'bg-gray-100 text-gray-800';
			default: return 'bg-gray-100 text-gray-800';
		}
	}
</script>

<div class="p-6">
	<div class="mb-8">
		<div class="flex items-center justify-between mb-4">
			<div>
				<h2 class="text-2xl font-bold text-gray-900 mb-2">📋 Advanced Form Management</h2>
				<p class="text-gray-600">MS Access-style database features with relationships, workflows, and advanced reporting</p>
			</div>
			<div class="flex space-x-3">
				<button
					class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors flex items-center space-x-2"
					on:click={handleDatabaseRelations}
				>
					<span>🔗</span>
					<span>Database Relations</span>
				</button>
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-2"
					on:click={handleAdvancedReports}
				>
					<span>📊</span>
					<span>Advanced Reports</span>
				</button>
				<button
					class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors flex items-center space-x-2"
					on:click={handleWorkflowDesigner}
				>
					<span>⚡</span>
					<span>Workflow Designer</span>
				</button>
			</div>
		</div>

		<!-- Feature Highlights -->
		<div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
			<div class="bg-blue-50 p-4 rounded-lg">
				<div class="flex items-center space-x-2 mb-2">
					<span class="text-blue-600">🔗</span>
					<h3 class="font-medium text-blue-900">Relationships</h3>
				</div>
				<p class="text-sm text-blue-700">One-to-one, one-to-many, many-to-many database relationships</p>
			</div>
			<div class="bg-green-50 p-4 rounded-lg">
				<div class="flex items-center space-x-2 mb-2">
					<span class="text-green-600">📊</span>
					<h3 class="font-medium text-green-900">Advanced Queries</h3>
				</div>
				<p class="text-sm text-green-700">Visual query builder with joins, filters, and aggregations</p>
			</div>
			<div class="bg-purple-50 p-4 rounded-lg">
				<div class="flex items-center space-x-2 mb-2">
					<span class="text-purple-600">⚡</span>
					<h3 class="font-medium text-purple-900">Workflow Automation</h3>
				</div>
				<p class="text-sm text-purple-700">Approval chains, notifications, and automated processes</p>
			</div>
			<div class="bg-orange-50 p-4 rounded-lg">
				<div class="flex items-center space-x-2 mb-2">
					<span class="text-orange-600">📈</span>
					<h3 class="font-medium text-orange-900">Smart Reporting</h3>
				</div>
				<p class="text-sm text-orange-700">Dashboards, charts, and export to Excel/PDF</p>
			</div>
		</div>
	</div>

	{#if forms.length === 0}
		<!-- Empty state -->
		<div class="text-center py-12">
			<div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
				<svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
				</svg>
			</div>
			<h3 class="text-lg font-medium text-gray-900 mb-2">No forms yet</h3>
			<p class="text-gray-500 mb-4">Create your first form to start collecting data</p>
			<button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
				Create Your First Form
			</button>
		</div>
	{:else}
		<!-- Forms grid with native drag and drop -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6">
			{#each forms as form, index (form.id)}
				<div
					class="form-card bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow cursor-move"
					draggable="true"
					on:dragstart={(e) => handleDragStart(e, form, index)}
					on:dragover={handleDragOver}
					on:dragenter={(e) => handleDragEnter(e, index)}
					on:dragleave={handleDragLeave}
					on:drop={(e) => handleDrop(e, index)}
					on:dragend={handleDragEnd}
					role="button"
					aria-label="Drag form '{form.name}' to reorder"
					tabindex="0"
					on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); /* Could trigger drag start */ } }}
				>
					<!-- Drag handle -->
					<div class="flex items-center justify-between mb-4">
						<div class="flex items-center space-x-2">
							<svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16"></path>
							</svg>
							<span class="text-xs text-gray-500">Drag to reorder</span>
						</div>
						<span class="px-2 py-1 text-xs font-medium rounded-full {getStatusColor(form.status)}">
							{form.status}
						</span>
					</div>

							<div class="flex items-start justify-between mb-4">
								<div class="flex-1">
									<h3 class="text-lg font-semibold text-gray-900 mb-1">{form.name}</h3>
									<p class="text-sm text-gray-600 mb-2">{form.description}</p>
									<div class="flex items-center space-x-4 text-sm text-gray-500">
										<span>{form.responses} responses</span>
										<span>Created {form.createdAt.toLocaleDateString()}</span>
									</div>
								</div>
							</div>

							<div class="flex items-center justify-between">
								<div class="flex space-x-2">
									<button
										class="px-3 py-1 text-sm bg-blue-50 text-blue-700 rounded hover:bg-blue-100 transition-colors"
										on:click={() => handleViewResponses(form)}
									>
										View Responses
									</button>
								<button
									class="px-3 py-1 text-sm bg-gray-50 text-gray-700 rounded hover:bg-gray-100 transition-colors"
									on:click={() => handleEdit(form)}
								>
									Edit
								</button>
								<button
									class="px-3 py-1 text-sm bg-indigo-50 text-indigo-700 rounded hover:bg-indigo-100 transition-colors"
									on:click={() => handleAdvanced(form)}
								>
									Advanced
								</button>
								</div>
								<button
									class="p-1 text-gray-400 hover:text-red-600 transition-colors"
									on:click={() => handleDelete(form)}
									title="Delete form"
								>
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
									</svg>
								</button>
							</div>
						</div>
					{/each}
				</div>

		<!-- Stats summary -->
		<div class="mt-8 bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			<h3 class="text-lg font-semibold text-gray-900 mb-4">Summary</h3>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
				<div class="text-center">
					<div class="text-2xl font-bold text-blue-600">{forms.length}</div>
					<div class="text-sm text-gray-600">Total Forms</div>
				</div>
				<div class="text-center">
					<div class="text-2xl font-bold text-green-600">{forms.reduce((sum, f) => sum + f.responses, 0)}</div>
					<div class="text-sm text-gray-600">Total Responses</div>
				</div>
				<div class="text-center">
					<div class="text-2xl font-bold text-purple-600">{forms.filter(f => f.status === 'active').length}</div>
					<div class="text-sm text-gray-600">Active Forms</div>
				</div>
			</div>
		</div>
	{/if}
</div>
