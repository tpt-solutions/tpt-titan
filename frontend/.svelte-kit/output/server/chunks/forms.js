import { c as create_ssr_component, f as createEventDispatcher, g as contacts, b as escape, d as each, a as add_attribute, s as subscribe } from "./calendar.js";
import { f as formulaBar, s as selectedCell, a as selectedCells, i as isFormulaRangeSelecting } from "./spreadsheet.js";
const EventForm = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { event = null } = $$props;
  let { selectedDate = null } = $$props;
  let { calendars = [] } = $$props;
  let { onClose } = $$props;
  createEventDispatcher();
  let formData = {
    calendar_id: "",
    title: "",
    description: "",
    location: "",
    start_time: "",
    end_time: "",
    is_all_day: false,
    attendee_ids: []
  };
  let errors = {};
  let availableContacts = [];
  let availableSTTModels = [];
  let schedulingSuggestions = [];
  contacts.subscribe((value) => availableContacts = value);
  if (event) {
    formData = {
      calendar_id: event.calendar_id,
      title: event.title || "",
      description: event.description || "",
      location: event.location || "",
      start_time: event.start_time ? new Date(event.start_time).toISOString().slice(0, 16) : "",
      end_time: event.end_time ? new Date(event.end_time).toISOString().slice(0, 16) : "",
      is_all_day: event.is_all_day || false,
      reminder_minutes: event.reminder_minutes || 15,
      attendee_ids: event.attendees ? event.attendees.map((a) => a.contact_id) : []
    };
  } else {
    const defaultCalendar = calendars.find((c) => c.is_default) || calendars[0];
    const now = selectedDate ? new Date(selectedDate) : /* @__PURE__ */ new Date();
    const startTime = new Date(now);
    startTime.setHours(startTime.getHours() + 1, 0, 0, 0);
    const endTime = new Date(startTime);
    endTime.setHours(endTime.getHours() + 1);
    formData = {
      calendar_id: defaultCalendar ? defaultCalendar.id : "",
      title: "",
      description: "",
      location: "",
      start_time: startTime.toISOString().slice(0, 16),
      end_time: endTime.toISOString().slice(0, 16),
      is_all_day: false,
      reminder_minutes: 15,
      attendee_ids: []
    };
  }
  if ($$props.event === void 0 && $$bindings.event && event !== void 0) $$bindings.event(event);
  if ($$props.selectedDate === void 0 && $$bindings.selectedDate && selectedDate !== void 0) $$bindings.selectedDate(selectedDate);
  if ($$props.calendars === void 0 && $$bindings.calendars && calendars !== void 0) $$bindings.calendars(calendars);
  if ($$props.onClose === void 0 && $$bindings.onClose && onClose !== void 0) $$bindings.onClose(onClose);
  calendars.find((c) => c.id === formData.calendar_id);
  return `  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"> <div class="relative top-20 mx-auto p-5 border w-full max-w-2xl shadow-lg rounded-md bg-white dark:bg-gray-800"><div class="mt-3 max-h-[80vh] overflow-y-auto"> <div class="flex items-center justify-between mb-4"><h3 class="text-lg font-medium text-gray-900 dark:text-white">${escape(event ? "Edit Event" : "Create New Event")}</h3> <button class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" data-svelte-h="svelte-jql23q"><svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg></button></div>  <form class="space-y-4"> ${errors.general ? `<div class="bg-red-50 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-600 dark:text-red-200 px-4 py-3 rounded">${escape(errors.general)}</div>` : ``}  <div><label for="calendar_id" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-1rv34no">Calendar</label> <select id="calendar_id" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"><option value="" data-svelte-h="svelte-lkzcwb">Select a calendar</option>${each(calendars, (calendar) => {
    return `<option${add_attribute("value", calendar.id, 0)}>${escape(calendar.name)}</option>`;
  })}</select> ${errors.calendar_id ? `<p class="mt-1 text-sm text-red-600 dark:text-red-400">${escape(errors.calendar_id)}</p>` : ``}</div>  <div><label for="title" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-gweqb0">Event Title *</label> <input id="title" type="text" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="Meeting with team"${add_attribute("value", formData.title, 0)}> ${errors.title ? `<p class="mt-1 text-sm text-red-600 dark:text-red-400">${escape(errors.title)}</p>` : ``}</div>  <div><label for="description" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-14jq8xi">Description</label> <textarea id="description" rows="3" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="Event details...">${escape(formData.description || "")}</textarea></div>  <div><label for="location" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-1n81a8g">Location</label> <input id="location" type="text" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="Conference Room A"${add_attribute("value", formData.location, 0)}></div>  <div class="flex items-center"><input id="is_all_day" type="checkbox" class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"${add_attribute("checked", formData.is_all_day, 1)}> <label for="is_all_day" class="ml-2 block text-sm text-gray-900 dark:text-white" data-svelte-h="svelte-11jihi5">All day event</label></div>  <div class="grid grid-cols-1 md:grid-cols-2 gap-4"><div><label for="start_time" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-qeumll">Start Time *</label> <input id="start_time" type="datetime-local" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"${add_attribute("value", formData.start_time, 0)}> ${errors.start_time ? `<p class="mt-1 text-sm text-red-600 dark:text-red-400">${escape(errors.start_time)}</p>` : ``}</div> <div><label for="end_time" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-1vetqtj">End Time *</label> <input id="end_time" type="datetime-local" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"${add_attribute("value", formData.end_time, 0)}> ${errors.end_time ? `<p class="mt-1 text-sm text-red-600 dark:text-red-400">${escape(errors.end_time)}</p>` : ``}</div></div>  <div><label for="reminder_minutes" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-1uk10qu">Reminder</label> <select id="reminder_minutes" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"><option${add_attribute("value", 0, 0)} data-svelte-h="svelte-1hlrlix">No reminder</option><option${add_attribute("value", 5, 0)} data-svelte-h="svelte-cypo54">5 minutes before</option><option${add_attribute("value", 15, 0)} data-svelte-h="svelte-3z5pyw">15 minutes before</option><option${add_attribute("value", 30, 0)} data-svelte-h="svelte-1uu1k4o">30 minutes before</option><option${add_attribute("value", 60, 0)} data-svelte-h="svelte-1c1ua1k">1 hour before</option><option${add_attribute("value", 1440, 0)} data-svelte-h="svelte-v8w9o1">1 day before</option></select></div>  ${availableContacts.length > 0 ? `<div><label id="attendees-label" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2" data-svelte-h="svelte-js8uf9">Attendees</label> <div role="group" aria-labelledby="attendees-label" class="max-h-40 overflow-y-auto border border-gray-300 dark:border-gray-600 rounded-md p-2 space-y-2">${each(availableContacts, (contact) => {
    return `<label class="flex items-center space-x-2 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 p-1 rounded"><input type="checkbox" ${formData.attendee_ids.includes(contact.id) ? "checked" : ""} class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded" aria-label="${"Add " + escape(
      contact.first_name || contact.last_name ? `${contact.first_name || ""} ${contact.last_name || ""}`.trim() : "Unnamed Contact",
      true
    ) + " as attendee"}"> <div class="flex-1 min-w-0"><div class="text-sm font-medium text-gray-900 dark:text-white truncate">${escape(contact.first_name || contact.last_name ? `${contact.first_name || ""} ${contact.last_name || ""}`.trim() : "Unnamed Contact")}</div> ${contact.email ? `<div class="text-xs text-gray-500 dark:text-gray-400 truncate">${escape(contact.email)} </div>` : ``}</div> </label>`;
  })}</div></div>` : ``}  ${!event ? `<div class="border-t border-gray-200 dark:border-gray-600 pt-4"><h4 class="text-sm font-medium text-gray-900 dark:text-white mb-3" data-svelte-h="svelte-1uy9u6l">🤖 AI Assistant</h4> <div class="grid grid-cols-1 md:grid-cols-2 gap-3">${availableSTTModels.length > 0 ? `<button type="button" ${""}${add_attribute(
    "aria-label",
    "Start voice event creation",
    0
  )} class="flex items-center justify-center px-4 py-3 border border-purple-300 dark:border-purple-600 rounded-md hover:bg-purple-50 dark:hover:bg-purple-900 transition-colors disabled:opacity-50"><div class="flex items-center space-x-2"><svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg> <span class="text-sm font-medium text-purple-700 dark:text-purple-300">${escape("Voice Create Event")}</span></div></button>` : ``} <button type="button" ${""}${add_attribute(
    "aria-label",
    "Get smart scheduling suggestions",
    0
  )} class="flex items-center justify-center px-4 py-3 border border-blue-300 dark:border-blue-600 rounded-md hover:bg-blue-50 dark:hover:bg-blue-900 transition-colors disabled:opacity-50"><div class="flex items-center space-x-2"><svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg> <span class="text-sm font-medium text-blue-700 dark:text-blue-300">${escape("Smart Scheduling")}</span></div></button></div>  ${``}  ${schedulingSuggestions.length > 0 ? `<div class="mt-3 space-y-2"><h5 class="text-sm font-medium text-gray-900 dark:text-white" data-svelte-h="svelte-a9b01h">Suggested Times:</h5> ${each(schedulingSuggestions, (suggestion) => {
    return `<div class="flex items-center justify-between p-2 bg-blue-50 dark:bg-blue-900 border border-blue-200 dark:border-blue-700 rounded"><div><div class="text-sm font-medium text-blue-900 dark:text-blue-100">${escape(new Date(suggestion.start_time).toLocaleDateString())} at ${escape(new Date(suggestion.start_time).toLocaleTimeString())}</div> <div class="text-xs text-blue-700 dark:text-blue-300">${escape(suggestion.duration)} minutes • ${escape(suggestion.reason)} </div></div> <button type="button" class="px-3 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-1gd84mn">Use This</button> </div>`;
  })}</div>` : ``}</div>` : ``}  <div class="flex justify-end space-x-3 pt-4"><button type="button" class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2" ${""}>Cancel</button> <button type="submit" ${""} class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed">${`${escape(event ? "Update Event" : "Create Event")}`}</button></div></form>  ${``}</div></div></div>`;
});
const ContactForm = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { contact = null } = $$props;
  let { onClose } = $$props;
  createEventDispatcher();
  let formData = {
    first_name: "",
    last_name: "",
    emails: [{ type: "work", value: "" }],
    phones: [{ type: "mobile", value: "" }],
    company: "",
    position: "",
    notes: ""
  };
  let errors = {};
  if (contact) {
    formData = {
      first_name: contact.first_name || "",
      last_name: contact.last_name || "",
      emails: contact.emails && contact.emails.length > 0 ? contact.emails : [{ type: "work", value: "" }],
      phones: contact.phones && contact.phones.length > 0 ? contact.phones : [{ type: "mobile", value: "" }],
      company: contact.company || "",
      position: contact.position || "",
      notes: contact.notes || ""
    };
  }
  if ($$props.contact === void 0 && $$bindings.contact && contact !== void 0) $$bindings.contact(contact);
  if ($$props.onClose === void 0 && $$bindings.onClose && onClose !== void 0) $$bindings.onClose(onClose);
  return `  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" role="presentation"> <div class="relative top-20 mx-auto p-5 border w-full max-w-lg shadow-lg rounded-md bg-white dark:bg-gray-800" role="dialog" aria-modal="true" aria-labelledby="contact-form-title"><div class="mt-3"> <div class="flex items-center justify-between mb-4"><h3 id="contact-form-title" class="text-lg font-medium text-gray-900 dark:text-white">${escape(contact ? "Edit Contact" : "Add New Contact")}</h3> <button class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" data-svelte-h="svelte-jql23q"><svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg></button></div>  <form class="space-y-4"> ${errors.general ? `<div class="bg-red-50 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-600 dark:text-red-200 px-4 py-3 rounded">${escape(errors.general)}</div>` : ``}  <div class="grid grid-cols-2 gap-4"><div><label for="first_name" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-1x4a1fb">First Name</label> <input id="first_name" type="text" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="John"${add_attribute("value", formData.first_name, 0)}></div> <div><label for="last_name" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-149cc5r">Last Name</label> <input id="last_name" type="text" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="Doe"${add_attribute("value", formData.last_name, 0)}></div></div> ${errors.name ? `<p class="text-sm text-red-600 dark:text-red-400">${escape(errors.name)}</p>` : ``}  <div><div class="flex items-center justify-between mb-2"><label id="email-label" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-1wsmn6t">Email Addresses</label> <button type="button" class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300" aria-label="Add another email address" data-svelte-h="svelte-xkbhh8">+ Add Email</button></div> <div role="group" aria-labelledby="email-label">${each(formData.emails, (email, index) => {
    return `<div class="flex gap-2 mb-2"><label for="${"email-type-" + escape(index, true)}" class="sr-only">Email ${escape(index + 1)} type</label> <select id="${"email-type-" + escape(index, true)}" class="flex-shrink-0 w-20 px-2 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500 text-sm" aria-label="${"Type for email " + escape(index + 1, true)}"><option value="work" data-svelte-h="svelte-1llns8i">Work</option><option value="personal" data-svelte-h="svelte-45abke">Personal</option><option value="other" data-svelte-h="svelte-902jce">Other</option></select> <label for="${"email-value-" + escape(index, true)}" class="sr-only">Email ${escape(index + 1)} address</label> <input id="${"email-value-" + escape(index, true)}" type="email" class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="john.doe@example.com" aria-label="${"Email " + escape(index + 1, true) + " address"}"${add_attribute("value", email.value, 0)}> ${formData.emails.length > 1 ? `<button type="button" class="flex-shrink-0 px-2 py-2 text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300" title="Remove email" aria-label="${"Remove email " + escape(index + 1, true)}" data-svelte-h="svelte-1pwocsw">✕
									</button>` : ``}</div> ${errors[`email_${index}`] ? `<p class="text-sm text-red-600 dark:text-red-400 mb-2">${escape(errors[`email_${index}`])}</p>` : ``}`;
  })}</div></div>  <div><div class="flex items-center justify-between mb-2"><label id="phone-label" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-93w69v">Phone Numbers</label> <button type="button" class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300" aria-label="Add another phone number" data-svelte-h="svelte-mg6qmf">+ Add Phone</button></div> <div role="group" aria-labelledby="phone-label">${each(formData.phones, (phone, index) => {
    return `<div class="flex gap-2 mb-2"><label for="${"phone-type-" + escape(index, true)}" class="sr-only">Phone ${escape(index + 1)} type</label> <select id="${"phone-type-" + escape(index, true)}" class="flex-shrink-0 w-20 px-2 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500 text-sm" aria-label="${"Type for phone " + escape(index + 1, true)}"><option value="mobile" data-svelte-h="svelte-87cmp6">Mobile</option><option value="work" data-svelte-h="svelte-1llns8i">Work</option><option value="home" data-svelte-h="svelte-1xhv3dq">Home</option><option value="other" data-svelte-h="svelte-902jce">Other</option></select> <label for="${"phone-value-" + escape(index, true)}" class="sr-only">Phone ${escape(index + 1)} number</label> <input id="${"phone-value-" + escape(index, true)}" type="tel" class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="+1 (555) 123-4567" aria-label="${"Phone " + escape(index + 1, true) + " number"}"${add_attribute("value", phone.value, 0)}> ${formData.phones.length > 1 ? `<button type="button" class="flex-shrink-0 px-2 py-2 text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300" title="Remove phone" aria-label="${"Remove phone " + escape(index + 1, true)}" data-svelte-h="svelte-1ymtm20">✕
									</button>` : ``} </div>`;
  })}</div></div>  <div class="grid grid-cols-2 gap-4"><div><label for="company" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-dlxxys">Company</label> <input id="company" type="text" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="Acme Corp"${add_attribute("value", formData.company, 0)}></div> <div><label for="position" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-ui71aq">Position</label> <input id="position" type="text" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="Software Engineer"${add_attribute("value", formData.position, 0)}></div></div>  <div><label for="notes" class="block text-sm font-medium text-gray-700 dark:text-gray-300" data-svelte-h="svelte-i0n5lc">Notes</label> <textarea id="notes" rows="3" class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500" placeholder="Additional notes about this contact...">${escape(formData.notes || "")}</textarea></div>  <div class="flex justify-end space-x-3 pt-4"><button type="button" class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2" ${""}>Cancel</button> <button type="submit" ${""} class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed">${`${escape(contact ? "Update Contact" : "Create Contact")}`}</button></div></form></div></div></div>`;
});
function getStatusColor(status) {
  switch (status) {
    case "active":
      return "bg-green-100 text-green-800";
    case "draft":
      return "bg-yellow-100 text-yellow-800";
    case "archived":
      return "bg-gray-100 text-gray-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
}
const FormList = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { forms = [] } = $$props;
  createEventDispatcher();
  if ($$props.forms === void 0 && $$bindings.forms && forms !== void 0) $$bindings.forms(forms);
  return `<div class="p-6"><div class="mb-8"><div class="flex items-center justify-between mb-4"><div data-svelte-h="svelte-7fo0d8"><h2 class="text-2xl font-bold text-gray-900 mb-2">📋 Advanced Form Management</h2> <p class="text-gray-600">MS Access-style database features with relationships, workflows, and advanced reporting</p></div> <div class="flex space-x-3"><button class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors flex items-center space-x-2" data-svelte-h="svelte-ordjzq"><span>🔗</span> <span>Database Relations</span></button> <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-2" data-svelte-h="svelte-1ksr4br"><span>📊</span> <span>Advanced Reports</span></button> <button class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors flex items-center space-x-2" data-svelte-h="svelte-1jbs3o7"><span>⚡</span> <span>Workflow Designer</span></button></div></div>  <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6" data-svelte-h="svelte-7iqe4v"><div class="bg-blue-50 p-4 rounded-lg"><div class="flex items-center space-x-2 mb-2"><span class="text-blue-600">🔗</span> <h3 class="font-medium text-blue-900">Relationships</h3></div> <p class="text-sm text-blue-700">One-to-one, one-to-many, many-to-many database relationships</p></div> <div class="bg-green-50 p-4 rounded-lg"><div class="flex items-center space-x-2 mb-2"><span class="text-green-600">📊</span> <h3 class="font-medium text-green-900">Advanced Queries</h3></div> <p class="text-sm text-green-700">Visual query builder with joins, filters, and aggregations</p></div> <div class="bg-purple-50 p-4 rounded-lg"><div class="flex items-center space-x-2 mb-2"><span class="text-purple-600">⚡</span> <h3 class="font-medium text-purple-900">Workflow Automation</h3></div> <p class="text-sm text-purple-700">Approval chains, notifications, and automated processes</p></div> <div class="bg-orange-50 p-4 rounded-lg"><div class="flex items-center space-x-2 mb-2"><span class="text-orange-600">📈</span> <h3 class="font-medium text-orange-900">Smart Reporting</h3></div> <p class="text-sm text-orange-700">Dashboards, charts, and export to Excel/PDF</p></div></div></div> ${forms.length === 0 ? ` <div class="text-center py-12" data-svelte-h="svelte-1j2p0xu"><div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4"><svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg></div> <h3 class="text-lg font-medium text-gray-900 mb-2">No forms yet</h3> <p class="text-gray-500 mb-4">Create your first form to start collecting data</p> <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">Create Your First Form</button></div>` : ` <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6">${each(forms, (form, index) => {
    return `<div class="form-card bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow cursor-move" draggable="true" role="button" aria-label="${"Drag form '" + escape(form.name, true) + "' to reorder"}" tabindex="0"> <div class="flex items-center justify-between mb-4"><div class="flex items-center space-x-2" data-svelte-h="svelte-14sakh3"><svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16"></path></svg> <span class="text-xs text-gray-500">Drag to reorder</span></div> <span class="${"px-2 py-1 text-xs font-medium rounded-full " + escape(getStatusColor(form.status), true)}">${escape(form.status)} </span></div> <div class="flex items-start justify-between mb-4"><div class="flex-1"><h3 class="text-lg font-semibold text-gray-900 mb-1">${escape(form.name)}</h3> <p class="text-sm text-gray-600 mb-2">${escape(form.description)}</p> <div class="flex items-center space-x-4 text-sm text-gray-500"><span>${escape(form.responses)} responses</span> <span>Created ${escape(form.createdAt.toLocaleDateString())}</span></div> </div></div> <div class="flex items-center justify-between"><div class="flex space-x-2"><button class="px-3 py-1 text-sm bg-blue-50 text-blue-700 rounded hover:bg-blue-100 transition-colors" data-svelte-h="svelte-x9fdrc">View Responses</button> <button class="px-3 py-1 text-sm bg-gray-50 text-gray-700 rounded hover:bg-gray-100 transition-colors" data-svelte-h="svelte-1quezoi">Edit</button> <button class="px-3 py-1 text-sm bg-indigo-50 text-indigo-700 rounded hover:bg-indigo-100 transition-colors" data-svelte-h="svelte-n526i1">Advanced
								</button></div> <button class="p-1 text-gray-400 hover:text-red-600 transition-colors" title="Delete form" data-svelte-h="svelte-12ae6pm"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg> </button></div> </div>`;
  })}</div>  <div class="mt-8 bg-white rounded-lg shadow-sm border border-gray-200 p-6"><h3 class="text-lg font-semibold text-gray-900 mb-4" data-svelte-h="svelte-cwl39c">Summary</h3> <div class="grid grid-cols-1 md:grid-cols-3 gap-4"><div class="text-center"><div class="text-2xl font-bold text-blue-600">${escape(forms.length)}</div> <div class="text-sm text-gray-600" data-svelte-h="svelte-je22rm">Total Forms</div></div> <div class="text-center"><div class="text-2xl font-bold text-green-600">${escape(forms.reduce((sum, f) => sum + f.responses, 0))}</div> <div class="text-sm text-gray-600" data-svelte-h="svelte-gr9orx">Total Responses</div></div> <div class="text-center"><div class="text-2xl font-bold text-purple-600">${escape(forms.filter((f) => f.status === "active").length)}</div> <div class="text-sm text-gray-600" data-svelte-h="svelte-ikhc2e">Active Forms</div></div></div></div>`}</div>`;
});
function getCellId(row, col) {
  let colStr = "";
  let c = col + 1;
  while (c > 0) {
    c--;
    colStr = String.fromCharCode(65 + c % 26) + colStr;
    c = Math.floor(c / 26);
  }
  return `${colStr}${row + 1}`;
}
function getCellStyle(format) {
  if (!format) return "";
  const styles = [];
  if (format.bold) styles.push("font-weight: bold");
  if (format.italic) styles.push("font-style: italic");
  if (format.underline) styles.push("text-decoration: underline");
  if (format.fontSize) styles.push(`font-size: ${format.fontSize}px`);
  if (format.color) styles.push(`color: ${format.color}`);
  if (format.backgroundColor) styles.push(`background-color: ${format.backgroundColor}`);
  if (format.align) styles.push(`text-align: ${format.align}`);
  if (format.borderTop) styles.push("border-top: 1px solid #ccc");
  if (format.borderBottom) styles.push("border-bottom: 1px solid #ccc");
  if (format.borderLeft) styles.push("border-left: 1px solid #ccc");
  if (format.borderRight) styles.push("border-right: 1px solid #ccc");
  if (format.borderAll) styles.push("border: 1px solid #ccc");
  return styles.join("; ");
}
const css = {
  code: "input.svelte-qenbhy:focus{box-shadow:0 0 0 2px rgba(59, 130, 246, 0.5)}",
  map: `{"version":3,"file":"FormulaBar.svelte","sources":["FormulaBar.svelte"],"sourcesContent":["<script>\\r\\n\\timport { createEventDispatcher } from 'svelte';\\r\\n\\timport {\\r\\n\\t\\tselectedCell,\\r\\n\\t\\tformulaBar,\\r\\n\\t\\tshowFormulaHelp,\\r\\n\\t\\tisFormulaRangeSelecting,\\r\\n\\t\\tformulaBeingEdited,\\r\\n\\t\\tselectedCells,\\r\\n\\t\\tselectCellRange,\\r\\n\\t\\tselectSingleCell\\r\\n\\t} from '../stores/spreadsheet-store.js';\\r\\n\\r\\n\\timport { getCellId, getCellFromId } from '../utils/spreadsheet-utils.js';\\r\\n\\r\\n\\tconst dispatch = createEventDispatcher();\\r\\n\\r\\n\\tlet nameBoxValue = '';\\r\\n\\r\\n\\t// Update name box when selected cell changes\\r\\n\\t$: if ($selectedCell) {\\r\\n\\t\\tconst selectedCount = $selectedCells.size;\\r\\n\\t\\tif (selectedCount > 1) {\\r\\n\\t\\t\\t// Show range if multiple cells selected\\r\\n\\t\\t\\tconst cells = Array.from($selectedCells).sort();\\r\\n\\t\\t\\tnameBoxValue = \`\${cells[0]}:\${cells[cells.length - 1]}\`;\\r\\n\\t\\t} else {\\r\\n\\t\\t\\tnameBoxValue = getCellId($selectedCell.row, $selectedCell.col);\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\t// Handle name box input for navigation\\r\\n\\tfunction handleNameBoxKeyDown(event) {\\r\\n\\t\\tif (event.key === 'Enter') {\\r\\n\\t\\t\\tnavigateToCell(nameBoxValue);\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\tfunction navigateToCell(reference) {\\r\\n\\t\\t// Parse cell reference (e.g., \\"A1\\", \\"B2:D10\\")\\r\\n\\t\\tconst rangeMatch = reference.match(/^([A-Z]+)(\\\\d+):([A-Z]+)(\\\\d+)$/);\\r\\n\\t\\tconst singleMatch = reference.match(/^([A-Z]+)(\\\\d+)$/);\\r\\n\\r\\n\\t\\tif (rangeMatch) {\\r\\n\\t\\t\\t// Range selection\\r\\n\\t\\t\\tconst startCol = rangeMatch[1].charCodeAt(0) - 65;\\r\\n\\t\\t\\tconst startRow = parseInt(rangeMatch[2]) - 1;\\r\\n\\t\\t\\tconst endCol = rangeMatch[3].charCodeAt(0) - 65;\\r\\n\\t\\t\\tconst endRow = parseInt(rangeMatch[4]) - 1;\\r\\n\\r\\n\\t\\t\\tif (isValidCell(startRow, startCol) && isValidCell(endRow, endCol)) {\\r\\n\\t\\t\\t\\tselectCellRange(startRow, startCol, endRow, endCol);\\r\\n\\t\\t\\t\\tselectedCell.set({ row: startRow, col: startCol });\\r\\n\\t\\t\\t\\tdispatch('action', { action: 'selectCell', row: startRow, col: startCol });\\r\\n\\t\\t\\t}\\r\\n\\t\\t} else if (singleMatch) {\\r\\n\\t\\t\\t// Single cell selection\\r\\n\\t\\t\\tconst col = singleMatch[1].charCodeAt(0) - 65;\\r\\n\\t\\t\\tconst row = parseInt(singleMatch[2]) - 1;\\r\\n\\r\\n\\t\\t\\tif (isValidCell(row, col)) {\\r\\n\\t\\t\\t\\tselectSingleCell(row, col);\\r\\n\\t\\t\\t\\tdispatch('action', { action: 'selectCell', row, col });\\r\\n\\t\\t\\t}\\r\\n\\t\\t}\\r\\n\\t}\\r\\n\\r\\n\\tfunction isValidCell(row, col) {\\r\\n\\t\\treturn row >= 0 && row < 100 && col >= 0 && col < 26;\\r\\n\\t}\\r\\n\\r\\n\\t// Handle formula input to enable range selection for ANY function that needs ranges\\r\\n\\r\\n\\tfunction handleFormulaInput(event) {\\r\\n\\t\\tconst value = event.target.value;\\r\\n\\r\\n\\t\\t// Check if user is typing any function that expects a range parameter\\r\\n\\t\\t// Pattern: =FUNCTIONNAME( with no closing ) yet\\r\\n\\t\\tconst functionMatch = value.toUpperCase().match(/^=([A-Z_]+)\\\\([^)]*$/);\\r\\n\\r\\n\\t\\tif (functionMatch) {\\r\\n\\t\\t\\tconst functionName = functionMatch[1];\\r\\n\\r\\n\\t\\t\\t// Common functions that expect ranges (can be extended)\\r\\n\\t\\t\\tconst rangeFunctions = [\\r\\n\\t\\t\\t\\t'SUM', 'AVERAGE', 'AVG', 'COUNT', 'MIN', 'MAX',\\r\\n\\t\\t\\t\\t'STDEV', 'STDEVP', 'VAR', 'VARP',\\r\\n\\t\\t\\t\\t'MEDIAN', 'MODE', 'RANK',\\r\\n\\t\\t\\t\\t'PRODUCT', 'GEOMEAN', 'HARMEAN',\\r\\n\\t\\t\\t\\t'SUBTOTAL', 'AGGREGATE'\\r\\n\\t\\t\\t];\\r\\n\\r\\n\\t\\t\\tif (rangeFunctions.includes(functionName)) {\\r\\n\\t\\t\\t\\t// Extract the partial formula up to the opening parenthesis\\r\\n\\t\\t\\t\\tconst openParenIndex = value.indexOf('(');\\r\\n\\t\\t\\t\\tformulaBeingEdited.set(value.substring(0, openParenIndex + 1)); // \\"=FUNCTION(\\"\\r\\n\\t\\t\\t\\tisFormulaRangeSelecting.set(true);\\r\\n\\t\\t\\t\\tconsole.log('Formula range selection enabled for:', functionName);\\r\\n\\t\\t\\t\\treturn;\\r\\n\\t\\t\\t}\\r\\n\\t\\t}\\r\\n\\r\\n\\t\\t// Reset formula range selection if formula is completed or changed\\r\\n\\t\\tisFormulaRangeSelecting.set(false);\\r\\n\\t\\tformulaBeingEdited.set('');\\r\\n\\t}\\r\\n\\r\\n\\t// Handle Enter key in formula bar\\r\\n\\tfunction handleFormulaKeyDown(event) {\\r\\n\\t\\tif (event.key === 'Enter') {\\r\\n\\t\\t\\tevent.preventDefault();\\r\\n\\t\\t\\t// Update the cell value and formula bar\\r\\n\\t\\t\\tconst { cellInput, spreadsheetData } = require('../stores/spreadsheet-store.js');\\r\\n\\t\\t\\tif ($selectedCell) {\\r\\n\\t\\t\\t\\tspreadsheetData.update(data => {\\r\\n\\t\\t\\t\\t\\tif (!data[$selectedCell.row]) data[$selectedCell.row] = [];\\r\\n\\t\\t\\t\\t\\tdata[$selectedCell.row][$selectedCell.col] = $formulaBar;\\r\\n\\t\\t\\t\\t\\treturn [...data];\\r\\n\\t\\t\\t\\t});\\r\\n\\t\\t\\t\\tcellInput.set($formulaBar);\\r\\n\\t\\t\\t}\\r\\n\\t\\t} else if (event.key === 'Escape') {\\r\\n\\t\\t\\t// Cancel range selection\\r\\n\\t\\t\\tisFormulaRangeSelecting.set(false);\\r\\n\\t\\t\\tformulaBeingEdited.set('');\\r\\n\\t\\t}\\r\\n\\t}\\r\\n<\/script>\\r\\n\\r\\n<!-- Formula bar with Name Box -->\\r\\n<div class=\\"flex items-center px-4 py-2 border-b border-gray-200 bg-white\\">\\r\\n\\t<!-- Name Box -->\\r\\n\\t<div class=\\"flex items-center mr-4\\">\\r\\n\\t\\t<span class=\\"text-xs text-gray-500 mr-2\\">Name Box</span>\\r\\n\\t\\t<input\\r\\n\\t\\t\\ttype=\\"text\\"\\r\\n\\t\\t\\tbind:value={nameBoxValue}\\r\\n\\t\\t\\tclass=\\"w-24 px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-medium\\"\\r\\n\\t\\t\\tplaceholder=\\"A1\\"\\r\\n\\t\\t\\ton:keydown={handleNameBoxKeyDown}\\r\\n\\t\\t\\ttitle=\\"Type cell reference (e.g., A1 or B2:D10) and press Enter\\"\\r\\n\\t\\t/>\\r\\n\\t</div>\\r\\n\\r\\n\\t<div class=\\"w-px h-6 bg-gray-300 mx-2\\"></div>\\r\\n\\r\\n\\t<!-- Formula Input -->\\r\\n\\t<span class=\\"text-sm font-medium text-gray-600 mr-2\\">fx</span>\\r\\n\\t<div class=\\"flex-1 relative\\">\\r\\n\\r\\n\\t\\t<input\\r\\n\\t\\t\\tbind:value={$formulaBar}\\r\\n\\t\\t\\tplaceholder=\\"Enter formula or value\\"\\r\\n\\t\\t\\tclass=\\"w-full px-3 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent\\"\\r\\n\\t\\t\\ton:input={handleFormulaInput}\\r\\n\\t\\t\\ton:keydown={handleFormulaKeyDown}\\r\\n\\t\\t/>\\r\\n\\t\\t{#if $isFormulaRangeSelecting}\\r\\n\\t\\t\\t<div class=\\"absolute right-2 top-1/2 transform -translate-y-1/2 text-xs text-green-600 font-medium bg-green-100 px-2 py-1 rounded\\">\\r\\n\\t\\t\\t\\tClick and drag to select range\\r\\n\\t\\t\\t</div>\\r\\n\\t\\t{/if}\\r\\n\\t</div>\\r\\n</div>\\r\\n\\r\\n<style>\\r\\n\\tinput:focus {\\r\\n\\t\\tbox-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5);\\r\\n\\t}\\r\\n</style>\\r\\n"],"names":[],"mappings":"AAsKC,mBAAK,MAAO,CACX,UAAU,CAAE,CAAC,CAAC,CAAC,CAAC,CAAC,CAAC,GAAG,CAAC,KAAK,EAAE,CAAC,CAAC,GAAG,CAAC,CAAC,GAAG,CAAC,CAAC,GAAG,CAC7C"}`
};
const FormulaBar = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $formulaBar, $$unsubscribe_formulaBar;
  let $selectedCell, $$unsubscribe_selectedCell;
  let $selectedCells, $$unsubscribe_selectedCells;
  let $isFormulaRangeSelecting, $$unsubscribe_isFormulaRangeSelecting;
  $$unsubscribe_formulaBar = subscribe(formulaBar, (value) => $formulaBar = value);
  $$unsubscribe_selectedCell = subscribe(selectedCell, (value) => $selectedCell = value);
  $$unsubscribe_selectedCells = subscribe(selectedCells, (value) => $selectedCells = value);
  $$unsubscribe_isFormulaRangeSelecting = subscribe(isFormulaRangeSelecting, (value) => $isFormulaRangeSelecting = value);
  createEventDispatcher();
  let nameBoxValue = "";
  $$result.css.add(css);
  {
    if ($selectedCell) {
      const selectedCount = $selectedCells.size;
      if (selectedCount > 1) {
        const cells = Array.from($selectedCells).sort();
        nameBoxValue = `${cells[0]}:${cells[cells.length - 1]}`;
      } else {
        nameBoxValue = getCellId($selectedCell.row, $selectedCell.col);
      }
    }
  }
  $$unsubscribe_formulaBar();
  $$unsubscribe_selectedCell();
  $$unsubscribe_selectedCells();
  $$unsubscribe_isFormulaRangeSelecting();
  return ` <div class="flex items-center px-4 py-2 border-b border-gray-200 bg-white"> <div class="flex items-center mr-4"><span class="text-xs text-gray-500 mr-2" data-svelte-h="svelte-nyjpbf">Name Box</span> <input type="text" class="w-24 px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-medium svelte-qenbhy" placeholder="A1" title="Type cell reference (e.g., A1 or B2:D10) and press Enter"${add_attribute("value", nameBoxValue, 0)}></div> <div class="w-px h-6 bg-gray-300 mx-2"></div>  <span class="text-sm font-medium text-gray-600 mr-2" data-svelte-h="svelte-3jwhpc">fx</span> <div class="flex-1 relative"><input placeholder="Enter formula or value" class="w-full px-3 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent svelte-qenbhy"${add_attribute("value", $formulaBar, 0)}> ${$isFormulaRangeSelecting ? `<div class="absolute right-2 top-1/2 transform -translate-y-1/2 text-xs text-green-600 font-medium bg-green-100 px-2 py-1 rounded" data-svelte-h="svelte-1v6alhh">Click and drag to select range</div>` : ``}</div> </div>`;
});
export {
  ContactForm as C,
  EventForm as E,
  FormulaBar as F,
  getCellStyle as a,
  FormList as b,
  getCellId as g
};
