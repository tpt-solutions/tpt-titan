<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import {
		getSpeechModels,
		createSpeechModel,
		textToSpeech,
		speechToText,
		getSpeechRequestStatus,
		getSpeechSettings,
		updateSpeechSettings,
		getSpeechHistory,
		formatApiError
	} from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let tab = 'tts';
	let error = null;
	let loading = false;

	// Models
	let ttsModels = [];
	let sttModels = [];
	let selectedTtsModel = '';
	let selectedSttModel = '';

	// TTS
	let ttsText = 'Hello, this is a test of the text to speech system.';
	let ttsVoice = 'alloy';
	let ttsLanguage = 'en';
	let ttsFormat = 'mp3';
	let ttsSpeed = 1.0;
	let ttsAudioEl = null;
	let ttsProcessing = false;

	// STT
	let sttFile = null;
	let sttResult = '';
	let sttProcessing = false;

	// Settings
	let settings = null;
	let settingsSaved = false;

	// History
	let history = [];

	onMount(loadModels);

	async function loadModels() {
		loading = true;
		error = null;
		try {
			const [t, s] = await Promise.all([getSpeechModels('tts'), getSpeechModels('stt')]);
			ttsModels = t.models || [];
			sttModels = s.models || [];
			selectedTtsModel = ttsModels[0]?.id || '';
			selectedSttModel = sttModels[0]?.id || '';
		} catch (err) {
			error = formatApiError(err);
		} finally {
			loading = false;
		}
	}

	async function doTTS() {
		if (!selectedTtsModel || !ttsText.trim()) return;
		ttsProcessing = true;
		error = null;
		try {
			const res = await textToSpeech({
				text: ttsText,
				model_id: selectedTtsModel,
				voice: ttsVoice,
				language: ttsLanguage,
				format: ttsFormat,
				speed: ttsSpeed
			});
			const status = await pollStatus(res.request_id);
			if (status?.audio_data) {
				const audioSrc = `data:audio/${ttsFormat};base64,${status.audio_data}`;
				if (ttsAudioEl) ttsAudioEl.src = audioSrc;
			} else {
				error = 'Audio was not produced. ' + (status?.error || '');
			}
		} catch (err) {
			error = formatApiError(err);
		} finally {
			ttsProcessing = false;
		}
	}

	async function doSTT() {
		if (!selectedSttModel || !sttFile) return;
		sttProcessing = true;
		error = null;
		sttResult = '';
		try {
			const fd = new FormData();
			fd.append('audio', sttFile);
			fd.append('model_id', selectedSttModel);
			fd.append('language', ttsLanguage);
			fd.append('format', sttFile.name.split('.').pop() || 'wav');
			const res = await speechToText(fd);
			const status = await pollStatus(res.request_id);
			sttResult = status?.transcript || '';
			if (!sttResult && status?.error) error = status.error;
		} catch (err) {
			error = formatApiError(err);
		} finally {
			sttProcessing = false;
		}
	}

	async function pollStatus(requestId, attempts = 20) {
		for (let i = 0; i < attempts; i++) {
			await new Promise(r => setTimeout(r, 1000));
			const status = await getSpeechRequestStatus(requestId);
			if (status.status === 'completed' || status.status === 'failed') return status;
		}
		return await getSpeechRequestStatus(requestId);
	}

	async function loadSettings() {
		try {
			const r = await getSpeechSettings();
			settings = r.settings || {};
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function saveSettings() {
		try {
			settingsSaved = false;
			await updateSpeechSettings(settings);
			settingsSaved = true;
			setTimeout(() => (settingsSaved = false), 3000);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function loadHistory() {
		try {
			const r = await getSpeechHistory({ limit: 30 });
			history = r.requests || [];
		} catch (err) {
			error = formatApiError(err);
		}
	}

	function changeTab(t) {
		tab = t;
		if (t === 'settings') loadSettings();
		if (t === 'history') loadHistory();
	}

	function onFile(e) {
		sttFile = e.target.files[0];
	}
</script>

<svelte:head>
	<title>Speech - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<div class="mb-6">
		<h1 class="text-3xl font-bold text-gray-900">Speech</h1>
		<p class="text-sm text-gray-500 mt-1">Text-to-speech and speech-to-text</p>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 border border-red-200 rounded-lg px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	<div class="flex flex-wrap gap-1 border-b border-gray-200 mb-6">
		{#each ['tts', 'stt', 'settings', 'history'] as t}
			<button
				class="px-4 py-2 text-sm font-medium border-b-2 transition-colors {tab === t ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-600 hover:text-gray-900'}"
				on:click={() => changeTab(t)}
			>
				{t.toUpperCase()}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-16 text-gray-400"><span>Loading models…</span></div>
	{:else if tab === 'tts'}
		<div class="bg-white border border-gray-200 rounded-lg p-6 space-y-4">
			<label class="block text-sm">
				Model
				<select bind:value={selectedTtsModel} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded">
					{#each ttsModels as m}<option value={m.id}>{m.name} ({m.provider})</option>{:else}<option value="">No TTS models available</option>{/each}
				</select>
			</label>
			<label class="block text-sm">
				Text
				<textarea bind:value={ttsText} rows="4" class="mt-1 w-full px-3 py-2 border border-gray-300 rounded"></textarea>
			</label>
			<div class="grid grid-cols-2 md:grid-cols-4 gap-3">
				<label class="text-sm">Voice<input bind:value={ttsVoice} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" /></label>
				<label class="text-sm">Language<input bind:value={ttsLanguage} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" /></label>
				<label class="text-sm">Format
					<select bind:value={ttsFormat} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded">
						<option value="mp3">mp3</option><option value="wav">wav</option><option value="ogg">ogg</option>
					</select>
				</label>
				<label class="text-sm">Speed
					<input type="number" step="0.1" min="0.5" max="2" bind:value={ttsSpeed} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
				</label>
			</div>
			<button on:click={doTTS} disabled={ttsProcessing || !selectedTtsModel} class="bg-blue-600 hover:bg-blue-700 disabled:opacity-40 text-white px-4 py-2 rounded text-sm">
				{ttsProcessing ? 'Generating…' : 'Generate speech'}
			</button>
			<audio bind:this={ttsAudioEl} controls class="w-full"></audio>
		</div>
	{:else if tab === 'stt'}
		<div class="bg-white border border-gray-200 rounded-lg p-6 space-y-4">
			<label class="block text-sm">
				Model
				<select bind:value={selectedSttModel} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded">
					{#each sttModels as m}<option value={m.id}>{m.name} ({m.provider})</option>{:else}<option value="">No STT models available</option>{/each}
				</select>
			</label>
			<label class="block text-sm">
				Audio file
				<input type="file" accept="audio/*" on:change={onFile} class="mt-1 w-full text-sm" />
			</label>
			<button on:click={doSTT} disabled={sttProcessing || !selectedSttModel || !sttFile} class="bg-blue-600 hover:bg-blue-700 disabled:opacity-40 text-white px-4 py-2 rounded text-sm">
				{sttProcessing ? 'Transcribing…' : 'Transcribe'}
			</button>
			{#if sttResult}
				<div class="mt-4">
					<p class="text-sm font-medium text-gray-700 mb-1">Transcript</p>
					<div class="bg-gray-50 border border-gray-200 rounded p-3 text-sm whitespace-pre-wrap">{sttResult}</div>
				</div>
			{/if}
		</div>
	{:else if tab === 'settings'}
		{#if settings}
			<div class="bg-white border border-gray-200 rounded-lg p-6 max-w-xl">
				{#if settingsSaved}<div class="bg-green-50 border border-green-200 rounded px-3 py-2 text-sm text-green-700 mb-3">Settings saved</div>{/if}
				{#each Object.keys(settings) as key}
					<label class="flex items-center gap-2 text-sm mb-2">
						<input type="checkbox" bind:checked={settings[key]} class="rounded" />
						{key.replace(/_/g, ' ')}
					</label>
				{/each}
				<button on:click={saveSettings} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm mt-2">Save settings</button>
			</div>
		{/if}
	{:else if tab === 'history'}
		<div class="bg-white border border-gray-200 rounded-lg divide-y divide-gray-100 max-h-[70vh] overflow-auto">
			{#each history as req}
				<div class="px-4 py-2 text-sm">
					<div class="flex items-center gap-2">
						<span class="inline-flex px-2 py-0.5 rounded-full text-xs bg-gray-100 text-gray-600">{req.request_type}</span>
						<span class="inline-flex px-2 py-0.5 rounded-full text-xs {req.status === 'completed' ? 'bg-green-100 text-green-700' : req.status === 'failed' ? 'bg-red-100 text-red-700' : 'bg-yellow-100 text-yellow-700'}">{req.status}</span>
					</div>
					{#if req.output_text}<p class="text-gray-500 text-xs mt-1">{req.output_text}</p>{/if}
				</div>
			{:else}
				<div class="px-4 py-8 text-center text-gray-400">No history yet</div>
			{/each}
		</div>
	{/if}
</div>
