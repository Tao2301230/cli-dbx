# Shared host and driver modules

The user approved this architecture on 2026-10-13 after reviewing the DuckDB comparison. Implement on `codex/dbx-shared-host-modules`, starting from `b5b7c1b`.

DBX provides one Go runtime, shared packages and database execution code. PostgreSQL, Oracle, DM and SQL Server are independently distributed native Go modules. Modules load when selected; users do not need a Go SDK. The first implementation pins Go 1.26.8 and binds a module to the exact host binary and target platform. Module files are trusted executable code, not a sandbox.

Preserve the existing process protocol and cancellation boundary by launching the same DBX executable in an internal module-worker mode. Each worker loads one module and serves one complete operation, including whole transactions/imports. There is one host file on disk; workers have separate runtime state in memory. This preserves the existing process driver's hard cancellation and crash isolation while removing separate driver executables from module distribution.

Keep existing native-process and ODBC configurations working. Add explicit `driver = "module"`, using `driver_path` for a `.dbmod` or the existing driver directory selection with `<engine>.dbmod`. Reject invalid module metadata, host/platform/toolchain mismatches, symlinks and checksum mismatches before executing module code. Passwords remain on private stdin, never worker arguments or module manifests. CLI policy checks remain before dispatch.

Reuse the archived patched goloader source with provenance/license. The builder prepares a private Go SDK copy if compiler internals are needed; do not mutate the installed SDK. Module entrypoints register their own database/sql driver; shared dbcore performs operations. Root source does not import four vendor drivers. Build dependencies may include them for module packaging.

Optimize the raw format: compact symbol names and relocation records, perform provably safe internal prelinking, preserve required GC/type/stack/init metadata. Do not count gzip compression as elimination of raw overhead. Measure host and each module in bytes and MiB, with raw and compressed values distinguished. Never claim all relocation reservations are removable.

Validation: format round-trip/rejection/size tests; safe prelink tests; manifest/discovery and policy tests; macOS arm64 host+four-module build and load/register/GC smoke; PostgreSQL wire-fixture query if feasible; Windows amd64 cross-build only until a Windows runner is available. Real database compatibility and Windows runtime support must not be claimed without execution. No unloading while Go runtime references survive.
