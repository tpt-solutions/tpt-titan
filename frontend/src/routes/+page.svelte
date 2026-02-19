<script>
	import AppSelector from '$lib/components/AppSelector.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	// Accept framework-provided props to avoid warnings
	export const data = null;

	export const form = null;

	let showAppSelector = true;

	const appRoutes = {
		spreadsheet: '/spreadsheet',
		forms: '/forms',
		editor: '/editor',
		tasks: '/tasks',
		calendar: '/calendar',
		contacts: '/contacts',
		email: '/email',
		database: '/database'
	};

	onMount(() => {
		// Check if user has a saved default app preference
		const savedApp = localStorage.getItem('tpt-default-app');
		const urlParams = new URLSearchParams(window.location.search);
		const forceSelector = urlParams.has('select-app');

		if (savedApp && !forceSelector && appRoutes[savedApp]) {
			// Auto-redirect to saved default app
			goto(appRoutes[savedApp]);
			showAppSelector = false;
		}
	});
</script>

<svelte:head>
	<title>TPT Titan - Open Source Office Suite</title>
	<meta name="description" content="Complete open source alternative to Microsoft Office 365" />
</svelte:head>

{#if showAppSelector}
	<div class="min-h-screen bg-gray-50 py-8">
		<div class="container mx-auto px-4">
			<header class="text-center mb-8">
				<h1 class="text-5xl font-bold text-gray-900 mb-4">
					Welcome to <span class="text-blue-600">TPT Titan</span>
				</h1>
				<p class="text-xl text-gray-600 max-w-2xl mx-auto mb-6">
					A complete open source alternative to Microsoft Office 365.
					Built with modern web technologies for productivity, collaboration, and privacy.
				</p>
			</header>

			<AppSelector />

			<div class="text-center mt-8">
				<div class="bg-gray-50 border border-gray-200 rounded-lg p-6">
					<h3 class="text-lg font-semibold text-gray-900 mb-2">🚀 Performance Optimized</h3>
					<p class="text-gray-600 mb-4">
						TPT Titan uses lazy loading to start up to 7x faster. Choose your default app above to begin.
					</p>
					<div class="flex justify-center space-x-4">
						<a href="/docs" class="text-blue-600 hover:text-blue-800 underline">
							View Documentation
						</a>
						<a href="https://github.com/tpt-titan/tpt-titan" class="text-gray-600 hover:text-gray-800 underline">
							GitHub Repository
						</a>
						<a href="?select-app" class="text-gray-600 hover:text-gray-800 underline">
							Change Default App
						</a>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
