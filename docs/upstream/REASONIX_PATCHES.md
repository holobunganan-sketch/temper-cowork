# REASONIX_PATCHES.md — Reasonix 上游补丁登记

> 当 Temper 需要最小修改 Reasonix-owned source 时,必须在此登记并写 regression test。

## 规则

- 只有 Temper adapter 无法解决的问题,才允许最小修改 Reasonix-owned source。
- 每个 patch 必须:描述问题 → 修改内容 → regression test 位置 → 与 upstream 的 diff 摘要。
- 目标:尽量零 patch;patch 越多,sync 成本越高。

## 已登记 Patch

### P1:check-motion-ci-contract.mjs 适配 Temper CI 结构

- 文件:desktop/frontend/scripts/check-motion-ci-contract.mjs
- 提交:900b452
- 问题:契约脚本校验 Reasonix 的 CI job 结构(desktop/desktop-macos/desktop-windows/lint/site),Temper 重组 CI 为 test/sdk/frontend/desktop 后,`test:motion` 在 CI 中必然失败。
- 修改:job 结构断言改为校验 Temper 的 frontend job 运行 `test:motion`/`test:transcript`、desktop job 运行 `wails build`;保留全部与 CI 结构无关的守护(源文件无 test-only WebView2 instrumentation、retired path 不存在、脚本内容必须包含特定套件)。
- Regression:本地 `pnpm test:motion` 通过(9 passed,含契约自检)。
- 后续 sync 注意:upstream 若改动该脚本,需重新评估 job 结构断言。

### P2:appidentity.AppUserModelID 改为 Temper 身份

- 文件:internal/appidentity/identity.go、identity_test.go
- 提交:milestone/01-identity-isolation(PHASE B)
- 问题:Temper 是独立产品,Windows 任务栏/通知身份必须与 Reasonix 区分。
- 修改:`AppUserModelID` 从 `"Reasonix"` 改为 `"Temper.Cowork.Desktop"`(与 MSIX package identity 一致);同步更新 identity_test.go 断言。
- Regression:`go test ./internal/appidentity/` 通过。
- 后续 sync 注意:upstream 若改动该值需重新评估。

### P3:frontend bundle budget 微调容纳 Temper 品牌

- 文件:desktop/frontend/scripts/check-bundle-budget.mjs
- 提交:milestone/08-temper-ui(PHASE I)
- 问题:Temper 品牌 Design System 规则(~0.7 KiB gzip)使 initial JS gzip 从 422.9 增至 423.6,超出 Reasonix 硬预算 423.5。
- 修改:budget 423.5 → 424.5 KiB(注释说明 Temper 品牌注入),仍保留增长门禁。
- 后续 sync 注意:upstream 若收紧 budget 需重新评估。

## 同步提醒

在 REASONIX_SYNC.md 中跟踪:本分支与 upstream main-v2 的偏差(diff 计数)。任何 patch 都会增加后续 sync 的手动成本。
