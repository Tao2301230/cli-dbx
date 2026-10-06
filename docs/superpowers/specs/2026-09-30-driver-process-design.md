# DBX 独立驱动进程设计

用户已确认基线为 `86b18d3ce5241148f132af211753c3c2de69b071`；主程序内置 MySQL、SQLite，PostgreSQL、Oracle、达梦、SQL Server 以独立驱动提供。分支为 `codex/dbx-driver-process`。

## 执行边界

保留现有 CLI、配置、安全检查和输出 envelope。`internal/db` 作为路由层，内置引擎直接调用无驱动依赖的 `internal/dbcore`，四种外部引擎通过 `internal/driverprocess` 调用独立程序。`dbcore` 保存现有 SQL、元数据、导入、事务实现，不导入任何厂商驱动；只有内置入口和独立模块注册其对应驱动。

采用每个数据库操作启动一个子进程的实现。相比模拟 `database/sql/driver` 的逐语句 RPC，这个边界能复用现有执行代码，且使事务与整批导入天然保留在同一进程和连接中。不增加常驻守护进程或跨 CLI 调用的事务会话。

`drivers/postgres`、`drivers/oracle`、`drivers/dm`、`drivers/sqlserver` 各有独立 go.mod；Oracle 所需工具链只约束它自身。模块通过本地 replace 复用根模块的无厂商依赖执行代码。主程序不得传递依赖 pgx、Oracle、DM、SQL Server 客户端库。

## 协议与加载

协议版本固定为 1，通过 stdin/stdout 私有管道传输带长度限制的 JSON 帧。连接 DSN 和解密后的凭据只在请求体内。stdout 专用于协议，子进程 stderr 不直接回显，以免泄漏凭据。所有查询、写入、结构查询、导入及整批事务均覆盖；行数据分帧传输，标量显式保留 NULL、整数、浮点、字符串、二进制、时间等类型。

驱动根目录：连接 `driver_dir` > `DBX_DRIVER_DIR` > `~/.config/dbx/drivers`。每个引擎目录包含 `driver.json` 和可执行程序。清单字段为 `protocol_version`、`engine`、`version`、`os`、`arch`、`executable`、`sha256`。程序名必须为目录内普通文件名，拒绝目录逃逸和符号链接；启动前核验引擎、平台、协议和 SHA-256，不从 PATH 搜索、不自动下载或回退。

主程序执行前完成 SQL 分类、allow_actions、allow_tables、单语句和危险写入检查。受影响行数限制由主程序传入，并在驱动提交前检查；取不到影响行数时不能假定为零。驱动进程不是安全沙箱，驱动包是可信代码。

父进程超时或取消会终止并回收子进程；子进程也有操作超时与父管道关闭监听。写操作失败不自动重试；提交期间断线不能声称已回滚，应由用户查询最终状态。

## 交付与验证

保留主程序/连接器包，另生成四种驱动包及校验和。已有用户文档及 examples 不覆盖。补充配置、CLI help、连接器技能与开发说明。

自动验证包括内置 MySQL/SQLite 依赖保持、主程序不链接四种外部驱动、真实子进程 SQLite 测试替身、协议和摘要错误、取消回收、类型保真、查询/元数据/导入/导出、事务失败与行数超限回滚，以及安全策略拒绝时不启动驱动。四个独立模块分别测试和构建。没有真实数据库服务时，明确区分驱动构建/本地协议验证与真实数据库兼容验证。
