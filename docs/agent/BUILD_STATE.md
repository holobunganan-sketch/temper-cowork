# BUILD_STATE.md — Temper v0.3.0 构建状态

> 更新规则:每个 Milestone / Micro Task 完成后更新本文件,记录真实进度与验证证据。

## 当前阶段

**PHASE H — Artifact / Quality(进行中,milestone/07-artifact-quality)**

## 里程碑进度

| Phase | 内容 | 状态 | 备注 |
|-------|------|------|------|
| A | Bootstrap | DONE | main CI 全绿 |
| B | 身份与数据隔离 | DONE | PR #5 merged(361d731) |
| C | Reasonix 功能 Parity | DONE | PR #6 merged(df81c8c) |
| D | CoWork Store | DONE | PR #7 merged(e93d1a5) |
| E | Project + Chat | DONE | PR #8 merged(d7e3e10) |
| F | Formal Work | DONE | PR #9 merged(2f2e92f) |
| G | CoWork Tools / Evidence / Decision | DONE | PR #10 merged(d08a41f) |
| H | Artifact / Quality | IN_PROGRESS | milestone/07-artifact-quality |
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
- GitHub:origin = holobunganan-sketch/temper-cowork;upstream = esengine/DeepSeek-Reasonix
- main 分支 CI 全绿:test / sdk / frontend / desktop 全部 success

## 已验证(带证据)

- [x] Reasonix `go test ./...`(环境修正后) — 除 Windows symlink 特权相关测试外全绿
- [x] `go vet ./...` — 通过
- [x] golangci-lint v2.12.2 — 0 issues
- [x] desktop `go test -short .` — 通过(本地 + CI)
- [x] frontend typecheck / CI 测试脚本 / build — 通过
- [x] `wails build` — 成功
- [x] Temper CI 全绿(run 31722439468:test/sdk/frontend/desktop 全部 success)

## CI 修复记录(PHASE A)

1. workflow env 用 runner.temp 非法 → 改为 step 级 env(仅 root test job)
2. check-motion-ci-contract.mjs 适配 Temper CI 结构(P1 patch)
3. .wails-version CRLF → .gitattributes 强制 LF
4. desktop/build/windows/icon.ico 未被跟踪 → git add -f
5. desktop 打包脚本(desktop-build.sh / package-windows-desktop.sh / verify-windows-portable.sh)被 A06 误删 → 从 baseline 恢复
6. frontend build 必须先于 desktop tests(go:embed all:frontend/dist)
7. workflow 级 REASONIX_HOME 破坏 desktop 测试(93 failures)→ 只对 root test 步骤设置隔离

## 已知偏差与风险

- 远端含 1 个占位提交(链接仓库时创建),保留,不 force push。
- 本地 Windows 无 symlink 特权:repair/sessiontemp 部分测试本地失败(CI 不受影响)。
- composer-goal-toggle.test.tsx 的 3 个失败:upstream CI 不运行该文件,非 CI 门槛。
- dependabot 自动创建 3 个依赖更新 PR,不阻塞主流程。
- P1 patch:check-motion-ci-contract.mjs 适配 Temper CI 结构(见 REASONIX_PATCHES.md)。
