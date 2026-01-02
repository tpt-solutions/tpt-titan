<script>
	import { emailSearchQuery } from '$lib/stores';

	export let emailsList = [];
	export let currentFolderValue = 'inbox';
	export let selectedEmailData = null;
	export let handleEmailSelect;
	export let handleFolderChange;

	let searchQuery = '';
	let filteredEmails = [];

	emailSearchQuery.subscribe(value => {
		searchQuery = value;
		filterEmails();
	});

	// Update filtered emails when emailsList changes
	$: if (emailsList) {
		filterEmails();
	}

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
</script>

<div class="flex flex-col h-full">
	<!-- Search Bar -->
	<div class="p-4 border-b border-gray-200 dark:border-gray-700">
		<div class="relative">
			<input
				bind:value={searchQuery}
				on:input={() => emailSearchQuery.set(searchQuery)}
				type="text"
				placeholder="Search emails..."
				class="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
			>
			<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
				<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
				</svg>
			</div>
		</div>
	</div>

	<!-- Folders -->
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
		</div>
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
					<div
						on:click={() => handleEmailSelect(email)}
						class="p-4 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors {selectedEmailData?.id === email.id ? 'bg-blue-50 dark:bg-blue-900 border-r-4 border-blue-500' : ''}"
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
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
