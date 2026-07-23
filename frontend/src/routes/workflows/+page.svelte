<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import WorkflowBuilder from '$lib/components/WorkflowBuilder.svelte';
	import {
		getWorkflows,
		createWorkflow,
		deleteWorkflow,
		executeWorkflow,
		dryRunWorkflow,
		getWorkflowExecution,
		getWorkflowTemplates,
		createWorkflowFromTemplate,
		getWorkflowInsights,
		optimizeWorkflow,
		generateWorkflowFromPrompt,
		formatApiError
	} from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	/** @type {any[]} */
	let workflows = [];
	/** @type {any[]} */
	let templates = [];
	let view = 'list'; // 'list' | 'builder'
	/** @type {any} */
	let editingId = null;
	/** @type {string | null} */
	let error = null;
	let loading = false;
	/** @type {any} */
	let insights = null;
	/** @type {any} */
	let dryRunResult = null;

	// AI-authored workflow generation (draft only; user reviews before saving).
	let aiPrompt = '';
	/** @type {any} */
	let aiDraft = null;
	/** @type {string | null} */
	let aiGenError = null;
	let aiGenLoading = false;

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

	async function dryRun(wf) {
		error = null;
		dryRunResult = null;
		try {
			const res = await dryRunWorkflow(wf.id);
			const exec = await pollExecution(res.execution_id || res.id);
			dryRunResult = buildDryRunView(exec);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function generateWorkflow() {
		aiGenError = null;
		aiDraft = null;
		aiGenLoading = true;
		try {
			const res = await generateWorkflowFromPrompt(aiPrompt);
			aiDraft = res.draft;
		} catch (err) {
			aiGenError = formatApiError(err);
		} finally {
			aiGenLoading = false;
		}
	}

	// Save the reviewed AI draft as a new workflow (opens the builder pre-filled
	// with the generated canvas for final editing/approval).
	async function useAiDraft() {
		if (!aiDraft) return;
		try {
			await createWorkflow({
				name: aiPrompt.slice(0, 60) || 'AI-generated workflow',
				description: aiPrompt,
				canvas_data: JSON.stringify(aiDraft),
				trigger_type: 'manual',
				is_active: false
			});
			aiDraft = null;
			aiPrompt = '';
			await load();
		} catch (err) {
			aiGenError = formatApiError(err);
		}
	}

	async function pollExecution(executionID) {
		for (let i = 0; i < 20; i++) {
			const r = await getWorkflowExecution(executionID);
			const exec = r.execution;
			if (exec.status !== 'running') return exec;
			await new Promise((resolve) => setTimeout(resolve, 250));
		}
		return (await getWorkflowExecution(executionID)).execution;
	}

	function buildDryRunView(exec) {
		let nodeStates = {};
		let outputData = {};
		try {
			nodeStates = exec.node_states ? JSON.parse(exec.node_states) : {};
		} catch (e) {}
		try {
			outputData = exec.output_data ? JSON.parse(exec.output_data) : {};
		} catch (e) {}
		const previews = Object.entries(nodeStates)
			.filter(([, s]) => s && s.output && s.output.would_execute)
			.map(([nodeID, s]) => ({ nodeID, ...s.output }));
		return {
			status: exec.status,
			isDryRun: exec.is_dry_run,
			previews,
			outputData
		};
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

	/** @param {any} wf */
	async function showOptimization(wf) {
		error = null;
		try {
			/** @type {any} */
			const r = await optimizeWorkflow(wf.id);
			alert(JSON.stringify(r.optimization || r, null, 2));
		} catch (err) {
			error = formatApiError(/** @type {Error} */ (err));
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

	{#if view === 'list'}
		<div class="bg-white border border-indigo-200 rounded-lg p-4 mb-6">
			<h2 class="font-semibold text-gray-900 mb-1">Generate a workflow with AI</h2>
			<p class="text-xs text-gray-500 mb-3">Describe what you want in plain language. The AI returns a draft workflow you review before saving — nothing is created until you confirm.</p>
			<div class="flex gap-2">
				<textarea bind:value={aiPrompt} rows="2" placeholder="e.g. When a form is submitted with priority high, create a task and email the team" class="flex-1 px-3 py-2 border border-gray-300 rounded text-sm"></textarea>
				<button on:click={generateWorkflow} disabled={aiGenLoading || !aiPrompt.trim()} class="bg-indigo-600 hover:bg-indigo-700 disabled:opacity-40 text-white px-4 py-2 rounded text-sm self-start">
					{aiGenLoading ? 'Generating…' : 'Generate'}
				</button>
			</div>
			{#if aiGenError}
				<p class="text-sm text-red-600 mt-2">{aiGenError}</p>
			{/if}
			{#if aiDraft}
				<div class="mt-4 border border-gray-200 rounded p-3">
					<div class="flex items-center justify-between mb-2">
						<p class="text-sm font-medium text-gray-800">Draft ({aiDraft.nodes?.length ?? 0} nodes, {aiDraft.connections?.length ?? 0} connections) — review before saving</p>
						<button on:click={useAiDraft} class="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-xs">Use this draft</button>
					</div>
					<pre class="text-xs bg-gray-50 rounded p-2 overflow-x-auto max-h-72">{JSON.stringify(aiDraft, null, 2)}</pre>
				</div>
			{/if}
		</div>
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
						<button on:click={() => dryRun(wf)} class="text-xs text-amber-600 hover:underline">Dry run</button>
						<button on:click={() => showOptimization(wf)} class="text-xs text-purple-600 hover:underline">Optimize</button>
						<button on:click={() => remove(wf)} class="text-xs text-red-600 hover:underline">Delete</button>
						</div>
					{:else}
						<div class="text-center text-gray-400 py-8 col-span-2">No workflows yet. Click “New workflow” to build one.</div>
					{/each}
				</div>
			</div>

			{#if dryRunResult}
				<div class="md:col-span-2 bg-white border border-amber-200 rounded-lg p-4 mt-4">
					<div class="flex items-center justify-between mb-2">
						<h2 class="font-semibold text-gray-900">Dry run result</h2>
						<span class="inline-flex px-2 py-0.5 rounded-full text-xs {dryRunResult.status === 'completed' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}">{dryRunResult.status}</span>
					</div>
					<p class="text-xs text-gray-400 mb-3">No real side effects were performed — these are previews of what each action node would do.</p>
					{#if dryRunResult.previews.length}
						<div class="space-y-2">
							{#each dryRunResult.previews as p}
								<div class="border border-gray-200 rounded p-3 text-sm">
									<p class="font-medium text-gray-800">{p.node_name || p.nodeID}</p>
									<p class="text-xs text-gray-500">{p.connector} / {p.action}</p>
									<pre class="mt-1 text-xs bg-gray-50 rounded p-2 overflow-x-auto">{JSON.stringify(p.with_config, null, 2)}</pre>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-sm text-gray-500">No action nodes would have executed (no matching branch, or workflow is purely structural).</p>
					{/if}
				</div>
			{/if}

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
