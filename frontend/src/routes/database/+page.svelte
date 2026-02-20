<script>
	import { onMount } from 'svelte';
	import { apiGet } from '../../lib/api.js';
	import DatabaseSpreadsheet from '../../lib/components/DatabaseSpreadsheet.svelte';

	// Accept framework-provided props to avoid warnings
	export let data = null;
	export let form = null;
	export let params = null;

	let tables = [];

	let selectedTable = null;
	let loading = true;
	let error = null;

	onMount(async () => {
		try {
			const response = await apiGet('/database/tables');
			tables = response.data.tables;
		} catch (err) {
			error = err.message || 'Failed to load tables';
			console.error('Tables error:', err);
		} finally {
			loading = false;
		}
	});

	function selectTable(tableName) {
		selectedTable = tableName;
	}
</script>

<svelte:head>
	<title>Database Editor - TPT Titan</title>
</svelte:head>

<div class="h-screen flex">
	<!-- Sidebar with table list -->
	<div class="w-64 bg-white border-r border-gray-200 flex flex-col">
		<div class="p-4 border-b border-gray-200">
			<h2 class="text-lg font-semibold text-gray-900">Database Tables</h2>
			<p class="text-sm text-gray-600 mt-1">Select a table to edit</p>
		</div>

		{#if loading}
			<div class="flex-1 flex items-center justify-center">
				<div class="text-center">
					<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 mx-auto"></div>
					<p class="mt-2 text-sm text-gray-600">Loading tables...</p>
				</div>
			</div>
		{:else if error}
			<div class="p-4">
				<div class="bg-red-50 border border-red-200 rounded p-3">
					<p class="text-sm text-red-700">{error}</p>
				</div>
			</div>
		{:else}
			<div class="flex-1 overflow-y-auto">
				{#each tables as table}
					<button
						on:click={() => selectTable(table)}
						class="w-full text-left px-4 py-3 hover:bg-gray-50 border-b border-gray-100 transition-colors
							   {selectedTable === table ? 'bg-blue-50 border-blue-200' : ''}"
					>
						<div class="flex items-center space-x-2">
							<svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
							</svg>
							<span class="text-sm font-medium text-gray-900">{table}</span>
						</div>
					</button>
				{/each}

				{#if tables.length === 0}
					<div class="p-4 text-center text-gray-500">
						<p class="text-sm">No tables found</p>
					</div>
				{/if}
			</div>
		{/if}

		<!-- Footer info -->
		<div class="p-4 border-t border-gray-200 bg-gray-50">
			<div class="text-xs text-gray-600">
				<p><strong>Database Editor</strong></p>
				<p>Spreadsheet-like editing with database constraints</p>
			</div>
		</div>
	</div>

	<!-- Main content area -->
	<div class="flex-1 flex flex-col">
    {#if selectedTable}
      <DatabaseSpreadsheet tableName={selectedTable} />
		{:else}
			<!-- Welcome screen -->
			<div class="flex-1 flex items-center justify-center bg-gray-50">
				<div class="text-center max-w-md">
					<div class="w-24 h-24 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-6">
						<svg class="w-12 h-12 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
						</svg>
					</div>
					<h1 class="text-2xl font-bold text-gray-900 mb-2">Database Editor</h1>
					<p class="text-gray-600 mb-6">
						Edit your database tables with the ease of a spreadsheet while maintaining all database rules and constraints.
					</p>
					<div class="text-left bg-white p-4 rounded-lg shadow-sm border border-gray-200">
						<h3 class="font-semibold text-gray-900 mb-2">Features:</h3>
						<ul class="text-sm text-gray-600 space-y-1">
							<li>✅ Click-to-edit cells like Excel</li>
							<li>✅ Real-time validation against database constraints</li>
							<li>✅ Foreign key relationship indicators</li>
							<li>✅ Data type validation</li>
							<li>✅ Optimistic updates with rollback</li>
							<li>✅ Pagination for large tables</li>
						</ul>
					</div>
					<p class="text-sm text-gray-500 mt-4">
						Select a table from the sidebar to get started
					</p>
				</div>
			</div>
		{/if}
	</div>
</div>
