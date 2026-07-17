<script>
	import { onMount } from 'svelte';
	import {
		getMCPServers,
		createMCPServer,
		deleteMCPServer,
		testMCPServer,
		formatApiError
	} from '$lib/api.js';

	let servers = [];
	let loading = false;
	let error = null;
	let testing = false;
	let testResult = null;

	// form state
	let name = '';
	let url = '';
	let transport = 'http';
	let authType = 'none';
	let token = '';

	onMount(load);

	async function load() {
		loading = true;
		error = null;
		try {
			const r = await getMCPServers();
			servers = r.servers || [];
		} catch (err) {
			error = formatApiError(err);
		} finally {
			loading = false;
		}
	}

	async function create() {
		error = null;
		if (!name || !url) {
			error = 'Name and URL are required.';
			return;
		}
		try {
			await createMCPServer({ name, url, transport, auth_type: authType, token });
			name = url = token = '';
			authType = 'none';
			transport = 'http';
			await load();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function remove(s) {
		if (!confirm(`Delete MCP server "${s.name}"?`)) return;
		error = null;
		try {
			await deleteMCPServer(s.id);
			await load();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function test() {
		error = null;
		testing = true;
		testResult = null;
		try {
			const r = await testMCPServer({ url, transport, auth_type: authType, token });
			testResult = r;
		} catch (err) {
			testResult = { ok: false, error: formatApiError(err) };
		} finally {
			testing = false;
		}
	}
</script>

<svelte:head>
	<title>MCP Servers - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<div class="mb-6">
		<h1 class="text-3xl font-bold text-gray-900">MCP Servers</h1>
		<p class="text-sm text-gray-500 mt-1">
			Connect external systems (including the sibling tpt-free-erp project) via the
			Model Context Protocol. Each exposed tool becomes a workflow connector named
			<code class="bg-gray-100 px-1 rounded">mcp.&lt;server&gt;.&lt;tool&gt;</code>.
		</p>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 border border-red-200 rounded-lg px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	<div class="bg-white border border-gray-200 rounded-lg p-4 mb-6">
		<h2 class="font-semibold text-gray-900 mb-3">Add server</h2>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
			<div>
				<label class="block text-xs font-medium text-gray-600 mb-1">Name</label>
				<input class="w-full border border-gray-300 rounded px-3 py-2 text-sm" bind:value={name} placeholder="My ERP" />
			</div>
			<div>
				<label class="block text-xs font-medium text-gray-600 mb-1">URL</label>
				<input class="w-full border border-gray-300 rounded px-3 py-2 text-sm" bind:value={url} placeholder="https://mcp.example.com/jsonrpc" />
			</div>
			<div>
				<label class="block text-xs font-medium text-gray-600 mb-1">Transport</label>
				<select class="w-full border border-gray-300 rounded px-3 py-2 text-sm" bind:value={transport}>
					<option value="http">http</option>
					<option value="streamable-http">streamable-http</option>
				</select>
			</div>
			<div>
				<label class="block text-xs font-medium text-gray-600 mb-1">Auth</label>
				<select class="w-full border border-gray-300 rounded px-3 py-2 text-sm" bind:value={authType}>
					<option value="none">none</option>
					<option value="bearer">bearer token</option>
				</select>
			</div>
			{#if authType === 'bearer'}
				<div class="md:col-span-2">
					<label class="block text-xs font-medium text-gray-600 mb-1">Token (encrypted at rest)</label>
					<input type="password" class="w-full border border-gray-300 rounded px-3 py-2 text-sm" bind:value={token} placeholder="••••••••" />
				</div>
			{/if}
		</div>
		<div class="flex gap-2 mt-3">
			<button on:click={create} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Add server</button>
			<button on:click={test} disabled={testing} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm disabled:opacity-50">
				{testing ? 'Testing…' : 'Test connection'}
			</button>
		</div>
		{#if testResult}
			<div class="mt-3 text-sm {testResult.ok ? 'text-green-700' : 'text-red-700'}">
				{#if testResult.ok}
					Connected — {testResult.tool_count} tool(s) available.
				{:else}
					Connection failed: {testResult.error}
				{/if}
			</div>
		{/if}
	</div>

	<h2 class="font-semibold text-gray-900 mb-2">Configured servers</h2>
	{#if loading}
		<p class="text-gray-400 text-sm">Loading…</p>
	{:else if servers.length === 0}
		<p class="text-gray-400 text-sm">No MCP servers configured yet.</p>
	{:else}
		<div class="space-y-2">
			{#each servers as s}
				<div class="bg-white border border-gray-200 rounded-lg p-4 flex items-start gap-3">
					<div class="flex-1">
						<p class="font-medium text-gray-900">{s.name}</p>
						<p class="text-xs text-gray-400">{s.url}</p>
						<p class="text-xs text-gray-400 mt-1">{s.transport} · auth: {s.auth_type} · {s.is_active ? 'active' : 'inactive'}</p>
					</div>
					<button on:click={() => remove(s)} class="text-xs text-red-600 hover:underline">Delete</button>
				</div>
			{/each}
		</div>
	{/if}

	<p class="text-xs text-gray-400 mt-6">
		After adding a server, its tools appear as connectors in the workflow builder's action node dropdown (under the
		<code class="bg-gray-100 px-1 rounded">mcp.</code> prefix). Pass the tool's arguments as a JSON object in the node's
		parameters field.
	</p>
</div>
