
// import test from '/info/testimport.js'

// test()

var game;
var socket;

class Game {

    

}

export function STARTGAME() {
    game = new Game();
    socket = new WebSocket(`${wsHost}/ws?type=${gameType}&room=${roomId}`);
    socket.onmessage = (e) => game.handleUpdate(e);
}

export function TERMINATE() {
    game.terminate();
}

