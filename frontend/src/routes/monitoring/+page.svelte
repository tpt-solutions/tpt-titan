<script>
	import { onMount, onDestroy } from 'svelte';
	import {
		getMonitoringMetrics,
		getMonitoringHealth,
		getMonitoringAlerts,
		getWebhookDeliveryLogs,
		formatApiError
	} from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let metrics = null;
	let health = null;
	let alerts = [];
	let error = null;
	let loading = false;
	let timer = null;
	let deliveryLogs = [];

	onMount(refresh);
	onDestroy(() => clearInterval(timer));

	async function refresh() {
		loading = true;
		error = null;
		try {
			const [m, h, a] = await Promise.all([
				getMonitoringMetrics(),
				getMonitoringHealth(),
				getMonitoringAlerts(20)
			]);
			metrics = m;
			health = h;
			alerts = a.alerts || [];
			try {
				const logs = await getWebhookDeliveryLogs();
				deliveryLogs = logs.logs || [];
			} catch (e) {
				deliveryLogs = [];
			}
		} catch (err) {
			error = formatApiError(err);
		} finally {
			loading = false;
		}
	}

	function startAuto() {
		if (timer) clearInterval(timer);
		timer = setInterval(refresh, 10000);
	}

	function stopAuto() {
		clearInterval(timer);
		timer = null;
	}

	function fmtBytes(n) {
		if (n == null) return '—';
		const u = ['B', 'KB', 'MB', 'GB', 'TB'];
		let i = 0;
		let v = n;
		while (v >= 1024 && i < u.length - 1) { v /= 1024; i++; }
		return `${v.toFixed(2)} ${u[i]}`;
	}

	function fmtDuration(s) {
		if (s == null) return '—';
		const sec = Math.floor(s);
		const h = Math.floor(sec / 3600);
		const m = Math.floor((sec % 3600) / 60);
		const secs = sec % 60;
		return `${h}h ${m}m ${secs}s`;
	}

	$: app = metrics?.application;
	$: sys = metrics?.system;
	$: db = metrics?.database;
</script>

<svelte:head>
	<title>Monitoring - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-6xl">
	<div class="flex items-center justify-between mb-6">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Monitoring</h1>
			<p class="text-sm text-gray-500 mt-1">System health and performance</p>
		</div>
		<div class="flex items-center gap-2">
			<button on:click={refresh} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Refresh</button>
			<button on:click={timer ? stopAuto() : startAuto()} class="bg-gray-700 hover:bg-gray-800 text-white px-4 py-2 rounded text-sm">
				{timer ? 'Stop auto' : 'Auto (10s)'}
			</button>
		</div>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 border border-red-200 rounded-lg px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	{#if loading && !metrics}
		<div class="flex items-center justify-center py-16 text-gray-400"><span>Loading…</span></div>
	{:else if health}
		<div class="mb-6">
			<span class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium {
				health.status === 'healthy' ? 'bg-green-100 text-green-700' :
				health.status === 'degraded' ? 'bg-yellow-100 text-yellow-700' : 'bg-red-100 text-red-700'
			}">
				● {health.status}
			</span>
			<div class="mt-2 flex flex-wrap gap-2">
				{#each Object.entries(health.services || {}) as [svc, st]}
					<span class="inline-flex px-2 py-0.5 rounded text-xs bg-gray-100 text-gray-600">{svc}: {st}</span>
				{/each}
			</div>
		</div>

		<div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<p class="text-2xl font-bold text-gray-900">{fmtDuration(app?.uptime)}</p>
				<p class="text-sm text-gray-500">Uptime</p>
			</div>
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<p class="text-2xl font-bold text-gray-900">{(app?.total_requests ?? 0).toLocaleString()}</p>
				<p class="text-sm text-gray-500">Total requests</p>
			</div>
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<p class="text-2xl font-bold text-gray-900">{((app?.error_rate ?? 0) * 100).toFixed(2)}%</p>
				<p class="text-sm text-gray-500">Error rate</p>
			</div>
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<p class="text-2xl font-bold text-gray-900">{sys?.goroutines ?? 0}</p>
				<p class="text-sm text-gray-500">Goroutines</p>
			</div>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<h2 class="font-semibold text-gray-900 mb-2">System</h2>
				<dl class="text-sm divide-y divide-gray-100">
					<div class="flex justify-between py-1"><dt class="text-gray-500">CPU</dt><dd>{sys?.cpu_usage_percent?.toFixed(1) ?? '—'}%</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Memory</dt><dd>{sys?.memory_usage_percent?.toFixed(1) ?? '—'}%</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Disk</dt><dd>{sys?.disk_usage_percent?.toFixed(1) ?? '—'}%</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Mem used</dt><dd>{fmtBytes(sys?.memory_used_bytes)}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Heap alloc</dt><dd>{fmtBytes(sys?.heap_alloc_bytes)}</dd></div>
				</dl>
			</div>

			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<h2 class="font-semibold text-gray-900 mb-2">Application</h2>
				<dl class="text-sm divide-y divide-gray-100">
					<div class="flex justify-between py-1"><dt class="text-gray-500">Active users</dt><dd>{app?.active_users ?? 0}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">DB connections</dt><dd>{app?.database_connections ?? '—'}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Cache hit rate</dt><dd>{((app?.cache_hit_rate ?? 0) * 100).toFixed(1)}%</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Cache keys</dt><dd>{app?.cache_keys ?? 0}</dd></div>
				</dl>
			</div>

			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<h2 class="font-semibold text-gray-900 mb-2">Database</h2>
				<dl class="text-sm divide-y divide-gray-100">
					<div class="flex justify-between py-1"><dt class="text-gray-500">Total conns</dt><dd>{db?.total_connections ?? '—'}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Active conns</dt><dd>{db?.active_connections ?? '—'}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Idle conns</dt><dd>{db?.idle_connections ?? '—'}</dd></div>
				</dl>
			</div>
		</div>

		<div class="mt-6">
			<h2 class="font-semibold text-gray-900 mb-2">Recent alerts</h2>
			<div class="space-y-2">
				{#each alerts as a}
					<div class="bg-white border border-gray-200 rounded-lg p-3 flex items-center gap-3">
						<span class="inline-flex px-2 py-0.5 rounded-full text-xs {a.resolved ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}">{a.level}</span>
						<div class="flex-1">
							<p class="font-medium text-gray-900">{a.title}</p>
							<p class="text-xs text-gray-400">{a.message}</p>
						</div>
					</div>
				{:else}
					<div class="text-center text-gray-400 py-6">No alerts</div>
				{/each}
			</div>
		</div>

		<div class="mt-6">
			<h2 class="font-semibold text-gray-900 mb-2">Webhook delivery log</h2>
			<div class="overflow-x-auto bg-white border border-gray-200 rounded-lg">
				<table class="w-full text-sm">
					<thead class="bg-gray-50 text-gray-500 text-left">
						<tr>
							<th class="px-3 py-2">Time</th>
							<th class="px-3 py-2">Dir</th>
							<th class="px-3 py-2">Host</th>
							<th class="px-3 py-2">Method</th>
							<th class="px-3 py-2">Status</th>
							<th class="px-3 py-2">Connector</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-100">
						{#each deliveryLogs as l}
							<tr>
								<td class="px-3 py-2 text-gray-500">{new Date(l.created_at).toLocaleString()}</td>
								<td class="px-3 py-2">{l.direction}</td>
								<td class="px-3 py-2 max-w-xs truncate">{l.host || '—'}</td>
								<td class="px-3 py-2">{l.method || '—'}</td>
								<td class="px-3 py-2">{l.status_code || (l.error ? 'ERR' : '—')}</td>
								<td class="px-3 py-2">{l.connector}</td>
							</tr>
						{:else}
							<tr><td colspan="6" class="px-3 py-6 text-center text-gray-400">No webhook calls recorded yet</td></tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>
