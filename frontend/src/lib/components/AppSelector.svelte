<script>
	import { userPreferences } from '../stores.js';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let selectedApp = 'spreadsheet';
	let isLoading = false;

	const apps = [
		{
			id: 'spreadsheet',
			name: 'Spreadsheet',
			description: 'Powerful spreadsheet with formulas and charts',
			icon: '📊',
			color: 'blue',
			route: '/spreadsheet'
		},
		{
			id: 'forms',
			name: 'Forms',
			description: 'Create and manage forms with workflows',
			icon: '📋',
			color: 'green',
			route: '/forms'
		},
		{
			id: 'editor',
			name: 'Text Editor',
			description: 'Rich text editing with collaboration',
			icon: '📝',
			color: 'purple',
			route: '/editor'
		},
		{
			id: 'tasks',
			name: 'Tasks',
			description: 'Task management and project tracking',
			icon: '✅',
			color: 'yellow',
			route: '/tasks'
		},
		{
			id: 'calendar',
			name: 'Calendar',
			description: 'Calendar with scheduling and reminders',
			icon: '📅',
			color: 'red',
			route: '/calendar'
		},
		{
			id: 'contacts',
			name: 'Contacts',
			description: 'Contact management and organization',
			icon: '👥',
			color: 'indigo',
			route: '/contacts'
		},
		{
			id: 'email',
			name: 'Email',
			description: 'Email client with encryption',
			icon: '📧',
			color: 'pink',
			route: '/email'
		},
		{
			id: 'database',
			name: 'Database',
			description: 'Database management and queries',
			icon: '🗄️',
			color: 'orange',
			route: '/database'
		}
	];

	onMount(() => {
		// Load saved preference
		const saved = localStorage.getItem('tpt-default-app');
		if (saved) {
			selectedApp = saved;
		}
	});

	async function handleAppSelect(appId) {
		selectedApp = appId;
		isLoading = true;

		try {
			// Save preference
			localStorage.setItem('tpt-default-app', appId);

			// Update store
			userPreferences.update(prefs => ({
				...prefs,
				defaultView: appId
			}));

			// Navigate to selected app
			const app = apps.find(a => a.id === appId);
			if (app) {
				await goto(app.route);
			}
		} catch (error) {
			console.error('Error selecting app:', error);
		} finally {
			isLoading = false;
		}
	}

	function getColorClasses(color, isSelected) {
		const baseClasses = 'p-6 rounded-lg shadow-md transition-all duration-200 cursor-pointer border-2';
		if (isSelected) {
			switch (color) {
				case 'blue': return `${baseClasses} bg-blue-50 border-blue-500 text-blue-700`;
				case 'green': return `${baseClasses} bg-green-50 border-green-500 text-green-700`;
				case 'purple': return `${baseClasses} bg-purple-50 border-purple-500 text-purple-700`;
				case 'yellow': return `${baseClasses} bg-yellow-50 border-yellow-500 text-yellow-700`;
				case 'red': return `${baseClasses} bg-red-50 border-red-500 text-red-700`;
				case 'indigo': return `${baseClasses} bg-indigo-50 border-indigo-500 text-indigo-700`;
				case 'pink': return `${baseClasses} bg-pink-50 border-pink-500 text-pink-700`;
				case 'orange': return `${baseClasses} bg-orange-50 border-orange-500 text-orange-700`;
				default: return `${baseClasses} bg-gray-50 border-gray-500 text-gray-700`;
			}
		} else {
			return `${baseClasses} bg-white border-gray-200 hover:border-gray-300 text-gray-900 hover:shadow-lg`;
		}
	}
</script>

<div class="max-w-6xl mx-auto">
	<div class="text-center mb-8">
		<h2 class="text-3xl font-bold text-gray-900 mb-4">Choose Your Default App</h2>
		<p class="text-lg text-gray-600 max-w-2xl mx-auto">
			Select which application you'd like to load by default. This will significantly improve your startup time by only loading the app you use most.
		</p>
		<div class="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
			<p class="text-sm text-blue-800">
				<strong>Performance Tip:</strong> Choosing one app reduces initial load from 135MB to ~2-3MB, cutting startup time from 35+ seconds to under 5 seconds.
			</p>
		</div>
	</div>

	<div class="grid md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mb-8">
		{#each apps as app}
			<div
				class={getColorClasses(app.color, selectedApp === app.id)}
				on:click={() => handleAppSelect(app.id)}
				role="button"
				tabindex="0"
				on:keydown={(e) => e.key === 'Enter' && handleAppSelect(app.id)}
			>
				<div class="flex flex-col items-center text-center">
					<div class="text-4xl mb-3">{app.icon}</div>
					<h3 class="text-lg font-semibold mb-2">{app.name}</h3>
					<p class="text-sm opacity-80 leading-relaxed">{app.description}</p>
					{#if selectedApp === app.id}
						<div class="mt-3 px-3 py-1 bg-current text-white text-xs rounded-full">
							Default App
						</div>
					{/if}
				</div>
			</div>
		{/each}
	</div>

	<div class="text-center">
		<button
			class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold text-lg shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
			disabled={isLoading}
			on:click={() => handleAppSelect(selectedApp)}
		>
			{#if isLoading}
				<span class="flex items-center">
					<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Loading {apps.find(a => a.id === selectedApp)?.name}...
				</span>
			{:else}
				Launch {apps.find(a => a.id === selectedApp)?.name}
			{/if}
		</button>
		<p class="text-sm text-gray-500 mt-4">
			You can change your default app anytime from the app menu
		</p>
	</div>
</div>
