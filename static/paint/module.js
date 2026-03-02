
// import test from '/info/testimport.js'

// test()

class PaintGame {
    constructor() {
        this.canvas = document.getElementById("game-canvas");
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
        this.ctx = canvas.getContext('2d');

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
        if (this.socket) {
            this.socket.onmessage = null;
            this.socket.close();
        }
    }

    keydown(e) {
        if (e.key === 'Escape') {
            const overlay = document.getElementById('game-overlay');
            if (overlay) overlay.classList.toggle('hidden');
            return
        };

        let cmd = new Uint8Array(1);
        if (e.key === 'w' && !game.keyboard.has(e.key)) {cmd[0] = 0b001; game.keyboard.add(e.key);}
        if (e.key === 'a' && !game.keyboard.has(e.key)) {cmd[0] = 0b011; game.keyboard.add(e.key);}
        if (e.key === 's' && !game.keyboard.has(e.key)) {cmd[0] = 0b101; game.keyboard.add(e.key);}
        if (e.key === 'd' && !game.keyboard.has(e.key)) {cmd[0] = 0b111; game.keyboard.add(e.key);}
        if (cmd[0] > 0) this.socket.send(cmd);
    }
    keyup(e) {
        let cmd = new Uint8Array(1);
        cmd[0] = 0b1000;
        if (e.key === 'w' && game.keyboard.has(e.key)) {cmd[0] = 0b000; game.keyboard.delete(e.key);}
        if (e.key === 'a' && game.keyboard.has(e.key)) {cmd[0] = 0b010; game.keyboard.delete(e.key);}
        if (e.key === 's' && game.keyboard.has(e.key)) {cmd[0] = 0b100; game.keyboard.delete(e.key);}
        if (e.key === 'd' && game.keyboard.has(e.key)) {cmd[0] = 0b110; game.keyboard.delete(e.key);}
        if (cmd[0] != 0b1000) this.socket.send(cmd);
    }
}

let game;

export function STARTGAME() {
    game = new PaintGame();
    // Assuming wsHost, gameType, and roomId are globally available from injector
    game.socket = new WebSocket(`${wsHost}/ws?type=${gameType}&room=${roomId}`);
    game.socket.onmessage = (e) => game.handleUpdate(e);

    console.log("Starting Game")
    window.addEventListener('keydown', game.keydown)
    window.addEventListener("keyup", game.keyup);

}

export function TERMINATE() {
    if (game) game.terminate();
}
