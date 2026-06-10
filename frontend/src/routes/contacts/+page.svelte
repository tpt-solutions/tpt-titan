<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import ContactList from '$lib/components/ContactList.svelte';
	import ContactForm from '$lib/components/ContactForm.svelte';
	import { contacts } from '$lib/stores';

	// Accept framework-provided props to avoid warnings
	export let data = null;
	export let form = null;
	export let params = null;


	let showForm = false;
	let editingContact = null;
	let searchQuery = '';
	let isLoading = false;
	let loadError = null;

	onMount(async () => {
		await loadContacts();
	});

	async function loadContacts() {
		isLoading = true;
		loadError = null;
		try {
			const response = await fetch('/api/v1/contacts', {
				headers: {
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				}
			});

			if (response.ok) {
				const data = await response.json();
				contacts.set(data.contacts || []);
			} else if (response.status === 401) {
				goto('/auth/login');
			} else {
				loadError = `Failed to load contacts (${response.status})`;
			}
		} catch (error) {
			loadError = 'Could not connect to server. Please check your connection.';
			console.error('Failed to load contacts:', error);
		} finally {
			isLoading = false;
		}
	}

	async function searchContacts() {
		if (!searchQuery.trim()) {
			await loadContacts();
			return;
		}

		try {
			const response = await fetch(`/api/v1/contacts/search?q=${encodeURIComponent(searchQuery)}`, {
				headers: {
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				}
			});

			if (response.ok) {
				const data = await response.json();
				contacts.set(data.contacts || []);
			}
		} catch (error) {
			console.error('Failed to search contacts:', error);
		}
	}

	function handleCreateContact() {
		editingContact = null;
		showForm = true;
	}

	function handleEditContact(contact) {
		editingContact = contact;
		showForm = true;
	}

	function handleFormClose() {
		showForm = false;
		editingContact = null;
		loadContacts();
	}

	function handleSearchKeydown(event) {
		if (event.key === 'Enter') {
			searchContacts();
		}
	}
</script>

<svelte:head>
	<title>Contacts - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="flex justify-between items-center mb-8">
		<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Contacts</h1>
		<button
			on:click={handleCreateContact}
			class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
			</svg>
			Add Contact
		</button>
	</div>

	<!-- Search Bar -->
	<div class="mb-6">
		<div class="flex gap-2">
			<input
				bind:value={searchQuery}
				on:keydown={handleSearchKeydown}
				type="text"
				placeholder="Search contacts..."
				class="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
			>
			<button
				on:click={searchContacts}
				class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-lg flex items-center gap-2"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
				</svg>
				Search
			</button>
		</div>
	</div>

	<!-- Error banner -->
	{#if loadError}
		<div class="mb-6 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg px-4 py-3 flex items-center gap-3">
			<svg class="w-5 h-5 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
			</svg>
			<span class="text-sm text-red-700 dark:text-red-300">{loadError}</span>
			<button on:click={loadContacts} class="ml-auto text-sm text-red-600 dark:text-red-400 underline hover:no-underline">Retry</button>
		</div>
	{/if}

	<!-- Contact List -->
	{#if isLoading}
		<div class="flex items-center justify-center py-16 text-gray-400 dark:text-gray-500">
			<svg class="animate-spin w-8 h-8 mr-3" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
			</svg>
			<span>Loading contacts…</span>
		</div>
	{:else}
		<ContactList {handleEditContact} />
	{/if}

	<!-- Contact Form Modal -->
	{#if showForm}
		<ContactForm
			contact={editingContact}
			onClose={handleFormClose}
		/>
	{/if}
</div>
