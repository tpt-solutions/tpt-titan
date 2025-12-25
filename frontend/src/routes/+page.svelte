<script>
	// Page script
</script>

<svelte:head>
	<title>TPT Titan - Open Source Office Suite</title>
	<meta name="description" content="Complete open source alternative to Microsoft Office 365" />
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<header class="text-center mb-12">
		<h1 class="text-5xl font-bold text-gray-900 mb-4">
			Welcome to <span class="text-blue-600">TPT Titan</span>
		</h1>
		<p class="text-xl text-gray-600 max-w-2xl mx-auto">
			A complete open source alternative to Microsoft Office 365.
			Built with modern web technologies for productivity, collaboration, and privacy.
		</p>
	</header>

	<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8 mb-12">
		<div class="bg-white p-6 rounded-lg shadow-md">
			<div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mb-4">
				<svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
				</svg>
			</div>
			<h3 class="text-xl font-semibold mb-2">Text Editor</h3>
			<p class="text-gray-600">Rich text editing with collaboration features and export capabilities.</p>
		</div>

		<div class="bg-white p-6 rounded-lg shadow-md">
			<div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mb-4">
				<svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
				</svg>
			</div>
			<h3 class="text-xl font-semibold mb-2">Spreadsheet</h3>
			<p class="text-gray-600">Powerful spreadsheet application with formulas, charts, and data analysis.</p>
		</div>

		<div class="bg-white p-6 rounded-lg shadow-md">
			<div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mb-4">
				<svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m-9 0h10a2 2 0 012 2v10a2 2 0 01-2 2H6a2 2 0 01-2-2V6a2 2 0 012-2z"></path>
				</svg>
			</div>
			<h3 class="text-xl font-semibold mb-2">Presentations</h3>
			<p class="text-gray-600">Create stunning presentations with animations and collaborative editing.</p>
		</div>

		<div class="bg-white p-6 rounded-lg shadow-md">
			<div class="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center mb-4">
				<svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
				</svg>
			</div>
			<h3 class="text-xl font-semibold mb-2">Email Client</h3>
			<p class="text-gray-600">Full-featured email client with encryption and organization tools.</p>
		</div>

		<div class="bg-white p-6 rounded-lg shadow-md">
			<div class="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center mb-4">
				<svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3a1 1 0 011-1h6a1 1 0 011 1v4m-9 0h10a2 2 0 012 2v10a2 2 0 01-2 2H7a2 2 0 01-2-2V9a2 2 0 012-2z"></path>
				</svg>
			</div>
			<h3 class="text-xl font-semibold mb-2">Calendar</h3>
			<p class="text-gray-600">Comprehensive calendar with scheduling, reminders, and sharing.</p>
		</div>

		<div class="bg-white p-6 rounded-lg shadow-md">
			<div class="w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center mb-4">
				<svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
				</svg>
			</div>
			<h3 class="text-xl font-semibold mb-2">Contacts</h3>
			<p class="text-gray-600">Contact management with groups, search, and integration.</p>
		</div>
	</div>

	<div class="text-center">
		<div class="bg-blue-50 border border-blue-200 rounded-lg p-6 mb-8">
			<h2 class="text-2xl font-bold text-blue-900 mb-4">🚀 Getting Started</h2>
			<p class="text-blue-800 mb-4">
				TPT Titan is currently in active development. We're building a complete, privacy-focused office suite.
			</p>
			<div class="flex justify-center space-x-4">
				<a href="/docs" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors">
					View Documentation
				</a>
				<a href="https://github.com/tpt-titan/tpt-titan" class="bg-gray-800 text-white px-6 py-2 rounded-lg hover:bg-gray-900 transition-colors">
					GitHub Repository
				</a>
			</div>
		</div>
	</div>
</div>
