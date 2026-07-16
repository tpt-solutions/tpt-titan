

export const index = 6;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/calendar/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/6.BnMc8fd1.js","_app/immutable/chunks/B3_NYjGC.js","_app/immutable/chunks/IoB7A9BM.js","_app/immutable/chunks/DPmjPLTl.js","_app/immutable/chunks/DYkWwLq0.js","_app/immutable/chunks/2xjUDx4y.js"];
export const stylesheets = ["_app/immutable/assets/spreadsheet.DHjeLx6L.css","_app/immutable/assets/forms.C6TuIXSB.css"];
export const fonts = [];
