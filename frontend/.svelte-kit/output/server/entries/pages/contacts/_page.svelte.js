import { c as create_ssr_component, g as contacts, d as each, b as escape, a as add_attribute, v as validate_component } from "../../../chunks/calendar.js";
import { g as goto } from "../../../chunks/svelte-kit.js";
import { C as ContactForm } from "../../../chunks/forms.js";
function getContactDisplayName(contact) {
  const firstName = contact.first_name || "";
  const lastName = contact.last_name || "";
  const fullName = `${firstName} ${lastName}`.trim();
  return fullName || "Unnamed Contact";
}
function getContactInitials(contact) {
  const firstInitial = contact.first_name ? contact.first_name.charAt(0).toUpperCase() : "";
  const lastInitial = contact.last_name ? contact.last_name.charAt(0).toUpperCase() : "";
  return firstInitial + lastInitial || "?";
}
const ContactList = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { handleEditContact } = $$props;
  let contactList = [];
  contacts.subscribe((value) => {
    contactList = value;
  });
  if ($$props.handleEditContact === void 0 && $$bindings.handleEditContact && handleEditContact !== void 0) $$bindings.handleEditContact(handleEditContact);
  return `<div class="bg-white dark:bg-gray-800 rounded-lg shadow">${contactList.length === 0 ? `<div class="p-8 text-center" data-svelte-h="svelte-1hvcx3j"><svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg> <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No contacts</h3> <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Get started by adding your first contact.</p></div>` : `<div class="divide-y divide-gray-200 dark:divide-gray-700">${each(contactList, (contact) => {
    return `<div class="p-6 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"><div class="flex items-center justify-between"><div class="flex items-center space-x-4"> <div class="flex-shrink-0"><div class="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold text-lg">${escape(getContactInitials(contact))} </div></div>  <div class="flex-1 min-w-0"><div class="flex items-center space-x-3"><h3 class="text-lg font-medium text-gray-900 dark:text-white truncate">${escape(getContactDisplayName(contact))}</h3> ${contact.email ? `<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200" data-svelte-h="svelte-br2qkn">Has Email
										</span>` : ``}</div> <div class="mt-1 space-y-1">${contact.email ? `<p class="text-sm text-gray-500 dark:text-gray-400 flex items-center"><svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg> ${escape(contact.email)} </p>` : ``} ${contact.phone ? `<p class="text-sm text-gray-500 dark:text-gray-400 flex items-center"><svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg> ${escape(contact.phone)} </p>` : ``} ${contact.company ? `<p class="text-sm text-gray-500 dark:text-gray-400 flex items-center"><svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path></svg> ${escape(contact.company)} ${contact.position ? `• ${escape(contact.position)}` : ``} </p>` : ``}</div> </div></div>  <div class="flex items-center space-x-2"><button class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300 p-1" title="Edit contact" data-svelte-h="svelte-42091y"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg></button> <button class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 p-1" title="Delete contact" data-svelte-h="svelte-1fojmxb"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button> </div></div> </div>`;
  })}</div>`}</div>`;
});
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { params = null } = $$props;
  let { data = null } = $$props;
  let { form = null } = $$props;
  let showForm = false;
  let editingContact = null;
  let searchQuery = "";
  async function loadContacts() {
    try {
      const response = await fetch("/api/v1/contacts", {
        headers: {
          "Authorization": `Bearer ${localStorage.getItem("token")}`
        }
      });
      if (response.ok) {
        const data2 = await response.json();
        contacts.set(data2.contacts || []);
      } else if (response.status === 401) {
        goto("/auth/login");
      }
    } catch (error) {
      console.error("Failed to load contacts:", error);
    }
  }
  function handleEditContact(contact) {
    editingContact = contact;
    showForm = true;
  }
  function handleFormClose() {
    showForm = false;
    editingContact = null;
    loadContacts();
  }
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  return `${$$result.head += `<!-- HEAD_svelte-1m7vk1k_START -->${$$result.title = `<title>Contacts - TPT Titan</title>`, ""}<!-- HEAD_svelte-1m7vk1k_END -->`, ""} <div class="container mx-auto px-4 py-8"><div class="flex justify-between items-center mb-8"><h1 class="text-3xl font-bold text-gray-900 dark:text-white" data-svelte-h="svelte-qbjbfc">Contacts</h1> <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2" data-svelte-h="svelte-jwj0yw"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
			Add Contact</button></div>  <div class="mb-6"><div class="flex gap-2"><input type="text" placeholder="Search contacts..." class="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"${add_attribute("value", searchQuery, 0)}> <button class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-lg flex items-center gap-2" data-svelte-h="svelte-1v0alkm"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
				Search</button></div></div>  ${validate_component(ContactList, "ContactList").$$render($$result, { handleEditContact }, {}, {})}  ${showForm ? `${validate_component(ContactForm, "ContactForm").$$render(
    $$result,
    {
      contact: editingContact,
      onClose: handleFormClose
    },
    {},
    {}
  )}` : ``}</div>`;
});
export {
  Page as default
};
