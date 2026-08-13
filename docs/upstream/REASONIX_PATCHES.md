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

## 同步提醒

在 REASONIX_SYNC.md 中跟踪:本分支与 upstream main-v2 的偏差(diff 计数)。任何 patch 都会增加后续 sync 的手动成本。
