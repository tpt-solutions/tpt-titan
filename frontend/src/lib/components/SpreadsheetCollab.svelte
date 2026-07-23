<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import { apiGet, apiPost } from '$lib/api.js';

	export let spreadsheetId = null;

	let mode = 'local';
	let peers = [];
	let discovered = [];
	let status = null;
	let peerAddress = '';
	let loading = false;
	let error = '';
	let message = '';

	const sid = () => spreadsheetId || '00000000-0000-0000-0000-000000000000';

	async function loadAll() {
		if (!spreadsheetId) return;
		await Promise.all([loadMode(), loadPeers(), loadDiscovered(), loadStatus()]);
	}

	async function loadMode() {
		try { const res = await apiGet(`/spreadsheets/${sid()}/collab/mode`); mode = res.mode || 'local'; }
		catch (e) { error = e.message; }
	}
	async function loadPeers() {
		try { const res = await apiGet(`/spreadsheets/${sid()}/collab/peers`); peers = res.peers || []; }
		catch (e) { error = e.message; }
	}
	async function loadDiscovered() {
		try { const res = await apiGet(`/spreadsheets/${sid()}/collab/peers/discovered`); discovered = res.peers || []; }
		catch (e) { error = e.message; }
	}
	async function loadStatus() {
		try { const res = await apiGet(`/spreadsheets/${sid()}/collab/status`); status = res; }
		catch (e) { error = e.message; }
	}

	async function setMode(m) {
		loading = true; error = '';
		try { await apiPost(`/spreadsheets/${sid()}/collab/mode`, { mode: m }); mode = m; message = 'Mode set to ' + m; await loadStatus(); }
		catch (e) { error = e.message; } finally { loading = false; }
	}

	async function connect() {
		if (!peerAddress) { error = 'Enter a peer address'; return; }
		loading = true; error = '';
		try { await apiPost(`/spreadsheets/${sid()}/collab/peers/connect`, { address: peerAddress }); message = 'Connect request sent'; peerAddress = ''; await loadPeers(); }
		catch (e) { error = e.message; } finally { loading = false; }
	}

	async function sync() {
		loading = true; error = '';
		try { await apiPost(`/spreadsheets/${sid()}/collab/sync`, {}); message = 'Sync triggered'; await loadStatus(); }
		catch (e) { error = e.message; } finally { loading = false; }
	}

	$: if (spreadsheetId) loadAll();
</script>

<div class="p-4 space-y-4">
	{#if message}<div class="text-sm text-green-700 bg-green-50 rounded px-3 py-2">{message}</div>{/if}
	{#if error}<div class="text-sm text-red-700 bg-red-50 rounded px-3 py-2">{error}</div>{/if}

	<div class="border rounded-lg p-4">
		<h3 class="font-semibold mb-3">Collaboration Mode</h3>
		<div class="flex space-x-2">
			{#each ['local','p2p','server'] as m}
				<button class="px-3 py-1 rounded text-sm {mode === m ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700'}" on:click={() => setMode(m)}>{m}</button>
			{/each}
		</div>
	</div>

	<div class="border rounded-lg p-4">
		<h3 class="font-semibold mb-2">Connected Peers</h3>
		<ul class="text-sm space-y-1 mb-3">
			{#each peers as p}
				<li class="border-b pb-1">{p.address || p.id || 'peer'}</li>
			{/each}
			{#if !peers.length}<li class="text-gray-400">None connected</li>{/if}
		</ul>
		<input class="w-full border rounded px-2 py-1 text-sm mb-2" placeholder="Peer address" bind:value={peerAddress} />
		<button class="px-3 py-1 bg-blue-600 text-white rounded text-sm" on:click={connect} disabled={loading}>Connect to Peer</button>
	</div>

	<div class="border rounded-lg p-4">
		<h3 class="font-semibold mb-2">Discovered Peers</h3>
		<ul class="text-sm space-y-1 mb-3">
			{#each discovered as d}
				<li class="border-b pb-1">{d.address || d.id || 'peer'}</li>
			{/each}
			{#if !discovered.length}<li class="text-gray-400">None discovered</li>{/if}
		</ul>
	</div>

	<div class="border rounded-lg p-4">
		<h3 class="font-semibold mb-2">Status & Sync</h3>
		{#if status}<pre class="text-xs bg-gray-50 p-2 rounded mb-2">{JSON.stringify(status, null, 2)}</pre>{/if}
		<button class="px-3 py-1 bg-green-600 text-white rounded text-sm" on:click={sync} disabled={loading}>Sync with Peers</button>
	</div>
</div>
