
// import test from '/info/testimport.js'

// test()

var socket;
var ctx;

class PaintGame {
    constructor() {
        const canvas = document.getElementById("game-canvas");
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
        ctx = canvas.getContext('2d');
    }

    handleUpdate(e) {
        const data = JSON.parse(e.data);
        // Map stats to the Overlay UI
        if (data.PlayerCount !== undefined) {
            document.getElementById('stat-players').innerText = data.PlayerCount;
        }
        if (data.Status) {
            document.getElementById('stat-state').innerText = data.Status;
        }
    }

    terminate() {
        console.log("Terminating Game")
        window.removeEventListener('keydown', game.toggleUI);
        if (socket) {
            socket.onmessage = null;
            socket.close();
        }
    }

    toggleUI() {
        const overlay = document.getElementById('game-overlay');
        if (overlay) overlay.classList.toggle('hidden');
    }
}

let game;

export function STARTGAME() {
    game = new PaintGame();
    // Assuming wsHost, gameType, and roomId are globally available from injector
    socket = new WebSocket(`${wsHost}/ws?type=${gameType}&room=${roomId}`);
    socket.onmessage = (e) => game.handleUpdate(e);

    window.addEventListener('keydown', game.toggleUI)
    console.log("Starting Game")

}

export function TERMINATE() {
    if (game) game.terminate();
}
