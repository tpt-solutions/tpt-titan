<script>
	import Spreadsheet from '$lib/components/Spreadsheet.svelte';
	import { onMount } from 'svelte';

	let mode = 'simple'; // 'simple' or 'advanced'

	onMount(() => {
		// Initialize spreadsheet
		console.log('Spreadsheet page loaded');
	});
</script>

<svelte:head>
	<title>Spreadsheet - TPT Titan</title>
</svelte:head>

<div class="h-screen flex flex-col bg-white">
	<!-- Header -->
	<header class="flex items-center justify-between px-6 py-3 border-b border-gray-200 bg-gray-50">
		<div class="flex items-center space-x-4">
			<h1 class="text-xl font-semibold text-gray-900">Spreadsheet</h1>
			<div class="flex items-center space-x-2">
				<button
					class="px-3 py-1 text-sm rounded-md transition-colors {mode === 'simple' ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:bg-gray-100'}"
					on:click={() => mode = 'simple'}
				>
					Simple
				</button>
				<button
					class="px-3 py-1 text-sm rounded-md transition-colors {mode === 'advanced' ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:bg-gray-100'}"
					on:click={() => mode = 'advanced'}
				>
					Advanced
				</button>
			</div>
		</div>

		<div class="flex items-center space-x-2">
			<button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
				Save
			</button>
			<button class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors">
				Export
			</button>
		</div>
	</header>

	<!-- Toolbar -->
	<div class="flex items-center space-x-2 px-6 py-2 border-b border-gray-200 bg-white">
		<button class="p-2 hover:bg-gray-100 rounded" title="Bold">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 4h8a4 4 0 014 4v8a4 4 0 01-4 4H6a4 4 0 01-4-4V8a4 4 0 014-4z"></path>
			</svg>
		</button>
		<button class="p-2 hover:bg-gray-100 rounded" title="Italic">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
			</svg>
		</button>
		<div class="w-px h-6 bg-gray-300"></div>
		<button class="p-2 hover:bg-gray-100 rounded text-green-600" title="Insert Chart">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
			</svg>
		</button>
		<div class="flex-1"></div>
		<div class="text-sm text-gray-600">
			Formula: <span class="font-mono bg-gray-100 px-2 py-1 rounded">=SUM(A1:A10)</span>
		</div>
	</div>

	<!-- Main Spreadsheet Area -->
	<div class="flex-1 overflow-hidden">
		<Spreadsheet {mode} />
	</div>
</div>
