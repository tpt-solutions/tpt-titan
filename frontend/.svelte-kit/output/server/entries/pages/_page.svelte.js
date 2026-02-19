import { c as create_ssr_component, d as each, b as escape, a as add_attribute, v as validate_component } from "../../chunks/calendar.js";
import "@sveltejs/kit/internal";
import "../../chunks/svelte-kit.js";
import "@sveltejs/kit/internal/server";
function getColorClasses(color, isSelected) {
  const baseClasses = "p-6 rounded-lg shadow-md transition-all duration-200 cursor-pointer border-2";
  if (isSelected) {
    switch (color) {
      case "blue":
        return `${baseClasses} bg-blue-50 border-blue-500 text-blue-700`;
      case "green":
        return `${baseClasses} bg-green-50 border-green-500 text-green-700`;
      case "purple":
        return `${baseClasses} bg-purple-50 border-purple-500 text-purple-700`;
      case "yellow":
        return `${baseClasses} bg-yellow-50 border-yellow-500 text-yellow-700`;
      case "red":
        return `${baseClasses} bg-red-50 border-red-500 text-red-700`;
      case "indigo":
        return `${baseClasses} bg-indigo-50 border-indigo-500 text-indigo-700`;
      case "pink":
        return `${baseClasses} bg-pink-50 border-pink-500 text-pink-700`;
      case "orange":
        return `${baseClasses} bg-orange-50 border-orange-500 text-orange-700`;
      default:
        return `${baseClasses} bg-gray-50 border-gray-500 text-gray-700`;
    }
  } else {
    return `${baseClasses} bg-white border-gray-200 hover:border-gray-300 text-gray-900 hover:shadow-lg`;
  }
}
const AppSelector = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let selectedApp = "spreadsheet";
  const apps = [
    {
      id: "spreadsheet",
      name: "Spreadsheet",
      description: "Powerful spreadsheet with formulas and charts",
      icon: "📊",
      color: "blue",
      route: "/spreadsheet"
    },
    {
      id: "forms",
      name: "Forms",
      description: "Create and manage forms with workflows",
      icon: "📋",
      color: "green",
      route: "/forms"
    },
    {
      id: "editor",
      name: "Text Editor",
      description: "Rich text editing with collaboration",
      icon: "📝",
      color: "purple",
      route: "/editor"
    },
    {
      id: "tasks",
      name: "Tasks",
      description: "Task management and project tracking",
      icon: "✅",
      color: "yellow",
      route: "/tasks"
    },
    {
      id: "calendar",
      name: "Calendar",
      description: "Calendar with scheduling and reminders",
      icon: "📅",
      color: "red",
      route: "/calendar"
    },
    {
      id: "contacts",
      name: "Contacts",
      description: "Contact management and organization",
      icon: "👥",
      color: "indigo",
      route: "/contacts"
    },
    {
      id: "email",
      name: "Email",
      description: "Email client with encryption",
      icon: "📧",
      color: "pink",
      route: "/email"
    },
    {
      id: "database",
      name: "Database",
      description: "Database management and queries",
      icon: "🗄️",
      color: "orange",
      route: "/database"
    }
  ];
  return `<div class="max-w-6xl mx-auto"><div class="text-center mb-8" data-svelte-h="svelte-1tmx4k"><h2 class="text-3xl font-bold text-gray-900 mb-4">Choose Your Default App</h2> <p class="text-lg text-gray-600 max-w-2xl mx-auto">Select which application you&#39;d like to load by default. This will significantly improve your startup time by only loading the app you use most.</p> <div class="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg"><p class="text-sm text-blue-800"><strong>Performance Tip:</strong> Choosing one app reduces initial load from 135MB to ~2-3MB, cutting startup time from 35+ seconds to under 5 seconds.</p></div></div> <div class="grid md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mb-8">${each(apps, (app) => {
    return `<div${add_attribute("class", getColorClasses(app.color, selectedApp === app.id), 0)} role="button" tabindex="0"><div class="flex flex-col items-center text-center"><div class="text-4xl mb-3">${escape(app.icon)}</div> <h3 class="text-lg font-semibold mb-2">${escape(app.name)}</h3> <p class="text-sm opacity-80 leading-relaxed">${escape(app.description)}</p> ${selectedApp === app.id ? `<div class="mt-3 px-3 py-1 bg-current text-white text-xs rounded-full" data-svelte-h="svelte-gzlp40">Default App
						</div>` : ``}</div> </div>`;
  })}</div> <div class="text-center"><button class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold text-lg shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed" ${""}>${`Launch ${escape(apps.find((a) => a.id === selectedApp)?.name)}`}</button> <p class="text-sm text-gray-500 mt-4" data-svelte-h="svelte-310lrl">You can change your default app anytime from the app menu</p></div></div>`;
});
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const params = null;
  const data = null;
  const form = null;
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  return `${$$result.head += `<!-- HEAD_svelte-135yy8d_START -->${$$result.title = `<title>TPT Titan - Open Source Office Suite</title>`, ""}<meta name="description" content="Complete open source alternative to Microsoft Office 365"><!-- HEAD_svelte-135yy8d_END -->`, ""} ${`<div class="min-h-screen bg-gray-50 py-8"><div class="container mx-auto px-4"><header class="text-center mb-8" data-svelte-h="svelte-19mocrq"><h1 class="text-5xl font-bold text-gray-900 mb-4">Welcome to <span class="text-blue-600">TPT Titan</span></h1> <p class="text-xl text-gray-600 max-w-2xl mx-auto mb-6">A complete open source alternative to Microsoft Office 365.
					Built with modern web technologies for productivity, collaboration, and privacy.</p></header> ${validate_component(AppSelector, "AppSelector").$$render($$result, {}, {}, {})} <div class="text-center mt-8" data-svelte-h="svelte-1h066qf"><div class="bg-gray-50 border border-gray-200 rounded-lg p-6"><h3 class="text-lg font-semibold text-gray-900 mb-2">🚀 Performance Optimized</h3> <p class="text-gray-600 mb-4">TPT Titan uses lazy loading to start up to 7x faster. Choose your default app above to begin.</p> <div class="flex justify-center space-x-4"><a href="/docs" class="text-blue-600 hover:text-blue-800 underline">View Documentation</a> <a href="https://github.com/tpt-titan/tpt-titan" class="text-gray-600 hover:text-gray-800 underline">GitHub Repository</a> <a href="?select-app" class="text-gray-600 hover:text-gray-800 underline">Change Default App</a></div></div></div></div></div>`}`;
});
export {
  Page as default
};
