# SECURITY.md — Temper 安全策略

## 原则

- Permission/Sandbox 继承 Reasonix,管控所有工具动作。
- 高风险动作必须用户审批。
- Secret redaction:Evidence/Decision 持久化前必须脱敏。
- 私钥永不进入 Git。

## 数据隔离

- Temper 与 Reasonix 使用不同 Home(REASONIX_HOME 注入 %APPDATA%\Temper\runtime\)。
- 业务 DB:%APPDATA%\Temper\cowork\temper.db(WAL, foreign_keys ON, busy timeout, crash safe)。

## MSIX 签名安全

- 自签名证书 CN=Temper Development,公钥(Temper-Development.cer)可提交。
- 私钥/PFX 禁止提交;CI 通过 GitHub Secrets(TEMPER_MSIX_PFX_B64 / TEMPER_MSIX_PFX_PASSWORD)注入。
- 禁止在日志/输出中出现 PFX、Base64、password、private key。
- SignTool sign /fd SHA256 + SignTool verify /pa /v。

## Defender 策略

- 禁止:关闭 Defender 完成 Release、Defender exclusion、SmartScreen bypass、UPX、混淆、process injection、hidden PowerShell、恶意持久化。
- 顺序:source audit → raw Temper.exe scan → MSIX staging audit → signed MSIX scan。

## 报告

安全漏洞请通过仓库 issue 或维护者联系(与 Reasonix SECURITY.md 流程对齐)。
