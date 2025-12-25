<script>
	import FormBuilder from '$lib/components/FormBuilder.svelte';
	import FormList from '$lib/components/FormList.svelte';
	import { onMount } from 'svelte';

	let currentView = 'list'; // 'list' or 'builder'
	let selectedForm = null;
	let forms = [];

	onMount(() => {
		// Load user's forms
		loadForms();
	});

	function loadForms() {
		// Mock data for now - will connect to API later
		forms = [
			{
				id: 1,
				name: 'Customer Feedback Survey',
				description: 'Collect customer satisfaction data',
				responses: 24,
				createdAt: new Date('2024-01-15'),
				status: 'active'
			},
			{
				id: 2,
				name: 'Event Registration',
				description: 'Register attendees for company events',
				responses: 156,
				createdAt: new Date('2024-01-10'),
				status: 'active'
			},
			{
				id: 3,
				name: 'Job Application',
				description: 'Standard job application form',
				responses: 8,
				createdAt: new Date('2024-01-05'),
				status: 'draft'
			}
		];
	}

	function createNewForm() {
		selectedForm = null;
		currentView = 'builder';
	}

	function editForm(form) {
		selectedForm = form;
		currentView = 'builder';
	}

	function deleteForm(formId) {
		forms = forms.filter(f => f.id !== formId);
	}

	function saveForm(formData) {
		if (selectedForm) {
			// Update existing form
			forms = forms.map(f => f.id === selectedForm.id ? { ...f, ...formData } : f);
		} else {
			// Create new form
			const newForm = {
				id: Date.now(),
				...formData,
				responses: 0,
				createdAt: new Date(),
				status: 'draft'
			};
			forms = [...forms, newForm];
		}
		currentView = 'list';
		selectedForm = null;
	}

	function viewResponses(form) {
		// TODO: Navigate to responses view
		console.log('View responses for:', form);
	}
</script>

<svelte:head>
	<title>Forms & Templates - TPT Titan</title>
</svelte:head>

<div class="h-screen flex flex-col bg-gray-50">
	<!-- Header -->
	<header class="flex items-center justify-between px-6 py-3 border-b border-gray-200 bg-white">
		<div class="flex items-center space-x-4">
			<h1 class="text-xl font-semibold text-gray-900">Forms & Templates</h1>
			<span class="text-sm text-gray-500">MS Access-style database features</span>
		</div>

		<div class="flex items-center space-x-2">
			{#if currentView === 'list'}
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					on:click={createNewForm}
				>
					Create Form
				</button>
			{:else}
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
					on:click={() => { currentView = 'list'; selectedForm = null; }}
				>
					Back to Forms
				</button>
			{/if}
		</div>
	</header>

	<!-- Main Content -->
	<div class="flex-1 overflow-hidden">
		{#if currentView === 'list'}
			<FormList
				{forms}
				on:edit={e => editForm(e.detail)}
				on:delete={e => deleteForm(e.detail)}
				on:viewResponses={e => viewResponses(e.detail)}
			/>
		{:else}
			<FormBuilder
				form={selectedForm}
				on:save={e => saveForm(e.detail)}
				on:cancel={() => { currentView = 'list'; selectedForm = null; }}
			/>
		{/if}
	</div>
</div>
