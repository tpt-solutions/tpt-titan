<script>
	import { contacts } from '$lib/stores';

	export let emailAccounts = [];
	export let onClose;

	let mode = 'compose'; // 'compose' or 'add-account'
	let availableContacts = [];

	contacts.subscribe(value => availableContacts = value);

	// Initialize based on whether we have email accounts
	if (emailAccounts.length === 0) {
		mode = 'add-account';
	}

	// Compose mode state
	let composeForm = {
		account_id: emailAccounts.length > 0 ? emailAccounts[0].id : '',
		subject: '',
		content: '',
		html_content: '',
		recipient_emails: [],
		cc_emails: [],
		bcc_emails: []
	};

	// Add account mode state
	let accountForm = {
		email: '',
		provider: 'imap',
		server: '',
		port: 993,
		username: '',
		password: '',
		use_ssl: true
	};

	let isSubmitting = false;
	let errors = {};
	let recipientInput = '';
	let ccInput = '';
	let bccInput = '';

	// Provider presets
	const providerPresets = {
		gmail: {
			imap: { server: 'imap.gmail.com', port: 993 },
			smtp: { server: 'smtp.gmail.com', port: 587 }
		},
		outlook: {
			imap: { server: 'outlook.office365.com', port: 993 },
			smtp: { server: 'smtp-mail.outlook.com', port: 587 }
		},
		yahoo: {
			imap: { server: 'imap.mail.yahoo.com', port: 993 },
			smtp: { server: 'smtp.mail.yahoo.com', port: 587 }
		}
	};

	function validateComposeForm() {
		errors = {};

		if (!composeForm.account_id) {
			errors.account_id = 'Please select an email account';
		}

		if (!composeForm.content.trim()) {
			errors.content = 'Email content is required';
		}

		if (composeForm.recipient_emails.length === 0) {
			errors.recipients = 'At least one recipient is required';
		}

		return Object.keys(errors).length === 0;
	}

	function validateAccountForm() {
		errors = {};

		if (!accountForm.email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(accountForm.email)) {
			errors.email = 'Please enter a valid email address';
		}

		if (!accountForm.server) {
			errors.server = 'Server address is required';
		}

		if (!accountForm.username) {
			errors.username = 'Username is required';
		}

		if (!accountForm.password) {
			errors.password = 'Password is required';
		}

		return Object.keys(errors).length === 0;
	}

	function updateProviderPresets() {
		const emailDomain = accountForm.email.split('@')[1];
		let preset = null;

		if (emailDomain?.includes('gmail.com')) {
			preset = providerPresets.gmail;
		} else if (emailDomain?.includes('outlook.com') || emailDomain?.includes('office365.com')) {
			preset = providerPresets.outlook;
		} else if (emailDomain?.includes('yahoo.com')) {
			preset = providerPresets.yahoo;
		}

		if (preset) {
			const providerConfig = preset[accountForm.provider] || preset.imap;
			accountForm.server = providerConfig.server;
			accountForm.port = providerConfig.port;
		}
	}

	function addEmailToList(inputValue, listKey) {
		const email = inputValue.trim();
		if (email && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
			if (!composeForm[listKey].includes(email)) {
				composeForm[listKey] = [...composeForm[listKey], email];
			}
			if (listKey === 'recipient_emails') recipientInput = '';
			if (listKey === 'cc_emails') ccInput = '';
			if (listKey === 'bcc_emails') bccInput = '';
		}
	}

	function removeEmailFromList(email, listKey) {
		composeForm[listKey] = composeForm[listKey].filter(e => e !== email);
	}

	function addContactToRecipients(contact) {
		const email = contact.email;
		if (email && !composeForm.recipient_emails.includes(email)) {
			composeForm.recipient_emails = [...composeForm.recipient_emails, email];
		}
	}

	async function handleSubmit() {
		if (mode === 'compose') {
			if (!validateComposeForm()) return;

			isSubmitting = true;
			try {
				const response = await fetch('/api/v1/emails', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${localStorage.getItem('token')}`
					},
					body: JSON.stringify(composeForm)
				});

				if (response.ok) {
					handleClose();
				} else {
					const errorData = await response.json();
					errors.general = errorData.error || 'Failed to send email';
				}
			} catch (error) {
				console.error('Failed to send email:', error);
				errors.general = 'Network error. Please try again.';
			}
		} else {
			if (!validateAccountForm()) return;

			isSubmitting = true;
			try {
				const response = await fetch('/api/v1/email-accounts', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${localStorage.getItem('token')}`
					},
					body: JSON.stringify(accountForm)
				});

				if (response.ok) {
					// Switch to compose mode after adding account
					mode = 'compose';
					composeForm.account_id = (await response.json()).account.id;
					errors = {};
				} else {
					const errorData = await response.json();
					errors.general = errorData.error || 'Failed to add email account';
				}
			} catch (error) {
				console.error('Failed to add email account:', error);
				errors.general = 'Network error. Please try again.';
			}
		}

		isSubmitting = false;
	}

	function handleClose() {
		// Reset forms
		composeForm = {
			account_id: emailAccounts.length > 0 ? emailAccounts[0].id : '',
			subject: '',
			content: '',
			html_content: '',
			recipient_emails: [],
			cc_emails: [],
			bcc_emails: []
		};

		accountForm = {
			email: '',
			provider: 'imap',
			server: '',
			port: 993,
			username: '',
			password: '',
			use_ssl: true
		};

		errors = {};
		recipientInput = '';
		ccInput = '';
		bccInput = '';

		onClose();
	}

	function handleKeydown(event) {
		if (event.key === 'Escape') {
			handleClose();
		}
	}

	// Update presets when email changes
	$: if (accountForm.email) {
		updateProviderPresets();
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal backdrop -->
<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" on:click={handleClose}>
	<!-- Modal container -->
	<div class="relative top-20 mx-auto p-5 border w-full max-w-4xl shadow-lg rounded-md bg-white dark:bg-gray-800" on:click|stopPropagation>
		<div class="mt-3 max-h-[80vh] overflow-y-auto">
			<!-- Header -->
			<div class="flex items-center justify-between mb-4">
				<h3 class="text-lg font-medium text-gray-900 dark:text-white">
					{mode === 'compose' ? 'Compose Email' : 'Add Email Account'}
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

			{#if mode === 'add-account'}
				<!-- Add Email Account Form -->
				<form on:submit|preventDefault={handleSubmit} class="space-y-4">
					<!-- General Error -->
					{#if errors.general}
						<div class="bg-red-50 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-600 dark:text-red-200 px-4 py-3 rounded">
							{errors.general}
						</div>
					{/if}

					<!-- Email -->
					<div>
						<label for="account_email" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Email Address *
						</label>
						<input
							id="account_email"
							type="email"
							bind:value={accountForm.email}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							placeholder="your.email@example.com"
						>
						{#if errors.email}
							<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.email}</p>
						{/if}
					</div>

					<!-- Provider -->
					<div>
						<label for="provider" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Account Type
						</label>
						<select
							id="provider"
							bind:value={accountForm.provider}
							on:change={updateProviderPresets}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
						>
							<option value="imap">IMAP (Recommended)</option>
							<option value="pop3">POP3</option>
						</select>
					</div>

					<!-- Server and Port -->
					<div class="grid grid-cols-3 gap-4">
						<div class="col-span-2">
							<label for="server" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
								Server *
							</label>
							<input
								id="server"
								type="text"
								bind:value={accountForm.server}
								class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
								placeholder="imap.example.com"
							>
							{#if errors.server}
								<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.server}</p>
							{/if}
						</div>
						<div>
							<label for="port" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
								Port *
							</label>
							<input
								id="port"
								type="number"
								bind:value={accountForm.port}
								class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							>
						</div>
					</div>

					<!-- Username -->
					<div>
						<label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Username *
						</label>
						<input
							id="username"
							type="text"
							bind:value={accountForm.username}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							placeholder="your.username"
						>
						{#if errors.username}
							<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.username}</p>
						{/if}
					</div>

					<!-- Password -->
					<div>
						<label for="account_password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Password *
						</label>
						<input
							id="account_password"
							type="password"
							bind:value={accountForm.password}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
						>
						{#if errors.password}
							<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.password}</p>
						{/if}
					</div>

					<!-- SSL Option -->
					<div class="flex items-center">
						<input
							id="use_ssl"
							type="checkbox"
							bind:checked={accountForm.use_ssl}
							class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
						>
						<label for="use_ssl" class="ml-2 block text-sm text-gray-900 dark:text-white">
							Use SSL/TLS encryption
						</label>
					</div>

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
								Adding Account...
							{:else}
								Add Account
							{/if}
						</button>
					</div>
				</form>
			{:else}
				<!-- Compose Email Form -->
				<form on:submit|preventDefault={handleSubmit} class="space-y-4">
					<!-- General Error -->
					{#if errors.general}
						<div class="bg-red-50 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-600 dark:text-red-200 px-4 py-3 rounded">
							{errors.general}
						</div>
					{/if}

					<!-- Account Selection -->
					{#if emailAccounts.length > 1}
						<div>
							<label for="account_id" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
								From
							</label>
							<select
								id="account_id"
								bind:value={composeForm.account_id}
								class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							>
								{#each emailAccounts as account}
									<option value={account.id}>{account.display_name || account.email}</option>
								{/each}
							</select>
							{#if errors.account_id}
								<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.account_id}</p>
							{/if}
						</div>
					{/if}

					<!-- Recipients -->
					<div>
						<label for="recipients" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							To *
						</label>
						<div class="mt-1 space-y-2">
							<!-- Recipient Input -->
							<input
								id="recipients"
								type="email"
								bind:value={recipientInput}
								on:keydown={(e) => e.key === 'Enter' && addEmailToList(recipientInput, 'recipient_emails')}
								class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
								placeholder="recipient@example.com"
							>
							{#if errors.recipients}
								<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.recipients}</p>
							{/if}

							<!-- Recipient Tags -->
							{#if composeForm.recipient_emails.length > 0}
								<div class="flex flex-wrap gap-2">
									{#each composeForm.recipient_emails as email}
										<span class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
											{email}
											<button
												type="button"
												on:click={() => removeEmailFromList(email, 'recipient_emails')}
												class="ml-2 text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
											>
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
												</svg>
											</button>
										</span>
									{/each}
								</div>
							{/if}

							<!-- Contacts Quick Add -->
							{#if availableContacts.length > 0}
								<div class="border border-gray-200 dark:border-gray-600 rounded-md p-2">
									<p class="text-xs text-gray-600 dark:text-gray-400 mb-2">Quick add from contacts:</p>
									<div class="flex flex-wrap gap-1">
										{#each availableContacts.filter(c => c.email) as contact}
											<button
												type="button"
												on:click={() => addContactToRecipients(contact)}
												disabled={composeForm.recipient_emails.includes(contact.email)}
												class="text-xs px-2 py-1 rounded bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
											>
												{contact.first_name || contact.last_name || contact.email}
											</button>
										{/each}
									</div>
								</div>
							{/if}
						</div>
					</div>

					<!-- CC/BCC -->
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="cc" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
								CC
							</label>
							<input
								id="cc"
								type="email"
								bind:value={ccInput}
								on:keydown={(e) => e.key === 'Enter' && addEmailToList(ccInput, 'cc_emails')}
								class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
								placeholder="cc@example.com"
							>
							{#if composeForm.cc_emails.length > 0}
								<div class="flex flex-wrap gap-1 mt-1">
									{#each composeForm.cc_emails as email}
										<span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200">
											{email}
											<button
												type="button"
												on:click={() => removeEmailFromList(email, 'cc_emails')}
												class="ml-1 text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-300"
											>
												<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
												</svg>
											</button>
										</span>
									{/each}
								</div>
							{/if}
						</div>
						<div>
							<label for="bcc" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
								BCC
							</label>
							<input
								id="bcc"
								type="email"
								bind:value={bccInput}
								on:keydown={(e) => e.key === 'Enter' && addEmailToList(bccInput, 'bcc_emails')}
								class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
								placeholder="bcc@example.com"
							>
							{#if composeForm.bcc_emails.length > 0}
								<div class="flex flex-wrap gap-1 mt-1">
									{#each composeForm.bcc_emails as email}
										<span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200">
											{email}
											<button
												type="button"
												on:click={() => removeEmailFromList(email, 'bcc_emails')}
												class="ml-1 text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-300"
											>
												<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
												</svg>
											</button>
										</span>
									{/each}
								</div>
							{/if}
						</div>
					</div>

					<!-- Subject -->
					<div>
						<label for="subject" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Subject
						</label>
						<input
							id="subject"
							type="text"
							bind:value={composeForm.subject}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							placeholder="Email subject"
						>
					</div>

					<!-- Content -->
					<div>
						<label for="content" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Message *
						</label>
						<textarea
							id="content"
							bind:value={composeForm.content}
							rows="8"
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							placeholder="Compose your email..."
						></textarea>
						{#if errors.content}
							<p class="mt-1 text-sm text-red-600 dark:text-red-400">{errors.content}</p>
						{/if}
					</div>

					<!-- Actions -->
					<div class="flex justify-between pt-4">
						<button
							type="button"
							on:click={() => mode = 'add-account'}
							class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
						>
							Add Email Account
						</button>
						<div class="flex space-x-3">
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
									Sending...
								{:else}
									Send Email
								{/if}
							</button>
						</div>
					</div>
				</form>
			{/if}
		</div>
	</div>
</div>
