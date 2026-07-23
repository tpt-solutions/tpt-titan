<script>
// @ts-nocheck
	import { templates, templateCategories } from '$lib/stores.js';
	import { createEventDispatcher } from 'svelte';

	export let showTemplates = false;

	const dispatch = createEventDispatcher();

	let selectedCategory = 'all';
	let searchQuery = '';

	$: filteredTemplates = $templates.filter(template => {
		const matchesCategory = selectedCategory === 'all' || template.category.toLowerCase() === selectedCategory.toLowerCase();
		const matchesSearch = searchQuery === '' ||
			template.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			template.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
			template.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()));

		return matchesCategory && matchesSearch;
	});

	function selectTemplate(template) {
		dispatch('select', { template });
		showTemplates = false;
	}

	function createBlankSpreadsheet() {
		dispatch('blank');
		showTemplates = false;
	}
</script>

{#if showTemplates}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg w-full max-w-6xl max-h-[90vh] overflow-hidden flex flex-col">
			<!-- Header -->
			<div class="flex items-center justify-between p-6 border-b border-gray-200">
				<div>
					<h2 class="text-2xl font-bold text-gray-900">Choose a Template</h2>
					<p class="text-gray-600 mt-1">Start with a pre-built template or create a blank spreadsheet</p>
				</div>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => showTemplates = false}
				>
					×
				</button>
			</div>

			<!-- Search and Filters -->
			<div class="p-6 border-b border-gray-200 bg-gray-50">
				<div class="flex items-center justify-between mb-4">
					<!-- Category Filter -->
					<div class="flex space-x-1">
						{#each templateCategories as category}
							<button
								class="px-3 py-1 text-sm rounded-full transition-colors flex items-center space-x-1
									{selectedCategory === category.id
										? 'bg-blue-600 text-white'
										: 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'}"
								on:click={() => selectedCategory = category.id}
							>
								<span>{category.icon}</span>
								<span>{category.name}</span>
							</button>
						{/each}
					</div>

					<!-- Create Blank Button -->
					<button
						class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors flex items-center space-x-2"
						on:click={createBlankSpreadsheet}
					>
						<span>📄</span>
						<span>Blank Spreadsheet</span>
					</button>
				</div>

				<!-- Search -->
				<div class="relative">
					<input
						type="text"
						placeholder="Search templates..."
						bind:value={searchQuery}
						class="w-full px-4 py-2 pl-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
					/>
					<svg class="w-5 h-5 text-gray-400 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
					</svg>
				</div>
			</div>

			<!-- Template Grid -->
			<div class="flex-1 overflow-y-auto p-6">
				{#if filteredTemplates.length === 0}
					<div class="text-center py-12">
						<div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
							<svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6-4h6m2 5.291A7.962 7.962 0 0112 15c-2.34 0-4.29-.966-5.5-2.5"></path>
							</svg>
						</div>
						<h3 class="text-lg font-medium text-gray-900 mb-2">No templates found</h3>
						<p class="text-gray-500">Try adjusting your search or category filter</p>
					</div>
				{:else}
					<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
						{#each filteredTemplates as template}
							<button
								class="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-lg transition-all cursor-pointer group text-left w-full"
								on:click={() => selectTemplate(template)}
								type="button"
							>
								<!-- Template Icon and Category -->
								<div class="flex items-center justify-between mb-4">
									<div class="text-3xl">{template.icon}</div>
									<span class="px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded-full">
										{template.category}
									</span>
								</div>

								<!-- Template Info -->
								<div class="mb-4">
									<h3 class="text-lg font-semibold text-gray-900 mb-2 group-hover:text-blue-600 transition-colors">
										{template.name}
									</h3>
									<p class="text-sm text-gray-600 mb-3 line-clamp-2">
										{template.description}
									</p>

									<!-- Tags -->
									<div class="flex flex-wrap gap-1 mb-3">
										{#each template.tags.slice(0, 3) as tag}
											<span class="px-2 py-1 text-xs bg-blue-50 text-blue-600 rounded">
												{tag}
											</span>
										{/each}
										{#if template.tags.length > 3}
											<span class="px-2 py-1 text-xs bg-gray-50 text-gray-600 rounded">
												+{template.tags.length - 3}
											</span>
										{/if}
									</div>
								</div>

								<!-- Preview Button -->
								<button
									class="w-full px-3 py-2 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
									on:click|stopPropagation={() => dispatch('preview', { template })}
									type="button"
								>
									👁️ Preview
								</button>
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Footer -->
			<div class="border-t border-gray-200 p-6 bg-gray-50">
				<div class="flex items-center justify-between text-sm text-gray-600">
					<span>{filteredTemplates.length} templates available</span>
					<span>💡 Tip: Templates include formulas and formatting</span>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
