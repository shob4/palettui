const BASE_URL = 'http://localhost:8080'

export async function createPalette() {
    const res = await fetch(`${BASE_URL}/palette/create`);
    if (!res.ok) throw new Error(`Server error: ${res.status}`);
    return res.json();
}

export async function createPaletteFromBase(lockedColors) {
    const res = await fetch(`${BASE_URL}/palette/create-from-base`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ color_list: lockedColors })
    });
    if (!res.ok) throw new Error(`Server error: ${res.status}`);
    return res.json();
}

export async function createGradient(color1, color2, numColors,) {
    const res = await fetch(`${BASE_URL}/palette/gradient`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ color_list: [color1, color2], num: numColors})
    });
    if (!res.ok) throw new Error(`Server error: ${res.status}`);
    return res.json();
}
