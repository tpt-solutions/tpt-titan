<script>
	import '../app.css';
	import { page } from '$app/stores';

	// Accept framework-provided props to avoid warnings
	export const data = null;
	export const form = null;
	export const params = null;

	let mobileMenuOpen = false;

	$: currentPath = $page.url.pathname;

	function navClass(path) {
		const active = currentPath === path || (path !== '/' && currentPath.startsWith(path));
		return `px-3 py-2 rounded-md text-sm font-medium transition-colors ${active ? 'text-blue-600 bg-blue-50' : 'text-gray-700 hover:text-blue-600 hover:bg-gray-50'}`;
	}

	const navLinks = [
		{ href: '/',            label: 'Home' },
		{ href: '/editor',      label: 'Editor' },
		{ href: '/spreadsheet', label: 'Spreadsheet' },
		{ href: '/forms',       label: 'Forms' },
		{ href: '/database',    label: 'Database' },
		{ href: '/email',       label: 'Email' },
		{ href: '/calendar',    label: 'Calendar' },
		{ href: '/contacts',    label: 'Contacts' },
		{ href: '/chat',        label: 'Chat' },
		{ href: '/files',       label: 'Files' },
		{ href: '/tasks',        label: 'Tasks' },
		{ href: '/workflows',    label: 'Workflows' },
		{ href: '/speech',       label: 'Speech' },
		{ href: '/voice',        label: 'Voice' },
		{ href: '/math',         label: 'Math' },
		{ href: '/export',       label: 'Export' },
		{ href: '/monitoring',   label: 'Monitoring' },
		{ href: '/admin',        label: 'Admin' },
		{ href: '/plugins',     label: 'Plugins' },
		{ href: '/settings',    label: 'Settings' },
	];
</script>

<main class="min-h-screen bg-gray-50">
	<!-- Navigation Header -->
	<header class="bg-white shadow-sm border-b border-gray-200">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center h-16">
				<!-- Logo -->
				<div class="flex items-center">
					<a href="/" class="text-2xl font-bold text-blue-600 hover:text-blue-700 transition-colors">
						TPT Titan
					</a>
				</div>

				<!-- Desktop Navigation -->
				<nav class="hidden md:flex flex-wrap gap-1">
					{#each navLinks as link}
						<a href={link.href} class={navClass(link.href)}>
							{link.label}
						</a>
					{/each}
				</nav>

				<!-- Mobile menu button -->
				<div class="md:hidden">
					<button
						on:click={() => mobileMenuOpen = !mobileMenuOpen}
						aria-label="Toggle navigation menu"
						aria-expanded={mobileMenuOpen}
						class="text-gray-700 hover:text-blue-600 p-2 rounded-md"
					>
						{#if mobileMenuOpen}
							<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
							</svg>
						{:else}
							<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
							</svg>
						{/if}
					</button>
				</div>
			</div>
		</div>

		<!-- Mobile dropdown menu -->
		{#if mobileMenuOpen}
			<div class="md:hidden border-t border-gray-200 bg-white px-4 py-3 space-y-1">
				{#each navLinks as link}
					<a
						href={link.href}
						on:click={() => mobileMenuOpen = false}
						class="block {navClass(link.href)}"
					>
						{link.label}
					</a>
				{/each}
			</div>
		{/if}
	</header>

	<!-- Main Content -->
	<slot />
</main>

<style>
	/* Global styles can go here */
</style>
