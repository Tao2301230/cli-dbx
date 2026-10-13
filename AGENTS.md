# AGENTS.md

## 1. 项目概览

`dbx` 是一个带安全边界的 database CLI，不是任意 SQL 壳。它的目标是让人类和智能体都能在统一命令面下访问数据库，同时把高风险动作限制在 DBX 自身的控制面里。

当前核心能力：

- 连接管理：`conn list / show / test`
- 结构查看：`inspect schema / table / connection`
- 显式 SQL 动作：`query / update / schema / admin`
- 文件导入导出：`import` / `export`
- 批处理事务：`tx`

设计重点：

- 显式动作优先，减少大模型误用
- DBX 层能力白名单，独立于数据库账号权限
- 单语句默认、事务批处理明确、输出结构化

## 2. 技术栈

- 语言：Go
- CLI：标准库 `flag`
- 数据库访问：`database/sql`
- 驱动：
  - 主程序内置 MySQL：`go-sql-driver/mysql`，SQLite：`modernc.org/sqlite`
  - 独立进程：PostgreSQL、Oracle、达梦、SQL Server，各有 `drivers/<engine>/go.mod`
  - 此共享宿主开发分支统一固定 Go 1.26.8，宿主、模块和公共依赖必须一致
- 配置格式：TOML
- 构建与测试：`go build`、`go test`
- 发布：仓库内 `scripts/release/build.sh`

## 3. 架构设计

模块分层大致如下：

- `cmd/dbx`
  CLI 入口，仅负责进程启动和顶层错误退出。
- `internal/cli`
  命令分发、参数解析、帮助信息、错误 envelope。
- `internal/config`
  单连接配置文件加载、路径归一化、密码来源解析。
- `internal/textfile`
  标准库 UTF-8 文件读取，移除文件开头至多一个 BOM；正文中的 BOM 保持原样。
- `internal/secret`
  v1/v2 密文处理、系统凭据库适配和运行时密钥解析。
- `internal/conn`
  把配置解析成可执行的连接规格 `Spec`。
- `internal/action`
  动作模型：`query / update / schema / admin`。
- `internal/sqlclass` / `internal/sqlanalyzer`
  SQL 分类、单语句校验、危险写入检测。
- `internal/sqllex`
  共享引号、注释和括号层级识别；方言语法不确定时拒绝执行。
- `internal/db`
  路由 MySQL/SQLite 到内置执行器，四种外部引擎到进程客户端。
- `internal/dbcore`
  共享 SQL 执行、元数据、导入与整批事务，不导入厂商驱动。
- `internal/driverprocess`
  驱动包校验、私有管道协议、进程取消与回收。
- `internal/drivermodule`、`internal/dbmodformat`
  可选 Go 模块的发现、校验、动态装载及 DBXMOD2 文件格式；`shared.go` 保留宿主公共实现。
- `modules/postgres`、`modules/oracle`、`modules/dm`、`modules/sqlserver`
  `.dbmod` 的注册入口，共用根 go.mod 的工具链与依赖；区别于下述独立可执行驱动。
- `cmd/dbmod-pack`、`scripts/build-modules.py`、`third_party/goloader`
  离线模块打包、规范构建编排及带许可证和修补记录的 Go 对象加载器。
- `drivers/postgres`、`drivers/oracle`、`drivers/dm`、`drivers/sqlserver`
  独立模块及可执行入口，各自注册对应数据库驱动。
- `drivers/odbc`
  独立 ODBC 模块，按构建时固定的 oracle/dm/sqlserver 方言分别生成适配程序；只有此模块链接 Go ODBC 适配层。
- `internal/output`
  结构化输出、摘要、表格输出。

核心调用链：

1. CLI 解析命令与参数
2. 连接配置解析为 `conn.Spec`
3. SQL 被分类为 statement class，再映射到 action
4. DBX 先做 `allow_actions` 和安全策略校验
5. 校验通过后才进入数据库执行层；外部引擎此时才校验并启动独立驱动

## 4. 目录结构

- `cmd/dbx`
  可执行程序入口
- `internal/`
  核心实现
- `docs/`
  面向用户的补充文档
- `testdata/`
  配置示例
- `scripts/release/`
  打包与发布脚本

文档职责固定为：

- `README.md`
  操作手册、快速开始、简单验证
- `docs/beginner-guide.md`
  新手第一次上手
- `AGENTS.md`
  设计、重要功能、开发约定、维护流程

## 5. 数据结构

### 连接配置

根结构是单连接文件：

```toml
[connection]
engine = "mysql|sqlite|postgres|oracle|dm|sqlserver"
allow_actions = ["query", "update", "schema", "admin"]
```

`allow_actions` 是 DBX 层保护，不依赖数据库账号本身的授权能力。

### 密码加密

- `dbx-aes-gcm:v2:<key-id>:...` 是默认格式。数据库密码保存在配置密文中，每条密文的随机 AES-256-GCM 密钥保存在当前用户的系统凭据库。
- 系统凭据库固定使用 macOS Keychain、Windows Credential Manager 或 Linux Secret Service，不允许降级到环境变量、命令或密钥文件。
- v2 运行时静默解密；密钥缺失、凭据库不可用或密文损坏都必须失败关闭。
- `dbx-aes-gcm:v1:...` 仅为兼容保留，读取时仍由 passphrase provider 获取旧 master passphrase。
- CLI 只提供 `dbx secret encrypt`，不得增加 decrypt、export、show-key 等明文或密钥输出入口。
- `dbx secret encrypt` 接受零个或一个位置参数：零参数隐藏回显地读取密码，一个参数直接加密且不得读取 stdin。位置参数必须原样处理，不得 trim 或写入 DBX 输出与错误。
- 位置参数模式只用于接受一次明文暴露的迁移场景；明文会进入调用智能体、Shell history、命令审计或进程参数，不能宣称该模式对模型保密。
- v2 的保护目标是避免配置文件直接泄露明文，不抵御同一操作系统用户主动读取系统凭据库或进程内存。

### 动作模型

- `query`
  只读 SQL
- `update`
  DML，如 `insert / update / delete / merge`
- `schema`
  DDL，如 `create / alter / drop / truncate / rename`
- `admin`
  管理动作，如 `grant / revoke / set / vacuum`

### 事务计划

`tx` 读取 JSON 计划文件：

```json
{
  "steps": [
    {"action":"query","sql":"select id from users where id = 1"},
    {"action":"update","sql":"update users set active = 1 where id = 1","max_rows_affected":1}
  ]
}
```

`tx` v1 只允许 `query` 和 `update` 步骤。

计划读取使用 `DisallowUnknownFields` 并检查第二次解码必须返回 `io.EOF`，在数据库执行前拒绝第二个 JSON 文档和尾随垃圾。JSON 导入保留 `UseNumber`、精确数字绑定和 EOF 检查；不得为了迁移文件编码修复而退回 `json.Unmarshal` 到浮点数字。

### 文件编码与 SQLite 路径

配置、SQL、事务 JSON、CSV 和 JSON 导入使用 `internal/textfile` 读取，支持带或不带 UTF-8 BOM，不自动转换 UTF-16。Windows 测试同时隔离 `HOME` 和 `USERPROFILE`，TOML 测试路径使用安全引用。

SQLite 普通文件路径编码为 `file:` URI 时，只编码文件名部分并保留原始查询参数。显式 URI 不重复编码；保留 `mode=ro`、内存库、共享缓存与 `_pragma` 行为。打包器的 `--executable` 以明确的暂存区相对普通文件路径指定可执行权限，不依赖构建主机的文件模式；无此参数的既有行为保持不变。

### 输出 Envelope

对外返回统一 envelope，关键字段包括：

- `kind`
- `action`
- `class`
- `summary`
- `data`
- `next`
- `warnings`

`action` 是稳定的外部语义；`statement_class` 仍然保留作兼容信息。

## 6. API 定义

`dbx` 的主要外部接口是 CLI，而不是 HTTP API。

核心命令面：

- `dbx conn ...`
- `dbx inspect ...`
- `dbx query <conn> '<sql>'`
- `dbx update <conn> '<sql>'`
- `dbx schema <conn> '<sql>'`
- `dbx admin <conn> '<sql>'`
- `dbx tx <conn> --plan <path.json>`
- `dbx import file ...`
- `dbx export table ...`
- `dbx secret encrypt`

接口约定：

- 默认只允许单语句 SQL
- `query/update/schema/admin` 要求 SQL 类型和命令动作匹配
- `tx` 为单连接、单次调用、单事务

错误语义重点：

- `action_blocked`
  当前连接未放行该动作
- `action_mismatch`
  命令动作和 SQL 类型不一致
- `multiple_statements_blocked`
  一次传入了多语句
- `missing_where`
  危险写入缺少 `WHERE`
- `tx_unsupported_action`
  事务步骤使用了 `schema` 或 `admin`

## 7. 开发要点

- 新增数据库能力时，先明确它属于哪个 action，再补 CLI 和帮助文案。
- 不要把设计解释重新堆进 `README.md`；设计说明统一维护在 `AGENTS.md`。
- `README.md` 只保留操作说明、配置示例、简单测试和排查。
- `docs/beginner-guide.md` 只讲新手视角，不承载架构事实。
- `allow_actions` 是 DBX 层最重要的能力边界之一，任何变更都必须补测试。
- `tx` 的边界优先保持简单和可验证，不要直接扩展到跨进程会话事务。
- 对外词汇优先使用 `query / update / schema / admin / tx`，保持动作边界清晰。
- 用户可见能力变更需要同步更新：
  - CLI help
  - `README.md`
  - 相关测试

## 8. 开发流程

此分支的完整测试需要私有 Go 1.26.8 SDK；在新机器上先按文末「跨机器开发交接」构建并设置 SDK，再执行以下命令。仅安装原版 Go 后直接运行 `go test ./...`，会因缺少 `cmd/objfile` 内部包而失败。

本地开发常用命令：

```bash
go build -o ./dbx ./cmd/dbx
go test ./...
(cd drivers/postgres && go test ./...)
(cd drivers/oracle && GOTOOLCHAIN=go1.26.8 go test ./...)
(cd drivers/dm && go test ./...)
(cd drivers/sqlserver && go test ./...)
```

修改命令或输出时，至少检查：

```bash
./dbx --help
./dbx query --help
./dbx tx --help
```

发布流程当前由维护者手工执行：

1. 将根目录 Git 跟踪的 `VERSION` 更新为计划发布的版本（格式如 `v0.1.0`），并与代码和文档一并提交
2. 运行 `go test ./...`
3. 创建并推送与 `VERSION` 相同的 tag，例如 `v0.1.0`
4. 执行 `scripts/release/build.sh` 构建核心包，另执行 `scripts/build-drivers.sh` 构建四种独立驱动包；本机验证可用 `--host --engine postgres`
5. 校验 `dist/<version>/..._checksums.txt` 与各独立驱动包的 `.sha256`
6. 在 GitHub Release 上传产物

如果文档分工变化：

- 设计和约定改 `AGENTS.md`
- 用户操作和快速验证改 `README.md`
- 新手入门体验改 `docs/beginner-guide.md`

## 9. 已知约束与注意事项

- `dbx` 不是数据库权限系统的替代品；它是在数据库账号之外再加一层 DBX 侧保护。
- `allow_actions` 只能收紧 DBX 的能力边界，不能放大数据库账号本身没有的权限。
- `tx` 不支持跨连接事务，也不支持跨多次 CLI 调用的事务会话。
- `tx` 不支持 `schema` 或 `admin`，主要是为了保持回滚语义可预期。
- 生产态连接默认更保守；带 `prod` 特征的连接会阻止高风险动作配置。
- v2 密文绑定创建它的机器和操作系统用户，复制配置到其他机器后必须重新加密。
- Linux v2 依赖可用且已解锁的 Secret Service 登录集合；不可用时不回退到其他密码来源。


## Platform 连接器发布包

项目位于 `agent-platform-connectors/dbx`，`connector/` 是 connector.json 模板、cli.json、skills 与全部资源的唯一源码。`VERSION` 同时控制 CLI 和连接器发布版本；模板不维护 version，构建时生成。任何 CLI、清单或技能改动均需发布新版本。

现有 Shell/PowerShell 发布脚本额外生成 `dist/<version>/builtin.dbx_<version>_<os>_<arch>.zip`，内容位于 `runtime/connectors/builtin.dbx/`，包含声明、完整 skills 和目标平台 `bin/dbx`。核心连接器只随主程序提供 MySQL、SQLite，四种外部驱动单独分发；不得向 `connector/` 添加驱动二进制。Platform 按完整包摘要消费，不修改包内内容。原 CLI 独立发布包继续提供。包生成器检查 VERSION 满足 cli.json.minVersion（包含 SemVer 预发布版本比较）。


## 独立驱动进程契约

- `connection.driver_path` 精确指定驱动可执行文件，支持 `~`，相对路径以连接 TOML 所在目录为准。包内平铺 `dbx-driver-<engine>-<version-without-v>[.exe]` 与完整程序文件名加 `.json` 的清单，多个版本互不覆盖。未指定路径时按 `driver_dir` > `DBX_DRIVER_DIR` > `~/.config/dbx/drivers` 选定根目录，只找平铺的无版本 `dbx-driver-<engine>[.exe]`。不搜索版本、不从 PATH 搜索、不自动下载或回退。旧的引擎子目录布局需重新安装。
- `<可执行文件名>.json` 必须包含 `protocol_version`（当前 1）、`engine`、`version`、`os`、`arch`、`executable` 和小写十六进制 `sha256`。清单内程序名必须与选定文件名完全一致，程序和清单必须为普通文件，拒绝目录逃逸和符号链接；启动前检查引擎、平台、协议及摘要。驱动包是当前用户权限下的可信代码，不是安全沙箱。
- 每次数据库操作使用一个进程；JSON 帧有长度限制，stdout 仅传协议，数据行分帧发送。标量显式保留 NULL、整数、浮点、字符串、二进制和时间，避免 JSON 浮点精度或时间格式变化。
- 凭据和 DSN 仅走私有 stdin，不能加入子进程参数、日志或直接转发 stderr。安全策略、密钥解析、审计和对外 envelope 仍在 DBX 主程序中。
- 事务全部步骤和整批导入必须在同一驱动进程内完成；影响行数在 commit 前校验，无法取得行数时不能假定为零。不要引入跨 CLI 调用的事务会话。
- 父进程取消或超时终止并回收子进程，子进程同时使用操作超时及管道关闭监听。写操作不自动重试；提交期间失联不能承诺已回滚，须提示查询确认结果。
- 根模块不得传递链接 pgx、Oracle、DM 或 SQL Server 的厂商客户端库。四个嵌套模块分别运行测试；根 `go test ./...` 不覆盖它们。
- 打包工具 `scripts/release/package-driver/` 仅依赖标准库。Shell/PowerShell 入口生成独立归档及 `.sha256`，固定目录、时间、权限和次序。Oracle 工具链只在其模块构建时覆盖，不运行全局 `go env -w`。
- 自动测试覆盖子进程执行、类型保真、清单/摘要/协议错误、取消回收、失败与超限回滚、策略拒绝时不启动驱动；真实厂商数据库未连接时，不能用编译或本地替身测试宣称数据库兼容性已验证。

用户安装与操作说明见 [docs/driver-process.md](./docs/driver-process.md)。


## ODBC 独立进程

- 连接 `driver = "odbc"` 显式选择 ODBC；空值或 `native` 保持 Go 路线。当前 ODBC 仅开放 Oracle、达梦、SQL Server。
- `driver_path` 指向 DBX 的版本化 ODBC 适配程序；`odbc_driver` 指向另外安装的厂商库，两者支持 `~` 和 TOML 相对路径。Windows 按当前进程架构查找指向该 DLL 的已注册名称，不写注册表。
- ODBC 使用结构化字段和 `odbc_options`，禁止原始 DSN 和凭据/驱动/结构化目标字段覆盖。主程序解析密码后通过私有协议传给子进程，子进程组装 ODBC 连接串并加载库；输出不得序列化私有 ODBC 配置。
- `drivers/odbc/fork` 保留 ODBC 上游许可证和来源，维护数据类型及影响行数修补；它是嵌套模块，需单独运行测试。NUMERIC/DECIMAL 以精确文本传输，二进制与文本按 SQL 类型区分。
- ODBC 参数标记为 `?`，SQL 引号和元数据仍跟随 `engine`。事务、导入和影响行数保护与原有进程契约相同。厂商错误信息保守脱敏，同时保留不确定写入状态的提示。
- 清单 `adapter` 缺省或 `native` 代表 Go 路线，`odbc` 必须与连接的驱动选择一致。可选 `files` 记录同目录配套文件的 SHA-256；加载时验证普通文件、名称无目录跳转和摘要。资源不得重名覆盖程序或清单。
- `scripts/build-odbc-drivers.sh` / `.ps1` 单独构建 ODBC 变体。Mac 在匹配架构原生构建，复制 unixODBC/libltdl 并重定位、签名临时副本。配套库和许可证平铺且带发布版本，三个引擎同版本复用相同内容。主程序包不携带 ODBC 运行库。
- Linux ODBC 需本机 CGO/unixODBC；Windows 仅声明 x64 支持，使用系统管理器。厂商客户端不纳入 DBX 包，其资源和硬编码依赖沿用厂商安装要求。
- 所有 Go/ODBC 包保持不同文件名和清单；禁止因连接失败自动切换路线。协议暂保持 v1，新增连接字段向后兼容原有 Go 适配器。

安装和构建说明见 [ODBC 独立进程](./docs/odbc-driver-process.md)。

## 共享宿主模块（开发分支）

- `driver = "module"` 选择四种 Go 厂商驱动的 `.dbmod` 模块。`Spec.Module` 记录传输选择，`Spec.Driver` 保持 database/sql 的实际驱动名。原 native 进程和 ODBC 路线保持显式选择，不自动回退。
- `scripts/build-modules.py` 以 `-tags=dbmod -trimpath -buildvcs=false -ldflags=-w` 构建共享宿主及独立模块。禁止 `-s`，宿主符号表是离线打包和运行时链接的依赖。普通构建不包含模块加载器。
- 主程序源码不导入四个厂商驱动；`modules/<engine>` 是模块注册入口，Oracle 另导出连接器准备函数。构建依赖统一在根 go.mod 固定，原独立驱动的嵌套模块仍可构建。
- `internal/drivermodule/shared.go` 显式保留模块重复使用的通用证书、文本、Unicode、数值、JSON 和容器实现，供打包器引用宿主符号。类型注册在 `goloader.RegSymbol` 之后执行，不调用被保留的函数；禁止借此向宿主引入厂商驱动。迁移同时比较宿主增长、模块减少和整体未压缩净大小。
- 模块以 `<engine>.dbmod` 与 `<engine>.dbmod.json` 分发，按 driver_path 或已有目录优先级发现。清单固定 schema、Go 版本、平台、引擎、模块名及模块/宿主 SHA-256；父进程和 worker 都先验证。模块是当前用户信任的本机代码，不是沙箱。
- 父进程启动自身 `dbx --internal-module-worker`，通过原私有协议发送凭据和操作。四个模块复用同一份磁盘宿主，worker 的运行时内存各自独立。取消、回收、整批事务/导入和不重试语义沿用 driverprocess。
- worker 在同一 OS 线程完成 Go 模块重定位与发布；整个操作保留 CodeModule、类型及 GC 元数据，不支持主动卸载。厂商初始化不得解析 worker 命令行。
- DBXMOD2 使用紧凑编号、前缀字典、变长整数和稀疏 NOP 预留区描述；磁盘省略的 NOP 在内存中复原，不能宣称等量减少运行时内存。安全 prelink 只移除地址基址无关的重定位；保留 stackobject、deferreturn、接口标记和必要元数据。
- 体积报告分别记录未压缩模块和 gzip 对照。真实数据库兼容、TLS、Windows 执行须有实际验证证据，交叉构建不替代运行验证。
- 第三方加载器来源及本地修补记录在 `third_party/goloader/DBX-PATCHES.md`，构建只准备私有 SDK，不修改全局 Go 安装。

### 跨机器开发交接（2026-10-13）

**分支与用户目标。** 当前开发位于 `codex/dbx-shared-host-modules`，从 `b5b7c1b` 开始。目标是一份约 12 MB 级别的共享宿主，加上四个按需复制、按需加载的数据库模块，减少独立 Go 可执行驱动重复携带基础代码的磁盘开销。验收以未压缩文件大小为主；同时报告宿主增长、单模块减少和整体净变化。用户使用多个 terminal 分别启动 CLI，不要求共享常驻服务。当前实现仍启动自身的内部 worker，尚未改成 CLI 主进程直接执行数据库操作；不要将「不需要业务并发」解释为可以删除 Go GC、调度器、类型或栈元数据。

**已实现的内容与入口。**

- `internal/conn/resolve.go` 和 `internal/driverprocess/client.go`：解析 `driver = "module"`，在现有策略检查之后校验模块并启动同一个宿主文件；原 native/ODBC 路线保留。
- `cmd/dbx/module_worker.go`：加载所选模块并调用原 `driverprocess.Serve`；Oracle 初始化前标记全局 flag 已解析，避免厂商代码解析私有 worker 参数。`module_disabled.go` 保留未启用模块的普通构建。
- `internal/drivermodule/manifest.go`：绑定引擎、平台、Go 版本与确切宿主/模块 SHA-256；读取有大小限制的不可变快照，校验和解码使用同一份字节，不能退回到校验后重读原文件。
- `internal/drivermodule/load.go`：先注册宿主符号，再调用 `registerSharedHostTypes`；在锁定的 OS 线程上链接模块，整个操作保留加载句柄。Oracle 额外提供 `Prepare(*conn.Spec) error`，其他入口及厂商实现见 `modules/<engine>/main.go`。
- `internal/drivermodule/shared.go`：显式保留通用证书、文本、Unicode、数值、JSON 和容器函数，让四模块引用宿主实现；函数值只用于保留和注册类型，不执行证书、网络或用户查询。
- `internal/dbmodformat/`：DBXMOD2 前缀字典、变长整数、稀疏 NOP 与 PC 函数名称引用；解码逐字节恢复原布局，保留 GC、类型和栈信息。旧 DBXMOD1 只保留测试夹具，生产加载器不兼容旧格式。
- `third_party/goloader/link/prelink.go`：只提前解析可证明与加载基址无关的重定位，不移动代码，不删除 stackobject、deferreturn 或接口标记。修改此处须同时检查运行时 `relocate.go` 的公式。
- `scripts/build-modules.py` 和 `cmd/dbmod-pack`：先生成宿主，再按其符号闭包打包四模块。`--host` 必须与当前源码、SDK、依赖和参数重建的规范宿主字节一致；只匹配 Go 版本和平台不足以保证 ABI。低层 packer 是内部工具，不能当作任意宿主的插件转换器。

**新机器构建。** 克隆代码并切换到本分支后，在仓库根目录执行。需要原版 Go **1.26.8**、Python 3；macOS arm64 还需要 Xcode Command Line Tools。首次运行不要加 `--offline`，先让 Go 获取 go.mod/go.sum 固定的依赖；缓存齐全后才使用该选项。无需下载本次开发机的实验目录。

macOS/Linux shell（将 Go 路径替换成本机原版 SDK）：

```bash
python3 scripts/build-modules.py \
  --go /path/to/go1.26.8/bin/go \
  --out dist/modules
export GOROOT="$PWD/dist/modules/_build/go1.26.8"
export PATH="$GOROOT/bin:$PATH"
export GOTOOLCHAIN=local
```

Windows PowerShell 的原生构建入口如下，要求 PATH 中的 `go.exe` 已是 1.26.8；这套模块原生构建与执行流程仍待 Windows 机器验证：

```powershell
python scripts/build-modules.py --go (Get-Command go).Source --out dist/modules
$env:GOROOT = Join-Path $PWD "dist/modules/_build/go1.26.8"
$env:PATH = "$env:GOROOT\bin;$env:PATH"
$env:GOTOOLCHAIN = "local"
```

上述 SDK 路径适用于从原版 SDK 首次构建：脚本复制 SDK，并将 `src/cmd/internal` 复制到私有 SDK 的 `src/cmd/objfile`。若输入本身已是准备好的私有 SDK，脚本会直接复用，应将 GOROOT 指向实际输入 SDK。不要修改全局 SDK，不要运行 `go env -w`。完整模块打包要求在目标 OS/架构原生执行；`--goos windows --goarch amd64 --host-only` 只能交叉构建宿主，不能据此宣称四模块已可在 Windows 使用。

**接手后的验证。** 设置好上述私有 SDK 后运行：

```bash
go test ./...
go test -tags=dbmod ./...
(cd drivers/postgres && go test ./...)
(cd drivers/oracle && go test ./...)
(cd drivers/dm && go test ./...)
(cd drivers/sqlserver && go test ./...)
(cd third_party/goloader && go test ./...)
DBX_MODULE_HOST="$PWD/dist/modules/dbx" \
  python3 -B -m unittest discover -s scripts/tests -v
```

Windows 的 Python 测试先设置 `$env:DBX_MODULE_HOST = Join-Path $PWD "dist/modules/dbx.exe"`，再执行 `python -B -m unittest discover -s scripts/tests -v`；各嵌套模块测试在其目录执行。Python 全套当前共 8 项，包含规范宿主拒绝、公共符号保留/厂商隔离、四模块加载及 PostgreSQL 本地协议查询。PostgreSQL 替身需要监听本机 `127.0.0.1`，沙箱限制监听时应申请相应测试权限，不能把限制当作数据库失败。省略 `DBX_MODULE_HOST` 会跳过产物验证，不能把跳过称为四模块通过。额外的旧完整图测试依赖开发机历史归档，其他机器可以跳过；应以新构建产物的四模块加载检查和本地协议测试作为可复现的接手验证。

**当前大小基线。** macOS arm64、Go 1.26.8 正式构建的原始字节数如下，五文件合计 **26,966,874 bytes**（26.97 MB，25.72 MiB）；四份清单另占 1,226 bytes。MB 是 1,000,000 bytes，MiB 是 1,048,576 bytes。

| 文件 | 未压缩 bytes |
| --- | ---: |
| `dbx` | 12,981,506 |
| `postgres.dbmod` | 5,377,907 |
| `oracle.dbmod` | 2,797,165 |
| `dm.dbmod` | 4,017,077 |
| `sqlserver.dbmod` | 1,793,219 |

公共实现迁移使四模块合计减少 2,550,747 bytes、宿主增加 911,280 bytes，整体净减少 1,639,467 bytes。这不是已达体积下限或已达到内嵌增量的承诺。每轮优化重新构建并看 `dist/modules/sizes.json`；平台变化后必须重新测量，不能沿用本表作为 Windows/Linux 大小。

**已验证与后续工作。** 已通过 macOS arm64 四模块加载/注册/连接器构造/关闭/GC、CLI PostgreSQL 本地替身的精确整数/小数/中文传输及策略检查、根与四嵌套驱动测试、加载器测试和 8 项 Python 测试。Windows amd64 当前宿主及原 ODBC 适配器只通过交叉构建。尚未验证真实四库认证、TLS、真实事务/导入、Windows 模块原生执行、其他平台执行或运行时内存；macOS 原生 ODBC 测试还需要 `sql.h`/unixODBC。继续开发优先在有相应环境的机器补足这些证据，再逐项评估剩余重复符号/元数据的可共享性。

**协作资料与产物。** 安装与构建说明见 [共享宿主模块](./docs/shared-host-modules.md)，完整大小、摘要与验证边界见 [验证记录](./docs/shared-host-modules-validation.json)，设计与已完成计划见 [设计](./docs/superpowers/specs/2026-10-13-shared-host-modules-design.md) 和 [计划](./docs/superpowers/plans/2026-10-13-shared-host-modules.md)。打包只分发宿主、所需 `.dbmod` 及对应清单；`dist/`、`_build/`、私有 SDK、缓存与 `/private/tmp` 实验文件均不随 Git 提交，另一台机器应重新构建。重新构建/签名宿主后须重新生成匹配模块，不得混用旧模块。当前 VERSION 仍为 `v0.2.0-driver-process.4`，本开发尚未发布；正式发布前需按发布流程确定新版本并接入模块发布包。
