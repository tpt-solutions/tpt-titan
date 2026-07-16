import { c as create_ssr_component, b as escape, v as validate_component } from "../../../chunks/calendar.js";
import { T as TaskBoard } from "../../../chunks/tasks.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { data = null } = $$props;
  let { form = null } = $$props;
  let { params = null } = $$props;
  let tasks = [];
  let projects = [];
  function getProjectName(projectId) {
    const project = projects.find((p) => p.id === projectId);
    return project ? project.name : "No Project";
  }
  function getTasksByStatus(status) {
    return tasks.filter((t) => t.status === status);
  }
  function getTaskStats() {
    const total = tasks.length;
    const completed = tasks.filter((t) => t.status === "done").length;
    const overdue = tasks.filter((t) => t.dueDate && t.dueDate < /* @__PURE__ */ new Date() && t.status !== "done").length;
    const highPriority = tasks.filter((t) => t.priority === "high" && t.status !== "done").length;
    return { total, completed, overdue, highPriority };
  }
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  return `${$$result.head += `<!-- HEAD_svelte-1ms7gyv_START -->${$$result.title = `<title>Tasks &amp; Projects - TPT Titan</title>`, ""}<!-- HEAD_svelte-1ms7gyv_END -->`, ""} <div class="h-screen flex flex-col bg-gray-50"> <header class="flex items-center justify-between px-6 py-3 border-b border-gray-200 bg-white"><div class="flex items-center space-x-4" data-svelte-h="svelte-10irm90"><h1 class="text-xl font-semibold text-gray-900">Tasks &amp; Projects</h1> <span class="text-sm text-gray-500">AI-powered project management</span></div> <div class="flex items-center space-x-2">${`<button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors" data-svelte-h="svelte-1cbhgdk">New Task</button>`}</div></header>  ${``}  <div class="bg-white border-b border-gray-200 px-6 py-4"><div class="grid grid-cols-1 md:grid-cols-4 gap-4"><div class="text-center"><div class="text-2xl font-bold text-blue-600">${escape(getTaskStats().total)}</div> <div class="text-sm text-gray-600" data-svelte-h="svelte-ch92cb">Total Tasks</div></div> <div class="text-center"><div class="text-2xl font-bold text-green-600">${escape(getTaskStats().completed)}</div> <div class="text-sm text-gray-600" data-svelte-h="svelte-vrz6bu">Completed</div></div> <div class="text-center"><div class="text-2xl font-bold text-red-600">${escape(getTaskStats().overdue)}</div> <div class="text-sm text-gray-600" data-svelte-h="svelte-19nropf">Overdue</div></div> <div class="text-center"><div class="text-2xl font-bold text-orange-600">${escape(getTaskStats().highPriority)}</div> <div class="text-sm text-gray-600" data-svelte-h="svelte-1ttppej">High Priority</div></div></div></div>  <div class="flex-1 overflow-hidden">${`${`${validate_component(TaskBoard, "TaskBoard").$$render(
    $$result,
    {
      tasks,
      projects,
      getProjectName,
      getTasksByStatus
    },
    {},
    {}
  )}`}`}</div></div>`;
});
export {
  Page as default
};
