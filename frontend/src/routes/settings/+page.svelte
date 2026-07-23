<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { apiGet, apiPost, apiPut } from '$lib/api.js';
	import { userPreferences } from '$lib/stores.js';

	export let data = null;
	export let form = null;
	export let params = null;

	let activeTab = 'profile';

	// Profile
	let profile = { username: '', email: '', first_name: '', last_name: '' };
	let isLoadingProfile = false;
	let profileError = null;
	let profileSuccess = null;

	// Password change
	let passwordForm = { old_password: '', new_password: '', confirm_password: '' };
	let passwordError = null;
	let passwordSuccess = null;
	let isChangingPassword = false;

	// Preferences (local store)
	let prefs = {};
	userPreferences.subscribe(v => prefs = { ...v });

	// 2FA
	let twoFAStatus = null;
	let isLoadingTwoFA = false;

	onMount(async () => {
		await loadProfile();
		await loadTwoFAStatus();
	});

	async function loadProfile() {
		isLoadingProfile = true;
		profileError = null;
		try {
			const data = await apiGet('/auth/profile');
			profile = {
				username: data.user?.username || data.username || '',
				email: data.user?.email || data.email || '',
				first_name: data.user?.first_name || data.first_name || '',
				last_name: data.user?.last_name || data.last_name || '',
			};
		} catch (err) {
			if (err?.status === 401) { goto('/auth/login'); return; }
			profileError = 'Could not load profile.';
		} finally {
			isLoadingProfile = false;
		}
	}

	async function saveProfile() {
		profileError = null;
		profileSuccess = null;
		try {
			await apiPut('/auth/profile', {
				first_name: profile.first_name,
				last_name: profile.last_name,
				username: profile.username,
			});
			profileSuccess = 'Profile updated.';
			setTimeout(() => profileSuccess = null, 3000);
		} catch (err) {
			profileError = err?.message || 'Failed to update profile.';
		}
	}

	async function changePassword() {
		passwordError = null;
		passwordSuccess = null;
		if (passwordForm.new_password !== passwordForm.confirm_password) {
			passwordError = 'New passwords do not match.';
			return;
		}
		if (passwordForm.new_password.length < 8) {
			passwordError = 'Password must be at least 8 characters.';
			return;
		}
		isChangingPassword = true;
		try {
			await apiPost('/auth/change-password', {
				old_password: passwordForm.old_password,
				new_password: passwordForm.new_password,
			});
			passwordSuccess = 'Password changed successfully.';
			passwordForm = { old_password: '', new_password: '', confirm_password: '' };
			setTimeout(() => passwordSuccess = null, 3000);
		} catch (err) {
			passwordError = err?.message || 'Failed to change password.';
		} finally {
			isChangingPassword = false;
		}
	}

	async function loadTwoFAStatus() {
		isLoadingTwoFA = true;
		try {
			const data = await apiGet('/auth/2fa/status');
			twoFAStatus = data;
		} catch (err) {
			console.error('Failed to load 2FA status:', err);
		} finally {
			isLoadingTwoFA = false;
		}
	}

	function savePreferences() {
		userPreferences.set(prefs);
		if (typeof document !== 'undefined') {
			document.documentElement.classList.toggle('dark', prefs.theme === 'dark');
		}
	}

	const tabs = [
		{ id: 'profile', label: 'Profile' },
		{ id: 'security', label: 'Security' },
		{ id: 'preferences', label: 'Preferences' },
	];
</script>

<svelte:head>
	<title>Settings - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-3xl">
	<h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">Settings</h1>

	<!-- Tabs -->
	<div class="border-b border-gray-200 dark:border-gray-700 mb-8">
		<nav class="flex gap-1 -mb-px">
			{#each tabs as tab}
				<button
					on:click={() => activeTab = tab.id}
					class="px-4 py-3 text-sm font-medium border-b-2 transition-colors {activeTab === tab.id
						? 'border-blue-600 text-blue-600'
						: 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300'}"
				>
					{tab.label}
				</button>
			{/each}
		</nav>
	</div>

	<!-- Profile Tab -->
	{#if activeTab === 'profile'}
		<div class="space-y-6">
			{#if isLoadingProfile}
				<div class="flex items-center gap-2 text-gray-400 py-8 justify-center">
					<svg class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
					</svg>
					Loading…
				</div>
			{:else}
				{#if profileError}
					<div class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg px-4 py-3 text-sm text-red-700 dark:text-red-300">{profileError}</div>
				{/if}
				{#if profileSuccess}
					<div class="bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-700 rounded-lg px-4 py-3 text-sm text-green-700 dark:text-green-300">{profileSuccess}</div>
				{/if}

				<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-4">
					<h2 class="font-semibold text-gray-900 dark:text-white">Personal Information</h2>
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">First name</label>
							<input bind:value={profile.first_name} class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Last name</label>
							<input bind:value={profile.last_name} class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
						</div>
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Username</label>
						<input bind:value={profile.username} class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email</label>
						<input value={profile.email} disabled class="w-full px-3 py-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-900 text-gray-500 dark:text-gray-400 text-sm cursor-not-allowed" />
						<p class="mt-1 text-xs text-gray-400">Email cannot be changed here. Contact an administrator.</p>
					</div>
					<button on:click={saveProfile} class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors">
						Save Changes
					</button>
				</div>
			{/if}
		</div>
	{/if}

	<!-- Security Tab -->
	{#if activeTab === 'security'}
		<div class="space-y-6">
			<!-- Change Password -->
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-4">
				<h2 class="font-semibold text-gray-900 dark:text-white">Change Password</h2>
				{#if passwordError}
					<div class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg px-4 py-3 text-sm text-red-700 dark:text-red-300">{passwordError}</div>
				{/if}
				{#if passwordSuccess}
					<div class="bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-700 rounded-lg px-4 py-3 text-sm text-green-700 dark:text-green-300">{passwordSuccess}</div>
				{/if}
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Current password</label>
					<input type="password" bind:value={passwordForm.old_password} class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">New password</label>
					<input type="password" bind:value={passwordForm.new_password} class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Confirm new password</label>
					<input type="password" bind:value={passwordForm.confirm_password} class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
				</div>
				<button on:click={changePassword} disabled={isChangingPassword} class="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white text-sm rounded-lg transition-colors">
					{isChangingPassword ? 'Changing…' : 'Change Password'}
				</button>
			</div>

			<!-- 2FA -->
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-4">
				<h2 class="font-semibold text-gray-900 dark:text-white">Two-Factor Authentication</h2>
				{#if isLoadingTwoFA}
					<p class="text-sm text-gray-400">Loading…</p>
				{:else if twoFAStatus?.enabled}
					<div class="flex items-center gap-3">
						<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
							Enabled
						</span>
						<span class="text-sm text-gray-600 dark:text-gray-400">Authenticator app is active.</span>
					</div>
					<a href="/auth/2fa/disable" class="inline-block px-4 py-2 border border-red-300 text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 text-sm rounded-lg transition-colors">
						Disable 2FA
					</a>
				{:else}
					<p class="text-sm text-gray-600 dark:text-gray-400">Protect your account with a one-time code in addition to your password.</p>
					<a href="/auth/2fa/setup" class="inline-block px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors">
						Enable 2FA
					</a>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Preferences Tab -->
	{#if activeTab === 'preferences'}
		<div class="space-y-6">
			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-5">
				<h2 class="font-semibold text-gray-900 dark:text-white">Appearance</h2>

				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-300">Theme</p>
						<p class="text-xs text-gray-400">Light or dark interface</p>
					</div>
					<select bind:value={prefs.theme} class="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm focus:ring-2 focus:ring-blue-500">
						<option value="light">Light</option>
						<option value="dark">Dark</option>
						<option value="system">System</option>
					</select>
				</div>

				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-300">Enable animations</p>
						<p class="text-xs text-gray-400">Smooth transitions and effects</p>
					</div>
					<label class="relative inline-flex items-center cursor-pointer">
						<input type="checkbox" bind:checked={prefs.enableAnimations} class="sr-only peer" />
						<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
					</label>
				</div>

				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-300">Touch mode</p>
						<p class="text-xs text-gray-400">Larger tap targets for touch screens</p>
					</div>
					<label class="relative inline-flex items-center cursor-pointer">
						<input type="checkbox" bind:checked={prefs.touchMode} class="sr-only peer" />
						<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
					</label>
				</div>
			</div>

			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-5">
				<h2 class="font-semibold text-gray-900 dark:text-white">Regional</h2>

				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-300">Locale</p>
					</div>
					<select bind:value={prefs.locale} class="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm focus:ring-2 focus:ring-blue-500">
						<option value="en-US">English (US)</option>
						<option value="en-GB">English (UK)</option>
						<option value="en-NZ">English (NZ)</option>
						<option value="fr-FR">French</option>
						<option value="de-DE">German</option>
						<option value="es-ES">Spanish</option>
					</select>
				</div>

				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-300">Date format</p>
					</div>
					<select bind:value={prefs.dateFormat} class="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm focus:ring-2 focus:ring-blue-500">
						<option value="MM/DD/YYYY">MM/DD/YYYY</option>
						<option value="DD/MM/YYYY">DD/MM/YYYY</option>
						<option value="YYYY-MM-DD">YYYY-MM-DD</option>
					</select>
				</div>
			</div>

			<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-5">
				<h2 class="font-semibold text-gray-900 dark:text-white">Editor</h2>

				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-300">Auto-save</p>
						<p class="text-xs text-gray-400">Automatically save changes as you work</p>
					</div>
					<label class="relative inline-flex items-center cursor-pointer">
						<input type="checkbox" bind:checked={prefs.autoSave} class="sr-only peer" />
						<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
					</label>
				</div>

				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-700 dark:text-gray-300">Show formula bar</p>
						<p class="text-xs text-gray-400">Display formula bar in spreadsheet</p>
					</div>
					<label class="relative inline-flex items-center cursor-pointer">
						<input type="checkbox" bind:checked={prefs.showFormulaBar} class="sr-only peer" />
						<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
					</label>
				</div>
			</div>

			<button on:click={savePreferences} class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors">
				Save Preferences
			</button>
		</div>
	{/if}
</div>
