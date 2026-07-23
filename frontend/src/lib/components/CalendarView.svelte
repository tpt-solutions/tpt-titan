<script>
// @ts-nocheck
	import { onMount } from 'svelte';
	import { currentView, currentDate } from '$lib/stores';

	export let calendarList = [];
	export let eventList = [];
	export let view = 'month';
	export let date = new Date();
	export let handleCreateEvent;
	export let handleEditEvent;
	export let handleViewChange;
	export let handleDateChange;

	let currentDateObj = new Date(date);

	// Update local date when store changes
	$: currentDateObj = new Date(date);

	function getMonthData(displayDate) {
		const year = displayDate.getFullYear();
		const month = displayDate.getMonth();

		// Get first day of the month
		const firstDay = new Date(year, month, 1);
		// Get last day of the month
		const lastDay = new Date(year, month + 1, 0);

		// Get the day of week for first day (0 = Sunday, 1 = Monday, etc.)
		let startDay = firstDay.getDay();
		// Adjust for Monday start (optional - comment out if you want Sunday start)
		startDay = startDay === 0 ? 6 : startDay - 1;

		const weeks = [];
		let currentWeek = [];
		let dayCounter = 1 - startDay; // Start from previous month if needed

		// Create 6 weeks (maximum needed for any month)
		for (let week = 0; week < 6; week++) {
			currentWeek = [];

			for (let day = 0; day < 7; day++) {
				const currentDate = new Date(year, month, dayCounter);
				const isCurrentMonth = currentDate.getMonth() === month;

				// Get events for this day
				const dayEvents = eventList.filter(event => {
					const eventDate = new Date(event.start_time);
					return eventDate.toDateString() === currentDate.toDateString();
				});

				currentWeek.push({
					date: currentDate,
					dayNumber: dayCounter,
					isCurrentMonth,
					events: dayEvents,
					isToday: currentDate.toDateString() === new Date().toDateString()
				});

				dayCounter++;
			}

			weeks.push(currentWeek);

			// Stop if we've covered all days of the current month
			if (dayCounter > lastDay.getDate() && week >= 4) break;
		}

		return weeks;
	}

	function navigateMonth(direction) {
		const newDate = new Date(currentDateObj);
		newDate.setMonth(newDate.getMonth() + direction);
		handleDateChange(newDate);
	}

	function goToToday() {
		handleDateChange(new Date());
	}

	function formatMonthYear(date) {
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long' });
	}

	function handleDayClick(dayData) {
		if (dayData.isCurrentMonth) {
			handleCreateEvent(dayData.date);
		}
	}

	function handleEventClick(event) {
		handleEditEvent(event);
	}

	$: monthData = getMonthData(currentDateObj);
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow">
	<!-- Calendar Header -->
	<div class="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
		<div class="flex items-center space-x-4">
			<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
				{formatMonthYear(currentDateObj)}
			</h2>
			<button
				on:click={goToToday}
				class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
			>
				Today
			</button>
		</div>

		<div class="flex items-center space-x-2">
			<!-- View Toggle (Month only for now) -->
			<div class="flex bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
				<button
					class="px-3 py-1 text-sm rounded-md {view === 'month' ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-white shadow-sm' : 'text-gray-600 dark:text-gray-300'}"
					on:click={() => handleViewChange('month')}
				>
					Month
				</button>
			</div>

			<!-- Navigation -->
			<div class="flex space-x-1" role="group" aria-label="Month navigation">
				<button
					on:click={() => navigateMonth(-1)}
					on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); navigateMonth(-1); } }}
					class="p-2 text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700"
					aria-label="Previous month"
					type="button"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
					</svg>
				</button>
				<button
					on:click={() => navigateMonth(1)}
					on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); navigateMonth(1); } }}
					class="p-2 text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700"
					aria-label="Next month"
					type="button"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
					</svg>
				</button>
			</div>
		</div>
	</div>

	<!-- Calendar Grid -->
	<div class="p-4">
		<!-- Day Headers -->
		<div class="grid grid-cols-7 gap-1 mb-2">
			{#each ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] as dayName}
				<div class="p-2 text-center text-sm font-medium text-gray-600 dark:text-gray-400">
					{dayName}
				</div>
			{/each}
		</div>

		<!-- Month Grid -->
		<div class="grid grid-cols-7 gap-1">
			{#each monthData as week}
				{#each week as dayData}
					<button
						class="min-h-[120px] p-2 border border-gray-200 dark:border-gray-600 rounded-md cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors {dayData.isCurrentMonth ? 'bg-white dark:bg-gray-800' : 'bg-gray-50 dark:bg-gray-700'} {dayData.isToday ? 'ring-2 ring-blue-500' : ''} w-full text-left"
						on:click={() => handleDayClick(dayData)}
						on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); handleDayClick(dayData); } }}
						aria-label="{dayData.isCurrentMonth ? formatMonthYear(dayData.date).split(' ')[0] + ' ' + dayData.date.getDate() + ', ' + dayData.date.getFullYear() + (dayData.events.length > 0 ? ', ' + dayData.events.length + ' event' + (dayData.events.length !== 1 ? 's' : '') : ', no events') : 'Not in current month'}"
						type="button"
					>
						<div class="flex justify-between items-start mb-1">
							<span class="text-sm font-medium {dayData.isCurrentMonth ? 'text-gray-900 dark:text-white' : 'text-gray-400 dark:text-gray-500'}">
								{dayData.date.getDate()}
							</span>
							{#if dayData.isToday}
								<span class="text-xs bg-blue-500 text-white px-1 py-0.5 rounded">
									Today
								</span>
							{/if}
						</div>

						<!-- Events for this day -->
						<div class="space-y-1 max-h-[80px] overflow-hidden">
							{#each dayData.events.slice(0, 3) as event}
								<div
									class="text-xs p-1 rounded truncate cursor-pointer hover:opacity-80 transition-opacity"
									style="background-color: {event.calendar_color}20; border-left: 3px solid {event.calendar_color};"
									on:click|stopPropagation={() => handleEventClick(event)}
									title="{event.title} - {new Date(event.start_time).toLocaleTimeString()}"
								>
									<div class="font-medium text-gray-900 dark:text-white truncate">
										{event.title}
									</div>
									{#if !event.is_all_day}
										<div class="text-gray-600 dark:text-gray-400">
											{new Date(event.start_time).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}
										</div>
									{/if}
								</div>
							{/each}
							{#if dayData.events.length > 3}
								<div class="text-xs text-gray-500 dark:text-gray-400 pl-1">
									+{dayData.events.length - 3} more
								</div>
							{/if}
						</div>
					</button>
				{/each}
			{/each}
		</div>
	</div>

	<!-- Calendar Legend -->
	{#if calendarList.length > 0}
		<div class="border-t border-gray-200 dark:border-gray-700 p-4">
			<div class="flex flex-wrap gap-4">
				{#each calendarList as calendar}
					<div class="flex items-center space-x-2">
						<div
							class="w-4 h-4 rounded"
							style="background-color: {calendar.color}"
						></div>
						<span class="text-sm text-gray-700 dark:text-gray-300">
							{calendar.name}
						</span>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>
