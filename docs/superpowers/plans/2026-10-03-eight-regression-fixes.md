# Eight Regression Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复已确认的 F01–F08，并保留可复查的回归证据。

**Architecture:** 保持 CLI 控制面、共享执行层和独立驱动协议。SQL 词法分析统一识别字面量与层级；数据列归一在 dbcore；CLI负责参数预检、数值读取、脱敏及执行通路选择。

**Tech Stack:** Go 1.22 根模块、标准库、现有 SQLite 与私有驱动进程。

## Global Constraints

- 规格：`docs/superpowers/specs/2026-10-03-eight-regression-fixes-design.md`。
- 仓库：`/Users/bantao/Documents/pworkspace/deploy_dev/cli-dbx`；不改变现有 Git index，不提交、发布或下载新依赖。
- 本轮日志：`/private/tmp/cli-dbx-fix-20261003-d44dd2tb`；缓存 `GOCACHE=/private/tmp/cli-dbx-review-gocache`；离线 `GOPROXY=off`。
- 每个任务先写并运行失败测试；只编辑所属文件，采用现有测试辅助函数或独立命名，避免共享文件冲突。
- 原始手册及失败报告保留；文档更新由主代理统一负责。

---

### Task 1: SQL safety F01–F03

**Files:** `internal/sqlclass/`、`internal/sqlanalyzer/`；必要时新增 `internal/sqllex/`；CLI端到端测试新建 `internal/cli/sql_safety_regression_test.go`。

**Interfaces:** 保持 `Classify`、`HasUnsafeWrite`、`Analyze`、`AnalyzeForEngine` 的签名及 Analysis 字段。CLI通过现有 `enforcePolicy` 消费分析结果；修改CLI实现需先通知主代理。

- [ ] 重现写入CTE、EXPLAIN ANALYZE DELETE、设置PRAGMA、tx隐藏写入、内层WHERE、字符串WHERE、字符串注释标记及USING来源表的失败；补正常只读CTE和引号/注释边界。
- [ ] 运行 `GOCACHE=/private/tmp/cli-dbx-review-gocache GOPROXY=off go test ./internal/sqlclass ./internal/sqlanalyzer ./internal/cli -run 'SafetyRegression|Classify|Analyze|Unsafe' -count=1`，记录 RED 日志。
- [ ] 用共享词法分词和层级检查实现规格行为；不能只针对已知SQL字符串写黑名单。复杂写入CTE需要保守拒绝，避免按外层SELECT的行数掩盖真实影响行数。
- [ ] 运行相关测试和已有 SQLServer 方言测试，记录 GREEN；写 `results/sql-safety.md`，含改动、测试命令及限制。

### Task 2: Lossless duplicate columns F05

**Files:** `internal/dbcore/db.go`、`internal/dbcore/db_test.go` 或专用新测试文件；如需要进程端到端，新增测试文件而不编辑其他任务文件。

**Interfaces:** 保持 `QueryResult`、`Column` 和进程 v1。在构造Columns时选择唯一名字，rows及stats统一使用该名字。重复列后缀不得占用查询原本存在的唯一名字。

- [ ] SQLite真实查询 `SELECT 1 AS id, 2 AS id` 必须保留两个值；再覆盖 `id/id/id_2`、重复三列和无重复查询。
- [ ] 运行 `GOCACHE=/private/tmp/cli-dbx-review-gocache GOPROXY=off go test ./internal/dbcore -count=1`，记录 RED。
- [ ] 在共享结果列归一处实现确定的名字选择；不能只修JSON渲染，不能忽略stats或Column顺序。
- [ ] 验证普通查询、分页和事务查询未回归，记录 GREEN 和 `results/duplicate-columns.md`。

### Task 3: CLI F04/F06/F07/F08

**Files:** `internal/cli/app.go`、专用 `internal/cli/fix_regression_test.go`；必要时 `internal/output/output.go` 及其独立测试。

**Interfaces:** 使用现有 `db.Execute` 非事务通路执行 admin，其他写入继续 `db.ExecuteGuarded`；整数最终以进程支持的精确标量传给db层。SQL安全任务不修改本实现文件。

- [ ] 添加假SQLServer改密参数脱敏、JSON大整数导入、SQLite VACUUM，以及非法format无写入副作用的失败测试。
- [ ] 运行 `GOCACHE=/private/tmp/cli-dbx-review-gocache GOPROXY=off go test ./internal/cli -run 'FixRegression' -count=1`，记录 RED。
- [ ] JSON用UseNumber读取并做精确转换，保持原json.Unmarshal对完整文档的校验；补边界int64、普通小数和非法输入。
- [ ] DSN与错误输出统一保护已支持的密码参数；使用假密码覆盖URL编码和大小写，不增加真实凭据访问。
- [ ] admin走现有非事务执行；格式校验在任何DB/输出文件副作用前完成，分别支持JSON/table结果和CSV/JSON导出文件格式。
- [ ] 运行相关测试和原CLI全套，记录 GREEN 和 `results/cli-fixes.md`。

### Task 4: Integration, review and documentation

**Files:** 主代理维护 `README.md`、相关help、修复验证报告、手动回归补充及新验证材料；不得改原始失败证据。

- [ ] 审读三份任务报告及相对于本轮baseline的diff，复查词法分支、精度边界及admin/格式副作用。
- [ ] 运行根 `go test -count=1 ./...`、`go test -race -count=1 ./...`、`go vet ./...`；分别测试全部嵌套驱动模块。
- [ ] 构建新CLI，并用隔离数据重跑原正常smoke与八类失败的修复预期。另检查进程结果列和JSON整数传输。
- [ ] 对本轮整体diff做独立代码审查；发现相关问题则补测试、修复并重新验证。
- [ ] 新建修复报告及手动回归补充，保留原手册；校验原手册SHA和Git index不变；最终报告测试范围和未验证环境。
