<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { getForms } from '$lib/api.js';

	/** @type {any[]} */
	export let forms = [];
	export let showModal = false;

	const dispatch = createEventDispatcher();

	/** @type {HTMLCanvasElement|null} */
	let canvas = null;
	/** @type {CanvasRenderingContext2D|null} */
	let ctx = null;
	/** @type {any} */
	let workflow = {
		id: null,
		name: '',
		description: '',
		formId: '',
		trigger: 'on_submit',
		isActive: false,
		steps: []
	};

	/** @type {any} */
	let selectedStep = null;
	let isDrawingConnection = false;
	/** @type {any} */
	let connectionStart = null;
	/** @type {any} */
	let connectionEnd = null;
	/** @type {any} */
	let draggedStep = null;
	let dragOffset = { x: 0, y: 0 };

	// Step types available
	const stepTypes = [
		{ id: 'approval', name: 'Approval', icon: '✅', description: 'Require approval from a user' },
		{ id: 'notification', name: 'Notification', icon: '📧', description: 'Send email or in-app notification' },
		{ id: 'assignment', name: 'Assignment', icon: '👤', description: 'Assign task to a user' },
		{ id: 'condition', name: 'Condition', icon: '🔀', description: 'Branch workflow based on condition' },
		{ id: 'action', name: 'Action', icon: '⚡', description: 'Execute custom action' }
	];

	// Trigger types
	const triggerTypes = [
		{ id: 'on_submit', name: 'On Form Submit', description: 'When a form is submitted' },
		{ id: 'on_update', name: 'On Form Update', description: 'When a form response is updated' },
		{ id: 'on_approve', name: 'On Approval', description: 'When an approval step completes' },
		{ id: 'scheduled', name: 'Scheduled', description: 'Run on a schedule' }
	];

	onMount(async () => {
		if (showModal) {
			await initializeCanvas();
		}
	});

	$: if (showModal && !canvas) {
		initializeCanvas();
	}

	async function initializeCanvas() {
		await new Promise(resolve => setTimeout(resolve, 100));
		canvas = /** @type {HTMLCanvasElement|null} */ (document.getElementById('workflow-canvas'));
		if (canvas) {
			ctx = canvas.getContext('2d');
			resizeCanvas();
			drawCanvas();
		}
	}

	function resizeCanvas() {
		if (canvas) {
			const container = canvas.parentElement;
			if (container) {
				canvas.width = container.clientWidth;
				canvas.height = container.clientHeight;
			}
		}
	}

	function drawCanvas() {
		if (!ctx || !canvas) return;

		// Clear canvas
		ctx.clearRect(0, 0, canvas.width, canvas.height);

		// Draw connections first
		workflow.steps.forEach((/** @type {any} */ step) => {
			if (step.nextStepId) {
				const nextStep = workflow.steps.find((/** @type {any} */ s) => s.id === step.nextStepId);
				if (nextStep) {
					drawConnection(step, nextStep, '#10B981', 'next');
				}
			}
			if (step.altStepId) {
				const altStep = workflow.steps.find((/** @type {any} */ s) => s.id === step.altStepId);
				if (altStep) {
					drawConnection(step, altStep, '#EF4444', 'alt');
				}
			}
		});

		// Draw steps
		workflow.steps.forEach((/** @type {any} */ step) => drawStep(step));
	}

	/** @param {any} step */
	function drawStep(step) {
		if (!ctx) return;

		const stepType = stepTypes.find(t => t.id === step.type);
		const isSelected = selectedStep?.id === step.id;

		// Step background
		ctx.fillStyle = isSelected ? '#3B82F6' : '#FFFFFF';
		ctx.strokeStyle = isSelected ? '#1D4ED8' : '#D1D5DB';
		ctx.lineWidth = isSelected ? 3 : 2;

		// Rounded rectangle
		const radius = 8;
		ctx.beginPath();
		ctx.roundRect(step.x, step.y, 160, 80, radius);
		ctx.fill();
		ctx.stroke();

		// Step icon and title
		ctx.fillStyle = isSelected ? '#FFFFFF' : '#374151';
		ctx.font = 'bold 12px Arial';
		ctx.fillText(stepType?.icon || '⚡', step.x + 10, step.y + 25);

		ctx.font = '12px Arial';
		ctx.fillText(step.name || stepType?.name || 'Step', step.x + 35, step.y + 25);

		// Step type label
		ctx.fillStyle = isSelected ? '#E0F2FE' : '#6B7280';
		ctx.font = '10px Arial';
		ctx.fillText(step.type.toUpperCase(), step.x + 10, step.y + 45);

		// Connection points
		ctx.fillStyle = '#10B981';
		ctx.beginPath();
		ctx.arc(step.x + 160, step.y + 25, 6, 0, 2 * Math.PI);
		ctx.fill(); // Next step connection point

		if (step.type === 'condition' || step.type === 'approval') {
			ctx.fillStyle = '#EF4444';
			ctx.beginPath();
			ctx.arc(step.x + 160, step.y + 55, 6, 0, 2 * Math.PI);
			ctx.fill(); // Alternative path connection point
		}
	}

	/**
	 * @param {any} fromStep
	 * @param {any} toStep
	 * @param {string} color
	 * @param {string} type
	 */
	function drawConnection(fromStep, toStep, color, type) {
		if (!ctx) return;

		const startX = fromStep.x + 160;
		const startY = fromStep.y + (type === 'alt' ? 55 : 25);
		const endX = toStep.x;
		const endY = toStep.y + 40;

		// Draw curved line
		ctx.strokeStyle = color;
		ctx.lineWidth = 3;
		ctx.beginPath();
		ctx.moveTo(startX, startY);

		// Create curved connection
		const midX = (startX + endX) / 2;
		const midY = Math.min(startY, endY) - 30;

		ctx.quadraticCurveTo(midX, midY, endX, endY);
		ctx.stroke();

		// Draw arrow
		const angle = Math.atan2(endY - midY, endX - midX);
		const arrowLength = 15;

		ctx.beginPath();
		ctx.moveTo(endX, endY);
		ctx.lineTo(
			endX - arrowLength * Math.cos(angle - Math.PI / 6),
			endY - arrowLength * Math.sin(angle - Math.PI / 6)
		);
		ctx.moveTo(endX, endY);
		ctx.lineTo(
			endX - arrowLength * Math.cos(angle + Math.PI / 6),
			endY - arrowLength * Math.sin(angle + Math.PI / 6)
		);
		ctx.stroke();

		// Label the connection
		ctx.fillStyle = color;
		ctx.font = 'bold 10px Arial';
		const labelX = midX - 10;
		const labelY = midY - 10;
		ctx.fillText(type === 'alt' ? 'NO/REJECT' : 'YES/APPROVE', labelX, labelY);
	}

	/** @param {MouseEvent} event */
	function handleCanvasMouseDown(event) {
		if (!canvas) return;
		const rect = canvas.getBoundingClientRect();
		const x = event.clientX - rect.left;
		const y = event.clientY - rect.top;

		// Check if clicking on a step
		for (const step of workflow.steps) {
			if (x >= step.x && x <= step.x + 160 && y >= step.y && y <= step.y + 80) {
				selectedStep = step;
				draggedStep = step;
				dragOffset.x = x - step.x;
				dragOffset.y = y - step.y;
				break;
			}
		}

		// Check if clicking on connection points
		for (const step of workflow.steps) {
			// Next step connection point
			const nextDist = Math.sqrt((x - (step.x + 160)) ** 2 + (y - (step.y + 25)) ** 2);
			if (nextDist <= 8) {
				isDrawingConnection = true;
				connectionStart = { stepId: step.id, type: 'next' };
				break;
			}

			// Alternative path connection point
			if (step.type === 'condition' || step.type === 'approval') {
				const altDist = Math.sqrt((x - (step.x + 160)) ** 2 + (y - (step.y + 55)) ** 2);
				if (altDist <= 8) {
					isDrawingConnection = true;
					connectionStart = { stepId: step.id, type: 'alt' };
					break;
				}
			}
		}

		// If not clicking on anything, deselect
		if (!draggedStep && !isDrawingConnection) {
			selectedStep = null;
		}

		drawCanvas();
	}

	/** @param {MouseEvent} event */
	function handleCanvasMouseMove(event) {
		if (draggedStep) {
			if (!canvas) return;
			const rect = canvas.getBoundingClientRect();
			const x = event.clientX - rect.left;
			const y = event.clientY - rect.top;

			draggedStep.x = x - dragOffset.x;
			draggedStep.y = y - dragOffset.y;

			drawCanvas();
		}

		if (isDrawingConnection) {
			// Could draw temporary connection line here
			drawCanvas();
		}
	}

	/** @param {MouseEvent} event */
	function handleCanvasMouseUp(event) {
		if (draggedStep) {
			draggedStep = null;
		}

		if (isDrawingConnection && connectionStart) {
			if (!canvas) return;
			const rect = canvas.getBoundingClientRect();
			const x = event.clientX - rect.left;
			const y = event.clientY - rect.top;

			// Check if dropped on another step
			for (const step of workflow.steps) {
				if (step.id !== connectionStart.stepId &&
					x >= step.x && x <= step.x + 160 &&
					y >= step.y && y <= step.y + 80) {

					// Create connection
					const sourceStep = workflow.steps.find((/** @type {any} */ s) => s.id === connectionStart.stepId);
					if (connectionStart.type === 'next') {
						sourceStep.nextStepId = step.id;
					} else {
						sourceStep.altStepId = step.id;
					}
					break;
				}
			}

			isDrawingConnection = false;
			connectionStart = null;
			drawCanvas();
		}
	}

	/** @param {any} stepType */
	function addStep(stepType) {
		const step = {
			id: Date.now() + Math.random(),
			type: stepType,
			name: stepTypes.find(t => t.id === stepType)?.name || 'New Step',
			x: Math.random() * 400 + 100,
			y: Math.random() * 300 + 100,
			order: workflow.steps.length,
			config: {}
		};

		workflow.steps = [...workflow.steps, step];
		selectedStep = step;
		drawCanvas();
	}

	/** @param {any} stepId */
	function deleteStep(stepId) {
		workflow.steps = workflow.steps.filter((/** @type {any} */ s) => s.id !== stepId);

		// Remove connections to this step
		workflow.steps.forEach((/** @type {any} */ step) => {
			if (step.nextStepId === stepId) step.nextStepId = null;
			if (step.altStepId === stepId) step.altStepId = null;
		});

		if (selectedStep?.id === stepId) {
			selectedStep = null;
		}

		drawCanvas();
	}

	/**
	 * @param {any} property
	 * @param {any} value
	 */
	function updateStepConfig(property, value) {
		if (selectedStep) {
			selectedStep.config = { ...selectedStep.config, [property]: value };
			workflow.steps = [...workflow.steps]; // Trigger reactivity
		}
	}

	async function saveWorkflow() {
		try {
			const workflowData = {
				name: workflow.name,
				description: workflow.description,
				form_id: workflow.formId,
				trigger: workflow.trigger,
				is_active: workflow.isActive,
				steps: workflow.steps.map((/** @type {any} */ step) => ({
					id: step.id,
					name: step.name,
					type: step.type,
					order: step.order,
					config: step.config,
					next_step_id: step.nextStepId || null,
					alt_step_id: step.altStepId || null
				}))
			};

			const response = await fetch('/api/v1/workflows', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify(workflowData)
			});

			if (response.ok) {
				alert('Workflow saved successfully!');
				closeModal();
			} else {
				const error = await response.json();
				alert('Failed to save workflow: ' + (error.error || 'Unknown error'));
			}
		} catch (/** @type {any} */ error) {
			console.error('Failed to save workflow:', error);
			alert('Failed to save workflow: ' + error.message);
		}
	}

	function resetWorkflow() {
		workflow = /** @type {any} */ ({
			id: null,
			name: '',
			description: '',
			formId: '',
			trigger: 'on_submit',
			isActive: false,
			steps: []
		});
		selectedStep = null;
		drawCanvas();
	}

	function closeModal() {
		showModal = false;
		resetWorkflow();
		dispatch('close');
	}
</script>

{#if showModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-7xl max-h-[90vh] overflow-hidden flex flex-col">
			<!-- Header -->
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<div>
					<h2 class="text-2xl font-bold text-gray-900">⚡ Workflow Designer</h2>
					<p class="text-gray-600 mt-1">Create automated approval chains and business processes</p>
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
					<!-- Add Step Buttons -->
					{#each stepTypes as stepType}
						<button
							class="px-3 py-2 bg-white border border-gray-300 rounded hover:bg-gray-50 text-sm flex items-center space-x-2"
							on:click={() => addStep(stepType.id)}
						>
							<span>{stepType.icon}</span>
							<span>{stepType.name}</span>
						</button>
					{/each}
				</div>

				<div class="flex items-center space-x-4">
					{#if selectedStep}
						<button
							class="px-3 py-2 bg-red-600 text-white rounded hover:bg-red-700 text-sm"
							on:click={() => deleteStep(selectedStep.id)}
						>
							🗑️ Delete Step
						</button>
					{/if}
					<button
						class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
						on:click={saveWorkflow}
						disabled={!workflow.name || !workflow.formId}
					>
						💾 Save Workflow
					</button>
				</div>
			</div>

			<!-- Main Content -->
			<div class="flex-1 flex overflow-hidden">
				<!-- Canvas Area -->
				<div class="flex-1 relative bg-gray-100">
					<canvas
						id="workflow-canvas"
						class="w-full h-full cursor-move"
						on:mousedown={handleCanvasMouseDown}
						on:mousemove={handleCanvasMouseMove}
						on:mouseup={handleCanvasMouseUp}
					></canvas>

					{#if workflow.steps.length === 0}
						<div class="absolute inset-0 flex items-center justify-center">
							<div class="text-center">
								<div class="w-24 h-24 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-4">
									<svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
									</svg>
								</div>
								<h3 class="text-xl font-medium text-gray-900 mb-2">Start Building Your Workflow</h3>
								<p class="text-gray-500 mb-6">Add steps from the toolbar above and connect them to create automated processes.</p>
								<div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-left">
									<div class="bg-blue-50 p-4 rounded-lg">
										<h4 class="font-medium text-blue-900 mb-2">✅ Approvals</h4>
										<p class="text-sm text-blue-700">Require sign-off from managers or team members</p>
									</div>
									<div class="bg-green-50 p-4 rounded-lg">
										<h4 class="font-medium text-green-900 mb-2">📧 Notifications</h4>
										<p class="text-sm text-green-700">Send emails or in-app alerts automatically</p>
									</div>
									<div class="bg-purple-50 p-4 rounded-lg">
										<h4 class="font-medium text-purple-900 mb-2">🔀 Conditions</h4>
										<p class="text-sm text-purple-700">Branch workflows based on form responses</p>
									</div>
								</div>
							</div>
						</div>
					{/if}

					<!-- Instructions Overlay -->
					{#if workflow.steps.length > 0 && !selectedStep}
						<div class="absolute bottom-4 left-4 bg-white p-4 rounded-lg shadow-lg border border-gray-200 max-w-sm">
							<h4 class="font-medium text-gray-900 mb-2">How to Use:</h4>
							<ul class="text-sm text-gray-600 space-y-1">
								<li>• Click and drag steps to reposition them</li>
								<li>• Click connection points (colored dots) to link steps</li>
								<li>• Select a step to configure its properties</li>
								<li>• Green connections = approval path</li>
								<li>• Red connections = rejection/alternative path</li>
							</ul>
						</div>
					{/if}
				</div>

				<!-- Sidebar - Workflow Settings & Step Properties -->
				<div class="w-80 bg-white border-l border-gray-200 p-4 overflow-y-auto">
					{#if selectedStep}
						<!-- Step Properties -->
						<div>
							<h3 class="text-lg font-semibold text-gray-900 mb-4">Step Properties</h3>

							<div class="space-y-4">
								<div>
									<label for="step-name" class="block text-sm font-medium text-gray-700 mb-1">Step Name</label>
									<input
										id="step-name"
										type="text"
										bind:value={selectedStep.name}
										placeholder="Enter step name"
										class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
									/>
								</div>

								<div>
									<label class="block text-sm font-medium text-gray-700 mb-1">Step Type</label>
									<div class="px-3 py-2 bg-gray-100 rounded text-sm">
										{stepTypes.find(t => t.id === selectedStep.type)?.name || selectedStep.type}
									</div>
								</div>

								<!-- Type-specific configuration -->
								{#if selectedStep.type === 'approval'}
									<div>
										<label for="assigned-to" class="block text-sm font-medium text-gray-700 mb-1">Assign To</label>
										<select
											id="assigned-to"
											bind:value={selectedStep.config.assigned_to}
											on:change={(e) => updateStepConfig('assigned_to', /** @type {HTMLInputElement} */ (e.target).value)}
											class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
										>
											<option value="">Select user...</option>
											<option value="manager">Manager</option>
											<option value="admin">Administrator</option>
											<!-- Would load actual users from API -->
										</select>
									</div>
								{:else if selectedStep.type === 'notification'}
									<div class="space-y-3">
										<div>
											<label for="notif-type" class="block text-sm font-medium text-gray-700 mb-1">Notification Type</label>
											<select
												id="notif-type"
												bind:value={selectedStep.config.type}
												on:change={(e) => updateStepConfig('type', /** @type {HTMLInputElement} */ (e.target).value)}
												class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
											>
												<option value="email">Email</option>
												<option value="in_app">In-App Notification</option>
											</select>
										</div>

										<div>
											<label for="recipients" class="block text-sm font-medium text-gray-700 mb-1">Recipients</label>
											<input
												id="recipients"
												type="text"
												bind:value={selectedStep.config.recipients}
												on:input={(e) => updateStepConfig('recipients', /** @type {HTMLInputElement} */ (e.target).value)}
												placeholder="user@example.com, manager"
												class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
											/>
										</div>

										<div>
											<label for="template" class="block text-sm font-medium text-gray-700 mb-1">Template</label>
											<select
												id="template"
												bind:value={selectedStep.config.template_id}
												on:change={(e) => updateStepConfig('template_id', /** @type {HTMLInputElement} */ (e.target).value)}
												class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
											>
												<option value="">Select template...</option>
												<option value="approval-request">Approval Request</option>
												<option value="status-update">Status Update</option>
											</select>
										</div>
									</div>
								{:else if selectedStep.type === 'condition'}
									<div>
										<label for="condition" class="block text-sm font-medium text-gray-700 mb-1">Condition</label>
										<input
											id="condition"
											type="text"
											bind:value={selectedStep.config.condition}
											on:input={(e) => updateStepConfig('condition', /** @type {HTMLInputElement} */ (e.target).value)}
											placeholder="e.g., status == 'approved'"
											class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
										/>
									</div>
								{:else if selectedStep.type === 'action'}
									<div class="space-y-3">
										<div>
											<label for="action-type" class="block text-sm font-medium text-gray-700 mb-1">Action Type</label>
											<select
												id="action-type"
												bind:value={selectedStep.config.action_type}
												on:change={(e) => updateStepConfig('action_type', /** @type {HTMLInputElement} */ (e.target).value)}
												class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
											>
												<option value="update_field">Update Field</option>
												<option value="create_record">Create Record</option>
												<option value="webhook">Call Webhook</option>
											</select>
										</div>

										{#if selectedStep.config.action_type === 'webhook'}
											<div>
												<label for="webhook-url" class="block text-sm font-medium text-gray-700 mb-1">Webhook URL</label>
												<input
													id="webhook-url"
													type="url"
													bind:value={selectedStep.config.url}
													on:input={(e) => updateStepConfig('url', /** @type {HTMLInputElement} */ (e.target).value)}
													placeholder="https://api.example.com/webhook"
													class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
												/>
											</div>
										{/if}
									</div>
								{/if}
							</div>
						</div>
					{:else}
						<!-- Workflow Settings -->
						<div>
							<h3 class="text-lg font-semibold text-gray-900 mb-4">Workflow Settings</h3>

							<div class="space-y-4">
								<div>
									<label for="workflow-name" class="block text-sm font-medium text-gray-700 mb-1">Workflow Name</label>
									<input
										id="workflow-name"
										type="text"
										bind:value={workflow.name}
										placeholder="Enter workflow name"
										class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
									/>
								</div>

								<div>
									<label for="workflow-description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
									<textarea
										id="workflow-description"
										bind:value={workflow.description}
										placeholder="Describe this workflow"
										rows="3"
										class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
									></textarea>
								</div>

								<div>
									<label for="form-select" class="block text-sm font-medium text-gray-700 mb-1">Target Form</label>
									<select
										id="form-select"
										bind:value={workflow.formId}
										class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
									>
										<option value="">Select form</option>
										{#each forms as form}
											<option value={form.id}>{form.name}</option>
										{/each}
									</select>
								</div>

								<div>
									<label for="trigger-select" class="block text-sm font-medium text-gray-700 mb-1">Trigger</label>
									<select
										id="trigger-select"
										bind:value={workflow.trigger}
										class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
									>
										{#each triggerTypes as trigger}
											<option value={trigger.id}>{trigger.name}</option>
										{/each}
									</select>
									<p class="text-xs text-gray-500 mt-1">
										{triggerTypes.find(t => t.id === workflow.trigger)?.description}
									</p>
								</div>

								<div class="flex items-center">
									<input
										type="checkbox"
										id="workflow-active"
										bind:checked={workflow.isActive}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									/>
									<label for="workflow-active" class="ml-2 text-sm text-gray-700">
										Activate workflow
									</label>
								</div>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}
