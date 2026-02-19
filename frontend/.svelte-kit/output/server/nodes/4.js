

export const index = 4;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/auth/register/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/4.DYRSJqje.js","_app/immutable/chunks/B9R0deLA.js","_app/immutable/chunks/CJe9_yH9.js"];
export const stylesheets = [];
export const fonts = [];
