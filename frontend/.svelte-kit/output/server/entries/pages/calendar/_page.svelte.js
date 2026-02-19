import { c as create_ssr_component, n as calendars, p as events, q as currentView, t as currentDate, v as validate_component, C as CalendarView } from "../../../chunks/calendar.js";
import "@sveltejs/kit/internal";
import "../../../chunks/svelte-kit.js";
import "@sveltejs/kit/internal/server";
import { E as EventForm } from "../../../chunks/forms.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { params = null } = $$props;
  let { data = null } = $$props;
  let { form = null } = $$props;
  let showEventForm = false;
  let editingEvent = null;
  let selectedDate = null;
  let calendarList = [];
  let eventList = [];
  let view = "month";
  let date = /* @__PURE__ */ new Date();
  calendars.subscribe((value) => calendarList = value);
  events.subscribe((value) => eventList = value);
  currentView.subscribe((value) => view = value);
  currentDate.subscribe((value) => date = value);
  async function loadEvents() {
    const startDate = new Date(date);
    startDate.setDate(1);
    startDate.setMonth(startDate.getMonth() - 1);
    const endDate = new Date(date);
    endDate.setMonth(endDate.getMonth() + 2);
    try {
      const response = await fetch(`/api/v1/events?start=${startDate.toISOString()}&end=${endDate.toISOString()}`, {
        headers: {
          "Authorization": `Bearer ${localStorage.getItem("token")}`
        }
      });
      if (response.ok) {
        const data2 = await response.json();
        events.set(data2.events || []);
      }
    } catch (error) {
      console.error("Failed to load events:", error);
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
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  return `${$$result.head += `<!-- HEAD_svelte-1qavwmb_START -->${$$result.title = `<title>Calendar - TPT Titan</title>`, ""}<!-- HEAD_svelte-1qavwmb_END -->`, ""} <div class="container mx-auto px-4 py-8"><div class="flex justify-between items-center mb-8"><h1 class="text-3xl font-bold text-gray-900 dark:text-white" data-svelte-h="svelte-3a3xx5">Calendar</h1> <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2" data-svelte-h="svelte-1b56j32"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
			New Event</button></div>  ${validate_component(CalendarView, "CalendarView").$$render(
    $$result,
    {
      calendarList,
      eventList,
      view,
      date,
      handleCreateEvent,
      handleEditEvent,
      handleViewChange,
      handleDateChange
    },
    {},
    {}
  )}  ${showEventForm ? `${validate_component(EventForm, "EventForm").$$render(
    $$result,
    {
      event: editingEvent,
      selectedDate,
      calendars: calendarList,
      onClose: handleFormClose
    },
    {},
    {}
  )}` : ``}</div>`;
});
export {
  Page as default
};
