<script>
	import TaskBoard from '$lib/components/TaskBoard.svelte';
	import TaskForm from '$lib/components/TaskForm.svelte';
	import { onMount } from 'svelte';
	import { apiGet, apiPost, apiPut, apiDelete } from '$lib/api.js';

	// Accept framework-provided props to avoid warnings
	export let data = null;
	export let form = null;
	export let params = null;


	let currentView = 'board'; // 'board' or 'form'
	let selectedTask = null;
	let tasks = [];
	let projects = [];
	let isLoading = false;
	let loadError = null;

	onMount(() => {
		loadTasks();
		loadProjects();
	});

	async function loadTasks() {
		isLoading = true;
		loadError = null;
		try {
			const res = await apiGet('/tasks');
			tasks = (res.tasks || []).map(t => ({
				...t,
				// Backend returns ISO strings; keep dueDate as-is for display
				dueDate: t.dueDate ? new Date(t.dueDate) : null,
				createdAt: t.createdAt ? new Date(t.createdAt) : null,
				subtasks: (t.subtasks || []).map(st => ({
					id: st.id,
					title: st.title,
					completed: !!st.completed
				}))
			}));
		} catch (err) {
			loadError = 'Could not load tasks. Please check your connection.';
		} finally {
			isLoading = false;
		}
	}

	async function loadProjects() {
		try {
			const res = await apiGet('/tasks/projects');
			projects = (res.projects || []).map(p => ({
				id: p.id,
				name: p.name,
				color: p.color || 'blue'
			}));
		} catch (err) {
			console.error('Failed to load projects:', err);
		}
	}

	function createNewTask() {
		selectedTask = null;
		currentView = 'form';
	}

	function editTask(task) {
		selectedTask = task;
		currentView = 'form';
	}

	async function saveTask(taskData) {
		try {
			if (selectedTask) {
				const res = await apiPut(`/tasks/${selectedTask.id}`, taskData);
				const updated = res.task;
				tasks = tasks.map(t => t.id === selectedTask.id ? normalizeTask(updated) : t);
			} else {
				const res = await apiPost('/tasks', taskData);
				tasks = [normalizeTask(res.task), ...tasks];
			}
		} catch (err) {
			loadError = 'Failed to save task. Please try again.';
			return;
		}
		currentView = 'board';
		selectedTask = null;
	}

	async function deleteTask(taskId) {
		try {
			await apiDelete(`/tasks/${taskId}`);
			tasks = tasks.filter(t => t.id !== taskId);
		} catch (err) {
			loadError = 'Failed to delete task. Please try again.';
		}
	}

	async function updateTaskStatus(taskId, newStatus) {
		const task = tasks.find(t => t.id === taskId);
		if (!task) return;
		// Optimistic update
		tasks = tasks.map(t => t.id === taskId ? { ...t, status: newStatus } : t);
		try {
			const res = await apiPut(`/tasks/${taskId}/status`, { status: newStatus });
			tasks = tasks.map(t => t.id === taskId ? normalizeTask(res.task) : t);
		} catch (err) {
			// Revert on failure
			tasks = tasks.map(t => t.id === taskId ? { ...t, status: task.status } : t);
		}
	}

	function normalizeTask(t) {
		return {
			...t,
			dueDate: t.dueDate ? new Date(t.dueDate) : null,
			createdAt: t.createdAt ? new Date(t.createdAt) : null,
			subtasks: (t.subtasks || []).map(st => ({
				id: st.id,
				title: st.title,
				completed: !!st.completed
			}))
		};
	}

	function getProjectName(projectId) {
		const project = projects.find(p => p.id === projectId);
		return project ? project.name : 'No Project';
	}

	function getTasksByStatus(status) {
		return tasks.filter(t => t.status === status);
	}

	function getTaskStats() {
		const total = tasks.length;
		const completed = tasks.filter(t => t.status === 'done').length;
		const overdue = tasks.filter(t => t.dueDate && t.dueDate < new Date() && t.status !== 'done').length;
		const highPriority = tasks.filter(t => t.priority === 'high' && t.status !== 'done').length;

		return { total, completed, overdue, highPriority };
	}
</script>

<svelte:head>
	<title>Tasks & Projects - TPT Titan</title>
</svelte:head>

<div class="h-screen flex flex-col bg-gray-50">
	<!-- Header -->
	<header class="flex items-center justify-between px-6 py-3 border-b border-gray-200 bg-white">
		<div class="flex items-center space-x-4">
			<h1 class="text-xl font-semibold text-gray-900">Tasks & Projects</h1>
			<span class="text-sm text-gray-500">AI-powered project management</span>
		</div>

		<div class="flex items-center space-x-2">
			{#if currentView === 'board'}
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					on:click={createNewTask}
				>
					New Task
				</button>
			{:else}
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
					on:click={() => { currentView = 'board'; selectedTask = null; }}
				>
					Back to Board
				</button>
			{/if}
		</div>
	</header>

	<!-- Error banner -->
	{#if loadError}
		<div class="bg-red-50 border-b border-red-200 px-6 py-3 flex items-center gap-3">
			<span class="text-sm text-red-700">{loadError}</span>
			<button on:click={loadTasks} class="ml-auto text-sm text-red-600 underline hover:no-underline">Retry</button>
		</div>
	{/if}

	<!-- Stats Bar -->
	<div class="bg-white border-b border-gray-200 px-6 py-4">
		<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
			<div class="text-center">
				<div class="text-2xl font-bold text-blue-600">{getTaskStats().total}</div>
				<div class="text-sm text-gray-600">Total Tasks</div>
			</div>
			<div class="text-center">
				<div class="text-2xl font-bold text-green-600">{getTaskStats().completed}</div>
				<div class="text-sm text-gray-600">Completed</div>
			</div>
			<div class="text-center">
				<div class="text-2xl font-bold text-red-600">{getTaskStats().overdue}</div>
				<div class="text-sm text-gray-600">Overdue</div>
			</div>
			<div class="text-center">
				<div class="text-2xl font-bold text-orange-600">{getTaskStats().highPriority}</div>
				<div class="text-sm text-gray-600">High Priority</div>
			</div>
		</div>
	</div>

	<!-- Main Content -->
	<div class="flex-1 overflow-hidden">
		{#if isLoading}
			<div class="flex items-center justify-center h-full text-gray-400">
				<span>Loading tasks…</span>
			</div>
		{:else if currentView === 'board'}
			<TaskBoard
				{tasks}
				{projects}
				{getProjectName}
				{getTasksByStatus}
				on:edit={e => editTask(e.detail)}
				on:delete={e => deleteTask(e.detail)}
				on:updateStatus={e => updateTaskStatus(e.detail.taskId, e.detail.status)}
			/>
		{:else}
			<TaskForm
				task={selectedTask}
				{projects}
				on:save={e => saveTask(e.detail)}
				on:cancel={() => { currentView = 'board'; selectedTask = null; }}
			/>
		{/if}
	</div>
</div>
