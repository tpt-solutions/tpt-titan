<!-- frontend/src/lib/components/TextEditorModals.svelte -->
<script>
// @ts-nocheck
	import { createEventDispatcher } from 'svelte';

	// Document list modal
	export let showDocumentList = false;
	export let documentList = [];
	export let isLoading = false;

	// Version history modal
	export let showVersionHistory = false;
	export let versionHistory = [];

	// Math help modal
	export let showMathHelp = false;

	// Find & replace modal
	export let showFindReplaceDialog = false;

	// AI assistant modal
	export let showAIAssistant = false;
	export let aiSuggestions = [];

	// Document summary modal
	export let documentSummary = '';

	// Text analysis modal
	export let showTextAnalysis = false;
	export let textAnalysis = null;

	// Find/replace state
	export let findText = '';
	export let replaceText = '';
	export let findResults = [];
	export let currentFindIndex = -1;
	export let isCaseSensitive = false;

	const dispatch = createEventDispatcher();

</script>

<!-- ── Document List Modal ─────────────────────────────────────── -->
{#if showDocumentList}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-md max-h-96 overflow-y-auto">
			<h3 class="text-lg font-semibold mb-4">Open Document</h3>
			{#if documentList.length === 0}
				<p class="text-gray-500">No documents found.</p>
			{:else}
				<div class="space-y-2">
					{#each documentList as doc}
						<button
							class="w-full text-left p-3 border border-gray-200 rounded hover:bg-gray-50"
							on:click={() => dispatch('loadDocument', doc)}
							disabled={isLoading}
						>
							<div class="font-medium">{doc.title}</div>
							<div class="text-sm text-gray-500">
								v{doc.version} • {new Date(doc.updated_at).toLocaleDateString()}
							</div>
						</button>
					{/each}
				</div>
			{/if}
			<div class="mt-4 flex justify-end">
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => dispatch('closeDocumentList')}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- ── Version History Modal ───────────────────────────────────── -->
{#if showVersionHistory}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-lg max-h-96 overflow-y-auto">
			<h3 class="text-lg font-semibold mb-4">Version History</h3>
			<div class="space-y-2">
				{#each versionHistory as version}
					<div class="flex items-center justify-between p-3 border border-gray-200 rounded">
						<div>
							<div class="font-medium">Version {version.version}</div>
							<div class="text-sm text-gray-500">
								{new Date(version.created_at).toLocaleString()}
								{#if version.is_active}
									<span class="ml-2 text-green-600">(Current)</span>
								{/if}
							</div>
						</div>
						{#if !version.is_active}
							<button
								class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
								on:click={() => dispatch('restoreVersion', version)}
							>
								Restore
							</button>
						{/if}
					</div>
				{/each}
			</div>
			<div class="mt-4 flex justify-end">
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => dispatch('closeVersionHistory')}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- ── Math Help Modal ─────────────────────────────────────────── -->
{#if showMathHelp}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">Natural Math Input - Better Than LaTeX!</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => dispatch('closeMathHelp')}
				>
					×
				</button>
			</div>

			<div class="space-y-6">
				<div class="bg-blue-50 p-4 rounded-lg">
					<h4 class="text-lg font-medium text-blue-900 mb-2">✨ Revolutionary Math Input</h4>
					<p class="text-blue-800">
						Forget complex LaTeX syntax! TPT Titan understands natural language math expressions.
						Simply type what you would say out loud, and watch it transform into beautiful mathematical notation.
					</p>
				</div>

				<!-- Basic Examples -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">📚 Basic Expressions</h4>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div class="bg-gray-50 p-3 rounded">
							<div class="text-sm text-gray-600 mb-1">Type this:</div>
							<code class="bg-white px-2 py-1 rounded text-sm">fraction a over b</code>
							<div class="text-sm text-gray-600 mt-1">Becomes:</div>
							<div class="font-serif text-lg">a/b</div>
						</div>
						<div class="bg-gray-50 p-3 rounded">
							<div class="text-sm text-gray-600 mb-1">Type this:</div>
							<code class="bg-white px-2 py-1 rounded text-sm">square root of x</code>
							<div class="text-sm text-gray-600 mt-1">Becomes:</div>
							<div class="font-serif text-lg">√x</div>
						</div>
						<div class="bg-gray-50 p-3 rounded">
							<div class="text-sm text-gray-600 mb-1">Type this:</div>
							<code class="bg-white px-2 py-1 rounded text-sm">pi times r squared</code>
							<div class="text-sm text-gray-600 mt-1">Becomes:</div>
							<div class="font-serif text-lg">π × r²</div>
						</div>
						<div class="bg-gray-50 p-3 rounded">
							<div class="text-sm text-gray-600 mb-1">Type this:</div>
							<code class="bg-white px-2 py-1 rounded text-sm">alpha beta gamma</code>
							<div class="text-sm text-gray-600 mt-1">Becomes:</div>
							<div class="font-serif text-lg">αβγ</div>
						</div>
					</div>
				</div>

				<!-- Advanced Examples -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">🔬 Advanced Mathematics</h4>
					<div class="space-y-4">
						<div class="bg-purple-50 p-4 rounded">
							<div class="text-sm text-purple-700 mb-2">Integrals:</div>
							<code class="bg-white px-3 py-2 rounded text-sm block mb-2">integral from 0 to infinity of e to the power of negative x squared dx</code>
							<div class="text-sm text-purple-700 mb-1">Becomes:</div>
							<div class="font-serif text-xl">∫₀^∞ e^(-x²) dx</div>
						</div>
						<div class="bg-green-50 p-4 rounded">
							<div class="text-sm text-green-700 mb-2">Summations:</div>
							<code class="bg-white px-3 py-2 rounded text-sm block mb-2">sum from i equals 1 to n of x sub i</code>
							<div class="text-sm text-green-700 mb-1">Becomes:</div>
							<div class="font-serif text-xl">∑_&#123;i=1&#125;^n x_i</div>
						</div>

						<div class="bg-orange-50 p-4 rounded">
							<div class="text-sm text-orange-700 mb-2">Complex fractions:</div>
							<code class="bg-white px-3 py-2 rounded text-sm block mb-2">fraction x plus y over z minus w</code>
							<div class="text-sm text-orange-700 mb-1">Becomes:</div>
							<div class="font-serif text-xl">x+y/z-w</div>
						</div>
					</div>
				</div>

				<!-- Quick Reference -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">🚀 Quick Reference</h4>
					<div class="grid grid-cols-2 md:grid-cols-3 gap-3 text-sm">
						<div><strong>fractions:</strong> "fraction a over b"</div>
						<div><strong>roots:</strong> "square root of", "cube root of"</div>
						<div><strong>integrals:</strong> "integral of", "integral from a to b of"</div>
						<div><strong>summations:</strong> "sum from i=1 to n of"</div>
						<div><strong>greek:</strong> "alpha", "beta", "gamma", "pi", "sigma"</div>
						<div><strong>operators:</strong> "times", "divided by", "plus or minus"</div>
						<div><strong>superscripts:</strong> "squared", "cubed", "to the power of"</div>
						<div><strong>subscripts:</strong> "sub", "to the"</div>
						<div><strong>symbols:</strong> "therefore", "because", "infinity"</div>
					</div>
				</div>

				<!-- Try It Section -->
				<div class="bg-yellow-50 p-4 rounded">
					<h4 class="text-lg font-medium text-yellow-900 mb-2">🎯 Try It Now!</h4>
					<p class="text-yellow-800 mb-3">
						Create a new math block and try typing any of these expressions:
					</p>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm">
						<button class="text-left p-2 bg-white rounded hover:bg-gray-50" on:click={() => dispatch('closeMathHelp')}>
							"fraction x squared plus y squared over z"
						</button>
						<button class="text-left p-2 bg-white rounded hover:bg-gray-50" on:click={() => dispatch('closeMathHelp')}>
							"integral from negative infinity to infinity of e to the power of negative x squared over square root of pi dx"
						</button>
						<button class="text-left p-2 bg-white rounded hover:bg-gray-50" on:click={() => dispatch('closeMathHelp')}>
							"sum from k equals 0 to infinity of fraction 1 over k factorial"
						</button>
						<button class="text-left p-2 bg-white rounded hover:bg-gray-50" on:click={() => dispatch('closeMathHelp')}>
							"limit as x approaches 0 of fraction sin x over x equals 1"
						</button>
					</div>
				</div>

				<!-- Export Options -->
				<div>
					<h4 class="text-lg font-medium text-gray-900 mb-3">📤 Export Your Math</h4>
					<p class="text-gray-600 mb-3">
						Once you've created beautiful math expressions, export them to various formats:
					</p>
					<div class="grid grid-cols-2 md:grid-cols-4 gap-3">
						<div class="text-center p-3 bg-red-50 rounded">
							<div class="text-red-600 font-medium">PDF</div>
							<div class="text-xs text-red-600">Vector graphics</div>
						</div>
						<div class="text-center p-3 bg-blue-50 rounded">
							<div class="text-blue-600 font-medium">LaTeX</div>
							<div class="text-xs text-blue-600">Source code</div>
						</div>
						<div class="text-center p-3 bg-green-50 rounded">
							<div class="text-green-600 font-medium">MathML</div>
							<div class="text-xs text-green-600">Web standard</div>
						</div>
						<div class="text-center p-3 bg-purple-50 rounded">
							<div class="text-purple-600 font-medium">SVG</div>
							<div class="text-xs text-purple-600">Scalable images</div>
						</div>
					</div>
				</div>
			</div>

			<div class="mt-6 flex justify-end">
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
					on:click={() => dispatch('closeMathHelp')}
				>
					Got it!
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- ── Find & Replace Modal ────────────────────────────────────── -->
{#if showFindReplaceDialog}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-lg">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">Find & Replace</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => dispatch('closeFindReplace')}
				>
					×
				</button>
			</div>

			<div class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Find</label>
					<div class="flex space-x-2">
						<input
							type="text"
							bind:value={findText}
							placeholder="Enter text to find..."
							class="flex-1 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
							on:keydown={(e) => {
								if (e.key === 'Enter') dispatch('performFind');
								if (e.key === 'Escape') dispatch('closeFindReplace');
							}}
						/>
						<button
							class="px-3 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200"
							class:bg-blue-100={isCaseSensitive}
							class:text-blue-700={isCaseSensitive}
							on:click={() => dispatch('toggleCaseSensitive')}
							title="Case sensitive"
						>
							Aa
						</button>
					</div>
				</div>

				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Replace with</label>
					<input
						type="text"
						bind:value={replaceText}
						placeholder="Enter replacement text..."
						class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
						on:keydown={(e) => {
							if (e.key === 'Enter') dispatch('performReplace');
							if (e.key === 'Escape') dispatch('closeFindReplace');
						}}
					/>
				</div>

				{#if findResults.length > 0}
					<div class="flex items-center justify-between bg-blue-50 p-3 rounded">
						<div class="text-sm text-blue-800">
							Found {findResults.length} occurrence{findResults.length !== 1 ? 's' : ''}
							{#if currentFindIndex >= 0}
								<span class="font-medium">(currently at {currentFindIndex + 1})</span>
							{/if}
						</div>
						<div class="flex space-x-2">
							<button
								class="px-2 py-1 text-sm bg-white text-gray-700 rounded hover:bg-gray-100"
								on:click={() => dispatch('findPrevious')}
								title="Previous (Shift+Enter)"
							>
								◀ Previous
							</button>
							<button
								class="px-2 py-1 text-sm bg-white text-gray-700 rounded hover:bg-gray-100"
								on:click={() => dispatch('findNext')}
								title="Next (Enter)"
							>
								Next ▶
							</button>
						</div>
					</div>
				{:else if findText}
					<div class="text-sm text-gray-500 bg-gray-50 p-3 rounded">
						No matches found
					</div>
				{/if}
			</div>

			<div class="mt-6 flex justify-end space-x-3">
				<button
					class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50"
					on:click={() => dispatch('closeFindReplace')}
				>
					Close
				</button>
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
					on:click={() => dispatch('performFind')}
					disabled={!findText}
				>
					Find
				</button>
				{#if findResults.length > 0}
					<button
						class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
						on:click={() => dispatch('performReplace')}
					>
						Replace
					</button>
					<button
						class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700"
						on:click={() => dispatch('replaceAll')}
					>
						Replace All
					</button>
				{/if}
			</div>
		</div>
	</div>
{/if}


<!-- ── AI Writing Assistant Modal ─────────────────────────────── -->
{#if showAIAssistant}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">✨ AI Writing Suggestions</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => dispatch('closeAIAssistant')}
				>
					×
				</button>
			</div>

			{#if aiSuggestions.length === 0}
				<div class="text-center py-8">
					<div class="text-4xl mb-4">💭</div>
					<p class="text-gray-600">Generating suggestions...</p>
				</div>
			{:else}
				<div class="space-y-4">
					{#each aiSuggestions as suggestion, index}
						<div class="border border-gray-200 rounded-lg p-4">
							<div class="flex items-start justify-between mb-2">
								<h4 class="font-medium text-gray-900">Suggestion {index + 1}</h4>
								<button
									class="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
									on:click={() => dispatch('applyAISuggestion', suggestion)}
								>
									Apply
								</button>
							</div>
							<p class="text-gray-700 whitespace-pre-wrap">{suggestion}</p>
						</div>
					{/each}
				</div>
			{/if}

			<div class="mt-6 flex justify-end">
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => dispatch('closeAIAssistant')}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- ── Document Summary Modal ─────────────────────────────────── -->
{#if documentSummary}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-3xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">📄 Document Summary</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => dispatch('closeSummary')}
				>
					×
				</button>
			</div>

			<div class="prose prose-sm max-w-none">
				<div class="whitespace-pre-wrap text-gray-700">{documentSummary}</div>
			</div>

			<div class="mt-6 flex justify-end space-x-3">
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
					on:click={() => navigator.clipboard.writeText(documentSummary)}
				>
					Copy to Clipboard
				</button>
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => dispatch('closeSummary')}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- ── Text Analysis Modal ─────────────────────────────────────── -->
{#if showTextAnalysis}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-semibold text-gray-900">📊 Text Analysis</h3>
				<button
					class="text-gray-400 hover:text-gray-600 text-2xl"
					on:click={() => dispatch('closeTextAnalysis')}
				>
					×
				</button>
			</div>

			{#if textAnalysis}
				<div class="grid grid-cols-2 gap-6">
					<div class="space-y-4">
						<h4 class="font-medium text-gray-900">Statistics</h4>
						<div class="space-y-2 text-sm">
							<div class="flex justify-between">
								<span class="text-gray-600">Words:</span>
								<span class="font-medium">{textAnalysis.word_count || 0}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Characters:</span>
								<span class="font-medium">{textAnalysis.char_count || 0}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Sentences:</span>
								<span class="font-medium">{textAnalysis.sentence_count || 0}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Paragraphs:</span>
								<span class="font-medium">{textAnalysis.paragraph_count || 0}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Reading Time:</span>
								<span class="font-medium">{textAnalysis.reading_time || 'N/A'}</span>
							</div>
						</div>
					</div>

					<div class="space-y-4">
						<h4 class="font-medium text-gray-900">Readability</h4>
						<div class="space-y-2 text-sm">
							<div class="flex justify-between">
								<span class="text-gray-600">Grade Level:</span>
								<span class="font-medium">{textAnalysis.grade_level || 'N/A'}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Readability:</span>
								<span class="font-medium">{textAnalysis.readability_score || 'N/A'}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-600">Complexity:</span>
								<span class="font-medium">{textAnalysis.complexity || 'N/A'}</span>
							</div>
						</div>
					</div>
				</div>

				{#if textAnalysis.key_phrases && textAnalysis.key_phrases.length > 0}
					<div class="mt-6">
						<h4 class="font-medium text-gray-900 mb-3">Key Phrases</h4>
						<div class="flex flex-wrap gap-2">
							{#each textAnalysis.key_phrases as phrase}
								<span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
									{phrase}
								</span>
							{/each}
						</div>
					</div>
				{/if}

				{#if textAnalysis.sentiment}
					<div class="mt-6">
						<h4 class="font-medium text-gray-900 mb-3">Sentiment Analysis</h4>
						<div class="flex items-center space-x-4">
							<div class="flex-1 bg-gray-200 rounded-full h-2">
								<div
									class="bg-green-500 h-2 rounded-full"
									style="width: {textAnalysis.sentiment.positive || 0}%"
								></div>
							</div>
							<span class="text-sm text-gray-600">
								Positive: {textAnalysis.sentiment.positive || 0}%
							</span>
						</div>
						<div class="flex items-center space-x-4 mt-2">
							<div class="flex-1 bg-gray-200 rounded-full h-2">
								<div
									class="bg-red-500 h-2 rounded-full"
									style="width: {textAnalysis.sentiment.negative || 0}%"
								></div>
							</div>
							<span class="text-sm text-gray-600">
								Negative: {textAnalysis.sentiment.negative || 0}%
							</span>
						</div>
					</div>
				{/if}
			{:else}
				<div class="text-center py-8">
					<p class="text-gray-600">No analysis data available</p>
				</div>
			{/if}

			<div class="mt-6 flex justify-end">
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
					on:click={() => dispatch('closeTextAnalysis')}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.prose {
		color: #374151;
	}
</style>
