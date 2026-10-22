// TODO 
// add add and remove colors
// info button with popup
// spacebar generates new colors
import { createPalette, createPaletteFromBase, createGradient } from './api.js'
import { renderPalette, renderGradient } from './render.js'

let currentColors = [];
let currentGradient = [];

// creates a new palette on page loading
document.addEventListener('DOMContentLoaded', async () => {
    try {
        const palette = await createPalette();

        if (!Array.isArray(palette.color_list)) {
            throw new Error('Response missing color_list: ' + JSON.stringify(palette))
        }

        currentColors = palette.color_list;
        renderPalette(currentColors);
    } catch (err) {
        console.error(err);
    }
});

// creates a new palette if no colors are locked, otherwise generates a palette
// of the same length based on the locked colors
document.getElementById('generate-btn').addEventListener('click', async () => {
    try {
        const locked = currentColors.filter(c => c.locked);

        if (locked.length === currentColors.length) {
            return
        }

        const palette = locked.length > 0
            ? await createPaletteFromBase(currentColors)
            : await createPalette();

        if (!Array.isArray(palette.color_list)) {
            throw new Error('Response missing color_list: ' + JSON.stringify(palette))
        }

        currentColors = palette.color_list;
        renderPalette(currentColors);
    } catch (err) {
        console.error(err);
    }
});

// creates a gradient based on two colors and a length
document.getElementById('generate-gradient').addEventListener('click', async () => {
    try {
        const color1 = document.getElementById('first-color');
        const color2 = document.getElementById('second-color');
        const numColors = document.getElementById('num-colors');

        const palette = await createGradient(color1.value, color2.value, Number(numColors.value));

        if (!Array.isArray(palette.color_list)) {
            throw new Error('Response missing color_list: ' + JSON.stringify(palette))
        }

        currentGradient = palette.color_list;
        renderGradient(currentGradient);
    } catch (err) {
        console.error(err);
    }
})
