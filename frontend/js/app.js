import { createPalette, createPaletteFromBase } from './api.js'

let currentColors = [];

document.getElementById('generate-btn').addEventListener('click', async() => {
    try {
        const locked = currentColors.filter(c => c.locked);

        console.log('from-site request: ', currentColors)
        const palette = locked.length > 0
            ? await createPaletteFromBase(currentColors)
            : await createPalette();
        console.log('from-base response: ', palette);

        if (!Array.isArray(palette.color_list)) {
            throw new Error('Response missing color_list: ' + JSON.stringify(palette))
        }
        
        currentColors = palette.color_list;
        renderPalette(currentColors);
    } catch (err) {
        console.error(err);
    }
});

function renderPalette(colors) {
    const container = document.getElementById('palette');
    container.innerHTML = ''; // clear swatches
    colors.forEach((c, i) => {
        const item = document.createElement('div');
        item.className = 'palette-item';

        const swatch = document.createElement('div');
        swatch.className = 'swatch';
        swatch.style.backgroundColor = `rgb(${c.rgb.r}, ${c.rgb.g}, ${c.rgb.b})`;
        swatch.title = c.name.name; 

        const lock = document.createElement('button');
        lock.className = 'lock';
        lock.innerText = c.locked ? 'unlock' : 'lock';

        lock.addEventListener('click', () => {
            currentColors[i].locked = !currentColors[i].locked;
            lock.textContent = currentColors[i].locked ? 'unlock' : 'lock';
            item.classList.toggle('locked', currentColors[i].locked);
        });

        item.classList.toggle('locked', c.locked);
        item.append(swatch, lock)
        container.appendChild(item);
    });
}
