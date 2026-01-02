<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { calendars, events, currentView, currentDate } from '$lib/stores';
	import CalendarView from '$lib/components/CalendarView.svelte';
	import EventForm from '$lib/components/EventForm.svelte';

	let showEventForm = false;
	let editingEvent = null;
	let selectedDate = null;

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
		try {
			const response = await fetch('/api/v1/calendars', {
				headers: {
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				}
			});

			if (response.ok) {
				const data = await response.json();
				calendars.set(data.calendars || []);
			} else if (response.status === 401) {
				goto('/auth/login');
			}
		} catch (error) {
			console.error('Failed to load calendars:', error);
		}
	}

	async function loadEvents() {
		const startDate = new Date(date);
		startDate.setDate(1);
		startDate.setMonth(startDate.getMonth() - 1); // Start from previous month

		const endDate = new Date(date);
		endDate.setMonth(endDate.getMonth() + 2); // End at next month

		try {
			const response = await fetch(`/api/v1/events?start=${startDate.toISOString()}&end=${endDate.toISOString()}`, {
				headers: {
					'Authorization': `Bearer ${localStorage.getItem('token')}`
				}
			});

			if (response.ok) {
				const data = await response.json();
				events.set(data.events || []);
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
</div>
