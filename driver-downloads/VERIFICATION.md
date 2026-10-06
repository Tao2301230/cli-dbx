# 下载目录验证记录（2026-10-04）

来源基线 86b18d3ce5241148f132af211753c3c2de69b071，当前 .3 从未提交工作树重建，包含已有修复；不是另一个 Git 提交。此工作未提交、暂存或推送 Git，也未安装厂商驱动或启动数据库服务。HEAD 与原暂存树未变，原有文件除新增忽略 driver-downloads/ 的 .gitignore 外全部摘要一致。

## 已保存

- 最新 Go 驱动24包：PostgreSQL、Oracle、DM、SQLServer × darwin/linux/windows × amd64/arm64。
- 最新 ODBC适配器15包：Oracle、DM、SQLServer × MacIntel/MacARM/LinuxAMD64/LinuxARM64/Windowsx64。
- 历史DBX包30：.1八包、.2八包、重建前.3十四包；历史字节保留。
- 厂商原始包24：Oracle10，Microsoft12，Navicat分发DM2。
- 完整本地库副本3：Oracle23.26.2、SQLServer18.7.1.1、OceanBase2.0.8.3 MacARM64；内部相对符号链接及资源保留。
- OceanBase源码2、unixODBC/libtool运行库源码3、Debian12构建依赖4、DBX源码快照1、移植工具包1。

合计107个制品。文件准确字节数、SHA256和来源见index.json；体积对比见SIZES.md。许可证副本、说明和索引也计入全目录SHA256SUMS。

## 验证范围

- 目录工具6项测试通过，包含结构字段防覆盖、索引完整性/一致性、清单名称、平台和摘要检查。Go打包器和进程模块测试通过。
- 所有69个DBX压缩包均检查归档路径、清单、程序和资源摘要、PE/ELF/Mach-O平台架构。独立审查对89个源码/厂商/DBX tar/zip包检查路径、重复成员和链接逃逸，未发现问题；新增源码和工具包逐成员字节/路径验证。
- 36个Mac ODBC配套库架构检查通过；Mac适配包加载路径仅使用同目录或系统库，代码签名校验通过。最新Intel3个适配器解压后通过Rosetta启动及协议错误响应检查。
- 最新MacARM64四个Go驱动由新编CLI实际加载，验证相对路径含空格、子进程协议、密码脱敏和预期连接错误；目标是本地关闭端口，不代表厂商数据库连接验收。
- 最新MacARM64 ODBC运行库解压到新目录后，SQLServer适配器使用真实SQLite ODBC库验证连接、中文、二进制/空二进制、绑定导入、事务提交/回滚、导入超限回滚、策略限制和脱敏。Oracle/DM验证程序启动和缺少vendor库的受控失败。SQLite替身不验证SQLServer/Oracle/DM的真实数据库兼容性。
- Oracle官方包全部与发布SHA256/大小匹配，ZIP CRC及只读DMG文件结构检查通过。Microsoft Mac包与官方SHA256匹配，其余包检查格式和厂商平台/版本信息；未运行安装器。DM包来源为Navicat官方支持页面，Windows/Linux实际库为x64；下载URL未发布完整版本。
- Linux ODBC交叉构建ELF解释器、共享依赖和GLIBC版本检查通过，要求glibc>=2.34及匹配libodbc.so.2。Windows/Linux未在目标系统运行，未连接任何真实厂商数据库。

## 尚未提供或需匹配的材料

WindowsARM64 ODBC适配器未提供；Go驱动ARM64包不能代替它。DM未取得验证过的Mac/ARM64厂商ODBC包；已保存的Navicat包没有独立许可文件，Linux还需匹配libcrypto.so.3，不能据保存状态认定可公开再分发或直接运行。OracleBasic/ODBC须选同平台同版本，Windows厂商库需安装注册。OceanBase保留完整构建及源码材料，当前DBX不开放该ODBC连接。
