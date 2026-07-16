# TPT Titan — TODO

Tracking list generated from a full-project audit (stubs/TODOs, frontend coverage, security, adoption tooling) on 2026-07-16. Items are grouped by area and roughly ordered by severity within each group. This file will be updated as remaining audit passes (code-level TODO/stub inventory, security review) complete.

**Update 2026-07-16:** Commit `d5fa6ee` wired up the previously dead route groups (admin, document export, spreadsheet chart/collab/excel/formula, all six advanced form modules) and added a P2P collaboration service. Those routes are now reachable — see the code-level stub items below for which of them still return mock/placeholder data rather than real logic. Tasks and File Sync remain fully unwired.

**Update 2026-07-16 (later, uncommitted):** Rewrote the spreadsheet formula arithmetic evaluator (`backend/services/spreadsheet_math.go`) — the old implementation only handled a single `+`/`-` split between exactly two terms (no operator precedence, no multi-term expressions, no real range expansion). Replaced with a recursive-descent parser supporting `+ - * /` with correct precedence, parentheses, nested function calls, and proper `A1:B3`-style range expansion into cell lists; cell references are now case-insensitive and function arguments are trimmed. Also added argument-count/nil guards to `math_power.go`, `math_rounding.go`, `math_trigonometric.go` (previously panicked on missing args) and nil-`DB` guards to `ai.go`/`model_service.go` `GetAvailableModels`/`CheckForUpgrades` (previously panicked if called before DB init). Tests updated to match corrected behavior. Not yet committed — see `git status`.

## Critical — Frontend UIs with no backend at all (fully non-functional)

- [x] Tasks: `frontend/src/routes/tasks/+page.svelte` + `TaskBoard.svelte`/`TaskForm.svelte` (459 lines) — was fully non-functional. Added `backend/models/task.go` (Task/Project/Subtask GORM models), `backend/services/task_service.go` (CRUD over `*sql.DB`), and `backend/routes/tasks.go` with `/tasks` (+ projects) routes wired in `server.go`. Frontend now fetches/creates/updates/deletes tasks via the API. AI prediction/suggestion buttons in TaskForm still call non-existent `/tasks/predict-*` endpoints and fail silently (optional, out of scope).
- [x] File Sync: `frontend/src/routes/files/+page.svelte` calls `/filesync/folders`, `/filesync/status`, `/filesync/sync/:id`. Added `backend/routes/filesync.go` wiring to the existing `services/filesync.go` (`GetSyncFolders`, `CreateSyncFolder`, `GetSyncStatus`, `SyncFolder`), registered under `/filesync` in `server.go`. UI is now functional.
- [x] Plugin system: `backend/services/plugin_system.go` now reachable — added `backend/routes/plugins.go` with `/plugins` (list/stats/enable/disable/unload/settings) routes + `InitPluginService`, plus a frontend page `frontend/src/routes/plugins/+page.svelte` and nav entry. Also fixed the `GetPluginSettings` regression (was returning `nil, nil`; now reads from DB).

## High — Backend features with no frontend UI at all

- [ ] Admin console — backend routes are now wired (`server.go:547-560`), but no `/admin` frontend route exists yet.
- [ ] Document export — ~10 backend endpoints, zero frontend calls.
- [ ] Speech (TTS/STT) — 8 registered endpoints, no dedicated UI.
- [ ] Voice notes/annotations — 9 registered endpoints (`voice.go`), no frontend route or component.
- [ ] Math (expressions, canvas, export, recognition, templates) — ~19 registered endpoints, zero frontend page.
- [ ] Workflows — `WorkflowBuilder.svelte` / `WorkflowDesignerModal.svelte` exist but are not mounted under any route in `frontend/src/routes`.
- [ ] Monitoring/metrics dashboard — health/metrics endpoints exist, no admin-facing dashboard page.

## Medium — Partial frontend coverage

- [ ] Forms: only basic CRUD + responses covered; advanced form modules (templates/relationships/reports/query builder/email distribution/workflow — now reachable via `formsAdvancedGroup` in `server.go:259-297`) have no frontend UI yet.
- [ ] Spreadsheets: chart/Excel import-export/collab routes are now mounted (`server.go:206-243`) but still have no frontend UI, and several return mock data server-side (see code-level stubs below).
- [ ] Documents/editor: `editor/+page.svelte` is a thin 24-line stub delegating to a component — verify AI processing/analysis endpoint coverage.
- [ ] Various components contain TODO/placeholder comments: `EmailComposer.svelte`, `WorkflowBuilder.svelte`, `FormBuilderCanvas.svelte` — audit and finish.
- [ ] `frontend/src/lib/api.js` functions that call spreadsheet chart/Excel endpoints (`evaluateFormula`, `getChartSuggestions`, `createChart`, `getSpreadsheetCharts`, `exportSpreadsheetToExcel`, `importExcelToSpreadsheet`, lines ~370-460) now hit live routes — verify they're actually invoked from the UI, since as of the last audit they were unused/dead client-side.

## Adoption / Onboarding Tooling Gaps

- [ ] No first-run admin/setup wizard in the frontend (create first admin account, configure SMTP/DB from the browser instead of hand-editing `.env`).
- [ ] No CI/CD — add `.github/workflows/` for automated build/test/lint on PRs, and for producing the binary releases the README promises.
- [ ] No backup/restore CLI or admin action, despite a backup/recovery service being referenced in docs — expose it.
- [ ] No update/migration tooling beyond `git pull` + restart — no versioned migration runner or upgrade guidance for breaking schema changes.
- [ ] No demo/seed data or `--seed` flag for evaluators to try the product without manual setup.
- [ ] No admin CLI tool (e.g. `tpt-titan admin create-user`) for first-deploy user/role management outside the web UI.
- [ ] `scripts/install.sh` downloads a versioned release tarball that doesn't exist yet in Releases — currently aspirational/non-functional until releases are published.
- [ ] Systemd unit is generated inline by `scripts/install.sh` rather than checked in as an inspectable template (e.g. `deploy/tpt-titan.service`) — check one in for review before running as root.
- [x] Root-level doc clutter/rot: Deleted stale duplicates `TODO - Copy.md`, `TODO 1260108.md`, `TODO 1260113.md`, `TODO_LAYOUT_IMPROVEMENTS.md` (consolidated into this tracked list). `TODO 1260716.md` was reviewed (a completed spreadsheet-layout checklist) and retained.
- [ ] No Makefile/task runner for a single `make setup` / `make dev` entry point across backend/frontend/desktop.
- [ ] `docs/installation.md` env-var example block doesn't match actual `.env.example` keys (`DATABASE_URL`/`ENCRYPTION_KEY` vs `DB_TYPE`/`DB_PATH`) — reconcile so self-hosters aren't misled.
- [ ] No Caddy reverse-proxy example (only Nginx) — Caddy is a common low-friction choice for automatic HTTPS.

## Code-level stubs & mocked logic (highest-impact first)

- [ ] **Admin panel now reachable but still mocked**: `backend/internal/server/server.go:547-560` admin route group is now wired up, but there is still no `/admin` frontend at all (see High section above). Additionally `GetSystemSettings` (`admin.go:584`) returns hardcoded mock settings and `UpdateSystemSettings` (`admin.go:609-622`) never persists — fix the mock logic now that routes are live.
- [ ] **Plugin system completely unreachable**: `backend/services/plugin_system.go` has full load/unload/event-bus/settings logic but still zero HTTP routes (no route file exists for it) and zero frontend UI. Also `GetPluginSettings` (`plugin_system.go:334-343`) always returns `nil, nil` instead of calling `getPluginSettingsFromDB` — the correct implementation (`getPluginSettingsInternal:346-357`) is orphaned dead code, looks like a regression. `downloadPlugin` (`plugin_system.go:703-707`) always errors "not implemented", and `validatePluginSettings` (`:698-701`) is a no-op.
- [ ] **User encryption salt discarded at registration**: `backend/routes/auth/auth.go:309-324` (called from `auth.go:167`) derives a per-user encryption salt via `NewKeyManager` but never stores it (`_ = salt // TODO: Store in user preferences table`) — silently breaks recovery of per-user encryption keys after initial setup. High severity, no user-facing error.
- [ ] **Forms backend is 100% mocked**: `backend/routes/forms.go:52-299` — `GetForms/GetForm/CreateForm/UpdateForm/DeleteForm/GetFormResponses/SubmitFormResponse` are hardcoded mock data with zero DB reads/writes. Forms never actually save/update/delete/store responses despite a real-looking API and a working-looking frontend.
- [ ] **AI job queue fabricates results**: `backend/services/ai_job_queue.go:546-620` — `processDocumentAnalysis`, `processEmailCategorization`, `processSpeechSynthesis`, `processWorkflowOptimization` all `time.Sleep()` then return the same canned fake result regardless of input (e.g. always the same word count, always the same fake audio URL, always the same 3 "optimization suggestions"). Called live from the dispatcher (`ai_job_queue.go:380-386`).
- [ ] **Calendar reminders silently no-op**: `backend/services/calendar_notifications.go:283-304` — `sendSMSReminder`, `sendPushReminder`, `sendInAppReminder` all return `nil` (success) without sending anything. Also SMS reminder wiring uses the user's email address as if it were a phone number (`:232`). Users believe reminders are configured but never receive them.
- [ ] **Email attachment pipeline entirely stubbed**: `backend/services/email_attachments.go:477-534` — `parseEmailParts` always returns empty; `isAttachmentPart` always `false`; filename/content-type/part-data extraction all return hardcoded defaults; `generateThumbnail` is a no-op passthrough; `StorageService.SaveFile/GetFile/DeleteFile` (`:521-534`) are all no-ops. Also virus scanning (`:460-466`) always marks every attachment `IsSafe: true` via a `"simulated_scanner"` — no real scanning occurs. `backend/models/email.go:178` always reports `HasAttachments: false` regardless of actual content.
- [ ] **Equation/handwriting rendering fabricated**: `backend/routes/math_canvas_routes.go:88-114` `GenerateEquationImage` returns literal placeholder text claiming to be a PNG/SVG/PDF instead of rendering LaTeX. `backend/services/handwriting_recognition.go:394-398` `buildComplexExpression` always returns the same hardcoded `"x^{2} + y^{2} = z^{2}"` regardless of input strokes; `:588,701-704` SVG/PNG/PDF generation return placeholder text; `:353-356` `strokesIntersect` always returns `true`.
- [ ] **Spreadsheet charts/Excel export ignore real data**: `backend/routes/spreadsheet_chart_routes.go:71-84,96-106` — `CreateChart` never writes to DB; `GetCharts` ignores the spreadsheet ID and always returns one hardcoded "Sales Data" chart. `backend/routes/spreadsheet_excel_routes.go:62-85` `ExportExcel` always exports the same hardcoded 3-row mock sheet regardless of which spreadsheet was requested.
- [ ] **Document PDF export is fake**: `backend/routes/document_export.go:78-87` `exportToPDF` returns a JSON message `"PDF export functionality would be implemented here"` with HTTP 200 instead of a binary PDF.
- [ ] **Backup checksum is a hardcoded string**: `backend/services/backup.go:711-715` `calculateFileChecksum` always returns `"checksum_placeholder"` instead of a real SHA-256 hash — backup integrity is never actually verified.
- [ ] **Local/offline speech-to-text mostly dead**: `backend/services/speech_local_stt.go` — Windows and macOS transcription always error "not yet implemented"; Kaldi/Julius paths write the audio then unconditionally error; only PocketSphinx on Linux can work, and only with hardcoded model paths pre-installed. (Errors are honestly surfaced rather than faked.)
- [ ] **AI provider/speech API keys stored in plaintext**: `backend/routes/ai.go:80` and `backend/routes/speech.go:100` both have `// TODO: Encrypt this` on stored API keys.
- [ ] **Biometric face-ID recovery is fake if ever exposed**: `backend/utils/crypto.go:227-240` `compareFaceTemplates` is a raw byte-equality check, not real biometric matching. Currently dormant (only exercised in tests, not wired to a route) — leave unwired or implement properly before ever exposing it.
- [ ] **AI hardware/model-upgrade logic is placeholder math**: `backend/services/ai.go:270-298`, `backend/services/hardware_service.go:24`, `backend/services/model_service.go:247` — `checkOpenRouterUpgrades` never queries OpenRouter; `DetectHardware` is a stub; `isUpgradeCandidate` checks compatibility against a hardcoded `1.0` size instead of real model size. `backend/routes/ai.go:288-297` `ApplyModelUpgrade` is an honestly-labeled placeholder (returns 202 "not yet implemented").
- [ ] **Workflow scheduling can duplicate-fire**: `backend/services/workflow_service.go:376-388` — rescheduling a workflow never removes the old cron job (`// TODO: Implement job removal by workflow ID`), so re-scheduling multiple times causes duplicate executions (e.g. duplicate emails/notifications).
- [ ] **Workflow AI optimization never finds anything**: `backend/services/workflow_ai_service.go:663-682` — `findLongSequentialChains`, `findRedundantOperations`, `analyzeErrorHandling`, `findUnusedDataFlow` are placeholders that always return empty/false.
- [ ] **Task integration returns IDs not tied to real rows**: `backend/services/task_integration.go:16,261,480,506,577` — `formService` field is a placeholder interface; `taskID := uuid.New()` appears 3× as a stand-in for the ID that should come back from actual task creation, risking orphaned/mismatched references.
- [ ] Document analysis metadata is wrong: `backend/routes/documents.go:502` always saves filename as literal `"uploaded_file"`; `:531` hardcodes `EnableLocalAI: true` regardless of user config; `:585` `ProcessingTime` always `0`.
- [ ] Monitoring "active users" metric always reports 0 (`backend/services/monitor.go:240`); DB optimizer's slow-query analysis returns two hardcoded fake queries instead of using `pg_stat_statements` (`backend/services/db_optimizer.go:67-89`).
- [ ] Calendar group-sharing permission check always errors "not implemented" (`backend/services/calendar_sharing.go:685-688`) — confirm whether "share with group" is exposed in the UI before treating as low priority.
- [ ] DOCX export silently drops headers/footers (`backend/services/docx_export.go:569`).
- [ ] Export format gaps in forms: PDF/JSON export in `form_reporting.go:614-630` return placeholder text instead of real files; Excel export silently substitutes CSV; `form_query_builder.go:585-590` JSON export errors "not implemented".
- [ ] Frontend dead buttons/UX gaps in Forms page (`frontend/src/routes/forms/+page.svelte`): `viewResponses()` (`:103`) only logs to console — clicking does nothing visible; drag-and-drop reordering (`:109`) isn't persisted to backend, resets on reload.
- [ ] `frontend/src/lib/components/AdvancedReportsModal.svelte:444` — at least one report type shows "This report type visualization is coming soon." (honestly labeled, still incomplete).

## Security audit

- [ ] Findings pending — security-review pass in progress; append results here once complete.

## Innovation / UX ideas

- [ ] To be proposed once the audits above are complete.
