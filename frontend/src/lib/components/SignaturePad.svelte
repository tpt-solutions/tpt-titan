<!-- frontend/src/lib/components/SignaturePad.svelte -->
<script>
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';

	export let value = ''; // Base64 encoded signature image
	export let width = 400;
	export let height = 200;
	export let backgroundColor = '#ffffff';
	export let penColor = '#000000';
	export let penSize = 2;
	export let disabled = false;
	export let required = false;
	export let placeholder = 'Sign here';

	const dispatch = createEventDispatcher();

	/** @type {HTMLCanvasElement} */
	let canvas;
	/** @type {CanvasRenderingContext2D} */
	let ctx;
	let isDrawing = false;
	let hasSignature = false;
	/** @type {{x: number, y: number}[]} */
	let points = [];

	onMount(() => {
		if (canvas) {
			ctx = /** @type {CanvasRenderingContext2D} */ (canvas.getContext('2d'));
			setupCanvas();
			
			// Load existing signature if provided
			if (value) {
				loadSignature(value);
			}
		}
	});

	function setupCanvas() {
		// Handle high DPI displays
		const dpr = window.devicePixelRatio || 1;
		const rect = canvas.getBoundingClientRect();
		
		canvas.width = width * dpr;
		canvas.height = height * dpr;
		
		ctx.scale(dpr, dpr);
		ctx.lineCap = 'round';
		ctx.lineJoin = 'round';
		ctx.strokeStyle = penColor;
		ctx.lineWidth = penSize;
		
		// Fill background
		ctx.fillStyle = backgroundColor;
		ctx.fillRect(0, 0, width, height);
	}

	/** @param {MouseEvent | TouchEvent} e */
	function getPoint(e) {
		const rect = canvas.getBoundingClientRect();
		const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
		const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;
		
		return {
			x: clientX - rect.left,
			y: clientY - rect.top
		};
	}

	/** @param {MouseEvent | TouchEvent} e */
	function startDrawing(e) {
		if (disabled) return;
		e.preventDefault();
		
		isDrawing = true;
		const point = getPoint(e);
		points = [point];
		
		ctx.beginPath();
		ctx.moveTo(point.x, point.y);
	}

	/** @param {MouseEvent | TouchEvent} e */
	function draw(e) {
		if (!isDrawing || disabled) return;
		e.preventDefault();
		
		const point = getPoint(e);
		points.push(point);
		
		ctx.lineTo(point.x, point.y);
		ctx.stroke();
		
		hasSignature = true;
	}

	function stopDrawing() {
		if (!isDrawing) return;
		
		isDrawing = false;
		ctx.closePath();
		
		// Save signature
		saveSignature();
	}

	function saveSignature() {
		if (!canvas) return;
		
		// Create a data URL of the signature
		const dataUrl = canvas.toDataURL('image/png');
		value = dataUrl;
		dispatch('change', { value: dataUrl });
	}

	/** @param {string} dataUrl */
	function loadSignature(dataUrl) {
		if (!ctx) return;
		
		const img = new Image();
		img.onload = () => {
			ctx.drawImage(img, 0, 0, width, height);
			hasSignature = true;
		};
		img.src = dataUrl;
	}

	function clear() {
		if (!ctx || disabled) return;
		
		ctx.fillStyle = backgroundColor;
		ctx.fillRect(0, 0, width, height);
		
		hasSignature = false;
		value = '';
		points = [];
		
		dispatch('clear');
		dispatch('change', { value: '' });
	}

	function undo() {
		if (points.length === 0 || disabled) return;
		
		// Remove last stroke (simplified - removes last 10 points)
		points = points.slice(0, -10);
		
		// Redraw
		ctx.fillStyle = backgroundColor;
		ctx.fillRect(0, 0, width, height);
		
		if (points.length > 0) {
			ctx.beginPath();
			ctx.moveTo(points[0].x, points[0].y);
			
			for (let i = 1; i < points.length; i++) {
				ctx.lineTo(points[i].x, points[i].y);
			}
			
			ctx.stroke();
			hasSignature = true;
			saveSignature();
		} else {
			hasSignature = false;
			value = '';
			dispatch('change', { value: '' });
		}
	}

	// Touch event handlers for mobile
	/** @param {TouchEvent} e */
	function handleTouchStart(e) {
		startDrawing(e);
	}

	/** @param {TouchEvent} e */
	function handleTouchMove(e) {
		draw(e);
	}

	/** @param {TouchEvent} e */
	function handleTouchEnd(e) {
		stopDrawing();
	}

	// Mouse event handlers
	/** @param {MouseEvent} e */
	function handleMouseDown(e) {
		startDrawing(e);
	}

	/** @param {MouseEvent} e */
	function handleMouseMove(e) {
		draw(e);
	}

	/** @param {MouseEvent} e */
	function handleMouseUp(e) {
		stopDrawing();
	}

	/** @param {MouseEvent} e */
	function handleMouseLeave(e) {
		stopDrawing();
	}

	// Export signature as different formats
	function exportAsPNG() {
		if (!hasSignature) return null;
		return canvas.toDataURL('image/png');
	}

	function exportAsJPEG(quality = 0.9) {
		if (!hasSignature) return null;
		return canvas.toDataURL('image/jpeg', quality);
	}

	function exportAsSVG() {
		if (!hasSignature) return null;
		
		// Create SVG from points
		let pathData = '';
		if (points.length > 0) {
			pathData = `M ${points[0].x} ${points[0].y}`;
			for (let i = 1; i < points.length; i++) {
				pathData += ` L ${points[i].x} ${points[i].y}`;
			}
		}
		
		const svg = `
			<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
				<rect width="100%" height="100%" fill="${backgroundColor}"/>
				<path d="${pathData}" stroke="${penColor}" stroke-width="${penSize}" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
			</svg>
		`;
		
		return 'data:image/svg+xml;base64,' + btoa(svg);
	}

	// Validate if signature exists
	function validate() {
		if (required && !hasSignature) {
			return { valid: false, error: 'Signature is required' };
		}
		return { valid: true, error: null };
	}

	// Reactive updates
	$: if (ctx && penColor) {
		ctx.strokeStyle = penColor;
	}

	$: if (ctx && penSize) {
		ctx.lineWidth = penSize;
	}

	$: if (value && !hasSignature && canvas) {
		loadSignature(value);
	}
</script>

<div class="signature-pad-container" class:disabled>
	<div 
		class="signature-pad"
		style="width: {width}px; height: {height}px; background-color: {backgroundColor};"
	>
		<canvas
			bind:this={canvas}
			style="width: {width}px; height: {height}px; touch-action: none;"
			on:mousedown={handleMouseDown}
			on:mousemove={handleMouseMove}
			on:mouseup={handleMouseUp}
			on:mouseleave={handleMouseLeave}
			on:touchstart={handleTouchStart}
			on:touchmove={handleTouchMove}
			on:touchend={handleTouchEnd}
		></canvas>
		
		{#if !hasSignature && placeholder}
			<div class="placeholder">{placeholder}</div>
		{/if}
		
		{#if disabled}
			<div class="disabled-overlay">Disabled</div>
		{/if}
	</div>
	
	<div class="toolbar">
		<button 
			type="button" 
			class="btn btn-clear" 
			on:click={clear}
			disabled={!hasSignature || disabled}
			title="Clear signature"
		>
			🗑️ Clear
		</button>
		
		<button 
			type="button" 
			class="btn btn-undo" 
			on:click={undo}
			disabled={points.length === 0 || disabled}
			title="Undo last stroke"
		>
			↩️ Undo
		</button>
		
		{#if hasSignature}
			<span class="status">✓ Signed</span>
		{:else}
			<span class="status unsigned">✗ Not signed</span>
		{/if}
		
		{#if required}
			<span class="required-indicator">* Required</span>
		{/if}
	</div>
	
	<div class="pen-controls">
		<label>
			Pen Color:
			<input 
				type="color" 
				bind:value={penColor} 
				disabled={disabled}
			/>
		</label>
		
		<label>
			Pen Size:
			<input 
				type="range" 
				min="1" 
				max="10" 
				bind:value={penSize}
				disabled={disabled}
			/>
			<span>{penSize}px</span>
		</label>
	</div>
</div>

<style>
	.signature-pad-container {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.signature-pad {
		position: relative;
		border: 2px solid #ddd;
		border-radius: 8px;
		overflow: hidden;
		cursor: crosshair;
	}

	.signature-pad.disabled {
		cursor: not-allowed;
		opacity: 0.6;
	}

	canvas {
		display: block;
		touch-action: none;
		user-select: none;
	}

	.placeholder {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		color: #999;
		font-size: 18px;
		font-style: italic;
		pointer-events: none;
	}

	.disabled-overlay {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.1);
		display: flex;
		align-items: center;
		justify-content: center;
		color: #666;
		font-weight: bold;
	}

	.toolbar {
		display: flex;
		gap: 10px;
		align-items: center;
		flex-wrap: wrap;
	}

	.btn {
		padding: 6px 12px;
		border: 1px solid #ddd;
		border-radius: 4px;
		background: #f5f5f5;
		cursor: pointer;
		font-size: 14px;
		transition: all 0.2s;
	}

	.btn:hover:not(:disabled) {
		background: #e5e5e5;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-clear {
		color: #dc3545;
	}

	.btn-undo {
		color: #007bff;
	}

	.status {
		margin-left: auto;
		font-size: 14px;
		color: #28a745;
		font-weight: 500;
	}

	.status.unsigned {
		color: #dc3545;
	}

	.required-indicator {
		color: #dc3545;
		font-size: 12px;
	}

	.pen-controls {
		display: flex;
		gap: 20px;
		align-items: center;
		flex-wrap: wrap;
		padding: 10px;
		background: #f8f9fa;
		border-radius: 4px;
	}

	.pen-controls label {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 14px;
		color: #666;
	}

	.pen-controls input[type="color"] {
		width: 40px;
		height: 30px;
		border: 1px solid #ddd;
		border-radius: 4px;
		cursor: pointer;
	}

	.pen-controls input[type="range"] {
		width: 100px;
	}
</style>
