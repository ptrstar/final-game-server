const env = 'dev';

const gameType = "wasd";
const host = env == 'dev' ? "localhost:8080" : "play.trojanos.ch";
const httpHost = env == 'dev' ? `http://${host}` : `https://${host}`;
const wsHost = env == 'dev' ? `ws://${host}` : `wss://${host}`;

(async function() {
    
    const resources = [
        `/${gameType}/game.js`,
    ];

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
        await loadScript(`/${gameType}/snippets/basehtml.js`);
        
        if (window.__GAMEHTML__) {
            viewport.innerHTML = window.__GAMEHTML__;
            delete window.__GAMEHTML__;
        } else {
            throw new Error("Snippet data missing after load.");
        }

        // Load all gameresources
        await Promise.all(resources.map(loadScript));

        console.log(`[${gameType}] Loaded all resources. All systems go.`);

    } catch (err) {
        viewport.innerHTML = `<div style="color:red;">Load Error: ${err.message}</div>`;
        console.error(err);
    }
})();
