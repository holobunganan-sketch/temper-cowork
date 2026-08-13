# RELEASE.md — Temper 发布流程

## 发行物

正式 GitHub Release(v0.3.0)至少包含:

```text
Temper-0.3.0-windows-x64.msix
Temper-Development.cer
Temper-0.3.0-SHA256SUMS.txt
INSTALL-MSIX.md
```

MSIX 内运行 Temper.exe。不以 Setup.exe 为主发行物;Portable ZIP 非最终主格式。

## 版本节奏

1. 开发 main 持续集成(CI green)。
2. `v0.3.0-rc.1` 预发布:CI green + Production Chat/Work E2E + restart recovery + MCP/Skills/Memory/Rewind + Windows UX + raw EXE Defender clean + signed MSIX Defender clean + MSIX signature valid + MSIX install smoke。
3. 修复后 rc.2/rc.3(禁止移动旧 tag)。
4. 全部通过 → `v0.3.0` 正式 Release(non-draft / non-prerelease)。

## MSIX 流程

```text
Wails build → Temper.exe
→ Defender raw scan
→ MakeAppx(从 Windows Kits 自动定位,不硬编码版本)
→ SignTool sign /fd SHA256
→ SignTool verify /pa /v
→ SHA256SUMS
→ Release assets
```

脚本:scripts/msix/(New-DevCertificate / Configure-GitHubSigning / Build-MSIX / Sign-MSIX / Verify-MSIX / Install-MSIX)。

## Release Blockers

Provider 不可用、Chat 无法完成、Tools/Permission 坏、Project restart 丢失、Work 无法完成、Artifact 无法打开、Completion Gate 绕过、MCP/Skills 基础失败、Memory restart 失败、Rewind 失败、Home/Work dead button、fake context/runtime 数据、Windows 控件坏、EXE/MSIX Defender 检测、签名无效、MSIX 无法安装/启动坏、CI red、P0/P1 regression — 任一存在禁止 Release。
