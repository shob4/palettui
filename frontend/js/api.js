export async function createPalette() {
    const res = await fetch('http://localhost:8080/palette/create');
    if (!res.ok) throw new Error(`Server error: ${res.status}`);
    return res.json();
}
