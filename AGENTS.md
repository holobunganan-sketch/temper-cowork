# Temper — Agent Guide

Temper 是构建在 Reasonix Runtime 之上的 Windows 桌面 CoWork 应用。

## 项目身份

- 产品:Temper
- 版本:0.3.0-dev
- 标语:Shape intent. Ship work.
- 仓库:git@github.com:holobunganan-sketch/temper-cowork.git (origin)
- Reasonix upstream:https://github.com/esengine/DeepSeek-Reasonix.git (upstream, 仅 fetch/compare/reference)

## 最高规范

- 唯一最高工程规范:docs/spec/TEMPER_V0.3.0_MASTER_BUILD.md
- 禁止重新设计产品;禁止再造第二套 Runtime;禁止 merge Reasonix Git 历史。

## 施工纪律(DeepSeek V4 Flash)

- 单个 Micro Task 目标上下文 <= 50K tokens,硬上限 <= 90K
- Production files <= 6 / Test files <= 4 / Docs files <= 3
- 先读 EXECUTOR_RULES.md,再读 BUILD_STATE.md / CURRENT_TASK.md

## 常用命令

- Go 测试:`go test ./...`(根模块);`cd desktop && go test -short .`
- lint:`make lint`(需 golangci-lint v2.12.2)
- Frontend:`cd desktop/frontend && pnpm typecheck && pnpm build`
- Wails build:`cd desktop && wails build`
- 环境变量:Windows 本地测试需 `GIT_CEILING_DIRECTORIES=C:\Users\ZhouNan` 与 `TMP/TEMP` 指向无 .git 祖先目录;前端测试需 NODE_OPTIONS 注入 en-US locale(见 EXECUTOR_RULES.md)

## Git 规则

- origin = Temper,upstream = Reasonix,禁止弄反
- 禁止:git reset --hard / git clean -fdx / git push --force / git tag -f
- 开发用 milestone/<number>-<slug> 分支 + PR + CI + merge
