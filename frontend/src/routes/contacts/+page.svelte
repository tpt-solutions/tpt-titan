<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import ContactList from '$lib/components/ContactList.svelte';
	import ContactForm from '$lib/components/ContactForm.svelte';
	import { contacts } from '$lib/stores';

	// Accept framework-provided props to avoid warnings
	export let params = null;
	export let data = null;
	export let form = null;

	let showForm = false;
	let editingContact = null;
	let searchQuery = '';

	onMount(async () => {
		await loadContacts();
	});

	async function loadContacts() {
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
			}
		} catch (error) {
			console.error('Failed to load contacts:', error);
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

	<!-- Contact List -->
	<ContactList {handleEditContact} />

	<!-- Contact Form Modal -->
	{#if showForm}
		<ContactForm
			contact={editingContact}
			onClose={handleFormClose}
		/>
	{/if}
</div>
