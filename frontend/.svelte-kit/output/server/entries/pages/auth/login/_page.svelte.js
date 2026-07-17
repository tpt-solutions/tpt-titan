import { c as create_ssr_component, a as add_attribute } from "../../../../chunks/calendar.js";
import "@sveltejs/kit/internal";
import "../../../../chunks/svelte-kit.js";
import "@sveltejs/kit/internal/server";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { data = null } = $$props;
  let { form = null } = $$props;
  let { params = null } = $$props;
  let email = "";
  let password = "";
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  return ` <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8"><div class="max-w-md w-full space-y-8"><div data-svelte-h="svelte-1alsund"><div class="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-blue-100"><svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg></div> <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">Sign in to TPT Titan</h2> <p class="mt-2 text-center text-sm text-gray-600">Your complete, private office suite</p></div> <form class="mt-8 space-y-6"><div class="rounded-md shadow-sm -space-y-px"><div><label for="email" class="sr-only" data-svelte-h="svelte-o70mpp">Email address</label> <input id="email" name="email" type="email" autocomplete="email" required class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm" placeholder="Email address"${add_attribute("value", email, 0)}></div> <div><label for="password" class="sr-only" data-svelte-h="svelte-15hhkob">Password</label> <input id="password" name="password" type="password" autocomplete="current-password" required class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm" placeholder="Password"${add_attribute("value", password, 0)}></div></div> ${``} <div><button type="submit" ${""} class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed">${`Sign in`}</button></div> <div class="text-center" data-svelte-h="svelte-zmc1dg"><a href="/auth/register" class="font-medium text-blue-600 hover:text-blue-500">Don&#39;t have an account? Sign up</a></div> <div class="text-center text-xs text-gray-500" data-svelte-h="svelte-1165enz"><a href="/setup" class="font-medium text-gray-500 hover:text-gray-700">First-time deploy? Run the setup wizard</a></div> <div class="text-center text-xs text-gray-500" data-svelte-h="svelte-owic3g"><p>🔒 End-to-end encrypted • Privacy-first • Open source</p></div></form></div></div>`;
});
export {
  Page as default
};
