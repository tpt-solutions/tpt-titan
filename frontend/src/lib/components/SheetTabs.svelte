<script>
	import { createEventDispatcher } from 'svelte';
	import {
		sheets,
		activeSheetId,
		sheetCounter,
		hasUnsavedChanges
	} from '../stores/spreadsheet-store.js';

	const dispatch = createEventDispatcher();

	let editingSheetId = null;
	let editingSheetName = '';

	function dispatchAction(action, data = {}) {
		dispatch('action', { action, ...data });
	}

	function handleAddSheet() {
		sheetCounter.update(c => c + 1);
		const newCounter = $sheetCounter;
		const newSheet = {
			id: `sheet${Date.now()}`,
			name: `Sheet${newCounter + 1}`,
			data: Array(100).fill().map(() => Array(26).fill(''))
		};

		sheets.update(s => [...s, newSheet]);
		activeSheetId.set(newSheet.id);
		hasUnsavedChanges.set(true);
	}

	function handleSwitchSheet(sheetId) {
		activeSheetId.set(sheetId);
		dispatchAction('switchSheet', { sheetId });
	}

	function handleDeleteSheet(sheetId, event) {
		event.stopPropagation();

		if ($sheets.length <= 1) {
			alert('Cannot delete the last sheet. At least one sheet must remain.');
			return;
		}

		const sheet = $sheets.find(s => s.id === sheetId);
		const confirmed = confirm(`Are you sure you want to delete "${sheet.name}"?`);
		if (!confirmed) return;

		sheets.update(s => s.filter(sheet => sheet.id !== sheetId));

		// If we deleted the active sheet, switch to the first available
		if ($activeSheetId === sheetId) {
			const remainingSheets = $sheets.filter(s => s.id !== sheetId);
			if (remainingSheets.length > 0) {
				activeSheetId.set(remainingSheets[0].id);
			}
		}

		hasUnsavedChanges.set(true);
	}

	function startRename(sheet, event) {
		event.stopPropagation();
		editingSheetId = sheet.id;
		editingSheetName = sheet.name;
	}

	function handleRename(sheetId) {
		if (editingSheetName.trim()) {
			sheets.update(s => s.map(sheet =>
				sheet.id === sheetId
					? { ...sheet, name: editingSheetName.trim() }
					: sheet
			));
			hasUnsavedChanges.set(true);
		}
		editingSheetId = null;
		editingSheetName = '';
	}

	function handleRenameKeyDown(event, sheetId) {
		if (event.key === 'Enter') {
			handleRename(sheetId);
		} else if (event.key === 'Escape') {
			editingSheetId = null;
			editingSheetName = '';
		}
	}
</script>

<!-- Sheet Tabs -->
<div class="bg-gray-100 border-t border-gray-300 flex items-center">
	<!-- Sheet Tabs -->
	<div class="flex items-center overflow-x-auto">
		{#each $sheets as sheet}
			<div
				class="group relative flex items-center min-w-24 max-w-40 px-3 py-2 border-r border-gray-300 cursor-pointer transition-colors
					{$activeSheetId === sheet.id
						? 'bg-white border-t-2 border-t-blue-500 text-gray-900 font-medium'
						: 'bg-gray-100 hover:bg-gray-200 text-gray-600'}"
				on:click={() => handleSwitchSheet(sheet.id)}
			>
				{#if editingSheetId === sheet.id}
					<input
						type="text"
						bind:value={editingSheetName}
						class="w-full px-1 py-0.5 text-sm border border-blue-500 rounded focus:outline-none"
						on:blur={() => handleRename(sheet.id)}
						on:keydown={(e) => handleRenameKeyDown(e, sheet.id)}
						autofocus
					/>
				{:else}
					<span class="truncate text-sm flex-1">{sheet.name}</span>

					<!-- Sheet Actions (visible on hover) -->
					<div class="flex items-center space-x-1 ml-2 opacity-0 group-hover:opacity-100 transition-opacity">
						<button
							class="p-0.5 rounded hover:bg-gray-300 text-gray-500"
							on:click={(e) => startRename(sheet, e)}
							title="Rename Sheet"
						>
							✏️
						</button>
						{#if $sheets.length > 1}
							<button
								class="p-0.5 rounded hover:bg-red-100 text-red-500"
								on:click={(e) => handleDeleteSheet(sheet.id, e)}
								title="Delete Sheet"
							>
								×
							</button>
						{/if}
					</div>
				{/if}
			</div>
		{/each}
	</div>

	<!-- Add Sheet Button -->
	<button
		class="flex items-center justify-center w-8 h-8 mx-2 rounded hover:bg-gray-200 text-gray-600 transition-colors"
		on:click={handleAddSheet}
		title="Add New Sheet"
	>
		+
	</button>
</div>
