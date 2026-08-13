# RESUME_PROMPT.md — 恢复提示

> 新 Session 从这里开始。按顺序读取,然后从 CURRENT_TASK.md 的"下一步"继续。

## 上下文快速恢复

1. 本项目从零构建 **Temper v0.3.0**(Windows 桌面 CoWork 应用),唯一最高规范:
   `docs/spec/TEMPER_V0.3.0_MASTER_BUILD.md`(由 TEMPER_V0.3.0_MASTER_BUILD_FINAL.md 复制而来)。
2. 架构固定:Reasonix Runtime + Temper CoWork Layer + Temper Desktop UI(React + TS + Wails)。
3. Reasonix baseline:main-v2 @ `49f24d19702c9542ab50500d590237dc872c4d58`,已验证全部 CI 检查通过。
4. 执行顺序:A Bootstrap → B 身份隔离 → C Parity → D CoWork Store → E Project/Chat → F Work → G Evidence/Decision → H Artifact/Quality → I/J UI → K Observability → L i18n → M E2E → N Defender → O MSIX → P Release Pipeline → Q RC/v0.3.0。
5. 禁止:重新设计产品、再造第二套 Runtime、merge Reasonix Git 历史、force push、reset --hard、clean -fdx。

## 必读文件(按序)

- AGENTS.md
- docs/agent/EXECUTOR_RULES.md
- docs/agent/BUILD_STATE.md
- docs/agent/CURRENT_TASK.md
- 当前任务相关源码与测试

## 环境速查

- Go: `C:\Myfolder\.toolchain\go\bin`;Node: `C:\Myfolder\.toolchain\node-v26.7.0-win-x64`;已加入 User PATH(新 shell 生效)。
- Go 测试环境变量:`GIT_CEILING_DIRECTORIES=C:\Users\ZhouNan` + `TMP/TEMP=C:\Myfolder\.testtmp`。
- 前端测试:`NODE_OPTIONS="--import=file:///C:/Myfolder/.testenv/fix-locale.mjs"`。
- GOPROXY=https://goproxy.cn,direct(已设 User 级)。
- Windows 网络:HTTPS 直连 github.com 不稳定,需代理 `http://127.0.0.1:7897`(ProxyEnable=0 但代理进程运行中;git 走 SSH 正常)。

## 下一个动作

见 docs/agent/CURRENT_TASK.md → 完成 A07 剩余控制文件 → A08 upstream → A09-A11 → CI green → PHASE B。
