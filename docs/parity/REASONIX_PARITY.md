# REASONIX_PARITY.md — Reasonix 功能 Parity 追踪

> PHASE C。逐项核对 Reasonix 能力在 Temper 中的继承/接线状态。
> 状态枚举:UPSTREAM_AVAILABLE(存在,未接线)/ TEMPER_WIRED(已接线)/ ADVANCED_ONLY / CLI_ONLY / NOT_WIRED / OUT_OF_SCOPE_WITH_REASON。

## 追踪表

| # | 能力 | Upstream owner | Upstream test | Upstream desktop surface | Temper surface | Adapter | Automated test | Manual smoke | Status |
|---|------|----------------|---------------|--------------------------|----------------|---------|----------------|--------------|--------|
| 1 | Providers / Models | internal/provider(11 files) | internal/provider/*_test.go | Settings→Providers;模型切换 | Settings→Providers(继承) | —(直接继承) | go test ./internal/provider/ | 待 PHASE M | TEMPER_WIRED |
| 2 | Agent / Controller | internal/agent(125), internal/control | 大量 * _test.go | 主 Chat 运行 | Chat 运行(继承) | — | go test ./internal/agent/ | 待 PHASE M | TEMPER_WIRED |
| 3 | Plan | internal/plancontract, desktop 前端 | plancontract tests | Plan 展示/切换 | Chat/Work Plan 展示(继承) | — | go test ./... | 待 PHASE J | TEMPER_WIRED |
| 4 | Goal | internal/goaleval, internal/taskintent | goaleval tests | Goal 设置/展示 | Chat Goal(继承) | — | go test ./internal/goaleval/ | 待 PHASE J | TEMPER_WIRED |
| 5 | Todo | internal/taskcatalog, TodoPanel.tsx | taskcatalog tests | TodoPanel | TodoPanel(继承) | — | go test ./internal/taskcatalog/ | 待 PHASE J | TEMPER_WIRED |
| 6 | Planner / Subagent | internal/agent(agent.go, capability_gate.go), desktop/subagents_app.go | subagents_app tests | Subagent 管理 | Subagent(继承) | — | go test ./... | 待 PHASE C smoke | TEMPER_WIRED |
| 7 | Context / Compaction / Cache | internal/agent(agent.go, compact.go), desktop/context_maintenance.go | agent/compact tests | ContextPanel/CompactRatio | ContextPanel(继承) | — | go test ./... | 待 PHASE K | TEMPER_WIRED |
| 8 | History / Memory | internal/history, internal/memory(17) | history/memory tests | MemoryPanel/MemoryCitations | MemoryPanel(继承) | — | go test ./internal/memory/ | 待 PHASE J | TEMPER_WIRED |
| 9 | Tools | internal/tool(registry), builtin | tool tests | Tool 卡片/审批 | Chat Tool 卡片(继承) | — | go test ./internal/tool/ | 待 PHASE G | TEMPER_WIRED |
| 10 | Permissions / Sandbox | internal/permission(5), internal/sandbox(10) | permission/sandbox tests | 审批 Modal/ApprovalModal | 审批(继承) | — | go test ./internal/permission/ | 待 PHASE J | TEMPER_WIRED |
| 11 | MCP | internal/mcpregistry, internal/mcplaunch | mcpregistry tests | MCP Settings | MCP Settings(继承) | — | go test ./... | 待 PHASE Q smoke | TEMPER_WIRED |
| 12 | Skills / Plugins / Extensions | internal/skill(7), internal/plugin(25), internal/extension(28) | skill/plugin/extension tests | Skill/Plugin/Extension 面板 | 继承 | — | go test ./... | 待 PHASE Q smoke | TEMPER_WIRED |
| 13 | Checkpoint / Rewind / Recovery | internal/checkpoint(13), internal/recovery(11) | checkpoint/recovery tests | UndoRewindBanner/RecoveryLineageDialog | 继承 | — | go test ./internal/recovery/ | 待 PHASE M | TEMPER_WIRED |
| 14 | Serve | internal/serve(7) | serve tests | Server 面板 | Advanced→Serve(继承) | — | go test ./internal/serve/ | 待 PHASE J | ADVANCED_ONLY |
| 15 | ACP | internal/acp(10) | acp tests | 无桌面组件(0 files) | Advanced→ACP(待接线) | 需 adapter | go test ./internal/acp/ | — | NOT_WIRED |
| 16 | Remote SSH | internal/remote(11), RemoteHostsPage.tsx | remote tests | RemoteHostsPage/RemotePanel | 继承 | — | go test ./internal/remote/ | 待 PHASE J | TEMPER_WIRED |
| 17 | Diagnostics | internal/capdiag, DiagnosticsSettingsPage.tsx | capdiag tests | DiagnosticsSettingsPage | 继承 | — | go test ./... | 待 PHASE K | TEMPER_WIRED |

## 结论

- Reasonix 运行时能力全部继承,无需再造第二套 Runtime(符合 Master 架构约束)。
- 16/17 项已由 Reasonix 桌面 UI 直接覆盖(组件已存在),Temper 复用;其中 Serve 归类 ADVANCED_ONLY。
- **ACP(1 项)** 无桌面 UI 绑定,需 Temper adapter 接线到 Advanced 页 → PHASE J 处理。
- 全部能力已通过 upstream 测试覆盖;Temper 的回归验证依赖 `go test ./...` 与 desktop 测试(CI 已跑)。

## 待办

- [ ] ACP:在 Advanced 页接入 ACP 面板(Temper adapter)
- [ ] 每项能力在 PHASE M(Production E2E)补 Manual smoke
