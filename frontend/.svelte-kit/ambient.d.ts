
// this file is generated — do not edit it


/// <reference types="@sveltejs/kit" />

/**
 * Environment variables [loaded by Vite](https://vitejs.dev/guide/env-and-mode.html#env-files) from `.env` files and `process.env`. Like [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), this module cannot be imported into client-side code. This module only includes variables that _do not_ begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) _and do_ start with [`config.kit.env.privatePrefix`](https://svelte.dev/docs/kit/configuration#env) (if configured).
 * 
 * _Unlike_ [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), the values exported from this module are statically injected into your bundle at build time, enabling optimisations like dead code elimination.
 * 
 * ```ts
 * import { API_KEY } from '$env/static/private';
 * ```
 * 
 * Note that all environment variables referenced in your code should be declared (for example in an `.env` file), even if they don't have a value until the app is deployed:
 * 
 * ```
 * MY_FEATURE_FLAG=""
 * ```
 * 
 * You can override `.env` values from the command line like so:
 * 
 * ```sh
 * MY_FEATURE_FLAG="enabled" npm run dev
 * ```
 */
declare module '$env/static/private' {
	export const AGENT: string;
	export const ALLUSERSPROFILE: string;
	export const ANDROID_HOME: string;
	export const ANDROID_SDK_ROOT: string;
	export const APPDATA: string;
	export const ChocolateyInstall: string;
	export const ChocolateyLastPathUpdate: string;
	export const CHROME_CRASHPAD_PIPE_NAME: string;
	export const CLASSPATH: string;
	export const CLAUDE_AGENT_SDK_VERSION: string;
	export const COLOR: string;
	export const CommonProgramFiles: string;
	export const CommonProgramW6432: string;
	export const COMPUTERNAME: string;
	export const ComSpec: string;
	export const DriverData: string;
	export const EDITOR: string;
	export const EFC_10012_1262719628: string;
	export const EFC_10012_1592913036: string;
	export const EFC_10012_2283032206: string;
	export const EFC_10012_2775293581: string;
	export const EFC_10012_2946480783: string;
	export const EFC_10012_3789132940: string;
	export const EFC_10012_4126798990: string;
	export const ELECTRON_RUN_AS_NODE: string;
	export const FPS_BROWSER_APP_PROFILE_STRING: string;
	export const FPS_BROWSER_USER_PROFILE_STRING: string;
	export const GOPATH: string;
	export const HOME: string;
	export const HOMEDRIVE: string;
	export const HOMEPATH: string;
	export const INIT_CWD: string;
	export const JAVA_HOME: string;
	export const KILO: string;
	export const KILOCODE_EDITOR_NAME: string;
	export const KILOCODE_FEATURE: string;
	export const KILOCODE_VERSION: string;
	export const KILO_APP_NAME: string;
	export const KILO_APP_VERSION: string;
	export const KILO_CLIENT: string;
	export const KILO_DISABLE_CHANNEL_DB: string;
	export const KILO_DISABLE_CLAUDE_CODE: string;
	export const KILO_EDITOR_NAME: string;
	export const KILO_ENABLE_QUESTION_TOOL: string;
	export const KILO_MACHINE_ID: string;
	export const KILO_PARENT_PID: string;
	export const KILO_PID: string;
	export const KILO_PLATFORM: string;
	export const KILO_PROCESS_ROLE: string;
	export const KILO_RUN_ID: string;
	export const KILO_SERVER_PASSWORD: string;
	export const KILO_TELEMETRY_LEVEL: string;
	export const KILO_TREE_SITTER_WASM_DIR: string;
	export const KILO_VSCODE_VERSION: string;
	export const LOCALAPPDATA: string;
	export const LOGONSERVER: string;
	export const MIMALLOC_PURGE_DELAY: string;
	export const MOZ_PLUGIN_PATH: string;
	export const NODE: string;
	export const NoDefaultCurrentDirectoryInExePath: string;
	export const NODE_ENV: string;
	export const NODE_USE_SYSTEM_CA: string;
	export const npm_command: string;
	export const npm_config_cache: string;
	export const npm_config_globalconfig: string;
	export const npm_config_global_prefix: string;
	export const npm_config_init_module: string;
	export const npm_config_local_prefix: string;
	export const npm_config_node_gyp: string;
	export const npm_config_noproxy: string;
	export const npm_config_npm_version: string;
	export const npm_config_prefix: string;
	export const npm_config_userconfig: string;
	export const npm_config_user_agent: string;
	export const npm_execpath: string;
	export const npm_lifecycle_event: string;
	export const npm_lifecycle_script: string;
	export const npm_node_execpath: string;
	export const npm_package_json: string;
	export const npm_package_name: string;
	export const npm_package_version: string;
	export const NUMBER_OF_PROCESSORS: string;
	export const NVM_HOME: string;
	export const NVM_SYMLINK: string;
	export const OneDrive: string;
	export const OneDriveCommercial: string;
	export const OPENCODE: string;
	export const OS: string;
	export const Path: string;
	export const PATHEXT: string;
	export const PNPM_HOME: string;
	export const POWERSHELL_DISTRIBUTION_CHANNEL: string;
	export const PROCESSOR_ARCHITECTURE: string;
	export const PROCESSOR_IDENTIFIER: string;
	export const PROCESSOR_LEVEL: string;
	export const PROCESSOR_REVISION: string;
	export const ProgramData: string;
	export const ProgramFiles: string;
	export const ProgramW6432: string;
	export const PROMPT: string;
	export const PSModulePath: string;
	export const PUBLIC: string;
	export const SESSIONNAME: string;
	export const SystemDrive: string;
	export const SystemRoot: string;
	export const TEMP: string;
	export const TMP: string;
	export const USERDOMAIN: string;
	export const USERDOMAIN_ROAMINGPROFILE: string;
	export const USERNAME: string;
	export const USERPROFILE: string;
	export const VSCODE_CODE_CACHE_PATH: string;
	export const VSCODE_CRASH_REPORTER_PROCESS_TYPE: string;
	export const VSCODE_CWD: string;
	export const VSCODE_ESM_ENTRYPOINT: string;
	export const VSCODE_HANDLES_UNCAUGHT_ERRORS: string;
	export const VSCODE_IPC_HOOK: string;
	export const VSCODE_L10N_BUNDLE_LOCATION: string;
	export const VSCODE_NLS_CONFIG: string;
	export const VSCODE_PID: string;
	export const windir: string;
	export const ZES_ENABLE_SYSMAN: string;
	export const __COMPAT_LAYER: string;
}

/**
 * Similar to [`$env/static/private`](https://svelte.dev/docs/kit/$env-static-private), except that it only includes environment variables that begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) (which defaults to `PUBLIC_`), and can therefore safely be exposed to client-side code.
 * 
 * Values are replaced statically at build time.
 * 
 * ```ts
 * import { PUBLIC_BASE_URL } from '$env/static/public';
 * ```
 */
declare module '$env/static/public' {
	
}

/**
 * This module provides access to runtime environment variables, as defined by the platform you're running on. For example if you're using [`adapter-node`](https://github.com/sveltejs/kit/tree/main/packages/adapter-node) (or running [`vite preview`](https://svelte.dev/docs/kit/cli)), this is equivalent to `process.env`. This module only includes variables that _do not_ begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) _and do_ start with [`config.kit.env.privatePrefix`](https://svelte.dev/docs/kit/configuration#env) (if configured).
 * 
 * This module cannot be imported into client-side code.
 * 
 * ```ts
 * import { env } from '$env/dynamic/private';
 * console.log(env.DEPLOYMENT_SPECIFIC_VARIABLE);
 * ```
 * 
 * > [!NOTE] In `dev`, `$env/dynamic` always includes environment variables from `.env`. In `prod`, this behavior will depend on your adapter.
 */
declare module '$env/dynamic/private' {
	export const env: {
		AGENT: string;
		ALLUSERSPROFILE: string;
		ANDROID_HOME: string;
		ANDROID_SDK_ROOT: string;
		APPDATA: string;
		ChocolateyInstall: string;
		ChocolateyLastPathUpdate: string;
		CHROME_CRASHPAD_PIPE_NAME: string;
		CLASSPATH: string;
		CLAUDE_AGENT_SDK_VERSION: string;
		COLOR: string;
		CommonProgramFiles: string;
		CommonProgramW6432: string;
		COMPUTERNAME: string;
		ComSpec: string;
		DriverData: string;
		EDITOR: string;
		EFC_10012_1262719628: string;
		EFC_10012_1592913036: string;
		EFC_10012_2283032206: string;
		EFC_10012_2775293581: string;
		EFC_10012_2946480783: string;
		EFC_10012_3789132940: string;
		EFC_10012_4126798990: string;
		ELECTRON_RUN_AS_NODE: string;
		FPS_BROWSER_APP_PROFILE_STRING: string;
		FPS_BROWSER_USER_PROFILE_STRING: string;
		GOPATH: string;
		HOME: string;
		HOMEDRIVE: string;
		HOMEPATH: string;
		INIT_CWD: string;
		JAVA_HOME: string;
		KILO: string;
		KILOCODE_EDITOR_NAME: string;
		KILOCODE_FEATURE: string;
		KILOCODE_VERSION: string;
		KILO_APP_NAME: string;
		KILO_APP_VERSION: string;
		KILO_CLIENT: string;
		KILO_DISABLE_CHANNEL_DB: string;
		KILO_DISABLE_CLAUDE_CODE: string;
		KILO_EDITOR_NAME: string;
		KILO_ENABLE_QUESTION_TOOL: string;
		KILO_MACHINE_ID: string;
		KILO_PARENT_PID: string;
		KILO_PID: string;
		KILO_PLATFORM: string;
		KILO_PROCESS_ROLE: string;
		KILO_RUN_ID: string;
		KILO_SERVER_PASSWORD: string;
		KILO_TELEMETRY_LEVEL: string;
		KILO_TREE_SITTER_WASM_DIR: string;
		KILO_VSCODE_VERSION: string;
		LOCALAPPDATA: string;
		LOGONSERVER: string;
		MIMALLOC_PURGE_DELAY: string;
		MOZ_PLUGIN_PATH: string;
		NODE: string;
		NoDefaultCurrentDirectoryInExePath: string;
		NODE_ENV: string;
		NODE_USE_SYSTEM_CA: string;
		npm_command: string;
		npm_config_cache: string;
		npm_config_globalconfig: string;
		npm_config_global_prefix: string;
		npm_config_init_module: string;
		npm_config_local_prefix: string;
		npm_config_node_gyp: string;
		npm_config_noproxy: string;
		npm_config_npm_version: string;
		npm_config_prefix: string;
		npm_config_userconfig: string;
		npm_config_user_agent: string;
		npm_execpath: string;
		npm_lifecycle_event: string;
		npm_lifecycle_script: string;
		npm_node_execpath: string;
		npm_package_json: string;
		npm_package_name: string;
		npm_package_version: string;
		NUMBER_OF_PROCESSORS: string;
		NVM_HOME: string;
		NVM_SYMLINK: string;
		OneDrive: string;
		OneDriveCommercial: string;
		OPENCODE: string;
		OS: string;
		Path: string;
		PATHEXT: string;
		PNPM_HOME: string;
		POWERSHELL_DISTRIBUTION_CHANNEL: string;
		PROCESSOR_ARCHITECTURE: string;
		PROCESSOR_IDENTIFIER: string;
		PROCESSOR_LEVEL: string;
		PROCESSOR_REVISION: string;
		ProgramData: string;
		ProgramFiles: string;
		ProgramW6432: string;
		PROMPT: string;
		PSModulePath: string;
		PUBLIC: string;
		SESSIONNAME: string;
		SystemDrive: string;
		SystemRoot: string;
		TEMP: string;
		TMP: string;
		USERDOMAIN: string;
		USERDOMAIN_ROAMINGPROFILE: string;
		USERNAME: string;
		USERPROFILE: string;
		VSCODE_CODE_CACHE_PATH: string;
		VSCODE_CRASH_REPORTER_PROCESS_TYPE: string;
		VSCODE_CWD: string;
		VSCODE_ESM_ENTRYPOINT: string;
		VSCODE_HANDLES_UNCAUGHT_ERRORS: string;
		VSCODE_IPC_HOOK: string;
		VSCODE_L10N_BUNDLE_LOCATION: string;
		VSCODE_NLS_CONFIG: string;
		VSCODE_PID: string;
		windir: string;
		ZES_ENABLE_SYSMAN: string;
		__COMPAT_LAYER: string;
		[key: `PUBLIC_${string}`]: undefined;
		[key: `${string}`]: string | undefined;
	}
}

/**
 * Similar to [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), but only includes variables that begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) (which defaults to `PUBLIC_`), and can therefore safely be exposed to client-side code.
 * 
 * Note that public dynamic environment variables must all be sent from the server to the client, causing larger network requests — when possible, use `$env/static/public` instead.
 * 
 * ```ts
 * import { env } from '$env/dynamic/public';
 * console.log(env.PUBLIC_DEPLOYMENT_SPECIFIC_VARIABLE);
 * ```
 */
declare module '$env/dynamic/public' {
	export const env: {
		[key: `PUBLIC_${string}`]: string | undefined;
	}
}
