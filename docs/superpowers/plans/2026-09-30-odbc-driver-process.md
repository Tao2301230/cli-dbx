# ODBC 独立驱动实施计划与验收

依据：[设计说明](../specs/2026-09-30-odbc-driver-process-design.md)。分支 `codex/dbx-driver-process`，从已有扁平驱动版本 `4f7d624` 继续，原实验基线为 `86b18d3ce5241148f132af211753c3c2de69b071`。

## 实施步骤

- [x] 连接配置增加显式 ODBC 路线和厂商库路径，校验选项，凭据通过私有管道传递。
- [x] 增加 Oracle、达梦、SQL Server 独立 ODBC 模块，保留原有 Go 路线。
- [x] ODBC 使用问号参数，保留数据库方言、动作权限、事务和行数限制。
- [x] 清单增加适配类型和配套文件摘要校验，平铺文件且版本互不覆盖。
- [x] Mac 打包 unixODBC/libltdl，重定位并签名副本；Windows 使用系统管理器。
- [x] 同步配置示例、帮助、安装文档和版本号。
- [x] 完成独立审查，修复精确数值、二进制类型、未知影响行数、目标别名和自动重试选项问题。
- [x] 最终重建、回归与本地提交。

## 已执行验证

- 根模块完整测试、独立模块测试；主程序依赖不包含 ODBC 或四种外置厂商驱动。
- Mac ARM64 / Windows x64 主程序、Platform 包、四种 Go 驱动及三种 ODBC 驱动构建。
- Mac ODBC 归档平铺 11 个文件，运行库和许可证摘要相同，可共存；Windows 归档平铺两个文件。
- Mac 解压到含空格路径后移动目录，移除 DYLD 环境覆盖，通过真实 unixODBC 加载测试用 SQLite ODBC 库，完成中文查询、ASCII/空值/非 UTF-8 二进制列保真、参数化导入、事务提交/失败回滚、超限导入回滚和多语句拒绝。
- 测试用 SQLite ODBC 源码的旧构建脚本不识别 ARM64，临时更新 config.guess/config.sub；其连接串解析不支持 DATABASE 花括号，临时测试副本增加去除外围花括号的处理。测试依赖未进入产品或发布包，DBX 的连接串转义保持完整。

## 验证边界

替身测试验证进程协议、ODBC 调用和控制面，不证明 Oracle、达梦、SQL Server 实际服务、厂商库、TLS 或认证兼容性。Windows 只完成交叉构建和归档检查，未运行 Windows 厂商库或 PowerShell 构建脚本；Linux 未在本机验证。厂商客户端不随包分发，仍需按厂商要求安装其资源和依赖。

- 最终额外运行 ODBC fork 的独立测试：SQL 类型元数据、高精度数字文本、未知/溢出影响行数均通过；SQLMoreResults 错误继续上抛。进程和共享数据库层 race 测试通过。
