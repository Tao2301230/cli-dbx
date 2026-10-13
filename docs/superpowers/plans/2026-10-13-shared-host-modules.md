# Shared Host Modules Implementation Plan

> **For agentic workers:** Use subagent-driven-development with independent file ownership and final code review. The user has approved execution in this conversation.

**Goal:** Build DBX once and distribute four compact, independently loadable Go database driver modules.

**Architecture:** A pinned Go host loads prebuilt native modules; its internal worker reuses the existing private protocol and database executor. Offline packaging removes redundant runtime-linker descriptions where safe.

**Tech Stack:** Go 1.26.8, patched pkujhd/goloader at 2695cfc5350631cc3becc1409162b12a885471b5, Python build orchestration.

## Global Constraints

- Host and module use identical Go toolchain, target and common dependencies; bind modules to the exact host SHA-256.
- No vendor driver imports in root application source; independent modules contain vendor implementations.
- Keep CLI policy checks, redaction, transaction scope and process cancellation semantics.
- Raw module bytes are the size acceptance metric; report gzip separately.
- Preserve GC, type and stack metadata and keep loaded modules alive until worker exit.
- macOS execution and Windows cross-compilation are distinct validation claims.

## Task 1: Compact format and offline linking

Files: `internal/dbmodformat/`, `third_party/goloader/` (separate ownership for prelink changes).

Interfaces: retain `Encode(*link.Linker, io.Writer) error` and `Decode(io.Reader) (*link.Linker, error)`; add optional offline optimization API with statistics.

- [x] Add failing round-trip, truncated/corrupt format, deterministic encoding and raw size regression tests.
- [x] Implement prefix-compressed symbol dictionary and packed relocation records while keeping executable data unchanged.
- [x] Prove and implement safe internal relocation resolution separately; preserve unsupported relocation types and necessary jump reserves.
- [x] Run tests using `GOROOT=<private Go 1.26.8> GOTOOLCHAIN=local GOCACHE=/private/tmp/dbx-driver-process-go-cache go test ./internal/dbmodformat`.

## Task 2: Reproducible build and module entrypoints

Files: `cmd/dbmod-pack/`, `modules/`, `scripts/build-modules.py`, build tests and dependency pins.

Interfaces: packer receives final host path, module object/dependency list, engine and output path. Emit raw `<engine>.dbmod` plus JSON manifest containing schema, engine, Go version, OS/arch, host SHA-256 and module SHA-256. Export module symbol `main.Register` with signature `func()`.

- [x] Test argument/manifest output and deterministic packaging failures before implementation.
- [x] Reuse archived dependency closure building with exact final host symbols.
- [x] Build each entrypoint from pinned dependencies; PostgreSQL/Oracle/DM/SQL Server register existing Go drivers.
- [x] Provide SDK preparation in a private output directory and explicit `--go`/`--out` controls.
- [x] Build host first, then package modules; emit actual size report in the output directory.

## Task 3: Host integration

Files: `internal/drivermodule/`, `internal/driverprocess/client.go`, configuration validation, `cmd/dbx/`, documentation.

Interfaces: module discovery follows existing driver roots, explicit module path wins. Worker is the current executable with internal arguments containing only engine and module path; request credentials travel on the existing pipe. Parent validates manifest before launching; worker validates again before loading.

- [x] Add failing discovery, hash/host/toolchain/platform mismatch and config tests.
- [x] Add a common process call path accepting executable plus non-secret arguments.
- [x] Implement worker validation, decode, symbol registration, locked OS-thread load and `main.Register` call; then invoke existing `driverprocess.Serve`.
- [x] Keep legacy builds available when module runtime is not selected; enable loader through an explicit build tag if needed.
- [x] Verify policy refusal never loads a module; existing driverprocess tests continue passing.

## Task 4: End-to-end validation and review

- [x] Run root tests and loader/format tests under the pinned toolchain.
- [x] Build real macOS host and four raw modules; smoke load/register/GC with no real database claims.
- [x] Validate PostgreSQL protocol query through the actual CLI against a local fixture.
- [x] Cross-build the Windows amd64 host; module packaging requires a native Windows environment and runtime execution remains unverified.
- [x] Record raw size comparison, remaining limitations and reproducible commands in `docs/shared-host-modules.md`.
- [x] Review the complete diff, fix material findings, re-run affected checks.

## Task 5: Move repeated generic implementations into the host

User-directed refinement: shrink the four optional database modules by retaining
their repeated generic implementations in DBX. Preserve vendor isolation,
existing process behavior, platform/toolchain checks and database capabilities.

Files: `internal/drivermodule/shared.go`, `internal/drivermodule/load.go`,
`scripts/tests/test_module_host.py`, `docs/shared-host-modules-validation.json`.

- [x] Measure local symbols in all four current module graphs and test shared
  host references with an isolated Go build overlay before editing production.
- [x] Add an executable symbol regression test; verify the previous host fails
  because it omits the six representative common implementations.
- [x] Retain certificate, text, number, JSON and container functions through
  `registerSharedHostTypes(map[string]uintptr)` after `goloader.RegSymbol`.
- [x] Rebuild using `scripts/build-modules.py` and run the shared-symbol test,
  all four offline load checks, PostgreSQL CLI protocol fixture and root tests.
- [x] Record actual raw bytes before/after, host growth and combined savings;
  review the refinement without claiming live database or Windows execution.
