

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.DvTcSTwJ.js","_app/immutable/chunks/B3_NYjGC.js","_app/immutable/chunks/IoB7A9BM.js"];
export const stylesheets = [];
export const fonts = [];
