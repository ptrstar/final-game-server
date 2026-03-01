const env = 'dev';

const params = new URLSearchParams(window.location.search)
var gameType = params.get('gameType') || 'info';
var roomId = params.get('roomId') || '';
const host = env == 'dev' ? "localhost:8080" : "play.trojanos.ch";
const httpHost = env == 'dev' ? `http://${host}` : `https://${host}`;
const wsHost = env == 'dev' ? `ws://${host}` : `wss://${host}`;

var activeModule;

async function game_injector(t, id) {

    if (activeModule) {
        activeModule.TERMINATE()
    }

    gameType = t
    roomId = id

    const viewport = document.getElementById('game-viewport');
    if (!viewport) return console.error("Missing #game-viewport");

    // Loading State
    viewport.innerHTML = `<div style="color:gray; font-style:italic;">Initializing ${gameType} Injector...</div>`;

    // CORS-proof loader
    const loadScript = (url) => new Promise((resolve, reject) => {
        const s = document.createElement('script');
        s.src = url.startsWith('http') ? url : `${httpHost}${url}`;
        s.async = true;
        s.onload = resolve;
        s.onerror = () => reject(new Error(`Failed to load: ${url}`));
        document.head.appendChild(s);
    });

    try {
        // First load and inject html
        await loadScript(`/${gameType}/snippets/base.js`);
        
        if (window.__GAMEHTML__) {
            viewport.innerHTML = window.__GAMEHTML__;
            delete window.__GAMEHTML__;
        } else {
            throw new Error("Snippet HTML data missing after load.");
        }

        // Load all gameresources
        const gameModule = await import(`${httpHost}/${gameType}/module.js`);

        console.log(`game_injector() -> Module loaded. Starting engine...`);

        // 3. Call the exported StartGame function
        if (gameModule.STARTGAME) {
            gameModule.STARTGAME();
            activeModule = gameModule;
        } else {
            console.error("Module does not export StartGame()");
        }

    } catch (err) {
        viewport.innerHTML = `<div style="color:red;">Load Error: ${err.message}</div>`;
        console.error(err);
    }
};

game_injector(gameType, roomId)
