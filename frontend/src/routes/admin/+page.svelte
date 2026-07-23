<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import {
		getAdminStats,
		getAdminUsers,
		updateAdminUserStatus,
		deleteAdminUser,
		getAdminLogs,
		getAdminDatabaseStats,
		runAdminDatabaseMaintenance,
		getAdminSecurityAlerts,
		resolveAdminSecurityAlert,
		getAdminSettings,
		updateAdminSettings,
		getOutboundDomainAllowlist,
		updateOutboundDomainAllowlist,
		formatApiError
	} from '$lib/api.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let tab = 'overview';
	let loading = false;
	let error = null;

	// Overview
	let stats = null;
	let dbStats = null;

	// Users
	let users = [];
	let userPagination = null;
	let userSearch = '';
	let userStatus = '';
	let userPage = 1;

	// Logs
	let logs = [];

	// Security
	let alerts = [];

	// Settings
	let settings = null;
	let settingsSaved = false;

	// Outbound domain allowlist
	let outboundDomains = '';
	let outboundSaved = false;
	let outboundError = null;

	onMount(load);

	async function load() {
		loading = true;
		error = null;
		try {
			if (tab === 'overview') {
				const [s, d] = await Promise.all([getAdminStats(), getAdminDatabaseStats()]);
				stats = s;
				dbStats = d;
			} else if (tab === 'users') {
				await loadUsers();
			} else if (tab === 'logs') {
				const r = await getAdminLogs({ limit: 100 });
				logs = r.logs || [];
			} else if (tab === 'security') {
				const r = await getAdminSecurityAlerts(50);
				alerts = r.alerts || [];
			} else if (tab === 'settings') {
				const r = await getAdminSettings();
				settings = r;
				try {
					const ol = await getOutboundDomainAllowlist();
					outboundDomains = (ol.domains || []).join('\n');
				} catch (e) {
					outboundDomains = '';
				}
			}
		} catch (err) {
			error = formatApiError(err);
		} finally {
			loading = false;
		}
	}

	async function loadUsers() {
		const r = await getAdminUsers({ page: userPage, limit: 25, search: userSearch, status: userStatus });
		users = r.users || [];
		userPagination = r.pagination || null;
	}

	function changeTab(t) {
		tab = t;
		load();
	}

	async function searchUsers() {
		userPage = 1;
		await loadUsers();
	}

	async function toggleUser(user) {
		try {
			await updateAdminUserStatus(user.id, !user.is_active);
			user.is_active = !user.is_active;
			users = users.slice();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function removeUser(user) {
		if (!confirm(`Delete user "${user.username}" and all their data? This cannot be undone.`)) return;
		try {
			await deleteAdminUser(user.id);
			users = users.filter(u => u.id !== user.id);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function runMaintenance() {
		if (!confirm('Run VACUUM ANALYZE on core tables?')) return;
		try {
			await runAdminDatabaseMaintenance();
			await load();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function resolveAlert(alert) {
		try {
			await resolveAdminSecurityAlert(alert.id);
			alert.resolved = true;
			alerts = alerts.slice();
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function saveSettings() {
		try {
			settingsSaved = false;
			await updateAdminSettings(settings);
			settingsSaved = true;
			setTimeout(() => (settingsSaved = false), 3000);
		} catch (err) {
			error = formatApiError(err);
		}
	}

	async function saveOutboundDomains() {
		try {
			outboundSaved = false;
			outboundError = null;
			const domains = outboundDomains
				.split('\n')
				.map((d) => d.trim())
				.filter((d) => d.length > 0);
			await updateOutboundDomainAllowlist(domains);
			outboundSaved = true;
			setTimeout(() => (outboundSaved = false), 3000);
		} catch (err) {
			outboundError = formatApiError(err);
		}
	}

	function fmtBytes(n) {
		if (n == null) return '—';
		const u = ['B', 'KB', 'MB', 'GB', 'TB'];
		let i = 0;
		let v = n;
		while (v >= 1024 && i < u.length - 1) {
			v /= 1024;
			i++;
		}
		return `${v.toFixed(2)} ${u[i]}`;
	}
</script>

<svelte:head>
	<title>Admin Console - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-6xl">
	<div class="flex items-center justify-between mb-6">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Admin Console</h1>
			<p class="text-sm text-gray-500 mt-1">System management and oversight</p>
		</div>
	</div>

	{#if error}
		<div class="mb-6 bg-red-50 border border-red-200 rounded-lg px-4 py-3 text-sm text-red-700">{error}</div>
	{/if}

	<div class="flex flex-wrap gap-1 border-b border-gray-200 mb-6">
		{#each ['overview', 'users', 'logs', 'security', 'settings'] as t}
			<button
				class="px-4 py-2 text-sm font-medium border-b-2 transition-colors {tab === t ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-600 hover:text-gray-900'}"
				on:click={() => changeTab(t)}
			>
				{t.charAt(0).toUpperCase() + t.slice(1)}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-16 text-gray-400"><span>Loading…</span></div>
	{:else if tab === 'overview'}
		<div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<p class="text-2xl font-bold text-gray-900">{stats?.users?.total ?? 0}</p>
				<p class="text-sm text-gray-500">Users (total)</p>
			</div>
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<p class="text-2xl font-bold text-green-600">{stats?.users?.active ?? 0}</p>
				<p class="text-sm text-gray-500">Active (30d)</p>
			</div>
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<p class="text-2xl font-bold text-gray-900">{stats?.content?.documents ?? 0}</p>
				<p class="text-sm text-gray-500">Documents</p>
			</div>
			<div class="bg-white border border-gray-200 rounded-lg p-4">
				<p class="text-2xl font-bold text-gray-900">{stats?.content?.tasks ?? 0}</p>
				<p class="text-sm text-gray-500">Tasks</p>
			</div>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div class="bg-white border border-gray-200 rounded-lg p-5">
				<h2 class="font-semibold text-gray-900 mb-3">Content Overview</h2>
				<dl class="text-sm divide-y divide-gray-100">
					<div class="flex justify-between py-1"><dt class="text-gray-500">Documents</dt><dd class="font-medium">{stats?.content?.documents ?? 0}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Emails</dt><dd class="font-medium">{stats?.content?.emails ?? 0}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Chat messages</dt><dd class="font-medium">{stats?.content?.chat_messages ?? 0}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Tasks</dt><dd class="font-medium">{stats?.content?.tasks ?? 0}</dd></div>
				</dl>
			</div>

			<div class="bg-white border border-gray-200 rounded-lg p-5">
				<div class="flex justify-between items-center mb-3">
					<h2 class="font-semibold text-gray-900">Database</h2>
					<button on:click={runMaintenance} class="text-xs bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded">Run VACUUM</button>
				</div>
				<dl class="text-sm divide-y divide-gray-100">
					<div class="flex justify-between py-1"><dt class="text-gray-500">Tables</dt><dd class="font-medium">{dbStats?.table_count ?? '—'}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Total size</dt><dd class="font-medium">{fmtBytes((dbStats?.total_size_mb || 0) * 1024 * 1024)}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Index size</dt><dd class="font-medium">{fmtBytes((dbStats?.index_size_mb || 0) * 1024 * 1024)}</dd></div>
					<div class="flex justify-between py-1"><dt class="text-gray-500">Active conns</dt><dd class="font-medium">{dbStats?.active_connections ?? '—'}</dd></div>
				</dl>
			</div>
		</div>
	{:else if tab === 'users'}
		<div class="flex flex-wrap gap-2 mb-4">
			<input
				type="text"
				bind:value={userSearch}
				placeholder="Search username, email, name…"
				class="flex-1 min-w-[200px] px-3 py-2 border border-gray-300 rounded text-sm"
				on:keydown={(e) => e.key === 'Enter' && searchUsers()}
			/>
			<select bind:value={userStatus} on:change={searchUsers} class="px-3 py-2 border border-gray-300 rounded text-sm">
				<option value="">All statuses</option>
				<option value="active">Active</option>
				<option value="inactive">Inactive</option>
				<option value="verified">Verified</option>
				<option value="unverified">Unverified</option>
				<option value="locked">Locked</option>
			</select>
			<button on:click={searchUsers} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Search</button>
		</div>

		<div class="bg-white border border-gray-200 rounded-lg overflow-hidden">
			<table class="w-full text-sm">
				<thead class="bg-gray-50 text-gray-600">
					<tr>
						<th class="text-left px-4 py-2">User</th>
						<th class="text-left px-4 py-2">Verified</th>
						<th class="text-left px-4 py-2">Status</th>
						<th class="text-left px-4 py-2">Created</th>
						<th class="text-right px-4 py-2">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-100">
					{#each users as u}
						<tr>
							<td class="px-4 py-2">
								<div class="font-medium text-gray-900">{u.username}</div>
								<div class="text-gray-400 text-xs">{u.email}</div>
							</td>
							<td class="px-4 py-2">{u.is_verified ? '✅' : '❌'}</td>
							<td class="px-4 py-2">
								<span class="inline-flex px-2 py-0.5 rounded-full text-xs {u.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'}">
									{u.is_active ? 'Active' : 'Inactive'}
								</span>
							</td>
							<td class="px-4 py-2 text-gray-500 text-xs">{u.created_at ? u.created_at.slice(0, 10) : '—'}</td>
							<td class="px-4 py-2 text-right whitespace-nowrap">
								<button on:click={() => toggleUser(u)} class="text-xs text-blue-600 hover:underline mr-2">
									{u.is_active ? 'Deactivate' : 'Activate'}
								</button>
								<button on:click={() => removeUser(u)} class="text-xs text-red-600 hover:underline">Delete</button>
							</td>
						</tr>
					{:else}
						<tr><td colspan="5" class="px-4 py-8 text-center text-gray-400">No users found</td></tr>
					{/each}
				</tbody>
			</table>
		</div>

		{#if userPagination && userPagination.totalPages > 1}
			<div class="flex items-center justify-center gap-3 mt-4 text-sm">
				<button
					on:click={() => { userPage = Math.max(1, userPage - 1); loadUsers(); }}
					disabled={userPage <= 1}
					class="px-3 py-1 border border-gray-300 rounded disabled:opacity-40"
				>Prev</button>
				<span class="text-gray-500">Page {userPage} / {userPagination.totalPages}</span>
				<button
					on:click={() => { userPage = Math.min(userPagination.totalPages, userPage + 1); loadUsers(); }}
					disabled={userPage >= userPagination.totalPages}
					class="px-3 py-1 border border-gray-300 rounded disabled:opacity-40"
				>Next</button>
			</div>
		{/if}
	{:else if tab === 'logs'}
		<div class="bg-white border border-gray-200 rounded-lg divide-y divide-gray-100 max-h-[70vh] overflow-auto">
			{#each logs as log}
				<div class="px-4 py-2 text-sm">
					<div class="flex items-center gap-2">
						<span class="text-xs text-gray-400">{log.timestamp}</span>
						<span class="inline-flex px-2 py-0.5 rounded-full text-xs bg-gray-100 text-gray-600">{log.event_type}</span>
						<span class="inline-flex px-2 py-0.5 rounded-full text-xs bg-blue-100 text-blue-700">{log.severity}</span>
						{#if log.action}<span class="text-xs text-gray-500">{log.action}</span>{/if}
					</div>
					{#if log.details}<p class="text-gray-500 text-xs mt-1">{log.details}</p>{/if}
				</div>
			{:else}
				<div class="px-4 py-8 text-center text-gray-400">No log entries</div>
			{/each}
		</div>
	{:else if tab === 'security'}
		<div class="space-y-3">
			{#each alerts as alert}
				<div class="bg-white border border-gray-200 rounded-lg p-4 flex items-start gap-3">
					<span class="mt-1 inline-flex px-2 py-0.5 rounded-full text-xs {alert.resolved ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}">
						{alert.severity}
					</span>
					<div class="flex-1">
						<p class="font-medium text-gray-900">{alert.event_type}</p>
						<p class="text-sm text-gray-500">{alert.details}</p>
						<p class="text-xs text-gray-400 mt-1">{alert.created_at}</p>
					</div>
					{#if !alert.resolved}
						<button on:click={() => resolveAlert(alert)} class="text-xs text-blue-600 hover:underline">Resolve</button>
					{/if}
				</div>
			{:else}
				<div class="px-4 py-8 text-center text-gray-400">No security alerts</div>
			{/each}
		</div>
	{:else if tab === 'settings'}
		{#if settings}
			<div class="bg-white border border-gray-200 rounded-lg p-6 max-w-2xl space-y-6">
				{#if settingsSaved}
					<div class="bg-green-50 border border-green-200 rounded px-3 py-2 text-sm text-green-700">Settings saved</div>
				{/if}

				<section>
					<h2 class="font-semibold text-gray-900 mb-2">Security</h2>
					<div class="grid grid-cols-2 gap-3">
						<label class="text-sm">Max login attempts
							<input type="number" bind:value={settings.security.max_login_attempts} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
						</label>
						<label class="text-sm">Lockout (minutes)
							<input type="number" bind:value={settings.security.lockout_duration_minutes} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
						</label>
						<label class="text-sm">Min password length
							<input type="number" bind:value={settings.security.password_min_length} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
						</label>
						<label class="text-sm">Session timeout (hours)
							<input type="number" bind:value={settings.security.session_timeout_hours} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
						</label>
					</div>
				</section>

				<section>
					<h2 class="font-semibold text-gray-900 mb-2">Storage</h2>
					<div class="grid grid-cols-2 gap-3">
						<label class="text-sm">Max file size (MB)
							<input type="number" bind:value={settings.storage.max_file_size_mb} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
						</label>
						<label class="text-sm">Backup retention (days)
							<input type="number" bind:value={settings.storage.backup_retention_days} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
						</label>
						<label class="text-sm">Storage quota (GB)
							<input type="number" bind:value={settings.storage.storage_quota_gb} class="mt-1 w-full px-3 py-2 border border-gray-300 rounded" />
						</label>
					</div>
				</section>

				<section>
					<h2 class="font-semibold text-gray-900 mb-2">Features</h2>
					<div class="space-y-2">
						{#each Object.keys(settings.features) as key}
							<label class="flex items-center gap-2 text-sm">
								<input type="checkbox" bind:checked={settings.features[key]} class="rounded" />
								{key.replace(/_/g, ' ')}
							</label>
						{/each}
					</div>
				</section>

				<section>
					<h2 class="font-semibold text-gray-900 mb-2">Outbound domain allowlist</h2>
					<p class="text-xs text-gray-500 mb-2">
						Restrict which external domains workflows may call via the <code>http.request</code> connector.
						One domain per line (e.g. <code>api.example.com</code> or <code>*.example.com</code>). Leave empty to allow all public destinations (SSRF protection still applies).
					</p>
					<textarea bind:value={outboundDomains} rows="4" placeholder={"api.example.com\n*.other.com"} class="w-full px-3 py-2 border border-gray-300 rounded font-mono text-sm"></textarea>
					{#if outboundSaved}
						<p class="text-sm text-green-700 mt-1">Outbound domain allowlist saved</p>
					{/if}
					{#if outboundError}
						<p class="text-sm text-red-600 mt-1">{outboundError}</p>
					{/if}
					<button on:click={saveOutboundDomains} class="mt-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Save allowlist</button>
				</section>

				<button on:click={saveSettings} class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm">Save settings</button>
			</div>
		{/if}
	{/if}
</div>
