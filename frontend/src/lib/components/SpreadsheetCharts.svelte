<script>
	import { getChartSuggestions, createChart, getSpreadsheetCharts } from '$lib/api.js';

	export let selectedRange = 'A1:B5';
	export let spreadsheetId = null;

	let charts = [];
	let suggestions = [];
	let loading = false;
	let error = '';
	let message = '';

	let form = {
		chart_type: 'bar',
		data_range: selectedRange,
		title: '',
		x_axis_label: '',
		y_axis_label: '',
		data: {}
	};

	async function loadCharts() {
		if (!spreadsheetId) return;
		try { const res = await getSpreadsheetCharts(spreadsheetId); charts = res.charts || []; }
		catch (e) { error = e.message; }
	}

	async function suggest() {
		loading = true; error = ''; message = '';
		try {
			const res = await getChartSuggestions({ range: form.data_range, data: form.data, data_types: {} });
			suggestions = res.suggestions || [];
			message = suggestions.length ? 'Suggestions ready' : 'No suggestions';
		} catch (e) { error = e.message; } finally { loading = false; }
	}

	async function create() {
		loading = true; error = ''; message = '';
		try {
			const payload = { ...form, spreadsheet_id: spreadsheetId || '00000000-0000-0000-0000-000000000000' };
			await createChart(payload);
			message = 'Chart created';
			form = { chart_type: 'bar', data_range: selectedRange, title: '', x_axis_label: '', y_axis_label: '', data: {} };
			await loadCharts();
		} catch (e) { error = e.message; } finally { loading = false; }
	}

	$: if (spreadsheetId) loadCharts();
</script>

<div class="p-4 space-y-4">
	{#if message}<div class="text-sm text-green-700 bg-green-50 rounded px-3 py-2">{message}</div>{/if}
	{#if error}<div class="text-sm text-red-700 bg-red-50 rounded px-3 py-2">{error}</div>{/if}

	<div class="border rounded-lg p-4">
		<h3 class="font-semibold mb-3">Create Chart</h3>
		<div class="grid grid-cols-2 gap-3 text-sm">
			<select class="border rounded px-2 py-1" bind:value={form.chart_type}>
				<option value="bar">Bar</option>
				<option value="line">Line</option>
				<option value="pie">Pie</option>
				<option value="scatter">Scatter</option>
				<option value="area">Area</option>
			</select>
			<input class="border rounded px-2 py-1" placeholder="Data range (e.g. A1:B5)" bind:value={form.data_range} />
			<input class="border rounded px-2 py-1" placeholder="Title" bind:value={form.title} />
			<input class="border rounded px-2 py-1" placeholder="X axis label" bind:value={form.x_axis_label} />
			<input class="border rounded px-2 py-1" placeholder="Y axis label" bind:value={form.y_axis_label} />
		</div>
		<div class="flex space-x-2 mt-3">
			<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={create} disabled={loading}>Create Chart</button>
			<button class="px-3 py-1 bg-gray-600 text-white rounded text-sm" on:click={suggest} disabled={loading}>Suggest</button>
		</div>
	</div>

	{#if suggestions.length}
		<div class="border rounded-lg p-4">
			<h3 class="font-semibold mb-2">Suggestions</h3>
			<ul class="text-sm space-y-1">
				{#each suggestions as s}
					<li class="border-b pb-1">{s.chart_type || s.type}: {s.title || s.reason || 'suggested'}</li>
				{/each}
			</ul>
		</div>
	{/if}

	<div class="border rounded-lg p-4">
		<h3 class="font-semibold mb-2">Existing Charts</h3>
		{#if charts.length === 0}
			<p class="text-sm text-gray-500">No charts yet.</p>
		{:else}
			<ul class="text-sm space-y-1">
				{#each charts as c}
					<li class="border-b pb-1">{c.title} <span class="text-gray-400">({c.type} · {c.data_range})</span></li>
				{/each}
			</ul>
		{/if}
	</div>
</div>
