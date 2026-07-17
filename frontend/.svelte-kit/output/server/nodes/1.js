

export const index = 1;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/error.svelte.js')).default;
export const imports = ["_app/immutable/nodes/1.BZbHQ_9D.js","_app/immutable/chunks/DX9WGBLo.js","_app/immutable/chunks/D_yhzppV.js"];
export const stylesheets = [];
export const fonts = [];
