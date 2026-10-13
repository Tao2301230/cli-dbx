# 共享宿主与独立 Go 驱动模块

此开发分支提供 `driver = "module"`：磁盘保存一份 `dbx`，PostgreSQL、Oracle、达梦和 SQL Server 各自分发一个 `.dbmod` 文件及校验清单。用户机器无需 Go SDK，也不需要编译或安装数据库客户端。

当前模块绑定**完全相同的宿主二进制**、Go 1.26.8 和目标平台。更新或重新签名宿主后，需要重新生成匹配模块。此版本尚未提供跨宿主版本稳定 ABI。

## 构建

开发机需要 Go 1.26.8、Python 3；macOS arm64 还需要 Xcode Command Line Tools。构建脚本会在输出目录准备私有 SDK，暴露加载器需要的编译器内部包，不修改全局 Go 安装。已经准备好的私有 SDK 可复用。

```sh
python3 scripts/build-modules.py \
  --go /path/to/go1.26.8/bin/go \
  --out dist/modules
```

默认构建宿主及四个模块；可重复 `--engine postgres --engine oracle` 只选择所需模块。`--offline` 要求依赖已经缓存。

产物中需要分发的文件是：

```text
dbx                       # Windows 为 dbx.exe
postgres.dbmod
postgres.dbmod.json
oracle.dbmod
oracle.dbmod.json
dm.dbmod
dm.dbmod.json
sqlserver.dbmod
sqlserver.dbmod.json
sizes.json                # 构建体积报告
```

`_build/` 是 SDK、对象文件、打包器及缓存，**不随产品分发**。每个模块及清单可独立复制，无需复制其他模块。

构建选项固定为 `-tags=dbmod -trimpath -buildvcs=false -ldflags=-w`。不要加 `-s`：宿主符号表用于模块链接。普通 `go build ./cmd/dbx` 保留原 native/ODBC 进程模式，不包含模块加载器。

需要复用已有宿主时，使用 `--host /path/to/dbx`。脚本会先从当前源码、依赖和相同构建参数重新生成宿主，只有 SHA-256 完全一致才继续打包，避免同名 Go 符号背后的类型或调用约定不同。`cmd/dbmod-pack` 是内部构建步骤，要求宿主、对象和依赖归档来自同一次规范构建，不是任意 Go 可执行文件的插件转换器。

模块完整打包当前要求在目标平台原生执行。可以只交叉编译 Windows 宿主：

```sh
python3 scripts/build-modules.py \
  --go /path/to/go1.26.8/bin/go \
  --goos windows --goarch amd64 --host-only \
  --out dist/modules-windows
```

## 使用

将需要的模块及其 `.json` 清单复制到驱动目录，并使用本次构建生成的 `dbx`：

```toml
[connection]
engine = "postgres"
driver = "module"
driver_path = "drivers/postgres.dbmod"
host = "localhost"
port = 5432
database = "app"
user = "app"
allow_actions = ["query"]
```

密码配置沿用现有连接文件规则。`driver_path` 支持 `~` 和相对于 TOML 的路径。省略它时，按 `driver_dir`、`DBX_DRIVER_DIR`、`~/.config/dbx/drivers` 的优先级选择目录，再查找 `<engine>.dbmod`。目录中可同时放置四个模块，仅加载此次操作所选的模块。

```sh
dist/modules/dbx conn test remote
dist/modules/dbx query remote 'select 1'
```

原有 `driver = "native"`（或空值）以及 `driver = "odbc"` 配置仍可使用，不会自动切换为模块模式。

## 磁盘、内存与校验

父进程启动**同一个 dbx 文件**的内部 worker。驱动不再各自携带一份宿主程序；每个 worker 的 Go 运行时内存仍然独立。每次操作的连接、整批事务或导入都留在一个 worker，沿用原有私有管道、凭据传输、超时终止和回收机制。

DBXMOD2 使用前缀字典、紧凑重定位记录、提前链接和稀疏元数据描述。省略的 NOP 预留空间在加载时按原偏移重建，因此磁盘减少量不能直接当作运行时 RAM 减少量。类型、GC、栈和初始化信息继续保留；模块到 worker 退出才释放，不支持操作中卸载。

主进程和 worker 都校验平台、工具链、引擎、文件名以及宿主/模块摘要。worker 解码的是已校验的不可变字节快照，避免校验后文件改写影响执行。模块和清单必须是普通文件。摘要用于一致性检查，不是发布者签名；只运行可信来源的模块。

## 验证方法与边界

普通单测及模块单测需使用已准备的私有 SDK：

```sh
export GOROOT=/path/to/private/go1.26.8
export PATH="$GOROOT/bin:$PATH"
export GOTOOLCHAIN=local
go test ./...
go test -tags=dbmod ./...
python3 -m unittest discover -s scripts/tests -v
```

四个模块可分别做不联网的构造和 GC 检查：

```sh
dist/modules/dbx --internal-module-check postgres dist/modules/postgres.dbmod
dist/modules/dbx --internal-module-check oracle dist/modules/oracle.dbmod
dist/modules/dbx --internal-module-check dm dist/modules/dm.dbmod
dist/modules/dbx --internal-module-check sqlserver dist/modules/sqlserver.dbmod
```

这些检查只验证加载、注册、连接器构造、关闭和 GC，输出 `database_connected: false`。它们不证明真实数据库认证、TLS、事务或全部驱动功能兼容。

实际 CLI 与本地 PostgreSQL 协议替身的集成检查：

```sh
DBX_MODULE_HOST="$PWD/dist/modules/dbx" \
  python3 -m unittest discover -s scripts/tests -p test_module_host.py -v
```

此测试覆盖策略拒绝、真正的 pgx 模块查询、精确小数/整数/中文传输和连接关闭，仍不是连接真实 PostgreSQL。Windows 当前只有宿主交叉构建证据，模块打包和实际运行需要 Windows 环境验证。

最终体积和验证记录见同目录的 `shared-host-modules-validation.json`。体积以未压缩 `.dbmod` 为主，gzip 单独列作下载对照；未承诺所有驱动已达到内嵌增量。
