# ARCHITECTURE.md — Temper 架构

## 固定架构

```text
Temper Desktop
React + TypeScript + Wails
        │
Temper CoWork Layer   (internal/temper/**)
Project / Chat / Work
Task Contract
Evidence / Decision
Artifact / Validation
Quality / Completion
        │
Reasonix Runtime       (继承复用,禁止再造)
Provider / Agent / Controller
Plan / Goal / Todo
Planner / Executor / Subagent
Context / Compaction / Cache
Memory / History
Tools / Permission / Sandbox
MCP / Skills / Plugins / Extensions
Checkpoint / Rewind / Recovery
Serve / ACP / Remote SSH
```

## 数据位置

- Runtime Home:`%APPDATA%\Temper\runtime\`
- Cache:`%LOCALAPPDATA%\Temper\cache\`
- CoWork DB:`%APPDATA%\Temper\cowork\temper.db`(SQLite, WAL)
- 通过 REASONIX_HOME / REASONIX_STATE_HOME / REASONIX_CACHE_HOME 注入隔离。

## 集成原则

1. Reasonix 已有 → 找源码/测试/Desktop 调用 → 复用。
2. 接口不适合 Temper UI → 建 Temper adapter(internal/temper/**)。
3. Reasonix 缺 CoWork 语义 → 写 internal/temper/**。
4. 只有 adapter 无法解决 → 最小修改 Reasonix-owned source + regression test + 登记 REASONIX_PATCHES.md。

## CoWork 数据模型(Temper-owned)

- projects / works / work_events / artifacts / evidence / decisions / acceptance_results / quality_runs
- 详见 PHASE D 与 internal/temper/**。
