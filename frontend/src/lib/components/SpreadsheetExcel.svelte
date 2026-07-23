<script>
// @ts-nocheck
	import { exportSpreadsheetToExcel, importExcelToSpreadsheet } from '$lib/api.js';

	export let spreadsheetId = null;

	let importFile = null;
	let loading = false;
	let error = '';
	let message = '';
	let formats = [];

	async function handleExport() {
		loading = true; error = ''; message = '';
		try {
			const blob = await exportSpreadsheetToExcel(spreadsheetId || '00000000-0000-0000-0000-000000000000', {});
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = 'spreadsheet.xlsx';
			a.click();
			URL.revokeObjectURL(url);
			message = 'Export downloaded';
		} catch (e) { error = e.message; } finally { loading = false; }
	}

	async function handleImport() {
		if (!importFile) { error = 'Choose a file first'; return; }
		loading = true; error = ''; message = '';
		try {
			const fd = new FormData();
			fd.append('file', importFile);
			const res = await importExcelToSpreadsheet(fd);
			message = 'Imported ' + ((res.sheets || []).length || 0) + ' sheet(s)';
		} catch (e) { error = e.message; } finally { loading = false; }
	}

	function onFileSelected(e) {
		importFile = e.target.files && e.target.files[0];
	}
</script>

<div class="p-4 space-y-4">
	{#if message}<div class="text-sm text-green-700 bg-green-50 rounded px-3 py-2">{message}</div>{/if}
	{#if error}<div class="text-sm text-red-700 bg-red-50 rounded px-3 py-2">{error}</div>{/if}

	<div class="border rounded-lg p-4">
		<h3 class="font-semibold mb-3">Excel Export</h3>
		<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={handleExport} disabled={loading}>Export to Excel (.xlsx)</button>
	</div>

	<div class="border rounded-lg p-4">
		<h3 class="font-semibold mb-3">Excel Import</h3>
		<input type="file" accept=".xlsx,.xls" on:change={onFileSelected} class="text-sm mb-2" />
		<button class="px-3 py-1 bg-green-600 text-white rounded text-sm" on:click={handleImport} disabled={loading}>Import from Excel</button>
	</div>

	{#if formats.length}
		<div class="border rounded-lg p-4 text-sm">
			<h3 class="font-semibold mb-2">Supported formats</h3>
			<pre class="text-xs bg-gray-50 p-2 rounded">{JSON.stringify(formats, null, 2)}</pre>
		</div>
	{/if}
</div>
