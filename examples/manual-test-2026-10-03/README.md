# 2026-10-03 手动测试配套文件

按 [独立驱动版手册](../../docs/数据库连接手动测试手册-2026-10-03-独立驱动版.md) 将整个目录复制到新的测试工作目录。

- connections：可直接用的SQLite配置。
- templates：远程连接模板；Windows按手册补.exe；编辑后再复制到connections。
- policy：逐个用单文件--config执行的配置负例及策略样例。
- fixtures / sql / plans：测试数据和SQL/事务输入。
- regressions：已知失败的合成配置，只允许指定dry-run/一次性SQLite测试。
- priority：配置来源优先级样例。
- sqlite-cases.json / regression-cases.json：手册命令及预期。
- driver-negative.py：在临时副本校验7类错误驱动包，Python3可选工具。
- verification：本地测试生成的运行记录和导出结果，不随源码提交。

不要提交后续填写的真实连接密码/密文，也不要将测试数据文件加入源码。
