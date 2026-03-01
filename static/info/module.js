import test from '/info/testimport.js'

test()

var game;
var socket;

class GameInfoLobby {
    constructor() {
        this.container = document.getElementById('game-list');
        this.statusBadge = document.getElementById('connection-status');
        // Change grid to a single column to hold the horizontal rows
        this.container.className = "flex flex-col gap-12"; 
        console.log("Lobby Monitor Started");
    }

    terminate() {
        console.log("Terminating Game")
        if (socket) {
            socket.onmessage = null;
            socket.close();
        }
    }

    updateStatus() {
        if (!socket || !this.statusBadge) return;
        const states = {
            [WebSocket.CONNECTING]: { text: 'Connecting...', class: 'bg-yellow-500/20 text-yellow-400 border-yellow-500/50' },
            [WebSocket.OPEN]: { text: 'Connected', class: 'bg-emerald-500/20 text-emerald-400 border-emerald-500/50' },
            [WebSocket.CLOSING]: { text: 'Closing...', class: 'bg-slate-500/20 text-slate-400 border-slate-500/50' },
            [WebSocket.CLOSED]: { text: 'Disconnected', class: 'bg-red-500/20 text-red-400 border-red-500/50' }
        };
        const current = states[socket.readyState];
        this.statusBadge.innerText = current.text;
        this.statusBadge.onclick = () => window.location.reload();
        this.statusBadge.className = `px-3 py-1 rounded-full text-sm border ${current.class}`;
    }

    async handleUpdate(event) {
        try {
            const text = await event.data.text();
            const rooms = JSON.parse(text);
            this.updateStatus();
            this.smartRender(rooms);
        } catch (e) {
            console.error("Parse error", e);
        }
    }

    smartRender(rooms) {
        if (!Array.isArray(rooms)) return;

        const grouped = rooms.reduce((acc, room) => {
            if (!acc[room.type]) acc[room.type] = [];
            acc[room.type].push(room);
            return acc;
        }, {});

        // 1. Manage Type Sections
        Object.entries(grouped).forEach(([type, items]) => {
            let section = document.getElementById(`section-${type}`);
            if (!section) {
                section = this.createSectionElement(type);
                this.container.appendChild(section);
            }
            
            const scrollContainer = section.querySelector('.scroll-container');
            const foundIds = new Set();

            // 2. Manage Individual Cards
            items.forEach(room => {
                foundIds.add(room.id);
                let card = document.getElementById(`card-${room.id}`);
                if (!card) {
                    card = document.createElement('div');
                    card.id = `card-${room.id}`;
                    card.className = "flex-none w-72 snap-start";
                    scrollContainer.appendChild(card);
                }
                // Only update content, preserving the element/scroll position
                card.innerHTML = this.getCardTemplate(room);
            });

            // 3. Remove stale cards (rooms that closed)
            Array.from(scrollContainer.children).forEach(child => {
                const id = child.id.replace('card-', '');
                if (!foundIds.has(id)) child.remove();
            });
        });
    }

    createSectionElement(type) {
        const div = document.createElement('section');
        div.id = `section-${type}`;
        div.innerHTML = `
            <div class="flex items-center mb-4 px-2">
                <h2 class="text-2xl font-bold capitalize text-slate-200">${type} Games</h2>
            </div>
            <div class="scroll-container flex overflow-x-auto pb-6 gap-6 scrollbar-hide snap-x"></div>
        `;
        return div;
    }

    getCardTemplate(room) {
        const isFull = room.player_count >= room.capacity;
        const canJoin = room.can_join && !isFull;
        const percent = Math.min(100, (room.player_count / room.capacity) * 100);

        return `
            <div class="bg-slate-800/50 border border-slate-700/50 rounded-2xl p-6 
                        hover:bg-slate-800 hover:border-blue-500/50 transition-all duration-300 group shadow-lg">
                <div class="flex justify-between items-center mb-6">
                    <span class="text-blue-400 font-bold">${room.type.toUpperCase()}</span>
                    <span class="text-[10px] font-mono text-slate-500">${room.id.slice(-6)}</span>
                </div>
                <h3 class="text-lg font-bold text-slate-100 mb-1">Room ${room.id.split('-').pop()}</h3>
                <div class="flex items-center gap-2 mb-6">
                    <div class="h-2 w-2 rounded-full ${canJoin ? 'bg-emerald-500 animate-pulse' : 'bg-red-500'}"></div>
                    <span class="text-xs text-slate-400 italic">${room.status}</span>
                </div>
                <div class="space-y-4">
                    <div class="flex justify-between text-sm">
                        <span class="text-slate-500">Players</span>
                        <span class="text-slate-200 font-semibold">${room.player_count}/${room.capacity}</span>
                    </div>
                    <div class="w-full bg-slate-900 rounded-full h-1.5 overflow-hidden">
                        <div class="bg-blue-500 h-full transition-all duration-500" style="width: ${percent}%"></div>
                    </div>
                    <button ${!canJoin ? 'disabled' : ''} 
                        onclick="game_injector('${room.type}','${room.id}')"
                        class="w-full py-2.5 rounded-xl font-bold text-sm transition-all
                        ${canJoin ? 'bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-900/20' : 'bg-slate-700 text-slate-500 cursor-not-allowed'}">
                        ${canJoin ? 'Connect' : 'Full'}
                    </button>
                </div>
            </div>`;
    }
}

export function STARTGAME() {
    game = new GameInfoLobby();
    socket = new WebSocket(`${wsHost}/ws?type=${gameType}&room=${roomId}`);
    socket.onopen = () => game.updateStatus();
    socket.onclose = () => game.updateStatus();
    socket.onerror = () => game.updateStatus();
    socket.onmessage = (e) => game.handleUpdate(e);
}

export function TERMINATE() {
    game.terminate();
}
