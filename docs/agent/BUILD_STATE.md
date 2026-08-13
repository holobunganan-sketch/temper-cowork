# BUILD_STATE.md — Temper v0.3.0 构建状态

> 更新规则:每个 Milestone / Micro Task 完成后更新本文件,记录真实进度与验证证据。

## 当前阶段

**PHASE A — Bootstrap(进行中)**

## 里程碑进度

| Phase | 内容 | 状态 | 备注 |
|-------|------|------|------|
| A | Bootstrap(基线导入/控制文件/CI/首推) | IN_PROGRESS | Reasonix baseline 已验证通过 |
| B | 身份与数据隔离 | PENDING | |
| C | Reasonix 功能 Parity | PENDING | |
| D | CoWork Store | PENDING | |
| E | Project + Chat | PENDING | |
| F | Formal Work | PENDING | |
| G | CoWork Tools / Evidence / Decision | PENDING | |
| H | Artifact / Quality | PENDING | |
| I | Temper UI 基础 | PENDING | |
| J | Chat / Work / Advanced UI | PENDING | |
| K | Runtime Observability | PENDING | |
| L | i18n / Windows UX | PENDING | |
| M | Production E2E | PENDING | |
| N | Defender | PENDING | |
| O | MSIX Packaging | PENDING | |
| P | GitHub MSIX Release | PENDING | |
| Q | RC / v0.3.0 正式 Release | PENDING | |

## 关键事实

- REASONIX_BASELINE_SHA = `49f24d19702c9542ab50500d590237dc872c4d58`(main-v2, 2026-08-13)
- Wails v2.13.0;Go 1.26.5;Node 26.7.0;pnpm 10.34.5
- 本地工具链位于 `C:\Myfolder\.toolchain\`(Go / Node)
- GitHub:origin = holobunganan-sketch/temper-cowork(空仓库,含 1 个占位提交 d66318b)
- 本地已有 .git(上一轮链接仓库时创建),继续沿用

## 已验证(带证据)

- [x] Reasonix `go test ./...`(环境修正后) — 除 Windows symlink 特权相关测试外全绿
- [x] `go vet ./...` — 通过
- [x] golangci-lint v2.12.2 — 0 issues
- [x] desktop `go test -short .` — 通过(217s)
- [x] frontend typecheck / CI 测试脚本 / build — 通过
- [x] `wails build` — 成功,产出 reasonix-desktop.exe

## 已知偏差与风险

- 远端含 1 个占位提交(本会话链接仓库时创建,内容为旧版文档),保留,不 force push。
- 本地 Windows 无 symlink 特权:repair/sessiontemp 部分测试本地失败(CI 不受影响)。
- composer-goal-toggle.test.tsx 的 3 个失败:upstream CI 不运行该文件,非 CI 门槛。
