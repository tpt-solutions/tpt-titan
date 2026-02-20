

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.DL0GskjO.js","_app/immutable/chunks/DED9TEO-.js","_app/immutable/chunks/BAYF4ZET.js"];
export const stylesheets = [];
export const fonts = [];
