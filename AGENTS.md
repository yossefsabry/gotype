# AGENTS

## Purpose
- This repo is a Go TUI typing app (tcell-based) with a strong focus on speed and responsiveness.
- Prefer many small files over large files; keep files tight and cohesive.
- Performance is more important than cosmetic changes; favor low-allocation paths.
- Concurrency is encouraged when it improves responsiveness or I/O throughput.
- Keep changes incremental and easy to review.

## Repo overview
- Module: github.com/yossefsabry/gotype
- Entry point: `main.go` calls `app.Run()`.
- Core app logic: `internal/app` (model, input, layout, rendering, persistence).
- Persistence layer: `internal/storage` (JSON load/save, path utilities).
- UI runtime: `tcell/v2`.
- Go version: 1.24.6 (from `go.mod`).

## Build / run / lint / test
- Build: `go build ./...`
- Run locally: `go run .`
- Format (always gofmt): `gofmt -w path/to/file.go`
- Vet (lint-ish): `go vet ./...`
- Tests (all packages): `go test ./...`
- Single package tests: `go test ./internal/app`
- Single test name (example): `go test ./internal/app -run TestName -count=1`
- Single test with pattern: `go test ./internal/app -run 'TestFoo/Case1' -count=1`
- Race check when touching concurrency: `go test ./... -race`
- Benchmarks: `go test ./internal/app -run '^$' -bench . -benchmem`
- No additional lint tools are configured; keep to gofmt/go vet unless instructed.

## Code style (Go)
- Follow gofmt and standard Go layout; no manual alignment.
- Keep functions small and focused; split into new files instead of growing files.
- Prefer explicit early returns over deep nesting.
- Keep logic close to data ownership (e.g., Model methods mutate model fields).
- Avoid global state; prefer structs + methods.
- Use ASCII-only in source unless the file already uses Unicode.

## Imports
- Group imports: standard library, blank line, third-party, blank line, local module.
- Use explicit import paths; no dot imports.
- Avoid unused imports; gofmt/goimports should clean them.
- Prefer stable, minimal import sets to keep compile time low.

## Formatting and layout
- Use gofmt for all Go files.
- Keep line lengths reasonable, but prefer clarity over manual wrapping.
- Use short variable names only for tight scopes (i, j, now, ok).
- Keep constants near their usage (e.g., layout or timing constants in the same file).
- Prefer small helper functions over long inline blocks.

## Types and data structures
- Use structs to group related state (App, Model, Renderer, Layout).
- Prefer value types for small structs; pointers for shared or mutable state.
- For slices, preallocate when size is known or can be estimated.
- Use time.Duration for time values; store seconds only at serialization boundaries.
- Keep JSON tags stable; changes to `internal/storage` affect saved data.

## Naming conventions
- Exported types/functions: PascalCase.
- Unexported names: camelCase.
- Use clear domain names (Model, Renderer, Persister, Layout, Generator).
- Boolean names read as predicates (finished, running, needsRender).
- Keep method names verb-based (Reset, Update, Render, Save, Load).

## Error handling
- Return errors immediately; avoid swallowing errors silently.
- Wrap errors only when adding actionable context; otherwise return as-is.
- For errors you intentionally ignore, assign to `_` and add a short comment.
- Prefer sentinel handling (os.IsNotExist) only when necessary.
- UI code should fail fast in `main` and exit with non-zero on fatal errors.

## Concurrency and performance
- Use goroutines/channels to avoid blocking UI (see `app.loop`, `Persister`).
- Favor buffered channels for UI/event pipelines to prevent stalls.
- Avoid goroutine leaks: signal with done/quit channels and close when exiting.
- Keep hot loops allocation-free; reuse slices and structs when possible.
- Prefer tickers for periodic updates; stop them when done.
- Avoid fmt.Sprintf in tight loops; build strings once per frame if possible.
- Do not spawn goroutines inside render paths unless strictly necessary.

## Event loop and input
- Poll events in a dedicated goroutine and send through a channel.
- Close the event channel if PollEvent returns nil.
- Use one `time.Now()` per input handling path and pass it down.
- Handle key events through `Model.HandleKey`; click events through `Model.HandleClick`.
- Normalize runes (e.g., upper to lower) before storing `LastKey`.

## UI and rendering (tcell)
- Rendering must be deterministic and fast; minimize full-screen clears.
- Keep screen writes bounded with bounds checks.
- Use `Layout` to compute positions; avoid inline geometry scattered across files.
- Respect model state for cursor, stats, and highlight styles.
- Sync theme changes through `Renderer.syncTheme`.
- Favor `tcell.Style` reuse via `Styles` instead of recomputing colors.

## Layout and text
- Recalculate layout on resize and mode changes.
- Keep text wrapping logic in `textwrap.go` helpers.
- Avoid reflowing text on every keystroke unless required.
- When extending text, use `Generator.Extend` to keep word spacing consistent.

## Themes and styles
- Themes live in `themes.go`; add new themes as new entries in `themeOptions`.
- Keep theme IDs stable; they are persisted in preferences.
- Use `ThemeLabel` for display text; keep labels short.
- Maintain consistent contrast for cursor, key active, and error colors.

## Persistence and storage
- Storage is JSON under a default path (see `internal/storage`).
- Always keep BestScores non-nil maps after load.
- Use atomic write patterns (temp file + rename).
- Persist preferences on meaningful changes to avoid excess writes.
- Avoid blocking UI when saving; use the Persister channel.

## Testing guidance
- There are currently no tests; add `_test.go` next to the file under test.
- Keep unit tests small and pure; avoid UI integration tests unless needed.
- Prefer table-driven tests for options and layout calculations.
- Add benchmarks for hot paths (text wrapping, rendering helpers).
- Use `-count=1` when you need to bypass test caching.

## File organization
- Keep new features split across small focused files.
- For new concerns, create a new file under `internal/app` or `internal/storage`.
- Prefer adding a file rather than enlarging an existing one beyond ~200 lines.
- Keep package-level vars limited and localized (e.g., word lists).

## Dependencies
- Primary external dependency is `github.com/gdamore/tcell/v2`.
- Add new dependencies only when they are lightweight and clearly justified.
- Keep dependency footprint small to preserve startup time.

## Logging and diagnostics
- Keep logging minimal; use stderr only for fatal errors.
- Avoid log spam inside render or hot loops.
- Use debug prints only temporarily and remove before committing.
- Keep user-facing errors short and actionable.

## Git hygiene
- Do not commit unless explicitly requested.
- Avoid rewriting history or force pushes.
- Keep changes scoped; avoid unrelated formatting churn.
- Update this file when tooling or conventions change.

## Project rules (owner preferences)
- Small files are preferred over large files.
- Speed and performance are top priorities.
- Concurrency (goroutines) is encouraged when it improves responsiveness.
- Avoid heavy abstractions that add overhead without clear payoff.

## Cursor/Copilot rules
- No `.cursor/rules`, `.cursorrules`, or `.github/copilot-instructions.md` found.
