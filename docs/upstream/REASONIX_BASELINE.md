# REASONIX_BASELINE.md — Reasonix 基线记录

## 基线信息

- Upstream:https://github.com/esengine/DeepSeek-Reasonix.git
- 分支:main-v2
- Baseline SHA:`49f24d19702c9542ab50500d590237dc872c4d58`
- 获取日期:2026-08-13
- License:MIT

## 基线验证结果

| 检查 | 结果 | 备注 |
|------|------|------|
| go test ./... | PASS(环境修正后) | 本地 Windows 无 symlink 特权,repair/sessiontemp 部分用例本地失败(CI 服务账户不受影响) |
| go vet ./... | PASS | |
| golangci-lint v2.12.2 | PASS | 0 issues |
| desktop go test -short | PASS | 217s |
| frontend typecheck | PASS | |
| frontend CI 测试脚本 | PASS | terminal/task-monitor/workspace/usage-stats/settings-responsive/transcript/motion/composer-menu/remote |
| frontend build | PASS | bundle budget 通过 |
| wails build | PASS | 产出 reasonix-desktop.exe |

## 本地环境修正(非 upstream 问题)

1. 用户主目录存在 `.git`(unborn)导致 workspacelease/worktree 测试失败 → `GIT_CEILING_DIRECTORIES=C:\Users\ZhouNan` + TMP/TEMP 指向干净目录。
2. Node 26 全局 navigator.language=zh-CN 导致前端测试英文文案失败 → NODE_OPTIONS 注入 en-US。
3. Windows symlink 特权 → repair/sessiontemp 部分用例本地失败。
4. composer-goal-toggle.test.tsx 3 个失败:upstream CI 不运行该文件(不在任何 CI pnpm script 中),非 CI 门槛。

## 更新策略

Reasonix upstream 仅用于 fetch/compare/reference/manual sync。禁止 merge 其 Git 历史。手动同步时更新本文件与 REASONIX_SYNC.md。
