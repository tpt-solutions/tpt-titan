

export const index = 1;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/error.svelte.js')).default;
export const imports = ["_app/immutable/nodes/1.yOV0YdJq.js","_app/immutable/chunks/CJe9_yH9.js","_app/immutable/chunks/B9R0deLA.js"];
export const stylesheets = [];
export const fonts = [];
