<script>
// @ts-nocheck
	import { contacts } from '$lib/stores';

	export let handleEditContact;

	let contactList = [];

	contacts.subscribe(value => {
		contactList = value;
	});

	async function handleDeleteContact(contact) {
		if (!confirm(`Are you sure you want to delete ${contact.first_name || contact.last_name || 'this contact'}?`)) {
			return;
		}

		try {
			const response = await fetch(`/api/v1/contacts/${contact.id}`, {
				method: 'DELETE',
				headers: {
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				}
			});

			if (response.ok) {
				// Remove from store
				contacts.update(current => current.filter(c => c.id !== contact.id));
			} else {
				alert('Failed to delete contact');
			}
		} catch (error) {
			console.error('Failed to delete contact:', error);
			alert('Failed to delete contact');
		}
	}

	function getContactDisplayName(contact) {
		const firstName = contact.first_name || '';
		const lastName = contact.last_name || '';
		const fullName = `${firstName} ${lastName}`.trim();
		return fullName || 'Unnamed Contact';
	}

	function getContactInitials(contact) {
		const firstInitial = contact.first_name ? contact.first_name.charAt(0).toUpperCase() : '';
		const lastInitial = contact.last_name ? contact.last_name.charAt(0).toUpperCase() : '';
		return firstInitial + lastInitial || '?';
	}
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow">
	{#if contactList.length === 0}
		<div class="p-8 text-center">
			<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
			</svg>
			<h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No contacts</h3>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Get started by adding your first contact.</p>
		</div>
	{:else}
		<div class="divide-y divide-gray-200 dark:divide-gray-700">
			{#each contactList as contact}
				<div class="p-6 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
					<div class="flex items-center justify-between">
						<div class="flex items-center space-x-4">
							<!-- Avatar -->
							<div class="flex-shrink-0">
								<div class="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold text-lg">
									{getContactInitials(contact)}
								</div>
							</div>

							<!-- Contact Info -->
							<div class="flex-1 min-w-0">
								<div class="flex items-center space-x-3">
									<h3 class="text-lg font-medium text-gray-900 dark:text-white truncate">
										{getContactDisplayName(contact)}
									</h3>
									{#if contact.email}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
											Has Email
										</span>
									{/if}
								</div>

								<div class="mt-1 space-y-1">
									{#if contact.email}
										<p class="text-sm text-gray-500 dark:text-gray-400 flex items-center">
											<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
											</svg>
											{contact.email}
										</p>
									{/if}
									{#if contact.phone}
										<p class="text-sm text-gray-500 dark:text-gray-400 flex items-center">
											<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/>
											</svg>
											{contact.phone}
										</p>
									{/if}
									{#if contact.company}
										<p class="text-sm text-gray-500 dark:text-gray-400 flex items-center">
											<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
											</svg>
											{contact.company}
											{#if contact.position}
												• {contact.position}
											{/if}
										</p>
									{/if}
								</div>
							</div>
						</div>

						<!-- Actions -->
						<div class="flex items-center space-x-2">
							<button
								on:click={() => handleEditContact(contact)}
								class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300 p-1"
								title="Edit contact"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
								</svg>
							</button>
							<button
								on:click={() => handleDeleteContact(contact)}
								class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 p-1"
								title="Delete contact"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
								</svg>
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
