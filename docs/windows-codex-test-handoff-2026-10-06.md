# DBX 独立驱动 Windows 测试交接

交接日期：2026-10-06。接收方：Windows 上的 Codex。目标是在 Windows 宿主机验证 DBX 主程序、独立 Go 驱动、ODBC 适配程序，再让当前 CuteJ 应用调用同一测试程序。测试结果应区分本地自检、真实数据库连接、CuteJ 智能体调用三个阶段。

## 1. 先确认交接版本

- Git 基线：`86b18d3ce5241148f132af211753c3c2de69b071`；当前分支：`codex/dbx-driver-process`。
- 测试版本：`v0.2.0-driver-process.3`。本功能及此前修复位于上述基线之后的功能提交；仅 checkout 基线不能取得本版功能。随附 20261004 源码快照保留当时未提交工作树的构建来源记录。
- 用户于 2026-10-06 授权在本分支提交本次相关改动。接收端继续测试时，未经用户授权不得 commit、stage、push、merge 或创建 tag；测试产物保存在独立目录。
- 主程序只内置 MySQL、SQLite。PostgreSQL、Oracle、达梦、SQL Server 通过独立 Go 驱动进程连接。Oracle、达梦、SQL Server 还各有独立 ODBC 适配程序。
- 驱动目录采用扁平结构。连接通过 `driver_path` 精确选版本，程序与同名 JSON 必须成对保存。

本次应使用 `driver-downloads` 中不含 `snapshot` 的 `.3` 包。`.1` 是旧目录布局，`.2` 是历史版本；重建前 `.3` 快照仅供追溯。同版本新旧包内部文件名相同，不能混合解压。

## 2. 发送到 Windows 的材料

Mac 端来源目录：`/Users/bantao/Documents/pworkspace/deploy_dev/cli-dbx/driver-downloads`。这些文件尚未发布为网络下载地址，Windows 端需要先收到实际文件；Mac 绝对路径不能在 Windows 直接使用。下载索引里的 `download_url` 当前为相对文件名。

最小交接材料：

1. 本文。
2. `dbx-sources_v0.2.0-driver-process.3_worktree-20261004.tar.gz` 及 `.sha256`；其中有完整测试所需源码、锁定依赖、构建脚本和 `SOURCE.json`。
3. `index.json`、`SHA256SUMS`，以及所选平台压缩包旁的 `.sha256`。
4. 下表当前架构的四个 Go 驱动；测试 ODBC 时再传三个 x64 ODBC 适配包。
5. 需要哪种 ODBC，再传对应厂商包。运行 Go 驱动不要求安装这些厂商 ODBC 库。

| 路线 | Windows 包名，`<arch>` 替换为 `amd64` 或 `arm64` |
|---|---|
| PostgreSQL Go | `dbx-driver-postgres_v0.2.0-driver-process.3_windows_<arch>.zip` |
| Oracle Go | `dbx-driver-oracle_v0.2.0-driver-process.3_windows_<arch>.zip` |
| 达梦 Go | `dbx-driver-dm_v0.2.0-driver-process.3_windows_<arch>.zip` |
| SQL Server Go | `dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_<arch>.zip` |
| Oracle ODBC | `dbx-driver-oracle-odbc_v0.2.0-driver-process.3_windows_amd64.zip` |
| 达梦 ODBC | `dbx-driver-dm-odbc_v0.2.0-driver-process.3_windows_amd64.zip` |
| SQL Server ODBC | `dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_windows_amd64.zip` |

源码快照 SHA256：

```text
028e60cc7dab81929d232c364171fad78f1244e26a251fd91d2217fe3b209e9f
```

可选厂商 Windows x64 原始包：

| 厂商 | 文件 | 使用说明 |
|---|---|---|
| Oracle Basic | `oracle-basic_23.26.3.0.0_windows_x64.zip` | 与同版本 ODBC 包合并至同一 Instant Client 目录 |
| Oracle ODBC | `oracle-odbc_23.26.3.0.0_windows_x64.zip` | 保留 Basic 依赖、许可和注册工具 |
| SQL Server 18 | `sqlserver-msodbcsql18-18.7.1.1-windows-x64.msi` | 按厂商安装器安装，版本 18 的 TLS 选项需匹配服务 |
| SQL Server 17 | `sqlserver-msodbcsql17-17.11.1.1-windows-x64.msi` | 可另行测试，不因连接失败自动切换版本 |
| 达梦 | `dm-navicat-dameng_odbc_win.zip` | Navicat 分发的 x64 包，DLL 产品版本 8.0.0.0；未附独立许可文件 |

达梦包不是已验证的厂商公开独立发行包，URL 未标注完整构建版本；优先使用目标机器已正式安装的达梦 ODBC。保留原包及来源说明，不能把保存状态解释为公开再分发许可。OceanBase 的材料目前用于保存和移植，本版 DBX 不开放 OceanBase ODBC 测试。

## 3. 当前已有验证和 Windows 待验证项

2026-10-04 的目录记录包含 107 个制品：最新 Go 24 包、ODBC 15 包，其余为历史包、厂商包、源码及依赖。包清单、摘要和 PE/ELF/Mach-O 架构已检查。

Mac ARM64 已验证四个 Go 子进程启动和脱敏；ODBC 运行库搬迁后，使用 SQLite ODBC 替身验证了中文、二进制、导入、事务和回滚。Mac Intel ODBC 已通过 Rosetta 启动。Mac CuteJ 工作区内的 SQLite CLI 自检也通过；这不等于 CuteJ 智能体已执行测试。

**Windows 尚未实机运行或连接真实数据库。** 接收端需要补上实际启动、Credential Manager、厂商 ODBC 注册、真实连接、数据类型、事务、取消和 CuteJ 调用的证据。SQLite 替身的通过结果不能替代 Oracle、达梦或 SQL Server 实库结果。

## 4. 架构和运行环境

默认建议在 Windows x64 宿主机执行 PowerShell。普通 Intel/AMD 64 位机器对应 `amd64`。

| 主程序 | Go 驱动 | ODBC 适配程序 | 测试范围 |
|---|---|---|---|
| Windows amd64 | windows/amd64 | windows/amd64 | 四种 Go 路线和三种 ODBC 路线 |
| Windows arm64 | windows/arm64 | 未提供 | 四种 Go 路线 |

主程序、驱动、厂商 DLL 必须使用一致架构。ARM64 Windows 若要通过 x64 模拟测试 ODBC，需要一整套 x64 主程序、适配器、厂商 DLL，且先确认该系统能够执行它们；不能用 ARM64 主程序混载 x64 适配器。主程序会按自身的 Go 运行架构校验驱动清单。

CuteJ 若使用 WSL/Linux 容器，执行环境应单独按 Linux 准备；不能将 Windows EXE、Windows Credential Manager 或 Windows ODBC 注册视为容器能力。先确认 CuteJ 的本次调用是 Windows 宿主机工具。

## 5. 校验和安装测试目录

将交接文件放在例如 `C:\dbx-incoming`。不要假设 CuteJ 内置或 PATH 中的 dbx 就是本版。每条原生命令执行后检查 `$LASTEXITCODE`；PowerShell 的 `$ErrorActionPreference` 不能代替此检查。

```powershell
$ErrorActionPreference = 'Stop'
[Console]::OutputEncoding = [Text.UTF8Encoding]::new($false)
$OutputEncoding = [Console]::OutputEncoding
$Incoming = 'C:\dbx-incoming'
$TestRoot = Join-Path $env:USERPROFILE 'dbx-windows-test'
$Arch = 'amd64' # ARM64 原生测试改成 arm64，同时跳过 ODBC
$Version = 'v0.2.0-driver-process.3'
$DriverDir = Join-Path $TestRoot 'drivers'
$ConfigDir = Join-Path $TestRoot 'config'
if (Test-Path -LiteralPath $TestRoot) { throw '测试目录已存在，先检查已有材料或选择新的目录名' }
New-Item -ItemType Directory -Path $DriverDir,$ConfigDir | Out-Null

# 先校验全部收到的旁文件摘要；只校验实际收到的文件。
$Sidecars = @(Get-ChildItem -LiteralPath $Incoming -Filter '*.sha256' -File)
if ($Sidecars.Count -eq 0) { throw '没有收到 .sha256 文件' }
foreach ($Sidecar in $Sidecars) {
    $Target = Join-Path $Incoming ($Sidecar.Name -replace '\.sha256$','')
    $Expected = ((Get-Content -LiteralPath $Sidecar.FullName -Raw).Trim() -split '\s+')[0]
    $Actual = (Get-FileHash -Algorithm SHA256 -LiteralPath $Target).Hash
    if ($Actual.ToLowerInvariant() -ne $Expected.ToLowerInvariant()) { throw "摘要不符：$Target" }
}

foreach ($Engine in 'postgres','oracle','dm','sqlserver') {
    $Zip = Join-Path $Incoming "dbx-driver-${Engine}_${Version}_windows_${Arch}.zip"
    Expand-Archive -LiteralPath $Zip -DestinationPath $DriverDir
}
if ($Arch -eq 'amd64') {
    foreach ($Engine in 'oracle','dm','sqlserver') {
        $Zip = Join-Path $Incoming "dbx-driver-${Engine}-odbc_${Version}_windows_amd64.zip"
        # 未接收 ODBC 包时，记录跳过；不要将缺包记录成通过。
        if (Test-Path -LiteralPath $Zip) { Expand-Archive -LiteralPath $Zip -DestinationPath $DriverDir }
    }
}
```

Windows 包解压后是同一目录下的 `.exe` 和 `.exe.json`，例如 `dbx-driver-postgres-0.2.0-driver-process.3.exe` 及其同名加 `.json` 的清单。不要改名或拆开清单，不要通过改摘要绕过校验，也不要直接运行驱动程序当成 CLI；驱动 stdout 使用 DBX 私有管道协议。

本次安装至测试目录，所以连接必须显式设置 `driver_path`。默认无版本驱动查找根目录依次为连接的 `driver_dir`、`DBX_DRIVER_DIR`、`~/.config/dbx/drivers`；Windows 的 `~` 指当前用户目录。默认加载器不扫描或选择版本，也不从 PATH 找驱动。

## 6. 从交接源码构建 Windows 主程序

**必须拿到包含本功能的提交或本次源码快照。** 原 `dist` 里的同名 `.3` 主程序不能仅凭版本字符串被视为与 2026-10-04 驱动匹配的工作树产物。最稳妥的交接方式是从快照编译主程序，继续使用已校验的独立驱动包。

需安装可满足模块依赖的 Go 工具链，并允许获取 `go.mod/go.sum` 锁定依赖。根模块声明 Go 1.22；仅在重编 Oracle 独立驱动时需要 Go 1.26.6。源码快照不含离线 Go 模块缓存或工具链。

```powershell
$SourceDir = Join-Path $TestRoot 'source'
New-Item -ItemType Directory -Path $SourceDir | Out-Null
& tar -xzf (Join-Path $Incoming 'dbx-sources_v0.2.0-driver-process.3_worktree-20261004.tar.gz') -C $SourceDir
if ($LASTEXITCODE -ne 0) { throw '源码解压失败；确认 Windows tar 可用' }
$Repo = Join-Path $SourceDir 'cli-dbx'
$Dbx = Join-Path $TestRoot 'dbx.exe'
Push-Location $Repo
$SavedGoEnv = @{}
foreach ($Name in 'GOOS','GOARCH','CGO_ENABLED') { $SavedGoEnv[$Name] = [Environment]::GetEnvironmentVariable($Name,'Process') }
try {
    $env:GOOS = 'windows'; $env:GOARCH = $Arch; $env:CGO_ENABLED = '0'
    & go version
    if ($LASTEXITCODE -ne 0) { throw 'Go 不可用' }
    $BuildFlags = '-s -w -X github.com/linlay/cli-dbx/internal/buildinfo.Version=v0.2.0-driver-process.3 -X github.com/linlay/cli-dbx/internal/buildinfo.Commit=86b18d3-worktree -X github.com/linlay/cli-dbx/internal/buildinfo.BuildTime=2026-10-06'
    & go build -trimpath -buildvcs=false -ldflags $BuildFlags -o $Dbx ./cmd/dbx
    if ($LASTEXITCODE -ne 0) { throw '主程序编译失败' }
} finally {
    foreach ($Name in $SavedGoEnv.Keys) { [Environment]::SetEnvironmentVariable($Name,$SavedGoEnv[$Name],'Process') }
    Pop-Location
}
& $Dbx --version
if ($LASTEXITCODE -ne 0) { throw "version failed" }
& $Dbx --help
if ($LASTEXITCODE -ne 0) { throw "help failed" }
Get-FileHash -Algorithm SHA256 -LiteralPath $Dbx
```

`86b18d3-worktree` 标识本快照生成时的工作树来源，不是一个新提交；后续提交不改变该历史快照的字节或摘要。记录实际构建工具链、目标架构和程序 SHA256。若需重建独立驱动，快照内可运行 `scripts/build-drivers.ps1 -TargetOS windows -TargetArch amd64 -OutputDir <新目录>`；ODBC 使用 `scripts/build-odbc-drivers.ps1 -OutputDir <新目录>`，后者目前固定构建 x64。不要覆盖原交接包。

## 7. 阶段一 SQLite 自检

该阶段不依赖外部数据库，用新创建的隔离数据库验证主程序。在 Windows PowerShell 5.1 和 PowerShell 7 均使用 UTF-8 无 BOM 写 TOML/JSON/CSV。

```powershell
$Utf8NoBom = [Text.UTF8Encoding]::new($false)
$SmokeConfig = Join-Path $ConfigDir 'smoke.toml'
$SmokeToml = @'
[connection]
engine = "sqlite"
path = "./smoke.sqlite"
allow_actions = ["query", "update", "schema"]
'@
[IO.File]::WriteAllText($SmokeConfig,$SmokeToml,$Utf8NoBom)

& $Dbx conn list --config $ConfigDir
if ($LASTEXITCODE -ne 0) { throw "conn list failed" }
& $Dbx conn test smoke --config $ConfigDir
if ($LASTEXITCODE -ne 0) { throw "conn test failed" }
& $Dbx schema smoke 'create table dbx_probe(id integer primary key, note text)' --config $ConfigDir
if ($LASTEXITCODE -ne 0) { throw "create table failed" }
& $Dbx update smoke "insert into dbx_probe(id,note) values (1,'中文测试'),(2,'driver-process.3')" --config $ConfigDir
if ($LASTEXITCODE -ne 0) { throw "insert failed" }
& $Dbx query smoke 'select * from dbx_probe order by id' --config $ConfigDir
if ($LASTEXITCODE -ne 0) { throw "query failed" }
& $Dbx inspect table smoke dbx_probe --config $ConfigDir
if ($LASTEXITCODE -ne 0) { throw "inspect table failed" }
```

成功命令退出码为 0，JSON envelope 的 `ok` 为 `true`；查询应返回两行，中文保持原样，表字段为 `id` 和 `note`。`smoke.sqlite` 位于配置 TOML 所在目录，因为 SQLite 相对路径以该目录解析。

建议继续在该隔离 SQLite 库测试：

| 测试 | 方法 | 通过标准 |
|---|---|---|
| 单语句保护 | `query smoke 'select 1; select 2'` | 非零退出码，`multiple_statements_blocked` |
| 命令动作匹配 | 用 `query` 提交 `insert` | `action_mismatch`，数据不变 |
| 提交与回滚 | `tx` 先成功两步，再让第二步访问不存在的表 | 成功事务写入；失败事务第一步也不落库 |
| 导入行数上限 | CSV 含两条新记录，设置 `--max-rows-affected 1` | 整批失败，两条都不落库 |
| 文件输出 | `export table` 输出到测试目录，再核对内容 | 行数、中文与原数据一致 |

完整语法以本次 `$Dbx tx --help`、`$Dbx import --help`、`$Dbx export --help` 为准。每项失败测试后重新查询表，确认没有部分提交。

## 8. 阶段二真实数据库及独立驱动

接收端需要从用户处取得：数据库类型、服务端版本、可达地址和端口、测试库或 Oracle 服务名、DM schema、测试账号，以及允许执行的测试范围。未提供某种服务时，将该项记录为“未测试／缺少服务”，不能按编译结果判为通过。

密码由用户在 Windows 交互终端执行 `& $Dbx secret encrypt` 后填入连接 TOML。v2 密文依赖当前 Windows 用户的 Credential Manager，Mac 配置中的密文不能直接迁移；换用户、服务账号或机器后也需要重新加密。不要把明文密码放进交接文档、聊天或测试报告。`password.env`、`password.cmd`、`password.file` 默认禁用。

真实连接默认只开放 `query`。下面示例保存为 `$ConfigDir\pg-native.toml`，占位符需替换。Windows TOML 路径可使用单引号保存反斜杠；相对 `driver_path` 以 TOML 所在目录解析。

```toml
[connection]
engine = "postgres"
driver = "native"
driver_path = '../drivers/dbx-driver-postgres-0.2.0-driver-process.3.exe'
host = "REPLACE_HOST"
port = 5432
user = "REPLACE_USER"
database = "REPLACE_DATABASE"
password = "dbx-aes-gcm:v2:REPLACE_WITH_WINDOWS_ENCRYPTED_VALUE"
timeout = "5s"
allow_actions = ["query"]
```

```powershell
& $Dbx conn test pg-native --config $ConfigDir
& $Dbx query pg-native 'select version()' --config $ConfigDir
& $Dbx inspect schema pg-native --config $ConfigDir
```

其他配置按表修改，并各自保存为与连接名相同的 `.toml`：

| 引擎 | 默认端口 | 程序名 | 数据库字段 | 最小探测 SQL |
|---|---:|---|---|---|
| mysql | 3306 | 不需要 driver_path | database | `select 1` |
| postgres | 5432 | `dbx-driver-postgres-0.2.0-driver-process.3.exe` | database | `select 1` |
| oracle | 1521 | `dbx-driver-oracle-0.2.0-driver-process.3.exe` | database 填服务名 | `select 1 from dual` |
| dm | 5236 | `dbx-driver-dm-0.2.0-driver-process.3.exe` | 使用 schema，不依赖 database 拼地址 | `select 1 from dual` |
| sqlserver | 1433 | `dbx-driver-sqlserver-0.2.0-driver-process.3.exe` | database | `select 1` |

每个可用服务依次验收：连接、元数据、查询、中文、NULL、超过 2^53 的整数、精确 decimal、二进制与空二进制、时间值。只在用户指定测试库内进行建表、导入、事务及回滚测试，并显式增加相应 `allow_actions`。已有生产配置不得被这次测试覆盖。

独立驱动的失败边界也需要验证，但使用复制出来的测试文件：

| 条件 | 预期 |
|---|---|
| driver_path 不存在或清单缺失 | `driver_missing`，不启动其他版本 |
| 驱动架构、引擎或 adapter 不匹配 | `driver_invalid` |
| 协议版本不匹配 | `driver_protocol_mismatch` |
| 程序摘要被破坏 | `driver_integrity_failed`，不执行破坏后的程序 |
| 超时或取消 | 子进程结束，无持续残留；不能把未完成写入报告成成功 |
| 真实连接错误 | 保留受控错误码，输出不泄露密码／完整凭据 DSN |

## 9. ODBC 的 Windows 验收

本版 Windows ODBC 使用系统 Driver Manager，不安装或携带 unixODBC。先按厂商安装流程准备完整客户端，再找实际注册的 DLL 路径。

DBX 适配器在其进程对应的注册表视图中读取 HKCU/HKLM 的 `SOFTWARE\ODBC\ODBCINST.INI`，比较注册项 `Driver` DLL 路径并取得注册名称。**只有 DLL 文件存在不够，`odbc_driver` 必须指向已注册的同架构 DLL。** 不手工写注册表来绕过厂商安装，也不向 `driver_path` 填 DLL 或驱动名称。

- SQL Server：使用相应 Microsoft MSI 安装。可在 Windows 的 64 位 ODBC 管理器查看驱动和版本。ODBC 18 的加密／证书选项应按测试服务器填写；不要为了通过测试默认关闭 TLS。
- Oracle：同版本 Basic 和 ODBC 合并，保留完整依赖及许可；按厂商说明准备 VC++ Runtime、PATH，并使用随包 `odbc_install.exe` 注册。厂商库目录保持自身布局。
- 达梦：优先使用正式 DM 客户端已注册的 `dodbc.dll`。若只有 Navicat 压缩包，先确认厂商安装及注册方法；不能只复制 DLL 后认定安装完成。

例如保存为 `sqlserver-odbc.toml`：

```toml
[connection]
engine = "sqlserver"
driver = "odbc"
driver_path = '../drivers/dbx-driver-sqlserver-odbc-0.2.0-driver-process.3.exe'
odbc_driver = 'C:\REPLACE_VENDOR_INSTALL_DIR\msodbcsql18.dll'
host = "REPLACE_HOST"
port = 1433
user = "REPLACE_USER"
database = "REPLACE_DATABASE"
password = "dbx-aes-gcm:v2:REPLACE_WITH_WINDOWS_ENCRYPTED_VALUE"
timeout = "5s"
allow_actions = ["query"]
```

DLL 的实际路径以厂商注册项为准，示例不是默认安装路径。Oracle、DM 分别换 engine、ODBC 适配程序、端口和库路径；Oracle database 仍为服务名，DM 使用 schema。通过 `odbc_options` 传递允许的非凭据选项，不能覆盖目标、凭据或驱动字段，ODBC 模式也不接受原始 `dsn`/`dsn_env`。

对同一可用数据库分别执行 native 和 odbc 的连接、元数据及数据类型测试；分别记录结果，不将一种路线通过推广到另一种，也不因失败自动回退。

## 10. 阶段三 CuteJ 应用调用

Mac 端曾查到 CuteJ 内置 DBX 为 v0.1.2；Windows 端必须重新核对自己的版本和程序路径，不能推断其版本相同。CuteJ／Agent Platform 挂载连接器后，会把连接器 bin 加入 PATH；系统终端里的 dbx 与智能体调用的 dbx 可能不是同一个程序。

先请 CuteJ 报告当前工作区、宿主机执行环境及 `dbx --version`，并定位实际使用的程序。用本次 `$TestRoot\dbx.exe` 的完整路径验证，避免 PATH 混淆。默认 CuteJ 数据根为用户目录下的 `.cutej`，Windows 还可能通过注册表 DataRoot 更改位置，因此不要硬编码实际工作区。

测试目录应位于该智能体允许读写的工作区内，或通过其正常的工作区权限配置让它访问。不要手工编辑 Platform 自动生成的 `ru-agents`，也不要用覆盖旧内置程序作为第一次自检步骤。

发给 CuteJ 的任务，可用 Windows Codex 确认后的实际路径替换占位符：

```text
请通过 Windows 宿主机工具执行下面这版 DBX：
程序：<实际测试目录>\dbx.exe
配置目录：<实际测试目录>\config
依次执行 --version、conn list、conn test smoke、
query smoke "select * from dbx_probe order by id"、inspect table smoke dbx_probe。
除 --version 外，所有命令都显式指定 --config <实际配置目录>。
汇报实际程序路径、版本、连接结果、查询行数和中文内容、表字段。
```

通过标准：CuteJ 实际工具调用使用新程序，版本正确，SQLite 连接成功，查询两行且中文正确。保留命令和结构化结果的证据。Codex 自己执行命令通过不能替代 CuteJ 智能体实际调用通过。

后续要验证默认配置发现时，再在持久的 agent 私有配置来源中按 Platform 配置流程添加连接，让 Platform 生成运行时配置。DBX 默认先查 `$AP_AGENT_CONFIG_HOME/dbx`，再回退到当前用户 `~/.config/dbx`；显式 `--config` 只读取指定文件或目录。重名或配置损坏不能触发意外回退。`AP_AGENT_CONFIG_HOME` 是 Platform 保留变量，不通过 agent 自定义环境覆盖它。

## 11. 测试报告和继续工作的范围

Windows Codex 请保存报告 `windows-dbx-test-report.md`，至少包含：

- Windows／PowerShell／Go 版本，CPU 和主程序运行架构，CuteJ 版本及测试工作区。
- 源码快照 SHA256、主程序 SHA256、选用驱动的文件名和 SHA256、厂商 ODBC 版本及注册状态。
- 每条测试的命令、退出码、关键结构化结果，以及“通过／失败／未测试／受环境阻塞”。
- 各数据库 native、odbc 分别报告，记录真实服务版本及测试范围；密码、Token、完整凭据 DSN 和真实密钥不进入报告。
- 子进程回收、类型保真、失败事务和导入超限是否整批回滚；失败项附最小复现和脱敏日志。
- CuteJ 是否实际调用了测试程序、是否使用预期配置，原有应用和连接是否仍正常。

开始时先完成 SQLite 自检及驱动安装校验，再处理实际数据库。对缺少的服务、账号、厂商库或 Go 网络依赖，说明所缺条件并继续可以独立完成的检查。发现代码问题时先复现和定位；保持用户原有变更，未经授权不提交 Git。报告测试结果与修改建议，未经用户授权不替换应用内置连接器或发布版本。

## 12. 给 Windows Codex 的开场指令

```text
请读取这份 DBX Windows 测试交接文档，接着完成 Windows 实机和 CuteJ 调用验收。
接收端未经用户授权不要提交 Git。代码来自 86b18d3ce5241148f132af211753c3c2de69b071 基线后的功能提交，
请使用该功能提交或随附 20261004 源码快照，以及不带 snapshot 的 .3 Windows 驱动；不能只 checkout 基线。
先核对架构和 SHA256，从快照构建主程序，在独立目录完成 SQLite 自检。
再按已提供的真实数据库条件测试 native 和 ODBC；没有服务的项目明确记录未测试。
最后核对 CuteJ 实际工具调用的程序路径、版本和配置，输出 windows-dbx-test-report.md。
Windows ARM64 的 ODBC 适配器未提供；Mac 密文不可迁移，Windows 密码需在本机用户下重新加密。
```
