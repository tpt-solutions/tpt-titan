import { c as create_ssr_component, f as createEventDispatcher, d as each, b as escape, a as add_attribute, v as validate_component } from "../../../chunks/calendar.js";
import { b as FormList } from "../../../chunks/forms.js";
function extractFieldNames(form) {
  if (form.fields && Array.isArray(form.fields)) {
    return form.fields.map((f) => f.label || f.name || f.id);
  }
  return ["id", "name", "created_at"];
}
function getRandomColor() {
  const colors = ["#3B82F6", "#10B981", "#F59E0B", "#EF4444", "#8B5CF6", "#06B6D4"];
  return colors[Math.floor(Math.random() * colors.length)];
}
const DatabaseRelationsModal = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { forms = [] } = $$props;
  let { showModal = false } = $$props;
  createEventDispatcher();
  let relationships = [];
  let canvas = null;
  let ctx = null;
  let nodes = [];
  let connections = [];
  async function initializeCanvas() {
    await new Promise((resolve) => setTimeout(resolve, 100));
    canvas = document.getElementById("relationship-canvas");
    if (canvas) {
      ctx = canvas.getContext("2d");
      resizeCanvas();
      setupNodes();
      drawCanvas();
    }
  }
  function resizeCanvas() {
    if (canvas) {
      const container = canvas.parentElement;
      canvas.width = container.clientWidth;
      canvas.height = container.clientHeight;
    }
  }
  function setupNodes() {
    nodes = forms.map((form, index) => ({
      id: form.id,
      name: form.name,
      x: 100 + index % 3 * 250,
      y: 100 + Math.floor(index / 3) * 200,
      width: 180,
      height: 100,
      color: getRandomColor(),
      fields: extractFieldNames(form)
    }));
  }
  function drawCanvas() {
    if (!ctx || !canvas) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    connections.forEach((connection) => drawConnection(connection));
    nodes.forEach((node) => drawNode(node));
  }
  function drawNode(node) {
    if (!ctx) return;
    ctx.fillStyle = node.color;
    ctx.fillRect(node.x, node.y, node.width, node.height);
    ctx.strokeStyle = "#374151";
    ctx.lineWidth = 2;
    ctx.strokeRect(node.x, node.y, node.width, node.height);
    ctx.fillStyle = "white";
    ctx.font = "bold 14px Arial";
    ctx.fillText(node.name, node.x + 10, node.y + 25);
    ctx.fillStyle = "#F3F4F6";
    ctx.font = "12px Arial";
    node.fields.slice(0, 4).forEach((field, index) => {
      ctx.fillText(`• ${field}`, node.x + 10, node.y + 45 + index * 15);
    });
    if (node.fields.length > 4) {
      ctx.fillText(`... +${node.fields.length - 4} more`, node.x + 10, node.y + 45 + 4 * 15);
    }
  }
  function drawConnection(connection) {
    if (!ctx) return;
    const sourceNode = nodes.find((n) => n.id === connection.sourceId);
    const targetNode = nodes.find((n) => n.id === connection.targetId);
    if (!sourceNode || !targetNode) return;
    const sourceX = sourceNode.x + sourceNode.width;
    const sourceY = sourceNode.y + sourceNode.height / 2;
    const targetX = targetNode.x;
    const targetY = targetNode.y + targetNode.height / 2;
    ctx.strokeStyle = connection.color || "#6B7280";
    ctx.lineWidth = 3;
    ctx.beginPath();
    ctx.moveTo(sourceX, sourceY);
    ctx.lineTo(targetX, targetY);
    ctx.stroke();
    const angle = Math.atan2(targetY - sourceY, targetX - sourceX);
    const arrowLength = 15;
    ctx.beginPath();
    ctx.moveTo(targetX, targetY);
    ctx.lineTo(targetX - arrowLength * Math.cos(angle - Math.PI / 6), targetY - arrowLength * Math.sin(angle - Math.PI / 6));
    ctx.moveTo(targetX, targetY);
    ctx.lineTo(targetX - arrowLength * Math.cos(angle + Math.PI / 6), targetY - arrowLength * Math.sin(angle + Math.PI / 6));
    ctx.stroke();
    const midX = (sourceX + targetX) / 2;
    const midY = (sourceY + targetY) / 2;
    ctx.fillStyle = "#374151";
    ctx.font = "bold 12px Arial";
    ctx.fillText(connection.type, midX - 20, midY - 5);
  }
  async function loadRelationships() {
    try {
      const allRelationships = [];
      for (const form of forms) {
        try {
          const response = await fetch(`/api/v1/form-relationships/${form.id}`, {
            headers: {
              "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
          });
          if (response.ok) {
            const data = await response.json();
            allRelationships.push(...data.relationships);
          }
        } catch (error) {
          console.warn(`Failed to load relationships for form ${form.id}:`, error);
        }
      }
      relationships = allRelationships;
      setupConnections();
      drawCanvas();
    } catch (error) {
      console.error("Failed to load relationships:", error);
    }
  }
  function setupConnections() {
    connections = relationships.map((rel) => {
      forms.find((f) => f.id === rel.source_form_id);
      forms.find((f) => f.id === rel.target_form_id);
      let color = "#6B7280";
      let type = rel.relationship_type;
      switch (rel.relationship_type) {
        case "one-to-one":
          color = "#10B981";
          type = "1:1";
          break;
        case "one-to-many":
          color = "#3B82F6";
          type = "1:N";
          break;
        case "many-to-many":
          color = "#F59E0B";
          type = "N:N";
          break;
      }
      return {
        sourceId: rel.source_form_id,
        targetId: rel.target_form_id,
        type,
        color,
        relationship: rel
      };
    });
  }
  if ($$props.forms === void 0 && $$bindings.forms && forms !== void 0) $$bindings.forms(forms);
  if ($$props.showModal === void 0 && $$bindings.showModal && showModal !== void 0) $$bindings.showModal(showModal);
  {
    if (showModal && !canvas) {
      initializeCanvas();
      loadRelationships();
    }
  }
  return `${showModal ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg w-full max-w-7xl max-h-[90vh] overflow-hidden flex flex-col"> <div class="flex items-center justify-between p-6 border-b border-gray-200"><div data-svelte-h="svelte-3t65yv"><h2 class="text-2xl font-bold text-gray-900">🔗 Database Relations Manager</h2> <p class="text-gray-600 mt-1">Create and manage relationships between your forms</p></div> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-1x0912u">×</button></div>  <div class="flex items-center justify-between p-4 bg-gray-50 border-b border-gray-200"><div class="flex items-center space-x-3"><button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors" data-svelte-h="svelte-9empqs">➕ Create Relationship</button> <button class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors" data-svelte-h="svelte-14qn5xt">🔍 Create Lookup Field</button></div> <div class="flex items-center space-x-4 text-sm" data-svelte-h="svelte-201dct"><div class="flex items-center space-x-2"><div class="w-4 h-4 bg-green-500 rounded"></div> <span>One-to-One</span></div> <div class="flex items-center space-x-2"><div class="w-4 h-4 bg-blue-500 rounded"></div> <span>One-to-Many</span></div> <div class="flex items-center space-x-2"><div class="w-4 h-4 bg-yellow-500 rounded"></div> <span>Many-to-Many</span></div></div></div>  <div class="flex-1 flex overflow-hidden"> <div class="flex-1 relative bg-gray-100"><canvas id="relationship-canvas" class="w-full h-full cursor-move"></canvas> ${nodes.length === 0 ? `<div class="absolute inset-0 flex items-center justify-center" data-svelte-h="svelte-e9unew"><div class="text-center"><div class="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-4"><svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"></path></svg></div> <h3 class="text-lg font-medium text-gray-900 mb-2">No forms to relate</h3> <p class="text-gray-500">Create some forms first to establish relationships</p></div></div>` : ``}</div>  <div class="w-80 bg-white border-l border-gray-200 p-4 overflow-y-auto"><h3 class="text-lg font-semibold text-gray-900 mb-4" data-svelte-h="svelte-1qb1i0n">Relationships</h3> ${relationships.length === 0 ? `<div class="text-center py-8 text-gray-500" data-svelte-h="svelte-1tp15mt"><p>No relationships created yet.</p> <p class="text-sm mt-1">Click &quot;Create Relationship&quot; to get started.</p></div>` : `<div class="space-y-3">${each(relationships, (rel) => {
    return `<div class="border border-gray-200 rounded-lg p-3 hover:bg-gray-50"><div class="flex items-center justify-between mb-2"><h4 class="font-medium text-gray-900">${escape(rel.name)}</h4> <span class="${"px-2 py-1 text-xs rounded-full " + escape(
      rel.relationship_type === "one-to-one" ? "bg-green-100 text-green-800" : rel.relationship_type === "one-to-many" ? "bg-blue-100 text-blue-800" : "bg-yellow-100 text-yellow-800",
      true
    )}">${escape(rel.relationship_type)} </span></div> <div class="text-sm text-gray-600"><div>${escape(forms.find((f) => f.id === rel.source_form_id)?.name)} → ${escape(forms.find((f) => f.id === rel.target_form_id)?.name)}</div> <div class="mt-1 text-xs">${escape(rel.cascade_delete ? "🗑️ Cascade Delete" : "")} ${escape(rel.cascade_update ? "🔄 Cascade Update" : "")} </div></div> </div>`;
  })}</div>`}</div></div></div></div>` : ``}  ${``}  ${``}`;
});
const AdvancedReportsModal = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { forms = [] } = $$props;
  let { showModal = false } = $$props;
  createEventDispatcher();
  let selectedForm = "";
  let selectedReportType = "summary";
  let dateRange = { start: "", end: "" };
  let filters = [];
  const reportTypes = [
    {
      id: "summary",
      name: "Response Summary",
      description: "Overview of form responses"
    },
    {
      id: "detailed",
      name: "Detailed Responses",
      description: "Individual response details"
    },
    {
      id: "trends",
      name: "Response Trends",
      description: "Response patterns over time"
    },
    {
      id: "completion",
      name: "Completion Rates",
      description: "Form completion analytics"
    },
    {
      id: "field-analysis",
      name: "Field Analysis",
      description: "Analysis of specific fields"
    },
    {
      id: "cross-tab",
      name: "Cross-tabulation",
      description: "Cross-reference multiple fields"
    }
  ];
  const filterTypes = [
    {
      id: "date_range",
      name: "Date Range",
      description: "Filter by submission date"
    },
    {
      id: "field_value",
      name: "Field Value",
      description: "Filter by field responses"
    },
    {
      id: "completion_status",
      name: "Completion Status",
      description: "Complete vs incomplete responses"
    }
  ];
  function getFormFields(formId) {
    const form = forms.find((f) => f.id === formId);
    if (form && form.fields) {
      return form.fields.map((f) => ({
        id: f.id,
        label: f.label || f.name || f.id
      }));
    }
    return [];
  }
  if ($$props.forms === void 0 && $$bindings.forms && forms !== void 0) $$bindings.forms(forms);
  if ($$props.showModal === void 0 && $$bindings.showModal && showModal !== void 0) $$bindings.showModal(showModal);
  {
    if (showModal && !selectedForm && forms.length > 0) {
      selectedForm = forms[0].id;
    }
  }
  return `${showModal ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg w-full max-w-6xl max-h-[90vh] overflow-hidden flex flex-col"> <div class="flex items-center justify-between p-6 border-b border-gray-200"><div data-svelte-h="svelte-1fxcljd"><h2 class="text-2xl font-bold text-gray-900">📊 Advanced Reports</h2> <p class="text-gray-600 mt-1">Generate comprehensive reports and analytics for your forms</p></div> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-1x0912u">×</button></div>  <div class="flex-1 flex overflow-hidden"> <div class="w-80 bg-gray-50 border-r border-gray-200 p-4 overflow-y-auto"><h3 class="text-lg font-semibold text-gray-900 mb-4" data-svelte-h="svelte-12mca9e">Report Configuration</h3> <div class="space-y-4"> <div><label for="report-form" class="block text-sm font-medium text-gray-700 mb-1" data-svelte-h="svelte-dvcuqi">Select Form</label> <select id="report-form" class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500">${each(forms, (form) => {
    return `<option${add_attribute("value", form.id, 0)}>${escape(form.name)}</option>`;
  })}</select></div>  <div><label for="report-type" class="block text-sm font-medium text-gray-700 mb-1" data-svelte-h="svelte-of7vnk">Report Type</label> <select id="report-type" class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500">${each(reportTypes, (type) => {
    return `<option${add_attribute("value", type.id, 0)}>${escape(type.name)}</option>`;
  })}</select> <p class="text-xs text-gray-500 mt-1">${escape(reportTypes.find((t) => t.id === selectedReportType)?.description)}</p></div>  <div class="grid grid-cols-2 gap-2"><div><label for="date-start" class="block text-sm font-medium text-gray-700 mb-1" data-svelte-h="svelte-le96de">Start Date</label> <input id="date-start" type="date" class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"${add_attribute("value", dateRange.start, 0)}></div> <div><label for="date-end" class="block text-sm font-medium text-gray-700 mb-1" data-svelte-h="svelte-1a9hjra">End Date</label> <input id="date-end" type="date" class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"${add_attribute("value", dateRange.end, 0)}></div></div>  <div><div class="flex items-center justify-between mb-2"><label class="block text-sm font-medium text-gray-700" data-svelte-h="svelte-182h465">Filters</label> <button class="px-2 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700" data-svelte-h="svelte-9hr51r">+ Add Filter</button></div> <div class="space-y-2 max-h-40 overflow-y-auto">${each(filters, (filter) => {
    return `<div class="border border-gray-200 rounded p-3 bg-white"><div class="flex items-center justify-between mb-2"><select class="text-sm border border-gray-300 rounded px-2 py-1">${each(filterTypes, (type) => {
      return `<option${add_attribute("value", type.id, 0)}>${escape(type.name)}</option>`;
    })}</select> <button class="text-red-500 hover:text-red-700 text-sm" data-svelte-h="svelte-1raz5va">✕
											</button></div> ${filter.type === "field_value" ? `<div class="grid grid-cols-3 gap-2"><select class="text-sm border border-gray-300 rounded px-2 py-1"><option value="" data-svelte-h="svelte-pel8ws">Field</option>${each(getFormFields(selectedForm), (field) => {
      return `<option${add_attribute("value", field.id, 0)}>${escape(field.label)}</option>`;
    })}</select> <select class="text-sm border border-gray-300 rounded px-2 py-1"><option value="equals" data-svelte-h="svelte-e1mtai">Equals</option><option value="contains" data-svelte-h="svelte-3z0d3y">Contains</option><option value="greater_than" data-svelte-h="svelte-9v2ja5">Greater Than</option><option value="less_than" data-svelte-h="svelte-15v009j">Less Than</option></select> <input type="text" placeholder="Value" class="text-sm border border-gray-300 rounded px-2 py-1"${add_attribute("value", filter.value, 0)}> </div>` : ``} </div>`;
  })}</div></div>  <button class="w-full px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50" ${!selectedForm ? "disabled" : ""}>${escape("Generate Report")}</button></div></div>  <div class="flex-1 p-6 overflow-y-auto">${` <div class="text-center py-16" data-svelte-h="svelte-1fxhkry"><div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4"><svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg></div> <h3 class="text-xl font-medium text-gray-900 mb-2">Generate Your First Report</h3> <p class="text-gray-500 mb-6">Select a form and report type, then click &quot;Generate Report&quot; to get started.</p> <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-left"><div class="bg-blue-50 p-4 rounded-lg"><h4 class="font-medium text-blue-900 mb-2">📊 Summary Reports</h4> <p class="text-sm text-blue-700">Overview statistics and completion rates</p></div> <div class="bg-green-50 p-4 rounded-lg"><h4 class="font-medium text-green-900 mb-2">📋 Detailed Responses</h4> <p class="text-sm text-green-700">Individual response data and filtering</p></div> <div class="bg-purple-50 p-4 rounded-lg"><h4 class="font-medium text-purple-900 mb-2">📈 Analytics</h4> <p class="text-sm text-purple-700">Trends, patterns, and insights</p></div></div></div>`}</div></div></div></div>` : ``}`;
});
const WorkflowDesignerModal = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { forms = [] } = $$props;
  let { showModal = false } = $$props;
  createEventDispatcher();
  let canvas = null;
  let ctx = null;
  let workflow = {
    name: "",
    trigger: "on_submit",
    isActive: false,
    steps: []
  };
  let selectedStep = null;
  const stepTypes = [
    {
      id: "approval",
      name: "Approval",
      icon: "✅",
      description: "Require approval from a user"
    },
    {
      id: "notification",
      name: "Notification",
      icon: "📧",
      description: "Send email or in-app notification"
    },
    {
      id: "assignment",
      name: "Assignment",
      icon: "👤",
      description: "Assign task to a user"
    },
    {
      id: "condition",
      name: "Condition",
      icon: "🔀",
      description: "Branch workflow based on condition"
    },
    {
      id: "action",
      name: "Action",
      icon: "⚡",
      description: "Execute custom action"
    }
  ];
  const triggerTypes = [
    {
      id: "on_submit",
      name: "On Form Submit",
      description: "When a form is submitted"
    },
    {
      id: "on_update",
      name: "On Form Update",
      description: "When a form response is updated"
    },
    {
      id: "on_approve",
      name: "On Approval",
      description: "When an approval step completes"
    },
    {
      id: "scheduled",
      name: "Scheduled",
      description: "Run on a schedule"
    }
  ];
  async function initializeCanvas() {
    await new Promise((resolve) => setTimeout(resolve, 100));
    canvas = document.getElementById("workflow-canvas");
    if (canvas) {
      ctx = canvas.getContext("2d");
      resizeCanvas();
      drawCanvas();
    }
  }
  function resizeCanvas() {
    if (canvas) {
      const container = canvas.parentElement;
      canvas.width = container.clientWidth;
      canvas.height = container.clientHeight;
    }
  }
  function drawCanvas() {
    if (!ctx || !canvas) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    workflow.steps.forEach((step) => {
      if (step.nextStepId) {
        const nextStep = workflow.steps.find((s) => s.id === step.nextStepId);
        if (nextStep) {
          drawConnection(step, nextStep, "#10B981", "next");
        }
      }
      if (step.altStepId) {
        const altStep = workflow.steps.find((s) => s.id === step.altStepId);
        if (altStep) {
          drawConnection(step, altStep, "#EF4444", "alt");
        }
      }
    });
    workflow.steps.forEach((step) => drawStep(step));
  }
  function drawStep(step) {
    if (!ctx) return;
    const stepType = stepTypes.find((t) => t.id === step.type);
    const isSelected = selectedStep?.id === step.id;
    ctx.fillStyle = isSelected ? "#3B82F6" : "#FFFFFF";
    ctx.strokeStyle = isSelected ? "#1D4ED8" : "#D1D5DB";
    ctx.lineWidth = isSelected ? 3 : 2;
    const radius = 8;
    ctx.beginPath();
    ctx.roundRect(step.x, step.y, 160, 80, radius);
    ctx.fill();
    ctx.stroke();
    ctx.fillStyle = isSelected ? "#FFFFFF" : "#374151";
    ctx.font = "bold 12px Arial";
    ctx.fillText(stepType?.icon || "⚡", step.x + 10, step.y + 25);
    ctx.font = "12px Arial";
    ctx.fillText(step.name || stepType?.name || "Step", step.x + 35, step.y + 25);
    ctx.fillStyle = isSelected ? "#E0F2FE" : "#6B7280";
    ctx.font = "10px Arial";
    ctx.fillText(step.type.toUpperCase(), step.x + 10, step.y + 45);
    ctx.fillStyle = "#10B981";
    ctx.beginPath();
    ctx.arc(step.x + 160, step.y + 25, 6, 0, 2 * Math.PI);
    ctx.fill();
    if (step.type === "condition" || step.type === "approval") {
      ctx.fillStyle = "#EF4444";
      ctx.beginPath();
      ctx.arc(step.x + 160, step.y + 55, 6, 0, 2 * Math.PI);
      ctx.fill();
    }
  }
  function drawConnection(fromStep, toStep, color, type) {
    if (!ctx) return;
    const startX = fromStep.x + 160;
    const startY = fromStep.y + (type === "alt" ? 55 : 25);
    const endX = toStep.x;
    const endY = toStep.y + 40;
    ctx.strokeStyle = color;
    ctx.lineWidth = 3;
    ctx.beginPath();
    ctx.moveTo(startX, startY);
    const midX = (startX + endX) / 2;
    const midY = Math.min(startY, endY) - 30;
    ctx.quadraticCurveTo(midX, midY, endX, endY);
    ctx.stroke();
    const angle = Math.atan2(endY - midY, endX - midX);
    const arrowLength = 15;
    ctx.beginPath();
    ctx.moveTo(endX, endY);
    ctx.lineTo(endX - arrowLength * Math.cos(angle - Math.PI / 6), endY - arrowLength * Math.sin(angle - Math.PI / 6));
    ctx.moveTo(endX, endY);
    ctx.lineTo(endX - arrowLength * Math.cos(angle + Math.PI / 6), endY - arrowLength * Math.sin(angle + Math.PI / 6));
    ctx.stroke();
    ctx.fillStyle = color;
    ctx.font = "bold 10px Arial";
    const labelX = midX - 10;
    const labelY = midY - 10;
    ctx.fillText(type === "alt" ? "NO/REJECT" : "YES/APPROVE", labelX, labelY);
  }
  if ($$props.forms === void 0 && $$bindings.forms && forms !== void 0) $$bindings.forms(forms);
  if ($$props.showModal === void 0 && $$bindings.showModal && showModal !== void 0) $$bindings.showModal(showModal);
  {
    if (showModal && !canvas) {
      initializeCanvas();
    }
  }
  return `${showModal ? `<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"><div class="bg-white rounded-lg w-full max-w-7xl max-h-[90vh] overflow-hidden flex flex-col"> <div class="flex items-center justify-between p-6 border-b border-gray-200"><div data-svelte-h="svelte-1pc4cch"><h2 class="text-2xl font-bold text-gray-900">⚡ Workflow Designer</h2> <p class="text-gray-600 mt-1">Create automated approval chains and business processes</p></div> <button class="text-gray-400 hover:text-gray-600 text-2xl" data-svelte-h="svelte-1x0912u">×</button></div>  <div class="flex items-center justify-between p-4 bg-gray-50 border-b border-gray-200"><div class="flex items-center space-x-3"> ${each(stepTypes, (stepType) => {
    return `<button class="px-3 py-2 bg-white border border-gray-300 rounded hover:bg-gray-50 text-sm flex items-center space-x-2"><span>${escape(stepType.icon)}</span> <span>${escape(stepType.name)}</span> </button>`;
  })}</div> <div class="flex items-center space-x-4">${``} <button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700" ${"disabled"}>💾 Save Workflow</button></div></div>  <div class="flex-1 flex overflow-hidden"> <div class="flex-1 relative bg-gray-100"><canvas id="workflow-canvas" class="w-full h-full cursor-move"></canvas> ${workflow.steps.length === 0 ? `<div class="absolute inset-0 flex items-center justify-center" data-svelte-h="svelte-1qrs0bh"><div class="text-center"><div class="w-24 h-24 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-4"><svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg></div> <h3 class="text-xl font-medium text-gray-900 mb-2">Start Building Your Workflow</h3> <p class="text-gray-500 mb-6">Add steps from the toolbar above and connect them to create automated processes.</p> <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-left"><div class="bg-blue-50 p-4 rounded-lg"><h4 class="font-medium text-blue-900 mb-2">✅ Approvals</h4> <p class="text-sm text-blue-700">Require sign-off from managers or team members</p></div> <div class="bg-green-50 p-4 rounded-lg"><h4 class="font-medium text-green-900 mb-2">📧 Notifications</h4> <p class="text-sm text-green-700">Send emails or in-app alerts automatically</p></div> <div class="bg-purple-50 p-4 rounded-lg"><h4 class="font-medium text-purple-900 mb-2">🔀 Conditions</h4> <p class="text-sm text-purple-700">Branch workflows based on form responses</p></div></div></div></div>` : ``}  ${workflow.steps.length > 0 && !selectedStep ? `<div class="absolute bottom-4 left-4 bg-white p-4 rounded-lg shadow-lg border border-gray-200 max-w-sm" data-svelte-h="svelte-io12vn"><h4 class="font-medium text-gray-900 mb-2">How to Use:</h4> <ul class="text-sm text-gray-600 space-y-1"><li>• Click and drag steps to reposition them</li> <li>• Click connection points (colored dots) to link steps</li> <li>• Select a step to configure its properties</li> <li>• Green connections = approval path</li> <li>• Red connections = rejection/alternative path</li></ul></div>` : ``}</div>  <div class="w-80 bg-white border-l border-gray-200 p-4 overflow-y-auto">${` <div><h3 class="text-lg font-semibold text-gray-900 mb-4" data-svelte-h="svelte-12ee41u">Workflow Settings</h3> <div class="space-y-4"><div><label for="workflow-name" class="block text-sm font-medium text-gray-700 mb-1" data-svelte-h="svelte-3sv0mg">Workflow Name</label> <input id="workflow-name" type="text" placeholder="Enter workflow name" class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"${add_attribute("value", workflow.name, 0)}></div> <div><label for="workflow-description" class="block text-sm font-medium text-gray-700 mb-1" data-svelte-h="svelte-1pzymp3">Description</label> <textarea id="workflow-description" placeholder="Describe this workflow" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500">${escape("")}</textarea></div> <div><label for="form-select" class="block text-sm font-medium text-gray-700 mb-1" data-svelte-h="svelte-10diir1">Target Form</label> <select id="form-select" class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500"><option value="" data-svelte-h="svelte-17oez36">Select form</option>${each(forms, (form) => {
    return `<option${add_attribute("value", form.id, 0)}>${escape(form.name)}</option>`;
  })}</select></div> <div><label for="trigger-select" class="block text-sm font-medium text-gray-700 mb-1" data-svelte-h="svelte-ygud1w">Trigger</label> <select id="trigger-select" class="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500">${each(triggerTypes, (trigger) => {
    return `<option${add_attribute("value", trigger.id, 0)}>${escape(trigger.name)}</option>`;
  })}</select> <p class="text-xs text-gray-500 mt-1">${escape(triggerTypes.find((t) => t.id === workflow.trigger)?.description)}</p></div> <div class="flex items-center"><input type="checkbox" id="workflow-active" class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"${add_attribute("checked", workflow.isActive, 1)}> <label for="workflow-active" class="ml-2 text-sm text-gray-700" data-svelte-h="svelte-1y2xp7b">Activate workflow</label></div></div></div>`}</div></div></div></div>` : ``}`;
});
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  const params = null;
  const data = null;
  const form = null;
  let forms = [];
  let showDatabaseRelationsModal = false;
  let showAdvancedReportsModal = false;
  let showWorkflowDesignerModal = false;
  if ($$props.params === void 0 && $$bindings.params && params !== void 0) $$bindings.params(params);
  if ($$props.data === void 0 && $$bindings.data && data !== void 0) $$bindings.data(data);
  if ($$props.form === void 0 && $$bindings.form && form !== void 0) $$bindings.form(form);
  return `${$$result.head += `<!-- HEAD_svelte-1wh74jd_START -->${$$result.title = `<title>Forms &amp; Templates - TPT Titan</title>`, ""}<!-- HEAD_svelte-1wh74jd_END -->`, ""} <div class="h-screen flex flex-col bg-gray-50"> <header class="flex items-center justify-between px-6 py-3 border-b border-gray-200 bg-white"><div class="flex items-center space-x-4" data-svelte-h="svelte-16kmn8v"><h1 class="text-xl font-semibold text-gray-900">Forms &amp; Templates</h1> <span class="text-sm text-gray-500">MS Access-style database features</span></div> <div class="flex items-center space-x-2">${`<button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors" data-svelte-h="svelte-2ndb8s">Create Form</button>`}</div></header>  <div class="flex-1 overflow-hidden">${`${validate_component(FormList, "FormList").$$render($$result, { forms }, {}, {})}`}</div></div>  ${validate_component(DatabaseRelationsModal, "DatabaseRelationsModal").$$render(
    $$result,
    {
      forms,
      showModal: showDatabaseRelationsModal
    },
    {},
    {}
  )}  ${validate_component(AdvancedReportsModal, "AdvancedReportsModal").$$render(
    $$result,
    {
      forms,
      showModal: showAdvancedReportsModal
    },
    {},
    {}
  )}  ${validate_component(WorkflowDesignerModal, "WorkflowDesignerModal").$$render(
    $$result,
    {
      forms,
      showModal: showWorkflowDesignerModal
    },
    {},
    {}
  )}`;
});
export {
  Page as default
};
