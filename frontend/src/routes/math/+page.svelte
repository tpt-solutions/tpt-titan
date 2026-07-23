<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import {
		validateExpression,
		optimizeExpression,
		convertExpression,
		getMathFunctions,
		getMathSymbols,
		getMathConstants,
		getMathTheorems,
		getEquationTemplates,
		searchEquations,
		getEquationTemplateCategories,
		saveEquationTemplate,
		generateEquationImage,
		exportEquation,
		batchExportEquations,
		recognizeEquationFromImage,
		formatApiError
	} from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let tab = 'expressions';
	let error = null;

	// Expressions
	let expr = 'x^2 + 2*x + 1';
	let exprResult = '';
	let exprOptimized = '';
	let convertFrom = 'text';
	let convertTo = 'latex';

	// Reference
	let functions = [];
	let symbols = [];
	let constants = [];
	let theorems = [];
	let refCategory = '';

	// Templates
	let templates = [];
	let templateCategories = [];
	let templateSearch = '';
	let newTemplateName = '';
	let newTemplateLatex = '';

	// Canvas image
	let canvasLatex = 'E = mc^2';
	let canvasFormat = 'svg';

	// Export
	let exportLatex = '\\int_{0}^{1} x^2 \\, dx';
	let exportFormat = 'latex';
	let exportPreview = '';

	// Recognition
	let recImage = null;
	let recResult = '';

	async function runValidate() {
		error = null;
		try {
			const r = await validateExpression(expr);
			exprResult = (r.valid ? '✅ Valid' : '❌ Invalid') + ' — ' + r.message;
		} catch (e) { error = formatApiError(e); }
	}

	async function runOptimize() {
		error = null;
		try {
			const r = await optimizeExpression(expr);
			exprOptimized = r.optimized;
		} catch (e) { error = formatApiError(e); }
	}

	async function runConvert() {
		error = null;
		try {
			const r = await convertExpression(expr, convertFrom, convertTo);
			exprResult = `${convertFrom} → ${convertTo}: ${r.converted}`;
		} catch (e) { error = formatApiError(e); }
	}

	async function loadReference() {
		error = null;
		try {
			const [f, s, c, t] = await Promise.all([
				getMathFunctions(),
				getMathSymbols(),
				getMathConstants(),
				getMathTheorems()
			]);
			functions = f.functions || [];
			symbols = s.symbols || [];
			constants = c.constants || [];
			theorems = t.theorems || [];
		} catch (e) { error = formatApiError(e); }
	}

	async function loadTemplates() {
		error = null;
		try {
			const [t, c] = await Promise.all([
				getEquationTemplates({ search: templateSearch }),
				getEquationTemplateCategories()
			]);
			templates = t.templates || [];
			templateCategories = c.categories || [];
		} catch (e) { error = formatApiError(e); }
	}

	async function saveTemplate() {
		if (!newTemplateName.trim() || !newTemplateLatex.trim()) return;
		error = null;
		try {
			await saveEquationTemplate({ name: newTemplateName, latex: newTemplateLatex, category: 'custom' });
			newTemplateName = '';
			newTemplateLatex = '';
			await loadTemplates();
		} catch (e) { error = formatApiError(e); }
	}

	async function doCanvasImage() {
		error = null;
		try {
			const blob = await generateEquationImage({ latex: canvasLatex, format: canvasFormat });
			exportPreview = URL.createObjectURL(blob);
		} catch (e) { error = formatApiError(e); }
	}

	async function doExport() {
		error = null;
		try {
			const blob = await exportEquation({ latex: exportLatex, name: 'equation' }, exportFormat);
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `equation.${exportFormat}`;
			a.click();
			URL.revokeObjectURL(url);
		} catch (e) { error = formatApiError(e); }
	}

	async function doBatchExport() {
		error = null;
		try {
			const exprs = expr.split('\n').filter(Boolean).map(latex => ({ latex, name: 'eq' }));
			const r = await batchExportEquations(exprs, exportFormat);
			alert(`Exported ${r.count || 0} equations as ${exportFormat}`);
		} catch (e) { error = formatApiError(e); }
	}

	async function doRecognize() {
		if (!recImage) { error = 'Choose an image first'; return; }
		error = null;
		try {
			const fd = new FormData();
			fd.append('image', recImage);
			fd.append('format', recImage.name.split('.').pop() || 'png');
			const r = await recognizeEquationFromImage(fd);
			recResult = r.latex || r.text || JSON.stringify(r, null, 2);
		} catch (e) { error = formatApiError(e); }
	}

	function changeTab(t) {
		tab = t;
		if (t === 'reference') loadReference();
		if (t === 'templates') loadTemplates();
	}

	onMount(() => { runValidate(); });
</script>

<svelte:head>
	<title>Math - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-5xl">
	<div class="mb-6">
		<h1 class="text-3xl font-bold text-gray-900">Math</h1>
		<p class="text-sm text-gray-500 mt-1">Natural math input, recognition, and export</p>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 border border-red-200 rounded-lg px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	<div class="flex flex-wrap gap-1 border-b border-gray-200 mb-6 overflow-x-auto">
		{#each ['expressions', 'reference', 'templates', 'canvas', 'export', 'recognition'] as t}
			<button
				class="px-3 py-2 text-sm font-medium border-b-2 transition-colors whitespace-nowrap {tab === t ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-600 hover:text-gray-900'}"
				on:click={() => changeTab(t)}
			>
				{t.charAt(0).toUpperCase() + t.slice(1)}
			</button>
		{/each}
	</div>

	{#if tab === 'expressions'}
		<div class="bg-white border border-gray-200 rounded-lg p-6 space-y-4">
			<label class="block text-sm">Expression
				<textarea bind:value={expr} rows="2" class="mt-1 w-full px-3 py-2 border border-gray-300 rounded font-mono"></textarea>
			</label>
			<div class="flex flex-wrap gap-2">
				<button on:click={runValidate} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Validate</button>
				<button on:click={runOptimize} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Optimize</button>
			</div>
			{#if exprResult}<div class="bg-gray-50 border border-gray-200 rounded p-3 text-sm">{exprResult}</div>{/if}
			{#if exprOptimized}<div class="bg-gray-50 border border-gray-200 rounded p-3 text-sm font-mono">Optimized: {exprOptimized}</div>{/if}

			<hr class="border-gray-200" />
			<h2 class="font-semibold text-gray-900">Convert format</h2>
			<div class="grid grid-cols-2 md:grid-cols-4 gap-3">
				<select bind:value={convertFrom} class="px-3 py-2 border border-gray-300 rounded text-sm">
					<option value="text">text</option><option value="latex">latex</option><option value="mathml">mathml</option>
				</select>
				<select bind:value={convertTo} class="px-3 py-2 border border-gray-300 rounded text-sm">
					<option value="latex">latex</option><option value="mathml">mathml</option><option value="svg">svg</option>
				</select>
				<button on:click={runConvert} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm">Convert</button>
			</div>
		</div>
	{:else if tab === 'reference'}
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div class="space-y-4">
				<div class="bg-white border border-gray-200 rounded-lg p-4">
					<h2 class="font-semibold text-gray-900 mb-2">Functions</h2>
					<div class="space-y-1 text-sm max-h-64 overflow-auto">
						{#each functions as f}<div><span class="font-mono font-medium">{f.name}</span> <span class="text-gray-400">— {f.description}</span></div>{/each}
					</div>
				</div>
				<div class="bg-white border border-gray-200 rounded-lg p-4">
					<h2 class="font-semibold text-gray-900 mb-2">Constants</h2>
					<div class="space-y-1 text-sm">
						{#each constants as c}<div><span class="font-mono font-medium">{c.symbol}</span> = {c.value} <span class="text-gray-400">— {c.description}</span></div>{/each}
					</div>
				</div>
			</div>
			<div class="space-y-4">
				<div class="bg-white border border-gray-200 rounded-lg p-4">
					<h2 class="font-semibold text-gray-900 mb-2">Symbols</h2>
					<div class="flex flex-wrap gap-2 text-lg">
						{#each symbols as s}<span class="px-2 py-1 bg-gray-50 rounded" title={s.name}>{s.unicode}</span>{/each}
					</div>
				</div>
				<div class="bg-white border border-gray-200 rounded-lg p-4">
					<h2 class="font-semibold text-gray-900 mb-2">Theorems</h2>
					<div class="space-y-2 text-sm">
						{#each theorems as th}
							<div><span class="font-medium">{th.name}</span> <span class="font-mono text-gray-500">{th.formula}</span></div>
						{/each}
					</div>
				</div>
			</div>
		</div>
	{:else if tab === 'templates'}
		<div class="bg-white border border-gray-200 rounded-lg p-5 mb-4 space-y-3">
			<input bind:value={templateSearch} placeholder="Search templates…" on:keydown={(e) => e.key === 'Enter' && loadTemplates()} class="w-full px-3 py-2 border border-gray-300 rounded text-sm" />
			<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
				<input bind:value={newTemplateName} placeholder="New template name" class="px-3 py-2 border border-gray-300 rounded text-sm" />
				<input bind:value={newTemplateLatex} placeholder="LaTeX (e.g. x^2)" class="px-3 py-2 border border-gray-300 rounded text-sm" />
			</div>
			<button on:click={saveTemplate} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Save template</button>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-3 gap-3">
			{#each templates as t}
				<div class="bg-white border border-gray-200 rounded-lg p-3">
					<p class="font-medium text-gray-900">{t.name}</p>
					<p class="font-mono text-sm text-gray-600 break-all">{t.latex || t.formula || ''}</p>
					{#if t.category}<p class="text-xs text-gray-400 mt-1">{t.category}</p>{/if}
				</div>
			{:else}
				<div class="col-span-3 text-center text-gray-400 py-8">No templates found</div>
			{/each}
		</div>
	{:else if tab === 'canvas'}
		<div class="bg-white border border-gray-200 rounded-lg p-6 space-y-4">
			<label class="block text-sm">LaTeX expression
				<input bind:value={canvasLatex} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded font-mono" />
			</label>
			<div class="flex items-end gap-3">
				<label class="text-sm">Format
					<select bind:value={canvasFormat} class="mt-1 block px-3 py-2 border border-gray-300 rounded">
						<option value="png">png</option><option value="svg">svg</option><option value="pdf">pdf</option>
					</select>
				</label>
				<button on:click={doCanvasImage} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Generate image</button>
			</div>
			{#if exportPreview}
				<div class="border border-gray-200 rounded p-4 bg-gray-50">
					{#if canvasFormat === 'svg'}
						<!-- svg blob rendered as object -->
						<object data={exportPreview} type="image/svg+xml" class="max-w-full h-40"></object>
					{:else}
						<img src={exportPreview} alt="equation" class="max-w-full h-40" />
					{/if}
				</div>
			{/if}
		</div>
	{:else if tab === 'export'}
		<div class="bg-white border border-gray-200 rounded-lg p-6 space-y-4">
			<label class="block text-sm">LaTeX expression
				<input bind:value={exportLatex} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded font-mono" />
			</label>
			<div class="flex items-end gap-3">
				<label class="text-sm">Format
					<select bind:value={exportFormat} class="mt-1 block px-3 py-2 border border-gray-300 rounded">
						<option value="latex">latex</option><option value="mathml">mathml</option><option value="svg">svg</option><option value="png">png</option><option value="pdf">pdf</option>
					</select>
				</label>
				<button on:click={doExport} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Export</button>
				<button on:click={doBatchExport} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm">Batch (from expressions)</button>
			</div>
		</div>
	{:else if tab === 'recognition'}
		<div class="bg-white border border-gray-200 rounded-lg p-6 space-y-4">
			<label class="block text-sm">Equation image
				<input type="file" accept="image/*" on:change={(e) => (recImage = e.target.files[0])} class="mt-1 w-full text-sm" />
			</label>
			<button on:click={doRecognize} disabled={!recImage} class="bg-blue-600 hover:bg-blue-700 disabled:opacity-40 text-white px-4 py-2 rounded text-sm">Recognize</button>
			{#if recResult}
				<div class="bg-gray-50 border border-gray-200 rounded p-3 text-sm font-mono break-all">{recResult}</div>
			{/if}
		</div>
	{/if}
</div>
