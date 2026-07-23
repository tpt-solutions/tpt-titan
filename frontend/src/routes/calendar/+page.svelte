<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { calendars, events, currentView, currentDate } from '$lib/stores';
	import CalendarView from '$lib/components/CalendarView.svelte';
	import EventForm from '$lib/components/EventForm.svelte';

	// Accept framework-provided props to avoid warnings
	export let data = null;
	export let form = null;
	export let params = null;


	let showEventForm = false;
	let editingEvent = null;
	let selectedDate = null;
	let isLoading = false;
	let loadError = null;

	let calendarList = [];
	let eventList = [];
	let view = 'month';
	let date = new Date();

	calendars.subscribe(value => calendarList = value);
	events.subscribe(value => eventList = value);
	currentView.subscribe(value => view = value);
	currentDate.subscribe(value => date = value);

	onMount(async () => {
		await loadCalendars();
		await loadEvents();
	});

	async function loadCalendars() {
		isLoading = true;
		loadError = null;
		try {
			const response = await fetch('/api/v1/calendars', {
				headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
			});
			if (response.ok) {
				const data = await response.json();
				calendars.set(data.calendars || []);
			} else if (response.status === 401) {
				goto('/auth/login');
			} else {
				loadError = `Failed to load calendars (${response.status})`;
			}
		} catch (error) {
			loadError = 'Could not connect to server. Please check your connection.';
			console.error('Failed to load calendars:', error);
		} finally {
			isLoading = false;
		}
	}

	async function loadEvents() {
		const startDate = new Date(date);
		startDate.setDate(1);
		startDate.setMonth(startDate.getMonth() - 1);
		const endDate = new Date(date);
		endDate.setMonth(endDate.getMonth() + 2);

		try {
			const response = await fetch(
				`/api/v1/events?start=${startDate.toISOString()}&end=${endDate.toISOString()}`,
				{ headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` } }
			);
			if (response.ok) {
				const data = await response.json();
				events.set(data.events || []);
			} else if (response.status === 401) {
				goto('/auth/login');
			}
		} catch (error) {
			console.error('Failed to load events:', error);
		}
	}

	function handleCreateEvent(selectedDateTime = null) {
		editingEvent = null;
		selectedDate = selectedDateTime;
		showEventForm = true;
	}

	function handleEditEvent(event) {
		editingEvent = event;
		selectedDate = null;
		showEventForm = true;
	}

	function handleFormClose() {
		showEventForm = false;
		editingEvent = null;
		selectedDate = null;
		loadEvents();
	}

	function handleViewChange(newView) {
		currentView.set(newView);
	}

	function handleDateChange(newDate) {
		currentDate.set(newDate);
		date = newDate;
		loadEvents();
	}
</script>

<svelte:head>
	<title>Calendar - TPT Titan</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="flex justify-between items-center mb-8">
		<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Calendar</h1>
		<button
			on:click={() => handleCreateEvent()}
			class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
			</svg>
			New Event
		</button>
	</div>

	<!-- Error banner -->
	{#if loadError}
		<div class="mb-6 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg px-4 py-3 flex items-center gap-3">
			<svg class="w-5 h-5 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
			</svg>
			<span class="text-sm text-red-700 dark:text-red-300">{loadError}</span>
			<button on:click={loadCalendars} class="ml-auto text-sm text-red-600 dark:text-red-400 underline hover:no-underline">Retry</button>
		</div>
	{/if}

	<!-- Loading state -->
	{#if isLoading}
		<div class="flex items-center justify-center py-16 text-gray-400 dark:text-gray-500">
			<svg class="animate-spin w-8 h-8 mr-3" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
			</svg>
			<span>Loading calendar…</span>
		</div>
	{:else}

	<!-- Calendar View -->
	<CalendarView
		{calendarList}
		{eventList}
		{view}
		{date}
		{handleCreateEvent}
		{handleEditEvent}
		{handleViewChange}
		{handleDateChange}
	/>

	<!-- Event Form Modal -->
	{#if showEventForm}
		<EventForm
			event={editingEvent}
			selectedDate={selectedDate}
			calendars={calendarList}
			onClose={handleFormClose}
		/>
	{/if}

	{/if}
</div>
