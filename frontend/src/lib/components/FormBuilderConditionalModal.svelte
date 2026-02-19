<!-- frontend/src/lib/components/FormBuilderConditionalModal.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';

	export let selectedField = null;
	export let allFields = [];

	const dispatch = createEventDispatcher();

	function addRule() {
		if (!selectedField) return;
		if (!selectedField.conditionalLogic) selectedField.conditionalLogic = [];
		selectedField.conditionalLogic.push({ field: '', operator: 'equals', value: '', action: 'show' });
		dispatch('updated');
	}

	function updateRule(index, property, value) {
		if (selectedField?.conditionalLogic?.[index]) {
			selectedField.conditionalLogic[index][property] = value;
			dispatch('updated');
		}
	}

	function removeRule(index) {
		if (selectedField?.conditionalLogic) {
			selectedField.conditionalLogic.splice(index, 1);
			dispatch('updated');
		}
	}

	$: otherFields = allFields.filter(f => f.id !== selectedField?.id);
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
	<div class="bg-white rounded-lg w-full max-w-2xl max-h-[90vh] overflow-hidden flex flex-col">
		<div class="flex items-center justify-between p-6 border-b border-gray-200">
			<div>
				<h2 class="text-2xl font-bold text-gray-900">Conditional Logic</h2>
				<p class="text-gray-600 mt-1">Show/hide fields based on responses to other fields</p>
			</div>
			<button class="text-gray-400 hover:text-gray-600 text-2xl" on:click={() => dispatch('close')}>×</button>
		</div>

		<div class="flex-1 overflow-y-auto p-6">
			{#if selectedField}
				<div class="space-y-6">
					<div class="bg-blue-50 p-4 rounded-lg">
						<h3 class="font-medium text-blue-900 mb-2">Field: {selectedField.label}</h3>
						<p class="text-sm text-blue-700">Set up rules to show or hide this field based on responses to other questions.</p>
					</div>

					{#if selectedField.conditionalLogic && selectedField.conditionalLogic.length > 0}
						<div class="space-y-4">
							{#each selectedField.conditionalLogic as rule, index}
								<div class="border border-gray-200 rounded-lg p-4">
									<div class="flex items-center justify-between mb-3">
										<h4 class="font-medium text-gray-900">Rule {index + 1}</h4>
										<button class="text-red-500 hover:text-red-700" on:click={() => removeRule(index)}>Remove</button>
									</div>
									<div class="grid grid-cols-1 md:grid-cols-4 gap-3">
										<div>
											<label for="rule-field-{index}" class="block text-sm font-medium text-gray-700 mb-1">If field</label>
											<select
												id="rule-field-{index}"
												bind:value={rule.field}
												on:change={(e) => updateRule(index, 'field', e.target.value)}
												class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
											>
												<option value="">Select field...</option>
												{#each otherFields as field}
													<option value={field.id}>{field.label}</option>
												{/each}
											</select>
										</div>
										<div>
											<label for="rule-operator-{index}" class="block text-sm font-medium text-gray-700 mb-1">Operator</label>
											<select
												id="rule-operator-{index}"
												bind:value={rule.operator}
												on:change={(e) => updateRule(index, 'operator', e.target.value)}
												class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
											>
												<option value="equals">Equals</option>
												<option value="not_equals">Not equals</option>
												<option value="contains">Contains</option>
												<option value="not_contains">Does not contain</option>
												<option value="greater_than">Greater than</option>
												<option value="less_than">Less than</option>
											</select>
										</div>
										<div>
											<label for="rule-value-{index}" class="block text-sm font-medium text-gray-700 mb-1">Value</label>
											<input
												id="rule-value-{index}"
												type="text"
												bind:value={rule.value}
												on:input={(e) => updateRule(index, 'value', e.target.value)}
												placeholder="Enter value..."
												class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
											/>
										</div>
										<div>
											<label for="rule-action-{index}" class="block text-sm font-medium text-gray-700 mb-1">Then</label>
											<select
												id="rule-action-{index}"
												bind:value={rule.action}
												on:change={(e) => updateRule(index, 'action', e.target.value)}
												class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
											>
												<option value="show">Show field</option>
												<option value="hide">Hide field</option>
											</select>
										</div>
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<div class="text-center py-8 text-gray-500">
							<p>No conditional rules set up yet.</p>
							<p class="text-sm mt-1">Add a rule below to show/hide this field based on other responses.</p>
						</div>
					{/if}

					<div class="flex justify-center">
						<button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors" on:click={addRule}>
							➕ Add Rule
						</button>
					</div>
				</div>
			{/if}
		</div>

		<div class="border-t border-gray-200 p-6 bg-gray-50">
			<div class="flex justify-end space-x-3">
				<button class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50" on:click={() => dispatch('close')}>
					Cancel
				</button>
				<button class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700" on:click={() => dispatch('close')}>
					Save Logic
				</button>
			</div>
		</div>
	</div>
</div>
