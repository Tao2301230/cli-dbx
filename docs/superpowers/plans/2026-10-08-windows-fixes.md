# Windows Fixes Implementation Plan

> **For agentic workers:** Use test-driven-development and verification-before-completion. Independent tasks may run in parallel with exclusive file ownership. No Git mutations in this implementation turn.

**Goal:** Migrate the usable Windows fixes into the independent-driver branch, preserving its safety boundaries and dependency split.

**Architecture:** A standard-library BOM reader sits before existing decoders. SQLite filesystem names are encoded into file URIs while retaining existing query semantics. Release packing accepts explicit executable names independently of host permissions.

**Tech Stack:** Go 1.22 root module, existing nested Go/ODBC modules, Python 3, Shell and PowerShell.

## Global Constraints

- Keep MySQL/SQLite built in and all other adapters in independent processes.
- Retain protocol v1, `driver_path`, manifest validation, private credentials, SQL dialect checks, import limits, `UseNumber` and EOF checks.
- Do not change `go.mod`, `go.sum`, driver modules or Git history/index.
- Preserve `mode=ro`, memory/cache and `_pragma` behavior when fixing SQLite filenames.
- Use `/private/tmp` for builds and test outputs; do not add generated binaries or logs.

---

### Task 1: UTF-8 input and strict JSON

**Files:** Create `internal/textfile/utf8.go`, `internal/textfile/utf8_test.go`; modify `internal/config/config.go`, config tests, `internal/cli/app.go`, CLI tests.

**Interfaces:** `Reader(io.Reader) *bufio.Reader` and `ReadFile(string) ([]byte, error)` remove one optional leading BOM. Consumers retain their decoders and all policy checks.

- [x] Add failing tests to existing readers for BOM SQL, transaction JSON, CSV headers and JSON import. Include actual large-integer import and transaction rejection before any write.
- [x] Run `env GOCACHE=/private/tmp/dbx-driver-process-go-cache GOPROXY=off go test -count=1 ./internal/config ./internal/cli` and record expected BOM/trailing-content failures.
- [x] Add the minimal BOM helper using `Peek(3)`/`Discard(3)` and `bytes.TrimPrefix`; replace only input readers. Add a second TX decode requiring `io.EOF`.
- [x] Verify empty/short/fragmented readers, embedded and double BOM, unknown fields, trailing objects/junk and integer preservation. Apply `USERPROFILE` isolation and path quoting to config/CLI test helpers.
- [x] Run the relevant package tests and report changed files plus red/green evidence.

### Task 2: SQLite paths with compatible URI options

**Files:** Modify `internal/conn/resolve.go`; add conn/SQLite integration regression tests.

**Interfaces:** `buildSQLiteDSN(Spec) string` keeps explicit URI and memory forms; encodes literal filesystem path separately from the existing raw query.

- [x] Add failing tests that create, close and independently reopen files containing `#`, `%23`, `%20`, Chinese and spaces; verify no substitute file is created.
- [x] Add behavior checks that `path=...?...mode=ro` rejects writes and `mode=memory&cache=shared` retains memory/shared semantics. Include `_pragma` and explicit URI controls.
- [x] Run focused conn/db tests and observe the path failures on the existing implementation.
- [x] Encode filesystem path characters into a file URI without double-encoding explicit URIs or dropping query parameters; preserve Windows drive and UNC paths.
- [x] Run conn/db tests and report the actual path and read-only behaviors verified.

### Task 3: Explicit executable permissions in release archives

**Files:** Modify `scripts/reproducible-package.py`, `scripts/test-reproducible-package.py`, `scripts/release/build.sh`, `scripts/release/build.ps1`.

**Interfaces:** `pack(stage, output, executables=())`; repeatable CLI `--executable PATH`; existing callers retain defaults.

- [x] Add a failing test packaging a `dbx` file with mode `0644`; assert `0755` in both TAR and ZIP, document files `0644`, and identical bytes after changing source mtimes.
- [x] Run `PYTHONDONTWRITEBYTECODE=1 python3 -B scripts/test-reproducible-package.py -v` and observe the permissions/API failure.
- [x] Implement validated stage-relative executable names and pass `$binary_name`/`$binaryName` from the two core release scripts. Reject invalid/missing paths before creating output.
- [x] Run Python tests and existing Go release/connector tests.

### Task 4: Documentation, examples and final integration

**Files:** Modify `README.md`, `AGENTS.md`, `docs/beginner-guide.md`, `internal/app/app.go`, `internal/app/commands.go`, `internal/cli/help.go`, app tests, native/ODBC examples, `VERSION`, `connector/cli.json` and relevant connector references.

- [x] Apply the attachment's Windows-compatible app test path/environment fixes, keeping current command behavior.
- [x] Document UTF-8 BOM, explicit UTF-8 PowerShell writes, SQLite path/URI behavior and strict transaction JSON. Correct examples to use independent `driver_path` and vendor `odbc_driver`.
- [x] Set `VERSION` to `v0.2.0-driver-process.4` and connector `minVersion` to `0.2.0-driver-process.4`; keep references to already-built `.3` drivers explicit.
- [x] Run root and nested-module tests, Python packaging/catalog tests, build the CLI, inspect required help and run BOM/SQLite smoke workflows.
- [x] Cross-compile core for Windows AMD64/ARM64 and check core dependency closure excludes the four vendor drivers and ODBC.
- [x] Review the combined working-tree diff, resolve findings and leave all modifications uncommitted.


## Integration verification before Git review (2026-10-08)

- Root: `env GOCACHE=/private/tmp/dbx-driver-process-go-cache GOPROXY=off go test -count=1 ./...` passed after integration.
- Nested modules: postgres, sqlserver, dm compiled; oracle tests passed with the existing Go 1.26.6 toolchain; ODBC adapter and fork tests passed using the existing temporary unixODBC/libtool installation. No live vendor database was exercised.
- Python: reproducible-package 5 tests and driver-download catalog 6 tests passed; `bash -n scripts/release/build.sh` passed. PowerShell runtime was unavailable.
- CLI: native ARM64 build passed. Actual BOM TOML/SQL/CSV/JSON/TX workflow preserved embedded FEFF and integer 9007199254740993; invalid transaction suffixes were rejected without writing; the literal Chinese/#/%23/%20 filename reopened with Python SQLite and mode=ro rejected writes. Public help and version output verified.
- Windows: main executable cross-compilation passed for AMD64 and ARM64; SQLite conn/db test binaries cross-compiled for AMD64. These are build checks, not Windows runtime tests; UNC requires a Windows share for runtime validation.
- Architecture: root module remains Go 1.22; go.mod, go.sum, driver modules, private protocol, secret handling and Git ignore rules have no changes. Main dependency closure includes MySQL/SQLite and excludes PG/Oracle/DM/SQL Server/ODBC vendor drivers.
- App test fixtures already used portable temporary directories and quoted TOML paths, so no additional app fixture change was required. Migrated test fixes are in config/CLI.
- Independent code review found no actionable defects; `git diff --check` passed. All changes remain uncommitted and unstaged. No release artifacts were added to the repository.

## Pre-commit review follow-up (2026-10-08)

- A later Git review reproduced two test fixture issues: JSON import columns were compared in map iteration order (15 failures in 100 repetitions), and SQLite test paths were interpolated into a TOML literal without escaping single quotes.
- JSON regression assertions now sort only JSON columns before comparison; CSV still checks header order. The SQLite fixture uses quoted TOML paths and includes a filename containing a single quote. The new filename failed before the fixture fix and passed afterward.
- Fresh checks passed with Go 1.26.4 on macOS 26.5 ARM64: the JSON regression repeated 100 times, SQLite path tests with a single quote in TMPDIR, root `go test -count=1 ./...`, root race tests, five reproducible-package tests, six catalog unit tests, Shell syntax and `git diff --check`. The updated SQLite test program also cross-compiled for Windows AMD64; it was not run on Windows.
- README now describes `.3` driver packages as existing local builds, without claiming a network release. `.4` release archives were not generated or published.
- The current local download catalog still fails full verification because the PORTING.md digest differs from SHA256SUMS. Catalog unit tests passing does not resolve this separate distribution issue. The six existing ODBC fork unsafe.Pointer diagnostics also remain outside this change.
- Windows native execution, PowerShell, UNC and real vendor databases remain unverified. This follow-up prepares local source commits only; it does not modify ignore rules, existing archives, vendor installations or database credentials.
