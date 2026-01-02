<script>
	export let email = null;
	export let onClose;

	async function markAsRead(isRead) {
		if (!email) return;

		try {
			const response = await fetch(`/api/v1/emails/${email.id}/read`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({ is_read: isRead })
			});

			if (response.ok) {
				// Update local email object
				email.is_read = isRead;
			}
		} catch (error) {
			console.error('Failed to mark email as read:', error);
		}
	}

	async function toggleStar() {
		if (!email) return;

		try {
			const response = await fetch(`/api/v1/emails/${email.id}/star`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({ is_starred: !email.is_starred })
			});

			if (response.ok) {
				// Update local email object
				email.is_starred = !email.is_starred;
			}
		} catch (error) {
			console.error('Failed to star email:', error);
		}
	}

	async function moveToFolder(folder) {
		if (!email) return;

		try {
			const response = await fetch(`/api/v1/emails/${email.id}/move`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({ folder })
			});

			if (response.ok) {
				// Update local email object
				email.folder = folder;
				// Close viewer since email moved to different folder
				onClose();
			}
		} catch (error) {
			console.error('Failed to move email:', error);
		}
	}

	function formatDate(dateString) {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatEmailList(emails) {
		if (!emails || emails.length === 0) return '';
		return emails.join(', ');
	}

	// Mark as read when email is viewed
	$: if (email && !email.is_read) {
		markAsRead(true);
	}
</script>

{#if email}
	<div class="flex flex-col h-full bg-white dark:bg-gray-800">
		<!-- Email Header -->
		<div class="border-b border-gray-200 dark:border-gray-700 p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="flex items-center space-x-4">
					<button
						on:click={onClose}
						class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
					>
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
						</svg>
					</button>
					<h2 class="text-xl font-semibold text-gray-900 dark:text-white truncate">
						{email.subject || '(No subject)'}
					</h2>
				</div>

				<div class="flex items-center space-x-2">
					<!-- Star Button -->
					<button
						on:click={toggleStar}
						class="p-2 text-gray-400 hover:text-yellow-500 transition-colors"
						title={email.is_starred ? 'Remove star' : 'Star email'}
					>
						<svg class="w-5 h-5 {email.is_starred ? 'text-yellow-500 fill-current' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"/>
						</svg>
					</button>

					<!-- Move to Folder Dropdown -->
					<div class="relative">
						<select
							on:change={(e) => moveToFolder(e.target.value)}
							class="text-sm border border-gray-300 dark:border-gray-600 rounded px-3 py-1 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
						>
							<option value="">Move to...</option>
							<option value="inbox">Inbox</option>
							<option value="archive">Archive</option>
							<option value="trash">Trash</option>
						</select>
					</div>
				</div>
			</div>

			<!-- Email Metadata -->
			<div class="space-y-2">
				<div class="flex items-center justify-between">
					<div class="flex items-center space-x-4">
						<div class="flex-shrink-0">
							<div class="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold">
								{(email.sender_name || email.sender_email).charAt(0).toUpperCase()}
							</div>
						</div>
						<div>
							<h3 class="text-lg font-medium text-gray-900 dark:text-white">
								{email.sender_name || email.sender_email}
							</h3>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								{email.sender_email}
							</p>
						</div>
					</div>
					<div class="text-sm text-gray-500 dark:text-gray-400">
						{formatDate(email.received_at || email.sent_at)}
					</div>
				</div>

				{#if email.recipient_emails && email.recipient_emails.length > 0}
					<div>
						<span class="text-sm font-medium text-gray-700 dark:text-gray-300">To:</span>
						<span class="text-sm text-gray-900 dark:text-white ml-2">
							{formatEmailList(email.recipient_emails)}
						</span>
					</div>
				{/if}

				{#if email.cc_emails && email.cc_emails.length > 0}
					<div>
						<span class="text-sm font-medium text-gray-700 dark:text-gray-300">CC:</span>
						<span class="text-sm text-gray-900 dark:text-white ml-2">
							{formatEmailList(email.cc_emails)}
						</span>
					</div>
				{/if}
			</div>
		</div>

		<!-- Email Content -->
		<div class="flex-1 overflow-y-auto p-6">
			<div class="max-w-none">
				{#if email.html_content}
					<!-- HTML Content -->
					<div class="prose dark:prose-invert max-w-none">
						{@html email.html_content}
					</div>
				{:else if email.content}
					<!-- Plain Text Content -->
					<div class="whitespace-pre-wrap text-gray-900 dark:text-white font-mono text-sm leading-relaxed">
						{email.content}
					</div>
				{:else}
					<div class="text-gray-500 dark:text-gray-400 italic">
						No content available
					</div>
				{/if}
			</div>
		</div>

		<!-- Attendees (if this is an event-related email) -->
		{#if email.attendees && email.attendees.length > 0}
			<div class="border-t border-gray-200 dark:border-gray-700 p-6">
				<h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Attendees</h4>
				<div class="flex flex-wrap gap-2">
					{#each email.attendees as attendee}
						<div class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
							{attendee.name}
							{#if attendee.status !== 'pending'}
								<span class="ml-2 text-xs">
									({attendee.status})
								</span>
							{/if}
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
{/if}
