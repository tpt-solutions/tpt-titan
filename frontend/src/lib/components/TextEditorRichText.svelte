<!-- frontend/src/lib/components/TextEditorRichText.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';

	export let editorElement = null;

	const dispatch = createEventDispatcher();

	let selectedFontFamily = 'Arial';
	let selectedFontSize = '12pt';
	let isBold = false;
	let isItalic = false;
	let isUnderline = false;
	let textAlign = 'left';

	function applyBold() {
		document.execCommand('bold');
		updateFormattingState();
	}

	function applyItalic() {
		document.execCommand('italic');
		updateFormattingState();
	}

	function applyUnderline() {
		document.execCommand('underline');
		updateFormattingState();
	}

	function applyFontFamily() {
		document.execCommand('fontName', false, selectedFontFamily);
	}

	function applyFontSize() {
		document.execCommand('fontSize', false, selectedFontSize.replace('pt', ''));
	}

	function applyAlignLeft() {
		document.execCommand('justifyLeft');
		textAlign = 'left';
	}

	function applyAlignCenter() {
		document.execCommand('justifyCenter');
		textAlign = 'center';
	}

	function applyAlignRight() {
		document.execCommand('justifyRight');
		textAlign = 'right';
	}

	function insertLink() {
		const url = prompt('Enter URL:');
		if (url) {
			document.execCommand('createLink', false, url);
		}
	}

	function updateFormattingState() {
		if (!editorElement) return;

		isBold = document.queryCommandState('bold');
		isItalic = document.queryCommandState('italic');
		isUnderline = document.queryCommandState('underline');

		selectedFontFamily = document.queryCommandValue('fontName') || 'Arial';
		const fontSizeValue = document.queryCommandValue('fontSize');
		selectedFontSize = fontSizeValue ? `${fontSizeValue}pt` : '12pt';

		if (document.queryCommandState('justifyLeft')) textAlign = 'left';
		else if (document.queryCommandState('justifyCenter')) textAlign = 'center';
		else if (document.queryCommandState('justifyRight')) textAlign = 'right';
	}

	function handleRichTextKeyDown(event) {
		if (event.ctrlKey || event.metaKey) {
			switch (event.key.toLowerCase()) {
				case 'b':
					event.preventDefault();
					applyBold();
					break;
				case 'i':
					event.preventDefault();
					applyItalic();
					break;
				case 'u':
					event.preventDefault();
					applyUnderline();
					break;
				case 'l':
					event.preventDefault();
					applyAlignLeft();
					break;
				case 'e':
					event.preventDefault();
					applyAlignCenter();
					break;
				case 'r':
					event.preventDefault();
					applyAlignRight();
					break;
				case 'k':
					event.preventDefault();
					insertLink();
					break;
			}
		}
	}

	function handleContentChange() {
		dispatch('change');
	}
</script>

<div class="border border-gray-300 rounded-lg overflow-hidden">
	<!-- Formatting Toolbar -->
	<div class="bg-gray-50 px-4 py-2 border-b border-gray-200">
		<div class="flex items-center space-x-2 text-sm">
			<select
				bind:value={selectedFontFamily}
				on:change={applyFontFamily}
				class="border border-gray-300 rounded px-2 py-1"
			>
				<option value="Arial">Arial</option>
				<option value="Times New Roman">Times New Roman</option>
				<option value="Inter">Inter</option>
				<option value="Georgia">Georgia</option>
				<option value="Verdana">Verdana</option>
			</select>

			<select
				bind:value={selectedFontSize}
				on:change={applyFontSize}
				class="border border-gray-300 rounded px-2 py-1"
			>
				<option value="12pt">12pt</option>
				<option value="14pt">14pt</option>
				<option value="16pt">16pt</option>
				<option value="18pt">18pt</option>
				<option value="24pt">24pt</option>
				<option value="32pt">32pt</option>
			</select>

			<div class="w-px h-6 bg-gray-300"></div>

			<button
				on:click={applyBold}
				class="p-1 hover:bg-white rounded {isBold ? 'bg-blue-100' : ''}"
				title="Bold"
			>
				<strong>B</strong>
			</button>

			<button
				on:click={applyItalic}
				class="p-1 hover:bg-white rounded {isItalic ? 'bg-blue-100' : ''}"
				title="Italic"
			>
				<em>I</em>
			</button>

			<button
				on:click={applyUnderline}
				class="p-1 hover:bg-white rounded {isUnderline ? 'bg-blue-100' : ''}"
				title="Underline"
			>
				<u>U</u>
			</button>

			<div class="w-px h-6 bg-gray-300"></div>

			<button
				on:click={applyAlignLeft}
				class="p-1 hover:bg-white rounded {textAlign === 'left' ? 'bg-blue-100' : ''}"
				title="Align Left"
			>
				⬅️
			</button>

			<button
				on:click={applyAlignCenter}
				class="p-1 hover:bg-white rounded {textAlign === 'center' ? 'bg-blue-100' : ''}"
				title="Align Center"
			>
				⬌
			</button>

			<button
				on:click={applyAlignRight}
				class="p-1 hover:bg-white rounded {textAlign === 'right' ? 'bg-blue-100' : ''}"
				title="Align Right"
			>
				➡️
			</button>

			<div class="w-px h-6 bg-gray-300"></div>

			<button
				on:click={insertLink}
				class="p-1 hover:bg-white rounded"
				title="Insert Link"
			>
				🔗
			</button>

			<button
				on:click={() => dispatch('openFindReplace')}
				class="p-1 hover:bg-white rounded"
				title="Find & Replace"
			>
				🔍
			</button>
		</div>
	</div>

	<!-- Editable Content Area -->
	<div
		bind:this={editorElement}
		contenteditable="true"
		class="min-h-96 p-4 focus:outline-none"
		placeholder="Start writing..."
		on:input={handleContentChange}
		on:keydown={handleRichTextKeyDown}
		on:mouseup={updateFormattingState}
		on:keyup={updateFormattingState}
	>
		<h1>Welcome to TPT Rich Text Editor</h1>
		<p>This is a traditional word processor interface with all the formatting options you expect.</p>
		<p>You can also insert <strong>natural math expressions</strong> like "integral from 0 to π of sin(x) dx" which will be automatically converted to proper mathematical notation.</p>
	</div>
</div>
