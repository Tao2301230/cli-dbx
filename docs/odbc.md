# 厂商 ODBC 驱动

Oracle、SQL Server、OceanBase 在 `~/.dbx/drivers/` 下的版本、安装方式、目录树和验证结果见 [macOS ARM64 安装记录](odbc-installation-macos.md)。

当前支持范围：Windows x64、macOS M 系列（ARM64）；Linux 范围不变。OceanBase 仅支持 MySQL 模式；TDSQL 仅考虑 MySQL／PostgreSQL 兼容模式，分别复用内置驱动。

## 文档分工

- Skill 只维护运行时调用、连接配置、支持范围、安全边界和错误处理。
- 驱动仓库与固定版本、Go/CGO 要求、构建标签、打包和动态库重定位统一维护在本文；历史方案比较见 [技术选型](数据库连接器技术选型.md)。
- 不把自定义驱动协议、`driver.json` 等已弃用设计放入运行时 Skill。
- 修改技能时核对 `connector/skills/builtin-dbx/` 发布源码与 `../zenmind-env/skills-center/builtin-dbx/` 的操作说明；保留各自元数据，不把开发文档打进技能。

## 调用关系

```text
DBX
├── 纯 Go 驱动方式
│   └── database/sql → 编译进 DBX 的 Go 驱动 → 数据库
└── ODBC 方式
    └── database/sql → Go ODBC 适配层 → ODBC 驱动管理器
        └── 外部厂商 ODBC 动态库及其依赖 → 数据库
```

DBX 不依赖自定义驱动程序。用户安装厂商原包；不存在 `dbx-driver-*` 协议、独立 Go 驱动模块或 `driver.json`。

当前项目中的对应路径：

```text
cli-dbx/
├── internal/config/                 # 读取 driver_dir、odbc_driver 等配置
├── internal/conn/native.go           # 统一组装所有内置 Go 驱动的连接参数
├── internal/conn/odbc.go             # 组装厂商 ODBC 连接参数
├── internal/sqlanalyzer/             # 执行前的 SQL 安全检查
├── internal/db/                     # database/sql 执行、事务、导入导出
│   ├── db.go                        # 连接、执行、结构查询与导入导出；原有函数内扩充数据库分支
│   └── db_test.go                   # 数据库执行层测试
└── internal/odbcdriver/              # 检索厂商库、引入 alexbrainman/odbc 并检查构建支持
    └── 系统 ODBC 驱动管理器          # Windows 系统组件；macOS 随包提供 unixODBC，Linux 自行安装
        └── 厂商驱动目录              # DBX 编译包外的 DLL/dylib/so 及配套依赖

cli-dbx/connector/skills/builtin-dbx/  # 连接器发布的技能源码
├── SKILL.md                         # 智能体使用 DBX 的入口说明
└── references/                      # 配置、命令、缺失驱动处理等说明
```

树图末尾的系统组件和厂商目录是运行依赖，不是仓库子目录。

## 纯 Go 与 ODBC 对比

| 对比项 | 纯 Go 驱动 | ODBC |
|---|---|---|
| 连接方式 | Go 驱动实现数据库网络协议，直接连接 | Go 适配层通过驱动管理器调用厂商驱动 |
| 驱动位置 | 通常编译进 DBX 可执行文件 | 厂商动态库安装在 DBX 外部 |
| 用户安装 | 通常只需 DBX；特殊认证等功能可能另有依赖 | macOS 发布包自带管理器，Windows 使用系统管理器；另装厂商驱动及其依赖 |
| 编译环境 | 驱动本身通常无需 C 编译器，便于交叉编译 | 取决于适配层；当前 DBX 在 macOS/Linux 需要 CGO 和 unixODBC 开发库 |
| 跨平台 | 取决于驱动实现，纯 Go 不保证支持所有系统 | 取决于适配层、驱动管理器和厂商库是否均支持目标系统、架构 |
| 增加数据库 | 增加 Go 依赖，通常需要重新编译 DBX | 可复用 ODBC 接入层，仍需适配连接参数、方言和元数据 |
| 驱动升级 | 通常重新编译并发布 DBX | 可独立更新外部厂商驱动，并验证兼容性 |
| 包体积 | 驱动代码增加 DBX 大小 | 厂商库不进入 DBX；额外安装包仍占空间，总体积未必更小 |
| 按需安装 | 普通 Go 源码包不能放入目录后被现有 DBX 自动加载 | 适合用户安装厂商库、DBX 按配置路径检索 |

OceanBase MySQL 模式的模板填写 `engine = "mysql"`，复用 DBX 内置的 `go-sql-driver/mysql`，不检索 ODBC 库，不需要额外安装或启用 ODBC 构建。Oracle、达梦、SQL Server 默认使用内置 Go 驱动；显式 `driver = "odbc"` 可使用厂商库。OceanBase Oracle 模式暂不支持。DBX 只检索和调用，安装由用户或智能体通过外部工具完成。安装 unixODBC 只提供驱动管理器，不等于安装了数据库厂商驱动。

GORM 是上层 ORM，底层仍需要 Go 驱动或其他数据库适配器，不是与纯 Go、ODBC 并列的第三种通信驱动。ODBC 也不会自动解决大表 JOIN 的执行效率问题，仍需依据索引、执行计划和 SQL 设计优化。

## 新增数据库的驱动与平台支持

记录日期：2026-09-23。下表比较厂商与驱动方案，不依据某台开发机；“未确认”表示公开资料不足，不等于技术上绝对无法实现。Windows 支持不自动包含 ARM64，具体以厂商版本为准。

| 数据库 | Go 接入方式 | 厂商 ODBC | macOS | Windows |
|---|---|---|---|---|
| Oracle | 官方纯 Go `oracle/go-oracledb`；社区纯 Go `go-ora`；`godror` 经 OCI；ODBC | 有，Instant Client ODBC | 有 Intel、Apple Silicon ODBC 包；纯 Go 方案按版本验证 | 有 ODBC；也可使用 Go 驱动 |
| 达梦 | 官方 `dm-go-driver`；ODBC；GORM 底层使用 DM 驱动 | 有，公开文档覆盖 Windows/Linux | 未确认官方 macOS ODBC 包；不能仅凭 Go 源码包认定兼容 macOS | 官方提供 DM Go、ODBC 接入说明 |
| SQL Server | 官方纯 Go `microsoft/go-mssqldb`；ODBC | 有，Microsoft ODBC Driver | Go 驱动可用；ODBC 支持 Intel，17.8 起支持 Apple Silicon | Go 驱动、ODBC 均有支持 |
| OceanBase MySQL 模式 | `go-sql-driver/mysql`；GORM MySQL；ODBC 需核对所选版本的模式支持 | 有 OceanBase Connector/ODBC | 可走纯 Go MySQL 驱动；ODBC 已从官方源码本地构建 ARM64 版，实库待测 | Go MySQL 驱动可用；ODBC 有 Windows 包，需核对模式 |

Oracle 官方 Go 驱动说明支持 Oracle 19c+。达梦的公开 Go/ODBC 文档覆盖 Windows/Linux，已有官方包缺少 macOS 实现的记录不能外推为未来所有版本均不支持。OceanBase ODBC 公开安装资料覆盖 Windows/Linux；Oracle 模式兼容不等于可以直接使用 Oracle 数据库驱动。

`oceanbase/go-oceanbase-driver` 仓库确实存在，但本次查阅的 README 没有明确声明 Oracle 模式支持，因此不能仅凭仓库名称认定其可替代 Oracle 模式 ODBC。

当前内置 Oracle 官方 `v26.0.1-beta`、微软 `go-mssqldb` 和社区 `gitee.com/chunanyong/dm`；不按 Oracle 服务端版本主动拦截，也不自动切换 ODBC。当前 DBX 的 ODBC 构建限制见下文“构建”：Windows x64 可用，Windows ARM64 暂不支持；macOS 发布包已启用 ODBC 并附带 unixODBC，Linux 需自行启用；两者仍需对应厂商库。

资料：[Oracle Go 驱动](https://github.com/oracle/go-oracledb)、[Oracle ARM64 ODBC 下载](https://www.oracle.com/database/technologies/instant-client/macos-arm64-downloads.html)、[达梦 DM 接口](https://eco.dameng.com/document/dm/zh-cn/app-dev/go_dm.html)、[达梦 Go ODBC 接口](https://eco.dameng.com/document/dm/zh-cn/app-dev/go_odbc.html)、[微软 Go 驱动](https://github.com/microsoft/go-mssqldb)、[微软 macOS ODBC](https://learn.microsoft.com/en-us/sql/connect/odbc/linux-mac/install-microsoft-odbc-driver-sql-server-macos)、[OceanBase ODBC](https://en.oceanbase.com/docs/obodbc-en)、[OceanBase Go 驱动仓库](https://github.com/oceanbase/go-oceanbase-driver)。

## 配置与发现

`engine` 始终填写数据库类型或兼容方言，`driver = "odbc"` 选择 ODBC 接入。厂商库通过 `odbc_driver` 指定，附加连接属性通过 `odbc_options` 提供。

- 目录优先级：连接的 `driver_dir` → `DBX_DRIVER_DIR` → `~/.config/dbx/drivers`。
- `odbc_driver` 指定目录内的相对文件路径；省略时匹配第一层常见厂商名称。多个不同版本同时存在时要求明确选择，同一文件的符号链接按一个候选处理。
- Oracle 匹配 `sqora` / `libsqora`，达梦匹配 `dodbc` / `libdodbc`，SQL Server 匹配 `msodbcsql` / `libmsodbcsql`；文件后缀随操作系统选择。使用内置驱动的连接不参与 ODBC 检索。
- 自动匹配未覆盖的厂商文件名，直接通过 `odbc_driver` 指定，不需要重命名厂商文件。
- Windows 从当前进程架构的 ODBC 注册信息中，查找指向该文件的驱动名称；未注册时提示使用厂商安装器。DBX 不写注册表。
- unixODBC 使用所选库的绝对路径作为 `DRIVER`。
- `odbc dir/list` 仅查看环境默认目录，列出候选文件并不等于已经连库验证。

缺失时返回 `driver_missing`，`market_package` 使用 `oracle-odbc`、`dameng-odbc`、`sqlserver-odbc` 作为市场检索标识。实际是否有包、包来源和签名由外部市场服务确认；这些标识不代表已发布的 DBX 自定义程序。

## 安装边界

DBX 不联网、不安装、不改名、不复制厂商驱动。智能体或用户通过外部工具获取可信厂商包，按厂商要求验证并安装，保留驱动附带的依赖库、资源文件和目录结构，再配置目录重试。只复制单个 DLL/动态库不保证可用。

Windows 的驱动名称注册、Oracle 客户端依赖、Linux 的 unixODBC 安装，以及厂商系统/架构支持范围，均遵循各自安装说明。数据库类型由连接配置决定，不根据连接失败自动切换。

## 构建

基础构建支持六种内置数据库，要求 Go ≥ 1.26.6（Oracle 官方 beta 驱动的要求）：

```bash
go build -o dbx ./cmd/dbx
```

Windows x64 构建直接包含 Go ODBC 适配层，使用系统 ODBC Driver Manager；不需要 CGO。厂商 DLL 不进入 DBX 编译包。Windows ARM64 不在当前支持和发布范围内。

### macOS M 系列发布包

在 Apple Silicon 构建机准备 Go、Xcode Command Line Tools、Python 3、unixODBC 和 libtool。运行：

```bash
bash scripts/release/build.sh --macos-binary /tmp/dbx-release/dbx
# 完整发布：bash scripts/release/build.sh
```

脚本默认从 Homebrew 的 unixodbc/libtool 前缀读取编译头文件和运行库；可用 `DBX_UNIXODBC_PREFIX`、`DBX_LIBTOOL_PREFIX` 指定其他安装前缀。不会自动安装构建依赖。macOS 包必须在 Apple Silicon 上构建，其他宿主的发布脚本只生成 Linux/Windows 包。

```text
独立 CLI 包                     Platform 连接器包
├── dbx                         ├── bin/dbx
├── lib/                        ├── bin/lib/
│   ├── libodbc.2.dylib           │   ├── libodbc.2.dylib
│   ├── libodbcinst.2.dylib       │   ├── libodbcinst.2.dylib
│   ├── libodbccr.2.dylib         │   ├── libodbccr.2.dylib
│   └── libltdl.7.dylib           │   └── libltdl.7.dylib
└── licenses/                   └── licenses/
```

用户保留完整目录即可，无需安装 Homebrew 或 unixODBC。运行库使用相对加载路径，不修改构建机已安装的库。随包保留许可证与作者信息，发布者应同时提供对应版本的库源码。厂商驱动及其依赖仍单独安装；厂商库若硬编码了其他库的绝对路径，仍需按厂商说明处理。

最低 macOS 版本取决于构建时所用的运行库和工具链；打包不会降低动态库已有的最低系统要求。需要支持旧 macOS 时，构建机应提供兼容版本的库。

### Linux 与自定义构建

Linux 发布包维持 `CGO_ENABLED=0`；需要 ODBC 时安装 unixODBC 开发及运行库，再执行：

```bash
CGO_ENABLED=1 go build -tags odbc -o dbx ./cmd/dbx
```

macOS 手动执行上述命令也可链接本机 unixODBC，但只有发布脚本会收集并重定位配套库。自定义基础构建未启用 ODBC 时仍返回 `odbc_unavailable`。

## 执行与安全

连接使用结构化字段和 `odbc_options`，不接受 ODBC 原始 `dsn` / `dsn_env` / `--dsn` 绕过所选驱动。原有数据库的连接串行为不变。`allow_actions`、`allow_tables` 和单语句检查在进入 ODBC 前执行。ODBC 写入无法获得有效影响行数时，受保护事务失败并回滚。

ODBC 查询取消最终依赖厂商驱动响应；所用 Go 库的连接建立和部分执行操作不保证被 Go context 硬中断。需要按厂商支持配置登录/查询超时，不能把 DBX 的 `timeout` 理解为所有原生调用的硬截止时间。

## 验证范围

已验证基础测试、Windows 交叉编译、macOS ODBC 构建，以及使用临时 SQLite ODBC 库的参数绑定、中文读写、事务提交和超限回滚。该测试仅验证 ODBC 调用链；四种厂商数据库的连接参数、元数据查询及版本兼容性仍需在实际部署环境验收。

## 来源

- [Go 官方数据库驱动列表](https://go.dev/wiki/SQLDrivers)、[Go ODBC 适配库](https://github.com/alexbrainman/odbc)。
- [Microsoft SQLDriverConnect](https://learn.microsoft.com/en-us/sql/odbc/reference/syntax/sqldriverconnect-function)、[unixODBC](https://www.unixodbc.org/odbcinst.html)。
- [Oracle ODBC 驱动](https://docs.oracle.com/en/database/oracle/oracle-database/19/adfns/odbc-driver.html)、[Oracle Instant Client 安装](https://www.oracle.com/database/technologies/releasenote-odbc-ic.html)。
- [达梦 ODBC 接口](https://eco.dameng.com/document/dm/zh-cn/app-dev/c_c%2B%2B_odbc.html)。
- [SQL Server macOS ODBC 安装](https://learn.microsoft.com/en-us/sql/connect/odbc/linux-mac/install-microsoft-odbc-driver-sql-server-macos)。
- [OceanBase Connector/ODBC](https://en.oceanbase.com/docs/common-oceanbase-database-10000000001231117)。
