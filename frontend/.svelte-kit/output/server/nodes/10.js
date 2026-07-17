

export const index = 10;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/editor/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/10.DWcPQGX4.js","_app/immutable/chunks/BC-VK_VN.js","_app/immutable/chunks/D_yhzppV.js","_app/immutable/chunks/CmsKOCeN.js","_app/immutable/chunks/CRAR0TG-.js","_app/immutable/chunks/B7aVzNxh.js","_app/immutable/chunks/Oeo97mMt.js"];
export const stylesheets = ["_app/immutable/assets/spreadsheet.DHjeLx6L.css","_app/immutable/assets/forms.C6TuIXSB.css","_app/immutable/assets/10.DzMnSV1g.css"];
export const fonts = [];
