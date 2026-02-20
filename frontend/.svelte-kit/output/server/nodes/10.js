

export const index = 10;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/forms/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/10.Cs30pxWL.js","_app/immutable/chunks/DED9TEO-.js","_app/immutable/chunks/DG8_rw9z.js","_app/immutable/chunks/Cq47Lcui.js","_app/immutable/chunks/Bx0C6IOw.js"];
export const stylesheets = ["_app/immutable/assets/spreadsheet.DHjeLx6L.css","_app/immutable/assets/forms.C6TuIXSB.css"];
export const fonts = [];
