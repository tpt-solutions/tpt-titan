<script>
	import { createEventDispatcher } from 'svelte';

	export let forms = [];

	const dispatch = createEventDispatcher();

	function handleEdit(form) {
		dispatch('edit', form);
	}

	function handleDelete(form) {
		if (confirm(`Are you sure you want to delete "${form.name}"? This action cannot be undone.`)) {
			dispatch('delete', form.id);
		}
	}

	function handleViewResponses(form) {
		dispatch('viewResponses', form);
	}

	function getStatusColor(status) {
		switch (status) {
			case 'active': return 'bg-green-100 text-green-800';
			case 'draft': return 'bg-yellow-100 text-yellow-800';
			case 'archived': return 'bg-gray-100 text-gray-800';
			default: return 'bg-gray-100 text-gray-800';
		}
	}
</script>

<div class="p-6">
	<div class="mb-6">
		<h2 class="text-2xl font-bold text-gray-900 mb-2">Your Forms</h2>
		<p class="text-gray-600">Create and manage forms with MS Access-style database features</p>
	</div>

	{#if forms.length === 0}
		<!-- Empty state -->
		<div class="text-center py-12">
			<div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
				<svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
				</svg>
			</div>
			<h3 class="text-lg font-medium text-gray-900 mb-2">No forms yet</h3>
			<p class="text-gray-500 mb-4">Create your first form to start collecting data</p>
			<button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
				Create Your First Form
			</button>
		</div>
	{:else}
		<!-- Forms grid -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each forms as form}
				<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow">
					<div class="flex items-start justify-between mb-4">
						<div class="flex-1">
							<h3 class="text-lg font-semibold text-gray-900 mb-1">{form.name}</h3>
							<p class="text-sm text-gray-600 mb-2">{form.description}</p>
							<div class="flex items-center space-x-4 text-sm text-gray-500">
								<span>{form.responses} responses</span>
								<span>Created {form.createdAt.toLocaleDateString()}</span>
							</div>
						</div>
						<span class="px-2 py-1 text-xs font-medium rounded-full {getStatusColor(form.status)}">
							{form.status}
						</span>
					</div>

					<div class="flex items-center justify-between">
						<div class="flex space-x-2">
							<button
								class="px-3 py-1 text-sm bg-blue-50 text-blue-700 rounded hover:bg-blue-100 transition-colors"
								on:click={() => handleViewResponses(form)}
							>
								View Responses
							</button>
							<button
								class="px-3 py-1 text-sm bg-gray-50 text-gray-700 rounded hover:bg-gray-100 transition-colors"
								on:click={() => handleEdit(form)}
							>
								Edit
							</button>
						</div>
						<button
							class="p-1 text-gray-400 hover:text-red-600 transition-colors"
							on:click={() => handleDelete(form)}
							title="Delete form"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
							</svg>
						</button>
					</div>
				</div>
			{/each}
		</div>

		<!-- Stats summary -->
		<div class="mt-8 bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			<h3 class="text-lg font-semibold text-gray-900 mb-4">Summary</h3>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
				<div class="text-center">
					<div class="text-2xl font-bold text-blue-600">{forms.length}</div>
					<div class="text-sm text-gray-600">Total Forms</div>
				</div>
				<div class="text-center">
					<div class="text-2xl font-bold text-green-600">{forms.reduce((sum, f) => sum + f.responses, 0)}</div>
					<div class="text-sm text-gray-600">Total Responses</div>
				</div>
				<div class="text-center">
					<div class="text-2xl font-bold text-purple-600">{forms.filter(f => f.status === 'active').length}</div>
					<div class="text-sm text-gray-600">Active Forms</div>
				</div>
			</div>
		</div>
	{/if}
</div>
