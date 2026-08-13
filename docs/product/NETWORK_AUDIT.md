# NETWORK_AUDIT.md — Temper 网络审计

## 审计结论(继承 Reasonix 默认 + Temper 覆写)

| 网络行为 | 策略 | 说明 |
|----------|------|------|
| Provider network(API 调用) | RETAIN | 用户配置的模型服务 |
| User Web tools | RETAIN | 用户主动触发的网页/搜索 |
| User MCP servers | RETAIN | 用户配置 |
| Reasonix telemetry | DISABLE | 默认关闭;dev 版本由 version=="dev" 保护,正式版由 Settings 默认 off |
| Reasonix crash upload | DISABLE | 默认关闭(REASONIX_CRASH_UPLOAD=0) |
| Reasonix updater | REPLACE/DISABLE | Temper 自管更新通道;REASONIX_UPDATE_DISABLE=1,正式版由 Settings 默认 off |
| Reasonix release URLs | REPLACE | 指向 Temper 发布(MSIX) |
| Temper product telemetry | OFF(默认) | 远程产品遥测默认关闭,用户可显式开启 |

## 原则

- 用户可见、可配置;无用户同意不发遥测。
- 网络行为在 Settings → Network 中可审计(继承 Reasonix Diagnostics)。

## 检查方式

- 通过 Reasonix Diagnostics / Network 面板核对实际连接。
- 新增任何出站连接前,先在本文档登记。

## 实施位置(PHASE B)

- `desktop/temper_identity.go`:ApplyTemperIdentity() 在 main() 最早期注入
  REASONIX_HOME/STATE/CACHE 与 telemetry/update 环境变量。
- 正式 Release(v0.3.0)时,Settings 中 CheckUpdates 与 Metrics 默认 off
  由发布配置保证(见 docs/product/RELEASE.md)。
