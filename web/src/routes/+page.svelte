<script lang="ts">
  interface Widget {
    id: string;
    title: string;
    content: string;
    colSpan: string;
  }

  // Svelte 5 Runes: Use $state() for reactive variables
  let widgets = $state<Widget[]>([
    { id: '1', title: 'Tasks & To-Dos', content: 'No pending tasks for today.', colSpan: 'col-span-1 md:col-span-2' },
    { id: '2', title: 'System Status', content: 'Proxmox & Go Backend: Healthy 🚀', colSpan: 'col-span-1' },
    { id: '3', title: 'Health & Macros', content: 'Calories: 2,100 / 2,500 kcal', colSpan: 'col-span-1' },
    { id: '4', title: 'Quick Notes / Wiki', content: 'Drafting LifeOS architecture specs...', colSpan: 'col-span-1 md:col-span-2' }
  ]);

  let draggedIndex = $state<number | null>(null);

  function handleDragStart(index: number) {
    draggedIndex = index;
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
  }

  function handleDrop(targetIndex: number) {
    if (draggedIndex === null || draggedIndex === targetIndex) return;
    
    // Reorder array directly (Svelte 5 $state tracks mutations)
    const item = widgets[draggedIndex];
    widgets.splice(draggedIndex, 1);
    widgets.splice(targetIndex, 0, item);
    
    draggedIndex = null;
  }
</script>

<main class="min-h-screen bg-slate-950 text-slate-100 p-6 md:p-10">
  <!-- Header (Fixed className -> class) -->
  <header class="mb-8 flex justify-between items-center">
    <div>
      <h1 class="text-3xl font-extrabold tracking-tight">Life OS</h1>
      <p class="text-sm text-slate-400 mt-1">Local-first modular command center</p>
    </div>
    <div class="text-xs bg-slate-900 border border-slate-800 px-3 py-1.5 rounded-full text-emerald-400 font-mono">
      ● Engine Live
    </div>
  </header>

  <!-- Responsive Dashboard Grid -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6 auto-rows-[minmax(160px,auto)]">
    {#each widgets as widget, index (widget.id)}
      <div
        role="region"
        aria-label={widget.title}
        draggable="true"
        ondragstart={() => handleDragStart(index)}
        ondragover={handleDragOver}
        ondrop={() => handleDrop(index)}
        class="{widget.colSpan} bg-slate-900/80 border border-slate-800/80 backdrop-blur-sm rounded-2xl p-5 flex flex-col justify-between shadow-xl cursor-grab active:cursor-grabbing hover:border-slate-700 transition-all duration-200 group"
      >
        <div class="flex items-center justify-between mb-3">
          <h2 class="font-semibold text-slate-200 tracking-wide">{widget.title}</h2>
          <div class="text-slate-600 group-hover:text-slate-400 transition-colors">
            ⠿
          </div>
        </div>

        <div class="text-slate-400 text-sm flex-grow flex items-center">
          <p>{widget.content}</p>
        </div>

        <div class="mt-4 pt-3 border-t border-slate-800/50 flex justify-end text-xs text-slate-500">
          <span>ID: {widget.id}</span>
        </div>
      </div>
    {/each}
  </div>
</main>

