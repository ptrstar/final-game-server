const socket = new WebSocket(`ws://${src}/ws?type=${game_type}&room=ANY`);
socket.binaryType = "arraybuffer";

class Game {
    constructor() {
        this.players = [];
        this.canvas = document.getElementById('game-canvas');
        this.uiLayer = document.getElementById('game-ui')
        this.ctx = this.canvas.getContext("2d");

        
        console.log("Starting Client Game");

        window.addEventListener('resize', this.upscaleCanvas);
        this.upscaleCanvas();
        this.main();
    }

    main() {

        this.gameUpdate();

        this.render();

        window.requestAnimationFrame(() => {game.main()});
    }

    gameUpdate() {
        this.players.forEach(p => {
            if (p.keyboard.has('w')) p.y -= 5;
            if (p.keyboard.has('a')) p.x -= 5;
            if (p.keyboard.has('s')) p.y += 5;
            if (p.keyboard.has('d')) p.x += 5;
        })
    }

    render() {
        this.ctx.fillStyle = '#43cae8';
        this.ctx.fillRect(0,0,this.canvas.clientWidth, this.canvas.clientHeight);

        this.players.forEach(p => {
            p.render(); 
        })
    }
    upscaleCanvas() {
        const dpi = window.devicePixelRatio || 1;
        const rect = this.canvas.getBoundingClientRect();
        this.canvas.width = rect.width * dpi;
        this.canvas.height = rect.height * dpi;
        this.ctx.scale(dpi, dpi);
    }

    handleUpdate(event) {
        const view = new DataView(event.data);
        const playerCount = view.getUint8(0);
        const newPlayers = [];
        for (let i = 0; i < playerCount; i++) {
            const offset = 1 + (i * 11);
            const update = {
                id: view.getUint8(offset),
                x: view.getUint32(offset + 1, true),
                y: view.getUint32(offset + 5, true),
                keyboard: view.getUint16(offset+9, true)
            };
            const p = this.players.find(p => p.id == update.id);
            if (!p) {
                newPlayers.push(new Player(update.id, update.x, update.y).setKeyboard(update.keyboard));
            } else {
                p.x = update.x;
                p.y = update.y;
                p.setKeyboard(update.keyboard);
                newPlayers.push(p);
            }
        }
        this.players = newPlayers;
    }
}

class Player {
    constructor(id, x, y) {
        this.id = id;
        this.x = x;
        this.y = y;
        this.keyboard = new Set(); 
        
        return this;
    }

    setKeyboard(kb) {
        this.keyboard = new Set();
        if (kb & 0b1) this.keyboard.add('w');
        if (kb & 0b10) this.keyboard.add('a');
        if (kb & 0b100) this.keyboard.add('s');
        if (kb & 0b1000) this.keyboard.add('d');

        return this;
    }

    render() {
        game.ctx.fillStyle = '#32a852';
        game.ctx.beginPath();
        game.ctx.arc(this.x, this.y, 20, 0, Math.PI * 2);
        game.ctx.fill();
        game.ctx.closePath();
    }
}





const keyboard = new Set();

// Start Client Game
const game = new Game();

socket.onmessage = (event) => {
    if (game) game.handleUpdate(event);
};

window.addEventListener("keydown", (e) => {
    if (e.key === 'Escape') window.toggleGameUI();

    let cmd = new Uint8Array(1);
    if (e.key === 'w' && !keyboard.has(e.key)) {cmd[0] = 0b001; keyboard.add(e.key);}
    if (e.key === 'a' && !keyboard.has(e.key)) {cmd[0] = 0b011; keyboard.add(e.key);}
    if (e.key === 's' && !keyboard.has(e.key)) {cmd[0] = 0b101; keyboard.add(e.key);}
    if (e.key === 'd' && !keyboard.has(e.key)) {cmd[0] = 0b111; keyboard.add(e.key);}
    if (cmd[0] > 0) socket.send(cmd);
});
window.addEventListener("keyup", (e) => {
    let cmd = new Uint8Array(1);
    cmd[0] = 0b1000;
    if (e.key === 'w' && keyboard.has(e.key)) {cmd[0] = 0b000; keyboard.delete(e.key);}
    if (e.key === 'a' && keyboard.has(e.key)) {cmd[0] = 0b010; keyboard.delete(e.key);}
    if (e.key === 's' && keyboard.has(e.key)) {cmd[0] = 0b100; keyboard.delete(e.key);}
    if (e.key === 'd' && keyboard.has(e.key)) {cmd[0] = 0b110; keyboard.delete(e.key);}
    if (cmd[0] != 0b1000) socket.send(cmd);
});
window.addEventListener('mousedown', (e) => {
    console.log(e.clientX, e.clientY)
})
window.toggleGameUI = (forceState) => {
    const isHidden = game.uiLayer.classList.contains('opacity-0');
    const newState = typeof forceState === 'boolean' ? !forceState : isHidden;
    
    if (newState) {
        game.uiLayer.classList.remove('opacity-0');
        game.uiLayer.classList.add('opacity-100');
    } else {
        game.uiLayer.classList.remove('opacity-100');
        game.uiLayer.classList.add('opacity-0');
    }
};