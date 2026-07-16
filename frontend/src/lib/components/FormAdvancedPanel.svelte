<script>
	import { onMount } from 'svelte';
	import {
		getFormTemplates,
		createFormTemplate,
		getFormTemplateCategories,
		useFormTemplate,
		getFormRelationships,
		createRelationship,
		getFormHierarchy,
		getFormReports,
		createReport,
		generateAdHocReport,
		getAvailableTables,
		buildSQL,
		executeVisualQuery,
		validateVisualQuery,
		getQuerySuggestions,
		saveQueryTemplate,
		getQueryTemplates,
		getEmailDistributions,
		createEmailDistribution,
		sendFormResponseEmail,
		createFormWorkflow,
		getPendingApprovals,
		createNotificationTemplate
	} from '$lib/api.js';

	export let form = null; // selected form object (needs .id)

	let activeTab = 'templates';
	let loading = false;
	let error = '';
	let message = '';

	// Templates
	let templates = [];
	let templateCategories = [];
	let newTemplate = { name: '', description: '', category: '', form_data: { fields: [] }, is_public: false, tags: [] };
	let templateSearch = '';
	let selectedTemplateId = '';
	let useTemplateName = '';

	// Relationships
	let relationships = [];
	let hierarchy = null;
	let newRelationship = { name: '', type: 'one_to_many', related_form_id: '', foreign_key: '', local_key: '' };
	let relatedForms = [];

	// Reports
	let reports = [];
	let newReport = { name: '', description: '', report_type: 'table', config: {} };
	let adHocConfig = { report_type: 'table', group_by: '', aggregate: '' };
	let adHocResult = null;

	// Query builder
	let tables = [];
	let queryElements = [];
	let builtSQL = '';
	let queryResult = null;
	let queryError = '';
	let savedQueries = [];
	let queryTemplatesName = '';
	let newElement = { type: 'select', table: '', field: '', alias: '' };

	// Email distribution
	let distributions = [];
	let newDistribution = { name: '', recipients: '', subject: '', message: '', trigger: 'immediate', is_active: true, include_data: false };
	let sendFormId = '';
	let sendRecipients = '';
	let sendSubject = '';

	// Workflow
	let pendingApprovals = [];
	let newNotificationTemplate = { name: '', subject: '', body: '', type: 'email' };
	let newWorkflow = { name: '', description: '', trigger: '', steps: [] };

	$: formId = form && (form.id || form.ID);

	onMount(async () => {
		if (!formId) return;
		await Promise.all([loadTemplates(), loadRelationships(), loadReports(), loadTables(), loadEmailDistributions(), loadPendingApprovals()]);
	});

	function flash(msg, isError = false) {
		if (isError) { error = msg; message = ''; } else { message = msg; error = ''; }
		setTimeout(() => { error = ''; message = ''; }, 4000);
	}

	// ── Templates ────────────────────────────────────────────────
	async function loadTemplates() {
		try {
			const params = {};
			if (templateSearch) params.search = templateSearch;
			const res = await getFormTemplates(params);
			templates = res.templates || [];
		} catch (e) { flash(e.message, true); }
	}
	async function loadCategories() {
		try { const res = await getFormTemplateCategories(); templateCategories = res.categories || []; } catch (e) { flash(e.message, true); }
	}
	async function saveTemplate() {
		loading = true;
		try { await createFormTemplate(newTemplate); flash('Template created'); newTemplate = { name: '', description: '', category: '', form_data: { fields: [] }, is_public: false, tags: [] }; await loadTemplates(); }
		catch (e) { flash(e.message, true); } finally { loading = false; }
	}
	async function applyTemplate() {
		if (!selectedTemplateId) return;
		loading = true;
		try { await useFormTemplate(selectedTemplateId, { name: useTemplateName || 'New Form from Template', description: '' }); flash('Form created from template'); }
		catch (e) { flash(e.message, true); } finally { loading = false; }
	}

	// ── Relationships ────────────────────────────────────────────
	async function loadRelationships() {
		try { const res = await getFormRelationships(formId); relationships = res.relationships || []; } catch (e) { flash(e.message, true); }
		try { const res = await getFormHierarchy(formId); hierarchy = res; } catch (e) {}
	}
	async function saveRelationship() {
		loading = true;
		try { await createRelationship(formId, { ...newRelationship }); flash('Relationship created'); newRelationship = { name: '', type: 'one_to_many', related_form_id: '', foreign_key: '', local_key: '' }; await loadRelationships(); }
		catch (e) { flash(e.message, true); } finally { loading = false; }
	}

	// ── Reports ──────────────────────────────────────────────────
	async function loadReports() {
		try { const res = await getFormReports(formId); reports = res.reports || []; } catch (e) { flash(e.message, true); }
	}
	async function saveReport() {
		loading = true;
		try { await createReport(formId, newReport); flash('Report created'); newReport = { name: '', description: '', report_type: 'table', config: {} }; await loadReports(); }
		catch (e) { flash(e.message, true); } finally { loading = false; }
	}
	async function runAdHoc() {
		loading = true;
		try { const res = await generateAdHocReport(formId, adHocConfig); adHocResult = res; flash('Ad-hoc report generated'); }
		catch (e) { flash(e.message, true); } finally { loading = false; }
	}

	// ── Query builder ────────────────────────────────────────────
	async function loadTables() {
		try { const res = await getAvailableTables(formId); tables = res.tables || []; } catch (e) { flash(e.message, true); }
	}
	async function addElement() {
		queryElements = [...queryElements, { ...newElement }];
		newElement = { type: 'select', table: '', field: '', alias: '' };
	}
	function removeElement(i) { queryElements = queryElements.filter((_, idx) => idx !== i); }
	async function buildQuery() {
		loading = true;
		try { const res = await buildSQL(formId, queryElements); builtSQL = res.sql; queryError = ''; }
		catch (e) { queryError = e.message; builtSQL = ''; } finally { loading = false; }
	}
	async function validateQuery() {
		try { const res = await validateVisualQuery(formId, queryElements); if (res.valid) flash('Query is valid'); else flash((res.errors || []).join(', '), true); }
		catch (e) { flash(e.message, true); }
	}
	async function runQuery() {
		loading = true;
		try { const res = await executeVisualQuery(formId, queryElements); queryResult = res; flash('Query executed'); }
		catch (e) { flash(e.message, true); } finally { loading = false; }
	}
	async function suggestQuery() {
		try { const res = await getQuerySuggestions(formId, queryElements); flash(((res.suggestions || []).map(s => s.description || s).join('; ') || 'No suggestions')); }
		catch (e) { flash(e.message, true); }
	}
	async function saveQuery() {
		if (!queryTemplatesName) { flash('Enter a name', true); return; }
		try { await saveQueryTemplate(formId, { name: queryTemplatesName, description: '', elements: queryElements }); flash('Query template saved'); queryTemplatesName = ''; }
		catch (e) { flash(e.message, true); }
	}
	async function loadSavedQueries() {
		try { const res = await getQueryTemplates(formId); savedQueries = res.templates || []; } catch (e) { flash(e.message, true); }
	}

	// ── Email distribution ───────────────────────────────────────
	async function loadEmailDistributions() {
		try { const res = await getEmailDistributions(formId); distributions = res.distributions || []; } catch (e) { flash(e.message, true); }
	}
	async function saveDistribution() {
		loading = true;
		try {
			const payload = { ...newDistribution, recipients: newDistribution.recipients.split(',').map(s => s.trim()).filter(Boolean) };
			await createEmailDistribution(formId, payload);
			flash('Email distribution created');
			newDistribution = { name: '', recipients: '', subject: '', message: '', trigger: 'immediate', is_active: true, include_data: false };
			await loadEmailDistributions();
		} catch (e) { flash(e.message, true); } finally { loading = false; }
	}
	async function sendEmail() {
		if (!sendFormId) { flash('Enter a response ID', true); return; }
		try { await sendFormResponseEmail(formId, sendFormId, { recipients: sendRecipients.split(',').map(s => s.trim()).filter(Boolean), subject: sendSubject }); flash('Email sent'); }
		catch (e) { flash(e.message, true); }
	}

	// ── Workflow ─────────────────────────────────────────────────
	async function loadPendingApprovals() {
		try { const res = await getPendingApprovals(formId); pendingApprovals = res.approvals || []; } catch (e) { flash(e.message, true); }
	}
	async function saveWorkflow() {
		loading = true;
		try { await createFormWorkflow(formId, { ...newWorkflow, steps: newWorkflow.steps.filter(Boolean).map(s => ({ name: s })) }); flash('Workflow created'); newWorkflow = { name: '', description: '', trigger: '', steps: [] }; }
		catch (e) { flash(e.message, true); } finally { loading = false; }
	}
	async function saveNotificationTemplate() {
		try { await createNotificationTemplate(formId, newNotificationTemplate); flash('Notification template created'); newNotificationTemplate = { name: '', subject: '', body: '', type: 'email' }; }
		catch (e) { flash(e.message, true); }
	}
</script>

<svelte:head>
	<title>Advanced Forms - TPT Titan</title>
</svelte:head>

<div class="h-full flex flex-col bg-white">
	{#if !formId}
		<div class="flex-1 flex items-center justify-center text-gray-500">Select a form to manage its advanced modules.</div>
	{:else}
		<div class="border-b border-gray-200 px-4 pt-3 flex items-center justify-between bg-gray-50">
			<div class="flex space-x-1 overflow-x-auto">
				{#each [['templates','Templates'],['relationships','Relationships'],['reports','Reports'],['query','Query Builder'],['email','Email'],['workflow','Workflow']] as [key,label]}
					<button
						class="px-3 py-2 text-sm rounded-t-md transition-colors {activeTab === key ? 'bg-white text-blue-700 border border-b-0 border-gray-200 font-medium' : 'text-gray-600 hover:bg-gray-100'}"
						on:click={() => activeTab = key}
					>{label}</button>
				{/each}
			</div>
		</div>

		{#if message}<div class="px-4 py-2 text-sm text-green-700 bg-green-50">{message}</div>{/if}
		{#if error}<div class="px-4 py-2 text-sm text-red-700 bg-red-50">{error}</div>{/if}

		<div class="flex-1 overflow-auto p-4">
			{#if loading}<div class="text-sm text-gray-500">Working…</div>{/if}

			{#if activeTab === 'templates'}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Browse Templates</h3>
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Search templates…" bind:value={templateSearch} on:input={loadTemplates} />
						<button class="text-xs text-blue-600 mb-2" on:click={loadCategories}>Load categories</button>
						<ul class="space-y-1 text-sm max-h-64 overflow-auto">
							{#each templates as t}
								<li class="flex items-center justify-between border-b pb-1">
									<span>{t.name} <span class="text-gray-400">({t.category})</span></span>
									<button class="text-blue-600 text-xs" on:click={() => { selectedTemplateId = t.id; useTemplateName = t.name + ' (copy)'; }}>Use</button>
								</li>
							{/each}
						</ul>
						{#if selectedTemplateId}
							<div class="mt-3 border-t pt-3">
								<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="New form name" bind:value={useTemplateName} />
								<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={applyTemplate}>Create form from template</button>
							</div>
						{/if}
					</div>
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Create Template</h3>
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Name" bind:value={newTemplate.name} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Description" bind:value={newTemplate.description} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Category" bind:value={newTemplate.category} />
						<button class="px-3 py-1 bg-green-600 text-white rounded text-sm" on:click={saveTemplate}>Save Template</button>
					</div>
				</div>

			{:else if activeTab === 'relationships'}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Relationships</h3>
						<ul class="space-y-1 text-sm max-h-64 overflow-auto mb-3">
							{#each relationships as r}
								<li class="border-b pb-1">{r.name || r.type} <span class="text-gray-400">({r.type})</span></li>
							{/each}
						</ul>
						<div class="border-t pt-3 space-y-2">
							<input class="w-full border rounded px-2 py-1 text-sm" placeholder="Relationship name" bind:value={newRelationship.name} />
							<select class="w-full border rounded px-2 py-1 text-sm" bind:value={newRelationship.type}>
								<option value="one_to_one">One to One</option>
								<option value="one_to_many">One to Many</option>
								<option value="many_to_many">Many to Many</option>
							</select>
							<input class="w-full border rounded px-2 py-1 text-sm" placeholder="Related form ID" bind:value={newRelationship.related_form_id} />
							<input class="w-full border rounded px-2 py-1 text-sm" placeholder="Foreign key" bind:value={newRelationship.foreign_key} />
							<input class="w-full border rounded px-2 py-1 text-sm" placeholder="Local key" bind:value={newRelationship.local_key} />
							<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={saveRelationship}>Add Relationship</button>
						</div>
					</div>
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Form Hierarchy</h3>
						<pre class="text-xs bg-gray-50 p-2 rounded overflow-auto max-h-80">{JSON.stringify(hierarchy, null, 2)}</pre>
					</div>
				</div>

			{:else if activeTab === 'reports'}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Reports</h3>
						<ul class="space-y-1 text-sm max-h-48 overflow-auto mb-3">
							{#each reports as r}
								<li class="border-b pb-1">{r.name} <span class="text-gray-400">({r.report_type})</span></li>
							{/each}
						</ul>
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Report name" bind:value={newReport.name} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Description" bind:value={newReport.description} />
						<select class="w-full border rounded px-2 py-1 text-sm mb-2" bind:value={newReport.report_type}>
							<option value="table">Table</option>
							<option value="chart">Chart</option>
							<option value="summary">Summary</option>
						</select>
						<button class="px-3 py-1 bg-green-600 text-white rounded text-sm" on:click={saveReport}>Create Report</button>
					</div>
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Ad-hoc Report</h3>
						<select class="w-full border rounded px-2 py-1 text-sm mb-2" bind:value={adHocConfig.report_type}>
							<option value="table">Table</option>
							<option value="chart">Chart</option>
							<option value="summary">Summary</option>
						</select>
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Group by" bind:value={adHocConfig.group_by} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Aggregate" bind:value={adHocConfig.aggregate} />
						<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={runAdHoc}>Generate</button>
						{#if adHocResult}
							<pre class="text-xs bg-gray-50 p-2 rounded overflow-auto max-h-48 mt-2">{JSON.stringify(adHocResult, null, 2)}</pre>
						{/if}
					</div>
				</div>

			{:else if activeTab === 'query'}
				<div class="space-y-3">
					<div class="flex items-end space-x-2 flex-wrap">
						<select class="border rounded px-2 py-1 text-sm" bind:value={newElement.table}>
							<option value="">Table…</option>
							{#each tables as t}<option value={t.name || t}>{t.name || t}</option>{/each}
						</select>
						<input class="border rounded px-2 py-1 text-sm" placeholder="Field" bind:value={newElement.field} />
						<select class="border rounded px-2 py-1 text-sm" bind:value={newElement.type}>
							<option value="select">select</option>
							<option value="where">where</option>
							<option value="order_by">order_by</option>
							<option value="group_by">group_by</option>
							<option value="join">join</option>
						</select>
						<button class="px-3 py-1 bg-gray-600 text-white rounded text-sm" on:click={addElement}>Add</button>
					</div>
					<ul class="space-y-1 text-sm">
						{#each queryElements as el, i}
							<li class="flex items-center justify-between border rounded px-2 py-1">
								<span><span class="text-gray-500">{el.type}</span> {el.table}.{el.field}</span>
								<button class="text-red-600 text-xs" on:click={() => removeElement(i)}>remove</button>
							</li>
						{/each}
					</ul>
					<div class="flex space-x-2 flex-wrap">
						<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={buildQuery}>Build SQL</button>
						<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={runQuery}>Execute</button>
						<button class="px-3 py-1 bg-gray-600 text-white rounded text-sm" on:click={validateQuery}>Validate</button>
						<button class="px-3 py-1 bg-gray-600 text-white rounded text-sm" on:click={suggestQuery}>Suggest</button>
						<button class="px-3 py-1 bg-green-600 text-white rounded text-sm" on:click={saveQuery}>Save Template</button>
						<button class="px-3 py-1 bg-gray-300 text-gray-700 rounded text-sm" on:click={loadSavedQueries}>Load Saved</button>
					</div>
					{#if builtSQL}<pre class="text-xs bg-gray-50 p-2 rounded">{builtSQL}</pre>{/if}
					{#if queryError}<div class="text-red-600 text-sm">{queryError}</div>{/if}
					{#if queryResult}
						<pre class="text-xs bg-gray-50 p-2 rounded overflow-auto max-h-64">{JSON.stringify(queryResult, null, 2)}</pre>
					{/if}
					{#if savedQueries.length}
						<div class="text-xs text-gray-500">Saved templates: {savedQueries.map(q => q.name).join(', ')}</div>
					{/if}
					<input class="border rounded px-2 py-1 text-sm" placeholder="Query template name" bind:value={queryTemplatesName} />
				</div>

			{:else if activeTab === 'email'}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Distributions</h3>
						<ul class="space-y-1 text-sm max-h-48 overflow-auto mb-3">
							{#each distributions as d}
								<li class="border-b pb-1">{d.name} <span class="text-gray-400">({d.trigger})</span></li>
							{/each}
						</ul>
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Name" bind:value={newDistribution.name} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Recipients (comma separated)" bind:value={newDistribution.recipients} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Subject" bind:value={newDistribution.subject} />
						<textarea class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Message" bind:value={newDistribution.message}></textarea>
						<select class="w-full border rounded px-2 py-1 text-sm mb-2" bind:value={newDistribution.trigger}>
							<option value="immediate">Immediate</option>
							<option value="daily">Daily</option>
							<option value="weekly">Weekly</option>
						</select>
						<label class="text-sm flex items-center space-x-2 mb-2">
							<input type="checkbox" bind:checked={newDistribution.include_data} /> Include response data
						</label>
						<button class="px-3 py-1 bg-green-600 text-white rounded text-sm" on:click={saveDistribution}>Create Distribution</button>
					</div>
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Send Response Email</h3>
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Response ID" bind:value={sendFormId} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Recipients (comma separated)" bind:value={sendRecipients} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Subject" bind:value={sendSubject} />
						<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={sendEmail}>Send</button>
					</div>
				</div>

			{:else if activeTab === 'workflow'}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Create Workflow</h3>
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Name" bind:value={newWorkflow.name} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Description" bind:value={newWorkflow.description} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Trigger" bind:value={newWorkflow.trigger} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Step (one per line)" on:input={(e) => newWorkflow.steps = e.target.value.split('\n')} />
						<button class="px-3 py-1 bg-green-600 text-white rounded text-sm" on:click={saveWorkflow}>Create Workflow</button>
					</div>
					<div class="border rounded-lg p-4">
						<h3 class="font-semibold mb-2">Pending Approvals</h3>
						<ul class="space-y-1 text-sm max-h-48 overflow-auto mb-3">
							{#each pendingApprovals as a}
								<li class="border-b pb-1">{a.name || a.id || 'Approval'}</li>
							{/each}
							{#if !pendingApprovals.length}<li class="text-gray-400">None</li>{/if}
						</ul>
						<h3 class="font-semibold mb-2">Notification Template</h3>
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Name" bind:value={newNotificationTemplate.name} />
						<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Subject" bind:value={newNotificationTemplate.subject} />
						<textarea class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Body" bind:value={newNotificationTemplate.body}></textarea>
						<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={saveNotificationTemplate}>Save Template</button>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>
