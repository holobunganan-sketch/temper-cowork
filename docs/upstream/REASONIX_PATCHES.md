# REASONIX_PATCHES.md — Reasonix 上游补丁登记

> 当 Temper 需要最小修改 Reasonix-owned source 时,必须在此登记并写 regression test。

## 规则

- 只有 Temper adapter 无法解决的问题,才允许最小修改 Reasonix-owned source。
- 每个 patch 必须:描述问题 → 修改内容 → regression test 位置 → 与 upstream 的 diff 摘要。
- 目标:尽量零 patch;patch 越多,sync 成本越高。

## 已登记 Patch

(暂无 — 目标保持零 patch)

## 同步提醒

在 REASONIX_SYNC.md 中跟踪:本分支与 upstream main-v2 的偏差(diff 计数)。任何 patch 都会增加后续 sync 的手动成本。
