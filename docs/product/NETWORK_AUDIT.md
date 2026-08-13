# NETWORK_AUDIT.md — Temper 网络审计

## 审计结论(继承 Reasonix 默认 + Temper 覆写)

| 网络行为 | 策略 | 说明 |
|----------|------|------|
| Provider network(API 调用) | RETAIN | 用户配置的模型服务 |
| User Web tools | RETAIN | 用户主动触发的网页/搜索 |
| User MCP servers | RETAIN | 用户配置 |
| Reasonix telemetry | DISABLE | 默认关闭 |
| Reasonix crash upload | DISABLE | 默认关闭 |
| Reasonix updater | REPLACE/DISABLE | Temper 使用自己的更新通道(或禁用) |
| Reasonix release URLs | REPLACE | 指向 Temper 发布 |
| Temper product telemetry | OFF(默认) | 远程产品遥测默认关闭 |

## 原则

- 用户可见、可配置;无用户同意不发遥测。
- 网络行为在 Settings → Network 中可审计(继承 Reasonix Diagnostics)。

## 检查方式

- 通过 Reasonix Diagnostics / Network 面板核对实际连接。
- 新增任何出站连接前,先在本文档登记。
