# Tauri v2 Migration — TODO

Migration plan: Tauri 1.5 → 2.x for `desktop/`. Scoped via [is-there-a-tauri-glowing-stallman.md](C:\Users\Phillip\.claude\plans\is-there-a-tauri-glowing-stallman.md) — no frontend Tauri API usage exists to port, and the Rust side's dialog/notification/shell integrations are unimplemented stubs, so this is mostly config/dependency/macro updates.

- [x] Bump `desktop/package.json`: `@tauri-apps/cli` → `^2`
- [x] Bump `desktop/src-tauri/Cargo.toml`: `tauri-build` → `"2"`, `tauri` → `{ version = "2", features = [] }` (drop unused `shell-open` feature)
- [x] Rewrite `desktop/src-tauri/tauri.conf.json` to v2 schema:
  - [x] Flatten `package.productName`/`package.version` → top-level `productName`/`version`
  - [x] `build.devPath` → `build.devUrl`, `build.distDir` → `build.frontendDist`
  - [x] Move `tauri.bundle` → top-level `bundle`, `tauri.windows` → `app.windows`, `tauri.security` → `app.security`
  - [x] Update `$schema` path
- [x] Add `desktop/src-tauri/capabilities/default.json` with the default window capability (`core:default`) — v2 requires an explicit capabilities file even with no permissions used
- [x] Verify `desktop/src-tauri/src/main.rs` compiles under v2 — check `use tauri::{AppHandle, State}` / unqualified `State` usage still resolves (v1's prelude may have covered this implicitly)
- [x] Leave dialog/notification/shell stubs as-is; note that completing them later requires `tauri-plugin-dialog` / `tauri-plugin-notification` / `tauri-plugin-shell` (Rust crate + `@tauri-apps/plugin-*` JS package)

## Verification
- [x] `cd desktop && npm install` — confirm `@tauri-apps/cli@2.x` installs cleanly
- [x] `cd desktop/src-tauri && cargo check` — confirm Rust compiles against `tauri = "2"`
- [ ] `cd desktop && npm run tauri-dev` — launch dev app, exercise all 9 commands (`greet`, `start_backend`, `stop_backend`, `get_backend_status`, `get_system_info`, `read_file`, `write_file`, `open_file_dialog`, `save_file_dialog`)
- [ ] `cd desktop && npm run tauri-build` — confirm release build succeeds and produces expected installer(s)

## Follow-up (not in scope of this migration, tracked separately)
- [ ] Bundle size/speed tuning: add `[profile.release]` (`strip`, `lto`, `opt-level`) to `desktop/src-tauri/Cargo.toml`, narrow `bundle.targets` from `"all"` to the single format actually shipped
