# REASONIX_SYNC.md — Reasonix 同步记录

## 同步策略

- Reasonix upstream 仅用于:fetch / compare / reference / manual sync。
- 禁止把 Reasonix Git 历史 merge 进 Temper。
- 手动同步流程:
  1. `git fetch upstream main-v2`
  2. 对比 `git diff REASONIX_BASELINE_SHA..upstream/main-v2 --stat`(只读参考)
  3. 评估变更是否影响 Temper 继承的 Runtime 能力
  4. 如需要,选择性 cherry-pick 文件内容(非 Git 历史)并跑全量回归
  5. 更新 REASONIX_BASELINE.md 的 SHA 与验证结果

## 同步历史

| 日期 | 操作 | 旧 SHA | 新 SHA | 结果 |
|------|------|--------|--------|------|
| 2026-08-13 | 初始基线固定 | — | 49f24d19702c9542ab50500d590237dc872c4d58 | 全量验证通过 |

## 偏差清单

- Temper 删除的 Reasonix 基础设施:site/ workers/ npm/ release-notes/ .signpath/ .goreleaser.yaml、release/signpath/npm 相关 workflows 与 scripts。
- Temper 新增:internal/temper/**、docs/product/**、docs/parity/**、packaging/msix/**、scripts/msix/**、Temper UI。
