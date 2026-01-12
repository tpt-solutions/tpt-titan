<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { contacts } from '$lib/stores';
	import SpeechService from '../services/speech.js';

	export let event = null; // null for create, object for edit
	export let selectedDate = null; // Pre-selected date for new events
	export let calendars = [];
	export let onClose;

	const dispatch = createEventDispatcher();

	let formData = {
		calendar_id: '',
		title: '',
		description: '',
		location: '',
		start_time: '',
		end_time: '',
		is_all_day: false,
		reminder_minutes: 15,
		attendee_ids: []
	};

	let isSubmitting = false;
	let errors = {};
	let availableContacts = [];

	// AI Features
	let showVoiceEventCreation = false;
	let voiceCreating = false;
	let voiceEventText = '';
	let availableSTTModels = [];
	let selectedSTTModel = null;
	let schedulingSuggestions = [];
	let isGeneratingSuggestions = false;

	contacts.subscribe(value => availableContacts = value);

	// Initialize form data
	if (event) {
		// Editing existing event
		formData = {
			calendar_id: event.calendar_id,
			title: event.title || '',
			description: event.description || '',
			location: event.location || '',
			start_time: event.start_time ? new Date(event.start_time).toISOString().slice(0, 16) : '',
			end_time: event.end_time ? new Date(event.end_time).toISOString().slice(0, 16) : '',
			is_all_day: event.is_all_day || false,
			reminder_minutes: event.reminder_minutes || 15,
			attendee_ids: event.attendees ? event.attendees.map(a => a.contact_id) : []
		};
	} else {
		// Creating new event
		const defaultCalendar = calendars.find(c => c.is_default) || calendars[0];
		const now = selectedDate ? new Date(selectedDate) : new Date();

		// Set default start time to next hour
		const startTime = new Date(now);
		startTime.setHours(startTime.getHours() + 1, 0, 0, 0);

		// Set default end time to 1 hour later
		const endTime = new Date(startTime);
		endTime.setHours(endTime.getHours() + 1);

		formData = {
			calendar_id: defaultCalendar ? defaultCalendar.id : '',
			title: '',
			description: '',
			location: '',
			start_time: startTime.toISOString().slice(0, 16),
			end_time: endTime.toISOString().slice(0, 16),
			is_all_day: false,
			reminder_minutes: 15,
			attendee_ids: []
		};
	}

	function validateForm() {
		errors = {};

		if (!formData.title.trim()) {
			errors.title = 'Event title is required';
		}

		if (!formData.calendar_id) {
			errors.calendar_id = 'Please select a calendar';
		}

		if (!formData.start_time) {
			errors.start_time = 'Start time is required';
		}

		if (!formData.end_time) {
			errors.end_time = 'End time is required';
		}

		if (formData.start_time && formData.end_time) {
			const start = new Date(formData.start_time);
			const end = new Date(formData.end_time);
			if (end <= start) {
				errors.end_time = 'End time must be after start time';
			}
		}

		return Object.keys(errors).length === 0;
	}

	async function handleSubmit() {
		if (!validateForm()) {
			return;
		}

		isSubmitting = true;

		try {
			// Convert datetime-local values to ISO strings
			const eventData = {
				...formData,
				start_time: new Date(formData.start_time).toISOString(),
				end_time: new Date(formData.end_time).toISOString(),
			};

			const method = event ? 'PUT' : 'POST';
			const url = event ? `/api/v1/events/${event.id}` : '/api/v1/events';

			const response = await fetch(url, {
				method,
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify(eventData)
			});

			if (response.ok) {
				dispatch('success');
				handleClose();
			} else {
				const errorData = await response.json();
				errors.general = errorData.error || 'Failed to save event';
			}
		} catch (error) {
			console.error('Failed to save event:', error);
			errors.general = 'Network error. Please try again.';
		} finally {
			isSubmitting = false;
		}
	}

	function handleClose() {
		formData = {
			calendar_id: '',
			title: '',
			description: '',
			location: '',
			start_time: '',
			end_time: '',
			is_all_day: false,
			reminder_minutes: 15,
			attendee_ids: []
		};
		errors = {};
		onClose();
	}

	function handleKeydown(event) {
		if (event.key === 'Escape') {
			handleClose();
		}
	}

	function toggleAttendee(contactId) {
		const index = formData.attendee_ids.indexOf(contactId);
		if (index > -1) {
			formData.attendee_ids.splice(index, 1);
		} else {
			formData.attendee_ids.push(contactId);
		}
		formData.attendee_ids = [...formData.attendee_ids]; // Trigger reactivity
	}

	$: selectedCalendar = calendars.find(c => c.id === formData.calendar_id);

	// Initialize AI features
	onMount(async () => {
		try {
			// Load available STT models for voice event creation
			availableSTTModels = await SpeechService.getAvailableModels('stt');

			// Set default STT model
			if (availableSTTModels.length > 0) {
				selectedSTTModel = availableSTTModels[0];
			}
		} catch (error) {
			console.error('Failed to initialize AI features for calendar:', error);
		}
	});

	// AI Functions
	async function startVoiceEventCreation() {
		if (!selectedSTTModel || voiceCreating) return;

		showVoiceEventCreation = true;
		voiceCreating = true;

		try {
			SpeechService.startRecording(
				async (audioBlob) => {
					try {
						let processedAudio = audioBlob;
						if (audioBlob.type !== 'audio/wav') {
							processedAudio = await SpeechService.convertToWav(audioBlob);
						}

						const result = await SpeechService.speechToText(
							await processedAudio.arrayBuffer(),
							selectedSTTModel.id,
							{ language: 'en' }
						);

						if (result && result.output_text) {
							voiceEventText = result.output_text;
							// Parse the voice input to create event data
							parseVoiceEventInput(result.output_text);
						}
					} catch (error) {
						console.error('Voice event creation error:', error);
						alert('Failed to transcribe voice input: ' + error.message);
					} finally {
						voiceCreating = false;
					}
				},
				(error) => {
					console.error('Voice event creation recording error:', error);
					alert('Voice recording failed: ' + error.message);
					voiceCreating = false;
				}
			);

			// Auto-stop after 45 seconds for event creation (longer for detailed input)
			setTimeout(() => {
				if (voiceCreating) {
					voiceCreating = false;
				}
			}, 45000);

		} catch (error) {
			console.error('Voice event creation setup error:', error);
			alert('Failed to start voice event creation: ' + error.message);
			voiceCreating = false;
		}
	}

	function parseVoiceEventInput(text) {
		// Simple natural language parsing for event creation
		// This would be enhanced with AI in a full implementation
		const lowerText = text.toLowerCase();

		// Extract event title (first sentence or key phrase)
		const sentences = text.split(/[.!?]+/).filter(s => s.trim());
		if (sentences.length > 0) {
			formData.title = sentences[0].trim();
		}

		// Extract date/time information
		const now = new Date();

		// Today
		if (lowerText.includes('today')) {
			const startTime = new Date(now);
			startTime.setHours(startTime.getHours() + 1, 0, 0, 0);
			formData.start_time = startTime.toISOString().slice(0, 16);

			const endTime = new Date(startTime);
			endTime.setHours(endTime.getHours() + 1);
			formData.end_time = endTime.toISOString().slice(0, 16);
		}
		// Tomorrow
		else if (lowerText.includes('tomorrow')) {
			const tomorrow = new Date(now);
			tomorrow.setDate(tomorrow.getDate() + 1);
			tomorrow.setHours(9, 0, 0, 0); // Default 9 AM
			formData.start_time = tomorrow.toISOString().slice(0, 16);

			const endTime = new Date(tomorrow);
			endTime.setHours(endTime.getHours() + 1);
			formData.end_time = endTime.toISOString().slice(0, 16);
		}
		// Next week
		else if (lowerText.includes('next week')) {
			const nextWeek = new Date(now);
			nextWeek.setDate(nextWeek.getDate() + 7);
			nextWeek.setHours(9, 0, 0, 0);
			formData.start_time = nextWeek.toISOString().slice(0, 16);

			const endTime = new Date(nextWeek);
			endTime.setHours(endTime.getHours() + 1);
			formData.end_time = endTime.toISOString().slice(0, 16);
		}

		// Extract duration
		if (lowerText.includes('30 minutes') || lowerText.includes('half hour')) {
			const start = new Date(formData.start_time);
			const end = new Date(start);
			end.setMinutes(end.getMinutes() + 30);
			formData.end_time = end.toISOString().slice(0, 16);
		} else if (lowerText.includes('2 hours') || lowerText.includes('two hours')) {
			const start = new Date(formData.start_time);
			const end = new Date(start);
			end.setHours(end.getHours() + 2);
			formData.end_time = end.toISOString().slice(0, 16);
		}

		// Extract location
		const locationPatterns = [
			/at (.+?)(?:\s|$)/i,
			/in (.+?)(?:\s|$)/i,
			/location (.+?)(?:\s|$)/i
		];

		for (const pattern of locationPatterns) {
			const match = text.match(pattern);
			if (match && match[1]) {
				formData.location = match[1].trim();
				break;
			}
		}

		// Extract description (remaining text)
		if (sentences.length > 1) {
			formData.description = sentences.slice(1).join('. ').trim();
		}

		// Trigger reactivity
		formData = { ...formData };
	}

	async function getSchedulingSuggestions() {
		isGeneratingSuggestions = true;

		try {
			const response = await fetch('/api/v1/calendar/scheduling-suggestions', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					current_events: [], // Would pass current calendar events
					preferences: {
						work_hours_start: 9,
						work_hours_end: 17,
						preferred_days: ['monday', 'tuesday', 'wednesday', 'thursday', 'friday']
					}
				})
			});

			if (response.ok) {
				const result = await response.json();
				schedulingSuggestions = result.suggestions || [];
			} else {
				console.warn('Failed to get scheduling suggestions');
			}
		} catch (error) {
			console.error('Error getting scheduling suggestions:', error);
		} finally {
			isGeneratingSuggestions = false;
		}
	}

	function applySchedulingSuggestion(suggestion) {
		formData.start_time = suggestion.start_time;
		formData.end_time = suggestion.end_time;
		formData = { ...formData };
		schedulingSuggestions = [];
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal backdrop -->
<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" on:click={handleClose}>
	<!-- Modal container -->
	<div class="relative top-20 mx-auto p-5 border w-full max-w-2xl shadow-lg rounded-md bg-white dark:bg-gray-800" on:click|stopPropagation>
		<div class="mt-3 max-h-[80vh] overflow-y-auto">
			<!-- Header -->
			<div class="flex items-center justify-between mb-4">
				<h3 class="text-lg font-medium text-gray-900 dark:text-white">
					{event ? 'Edit Event' : 'Create New Event'}
				</h3>
				<button
					on:click={handleClose}
					class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
				>
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
					</svg>
				</button>
			</div>

			<!-- Form -->
			<form on:submit|preventDefault={handleSubmit} class="space-y-4">
				<!-- General Error -->
				{#if errors.general}
					<div class="bg-red-50 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-600 dark:text-red-200 px-4 py-3 rounded">
						{errors.general}
					</div>
				{/if}

				<!-- Calendar Selection -->
				<div>
					<label for="calendar_id" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Calendar
					</label>
					<select
						id="calendar_id"
						bind:value={formData.calendar_id}
						class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
					>
						<option value="">Select a calendar</option>
						{#each calendars as calendar}
							<option value={calendar.id}>{calendar.name}</option>
						{/each}
					</select>
					{#if errors.calendar_id}
						<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.calendar_id}</p>
					{/if}
				</div>

				<!-- Title -->
				<div>
					<label for="title" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Event Title *
					</label>
					<input
						id="title"
						type="text"
						bind:value={formData.title}
						class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
						placeholder="Meeting with team"
					>
					{#if errors.title}
						<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.title}</p>
					{/if}
				</div>

				<!-- Description -->
				<div>
					<label for="description" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Description
					</label>
					<textarea
						id="description"
						bind:value={formData.description}
						rows="3"
						class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
						placeholder="Event details..."
					></textarea>
				</div>

				<!-- Location -->
				<div>
					<label for="location" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Location
					</label>
					<input
						id="location"
						type="text"
						bind:value={formData.location}
						class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
						placeholder="Conference Room A"
					>
				</div>

				<!-- All Day Toggle -->
				<div class="flex items-center">
					<input
						id="is_all_day"
						type="checkbox"
						bind:checked={formData.is_all_day}
						class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
					>
					<label for="is_all_day" class="ml-2 block text-sm text-gray-900 dark:text-white">
						All day event
					</label>
				</div>

				<!-- Date/Time Fields -->
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<label for="start_time" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Start Time *
						</label>
						<input
							id="start_time"
							type="datetime-local"
							bind:value={formData.start_time}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
						>
						{#if errors.start_time}
							<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.start_time}</p>
						{/if}
					</div>
					<div>
						<label for="end_time" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							End Time *
						</label>
						<input
							id="end_time"
							type="datetime-local"
							bind:value={formData.end_time}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
						>
						{#if errors.end_time}
							<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.end_time}</p>
						{/if}
					</div>
				</div>

				<!-- Reminder -->
				<div>
					<label for="reminder_minutes" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Reminder
					</label>
					<select
						id="reminder_minutes"
						bind:value={formData.reminder_minutes}
						class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
					>
						<option value={0}>No reminder</option>
						<option value={5}>5 minutes before</option>
						<option value={15}>15 minutes before</option>
						<option value={30}>30 minutes before</option>
						<option value={60}>1 hour before</option>
						<option value={1440}>1 day before</option>
					</select>
				</div>

				<!-- Attendees (if contacts available) -->
				{#if availableContacts.length > 0}
					<div>
						<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
							Attendees
						</label>
						<div class="max-h-40 overflow-y-auto border border-gray-300 dark:border-gray-600 rounded-md p-2 space-y-2">
							{#each availableContacts as contact}
								<label class="flex items-center space-x-2 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 p-1 rounded">
									<input
										type="checkbox"
										checked={formData.attendee_ids.includes(contact.id)}
										on:change={() => toggleAttendee(contact.id)}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
									<div class="flex-1 min-w-0">
										<div class="text-sm font-medium text-gray-900 dark:text-white truncate">
											{contact.first_name || contact.last_name ? `${contact.first_name || ''} ${contact.last_name || ''}`.trim() : 'Unnamed Contact'}
										</div>
										{#if contact.email}
											<div class="text-xs text-gray-500 dark:text-gray-400 truncate">
												{contact.email}
											</div>
										{/if}
									</div>
								</label>
							{/each}
						</div>
					</div>
				{/if}

				<!-- AI Features Section -->
				{#if !event}
					<div class="border-t border-gray-200 dark:border-gray-600 pt-4">
						<h4 class="text-sm font-medium text-gray-900 dark:text-white mb-3">🤖 AI Assistant</h4>

						<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
							{#if availableSTTModels.length > 0}
								<button
									type="button"
									on:click={startVoiceEventCreation}
									disabled={voiceCreating}
									class="flex items-center justify-center px-4 py-3 border border-purple-300 dark:border-purple-600 rounded-md hover:bg-purple-50 dark:hover:bg-purple-900 transition-colors disabled:opacity-50"
								>
									<div class="flex items-center space-x-2">
										<svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"/>
										</svg>
										<span class="text-sm font-medium text-purple-700 dark:text-purple-300">
											{voiceCreating ? 'Listening...' : 'Voice Create Event'}
										</span>
									</div>
								</button>
							{/if}

							<button
								type="button"
								on:click={getSchedulingSuggestions}
								disabled={isGeneratingSuggestions}
								class="flex items-center justify-center px-4 py-3 border border-blue-300 dark:border-blue-600 rounded-md hover:bg-blue-50 dark:hover:bg-blue-900 transition-colors disabled:opacity-50"
							>
								<div class="flex items-center space-x-2">
									<svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
									</svg>
									<span class="text-sm font-medium text-blue-700 dark:text-blue-300">
										{isGeneratingSuggestions ? 'Finding...' : 'Smart Scheduling'}
									</span>
								</div>
							</button>
						</div>

						<!-- Voice Event Creation Result -->
						{#if voiceEventText && !voiceCreating}
							<div class="mt-3 p-3 bg-purple-50 dark:bg-purple-900 border border-purple-200 dark:border-purple-700 rounded-md">
								<h5 class="text-sm font-medium text-purple-900 dark:text-purple-100 mb-2">Voice Input Processed:</h5>
								<p class="text-sm text-purple-800 dark:text-purple-200 italic">"{voiceEventText}"</p>
								<p class="text-xs text-purple-600 dark:text-purple-400 mt-1">Event details have been automatically filled based on your voice input.</p>
							</div>
						{/if}

						<!-- Scheduling Suggestions -->
						{#if schedulingSuggestions.length > 0}
							<div class="mt-3 space-y-2">
								<h5 class="text-sm font-medium text-gray-900 dark:text-white">Suggested Times:</h5>
								{#each schedulingSuggestions as suggestion}
									<div class="flex items-center justify-between p-2 bg-blue-50 dark:bg-blue-900 border border-blue-200 dark:border-blue-700 rounded">
										<div>
											<div class="text-sm font-medium text-blue-900 dark:text-blue-100">
												{new Date(suggestion.start_time).toLocaleDateString()} at {new Date(suggestion.start_time).toLocaleTimeString()}
											</div>
											<div class="text-xs text-blue-700 dark:text-blue-300">
												{suggestion.duration} minutes • {suggestion.reason}
											</div>
										</div>
										<button
											type="button"
											on:click={() => applySchedulingSuggestion(suggestion)}
											class="px-3 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700"
										>
											Use This
										</button>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				{/if}

				<!-- Actions -->
				<div class="flex justify-end space-x-3 pt-4">
					<button
						type="button"
						on:click={handleClose}
						class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
						disabled={isSubmitting}
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={isSubmitting}
						class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{#if isSubmitting}
							<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
							</svg>
							Saving...
						{:else}
							{event ? 'Update Event' : 'Create Event'}
						{/if}
					</button>
				</div>
			</form>

<!-- Voice Event Creation Modal -->
{#if showVoiceEventCreation}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">🎤 Voice Event Creation</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => { showVoiceEventCreation = false; voiceEventText = ''; }}
				>
					×
				</button>
			</div>

			<div class="space-y-4">
				<div class="text-center">
					{#if voiceCreating}
						<div class="text-4xl mb-4 animate-pulse">🎙️</div>
						<p class="text-gray-600">Speak your event details. For example:</p>
						<p class="text-sm text-gray-500 italic mt-2">"Schedule a team meeting tomorrow at 2 PM for 1 hour in Conference Room A"</p>
						<p class="text-sm text-gray-500 mt-1">Recording will stop automatically after 45 seconds.</p>
					{:else}
						<div class="text-4xl mb-4 text-green-500">✅</div>
						<p class="text-gray-600">Voice recording completed!</p>
					{/if}
				</div>

				{#if voiceEventText}
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-2">Transcribed Event Details</label>
						<div class="p-3 bg-gray-50 border border-gray-200 rounded-md">
							<p class="text-gray-700 italic">"{voiceEventText}"</p>
						</div>
						<p class="text-xs text-gray-500 mt-2">The event form has been automatically filled based on your voice input.</p>
					</div>

					<div class="flex justify-end space-x-3">
						<button
							on:click={() => { voiceEventText = ''; startVoiceEventCreation(); }}
							class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
						>
							Record Again
						</button>
						<button
							on:click={() => { showVoiceEventCreation = false; }}
							class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
						>
							Continue Editing
						</button>
					</div>
				{:else if !voiceCreating}
					<div class="text-center py-4">
						<p class="text-gray-500">No event details were transcribed. Try recording again.</p>
						<button
							on:click={startVoiceEventCreation}
							class="mt-3 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
						>
							Try Again
						</button>
					</div>
				{/if}

				<div class="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
					<h5 class="text-sm font-medium text-blue-900 mb-2">💡 Voice Commands Examples:</h5>
					<ul class="text-xs text-blue-800 space-y-1">
						<li>• "Schedule a meeting tomorrow at 3 PM for 2 hours"</li>
						<li>• "Team standup today at 9 AM in conference room"</li>
						<li>• "Client call next week at 2 PM for 30 minutes"</li>
						<li>• "Lunch with Sarah tomorrow at noon"</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
{/if}
		</div>
	</div>
</div>
