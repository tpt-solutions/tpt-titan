<script>
	import TaskBoard from '$lib/components/TaskBoard.svelte';
	import TaskForm from '$lib/components/TaskForm.svelte';
	import { onMount } from 'svelte';

	// Accept framework-provided props to avoid warnings
	export let data = null;
	export let form = null;

	let currentView = 'board'; // 'board' or 'form'
	let selectedTask = null;
	let tasks = [];
	let projects = [];

	onMount(() => {
		loadTasks();
		loadProjects();
	});

	function loadTasks() {
		// Mock data - will connect to encrypted API later
		tasks = [
			{
				id: 1,
				title: 'Review Q4 Sales Data',
				description: 'Analyze spreadsheet data and prepare quarterly report',
				status: 'in-progress',
				priority: 'high',
				projectId: 1,
				dueDate: new Date('2024-01-15'),
				assignedTo: 'John Doe',
				tags: ['analysis', 'quarterly'],
				createdAt: new Date('2024-01-01'),
				subtasks: [
					{ id: 1, title: 'Import sales data', completed: true },
					{ id: 2, title: 'Create pivot tables', completed: false },
					{ id: 3, title: 'Generate charts', completed: false }
				]
			},
			{
				id: 2,
				title: 'Update Customer Survey',
				description: 'Incorporate feedback from recent form responses',
				status: 'todo',
				priority: 'medium',
				projectId: 1,
				dueDate: new Date('2024-01-20'),
				assignedTo: 'Jane Smith',
				tags: ['forms', 'customer-feedback'],
				createdAt: new Date('2024-01-05'),
				subtasks: []
			},
			{
				id: 3,
				title: 'Prepare Presentation Slides',
				description: 'Create slides for quarterly business review',
				status: 'done',
				priority: 'high',
				projectId: 2,
				dueDate: new Date('2024-01-10'),
				assignedTo: 'Bob Wilson',
				tags: ['presentation', 'business-review'],
				createdAt: new Date('2024-01-02'),
				subtasks: [
					{ id: 1, title: 'Gather data', completed: true },
					{ id: 2, title: 'Create slides', completed: true },
					{ id: 3, title: 'Practice presentation', completed: true }
				]
			}
		];
	}

	function loadProjects() {
		projects = [
			{ id: 1, name: 'Q4 Business Analysis', color: 'blue', tasks: 8 },
			{ id: 2, name: 'Customer Experience', color: 'green', tasks: 12 },
			{ id: 3, name: 'Product Development', color: 'purple', tasks: 5 }
		];
	}

	function createNewTask() {
		selectedTask = null;
		currentView = 'form';
	}

	function editTask(task) {
		selectedTask = task;
		currentView = 'form';
	}

	function saveTask(taskData) {
		if (selectedTask) {
			// Update existing task
			tasks = tasks.map(t => t.id === selectedTask.id ? { ...t, ...taskData } : t);
		} else {
			// Create new task
			const newTask = {
				id: Date.now(),
				...taskData,
				status: 'todo',
				createdAt: new Date(),
				subtasks: []
			};
			tasks = [...tasks, newTask];
		}
		currentView = 'board';
		selectedTask = null;
	}

	function deleteTask(taskId) {
		tasks = tasks.filter(t => t.id !== taskId);
	}

	function updateTaskStatus(taskId, newStatus) {
		tasks = tasks.map(t => t.id === taskId ? { ...t, status: newStatus } : t);
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
		{#if currentView === 'board'}
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
