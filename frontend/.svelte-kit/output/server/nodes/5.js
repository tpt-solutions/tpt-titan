

export const index = 5;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/auth/register/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/5.FmCKMHd6.js","_app/immutable/chunks/D_yhzppV.js","_app/immutable/chunks/1brQZ6N_.js"];
export const stylesheets = [];
export const fonts = [];
