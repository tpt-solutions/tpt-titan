<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { api } from '$lib/api.js';

  export let workflowId = null; // null for new workflow, ID for editing
  export let initialData = null;

  let canvasRef;
  let selectedNode = null;
  let draggedNode = null;
  let isConnecting = false;
  let connectionStart = null;
  let tempConnection = null;

  // Workflow data
  let workflow = {
    id: null,
    name: '',
    description: '',
    category: '',
    triggerType: 'manual',
    schedule: '',
    isActive: true,
    canvasData: {
      nodes: [],
      connections: [],
      viewport: { x: 0, y: 0, zoom: 1 }
    }
  };

  // Available node types
  const nodeTypes = [
    {
      type: 'trigger',
      label: 'Trigger',
      icon: '🎯',
      color: '#4CAF50',
      category: 'triggers',
      description: 'Workflow starting point'
    },
    {
      type: 'action',
      label: 'Action',
      icon: '⚡',
      color: '#2196F3',
      category: 'actions',
      description: 'Execute an action'
    },
    {
      type: 'condition',
      label: 'Condition',
      icon: '🔀',
      color: '#FF9800',
      category: 'logic',
      description: 'Conditional logic'
    },
    {
      type: 'delay',
      label: 'Delay',
      icon: '⏱️',
      color: '#9C27B0',
      category: 'logic',
      description: 'Add time delay'
    },
    {
      type: 'notification',
      label: 'Notification',
      icon: '🔔',
      color: '#607D8B',
      category: 'actions',
      description: 'Send notification'
    }
  ];

  // Available connectors
  let connectors = [];

  const dispatch = createEventDispatcher();

  onMount(async () => {
    // Load available connectors
    try {
      const response = await api.get('/workflows/connectors');
      connectors = response.connectors || [];
    } catch (error) {
      console.error('Failed to load connectors:', error);
    }

    // Load existing workflow if editing
    if (workflowId) {
      await loadWorkflow(workflowId);
    } else if (initialData) {
      workflow = { ...workflow, ...initialData };
      if (workflow.canvasData && typeof workflow.canvasData === 'string') {
        workflow.canvasData = JSON.parse(workflow.canvasData);
      }
    }
  });

  async function loadWorkflow(id) {
    try {
      const response = await api.get(`/workflows/${id}`);
      workflow = response.workflow;
      if (workflow.canvasData && typeof workflow.canvasData === 'string') {
        workflow.canvasData = JSON.parse(workflow.canvasData);
      }
    } catch (error) {
      console.error('Failed to load workflow:', error);
      dispatch('error', { message: 'Failed to load workflow' });
    }
  }

  function addNode(nodeType, position = { x: 100, y: 100 }) {
    const node = {
      id: `node_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      type: nodeType.type,
      name: nodeType.label,
      position: { ...position },
      config: getDefaultConfig(nodeType.type),
      width: 200,
      height: 80
    };

    workflow.canvasData.nodes = [...workflow.canvasData.nodes, node];
    workflow.canvasData = { ...workflow.canvasData }; // Trigger reactivity
  }

  function getDefaultConfig(nodeType) {
    switch (nodeType) {
      case 'trigger':
        return {
          triggerType: 'manual',
          conditions: []
        };
      case 'action':
        return {
          connector: '',
          action: '',
          parameters: {}
        };
      case 'condition':
        return {
          field: '',
          operator: 'equals',
          value: ''
        };
      case 'delay':
        return {
          delay_seconds: 60
        };
      case 'notification':
        return {
          title: '',
          message: '',
          type: 'info'
        };
      default:
        return {};
    }
  }

  function updateNode(nodeId, updates) {
    const nodeIndex = workflow.canvasData.nodes.findIndex(n => n.id === nodeId);
    if (nodeIndex !== -1) {
      workflow.canvasData.nodes[nodeIndex] = {
        ...workflow.canvasData.nodes[nodeIndex],
        ...updates
      };
      workflow.canvasData = { ...workflow.canvasData };
    }
  }

  function deleteNode(nodeId) {
    workflow.canvasData.nodes = workflow.canvasData.nodes.filter(n => n.id !== nodeId);
    workflow.canvasData.connections = workflow.canvasData.connections.filter(
      c => c.from !== nodeId && c.to !== nodeId
    );
    workflow.canvasData = { ...workflow.canvasData };
  }

  function selectNode(node) {
    selectedNode = node;
  }

  function startConnection(nodeId, port) {
    isConnecting = true;
    connectionStart = { nodeId, port };
  }

  function completeConnection(nodeId, port) {
    if (!isConnecting || !connectionStart) return;

    // Don't connect to self
    if (connectionStart.nodeId === nodeId) return;

    const connection = {
      id: `conn_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      from: connectionStart.nodeId,
      to: nodeId,
      fromPort: connectionStart.port,
      toPort: port
    };

    workflow.canvasData.connections = [...workflow.canvasData.connections, connection];
    workflow.canvasData = { ...workflow.canvasData };

    isConnecting = false;
    connectionStart = null;
    tempConnection = null;
  }

  function deleteConnection(connectionId) {
    workflow.canvasData.connections = workflow.canvasData.connections.filter(
      c => c.id !== connectionId
    );
    workflow.canvasData = { ...workflow.canvasData };
  }

  function handleCanvasClick(event) {
    if (isConnecting) {
      isConnecting = false;
      connectionStart = null;
      tempConnection = null;
      return;
    }

    selectedNode = null;
  }

  function handleCanvasMouseMove(event) {
    if (isConnecting && connectionStart) {
      const rect = canvasRef.getBoundingClientRect();
      const x = event.clientX - rect.left - workflow.canvasData.viewport.x;
      const y = event.clientY - rect.top - workflow.canvasData.viewport.y;

      tempConnection = {
        from: connectionStart.nodeId,
        to: { x, y }
      };
    }
  }

  async function saveWorkflow() {
    try {
      const workflowData = {
        ...workflow,
        canvasData: JSON.stringify(workflow.canvasData)
      };

      let response;
      if (workflow.id) {
        response = await api.put(`/workflows/${workflow.id}`, workflowData);
      } else {
        response = await api.post('/workflows', workflowData);
      }

      workflow.id = response.workflow.id;
      dispatch('saved', { workflow: response.workflow });
    } catch (error) {
      console.error('Failed to save workflow:', error);
      dispatch('error', { message: 'Failed to save workflow' });
    }
  }

  async function executeWorkflow() {
    if (!workflow.id) {
      await saveWorkflow();
    }

    try {
      const response = await api.post(`/workflows/${workflow.id}/execute`);
      dispatch('executed', { execution: response });
    } catch (error) {
      console.error('Failed to execute workflow:', error);
      dispatch('error', { message: 'Failed to execute workflow' });
    }
  }

  function exportWorkflow() {
    const dataStr = JSON.stringify(workflow, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `${workflow.name || 'workflow'}.json`;
    link.click();
    URL.revokeObjectURL(url);
  }

  function importWorkflow(event) {
    const file = event.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const importedWorkflow = JSON.parse(e.target.result);
        workflow = { ...workflow, ...importedWorkflow };
        if (typeof workflow.canvasData === 'string') {
          workflow.canvasData = JSON.parse(workflow.canvasData);
        }
        dispatch('imported', { workflow });
      } catch (error) {
        console.error('Failed to import workflow:', error);
        dispatch('error', { message: 'Invalid workflow file' });
      }
    };
    reader.readAsText(file);
  }

  // Node positioning calculations
  $: nodePositions = workflow.canvasData.nodes.map(node => ({
    ...node,
    screenX: node.position.x + workflow.canvasData.viewport.x,
    screenY: node.position.y + workflow.canvasData.viewport.y
  }));
</script>

<div class="workflow-builder">
  <!-- Toolbar -->
  <div class="toolbar">
    <div class="workflow-info">
      <input
        type="text"
        bind:value={workflow.name}
        placeholder="Workflow name"
        class="workflow-name"
      />
      <select bind:value={workflow.category} class="category-select">
        <option value="">Select category</option>
        <option value="business">Business</option>
        <option value="personal">Personal</option>
        <option value="automation">Automation</option>
        <option value="integration">Integration</option>
      </select>
    </div>

    <div class="toolbar-actions">
      <button on:click={saveWorkflow} class="btn btn-primary">
        💾 Save
      </button>
      <button on:click={executeWorkflow} class="btn btn-success">
        ▶️ Execute
      </button>
      <button on:click={exportWorkflow} class="btn btn-secondary">
        📤 Export
      </button>
      <label class="btn btn-secondary">
        📥 Import
        <input
          type="file"
          accept=".json"
          on:change={importWorkflow}
          style="display: none;"
        />
      </label>
    </div>
  </div>

  <div class="builder-content">
    <!-- Node Palette -->
    <div class="node-palette">
      <h4>Node Types</h4>
      {#each nodeTypes as nodeType}
        <div
          class="node-palette-item"
          draggable="true"
          on:dragstart={() => draggedNode = nodeType}
          style="background-color: {nodeType.color}20; border-color: {nodeType.color};"
        >
          <span class="node-icon">{nodeType.icon}</span>
          <div class="node-info">
            <div class="node-label">{nodeType.label}</div>
            <div class="node-description">{nodeType.description}</div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Canvas -->
    <div class="canvas-container">
      <div
        class="canvas"
        bind:this={canvasRef}
        on:click={handleCanvasClick}
        on:mousemove={handleCanvasMouseMove}
        on:dragover={(e) => e.preventDefault()}
        on:drop={(e) => {
          e.preventDefault();
          if (draggedNode) {
            const rect = canvasRef.getBoundingClientRect();
            const x = e.clientX - rect.left - workflow.canvasData.viewport.x;
            const y = e.clientY - rect.top - workflow.canvasData.viewport.y;
            addNode(draggedNode, { x, y });
            draggedNode = null;
          }
        }}
      >
        <!-- Grid background -->
        <div class="canvas-grid"></div>

        <!-- Connections -->
        <svg class="connections-layer">
          {#each workflow.canvasData.connections as connection}
            <WorkflowConnection
              {connection}
              nodes={nodePositions}
              on:delete={() => deleteConnection(connection.id)}
            />
          {/each}

          <!-- Temporary connection while dragging -->
          {#if tempConnection}
            <line
              x1={nodePositions.find(n => n.id === tempConnection.from)?.screenX + 100 || 0}
              y1={nodePositions.find(n => n.id === tempConnection.from)?.screenY + 40 || 0}
              x2={tempConnection.to.x}
              y2={tempConnection.to.y}
              stroke="#666"
              stroke-width="2"
              stroke-dasharray="5,5"
            />
          {/if}
        </svg>

        <!-- Nodes -->
        {#each nodePositions as node}
          <WorkflowNode
            {node}
            {selectedNode}
            {connectors}
            on:select={() => selectNode(node)}
            on:update={(e) => updateNode(node.id, e.detail)}
            on:delete={() => deleteNode(node.id)}
            on:startConnection={(e) => startConnection(node.id, e.detail)}
            on:completeConnection={(e) => completeConnection(node.id, e.detail)}
          />
        {/each}
      </div>
    </div>

    <!-- Properties Panel -->
    {#if selectedNode}
      <div class="properties-panel">
        <h4>Node Properties</h4>
        <WorkflowNodeProperties
          node={selectedNode}
          {connectors}
          on:update={(e) => updateNode(selectedNode.id, e.detail)}
        />
      </div>
    {/if}
  </div>
</div>

<style>
  .workflow-builder {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: #f5f5f5;
  }

  .toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: white;
    border-bottom: 1px solid #e0e0e0;
    gap: 1rem;
  }

  .workflow-info {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .workflow-name {
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 1.1rem;
    font-weight: 500;
  }

  .category-select {
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }

  .toolbar-actions {
    display: flex;
    gap: 0.5rem;
  }

  .btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    transition: all 0.2s;
  }

  .btn-primary {
    background: #007bff;
    color: white;
  }

  .btn-primary:hover {
    background: #0056b3;
  }

  .btn-success {
    background: #28a745;
    color: white;
  }

  .btn-success:hover {
    background: #1e7e34;
  }

  .btn-secondary {
    background: #6c757d;
    color: white;
  }

  .btn-secondary:hover {
    background: #545b62;
  }

  .builder-content {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .node-palette {
    width: 250px;
    background: white;
    border-right: 1px solid #e0e0e0;
    padding: 1rem;
    overflow-y: auto;
  }

  .node-palette h4 {
    margin: 0 0 1rem 0;
    color: #333;
  }

  .node-palette-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    margin-bottom: 0.5rem;
    border: 2px solid transparent;
    border-radius: 6px;
    cursor: grab;
    transition: all 0.2s;
  }

  .node-palette-item:hover {
    border-color: #007bff;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  .node-icon {
    font-size: 1.5rem;
  }

  .node-info {
    flex: 1;
  }

  .node-label {
    font-weight: 500;
    color: #333;
  }

  .node-description {
    font-size: 0.8rem;
    color: #666;
  }

  .canvas-container {
    flex: 1;
    position: relative;
    overflow: hidden;
  }

  .canvas {
    width: 100%;
    height: 100%;
    position: relative;
    background: white;
    cursor: crosshair;
  }

  .canvas-grid {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-image:
      linear-gradient(#e0e0e0 1px, transparent 1px),
      linear-gradient(90deg, #e0e0e0 1px, transparent 1px);
    background-size: 20px 20px;
  }

  .connections-layer {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: 1;
  }

  .properties-panel {
    width: 300px;
    background: white;
    border-left: 1px solid #e0e0e0;
    padding: 1rem;
    overflow-y: auto;
  }

  .properties-panel h4 {
    margin: 0 0 1rem 0;
    color: #333;
  }

  /* Responsive adjustments */
  @media (max-width: 1200px) {
    .properties-panel {
      width: 250px;
    }
  }

  @media (max-width: 768px) {
    .builder-content {
      flex-direction: column;
    }

    .node-palette {
      width: 100%;
      height: 200px;
      border-right: none;
      border-bottom: 1px solid #e0e0e0;
    }

    .properties-panel {
      width: 100%;
      height: 300px;
      border-left: none;
      border-top: 1px solid #e0e0e0;
    }
  }
</style>
