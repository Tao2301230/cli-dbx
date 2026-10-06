# DBX 可选数据库 Connector 设计

> 历史讨论记录：其中自定义 Connector、进程协议及 ODBC 自动回退方案不再采用。当前实现改为直接调用厂商 ODBC 驱动，见 [ODBC 说明](../../odbc.md)；本文不作为当前开发指令。

> 状态：已确认
>
> 基线：`cli-dbx` commit `b8b7534e3e9a910b385af544b6f6dab373493fd4`

## 目标

- `builtin.dbx` 保留 MySQL、PostgreSQL、SQLite。
- Oracle、达梦、SQL Server、OceanBase 分别作为用户按需安装的外置 Connector。
- 用户继续使用统一的 `dbx` 命令。
- 新增数据库尽量不要求修改或重新发布 `builtin.dbx`。

## 总体结构

```text
builtin.dbx             基础 CLI、Skill、安全策略、现有三个数据库
dbx.oracle              Oracle 外置扩展
dbx.dameng              达梦外置扩展
dbx.sqlserver           SQL Server 外置扩展
dbx.oceanbase           OceanBase MySQL/Oracle 双模式外置扩展
```

外置 Connector 只在用户需要时安装。选择外置 Connector 的 Agent 自动挂载 `builtin.dbx`；卸载扩展不影响基础 Connector。

## 职责边界

### `builtin.dbx`

- 命令解析和统一输出。
- 连接配置读取与密码解密。
- `allow_actions`、`allow_tables` 和 SQL 安全检查。
- MySQL、PostgreSQL、SQLite 的现有实现。
- 驱动注册表读取、协议握手和扩展进程调用。

### 外置 Connector

- 数据库驱动和 ODBC 适配。
- DSN、参数占位符和标识符引用。
- 连接测试、SQL 执行、事务及元数据查询。
- 数据库版本和能力报告。

外置 Connector 不重复携带 `builtin-dbx` Skill。

## 安装与运行位置

`builtin.dbx` 继续由 Agent Platform 内置资源提供。四个外置扩展安装在环境的 `connectors-center`：

```text
zenmind-env/
└── connectors-center/
    ├── dbx.oracle/
    │   ├── connector.json
    │   ├── cli.json
    │   └── bin/dbx-oracle
    ├── dbx.dameng/
    ├── dbx.sqlserver/
    └── dbx.oceanbase/
```

Windows 二进制使用 `.exe`。Agent 选中扩展后，Platform 从 `connectors-center` 将它与依赖的 `builtin.dbx` 一起物化到 `ru-agents/<agentKey>/connectors/`。`connectors-center` 是外置 Connector 的安装和分发目录，不是 Platform 内置资源目录。

各扩展使用独立 Go Module 和发布产物；基础 `cli-dbx/go.mod` 不声明、不下载、不编译四个可选驱动。

## Connector 声明

Agent Platform 的 Connector Manifest 增加依赖和 DBX 扩展信息：

```json
{
  "id": "dbx.oracle",
  "name": "DBX Oracle",
  "version": "1.0.0",
  "type": "cli",
  "auth_mode": null,
  "requires": {
    "connector": "builtin.dbx",
    "minVersion": "0.1.2",
    "protocol": "dbx-driver/v1"
  },
  "provides": {
    "dbxEngines": ["oracle"],
    "dialectProfile": "oracle",
    "capabilities": ["ping", "query", "execute", "transaction", "inspect"]
  }
}
```

Platform 校验依赖、版本、协议、引擎冲突和二进制存在性。运行时生成只包含已挂载扩展的驱动注册表，记录引擎、绝对路径、版本和协议。DBX 不扫描全局 PATH。

## 调用协议

- 协议名称：`dbx-driver/v1`。
- 每次 DBX 命令启动一次扩展进程，不使用常驻服务。
- 请求通过 stdin 发送 JSON，密码不得进入命令参数、环境变量或日志。
- stdout 只输出协议响应；查询行使用 JSON Lines 流式返回。
- stderr 只输出已脱敏的诊断信息。
- 调用前完成协议和能力握手。
- 事务在单次扩展进程内完成，不跨 CLI 调用保存会话。

缺失扩展、版本不兼容、能力不支持、ODBC 不可用、驱动加载失败、超时和数据库错误使用不同的稳定错误码。

## ODBC 与内置驱动

每个外置 Connector 支持：

```toml
backend = "auto" # auto | odbc | embedded
```

`auto` 的规则：

1. 查找与数据库及 CPU 架构匹配的 ODBC Driver；只有 Driver Manager 不算可用。
2. ODBC Driver 不存在、架构不匹配或无法加载时，使用 Connector 内置 Go 驱动。
3. ODBC 已加载但连接失败时不回退，直接报告认证、TLS 或网络错误。
4. 输出和审计记录实际后端、驱动名称及版本。

内置 Go 驱动在 CI 中预编译进外置 Connector，不在用户机器上编译。它会增加外置包体积，但不会增加基础 DBX 体积。

## 安全边界

- 基础 DBX 在调用扩展前执行动作、表范围和 SQL 安全检查。
- 扩展声明 `dialectProfile`，只能选择基础 DBX 已知的安全方言配置。
- 新方言存在特殊写入或锁定语法时，必须先升级基础 DBX 的分析规则。
- Platform 注册表使用已验证的绝对路径，防止 PATH 中同名程序劫持。
- 数据库只读账号仍是生产环境的最终权限边界。

## 后续数据库扩展

新增数据库时：

1. 实现 `dbx-driver/v1`。
2. 声明引擎名称、别名、方言配置、配置 Schema 和能力。
3. 生成 macOS、Windows 所需架构的预编译包。
4. 安装后由 Platform 注册，`dbx driver list/help/doctor` 展示能力和诊断结果。

协议足够且能复用已有安全方言配置时，无需修改 `builtin.dbx`。只有新增协议能力或安全方言规则时才升级基础 DBX。

## 迁移范围

- 以 `b8b7534` 为基础保留现有三个数据库。
- 将提交 `118d5e9` 中的 Oracle、达梦、SQL Server、OceanBase 驱动、DSN、元数据和方言执行代码移入四个扩展。
- 将通用且与数据库驱动无关的安全修复保留在基础 DBX。
- Agent Platform 从 `connectors-center` 加载外置扩展，并增加 Connector 依赖解析、DBX 扩展注册表和自动挂载。

## 验证

- 基础 DBX 的 MySQL、PostgreSQL、SQLite 回归测试。
- 使用虚拟扩展进程验证握手、JSON/JSONL、超时、取消和错误码。
- 验证缺失扩展、版本不兼容、重复引擎和路径劫持均失败关闭。
- 验证 Agent 选择扩展时自动挂载 `builtin.dbx`。
- 验证 ODBC 检测、embedded 回退以及“连接失败不回退”。
- 四个扩展完成 macOS/Windows 构建；对实际数据库执行连接、元数据、查询、事务和典型类型测试。

## 完成标准

- 基础 `go.mod` 不包含四个可选数据库驱动；基础 DBX 的发布二进制不超过 `b8b7534` 同参数构建结果 1 MiB。
- 未安装扩展时，MySQL、PostgreSQL、SQLite 的命令和输出保持兼容。
- 安装任一扩展后，用户仍通过统一 `dbx` 命令使用对应数据库。
- 新增一个兼容现有协议和安全方言的数据库时，不修改基础 DBX 代码。
