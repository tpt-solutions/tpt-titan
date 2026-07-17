

export const index = 4;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/auth/login/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/4.BCm0viZB.js","_app/immutable/chunks/D_yhzppV.js","_app/immutable/chunks/DX9WGBLo.js"];
export const stylesheets = [];
export const fonts = [];
