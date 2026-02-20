

export const index = 5;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/calendar/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/5.CJqC49hY.js","_app/immutable/chunks/DED9TEO-.js","_app/immutable/chunks/BAYF4ZET.js","_app/immutable/chunks/DG8_rw9z.js","_app/immutable/chunks/Cq47Lcui.js","_app/immutable/chunks/Bx0C6IOw.js"];
export const stylesheets = ["_app/immutable/assets/spreadsheet.DHjeLx6L.css","_app/immutable/assets/forms.C6TuIXSB.css"];
export const fonts = [];
