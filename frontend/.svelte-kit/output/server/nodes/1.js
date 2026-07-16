

export const index = 1;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/error.svelte.js')).default;
export const imports = ["_app/immutable/nodes/1.Dj01_uHZ.js","_app/immutable/chunks/D-oSuO4X.js","_app/immutable/chunks/D_yhzppV.js"];
export const stylesheets = [];
export const fonts = [];
