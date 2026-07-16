

export const index = 1;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/error.svelte.js')).default;
export const imports = ["_app/immutable/nodes/1.BECadB8n.js","_app/immutable/chunks/IoB7A9BM.js","_app/immutable/chunks/B3_NYjGC.js"];
export const stylesheets = [];
export const fonts = [];
