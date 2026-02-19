

export const index = 10;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/forms/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/10.Rx0sOoG3.js","_app/immutable/chunks/B9R0deLA.js","_app/immutable/chunks/lV6On81A.js","_app/immutable/chunks/DnTaMpYL.js","_app/immutable/chunks/BHtuMjBG.js"];
export const stylesheets = ["_app/immutable/assets/spreadsheet.m32kOtRp.css","_app/immutable/assets/forms.C6TuIXSB.css"];
export const fonts = [];
