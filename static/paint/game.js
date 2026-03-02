import Player from '/paint/player.js'

export default class PaintGame {
    constructor() {
        this.canvas = document.getElementById("game-canvas");
        this.ctx = this.canvas.getContext('2d');

        this.keyboard = new Set();
        
        this.players = []
    }

    main() {
        this.gameUpdate();

        this.render();

        window.requestAnimationFrame(this.main.bind(this));
    }

    gameUpdate() {
        
    }

    render() {
        this.ctx.fillStyle = '#43cae8';
        this.ctx.fillRect(0,0,this.canvas.clientWidth, this.canvas.clientHeight);

        this.players.forEach(p => {
            p.render(); 
        })
    }
    
    upscaleCanvas() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
        const dpi = window.devicePixelRatio || 1;
        const rect = this.canvas.getBoundingClientRect();
        this.canvas.width = rect.width * dpi;
        this.canvas.height = rect.height * dpi;
        this.ctx.scale(dpi, dpi);
    }

    async handleUpdate(e) {
        try {
            const text = await e.data.text();
            // TODO: do this in a sensible way
            const players = JSON.parse(text)
            this.players = []
            players.forEach(p => {
                this.players.push(new Player(p))
            });
            
        } catch (e) {
            console.error("Parse error", e);
        }
    }

    terminate() {
        console.log("Terminating Game")
        window.removeEventListener('resize', this.upscaleCanvas)
        window.removeEventListener('keydown', this.keydown);
        window.removeEventListener('keyup', this.keyup);
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
        if (e.key === 'w' && !this.keyboard.has(e.key)) {cmd[0] = 0b001; this.keyboard.add(e.key);}
        if (e.key === 'a' && !this.keyboard.has(e.key)) {cmd[0] = 0b011; this.keyboard.add(e.key);}
        if (e.key === 's' && !this.keyboard.has(e.key)) {cmd[0] = 0b101; this.keyboard.add(e.key);}
        if (e.key === 'd' && !this.keyboard.has(e.key)) {cmd[0] = 0b111; this.keyboard.add(e.key);}
        if (cmd[0] > 0) this.socket.send(cmd);
    }
    keyup(e) {
        let cmd = new Uint8Array(1);
        cmd[0] = 0b1000;
        if (e.key === 'w' && this.keyboard.has(e.key)) {cmd[0] = 0b000; this.keyboard.delete(e.key);}
        if (e.key === 'a' && this.keyboard.has(e.key)) {cmd[0] = 0b010; this.keyboard.delete(e.key);}
        if (e.key === 's' && this.keyboard.has(e.key)) {cmd[0] = 0b100; this.keyboard.delete(e.key);}
        if (e.key === 'd' && this.keyboard.has(e.key)) {cmd[0] = 0b110; this.keyboard.delete(e.key);}
        if (cmd[0] != 0b1000) this.socket.send(cmd);
    }
}