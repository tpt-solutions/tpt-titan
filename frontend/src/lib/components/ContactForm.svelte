<script>
// @ts-nocheck
	import { createEventDispatcher } from 'svelte';

	export let contact = null; // null for create, object for edit
	export let onClose;

	const dispatch = createEventDispatcher();

	let formData = {
		first_name: '',
		last_name: '',
		emails: [{ type: 'work', value: '' }],
		phones: [{ type: 'mobile', value: '' }],
		company: '',
		position: '',
		notes: ''
	};

	let isSubmitting = false;
	let errors = {};

	// Initialize form data if editing
	if (contact) {
		formData = {
			first_name: contact.first_name || '',
			last_name: contact.last_name || '',
			emails: contact.emails && contact.emails.length > 0 ? contact.emails : [{ type: 'work', value: '' }],
			phones: contact.phones && contact.phones.length > 0 ? contact.phones : [{ type: 'mobile', value: '' }],
			company: contact.company || '',
			position: contact.position || '',
			notes: contact.notes || ''
		};
	}

	function validateForm() {
		errors = {};

		if (!formData.first_name.trim() && !formData.last_name.trim()) {
			errors.name = 'At least first name or last name is required';
		}

		// Validate emails
		formData.emails.forEach((email, index) => {
			if (email.value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
				errors[`email_${index}`] = 'Please enter a valid email address';
			}
		});

		return Object.keys(errors).length === 0;
	}

	function addEmail() {
		formData.emails = [...formData.emails, { type: 'personal', value: '' }];
	}

	function removeEmail(index) {
		if (formData.emails.length > 1) {
			formData.emails = formData.emails.filter((_, i) => i !== index);
		}
	}

	function addPhone() {
		formData.phones = [...formData.phones, { type: 'home', value: '' }];
	}

	function removePhone(index) {
		if (formData.phones.length > 1) {
			formData.phones = formData.phones.filter((_, i) => i !== index);
		}
	}

	async function handleSubmit() {
		if (!validateForm()) {
			return;
		}

		isSubmitting = true;

		try {
			const method = contact ? 'PUT' : 'POST';
			const url = contact ? `/api/v1/contacts/${contact.id}` : '/api/v1/contacts';

			const response = await fetch(url, {
				method,
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify(formData)
			});

			if (response.ok) {
				dispatch('success');
				handleClose();
			} else {
				const errorData = await response.json();
				errors.general = errorData.error || 'Failed to save contact';
			}
		} catch (error) {
			console.error('Failed to save contact:', error);
			errors.general = 'Network error. Please try again.';
		} finally {
			isSubmitting = false;
		}
	}

	function handleClose() {
		formData = {
			first_name: '',
			last_name: '',
			emails: [{ type: 'work', value: '' }],
			phones: [{ type: 'mobile', value: '' }],
			company: '',
			position: '',
			notes: ''
		};
		errors = {};
		onClose();
	}

	function handleKeydown(event) {
		if (event.key === 'Escape') {
			handleClose();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal backdrop -->
<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" on:click={handleClose} role="presentation">
	<!-- Modal container -->
	<div class="relative top-20 mx-auto p-5 border w-full max-w-lg shadow-lg rounded-md bg-white dark:bg-gray-800" on:click|stopPropagation role="dialog" aria-modal="true" aria-labelledby="contact-form-title">
		<div class="mt-3">
			<!-- Header -->
			<div class="flex items-center justify-between mb-4">
				<h3 id="contact-form-title" class="text-lg font-medium text-gray-900 dark:text-white">
					{contact ? 'Edit Contact' : 'Add New Contact'}
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

				<!-- Name Fields -->
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="first_name" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							First Name
						</label>
						<input
							id="first_name"
							type="text"
							bind:value={formData.first_name}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							placeholder="John"
						>
					</div>
					<div>
						<label for="last_name" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Last Name
						</label>
						<input
							id="last_name"
							type="text"
							bind:value={formData.last_name}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							placeholder="Doe"
						>
					</div>
				</div>
				{#if errors.name}
					<p class="text-sm text-red-600 dark:text-red-400">{errors.name}</p>
				{/if}

				<!-- Emails -->
				<div>
					<div class="flex items-center justify-between mb-2">
						<label id="email-label" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Email Addresses
						</label>
						<button
							type="button"
							on:click={addEmail}
							class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
							aria-label="Add another email address"
						>
							+ Add Email
						</button>
					</div>
					<div role="group" aria-labelledby="email-label">
						{#each formData.emails as email, index}
							<div class="flex gap-2 mb-2">
								<label for="email-type-{index}" class="sr-only">Email {index + 1} type</label>
								<select
									id="email-type-{index}"
									bind:value={email.type}
									class="flex-shrink-0 w-20 px-2 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500 text-sm"
									aria-label="Type for email {index + 1}"
								>
									<option value="work">Work</option>
									<option value="personal">Personal</option>
									<option value="other">Other</option>
								</select>
								<label for="email-value-{index}" class="sr-only">Email {index + 1} address</label>
								<input
									id="email-value-{index}"
									type="email"
									bind:value={email.value}
									class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
									placeholder="john.doe@example.com"
									aria-label="Email {index + 1} address"
								>
								{#if formData.emails.length > 1}
									<button
										type="button"
										on:click={() => removeEmail(index)}
										class="flex-shrink-0 px-2 py-2 text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
										title="Remove email"
										aria-label="Remove email {index + 1}"
									>
										✕
									</button>
								{/if}
							</div>
							{#if errors[`email_${index}`]}
								<p class="text-sm text-red-600 dark:text-red-400 mb-2">{errors[`email_${index}`]}</p>
							{/if}
						{/each}
					</div>
				</div>

				<!-- Phones -->
				<div>
					<div class="flex items-center justify-between mb-2">
						<label id="phone-label" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Phone Numbers
						</label>
						<button
							type="button"
							on:click={addPhone}
							class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
							aria-label="Add another phone number"
						>
							+ Add Phone
						</button>
					</div>
					<div role="group" aria-labelledby="phone-label">
						{#each formData.phones as phone, index}
							<div class="flex gap-2 mb-2">
								<label for="phone-type-{index}" class="sr-only">Phone {index + 1} type</label>
								<select
									id="phone-type-{index}"
									bind:value={phone.type}
									class="flex-shrink-0 w-20 px-2 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500 text-sm"
									aria-label="Type for phone {index + 1}"
								>
									<option value="mobile">Mobile</option>
									<option value="work">Work</option>
									<option value="home">Home</option>
									<option value="other">Other</option>
								</select>
								<label for="phone-value-{index}" class="sr-only">Phone {index + 1} number</label>
								<input
									id="phone-value-{index}"
									type="tel"
									bind:value={phone.value}
									class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
									placeholder="+1 (555) 123-4567"
									aria-label="Phone {index + 1} number"
								>
								{#if formData.phones.length > 1}
									<button
										type="button"
										on:click={() => removePhone(index)}
										class="flex-shrink-0 px-2 py-2 text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
										title="Remove phone"
										aria-label="Remove phone {index + 1}"
									>
										✕
									</button>
								{/if}
							</div>
						{/each}
					</div>
				</div>

				<!-- Company and Position -->
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="company" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Company
						</label>
						<input
							id="company"
							type="text"
							bind:value={formData.company}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							placeholder="Acme Corp"
						>
					</div>
					<div>
						<label for="position" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Position
						</label>
						<input
							id="position"
							type="text"
							bind:value={formData.position}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
							placeholder="Software Engineer"
						>
					</div>
				</div>

				<!-- Notes -->
				<div>
					<label for="notes" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Notes
					</label>
					<textarea
						id="notes"
						bind:value={formData.notes}
						rows="3"
						class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
						placeholder="Additional notes about this contact..."
					></textarea>
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
							<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
							</svg>
							Saving...
						{:else}
							{contact ? 'Update Contact' : 'Create Contact'}
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
</div>
