# Flat Drivers Implementation Plan

> **For agentic workers:** Use subagent-driven-development for independent packaging and review tasks; keep edits in the existing experimental branch. Steps use checkbox syntax for tracking.

**Goal:** Install multiple driver versions in one flat directory and select them through `driver_path`.

**Architecture:** Normalize the new connection field in config, propagate it through `conn.Spec`, and load exactly its sibling manifest. Package each driver as two versioned root files; preserve the process protocol.

**Tech Stack:** Go, TOML, JSON, existing Go archive packager and shell/PowerShell entrypoints.

## Global Constraints

- Follow `docs/superpowers/specs/2026-09-30-flat-drivers-design.md`.
- Preserve unrelated untracked user files and the MySQL/SQLite builtin boundary.
- Retain existing archive formats; no installer or version discovery.
- Validation precedes local commit; no push or merge.

### Task 1: Driver loading and configuration

Files: `internal/config/config.go`, `internal/conn/resolve.go`, `internal/driverprocess/manifest.go`, `internal/cli/app.go`, and adjacent tests.

- [x] Add tests proving flat installation works and old layout fails; add `driver_path` config/resolve/CLI tests, multiple versions, no fallback, manifest redirection, symlink, directory, and oversized manifest cases.
- [x] Observe new tests fail before implementation.
- [x] Add `DriverPath string` with TOML key `driver_path`, normalize via `resolveResourcePath`, propagate to `conn.Spec` and driver error envelopes. Choose the explicit path or flat unversioned default, use `path + ".json"`, and require manifest executable to equal `filepath.Base(path)` before hashing and execution.
- [x] Run `GOCACHE=/private/tmp/dbx-driver-process-go-cache go test ./internal/config ./internal/conn ./internal/driverprocess ./internal/cli`.

### Task 2: Flat versioned archives

Files: `scripts/release/package-driver/main.go`, `scripts/release/package-driver/main_test.go` only.

Interface: `packageDriver(binary, outputDir, engine, version, goos, arch string) (string, error)` remains unchanged. Runtime consumes the full basename plus `.json` sidecar.

- [x] Update package assertions to require exactly two root entries. Set `exe := "dbx-driver-" + engine + "-" + strings.TrimPrefix(version, "v")`; add `.exe` on Windows, then manifest name `exe + ".json"`.
- [x] Observe tests fail; implement the filename and entry changes without changing archive names/formats.
- [x] Test extracting two versions into one directory preserves four distinct files with matching manifests, bytes, hashes and permissions. Preserve reproducibility and input validation tests.
- [x] Run `GOCACHE=/private/tmp/dbx-driver-process-go-cache go test ./scripts/release/package-driver`.

### Task 3: Documentation, release, integration

Files: `VERSION`, `connector/cli.json`, connector skill metadata/references, README, AGENTS, both CLI help surfaces, `docs/driver-process.md`.

- [x] Bump release to `v0.2.0-driver-process.2`; document flat Mac/Windows filenames, explicit version selection, relative paths, default unversioned fallback and migration from old layout.
- [x] Run fresh root and nested module tests, race tests for changed process loading, and `git diff --check`.
- [x] Rebuild core plus all drivers for darwin/arm64 and windows/amd64. Inspect archives and smoke test actual Mac drivers with `driver_path` against a closed local port, checking secrets are absent from output.
- [x] Review the complete change, resolve findings, save a local commit, report artifacts and remaining live-database validation limits.

## Validation record

- New flat process fixtures failed against the original per-engine loader. New config, resolution, explicit-version and CLI diagnostic tests also failed before their implementation.
- `GOCACHE=/private/tmp/dbx-driver-process-go-cache go test -count=1 ./...`: passed.
- `go test -race -count=1 ./internal/driverprocess`: passed with the same GOCACHE.
- `go test -count=1 ./...` in all four driver modules: passed (Oracle has adapter tests; other entrypoints report no test files).
- Packaging tests verify two-file root layout, multiple-version extraction, full-basename sidecars including `.exe.json`, manifest/binary digests, executable permissions, deterministic bytes and unchanged archive names.
- Core, Platform connector and all four drivers rebuilt for darwin/arm64 and windows/amd64 into `dist/v0.2.0-driver-process.2/`.
- Actual Mac driver smoke checks passed with absolute, home-relative and TOML-relative `driver_path`, deliberately conflicting root settings, and all engines installed flat together. Each reached its vendor connection attempt against localhost port 1 and returned a sanitized database error. Help/version checks passed.
- Windows archives and PE headers were checked along with manifest binding and hashes; Windows runtime and PowerShell execution were not performed.
- Independent whole-change review: spec and quality PASS, no actionable findings. `git diff --check` passed.
- No real vendor database compatibility testing, installation to the user's system, remote push or merge was performed.
