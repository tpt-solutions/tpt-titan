

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.D6QjlN0d.js","_app/immutable/chunks/D_yhzppV.js","_app/immutable/chunks/D-oSuO4X.js"];
export const stylesheets = [];
export const fonts = [];
