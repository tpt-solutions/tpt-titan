

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export const imports = ["_app/immutable/nodes/0.BzRfgvQB.js","_app/immutable/chunks/B3_NYjGC.js","_app/immutable/chunks/IoB7A9BM.js"];
export const stylesheets = ["_app/immutable/assets/0.uRYZ9iRx.css"];
export const fonts = [];
