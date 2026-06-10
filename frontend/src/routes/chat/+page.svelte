<script>
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { apiGet, apiPost } from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let rooms = [];
	let selectedRoom = null;
	let messages = [];
	let newMessage = '';
	let ws = null;
	let isLoadingRooms = false;
	let isLoadingMessages = false;
	let loadError = null;
	let showNewRoomForm = false;
	let newRoomName = '';
	let newRoomType = 'group';
	let messagesEl = null;

	onMount(async () => {
		await loadRooms();
	});

	onDestroy(() => {
		if (ws) ws.close();
	});

	async function loadRooms() {
		isLoadingRooms = true;
		loadError = null;
		try {
			const data = await apiGet('/chat/rooms');
			rooms = data.rooms || [];
		} catch (err) {
			if (err?.status === 401) { goto('/auth/login'); return; }
			loadError = 'Could not load chat rooms. Please check your connection.';
		} finally {
			isLoadingRooms = false;
		}
	}

	async function selectRoom(room) {
		selectedRoom = room;
		isLoadingMessages = true;
		messages = [];
		try {
			const data = await apiGet(`/chat/rooms/${room.id}/messages`);
			messages = data.messages || [];
			scrollToBottom();
		} catch (err) {
			console.error('Failed to load messages:', err);
		} finally {
			isLoadingMessages = false;
		}
		connectWebSocket(room.id);
	}

	function connectWebSocket(roomId) {
		if (ws) ws.close();
		const token = localStorage.getItem('token');
		const protocol = location.protocol === 'https:' ? 'wss' : 'ws';
		ws = new WebSocket(`${protocol}://${location.host}/api/v1/ws?token=${token}`);

		ws.onmessage = (event) => {
			try {
				const msg = JSON.parse(event.data);
				if (msg.room_id === roomId) {
					messages = [...messages, msg];
					scrollToBottom();
				}
			} catch {}
		};

		ws.onclose = () => { ws = null; };
	}

	async function sendMessage() {
		if (!newMessage.trim() || !selectedRoom) return;
		const content = newMessage.trim();
		newMessage = '';
		try {
			await apiPost(`/chat/rooms/${selectedRoom.id}/messages`, { content });
		} catch (err) {
			console.error('Failed to send message:', err);
			newMessage = content;
		}
	}

	function handleKeydown(event) {
		if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			sendMessage();
		}
	}

	async function createRoom() {
		if (!newRoomName.trim()) return;
		try {
			await apiPost('/chat/rooms', { name: newRoomName.trim(), type: newRoomType });
			newRoomName = '';
			showNewRoomForm = false;
			await loadRooms();
		} catch (err) {
			console.error('Failed to create room:', err);
		}
	}

	function scrollToBottom() {
		setTimeout(() => {
			if (messagesEl) messagesEl.scrollTop = messagesEl.scrollHeight;
		}, 0);
	}

	function formatTime(ts) {
		if (!ts) return '';
		return new Date(ts).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
	}
</script>

<svelte:head>
	<title>Chat - TPT Titan</title>
</svelte:head>

<div class="h-screen flex flex-col">
	<div class="flex flex-1 overflow-hidden">
		<!-- Sidebar: room list -->
		<aside class="w-72 bg-gray-900 text-white flex flex-col flex-shrink-0">
			<div class="px-4 py-4 border-b border-gray-700 flex items-center justify-between">
				<h1 class="text-lg font-semibold">Chat</h1>
				<button
					on:click={() => showNewRoomForm = !showNewRoomForm}
					aria-label="New room"
					class="text-gray-400 hover:text-white transition-colors"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
					</svg>
				</button>
			</div>

			{#if showNewRoomForm}
				<div class="px-4 py-3 border-b border-gray-700 space-y-2">
					<input
						bind:value={newRoomName}
						placeholder="Room name"
						class="w-full px-3 py-1.5 text-sm bg-gray-700 rounded text-white placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-blue-500"
					/>
					<select bind:value={newRoomType} class="w-full px-3 py-1.5 text-sm bg-gray-700 rounded text-white focus:outline-none">
						<option value="group">Group</option>
						<option value="channel">Channel</option>
					</select>
					<div class="flex gap-2">
						<button on:click={createRoom} class="flex-1 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 rounded text-white transition-colors">Create</button>
						<button on:click={() => showNewRoomForm = false} class="flex-1 py-1.5 text-sm bg-gray-600 hover:bg-gray-500 rounded text-white transition-colors">Cancel</button>
					</div>
				</div>
			{/if}

			<div class="flex-1 overflow-y-auto py-2">
				{#if loadError}
					<div class="px-4 py-3 text-sm text-red-400">
						{loadError}
						<button on:click={loadRooms} class="block mt-1 underline hover:no-underline">Retry</button>
					</div>
				{:else if isLoadingRooms}
					<div class="px-4 py-4 flex items-center gap-2 text-gray-400 text-sm">
						<svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
						</svg>
						Loading…
					</div>
				{:else if rooms.length === 0}
					<p class="px-4 py-3 text-sm text-gray-400">No rooms yet. Create one above.</p>
				{:else}
					{#each rooms as room}
						<button
							on:click={() => selectRoom(room)}
							class="w-full text-left px-4 py-3 hover:bg-gray-700 transition-colors flex items-center gap-3 {selectedRoom?.id === room.id ? 'bg-gray-700' : ''}"
						>
							<span class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-sm font-medium flex-shrink-0">
								{room.name?.[0]?.toUpperCase() ?? '#'}
							</span>
							<div class="min-w-0">
								<p class="text-sm font-medium truncate">{room.name}</p>
								<p class="text-xs text-gray-400 capitalize">{room.type || 'group'}</p>
							</div>
						</button>
					{/each}
				{/if}
			</div>
		</aside>

		<!-- Main: message area -->
		<div class="flex-1 flex flex-col bg-white dark:bg-gray-900 overflow-hidden">
			{#if !selectedRoom}
				<div class="flex-1 flex items-center justify-center text-gray-400 dark:text-gray-500">
					<div class="text-center">
						<svg class="mx-auto w-12 h-12 mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"/>
						</svg>
						<p class="text-sm">Select a room to start chatting</p>
					</div>
				</div>
			{:else}
				<!-- Room header -->
				<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center gap-3">
					<span class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-sm font-medium text-white">
						{selectedRoom.name?.[0]?.toUpperCase() ?? '#'}
					</span>
					<div>
						<h2 class="font-semibold text-gray-900 dark:text-white">{selectedRoom.name}</h2>
						<p class="text-xs text-gray-500 dark:text-gray-400 capitalize">{selectedRoom.type || 'group'}</p>
					</div>
				</div>

				<!-- Messages -->
				<div class="flex-1 overflow-y-auto px-6 py-4 space-y-4" bind:this={messagesEl}>
					{#if isLoadingMessages}
						<div class="flex items-center justify-center py-8 text-gray-400">
							<svg class="animate-spin w-6 h-6 mr-2" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
							</svg>
							Loading messages…
						</div>
					{:else if messages.length === 0}
						<p class="text-center text-sm text-gray-400 dark:text-gray-500 py-8">No messages yet. Say hello!</p>
					{:else}
						{#each messages as msg}
							<div class="flex items-start gap-3">
								<span class="w-8 h-8 rounded-full bg-gray-300 dark:bg-gray-600 flex items-center justify-center text-xs font-medium flex-shrink-0 text-gray-700 dark:text-gray-300">
									{msg.sender_name?.[0]?.toUpperCase() ?? '?'}
								</span>
								<div>
									<div class="flex items-baseline gap-2">
										<span class="text-sm font-medium text-gray-900 dark:text-white">{msg.sender_name || 'Unknown'}</span>
										<span class="text-xs text-gray-400">{formatTime(msg.created_at)}</span>
									</div>
									<p class="text-sm text-gray-700 dark:text-gray-300 mt-0.5">{msg.content}</p>
								</div>
							</div>
						{/each}
					{/if}
				</div>

				<!-- Input -->
				<div class="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
					<div class="flex items-end gap-3">
						<textarea
							bind:value={newMessage}
							on:keydown={handleKeydown}
							placeholder="Type a message… (Enter to send, Shift+Enter for newline)"
							rows="1"
							class="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none text-sm"
						></textarea>
						<button
							on:click={sendMessage}
							disabled={!newMessage.trim()}
							aria-label="Send message"
							class="p-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-40 text-white rounded-lg transition-colors flex-shrink-0"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
							</svg>
						</button>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
