import { c as create_ssr_component, f as createEventDispatcher, d as each, b as escape } from "./calendar.js";
function getPriorityColor(priority) {
  switch (priority) {
    case "high":
      return "bg-red-100 text-red-800 border-red-200";
    case "medium":
      return "bg-yellow-100 text-yellow-800 border-yellow-200";
    case "low":
      return "bg-green-100 text-green-800 border-green-200";
    default:
      return "bg-gray-100 text-gray-800 border-gray-200";
  }
}
function formatDate(date) {
  if (!date) return "";
  return new Date(date).toLocaleDateString();
}
function isOverdue(dueDate, status) {
  if (!dueDate || status === "done") return false;
  return new Date(dueDate) < /* @__PURE__ */ new Date();
}
function getSubtaskProgress(subtasks) {
  if (!subtasks || subtasks.length === 0) return null;
  const completed = subtasks.filter((st) => st.completed).length;
  return { completed, total: subtasks.length };
}
const TaskBoard = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { tasks = [] } = $$props;
  let { projects = [] } = $$props;
  let { getProjectName } = $$props;
  let { getTasksByStatus } = $$props;
  createEventDispatcher();
  let draggedTask = null;
  let draggedOverColumn = null;
  function getProjectColor(projectId) {
    const project = projects.find((p) => p.id === projectId);
    if (!project) return "bg-gray-100";
    const colors = {
      blue: "bg-blue-100 border-blue-200",
      green: "bg-green-100 border-green-200",
      purple: "bg-purple-100 border-purple-200",
      red: "bg-red-100 border-red-200",
      yellow: "bg-yellow-100 border-yellow-200"
    };
    return colors[project.color] || "bg-gray-100 border-gray-200";
  }
  const columns = [
    {
      id: "todo",
      title: "To Do",
      color: "bg-gray-100"
    },
    {
      id: "in-progress",
      title: "In Progress",
      color: "bg-blue-100"
    },
    {
      id: "review",
      title: "Review",
      color: "bg-yellow-100"
    },
    {
      id: "done",
      title: "Done",
      color: "bg-green-100"
    }
  ];
  if ($$props.tasks === void 0 && $$bindings.tasks && tasks !== void 0) $$bindings.tasks(tasks);
  if ($$props.projects === void 0 && $$bindings.projects && projects !== void 0) $$bindings.projects(projects);
  if ($$props.getProjectName === void 0 && $$bindings.getProjectName && getProjectName !== void 0) $$bindings.getProjectName(getProjectName);
  if ($$props.getTasksByStatus === void 0 && $$bindings.getTasksByStatus && getTasksByStatus !== void 0) $$bindings.getTasksByStatus(getTasksByStatus);
  return `<div class="h-full p-6 overflow-x-auto"><div class="flex space-x-6 min-w-max">${each(columns, (column) => {
    return `<div class="flex-1 min-w-80"> <div class="flex items-center justify-between mb-4"><div class="flex items-center space-x-2"><div class="${"w-3 h-3 rounded-full " + escape(column.color, true)}"></div> <h3 class="text-lg font-semibold text-gray-900">${escape(column.title)}</h3> <span class="bg-gray-200 text-gray-700 px-2 py-1 rounded-full text-sm">${escape(getTasksByStatus(column.id).length)} </span></div> ${column.id === "todo" ? `<button class="text-gray-400 hover:text-gray-600" data-svelte-h="svelte-i4b1nd"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg> </button>` : ``}</div>  <div class="${"min-h-96 p-4 rounded-lg border-2 border-dashed transition-colors " + escape(
      draggedOverColumn === column.id ? "border-blue-400 bg-blue-50" : "border-gray-300",
      true
    )}"><div class="space-y-3">${each(getTasksByStatus(column.id), (task) => {
      return ` <div class="${"bg-white border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow cursor-move " + escape(draggedTask?.id === task.id ? "opacity-50" : "", true)}" draggable="true"> <div class="flex items-start justify-between mb-3"><div class="flex-1"><h4 class="font-medium text-gray-900 mb-1">${escape(task.title)}</h4> ${task.description ? `<p class="text-sm text-gray-600 line-clamp-2">${escape(task.description)}</p>` : ``}</div> <div class="flex space-x-1 ml-2"><button class="text-gray-400 hover:text-gray-600 p-1" title="Edit task" data-svelte-h="svelte-b1h99r"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg></button> <button class="text-gray-400 hover:text-red-600 p-1" title="Delete task" data-svelte-h="svelte-1rgza6s"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button> </div></div>  <div class="flex items-center justify-between text-xs text-gray-500 mb-3"><div class="flex items-center space-x-3">${task.assignedTo ? `<span class="flex items-center"><svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg> ${escape(task.assignedTo)} </span>` : ``} ${task.dueDate ? `<span class="${"flex items-center " + escape(
        isOverdue(task.dueDate, task.status) ? "text-red-600" : "",
        true
      )}"><svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg> ${escape(formatDate(task.dueDate))} </span>` : ``}</div> ${task.priority ? `<span class="${"px-2 py-1 rounded text-xs font-medium " + escape(getPriorityColor(task.priority), true)}">${escape(task.priority)} </span>` : ``}</div>  ${task.projectId ? `<div class="mb-3"><span class="${"inline-flex items-center px-2 py-1 rounded-full text-xs font-medium " + escape(getProjectColor(task.projectId), true)}">${escape(getProjectName(task.projectId))}</span> </div>` : ``}  ${task.tags && task.tags.length > 0 ? `<div class="flex flex-wrap gap-1 mb-3">${each(task.tags, (tag) => {
        return `<span class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-gray-100 text-gray-800">${escape(tag)} </span>`;
      })} </div>` : ``}  ${task.subtasks && task.subtasks.length > 0 ? `<div class="mb-3"><div class="flex items-center justify-between text-xs text-gray-600 mb-1"><span data-svelte-h="svelte-ogpjuy">Subtasks</span> <span>${escape(getSubtaskProgress(task.subtasks).completed)}/${escape(getSubtaskProgress(task.subtasks).total)}</span></div> <div class="w-full bg-gray-200 rounded-full h-1.5"><div class="bg-blue-600 h-1.5 rounded-full transition-all duration-300" style="${"width: " + escape(Math.round(getSubtaskProgress(task.subtasks).completed / getSubtaskProgress(task.subtasks).total * 100), true) + "%"}"></div></div> </div>` : ``}  ${task.status === "in-progress" && Math.random() > 0.7 ? `<div class="mt-3 p-2 bg-blue-50 border border-blue-200 rounded text-xs text-blue-800" data-svelte-h="svelte-1a1xzif">💡 <strong>AI Insight:</strong> This task is 2 days overdue. Consider breaking it into smaller subtasks.
					</div>` : ``} </div>`;
    })}  ${getTasksByStatus(column.id).length === 0 ? `<div class="text-center py-8 text-gray-400"><svg class="w-12 h-12 mx-auto mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path></svg> <p class="text-sm">No tasks in ${escape(column.title.toLowerCase())}</p> </div>` : ``} </div></div> </div>`;
  })}</div></div>`;
});
export {
  TaskBoard as T
};
