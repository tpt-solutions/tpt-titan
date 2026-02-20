

export const index = 1;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/error.svelte.js')).default;
export const imports = ["_app/immutable/nodes/1.CkB7y4_l.js","_app/immutable/chunks/BAYF4ZET.js","_app/immutable/chunks/DED9TEO-.js"];
export const stylesheets = [];
export const fonts = [];
