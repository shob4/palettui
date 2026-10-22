// TODO 
// add add and remove colors
// info button with popup
// spacebar generates new colors
// move buttons into swatches
// buttons only visible on mouseover
// show hex and name with color picker button
import { createMonochrome, createColor, createPalette, createPaletteFromBase, createGradient } from './api.js'

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

// -------- render --------

const MAX_CACHED = 50;
const shadeCache = new Map();

function cacheShades(key, shades) {
    shadeCache.delete(key);
    shadeCache.set(key, shades);
    if (shadeCache.size > MAX_CACHED) {
        const oldest = shadeCache.keys().next().value;
        shadeCache.delete(oldest);
    }
}

export function renderGradient(colors) {
    const container = document.getElementById('gradient');
    container.innerHTML = ''; // clear swatches
    colors.forEach((c, i) => {
        if (!c.locked) {
            shadeCache[i] = [];
        }
        const item = document.createElement('div');
        item.className = 'palette-item';

        const swatch = renderSwatch(c);

        item.append(swatch)
        container.appendChild(item);
    });
}

export function renderPalette(colors) {
    const container = document.getElementById('palette');
    container.innerHTML = ''; // clear swatches
    colors.forEach((c, i) => {
        const item = document.createElement('div');
        item.className = 'palette-item';

        const swatch = renderSwatch(c);

        const lock = document.createElement('button');
        lock.className = 'lock';
        lock.innerText = c.locked ? 'unlock' : 'lock';

        lock.addEventListener('click', () => {
            currentColors[i].locked = !currentColors[i].locked;
            lock.textContent = currentColors[i].locked ? 'unlock' : 'lock';
            item.classList.toggle('locked', currentColors[i].locked);
        });

        const shadeList = document.createElement('div');
        shadeList.className = 'shade-list';
        shadeList.hidden = true;

        const shadeButton = document.createElement('button');
        shadeButton.className = 'shade-button';
        shadeButton.textContent = 'shades';

        shadeButton.addEventListener('click', async () => {
            if (!shadeList.hidden) {
                shadeList.hidden = true;
                return;
            }
            try {
                const key = toHex(colors[i]);
                let shades = shadeCache.get(key);
                if (shades) {
                    cacheShades(key, shades);
                } else {
                    shades = await createMonochrome(colors[i]);
                    cacheShades(key, shades);
                }
                shadeList.replaceChildren(...shades.map(shade => {
                    const s = document.createElement('div');
                    s.className = 'shade';
                    s.style.backgroundColor = toRgb(shade);
                    s.title = toHex(shade);
                    s.addEventListener('click', () => {
                        colors[i] = { ...shade, locked: colors[i].locked };
                        swatch.style.backgroundColor = toRgb(shade);
                        swatch.title = shade.name.name;
                        shadeList.hidden = true;
                    })
                    return s;
                }));
                shadeList.hidden = false;
            } catch (err) {
                console.error(err);
            }
        });

        const picker = document.createElement('input');
        picker.type = 'color'
        picker.className = 'picker-input';
        picker.value = toHex(c);

        const pickerButton = document.createElement('button');
        pickerButton.textContent = picker.value;

        pickerButton.addEventListener('click', () => {
            picker.value = toHex(colors[i]);
            if (picker.showPicker) picker.showPicker();
            else picker.click();
        })

        picker.addEventListener('change', async () => {
            try {
                const picked = await createColor(picker.value);
                colors[i] = { ...picked, locked: colors[i].locked };
                swatch.style.backgroundColor = toRgb(picked);
                swatch.title = picked.name.name;
            } catch (err) {
                console.error(err);
            }
        })

        item.classList.toggle('locked', c.locked);
        item.append(swatch, shadeList, lock, shadeButton, pickerButton, picker);
        container.appendChild(item);
    });
}

export function renderSwatch(color) {
    const swatch = document.createElement('div');
    swatch.className = 'swatch';
    swatch.style.backgroundColor = `rgb(${color.rgb.r}, ${color.rgb.g}, ${color.rgb.b})`;
    swatch.title = color.name.name;
    return swatch;
}

const toHex = c => '#' + c.hex.h.toString(16).padStart(6, '0');
const toRgb = c => `rgb(${c.rgb.r}, ${c.rgb.g}, ${c.rgb.b})`;
