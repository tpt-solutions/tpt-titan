<script>
	import { createEventDispatcher } from 'svelte';
	import { getAISettings, updateAISettings, getSpeechSettings, updateSpeechSettings } from '../api.js';

	export let userId = null;

	const dispatch = createEventDispatcher();

	let aiSettings = {
		enableAIFeatures: true,
		enableOCR: true,
		enableSpeech: true,
		enableWorkflows: true,
		enableLocalAI: true,
		enableCloudAI: false,
		defaultProvider: 'ollama',
		maxConcurrentRequests: 3,
		requestTimeout: 30,
		hardwareAcceleration: true,
		lowPowerMode: false,
		apiKeys: {
			openai: '',
			elevenlabs: '',
			replicate: '',
			assemblyai: '',
			deepgram: ''
		}
	};

	let speechSettings = {
		enableTTS: true,
		enableSTT: true,
		defaultTTSModel: '',
		defaultSTTModel: '',
		defaultVoice: 'alloy',
		defaultLanguage: 'en',
		ttsSpeed: 1.0,
		ttsVolume: 1.0,
		sttLanguage: 'en',
		autoPlayTTS: false,
		showSTTTranscript: true,
		keyboardShortcut: 'ctrl+shift+s'
	};

	let isLoading = true;
	let isSaving = false;
	let activeTab = 'general';

	// Load settings on mount
	async function loadSettings() {
		try {
			isLoading = true;

			// Load AI settings
			const aiResponse = await getAISettings();
			if (aiResponse) {
				aiSettings = { ...aiSettings, ...aiResponse };
			}

			// Load speech settings
			const speechResponse = await getSpeechSettings();
			if (speechResponse) {
				speechSettings = { ...speechSettings, ...speechResponse };
			}

		} catch (error) {
			console.error('Failed to load AI settings:', error);
		} finally {
			isLoading = false;
		}
	}

	// Save AI settings
	async function saveAISettings() {
		try {
			isSaving = true;
			await updateAISettings(aiSettings);
			dispatch('saved', { type: 'ai', settings: aiSettings });
		} catch (error) {
			console.error('Failed to save AI settings:', error);
			alert('Failed to save AI settings: ' + error.message);
		} finally {
			isSaving = false;
		}
	}

	// Save speech settings
	async function saveSpeechSettings() {
		try {
			isSaving = true;
			await updateSpeechSettings(speechSettings);
			dispatch('saved', { type: 'speech', settings: speechSettings });
		} catch (error) {
			console.error('Failed to save speech settings:', error);
			alert('Failed to save speech settings: ' + error.message);
		} finally {
			isSaving = false;
		}
	}

	// Test API key
	async function testAPIKey(provider) {
		const key = aiSettings.apiKeys[provider];
		if (!key) {
			alert('Please enter an API key first');
			return;
		}

		try {
			// This would call a test endpoint for each provider
			alert(`Testing ${provider} API key... (This would normally make a test API call)`);
		} catch (error) {
			alert(`API key test failed: ${error.message}`);
		}
	}

	// Initialize on mount
	import { onMount } from 'svelte';
	onMount(() => {
		loadSettings();
	});
</script>

<div class="max-w-4xl mx-auto p-6">
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-900 mb-2">AI Feature Settings</h1>
		<p class="text-gray-600">Configure AI features, providers, and preferences for your TPT Titan experience.</p>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<span class="ml-3 text-gray-600">Loading settings...</span>
		</div>
	{:else}
		<!-- Tab Navigation -->
		<div class="border-b border-gray-200 mb-6">
			<nav class="flex space-x-8">
				<button
					class="py-2 px-1 border-b-2 font-medium text-sm {activeTab === 'general' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
					on:click={() => activeTab = 'general'}
				>
					General
				</button>
				<button
					class="py-2 px-1 border-b-2 font-medium text-sm {activeTab === 'providers' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
					on:click={() => activeTab = 'providers'}
				>
					AI Providers
				</button>
				<button
					class="py-2 px-1 border-b-2 font-medium text-sm {activeTab === 'speech' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
					on:click={() => activeTab = 'speech'}
				>
					Speech
				</button>
				<button
					class="py-2 px-1 border-b-2 font-medium text-sm {activeTab === 'advanced' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
					on:click={() => activeTab = 'advanced'}
				>
					Advanced
				</button>
			</nav>
		</div>

		<!-- Tab Content -->
		<div class="space-y-6">
			<!-- General Settings -->
			{#if activeTab === 'general'}
				<div class="bg-white shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200">
						<h3 class="text-lg font-medium text-gray-900">General AI Settings</h3>
						<p class="text-sm text-gray-600">Control which AI features are enabled and how they behave.</p>
					</div>

					<div class="px-6 py-4 space-y-6">
						<!-- Master Toggle -->
						<div class="flex items-center justify-between">
							<div>
								<h4 class="text-sm font-medium text-gray-900">Enable AI Features</h4>
								<p class="text-sm text-gray-600">Master switch for all AI functionality</p>
							</div>
							<label class="relative inline-flex items-center cursor-pointer">
								<input
									type="checkbox"
									bind:checked={aiSettings.enableAIFeatures}
									class="sr-only peer"
								>
								<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
							</label>
						</div>

						{#if aiSettings.enableAIFeatures}
							<!-- Feature Toggles -->
							<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
								<div class="space-y-4">
									<h4 class="text-sm font-medium text-gray-900">Document Processing</h4>
									<div class="flex items-center justify-between">
										<span class="text-sm text-gray-600">OCR & Document Analysis</span>
										<input
											type="checkbox"
											bind:checked={aiSettings.enableOCR}
											class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
										>
									</div>
									<div class="flex items-center justify-between">
										<span class="text-sm text-gray-600">Workflow Automation</span>
										<input
											type="checkbox"
											bind:checked={aiSettings.enableWorkflows}
											class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
										>
									</div>
								</div>

								<div class="space-y-4">
									<h4 class="text-sm font-medium text-gray-900">Speech Services</h4>
									<div class="flex items-center justify-between">
										<span class="text-sm text-gray-600">Text-to-Speech & Speech-to-Text</span>
										<input
											type="checkbox"
											bind:checked={aiSettings.enableSpeech}
											class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
										>
									</div>
								</div>
							</div>

							<!-- Processing Mode -->
							<div class="space-y-4">
								<h4 class="text-sm font-medium text-gray-900">AI Processing Mode</h4>
								<div class="space-y-2">
									<div class="flex items-center">
										<input
											id="local-ai"
											name="ai-mode"
											type="radio"
											bind:group={aiSettings.defaultProvider}
											value="local"
											class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
										>
										<label for="local-ai" class="ml-3 text-sm text-gray-700">
											<span class="font-medium">Local AI Only</span>
											<span class="block text-xs text-gray-500">Use only locally installed models (privacy-focused, no internet required)</span>
										</label>
									</div>
									<div class="flex items-center">
										<input
											id="hybrid-ai"
											name="ai-mode"
											type="radio"
											bind:group={aiSettings.defaultProvider}
											value="hybrid"
											class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
										>
										<label for="hybrid-ai" class="ml-3 text-sm text-gray-700">
											<span class="font-medium">Hybrid Mode</span>
											<span class="block text-xs text-gray-500">Prefer local models, fall back to cloud services when needed</span>
										</label>
									</div>
									<div class="flex items-center">
										<input
											id="cloud-ai"
											name="ai-mode"
											type="radio"
											bind:group={aiSettings.defaultProvider}
											value="cloud"
											class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
										>
										<label for="cloud-ai" class="ml-3 text-sm text-gray-700">
											<span class="font-medium">Cloud Services</span>
											<span class="block text-xs text-gray-500">Use cloud AI services for best performance and features</span>
										</label>
									</div>
								</div>
							</div>
						{/if}
					</div>

					<div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
						<button
							on:click={saveAISettings}
							disabled={isSaving}
							class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition-colors"
						>
							{isSaving ? 'Saving...' : 'Save General Settings'}
						</button>
					</div>
				</div>
			{/if}

			<!-- AI Providers Settings -->
			{#if activeTab === 'providers'}
				<div class="bg-white shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200">
						<h3 class="text-lg font-medium text-gray-900">AI Provider Configuration</h3>
						<p class="text-sm text-gray-600">Configure API keys and preferences for different AI service providers.</p>
					</div>

					<div class="px-6 py-4 space-y-6">
						<!-- Local AI -->
						<div class="border border-gray-200 rounded-lg p-4">
							<div class="flex items-center justify-between mb-4">
								<div>
									<h4 class="text-sm font-medium text-gray-900">Local AI (Ollama)</h4>
									<p class="text-sm text-gray-600">Run AI models locally on your machine</p>
								</div>
								<div class="flex items-center space-x-2">
									<span class="text-xs text-green-600 bg-green-100 px-2 py-1 rounded">Recommended</span>
									<input
										type="checkbox"
										bind:checked={aiSettings.enableLocalAI}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
								</div>
							</div>
							<div class="text-sm text-gray-600">
								<p><strong>Benefits:</strong> Private, no API costs, works offline</p>
								<p><strong>Requirements:</strong> Sufficient RAM (8GB+ recommended), Ollama installed</p>
							</div>
						</div>

						<!-- Cloud Providers -->
						<div class="space-y-4">
							<h4 class="text-sm font-medium text-gray-900">Cloud AI Providers</h4>

							<!-- OpenAI -->
							<div class="border border-gray-200 rounded-lg p-4">
								<div class="flex items-center justify-between mb-4">
									<div>
										<h4 class="text-sm font-medium text-gray-900">OpenAI</h4>
										<p class="text-sm text-gray-600">GPT models for text generation, Whisper for speech</p>
									</div>
									<input
										type="checkbox"
										bind:checked={aiSettings.enableCloudAI}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
								</div>
								<div class="space-y-3">
									<div>
										<label class="block text-xs font-medium text-gray-700 mb-1">API Key</label>
										<div class="flex gap-2">
											<input
												type="password"
												bind:value={aiSettings.apiKeys.openai}
												placeholder="sk-..."
												class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
											>
											<button
												on:click={() => testAPIKey('openai')}
												class="px-3 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 text-sm"
											>
												Test
											</button>
										</div>
									</div>
								</div>
							</div>

							<!-- ElevenLabs -->
							<div class="border border-gray-200 rounded-lg p-4">
								<div class="flex items-center justify-between mb-4">
									<div>
										<h4 class="text-sm font-medium text-gray-900">ElevenLabs</h4>
										<p class="text-sm text-gray-600">High-quality text-to-speech voices</p>
									</div>
									<input
										type="checkbox"
										checked={true}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
										disabled
									>
								</div>
								<div class="space-y-3">
									<div>
										<label class="block text-xs font-medium text-gray-700 mb-1">API Key</label>
										<div class="flex gap-2">
											<input
												type="password"
												bind:value={aiSettings.apiKeys.elevenlabs}
												placeholder="Enter ElevenLabs API key"
												class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
											>
											<button
												on:click={() => testAPIKey('elevenlabs')}
												class="px-3 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 text-sm"
											>
												Test
											</button>
										</div>
									</div>
								</div>
							</div>

							<!-- Replicate -->
							<div class="border border-gray-200 rounded-lg p-4">
								<div class="flex items-center justify-between mb-4">
									<div>
										<h4 class="text-sm font-medium text-gray-900">Replicate</h4>
										<p class="text-sm text-gray-600">Open-source AI models via API</p>
									</div>
									<input
										type="checkbox"
										checked={true}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
										disabled
									>
								</div>
								<div class="space-y-3">
									<div>
										<label class="block text-xs font-medium text-gray-700 mb-1">API Key</label>
										<div class="flex gap-2">
											<input
												type="password"
												bind:value={aiSettings.apiKeys.replicate}
												placeholder="Enter Replicate API key"
												class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
											>
											<button
												on:click={() => testAPIKey('replicate')}
												class="px-3 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 text-sm"
											>
												Test
											</button>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>

					<div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
						<button
							on:click={saveAISettings}
							disabled={isSaving}
							class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition-colors"
						>
							{isSaving ? 'Saving...' : 'Save Provider Settings'}
						</button>
					</div>
				</div>
			{/if}

			<!-- Speech Settings -->
			{#if activeTab === 'speech'}
				<div class="bg-white shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200">
						<h3 class="text-lg font-medium text-gray-900">Speech Settings</h3>
						<p class="text-sm text-gray-600">Configure text-to-speech and speech-to-text preferences.</p>
					</div>

					<div class="px-6 py-4 space-y-6">
						<!-- Speech Feature Toggles -->
						<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
							<div class="space-y-4">
								<h4 class="text-sm font-medium text-gray-900">Speech Features</h4>
								<div class="flex items-center justify-between">
									<span class="text-sm text-gray-600">Enable Text-to-Speech</span>
									<input
										type="checkbox"
										bind:checked={speechSettings.enableTTS}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
								</div>
								<div class="flex items-center justify-between">
									<span class="text-sm text-gray-600">Enable Speech-to-Text</span>
									<input
										type="checkbox"
										bind:checked={speechSettings.enableSTT}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
								</div>
							</div>

							<div class="space-y-4">
								<h4 class="text-sm font-medium text-gray-900">TTS Preferences</h4>
								<div>
									<label class="block text-xs font-medium text-gray-700 mb-1">Default Voice</label>
									<select
										bind:value={speechSettings.defaultVoice}
										class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
									>
										<option value="alloy">Alloy (OpenAI)</option>
										<option value="echo">Echo (OpenAI)</option>
										<option value="fable">Fable (OpenAI)</option>
										<option value="onyx">Onyx (OpenAI)</option>
										<option value="nova">Nova (OpenAI)</option>
										<option value="shimmer">Shimmer (OpenAI)</option>
									</select>
								</div>
								<div>
									<label class="block text-xs font-medium text-gray-700 mb-1">Language</label>
									<select
										bind:value={speechSettings.defaultLanguage}
										class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
									>
										<option value="en">English</option>
										<option value="es">Spanish</option>
										<option value="fr">French</option>
										<option value="de">German</option>
										<option value="it">Italian</option>
										<option value="pt">Portuguese</option>
									</select>
								</div>
							</div>
						</div>

						<!-- TTS Settings -->
						<div class="space-y-4">
							<h4 class="text-sm font-medium text-gray-900">TTS Settings</h4>
							<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
								<div>
									<label class="block text-xs font-medium text-gray-700 mb-1">Speed ({speechSettings.ttsSpeed}x)</label>
									<input
										type="range"
										min="0.5"
										max="2.0"
										step="0.1"
										bind:value={speechSettings.ttsSpeed}
										class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
									>
								</div>
								<div>
									<label class="block text-xs font-medium text-gray-700 mb-1">Volume ({Math.round(speechSettings.ttsVolume * 100)}%)</label>
									<input
										type="range"
										min="0.1"
										max="1.0"
										step="0.1"
										bind:value={speechSettings.ttsVolume}
										class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
									>
								</div>
								<div class="flex items-center">
									<input
										type="checkbox"
										bind:checked={speechSettings.autoPlayTTS}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
									<label class="ml-2 text-xs font-medium text-gray-700">Auto-play TTS</label>
								</div>
							</div>
						</div>

						<!-- STT Settings -->
						<div class="space-y-4">
							<h4 class="text-sm font-medium text-gray-900">STT Settings</h4>
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
								<div>
									<label class="block text-xs font-medium text-gray-700 mb-1">Language</label>
									<select
										bind:value={speechSettings.sttLanguage}
										class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
									>
										<option value="en">English</option>
										<option value="es">Spanish</option>
										<option value="fr">French</option>
										<option value="de">German</option>
										<option value="it">Italian</option>
										<option value="pt">Portuguese</option>
									</select>
								</div>
								<div class="flex items-center">
									<input
										type="checkbox"
										bind:checked={speechSettings.showSTTTranscript}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
									<label class="ml-2 text-xs font-medium text-gray-700">Show transcripts</label>
								</div>
							</div>
						</div>

						<!-- Voice Commands -->
						<div class="space-y-4">
							<h4 class="text-sm font-medium text-gray-900">Voice Commands</h4>
							<div>
								<label class="block text-xs font-medium text-gray-700 mb-1">Keyboard Shortcut</label>
								<input
									type="text"
									bind:value={speechSettings.keyboardShortcut}
									placeholder="ctrl+shift+s"
									class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
								>
								<p class="text-xs text-gray-500 mt-1">Shortcut to activate voice commands</p>
							</div>
						</div>
					</div>

					<div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
						<button
							on:click={saveSpeechSettings}
							disabled={isSaving}
							class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition-colors"
						>
							{isSaving ? 'Saving...' : 'Save Speech Settings'}
						</button>
					</div>
				</div>
			{/if}

			<!-- Advanced Settings -->
			{#if activeTab === 'advanced'}
				<div class="bg-white shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200">
						<h3 class="text-lg font-medium text-gray-900">Advanced Settings</h3>
						<p class="text-sm text-gray-600">Performance tuning and resource management options.</p>
					</div>

					<div class="px-6 py-4 space-y-6">
						<!-- Performance Settings -->
						<div class="space-y-4">
							<h4 class="text-sm font-medium text-gray-900">Performance</h4>
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
								<div>
									<label class="block text-xs font-medium text-gray-700 mb-1">Max Concurrent Requests</label>
									<select
										bind:value={aiSettings.maxConcurrentRequests}
										class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
									>
										<option value={1}>1 (Conservative)</option>
										<option value={3}>3 (Balanced)</option>
										<option value={5}>5 (Performance)</option>
										<option value={10}>10 (High Performance)</option>
									</select>
								</div>
								<div>
									<label class="block text-xs font-medium text-gray-700 mb-1">Request Timeout (seconds)</label>
									<input
										type="number"
										bind:value={aiSettings.requestTimeout}
										min="10"
										max="300"
										class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-blue-500 focus:border-blue-500"
									>
								</div>
							</div>
						</div>

						<!-- Hardware Settings -->
						<div class="space-y-4">
							<h4 class="text-sm font-medium text-gray-900">Hardware</h4>
							<div class="space-y-3">
								<div class="flex items-center justify-between">
									<div>
										<span class="text-sm text-gray-600">Hardware Acceleration</span>
										<p class="text-xs text-gray-500">Use GPU/TPU when available</p>
									</div>
									<input
										type="checkbox"
										bind:checked={aiSettings.hardwareAcceleration}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
								</div>
								<div class="flex items-center justify-between">
									<div>
										<span class="text-sm text-gray-600">Low Power Mode</span>
										<p class="text-xs text-gray-500">Reduce resource usage (slower processing)</p>
									</div>
									<input
										type="checkbox"
										bind:checked={aiSettings.lowPowerMode}
										class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
									>
								</div>
							</div>
						</div>

						<!-- System Information -->
						<div class="space-y-4">
							<h4 class="text-sm font-medium text-gray-900">System Information</h4>
							<div class="bg-gray-50 rounded-lg p-4">
								<div class="grid grid-cols-2 gap-4 text-sm">
									<div>
										<span class="font-medium text-gray-700">Browser:</span>
										<span class="text-gray-600 ml-2">{typeof window !== 'undefined' ? navigator.userAgent.split(' ').pop() : 'Unknown'}</span>
									</div>
									<div>
										<span class="font-medium text-gray-700">Platform:</span>
										<span class="text-gray-600 ml-2">{typeof window !== 'undefined' ? navigator.platform : 'Unknown'}</span>
									</div>
									<div>
										<span class="font-medium text-gray-700">WebGL:</span>
										<span class="text-gray-600 ml-2">Supported</span>
									</div>
									<div>
										<span class="font-medium text-gray-700">Local Storage:</span>
										<span class="text-gray-600 ml-2">Available</span>
									</div>
								</div>
							</div>
						</div>
					</div>

					<div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
						<button
							on:click={saveAISettings}
							disabled={isSaving}
							class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition-colors"
						>
							{isSaving ? 'Saving...' : 'Save Advanced Settings'}
						</button>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	/* Custom range slider styling */
	input[type="range"] {
		-webkit-appearance: none;
		appearance: none;
		background: transparent;
		cursor: pointer;
	}

	input[type="range"]::-webkit-slider-thumb {
		-webkit-appearance: none;
		appearance: none;
		height: 20px;
		width: 20px;
		border-radius: 50%;
		background: #2563eb;
		cursor: pointer;
	}

	input[type="range"]::-moz-range-thumb {
		height: 20px;
		width: 20px;
		border-radius: 50%;
		background: #2563eb;
		cursor: pointer;
		border: none;
	}
</style>
