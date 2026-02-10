# Understanding gotype (step-by-step)

This guide walks the repo in a deliberate reading order. Each step
names the file, its purpose, and the functions it defines.

1. `AGENTS.md` - Repo rules and performance priorities for
contributions; Functions: none.
2. `README.md` - User-facing overview, install, usage, and build
instructions; Functions: none.
3. `dooooit` - Improvement plan and phased TODOs for future work;
Functions: none.
4. `go.mod` - Module name and dependency list (Go 1.24.6, tcell);
Functions: none.
5. `go.sum` - Dependency checksums for reproducible builds;
Functions: none.
6. `main.go` - Entry point; Functions: `main` (calls `app.Run`,
prints error, exits non-zero on failure).
7. `internal/app/doc.go` - Package comment for `app`; Functions:
none.
8. `internal/storage/doc.go` - Package comment for `storage`;
Functions: none.
9. `internal/app/app.go` - App wiring and event loop; Functions:
`Run` (init screen, model, layout, persister), `(*App).loop` (poll
events, tick, render), `(*App).syncPersistence` (save prefs and best
scores, finalize results).
10. `internal/app/input.go` - Input routing; Functions:
`(*Model).HandleKey` (keyboard actions and typing),
`(*Model).registerKey` (normalize and store last key),
`(*Model).HandleClick` (hit-test regions), `(*Model).applyRegion`
(toggle options, modes, selectors, themes).
11. `internal/app/model.go` - Core state and mechanics; Types:
`Mode`, `Options`, `Timer`, `Stats`, `Text`, `UIState`, `Model`;
Functions: `NewModel`, `(*Model).Reset`, `(*Model).StartTimer`,
`(*Model).Update`, `(*Model).AddRune`, `(*Model).Backspace`,
`(*Model).BackspaceWord`, `(*Model).UpdateDerived`,
`(*Model).SetMessage`, `(*Model).ensureTarget`,
`(*Model).removeTypedRange`, `(*Model).WordsLeft`,
`(*Model).SetTheme`.
12. `internal/app/model_helpers.go` - Small helpers for model
bookkeeping; Functions: `(*Model).elapsedForStats`,
`(*Model).focusActive`, `(*Model).syncLayoutFocus`,
`(*Model).resetMistakes`, `(*Model).recordMistake`,
`(*Model).recalculateStreak`.
13. `internal/app/persist.go` - Persistence glue and score
comparison; Types: `Persister`; Functions: `NewPersister`,
`(*Persister).loop`, `(*Persister).Save`, `(*Persister).Close`,
`loadPersistedData`, `preferencesFromModel`, `applyPreferences`,
`scoreKey`, `updateBestScore`, `modeToString`, `modeFromString`.
14. `internal/storage/types.go` - JSON data shapes for disk; Types:
`Preferences`, `BestScore`, `Data`; Functions: none.
15. `internal/storage/store.go` - JSON load/save with atomic temp
file; Functions: `Load`, `Save`.
16. `internal/storage/path.go` - Default config path builder;
Functions: `DefaultPath`.
17. `internal/app/layout.go` - UI geometry and clickable regions;
Types: `Region`, `Layout`; Functions: `Region.Contains`,
`(*Layout).Recalculate`, `labelForRegion`; Data: `regionLabels`,
`modeOrder`.
18. `internal/app/selectors.go` - Time/word presets for the top
bar; Type: `SelectorOption`; Data: `selectorOptions`,
`selectorOrder`; Functions: `selectorByID`, `selectorLabel`.
19. `internal/app/themes.go` - Theme catalog and lookup; Type:
`Theme`; Data: `themeOptions`, `themeRegionPrefix`; Functions:
`ThemeOptions`, `DefaultThemeID`, `ThemeByID`, `ThemeLabel`,
`ThemeRegionID`, `ThemeIDFromRegion`.
20. `internal/app/style.go` - Convert theme colors to tcell styles;
Type: `Styles`; Functions: `NewStyles`, `hexColor`.
21. `internal/app/render.go` - Renderer core and screen
primitives; Type: `Renderer`; Functions: `NewRenderer`,
`(*Renderer).Render`, `(*Renderer).syncTheme`,
`(*Renderer).fillScreen`, `(*Renderer).fillLine`,
`(*Renderer).panelStyle`, `(*Renderer).drawString`,
`(*Renderer).setContent`.
22. `internal/app/render_panels.go` - Top bar, stats, and footer
text; Functions: `(*Renderer).drawTopBar`,
`(*Renderer).drawThemeMenu`, `(*Renderer).drawStats`,
`(*Renderer).drawFocusStatus`, `(*Renderer).drawFooter`,
`(*Renderer).styleForRegion`, `formatDuration`.
23. `internal/app/render_text.go` - Main typing text rendering;
Functions: `(*Renderer).drawText`, `(*Renderer).drawLine`,
`(*Renderer).centeredLineX`, `lineVisualWidth`,
`(*Renderer).drawRuneBlock`; Constant: `usePlainText` (switches
big glyphs off).
24. `internal/app/bigtext.go` - 3x3 glyphs for scaled text; Type:
`glyph`; Data: `bigGlyphs`; Functions: `glyphForRune`; Constant:
`textScale`.
25. `internal/app/textwrap.go` - Word-based line breaking; Types:
`Line`, `wordRange`; Functions: `buildLines`, `collectWords`,
`lineIndexFor`, `defaultStartLine`; Constants: `maxWordsPerLine`,
`maxVisibleLines`.
26. `internal/app/linecache.go` - Cached line layout per width;
Type: `LineCache`; Functions: `(*Model).bumpTargetVersion`,
`(*Model).linesForWidth`.
27. `internal/app/keyboard.go` - Logical keyboard layout; Type:
`Key`; Data: `keyboardRows`; Functions: `newKeyRow`, `newKey`;
Constants: `keyboardKeyWidth`, `keyboardSpaceWide`.
28. `internal/app/render_keyboard.go` - Keyboard rendering and
adaptive width math; Functions: `(*Renderer).drawKeyboard`,
`keyboardHeight`, `keyboardRowLayout`, `minKeyWidth`,
`sumKeyWidths`, `(*Renderer).drawKey`,
`(*Renderer).keyboardStartY`; Constants: `keyboardKeyGap`,
`keyboardMinKeyGap`, `keyboardRowGap`, `keyboardFooterGap`,
`keyboardMinKeyWidth`, `keyboardMinKeyPadding`.
29. `internal/app/results.go` - Results and comparison logic;
Type: `ResultsState`; Functions: `(*Model).ResetResults`,
`(*Model).FinalizeResults`, `isBetter`, `isWorse`.
30. `internal/app/render_results.go` - Renders results and
best-score line; Functions: `(*Renderer).drawResults`.
31. `internal/app/stats.go` - WPM history and consistency; Type:
`StatsHistory`; Functions: `(*StatsHistory).Reset`,
`(*StatsHistory).Record`, `(*StatsHistory).StdDev`.
32. `internal/app/review.go` - Review scrolling after finish;
Functions: `(*Model).ResetReview`, `(*Model).InitReviewStart`,
`(*Model).ScrollReview`, `(*Model).ReviewTop`,
`(*Model).ReviewBottom`.
33. `internal/app/normalize.go` - Rune normalization helper;
Functions: `normalizeRune`.
34. `internal/app/words.go` - Word and punctuation generator;
Type: `Generator`; Functions: `NewGenerator`,
`(*Generator).Build`, `(*Generator).Extend`,
`(*Generator).nextWord`; Data: `defaultPunct`, `defaultWords`.
35. `internal/app/render_test.go` - Render benchmark using a
simulation screen; Functions: `BenchmarkRender`.
36. `internal/app/words_test.go` - Generator benchmark; Functions:
`BenchmarkGeneratorBuild`.
37. `internal/app/textwrap_test.go` - Line wrapping benchmark;
Functions: `BenchmarkBuildLines`.
38. `images/UI.png` - Screenshot of the UI; Functions: none.
