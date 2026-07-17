<script>
  import { createEventDispatcher } from 'svelte';

  export let connection;
  export let nodes = [];

  const dispatch = createEventDispatcher();

  // Calculate connection path between nodes
  $: fromNode = nodes.find(n => n.id === connection.from);
  $: toNode = nodes.find(n => n.id === connection.to);

  $: path = calculatePath(fromNode, toNode, connection);

  function calculatePath(fromNode, toNode, connection) {
    if (!fromNode || !toNode) return '';

    // Calculate connection points
    const fromX = (fromNode.screenX || fromNode.position.x) + fromNode.width;
    const fromY = (fromNode.screenY || fromNode.position.y) + (fromNode.height / 2);
    const toX = (toNode.screenX || toNode.position.x);
    const toY = (toNode.screenY || toNode.position.y) + (toNode.height / 2);

    // Create curved path
    const midX = (fromX + toX) / 2;
    const dx = Math.abs(toX - fromX);
    const curvature = Math.min(dx * 0.3, 100); // Adaptive curvature

    return `M ${fromX} ${fromY} C ${fromX + curvature} ${fromY}, ${toX - curvature} ${toY}, ${toX} ${toY}`;
  }

  function handleClick(event) {
    event.stopPropagation();
    // Could add connection selection here
  }

  function handleDelete(event) {
    event.stopPropagation();
    dispatch('delete');
  }
</script>

<g class="workflow-connection" on:click={handleClick}>
  <!-- Connection path -->
  <path
    d={path}
    stroke="#666"
    stroke-width="2"
    fill="none"
    marker-end="url(#arrowhead)"
    class="connection-path"
  />

  <!-- Invisible wider path for easier clicking -->
  <path
    d={path}
    stroke="transparent"
    stroke-width="10"
    fill="none"
    class="connection-hit-area"
  />

  <!-- Delete button (shown on hover) -->
  {#if path}
    {@const midPoint = getMidPoint(fromNode, toNode)}
    {#if connection.fromPort && connection.fromPort !== 'output'}
      <rect
        x={midPoint.x - 18}
        y={midPoint.y - 24}
        width="36"
        height="16"
        rx="8"
        fill={connection.fromPort === 'true' ? '#e8f5e9' : '#ffebee'}
        stroke={connection.fromPort === 'true' ? '#4CAF50' : '#f44336'}
        stroke-width="1"
        class="port-label-bg"
      />
      <text
        x={midPoint.x}
        y={midPoint.y - 16}
        text-anchor="middle"
        dominant-baseline="middle"
        font-size="10"
        font-weight="600"
        fill={connection.fromPort === 'true' ? '#2e7d32' : '#c62828'}
        class="port-label-text"
      >{connection.fromPort}</text>
    {/if}
    <circle
      cx={midPoint.x}
      cy={midPoint.y}
      r="8"
      fill="white"
      stroke="#666"
      stroke-width="1"
      class="delete-handle"
      on:click={handleDelete}
    />
    <text
      x={midPoint.x}
      y={midPoint.y + 1}
      text-anchor="middle"
      dominant-baseline="middle"
      font-size="10"
      fill="#666"
      class="delete-icon"
      on:click={handleDelete}
    >×</text>
  {/if}
</g>

<!-- Arrow marker definition -->
<defs>
  <marker
    id="arrowhead"
    markerWidth="10"
    markerHeight="7"
    refX="9"
    refY="3.5"
    orient="auto"
  >
    <polygon
      points="0 0, 10 3.5, 0 7"
      fill="#666"
    />
  </marker>
</defs>

<style>
  .workflow-connection {
    cursor: pointer;
  }

  .connection-path {
    transition: stroke 0.2s;
  }

  .workflow-connection:hover .connection-path {
    stroke: #007bff;
    stroke-width: 3;
  }

  .delete-handle {
    opacity: 0;
    transition: opacity 0.2s;
    pointer-events: none;
  }

  .delete-icon {
    opacity: 0;
    pointer-events: none;
    font-weight: bold;
    user-select: none;
  }

  .workflow-connection:hover .delete-handle,
  .workflow-connection:hover .delete-icon {
    opacity: 1;
    pointer-events: all;
  }

  .delete-handle:hover {
    fill: #ff4757;
    stroke: #ff4757;
  }

  .delete-icon:hover {
    fill: white;
  }
</style>

<script>
  function getMidPoint(fromNode, toNode) {
    if (!fromNode || !toNode) return { x: 0, y: 0 };

    const fromX = (fromNode.screenX || fromNode.position.x) + fromNode.width;
    const fromY = (fromNode.screenY || fromNode.position.y) + (fromNode.height / 2);
    const toX = (toNode.screenX || toNode.position.x);
    const toY = (toNode.screenY || toNode.position.y) + (toNode.height / 2);

    return {
      x: (fromX + toX) / 2,
      y: (fromY + toY) / 2
    };
  }
</script>
