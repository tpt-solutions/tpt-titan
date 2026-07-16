import { c as create_ssr_component, s as subscribe, d as each, a as add_attribute, b as escape } from "../../chunks/calendar.js";
import { p as page } from "../../chunks/svelte-kit.js";
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let currentPath;
  let $page, $$unsubscribe_page;
  $$unsubscribe_page = subscribe(page, (value) => $page = value);
  const data = null;
  const form = null;
  const params = null;
  let mobileMenuOpen = false;
  function navClass(path) {
    const active = currentPath === path || path !== "/" && currentPath.startsWith(path);
    return `px-3 py-2 rounded-md text-sm font-medium transition-colors ${active ? "text-blue-600 bg-blue-50" : "text-gray-700 hover:text-blue-600 hover:bg-gray-50"}`;
  }
  const navLinks = [
    { href: "/", label: "Home" },
    { href: "/editor", label: "Editor" },
    {
      href: "/spreadsheet",
      label: "Spreadsheet"
    },
    { href: "/forms", label: "Forms" },
    { href: "/database", label: "Database" },
    { href: "/email", label: "Email" },
    { href: "/calendar", label: "Calendar" },
    { href: "/contacts", label: "Contacts" },
    { href: "/chat", label: "Chat" },
    { href: "/files", label: "Files" },
    { href: "/tasks", label: "Tasks" },
    { href: "/workflows", label: "Workflows" },
    { href: "/speech", label: "Speech" },
    { href: "/voice", label: "Voice" },
    { href: "/math", label: "Math" },
    { href: "/export", label: "Export" },
    { href: "/monitoring", label: "Monitoring" },
    { href: "/admin", label: "Admin" },
    { href: "/plugins", label: "Plugins" },
    { href: "/settings", label: "Settings" }
  ];
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  currentPath = $page.url.pathname;
  $$unsubscribe_page();
  return `<main class="min-h-screen bg-gray-50"> <header class="bg-white shadow-sm border-b border-gray-200"><div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"><div class="flex justify-between items-center h-16"> <div class="flex items-center" data-svelte-h="svelte-o9r20v"><a href="/" class="text-2xl font-bold text-blue-600 hover:text-blue-700 transition-colors">TPT Titan</a></div>  <nav class="hidden md:flex flex-wrap gap-1">${each(navLinks, (link) => {
    return `<a${add_attribute("href", link.href, 0)}${add_attribute("class", navClass(link.href), 0)}>${escape(link.label)} </a>`;
  })}</nav>  <div class="md:hidden"><button aria-label="Toggle navigation menu"${add_attribute("aria-expanded", mobileMenuOpen, 0)} class="text-gray-700 hover:text-blue-600 p-2 rounded-md">${`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>`}</button></div></div></div>  ${``}</header>  ${slots.default ? slots.default({}) : ``} </main>`;
});
export {
  Layout as default
};
