<script>
	import { onMount } from 'svelte';
	import {
		getDocumentExportFormats,
		getDOCXTemplates,
		getDOCXFeatures,
		exportDocument,
		convertDocument,
		batchExportDocuments,
		getDocumentStatistics,
		validateDocumentContent,
		formatApiError
	} from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let formats = [];
	let docxTemplates = [];
	let docxFeatures = null;
	let error = null;
	let loading = false;

	let title = '';
	let content = '# My Document\n\nThis is a paragraph of sample content.\n\n## Section\n\nAnother paragraph.';
	let format = 'docx';
	let docId = '';
	let stats = null;

	onMount(load);

	async function load() {
		loading = true;
		error = null;
		try {
			const [f, t, ft] = await Promise.all([
				getDocumentExportFormats(),
				getDOCXTemplates(),
				getDOCXFeatures()
			]);
			formats = Object.entries(f.formats || {}).map(([k, v]) => ({ key: k, ...v }));
			docxTemplates = t.templates || [];
			docxFeatures = ft;
		} catch (err) {
			error = formatApiError(err);
		} finally {
			loading = false;
		}
	}

	async function doExport() {
		error = null;
		try {
			if (docId.trim()) {
				const blob = await exportDocument(docId.trim(), { content, format, title });
				downloadBlob(blob, `${title || 'document'}.${format}`);
			} else {
				const blob = await exportDocument('0', { content, format, title });
				downloadBlob(blob, `${title || 'document'}.${format}`);
			}
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function doConvert() {
		error = null;
		try {
			const blob = await convertDocument({ content, from_format: 'markdown', to_format: format, title });
			downloadBlob(blob, `converted.${format}`);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function doBatch() {
		error = null;
		try {
			const lines = content.split('\n\n').map((c, i) => ({ content: c, title: `${title || 'doc'}-${i + 1}` }));
			const r = await batchExportDocuments({ documents: lines, format });
			alert(`Batch exported ${r.count || 0} documents as ${format}`);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function showStats() {
		error = null;
		try {
			stats = await getDocumentStatistics(content);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function doValidate() {
		error = null;
		try {
			const r = await validateDocumentContent(content, format);
			stats = { ...(stats || {}), valid: r.valid, validation_errors: r.errors };
		} catch (err) {
			error = formatApiError(err);
		}
	}

	function downloadBlob(blob, name) {
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = name;
		a.click();
		URL.revokeObjectURL(url);
	}
</script>

<svelte:head>
	<title>Export - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<div class="mb-6">
		<h1 class="text-3xl font-bold text-gray-900">Document Export</h1>
		<p class="text-sm text-gray-500 mt-1">Convert and export documents to multiple formats</p>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 border border-red-200 rounded-lg px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	{#if loading}
		<div class="flex items-center justify-center py-16 text-gray-400"><span>Loading formats…</span></div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div class="space-y-4">
				<div class="bg-white border border-gray-200 rounded-lg p-5 space-y-3">
					<label class="block text-sm">Title
						<input bind:value={title} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
					</label>
					<label class="block text-sm">Document ID (optional)
						<input bind:value={docId} placeholder="Leave blank to export inline content" class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
					</label>
					<label class="block text-sm">Content
						<textarea bind:value={content} rows="10" class="mt-1 w-full px-3 py-2 border border-gray-300 rounded font-mono text-sm"></textarea>
					</label>
					<label class="block text-sm">Format
						<select bind:value={format} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded">
							{#each formats as f}<option value={f.key}>{f.name}</option>{/each}
						</select>
					</label>
					<div class="flex flex-wrap gap-2">
						<button on:click={doExport} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Export</button>
						<button on:click={doConvert} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm">Convert</button>
						<button on:click={doBatch} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm">Batch</button>
					</div>
					<div class="flex flex-wrap gap-2">
						<button on:click={showStats} class="text-xs text-blue-600 hover:underline">Statistics</button>
						<button on:click={doValidate} class="text-xs text-blue-600 hover:underline">Validate</button>
					</div>
					{#if stats}
						<div class="bg-gray-50 border border-gray-200 rounded p-3 text-sm space-y-1">
							<div>Words: {stats.word_count ?? '—'} · Characters: {stats.character_count ?? '—'}</div>
							<div>Headings: {stats.heading_count ?? '—'} · Paragraphs: {stats.paragraph_count ?? '—'}</div>
							{#if stats.reading_time != null}<div>Reading time: {stats.reading_time} min</div>{/if}
							{#if stats.valid != null}
								<div class={stats.valid ? 'text-green-700' : 'text-red-700'}>
									{stats.valid ? 'Valid' : 'Invalid'}: {(stats.validation_errors || []).join(', ')}
								</div>
							{/if}
						</div>
					{/if}
				</div>
			</div>

			<div class="space-y-4">
				<div class="bg-white border border-gray-200 rounded-lg p-5">
					<h2 class="font-semibold text-gray-900 mb-2">DOCX Templates</h2>
					<ul class="text-sm divide-y divide-gray-100">
						{#each docxTemplates as t}
							<li class="py-2"><span class="font-medium">{t.name}</span> <span class="text-gray-400">— {t.description}</span></li>
						{/each}
					</ul>
				</div>
				{#if docxFeatures}
					<div class="bg-white border border-gray-200 rounded-lg p-5">
						<h2 class="font-semibold text-gray-900 mb-2">DOCX Features</h2>
						<pre class="text-xs text-gray-600 whitespace-pre-wrap">{JSON.stringify(docxFeatures, null, 2)}</pre>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
