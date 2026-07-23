<script>
// @ts-nocheck
  import { createEventDispatcher } from 'svelte';

  export let node;
  export let selectedNode;
  export let connectors = [];

  let isDragging = false;
  let dragOffset = { x: 0, y: 0 };
  let nodeElement;

  const dispatch = createEventDispatcher();

  // Node type configurations
  const nodeConfigs = {
    trigger: {
      color: '#4CAF50',
      icon: '🎯',
      inputs: [],
      outputs: ['output']
    },
    action: {
      color: '#2196F3',
      icon: '⚡',
      inputs: ['input'],
      outputs: ['output']
    },
    condition: {
      color: '#FF9800',
      icon: '🔀',
      inputs: ['input'],
      outputs: ['true', 'false']
    },
    delay: {
      color: '#9C27B0',
      icon: '⏱️',
      inputs: ['input'],
      outputs: ['output']
    },
    notification: {
      color: '#607D8B',
      icon: '🔔',
      inputs: ['input'],
      outputs: ['output']
    }
  };

  $: config = nodeConfigs[node.type] || nodeConfigs.action;
  $: isSelected = selectedNode && selectedNode.id === node.id;

  function handleMouseDown(event) {
    if (event.target.classList.contains('connection-point')) return;

    isDragging = true;
    const rect = nodeElement.getBoundingClientRect();
    dragOffset.x = event.clientX - rect.left;
    dragOffset.y = event.clientY - rect.top;

    dispatch('select');

    event.preventDefault();
  }

  function handleMouseMove(event) {
    if (!isDragging) return;

    const canvas = nodeElement.closest('.canvas');
    const canvasRect = canvas.getBoundingClientRect();

    const newX = event.clientX - canvasRect.left - dragOffset.x;
    const newY = event.clientY - canvasRect.top - dragOffset.y;

    dispatch('update', {
      position: { x: newX, y: newY }
    });
  }

  function handleMouseUp() {
    isDragging = false;
  }

  function handleConnectionStart(port, event) {
    event.stopPropagation();
    dispatch('startConnection', port);
  }

  function handleConnectionEnd(port, event) {
    event.stopPropagation();
    dispatch('completeConnection', port);
  }

  function deleteNode(event) {
    event.stopPropagation();
    dispatch('delete');
  }

  // Global mouse event listeners for dragging
  if (typeof window !== 'undefined') {
    window.addEventListener('mousemove', handleMouseMove);
    window.addEventListener('mouseup', handleMouseUp);
  }
</script>

<div
  class="workflow-node"
  class:selected={isSelected}
  bind:this={nodeElement}
  style="
    left: {node.screenX || node.position.x}px;
    top: {node.screenY || node.position.y}px;
    border-color: {config.color};
    background-color: {config.color}15;
  "
  on:mousedown={handleMouseDown}
>
  <!-- Node Header -->
  <div class="node-header">
    <span class="node-icon">{config.icon}</span>
    <span class="node-title">{node.name || node.type}</span>
    <button class="delete-btn" on:click={deleteNode}>×</button>
  </div>

  <!-- Connection Points -->
  <div class="connection-points">
    <!-- Input points -->
    {#each config.inputs as input}
      <div
        class="connection-point input-point"
        data-port={input}
        on:mousedown={(e) => handleConnectionStart(input, e)}
        on:mouseup={(e) => handleConnectionEnd(input, e)}
      ></div>
    {/each}

    <!-- Output points -->
    {#each config.outputs as output}
      <div
        class="connection-point output-point"
        class:true-port={output === 'true'}
        class:false-port={output === 'false'}
        data-port={output}
        on:mousedown={(e) => handleConnectionStart(output, e)}
        on:mouseup={(e) => handleConnectionEnd(output, e)}
      >
        {#if node.type === 'condition'}
          <span class="port-tag">{output === 'true' ? 'T' : output === 'false' ? 'F' : ''}</span>
        {/if}
      </div>
    {/each}
  </div>

  <!-- Node Content Preview -->
  <div class="node-content">
    {#if node.type === 'action' && node.config?.connector}
      <div class="connector-preview">
        {connectors.find(c => c.id === node.config.connector)?.name || 'Select connector'}
      </div>
    {:else if node.type === 'condition'}
      <div class="condition-preview">
        {node.config?.field || 'field'} {node.config?.operator || '=='} {node.config?.value || 'value'}
      </div>
    {:else if node.type === 'delay'}
      <div class="delay-preview">
        {node.config?.delay_seconds || 60}s delay
      </div>
    {:else if node.type === 'notification'}
      <div class="notification-preview">
        {node.config?.title || 'Notification'}
      </div>
    {:else}
      <div class="default-preview">
        Configure node settings
      </div>
    {/if}
  </div>
</div>

<style>
  .workflow-node {
    position: absolute;
    min-width: 200px;
    min-height: 80px;
    background: white;
    border: 2px solid #ddd;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    cursor: move;
    user-select: none;
    transition: all 0.2s;
    z-index: 10;
  }

  .workflow-node:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .workflow-node.selected {
    border-color: #007bff;
    box-shadow: 0 4px 16px rgba(0, 123, 255, 0.2);
  }

  .node-header {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    background: rgba(255, 255, 255, 0.9);
    border-bottom: 1px solid #eee;
    border-radius: 6px 6px 0 0;
  }

  .node-icon {
    font-size: 1.2rem;
    margin-right: 8px;
  }

  .node-title {
    flex: 1;
    font-weight: 500;
    font-size: 0.9rem;
    color: #333;
  }

  .delete-btn {
    background: none;
    border: none;
    color: #999;
    cursor: pointer;
    font-size: 1.2rem;
    padding: 0;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.2s;
  }

  .delete-btn:hover {
    background: #ff4757;
    color: white;
  }

  .connection-points {
    position: relative;
    height: 0;
  }

  .connection-point {
    position: absolute;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    border: 2px solid #666;
    background: white;
    cursor: pointer;
    transition: all 0.2s;
    z-index: 20;
  }

  .connection-point:hover {
    transform: scale(1.2);
    border-color: #007bff;
  }

  .input-point {
    left: -7px;
    top: 50%;
    transform: translateY(-50%);
  }

  .output-point {
    right: -7px;
    top: 50%;
    transform: translateY(-50%);
  }

  .output-point.true-port {
    border-color: #4CAF50;
    background: #e8f5e9;
  }

  .output-point.false-port {
    border-color: #f44336;
    background: #ffebee;
  }

  .port-tag {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 9px;
    font-weight: 700;
    line-height: 1;
    pointer-events: none;
  }

  .true-port .port-tag {
    color: #2e7d32;
  }

  .false-port .port-tag {
    color: #c62828;
  }

  .node-content {
    padding: 8px 12px;
    font-size: 0.8rem;
    color: #666;
  }

  .connector-preview {
    font-style: italic;
  }

  .condition-preview {
    font-family: monospace;
    background: rgba(255, 152, 0, 0.1);
    padding: 2px 4px;
    border-radius: 3px;
  }

  .delay-preview {
    color: #9C27B0;
    font-weight: 500;
  }

  .notification-preview {
    color: #607D8B;
    font-weight: 500;
  }

  .default-preview {
    font-style: italic;
    color: #999;
  }

  /* Multiple connection points layout */
  .workflow-node[data-inputs="2"] .input-point:nth-child(1) { top: 25%; }
  .workflow-node[data-inputs="2"] .input-point:nth-child(2) { top: 75%; }

  .workflow-node[data-outputs="2"] .output-point:nth-child(1) { top: 25%; }
  .workflow-node[data-outputs="2"] .output-point:nth-child(2) { top: 75%; }
</style>
