import { createPalette } from './api.js'
document.getElementById('generate-btn').addEventListener('click', async() => {
    try {
        const palette = await createPalette();
        console.log(palette.color_list);
        renderPalette(palette.color_list);
    } catch (err) {
        console.error(err);
    }
});

function renderPalette(colors) {
    const container = document.getElementById('palette');
    container.innerHTML = ''; // clear swatches
    colors.forEach(c => {
        const swatch = document.createElement('div');
        swatch.className = 'swatch';
        swatch.style.backgroundColor = `rgb(${c.rgb.r}, ${c.rgb.g}, ${c.rgb.b})`;
        swatch.title = c.name.name; 
        container.appendChild(swatch);
    });
}
