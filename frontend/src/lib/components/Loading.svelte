<script>
	export let message = 'Loading...';
	export let size = 'medium'; // 'small', 'medium', 'large'
</script>

<div class="flex flex-col items-center justify-center p-8">
	<div class="flex items-center space-x-2 mb-4">
		<svg
			class="animate-spin {size === 'small' ? 'h-4 w-4' : size === 'large' ? 'h-8 w-8' : 'h-6 w-6'} text-blue-600"
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
		>
			<circle
				class="opacity-25"
				cx="12"
				cy="12"
				r="10"
				stroke="currentColor"
				stroke-width="4"
			></circle>
			<path
				class="opacity-75"
				fill="currentColor"
				d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
			></path>
		</svg>
		<span class="text-gray-600 {size === 'small' ? 'text-sm' : size === 'large' ? 'text-lg' : 'text-base'}">
			{message}
		</span>
	</div>

	{#if size === 'large'}
		<div class="w-64 bg-gray-200 rounded-full h-2 mb-4">
			<div class="bg-blue-600 h-2 rounded-full animate-pulse" style="width: 60%"></div>
		</div>
		<p class="text-sm text-gray-500 text-center">
			Loading application components...
		</p>
	{/if}
</div>
