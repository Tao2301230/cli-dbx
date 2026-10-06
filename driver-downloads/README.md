# DBX 独立驱动下载入口

此目录的压缩包按引擎、版本、操作系统和架构平铺；index.json 可供下载服务或智能体读取。

选择 adapter=native 获取 Go 驱动；adapter=odbc 是 DBX 适配程序，还需匹配的 vendor 包。
当前版本包由未提交工作树构建，来源和日期见索引。历史快照保留原始字节。
`.1` 是旧引擎子目录布局，不适用于当前扁平加载规则；`.2`/`.3` 按 driver_path 精确选择。
厂商原始包原样保留；随包存在的许可证一并保留，缺失情况见索引。保存文件不等于实库验收通过。
OceanBase 库仅作为保存和移植材料，当前 DBX 未接入 OceanBase ODBC。
Windows 厂商 ODBC 使用正式安装器注册；Linux 包按发行版、libc、OpenSSL 和 unixODBC 依赖选用。

使用前执行 `shasum -a 256 -c SHA256SUMS`（Linux: `sha256sum -c SHA256SUMS`）；
Windows 可用 `Get-FileHash -Algorithm SHA256 <文件>` 与 index.json 比对。

详细移植、来源和平台限制见 PORTING.md 与 index.json。

| 文件 | 类别 | 平台/架构 | 大小 MiB |
|---|---|---|---:|
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_darwin_amd64.tar.gz](dbx-driver-dm-odbc_v0.2.0-driver-process.3_darwin_amd64.tar.gz) | dbx-driver | darwin/amd64 | 1.61 |
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_darwin_arm64.tar.gz](dbx-driver-dm-odbc_v0.2.0-driver-process.3_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 1.45 |
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-dm-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | dbx-driver | darwin/arm64 | 1.45 |
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_linux_amd64.tar.gz](dbx-driver-dm-odbc_v0.2.0-driver-process.3_linux_amd64.tar.gz) | dbx-driver | linux/amd64 | 1.46 |
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_linux_arm64.tar.gz](dbx-driver-dm-odbc_v0.2.0-driver-process.3_linux_arm64.tar.gz) | dbx-driver | linux/arm64 | 1.34 |
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_windows_amd64.zip](dbx-driver-dm-odbc_v0.2.0-driver-process.3_windows_amd64.zip) | dbx-driver | windows/amd64 | 1.30 |
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-dm-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | dbx-driver | windows/amd64 | 1.30 |
| [dbx-driver-dm_v0.2.0-driver-process.1_darwin_arm64.tar.gz](dbx-driver-dm_v0.2.0-driver-process.1_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 3.19 |
| [dbx-driver-dm_v0.2.0-driver-process.1_windows_amd64.zip](dbx-driver-dm_v0.2.0-driver-process.1_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.47 |
| [dbx-driver-dm_v0.2.0-driver-process.2_darwin_arm64.tar.gz](dbx-driver-dm_v0.2.0-driver-process.2_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 3.19 |
| [dbx-driver-dm_v0.2.0-driver-process.2_windows_amd64.zip](dbx-driver-dm_v0.2.0-driver-process.2_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.47 |
| [dbx-driver-dm_v0.2.0-driver-process.3_darwin_amd64.tar.gz](dbx-driver-dm_v0.2.0-driver-process.3_darwin_amd64.tar.gz) | dbx-driver | darwin/amd64 | 3.42 |
| [dbx-driver-dm_v0.2.0-driver-process.3_darwin_arm64.tar.gz](dbx-driver-dm_v0.2.0-driver-process.3_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 3.19 |
| [dbx-driver-dm_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-dm_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | dbx-driver | darwin/arm64 | 3.19 |
| [dbx-driver-dm_v0.2.0-driver-process.3_linux_amd64.tar.gz](dbx-driver-dm_v0.2.0-driver-process.3_linux_amd64.tar.gz) | dbx-driver | linux/amd64 | 3.42 |
| [dbx-driver-dm_v0.2.0-driver-process.3_linux_arm64.tar.gz](dbx-driver-dm_v0.2.0-driver-process.3_linux_arm64.tar.gz) | dbx-driver | linux/arm64 | 3.13 |
| [dbx-driver-dm_v0.2.0-driver-process.3_windows_amd64.zip](dbx-driver-dm_v0.2.0-driver-process.3_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.47 |
| [dbx-driver-dm_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-dm_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | dbx-driver | windows/amd64 | 3.47 |
| [dbx-driver-dm_v0.2.0-driver-process.3_windows_arm64.zip](dbx-driver-dm_v0.2.0-driver-process.3_windows_arm64.zip) | dbx-driver | windows/arm64 | 3.14 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_darwin_amd64.tar.gz](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_darwin_amd64.tar.gz) | dbx-driver | darwin/amd64 | 1.61 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_darwin_arm64.tar.gz](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 1.45 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | dbx-driver | darwin/arm64 | 1.45 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_linux_amd64.tar.gz](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_linux_amd64.tar.gz) | dbx-driver | linux/amd64 | 1.46 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_linux_arm64.tar.gz](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_linux_arm64.tar.gz) | dbx-driver | linux/arm64 | 1.34 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_windows_amd64.zip](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_windows_amd64.zip) | dbx-driver | windows/amd64 | 1.30 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | dbx-driver | windows/amd64 | 1.30 |
| [dbx-driver-oracle_v0.2.0-driver-process.1_darwin_arm64.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.1_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 2.89 |
| [dbx-driver-oracle_v0.2.0-driver-process.1_windows_amd64.zip](dbx-driver-oracle_v0.2.0-driver-process.1_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.22 |
| [dbx-driver-oracle_v0.2.0-driver-process.2_darwin_arm64.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.2_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 2.89 |
| [dbx-driver-oracle_v0.2.0-driver-process.2_windows_amd64.zip](dbx-driver-oracle_v0.2.0-driver-process.2_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.22 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_darwin_amd64.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.3_darwin_amd64.tar.gz) | dbx-driver | darwin/amd64 | 3.17 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_darwin_arm64.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.3_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 2.89 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | dbx-driver | darwin/arm64 | 2.89 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_linux_amd64.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.3_linux_amd64.tar.gz) | dbx-driver | linux/amd64 | 3.16 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_linux_arm64.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.3_linux_arm64.tar.gz) | dbx-driver | linux/arm64 | 2.83 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_windows_amd64.zip](dbx-driver-oracle_v0.2.0-driver-process.3_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.22 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-oracle_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | dbx-driver | windows/amd64 | 3.22 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_windows_arm64.zip](dbx-driver-oracle_v0.2.0-driver-process.3_windows_arm64.zip) | dbx-driver | windows/arm64 | 2.84 |
| [dbx-driver-postgres_v0.2.0-driver-process.1_darwin_arm64.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.1_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 3.12 |
| [dbx-driver-postgres_v0.2.0-driver-process.1_windows_amd64.zip](dbx-driver-postgres_v0.2.0-driver-process.1_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.44 |
| [dbx-driver-postgres_v0.2.0-driver-process.2_darwin_arm64.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.2_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 3.12 |
| [dbx-driver-postgres_v0.2.0-driver-process.2_windows_amd64.zip](dbx-driver-postgres_v0.2.0-driver-process.2_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.44 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_darwin_amd64.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.3_darwin_amd64.tar.gz) | dbx-driver | darwin/amd64 | 3.37 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_darwin_arm64.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.3_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 3.13 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | dbx-driver | darwin/arm64 | 3.13 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_linux_amd64.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.3_linux_amd64.tar.gz) | dbx-driver | linux/amd64 | 3.35 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_linux_arm64.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.3_linux_arm64.tar.gz) | dbx-driver | linux/arm64 | 3.03 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_windows_amd64.zip](dbx-driver-postgres_v0.2.0-driver-process.3_windows_amd64.zip) | dbx-driver | windows/amd64 | 3.44 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-postgres_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | dbx-driver | windows/amd64 | 3.44 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_windows_arm64.zip](dbx-driver-postgres_v0.2.0-driver-process.3_windows_arm64.zip) | dbx-driver | windows/arm64 | 3.06 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_darwin_amd64.tar.gz](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_darwin_amd64.tar.gz) | dbx-driver | darwin/amd64 | 1.61 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_darwin_arm64.tar.gz](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 1.45 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | dbx-driver | darwin/arm64 | 1.45 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_linux_amd64.tar.gz](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_linux_amd64.tar.gz) | dbx-driver | linux/amd64 | 1.46 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_linux_arm64.tar.gz](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_linux_arm64.tar.gz) | dbx-driver | linux/arm64 | 1.34 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_windows_amd64.zip](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_windows_amd64.zip) | dbx-driver | windows/amd64 | 1.30 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | dbx-driver | windows/amd64 | 1.30 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.1_darwin_arm64.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.1_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 2.59 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.1_windows_amd64.zip](dbx-driver-sqlserver_v0.2.0-driver-process.1_windows_amd64.zip) | dbx-driver | windows/amd64 | 2.75 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.2_darwin_arm64.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.2_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 2.59 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.2_windows_amd64.zip](dbx-driver-sqlserver_v0.2.0-driver-process.2_windows_amd64.zip) | dbx-driver | windows/amd64 | 2.75 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_darwin_amd64.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.3_darwin_amd64.tar.gz) | dbx-driver | darwin/amd64 | 2.78 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_darwin_arm64.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.3_darwin_arm64.tar.gz) | dbx-driver | darwin/arm64 | 2.59 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | dbx-driver | darwin/arm64 | 2.59 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_linux_amd64.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.3_linux_amd64.tar.gz) | dbx-driver | linux/amd64 | 2.71 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_linux_arm64.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.3_linux_arm64.tar.gz) | dbx-driver | linux/arm64 | 2.47 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_amd64.zip](dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_amd64.zip) | dbx-driver | windows/amd64 | 2.75 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | dbx-driver | windows/amd64 | 2.75 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_arm64.zip](dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_arm64.zip) | dbx-driver | windows/arm64 | 2.48 |
| [dbx-porting-kit_20261004.tar.gz](dbx-porting-kit_20261004.tar.gz) | porting-kit | source/all | 0.02 |
| [dbx-sources_v0.2.0-driver-process.3_worktree-20261004.tar.gz](dbx-sources_v0.2.0-driver-process.3_worktree-20261004.tar.gz) | dbx-source | source/all | 0.16 |
| [dm-navicat-dameng_odbc_linux.tar.gz](dm-navicat-dameng_odbc_linux.tar.gz) | vendor-installer | linux/amd64 | 10.80 |
| [dm-navicat-dameng_odbc_win.zip](dm-navicat-dameng_odbc_win.zip) | vendor-installer | windows/amd64 | 8.99 |
| [libtool-source_2.4.3.tar.xz](libtool-source_2.4.3.tar.xz) | runtime-source | source/all | 0.90 |
| [libtool-source_2.6.2.tar.xz](libtool-source_2.6.2.tar.xz) | runtime-source | source/all | 1.04 |
| [oceanbase-ob-connector-odbc-source.tar.gz](oceanbase-ob-connector-odbc-source.tar.gz) | vendor-source | source/all | 2.31 |
| [oceanbase-obconnector-c-source.tar.gz](oceanbase-obconnector-c-source.tar.gz) | vendor-source | source/all | 0.92 |
| [oceanbase-odbc-bundle_2.0.8.3_darwin_arm64.tar.gz](oceanbase-odbc-bundle_2.0.8.3_darwin_arm64.tar.gz) | vendor-bundle | darwin/arm64 | 3.61 |
| [oracle-basic_19.16.0.0.0_macos_x64.dmg](oracle-basic_19.16.0.0.0_macos_x64.dmg) | vendor-installer | darwin/amd64 | 81.33 |
| [oracle-basic_23.26.2.0.0_macos_arm64.dmg](oracle-basic_23.26.2.0.0_macos_arm64.dmg) | vendor-installer | darwin/arm64 | 130.50 |
| [oracle-basic_23.26.3.0.0_linux_arm64.zip](oracle-basic_23.26.3.0.0_linux_arm64.zip) | vendor-installer | linux/arm64 | 121.99 |
| [oracle-basic_23.26.3.0.0_linux_x64.zip](oracle-basic_23.26.3.0.0_linux_x64.zip) | vendor-installer | linux/amd64 | 129.81 |
| [oracle-basic_23.26.3.0.0_windows_x64.zip](oracle-basic_23.26.3.0.0_windows_x64.zip) | vendor-installer | windows/amd64 | 137.13 |
| [oracle-odbc-bundle_23.26.2.0.0_darwin_arm64.tar.gz](oracle-odbc-bundle_23.26.2.0.0_darwin_arm64.tar.gz) | vendor-bundle | darwin/arm64 | 117.16 |
| [oracle-odbc_19.16.0.0.0_macos_x64.dmg](oracle-odbc_19.16.0.0.0_macos_x64.dmg) | vendor-installer | darwin/amd64 | 2.30 |
| [oracle-odbc_23.26.2.0.0_macos_arm64.dmg](oracle-odbc_23.26.2.0.0_macos_arm64.dmg) | vendor-installer | darwin/arm64 | 1.01 |
| [oracle-odbc_23.26.3.0.0_linux_arm64.zip](oracle-odbc_23.26.3.0.0_linux_arm64.zip) | vendor-installer | linux/arm64 | 0.79 |
| [oracle-odbc_23.26.3.0.0_linux_x64.zip](oracle-odbc_23.26.3.0.0_linux_x64.zip) | vendor-installer | linux/amd64 | 0.78 |
| [oracle-odbc_23.26.3.0.0_windows_x64.zip](oracle-odbc_23.26.3.0.0_windows_x64.zip) | vendor-installer | windows/amd64 | 1.01 |
| [sqlserver-msodbcsql17-17.11.1.1-macos-amd64.tar.gz](sqlserver-msodbcsql17-17.11.1.1-macos-amd64.tar.gz) | vendor-installer | darwin/amd64 | 0.80 |
| [sqlserver-msodbcsql17-17.11.1.1-macos-arm64.tar.gz](sqlserver-msodbcsql17-17.11.1.1-macos-arm64.tar.gz) | vendor-installer | darwin/arm64 | 0.72 |
| [sqlserver-msodbcsql17-17.11.1.1-rhel9-x86_64.rpm](sqlserver-msodbcsql17-17.11.1.1-rhel9-x86_64.rpm) | vendor-installer | linux/amd64 | 0.77 |
| [sqlserver-msodbcsql17-17.11.1.1-ubuntu24.04-amd64.deb](sqlserver-msodbcsql17-17.11.1.1-ubuntu24.04-amd64.deb) | vendor-installer | linux/amd64 | 0.74 |
| [sqlserver-msodbcsql17-17.11.1.1-windows-x64.msi](sqlserver-msodbcsql17-17.11.1.1-windows-x64.msi) | vendor-installer | windows/amd64 | 4.46 |
| [sqlserver-msodbcsql18-18.7.1.1-macos-amd64.tar.gz](sqlserver-msodbcsql18-18.7.1.1-macos-amd64.tar.gz) | vendor-installer | darwin/amd64 | 0.83 |
| [sqlserver-msodbcsql18-18.7.1.1-macos-arm64.tar.gz](sqlserver-msodbcsql18-18.7.1.1-macos-arm64.tar.gz) | vendor-installer | darwin/arm64 | 0.74 |
| [sqlserver-msodbcsql18-18.7.1.1-rhel9-aarch64.rpm](sqlserver-msodbcsql18-18.7.1.1-rhel9-aarch64.rpm) | vendor-installer | linux/arm64 | 0.85 |
| [sqlserver-msodbcsql18-18.7.1.1-rhel9-x86_64.rpm](sqlserver-msodbcsql18-18.7.1.1-rhel9-x86_64.rpm) | vendor-installer | linux/amd64 | 0.92 |
| [sqlserver-msodbcsql18-18.7.1.1-ubuntu24.04-amd64.deb](sqlserver-msodbcsql18-18.7.1.1-ubuntu24.04-amd64.deb) | vendor-installer | linux/amd64 | 0.76 |
| [sqlserver-msodbcsql18-18.7.1.1-ubuntu24.04-arm64.deb](sqlserver-msodbcsql18-18.7.1.1-ubuntu24.04-arm64.deb) | vendor-installer | linux/arm64 | 0.65 |
| [sqlserver-msodbcsql18-18.7.1.1-windows-x64.msi](sqlserver-msodbcsql18-18.7.1.1-windows-x64.msi) | vendor-installer | windows/amd64 | 9.62 |
| [sqlserver-odbc-bundle_18.7.1.1_darwin_arm64.tar.gz](sqlserver-odbc-bundle_18.7.1.1_darwin_arm64.tar.gz) | vendor-bundle | darwin/arm64 | 0.74 |
| [unixODBC-debian12-libodbc2_2.3.11-2+deb12u1_amd64.deb](unixODBC-debian12-libodbc2_2.3.11-2+deb12u1_amd64.deb) | build-dependency | linux/amd64 | 0.14 |
| [unixODBC-debian12-libodbc2_2.3.11-2+deb12u1_arm64.deb](unixODBC-debian12-libodbc2_2.3.11-2+deb12u1_arm64.deb) | build-dependency | linux/arm64 | 0.13 |
| [unixODBC-debian12-unixodbc-dev_2.3.11-2+deb12u1_amd64.deb](unixODBC-debian12-unixodbc-dev_2.3.11-2+deb12u1_amd64.deb) | build-dependency | linux/amd64 | 0.22 |
| [unixODBC-debian12-unixodbc-dev_2.3.11-2+deb12u1_arm64.deb](unixODBC-debian12-unixodbc-dev_2.3.11-2+deb12u1_arm64.deb) | build-dependency | linux/arm64 | 0.20 |
| [unixODBC-source_2.3.14.tar.gz](unixODBC-source_2.3.14.tar.gz) | runtime-source | source/all | 1.68 |

## 未提供的组合

- dm windows/amd64: Dameng official public download page currently lists standalone JDBC, DotNet and R2DBC drivers only; no standalone ODBC package is listed. Official C/C++ documentation instructs obtaining ODBC from a full DM8 database installation (drivers/odbc or bin/libdodbc.so). Full database installer was not downloaded for this driver collection.
- dm linux/amd64: Dameng official public download page currently lists standalone JDBC, DotNet and R2DBC drivers only; no standalone ODBC package is listed. Official C/C++ documentation instructs obtaining ODBC from a full DM8 database installation (drivers/odbc or bin/libdodbc.so). Full database installer was not downloaded for this driver collection.
- dm darwin/arm64: No macOS platform or macOS ODBC archive is listed in the official public DM8 download catalog inspected 2026-10-04. Official ODBC documentation describes Windows and Linux and supplies no macOS ODBC download. This records absence of a verified public first-party package, not a universal claim about support.
- dm darwin/amd64: No macOS platform or macOS ODBC archive is listed in the official public DM8 download catalog inspected 2026-10-04. Official ODBC documentation describes Windows and Linux and supplies no macOS ODBC download. This records absence of a verified public first-party package, not a universal claim about support.
- oracle windows/arm64: 当前 DBX ODBC 适配器只构建 Windows x64，未声明 Windows ARM64 支持。
- dm windows/arm64: 当前 DBX ODBC 适配器只构建 Windows x64，未声明 Windows ARM64 支持。
- sqlserver windows/arm64: 当前 DBX ODBC 适配器只构建 Windows x64，未声明 Windows ARM64 支持。
- dm windows/arm64: 本次只取得 Navicat 分发的 x64 达梦 ODBC 包；未取得经验证的 ARM64 包。
- dm linux/arm64: 本次只取得 Navicat 分发的 x64 达梦 ODBC 包；未取得经验证的 ARM64 包。
