<script>
// @ts-nocheck
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import { apiPost } from '$lib/api.js';

  export let accept = '*/*';
  export let multiple = false;
  export let maxSize = 10 * 1024 * 1024; // 10MB default
  export let showAIOption = true;
  export let analysisType = 'ocr'; // 'ocr', 'invoice', 'receipt', 'business_card', 'contract'

  let fileInput;
  let selectedFiles = [];
  let isUploading = false;
  let uploadProgress = {};
  let enableAIProcessing = true;
  let wsConnection = null;
  let reconnectAttempts = 0;
  let maxReconnectAttempts = 5;

  const dispatch = createEventDispatcher();

  // WebSocket connection management
  function connectWebSocket() {
    if (wsConnection && wsConnection.readyState === WebSocket.OPEN) {
      return;
    }

    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/api/v1/ws`;

      wsConnection = new WebSocket(wsUrl);

      wsConnection.onopen = () => {
        console.log('WebSocket connected for document processing updates');
        reconnectAttempts = 0;
      };

      wsConnection.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          handleWebSocketMessage(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      wsConnection.onclose = () => {
        console.log('WebSocket connection closed');
        if (reconnectAttempts < maxReconnectAttempts) {
          reconnectAttempts++;
          setTimeout(connectWebSocket, 2000 * reconnectAttempts);
        }
      };

      wsConnection.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    } catch (error) {
      console.error('Failed to connect WebSocket:', error);
    }
  }

  function handleWebSocketMessage(message) {
    switch (message.type) {
      case 'document_processing':
        handleDocumentProcessing(message);
        break;
      case 'document_processed':
        handleDocumentProcessed(message);
        break;
      case 'document_failed':
        handleDocumentFailed(message);
        break;
    }
  }

  function handleDocumentProcessing(message) {
    const { document_id, status, message: statusMessage } = message.data;
    if (uploadProgress[document_id]) {
      uploadProgress[document_id] = {
        ...uploadProgress[document_id],
        status: 'processing',
        progress: 75,
        processingMessage: statusMessage || 'Processing with AI...'
      };
      uploadProgress = { ...uploadProgress }; // Trigger reactivity
    }
  }

  function handleDocumentProcessed(message) {
    const {
      document_id,
      analysis_id,
      status,
      text_content,
      confidence,
      pages,
      fields_count,
      tables_count,
      message: successMessage
    } = message.data;

    if (uploadProgress[document_id]) {
      uploadProgress[document_id] = {
        ...uploadProgress[document_id],
        status: 'completed',
        progress: 100,
        analysisId: analysis_id,
        textContent: text_content,
        confidence: confidence,
        pages: pages,
        fieldsCount: fields_count,
        tablesCount: tables_count,
        processingMessage: successMessage || 'AI processing completed successfully'
      };
      uploadProgress = { ...uploadProgress }; // Trigger reactivity

      // Dispatch success event with AI results
      dispatch('aiProcessingComplete', {
        documentId: document_id,
        analysisId: analysis_id,
        textContent: text_content,
        confidence: confidence,
        pages: pages,
        fieldsCount: fields_count,
        tablesCount: tables_count
      });
    }
  }

  function handleDocumentFailed(message) {
    const { document_id, error } = message.data;
    if (uploadProgress[document_id]) {
      uploadProgress[document_id] = {
        ...uploadProgress[document_id],
        status: 'failed',
        progress: 0,
        error: error || 'AI processing failed'
      };
      uploadProgress = { ...uploadProgress }; // Trigger reactivity
    }
  }

  // Connect to WebSocket when component mounts
  onMount(() => {
    if (showAIOption) {
      connectWebSocket();
    }
  });

  onDestroy(() => {
    if (wsConnection) {
      wsConnection.close();
    }
  });

  function handleFileSelect(event) {
    const files = Array.from(event.target.files);
    selectedFiles = files.filter(file => {
      if (file.size > maxSize) {
        alert(`File ${file.name} is too large. Maximum size is ${maxSize / (1024 * 1024)}MB.`);
        return false;
      }
      return true;
    });
  }

  async function uploadFiles() {
    if (selectedFiles.length === 0) return;

    isUploading = true;
    uploadProgress = {};

    try {
      for (const file of selectedFiles) {
        uploadProgress[file.name] = { status: 'uploading', progress: 0 };

        // Read file as base64
        const base64Data = await readFileAsBase64(file);
        uploadProgress[file.name].progress = 25;

        // Prepare upload data
        const uploadData = {
          title: file.name,
          file_name: file.name,
          file_data: base64Data,
          file_type: getFileType(file),
          process_with_ai: showAIOption && enableAIProcessing,
          analysis_type: analysisType
        };

        uploadProgress[file.name].progress = 50;

        // Upload file
        const response = await apiPost('/documents/upload', uploadData);
        uploadProgress[file.name] = {
          status: 'completed',
          progress: 100,
          result: response
        };

        // Dispatch success event
        dispatch('uploadSuccess', {
          file: file,
          response: response,
          aiProcessed: enableAIProcessing
        });

      }
    } catch (error) {
      console.error('Upload failed:', error);
      selectedFiles.forEach(file => {
        uploadProgress[file.name] = {
          status: 'error',
          progress: 0,
          error: error.message
        };
      });

      dispatch('uploadError', { error });
    } finally {
      isUploading = false;
    }
  }

  function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result.split(',')[1]); // Remove data URL prefix
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  }

  function getFileType(file) {
    const extension = file.name.split('.').pop().toLowerCase();
    switch (extension) {
      case 'pdf':
        return 'pdf';
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
      case 'bmp':
      case 'tiff':
        return 'image';
      case 'doc':
      case 'docx':
        return 'doc';
      default:
        return 'unknown';
    }
  }

  function removeFile(index) {
    selectedFiles = selectedFiles.filter((_, i) => i !== index);
    if (selectedFiles.length === 0) {
      uploadProgress = {};
    }
  }

  function clearAll() {
    selectedFiles = [];
    uploadProgress = {};
    if (fileInput) fileInput.value = '';
  }

  $: canUpload = selectedFiles.length > 0 && !isUploading;
</script>

<div class="document-upload">
  <div class="upload-area">
    <input
      bind:this={fileInput}
      type="file"
      {accept}
      {multiple}
      on:change={handleFileSelect}
      class="file-input"
      id="file-input"
    />

    <label for="file-input" class="upload-zone">
      <div class="upload-icon">📎</div>
      <div class="upload-text">
        <strong>Click to select files</strong> or drag and drop
      </div>
      <div class="upload-hint">
        Maximum file size: {maxSize / (1024 * 1024)}MB
        {#if accept !== '*/*'}
          <br>Accepted types: {accept}
        {/if}
      </div>
    </label>
  </div>

  {#if showAIOption}
    <div class="ai-options">
      <label class="ai-toggle">
        <input
          type="checkbox"
          bind:checked={enableAIProcessing}
          disabled={isUploading}
        />
        <span class="toggle-slider"></span>
        Enable AI Document Processing
      </label>

      {#if enableAIProcessing}
        <div class="analysis-type">
          <label for="analysis-select">Analysis Type:</label>
          <select
            id="analysis-select"
            bind:value={analysisType}
            disabled={isUploading}
          >
            <option value="ocr">OCR (Text Extraction)</option>
            <option value="invoice">Invoice Analysis</option>
            <option value="receipt">Receipt Analysis</option>
            <option value="business_card">Business Card</option>
            <option value="contract">Contract Analysis</option>
          </select>
        </div>
      {/if}
    </div>
  {/if}

  {#if selectedFiles.length > 0}
    <div class="file-list">
      <h4>Selected Files ({selectedFiles.length})</h4>
      {#each selectedFiles as file, index}
        <div class="file-item">
          <div class="file-info">
            <span class="file-name">{file.name}</span>
            <span class="file-size">({(file.size / 1024).toFixed(1)} KB)</span>
          </div>

          {#if uploadProgress[file.name]}
            <div class="progress-bar">
              <div
                class="progress-fill"
                style="width: {uploadProgress[file.name].progress}%"
              ></div>
            </div>
            <div class="progress-status">
              {#if uploadProgress[file.name].status === 'uploading'}
                Uploading...
              {:else if uploadProgress[file.name].status === 'processing'}
                🤖 {uploadProgress[file.name].processingMessage || 'Processing with AI...'}
              {:else if uploadProgress[file.name].status === 'completed'}
                ✅ Completed
              {:else if uploadProgress[file.name].status === 'error'}
                ❌ Error: {uploadProgress[file.name].error}
              {/if}
            </div>
          {:else}
            <button
              class="remove-btn"
              on:click={() => removeFile(index)}
              disabled={isUploading}
            >
              ✕
            </button>
          {/if}
        </div>
      {/each}
    </div>

    <div class="upload-actions">
      <button
        class="upload-btn"
        on:click={uploadFiles}
        disabled={!canUpload}
      >
        {#if isUploading}
          Uploading...
        {:else}
          Upload {selectedFiles.length} File{selectedFiles.length > 1 ? 's' : ''}
        {/if}
      </button>

      <button
        class="clear-btn"
        on:click={clearAll}
        disabled={isUploading}
      >
        Clear All
      </button>
    </div>
  {/if}
</div>

<style>
  .document-upload {
    max-width: 600px;
    margin: 0 auto;
  }

  .upload-area {
    position: relative;
    margin-bottom: 1rem;
  }

  .file-input {
    position: absolute;
    opacity: 0;
    width: 100%;
    height: 100%;
    cursor: pointer;
  }

  .upload-zone {
    display: block;
    border: 2px dashed #ccc;
    border-radius: 8px;
    padding: 2rem;
    text-align: center;
    cursor: pointer;
    transition: border-color 0.3s;
    background: #fafafa;
  }

  .upload-zone:hover {
    border-color: #007bff;
    background: #f0f8ff;
  }

  .upload-icon {
    font-size: 3rem;
    margin-bottom: 0.5rem;
  }

  .upload-text {
    margin-bottom: 0.5rem;
    color: #666;
  }

  .upload-hint {
    font-size: 0.9rem;
    color: #888;
  }

  .ai-options {
    margin-bottom: 1rem;
    padding: 1rem;
    background: #f8f9fa;
    border-radius: 6px;
    border: 1px solid #e9ecef;
  }

  .ai-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 500;
    cursor: pointer;
  }

  .analysis-type {
    margin-top: 0.75rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .analysis-type select {
    padding: 0.25rem 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    background: white;
  }

  .file-list {
    margin-bottom: 1rem;
  }

  .file-list h4 {
    margin: 0 0 0.5rem 0;
    color: #333;
  }

  .file-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.75rem;
    border: 1px solid #e9ecef;
    border-radius: 6px;
    margin-bottom: 0.5rem;
    background: white;
  }

  .file-info {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .file-name {
    font-weight: 500;
    color: #333;
  }

  .file-size {
    font-size: 0.9rem;
    color: #666;
  }

  .progress-bar {
    flex: 1;
    height: 8px;
    background: #e9ecef;
    border-radius: 4px;
    margin: 0 1rem;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: #007bff;
    transition: width 0.3s ease;
    border-radius: 4px;
  }

  .progress-status {
    font-size: 0.9rem;
    color: #666;
    min-width: 100px;
    text-align: right;
  }

  .remove-btn {
    background: #dc3545;
    color: white;
    border: none;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    transition: background-color 0.2s;
  }

  .remove-btn:hover:not(:disabled) {
    background: #c82333;
  }

  .remove-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .upload-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: center;
  }

  .upload-btn, .clear-btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .upload-btn {
    background: #007bff;
    color: white;
  }

  .upload-btn:hover:not(:disabled) {
    background: #0056b3;
  }

  .upload-btn:disabled {
    background: #6c757d;
    cursor: not-allowed;
  }

  .clear-btn {
    background: #6c757d;
    color: white;
  }

  .clear-btn:hover:not(:disabled) {
    background: #545b62;
  }

  .clear-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Toggle switch styles */
  .toggle-slider {
    position: relative;
    display: inline-block;
    width: 44px;
    height: 24px;
    background-color: #ccc;
    border-radius: 24px;
    transition: 0.4s;
    cursor: pointer;
  }

  .ai-toggle input:checked + .toggle-slider {
    background-color: #007bff;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    border-radius: 50%;
    transition: 0.4s;
  }

  .ai-toggle input:checked + .toggle-slider:before {
    transform: translateX(20px);
  }

  .ai-toggle input {
    opacity: 0;
    width: 0;
    height: 0;
  }
</style>
