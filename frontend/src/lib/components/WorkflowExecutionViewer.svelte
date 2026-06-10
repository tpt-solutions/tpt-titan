<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { apiGet } from '$lib/api.js';

  export let executionId = null;
  export let workflowId = null;

  let execution = null;
  let executions = [];
  let isLoading = false;
  let selectedExecution = null;
  let autoRefresh = false;
  let refreshInterval = null;

  const dispatch = createEventDispatcher();

  onMount(async () => {
    if (executionId) {
      await loadExecution(executionId);
    } else if (workflowId) {
      await loadWorkflowExecutions();
    }

    // Auto-refresh for running executions
    if (autoRefresh) {
      startAutoRefresh();
    }
  });

  async function loadExecution(id) {
    isLoading = true;
    try {
      const response = await apiGet(`/workflows/executions/${id}`);
      execution = response.execution;
      selectedExecution = execution;
    } catch (error) {
      console.error('Failed to load execution:', error);
      dispatch('error', { message: 'Failed to load execution details' });
    } finally {
      isLoading = false;
    }
  }

  async function loadWorkflowExecutions() {
    isLoading = true;
    try {
      const response = await apiGet(`/workflows/${workflowId}/executions`);
      executions = response.executions || [];
    } catch (error) {
      console.error('Failed to load executions:', error);
      dispatch('error', { message: 'Failed to load workflow executions' });
    } finally {
      isLoading = false;
    }
  }

  function selectExecution(exec) {
    selectedExecution = exec;
    execution = exec;
  }

  function startAutoRefresh() {
    if (refreshInterval) clearInterval(refreshInterval);

    refreshInterval = setInterval(async () => {
      if (executionId && execution && ['running', 'pending'].includes(execution.status)) {
        await loadExecution(executionId);
      } else if (workflowId) {
        await loadWorkflowExecutions();
      }
    }, 2000); // Refresh every 2 seconds
  }

  function stopAutoRefresh() {
    if (refreshInterval) {
      clearInterval(refreshInterval);
      refreshInterval = null;
    }
  }

  function toggleAutoRefresh() {
    autoRefresh = !autoRefresh;
    if (autoRefresh) {
      startAutoRefresh();
    } else {
      stopAutoRefresh();
    }
  }

  function formatDuration(ms) {
    if (!ms) return '0s';

    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);

    if (hours > 0) {
      return `${hours}h ${minutes % 60}m ${seconds % 60}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`;
    } else {
      return `${seconds}s`;
    }
  }

  function getStatusColor(status) {
    switch (status) {
      case 'completed': return '#28a745';
      case 'running': return '#007bff';
      case 'pending': return '#ffc107';
      case 'failed': return '#dc3545';
      default: return '#6c757d';
    }
  }

  function getStatusIcon(status) {
    switch (status) {
      case 'completed': return '✅';
      case 'running': return '🔄';
      case 'pending': return '⏳';
      case 'failed': return '❌';
      default: return '❓';
    }
  }

  // Cleanup on destroy
  import { onDestroy } from 'svelte';
  onDestroy(() => {
    stopAutoRefresh();
  });
</script>

<div class="execution-viewer">
  <!-- Header -->
  <div class="viewer-header">
    <h3>Workflow Execution</h3>
    <div class="header-controls">
      <label class="auto-refresh-toggle">
        <input
          type="checkbox"
          bind:checked={autoRefresh}
          on:change={toggleAutoRefresh}
        />
        <span class="toggle-slider"></span>
        Auto-refresh
      </label>
      <button
        class="btn btn-secondary"
        on:click={workflowId ? loadWorkflowExecutions : () => loadExecution(executionId)}
        disabled={isLoading}
      >
        🔄 Refresh
      </button>
    </div>
  </div>

  {#if isLoading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Loading execution details...</p>
    </div>
  {:else if workflowId && executions.length > 0}
    <!-- Executions List -->
    <div class="executions-list">
      <h4>Recent Executions</h4>
      {#each executions as exec}
        <div
          class="execution-item"
          class:selected={selectedExecution && selectedExecution.id === exec.id}
          on:click={() => selectExecution(exec)}
        >
          <div class="execution-header">
            <span class="status-icon" style="color: {getStatusColor(exec.status)}">
              {getStatusIcon(exec.status)}
            </span>
            <span class="execution-id">#{exec.id.split('-')[0]}</span>
            <span class="execution-status" style="color: {getStatusColor(exec.status)}">
              {exec.status}
            </span>
            <span class="execution-time">
              {new Date(exec.started_at).toLocaleString()}
            </span>
          </div>

          {#if exec.duration}
            <div class="execution-duration">
              Duration: {formatDuration(exec.duration)}
            </div>
          {/if}

          {#if exec.error_message}
            <div class="execution-error">
              {exec.error_message}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}

  {#if selectedExecution || execution}
    {@const currentExecution = selectedExecution || execution}
    <!-- Execution Details -->
    <div class="execution-details">
      <div class="execution-summary">
        <div class="summary-header">
          <h4>Execution Details</h4>
          <div class="execution-meta">
            <span class="status-badge" style="background-color: {getStatusColor(currentExecution.status)}">
              {getStatusIcon(currentExecution.status)} {currentExecution.status}
            </span>
            <span class="execution-timestamp">
              Started: {new Date(currentExecution.started_at).toLocaleString()}
            </span>
            {#if currentExecution.completed_at}
              <span class="execution-timestamp">
                Completed: {new Date(currentExecution.completed_at).toLocaleString()}
              </span>
            {/if}
          </div>
        </div>

        <div class="summary-stats">
          <div class="stat-item">
            <span class="stat-label">Duration</span>
            <span class="stat-value">
              {#if currentExecution.duration}
                {formatDuration(currentExecution.duration)}
              {:else}
                Still running...
              {/if}
            </span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Workflow</span>
            <span class="stat-value">#{currentExecution.workflow_id.split('-')[0]}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Trigger</span>
            <span class="stat-value">{currentExecution.trigger_type}</span>
          </div>
        </div>

        {#if currentExecution.trigger_data}
          <div class="trigger-data">
            <h5>Trigger Data</h5>
            <pre class="json-data">{JSON.stringify(JSON.parse(currentExecution.trigger_data), null, 2)}</pre>
          </div>
        {/if}

        {#if currentExecution.error_message}
          <div class="error-message">
            <h5>Error</h5>
            <div class="error-content">{currentExecution.error_message}</div>
          </div>
        {/if}
      </div>

      <!-- Node Execution States -->
      {#if currentExecution.node_states}
        {@const nodeStates = typeof currentExecution.node_states === 'string'
          ? JSON.parse(currentExecution.node_states)
          : currentExecution.node_states}
        <div class="node-states">
          <h4>Node Execution Log</h4>
          <div class="states-list">
            {#each Object.entries(nodeStates) as [nodeId, state]}
              <div class="node-state-item" class:failed={state.status === 'failed'}>
                <div class="node-header">
                  <span class="node-id">{nodeId}</span>
                  <span class="node-status" style="color: {getStatusColor(state.status)}">
                    {state.status}
                  </span>
                  {#if state.started_at}
                    <span class="node-time">
                      {new Date(state.started_at).toLocaleTimeString()}
                    </span>
                  {/if}
                </div>

                {#if state.error}
                  <div class="node-error">{state.error}</div>
                {/if}

                {#if state.completed_at}
                  <div class="node-completion">
                    Completed at {new Date(state.completed_at).toLocaleTimeString()}
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Output Data -->
      {#if currentExecution.output_data}
        {@const outputData = typeof currentExecution.output_data === 'string'
          ? JSON.parse(currentExecution.output_data)
          : currentExecution.output_data}
        <div class="output-data">
          <h4>Output Data</h4>
          <pre class="json-data">{JSON.stringify(outputData, null, 2)}</pre>
        </div>
      {/if}
    </div>
  {:else if !isLoading}
    <div class="no-execution">
      <div class="empty-state">
        <div class="empty-icon">📋</div>
        <h4>No Execution Selected</h4>
        <p>Select an execution from the list to view details, or run the workflow to create a new execution.</p>
      </div>
    </div>
  {/if}
</div>

<style>
  .execution-viewer {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .viewer-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid #e0e0e0;
    background: #f8f9fa;
  }

  .viewer-header h3 {
    margin: 0;
    color: #333;
  }

  .header-controls {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .auto-refresh-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.9rem;
    cursor: pointer;
  }

  .btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    transition: all 0.2s;
  }

  .btn-secondary {
    background: #6c757d;
    color: white;
  }

  .btn-secondary:hover:not(:disabled) {
    background: #545b62;
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #007bff;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .executions-list {
    max-height: 300px;
    overflow-y: auto;
    border-bottom: 1px solid #e0e0e0;
  }

  .executions-list h4 {
    margin: 0;
    padding: 1rem 1.5rem;
    background: #f8f9fa;
    border-bottom: 1px solid #e0e0e0;
    color: #333;
  }

  .execution-item {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid #f0f0f0;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .execution-item:hover {
    background: #f8f9fa;
  }

  .execution-item.selected {
    background: #e3f2fd;
    border-left: 4px solid #007bff;
  }

  .execution-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.25rem;
  }

  .status-icon {
    font-size: 1.1rem;
  }

  .execution-id {
    font-weight: 500;
    color: #333;
  }

  .execution-status {
    font-size: 0.8rem;
    font-weight: 500;
    text-transform: uppercase;
    padding: 0.25rem 0.5rem;
    border-radius: 12px;
    background: rgba(0, 0, 0, 0.1);
  }

  .execution-time {
    margin-left: auto;
    font-size: 0.8rem;
    color: #666;
  }

  .execution-duration {
    font-size: 0.8rem;
    color: #666;
  }

  .execution-error {
    font-size: 0.8rem;
    color: #dc3545;
    margin-top: 0.25rem;
  }

  .execution-details {
    flex: 1;
    padding: 1.5rem;
    overflow-y: auto;
  }

  .execution-summary {
    margin-bottom: 2rem;
  }

  .summary-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 1rem;
  }

  .summary-header h4 {
    margin: 0;
    color: #333;
  }

  .execution-meta {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 0.25rem;
  }

  .status-badge {
    padding: 0.25rem 0.75rem;
    border-radius: 16px;
    color: white;
    font-size: 0.8rem;
    font-weight: 500;
    text-transform: uppercase;
  }

  .execution-timestamp {
    font-size: 0.8rem;
    color: #666;
  }

  .summary-stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .stat-item {
    text-align: center;
    padding: 0.75rem;
    background: #f8f9fa;
    border-radius: 6px;
  }

  .stat-label {
    display: block;
    font-size: 0.8rem;
    color: #666;
    margin-bottom: 0.25rem;
  }

  .stat-value {
    display: block;
    font-size: 1.1rem;
    font-weight: 600;
    color: #333;
  }

  .trigger-data, .error-message, .output-data {
    margin-top: 1.5rem;
  }

  .trigger-data h5, .error-message h5, .output-data h5, .node-states h4 {
    margin: 0 0 0.75rem 0;
    color: #333;
    font-size: 1rem;
  }

  .json-data {
    background: #f8f9fa;
    border: 1px solid #e9ecef;
    border-radius: 4px;
    padding: 1rem;
    font-size: 0.8rem;
    font-family: 'Monaco', 'Menlo', monospace;
    overflow-x: auto;
    white-space: pre-wrap;
  }

  .error-content {
    background: #f8d7da;
    border: 1px solid #f5c6cb;
    border-radius: 4px;
    padding: 1rem;
    color: #721c24;
  }

  .node-states {
    margin-top: 2rem;
  }

  .states-list {
    max-height: 400px;
    overflow-y: auto;
  }

  .node-state-item {
    padding: 0.75rem;
    margin-bottom: 0.5rem;
    border: 1px solid #e9ecef;
    border-radius: 4px;
    background: #f8f9fa;
  }

  .node-state-item.failed {
    border-color: #f5c6cb;
    background: #f8d7da;
  }

  .node-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.5rem;
  }

  .node-id {
    font-weight: 500;
    color: #333;
  }

  .node-status {
    font-size: 0.8rem;
    font-weight: 500;
    text-transform: uppercase;
    padding: 0.25rem 0.5rem;
    border-radius: 12px;
    background: rgba(0, 0, 0, 0.1);
  }

  .node-time {
    margin-left: auto;
    font-size: 0.8rem;
    color: #666;
  }

  .node-error {
    font-size: 0.8rem;
    color: #721c24;
    margin-top: 0.25rem;
  }

  .node-completion {
    font-size: 0.8rem;
    color: #155724;
    margin-top: 0.25rem;
  }

  .no-execution {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .empty-state {
    text-align: center;
    color: #666;
  }

  .empty-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
    opacity: 0.5;
  }

  .empty-state h4 {
    margin: 0 0 0.5rem 0;
    color: #333;
  }

  .empty-state p {
    margin: 0;
    max-width: 300px;
  }

  /* Toggle switch styles */
  .toggle-slider {
    position: relative;
    display: inline-block;
    width: 34px;
    height: 18px;
    background-color: #ccc;
    border-radius: 18px;
    transition: 0.4s;
  }

  .auto-refresh-toggle input:checked + .toggle-slider {
    background-color: #007bff;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 14px;
    width: 14px;
    left: 2px;
    bottom: 2px;
    background-color: white;
    border-radius: 50%;
    transition: 0.4s;
  }

  .auto-refresh-toggle input:checked + .toggle-slider:before {
    transform: translateX(16px);
  }

  .auto-refresh-toggle input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  /* Responsive adjustments */
  @media (max-width: 768px) {
    .viewer-header {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }

    .header-controls {
      justify-content: space-between;
    }

    .summary-header {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }

    .execution-meta {
      align-items: flex-start;
    }

    .summary-stats {
      grid-template-columns: 1fr;
    }

    .execution-header {
      flex-wrap: wrap;
    }

    .execution-time {
      margin-left: 0;
      margin-top: 0.25rem;
    }
  }
</style>
