# ODBC 安装记录：macOS ARM64

记录日期：2026-09-23。安装根目录：`~/.dbx/drivers/`。

| 数据库 | 已安装版本 | 安装方式 | 验证结果 |
|---|---|---|---|
| Oracle | Instant Client 23.26.2.0.0 | 官方 Basic + ODBC ARM64 DMG，合并保留完整文件 | 下载包 SHA-256 一致；加载与句柄检查通过；原包 ODBC 库签名检查报错 |
| SQL Server | Microsoft ODBC Driver 18.7.1.1 | 从已安装的微软官方 Homebrew 包复制完整目录 | 签名、加载与句柄检查通过 |
| OceanBase | Connector/ODBC 2.0.8.3 | 官方源码加本地编译适配，构建 ARM64 动态库及认证插件 | 编译、加载、unixODBC 调用和文件校验通过 |

三项均已被源码运行的 `dbx driver list`（现名 `dbx odbc list`） 发现。没有实际数据库账号和服务器信息，尚未验证登录、TLS 握手、SQL 执行及 OceanBase 两种租户模式。达梦未安装，未取得可确认的官方 macOS ARM64 ODBC 包。

## 目录与使用

```text
~/.dbx/drivers/
├── oracle/instantclient_23_26/           # Oracle Basic + ODBC 完整内容
├── sqlserver/18.7.1.1/                   # 微软完整包，含 lib、share 和许可证
├── oceanbase/2.0.8.3-macos-arm64/
│   ├── lib/mariadb/libobodbc.dylib
│   ├── lib/plugin/                      # 动态认证插件
│   ├── docs/                            # 许可证
│   ├── build-info/                      # 源码归档、补丁、复现脚本及验证代码
│   ├── SHA256SUMS
│   └── README.md
├── libsqora.dylib                        # → oracle/instantclient_23_26/libsqora.dylib.23.1
├── libmsodbcsql.18.dylib                  # → sqlserver/18.7.1.1/lib/libmsodbcsql.18.dylib
├── libobodbc.dylib                       # → oceanbase/2.0.8.3-macos-arm64/lib/mariadb/libobodbc.dylib
├── env.sh
└── README.md
```

根目录用符号链接供 DBX 第一层扫描，实际库和配套文件保留在各自子目录。驱动不进入 DBX 编译包。

在启动 DBX 的同一终端执行：

```bash
source "$HOME/.dbx/drivers/env.sh"
dbx odbc list
```

`env.sh` 的内容：

```bash
export DBX_DRIVER_DIR="$HOME/.dbx/drivers"
export DYLD_LIBRARY_PATH="/opt/homebrew/opt/unixodbc/lib${DYLD_LIBRARY_PATH:+:$DYLD_LIBRARY_PATH}"
```

此次安装通过环境变量使用新目录，源码默认目录仍是 `~/.config/dbx/drivers`。连接文件中已有的 `driver_dir` 优先级更高，需要指向实际安装根目录。桌面程序启动 DBX 时也需要传入对应环境变量。

DBX 必须使用启用 ODBC 的 macOS 构建，见 [构建说明](odbc.md#构建)。运行依赖沿用 Homebrew 的 unixODBC 2.3.14；OceanBase 另外依赖 OpenSSL 3.6.4，未将这两个运行库复制进驱动目录。

## Oracle

从 [Oracle 官方 ARM64 下载页](https://www.oracle.com/database/technologies/instant-client/macos-arm64-downloads.html)下载同版本的 Basic 和 ODBC 包。校验后只读挂载两个 DMG，将内容复制到同一 `instantclient_23_26` 目录，保留动态库、符号链接、资源和许可证。使用完整 Basic 包，保留其字符集支持。

| 安装镜像 | SHA-256 |
|---|---|
| `instantclient-basic-macos.arm64-23.26.2.0.0.dmg` | `c3c4fce37a557192322717c1a8422fe098fc829a41f382f57a288b24ad6ba11e` |
| `instantclient-odbc-macos.arm64-23.26.2.0.0.dmg` | `377205537e3a4a7770aa385d48ba4e67b04950c57b23be8875875f77eaa2c984` |

ODBC 库引用了厂商构建环境中的 unixODBC 路径，因此通过上述 `DYLD_LIBRARY_PATH` 解析到本机 Homebrew 运行库。

`codesign --verify` 对 ODBC 库报 `invalid signature`；只读安装镜像内的原始库同样报错，且与安装副本哈希一致。未修改或重新签名此库。动态加载和 ODBC 环境句柄申请、释放通过，但不能视为厂商签名验证通过，签名问题需厂商确认。

## SQL Server

来源为微软维护的 `microsoft/mssql-release` Homebrew tap，按[微软官方说明](https://learn.microsoft.com/en-us/sql/connect/odbc/linux-mac/install-microsoft-odbc-driver-sql-server-macos)安装。本次复用已有的 `/opt/homebrew/Cellar/msodbcsql18/18.7.1.1`，完整复制至 `~/.dbx/drivers/sqlserver/18.7.1.1`。

保留 `lib/libmsodbcsql.18.dylib`、`share/msodbcsql18/resources` 及许可证等文件，不只复制动态库。原 Homebrew 安装保持不变；DBX 通过新目录中的路径加载副本。签名校验、动态加载、`SQLDriverConnectW` 导出检查及环境句柄申请、释放均通过。

## OceanBase

本次使用官方源码构建，产物是本地适配版，**不是厂商发布或认证的 macOS 安装包**。

| 项目 | 来源／版本 |
|---|---|
| ODBC 源码 | [oceanbase/ob-connector-odbc](https://github.com/oceanbase/ob-connector-odbc)，提交 `243bdaf858535a87384b75708746a2cb46b9397a` |
| C 客户端依赖 | [oceanbase/obconnector-c](https://github.com/oceanbase/obconnector-c)，提交 `2618276f5dd11a512aab21518994fef21c197816` |
| 构建工具 | Apple Clang 17、CMake 4.4.3，Release，ARM64 |
| 链接方式 | C 客户端静态链接进驱动；unixODBC、OpenSSL 动态链接，保留 TLS 支持 |

编译适配包括：排除 Linux 专用参数；跳过引用旧目标的 macOS PKG 打包及 iODBC 后处理；指定客户端静态库并补齐 OpenSSL 链接；补系统头文件、已有函数声明和缺失的 `void` 返回类型。未修改 DBX 或 Skill 源码。

复现材料保存在 `~/.dbx/drivers/oceanbase/2.0.8.3-macos-arm64/build-info/`：

- `ob-connector-odbc-source.tar.gz`、`obconnector-c-source.tar.gz`：原始源码。
- `macos-arm64.patch`：本地适配补丁。
- `rebuild.sh`：接收“新的构建目录”和“最终安装绝对路径”，在构建目录的 `runtime/` 生成安装文件。
- `check_odbc.c`：经 unixODBC 调用驱动的验证程序。

认证插件的默认路径编译为本次安装目录中的 `lib/plugin`；移动目录后需指定厂商 `PLUGIN_DIR` 选项或重新构建。整个包仍依赖 `/opt/homebrew/opt/unixodbc` 和 `/opt/homebrew/opt/openssl@3`。

验证时访问本机 `127.0.0.1:1`，收到 OceanBase 驱动返回的连接失败错误，证明请求经过驱动管理器进入了厂商驱动。安装文件哈希、动态库和认证插件加载、本地代码签名完整性检查通过；本地签名不代表厂商签名或公证。

该安装记录仅供历史参考。当前 OceanBase MySQL 模式使用内置 MySQL 驱动，无需此库；Oracle 模式暂不支持，`odbc list` 不再自动识别 OceanBase ODBC 库。
