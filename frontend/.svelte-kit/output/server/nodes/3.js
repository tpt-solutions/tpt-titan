

export const index = 3;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/auth/login/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/3.L1GY-YhN.js","_app/immutable/chunks/DED9TEO-.js","_app/immutable/chunks/BAYF4ZET.js"];
export const stylesheets = [];
export const fonts = [];
