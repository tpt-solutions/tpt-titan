

export const index = 3;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/auth/login/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/3.CK96d3jA.js","_app/immutable/chunks/B9R0deLA.js","_app/immutable/chunks/CJe9_yH9.js"];
export const stylesheets = [];
export const fonts = [];
