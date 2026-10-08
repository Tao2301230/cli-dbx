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
  - Oracle 使用 Go 1.26.6，仅独立模块受此工具链要求约束；根模块保持 Go 1.22
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

本地开发常用命令：

```bash
go build -o ./dbx ./cmd/dbx
go test ./...
(cd drivers/postgres && go test ./...)
(cd drivers/oracle && GOTOOLCHAIN=go1.26.6 go test ./...)
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
