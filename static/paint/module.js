import PaintGame from '/paint/game.js'

export function STARTGAME() {
    window.game = new PaintGame();
    // Assuming wsHost, gameType, and roomId are globally available from injector
    window.game.socket = new WebSocket(`${wsHost}/ws?type=${gameType}&room=${roomId}`);
    window.game.socket.onmessage = window.game.handleUpdate.bind(window.game);

    console.log("Starting Game")
    window.addEventListener('keydown', window.game.keydown.bind(window.game))
    window.addEventListener("keyup", window.game.keyup.bind(window.game));
    window.addEventListener('resize', window.game.upscaleCanvas.bind(window.game));
    game.upscaleCanvas();
    window.requestAnimationFrame(window.game.main.bind(window.game))

}

export function TERMINATE() {
    if (window.game) window.game.terminate();
}
