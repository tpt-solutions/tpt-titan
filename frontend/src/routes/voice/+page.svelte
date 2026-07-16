<script>
	import { onMount } from 'svelte';
	import {
		getVoiceNotes,
		getVoiceNote,
		createVoiceNote,
		deleteVoiceNote,
		getVoiceAnnotations,
		createVoiceAnnotation,
		deleteVoiceAnnotation,
		formatApiError
	} from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let tab = 'notes';
	let error = null;
	let loading = false;

	// Notes
	let notes = [];
	let selectedNote = null;
	let noteTitle = '';
	let noteTags = '';
	let recording = false;
	let recordedChunks = [];
	let mediaRecorder = null;
	let recordSeconds = 0;
	let recordTimer = null;
	let audioEl = null;

	// Annotations
	let annotations = [];
	let annContentType = 'document';
	let annContentId = '';
	let annTitle = '';
	let annRecording = false;
	let annChunks = [];
	let annSeconds = 0;
	let annTimer = null;

	function arrayBufferToBase64(buffer) {
		let binary = '';
		const bytes = new Uint8Array(buffer);
		for (let i = 0; i < bytes.byteLength; i++) binary += String.fromCharCode(bytes[i]);
		return btoa(binary);
	}

	async function loadNotes() {
		loading = true;
		error = null;
		try {
			const r = await getVoiceNotes({ limit: 100 });
			notes = r.notes || [];
		} catch (err) {
			error = formatApiError(err);
		} finally {
			loading = false;
		}
	}

	async function openNote(note) {
		error = null;
		try {
			const full = await getVoiceNote(note.id);
			selectedNote = full;
			if (audioEl && full.audio_data) {
				audioEl.src = `data:audio/${full.audio_format};base64,${full.audio_data}`;
			}
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function startRecording(target) {
		try {
			const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
			const rec = new MediaRecorder(stream);
			const chunks = target === 'note' ? (recordedChunks = []) : (annChunks = []);
			rec.ondataavailable = (e) => chunks.push(e.data);
			rec.onstop = () => {
				stream.getTracks().forEach(t => t.stop());
			};
			rec.start();
			if (target === 'note') {
				mediaRecorder = rec;
				recording = true;
				recordSeconds = 0;
				recordTimer = setInterval(() => recordSeconds++, 1000);
			} else {
				mediaRecorder = rec;
				annRecording = true;
				annSeconds = 0;
				annTimer = setInterval(() => annSeconds++, 1000);
			}
		} catch (err) {
			error = 'Microphone access denied: ' + err.message;
		}
	}

	function stopRecording(target) {
		if (mediaRecorder && mediaRecorder.state !== 'inactive') mediaRecorder.stop();
		if (target === 'note') {
			recording = false;
			clearInterval(recordTimer);
		} else {
			annRecording = false;
			clearInterval(annTimer);
		}
	}

	async function saveNote() {
		if (!noteTitle.trim()) { error = 'Title is required'; return; }
		if (recordedChunks.length === 0) { error = 'Record audio first'; return; }
		error = null;
		try {
			const blob = new Blob(recordedChunks);
			const buf = await blob.arrayBuffer();
			const res = await createVoiceNote({
				title: noteTitle,
				tags: noteTags.split(',').map(t => t.trim()).filter(Boolean),
				audio_data: arrayBufferToBase64(buf),
				audio_format: blob.type.split('/')[1] || 'webm',
				duration: recordSeconds
			});
			noteTitle = '';
			noteTags = '';
			recordedChunks = [];
			recordSeconds = 0;
			await loadNotes();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function removeNote(note) {
		if (!confirm('Delete this voice note?')) return;
		try {
			await deleteVoiceNote(note.id);
			notes = notes.filter(n => n.id !== note.id);
			if (selectedNote?.id === note.id) selectedNote = null;
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function loadAnnotations() {
		if (!annContentId.trim()) { error = 'Content ID is required to list annotations'; return; }
		error = null;
		try {
			const r = await getVoiceAnnotations({ content_type: annContentType, content_id: annContentId });
			annotations = r.annotations || [];
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function saveAnnotation() {
		if (!annTitle.trim() || !annContentId.trim()) { error = 'Title and content ID are required'; return; }
		if (annChunks.length === 0) { error = 'Record audio first'; return; }
		error = null;
		try {
			const blob = new Blob(annChunks);
			const buf = await blob.arrayBuffer();
			await createVoiceAnnotation({
				content_type: annContentType,
				content_id: annContentId,
				title: annTitle,
				audio_data: arrayBufferToBase64(buf),
				audio_format: blob.type.split('/')[1] || 'webm',
				duration: annSeconds
			});
			annTitle = '';
			annChunks = [];
			annSeconds = 0;
			await loadAnnotations();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function removeAnnotation(a) {
		if (!confirm('Delete this annotation?')) return;
		try {
			await deleteVoiceAnnotation(a.id);
			annotations = annotations.filter(x => x.id !== a.id);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	onMount(loadNotes);
</script>

<svelte:head>
	<title>Voice - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-5xl">
	<div class="mb-6">
		<h1 class="text-3xl font-bold text-gray-900">Voice</h1>
		<p class="text-sm text-gray-500 mt-1">Voice notes and annotations</p>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 border border-red-200 rounded-lg px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	<div class="flex flex-wrap gap-1 border-b border-gray-200 mb-6">
		{#each ['notes', 'annotations'] as t}
			<button
				class="px-4 py-2 text-sm font-medium border-b-2 transition-colors {tab === t ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-600 hover:text-gray-900'}"
				on:click={() => (tab = t)}
			>
				{t.charAt(0).toUpperCase() + t.slice(1)}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-16 text-gray-400"><span>Loading…</span></div>
	{:else if tab === 'notes'}
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div>
				<div class="bg-white border border-gray-200 rounded-lg p-5 space-y-3">
					<h2 class="font-semibold text-gray-900">Record a note</h2>
					<input bind:value={noteTitle} placeholder="Title" class="w-full px-3 py-2 border border-gray-300 rounded text-sm" />
					<input bind:value={noteTags} placeholder="Tags (comma separated)" class="w-full px-3 py-2 border border-gray-300 rounded text-sm" />
					<div class="flex items-center gap-3">
						{#if !recording}
							<button on:click={() => startRecording('note')} class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded text-sm">● Record</button>
						{:else}
							<button on:click={() => stopRecording('note')} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm">■ Stop ({recordSeconds}s)</button>
						{/if}
						{#if recordedChunks.length}
							<button on:click={saveNote} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Save</button>
						{/if}
					</div>
				</div>

				<h2 class="font-semibold text-gray-900 mt-6 mb-2">Your notes</h2>
				<div class="space-y-2">
					{#each notes as note}
						<div class="bg-white border border-gray-200 rounded-lg p-3 flex items-center gap-3">
							<button on:click={() => openNote(note)} class="flex-1 text-left">
								<p class="font-medium text-gray-900">{note.title}</p>
								<p class="text-xs text-gray-400">{note.duration ? note.duration + 's · ' : ''}{note.tags?.join(', ')}</p>
							</button>
							<button on:click={() => removeNote(note)} class="text-red-600 hover:underline text-xs">Delete</button>
						</div>
					{:else}
						<div class="text-center text-gray-400 py-8">No voice notes yet</div>
					{/each}
				</div>
			</div>

			<div>
				{#if selectedNote}
					<div class="bg-white border border-gray-200 rounded-lg p-5">
						<h2 class="font-semibold text-gray-900 mb-1">{selectedNote.title}</h2>
						<p class="text-xs text-gray-400 mb-3">{selectedNote.content}</p>
						<audio bind:this={audioEl} controls class="w-full"></audio>
					</div>
				{:else}
					<div class="bg-gray-50 border border-dashed border-gray-300 rounded-lg p-8 text-center text-gray-400">
						Select a note to play it back
					</div>
				{/if}
			</div>
		</div>
	{:else if tab === 'annotations'}
		<div class="bg-white border border-gray-200 rounded-lg p-5 space-y-3 mb-6">
			<h2 class="font-semibold text-gray-900">Record an annotation</h2>
			<div class="grid grid-cols-2 md:grid-cols-3 gap-3">
				<label class="text-sm">Content type
					<select bind:value={annContentType} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded">
						<option value="document">document</option>
						<option value="task">task</option>
						<option value="email">email</option>
						<option value="calendar">calendar</option>
						<option value="contact">contact</option>
					</select>
				</label>
				<label class="text-sm md:col-span-2">Content ID
					<input bind:value={annContentId} placeholder="UUID of the content" class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
				</label>
				<label class="text-sm md:col-span-3">Title
					<input bind:value={annTitle} placeholder="Annotation title" class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
				</label>
			</div>
			<div class="flex items-center gap-3">
				{#if !annRecording}
					<button on:click={() => startRecording('ann')} class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded text-sm">● Record</button>
				{:else}
					<button on:click={() => stopRecording('ann')} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm">■ Stop ({annSeconds}s)</button>
				{/if}
				{#if annChunks.length}
					<button on:click={saveAnnotation} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Save annotation</button>
				{/if}
			</div>
			<button on:click={loadAnnotations} class="text-xs text-blue-600 hover:underline">List annotations for this content ID</button>
		</div>

		<div class="space-y-2">
			{#each annotations as a}
				<div class="bg-white border border-gray-200 rounded-lg p-3 flex items-center gap-3">
					<div class="flex-1">
						<p class="font-medium text-gray-900">{a.title}</p>
						<p class="text-xs text-gray-400">{a.duration ? a.duration + 's' : ''}</p>
					</div>
					<button on:click={() => removeAnnotation(a)} class="text-red-600 hover:underline text-xs">Delete</button>
				</div>
			{:else}
				<div class="text-center text-gray-400 py-8">No annotations for this content</div>
			{/each}
		</div>
	{/if}
</div>
