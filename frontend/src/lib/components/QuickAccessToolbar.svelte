<script>
// @ts-nocheck
	import { createEventDispatcher } from 'svelte';
	import {
		canUndo,
		canRedo,
		hasUnsavedChanges,
		workbookName
	} from '../stores/spreadsheet-store.js';

	const dispatch = createEventDispatcher();

	function dispatchAction(action, data = {}) {
		dispatch('action', { action, ...data });
	}

	function handleNew() {
		dispatchAction('newSpreadsheet');
	}

	function handleOpen() {
		dispatchAction('openSpreadsheet');
	}

	function handleSave() {
		dispatchAction('save');
	}

	function handleUndo() {
		dispatchAction('undo');
	}

	function handleRedo() {
		dispatchAction('redo');
	}
</script>

<!-- Quick Access Toolbar -->
<div class="bg-gray-900 text-white px-3 py-1.5 flex items-center space-x-1">
	<!-- App Logo/Title -->
	<div class="flex items-center space-x-2 mr-4">
		<span class="text-lg">📊</span>
		<span class="font-semibold text-sm">TPT Titan</span>
	</div>

	<!-- Quick Actions -->
	<div class="flex items-center space-x-1">
		<!-- New -->
		<button
			class="p-1.5 rounded hover:bg-gray-700 transition-colors"
			on:click={handleNew}
			title="New Spreadsheet (Ctrl+N)"
		>
			📄
		</button>

		<!-- Open -->
		<button
			class="p-1.5 rounded hover:bg-gray-700 transition-colors"
			on:click={handleOpen}
			title="Open Spreadsheet (Ctrl+O)"
		>
			📂
		</button>

		<!-- Save -->
		<button
			class="p-1.5 rounded hover:bg-gray-700 transition-colors relative"
			on:click={handleSave}
			title="Save (Ctrl+S)"
		>
			💾
			{#if $hasUnsavedChanges}
				<span class="absolute -top-0.5 -right-0.5 w-2 h-2 bg-orange-500 rounded-full"></span>
			{/if}
		</button>

		<div class="w-px h-5 bg-gray-600 mx-2"></div>

		<!-- Undo -->
		<button
			class="p-1.5 rounded hover:bg-gray-700 transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
			on:click={handleUndo}
			disabled={!$canUndo}
			title="Undo (Ctrl+Z)"
		>
			↶
		</button>

		<!-- Redo -->
		<button
			class="p-1.5 rounded hover:bg-gray-700 transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
			on:click={handleRedo}
			disabled={!$canRedo}
			title="Redo (Ctrl+Y)"
		>
			↷
		</button>
	</div>

	<!-- Workbook Name (centered) -->
	<div class="flex-1 text-center">
		<span class="text-sm text-gray-300 truncate max-w-md inline-block">
			{$workbookName}
			{#if $hasUnsavedChanges}
				<span class="text-orange-400">●</span>
			{/if}
		</span>
	</div>

	<!-- Right side actions -->
	<div class="flex items-center space-x-1">
		<button
			class="p-1.5 rounded hover:bg-gray-700 transition-colors"
			on:click={() => dispatchAction('share')}
			title="Share"
		>
			👥
		</button>
	</div>
</div>
