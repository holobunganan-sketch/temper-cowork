# Temper v0.3.0 — 从空仓库开始的 Master Build

> 执行模型：DeepSeek V4 Flash  
> 本地起点：项目目录中仅存在本文件  
> GitHub 远端：`git@github.com:holobunganan-sketch/temper-cowork.git`  
> 远端状态：全新空仓库  
> Reasonix 参考源码：`https://github.com/esengine/DeepSeek-Reasonix.git`，`main-v2`  
> 产品：Temper  
> 目标版本：`v0.3.0`  
> 平台：Windows 11 x64  
> 最终正式发行格式：MSIX  
> 产品理念：`Shape intent. Ship work.`

---

## 0. 本文件的地位

本文件是 Temper v0.3.0 的最高工程规范。

旧 Temper 仓库历史已经被用户删除。禁止恢复旧历史、禁止把旧仓库当开发基础、禁止合并旧 Git 历史。

当前本地只允许假定存在：

```text
TEMPER_V0.3.0_MASTER_BUILD_FINAL.md
```

当前 GitHub 远端是空仓库：

```text
git@github.com:holobunganan-sketch/temper-cowork.git
```

本次任务从零开始。

---

## 1. 最终交付

正式 GitHub Release：

```text
v0.3.0
```

正式 Windows Release 至少包含：

```text
Temper-0.3.0-windows-x64.msix
Temper-Development.cer
Temper-0.3.0-SHA256SUMS.txt
INSTALL-MSIX.md
```

MSIX 内真正运行：

```text
Temper.exe
```

最终用户应能完成：

```text
信任 Temper-Development.cer
→ 安装 MSIX
→ 从开始菜单启动 Temper
→ 配置 Provider / Model
→ 创建 Project
→ Quick Chat / Project Chat
→ 创建正式 Work
→ 使用 Tools / MCP / Skills
→ 审批高风险动作
→ 生成 Artifact
→ 查看 Evidence / Decision
→ Validation / Review
→ Completion Gate
→ Work 完成
→ 重启后恢复
```

---

## 2. MSIX 签名固定方案

MSIX 正常部署需要签名。

v0.3.0 不购买商业证书，使用免费开发自签名代码签名证书：

```text
Subject / Publisher:
CN=Temper Development
```

首次在本机生成：

```text
CurrentUser\My
Code Signing EKU
DigitalSignature
建议有效期 5 年
private key exportable
```

公钥导出：

```text
packaging/msix/Temper-Development.cer
```

允许提交到 Git。

私钥禁止提交。

GitHub Actions 需要签名时：

```text
临时导出 PFX
→ 生成随机强密码
→ PFX Base64
→ gh secret set TEMPER_MSIX_PFX_B64
→ gh secret set TEMPER_MSIX_PFX_PASSWORD
→ 删除临时 PFX
```

禁止把 PFX、密码、Base64 私钥写入代码或日志。

测试机器安装 MSIX 前，需要将公钥证书信任到：

```text
LocalMachine\TrustedPeople
```

---

## 3. 产品架构

固定：

```text
Temper Desktop
React + TypeScript + Wails
        │
Temper CoWork Layer
Project / Chat / Work
Task Contract
Evidence / Decision
Artifact / Validation
Quality / Completion
        │
Reasonix Runtime
Provider / Agent / Controller
Plan / Goal / Todo
Planner / Executor / Subagent
Context / Compaction / Cache
Memory / History
Tools / Permission / Sandbox
MCP / Skills / Plugins / Extensions
Checkpoint / Rewind / Recovery
Serve / ACP / Remote SSH
```

---

## 4. Reasonix 使用原则

Reasonix 是 Temper 的成熟 Runtime 基线。

以下能力禁止 Temper 再造第二套：

```text
Provider
Model Registry
Streaming
Tool Calling
Reasoning Effort
Agent Loop
Controller
Plan
Goal
Todo
Planner
Executor
Subagent
Context
Compaction
Cache
History
Memory
Remember / Forget
File Tools
Shell
Web
Permission
Sandbox
MCP
Skills
Plugins
Extensions
Checkpoint
Rewind
Recovery
Serve
ACP
Remote SSH
Diagnostics
```

固定开发规则：

```text
Reasonix 已有
→ 找源码
→ 找测试
→ 找 Desktop/CLI 调用
→ 复用

接口不适合 Temper UI
→ 建 Temper adapter

Reasonix 缺少 CoWork 语义
→ 写 internal/temper/**

只有 adapter 无法解决
→ 最小修改 Reasonix-owned source
→ 写 regression test
→ 登记 docs/upstream/REASONIX_PATCHES.md
```

Reasonix 的功能和细节参考，视觉 UI 不复制。

---

## 5. Temper 自己负责

Temper-owned：

```text
Project metadata
Chat product surface
Formal Work
Work lifecycle
Task Contract compiler
Acceptance Criteria
Evidence
Decision
Artifact Registry
Artifact Validation
Quality Gate
Completion Gate
Temper Design System
Temper Desktop UI
Windows MSIX release
```

Temper 不再创建第二套大型 Work Graph Planner；任务拆解、Goal、Todo、Planner、Subagent 使用 Reasonix。

---

## 6. DeepSeek V4 Flash 施工规则

Master 只完整读取一次。

立即创建：

```text
AGENTS.md
docs/agent/EXECUTOR_RULES.md
docs/agent/BUILD_STATE.md
docs/agent/CURRENT_TASK.md
docs/agent/RESUME_PROMPT.md
```

后续新 Session 只读：

```text
AGENTS.md
EXECUTOR_RULES.md
BUILD_STATE.md
CURRENT_TASK.md
当前任务相关源码
当前任务相关测试
```

单个 Micro Task：

```text
目标上下文 <= 50K tokens
建议 <= 70K
硬上限 <= 90K

Production files <= 6
Test files <= 4
Docs files <= 3
建议净代码变更 <= 800 lines
```

超过上限：

```text
完成当前安全点
更新 state
compact/reset
继续
```

常规 reasoning effort：

```text
high
```

以下使用 max：

```text
连续两次失败
跨 Go/Wails/React bug
Context/Recovery bug
Windows-only bug
MSIX bug
CI 与本地不一致
Defender detection
Release failure
```

每个 Micro Task 固定：

```text
Locate owner
→ Read implementation
→ Read tests
→ Failing test/reproduction
→ Minimal implementation
→ Target tests
→ Regression tests
→ git diff review
→ state
→ commit
```

禁止：

```text
先写巨大代码再测
删除测试
skip
弱化 assertion
提高 timeout 掩盖死锁
吞掉 error
Mock 替代 production
```

---

## 7. Git 关系

最终：

```text
origin:
git@github.com:holobunganan-sketch/temper-cowork.git

upstream:
https://github.com/esengine/DeepSeek-Reasonix.git
```

Reasonix upstream 只用于：

```text
fetch
compare
reference
manual sync
```

禁止把 Reasonix Git 历史 merge 到 Temper。

禁止：

```bash
git reset --hard
git clean -fdx
git push --force
git tag -f
```

---

# PHASE A — 从空目录建立项目

## A01 本地检查

执行：

```powershell
Get-ChildItem -Force
```

当前目录只应有：

```text
TEMPER_V0.3.0_MASTER_BUILD_FINAL.md
```

如果发现未知源码、旧 `.git` 或其他重要文件：

```text
STOP
```

不要删除。

## A02 检查空远端

```bash
gh auth status
gh repo view holobunganan-sketch/temper-cowork
git ls-remote git@github.com:holobunganan-sketch/temper-cowork.git
```

必须能访问且没有 refs。

远端非空则 STOP，禁止覆盖。

## A03 建 Reasonix Reference

在项目同级目录创建：

```text
../temper-reasonix-reference
```

执行：

```bash
git clone --branch main-v2 --single-branch   https://github.com/esengine/DeepSeek-Reasonix.git   ../temper-reasonix-reference
```

记录：

```bash
git -C ../temper-reasonix-reference rev-parse HEAD
```

实际结果记为：

```text
REASONIX_BASELINE_SHA
```

禁止硬编码本文件生成时的 Reasonix SHA。

## A04 验证 Reasonix baseline

先读：

```text
go.mod
.wails-version
Makefile
CONTRIBUTING.md
desktop/README.md
desktop/frontend/package.json
.github/workflows/ci.yml
```

使用真实 pin。

至少运行：

```bash
go test ./...
go vet ./...
make lint
make desktop-test-short
```

Frontend：

```bash
cd desktop/frontend
pnpm install --frozen-lockfile
pnpm typecheck
pnpm test
pnpm build
```

Desktop：

```bash
cd desktop
wails build
```

Reasonix baseline 明显损坏时，先检查 upstream/main-v2 是否已有修复；更新 reference 到新 HEAD 后重新全测一次。仍损坏才 STOP。

## A05 导入 Reasonix Source Snapshot

不要复制 Reasonix `.git`。

从 pinned SHA：

```bash
git -C ../temper-reasonix-reference archive   --format=tar   -o ../temper-reasonix-source.tar   <REASONIX_BASELINE_SHA>
```

解压到当前项目目录，同时保留本 Master。

## A06 清理 Reasonix 仓库基础设施

保留 Runtime 源码。

删除/不采用 Reasonix 自己的：

```text
.signpath/
site/
workers/
npm/
release-notes/
Reasonix release workflows
Reasonix website deploy
Reasonix R2 deploy
Reasonix npm publish
Reasonix issue/release automation
Reasonix sponsor assets
```

`.github/workflows/` 重新建立 Temper：

```text
ci.yml
release-msix.yml
```

不得误删 Runtime/build 真正依赖文件。

## A07 项目控制文件

创建：

```text
AGENTS.md
docs/spec/TEMPER_V0.3.0_MASTER_BUILD.md

docs/agent/EXECUTOR_RULES.md
docs/agent/BUILD_STATE.md
docs/agent/CURRENT_TASK.md
docs/agent/RESUME_PROMPT.md

docs/upstream/REASONIX_BASELINE.md
docs/upstream/REASONIX_PATCHES.md
docs/upstream/REASONIX_SYNC.md

docs/product/PRODUCT.md
docs/product/ARCHITECTURE.md
docs/product/QUALITY.md
docs/product/SECURITY.md
docs/product/NETWORK_AUDIT.md
docs/product/RELEASE.md

docs/parity/REASONIX_PARITY.md
docs/release/
```

把本文件完整复制为：

```text
docs/spec/TEMPER_V0.3.0_MASTER_BUILD.md
```

## A08 Git init

```bash
git init -b main
git remote add origin git@github.com:holobunganan-sketch/temper-cowork.git
git remote add upstream https://github.com/esengine/DeepSeek-Reasonix.git
git remote -v
```

## A09 License / Attribution

继续 MIT。

保留 Reasonix attribution。

创建：

```text
THIRD_PARTY_NOTICES.md
```

记录：

```text
Reasonix upstream
MIT
REASONIX_BASELINE_SHA
Temper incorporates/modifies Reasonix components
```

## A10 初始 CI

`.github/workflows/ci.yml` 至少：

```text
Go test
Go vet
lint
Desktop short tests
Frontend typecheck
Frontend tests
Frontend build
Windows Wails build
```

CI 使用隔离：

```text
REASONIX_HOME = runner temp
REASONIX_STATE_HOME = runner temp
REASONIX_CACHE_HOME = runner temp
REASONIX_TELEMETRY = 0
DO_NOT_TRACK = 1
```

## A11 First push

```bash
git status
git diff --check
git add .
git commit -m "chore: establish Temper v0.3 Reasonix runtime baseline"
git push -u origin main
```

等待 GitHub CI Green。

---

# PHASE B — Temper 身份与数据隔离

创建：

```text
milestone/01-identity-isolation
```

用户可见：

```text
Temper
0.3.0-dev
Shape intent. Ship work.
```

Windows：

```text
Temper.exe
```

内部 Go module / `reasonix/internal/*` 暂时保留，禁止全仓库机械重命名。

Temper Runtime Home：

```text
%APPDATA%\Temper\runtime\
%LOCALAPPDATA%\Temper\cache\
%APPDATA%\Temper\cowork\
```

必须在 Reasonix Boot / Config load 之前设置或注入：

```text
REASONIX_HOME
REASONIX_STATE_HOME
REASONIX_CACHE_HOME
```

测试两个不同 Home，确保 Temper 不读取正式 Reasonix Provider、Session、Memory、Plugin、Cache。

网络审计：

```text
Provider network       RETAIN
User Web tools         RETAIN
User MCP               RETAIN
Reasonix telemetry     DISABLE
Reasonix crash upload  DISABLE
Reasonix updater       REPLACE/DISABLE
Reasonix release URLs  REPLACE
```

Temper 默认 remote product telemetry OFF。

---

# PHASE C — Reasonix 功能 Parity

在大规模 UI 前完成：

```text
docs/parity/REASONIX_PARITY.md
```

至少覆盖：

```text
Providers / Models
Agent / Controller
Plan / Goal / Todo
Planner / Subagent
Context / Compaction / Cache
History / Memory
Tools
Permissions / Sandbox
MCP
Skills / Plugins / Extensions
Checkpoint / Rewind / Recovery
Serve / ACP / Remote SSH
Diagnostics
```

每项记录：

```text
Upstream owner
Upstream test
Upstream desktop surface
Temper surface
Adapter
Automated test
Manual smoke
Status
```

状态只能：

```text
UPSTREAM_AVAILABLE
TEMPER_WIRED
ADVANCED_ONLY
CLI_ONLY
NOT_WIRED
OUT_OF_SCOPE_WITH_REASON
```

---

# PHASE D — Temper CoWork Store

Temper-owned：

```text
internal/temper/
```

建议：

```text
domain/
store/
project/
work/
evidence/
artifact/
quality/
tools/
```

业务 DB：

```text
%APPDATA%\Temper\cowork\temper.db
```

Tables：

```text
schema_migrations
projects
works
work_events
artifacts
evidence
decisions
acceptance_results
quality_runs
```

要求：

```text
WAL
foreign_keys ON
busy timeout
transactions
UTC
canonical paths
migration
crash safe
```

测试：

```text
fresh
migration
rollback
restart
busy
concurrency
corruption
future schema
```

---

# PHASE E — Project + Chat

Project identity：

```text
real workspace root
```

优先复用 Reasonix project/workspace registration。

支持：

```text
Add existing folder
Open
Search
Recent
Reveal
Remove from Temper
```

Remove 不删除真实 workspace。

Chat = Reasonix Session：

```text
Quick Chat
Project Chat
```

必须继承：

```text
model switch
Plan
Goal
Todo
Tools
Permission
MCP
Skills
Memory
Rewind
```

Temper 不复制 Reasonix transcript。

---

# PHASE F — Formal Work

Work fields：

```text
id
project_id
title
goal
status
reasonix_session_ref
model_ref
quality_profile
created_at
updated_at
started_at
completed_at
final_artifact_id
```

Status：

```text
draft
ready
running
waiting_user
blocked
reviewing
validating
completed
failed
cancelled
```

Work Form：

```text
Project
Goal
Materials
Deliverable
Audience
Constraints
Acceptance Criteria
Source Policy
Pause Policy
Quality
Model
```

编译为 Reasonix Task Contract：

```text
Context
Request
Output format
Constraints
Acceptance criteria
Pause policy
```

禁止第二个 Prompt Runtime。

---

# PHASE G — CoWork Tools / Evidence / Decision

通过 Reasonix Tool Registry 注册：

```text
temper_record_evidence
temper_record_decision
temper_register_artifact
temper_set_final_artifact
temper_report_validation
temper_complete_work
```

每个：

```text
strict JSON schema
typed input
structured output
error codes
tests
```

Evidence：

```text
summary
source_type
source_ref
supports
timestamp
```

Decision：

```text
decision
rationale
alternatives
evidence_ids
timestamp
```

持久化前 secret redaction。

---

# PHASE H — Artifact / Quality

Artifact：

```text
real workspace file + metadata
```

保存：

```text
id
project_id
work_id
relative_path
kind
title
description
sha256
size
validation
is_final
timestamps
```

正式稳定：

```text
MD
TXT
JSON
CSV
HTML
SVG
PNG
JPEG
source code
```

DOCX/XLSX/PPTX/PDF 只有真实 renderer + validator + golden test + Windows open smoke 全通过才标 SUPPORTED。

Acceptance：

```text
pending
pass
fail
uncertain
```

Host Completion Gate 检查：

```text
Task Contract exists
required deliverable exists
final artifact exists
file exists
hash current
validation no fail
no blocking error
acceptance evaluated
```

Reviewer 使用 Reasonix read-only subagent，但不能直接标 completed。

---

# PHASE I — Temper UI

不复制 Reasonix 视觉。

方向：

```text
Quiet Precision
Premium Desktop Workbench
Dense but readable
Low visual noise
Windows productivity
```

使用 Native Windows title bar。

主布局：

```text
Navigation | Main Canvas | Inspector
Runtime Strip at bottom
```

一级：

```text
Home
Chat
Work
Projects
Artifacts
Advanced
Settings
```

Inspector：

```text
Overview
Context
Evidence
Artifacts
Quality
Runtime
```

无 handler 的按钮不能 enabled。

禁止 dead button / console.log only / fake coming-soon action。

---

# PHASE J — Chat / Work / Advanced UI

Chat 必须支持真实：

```text
streaming
reasoning state
tool cards
tool progress
permissions
diff
cancel
steering
model switch
Plan
Goal
Todo
memory citations
web sources
rewind
history
```

不显示私有 chain-of-thought。

Work：

```text
Header
Status
Goal
Task Contract
Plan/Todo
Timeline
Transcript
Evidence
Decisions
Artifacts
Validation
Quality
Completion
```

Project：

```text
Overview
Chats
Works
Artifacts
Instructions
Memory
Settings
```

Settings：

```text
General
Appearance
Language
Providers
Models
Agent
Context
Permissions
Sandbox
MCP
Skills
Memory
Updates
About
```

Advanced：

```text
Sessions
Plugins
Extensions
Hooks
Commands
Themes
Diagnostics
Remote SSH
ACP
Serve
Developer
```

Reasonix config 必须读写真实 contract。

---

# PHASE K — Runtime Observability

底部真实显示：

```text
Model
Mode
Context used/window
Compact threshold
Cache
Input/output tokens
Current tool
Stage
Elapsed
```

没有真实数据则显示 `—`，禁止造数。

Context Inspector：

```text
window
used
remaining
compact ratio
last compaction
history
memory recall
cache
```

---

# PHASE L — i18n / Windows UX

正式：

```text
zh-CN
en-US
```

全部用户文案走 dictionary，并做 key parity test。

Viewport：

```text
760x480
1024x640
1240x720
1440x900
1920x1080
```

DPI：

```text
100%
125%
150%
```

验证：

```text
Minimize
Maximize
Restore
Close
Resize
IME
Clipboard
File picker
Explorer reveal
Dark/light system
multi-monitor when available
```

---

# PHASE M — Production E2E

Mock frontend 不是完成证据。

必须从真实：

```text
Temper.exe
Wails
embedded frontend
Go binding
Reasonix runtime
Temper DB
filesystem
restart
```

验证。

Chat：

```text
stream
read
permission
write
diff
rewind
```

Work fixture：

```text
temp-workspace/notes.md
→ summary.md
→ Evidence
→ Artifact
→ Validation
→ Final
→ Completion
```

重启后 Project / Work / Session / Artifact 必须继续存在。

---

# PHASE N — Defender

不能关闭 Defender 作为 Release 条件。

顺序：

```text
source audit
→ raw Temper.exe scan
→ MSIX staging audit
→ signed MSIX scan
```

raw EXE detected：

```text
先修 EXE
禁止用 MSIX 掩盖
```

禁止：

```text
Defender exclusion
SmartScreen bypass
obfuscation
UPX
process injection
hidden PowerShell
persistence tricks
```

---

# PHASE O — MSIX Packaging

目录：

```text
packaging/msix/
  AppxManifest.xml.template
  assets/
  Temper-Development.cer
  README.md

scripts/msix/
  New-DevCertificate.ps1
  Configure-GitHubSigning.ps1
  Build-MSIX.ps1
  Sign-MSIX.ps1
  Verify-MSIX.ps1
  Install-MSIX.ps1
```

Package identity 固定：

```text
Name: Temper.Cowork.Desktop
Application Id: Temper
DisplayName: Temper
Publisher: CN=Temper Development
PublisherDisplayName: Temper
Architecture: x64
SemVer: 0.3.0
MSIX Version: 0.3.0.0
```

Manifest：

```text
TargetDeviceFamily = Windows.Desktop
MinVersion = 10.0.19041.0
full-trust packaged classic desktop app
Executable = Temper.exe
runFullTrust capability
```

使用 Windows SDK：

```text
MakeAppx.exe
SignTool.exe
```

脚本自动从：

```text
${env:ProgramFiles(x86)}\Windows Kits\10\bin
```

找最高版本 x64 工具，禁止硬编码 SDK 版本。

Temper.exe 在 MSIX 中禁止写安装目录；可变数据只写 `%APPDATA%\Temper`、`%LOCALAPPDATA%\Temper` 或 workspace。

---

## O01 Build-MSIX.ps1

固定：

```text
clean staging
→ Wails production build
→ Temper.exe
→ Defender raw scan if available
→ copy exe/assets
→ render AppxManifest version
→ MakeAppx pack
→ unsigned MSIX
```

临时文件：

```text
Temper-0.3.0-windows-x64-unsigned.msix
```

禁止 Release unsigned 包。

## O02 New-DevCertificate.ps1

寻找：

```text
CN=Temper Development
Code Signing EKU
HasPrivateKey
```

有则复用。

无则创建：

```text
New-SelfSignedCertificate
CurrentUser\My
DigitalSignature
Code Signing EKU
5-year validity
exportable private key
```

导出公钥：

```text
packaging/msix/Temper-Development.cer
```

## O03 Configure-GitHubSigning.ps1

自动：

```text
read local cert
→ random strong PFX password
→ temp PFX
→ Base64
→ gh secret set TEMPER_MSIX_PFX_B64
→ gh secret set TEMPER_MSIX_PFX_PASSWORD
→ delete temp PFX
→ clear password variables
```

禁止日志输出 secret。

## O04 Sign-MSIX.ps1

本地优先 certificate store thumbprint。

CI 使用 temporary PFX from GitHub Secrets。

固定：

```text
SignTool sign /fd SHA256
SignTool verify /pa /v
```

正式输出：

```text
Temper-0.3.0-windows-x64.msix
```

签名失败直接阻断 Release。

## O05 Install-MSIX.ps1

不得偷偷修改 Windows trust store。

如果 CER 未信任：

```text
明确提示用户以管理员身份信任
Temper-Development.cer
到 LocalMachine\TrustedPeople
```

然后：

```powershell
Add-AppxPackage <signed-msix>
```

验证：

```text
installed
Start Menu Temper exists
launch works
```

---

# PHASE P — GitHub MSIX Release

`.github/workflows/release-msix.yml` 使用：

```text
windows-latest
```

步骤：

```text
checkout
setup pinned Go
setup Node/pnpm
setup matching Wails
tests
frontend build
Wails build
MakeAppx
decode temporary PFX
SignTool sign SHA256
SignTool verify
SHA256
upload
```

结束删除临时 PFX。

Release assets：

```text
Temper-0.3.0-windows-x64.msix
Temper-Development.cer
Temper-0.3.0-SHA256SUMS.txt
INSTALL-MSIX.md
```

---

# PHASE Q — RC / Final Release

先：

```text
v0.3.0-rc.1
```

必须：

```text
CI green
Production Chat E2E
Production Work E2E
restart recovery
MCP
Skills
Memory
Rewind
Windows UX
raw EXE Defender clean
signed MSIX Defender clean
MSIX signature valid
MSIX install smoke
```

失败后：

```text
fix
commit
push
rc.2 / rc.3
```

禁止移动旧 tag。

通过后：

```text
v0.3.0
```

正式 Release non-draft / non-prerelease。

---

# Release Blockers

任意存在就禁止正式 Release：

```text
Provider setup broken
Chat cannot complete
Tools broken
Permission broken
Project restart loss
Work cannot complete
Artifact cannot open
Completion Gate bypass
MCP basic fail
Skills basic fail
Memory restart fail
Rewind fail
Home dead button
Work dead button
fake context/runtime data
Windows controls broken
Temper.exe Defender detected
MSIX Defender detected
MSIX signature invalid
MSIX cannot install
MSIX launches broken app
CI red
P0/P1 regression
```

---

# Git / Milestone 工作法

开发分支：

```text
milestone/<number>-<slug>
```

流程：

```text
main
→ branch
→ microtasks
→ tests
→ push
→ CI
→ PR
→ merge
→ compact/reset
→ next
```

Agent 有 GitHub 权限时自行：

```text
create PR
read CI
fix CI
merge PR
```

除硬阻塞外，不需要询问用户是否继续。

---

# 硬阻塞

只有以下允许询问用户：

```text
GitHub auth 失效
GitHub remote 不是空仓库
本地目录存在未知重要旧文件
Windows SDK 缺失且无法安装
Wails 必要环境无法安装
最终 MSIX 安装需要用户手工信任 CER
真实 Provider smoke 需要用户 API Key
Windows Security UI 需要人工动作
```

测试失败、编译失败、CI 失败、MSIX manifest 失败、SignTool 失败、UI bug、DB bug都自行修复。

---

# 每个 Milestone 状态闭环

结束时：

```text
tests
→ BUILD_STATE
→ CURRENT_TASK
→ git diff --check
→ commit
→ push
→ CI
→ PR/merge
→ compact/reset
```

如果 Reasonix 支持 `/compact`，Milestone 边界主动使用。

---

# 最终 Definition of Done

Runtime：

```text
Providers PASS
Models PASS
Agent PASS
Plan PASS
Goal PASS
Todo PASS
Subagents PASS
Context PASS
Compaction PASS
History PASS
Memory PASS
Tools PASS
Permissions PASS
Sandbox PASS
MCP PASS
Skills PASS
Rewind PASS
Recovery PASS
```

Temper：

```text
Project PASS
Chat PASS
Work PASS
Task Contract PASS
Evidence PASS
Decision PASS
Artifact PASS
Validation PASS
Quality PASS
Completion PASS
Restart PASS
```

UI：

```text
Temper identity PASS
zh-CN PASS
en-US PASS
Accessibility PASS
Responsive PASS
No dead actions PASS
Runtime observability PASS
```

Windows：

```text
Temper.exe PASS
Raw EXE Defender PASS
MSIX build PASS
MSIX signature PASS
MSIX Defender PASS
MSIX install PASS
Start Menu launch PASS
Uninstall PASS
User workspace preserved PASS
```

GitHub：

```text
main CI PASS
RC PASS
v0.3.0 tag PASS
Release PASS
MSIX asset PASS
CER asset PASS
SHA256 PASS
```

---

# 执行顺序

```text
A Bootstrap
→ B Identity / Isolation
→ C Reasonix Parity
→ D CoWork Store
→ E Project / Chat
→ F Work
→ G Evidence / Decision / Tools
→ H Artifact / Quality
→ I Temper UI
→ J Chat / Work / Advanced UI
→ K Runtime Observability
→ L i18n / Windows UX
→ M Production E2E
→ N Defender
→ O MSIX
→ P GitHub Release Pipeline
→ Q RC / v0.3.0
```

持续依赖 state 文件执行，不要一次把整个项目塞进上下文。

---

# 最终报告

只在正式 Release 成功后输出：

```text
TEMPER v0.3.0 RELEASED

Repository:
Main SHA:
Tag:
Release:

Reasonix baseline:
Upstream patches:

Windows:
MSIX:
SHA256:
Certificate:
Signature:
Defender EXE:
Defender MSIX:
Install smoke:
Uninstall smoke:

Runtime parity:
Providers:
Agent:
Context:
Memory:
Permissions:
MCP:
Skills:
Recovery:

Temper CoWork:
Project:
Chat:
Work:
Evidence:
Decision:
Artifact:
Quality:
Completion:

Tests:
Go:
Desktop:
Frontend:
E2E:
Windows:

Known limitations:
```

---

# 现在开始

当前目录只应存在：

```text
TEMPER_V0.3.0_MASTER_BUILD_FINAL.md
```

远端：

```text
git@github.com:holobunganan-sketch/temper-cowork.git
```

为空。

从 PHASE A 开始。

不要重新设计本 Master。

除硬阻塞外，不要问用户是否继续。

持续执行，直到 v0.3.0 MSIX Release 完成。
