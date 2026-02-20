import { c as create_ssr_component, x as emailAccounts, y as emails, z as selectedEmail, A as currentFolder, d as each, v as validate_component, b as escape } from "../../../chunks/calendar.js";
import "@sveltejs/kit/internal";
import "../../../chunks/svelte-kit.js";
import "@sveltejs/kit/internal/server";
import { E as EmailInbox, a as EmailViewer } from "../../../chunks/email.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { data = null } = $$props;
  let { form = null } = $$props;
  let { params = null } = $$props;
  let emailAccountsList = [];
  let emailsList = [];
  let selectedEmailData = null;
  let currentFolderValue = "inbox";
  emailAccounts.subscribe((value) => emailAccountsList = value);
  emails.subscribe((value) => emailsList = value);
  selectedEmail.subscribe((value) => selectedEmailData = value);
  currentFolder.subscribe((value) => currentFolderValue = value);
  async function loadEmails() {
    try {
      const params2 = new URLSearchParams();
      if (currentFolderValue) params2.append("folder", currentFolderValue);
      const response = await fetch(`/api/v1/emails?${params2}`, {
        headers: {
          "Authorization": `Bearer ${localStorage.getItem("token")}`
        }
      });
      if (response.ok) {
        const data2 = await response.json();
        emails.set(data2.emails || []);
      }
    } catch (error) {
      console.error("Failed to load emails:", error);
    }
  }
  function handleEmailSelect(email) {
    selectedEmail.set(email);
  }
  function handleFolderChange(folder) {
    currentFolder.set(folder);
    loadEmails();
    selectedEmail.set(null);
  }
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  {
    if (currentFolderValue) {
      loadEmails();
    }
  }
  return `${$$result.head += `<!-- HEAD_svelte-1ju8p3j_START -->${$$result.title = `<title>Email - TPT Titan</title>`, ""}<!-- HEAD_svelte-1ju8p3j_END -->`, ""} <div class="h-screen flex flex-col bg-gray-50 dark:bg-gray-900"> <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4"><div class="flex justify-between items-center"><h1 class="text-2xl font-bold text-gray-900 dark:text-white" data-svelte-h="svelte-1xsf2q4">Email</h1> <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2" data-svelte-h="svelte-14ymtr5"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
				Compose</button></div>  ${emailAccountsList.length > 0 ? `<div class="mt-4 flex items-center space-x-4"><span class="text-sm text-gray-600 dark:text-gray-400" data-svelte-h="svelte-1nlzdx1">Accounts:</span> <div class="flex space-x-2">${each(emailAccountsList, (account) => {
    return `<span class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">${escape(account.display_name || account.email)} </span>`;
  })}</div></div>` : `<div class="mt-4 text-center py-8"><svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg> <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white" data-svelte-h="svelte-10xt0ex">No email accounts</h3> <p class="mt-1 text-sm text-gray-500 dark:text-gray-400" data-svelte-h="svelte-oqfde8">Add an email account to get started.</p> <button class="mt-4 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg" data-svelte-h="svelte-1536vzg">Add Email Account</button></div>`}</div>  ${emailAccountsList.length > 0 ? `<div class="flex-1 flex overflow-hidden"> <div class="w-80 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col">${validate_component(EmailInbox, "EmailInbox").$$render(
    $$result,
    {
      emailsList,
      currentFolderValue,
      selectedEmailData,
      handleEmailSelect,
      handleFolderChange
    },
    {},
    {}
  )}</div>  <div class="flex-1 flex flex-col min-w-0">${`${selectedEmailData ? `${validate_component(EmailViewer, "EmailViewer").$$render(
    $$result,
    {
      email: selectedEmailData,
      onClose: () => selectedEmail.set(null)
    },
    {},
    {}
  )}` : ` <div class="flex-1 flex items-center justify-center bg-gray-50 dark:bg-gray-900" data-svelte-h="svelte-trn8af"><div class="text-center"><svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg> <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No email selected</h3> <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Choose an email from the list to view it.</p></div></div>`}`}</div></div>` : ``}</div>`;
});
export {
  Page as default
};
