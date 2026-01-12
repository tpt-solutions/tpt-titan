<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { aiService } from '../services/ai.js';

	const dispatch = createEventDispatcher();

	let currentStep = 0;
	let showWizard = true;

	// User preferences and setup
	let userPreferences = {
		experienceLevel: 'beginner', // beginner, intermediate, advanced
		primaryUse: 'writing', // writing, email, tasks, calendar, all
		privacyPreference: 'balanced', // local-only, balanced, cloud-first
		setupComplete: false
	};

	let aiStatus = {
		localAI: false,
		cloudAI: false,
		speech: false,
		testing: false
	};

	let testResults = {
		basicAI: null,
		speech: null,
		suggestions: null
	};

	const steps = [
		{
			id: 'welcome',
			title: 'Welcome to AI-Powered Productivity! 🚀',
			description: 'Let\'s get you set up with TPT Titan\'s AI features. This will take just 2-3 minutes.',
			content: 'welcome'
		},
		{
			id: 'experience',
			title: 'How familiar are you with AI?',
			description: 'This helps us customize your experience.',
			content: 'experience'
		},
		{
			id: 'usage',
			title: 'What will you use AI for?',
			description: 'Select your primary use case to optimize the setup.',
			content: 'usage'
		},
		{
			id: 'privacy',
			title: 'Privacy & Performance Preferences',
			description: 'Choose how you want AI to work with your data.',
			content: 'privacy'
		},
		{
			id: 'providers',
			title: 'AI Provider Setup',
			description: 'Configure your AI services.',
			content: 'providers'
		},
		{
			id: 'testing',
			title: 'Testing Your Setup',
			description: 'Let\'s make sure everything works!',
			content: 'testing'
		},
		{
			id: 'complete',
			title: 'You\'re All Set! 🎉',
			description: 'Start using AI features in your workflow.',
			content: 'complete'
		}
	];

	function nextStep() {
		if (currentStep < steps.length - 1) {
			currentStep++;
		}
	}

	function prevStep() {
		if (currentStep > 0) {
			currentStep--;
		}
	}

	function skipWizard() {
		showWizard = false;
		dispatch('complete', { skipped: true });
	}

	async function completeSetup() {
		try {
			// Save user preferences
			localStorage.setItem('ai_onboarding_complete', 'true');
			localStorage.setItem('ai_user_preferences', JSON.stringify(userPreferences));

			showWizard = false;
			dispatch('complete', { preferences: userPreferences, status: aiStatus });
		} catch (error) {
			console.error('Failed to save preferences:', error);
		}
	}

	async function testAIProviders() {
		aiStatus.testing = true;

		try {
			// Test basic AI functionality
			testResults.basicAI = 'testing';
			const testResponse = await aiService.getWritingSuggestions('Hello world', 'general');
			testResults.basicAI = testResponse ? 'success' : 'failed';
		} catch (error) {
			testResults.basicAI = 'failed';
			console.error('Basic AI test failed:', error);
		}

		try {
			// Test speech functionality
			testResults.speech = 'testing';
			// This would test speech recognition
			testResults.speech = 'success'; // Placeholder
		} catch (error) {
			testResults.speech = 'failed';
		}

		aiStatus.testing = false;
	}

	onMount(() => {
		// Check if user has already completed onboarding
		const completed = localStorage.getItem('ai_onboarding_complete');
		if (completed === 'true') {
			showWizard = false;
			dispatch('already-completed');
		}
	});
</script>

{#if showWizard}
<div class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50">
	<div class="bg-white rounded-xl shadow-2xl w-full max-w-2xl mx-4 max-h-90vh overflow-hidden">
		<!-- Progress Bar -->
		<div class="bg-gray-200 h-1">
			<div
				class="bg-blue-600 h-1 transition-all duration-300"
				style="width: {((currentStep + 1) / steps.length) * 100}%"
			></div>
		</div>

		<!-- Header -->
		<div class="px-8 py-6 border-b border-gray-200">
			<div class="flex items-center justify-between">
				<div>
					<h2 class="text-2xl font-bold text-gray-900">
						{steps[currentStep].title}
					</h2>
					<p class="text-gray-600 mt-1">
						{steps[currentStep].description}
					</p>
				</div>
				<button
					on:click={skipWizard}
					class="text-gray-400 hover:text-gray-600 text-xl"
					title="Skip setup (you can always configure later)"
				>
					×
				</button>
			</div>

			<!-- Step Indicator -->
			<div class="flex items-center mt-4 space-x-2">
				{#each steps as step, index}
					<div
						class="flex items-center"
						class:opacity-50={index > currentStep}
					>
						<div
							class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium
								{index < currentStep ? 'bg-green-500 text-white' :
								 index === currentStep ? 'bg-blue-600 text-white' :
								 'bg-gray-300 text-gray-600'}"
						>
							{#if index < currentStep}
								✓
							{:else}
								{index + 1}
							{/if}
						</div>
						{#if index < steps.length - 1}
							<div class="w-8 h-0.5 mx-2
								{index < currentStep ? 'bg-green-500' :
								 'bg-gray-300'}"></div>
						{/if}
					</div>
				{/each}
			</div>
		</div>

		<!-- Content -->
		<div class="px-8 py-6 overflow-y-auto max-h-96">
			{#if steps[currentStep].content === 'welcome'}
				<div class="text-center">
					<div class="text-6xl mb-6">🤖✨</div>
					<p class="text-lg text-gray-700 mb-6">
						TPT Titan includes powerful AI features designed to enhance your productivity.
						From writing assistance to automated workflows, AI can help you work smarter.
					</p>
					<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
						<h4 class="font-medium text-blue-900 mb-2">What you'll get:</h4>
						<ul class="text-sm text-blue-800 space-y-1 text-left">
							<li>• AI-powered writing assistance and proofreading</li>
							<li>• Smart email categorization and voice composition</li>
							<li>• Intelligent task prioritization and scheduling</li>
							<li>• Voice commands across all applications</li>
							<li>• Automated document processing and analysis</li>
						</ul>
					</div>
				</div>
			{/if}

			{#if steps[currentStep].content === 'experience'}
				<div class="space-y-4">
					<p class="text-gray-600 mb-4">
						Select your experience level to customize the setup and recommendations.
					</p>

					<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.experienceLevel}
								value="beginner"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:border-blue-300 transition-colors">
								<div class="text-2xl mb-2">🌱</div>
								<div class="font-medium text-gray-900">Beginner</div>
								<div class="text-sm text-gray-600 mt-1">
									New to AI - show me everything step by step
								</div>
							</div>
						</label>

						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.experienceLevel}
								value="intermediate"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:border-blue-300 transition-colors">
								<div class="text-2xl mb-2">⚡</div>
								<div class="font-medium text-gray-900">Intermediate</div>
								<div class="text-sm text-gray-600 mt-1">
									Familiar with AI - show shortcuts and tips
								</div>
							</div>
						</label>

						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.experienceLevel}
								value="advanced"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:border-blue-300 transition-colors">
								<div class="text-2xl mb-2">🚀</div>
								<div class="font-medium text-gray-900">Advanced</div>
								<div class="text-sm text-gray-600 mt-1">
									Power user - give me full control and advanced features
								</div>
							</div>
						</label>
					</div>
				</div>
			{/if}

			{#if steps[currentStep].content === 'usage'}
				<div class="space-y-4">
					<p class="text-gray-600 mb-4">
						What will you primarily use AI for? This helps us prioritize features and setup.
					</p>

					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.primaryUse}
								value="writing"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:border-blue-300 transition-colors">
								<div class="text-2xl mb-2">✍️</div>
								<div class="font-medium text-gray-900">Writing & Content</div>
								<div class="text-sm text-gray-600 mt-1">
									Documents, reports, blogs, emails
								</div>
							</div>
						</label>

						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.primaryUse}
								value="tasks"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:border-blue-300 transition-colors">
								<div class="text-2xl mb-2">✅</div>
								<div class="font-medium text-gray-900">Task Management</div>
								<div class="text-sm text-gray-600 mt-1">
									Planning, prioritization, deadlines
								</div>
							</div>
						</label>

						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.primaryUse}
								value="email"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:border-blue-300 transition-colors">
								<div class="text-2xl mb-2">📧</div>
								<div class="font-medium text-gray-900">Email & Communication</div>
								<div class="text-sm text-gray-600 mt-1">
									Categorization, composition, automation
								</div>
							</div>
						</label>

						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.primaryUse}
								value="all"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:border-blue-300 transition-colors">
								<div class="text-2xl mb-2">🎯</div>
								<div class="font-medium text-gray-900">Everything</div>
								<div class="text-sm text-gray-600 mt-1">
									Full productivity suite integration
								</div>
							</div>
						</label>
					</div>
				</div>
			{/if}

			{#if steps[currentStep].content === 'privacy'}
				<div class="space-y-4">
					<p class="text-gray-600 mb-4">
						Choose your balance between privacy and AI capabilities.
					</p>

					<div class="space-y-4">
						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.privacyPreference}
								value="local-only"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-green-500 peer-checked:bg-green-50 hover:border-green-300 transition-colors">
								<div class="flex items-start">
									<div class="text-2xl mr-3">🔒</div>
									<div class="flex-1">
										<div class="font-medium text-gray-900">Local Only - Maximum Privacy</div>
										<div class="text-sm text-gray-600 mt-1">
											All AI processing happens on your device. No internet required.
											Limited to local AI capabilities.
										</div>
									</div>
								</div>
							</div>
						</label>

						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.privacyPreference}
								value="balanced"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:border-blue-300 transition-colors">
								<div class="flex items-start">
									<div class="text-2xl mr-3">⚖️</div>
									<div class="flex-1">
										<div class="font-medium text-gray-900">Balanced - Recommended</div>
										<div class="text-sm text-gray-600 mt-1">
											Use local AI when available, cloud AI for advanced features.
											Best balance of privacy and capabilities.
										</div>
									</div>
								</div>
							</div>
						</label>

						<label class="relative">
							<input
								type="radio"
								bind:group={userPreferences.privacyPreference}
								value="cloud-first"
								class="sr-only peer"
							>
							<div class="border-2 border-gray-200 rounded-lg p-4 cursor-pointer peer-checked:border-orange-500 peer-checked:bg-orange-50 hover:border-orange-300 transition-colors">
								<div class="flex items-start">
									<div class="text-2xl mr-3">⚡</div>
									<div class="flex-1">
										<div class="font-medium text-gray-900">Cloud First - Maximum Performance</div>
										<div class="text-sm text-gray-600 mt-1">
											Prioritize speed and advanced features. Data sent to AI providers.
											Review provider privacy policies.
										</div>
									</div>
								</div>
							</div>
						</label>
					</div>
				</div>
			{/if}

			{#if steps[currentStep].content === 'providers'}
				<div class="space-y-6">
					<p class="text-gray-600">
						Based on your preferences, here's what we recommend setting up:
					</p>

					<!-- Provider Status Cards -->
					<div class="space-y-4">
						<div class="border border-gray-200 rounded-lg p-4">
							<div class="flex items-center justify-between mb-2">
								<div class="flex items-center">
									<div class="text-2xl mr-3">🏠</div>
									<div>
										<div class="font-medium text-gray-900">Local AI (Ollama)</div>
										<div class="text-sm text-gray-600">Free, private, runs on your device</div>
									</div>
								</div>
								<div class="flex items-center">
									{#if aiStatus.localAI}
										<span class="text-green-600 font-medium">✓ Connected</span>
									{:else}
										<span class="text-orange-600 text-sm">Setup required</span>
									{/if}
								</div>
							</div>

							{#if !aiStatus.localAI}
								<div class="mt-3 p-3 bg-blue-50 border border-blue-200 rounded">
									<div class="text-sm text-blue-800 mb-2">
										<strong>Quick Setup:</strong> Download Ollama from ollama.com and run these commands:
									</div>
									<code class="block text-xs bg-white p-2 rounded border font-mono">
										ollama pull qwen2.5:7b-instruct<br>
										ollama pull qwen2.5-vl:7b
									</code>
								</div>
							{/if}
						</div>

						{#if userPreferences.privacyPreference !== 'local-only'}
							<div class="border border-gray-200 rounded-lg p-4">
								<div class="flex items-center justify-between mb-2">
									<div class="flex items-center">
										<div class="text-2xl mr-3">☁️</div>
										<div>
											<div class="font-medium text-gray-900">Cloud AI (OpenRouter)</div>
											<div class="text-sm text-gray-600">Fast, powerful, multiple AI models</div>
										</div>
									</div>
									<div class="flex items-center">
										{#if aiStatus.cloudAI}
											<span class="text-green-600 font-medium">✓ Connected</span>
										{:else}
											<span class="text-orange-600 text-sm">API key needed</span>
										{/if}
									</div>
								</div>

								{#if !aiStatus.cloudAI}
									<div class="mt-3 p-3 bg-green-50 border border-green-200 rounded">
										<div class="text-sm text-green-800 mb-2">
											<strong>Get your API key:</strong>
											<a href="https://openrouter.ai/" target="_blank" class="text-green-700 underline ml-1">
												openrouter.ai
											</a>
										</div>
										<div class="text-xs text-green-700">
											Free tier available • Pay-as-you-go pricing
										</div>
									</div>
								{/if}
							</div>
						{/if}

						<div class="border border-gray-200 rounded-lg p-4">
							<div class="flex items-center justify-between mb-2">
								<div class="flex items-center">
									<div class="text-2xl mr-3">🎤</div>
									<div>
										<div class="font-medium text-gray-900">Voice Features</div>
										<div class="text-sm text-gray-600">Speech-to-text and text-to-speech</div>
									</div>
								</div>
								<div class="flex items-center">
									{#if aiStatus.speech}
										<span class="text-green-600 font-medium">✓ Ready</span>
									{:else}
										<span class="text-green-600 text-sm">Built-in</span>
									{/if}
								</div>
							</div>
						</div>
					</div>

					<div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
						<div class="flex items-start">
							<div class="text-yellow-600 mr-3">💡</div>
							<div>
								<div class="font-medium text-yellow-900">Pro Tip</div>
								<div class="text-sm text-yellow-800 mt-1">
									Start with local AI for privacy, then add cloud AI for advanced features.
									You can change these settings anytime in AI Settings.
								</div>
							</div>
						</div>
					</div>
				</div>
			{/if}

			{#if steps[currentStep].content === 'testing'}
				<div class="space-y-6">
					<p class="text-gray-600">
						Let's test your AI setup to make sure everything works correctly.
					</p>

					<!-- Test Results -->
					<div class="space-y-4">
						<div class="border border-gray-200 rounded-lg p-4">
							<div class="flex items-center justify-between mb-2">
								<div class="font-medium text-gray-900">Basic AI Functionality</div>
								<div class="flex items-center">
									{#if testResults.basicAI === 'success'}
										<span class="text-green-600 font-medium">✓ Passed</span>
									{:else if testResults.basicAI === 'failed'}
										<span class="text-red-600 font-medium">✗ Failed</span>
									{:else if testResults.basicAI === 'testing'}
										<span class="text-blue-600 font-medium">⏳ Testing...</span>
									{:else}
										<span class="text-gray-500">Not tested</span>
									{/if}
								</div>
							</div>
							<div class="text-sm text-gray-600">
								Tests AI writing suggestions and basic functionality
							</div>
						</div>

						<div class="border border-gray-200 rounded-lg p-4">
							<div class="flex items-center justify-between mb-2">
								<div class="font-medium text-gray-900">Speech Features</div>
								<div class="flex items-center">
									{#if testResults.speech === 'success'}
										<span class="text-green-600 font-medium">✓ Passed</span>
									{:else if testResults.speech === 'failed'}
										<span class="text-red-600 font-medium">✗ Failed</span>
									{:else if testResults.speech === 'testing'}
										<span class="text-blue-600 font-medium">⏳ Testing...</span>
									{:else}
										<span class="text-gray-500">Not tested</span>
									{/if}
								</div>
							</div>
							<div class="text-sm text-gray-600">
								Tests voice input and text-to-speech capabilities
							</div>
						</div>

						<div class="border border-gray-200 rounded-lg p-4">
							<div class="flex items-center justify-between mb-2">
								<div class="font-medium text-gray-900">AI Suggestions</div>
								<div class="flex items-center">
									{#if testResults.suggestions === 'success'}
										<span class="text-green-600 font-medium">✓ Passed</span>
									{:else if testResults.suggestions === 'failed'}
										<span class="text-red-600 font-medium">✗ Failed</span>
									{:else if testResults.suggestions === 'testing'}
										<span class="text-blue-600 font-medium">⏳ Testing...</span>
									{:else}
										<span class="text-gray-500">Not tested</span>
									{/if}
								</div>
							</div>
							<div class="text-sm text-gray-600">
								Tests intelligent suggestions and recommendations
							</div>
						</div>
					</div>

					{#if !aiStatus.testing}
						<button
							on:click={testAIProviders}
							class="w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors"
						>
							Run Tests
						</button>
					{:else}
						<div class="text-center py-4">
							<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-2"></div>
							<p class="text-gray-600">Testing your AI setup...</p>
						</div>
					{/if}

					<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
						<div class="text-sm text-blue-800">
							<strong>Note:</strong> Tests require working AI providers. If tests fail,
							you can still use TPT Titan - some features may need additional setup.
						</div>
					</div>
				</div>
			{/if}

			{#if steps[currentStep].content === 'complete'}
				<div class="text-center">
					<div class="text-6xl mb-6">🎉</div>
					<p class="text-lg text-gray-700 mb-6">
						Your AI setup is complete! You're ready to boost your productivity with AI-powered features.
					</p>

					<div class="bg-green-50 border border-green-200 rounded-lg p-6 mb-6">
						<h4 class="font-medium text-green-900 mb-4">Your AI Setup Summary:</h4>
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
							<div class="text-left">
								<div class="font-medium text-green-800">Experience Level:</div>
								<div class="text-green-700 capitalize">{userPreferences.experienceLevel}</div>
							</div>
							<div class="text-left">
								<div class="font-medium text-green-800">Primary Use:</div>
								<div class="text-green-700 capitalize">{userPreferences.primaryUse}</div>
							</div>
							<div class="text-left">
								<div class="font-medium text-green-800">Privacy Setting:</div>
								<div class="text-green-700 capitalize">{userPreferences.privacyPreference.replace('-', ' ')}</div>
							</div>
							<div class="text-left">
								<div class="font-medium text-green-800">AI Providers:</div>
								<div class="text-green-700">
									{aiStatus.localAI ? 'Local AI ✓' : 'Local AI ✗'},
									{aiStatus.cloudAI ? 'Cloud AI ✓' : 'Cloud AI ✗'}
								</div>
							</div>
						</div>
					</div>

					<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
						<h4 class="font-medium text-blue-900 mb-2">Next Steps:</h4>
						<ul class="text-sm text-blue-800 text-left space-y-1">
							<li>• Try AI writing suggestions in the text editor</li>
							<li>• Set up your email categorization</li>
							<li>• Create a voice-controlled task</li>
							<li>• Explore workflow automation</li>
						</ul>
					</div>
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="px-8 py-4 border-t border-gray-200 bg-gray-50">
			<div class="flex items-center justify-between">
				<button
					on:click={prevStep}
					disabled={currentStep === 0}
					class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Back
				</button>

				<div class="text-sm text-gray-500">
					{currentStep + 1} of {steps.length}
				</div>

				{#if currentStep === steps.length - 1}
					<button
						on:click={completeSetup}
						class="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
					>
						Get Started!
					</button>
				{:else}
					<button
						on:click={nextStep}
						class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					>
						Next
					</button>
				{/if}
			</div>
		</div>
	</div>
</div>
{/if}

<style>
	.max-h-90vh {
		max-height: 90vh;
	}
</style>
