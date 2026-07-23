<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import EmailInbox from '$lib/components/EmailInbox.svelte';
	import EmailViewer from '$lib/components/EmailViewer.svelte';
	import EmailComposer from '$lib/components/EmailComposer.svelte';
	import { emailAccounts, emails, selectedEmail, currentFolder } from '$lib/stores';

	// Accept framework-provided props to avoid warnings
	export let data = null;
	export let form = null;
	export let params = null;


	let showComposer = false;
	let emailAccountsList = [];
	let emailsList = [];
	let selectedEmailData = null;
	let currentFolderValue = 'inbox';
	let isLoadingAccounts = false;
	let isLoadingEmails = false;
	let loadError = null;

	emailAccounts.subscribe(value => emailAccountsList = value);
	emails.subscribe(value => emailsList = value);
	selectedEmail.subscribe(value => selectedEmailData = value);
	currentFolder.subscribe(value => currentFolderValue = value);

	onMount(async () => {
		await loadEmailAccounts();
		await loadEmails();
	});

	async function loadEmailAccounts() {
		isLoadingAccounts = true;
		loadError = null;
		try {
			const response = await fetch('/api/v1/email-accounts', {
				headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
			});
			if (response.ok) {
				const data = await response.json();
				emailAccounts.set(data.accounts || []);
			} else if (response.status === 401) {
				goto('/auth/login');
			} else {
				loadError = `Failed to load email accounts (${response.status})`;
			}
		} catch (error) {
			loadError = 'Could not connect to server. Please check your connection.';
			console.error('Failed to load email accounts:', error);
		} finally {
			isLoadingAccounts = false;
		}
	}

	async function loadEmails() {
		isLoadingEmails = true;
		try {
			const params = new URLSearchParams();
			if (currentFolderValue) params.append('folder', currentFolderValue);
			const response = await fetch(`/api/v1/emails?${params}`, {
				headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
			});
			if (response.ok) {
				const data = await response.json();
				emails.set(data.emails || []);
			} else if (response.status === 401) {
				goto('/auth/login');
			} else {
				console.error('Failed to load emails:', response.status);
			}
		} catch (error) {
			console.error('Failed to load emails:', error);
		} finally {
			isLoadingEmails = false;
		}
	}

	function handleComposeEmail() {
		showComposer = true;
	}

	function handleEmailSelect(email) {
		selectedEmail.set(email);
	}

	function handleFolderChange(folder) {
		currentFolder.set(folder);
		loadEmails();
		selectedEmail.set(null);
	}

	function handleComposerClose() {
		showComposer = false;
		loadEmails(); // Refresh emails after sending
	}

	$: if (currentFolderValue) {
		loadEmails();
	}
</script>

<svelte:head>
	<title>Email - TPT Titan</title>
</svelte:head>

<div class="h-screen flex flex-col bg-gray-50 dark:bg-gray-900">
	<!-- Email Header -->
	<div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
		<div class="flex justify-between items-center">
			<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Email</h1>
			<button
				on:click={handleComposeEmail}
				class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
				</svg>
				Compose
			</button>
		</div>

		<!-- Error banner -->
		{#if loadError}
			<div class="mt-3 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg px-4 py-3 flex items-center gap-3">
				<svg class="w-5 h-5 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
				</svg>
				<span class="text-sm text-red-700 dark:text-red-300">{loadError}</span>
				<button on:click={loadEmailAccounts} class="ml-auto text-sm text-red-600 dark:text-red-400 underline hover:no-underline">Retry</button>
			</div>
		{/if}

		<!-- Email Accounts -->
		{#if isLoadingAccounts}
			<div class="mt-4 flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
				<svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
				</svg>
				Loading email accounts…
			</div>
		{:else if emailAccountsList.length > 0}
			<div class="mt-4 flex items-center space-x-4">
				<span class="text-sm text-gray-600 dark:text-gray-400">Accounts:</span>
				<div class="flex space-x-2">
					{#each emailAccountsList as account}
						<span class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
							{account.display_name || account.email}
						</span>
					{/each}
				</div>
			</div>
		{:else}
			<div class="mt-4 text-center py-8">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No email accounts</h3>
				<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Add an email account to get started.</p>
				<button
					on:click={() => showComposer = true}
					class="mt-4 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
				>
					Add Email Account
				</button>
			</div>
		{/if}
	</div>

	<!-- Main Email Interface -->
	{#if isLoadingAccounts}
		<!-- Skeleton while accounts load -->
		<div class="flex-1 flex items-center justify-center">
			<div class="text-center text-gray-400 dark:text-gray-500">
				<svg class="animate-spin mx-auto w-8 h-8 mb-3" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
				</svg>
				<p class="text-sm">Loading…</p>
			</div>
		</div>
	{:else if emailAccountsList.length > 0}
		<div class="flex-1 flex overflow-hidden">
			<!-- Email Inbox Sidebar -->
			<div class="w-80 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col">
				<EmailInbox
					{emailsList}
					{currentFolderValue}
					{selectedEmailData}
					{handleEmailSelect}
					{handleFolderChange}
				/>
			</div>

			<!-- Email Viewer/Composer -->
			<div class="flex-1 flex flex-col min-w-0">
				{#if showComposer}
					<EmailComposer
						{emailAccountsList}
						onClose={handleComposerClose}
					/>
				{:else if selectedEmailData}
					<EmailViewer
						email={selectedEmailData}
						onClose={() => selectedEmail.set(null)}
					/>
				{:else}
					<!-- Empty state -->
					<div class="flex-1 flex items-center justify-center bg-gray-50 dark:bg-gray-900">
						<div class="text-center">
							<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
							</svg>
							<h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No email selected</h3>
							<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Choose an email from the list to view it.</p>
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
