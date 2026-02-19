<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { getForms } from '$lib/api.js';

	export let forms = [];
	export let showModal = false;

	const dispatch = createEventDispatcher();

	let relationships = [];
	let selectedRelationship = null;
	let showCreateRelationship = false;
	let showLookupCreator = false;
	let canvas = null;
	let ctx = null;
	let nodes = [];
	let connections = [];
	let isDragging = false;
	let draggedNode = null;
	let dragOffset = { x: 0, y: 0 };

	// Relationship creation
	let newRelationship = {
		sourceFormId: '',
		targetFormId: '',
		sourceField: '',
		targetField: '',
		relationshipType: 'one-to-many',
		name: '',
		description: '',
		cascadeDelete: false,
		cascadeUpdate: false
	};

	// Lookup field creation
	let newLookupField = {
		formId: '',
		fieldId: '',
		relatedFormId: '',
		relatedFieldId: '',
		displayField: '',
		allowMultiple: false,
		filterField: '',
		filterValue: ''
	};

	onMount(async () => {
		if (showModal) {
			await initializeCanvas();
			await loadRelationships();
		}
	});

	$: if (showModal && !canvas) {
		initializeCanvas();
		loadRelationships();
	}

	async function initializeCanvas() {
		await new Promise(resolve => setTimeout(resolve, 100)); // Wait for DOM
		canvas = document.getElementById('relationship-canvas');
		if (canvas) {
			ctx = canvas.getContext('2d');
			resizeCanvas();
			setupNodes();
			drawCanvas();
		}
	}

	function resizeCanvas() {
		if (canvas) {
			const container = canvas.parentElement;
			canvas.width = container.clientWidth;
			canvas.height = container.clientHeight;
		}
	}

	function setupNodes() {
		nodes = forms.map((form, index) => ({
			id: form.id,
			name: form.name,
			x: 100 + (index % 3) * 250,
			y: 100 + Math.floor(index / 3) * 200,
			width: 180,
			height: 100,
			color: getRandomColor(),
			fields: extractFieldNames(form)
		}));
	}

	function extractFieldNames(form) {
		// Extract field names from form data
		if (form.fields && Array.isArray(form.fields)) {
			return form.fields.map(f => f.label || f.name || f.id);
		}
		return ['id', 'name', 'created_at'];
	}


	function getRandomColor() {
		const colors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#06B6D4'];
		return colors[Math.floor(Math.random() * colors.length)];
	}

	function drawCanvas() {
		if (!ctx || !canvas) return;

		// Clear canvas
		ctx.clearRect(0, 0, canvas.width, canvas.height);

		// Draw connections first
		connections.forEach(connection => drawConnection(connection));

		// Draw nodes
		nodes.forEach(node => drawNode(node));
	}

	function drawNode(node) {
		if (!ctx) return;

		// Node background
		ctx.fillStyle = node.color;
		ctx.fillRect(node.x, node.y, node.width, node.height);

		// Node border
		ctx.strokeStyle = '#374151';
		ctx.lineWidth = 2;
		ctx.strokeRect(node.x, node.y, node.width, node.height);

		// Node title
		ctx.fillStyle = 'white';
		ctx.font = 'bold 14px Arial';
		ctx.fillText(node.name, node.x + 10, node.y + 25);

		// Node fields
		ctx.fillStyle = '#F3F4F6';
		ctx.font = '12px Arial';
		node.fields.slice(0, 4).forEach((field, index) => {
			ctx.fillText(`• ${field}`, node.x + 10, node.y + 45 + index * 15);
		});

		if (node.fields.length > 4) {
			ctx.fillText(`... +${node.fields.length - 4} more`, node.x + 10, node.y + 45 + 4 * 15);
		}
	}

	function drawConnection(connection) {
		if (!ctx) return;

		const sourceNode = nodes.find(n => n.id === connection.sourceId);
		const targetNode = nodes.find(n => n.id === connection.targetId);

		if (!sourceNode || !targetNode) return;

		const sourceX = sourceNode.x + sourceNode.width;
		const sourceY = sourceNode.y + sourceNode.height / 2;
		const targetX = targetNode.x;
		const targetY = targetNode.y + targetNode.height / 2;

		// Draw line
		ctx.strokeStyle = connection.color || '#6B7280';
		ctx.lineWidth = 3;
		ctx.beginPath();
		ctx.moveTo(sourceX, sourceY);
		ctx.lineTo(targetX, targetY);
		ctx.stroke();

		// Draw arrow
		const angle = Math.atan2(targetY - sourceY, targetX - sourceX);
		const arrowLength = 15;

		ctx.beginPath();
		ctx.moveTo(targetX, targetY);
		ctx.lineTo(
			targetX - arrowLength * Math.cos(angle - Math.PI / 6),
			targetY - arrowLength * Math.sin(angle - Math.PI / 6)
		);
		ctx.moveTo(targetX, targetY);
		ctx.lineTo(
			targetX - arrowLength * Math.cos(angle + Math.PI / 6),
			targetY - arrowLength * Math.sin(angle + Math.PI / 6)
		);
		ctx.stroke();

		// Draw relationship type label
		const midX = (sourceX + targetX) / 2;
		const midY = (sourceY + targetY) / 2;

		ctx.fillStyle = '#374151';
		ctx.font = 'bold 12px Arial';
		ctx.fillText(connection.type, midX - 20, midY - 5);
	}

	function handleCanvasMouseDown(event) {
		const rect = canvas.getBoundingClientRect();
		const x = event.clientX - rect.left;
		const y = event.clientY - rect.top;

		// Check if clicking on a node
		for (const node of nodes) {
			if (x >= node.x && x <= node.x + node.width &&
				y >= node.y && y <= node.y + node.height) {
				isDragging = true;
				draggedNode = node;
				dragOffset.x = x - node.x;
				dragOffset.y = y - node.y;
				break;
			}
		}
	}

	function handleCanvasMouseMove(event) {
		if (!isDragging || !draggedNode) return;

		const rect = canvas.getBoundingClientRect();
		const x = event.clientX - rect.left;
		const y = event.clientY - rect.top;

		draggedNode.x = x - dragOffset.x;
		draggedNode.y = y - dragOffset.y;

		drawCanvas();
	}

	function handleCanvasMouseUp() {
		isDragging = false;
		draggedNode = null;
	}

	async function loadRelationships() {
		try {
			// Load relationships for each form
			const allRelationships = [];
			for (const form of forms) {
				try {
					const response = await fetch(`/api/v1/form-relationships/${form.id}`, {
						headers: {
							'Authorization': `Bearer ${localStorage.getItem('token')}`
						}
					});
					if (response.ok) {
						const data = await response.json();
						allRelationships.push(...data.relationships);
					}
				} catch (error) {
					console.warn(`Failed to load relationships for form ${form.id}:`, error);
				}
			}

			relationships = allRelationships;
			setupConnections();
			drawCanvas();
		} catch (error) {
			console.error('Failed to load relationships:', error);
		}
	}

	function setupConnections() {
		connections = relationships.map(rel => {
			const sourceForm = forms.find(f => f.id === rel.source_form_id);
			const targetForm = forms.find(f => f.id === rel.target_form_id);

			let color = '#6B7280';
			let type = rel.relationship_type;

			switch (rel.relationship_type) {
				case 'one-to-one':
					color = '#10B981';
					type = '1:1';
					break;
				case 'one-to-many':
					color = '#3B82F6';
					type = '1:N';
					break;
				case 'many-to-many':
					color = '#F59E0B';
					type = 'N:N';
					break;
			}

			return {
				sourceId: rel.source_form_id,
				targetId: rel.target_form_id,
				type: type,
				color: color,
				relationship: rel
			};
		});
	}

	async function createRelationship() {
		try {
			const response = await fetch('/api/v1/form-relationships', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					name: newRelationship.name,
					description: newRelationship.description,
					source_form_id: newRelationship.sourceFormId,
					source_field: newRelationship.sourceField,
					target_form_id: newRelationship.targetFormId,
					target_field: newRelationship.targetField,
					relationship_type: newRelationship.relationshipType,
					cascade_delete: newRelationship.cascadeDelete,
					cascade_update: newRelationship.cascadeUpdate
				})
			});

			if (response.ok) {
				showCreateRelationship = false;
				resetNewRelationship();
				await loadRelationships();
			} else {
				const error = await response.json();
				alert('Failed to create relationship: ' + (error.error || 'Unknown error'));
			}
		} catch (error) {
			console.error('Failed to create relationship:', error);
			alert('Failed to create relationship: ' + error.message);
		}
	}

	function resetNewRelationship() {
		newRelationship = {
			sourceFormId: '',
			targetFormId: '',
			sourceField: '',
			targetField: '',
			relationshipType: 'one-to-many',
			name: '',
			description: '',
			cascadeDelete: false,
			cascadeUpdate: false
		};
	}

	async function createLookupField() {
		try {
			const response = await fetch('/api/v1/form-lookup-fields', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					form_id: newLookupField.formId,
					field_id: newLookupField.fieldId,
					related_form_id: newLookupField.relatedFormId,
					related_field_id: newLookupField.relatedFieldId,
					display_field: newLookupField.displayField,
					filter_field: newLookupField.filterField || null,
					filter_value: newLookupField.filterValue || null,
					allow_multiple: newLookupField.allowMultiple
				})
			});

			if (response.ok) {
				showLookupCreator = false;
				resetNewLookupField();
				alert('Lookup field created successfully!');
			} else {
				const error = await response.json();
				alert('Failed to create lookup field: ' + (error.error || 'Unknown error'));
			}
		} catch (error) {
			console.error('Failed to create lookup field:', error);
			alert('Failed to create lookup field: ' + error.message);
		}
	}

	function resetNewLookupField() {
		newLookupField = {
			formId: '',
			fieldId: '',
			relatedFormId: '',
			relatedFieldId: '',
			displayField: '',
			allowMultiple: false,
			filterField: '',
			filterValue: ''
		};
	}

	function getFormFields(formId) {
		const form = forms.find(f => f.id === formId);
		if (form && form.fields) {
			return form.fields.map(f => ({ id: f.id, label: f.label || f.name || f.id }));
		}
		return [];
	}

	function closeModal() {
		showModal = false;
		dispatch('close');
	}
</script>

{#if showModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-7xl max-h-[90vh] overflow-hidden flex flex-col">
			<!-- Header -->
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<div>
					<h2 class="text-2xl font-bold text-gray-900">🔗 Database Relations Manager</h2>
					<p class="text-gray-600 mt-1">Create and manage relationships between your forms</p>
				</div>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={closeModal}
				>
					×
				</button>
			</div>

			<!-- Toolbar -->
			<div class="flex items-center justify-between p-4 bg-gray-50 border-b border-gray-200">
				<div class="flex items-center space-x-3">
					<button
						class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
						on:click={() => showCreateRelationship = true}
					>
						➕ Create Relationship
					</button>
					<button
						class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
						on:click={() => showLookupCreator = true}
					>
						🔍 Create Lookup Field
					</button>
				</div>

				<div class="flex items-center space-x-4 text-sm">
					<div class="flex items-center space-x-2">
						<div class="w-4 h-4 bg-green-500 rounded"></div>
						<span>One-to-One</span>
					</div>
					<div class="flex items-center space-x-2">
						<div class="w-4 h-4 bg-blue-500 rounded"></div>
						<span>One-to-Many</span>
					</div>
					<div class="flex items-center space-x-2">
						<div class="w-4 h-4 bg-yellow-500 rounded"></div>
						<span>Many-to-Many</span>
					</div>
				</div>
			</div>

			<!-- Main Content -->
			<div class="flex-1 flex overflow-hidden">
				<!-- Canvas Area -->
				<div class="flex-1 relative bg-gray-100">
					<canvas
						id="relationship-canvas"
						class="w-full h-full cursor-move"
						on:mousedown={handleCanvasMouseDown}
						on:mousemove={handleCanvasMouseMove}
						on:mouseup={handleCanvasMouseUp}
					></canvas>

					{#if nodes.length === 0}
						<div class="absolute inset-0 flex items-center justify-center">
							<div class="text-center">
								<div class="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-4">
									<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"></path>
									</svg>
								</div>
								<h3 class="text-lg font-medium text-gray-900 mb-2">No forms to relate</h3>
								<p class="text-gray-500">Create some forms first to establish relationships</p>
							</div>
						</div>
					{/if}
				</div>

				<!-- Sidebar -->
				<div class="w-80 bg-white border-l border-gray-200 p-4 overflow-y-auto">
					<h3 class="text-lg font-semibold text-gray-900 mb-4">Relationships</h3>

					{#if relationships.length === 0}
						<div class="text-center py-8 text-gray-500">
							<p>No relationships created yet.</p>
							<p class="text-sm mt-1">Click "Create Relationship" to get started.</p>
						</div>
					{:else}
						<div class="space-y-3">
							{#each relationships as rel}
								<div class="border border-gray-200 rounded-lg p-3 hover:bg-gray-50">
									<div class="flex items-center justify-between mb-2">
										<h4 class="font-medium text-gray-900">{rel.name}</h4>
										<span class="px-2 py-1 text-xs rounded-full {rel.relationship_type === 'one-to-one' ? 'bg-green-100 text-green-800' : rel.relationship_type === 'one-to-many' ? 'bg-blue-100 text-blue-800' : 'bg-yellow-100 text-yellow-800'}">
											{rel.relationship_type}
										</span>
									</div>
									<div class="text-sm text-gray-600">
										<div>{forms.find(f => f.id === rel.source_form_id)?.name} → {forms.find(f => f.id === rel.target_form_id)?.name}</div>
										<div class="mt-1 text-xs">
											{rel.cascade_delete ? '🗑️ Cascade Delete' : ''} {rel.cascade_update ? '🔄 Cascade Update' : ''}
										</div>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Create Relationship Modal -->
{#if showCreateRelationship}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-lg">
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<h3 class="text-xl font-bold text-gray-900">Create New Relationship</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => { showCreateRelationship = false; resetNewRelationship(); }}
				>
					×
				</button>
			</div>

			<div class="p-6 space-y-4">
				<div>
					<label for="rel-name" class="block text-sm font-medium text-gray-700 mb-1">Relationship Name</label>
					<input
						id="rel-name"
						type="text"
						bind:value={newRelationship.name}
						placeholder="e.g., User has many Orders"
						class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>

				<div>
					<label for="rel-description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
					<input
						id="rel-description"
						type="text"
						bind:value={newRelationship.description}
						placeholder="Optional description"
						class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="source-form" class="block text-sm font-medium text-gray-700 mb-1">Source Form</label>
						<select
							id="source-form"
							bind:value={newRelationship.sourceFormId}
							class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
						>
							<option value="">Select source form</option>
							{#each forms as form}
								<option value={form.id}>{form.name}</option>
							{/each}
						</select>
					</div>

					<div>
						<label for="target-form" class="block text-sm font-medium text-gray-700 mb-1">Target Form</label>
						<select
							id="target-form"
							bind:value={newRelationship.targetFormId}
							class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
						>
							<option value="">Select target form</option>
							{#each forms as form}
								<option value={form.id}>{form.name}</option>
							{/each}
						</select>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="source-field" class="block text-sm font-medium text-gray-700 mb-1">Source Field</label>
						<select
							id="source-field"
							bind:value={newRelationship.sourceField}
							class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
							disabled={!newRelationship.sourceFormId}
						>
							<option value="">Select field</option>
							{#each getFormFields(newRelationship.sourceFormId) as field}
								<option value={field.id}>{field.label}</option>
							{/each}
						</select>
					</div>

					<div>
						<label for="target-field" class="block text-sm font-medium text-gray-700 mb-1">Target Field</label>
						<select
							id="target-field"
							bind:value={newRelationship.targetField}
							class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
							disabled={!newRelationship.targetFormId}
						>
							<option value="">Select field</option>
							{#each getFormFields(newRelationship.targetFormId) as field}
								<option value={field.id}>{field.label}</option>
							{/each}
						</select>
					</div>
				</div>

				<div>
					<label for="rel-type" class="block text-sm font-medium text-gray-700 mb-1">Relationship Type</label>
					<select
						id="rel-type"
						bind:value={newRelationship.relationshipType}
						class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
					>
						<option value="one-to-one">One-to-One (1:1)</option>
						<option value="one-to-many">One-to-Many (1:N)</option>
						<option value="many-to-many">Many-to-Many (N:N)</option>
					</select>
				</div>

				<div class="space-y-2">
					<div class="flex items-center">
						<input
							type="checkbox"
							id="cascade-delete"
							bind:checked={newRelationship.cascadeDelete}
							class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
						/>
						<label for="cascade-delete" class="ml-2 text-sm text-gray-700">
							Cascade Delete (delete related records when parent is deleted)
						</label>
					</div>

					<div class="flex items-center">
						<input
							type="checkbox"
							id="cascade-update"
							bind:checked={newRelationship.cascadeUpdate}
							class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
						/>
						<label for="cascade-update" class="ml-2 text-sm text-gray-700">
							Cascade Update (update related records when parent changes)
						</label>
					</div>
				</div>
			</div>

			<div class="flex justify-end space-x-3 p-6 border-t border-gray-200 bg-gray-50">
				<button
					class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50"
					on:click={() => { showCreateRelationship = false; resetNewRelationship(); }}
				>
					Cancel
				</button>
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
					on:click={createRelationship}
					disabled={!newRelationship.name || !newRelationship.sourceFormId || !newRelationship.targetFormId}
				>
					Create Relationship
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Create Lookup Field Modal -->
{#if showLookupCreator}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-lg">
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<h3 class="text-xl font-bold text-gray-900">Create Lookup Field</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => { showLookupCreator = false; resetNewLookupField(); }}
				>
					×
				</button>
			</div>

			<div class="p-6 space-y-4">
				<div>
					<label for="lookup-form" class="block text-sm font-medium text-gray-700 mb-1">Form to add lookup to</label>
					<select
						id="lookup-form"
						bind:value={newLookupField.formId}
						class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
					>
						<option value="">Select form</option>
						{#each forms as form}
							<option value={form.id}>{form.name}</option>
						{/each}
					</select>
				</div>

				<div>
					<label for="lookup-field" class="block text-sm font-medium text-gray-700 mb-1">Field Name</label>
					<input
						id="lookup-field"
						type="text"
						bind:value={newLookupField.fieldId}
						placeholder="e.g., customer_id"
						class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>

				<div>
					<label for="lookup-related-form" class="block text-sm font-medium text-gray-700 mb-1">Related Form</label>
					<select
						id="lookup-related-form"
						bind:value={newLookupField.relatedFormId}
						class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
					>
						<option value="">Select related form</option>
						{#each forms.filter(f => f.id !== newLookupField.formId) as form}
							<option value={form.id}>{form.name}</option>
						{/each}
					</select>
				</div>

				<div>
					<label for="lookup-related-field" class="block text-sm font-medium text-gray-700 mb-1">Related Field</label>
					<select
						id="lookup-related-field"
						bind:value={newLookupField.relatedFieldId}
						class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
						disabled={!newLookupField.relatedFormId}
					>
						<option value="">Select field to store</option>
						{#each getFormFields(newLookupField.relatedFormId) as field}
							<option value={field.id}>{field.label}</option>
						{/each}
					</select>
				</div>

				<div>
					<label for="lookup-display-field" class="block text-sm font-medium text-gray-700 mb-1">Display Field</label>
					<select
						id="lookup-display-field"
						bind:value={newLookupField.displayField}
						class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
						disabled={!newLookupField.relatedFormId}
					>
						<option value="">Select field to display</option>
						{#each getFormFields(newLookupField.relatedFormId) as field}
							<option value={field.label}>{field.label}</option>
						{/each}
					</select>
				</div>

				<div class="flex items-center">
					<input
						type="checkbox"
						id="lookup-multiple"
						bind:checked={newLookupField.allowMultiple}
						class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
					/>
					<label for="lookup-multiple" class="ml-2 text-sm text-gray-700">
						Allow multiple selections
					</label>
				</div>
			</div>

			<div class="flex justify-end space-x-3 p-6 border-t border-gray-200 bg-gray-50">
				<button
					class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50"
					on:click={() => { showLookupCreator = false; resetNewLookupField(); }}
				>
					Cancel
				</button>
				<button
					class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
					on:click={createLookupField}
					disabled={!newLookupField.formId || !newLookupField.fieldId || !newLookupField.relatedFormId}
				>
					Create Lookup Field
				</button>
			</div>
		</div>
	</div>
{/if}
