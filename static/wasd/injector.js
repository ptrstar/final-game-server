const game_type = "wasd";
const src = "localhost:8080";

const game_viewport = document.getElementById('game-viewport');
if (!game_viewport) console.error("element with id='game-viewport' not found");
game_viewport.innerHTML = `<div class="absolute inset-0 flex items-center justify-center text-gray-500 italic">Loading game injector: done. Loading game...</div>`

fetch(`/snippets/${game_type}.html`)
    .then(res => res.text())
    .then(html => {
        game_viewport.innerHTML = html;
        const script = document.createElement("script");
        script.src = `http://${src}/wasd/game.js`;
        document.body.appendChild(script); 
    });