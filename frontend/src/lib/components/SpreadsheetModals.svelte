<script>
	import {
		showFormulaHelp,
		showPivotTableBuilder,
		showDataPrepTools,
		showFormatDialog,
		showRibbonCustomizer,
		showFreezePaneDialog,
		showFindReplace,
		showContextMenu,
		contextMenuPosition,
		contextMenuItems
	} from '../stores/spreadsheet-store.js';

	// Dispatch events to parent
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	function dispatchAction(action, data = {}) {
		dispatch('action', { action, ...data });
	}

	// Context menu action handler
	function executeContextMenuAction(action) {
		dispatchAction(action.action, action.data || {});
		showContextMenu.set(false);
	}

	// Close all modals helper
	function closeAllModals() {
		showFormulaHelp.set(false);
		showPivotTableBuilder.set(false);
		showDataPrepTools.set(false);
		showFormatDialog.set(false);
		showRibbonCustomizer.set(false);
		showFreezePaneDialog.set(false);
		showFindReplace.set(false);
		showContextMenu.set(false);
	}

	// Handle window click to close context menu
	function handleWindowClick(event) {
		// Close context menu
		if ($showContextMenu) {
			const contextMenu = event.target.closest('.fixed.z-50');
			if (!contextMenu) {
				showContextMenu.set(false);
			}
		}
	}

	// Window event listeners
	import { onMount } from 'svelte';
	onMount(() => {
		window.addEventListener('click', handleWindowClick);
		return () => {
			window.removeEventListener('click', handleWindowClick);
		};
	});
</script>

<!-- Formula Help Modal -->
{#if $showFormulaHelp}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">🔢 Advanced Spreadsheet Functions</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showFormulaHelp.set(false)}
				>
					×
				</button>
			</div>

			<div class="space-y-6">
				<!-- Introduction -->
				<div class="bg-blue-50 p-4 rounded-lg">
					<h4 class="text-lg font-medium text-blue-900 mb-2">⚡ Powerful Mathematical Engine</h4>
					<p class="text-blue-800">
						TPT Titan includes a comprehensive mathematical function library with 50+ functions,
						advanced Excel compatibility, real-time collaboration, and automatic chart generation.
					</p>
				</div>

				<!-- Function Categories -->
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					<!-- Mathematical Functions -->
					<div>
						<h4 class="text-lg font-medium text-gray-900 mb-3">🧮 Mathematical Functions</h4>
						<div class="space-y-2 text-sm">
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">SUM(range)</code>
								<span class="text-gray-600">Add numbers</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">AVERAGE(range)</code>
								<span class="text-gray-600">Calculate mean</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">SIN(angle)</code>
								<span class="text-gray-600">Sine function</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">COS(angle)</code>
								<span class="text-gray-600">Cosine function</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">SQRT(number)</code>
								<span class="text-gray-600">Square root</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">POWER(base,exp)</code>
								<span class="text-gray-600">Power function</span>
							</div>
						</div>
					</div>

					<!-- Statistical Functions -->
					<div>
						<h4 class="text-lg font-medium text-gray-900 mb-3">📊 Statistical Functions</h4>
						<div class="space-y-2 text-sm">
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">MIN(range)</code>
								<span class="text-gray-600">Minimum value</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">MAX(range)</code>
								<span class="text-gray-600">Maximum value</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">COUNT(range)</code>
								<span class="text-gray-600">Count numbers</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">STDEV(range)</code>
								<span class="text-gray-600">Standard deviation</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">MEDIAN(range)</code>
								<span class="text-gray-600">Median value</span>
							</div>
							<div class="flex justify-between">
								<code class="bg-gray-100 px-2 py-1 rounded">MODE(range)</code>
								<span class="text-gray-600">Most frequent value</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Advanced Features -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">🚀 Advanced Features</h4>
					<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
						<div class="bg-purple-50 p-4 rounded">
							<h5 class="font-medium text-purple-900 mb-2">📈 Charts & Visualization</h5>
							<p class="text-sm text-purple-700">Automatic chart suggestions based on your data patterns. Support for bar, line, pie, scatter, and area charts.</p>
						</div>
						<div class="bg-green-50 p-4 rounded">
							<h5 class="font-medium text-green-900 mb-2">🔄 Real-time Collaboration</h5>
							<p class="text-sm text-green-700">Work with others simultaneously. See changes in real-time with conflict resolution and version control.</p>
						</div>
						<div class="bg-blue-50 p-4 rounded">
							<h5 class="font-medium text-blue-900 mb-2">📊 Excel Compatibility</h5>
							<p class="text-sm text-blue-700">Full .xlsx import/export with formulas, styles, multiple sheets, and formatting preservation.</p>
						</div>
					</div>
				</div>

				<!-- Quick Examples -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">💡 Quick Examples</h4>
					<div class="bg-gray-50 p-4 rounded">
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
							<div>
								<p class="font-medium mb-2">Basic Calculations:</p>
								<code class="block bg-white p-2 rounded mb-1">=SUM(A1:A10)</code>
								<code class="block bg-white p-2 rounded mb-1">=AVERAGE(B1:B20)</code>
								<code class="block bg-white p-2 rounded mb-1">=A1*B1 + C1</code>
							</div>
							<div>
								<p class="font-medium mb-2">Advanced Functions:</p>
								<code class="block bg-white p-2 rounded mb-1">=SIN(RADIANS(30))</code>
								<code class="block bg-white p-2 rounded mb-1">=SQRT(A1^2 + B1^2)</code>
								<code class="block bg-white p-2 rounded mb-1">=IF(A1>100, "High", "Low")</code>
							</div>
						</div>
					</div>
				</div>

				<!-- Export Options -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">📤 Export & Share</h4>
					<p class="text-gray-600 mb-3">
						Export your spreadsheets to various formats and share with collaborators:
					</p>
					<div class="grid grid-cols-2 md:grid-cols-4 gap-3">
						<div class="text-center p-3 bg-green-50 rounded">
							<div class="text-green-600 font-medium">Excel</div>
							<div class="text-xs text-green-600">.xlsx format</div>
						</div>
						<div class="text-center p-3 bg-blue-50 rounded">
							<div class="text-blue-600 font-medium">CSV</div>
							<div class="text-xs text-blue-600">Data export</div>
						</div>
						<div class="text-center p-3 bg-purple-50 rounded">
							<div class="text-purple-600 font-medium">PDF</div>
							<div class="text-xs text-purple-600">Printable format</div>
						</div>
						<div class="text-center p-3 bg-orange-50 rounded">
							<div class="text-orange-600 font-medium">Share</div>
							<div class="text-xs text-orange-600">Collaborate</div>
						</div>
					</div>
				</div>
			</div>

			<div class="mt-6 flex justify-end">
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
					on:click={() => showFormulaHelp.set(false)}
				>
					Got it!
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Pivot Table Builder Modal -->
{#if $showPivotTableBuilder}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-6xl max-h-[90vh] overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">📋 TPT Titan Pivot Table Builder</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showPivotTableBuilder.set(false)}
				>
					×
				</button>
			</div>

			<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
				<!-- Data Source Selection -->
				<div class="lg:col-span-1">
					<h4 class="text-lg font-medium text-gray-900 mb-4">📊 Data Source</h4>
					<div class="space-y-4">
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Select Range</label>
							<input
								type="text"
								placeholder="Click and drag on spreadsheet to select range"
								value="A1:D10"
								class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
								readonly
							/>
							<p class="text-xs text-gray-500 mt-1">💡 Open the pivot table builder and select cells visually</p>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Data has headers</label>
							<input type="checkbox" checked class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500" />
						</div>
					</div>

					<h4 class="text-lg font-medium text-gray-900 mb-4 mt-6">🎯 Pivot Configuration</h4>
					<div class="space-y-4">
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Rows</label>
							<select class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500">
								<option>Select field...</option>
								<option>Product</option>
								<option>Category</option>
								<option>Region</option>
							</select>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Columns</label>
							<select class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500">
								<option>Select field...</option>
								<option>Month</option>
								<option>Quarter</option>
								<option>Year</option>
							</select>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-2">Values</label>
							<select class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500">
								<option>Select field...</option>
								<option>Sales</option>
								<option>Revenue</option>
								<option>Quantity</option>
							</select>
						</div>
					</div>
				</div>

				<!-- Preview Area -->
				<div class="lg:col-span-2">
					<h4 class="text-lg font-medium text-gray-900 mb-4">👁️ Preview</h4>
					<div class="bg-gray-50 p-4 rounded-lg min-h-64">
						<table class="w-full text-sm">
							<thead>
								<tr class="bg-gray-200">
									<th class="p-2 text-left">Product</th>
									<th class="p-2 text-left">Q1</th>
									<th class="p-2 text-left">Q2</th>
									<th class="p-2 text-left">Q3</th>
									<th class="p-2 text-left">Q4</th>
									<th class="p-2 text-left">Total</th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td class="p-2 border">Widget A</td>
									<td class="p-2 border">12,500</td>
									<td class="p-2 border">15,200</td>
									<td class="p-2 border">18,900</td>
									<td class="p-2 border">14,300</td>
									<td class="p-2 border font-bold">60,900</td>
								</tr>
								<tr>
									<td class="p-2 border">Widget B</td>
									<td class="p-2 border">8,900</td>
									<td class="p-2 border">11,200</td>
									<td class="p-2 border">9,800</td>
									<td class="p-2 border">12,400</td>
									<td class="p-2 border font-bold">42,300</td>
								</tr>
								<tr class="bg-gray-200 font-bold">
									<td class="p-2 border">Total</td>
									<td class="p-2 border">21,400</td>
									<td class="p-2 border">26,400</td>
									<td class="p-2 border">28,700</td>
									<td class="p-2 border">26,700</td>
									<td class="p-2 border">103,200</td>
								</tr>
							</tbody>
						</table>
					</div>

					<div class="mt-4 flex justify-end space-x-3">
						<button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700">
							📊 Create Chart
						</button>
						<button class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700">
							✅ Insert Pivot Table
						</button>
					</div>
				</div>
			</div>

			<!-- Advanced Features -->
			<div class="mt-8">
				<h4 class="text-lg font-medium text-gray-900 mb-4">🚀 Advanced Pivot Features</h4>
				<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
					<div class="bg-blue-50 p-4 rounded">
						<h5 class="font-medium text-blue-900 mb-2">🎲 Calculated Fields</h5>
						<p class="text-sm text-blue-700">Add custom calculations like profit margins, growth rates, or percentage changes.</p>
					</div>
					<div class="bg-green-50 p-4 rounded">
						<h5 class="font-medium text-green-900 mb-2">🔍 Filters & Slicers</h5>
						<p class="text-sm text-green-700">Interactive filters to drill down into your data and focus on specific segments.</p>
					</div>
					<div class="bg-purple-50 p-4 rounded">
						<h5 class="font-medium text-purple-900 mb-2">📈 Multiple Aggregations</h5>
						<p class="text-sm text-purple-700">Sum, average, count, min, max, and custom aggregations for comprehensive analysis.</p>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Data Preparation Tools Modal -->
{#if $showDataPrepTools}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-5xl max-h-[90vh] overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">🔧 Data Preparation Tools</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showDataPrepTools.set(false)}
				>
					×
				</button>
			</div>

			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				<!-- Quick Actions -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-4">⚡ Quick Fixes</h4>
					<div class="space-y-3">
						<button class="w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 text-left transition-colors">
							<div class="flex items-center">
								<span class="text-2xl mr-3">🧹</span>
								<div>
									<h5 class="font-medium">Remove Empty Rows/Columns</h5>
									<p class="text-sm text-gray-600">Clean up your data by removing blank entries</p>
								</div>
							</div>
						</button>

						<button class="w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 text-left transition-colors">
							<div class="flex items-center">
								<span class="text-2xl mr-3">🔄</span>
								<div>
									<h5 class="font-medium">Convert to Table Format</h5>
									<p class="text-sm text-gray-600">Transform your data into proper table structure</p>
								</div>
							</div>
						</button>

						<button class="w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 text-left transition-colors">
							<div class="flex items-center">
								<span class="text-2xl mr-3">📋</span>
								<div>
									<h5 class="font-medium">Normalize Data Types</h5>
									<p class="text-sm text-gray-600">Ensure consistent data types across columns</p>
								</div>
							</div>
						</button>

						<button class="w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 text-left transition-colors">
							<div class="flex items-center">
								<span class="text-2xl mr-3">🔀</span>
								<div>
									<h5 class="font-medium">Unpivot Data</h5>
									<p class="text-sm text-gray-600">Convert wide tables to long format for analysis</p>
								</div>
							</div>
						</button>
					</div>
				</div>

				<!-- Advanced Tools -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-4">🛠️ Advanced Tools</h4>
					<div class="space-y-3">
						<div class="p-4 border border-gray-200 rounded-lg">
							<h5 class="font-medium mb-2">📊 Data Profiling</h5>
							<p class="text-sm text-gray-600 mb-3">Analyze your data structure, identify patterns, and get insights about data quality.</p>
							<div class="grid grid-cols-3 gap-2 text-xs">
								<div class="text-center p-2 bg-blue-50 rounded">
									<div class="font-medium">15</div>
									<div class="text-gray-600">Columns</div>
								</div>
								<div class="text-center p-2 bg-green-50 rounded">
									<div class="font-medium">1,247</div>
									<div class="text-gray-600">Rows</div>
								</div>
								<div class="text-center p-2 bg-yellow-50 rounded">
									<div class="font-medium">98.5%</div>
									<div class="text-gray-600">Complete</div>
								</div>
							</div>
						</div>

						<div class="p-4 border border-gray-200 rounded-lg">
							<h5 class="font-medium mb-2">🔄 Data Transformation</h5>
							<p class="text-sm text-gray-600 mb-3">Apply advanced transformations to prepare data for analysis.</p>
							<div class="space-y-2">
								<label class="flex items-center text-sm">
									<input type="checkbox" class="mr-2" />
									Convert text to proper case
								</label>
								<label class="flex items-center text-sm">
									<input type="checkbox" class="mr-2" />
									Remove duplicate rows
								</label>
								<label class="flex items-center text-sm">
									<input type="checkbox" class="mr-2" />
									Standardize date formats
								</label>
							</div>
						</div>

						<div class="p-4 border border-gray-200 rounded-lg">
							<h5 class="font-medium mb-2">🎯 Pivot-Ready Conversion</h5>
							<p class="text-sm text-gray-600 mb-3">One-click conversion to pivot table friendly format.</p>
							<div class="space-y-2">
								<button class="w-full px-3 py-2 bg-blue-600 text-white text-sm rounded hover:bg-blue-700">
									Convert to Long Format
								</button>
								<button class="w-full px-3 py-2 bg-green-600 text-white text-sm rounded hover:bg-green-700">
									Create Fact-Dimension Tables
								</button>
								<button class="w-full px-3 py-2 bg-purple-600 text-white text-sm rounded hover:bg-purple-700">
									Generate Pivot Template
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>

			<!-- Data Preview -->
			<div class="mt-8">
				<h4 class="text-lg font-medium text-gray-900 mb-4">👁️ Data Preview</h4>
				<div class="bg-gray-50 p-4 rounded-lg overflow-x-auto">
					<table class="w-full text-sm">
						<thead>
							<tr class="bg-gray-200">
								<th class="p-2 text-left">Product</th>
								<th class="p-2 text-left">Category</th>
								<th class="p-2 text-left">Jan</th>
								<th class="p-2 text-left">Feb</th>
								<th class="p-2 text-left">Mar</th>
								<th class="p-2 text-left">Total</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td class="p-2 border">Widget A</td>
								<td class="p-2 border">Electronics</td>
								<td class="p-2 border">1,200</td>
								<td class="p-2 border">1,500</td>
								<td class="p-2 border">1,800</td>
								<td class="p-2 border font-bold">4,500</td>
							</tr>
							<tr>
								<td class="p-2 border">Widget B</td>
								<td class="p-2 border">Electronics</td>
								<td class="p-2 border">900</td>
								<td class="p-2 border">1,100</td>
								<td class="p-2 border">1,200</td>
								<td class="p-2 border font-bold">3,200</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>

			<div class="mt-6 flex justify-end space-x-3">
				<button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" on:click={() => showDataPrepTools.set(false)}>
					Cancel
				</button>
				<button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
					Apply Changes
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Format Cells Dialog -->
{#if $showFormatDialog}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-lg">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">🎨 Format Cells</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showFormatDialog.set(false)}
				>
					×
				</button>
			</div>

			<div class="space-y-6">
				<!-- Number Formatting -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">📊 Number</h4>
					<div class="grid grid-cols-2 gap-3">
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-left"
							on:click={() => { dispatchAction('applyNumberFormat', { formatType: 'general' }); showFormatDialog.set(false); }}
						>
							<div class="font-medium">General</div>
							<div class="text-sm text-gray-500">No specific formatting</div>
						</button>
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-left"
							on:click={() => { dispatchAction('applyNumberFormat', { formatType: 'currency' }); showFormatDialog.set(false); }}
						>
							<div class="font-medium">Currency</div>
							<div class="text-sm text-gray-500">$1,234.56</div>
						</button>
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-left"
							on:click={() => { dispatchAction('applyNumberFormat', { formatType: 'percent' }); showFormatDialog.set(false); }}
						>
							<div class="font-medium">Percentage</div>
							<div class="text-sm text-gray-500">123.46%</div>
						</button>
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-left"
							on:click={() => { dispatchAction('applyNumberFormat', { formatType: 'decimal', decimals: 2 }); showFormatDialog.set(false); }}
						>
							<div class="font-medium">Number</div>
							<div class="text-sm text-gray-500">1,234.56</div>
						</button>
					</div>
				</div>

				<!-- Alignment -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">📐 Alignment</h4>
					<div class="grid grid-cols-3 gap-3">
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center"
							on:click={() => { dispatchAction('applyFormatting', { type: 'align', value: 'left' }); showFormatDialog.set(false); }}
						>
							<div class="text-2xl mb-2">⬅️</div>
							<div class="font-medium">Left</div>
						</button>
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center"
							on:click={() => { dispatchAction('applyFormatting', { type: 'align', value: 'center' }); showFormatDialog.set(false); }}
						>
							<div class="text-2xl mb-2">⬌</div>
							<div class="font-medium">Center</div>
						</button>
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center"
							on:click={() => { dispatchAction('applyFormatting', { type: 'align', value: 'right' }); showFormatDialog.set(false); }}
						>
							<div class="text-2xl mb-2">➡️</div>
							<div class="font-medium">Right</div>
						</button>
					</div>
				</div>

				<!-- Font Formatting -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">✏️ Font</h4>
					<div class="grid grid-cols-3 gap-3">
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center"
							on:click={() => { dispatchAction('applyFormatting', { type: 'bold' }); showFormatDialog.set(false); }}
						>
							<div class="text-2xl mb-2"><strong>B</strong></div>
							<div class="font-medium">Bold</div>
						</button>
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center"
							on:click={() => { dispatchAction('applyFormatting', { type: 'italic' }); showFormatDialog.set(false); }}
						>
							<div class="text-2xl mb-2"><em>I</em></div>
							<div class="font-medium">Italic</div>
						</button>
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center"
							on:click={() => { dispatchAction('applyFormatting', { type: 'underline' }); showFormatDialog.set(false); }}
						>
							<div class="text-2xl mb-2"><u>U</u></div>
							<div class="font-medium">Underline</div>
						</button>
					</div>
				</div>

				<!-- Border Options -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">🔲 Borders</h4>
					<div class="grid grid-cols-2 gap-3">
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center"
							on:click={() => { dispatchAction('applyBorder', { borderType: 'all' }); showFormatDialog.set(false); }}
						>
							<div class="text-2xl mb-2">🔲</div>
							<div class="font-medium">All Borders</div>
						</button>
						<button
							class="p-3 border border-gray-200 rounded hover:bg-blue-50 hover:border-blue-300 text-center"
							on:click={() => { dispatchAction('applyBorder', { borderType: 'outline' }); showFormatDialog.set(false); }}
						>
							<div class="text-2xl mb-2">▭</div>
							<div class="font-medium">Outline</div>
						</button>
					</div>
				</div>
			</div>

			<div class="mt-6 flex justify-end space-x-3">
				<button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" on:click={() => showFormatDialog.set(false)}>
					Cancel
				</button>
				<button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" on:click={() => showFormatDialog.set(false)}>
					Apply
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Context Menu -->
{#if $showContextMenu}
	<div class="fixed z-50" style="left: {$contextMenuPosition.x}px; top: {$contextMenuPosition.y}px;">
		<div class="bg-white border border-gray-300 rounded-lg shadow-lg py-1 min-w-48">
			{#each $contextMenuItems as item}
				{#if item.type === 'divider'}
					<div class="border-t border-gray-200 my-1"></div>
				{:else}
					<button
						class="w-full text-left px-4 py-2 text-sm hover:bg-gray-100 flex items-center justify-between {item.disabled ? 'opacity-50 cursor-not-allowed' : ''}"
						disabled={item.disabled}
						on:click={() => executeContextMenuAction(item)}
					>
						<span class="flex items-center">
							<span class="mr-3 text-lg">{item.icon}</span>
							<span>{item.label}</span>
						</span>
						{#if item.shortcut}
							<span class="text-xs text-gray-500 ml-4">{item.shortcut}</span>
						{/if}
					</button>
				{/if}
			{/each}
		</div>
	</div>
{/if}

<!-- Ribbon Customizer Modal (simplified placeholder) -->
{#if $showRibbonCustomizer}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
		<div class="bg-white rounded-lg w-full max-w-6xl max-h-[90vh] overflow-hidden">
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<h3 class="text-xl font-semibold text-gray-900">🎨 Customize Ribbon</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showRibbonCustomizer.set(false)}
				>
					×
				</button>
			</div>

			<div class="p-6">
				<div class="bg-blue-50 p-4 rounded-lg mb-4">
					<h5 class="font-medium text-blue-900 mb-2">📖 How to Customize</h5>
					<ul class="text-sm text-blue-800 space-y-1">
						<li>• <strong>Drag & Drop:</strong> Drag tools from the palette to add them to the ribbon</li>
						<li>• <strong>Remove Tools:</strong> Hover over ribbon tools and click the × to remove them</li>
						<li>• <strong>Switch Tabs:</strong> Click on different ribbon tabs to customize each one</li>
						<li>• <strong>Reset:</strong> Use "Reset to Default" to restore the original ribbon layout</li>
						<li>• <strong>Save:</strong> Your customizations are automatically saved to your browser</li>
					</ul>
				</div>

				<div class="flex justify-end space-x-3">
					<button class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700" on:click={() => showRibbonCustomizer.set(false)}>
						Close
					</button>
					<button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" on:click={() => showRibbonCustomizer.set(false)}>
						Save & Close
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
