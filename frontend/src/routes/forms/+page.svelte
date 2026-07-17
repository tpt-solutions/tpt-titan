<script>
	import FormBuilder from '$lib/components/FormBuilder.svelte';
	import FormList from '$lib/components/FormList.svelte';
	import FormAdvancedPanel from '$lib/components/FormAdvancedPanel.svelte';
	import DatabaseRelationsModal from '$lib/components/DatabaseRelationsModal.svelte';
	import AdvancedReportsModal from '$lib/components/AdvancedReportsModal.svelte';
	import WorkflowDesignerModal from '$lib/components/WorkflowDesignerModal.svelte';
	import { getForms, getFormResponses } from '$lib/api.js';
	import { onMount } from 'svelte';

	// Accept framework-provided props to avoid warnings
	export const data = null;
	export const form = null;
	export const params = null;


	let currentView = 'list'; // 'list', 'builder', or 'advanced'
	let selectedForm = null;
	let forms = [];
	let showDatabaseRelationsModal = false;
	let showAdvancedReportsModal = false;
	let showWorkflowDesignerModal = false;
	let showResponsesModal = false;
	let responsesForm = null;
	let responses = [];
	let responsesLoading = false;
	let responsesError = '';

	onMount(() => {
		// Load user's forms
		loadForms();
	});

	async function loadForms() {
		try {
			const response = await getForms();
			// Transform API response to match component expectations
			forms = response.forms.map(form => ({
				...form,
				createdAt: new Date(form.created_at || form.createdAt),
				responses: form.responses || 0 // Add responses count if not provided
			}));
		} catch (error) {
			console.warn('Backend not available, using demo data:', error.message);
			// Fallback to demo data when backend is not running
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
	}

	function createNewForm() {
		selectedForm = null;
		currentView = 'builder';
	}

	function editForm(form) {
		selectedForm = form;
		currentView = 'builder';
	}

	function openAdvanced(form) {
		selectedForm = form;
		currentView = 'advanced';
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

	async function viewResponses(form) {
		responsesForm = form;
		showResponsesModal = true;
		responsesLoading = true;
		responsesError = '';
		responses = [];
		try {
			const result = await getFormResponses(form.id);
			responses = result.responses || result || [];
		} catch (error) {
			responsesError = error.message || 'Failed to load responses';
		} finally {
			responsesLoading = false;
		}
	}

	function handleReorder(event) {
		forms = event.detail.forms;
		// TODO: Persist order to backend
		console.log('Forms reordered:', forms);
	}

	function handleAdvanced(event) {
		openAdvanced(event.detail);
	}

	function handleDatabaseRelations() {
		showDatabaseRelationsModal = true;
	}

	function handleAdvancedReports() {
		showAdvancedReportsModal = true;
	}

	function handleWorkflowDesigner() {
		showWorkflowDesignerModal = true;
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
					class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
					on:click={() => { if (selectedForm) openAdvanced(selectedForm); else alert('Select a form first'); }}
				>
					Advanced Modules
				</button>
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					on:click={createNewForm}
				>
					Create Form
				</button>
			{:else if currentView === 'advanced'}
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
					on:click={() => { currentView = 'list'; selectedForm = null; }}
				>
					Back to Forms
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
				on:reorder={handleReorder}
				on:openDatabaseRelations={handleDatabaseRelations}
				on:openAdvancedReports={handleAdvancedReports}
				on:openWorkflowDesigner={handleWorkflowDesigner}
				on:openAdvanced={handleAdvanced}
			/>
		{:else if currentView === 'advanced'}
			<FormAdvancedPanel form={selectedForm} />
		{:else}
			<FormBuilder
				form={selectedForm}
				on:save={e => saveForm(e.detail)}
				on:cancel={() => { currentView = 'list'; selectedForm = null; }}
			/>
		{/if}
	</div>
</div>

<!-- Database Relations Modal -->
<DatabaseRelationsModal
	{forms}
	showModal={showDatabaseRelationsModal}
	on:close={() => showDatabaseRelationsModal = false}
/>

<!-- Advanced Reports Modal -->
<AdvancedReportsModal
	{forms}
	showModal={showAdvancedReportsModal}
	on:close={() => showAdvancedReportsModal = false}
/>

<!-- Workflow Designer Modal -->
<WorkflowDesignerModal
	{forms}
	showModal={showWorkflowDesignerModal}
	on:close={() => showWorkflowDesignerModal = false}
/>

<!-- Form Responses Modal -->
{#if showResponsesModal}
	<div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
		<div class="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[80vh] flex flex-col">
			<div class="flex items-center justify-between px-6 py-4 border-b border-gray-200">
				<h2 class="text-lg font-semibold text-gray-900">
					Responses — {responsesForm?.name}
				</h2>
				<button
					class="text-gray-400 hover:text-gray-600"
					on:click={() => showResponsesModal = false}
				>
					✕
				</button>
			</div>
			<div class="flex-1 overflow-auto px-6 py-4">
				{#if responsesLoading}
					<p class="text-sm text-gray-500">Loading responses…</p>
				{:else if responsesError}
					<p class="text-sm text-red-600">{responsesError}</p>
				{:else if responses.length === 0}
					<p class="text-sm text-gray-500">No responses yet.</p>
				{:else}
					<ul class="space-y-3">
						{#each responses as response}
							<li class="border border-gray-200 rounded-lg p-3 text-sm">
								<pre class="whitespace-pre-wrap break-words text-gray-700">{JSON.stringify(response.data ?? response, null, 2)}</pre>
							</li>
						{/each}
					</ul>
				{/if}
			</div>
		</div>
	</div>
{/if}
