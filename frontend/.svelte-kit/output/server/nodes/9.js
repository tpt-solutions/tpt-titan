

export const index = 9;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/database/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/9.CAX_7g90.js","_app/immutable/chunks/D_yhzppV.js","_app/immutable/chunks/Di8qqXgd.js","_app/immutable/chunks/B7aVzNxh.js","_app/immutable/chunks/CV8_1icv.js"];
export const stylesheets = ["_app/immutable/assets/spreadsheet.DHjeLx6L.css","_app/immutable/assets/forms.C6TuIXSB.css"];
export const fonts = [];
