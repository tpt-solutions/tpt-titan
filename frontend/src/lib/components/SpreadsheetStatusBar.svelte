<script>
	import { createEventDispatcher } from 'svelte';
	import {
		statusBarInfo,
		zoomLevel,
		selectedCell,
		sheets,
		activeSheetId
	} from '../stores/spreadsheet-store.js';

	import { getCellId } from '../utils/spreadsheet-utils.js';

	const dispatch = createEventDispatcher();

	function dispatchAction(action, data = {}) {
		dispatch('action', { action, ...data });
	}

	function handleZoomIn() {
		zoomLevel.update(z => Math.min(z + 10, 200));
	}

	function handleZoomOut() {
		zoomLevel.update(z => Math.max(z - 10, 50));
	}

	function handleZoomReset() {
		zoomLevel.set(100);
	}

	function handleZoomChange(event) {
		const value = parseInt(event.target.value);
		if (!isNaN(value) && value >= 50 && value <= 200) {
			zoomLevel.set(value);
		}
	}

	// Format numbers for display
	function formatNumber(num) {
		if (num === null || num === undefined) return '-';
		if (Number.isInteger(num)) return num.toString();
		return num.toFixed(2);
	}
</script>

<!-- Status Bar -->
<div class="bg-gray-100 border-t border-gray-300 px-3 py-1.5 flex items-center justify-between text-sm">
	<!-- Left: Cell Info -->
	<div class="flex items-center space-x-4">
		{#if $selectedCell}
			<span class="text-gray-600 font-medium">
				{getCellId($selectedCell.row, $selectedCell.col)}
			</span>
		{/if}

		{#if $statusBarInfo.selectedCount > 0}
			<div class="flex items-center space-x-3 text-gray-600">
				<span>Count: <strong class="text-gray-800">{$statusBarInfo.selectedCount}</strong></span>
				<span>Sum: <strong class="text-gray-800">{formatNumber($statusBarInfo.sum)}</strong></span>
				<span>Avg: <strong class="text-gray-800">{formatNumber($statusBarInfo.average)}</strong></span>
				{#if $statusBarInfo.min !== null}
					<span>Min: <strong class="text-gray-800">{formatNumber($statusBarInfo.min)}</strong></span>
				{/if}
				{#if $statusBarInfo.max !== null}
					<span>Max: <strong class="text-gray-800">{formatNumber($statusBarInfo.max)}</strong></span>
				{/if}
			</div>
		{/if}
	</div>

	<!-- Center: Sheet Navigation -->
	<div class="flex items-center space-x-2">
		<span class="text-gray-500 text-xs">
			{$sheets.find(s => s.id === $activeSheetId)?.name || 'Sheet1'}
		</span>
	</div>

	<!-- Right: Zoom Controls -->
	<div class="flex items-center space-x-2">
		<button
			class="p-1 rounded hover:bg-gray-200 transition-colors text-gray-600"
			on:click={handleZoomOut}
			title="Zoom Out"
		>
			−
		</button>

		<div class="flex items-center space-x-1">
			<input
				type="number"
				value={$zoomLevel}
				min="50"
				max="200"
				class="w-14 px-1 py-0.5 text-center text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
				on:change={handleZoomChange}
			/>
			<span class="text-gray-600">%</span>
		</div>

		<button
			class="p-1 rounded hover:bg-gray-200 transition-colors text-gray-600"
			on:click={handleZoomIn}
			title="Zoom In"
		>
			+
		</button>

		<button
			class="px-2 py-1 text-xs rounded hover:bg-gray-200 transition-colors text-gray-600 ml-2"
			on:click={handleZoomReset}
			title="Reset Zoom"
		>
			100%
		</button>
	</div>
</div>
