# QUALITY.md — Temper 质量策略

## 测试纪律

- 每项功能:production path + automated tests + runtime smoke 全部成立才标 DONE。
- Regression:先建立 failing test/reproduction,再实现。
- 禁止:删除失败测试 / skip / 弱化 assertion / timeout 掩盖 deadlock / 吞 error / Mock 替代 production。

## 测试分层

| 层 | 位置 | 命令 |
|----|------|------|
| Go 单元/集成 | internal/** | go test ./... |
| Desktop Go | desktop/ | cd desktop && go test -short . |
| Frontend | desktop/frontend | pnpm typecheck / pnpm test:<script> / pnpm build |
| E2E | PHASE M | Temper.exe + Wails + DB + filesystem + restart |
| Windows | PHASE L/N/O | viewport/DPI/窗口、Defender、MSIX install/uninstall |

## 质量门(Quality Gate)

- CI 全绿(main)
- 无 P0/P1 regression
- 无 dead button / fake 数据 / console.log-only 动作

## 完成门(Completion Gate)

Host 检查:Task Contract exists → required deliverable exists → final artifact exists → file exists → hash current → validation no fail → no blocking error → acceptance evaluated。

## 环境注意

本地 Windows 测试需:Go 测试设 `GIT_CEILING_DIRECTORIES` + 干净 TMP/TEMP;前端测试注入 en-US locale。详见 EXECUTOR_RULES.md。
