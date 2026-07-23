<!-- frontend/src/lib/components/FormVersionHistoryModal.svelte -->
<script>
// @ts-nocheck
	import { createEventDispatcher } from 'svelte';
	
	export let show = false;
	export let versions = [];
	export let currentVersion = null;
	
	const dispatch = createEventDispatcher();
	
	function formatDate(dateString) {
		return new Date(dateString).toLocaleString();
	}
	
	function handleRestore(version) {
		if (confirm(`Restore to version ${version.version}? Current changes will be lost.`)) {
			dispatch('restore', version);
		}
	}
	
	function handleClose() {
		dispatch('close');
	}
	
	function handlePreview(version) {
		dispatch('preview', version);
	}
	
	function getFieldCount(version) {
		if (!version.data || !version.data.fields) return 0;
		return version.data.fields.length;
	}
</script>

{#if show}
	<div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
		<div class="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[80vh] flex flex-col">
			<!-- Header -->
			<div class="flex items-center justify-between px-6 py-4 border-b border-gray-200">
				<h2 class="text-xl font-semibold text-gray-900">Form Version History</h2>
				<button
					class="text-gray-400 hover:text-gray-600 transition-colors"
					on:click={handleClose}
				>
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
			
			<!-- Content -->
			<div class="flex-1 overflow-y-auto p-6">
				{#if versions.length === 0}
					<div class="text-center py-8 text-gray-500">
						<svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
						</svg>
						<p>No version history available</p>
						<p class="text-sm mt-1">Save your form to create versions</p>
					</div>
				{:else}
					<div class="space-y-3">
						{#each versions as version, index}
							<div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors {version.is_active ? 'bg-blue-50 border-blue-200' : ''}">
								<div class="flex items-start space-x-3">
									<div class="flex-shrink-0">
										<div class="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center text-gray-600 font-semibold">
											{versions.length - index}
										</div>
									</div>
									<div>
										<div class="flex items-center space-x-2">
											<span class="font-medium text-gray-900">Version {version.version}</span>
											{#if version.is_active}
												<span class="px-2 py-0.5 text-xs bg-blue-100 text-blue-700 rounded-full">Current</span>
											{/if}
										</div>
										<div class="text-sm text-gray-500 mt-1">
											{formatDate(version.created_at)}
										</div>
										<div class="text-xs text-gray-400 mt-1">
											By {version.created_by || 'Unknown'} • {getFieldCount(version)} field{getFieldCount(version) !== 1 ? 's' : ''}
										</div>
										{#if version.comment}
											<div class="text-sm text-gray-600 mt-2 italic">
												"{version.comment}"
											</div>
										{/if}
									</div>
								</div>
								
								<div class="flex items-center space-x-2">
									<button
										class="px-3 py-1.5 text-sm text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded transition-colors"
										on:click={() => handlePreview(version)}
									>
										Preview
									</button>
									{#if !version.is_active}
										<button
											class="px-3 py-1.5 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
											on:click={() => handleRestore(version)}
										>
											Restore
										</button>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
			
			<!-- Footer -->
			<div class="px-6 py-4 border-t border-gray-200 bg-gray-50 rounded-b-lg">
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-500">
						{versions.length} version{versions.length !== 1 ? 's' : ''} saved
					</span>
					<button
						class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 transition-colors"
						on:click={handleClose}
					>
						Close
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
