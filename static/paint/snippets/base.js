window.__GAMEHTML__ = `
<canvas id="game-canvas" class="fixed inset-0 z-0 bg-[#0f172a]"></canvas>

<!-- Overlay Menu -->
<div id="game-overlay" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-md transition-all">
    <div class="relative w-full max-w-md p-8 bg-slate-800/90 border border-slate-700 rounded-3xl shadow-2xl text-slate-200">
        
        <!-- Top Right Exit -->
        <button onclick="game_injector('info', '')" class="absolute top-4 right-4 p-2 bg-red-500/20 hover:bg-red-500 text-red-400 hover:text-white rounded-full transition-all">
           Exit Game 
        </button>

        <h1 class="text-4xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-emerald-400 mb-2">Paint</h1>
        
        <div id="game-stats" class="grid grid-cols-2 gap-4 mb-6 text-sm py-3 border-y border-slate-700/50">
            <div><span class="text-slate-400">Players:</span> <span id="stat-players" class="text-white font-bold">--</span></div>
            <div><span class="text-slate-400">Status:</span> <span id="stat-state" class="text-emerald-400 font-bold">Waiting</span></div>
        </div>

        <div class="space-y-4">
            <h3 class="font-bold text-blue-400 uppercase tracking-widest text-xs">Objective</h3>
            <p class="text-sm text-slate-300 leading-relaxed">Cover the canvas in your color. Most surface area wins!</p>
            
            <h3 class="font-bold text-blue-400 uppercase tracking-widest text-xs mt-4">Controls</h3>
            <div class="grid grid-cols-2 gap-2 text-xs">
                <div class="flex items-center gap-2"><kbd class="px-2 py-1 rounded bg-slate-700 text-white">WASD</kbd> <span>Move</span></div>
                <div class="flex items-center gap-2"><kbd class="px-2 py-1 rounded bg-slate-700 text-white">SPACE</kbd> <span>Fire</span></div>
            </div>

            <div class="mt-6 flex items-center gap-2 p-3 bg-slate-900/50 rounded-xl border border-slate-700/50 text-xs italic text-slate-400">
                <span class="animate-pulse">●</span> Press <kbd class="mx-1 px-1.5 py-0.5 rounded bg-slate-700 text-white not-italic">ESC</kbd> to toggle this menu
            </div>
        </div>
    </div>
</div>`;


