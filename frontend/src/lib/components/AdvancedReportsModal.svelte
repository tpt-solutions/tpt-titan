<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { getForms } from '$lib/api.js';

	export let forms = [];
	export let showModal = false;

	const dispatch = createEventDispatcher();

	let selectedForm = '';
	let selectedReportType = 'summary';
	let dateRange = {
		start: '',
		end: ''
	};
	let filters = [];
	let reportData = null;
	let isGenerating = false;
	let exportFormat = 'csv';

	// Report types
	const reportTypes = [
		{ id: 'summary', name: 'Response Summary', description: 'Overview of form responses' },
		{ id: 'detailed', name: 'Detailed Responses', description: 'Individual response details' },
		{ id: 'trends', name: 'Response Trends', description: 'Response patterns over time' },
		{ id: 'completion', name: 'Completion Rates', description: 'Form completion analytics' },
		{ id: 'field-analysis', name: 'Field Analysis', description: 'Analysis of specific fields' },
		{ id: 'cross-tab', name: 'Cross-tabulation', description: 'Cross-reference multiple fields' }
	];

	// Filter types
	const filterTypes = [
		{ id: 'date_range', name: 'Date Range', description: 'Filter by submission date' },
		{ id: 'field_value', name: 'Field Value', description: 'Filter by field responses' },
		{ id: 'completion_status', name: 'Completion Status', description: 'Complete vs incomplete responses' }
	];

	onMount(async () => {
		if (showModal) {
			// Initialize with first form if available
			if (forms.length > 0) {
				selectedForm = forms[0].id;
			}
		}
	});

	$: if (showModal && !selectedForm && forms.length > 0) {
		selectedForm = forms[0].id;
	}

	async function generateReport() {
		if (!selectedForm) return;

		isGenerating = true;
		reportData = null;

		try {
			const response = await fetch(`/api/v1/forms/${selectedForm}/reports`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					report_type: selectedReportType,
					date_range: dateRange,
					filters: filters,
					include_data: true
				})
			});

			if (response.ok) {
				reportData = await response.json();
			} else {
				const error = await response.json();
				alert('Failed to generate report: ' + (error.error || 'Unknown error'));
			}
		} catch (error) {
			console.error('Failed to generate report:', error);
			alert('Failed to generate report: ' + error.message);
		} finally {
			isGenerating = false;
		}
	}

	async function exportReport() {
		if (!reportData) return;

		try {
			const response = await fetch(`/api/v1/forms/${selectedForm}/reports/export`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				},
				body: JSON.stringify({
					report_type: selectedReportType,
					date_range: dateRange,
					filters: filters,
					format: exportFormat
				})
			});

			if (response.ok) {
				const blob = await response.blob();
				const url = URL.createObjectURL(blob);
				const a = document.createElement('a');
				a.href = url;
				a.download = `report_${selectedForm}_${new Date().toISOString().split('T')[0]}.${exportFormat}`;
				a.click();
				URL.revokeObjectURL(url);
			} else {
				const error = await response.json();
				alert('Failed to export report: ' + (error.error || 'Unknown error'));
			}
		} catch (error) {
			console.error('Failed to export report:', error);
			alert('Failed to export report: ' + error.message);
		}
	}

	function addFilter() {
		filters = [...filters, {
			id: Date.now(),
			type: 'field_value',
			field: '',
			operator: 'equals',
			value: ''
		}];
	}

	function removeFilter(filterId) {
		filters = filters.filter(f => f.id !== filterId);
	}

	function updateFilter(filterId, property, value) {
		filters = filters.map(f =>
			f.id === filterId ? { ...f, [property]: value } : f
		);
	}

	function getFormFields(formId) {
		const form = forms.find(f => f.id === formId);
		if (form && form.fields) {
			return form.fields.map(f => ({ id: f.id, label: f.label || f.name || f.id }));
		}
		return [];
	}

	function closeModal() {
		showModal = false;
		reportData = null;
		dispatch('close');
	}
</script>

{#if showModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-6xl max-h-[90vh] overflow-hidden flex flex-col">
			<!-- Header -->
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<div>
					<h2 class="text-2xl font-bold text-gray-900">📊 Advanced Reports</h2>
					<p class="text-gray-600 mt-1">Generate comprehensive reports and analytics for your forms</p>
				</div>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={closeModal}
				>
					×
				</button>
			</div>

			<!-- Main Content -->
			<div class="flex-1 flex overflow-hidden">
				<!-- Sidebar - Report Configuration -->
				<div class="w-80 bg-gray-50 border-r border-gray-200 p-4 overflow-y-auto">
					<h3 class="text-lg font-semibold text-gray-900 mb-4">Report Configuration</h3>

					<div class="space-y-4">
						<!-- Form Selection -->
						<div>
							<label for="report-form" class="block text-sm font-medium text-gray-700 mb-1">Select Form</label>
							<select
								id="report-form"
								bind:value={selectedForm}
								class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
							>
								{#each forms as form}
									<option value={form.id}>{form.name}</option>
								{/each}
							</select>
						</div>

						<!-- Report Type -->
						<div>
							<label for="report-type" class="block text-sm font-medium text-gray-700 mb-1">Report Type</label>
							<select
								id="report-type"
								bind:value={selectedReportType}
								class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
							>
								{#each reportTypes as type}
									<option value={type.id}>{type.name}</option>
								{/each}
							</select>
							<p class="text-xs text-gray-500 mt-1">
								{reportTypes.find(t => t.id === selectedReportType)?.description}
							</p>
						</div>

						<!-- Date Range -->
						<div class="grid grid-cols-2 gap-2">
							<div>
								<label for="date-start" class="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
								<input
									id="date-start"
									type="date"
									bind:value={dateRange.start}
									class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
							<div>
								<label for="date-end" class="block text-sm font-medium text-gray-700 mb-1">End Date</label>
								<input
									id="date-end"
									type="date"
									bind:value={dateRange.end}
									class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
						</div>

						<!-- Filters -->
						<div>
							<div class="flex items-center justify-between mb-2">
								<label class="block text-sm font-medium text-gray-700">Filters</label>
								<button
									class="px-2 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700"
									on:click={addFilter}
								>
									+ Add Filter
								</button>
							</div>

							<div class="space-y-2 max-h-40 overflow-y-auto">
								{#each filters as filter}
									<div class="border border-gray-200 rounded p-3 bg-white">
										<div class="flex items-center justify-between mb-2">
											<select
												bind:value={filter.type}
												on:change={(e) => updateFilter(filter.id, 'type', e.target.value)}
												class="text-sm border border-gray-300 rounded px-2 py-1"
											>
												{#each filterTypes as type}
													<option value={type.id}>{type.name}</option>
												{/each}
											</select>
											<button
												class="text-red-500 hover:text-red-700 text-sm"
												on:click={() => removeFilter(filter.id)}
											>
												✕
											</button>
										</div>

										{#if filter.type === 'field_value'}
											<div class="grid grid-cols-3 gap-2">
												<select
													bind:value={filter.field}
													on:change={(e) => updateFilter(filter.id, 'field', e.target.value)}
													class="text-sm border border-gray-300 rounded px-2 py-1"
												>
													<option value="">Field</option>
													{#each getFormFields(selectedForm) as field}
														<option value={field.id}>{field.label}</option>
													{/each}
												</select>
												<select
													bind:value={filter.operator}
													on:change={(e) => updateFilter(filter.id, 'operator', e.target.value)}
													class="text-sm border border-gray-300 rounded px-2 py-1"
												>
													<option value="equals">Equals</option>
													<option value="contains">Contains</option>
													<option value="greater_than">Greater Than</option>
													<option value="less_than">Less Than</option>
												</select>
												<input
													type="text"
													bind:value={filter.value}
													on:input={(e) => updateFilter(filter.id, 'value', e.target.value)}
													placeholder="Value"
													class="text-sm border border-gray-300 rounded px-2 py-1"
												/>
											</div>
										{/if}
									</div>
								{/each}
							</div>
						</div>

						<!-- Generate Button -->
						<button
							class="w-full px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
							on:click={generateReport}
							disabled={isGenerating || !selectedForm}
						>
							{isGenerating ? 'Generating...' : 'Generate Report'}
						</button>
					</div>
				</div>

				<!-- Main Content Area -->
				<div class="flex-1 p-6 overflow-y-auto">
					{#if reportData}
						<!-- Report Results -->
						<div class="mb-4 flex items-center justify-between">
							<h3 class="text-xl font-semibold text-gray-900">Report Results</h3>
							<div class="flex items-center space-x-2">
								<select
									bind:value={exportFormat}
									class="px-3 py-1 border border-gray-300 rounded text-sm"
								>
									<option value="csv">CSV</option>
									<option value="excel">Excel</option>
									<option value="pdf">PDF</option>
								</select>
								<button
									class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 text-sm"
									on:click={exportReport}
								>
									📊 Export Report
								</button>
							</div>
						</div>

						<!-- Report Content -->
						<div class="bg-white border border-gray-200 rounded-lg p-6">
							{#if selectedReportType === 'summary'}
								<!-- Summary Report -->
								<div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
									<div class="bg-blue-50 p-4 rounded-lg">
										<div class="text-2xl font-bold text-blue-600">{reportData.total_responses || 0}</div>
										<div class="text-sm text-blue-600">Total Responses</div>
									</div>
									<div class="bg-green-50 p-4 rounded-lg">
										<div class="text-2xl font-bold text-green-600">{reportData.completed_responses || 0}</div>
										<div class="text-sm text-green-600">Completed</div>
									</div>
									<div class="bg-yellow-50 p-4 rounded-lg">
										<div class="text-2xl font-bold text-yellow-600">{Math.round((reportData.completion_rate || 0) * 100)}%</div>
										<div class="text-sm text-yellow-600">Completion Rate</div>
									</div>
									<div class="bg-purple-50 p-4 rounded-lg">
										<div class="text-2xl font-bold text-purple-600">{reportData.average_time || 'N/A'}</div>
										<div class="text-sm text-purple-600">Avg. Time</div>
									</div>
								</div>

								{#if reportData.field_stats}
									<div>
										<h4 class="text-lg font-semibold text-gray-900 mb-4">Field Statistics</h4>
										<div class="overflow-x-auto">
											<table class="w-full border border-gray-200 rounded">
												<thead class="bg-gray-50">
													<tr>
														<th class="px-4 py-2 text-left border-b">Field</th>
														<th class="px-4 py-2 text-left border-b">Responses</th>
														<th class="px-4 py-2 text-left border-b">Completion %</th>
														<th class="px-4 py-2 text-left border-b">Avg Length</th>
													</tr>
												</thead>
												<tbody>
													{#each reportData.field_stats as stat}
														<tr class="border-b border-gray-200">
															<td class="px-4 py-2">{stat.field_name}</td>
															<td class="px-4 py-2">{stat.response_count}</td>
															<td class="px-4 py-2">{Math.round(stat.completion_rate * 100)}%</td>
															<td class="px-4 py-2">{stat.avg_length || 'N/A'}</td>
														</tr>
													{/each}
												</tbody>
											</table>
										</div>
									</div>
								{/if}

							{:else if selectedReportType === 'detailed'}
								<!-- Detailed Responses -->
								{#if reportData.responses && reportData.responses.length > 0}
									<div class="overflow-x-auto">
										<table class="w-full border border-gray-200 rounded">
											<thead class="bg-gray-50">
												<tr>
													<th class="px-4 py-2 text-left border-b">ID</th>
													<th class="px-4 py-2 text-left border-b">Submitted</th>
													{#each reportData.fields || [] as field}
														<th class="px-4 py-2 text-left border-b">{field.label}</th>
													{/each}
												</tr>
											</thead>
											<tbody>
												{#each reportData.responses as response}
													<tr class="border-b border-gray-200">
														<td class="px-4 py-2">{response.id}</td>
														<td class="px-4 py-2">{new Date(response.submitted_at).toLocaleString()}</td>
														{#each reportData.fields || [] as field}
															<td class="px-4 py-2">{response.responses?.[field.label] || '-'}</td>
														{/each}
													</tr>
												{/each}
											</tbody>
										</table>
									</div>
								{:else}
									<div class="text-center py-8 text-gray-500">
										<p>No responses found for the selected criteria.</p>
									</div>
								{/if}

							{:else if selectedReportType === 'trends'}
								<!-- Trends Report -->
								<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
									<div>
										<h4 class="text-lg font-semibold text-gray-900 mb-4">Response Trends</h4>
										{#if reportData.trends && reportData.trends.length > 0}
											<div class="bg-gray-50 p-4 rounded space-y-2">
												{#each reportData.trends as point}
													<div>
														<div class="flex justify-between text-sm text-gray-700">
															<span>{point.label}</span>
															<span>{point.count}</span>
														</div>
														<div class="w-full bg-gray-200 rounded h-2 mt-1">
															<div class="bg-blue-500 h-2 rounded" style="width: {Math.max(2, Math.round((point.count / (reportData.trend_max || point.count || 1)) * 100))}%"></div>
														</div>
													</div>
												{/each}
											</div>
										{:else}
											<div class="bg-gray-50 p-4 rounded">
												<p class="text-gray-600">No trend data returned for the selected criteria.</p>
												<p class="text-sm text-gray-500 mt-2">Export the raw data to analyze response patterns over time.</p>
											</div>
										{/if}
									</div>
									<div>
										<h4 class="text-lg font-semibold text-gray-900 mb-4">Peak Hours</h4>
										{#if reportData.peak_hours && reportData.peak_hours.length > 0}
											<div class="bg-gray-50 p-4 rounded space-y-2">
												{#each reportData.peak_hours as point}
													<div>
														<div class="flex justify-between text-sm text-gray-700">
															<span>{point.label}</span>
															<span>{point.count}</span>
														</div>
														<div class="w-full bg-gray-200 rounded h-2 mt-1">
															<div class="bg-purple-500 h-2 rounded" style="width: {Math.max(2, Math.round((point.count / (reportData.peak_max || point.count || 1)) * 100))}%"></div>
														</div>
													</div>
												{/each}
											</div>
										{:else}
											<div class="bg-gray-50 p-4 rounded">
												<p class="text-gray-600">No peak-hour data returned for the selected criteria.</p>
												<p class="text-sm text-gray-500 mt-2">Export the raw data to identify when users are most active.</p>
											</div>
										{/if}
									</div>
								</div>

							{:else}
								<!-- Generic data-driven view for completion / field-analysis / cross-tab -->
								{#if reportData.responses && reportData.responses.length > 0}
									<h4 class="text-lg font-semibold text-gray-900 mb-4">
										{reportTypes.find(t => t.id === selectedReportType)?.name} — {reportData.responses.length} responses
									</h4>
									<div class="overflow-x-auto">
										<table class="w-full border border-gray-200 rounded">
											<thead class="bg-gray-50">
												<tr>
													<th class="px-4 py-2 text-left border-b">ID</th>
													<th class="px-4 py-2 text-left border-b">Submitted</th>
													{#each reportData.fields || Object.keys(reportData.responses[0].responses || {}) as field}
														<th class="px-4 py-2 text-left border-b">{typeof field === 'string' ? field : field.label}</th>
													{/each}
												</tr>
											</thead>
											<tbody>
												{#each reportData.responses as response}
													<tr class="border-b border-gray-200">
														<td class="px-4 py-2">{response.id}</td>
														<td class="px-4 py-2">{response.submitted_at ? new Date(response.submitted_at).toLocaleString() : '-'}</td>
														{#each reportData.fields || Object.keys(response.responses || {}) as field}
															{@const key = typeof field === 'string' ? field : field.label}
															<td class="px-4 py-2">{response.responses?.[key] || '-'}</td>
														{/each}
													</tr>
												{/each}
											</tbody>
										</table>
									</div>
								{:else if reportData.field_stats}
									<h4 class="text-lg font-semibold text-gray-900 mb-4">
										{reportTypes.find(t => t.id === selectedReportType)?.name} — Field Statistics
									</h4>
									<div class="overflow-x-auto">
										<table class="w-full border border-gray-200 rounded">
											<thead class="bg-gray-50">
												<tr>
													<th class="px-4 py-2 text-left border-b">Field</th>
													<th class="px-4 py-2 text-left border-b">Responses</th>
													<th class="px-4 py-2 text-left border-b">Completion %</th>
												</tr>
											</thead>
											<tbody>
												{#each reportData.field_stats as stat}
													<tr class="border-b border-gray-200">
														<td class="px-4 py-2">{stat.field_name}</td>
														<td class="px-4 py-2">{stat.response_count}</td>
														<td class="px-4 py-2">{Math.round((stat.completion_rate || 0) * 100)}%</td>
													</tr>
												{/each}
											</tbody>
										</table>
									</div>
								{:else}
									<div class="text-center py-8 text-gray-500">
										<p>No data available for this report configuration.</p>
										<p class="text-sm mt-2">Try adjusting the filters or date range, or export the raw data.</p>
									</div>
								{/if}
							{/if}
						</div>
					{:else}
						<!-- No Report Generated Yet -->
						<div class="text-center py-16">
							<div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
								<svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
								</svg>
							</div>
							<h3 class="text-xl font-medium text-gray-900 mb-2">Generate Your First Report</h3>
							<p class="text-gray-500 mb-6">Select a form and report type, then click "Generate Report" to get started.</p>
							<div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-left">
								<div class="bg-blue-50 p-4 rounded-lg">
									<h4 class="font-medium text-blue-900 mb-2">📊 Summary Reports</h4>
									<p class="text-sm text-blue-700">Overview statistics and completion rates</p>
								</div>
								<div class="bg-green-50 p-4 rounded-lg">
									<h4 class="font-medium text-green-900 mb-2">📋 Detailed Responses</h4>
									<p class="text-sm text-green-700">Individual response data and filtering</p>
								</div>
								<div class="bg-purple-50 p-4 rounded-lg">
									<h4 class="font-medium text-purple-900 mb-2">📈 Analytics</h4>
									<p class="text-sm text-purple-700">Trends, patterns, and insights</p>
								</div>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}
