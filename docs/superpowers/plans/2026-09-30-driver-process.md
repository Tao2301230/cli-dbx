# DBX Driver Process Implementation Plan

> **For agentic workers:** Use subagent-driven-development for independent tasks and review the integrated change before completion.

**Goal:** Keep MySQL and SQLite built in; provide PostgreSQL, Oracle, DM and SQL Server as separately installed driver processes.

**Architecture:** Preserve CLI policy and result APIs. Route built-ins to driver-neutral dbcore and external engines through a versioned, typed stdio protocol. Driver modules register only their own vendor library.

**Tech Stack:** Go, database/sql, encoding/json, os/exec, SHA-256, existing vendor Go libraries.

## Global Constraints

- Baseline: 86b18d3ce5241148f132af211753c3c2de69b071.
- Built-in engines: mysql, sqlite. External engines: postgres, oracle, dm, sqlserver.
- Preserve user untracked files and examples; no main merge, push, tag or publication.
- Protocol version 1; no credentials in process arguments or stderr forwarding.
- Keep transaction and import execution within one child; row limits checked before commit.
- Root module must not depend on external vendor drivers. Oracle toolchain isolated in its module.

### Task 1: Shared execution and vendor SQL

- [x] Move the existing internal/db implementation to internal/dbcore without changing its public function/type names; remove driver blank imports.
- [x] Add a json-excluded `Connector driver.Connector` to conn.Spec and let dbcore.Open use sql.OpenDB when present.
- [x] Add tests first for vendor bind markers, qualified identifier quoting, row-count failure rollback, and SELECT INTO classification.
- [x] Port only required metadata SQL for oracle, dm and sqlserver from the existing main implementation; do not bring ODBC or vendor imports into dbcore.
- [x] Preserve baseline tests and add focused dbcore tests using SQLite or a controlled database/sql test driver.

### Task 2: Connection configuration and independent modules

- [x] Add tests for external engine aliases, structured DSNs and driver_dir precedence, then extend config/conn.
- [x] Create drivers/{postgres,oracle,dm,sqlserver}, each with its own go.mod, go.sum, adapter tests and main.
- [x] Each main calls `driverprocess.Serve(ctx, os.Stdin, os.Stdout, engine, prepare)`; prepare is `func(*conn.Spec) error` and may set Spec.Connector.
- [x] PostgreSQL uses pgx v5.7.2; SQL Server uses microsoft/go-mssqldb v1.8.0; DM uses gitee.com/chunanyong/dm v1.8.22; Oracle uses oracle/go-oracledb/v26 v26.0.1-beta with Go 1.26.6.
- [x] Child modules require github.com/linlay/cli-dbx v0.0.0 with `replace github.com/linlay/cli-dbx => ../..` and the Cobra local replace needed by root source builds.
- [x] Run tests/build separately in each nested module.

### Task 3: Protocol, verified launcher and routing

- [x] Write failing tests for typed cell round trips, bad versions/manifests/hashes, real child query, rollback and cancellation.
- [x] Implement Request/Response with bounded JSON frames, typed query rows and per-step transaction rows.
- [x] Implement `Serve(context.Context, io.Reader, io.Writer, string, func(*conn.Spec) error) error` and `Call(context.Context, conn.Spec, Request) (Response, error)`.
- [x] Implement root resolution and manifest validation, launch with empty credential arguments, enforce timeout and reap children. Do not retry failures.
- [x] Replace internal/db with type aliases and thin operation wrappers; blank import mysql/sqlite only here.
- [x] Adapt structured CLI errors and redact vendor connection passwords; run policy refusal tests before launch.
- [x] Remove pgx and its unused transitive dependencies from root go.mod/go.sum.

### Task 4: Packaging and docs

- [x] Add a standalone driver packaging utility and tests for manifest/hash/path/executable mode and deterministic archives.
- [x] Add driver release build entry points for shell and PowerShell; keep driver archives independent of the core connector package.
- [x] Update README, AGENTS, both help surfaces, connector skills and tracked driver configuration examples to explain installation and compatibility.
- [x] Use experimental release version v0.2.0-driver-process.1.

### Task 5: Integration and review

- [x] Run root go test ./..., protocol race tests, module tests/builds and help smoke tests.
- [x] Use go list -deps ./cmd/dbx to prove four external vendor packages are absent and mysql/sqlite remain.
- [x] Build stripped core and compare to baseline 13,415,010 bytes using the same local toolchain/options.
- [x] Exercise driver archive build and installed-manifest loading; record unavailable external services explicitly.
- [x] Independent review; fix important findings and rerun covering tests.

## Validation record — 2026-09-30

- Root `go test -count=1 ./...`: passed.
- Protocol/core race tests: passed.
- Each independent module `go test ./...`: passed (Oracle adapter tests; other modules compile their entrypoints).
- Four darwin/arm64 archives built with `scripts/build-drivers.sh --host`; archive and binary SHA-256, executable modes, installed directory loading, vendor connection-error path and credential redaction verified.
- Core help/query/tx/version smoke checks passed.
- Root dependency list contains mysql/sqlite and no pgx/oracle/DM/SQL Server driver packages.
- Same stripped local build: baseline 13,415,010 bytes; experimental core 9,673,554 bytes (about 28% smaller).
- Independent review findings addressed: SQL Server batches including nested comments and control commands; DM credential overrides; import row-limit rollback; uncertainty guidance after writes.
- No real vendor database services were available; success against real servers and PowerShell execution remain unverified.
