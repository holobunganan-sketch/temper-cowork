# THIRD_PARTY_NOTICES.md

## Reasonix(上游基线)

- Upstream: https://github.com/esengine/DeepSeek-Reasonix.git
- Branch: main-v2
- Baseline SHA: 49f24d19702c9542ab50500d590237dc872c4d58
- License: MIT (见 LICENSE)

Temper incorporates and modifies Reasonix components (Reasonix Runtime) as its execution baseline. Reasonix is a deepseek-native AI coding agent with a runtime covering Provider / Agent / Controller / Plan / Goal / Todo / Planner / Executor / Subagent / Context / Memory / Tools / Permissions / Sandbox / MCP / Skills / Plugins / Extensions / Checkpoint / Rewind / Recovery / Serve / ACP / Remote SSH.

## 其他第三方依赖

以 go.mod / go.sum 与 desktop/frontend/package.json / pnpm-lock.yaml 为准。Reasonix 上游 MIT 许可范围内的第三方库列表见其自身 LICENSE 与依赖声明;Temper 保留相同归属。

## 说明

- Temper 不包含 Reasonix 的发布/网站/R2/npm/SignPath 基础设施。
- Temper 桌面 UI 为独立设计,不复制 Reasonix 视觉。
- 本文件随构建持续更新;新引入第三方组件时须在此登记。
