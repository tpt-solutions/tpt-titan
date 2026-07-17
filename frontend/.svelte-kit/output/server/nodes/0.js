

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export const imports = ["_app/immutable/nodes/0.CICj2Vg2.js","_app/immutable/chunks/D_yhzppV.js","_app/immutable/chunks/DX9WGBLo.js"];
export const stylesheets = ["_app/immutable/assets/0.DLK1wheZ.css"];
export const fonts = [];
