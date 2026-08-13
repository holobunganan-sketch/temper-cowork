# EXECUTOR_RULES.md — Temper 施工执行规则

> 本文件是 Temper v0.3.0 施工的执行细则,服从 `docs/spec/TEMPER_V0.3.0_MASTER_BUILD.md`(最高规范)。

## 1. 上下文预算(DeepSeek V4 Flash)

- 目标上下文 <= 50K tokens,建议 <= 70K,硬上限 <= 90K。
- 超出预算:完成当前安全点 → 更新 BUILD_STATE/CURRENT_TASK → commit/push → compact/reset → 继续。
- 一次只处理一个明确问题。
- 在 Milestone 边界主动使用 `/compact`(若可用)。

## 2. 文件预算

- Production files <= 6,Test files <= 4,Docs files <= 3,建议净代码变更 <= 800 lines。

## 3. 每个 Micro Task 固定流程

```text
Locate owner
→ Read implementation
→ Read tests
→ Failing test/reproduction
→ Minimal implementation
→ Target tests
→ Regression tests
→ git diff review
→ state 更新
→ commit
```

## 4. 禁止事项

- 先写巨大代码再测
- 删除测试 / skip 测试 / 弱化 assertion
- 提高 timeout 掩盖 deadlock
- 吞掉 error
- Mock 替代 production path
- 只因为能 compile 就报告完成
- dead button / console.log-only 按钮 / 假功能 / Mock 冒充生产功能

## 5. Reasoning Effort

- 普通施工:high
- 复杂调试(连续两次失败 / 跨 Go/Wails/React bug / Context/Recovery bug / Windows-only bug / MSIX bug / CI 与本地不一致 / Defender detection / Release failure):max

## 6. 本地环境已知事项(Windows)

- 用户主目录存在 `.git`(unborn),会导致 Reasonix `workspacelease`/`worktree` 测试互相阻塞。
  - Go 测试必须设置:`GIT_CEILING_DIRECTORIES=C:\Users\ZhouNan`
  - 并把 `TMP`/`TEMP` 指向无 .git 祖先的目录(如 `C:\Myfolder\.testtmp`)。
- Node 26 全局 `navigator.language` 跟随系统中文 locale,导致前端测试硬编码英文文案失败。
  - 前端测试必须:`NODE_OPTIONS="--import=file:///C:/Myfolder/.testenv/fix-locale.mjs"`。
- Windows 非管理员无法创建 symlink(SeCreateSymbolicLinkPrivilege),`internal/repair`、`internal/sessiontemp` 的部分测试在本地失败,但 CI(windows-latest runner 服务账户)不受影响。
- upstream CI 不运行全量 `pnpm test`(run-tests.mjs);`composer-goal-toggle.test.tsx` 不在 CI 脚本中,其 3 个失败在本地 `pnpm test` 下可见,但非 CI 门槛。

## 7. 完成定义(Definition of Done)

功能只有在 **production path + automated tests + runtime smoke** 全部成立时才标 DONE。

## 8. GitHub 纪律

- 每个 Milestone:branch → microtasks → tests → commit → push → CI → PR → CI → merge main。
- 自行创建 PR、读取 CI、修复 CI、merge;除硬阻塞外不询问用户。
- 禁止:`git push --force` / `git tag -f` / `git reset --hard` / `git clean -fdx`。
