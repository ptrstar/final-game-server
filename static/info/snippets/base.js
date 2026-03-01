window.__GAMEHTML__ = `
<div class="w-full h-full bg-slate-900 text-white p-8 overflow-y-auto">
    <div class="max-w-6xl mx-auto">
        <header class="mb-10 flex justify-between items-center border-b border-slate-700 pb-6">
            <div>
                <h1 class="text-4xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-emerald-400">
                    Overview play.trojanos.ch
                </h1>
                <p class="text-slate-400 mt-2">Currently active game rooms on this instance.</p>
            </div>
            <div id="connection-status" class="px-3 py-1 rounded-full bg-emerald-500/20 text-emerald-400 text-sm border border-emerald-500/50">
                Connected
            </div>
        </header>
        <button
            onclick="game_injector('paint','')"
            class="p-2 mb-2 py-2.5 rounded-xl font-bold text-sm transition-all bg-violet-600 hover:bg-violet-500 text-white shadow-lg shadow-violet-900/20">
            Paint
        </button>
        <div id="game-list" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <!-- Cards will be injected here -->
        </div>
    </div>
</div>`;