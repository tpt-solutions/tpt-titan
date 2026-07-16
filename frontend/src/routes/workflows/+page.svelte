<script>
	import { onMount } from 'svelte';
	import WorkflowBuilder from '$lib/components/WorkflowBuilder.svelte';
	import {
		getWorkflows,
		createWorkflow,
		deleteWorkflow,
		executeWorkflow,
		getWorkflowTemplates,
		createWorkflowFromTemplate,
		getWorkflowInsights,
		optimizeWorkflow,
		formatApiError
	} from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let workflows = [];
	let templates = [];
	let view = 'list'; // 'list' | 'builder'
	let editingId = null;
	let error = null;
	let loading = false;
	let insights = null;

	onMount(load);

	async function load() {
		loading = true;
		error = null;
		try {
			const [w, t] = await Promise.all([getWorkflows(), getWorkflowTemplates()]);
			workflows = w.workflows || [];
			templates = t.templates || [];
		} catch (err) {
			error = formatApiError(err);
		} finally {
			loading = false;
		}
	}

	function createNew() {
		editingId = null;
		view = 'builder';
	}

	function editWorkflow(wf) {
		editingId = wf.id;
		view = 'builder';
	}

	function backToList() {
		view = 'list';
		editingId = null;
		load();
	}

	function onSaved(e) {
		backToList();
	}

	function onError(e) {
		error = e.detail?.message || 'Workflow error';
	}

	async function run(wf) {
		error = null;
		try {
			await executeWorkflow(wf.id);
			alert('Workflow execution started');
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function remove(wf) {
		if (!confirm(`Delete workflow "${wf.name}"?`)) return;
		error = null;
		try {
			await deleteWorkflow(wf.id);
			workflows = workflows.filter(w => w.id !== wf.id);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function useTemplate(tmpl) {
		error = null;
		try {
			await createWorkflowFromTemplate(tmpl.id);
			await load();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function loadInsights() {
		error = null;
		try {
			insights = await getWorkflowInsights();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function showOptimization(wf) {
		error = null;
		try {
			const r = await optimizeWorkflow(wf.id);
			alert(JSON.stringify(r.optimization || r, null, 2));
		} catch (err) {
			error = formatApiError(err);
		}
	}
</script>

<svelte:head>
	<title>Workflows - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-6xl">
	<div class="flex items-center justify-between mb-6">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Workflows</h1>
			<p class="text-sm text-gray-500 mt-1">Automate tasks with visual workflows</p>
		</div>
		<div class="flex gap-2">
			<button on:click={loadInsights} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm">Insights</button>
			<button on:click={createNew} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">+ New workflow</button>
		</div>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 border border-red-200 rounded-lg px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	{#if view === 'builder'}
		<div class="bg-white border border-gray-200 rounded-lg p-4 mb-6">
			<WorkflowBuilder workflowId={editingId} on:saved={onSaved} on:error={onError} on:executed={onSaved} />
			<div class="mt-4">
				<button on:click={backToList} class="text-xs text-gray-500 hover:underline">← Back to list</button>
			</div>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="md:col-span-2">
				<h2 class="font-semibold text-gray-900 mb-2">Your workflows</h2>
				<div class="space-y-2">
					{#each workflows as wf}
						<div class="bg-white border border-gray-200 rounded-lg p-4 flex items-center gap-3">
							<div class="flex-1">
								<p class="font-medium text-gray-900">{wf.name}</p>
								<p class="text-xs text-gray-400">{wf.description || wf.category || 'No description'}</p>
								<span class="inline-flex px-2 py-0.5 mt-1 rounded-full text-xs {wf.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'}">
									{wf.is_active ? 'Active' : 'Inactive'}
								</span>
							</div>
							<button on:click={() => editWorkflow(wf)} class="text-xs text-blue-600 hover:underline">Edit</button>
							<button on:click={() => run(wf)} class="text-xs text-green-600 hover:underline">Run</button>
							<button on:click={() => showOptimization(wf)} class="text-xs text-purple-600 hover:underline">Optimize</button>
							<button on:click={() => remove(wf)} class="text-xs text-red-600 hover:underline">Delete</button>
						</div>
					{:else}
						<div class="text-center text-gray-400 py-8 col-span-2">No workflows yet. Click “New workflow” to build one.</div>
					{/each}
				</div>
			</div>

			<div>
				<h2 class="font-semibold text-gray-900 mb-2">Templates</h2>
				<div class="space-y-2">
					{#each templates as t}
						<div class="bg-white border border-gray-200 rounded-lg p-3">
							<p class="font-medium text-gray-900">{t.name}</p>
							<p class="text-xs text-gray-400 mb-2">{t.description}</p>
							<button on:click={() => useTemplate(t)} class="text-xs text-blue-600 hover:underline">Use template</button>
						</div>
					{:else}
						<div class="text-center text-gray-400 py-4">No templates</div>
					{/each}
				</div>

				{#if insights}
					<div class="bg-white border border-gray-200 rounded-lg p-4 mt-4">
						<h2 class="font-semibold text-gray-900 mb-2">Insights</h2>
						<dl class="text-sm space-y-1">
							<div class="flex justify-between"><dt class="text-gray-500">Total workflows</dt><dd>{insights.summary?.total_workflows ?? 0}</dd></div>
							<div class="flex justify-between"><dt class="text-gray-500">Active</dt><dd>{insights.summary?.active_workflows ?? 0}</dd></div>
							<div class="flex justify-between"><dt class="text-gray-500">Total executions</dt><dd>{insights.summary?.total_executions ?? 0}</dd></div>
							<div class="flex justify-between"><dt class="text-gray-500">Success rate</dt><dd>{insights.summary?.success_rate?.toFixed(1) ?? 0}%</dd></div>
						</dl>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
