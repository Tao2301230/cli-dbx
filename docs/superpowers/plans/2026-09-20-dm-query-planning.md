# DM Query Planning Implementation Plan

> 历史讨论记录：其中自定义 Connector、进程协议及 ODBC 自动回退方案不再采用。当前实现改为直接调用厂商 ODBC 驱动，见 [ODBC 说明](../../odbc.md)；本文不作为当前开发指令。

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let DBX inspect DM indexes and obtain a normalized, non-executing DM estimated plan so the Skill can compare candidate SQL before large reads.

**Architecture:** Keep orchestration in the existing CLI and put DM catalog/plan behavior in `internal/db`. `inspect table` gains index metadata. A dedicated `explain` command validates one read-only statement, applies existing connection/table policy, calls DM `EXPLAIN FOR`, and returns normalized plan nodes plus advisory signals; it never runs the target query and never creates indexes.

**Tech Stack:** Go, `database/sql`, DM Go driver, Cobra, `go-sqlmock`, Markdown Skill references.

## Global Constraints

- Implement DM only in this phase; other engines return an explicit unsupported-engine error from `explain`.
- Use `ALL_INDEXES` and `ALL_IND_COLUMNS`, not DBA-only views.
- Accept only one read-only target statement and preserve `allow_actions` / `allow_tables` checks.
- Do not run `AUTOTRACE`, `TRACE`, or any actual-plan command.
- Do not automatically create indexes, add hints, or rewrite joins.
- Keep plan signals advisory; a clustered scan is not automatically a failure.

---

### Task 1: DM index metadata

**Files:**
- Modify: `internal/db/db.go`
- Modify: `internal/db/metadata_extended.go`
- Test: `internal/db/metadata_extended_test.go`

**Interfaces:**
- Produces: `IndexInfo` and `TableInfo.Indexes`.
- Produces: `loadDMIndexes(context.Context, *sql.DB, string, string) ([]IndexInfo, error)`.

- [ ] **Step 1: Write the failing tests**

Add a SQL-mock test that describes a DM table, returns two ordered index-column rows, and asserts one composite `IndexInfo`. Assert permission and iteration errors propagate.

- [ ] **Step 2: Run test to verify it fails**

Run: `GOCACHE=/tmp/cli-dbx-go-build-cache go test ./internal/db -run 'TestDMDescribeIncludesIndexes|TestDMIndexMetadataErrorsPropagate' -count=1`

Expected: FAIL because `TableInfo.Indexes` and DM index loading do not exist.

- [ ] **Step 3: Write minimal implementation**

Add:

```go
type IndexInfo struct {
    Name        string   `json:"name"`
    Columns     []string `json:"columns"`
    Unique      bool     `json:"unique"`
    Type        string   `json:"type,omitempty"`
    Partitioned bool     `json:"partitioned,omitempty"`
    Status      string   `json:"status,omitempty"`
}
```

Load visible DM indexes from `ALL_INDEXES` joined to `ALL_IND_COLUMNS`, scoped by both owner and table name, ordered by index name and column position. Call it only for `spec.Engine == "dm"` after constraints are loaded.

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/cli-dbx-go-build-cache go test ./internal/db -run 'TestDMDescribeIncludesIndexes|TestDMIndexMetadataErrorsPropagate' -count=1`

Expected: PASS.

### Task 2: DM estimated-plan adapter

**Files:**
- Create: `internal/db/explain.go`
- Create: `internal/db/explain_test.go`

**Interfaces:**
- Produces: `Explain(context.Context, conn.Spec, string) (ExplainPlan, error)`.
- Produces: `ExplainPlan`, `PlanNode`, and `PlanSignal` JSON models.

- [ ] **Step 1: Write the failing tests**

Use SQL mock to expect `EXPLAIN FOR SELECT ...`, return representative `CSCN2`, `HASH2 INNER JOIN`, `SSEK2`, and `BLKUP2` rows, and assert normalized table/index/row/cost/partition fields and advisory signals. Add unsupported-engine and row-iteration error tests.

- [ ] **Step 2: Run test to verify it fails**

Run: `GOCACHE=/tmp/cli-dbx-go-build-cache go test ./internal/db -run 'TestExplainDM|TestExplainRejectsUnsupportedEngine' -count=1`

Expected: FAIL because `Explain` does not exist.

- [ ] **Step 3: Write minimal implementation**

Execute `EXPLAIN FOR ` plus the validated target SQL with the connection timeout. Scan result columns by normalized column name so nullable values and compatible DM driver value types are handled. Emit signals for clustered scans, index seeks, lookups, joins, sorts, Cartesian operations, partition ranges, and non-empty `ADVICE_INFO`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/cli-dbx-go-build-cache go test ./internal/db -run 'TestExplainDM|TestExplainRejectsUnsupportedEngine' -count=1`

Expected: PASS.

### Task 3: CLI explain command and inspect output

**Files:**
- Modify: `internal/cli/app.go`
- Modify: `internal/cli/help.go`
- Modify: `internal/app/app.go`
- Modify: `internal/app/commands.go`
- Test: `internal/cli/app_test.go`
- Test: `internal/app/app_test.go`

**Interfaces:**
- Produces: `dbx explain <conn> <sql>` and `dbx explain file <conn> <path.sql>`.
- Extends `inspect table` JSON with `indexes`.

- [ ] **Step 1: Write failing command/help tests**

Assert root help lists `explain`, explain help documents estimated/non-executing DM behavior, malformed invocations fail with usage, and the legacy command rejects non-DM engines without executing SQL.

- [ ] **Step 2: Run test to verify it fails**

Run: `GOCACHE=/tmp/cli-dbx-go-build-cache go test ./internal/app ./internal/cli -run 'Explain|InspectTableIncludesIndexes' -count=1`

Expected: FAIL because the command is missing.

- [ ] **Step 3: Implement command and policy path**

Resolve the connection with query permission, parse the target SQL, require exactly one read-only statement, enforce existing table allow-list policy, skip bounded-result enforcement because DM does not execute the target SQL, call `db.Explain`, and return `kind=explain_plan` with normalized nodes and signals.

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/cli-dbx-go-build-cache go test ./internal/app ./internal/cli -run 'Explain|InspectTableIncludesIndexes' -count=1`

Expected: PASS.

### Task 4: Skill, mirrors, and verification

**Files:**
- Modify: `README.md`
- Modify: `docs/beginner-guide.md`
- Modify: platform Skill `SKILL.md` and `references/commands.md`, `references/output-and-behavior.md`, `references/safety-and-troubleshooting.md`
- Sync: `agent-platform-builtins/dbx`
- Sync: installed `~/.zenmind/skills-center/builtin-dbx`

**Interfaces:**
- Produces: DM workflow `inspect table -> explain candidate SQL -> compare signals -> bounded preview/bulk strategy`.

- [ ] **Step 1: Update documentation and Skill**

Document DM operations `CSCN2`, `SSEK2/CSEK2`, `BLKUP2`, `HASH JOIN`, `NEST LOOP`, `SORT`, `PSTART/PSTOP`, and `ADVICE_INFO`. State that fewer joins are not automatically faster and that DBX does not apply hints or indexes.

- [ ] **Step 2: Sync mirrors**

Copy CLI source to `agent-platform-builtins/dbx` and platform Skill source to the installed Skill directory.

- [ ] **Step 3: Run full verification**

Run CLI and mirror `go test ./...`, platform connector/builtin tests, `git diff --check`, mirror diffs, and DM-related cross builds for macOS/Windows amd64/arm64.

Expected: all pass and both mirrors are identical.
