
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

        this.keyboard = new Set();
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
        window.removeEventListener('keydown', game.keydown);
        window.removeEventListener('keyup', game.keyup);
        if (socket) {
            socket.onmessage = null;
            socket.close();
        }
    }

    keydown(e) {
        if (e.key === 'Escape') {
            const overlay = document.getElementById('game-overlay');
            if (overlay) overlay.classList.toggle('hidden');
            return
        };

        let cmd = new Uint8Array(1);
        if (e.key === 'w' && !this.keyboard.has(e.key)) {cmd[0] = 0b001; this.keyboard.add(e.key);}
        if (e.key === 'a' && !this.keyboard.has(e.key)) {cmd[0] = 0b011; this.keyboard.add(e.key);}
        if (e.key === 's' && !this.eyboard.has(e.key)) {cmd[0] = 0b101; this.keyboard.add(e.key);}
        if (e.key === 'd' && !this.keyboard.has(e.key)) {cmd[0] = 0b111; this.keyboard.add(e.key);}
        if (cmd[0] > 0) socket.send(cmd);
    }
    keyup(e) {
        let cmd = new Uint8Array(1);
        cmd[0] = 0b1000;
        if (e.key === 'w' && this.keyboard.has(e.key)) {cmd[0] = 0b000; this.keyboard.delete(e.key);}
        if (e.key === 'a' && this.keyboard.has(e.key)) {cmd[0] = 0b010; this.keyboard.delete(e.key);}
        if (e.key === 's' && this.keyboard.has(e.key)) {cmd[0] = 0b100; this.keyboard.delete(e.key);}
        if (e.key === 'd' && this.keyboard.has(e.key)) {cmd[0] = 0b110; this.keyboard.delete(e.key);}
        if (cmd[0] != 0b1000) socket.send(cmd);
    }
}

let game;

export function STARTGAME() {
    game = new PaintGame();
    // Assuming wsHost, gameType, and roomId are globally available from injector
    socket = new WebSocket(`${wsHost}/ws?type=${gameType}&room=${roomId}`);
    socket.onmessage = (e) => game.handleUpdate(e);

    console.log("Starting Game")
    window.addEventListener('keydown', game.keydown)
    window.addEventListener("keyup", game.keyup);

}

export function TERMINATE() {
    if (game) game.terminate();
}
