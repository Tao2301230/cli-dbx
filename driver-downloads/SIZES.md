# DBX 独立驱动体积

单位 MiB（1 MiB = 1,048,576 字节）。每格为 **可执行文件 / 完整压缩包**；Mac/Linux 使用 tar.gz，Windows 使用 zip。ODBC Mac 包包含 unixODBC/libltdl、清单和许可文件；Linux ODBC 包使用系统 libodbc.so.2，厂商库均另行保存。

## 当前 .3

| 引擎/路线 | Mac Intel | Mac ARM64 | Win x64 | Win ARM64 | Linux x64 | Linux ARM64 |
|---|---:|---:|---:|---:|---:|---:|
| postgres / native | 9.34 / 3.37 | 8.79 / 3.13 | 9.47 / 3.44 | 8.72 / 3.06 | 9.24 / 3.35 | 8.69 / 3.03 |
| oracle / native | 7.74 / 3.17 | 7.15 / 2.89 | 7.85 / 3.22 | 7.13 / 2.84 | 7.69 / 3.16 | 7.13 / 2.83 |
| dm / native | 8.20 / 3.42 | 7.65 / 3.19 | 8.30 / 3.47 | 7.59 / 3.14 | 8.15 / 3.42 | 7.56 / 3.13 |
| sqlserver / native | 6.86 / 2.78 | 6.37 / 2.59 | 6.77 / 2.75 | 6.17 / 2.48 | 6.64 / 2.71 | 6.19 / 2.47 |
| oracle / odbc | 3.20 / 1.61 | 2.98 / 1.45 | 3.08 / 1.30 | 未提供 | 3.44 / 1.46 | 3.22 / 1.34 |
| dm / odbc | 3.20 / 1.61 | 2.98 / 1.45 | 3.08 / 1.30 | 未提供 | 3.44 / 1.46 | 3.22 / 1.34 |
| sqlserver / odbc | 3.20 / 1.61 | 2.98 / 1.45 | 3.08 / 1.30 | 未提供 | 3.44 / 1.46 | 3.22 / 1.34 |

## 历史版本和重建前快照

每个原始包单独保存；.1 为旧子目录布局，当前加载器不直接支持。snapshot 不表示独立版本号，同版本包内文件名相同。

| 文件 | 可执行 MiB | 压缩包 MiB |
|---|---:|---:|
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-dm-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | 2.98 | 1.45 |
| [dbx-driver-dm-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-dm-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | 3.08 | 1.30 |
| [dbx-driver-dm_v0.2.0-driver-process.1_darwin_arm64.tar.gz](dbx-driver-dm_v0.2.0-driver-process.1_darwin_arm64.tar.gz) | 7.62 | 3.19 |
| [dbx-driver-dm_v0.2.0-driver-process.1_windows_amd64.zip](dbx-driver-dm_v0.2.0-driver-process.1_windows_amd64.zip) | 8.30 | 3.47 |
| [dbx-driver-dm_v0.2.0-driver-process.2_darwin_arm64.tar.gz](dbx-driver-dm_v0.2.0-driver-process.2_darwin_arm64.tar.gz) | 7.62 | 3.19 |
| [dbx-driver-dm_v0.2.0-driver-process.2_windows_amd64.zip](dbx-driver-dm_v0.2.0-driver-process.2_windows_amd64.zip) | 8.30 | 3.47 |
| [dbx-driver-dm_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-dm_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | 7.62 | 3.19 |
| [dbx-driver-dm_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-dm_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | 8.30 | 3.47 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | 2.98 | 1.45 |
| [dbx-driver-oracle-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-oracle-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | 3.08 | 1.30 |
| [dbx-driver-oracle_v0.2.0-driver-process.1_darwin_arm64.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.1_darwin_arm64.tar.gz) | 7.15 | 2.89 |
| [dbx-driver-oracle_v0.2.0-driver-process.1_windows_amd64.zip](dbx-driver-oracle_v0.2.0-driver-process.1_windows_amd64.zip) | 7.85 | 3.22 |
| [dbx-driver-oracle_v0.2.0-driver-process.2_darwin_arm64.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.2_darwin_arm64.tar.gz) | 7.15 | 2.89 |
| [dbx-driver-oracle_v0.2.0-driver-process.2_windows_amd64.zip](dbx-driver-oracle_v0.2.0-driver-process.2_windows_amd64.zip) | 7.85 | 3.22 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-oracle_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | 7.15 | 2.89 |
| [dbx-driver-oracle_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-oracle_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | 7.85 | 3.22 |
| [dbx-driver-postgres_v0.2.0-driver-process.1_darwin_arm64.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.1_darwin_arm64.tar.gz) | 8.79 | 3.12 |
| [dbx-driver-postgres_v0.2.0-driver-process.1_windows_amd64.zip](dbx-driver-postgres_v0.2.0-driver-process.1_windows_amd64.zip) | 9.46 | 3.44 |
| [dbx-driver-postgres_v0.2.0-driver-process.2_darwin_arm64.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.2_darwin_arm64.tar.gz) | 8.79 | 3.12 |
| [dbx-driver-postgres_v0.2.0-driver-process.2_windows_amd64.zip](dbx-driver-postgres_v0.2.0-driver-process.2_windows_amd64.zip) | 9.46 | 3.44 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-postgres_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | 8.79 | 3.13 |
| [dbx-driver-postgres_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-postgres_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | 9.47 | 3.44 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | 2.98 | 1.45 |
| [dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-sqlserver-odbc_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | 3.08 | 1.30 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.1_darwin_arm64.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.1_darwin_arm64.tar.gz) | 6.37 | 2.59 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.1_windows_amd64.zip](dbx-driver-sqlserver_v0.2.0-driver-process.1_windows_amd64.zip) | 6.77 | 2.75 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.2_darwin_arm64.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.2_darwin_arm64.tar.gz) | 6.37 | 2.59 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.2_windows_amd64.zip](dbx-driver-sqlserver_v0.2.0-driver-process.2_windows_amd64.zip) | 6.77 | 2.75 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz](dbx-driver-sqlserver_v0.2.0-driver-process.3_darwin_arm64_snapshot-20261004-before-rebuild.tar.gz) | 6.37 | 2.59 |
| [dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip](dbx-driver-sqlserver_v0.2.0-driver-process.3_windows_amd64_snapshot-20261004-before-rebuild.zip) | 6.77 | 2.75 |
