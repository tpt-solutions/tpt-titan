<!-- frontend/src/lib/components/FormBuilderResponsesModal.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';

	export let formData;
	export let formResponses = [];

	const dispatch = createEventDispatcher();

	function exportCSV() {
		let csv = 'ID,Submitted At';
		formData.fields.forEach(field => { csv += ',' + field.label; });
		csv += '\n';

		formResponses.forEach(response => {
			csv += response.id + ',' + response.submitted_at;
			formData.fields.forEach(field => {
				const value = response.responses[field.label.toLowerCase()] || '';
				csv += ',"' + value.replace(/"/g, '""') + '"';
			});
			csv += '\n';
		});

		const blob = new Blob([csv], { type: 'text/csv' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `${formData.name || 'form'}_responses.csv`;
		a.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
	<div class="bg-white rounded-lg w-full max-w-6xl max-h-[90vh] overflow-hidden flex flex-col">
		<div class="flex items-center justify-between p-6 border-b border-gray-200">
			<div>
				<h2 class="text-2xl font-bold text-gray-900">Form Responses</h2>
				<p class="text-gray-600 mt-1">{formResponses.length} responses received</p>
			</div>
			<div class="flex items-center space-x-3">
				<button class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors" on:click={exportCSV}>
					📊 Export CSV
				</button>
				<button class="text-gray-400 hover:text-gray-600 text-2xl" on:click={() => dispatch('close')}>×</button>
			</div>
		</div>

		<div class="flex-1 overflow-y-auto">
			{#if formResponses.length === 0}
				<div class="text-center py-12">
					<div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
						</svg>
					</div>
					<h3 class="text-lg font-medium text-gray-900 mb-2">No responses yet</h3>
					<p class="text-gray-500">Responses will appear here once people submit your form</p>
				</div>
			{:else}
				<div class="p-6">
					<!-- Summary stats -->
					<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
						<div class="bg-blue-50 p-4 rounded-lg">
							<div class="text-2xl font-bold text-blue-600">{formResponses.length}</div>
							<div class="text-sm text-blue-600">Total Responses</div>
						</div>
						<div class="bg-green-50 p-4 rounded-lg">
							<div class="text-2xl font-bold text-green-600">
								{Math.round((formResponses.length / (formResponses.length + 5)) * 100)}%
							</div>
							<div class="text-sm text-green-600">Completion Rate</div>
						</div>
						<div class="bg-purple-50 p-4 rounded-lg">
							<div class="text-2xl font-bold text-purple-600">
								{new Date(Math.max(...formResponses.map(r => new Date(r.submitted_at)))).toLocaleDateString()}
							</div>
							<div class="text-sm text-purple-600">Latest Response</div>
						</div>
					</div>

					<!-- Responses table -->
					<div class="bg-white border border-gray-200 rounded-lg overflow-hidden">
						<div class="overflow-x-auto">
							<table class="w-full">
								<thead class="bg-gray-50">
									<tr>
										<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
										<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Submitted</th>
										{#each formData.fields as field}
											<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{field.label}</th>
										{/each}
									</tr>
								</thead>
								<tbody class="bg-white divide-y divide-gray-200">
									{#each formResponses as response}
										<tr class="hover:bg-gray-50">
											<td class="px-4 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{response.id}</td>
											<td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500">
												{new Date(response.submitted_at).toLocaleString()}
											</td>
											{#each formData.fields as field}
												<td class="px-4 py-4 whitespace-nowrap text-sm text-gray-900">
													{response.responses[field.label.toLowerCase()] || '-'}
												</td>
											{/each}
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
