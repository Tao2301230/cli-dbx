# 安装、下载服务和移植

本仓库跟踪本目录的 DBX 说明文档。安装包、编译产物、厂商授权文件、索引和校验记录保留本地，按发布流程分发。README.md 和 SIZES.md 中的文件链接适用于包含完整产物的本地下载目录或下载服务。

## 固定入口

所有文件位于同一 driver-downloads 目录。托管此目录即可得到稳定下载入口；index.json 的 download_url 默认是相对文件名。需要完整 URL 时，用 `python3 scripts/release/driver-downloads.py driver-downloads --records driver-downloads/catalog-records.json --base-url https://你的域名/drivers` 重新生成；随后加 `--verify` 校验。本次只保存本地制品，未发布到网络。

## 三类文件

1. dbx-driver-*：DBX 独立程序。native 为 Go 路线，odbc 为调用厂商库的 ODBC 路线。
2. oracle/sqlserver/dm/oceanbase-*：厂商原始安装包或完整库副本，需要结合平台和实际数据库选用。Oracle Basic 与 ODBC 需同平台同版本合并在 Instant Client 目录。Windows 使用安装器完成驱动注册，不只是把 DLL 拷入 DBX 目录。
3. *source* / *porting-kit*：源码、补丁、锁定依赖和构建复现材料。OceanBase 本地版不是厂商发布二进制，当前 DBX 不支持该 ODBC 路线；保留其材料供后续移植，不新增连接能力。

包保存状态和 DBX 支持状态分别记录。厂商原包未改名内部文件、未改写库或运行安装器，已有许可文件随包保留；Navicat 分发的达梦包未附独立许可文件，索引不承诺公开再分发权限。Linux vendor 包还需匹配发行版、libc 和 OpenSSL，不能把同 CPU 的任意 .so 通用替换。

## 安装 DBX 驱动

选定目标系统及 CPU，校验 SHA256 后将整个 DBX 包解压到目标 drivers 目录。当前 .3 为平铺安装，版本通过 driver_path 指定；Mac .tar.gz、Windows .zip、Linux .tar.gz。

```toml
[connection]
engine = "oracle"
driver = "native"
driver_path = "~/.config/dbx/drivers/dbx-driver-oracle-0.2.0-driver-process.3"
# ODBC 时改成 driver="odbc"、程序名 -oracle-odbc-，另加 odbc_driver="厂商库绝对路径"
```

Windows 程序名末尾加 .exe。历史 .1 是旧的引擎子目录布局，仅为版本保存，不适用于当前加载器；.2 为平铺版。snapshot 是本次重建前的同版本旧制品，供追溯使用，推荐当前未带 snapshot 的 .3。相同版本不同构建的包内程序同名，不能同时解压到一个运行目录。

## Go 原生移植与编包

解压 dbx-sources 包，进入 cli-dbx。需联网获取 go.mod/go.sum 锁定依赖，源码包不是完整离线 Go 工具链或模块缓存。

```bash
scripts/build-drivers.sh --output-dir driver-downloads
# 也可以选择一种目标
scripts/build-drivers.sh --os linux --arch arm64 --engine postgres --output-dir driver-downloads
```

默认构建 darwin/linux/windows × amd64/arm64 四个引擎共24包。主模块语言基线 Go1.22；Oracle 模块当前依赖要求 Go1.26.6，脚本仅为该模块设置 GOTOOLCHAIN。Windows 本机可使用 scripts/build-drivers.ps1。

## ODBC 移植

Mac 原生构建需 unixODBC 和 libtool；脚本打包对应架构运行库，重写打包副本的加载路径并签名。Mac vendor 库仍有自己的依赖要求，随 DBX 运行库不自动修复 vendor 硬编码路径。

```bash
scripts/build-odbc-drivers.sh --host --output-dir driver-downloads
scripts/build-odbc-drivers.sh --os windows --arch amd64 --output-dir driver-downloads
# Linux 本机：先装该发行版 unixODBC 开发包，再原生编译
scripts/build-odbc-drivers.sh --os linux --arch amd64 --output-dir driver-downloads
```

porting-kit 另含 reproduce-macos-intel-odbc.sh，可在 Apple Silicon 使用 Xcode clang 编译 Intel 适配器及共享 unixODBC/libltdl。Intel 程序已通过 Rosetta 启动检查，未连接厂商数据库。

porting-kit 另含 Apple Silicon 宿主交叉编译 Linux 的 reproduce-linux-odbc.sh，锁定 Zig0.14.1 和 Debian12 unixODBC2.3.11 的官方校验值，无需安装到系统目录。这批 Linux ODBC 包需 glibc >=2.34 和匹配架构 libodbc.so.2。Alpine/musl不能直接使用；需本机重建。Windows ODBC 当前仅 x64；Go 驱动支持 Windows ARM64，不能据此推导 ODBC ARM64 支持。

Oracle MacIntel19.16 与 MacARM23.26、Linux/Windows23.26 是不同厂商版本，以每个记录为准。达梦已收录包目前为 x64；未取得可验证的官方 Mac/ARM64 ODBC 包。

## 校验和验证等级

```bash
shasum -a 256 -c SHA256SUMS
python3 scripts/release/driver-downloads.py driver-downloads --verify
```

DBX 包检查清单、归档路径、配套文件摘要和实际 PE/ELF/Mach-O 架构。原始 vendor 包记录官方摘要匹配、文件结构和只读检查结果；不把下载完成等同于真实连接通过。Mac 可执行 smoke 会记录在 catalog 验证报告；Windows/Linux 跨编译文件未在对应系统运行，真实服务、TLS、认证和 SQL 方言仍需目标环境验收。
