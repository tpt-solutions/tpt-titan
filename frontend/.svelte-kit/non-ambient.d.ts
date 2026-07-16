
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	export interface AppTypes {
		RouteId(): "/" | "/admin" | "/auth" | "/auth/login" | "/auth/register" | "/calendar" | "/chat" | "/contacts" | "/database" | "/editor" | "/email" | "/export" | "/files" | "/forms" | "/math" | "/monitoring" | "/plugins" | "/settings" | "/setup" | "/speech" | "/spreadsheet" | "/tasks" | "/voice" | "/workflows";
		RouteParams(): {
			
		};
		LayoutParams(): {
			"/": Record<string, never>;
			"/admin": Record<string, never>;
			"/auth": Record<string, never>;
			"/auth/login": Record<string, never>;
			"/auth/register": Record<string, never>;
			"/calendar": Record<string, never>;
			"/chat": Record<string, never>;
			"/contacts": Record<string, never>;
			"/database": Record<string, never>;
			"/editor": Record<string, never>;
			"/email": Record<string, never>;
			"/export": Record<string, never>;
			"/files": Record<string, never>;
			"/forms": Record<string, never>;
			"/math": Record<string, never>;
			"/monitoring": Record<string, never>;
			"/plugins": Record<string, never>;
			"/settings": Record<string, never>;
			"/setup": Record<string, never>;
			"/speech": Record<string, never>;
			"/spreadsheet": Record<string, never>;
			"/tasks": Record<string, never>;
			"/voice": Record<string, never>;
			"/workflows": Record<string, never>
		};
		Pathname(): "/" | "/admin" | "/admin/" | "/auth" | "/auth/" | "/auth/login" | "/auth/login/" | "/auth/register" | "/auth/register/" | "/calendar" | "/calendar/" | "/chat" | "/chat/" | "/contacts" | "/contacts/" | "/database" | "/database/" | "/editor" | "/editor/" | "/email" | "/email/" | "/export" | "/export/" | "/files" | "/files/" | "/forms" | "/forms/" | "/math" | "/math/" | "/monitoring" | "/monitoring/" | "/plugins" | "/plugins/" | "/settings" | "/settings/" | "/setup" | "/setup/" | "/speech" | "/speech/" | "/spreadsheet" | "/spreadsheet/" | "/tasks" | "/tasks/" | "/voice" | "/voice/" | "/workflows" | "/workflows/";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): "/favicon.ico" | "/site.webmanifest" | string & {};
	}
}