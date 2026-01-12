<script>
	import { createEventDispatcher } from 'svelte';
	import SpeechService from '../services/speech.js';

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

	// Speech functionality
	let isRecording = false;
	let availableSTTModels = [];
	let selectedSTTModel = null;
	let isListeningForCommands = false;

	// AI Features
	let taskSuggestions = [];
	let isGeneratingSuggestions = false;
	let priorityPrediction = null;
	let isPredictingPriority = false;
	let deadlinePrediction = null;
	let isPredictingDeadline = false;

	// Initialize speech functionality
	onMount(async () => {
		try {
			// Load available STT models for voice input
			availableSTTModels = await SpeechService.getAvailableModels('stt');

			// Set default STT model
			if (availableSTTModels.length > 0) {
				selectedSTTModel = availableSTTModels[0];
			}
		} catch (error) {
			console.error('Failed to initialize speech for tasks:', error);
		}
	});

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

	// Voice input for title field
	async function startVoiceInput() {
		if (!selectedSTTModel || isRecording) return;

		try {
			isRecording = true;

			// Start recording
			SpeechService.startRecording(
				async (audioBlob) => {
					try {
						// Convert to WAV if needed
						let processedAudio = audioBlob;
						if (audioBlob.type !== 'audio/wav') {
							processedAudio = await SpeechService.convertToWav(audioBlob);
						}

						// Send to speech-to-text API
						const result = await SpeechService.speechToText(
							await processedAudio.arrayBuffer(),
							selectedSTTModel.id,
							{ language: 'en' }
						);

						// Set the transcribed text as title
						if (result && result.output_text) {
							formData.title = result.output_text;
							formData = { ...formData }; // Trigger reactivity
						}
					} catch (error) {
						console.error('Voice input error:', error);
						alert('Failed to transcribe voice input: ' + error.message);
					} finally {
						isRecording = false;
					}
				},
				(error) => {
					console.error('Recording error:', error);
					alert('Voice recording failed: ' + error.message);
					isRecording = false;
				}
			);

			// Auto-stop after 5 seconds
			setTimeout(() => {
				if (isRecording) {
					// Note: In a real implementation, you'd need to stop the MediaRecorder
					isRecording = false;
				}
			}, 5000);

		} catch (error) {
			console.error('Voice input setup error:', error);
			alert('Failed to start voice input: ' + error.message);
			isRecording = false;
		}
	}

	// Voice commands for task creation
	async function startVoiceCommands() {
		if (!selectedSTTModel || isListeningForCommands) return;

		try {
			isListeningForCommands = true;

			SpeechService.startRecording(
				async (audioBlob) => {
					try {
						// Convert to WAV if needed
						let processedAudio = audioBlob;
						if (audioBlob.type !== 'audio/wav') {
							processedAudio = await SpeechService.convertToWav(audioBlob);
						}

						// Send to speech-to-text API
						const result = await SpeechService.speechToText(
							await processedAudio.arrayBuffer(),
							selectedSTTModel.id,
							{ language: 'en' }
						);

						// Process voice commands
						if (result && result.output_text) {
							processVoiceCommand(result.output_text.toLowerCase());
						}
					} catch (error) {
						console.error('Voice command error:', error);
						alert('Failed to process voice command: ' + error.message);
					} finally {
						isListeningForCommands = false;
					}
				},
				(error) => {
					console.error('Voice command recording error:', error);
					alert('Voice command recording failed: ' + error.message);
					isListeningForCommands = false;
				}
			);

			// Auto-stop after 8 seconds for commands
			setTimeout(() => {
				if (isListeningForCommands) {
					isListeningForCommands = false;
				}
			}, 8000);

		} catch (error) {
			console.error('Voice commands setup error:', error);
			alert('Failed to start voice commands: ' + error.message);
			isListeningForCommands = false;
		}
	}

	function processVoiceCommand(command) {
		// Parse natural language commands for task creation
		const patterns = {
			// Priority commands
			'make this high priority': () => formData.priority = 'high',
			'set priority to high': () => formData.priority = 'high',
			'make this medium priority': () => formData.priority = 'medium',
			'set priority to medium': () => formData.priority = 'medium',
			'make this low priority': () => formData.priority = 'low',
			'set priority to low': () => formData.priority = 'low',

			// Due date commands
			'due today': () => formData.dueDate = new Date().toISOString().split('T')[0],
			'due tomorrow': () => {
				const tomorrow = new Date();
				tomorrow.setDate(tomorrow.getDate() + 1);
				formData.dueDate = tomorrow.toISOString().split('T')[0];
			},
			'due next week': () => {
				const nextWeek = new Date();
				nextWeek.setDate(nextWeek.getDate() + 7);
				formData.dueDate = nextWeek.toISOString().split('T')[0];
			},

			// Tag commands
			'add tag (.+)': (match) => {
				const tag = match[1].trim();
				if (tag && !formData.tags.includes(tag)) {
					formData.tags = [...formData.tags, tag];
				}
			},
			'tag this as (.+)': (match) => {
				const tag = match[1].trim();
				if (tag && !formData.tags.includes(tag)) {
					formData.tags = [...formData.tags, tag];
				}
			},

			// Subtask commands
			'add subtask (.+)': (match) => {
				const subtaskTitle = match[1].trim();
				if (subtaskTitle) {
					addSubtaskFromVoice(subtaskTitle);
				}
			},
			'break this into (.+)': (match) => {
				const subtasks = match[1].trim().split(' and ');
				subtasks.forEach(title => {
					if (title.trim()) {
						addSubtaskFromVoice(title.trim());
					}
				});
			}
		};

		let commandProcessed = false;

		for (const [pattern, action] of Object.entries(patterns)) {
			const regex = new RegExp(pattern, 'i');
			const match = command.match(regex);

			if (match) {
				if (typeof action === 'function') {
					action(match);
				}
				commandProcessed = true;
				break;
			}
		}

		if (commandProcessed) {
			formData = { ...formData }; // Trigger reactivity
			alert('Voice command processed successfully!');
		} else {
			alert('Voice command not recognized. Try saying things like "set priority to high" or "due tomorrow"');
		}
	}

	function addSubtaskFromVoice(title) {
		const newSubtask = {
			id: Date.now() + Math.random(),
			title: title,
			completed: false
		};
		formData.subtasks = [...formData.subtasks, newSubtask];
	}

	// AI Functions
	async function predictPriority() {
		if (!formData.title.trim()) {
			alert('Please enter a task title first');
			return;
		}

		isPredictingPriority = true;

		try {
			const response = await fetch('/api/v1/tasks/predict-priority', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					title: formData.title,
					description: formData.description,
					tags: formData.tags
				})
			});

			if (response.ok) {
				const result = await response.json();
				priorityPrediction = result.prediction;
				formData.priority = result.predicted_priority;
				formData = { ...formData }; // Trigger reactivity
			} else {
				console.warn('Failed to predict priority');
			}
		} catch (error) {
			console.error('Error predicting priority:', error);
		} finally {
			isPredictingPriority = false;
		}
	}

	async function predictDeadline() {
		if (!formData.title.trim()) {
			alert('Please enter a task title first');
			return;
		}

		isPredictingDeadline = true;

		try {
			const response = await fetch('/api/v1/tasks/predict-deadline', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					title: formData.title,
					description: formData.description,
					priority: formData.priority,
					tags: formData.tags
				})
			});

			if (response.ok) {
				const result = await response.json();
				deadlinePrediction = result.prediction;
				formData.dueDate = result.predicted_date;
				formData = { ...formData }; // Trigger reactivity
			} else {
				console.warn('Failed to predict deadline');
			}
		} catch (error) {
			console.error('Error predicting deadline:', error);
		} finally {
			isPredictingDeadline = false;
		}
	}

	async function generateTaskSuggestions() {
		if (!formData.title.trim()) {
			alert('Please enter a task title first');
			return;
		}

		isGeneratingSuggestions = true;

		try {
			const response = await fetch('/api/v1/tasks/generate-suggestions', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					title: formData.title,
					description: formData.description,
					context: 'email_form' // Could be 'email', 'form', 'manual'
				})
			});

			if (response.ok) {
				const result = await response.json();
				taskSuggestions = result.suggestions || [];
			} else {
				console.warn('Failed to generate task suggestions');
			}
		} catch (error) {
			console.error('Error generating task suggestions:', error);
		} finally {
			isGeneratingSuggestions = false;
		}
	}

	function applyTaskSuggestion(suggestion) {
		if (suggestion.title) formData.title = suggestion.title;
		if (suggestion.description) formData.description = suggestion.description;
		if (suggestion.priority) formData.priority = suggestion.priority;
		if (suggestion.due_date) formData.dueDate = suggestion.due_date;
		if (suggestion.tags) formData.tags = [...new Set([...formData.tags, ...suggestion.tags])];
		if (suggestion.subtasks) {
			const newSubtasks = suggestion.subtasks.map(title => ({
				id: Date.now() + Math.random(),
				title: title,
				completed: false
			}));
			formData.subtasks = [...formData.subtasks, ...newSubtasks];
		}

		formData = { ...formData }; // Trigger reactivity
		taskSuggestions = [];
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
					<div class="flex gap-2">
						<input
							id="title"
							bind:value={formData.title}
							type="text"
							placeholder="Enter task title..."
							class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							required
						/>
						{#if availableSTTModels.length > 0}
							<button
								type="button"
								on:click={startVoiceInput}
								disabled={isRecording}
								class="px-3 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 disabled:opacity-50 transition-colors"
								title="Voice input for title"
							>
								{isRecording ? '🎙️' : '🎤'}
							</button>
						{/if}
					</div>
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

			<!-- Voice Commands Section -->
			{#if availableSTTModels.length > 0 && !task}
				<div class="bg-purple-50 border border-purple-200 rounded-lg p-4">
					<h4 class="text-sm font-medium text-purple-900 mb-3">🎙️ Voice Commands</h4>
					<p class="text-sm text-purple-800 mb-3">
						Use voice commands to quickly configure your task. Say things like "set priority to high", "due tomorrow", or "add tag urgent".
					</p>
					<div class="flex gap-2">
						<button
							type="button"
							on:click={startVoiceCommands}
							disabled={isListeningForCommands}
							class="px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 disabled:opacity-50 transition-colors text-sm"
						>
							{isListeningForCommands ? '🎙️ Listening...' : '🎤 Voice Commands'}
						</button>
					</div>
					<div class="mt-3 text-xs text-purple-700">
						<strong>Supported commands:</strong> priority settings, due dates, tags, subtasks
					</div>
				</div>
			{/if}

			<!-- AI Features Section -->
			{#if !task && formData.title}
				<div class="border-t border-gray-200 pt-4">
					<h4 class="text-sm font-medium text-gray-900 mb-3">🤖 AI Assistant</h4>

					<div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-4">
						<button
							type="button"
							on:click={predictPriority}
							disabled={isPredictingPriority}
							class="flex items-center justify-center px-3 py-2 border border-red-300 rounded-md hover:bg-red-50 transition-colors disabled:opacity-50"
						>
							<div class="flex items-center space-x-2">
								<span class="text-sm">🎯</span>
								<span class="text-xs font-medium text-red-700">
									{isPredictingPriority ? 'Predicting...' : 'Predict Priority'}
								</span>
							</div>
						</button>

						<button
							type="button"
							on:click={predictDeadline}
							disabled={isPredictingDeadline}
							class="flex items-center justify-center px-3 py-2 border border-green-300 rounded-md hover:bg-green-50 transition-colors disabled:opacity-50"
						>
							<div class="flex items-center space-x-2">
								<span class="text-sm">📅</span>
								<span class="text-xs font-medium text-green-700">
									{isPredictingDeadline ? 'Predicting...' : 'Predict Deadline'}
								</span>
							</div>
						</button>

						<button
							type="button"
							on:click={generateTaskSuggestions}
							disabled={isGeneratingSuggestions}
							class="flex items-center justify-center px-3 py-2 border border-blue-300 rounded-md hover:bg-blue-50 transition-colors disabled:opacity-50"
						>
							<div class="flex items-center space-x-2">
								<span class="text-sm">💡</span>
								<span class="text-xs font-medium text-blue-700">
									{isGeneratingSuggestions ? 'Generating...' : 'Smart Suggestions'}
								</span>
							</div>
						</button>
					</div>

					<!-- AI Predictions Display -->
					{#if priorityPrediction}
						<div class="mb-3 p-3 bg-red-50 border border-red-200 rounded-md">
							<h5 class="text-sm font-medium text-red-900 mb-1">🎯 Priority Prediction</h5>
							<p class="text-sm text-red-800">{priorityPrediction}</p>
						</div>
					{/if}

					{#if deadlinePrediction}
						<div class="mb-3 p-3 bg-green-50 border border-green-200 rounded-md">
							<h5 class="text-sm font-medium text-green-900 mb-1">📅 Deadline Prediction</h5>
							<p class="text-sm text-green-800">{deadlinePrediction}</p>
						</div>
					{/if}

					<!-- Task Suggestions -->
					{#if taskSuggestions.length > 0}
						<div class="space-y-3">
							<h5 class="text-sm font-medium text-gray-900">💡 Smart Suggestions</h5>
							{#each taskSuggestions as suggestion}
								<div class="p-3 bg-blue-50 border border-blue-200 rounded-md">
									<div class="flex items-start justify-between mb-2">
										<h6 class="text-sm font-medium text-blue-900">{suggestion.title || 'Suggested Task'}</h6>
										<button
											type="button"
											on:click={() => applyTaskSuggestion(suggestion)}
											class="px-2 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700"
										>
											Apply
										</button>
									</div>

									{#if suggestion.description}
										<p class="text-xs text-blue-800 mb-2">{suggestion.description}</p>
									{/if}

									<div class="flex flex-wrap gap-1 text-xs">
										{#if suggestion.priority}
											<span class="px-1 py-0.5 bg-red-100 text-red-700 rounded">Priority: {suggestion.priority}</span>
										{/if}
										{#if suggestion.due_date}
											<span class="px-1 py-0.5 bg-green-100 text-green-700 rounded">Due: {suggestion.due_date}</span>
										{/if}
										{#if suggestion.tags && suggestion.tags.length > 0}
											{#each suggestion.tags as tag}
												<span class="px-1 py-0.5 bg-purple-100 text-purple-700 rounded">#{tag}</span>
											{/each}
										{/if}
									</div>

									{#if suggestion.subtasks && suggestion.subtasks.length > 0}
										<div class="mt-2">
											<p class="text-xs text-blue-700 font-medium mb-1">Suggested subtasks:</p>
											<ul class="text-xs text-blue-800 space-y-0.5">
												{#each suggestion.subtasks as subtask}
													<li>• {subtask}</li>
												{/each}
											</ul>
										</div>
									{/if}
								</div>
							{/each}
						</div>
					{/if}
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
