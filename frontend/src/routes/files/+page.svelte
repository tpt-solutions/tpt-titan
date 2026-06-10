<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { apiGet, apiPost } from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let folders = [];
	let selectedFolder = null;
	let syncStatus = null;
	let isLoadingFolders = false;
	let isLoadingStatus = false;
	let loadError = null;
	let showAddFolderForm = false;
	let newFolderPath = '';
	let newFolderName = '';
	let isSyncing = false;

	onMount(async () => {
		await Promise.all([loadFolders(), loadSyncStatus()]);
	});

	async function loadFolders() {
		isLoadingFolders = true;
		loadError = null;
		try {
			const data = await apiGet('/filesync/folders');
			folders = data.folders || [];
		} catch (err) {
			if (err?.status === 401) { goto('/auth/login'); return; }
			loadError = 'Could not load sync folders. Please check your connection.';
		} finally {
			isLoadingFolders = false;
		}
	}

	async function loadSyncStatus() {
		isLoadingStatus = true;
		try {
			const data = await apiGet('/filesync/status');
			syncStatus = data.status || data;
		} catch (err) {
			console.error('Failed to load sync status:', err);
		} finally {
			isLoadingStatus = false;
		}
	}

	async function addFolder() {
		if (!newFolderPath.trim()) return;
		try {
			await apiPost('/filesync/folders', {
				name: newFolderName.trim() || newFolderPath.split(/[\\/]/).pop(),
				path: newFolderPath.trim(),
			});
			newFolderPath = '';
			newFolderName = '';
			showAddFolderForm = false;
			await loadFolders();
		} catch (err) {
			console.error('Failed to add folder:', err);
		}
	}

	async function syncFolder(folderId) {
		isSyncing = true;
		try {
			await apiPost(`/filesync/sync/${folderId}`, {});
			await loadSyncStatus();
		} catch (err) {
			console.error('Failed to sync folder:', err);
		} finally {
			isSyncing = false;
		}
	}

	function formatBytes(bytes) {
		if (!bytes) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
	}

	function formatDate(ts) {
		if (!ts) return '—';
		return new Date(ts).toLocaleString();
	}
</script>

<svelte:head>
	<title>Files - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-5xl">
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Files</h1>
			<p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Sync folders across devices</p>
		</div>
		<button
			on:click={() => showAddFolderForm = !showAddFolderForm}
			class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
			</svg>
			Add Folder
		</button>
	</div>

	<!-- Error banner -->
	{#if loadError}
		<div class="mb-6 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg px-4 py-3 flex items-center gap-3">
			<svg class="w-5 h-5 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
			</svg>
			<span class="text-sm text-red-700 dark:text-red-300">{loadError}</span>
			<button on:click={loadFolders} class="ml-auto text-sm text-red-600 dark:text-red-400 underline hover:no-underline">Retry</button>
		</div>
	{/if}

	<!-- Add Folder Form -->
	{#if showAddFolderForm}
		<div class="mb-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3">
			<h3 class="font-medium text-gray-900 dark:text-white">Add Sync Folder</h3>
			<div>
				<label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">Folder path</label>
				<input
					bind:value={newFolderPath}
					placeholder="/home/user/Documents"
					class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
				/>
			</div>
			<div>
				<label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">Display name (optional)</label>
				<input
					bind:value={newFolderName}
					placeholder="My Documents"
					class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
				/>
			</div>
			<div class="flex gap-2">
				<button on:click={addFolder} class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors">Add</button>
				<button on:click={() => showAddFolderForm = false} class="px-4 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 text-sm rounded-lg transition-colors">Cancel</button>
			</div>
		</div>
	{/if}

	<!-- Sync Status Summary -->
	{#if syncStatus && !isLoadingStatus}
		<div class="mb-6 grid grid-cols-3 gap-4">
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 text-center">
				<p class="text-2xl font-bold text-gray-900 dark:text-white">{syncStatus.total_folders ?? folders.length}</p>
				<p class="text-sm text-gray-500 dark:text-gray-400">Folders</p>
			</div>
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 text-center">
				<p class="text-2xl font-bold text-gray-900 dark:text-white">{syncStatus.total_files ?? '—'}</p>
				<p class="text-sm text-gray-500 dark:text-gray-400">Files synced</p>
			</div>
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 text-center">
				<p class="text-2xl font-bold text-gray-900 dark:text-white">{syncStatus.total_size ? formatBytes(syncStatus.total_size) : '—'}</p>
				<p class="text-sm text-gray-500 dark:text-gray-400">Total size</p>
			</div>
		</div>
	{/if}

	<!-- Folder List -->
	{#if isLoadingFolders}
		<div class="flex items-center justify-center py-16 text-gray-400 dark:text-gray-500">
			<svg class="animate-spin w-8 h-8 mr-3" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
			</svg>
			<span>Loading folders…</span>
		</div>
	{:else if folders.length === 0}
		<div class="text-center py-16">
			<svg class="mx-auto w-12 h-12 text-gray-300 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7a2 2 0 012-2h4l2 2h8a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V7z"/>
			</svg>
			<h3 class="text-gray-500 dark:text-gray-400 font-medium mb-1">No sync folders</h3>
			<p class="text-sm text-gray-400 dark:text-gray-500">Add a folder to start syncing files across your devices.</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each folders as folder}
				<div
					class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 flex items-center gap-4 hover:border-blue-300 dark:hover:border-blue-600 transition-colors cursor-pointer {selectedFolder?.id === folder.id ? 'border-blue-400 dark:border-blue-500' : ''}"
					on:click={() => selectedFolder = selectedFolder?.id === folder.id ? null : folder}
					on:keydown={(e) => e.key === 'Enter' && (selectedFolder = selectedFolder?.id === folder.id ? null : folder)}
					role="button"
					tabindex="0"
				>
					<div class="w-10 h-10 rounded-lg bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center flex-shrink-0">
						<svg class="w-5 h-5 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7a2 2 0 012-2h4l2 2h8a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V7z"/>
						</svg>
					</div>
					<div class="flex-1 min-w-0">
						<p class="font-medium text-gray-900 dark:text-white truncate">{folder.name || folder.path}</p>
						<p class="text-sm text-gray-500 dark:text-gray-400 truncate">{folder.path}</p>
					</div>
					<div class="text-right flex-shrink-0 text-sm text-gray-500 dark:text-gray-400">
						<p>Last sync: {formatDate(folder.last_sync)}</p>
						{#if folder.size}
							<p>{formatBytes(folder.size)}</p>
						{/if}
					</div>
					<button
						on:click|stopPropagation={() => syncFolder(folder.id)}
						disabled={isSyncing}
						aria-label="Sync folder"
						class="p-2 text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 disabled:opacity-40 transition-colors"
					>
						<svg class="w-5 h-5 {isSyncing ? 'animate-spin' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
						</svg>
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>
