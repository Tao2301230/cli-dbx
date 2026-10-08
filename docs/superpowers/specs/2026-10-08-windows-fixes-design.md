# Windows 修复迁移设计

用户已审阅 `cli-dbx_v0.1.5_source.zip` 的修复内容，并要求迁移可用部分、保持独立驱动路线。基线为本分支 `ac0e18d`；附件基线 `0401047` 仅作修复参考，不整体合并。

## 边界

- 主程序仍仅内置 MySQL、SQLite；四个 Go 驱动及 ODBC 适配程序仍为独立进程。根 `go.mod` 保持 Go 1.22，不引入厂商依赖。
- 保留 `driver_path`、扁平清单、协议 v1、凭据私有传输、方言 SQL 检查、导入行数限制和大整数精度保护。
- 保留本分支 Git 忽略规则；不导入旧版 `.gitignore`、内置 ODBC 命令或全内置模块。此次修改留在工作区，不自动 stage、commit 或发布。

## 修复

1. 增加标准库 `internal/textfile`，读取时只移除一个文件头 UTF-8 BOM。应用到 TOML、SQL 文件、事务 JSON、导入 CSV/JSON。正文中的 FEFF 保持原样；不自动转换 UTF-16。
2. 事务 JSON 在解析完计划后要求 EOF；第二份 JSON、尾随垃圾、未知字段均在执行前拒绝。JSON 导入继续使用 `UseNumber`、整数转换和整文档校验。
3. SQLite 普通文件名中的 `#`、`%`、中文和空格按真实文件名编码。采用转义文件路径后构造 `file:` URI 的方式，保留现有查询参数，尤其 `mode=ro`、`mode=memory`、`cache=shared` 和 `_pragma`。不能照搬附件的裸路径透传，否则会静默丢失只读语义。显式 `file:` URI 与内存配置维持现有行为。
4. Windows 测试同时隔离 `HOME`、`USERPROFILE`，写入 TOML 的路径按合法字符串转义。发布打包器显式指定主程序执行权限，不依赖宿主文件系统的 Unix 权限位。
5. 同步 CLI help、Windows 操作说明和独立驱动示例。主程序版本递增为 `v0.2.0-driver-process.4`，连接器最低版本同步；既有 `.3` 驱动包继续可由 `driver_path` 精确选用。

## 验证

先写回归并确认旧实现失败，再实现修复。覆盖 BOM/无 BOM、正文 FEFF、未知 JSON 字段、尾随内容、大整数入库、SQLite 特殊路径原文件重开、只读与内存行为、打包权限及重复打包字节一致性。最终运行根模块测试、嵌套驱动模块测试、Python 打包测试、Windows 交叉编译和本机 CLI 流程。交叉编译不表述为 Windows 实机通过。
