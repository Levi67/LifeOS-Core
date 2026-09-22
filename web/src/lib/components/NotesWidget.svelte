<!-- web/src/lib/components/NotesWidget.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';

	interface Note {
		id: number;
		title: string;
		content: string;
		pinned: boolean;
	}

	let notes = $state<Note[]>([]);
	let selectedNote = $state<Note | null>(null);

	async function fetchNotes() {
		try {
			const res = await fetch('/api/v1/notes');
			if (res.ok) {
				notes = await res.json();
				if (notes.length > 0 && !selectedNote) {
					selectedNote = notes[0];
				}
			}
		} catch (e) {
			console.error(e);
		}
	}

	async function createNote() {
		const res = await fetch('/api/v1/notes', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ title: 'Scratchpad', content: '' })
		});

		if (res.ok) {
			const note = await res.json();
			notes = [note, ...notes];
			selectedNote = note;
		}
	}

	async function updateNote() {
		if (!selectedNote) return;

		await fetch('/api/v1/notes', {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				id: selectedNote.id,
				title: selectedNote.title,
				content: selectedNote.content,
				pinned: selectedNote.pinned
			})
		});
	}

	onMount(fetchNotes);
</script>

<div class="flex h-full flex-col rounded-xl border border-neutral-800 bg-neutral-900/60 p-4 backdrop-blur-md">
	<div class="mb-3 flex items-center justify-between">
		<h3 class="font-medium text-neutral-200">Scratchpad</h3>
		<button
			onclick={createNote}
			class="rounded-md border border-neutral-800 bg-neutral-950 px-2 py-1 text-xs text-neutral-400 hover:text-neutral-100 transition"
		>
			+ Note
		</button>
	</div>

	{#if selectedNote}
		<div class="flex flex-1 flex-col gap-2">
			<input
				type="text"
				bind:value={selectedNote.title}
				oninput={updateNote}
				placeholder="Note title..."
				class="w-full bg-transparent font-medium text-xs text-neutral-200 border-b border-neutral-800 pb-1 focus:outline-none focus:border-neutral-600"
			/>
			<textarea
				bind:value={selectedNote.content}
				oninput={updateNote}
				placeholder="Type thoughts or logs here..."
				class="w-full flex-1 resize-none bg-transparent text-xs text-neutral-300 focus:outline-none"
			></textarea>
		</div>
	{:else}
		<p class="text-xs text-neutral-600">No active note selected.</p>
	{/if}
</div>