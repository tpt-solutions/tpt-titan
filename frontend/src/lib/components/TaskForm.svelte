<script>
	import { createEventDispatcher } from 'svelte';

	export let task = null; // Existing task to edit, or null for new
	export let projects = [];

	const dispatch = createEventDispatcher();

	let formData = {
		title: task?.title || '',
		description: task?.description || '',
		priority: task?.priority || 'medium',
		dueDate: task?.dueDate ? new Date(task.dueDate).toISOString().split('T')[0] : '',
		projectId: task?.projectId || null,
		assignedTo: task?.assignedTo || '',
		tags: task?.tags || [],
		subtasks: task?.subtasks || []
	};

	let newTag = '';
	let newSubtask = '';

	function addTag() {
		if (newTag.trim() && !formData.tags.includes(newTag.trim())) {
			formData.tags = [...formData.tags, newTag.trim()];
			newTag = '';
		}
	}

	function removeTag(tag) {
		formData.tags = formData.tags.filter(t => t !== tag);
	}

	function addSubtask() {
		if (newSubtask.trim()) {
			formData.subtasks = [...formData.subtasks, {
				id: Date.now() + Math.random(),
				title: newSubtask.trim(),
				completed: false
			}];
			newSubtask = '';
		}
	}

	function removeSubtask(subtaskId) {
		formData.subtasks = formData.subtasks.filter(st => st.id !== subtaskId);
	}

	function toggleSubtask(subtaskId) {
		formData.subtasks = formData.subtasks.map(st =>
			st.id === subtaskId ? { ...st, completed: !st.completed } : st
		);
	}

	function saveTask() {
		// Basic validation
		if (!formData.title.trim()) {
			alert('Please enter a task title');
			return;
		}

		// Convert due date back to Date object if provided
		const taskData = {
			...formData,
			dueDate: formData.dueDate ? new Date(formData.dueDate) : null
		};

		dispatch('save', taskData);
	}

	function cancel() {
		dispatch('cancel');
	}

	function handleKeyDown(event) {
		if (event.key === 'Enter' && (event.target === document.querySelector('input[placeholder="Add tag..."]'))) {
			event.preventDefault();
			addTag();
		} else if (event.key === 'Enter' && (event.target === document.querySelector('input[placeholder="Add subtask..."]'))) {
			event.preventDefault();
			addSubtask();
		}
	}
</script>

<svelte:window on:keydown={handleKeyDown} />

<div class="h-full overflow-y-auto bg-white">
	<div class="max-w-2xl mx-auto p-6">
		<div class="mb-6">
			<h2 class="text-2xl font-bold text-gray-900 mb-2">
				{task ? 'Edit Task' : 'Create New Task'}
			</h2>
			<p class="text-gray-600">
				{task ? 'Update task details and settings' : 'Add a new task to your project'}
			</p>
		</div>

		<form on:submit|preventDefault={saveTask} class="space-y-6">
			<!-- Basic Information -->
			<div class="space-y-4">
				<div>
					<label for="title" class="block text-sm font-medium text-gray-700 mb-1">
						Task Title *
					</label>
					<input
						id="title"
						bind:value={formData.title}
						type="text"
						placeholder="Enter task title..."
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						required
					/>
				</div>

				<div>
					<label for="description" class="block text-sm font-medium text-gray-700 mb-1">
						Description
					</label>
					<textarea
						id="description"
						bind:value={formData.description}
						rows="3"
						placeholder="Describe the task in detail..."
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 resize-none"
					></textarea>
				</div>
			</div>

			<!-- Task Settings -->
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div>
					<label for="priority" class="block text-sm font-medium text-gray-700 mb-1">
						Priority
					</label>
					<select
						id="priority"
						bind:value={formData.priority}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
					>
						<option value="low">Low</option>
						<option value="medium">Medium</option>
						<option value="high">High</option>
					</select>
				</div>

				<div>
					<label for="dueDate" class="block text-sm font-medium text-gray-700 mb-1">
						Due Date
					</label>
					<input
						id="dueDate"
						bind:value={formData.dueDate}
						type="date"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>

				<div>
					<label for="project" class="block text-sm font-medium text-gray-700 mb-1">
						Project
					</label>
					<select
						id="project"
						bind:value={formData.projectId}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
					>
						<option value={null}>No Project</option>
						{#each projects as project}
							<option value={project.id}>{project.name}</option>
						{/each}
					</select>
				</div>

				<div>
					<label for="assignedTo" class="block text-sm font-medium text-gray-700 mb-1">
						Assigned To
					</label>
					<input
						id="assignedTo"
						bind:value={formData.assignedTo}
						type="text"
						placeholder="Enter name or email..."
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>
			</div>

			<!-- Tags -->
			<div>
				<label class="block text-sm font-medium text-gray-700 mb-2">
					Tags
				</label>
				<div class="flex flex-wrap gap-2 mb-2">
					{#each formData.tags as tag}
						<span class="inline-flex items-center px-2 py-1 rounded-full text-sm bg-blue-100 text-blue-800">
							{tag}
							<button
								type="button"
								class="ml-1 text-blue-600 hover:text-blue-800"
								on:click={() => removeTag(tag)}
							>
								×
							</button>
						</span>
					{/each}
				</div>
				<div class="flex gap-2">
					<input
						bind:value={newTag}
						type="text"
						placeholder="Add tag..."
						class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						on:keydown={(e) => e.key === 'Enter' && addTag()}
					/>
					<button
						type="button"
						on:click={addTag}
						class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
					>
						Add
					</button>
				</div>
			</div>

			<!-- Subtasks -->
			<div>
				<label class="block text-sm font-medium text-gray-700 mb-2">
					Subtasks
				</label>

				<!-- Existing Subtasks -->
				{#if formData.subtasks.length > 0}
					<div class="space-y-2 mb-3">
						{#each formData.subtasks as subtask}
							<div class="flex items-center gap-2 p-2 border border-gray-200 rounded-md">
								<input
									type="checkbox"
									checked={subtask.completed}
									on:change={() => toggleSubtask(subtask.id)}
									class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
								/>
								<span class="flex-1 {subtask.completed ? 'line-through text-gray-500' : ''}">
									{subtask.title}
								</span>
								<button
									type="button"
									on:click={() => removeSubtask(subtask.id)}
									class="text-red-500 hover:text-red-700 p-1"
								>
									×
								</button>
							</div>
						{/each}
					</div>
				{/if}

				<!-- Add New Subtask -->
				<div class="flex gap-2">
					<input
						bind:value={newSubtask}
						type="text"
						placeholder="Add subtask..."
						class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						on:keydown={(e) => e.key === 'Enter' && addSubtask()}
					/>
					<button
						type="button"
						on:click={addSubtask}
						class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
					>
						Add
					</button>
				</div>
			</div>

			<!-- AI Suggestions (placeholder) -->
			{#if formData.title && !task}
				<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
					<h4 class="text-sm font-medium text-blue-900 mb-2">💡 AI Suggestions</h4>
					<div class="space-y-2 text-sm text-blue-800">
						<p><strong>Recommended priority:</strong> Based on your description, this seems like a {formData.priority} priority task.</p>
						{#if formData.description && formData.description.length > 50}
							<p><strong>Suggested subtasks:</strong> Consider breaking this into smaller tasks for better tracking.</p>
						{/if}
						{#if !formData.dueDate}
							<p><strong>Due date:</strong> Consider setting a realistic deadline to stay on track.</p>
						{/if}
					</div>
				</div>
			{/if}

			<!-- Form Actions -->
			<div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
				<button
					type="button"
					on:click={cancel}
					class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
				>
					{task ? 'Update Task' : 'Create Task'}
				</button>
			</div>
		</form>
	</div>
</div>
