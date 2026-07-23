<script>
// @ts-nocheck
	import Spreadsheet from '$lib/components/Spreadsheet.svelte';
	import TemplateSelector from '$lib/components/TemplateSelector.svelte';
	import SpreadsheetCharts from '$lib/components/SpreadsheetCharts.svelte';
	import SpreadsheetExcel from '$lib/components/SpreadsheetExcel.svelte';
	import SpreadsheetCollab from '$lib/components/SpreadsheetCollab.svelte';
	import { onMount } from 'svelte';

	// Accept framework-provided props to avoid warnings
	export const data = null;
	export const form = null;
	export const params = null;



	let mode = 'simple'; // 'simple' or 'advanced'
	let showTemplates = true; // Show template selector initially
	let selectedTemplate = null;
	let showAppMenu = false;
	let appMenuElement = null;
	let showSidePanel = false;
	let sidePanelTab = 'charts';
	let spreadsheetId = null;

	onMount(() => {
		// Initialize spreadsheet
		console.log('Spreadsheet page loaded');

		// Stable per-session spreadsheet id so charts/excel/collab endpoints have a target
		spreadsheetId = localStorage.getItem('tpt_spreadsheet_id') || (crypto.randomUUID ? crypto.randomUUID() : 'ss-' + Date.now());
		localStorage.setItem('tpt_spreadsheet_id', spreadsheetId);

		// Add click outside listener for app menu
		document.addEventListener('click', handleClickOutside);
	});

	function handleClickOutside(event) {
		if (showAppMenu && appMenuElement && !appMenuElement.contains(event.target) && !event.target.closest('.apps-button')) {
			showAppMenu = false;
		}
	}

	function handleTemplateSelect(event) {
		selectedTemplate = event.detail.template;
		showTemplates = false;
		console.log('Template selected:', selectedTemplate.name);
		console.log('Template data length:', selectedTemplate.data ? selectedTemplate.data.length : 'no data');
		console.log('Template styles:', selectedTemplate.styles ? 'has styles' : 'no styles');
	}

	function handleCreateBlank() {
		selectedTemplate = null;
		showTemplates = false;
		console.log('Blank spreadsheet created');
	}

	function handleTemplatePreview(event) {
		// Could show a preview modal here
		console.log('Preview template:', event.detail.template);
	}

	function toggleSidePanel(tab) {
		if (sidePanelTab === tab && showSidePanel) {
			showSidePanel = false;
		} else {
			sidePanelTab = tab;
			showSidePanel = true;
		}
	}

	// Handlers dispatched from the Spreadsheet component
	function handleExportExcel() { toggleSidePanel('excel'); }
	function handleImport() { toggleSidePanel('excel'); }
	function handleShare() { toggleSidePanel('collab'); }
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

		<!-- Navigation Menu -->
		<div class="flex items-center space-x-2">
			<!-- App Navigation -->
			<div class="relative">
				<button
					class="apps-button px-3 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors flex items-center space-x-2"
					on:click={() => showAppMenu = !showAppMenu}
				>
					<span>📱 Apps</span>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
					</svg>
				</button>

				<!-- App Menu Dropdown -->
				{#if showAppMenu}
					<div class="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-lg border border-gray-200 z-50" bind:this={appMenuElement}>
						<div class="py-2">
							<a href="/" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
								<span class="mr-3">🏠</span>
								<span>Home</span>
							</a>
							<div class="border-t border-gray-100 my-1"></div>
							<a href="/spreadsheet" class="flex items-center px-4 py-2 text-sm text-blue-600 bg-blue-50 transition-colors">
								<span class="mr-3">📊</span>
								<span>Spreadsheet</span>
							</a>
							<a href="/forms" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
								<span class="mr-3">📋</span>
								<span>Forms</span>
							</a>
							<a href="/editor" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
								<span class="mr-3">📝</span>
								<span>Text Editor</span>
							</a>
							<a href="/tasks" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
								<span class="mr-3">✅</span>
								<span>Tasks</span>
							</a>
							<div class="border-t border-gray-100 my-1"></div>
							<a href="/contacts" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
								<span class="mr-3">👥</span>
								<span>Contacts</span>
							</a>
							<a href="/calendar" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
								<span class="mr-3">📅</span>
								<span>Calendar</span>
							</a>
							<a href="/email" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
								<span class="mr-3">📧</span>
								<span>Email</span>
							</a>
						</div>
					</div>
				{/if}
			</div>

			<!-- Template Button -->
			<button
				class="px-3 py-2 bg-green-100 text-green-700 rounded-lg hover:bg-green-200 transition-colors flex items-center space-x-2"
				on:click={() => showTemplates = true}
			>
				<span>📋</span>
				<span>Templates</span>
			</button>

		<!-- Save & Export -->
		<button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
			Save
		</button>
		<button class="px-3 py-2 bg-amber-100 text-amber-700 rounded-lg hover:bg-amber-200 transition-colors" on:click={() => toggleSidePanel('charts')}>
			Charts
		</button>
		<button class="px-3 py-2 bg-emerald-100 text-emerald-700 rounded-lg hover:bg-emerald-200 transition-colors" on:click={() => toggleSidePanel('excel')}>
			Excel
		</button>
		<button class="px-3 py-2 bg-violet-100 text-violet-700 rounded-lg hover:bg-violet-200 transition-colors" on:click={() => toggleSidePanel('collab')}>
			Collaborate
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
		<button class="p-2 hover:bg-gray-100 rounded text-green-600" title="Insert Chart" on:click={() => toggleSidePanel('charts')}>
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
	<div class="flex-1 overflow-hidden flex">
		<div class="flex-1 overflow-hidden">
			<Spreadsheet {mode} {selectedTemplate}
				on:exportExcel={handleExportExcel}
				on:import={handleImport}
				on:share={handleShare}
			/>
		</div>
		{#if showSidePanel}
			<aside class="w-96 border-l border-gray-200 bg-white overflow-auto">
				<div class="flex items-center justify-between px-4 py-2 border-b border-gray-200">
					<div class="flex space-x-1 text-sm">
						{#each [['charts','Charts'],['excel','Excel'],['collab','Collab']] as [key,label]}
							<button class="px-2 py-1 rounded {sidePanelTab === key ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:bg-gray-100'}" on:click={() => sidePanelTab = key}>{label}</button>
						{/each}
					</div>
					<button class="text-gray-400 hover:text-gray-700" on:click={() => showSidePanel = false}>✕</button>
				</div>
				{#if sidePanelTab === 'charts'}
					<SpreadsheetCharts spreadsheetId={spreadsheetId} selectedRange="A1:B5" />
				{:else if sidePanelTab === 'excel'}
					<SpreadsheetExcel spreadsheetId={spreadsheetId} />
				{:else}
					<SpreadsheetCollab spreadsheetId={spreadsheetId} />
				{/if}
			</aside>
		{/if}
	</div>
</div>

<!-- Template Selector -->
<TemplateSelector
	{showTemplates}
	on:select={handleTemplateSelect}
	on:blank={handleCreateBlank}
	on:preview={handleTemplatePreview}
/>
