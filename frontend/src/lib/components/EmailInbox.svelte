<script>
	import { onMount } from 'svelte';
	import { emailSearchQuery } from '$lib/stores';
	import SpeechService from '../services/speech.js';

	export let emailsList = [];
	export let currentFolderValue = 'inbox';
	export let selectedEmailData = null;
	export let handleEmailSelect;
	export let handleFolderChange;

	let searchQuery = '';
	let filteredEmails = [];

	// AI Features
	let showAICategorization = false;
	let aiCategorizing = false;
	let emailCategories = {};
	let showVoiceCompose = false;
	let voiceComposing = false;
	let voiceComposeText = '';
	let availableSTTModels = [];
	let selectedSTTModel = null;

	emailSearchQuery.subscribe(value => {
		searchQuery = value;
		filterEmails();
	});

	// Update filtered emails when emailsList changes
	$: if (emailsList) {
		filterEmails();
	}

	// Initialize AI features
	onMount(async () => {
		try {
			// Load available STT models for voice composition
			availableSTTModels = await SpeechService.getAvailableModels('stt');

			// Set default STT model
			if (availableSTTModels.length > 0) {
				selectedSTTModel = availableSTTModels[0];
			}
		} catch (error) {
			console.error('Failed to initialize AI features for email:', error);
		}
	});

	function filterEmails() {
		if (!searchQuery.trim()) {
			filteredEmails = emailsList;
		} else {
			const query = searchQuery.toLowerCase();
			filteredEmails = emailsList.filter(email =>
				email.subject?.toLowerCase().includes(query) ||
				email.sender_name?.toLowerCase().includes(query) ||
				email.sender_email?.toLowerCase().includes(query) ||
				email.content?.toLowerCase().includes(query)
			);
		}
	}

	function formatDate(dateString) {
		if (!dateString) return '';

		const date = new Date(dateString);
		const now = new Date();
		const diffTime = Math.abs(now - date);
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

		if (diffDays === 1) {
			return 'Today';
		} else if (diffDays === 2) {
			return 'Yesterday';
		} else if (diffDays <= 7) {
			return date.toLocaleDateString('en-US', { weekday: 'short' });
		} else {
			return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
		}
	}

	function truncateText(text, maxLength = 100) {
		if (!text) return '';
		if (text.length <= maxLength) return text;
		return text.substring(0, maxLength) + '...';
	}

	// Folder statistics
	$: inboxCount = emailsList.filter(email => email.folder === 'inbox' && !email.is_read).length;
	$: starredCount = emailsList.filter(email => email.is_starred).length;
	$: sentCount = emailsList.filter(email => email.folder === 'sent').length;

	// AI Functions
	async function categorizeEmails() {
		if (emailsList.length === 0) return;

		aiCategorizing = true;
		showAICategorization = true;

		try {
			const response = await fetch('/api/v1/emails/categorize', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					email_ids: emailsList.map(email => email.id)
				})
			});

			if (response.ok) {
				const result = await response.json();
				emailCategories = result.categories || {};
			} else {
				console.warn('Failed to categorize emails');
			}
		} catch (error) {
			console.error('Error categorizing emails:', error);
			alert('Failed to categorize emails. Please try again.');
		} finally {
			aiCategorizing = false;
		}
	}

	async function convertEmailToTask(email) {
		try {
			const response = await fetch('/api/v1/emails/convert-to-task', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					email_id: email.id
				})
			});

			if (response.ok) {
				const result = await response.json();
				alert(`Task "${result.task.title}" has been created from this email.`);
			} else {
				console.warn('Failed to convert email to task');
				alert('Failed to convert email to task. Please try again.');
			}
		} catch (error) {
			console.error('Error converting email to task:', error);
			alert('Failed to convert email to task. Please try again.');
		}
	}

	async function startVoiceCompose() {
		if (!selectedSTTModel || voiceComposing) return;

		showVoiceCompose = true;
		voiceComposing = true;

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
							voiceComposeText = result.output_text;
						}
					} catch (error) {
						console.error('Voice compose error:', error);
						alert('Failed to transcribe voice input: ' + error.message);
					} finally {
						voiceComposing = false;
					}
				},
				(error) => {
					console.error('Voice compose recording error:', error);
					alert('Voice recording failed: ' + error.message);
					voiceComposing = false;
				}
			);

			// Auto-stop after 30 seconds for voice composition
			setTimeout(() => {
				if (voiceComposing) {
					voiceComposing = false;
				}
			}, 30000);

		} catch (error) {
			console.error('Voice compose setup error:', error);
			alert('Failed to start voice composition: ' + error.message);
			voiceComposing = false;
		}
	}

	function createEmailFromVoice() {
		// This would open the email composer with the voice text
		// For now, just copy to clipboard
		if (voiceComposeText.trim()) {
			navigator.clipboard.writeText(voiceComposeText);
			alert('Voice text copied to clipboard. You can now paste it into the email composer.');
			showVoiceCompose = false;
			voiceComposeText = '';
		}
	}

	function getCategoryColor(category) {
		const colors = {
			'work': 'bg-blue-100 text-blue-800',
			'personal': 'bg-green-100 text-green-800',
			'promotional': 'bg-yellow-100 text-yellow-800',
			'social': 'bg-purple-100 text-purple-800',
			'finance': 'bg-red-100 text-red-800',
			'news': 'bg-indigo-100 text-indigo-800'
		};
		return colors[category] || 'bg-gray-100 text-gray-800';
	}
</script>

<div class="flex flex-col h-full">
	<!-- Search Bar -->
	<div class="p-4 border-b border-gray-200 dark:border-gray-700">
		<div class="relative">
			<label for="email-search" class="sr-only">Search emails</label>
			<input
				id="email-search"
				bind:value={searchQuery}
				on:input={() => emailSearchQuery.set(searchQuery)}
				type="text"
				placeholder="Search emails..."
				class="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				aria-label="Search emails by subject, sender, or content"
			>
			<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none" aria-hidden="true">
				<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
				</svg>
			</div>
		</div>
	</div>

	<!-- Folders and AI Tools -->
	<div class="border-b border-gray-200 dark:border-gray-700">
		<div class="p-2 space-y-1">
			<button
				on:click={() => handleFolderChange('inbox')}
				class="w-full flex items-center justify-between px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors {currentFolderValue === 'inbox' ? 'bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300' : 'text-gray-700 dark:text-gray-300'}"
			>
				<div class="flex items-center space-x-2">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
					</svg>
					<span>Inbox</span>
				</div>
				{#if inboxCount > 0}
					<span class="bg-blue-600 text-white text-xs px-2 py-1 rounded-full min-w-[20px] text-center">
						{inboxCount}
					</span>
				{/if}
			</button>

			<button
				on:click={() => handleFolderChange('starred')}
				class="w-full flex items-center justify-between px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors {currentFolderValue === 'starred' ? 'bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300' : 'text-gray-700 dark:text-gray-300'}"
			>
				<div class="flex items-center space-x-2">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"/>
					</svg>
					<span>Starred</span>
				</div>
				{#if starredCount > 0}
					<span class="bg-yellow-600 text-white text-xs px-2 py-1 rounded-full min-w-[20px] text-center">
						{starredCount}
					</span>
				{/if}
			</button>

			<button
				on:click={() => handleFolderChange('sent')}
				class="w-full flex items-center justify-between px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors {currentFolderValue === 'sent' ? 'bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300' : 'text-gray-700 dark:text-gray-300'}"
			>
				<div class="flex items-center space-x-2">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
					</svg>
					<span>Sent</span>
				</div>
				{#if sentCount > 0}
					<span class="bg-gray-600 text-white text-xs px-2 py-1 rounded-full min-w-[20px] text-center">
						{sentCount}
					</span>
				{/if}
			</button>

			<!-- AI Tools Section -->
			<div class="pt-2 border-t border-gray-200 dark:border-gray-600">
				<div class="px-3 py-1">
					<span class="text-xs text-gray-500 font-medium">AI TOOLS</span>
				</div>

				<button
					on:click={categorizeEmails}
					disabled={aiCategorizing}
					class="w-full flex items-center px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-300 disabled:opacity-50"
				>
					<div class="flex items-center space-x-2">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
						</svg>
						<span>{aiCategorizing ? 'Categorizing...' : 'Smart Categorize'}</span>
					</div>
				</button>

				{#if availableSTTModels.length > 0}
					<button
						on:click={startVoiceCompose}
						disabled={voiceComposing}
						class="w-full flex items-center px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-300 disabled:opacity-50"
					>
						<div class="flex items-center space-x-2">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"/>
							</svg>
							<span>{voiceComposing ? 'Listening...' : 'Voice Compose'}</span>
						</div>
					</button>
		{/if}
	</div>
</div>

<!-- Voice Compose Modal -->
{#if showVoiceCompose}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">🎤 Voice Email Composition</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => { showVoiceCompose = false; voiceComposeText = ''; }}
				>
					×
				</button>
			</div>

			<div class="space-y-4">
				<div class="text-center">
					{#if voiceComposing}
						<div class="text-4xl mb-4 animate-pulse">🎙️</div>
						<p class="text-gray-600">Listening... Speak your email content.</p>
						<p class="text-sm text-gray-500 mt-2">Recording will stop automatically after 30 seconds.</p>
					{:else}
						<div class="text-4xl mb-4 text-green-500">✅</div>
						<p class="text-gray-600">Voice recording completed!</p>
					{/if}
				</div>

				{#if voiceComposeText}
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-2">Transcribed Text</label>
						<textarea
							bind:value={voiceComposeText}
							rows="6"
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
							placeholder="Your transcribed email text will appear here..."
						></textarea>
					</div>

					<div class="flex justify-end space-x-3">
						<button
							on:click={() => { voiceComposeText = ''; startVoiceCompose(); }}
							class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
						>
							Record Again
						</button>
						<button
							on:click={createEmailFromVoice}
							class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
						>
							Use This Text
						</button>
					</div>
				{:else if !voiceComposing}
					<div class="text-center py-4">
						<p class="text-gray-500">No text was transcribed. Try recording again.</p>
						<button
							on:click={startVoiceCompose}
							class="mt-3 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
						>
							Try Again
						</button>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}

<!-- AI Categorization Results Modal -->
{#if showAICategorization}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">🧠 AI Email Categorization</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => { showAICategorization = false; emailCategories = {}; }}
				>
					×
				</button>
			</div>

			{#if aiCategorizing}
				<div class="text-center py-8">
					<div class="text-4xl mb-4 animate-spin">🧠</div>
					<p class="text-gray-600">Analyzing your emails and categorizing them...</p>
				</div>
			{:else if Object.keys(emailCategories).length === 0}
				<div class="text-center py-8">
					<p class="text-gray-600">No categorization results available.</p>
				</div>
			{:else}
				<div class="space-y-6">
					{#each Object.entries(emailCategories) as [category, emails]}
						<div class="border border-gray-200 rounded-lg p-4">
							<h4 class="font-medium text-gray-900 mb-3 capitalize flex items-center space-x-2">
								<span class="px-2 py-1 rounded text-sm {getCategoryColor(category)}">
									{category}
								</span>
								<span>({emails.length} emails)</span>
							</h4>

							<div class="space-y-2">
								{#each emails as email}
									<div class="flex items-center justify-between p-2 bg-gray-50 rounded">
										<div class="flex-1 min-w-0">
											<h5 class="text-sm font-medium truncate">{email.subject || '(No subject)'}</h5>
											<p class="text-xs text-gray-600 truncate">{email.sender_name || email.sender_email}</p>
										</div>
										<div class="flex items-center space-x-2">
											<button
												on:click={() => convertEmailToTask(email)}
												class="px-2 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700"
												title="Convert to task"
											>
												📋 Task
											</button>
											<button
												on:click={() => handleEmailSelect(email)}
												class="px-2 py-1 text-xs bg-gray-600 text-white rounded hover:bg-gray-700"
												title="View email"
											>
												View
											</button>
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}
	</div>

	<!-- Email List -->
	<div class="flex-1 overflow-y-auto">
		{#if filteredEmails.length === 0}
			<div class="p-8 text-center">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No emails</h3>
				<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
					{searchQuery ? 'No emails match your search.' : `No emails in ${currentFolderValue}.`}
				</p>
			</div>
		{:else}
			<div class="divide-y divide-gray-200 dark:divide-gray-700">
				{#each filteredEmails as email}
					<button
						on:click={() => handleEmailSelect(email)}
						on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); handleEmailSelect(email); } }}
						class="w-full text-left p-4 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors {selectedEmailData?.id === email.id ? 'bg-blue-50 dark:bg-blue-900 border-r-4 border-blue-500' : ''} focus:outline-none focus:ring-2 focus:ring-blue-500"
						aria-label="Read email: {email.subject || '(No subject)'} from {email.sender_name || email.sender_email}{email.is_starred ? ', starred' : ''}{!email.is_read ? ', unread' : ''}"
						type="button"
					>
						<div class="flex items-start justify-between">
							<div class="flex-1 min-w-0">
								<div class="flex items-center space-x-2 mb-1">
									{#if email.is_starred}
										<svg class="w-4 h-4 text-yellow-500 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
											<path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"/>
										</svg>
									{/if}
									<h4 class="text-sm font-medium text-gray-900 dark:text-white truncate">
										{email.sender_name || email.sender_email}
									</h4>
									{#if !email.is_read}
										<span class="inline-block w-2 h-2 bg-blue-600 rounded-full flex-shrink-0"></span>
									{/if}
								</div>

								<h3 class="text-sm font-medium text-gray-900 dark:text-white truncate mb-1">
									{email.subject || '(No subject)'}
								</h3>

								<p class="text-sm text-gray-600 dark:text-gray-400 truncate">
									{truncateText(email.content, 120)}
								</p>
							</div>

							<div class="flex-shrink-0 ml-4 text-xs text-gray-500 dark:text-gray-400">
								{formatDate(email.received_at || email.sent_at)}
							</div>
						</div>
					</button>
				{/each}
			</div>
		{/if}
	</div>
</div>
