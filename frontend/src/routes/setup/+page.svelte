<script>
// @ts-nocheck
	import { goto } from '$app/navigation';

	// Accept framework-provided props to avoid warnings
	export let data = null;
	export let form = null;
	export let params = null;

	let step = 1;
	let loading = false;
	let checking = true;
	let error = '';
	let setupRequired = false;

	let admin = {
		email: '',
		username: '',
		password: '',
		confirmPassword: '',
		firstName: '',
		lastName: ''
	};

	let smtp = {
		enabled: false,
		host: '',
		port: '587',
		username: '',
		password: ''
	};

	async function checkStatus() {
		checking = true;
		error = '';
		try {
			const res = await fetch('/api/v1/setup/status');
			const body = await res.json();
			setupRequired = !!body.setup_required;
			if (!setupRequired) {
				goto('/auth/login');
			}
		} catch (err) {
			error = 'Could not reach the server. Is the backend running?';
		} finally {
			checking = false;
		}
	}

	function validateAdmin() {
		if (!admin.email || !admin.username || !admin.password) {
			error = 'Email, username and password are required.';
			return false;
		}
		if (admin.password !== admin.confirmPassword) {
			error = 'Passwords do not match.';
			return false;
		}
		if (admin.password.length < 8) {
			error = 'Password must be at least 8 characters.';
			return false;
		}
		if (admin.username.length < 3) {
			error = 'Username must be at least 3 characters.';
			return false;
		}
		return true;
	}

	async function finishSetup() {
		if (!validateAdmin()) return;

		loading = true;
		error = '';
		try {
			const payload = {
				email: admin.email.trim(),
				username: admin.username.trim(),
				password: admin.password,
				first_name: admin.firstName.trim(),
				last_name: admin.lastName.trim()
			};
			if (smtp.enabled) {
				payload.smtp = {
					host: smtp.host.trim(),
					port: smtp.port.trim(),
					username: smtp.username.trim(),
					password: smtp.password
				};
			}

			const res = await fetch('/api/v1/setup/complete', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(payload)
			});
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || 'Setup failed');
			}
			step = 3;
		} catch (err) {
			error = err.message || 'Setup failed. Please try again.';
		} finally {
			loading = false;
		}
	}

	checkStatus();
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<div class="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-blue-100">
				<svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
				</svg>
			</div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">TPT Titan Setup</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				Welcome! Let's get your private office suite configured.
			</p>
		</div>

		{#if checking}
			<p class="text-center text-sm text-gray-500">Checking setup status…</p>
		{:else if error && step !== 3}
			<div class="rounded-md bg-red-50 p-4">
				<p class="text-sm font-medium text-red-800">{error}</p>
			</div>
		{/if}

		{#if !checking && step === 1}
			<div class="space-y-4">
				<p class="text-sm text-gray-600">Step 1 of 2 — Create the first administrator account.</p>
				<div>
					<label class="block text-sm font-medium text-gray-700">Email address *</label>
					<input type="email" bind:value={admin.email} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="admin@yourdomain.com" />
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700">Username *</label>
					<input type="text" bind:value={admin.username} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="admin" />
				</div>
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="block text-sm font-medium text-gray-700">First name</label>
						<input type="text" bind:value={admin.firstName} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-700">Last name</label>
						<input type="text" bind:value={admin.lastName} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
					</div>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700">Password *</label>
					<input type="password" bind:value={admin.password} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="••••••••" />
					<p class="mt-1 text-xs text-gray-500">At least 8 characters.</p>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700">Confirm password *</label>
					<input type="password" bind:value={admin.confirmPassword} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="••••••••" />
				</div>
				<button on:click={() => { if (validateAdmin()) step = 2; }} class="w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
					Next: Email (optional)
				</button>
			</div>
		{:else if !checking && step === 2}
			<div class="space-y-4">
				<p class="text-sm text-gray-600">Step 2 of 2 — Optionally configure outbound email (SMTP). You can skip this and set it later in Settings.</p>
				<label class="flex items-center space-x-2">
					<input type="checkbox" bind:checked={smtp.enabled} class="rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
					<span class="text-sm text-gray-700">Configure SMTP now</span>
				</label>
				{#if smtp.enabled}
					<div>
						<label class="block text-sm font-medium text-gray-700">SMTP host</label>
						<input type="text" bind:value={smtp.host} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="smtp.gmail.com" />
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-700">SMTP port</label>
						<input type="text" bind:value={smtp.port} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="587" />
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-700">SMTP username</label>
						<input type="text" bind:value={smtp.username} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-700">SMTP password</label>
						<input type="password" bind:value={smtp.password} class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
					</div>
				{/if}
				<div class="flex space-x-3">
					<button on:click={() => step = 1} class="flex-1 py-2 px-4 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50">Back</button>
					<button on:click={finishSetup} disabled={loading} class="flex-1 flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50">
						{#if loading}Setting up…{:else}Finish setup{/if}
					</button>
				</div>
			</div>
		{:else if step === 3}
			<div class="rounded-md bg-green-50 p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-sm font-medium text-green-800">Setup complete!</h3>
						<p class="mt-2 text-sm text-green-700">Your administrator account is ready. You can now sign in.</p>
						<a href="/auth/login" class="mt-3 inline-block font-medium text-green-600 hover:text-green-500">Continue to sign in →</a>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
