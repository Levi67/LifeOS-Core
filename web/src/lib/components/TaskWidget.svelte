<!-- web/src/lib/components/TaskWidget.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';

	interface Task {
		id: number;
		title: string;
		completed: boolean;
	}

	let tasks = $state<Task[]>([]);
	let newTitle = $state('');

	async function fetchTasks() {
		try {
			const res = await fetch('/api/v1/tasks');
			if (res.ok) tasks = await res.json();
		} catch (e) {
			console.error(e);
		}
	}

	async function addTask(e: SubmitEvent) {
		e.preventDefault();
		if (!newTitle.trim()) return;

		const res = await fetch('/api/v1/tasks', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ title: newTitle.trim() })
		});

		if (res.ok) {
			const task = await res.json();
			tasks = [task, ...tasks];
			newTitle = '';
		}
	}

	async function toggleTask(task: Task) {
		const res = await fetch('/api/v1/tasks', {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ id: task.id, completed: !task.completed })
		});

		if (res.ok) {
			tasks = tasks.map((t) => (t.id === task.id ? { ...t, completed: !t.completed } : t));
		}
	}

	async function deleteTask(id: number) {
		const res = await fetch('/api/v1/tasks', {
			method: 'DELETE',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ id })
		});

		if (res.ok) {
			tasks = tasks.filter((t) => t.id !== id);
		}
	}

	onMount(fetchTasks);
</script>

<div class="flex h-full flex-col rounded-xl border border-neutral-800 bg-neutral-900/60 p-4 backdrop-blur-md">
	<div class="mb-3 flex items-center justify-between">
		<h3 class="font-medium text-neutral-200">Tasks</h3>
		<span class="text-xs text-neutral-500">{tasks.filter((t) => !t.completed).length} open</span>
	</div>

	<form onsubmit={addTask} class="mb-3 flex gap-2">
		<input
			type="text"
			bind:value={newTitle}
			placeholder="Add task..."
			class="w-full rounded-lg border border-neutral-800 bg-neutral-950 px-3 py-1.5 text-xs text-neutral-200 focus:border-neutral-600 focus:outline-none"
		/>
		<button
			type="submit"
			class="rounded-lg bg-neutral-100 px-3 py-1.5 text-xs font-semibold text-neutral-900 transition hover:bg-neutral-300"
		>
			+
		</button>
	</form>

	<div class="flex-1 overflow-y-auto space-y-1.5 pr-1">
		{#if tasks.length === 0}
			<p class="text-xs text-neutral-600">No tasks yet.</p>
		{:else}
			{#each tasks as task (task.id)}
				<div class="group flex items-center justify-between rounded-lg border border-neutral-800/40 bg-neutral-950/40 px-3 py-2 text-xs">
					<label class="flex items-center gap-2 cursor-pointer select-none">
						<input
							type="checkbox"
							checked={task.completed}
							onchange={() => toggleTask(task)}
							class="rounded border-neutral-700 bg-neutral-900 text-neutral-100 focus:ring-0"
						/>
						<span class={task.completed ? 'line-through text-neutral-600' : 'text-neutral-300'}>
							{task.title}
						</span>
					</label>
					<button
						onclick={() => deleteTask(task.id)}
						class="opacity-0 group-hover:opacity-100 text-neutral-500 hover:text-red-400 transition"
					>
						✕
					</button>
				</div>
			{/each}
		{/if}
	</div>
</div>