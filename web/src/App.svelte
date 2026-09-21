<!-- web/src/App.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import TaskWidget from './lib/components/TaskWidget.svelte';
	import NotesWidget from './lib/components/NotesWidget.svelte';

	// Live Clock State
	let currentTime = $state('');
	let currentDate = $state('');

	function updateTime() {
		const now = new Date();
		currentTime = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
		currentDate = now.toLocaleDateString([], { weekday: 'short', month: 'short', day: 'numeric' });
	}

	// Active Sidebar Tab State
	let tabs = $state(['Overview', 'Workspace', 'Planning', 'Server']);
	let activeTab = $state('Overview');

	onMount(() => {
		updateTime();
		const interval = setInterval(updateTime, 1000);
		return () => clearInterval(interval);
	});
</script>

<div class="flex h-screen w-screen flex-col overflow-hidden bg-neutral-950 text-neutral-100 select-none">
	<!-- TOP BAR -->
	<header class="flex h-14 items-center justify-between border-b border-neutral-800/80 px-4 bg-neutral-900/30">
		<!-- Left: Add Box / Tab Action -->
		<div class="flex items-center gap-2">
			<button class="flex items-center gap-2 rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-1.5 text-xs font-medium text-neutral-300 hover:bg-neutral-800 hover:text-white transition">
				<span class="text-sm font-bold">+</span> Add Widget
			</button>
		</div>

		<!-- Middle: Live System Time & Date -->
		<div class="flex flex-col items-center">
			<span class="font-mono text-sm font-semibold tracking-wider text-neutral-100">{currentTime}</span>
			<span class="text-[10px] text-neutral-500 uppercase tracking-widest">{currentDate}</span>
		</div>

		<!-- Right: Account, Settings, Notifications -->
		<div class="flex items-center gap-2">
			<button class="flex h-8 w-8 items-center justify-center rounded-full border border-neutral-800 bg-neutral-900 text-neutral-400 hover:text-neutral-100 transition">
				⚙️
			</button>
			<button class="flex h-8 w-8 items-center justify-center rounded-full border border-neutral-800 bg-neutral-900 text-neutral-400 hover:text-neutral-100 transition">
				🔔
			</button>
			<div class="ml-2 flex h-8 w-8 items-center justify-center rounded-full border border-indigo-500/30 bg-indigo-500/10 text-xs font-bold text-indigo-400">
				OS
			</div>
		</div>
	</header>

	<!-- MAIN WORKSPACE LAYOUT -->
	<div class="flex flex-1 overflow-hidden">
		<!-- SIDEBAR: Tabs Navigation -->
		<aside class="w-48 border-r border-neutral-800/80 bg-neutral-900/20 p-3 flex flex-col gap-2">
			<span class="px-2 text-[10px] font-semibold tracking-wider text-neutral-500 uppercase">Spaces</span>
			{#each tabs as tab}
				<button
					onclick={() => (activeTab = tab)}
					class="w-full text-left rounded-lg px-3 py-2 text-xs font-medium transition {activeTab === tab
						? 'bg-neutral-800 text-neutral-100 border border-neutral-700/50'
						: 'text-neutral-400 hover:bg-neutral-900 hover:text-neutral-200'}"
				>
					{tab}
				</button>
			{/each}
		</aside>

		<!-- MAIN CONTENT: 3x2 Grid -->
		<main class="flex-1 overflow-y-auto p-4">
			<div class="grid h-full grid-cols-3 grid-rows-2 gap-4">
				<!-- Box A: Tasks Module -->
				<TaskWidget />

				<!-- Box B: Notes / Scratchpad Module -->
				<NotesWidget />

				<!-- Box C: System Monitor Placeholder -->
				<div class="flex h-full flex-col justify-between rounded-xl border border-neutral-800/80 bg-neutral-900/40 p-4">
					<div class="flex items-center justify-between">
						<h3 class="font-medium text-neutral-300">System Monitor</h3>
						<span class="h-2 w-2 rounded-full bg-emerald-500"></span>
					</div>
					<div class="text-xs text-neutral-500">
						Host CPU & Memory Metrics (Ready for core.sysmon)
					</div>
				</div>

				<!-- Box D: Quick Launch Placeholder -->
				<div class="flex h-full flex-col justify-between rounded-xl border border-neutral-800/80 bg-neutral-900/40 p-4">
					<h3 class="font-medium text-neutral-300">Quick Launch</h3>
					<div class="text-xs text-neutral-500">
						HomeLab Bookmarks & Services
					</div>
				</div>

				<!-- Box E: Calendar / Schedule Placeholder -->
				<div class="flex h-full flex-col justify-between rounded-xl border border-neutral-800/80 bg-neutral-900/40 p-4">
					<h3 class="font-medium text-neutral-300">Calendar</h3>
					<div class="text-xs text-neutral-500">
						Upcoming events & schedule
					</div>
				</div>

				<!-- Box F: Quick Actions Placeholder -->
				<div class="flex h-full flex-col justify-between rounded-xl border border-neutral-800/80 bg-neutral-900/40 p-4">
					<h3 class="font-medium text-neutral-300">Quick Actions</h3>
					<div class="text-xs text-neutral-500">
						Run terminal commands & scripts
					</div>
				</div>
			</div>
		</main>
	</div>
</div>