<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import { apiGet, apiPost } from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let plugins = [];
	let stats = null;
	let isLoading = false;
	let error = null;

	onMount(load);

	async function load() {
		isLoading = true;
		error = null;
		try {
			const [p, s] = await Promise.all([
				apiGet('/plugins'),
				apiGet('/plugins/stats')
			]);
			plugins = p.plugins || [];
			stats = s.stats || null;
		} catch (err) {
			error = 'Could not load plugins. Please check your connection.';
		} finally {
			isLoading = false;
		}
	}

	async function toggle(id, enable) {
		try {
			if (enable) {
				await apiPost(`/plugins/${id}/enable`, {});
			} else {
				await apiPost(`/plugins/${id}/disable`, {});
			}
			plugins = plugins.map(pl => pl.id === id ? { ...pl, enabled: enable } : pl);
		} catch (err) {
			error = 'Failed to update plugin. Please try again.';
		}
	}

	async function uninstall(id) {
		if (!confirm('Uninstall this plugin? This cannot be undone.')) return;
		try {
			await apiPost(`/plugins/${id}/unload`, {});
			plugins = plugins.filter(pl => pl.id !== id);
		} catch (err) {
			error = 'Failed to uninstall plugin. Please try again.';
		}
	}
</script>

<svelte:head>
	<title>Plugins - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-5xl">
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Plugins</h1>
			<p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Extend TPT Titan with plugins</p>
		</div>
		<button
			on:click={load}
			class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm transition-colors"
		>
			Refresh
		</button>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg px-4 py-3 text-sm text-red-700 dark:text-red-300">
			{error}
		</div>
	{/if}

	{#if stats}
		<div class="mb-6 grid grid-cols-2 md:grid-cols-4 gap-4">
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 text-center">
				<p class="text-2xl font-bold text-gray-900 dark:text-white">{stats.total_plugins ?? 0}</p>
				<p class="text-sm text-gray-500 dark:text-gray-400">Total</p>
			</div>
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 text-center">
				<p class="text-2xl font-bold text-green-600">{stats.enabled_plugins ?? 0}</p>
				<p class="text-sm text-gray-500 dark:text-gray-400">Enabled</p>
			</div>
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 text-center">
				<p class="text-2xl font-bold text-gray-900 dark:text-white">{stats.total_hooks ?? 0}</p>
				<p class="text-sm text-gray-500 dark:text-gray-400">Hooks</p>
			</div>
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 text-center">
				<p class="text-2xl font-bold text-gray-900 dark:text-white">{stats.total_apis ?? 0}</p>
				<p class="text-sm text-gray-500 dark:text-gray-400">APIs</p>
			</div>
		</div>
	{/if}

	{#if isLoading}
		<div class="flex items-center justify-center py-16 text-gray-400">
			<span>Loading plugins…</span>
		</div>
	{:else if plugins.length === 0}
		<div class="text-center py-16">
			<svg class="mx-auto w-12 h-12 text-gray-300 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
			</svg>
			<h3 class="text-gray-500 dark:text-gray-400 font-medium mb-1">No plugins loaded</h3>
			<p class="text-sm text-gray-400 dark:text-gray-500">Place compiled <code>.so</code> plugins in the plugins directory and restart the server.</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each plugins as plugin}
				<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 flex items-center gap-4">
					<div class="flex-1 min-w-0">
						<p class="font-medium text-gray-900 dark:text-white truncate">
							{plugin.name}
							<span class="ml-2 text-xs text-gray-400">v{plugin.version}</span>
						</p>
						<p class="text-sm text-gray-500 dark:text-gray-400 truncate">{plugin.description}</p>
						{#if plugin.author}
							<p class="text-xs text-gray-400 mt-1">by {plugin.author}</p>
						{/if}
					</div>
					<div class="flex items-center gap-2 flex-shrink-0">
						<label class="inline-flex items-center cursor-pointer">
							<input
								type="checkbox"
								class="sr-only peer"
								checked={plugin.enabled}
								on:change={(e) => toggle(plugin.id, e.target.checked)}
							/>
							<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600 relative"></div>
						</label>
						<button
							on:click={() => uninstall(plugin.id)}
							class="p-2 text-gray-400 hover:text-red-600 transition-colors"
							title="Uninstall plugin"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
							</svg>
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
