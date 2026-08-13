# CURRENT_TASK.md — 当前任务

## 正在执行

**PHASE A — Bootstrap — DONE ✅**

全部 A01-A11 完成,main CI 全绿(run 31722439468:test/sdk/frontend/desktop success)。

## 下一步

**PHASE B — Temper 身份与数据隔离(milestone/01-identity-isolation)**

1. 创建分支 `milestone/01-identity-isolation`
2. 用户可见身份:Temper / 0.3.0-dev / Shape intent. Ship work. / Temper.exe
3. Runtime Home 隔离:%APPDATA%\Temper\runtime\、%LOCALAPPDATA%\Temper\cache\、%APPDATA%\Temper\cowork\
4. 在 Reasonix Boot/Config load 前注入 REASONIX_HOME/STATE/CACHE
5. 测试两个不同 Home 隔离
6. 网络审计:telemetry DISABLE、updater REPLACE/DISABLE、release URLs REPLACE
7. 测试 + commit + push + CI + PR + merge
