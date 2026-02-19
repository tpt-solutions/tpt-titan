<!-- frontend/src/lib/components/FormBuilderTemplatesModal.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	const formTemplates = [
		{
			id: 'contact',
			name: 'Contact Form',
			description: 'Basic contact information collection',
			icon: '📞',
			fields: [
				{ type: 'text',     label: 'Full Name',      properties: { required: true } },
				{ type: 'email',    label: 'Email Address',  properties: { required: true } },
				{ type: 'phone',    label: 'Phone Number',   properties: { required: false } },
				{ type: 'textarea', label: 'Message',        properties: { required: true, rows: 4 } }
			]
		},
		{
			id: 'survey',
			name: 'Customer Survey',
			description: 'Gather feedback from customers',
			icon: '📊',
			fields: [
				{ type: 'text',     label: 'Name',                    properties: { required: false } },
				{ type: 'email',    label: 'Email',                   properties: { required: false } },
				{ type: 'radio',    label: 'How satisfied are you?',  properties: { options: ['Very Satisfied', 'Satisfied', 'Neutral', 'Dissatisfied', 'Very Dissatisfied'], required: true } },
				{ type: 'textarea', label: 'Comments',                properties: { required: false, rows: 3 } }
			]
		},
		{
			id: 'registration',
			name: 'Event Registration',
			description: 'Register attendees for events',
			icon: '🎟️',
			fields: [
				{ type: 'text',     label: 'Full Name',       properties: { required: true } },
				{ type: 'email',    label: 'Email Address',   properties: { required: true } },
				{ type: 'phone',    label: 'Phone Number',    properties: { required: true } },
				{ type: 'select',   label: 'Ticket Type',     properties: { options: ['General Admission', 'VIP', 'Student'], required: true } },
				{ type: 'number',   label: 'Number of Tickets', properties: { required: true, min: 1, max: 10 } },
				{ type: 'textarea', label: 'Special Requests', properties: { required: false, rows: 2 } }
			]
		},
		{
			id: 'feedback',
			name: 'Product Feedback',
			description: 'Collect feedback on products or services',
			icon: '💬',
			fields: [
				{ type: 'text',     label: 'Product Name',           properties: { required: true } },
				{ type: 'rating',   label: 'Overall Rating',         properties: { required: true, maxRating: 5 } },
				{ type: 'radio',    label: 'Would you recommend?',   properties: { options: ['Yes', 'No', 'Maybe'], required: true } },
				{ type: 'textarea', label: 'What did you like?',     properties: { required: false, rows: 2 } },
				{ type: 'textarea', label: 'What could be improved?', properties: { required: false, rows: 2 } }
			]
		}
	];

	function apply(template) {
		dispatch('apply', template);
	}
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
	<div class="bg-white rounded-lg w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
		<div class="flex items-center justify-between p-6 border-b border-gray-200">
			<div>
				<h2 class="text-2xl font-bold text-gray-900">Choose a Form Template</h2>
				<p class="text-gray-600 mt-1">Start with a pre-built template and customise it for your needs</p>
			</div>
			<button class="text-gray-400 hover:text-gray-600 text-2xl" on:click={() => dispatch('close')}>×</button>
		</div>

		<div class="flex-1 overflow-y-auto p-6">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				{#each formTemplates as template}
					<div
						class="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-lg transition-all cursor-pointer group"
						on:click={() => apply(template)}
						on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); apply(template); } }}
						tabindex="0"
						role="button"
						aria-label="Select {template.name} template"
						aria-pressed="false"
					>
						<div class="flex items-center mb-4">
							<div class="text-3xl mr-3">{template.icon}</div>
							<div>
								<h3 class="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">{template.name}</h3>
								<p class="text-sm text-gray-600">{template.description}</p>
							</div>
						</div>
						<div class="space-y-2">
							{#each template.fields.slice(0, 4) as field}
								<div class="flex items-center text-sm text-gray-600">
									<span class="w-4 h-4 mr-2 bg-gray-200 rounded flex items-center justify-center text-xs">
										{field.type === 'text' ? '📝' : field.type === 'email' ? '📧' : field.type === 'radio' ? '⭕' : '📄'}
									</span>
									<span>{field.label}</span>
									{#if field.properties.required}
										<span class="text-red-500 ml-1">*</span>
									{/if}
								</div>
							{/each}
							{#if template.fields.length > 4}
								<div class="text-sm text-gray-500">+{template.fields.length - 4} more fields</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		</div>

		<div class="border-t border-gray-200 p-6 bg-gray-50">
			<div class="flex items-center justify-between text-sm text-gray-600">
				<span>{formTemplates.length} templates available</span>
				<span>💡 Templates can be fully customized after selection</span>
			</div>
		</div>
	</div>
</div>
