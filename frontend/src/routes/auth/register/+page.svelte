<script>
	import { goto } from '$app/navigation';

	// Accept framework-provided props to avoid warnings
	export let data = null;
	export let form = null;
	export let params = null;


	let formData = {
		email: '',
		username: '',
		password: '',
		confirmPassword: '',
		firstName: '',
		lastName: ''
	};

	let isLoading = false;
	let error = '';
	let success = false;

	async function handleRegister() {
		// Validation
		if (!formData.email || !formData.username || !formData.password) {
			error = 'Please fill in all required fields';
			return;
		}

		if (formData.password !== formData.confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		if (formData.password.length < 8) {
			error = 'Password must be at least 8 characters long';
			return;
		}

		if (formData.username.length < 3) {
			error = 'Username must be at least 3 characters long';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const response = await fetch('/api/v1/auth/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					email: formData.email.trim(),
					username: formData.username.trim(),
					password: formData.password,
					first_name: formData.firstName.trim(),
					last_name: formData.lastName.trim(),
				}),
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Registration failed');
			}

			// Store auth token
			localStorage.setItem('auth_token', data.token);
			localStorage.setItem('user', JSON.stringify(data.user));

			success = true;

			// Redirect to home after a brief delay
			setTimeout(() => {
				goto('/');
			}, 2000);

		} catch (err) {
			error = err.message || 'Registration failed. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	function handleKeyDown(event) {
		if (event.key === 'Enter' && event.ctrlKey) {
			handleRegister();
		}
	}
</script>

<svelte:window on:keydown={handleKeyDown} />

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<div class="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-green-100">
				<svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"></path>
				</svg>
			</div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
				Create your account
			</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				Join TPT Titan - your private, encrypted office suite
			</p>
		</div>

		{#if success}
			<div class="rounded-md bg-green-50 p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-sm font-medium text-green-800">Account created successfully!</h3>
						<p class="mt-2 text-sm text-green-700">Welcome to TPT Titan. Redirecting you to your dashboard...</p>
					</div>
				</div>
			</div>
		{:else}
			<form class="mt-8 space-y-6" on:submit|preventDefault={handleRegister}>
				<div class="space-y-4">
					<!-- Email -->
					<div>
						<label for="email" class="block text-sm font-medium text-gray-700">Email address *</label>
						<input
							id="email"
							name="email"
							type="email"
							autocomplete="email"
							required
							bind:value={formData.email}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
							placeholder="your@email.com"
						/>
					</div>

					<!-- Username -->
					<div>
						<label for="username" class="block text-sm font-medium text-gray-700">Username *</label>
						<input
							id="username"
							name="username"
							type="text"
							autocomplete="username"
							required
							bind:value={formData.username}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
							placeholder="yourusername"
						/>
						<p class="mt-1 text-xs text-gray-500">3-50 characters, letters and numbers only</p>
					</div>

					<!-- Name fields -->
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="firstName" class="block text-sm font-medium text-gray-700">First name</label>
							<input
								id="firstName"
								name="firstName"
								type="text"
								autocomplete="given-name"
								bind:value={formData.firstName}
								class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="John"
							/>
						</div>
						<div>
							<label for="lastName" class="block text-sm font-medium text-gray-700">Last name</label>
							<input
								id="lastName"
								name="lastName"
								type="text"
								autocomplete="family-name"
								bind:value={formData.lastName}
								class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Doe"
							/>
						</div>
					</div>

					<!-- Password -->
					<div>
						<label for="password" class="block text-sm font-medium text-gray-700">Password *</label>
						<input
							id="password"
							name="password"
							type="password"
							autocomplete="new-password"
							required
							bind:value={formData.password}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
							placeholder="••••••••"
						/>
						<p class="mt-1 text-xs text-gray-500">At least 8 characters</p>
					</div>

					<!-- Confirm Password -->
					<div>
						<label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirm password *</label>
						<input
							id="confirmPassword"
							name="confirmPassword"
							type="password"
							autocomplete="new-password"
							required
							bind:value={formData.confirmPassword}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
							placeholder="••••••••"
						/>
					</div>
				</div>

				{#if error}
					<div class="rounded-md bg-red-50 p-4">
						<div class="flex">
							<div class="flex-shrink-0">
								<svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
								</svg>
							</div>
							<div class="ml-3">
								<p class="text-sm font-medium text-red-800">{error}</p>
							</div>
						</div>
					</div>
				{/if}

				<div>
					<button
						type="submit"
						disabled={isLoading}
						class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{#if isLoading}
							<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Creating account...
						{:else}
							Create account
						{/if}
					</button>
				</div>

				<div class="text-center">
					<a href="/auth/login" class="font-medium text-green-600 hover:text-green-500">
						Already have an account? Sign in
					</a>
				</div>

				<div class="text-center text-xs text-gray-500">
					<p>🔐 Your data is encrypted and never shared</p>
					<p>📜 Licensed under AGPL v3.0 or later</p>
				</div>
			</form>
		{/if}
	</div>
</div>
