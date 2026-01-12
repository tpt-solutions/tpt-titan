<script>
  import { createEventDispatcher } from 'svelte';

  export let node;
  export let connectors = [];

  const dispatch = createEventDispatcher();

  // Available operators for conditions
  const operators = [
    { value: 'equals', label: 'Equals (=)' },
    { value: 'not_equals', label: 'Not Equals (!=)' },
    { value: 'greater_than', label: 'Greater Than (>)' },
    { value: 'less_than', label: 'Less Than (<)' },
    { value: 'contains', label: 'Contains' },
    { value: 'not_contains', label: 'Does Not Contain' }
  ];

  // Available notification types
  const notificationTypes = [
    { value: 'info', label: 'Info' },
    { value: 'success', label: 'Success' },
    { value: 'warning', label: 'Warning' },
    { value: 'error', label: 'Error' }
  ];

  function updateConfig(updates) {
    const newConfig = { ...node.config, ...updates };
    dispatch('update', { config: newConfig });
  }

  function updateNodeName(name) {
    dispatch('update', { name });
  }

  function handleConnectorChange(event) {
    const connectorId = event.target.value;
    updateConfig({ connector: connectorId });
  }

  function handleConditionChange(field, value) {
    updateConfig({ [field]: value });
  }

  function handleDelayChange(seconds) {
    updateConfig({ delay_seconds: parseInt(seconds) || 60 });
  }

  function handleNotificationChange(field, value) {
    updateConfig({ [field]: value });
  }
</script>

<div class="node-properties">
  <!-- Node Name -->
  <div class="property-group">
    <label class="property-label">Node Name</label>
    <input
      type="text"
      class="property-input"
      bind:value={node.name}
      on:input={(e) => updateNodeName(e.target.value)}
      placeholder="Enter node name"
    />
  </div>

  <!-- Node-specific properties -->
  {#if node.type === 'trigger'}
    <div class="property-group">
      <h5>Trigger Configuration</h5>
      <p class="help-text">This node starts the workflow execution.</p>

      <label class="property-label">Trigger Type</label>
      <select
        class="property-select"
        bind:value={node.config.triggerType}
        on:change={(e) => updateConfig({ triggerType: e.target.value })}
      >
        <option value="manual">Manual</option>
        <option value="scheduled">Scheduled</option>
        <option value="webhook">Webhook</option>
      </select>
    </div>

  {:else if node.type === 'action'}
    <div class="property-group">
      <h5>Action Configuration</h5>

      <label class="property-label">Connector</label>
      <select
        class="property-select"
        bind:value={node.config.connector}
        on:change={handleConnectorChange}
      >
        <option value="">Select a connector</option>
        {#each connectors as connector}
          <option value={connector.id}>{connector.name}</option>
        {/each}
      </select>

      {#if node.config.connector}
        {@const selectedConnector = connectors.find(c => c.id === node.config.connector)}
        {#if selectedConnector}
          <div class="connector-config">
            <label class="property-label">Action</label>
            <select
              class="property-select"
              bind:value={node.config.action}
              on:change={(e) => updateConfig({ action: e.target.value })}
            >
              <option value="">Select an action</option>
              <option value="send_email">Send Email</option>
              <option value="create_task">Create Task</option>
              <option value="update_calendar">Update Calendar</option>
              <option value="send_notification">Send Notification</option>
            </select>

            <!-- Action-specific parameters -->
            {#if node.config.action === 'send_email'}
              <div class="action-params">
                <label class="property-label">To</label>
                <input
                  type="email"
                  class="property-input"
                  bind:value={node.config.parameters.to}
                  on:input={(e) => updateConfig({
                    parameters: { ...node.config.parameters, to: e.target.value }
                  })}
                  placeholder="recipient@example.com"
                />

                <label class="property-label">Subject</label>
                <input
                  type="text"
                  class="property-input"
                  bind:value={node.config.parameters.subject}
                  on:input={(e) => updateConfig({
                    parameters: { ...node.config.parameters, subject: e.target.value }
                  })}
                  placeholder="Email subject"
                />

                <label class="property-label">Message</label>
                <textarea
                  class="property-textarea"
                  bind:value={node.config.parameters.body}
                  on:input={(e) => updateConfig({
                    parameters: { ...node.config.parameters, body: e.target.value }
                  })}
                  placeholder="Email body"
                  rows="3"
                ></textarea>
              </div>
            {:else if node.config.action === 'create_task'}
              <div class="action-params">
                <label class="property-label">Task Title</label>
                <input
                  type="text"
                  class="property-input"
                  bind:value={node.config.parameters.title}
                  on:input={(e) => updateConfig({
                    parameters: { ...node.config.parameters, title: e.target.value }
                  })}
                  placeholder="Task title"
                />

                <label class="property-label">Description</label>
                <textarea
                  class="property-textarea"
                  bind:value={node.config.parameters.description}
                  on:input={(e) => updateConfig({
                    parameters: { ...node.config.parameters, description: e.target.value }
                  })}
                  placeholder="Task description"
                  rows="2"
                ></textarea>
              </div>
            {/if}
          </div>
        {/if}
      {/if}
    </div>

  {:else if node.type === 'condition'}
    <div class="property-group">
      <h5>Condition Configuration</h5>

      <label class="property-label">Field</label>
      <input
        type="text"
        class="property-input"
        bind:value={node.config.field}
        on:input={(e) => handleConditionChange('field', e.target.value)}
        placeholder="field_name"
      />

      <label class="property-label">Operator</label>
      <select
        class="property-select"
        bind:value={node.config.operator}
        on:change={(e) => handleConditionChange('operator', e.target.value)}
      >
        {#each operators as op}
          <option value={op.value}>{op.label}</option>
        {/each}
      </select>

      <label class="property-label">Value</label>
      <input
        type="text"
        class="property-input"
        bind:value={node.config.value}
        on:input={(e) => handleConditionChange('value', e.target.value)}
        placeholder="comparison value"
      />

      <div class="condition-help">
        <p class="help-text">
          This condition will route execution to the "true" output if the condition is met,
          otherwise to the "false" output.
        </p>
      </div>
    </div>

  {:else if node.type === 'delay'}
    <div class="property-group">
      <h5>Delay Configuration</h5>

      <label class="property-label">Delay (seconds)</label>
      <input
        type="number"
        class="property-input"
        min="1"
        max="86400"
        bind:value={node.config.delay_seconds}
        on:input={(e) => handleDelayChange(e.target.value)}
      />

      <div class="delay-preview">
        {#if node.config.delay_seconds}
          {@const seconds = node.config.delay_seconds}
          {@const minutes = Math.floor(seconds / 60)}
          {@const remainingSeconds = seconds % 60}
          {@const hours = Math.floor(minutes / 60)}
          {@const remainingMinutes = minutes % 60}

          {#if hours > 0}
            Delay: {hours}h {remainingMinutes}m {remainingSeconds}s
          {:else if minutes > 0}
            Delay: {minutes}m {remainingSeconds}s
          {:else}
            Delay: {seconds}s
          {/if}
        {:else}
          Delay: 60s (default)
        {/if}
      </div>
    </div>

  {:else if node.type === 'notification'}
    <div class="property-group">
      <h5>Notification Configuration</h5>

      <label class="property-label">Title</label>
      <input
        type="text"
        class="property-input"
        bind:value={node.config.title}
        on:input={(e) => handleNotificationChange('title', e.target.value)}
        placeholder="Notification title"
      />

      <label class="property-label">Message</label>
      <textarea
        class="property-textarea"
        bind:value={node.config.message}
        on:input={(e) => handleNotificationChange('message', e.target.value)}
        placeholder="Notification message"
        rows="3"
      ></textarea>

      <label class="property-label">Type</label>
      <select
        class="property-select"
        bind:value={node.config.type}
        on:change={(e) => handleNotificationChange('type', e.target.value)}
      >
        {#each notificationTypes as type}
          <option value={type.value}>{type.label}</option>
        {/each}
      </select>
    </div>
  {/if}

  <!-- Node Position (for advanced users) -->
  <div class="property-group">
    <h5>Position</h5>
    <div class="position-inputs">
      <div>
        <label class="property-label">X</label>
        <input
          type="number"
          class="property-input small"
          bind:value={node.position.x}
          on:input={(e) => dispatch('update', {
            position: { ...node.position, x: parseInt(e.target.value) || 0 }
          })}
        />
      </div>
      <div>
        <label class="property-label">Y</label>
        <input
          type="number"
          class="property-input small"
          bind:value={node.position.y}
          on:input={(e) => dispatch('update', {
            position: { ...node.position, y: parseInt(e.target.value) || 0 }
          })}
        />
      </div>
    </div>
  </div>
</div>

<style>
  .node-properties {
    padding: 1rem;
  }

  .property-group {
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #e0e0e0;
  }

  .property-group:last-child {
    border-bottom: none;
    margin-bottom: 0;
  }

  .property-group h5 {
    margin: 0 0 0.5rem 0;
    color: #333;
    font-size: 1rem;
    font-weight: 600;
  }

  .property-label {
    display: block;
    margin-bottom: 0.25rem;
    font-size: 0.9rem;
    font-weight: 500;
    color: #555;
  }

  .property-input,
  .property-select,
  .property-textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 0.9rem;
    background: white;
    transition: border-color 0.2s;
  }

  .property-input:focus,
  .property-select:focus,
  .property-textarea:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
  }

  .property-input.small {
    width: 80px;
  }

  .property-textarea {
    resize: vertical;
    min-height: 60px;
  }

  .position-inputs {
    display: flex;
    gap: 1rem;
  }

  .position-inputs > div {
    flex: 1;
  }

  .help-text {
    font-size: 0.8rem;
    color: #666;
    margin: 0.5rem 0;
    line-height: 1.4;
  }

  .condition-help {
    margin-top: 1rem;
    padding: 0.75rem;
    background: rgba(255, 152, 0, 0.1);
    border-radius: 4px;
    border-left: 3px solid #ff9800;
  }

  .delay-preview {
    margin-top: 0.5rem;
    padding: 0.5rem;
    background: rgba(156, 39, 176, 0.1);
    border-radius: 4px;
    font-size: 0.9rem;
    color: #7b1fa2;
    font-weight: 500;
  }

  .connector-config {
    margin-top: 1rem;
    padding: 1rem;
    background: #f8f9fa;
    border-radius: 4px;
    border: 1px solid #e9ecef;
  }

  .action-params {
    margin-top: 1rem;
  }

  .action-params > * {
    margin-bottom: 1rem;
  }

  .action-params > *:last-child {
    margin-bottom: 0;
  }

  /* Responsive adjustments */
  @media (max-width: 768px) {
    .node-properties {
      padding: 0.75rem;
    }

    .position-inputs {
      flex-direction: column;
      gap: 0.5rem;
    }

    .property-input.small {
      width: 100%;
    }
  }
</style>
