<script>
	import '../app.css';
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

				<!-- Navigation -->
				<nav class="hidden md:flex space-x-8">
					<a href="/" class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
						Home
					</a>
					<a href="/spreadsheet" class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
						Spreadsheet
					</a>
					<a href="/forms" class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
						Forms
					</a>
					<a href="/editor" class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
						Text Editor
					</a>
					<a href="/contacts" class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
						Contacts
					</a>
					<a href="/calendar" class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
						Calendar
					</a>
					<a href="/email" class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
						Email
					</a>
					<a href="/tasks" class="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
						Tasks
					</a>
				</nav>

				<!-- Mobile menu button -->
				<div class="md:hidden">
					<button class="text-gray-700 hover:text-blue-600 p-2">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
						</svg>
					</button>
				</div>
			</div>
		</div>
	</header>

	<!-- Main Content -->
	<slot />
</main>

<style>
	/* Global styles can go here */
</style>
