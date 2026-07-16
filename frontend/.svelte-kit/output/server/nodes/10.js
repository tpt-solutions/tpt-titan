

export const index = 10;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/editor/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/10.DFMcCQqy.js","_app/immutable/chunks/Gqbyd49N.js","_app/immutable/chunks/B3_NYjGC.js","_app/immutable/chunks/CmsKOCeN.js","_app/immutable/chunks/2xjUDx4y.js","_app/immutable/chunks/DPmjPLTl.js","_app/immutable/chunks/DYkWwLq0.js"];
export const stylesheets = ["_app/immutable/assets/forms.C6TuIXSB.css","_app/immutable/assets/spreadsheet.DHjeLx6L.css","_app/immutable/assets/10.DzMnSV1g.css"];
export const fonts = [];
