

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export const imports = ["_app/immutable/nodes/0.rkayD7ut.js","_app/immutable/chunks/D_yhzppV.js","_app/immutable/chunks/D-oSuO4X.js"];
export const stylesheets = ["_app/immutable/assets/0.fFXixcFY.css"];
export const fonts = [];
