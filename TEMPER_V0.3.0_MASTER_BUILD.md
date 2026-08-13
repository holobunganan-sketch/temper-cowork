# Temper v0.3.0 — 从零构建总执行指令

> **目标执行模型：DeepSeek V4 Flash**
>
> **目标产品：Temper**
>
> **目标版本：v0.3.0**
>
> **新 GitHub 仓库：`holobunganan-sketch/temper`**
>
> **Reasonix 参考上游：`https://github.com/esengine/DeepSeek-Reasonix.git` / `main-v2`**
>
> **平台：Windows 11 x64**
>
> **产品理念：Shape intent. Ship work.**
>
> **总路线：以 Reasonix 的成熟 Runtime/功能行为为基础，建立一个全新、干净、独立的 Temper 仓库；完整保留 Reasonix 的通用 Agent 能力，在其上增加 CoWork 工作体系和 Temper 独立桌面 UI。**

---

# 0. 本文件是什么

本文件是 Temper v0.3.0 从零构建的最高工程规范。

旧项目：

```text
holobunganan-sketch/temper-cowork
```

只作为历史。

**不得修改、不得继续修、不得从其中 cherry-pick 大片旧实现。**

新的项目：

```text
holobunganan-sketch/temper
```

必须从一个全新的本地目录和全新的 Git 历史开始。

Reasonix 只作为：

```text
功能基线
Runtime 源码基线
行为参考
测试参考
工程细节参考
```

不再执行：

```text
Temper 历史与 Reasonix 历史 merge
git merge -s ours
把两个已有仓库强行拼接
```

---

# 1. 最终目标

Temper v0.3.0 必须成为一个真正可使用的 Windows CoWork 桌面 Agent。

最终用户能够完成：

```text
安装或解压
→ 启动 Temper
→ 配置 Provider / Model
→ 创建或打开 Project
→ Quick Chat
→ 创建正式 Work
→ Agent 使用文件、Shell、Web、MCP、Skills 等能力工作
→ 用户审批危险操作
→ 查看 Context / Cache / Runtime
→ 查看 Todo / Plan / Goal
→ 记录 Evidence / Decision
→ 得到 Artifact
→ Review / Validation
→ Completion Gate
→ 正式完成
→ 重启后继续
```

v0.3.0 最终至少发布：

```text
Temper-0.3.0-windows-x64-portable.zip
Temper-0.3.0-SHA256SUMS.txt
```

只有 Defender clean 时才额外发布：

```text
Temper-0.3.0-windows-x64-setup.exe
```

代码签名不是 v0.3.0 的前置条件。

MSIX 不是 v0.3.0 主发行格式。

---

# 2. 产品架构固定

Temper 固定采用：

```text
┌─────────────────────────────────────────────────────┐
│                  Temper Desktop                     │
│                                                     │
│ React + TypeScript + Wails                          │
│ Temper 独立 UI / Design System                      │
└───────────────────────┬─────────────────────────────┘
                        │
                 Temper Application
                        │
┌───────────────────────▼─────────────────────────────┐
│                Temper CoWork Layer                  │
│                                                     │
│ Project                                             │
│ Chat                                                │
│ Work                                                │
│ Task Contract                                       │
│ Evidence                                            │
│ Decision                                            │
│ Artifact                                            │
│ Quality Gate                                        │
│ Completion Gate                                     │
└───────────────────────┬─────────────────────────────┘
                        │
                Thin Adapter / Tools
                        │
┌───────────────────────▼─────────────────────────────┐
│                 Reasonix Kernel                     │
│                                                     │
│ Provider / Model                                    │
│ Agent / Controller                                  │
│ Plan / Goal / Todo                                  │
│ Planner / Executor / Subagent                       │
│ Context / Compaction / Cache                        │
│ Memory / History                                    │
│ Tools                                               │
│ Permission / Sandbox                                │
│ MCP                                                 │
│ Skills / Plugins / Extensions                       │
│ Checkpoint / Rewind / Recovery                      │
│ Serve / ACP / Remote SSH                            │
└─────────────────────────────────────────────────────┘
```

---

# 3. 最重要的架构规则

## 3.1 Reasonix 已有能力禁止重复实现

以下能力以 Reasonix 当前源码为事实来源：

```text
Provider
Model registry
OpenAI-compatible protocol
Anthropic-compatible protocol
Responses-compatible protocol
Streaming
Tool calling
Reasoning effort
Reasoning language

Agent loop
Controller
Plan
Goal
Task Contract semantics
Todo
Planner / Executor
Subagent
Background task

Context accounting
Compaction
Prompt/cache behavior
Usage
History
Memory
Remember
Forget

Built-in tools
Shell/Bash
File tools
Web
Permission
Sandbox

MCP
Plugin
Skill
Extension
Hook
Command
Theme

Checkpoint
Rewind
Fork
Recovery

Serve
ACP
Remote SSH
Diagnostics
```

如果 Reasonix 已经存在：

```text
先复用
再包装
最后才允许小范围修改
```

禁止创建第二套：

```text
Agent Loop
Provider Layer
Context Engine
Memory Engine
MCP Client
Permission Engine
Recovery Engine
Planner
```

---

## 3.2 Temper 自己负责的内容

Temper 新增：

```text
CoWork Project metadata
Formal Work
Work lifecycle
Work ↔ Reasonix session mapping
Task Contract form/compiler
Acceptance Criteria
Evidence Ledger
Decision Ledger
Artifact Registry
Artifact Validation
Reviewer flow
Quality Gate
Completion Gate
Delivery state
Temper UI
Temper Windows identity
Temper Release
```

---

## 3.3 不创建第二个 Work Graph Planner

Reasonix 已经拥有：

```text
Goal
Todo
Task catalog
Planner
Subagents
Completion / Delivery semantics
```

Temper v0.3.0 不再重新制造大型 Work Graph Runtime。

Temper Work 负责：

```text
合同
状态
证据
成果
质量
交付
```

任务拆解和实际 Agent 执行继续由 Reasonix 负责。

---

# 4. 从零仓库策略

目标新仓库：

```text
holobunganan-sketch/temper
```

旧仓库：

```text
holobunganan-sketch/temper-cowork
```

永久只读。

---

## 4.1 Reasonix 源码导入方式

不要 fork/merge 两套 Git 历史。

执行时：

1. 在临时目录 clone Reasonix `main-v2`；
2. 记录实际 SHA；
3. 在临时 Reasonix checkout 跑 baseline tests；
4. 将该 commit 的工作树导出到一个全新的 `temper` 目录；
5. 删除 Reasonix `.git`；
6. 清理 Reasonix 自己的 repository infrastructure；
7. 在 `temper` 目录执行新的 `git init -b main`；
8. 创建 Temper 第一个 commit；
9. 创建全新的 GitHub repository。

这样：

```text
Reasonix source
        ↓
Temper source baseline

Reasonix Git history
        ✗ 不进入 Temper

旧 Temper Git history
        ✗ 不进入 Temper
```

---

## 4.2 Future upstream sync

Temper 新仓库允许添加：

```text
upstream = https://github.com/esengine/DeepSeek-Reasonix.git
```

只用于：

```text
fetch
compare
inspect
manual port
```

禁止直接 merge unrelated upstream history。

维护：

```text
docs/upstream/REASONIX_BASELINE.md
docs/upstream/REASONIX_PATCHES.md
docs/upstream/REASONIX_SYNC.md
```

每次 upstream-owned 修改记录：

```text
Patch ID
Reasonix files
Reason
Temper behavior
Tests
Can remove later?
Merge risk
```

---

# 5. DeepSeek V4 Flash 专用执行协议

这是硬规则。

## 5.1 一个 Session 只执行一个 Milestone

```text
一个 Reasonix / DeepSeek Session
=
一个 Milestone
```

Milestone 完成：

```text
测试
→ BUILD_STATE
→ CURRENT_TASK
→ commit
→ push
→ CI
→ 报告
→ STOP
```

下一 Milestone 用新 Session。

---

## 5.2 Context 限制

单个 Micro Task：

```text
目标上下文 <= 50K tokens
建议上限 <= 70K
硬上限 <= 90K
```

超过 90K：

```text
停止继续读取
写状态
完成当前安全点
换新 Session
```

禁止：

```text
一次扫描整个 Reasonix
每轮重新读完整 Master
在一个 Context 内完成 10 个子系统
```

---

## 5.3 Micro Task 大小

一个 Micro Task：

```text
Production files <= 6
Test files <= 4
Docs files <= 3
净代码修改建议 <= 800 lines
```

若任务需要同时理解超过两个大型子系统：

```text
继续拆
```

---

## 5.4 DeepSeek reasoning

常规：

```text
thinking enabled
reasoning_effort = high
```

使用 `max`：

```text
连续两次失败
跨 Go / Wails / React 问题
Context / Recovery bug
upstream patch conflict
CI 与本地不一致
Windows-only bug
Defender detection
Release failure
```

---

## 5.5 每个任务固定流程

```text
Read state
→ Find owner
→ Find tests
→ Write/reproduce failing test
→ Minimal implementation
→ Target test
→ Related regression
→ Diff review
→ Commit
```

禁止：

```text
先写一大堆代码再测
顺手重构无关代码
删失败测试
弱化 assertion
加 timeout 掩盖 deadlock
```

---

# 6. Agent 状态文件

项目必须创建：

```text
AGENTS.md

docs/agent/BUILD_STATE.md
docs/agent/CURRENT_TASK.md
docs/agent/EXECUTOR_RULES.md
docs/agent/RESUME_PROMPT.md

docs/spec/TEMPER_V0.3.0_MASTER_BUILD.md

docs/upstream/REASONIX_BASELINE.md
docs/upstream/REASONIX_PATCHES.md
docs/upstream/REASONIX_SYNC.md

docs/product/PRODUCT.md
docs/product/ARCHITECTURE.md
docs/product/QUALITY.md
docs/product/SECURITY.md
docs/product/RELEASE.md

docs/parity/REASONIX_PARITY.md
docs/release/
```

---

## 6.1 新 Session 只读

```text
AGENTS.md
docs/agent/EXECUTOR_RULES.md
docs/agent/BUILD_STATE.md
docs/agent/CURRENT_TASK.md
当前任务相关代码
当前任务相关测试
```

除非 `CURRENT_TASK` 明确引用，否则禁止重新读完整 Master。

---

# 7. Git 规则

禁止：

```bash
git reset --hard
git clean -fdx
git push --force
git tag -f
```

main 必须可构建。

开发使用：

```text
milestone/<id>-<name>
```

每个 Milestone：

```text
branch from main
→ implementation
→ local tests
→ push
→ CI
→ PR
→ merge main
→ next milestone
```

Commit：

```text
chore:
feat:
fix:
test:
docs:
refactor:
build:
ci:
```

---

# 8. MILESTONE 00 — Preflight + Reasonix Baseline

只做基础验证。

---

## 8.1 环境

检查：

```powershell
$ErrorActionPreference = "Stop"

git --version
gh --version
gh auth status
go version
node --version
pnpm --version
```

缺少 pnpm：

```powershell
corepack enable
corepack prepare pnpm@10 --activate
```

Reasonix 的真实版本以：

```text
go.mod
.wails-version
desktop/frontend/package.json
pnpm-lock.yaml
```

为准。

禁止使用旧 Temper 的依赖版本表覆盖 Reasonix pins。

---

## 8.2 新目录安全

当前工作 parent directory 中：

```text
temper
```

如果不存在：

继续。

如果存在且为空：

允许。

如果存在且包含未知文件：

```text
STOP
```

禁止覆盖。

---

## 8.3 Clone 临时 Reasonix Reference

创建：

```text
.reasonix-reference
```

它不能位于最终 Temper Git repository 内。

执行：

```bash
git clone --branch main-v2 --single-branch https://github.com/esengine/DeepSeek-Reasonix.git .reasonix-reference
```

如果 clone 网络失败：

```text
重试 3 次
```

仍失败才报告 blocker。

记录：

```bash
git -C .reasonix-reference rev-parse HEAD
git -C .reasonix-reference log -5 --oneline
```

输出 SHA 写：

```text
PINNED_REASONIX_SHA
```

---

## 8.4 Reasonix Baseline Test

在 reference clone 中运行 Reasonix 自己的测试。

先读：

```text
Makefile
CONTRIBUTING.md
desktop/README.md
.github/workflows/ci.yml
```

至少：

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
cd ../..
```

Wails：

安装 `.wails-version` 对应 Wails。

运行：

```bash
cd desktop
wails build
cd ..
```

只有 baseline 足够可信才继续。

如果 Reasonix 当前 commit 自己失败：

```text
记录真实 failure
优先更新到新的 upstream main-v2 commit 重试一次
```

仍失败：

```text
STOP
```

不要把破损 baseline 导入 Temper。

---

# 9. MILESTONE 01 — Clean Temper Repository

---

## 9.1 导出 Reasonix Source Snapshot

从 pinned Reasonix checkout 导出：

```bash
git -C .reasonix-reference archive --format=tar -o ../reasonix-temper-source.tar <PINNED_REASONIX_SHA>
```

创建：

```text
temper
```

解压 source tar 到：

```text
temper/
```

确认：

```text
temper/.git
```

不存在。

---

## 9.2 删除 Reasonix repository-only infrastructure

删除或不导入：

```text
.github/workflows/*
.github/CODEOWNERS
.github/sponsor/
.signpath/
site/
workers/
npm/
release-notes/
benchmarks/
```

如果某项被 Runtime/build test 真正依赖：

```text
先证明依赖
再保留最小部分
```

保留：

```text
cmd/
internal/
desktop/
sdk/（如果 runtime/extension tests 需要）
必要 scripts/
必要 tools/
go.mod
go.sum
Makefile
.wails-version
pnpm-lock.yaml
LICENSE
```

---

## 9.3 初始化新 Git

```bash
cd temper
git init -b main
git remote add upstream https://github.com/esengine/DeepSeek-Reasonix.git
```

创建：

```text
docs/upstream/REASONIX_BASELINE.md
```

写：

```text
Repository
Branch
Pinned SHA
Imported date
Import method
Baseline test evidence
```

---

## 9.4 License

Temper 继续 MIT。

保留 Reasonix MIT copyright notice。

增加 Temper copyright。

创建：

```text
THIRD_PARTY_NOTICES.md
```

明确：

```text
Temper incorporates and modifies components from Reasonix.
Reasonix repository:
https://github.com/esengine/DeepSeek-Reasonix
License: MIT
Pinned baseline: <SHA>
```

---

## 9.5 创建新 GitHub Repo

确认：

```bash
gh repo view holobunganan-sketch/temper
```

如果返回 Not Found：

```bash
gh repo create holobunganan-sketch/temper \
  --public \
  --source=. \
  --remote=origin \
  --description "Temper — a local-first CoWork desktop agent. Shape intent. Ship work."
```

如果仓库意外已存在且非空：

```text
STOP
```

禁止覆盖。

---

## 9.6 Baseline commit

先创建最小 Temper CI：

```text
.github/workflows/ci.yml
```

至少：

```text
root Go test
root Go vet
desktop short test
frontend typecheck
frontend test
frontend build
Windows Wails build
```

然后：

```bash
git add .
git commit -m "chore: establish Temper Reasonix runtime baseline"
git push -u origin main
```

等待 CI Green。

---

# 10. MILESTONE 02 — Identity + Data + Network Isolation

创建：

```text
milestone/02-isolation
```

---

## 10.1 Temper Identity

用户可见：

```text
Temper
Version 0.3.0-dev
Shape intent. Ship work.
```

Windows：

```text
Temper.exe
```

Repository：

```text
holobunganan-sketch/temper
```

内部 Go module 暂时继续：

```text
reasonix
```

禁止为了品牌全仓库修改：

```text
reasonix/internal/*
```

---

## 10.2 Temper Runtime Home

Reasonix Runtime 必须使用 Temper 独立数据根。

Windows：

```text
%APPDATA%\Temper\runtime\
%LOCALAPPDATA%\Temper\cache\
```

Temper CoWork：

```text
%APPDATA%\Temper\cowork\
```

在 Reasonix Boot 之前设置/注入：

```text
REASONIX_HOME
REASONIX_STATE_HOME
REASONIX_CACHE_HOME
```

或复用 Reasonix 当前等价配置 seam。

---

## 10.3 Isolation Tests

建立两个假 Home：

```text
ReasonixHome/
TemperHome/
```

Reasonix Home 放：

```text
reasonix-only provider
reasonix-only session
reasonix-only memory
```

Temper Home 放：

```text
temper-only provider
temper-only session
temper-only memory
```

Temper 启动后只能读取 Temper Home。

必须测试：

```text
provider
credential
session
memory
plugin
cache
```

---

## 10.4 Network Audit

搜索：

```text
reasonix.io
crash.reasonix.io
update
telemetry
R2
Discord
website
release gateway
```

创建：

```text
docs/product/NETWORK_AUDIT.md
```

每项：

```text
RETAIN
REPLACE
DISABLE
DEV_ONLY
```

固定：

```text
Provider network                 RETAIN
User MCP                         RETAIN
User web tools                   RETAIN

Reasonix product telemetry       DISABLE
Reasonix crash upload            DISABLE
Reasonix updater                 REPLACE/DISABLE
Reasonix website/release URLs    REPLACE
```

Temper v0.3 默认：

```text
remote product telemetry OFF
```

---

# 11. MILESTONE 03 — Reasonix Capability Parity Contract

此 Milestone 不做大 UI。

创建：

```text
docs/parity/REASONIX_PARITY.md
```

逐项检查真实 upstream source、test、Desktop surface。

至少覆盖：

## Providers

```text
OpenAI-compatible
Anthropic-compatible
Responses-compatible
custom endpoint
model list
default model
mid-session switch
context window
reasoning effort
reasoning language
web search
usage
provider balance
API key save/delete
```

## Agent

```text
single agent
Plan
Goal
Task Contract semantics
Todo
Planner
Executor
Subagent
background work
cancel
steering
```

## Context

```text
window
compaction
manual compact
cache metrics
usage
history
memory
remember
forget
```

## Tools

```text
read
write
edit
move
grep
glob
bash
web
tool progress
diff
workspace preview
```

## Safety

```text
Ask
Auto
YOLO
allow
ask rules
deny
dynamic shell handling
sandbox root
allow_write
forbid_read
```

## MCP

```text
stdio
streamable HTTP
SSE
OAuth/PKCE if current upstream supports
tools
prompts
resources
roots
progress
```

## Extensibility

```text
Skills
Plugins
Extensions
Hooks
Commands
Themes
Runtime reload
```

## Recovery

```text
checkpoint
rewind code
rewind conversation
rewind both
fork
summarize
recovery
resume
```

## Advanced

```text
Serve
ACP
Remote SSH
Diagnostics
Crash reports
```

每项状态：

```text
UPSTREAM_AVAILABLE
TEMPER_WIRED
ADVANCED_ONLY
CLI_ONLY
NOT_WIRED
OUT_OF_SCOPE_WITH_REASON
```

---

# 12. MILESTONE 04 — Temper CoWork Store

创建：

```text
internal/temper/
```

结构建议：

```text
internal/temper/domain/
internal/temper/store/
internal/temper/project/
internal/temper/work/
internal/temper/artifact/
internal/temper/evidence/
internal/temper/quality/
internal/temper/tools/
```

---

## 12.1 Store

Reasonix 自己的 disposable catalog/projection DB 不作为 Temper business source of truth。

建立 Temper authoritative SQLite store。

建议：

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
quality_runs
acceptance_results
```

---

## 12.2 DB requirements

```text
WAL
foreign_keys=ON
busy_timeout
transactions
UTC
canonical paths
schema version
crash-safe writes
```

测试：

```text
fresh migration
upgrade
restart
rollback
busy
concurrent reads
corruption handling
future schema fail-safe
```

---

# 13. MILESTONE 05 — Project + Chat

---

## 13.1 Project

Project 主身份：

```text
workspace root
```

优先复用 Reasonix Desktop 现有 workspace/project registration。

Temper 添加：

```text
display name
favorite
last opened
Work catalog
Artifact catalog
```

支持：

```text
Add existing folder
Open
Recent
Search
Reveal
Remove from Temper
```

Remove：

```text
不删除真实 workspace
```

---

## 13.2 Chat

Chat = Reasonix Session。

两种：

```text
Quick Chat
Project Chat
```

Temper 不复制 transcript。

只存必要 reference：

```text
Reasonix session ID/path
Project root
```

Chat 必须完整支持 Reasonix：

```text
model switch
reasoning
Plan
Goal
Todo
tools
permissions
MCP
skills
memory
rewind
```

---

# 14. MILESTONE 06 — Formal Work

Work 是 Temper 的核心。

---

## 14.1 Work Model

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

每条 transition 有测试。

---

## 14.2 Work Creation

Basic：

```text
Project
Goal
Materials
Deliverable
Quality
Model
```

Advanced：

```text
Audience
Constraints
Acceptance Criteria
Source Policy
Pause Policy
```

---

## 14.3 Task Contract

Temper 编译：

```text
Context
Request
Output format
Constraints
Acceptance Criteria
Pause policy
```

交给 Reasonix Goal / Task Contract。

禁止第二个 Prompt Engine。

---

## 14.4 Runtime

```text
Create Work
→ Persist Draft
→ Compile Contract
→ Create Reasonix Session
→ Start Goal
→ Reasonix executes
→ Temper collects evidence/artifacts
→ Review
→ Validation
→ Completion Gate
```

---

# 15. MILESTONE 07 — Temper CoWork Tools

通过 Reasonix Tool Registry 正式注册：

```text
temper_record_evidence
temper_record_decision
temper_register_artifact
temper_set_final_artifact
temper_report_validation
temper_complete_work
```

每个 tool：

```text
stable JSON schema
typed input
structured output
error code
tests
```

---

## 15.1 Evidence

```text
summary
source_type
source_ref
supports
timestamp
```

source：

```text
workspace file
URL
MCP resource
Reasonix session/history
manual note
```

Secret redaction 后才能持久化。

---

## 15.2 Decision

```text
decision
rationale
alternatives
evidence_ids
timestamp
```

---

# 16. MILESTONE 08 — Artifact System

Artifact = 真实文件 + metadata。

---

## 16.1 Registry

记录：

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
validation_state
is_final
created_at
updated_at
```

---

## 16.2 Core supported formats

v0.3 正式保证：

```text
Markdown
TXT
JSON
CSV
HTML
SVG
PNG
JPEG
source code
```

Preview：

```text
Markdown render
text/code viewer
JSON viewer
CSV table
SVG/image
HTML source/controlled preview
```

大文件必须限制和虚拟化。

---

## 16.3 Office/PDF

可开发：

```text
DOCX
XLSX
PPTX
PDF
```

但只有满足：

```text
real renderer
structural validator
content validator
Windows open smoke
regression tests
```

才能在 UI 显示：

```text
SUPPORTED
```

否则：

```text
EXPERIMENTAL
```

不允许“按钮存在 = 功能完成”。

---

# 17. MILESTONE 09 — Quality + Completion

---

## 17.1 Acceptance Result

每一条 Acceptance Criterion：

```text
pending
pass
fail
uncertain
```

保存：

```text
criterion
result
evidence
validator
timestamp
```

---

## 17.2 Deterministic Gate

最少检查：

```text
Task Contract exists
required deliverable exists
final artifact exists
artifact file exists
hash current
no failed validation
no blocking error
acceptance criteria evaluated
```

---

## 17.3 Reviewer

复用 Reasonix read-only subagent。

Reviewer 输出：

```text
criterion
pass|fail|uncertain
evidence
issue
repair
```

Reviewer：

```text
不能直接把 Work 标 completed
```

---

## 17.4 Completion

只有 Host Gate 可执行：

```text
Work -> completed
```

模型调用：

```text
temper_complete_work
```

如果失败：

```text
返回具体 gap
保持 Work running/reviewing
继续修复
```

---

# 18. MILESTONE 10 — Temper Design System

UI 不参考 Reasonix 视觉。

只参考 Reasonix 功能和交互细节。

---

## 18.1 Visual Language

名称：

```text
Temper Forge
```

关键词：

```text
Quiet
Precise
Premium
Technical
Desktop-native
Dense but readable
```

---

## 18.2 Dark tokens

```css
--app-bg: #0B0D10;
--sidebar-bg: #0F1217;
--surface-1: #12161C;
--surface-2: #171C23;
--surface-3: #1D232C;

--border-soft: rgba(255,255,255,0.07);
--border-strong: rgba(255,255,255,0.13);

--text-1: #F3F5F7;
--text-2: #AAB2BD;
--text-3: #707A87;

--accent: #6D8CFF;
--accent-soft: rgba(109,140,255,0.15);

--success: #4CC38A;
--warning: #E7B75A;
--danger: #F16D76;
--info: #63B3ED;
```

Light theme 需要完整对应 tokens。

---

## 18.3 Typography

```text
Segoe UI Variable
Segoe UI
Microsoft YaHei UI
system-ui
```

Code：

```text
Cascadia Code
Consolas
```

---

## 18.4 Layout

Native Windows title bar。

主 UI：

```text
┌─────────────┬───────────────────────────────┬─────────────────┐
│ Navigation  │ Main Canvas                   │ Inspector       │
│             │                               │                 │
│ Home        │ Chat / Work / Project         │ Overview        │
│ Chat        │                               │ Context         │
│ Work        │                               │ Evidence        │
│ Projects    │                               │ Artifacts       │
│ Artifacts   │                               │ Runtime         │
│ Advanced    │                               │                 │
│ Settings    │                               │                 │
└─────────────┴───────────────────────────────┴─────────────────┘
│ Model · CTX · Cache · Stage · Tool · Tokens · Time             │
└──────────────────────────────────────────────────────────────────┘
```

尺寸：

```text
Navigation expanded: 220–236px
Navigation compact: 56px
Inspector: 320–360px
Runtime strip: 28–30px
```

---

## 18.5 Motion

```text
120–180 ms
```

只用于：

```text
hover
popover
panel
state transition
progress
```

禁止大量装饰 animation。

---

## 18.6 Components

至少：

```text
Button
IconButton
Input
Textarea
Select
Combobox
Tabs
Segmented
Badge
StatusDot
Tooltip
Popover
Menu
Dialog
Drawer
Toast
Progress
Skeleton
EmptyState
ErrorState
List
VirtualList
DataTable
CodeViewer
DiffViewer
ResizablePanel
CommandPalette
```

---

# 19. MILESTONE 11 — AppShell + Home

Home：

```text
New Work
New Chat
Continue Work
Active Works
Recent Projects
Recent Artifacts
Runtime Health
```

Runtime Health：

```text
Provider
Model
DB
MCP online
Skills
Context
Last error
```

没有真实 handler 的按钮：

```text
不能 enabled
```

绝对禁止：

```text
dead button
console.log only
coming soon fake action
```

---

# 20. MILESTONE 12 — Chat UI Parity

视觉重做。

行为完全复用 Reasonix。

必须覆盖：

```text
streaming
reasoning status
tool card
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
web search sources
rewind
session history
```

Timeline 不显示私有 chain-of-thought 正文。

允许：

```text
Thinking
Planning
Executing
Reviewing
```

---

# 21. MILESTONE 13 — Work UI

Work 页面：

```text
Header
Goal
Status
Task Contract
Plan/Todo
Execution Timeline
Transcript
Evidence
Decisions
Artifacts
Validation
Quality
Completion
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

运行状态必须来自真实 backend。

---

# 22. MILESTONE 14 — Projects + Artifacts UI

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

Artifacts：

```text
Search
Project filter
Work filter
Type
Final
Validation
Recent
Preview
Open
Reveal
Copy path
Set final
Continue Work
```

---

# 23. MILESTONE 15 — Settings + Advanced

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

直接读写 Reasonix actual config contract。

Temper 私有设置独立存储。

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

如果某能力仍主要由 CLI 操作：

```text
显示真实状态
显示 exact command
```

不要制作假的 GUI。

---

# 24. MILESTONE 16 — Runtime Observability

底部 Runtime Strip：

```text
Model
Mode
Context used/window
Compact threshold
Cache
Input/output tokens
Current tool
Work stage
Elapsed
```

只显示真实数据。

没有：

```text
显示 —
```

禁止虚假估算。

Context Inspector：

```text
context window
used
remaining
compact ratio
last compaction
history
memory recall
cache
```

---

# 25. MILESTONE 17 — Recovery / Memory / MCP / Skills Parity

逐条真实测试：

## Rewind

```text
modify file
checkpoint
rewind code
rewind conversation
rewind both
restart
```

## Memory

```text
remember
search
read
restart
forget
```

## MCP

使用本地 deterministic example：

```text
connect
tools/list
tool call
progress
reload/disconnect
```

## Skills

```text
discover
enable
invoke
disable
restart
```

---

# 26. MILESTONE 18 — i18n + Accessibility + Windows UX

语言：

```text
zh-CN
en-US
```

全部用户可见文字走 i18n。

Key parity test。

Accessibility：

```text
aria-label
focus
keyboard
dialogs
menus
command palette
contrast
```

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

Windows：

```text
Minimize
Maximize
Restore
Close
Resize
Native drag
IME
Clipboard
File picker
Explorer reveal
Dark/light system
```

---

# 27. MILESTONE 19 — Production-entry E2E

Mock frontend 不能作为完成证据。

必须验证：

```text
Temper.exe
Wails
embedded frontend
Go binding
Reasonix Controller
Temper DB
filesystem
restart
```

---

## 27.1 Chat E2E

确定性 fake provider：

```text
launch
configure test provider
create Project
Chat
stream
read file
writer permission
write file
diff
rewind
```

---

## 27.2 Work E2E

Workspace：

```text
temp-workspace/
  notes.md
```

任务：

```text
读取 notes.md
生成 summary.md
记录 Evidence
注册 Artifact
验证 Artifact
设为 Final
Completion Gate
```

验证：

```text
Work completed
summary.md exists
Evidence exists
Artifact exists
Final persists
restart
Project persists
Work persists
Session continues
Artifact persists
```

这是 Release Blocker。

---

# 28. MILESTONE 20 — Security Hardening

保持 Reasonix：

```text
Permission Gate
Sandbox
URL validation
secret redaction
MCP trust
path protections
```

Temper 新增测试：

```text
project path traversal
artifact escape
SQLite injection
secret in evidence
secret in decision
unsafe external open
invalid workspace
malformed Work data
```

---

# 29. MILESTONE 21 — Defender Hardening

这是 Release Blocker。

不关闭 Defender。

不依赖代码签名。

---

## 29.1 Source audit

搜索：

```text
Add-MpPreference
ExclusionPath
Set-ExecutionPolicy
CreateRemoteThread
schtasks
hidden PowerShell
process injection
credential dumping
startup persistence
registry persistence
download-and-execute
self modification
UPX
obfuscation
```

合法 Reasonix shell/remote 能力保留。

不要建立 v0.1 那种高风险字符串 blacklist。

---

## 29.2 Build raw EXE

先构建：

```text
Temper.exe
```

扫描它。

如果：

```text
Temper.exe detected
```

禁止继续打 installer。

先定位二进制/行为根因。

---

## 29.3 Portable

Raw EXE clean 后：

```text
ZIP
scan ZIP
extract
scan directory
```

都 clean 才通过。

---

## 29.4 Installer

只有：

```text
EXE clean
Portable clean
```

后才构建 installer。

Setup detected：

```text
不发布 Setup
```

不影响 Portable v0.3.0 Release。

---

# 30. MILESTONE 22 — Release Pipeline

CI：

```text
Go tests
Go vet
Reasonix compatibility tests
Desktop tests
Frontend tests
Frontend build
Wails Windows build
Temper E2E
```

GitHub CI 如果不能真实证明 Defender：

```text
DEFENDER_SCAN=NOT_RUN
```

禁止伪造 PASS。

---

# 31. RC

创建：

```text
v0.3.0-rc.1
```

RC 验收：

```text
CI green
Chat E2E
Work E2E
Restart recovery
MCP
Skills
Memory
Rewind
Windows UX
Temper.exe Defender clean
Portable Defender clean
SHA verified
```

失败：

```text
fix
test
commit
rc.2
```

禁止移动 tag。

---

# 32. Final Release

最新 RC 全部通过后：

```text
v0.3.0
```

最低 assets：

```text
Temper-0.3.0-windows-x64-portable.zip
Temper-0.3.0-SHA256SUMS.txt
```

可选：

```text
Temper-0.3.0-windows-x64-setup.exe
```

前提：

```text
Defender clean
```

---

# 33. Release Blockers

以下任一存在：

```text
Provider 配置失败
真实 Chat 无法跑完
Tool approval 不工作
文件操作不工作
Rewind 不工作
Project restart 丢失
Work 无法 completed
Artifact 无法打开
Completion Gate 可被模型绕过
MCP 基本测试失败
Skill 基本测试失败
Memory restart 失败
Context/compaction UI 是假数据
Home 有死按钮
Work 有死按钮
Windows controls 异常
生产 Temper.exe 无法启动
Temper.exe Defender detected
Portable Defender detected
P0/P1 regression
CI red
```

禁止正式 Release。

---

# 34. AGENTS.md 固定核心

AGENTS.md 控制在约 200 行内。

必须写：

```text
Temper v0.3 is Reasonix-runtime-first.
Old temper-cowork is read-only history.
Do not merge old Git histories.
Do not duplicate Reasonix Runtime.
Temper-owned code stays in internal/temper and frontend temper namespaces.
One milestone per DeepSeek session.
50K target / 90K hard context.
TDD/regression first.
No dead buttons.
No false DONE.
Windows-first.
Defender clean portable required.
```

---

# 35. BUILD_STATE

格式：

```markdown
# Temper v0.3.0 Build State

## Baseline
Reasonix repo:
Reasonix SHA:
Temper repo:
Main SHA:

## Current Milestone
ID:
Branch:
Status:

## Completed
- [x]

## In Progress
- [ ]

## Tests
- command:
  result:

## CI
run:
result:

## Windows
exe:
portable:
setup:
defender:

## Upstream patches
count:

## Next
Milestone:
Objective:
```

---

# 36. CURRENT_TASK

一次只写一个任务：

```markdown
# CURRENT_TASK

ID:

Objective:

Read:
- exact files

Modify:
- exact files

Do not touch:
- ...

Tests first:
- ...
- command

Acceptance:
- ...

Commit:
type(scope): message
```

---

# 37. 每个 Milestone 的固定闭环

```text
1. Read state
2. Create milestone branch
3. Inspect relevant Reasonix behavior
4. Write failing tests
5. Implement microtasks
6. Run target tests
7. Run regression tests
8. Run production smoke where applicable
9. Update parity
10. Update BUILD_STATE
11. git diff --check
12. commit
13. push
14. CI green
15. PR
16. merge main
17. STOP
```

---

# 38. 不允许的工作方式

禁止：

```text
一次把 v0.3.0 全做完
一口气改几十个文件
让 Context 接近 1M 再处理
每次 Session 读取完整 repo
模型自己宣布功能完成
页面只有视觉没有 backend
backend 有功能但 UI 假连接
前端生产路径使用 MockBridge
按钮点击只 console.log
关闭 Defender 测试发布
为了避报毒混淆/加壳
为了 CI 绿删除测试
全仓库替换 reasonix → temper
```

---

# 39. Definition of Done

## Runtime

```text
Provider              PASS
Models                PASS
Agent                 PASS
Plan                  PASS
Goal                  PASS
Subagent              PASS
Context               PASS
Compaction            PASS
Memory                PASS
History               PASS
Tools                 PASS
Permissions           PASS
Sandbox               PASS
MCP                   PASS
Skills                PASS
Extensions            AVAILABLE
Checkpoint/Rewind     PASS
Recovery              PASS
```

## CoWork

```text
Project               PASS
Chat                  PASS
Work                  PASS
Task Contract         PASS
Evidence              PASS
Decision              PASS
Artifact              PASS
Validation            PASS
Quality Gate          PASS
Completion Gate       PASS
Restart               PASS
```

## UI

```text
Temper visual identity
No dead actions
zh-CN
en-US
Accessibility
Responsive Windows layout
Runtime observability
```

## Release

```text
CI green
Production E2E
Windows smoke
Defender-clean EXE
Defender-clean Portable
RC
v0.3.0 Release
```

---

# 40. 第一次执行时只完成 MILESTONE 00 和 01

第一次 DeepSeek V4 Flash Session：

```text
MILESTONE 00
MILESTONE 01
```

完成后必须停止。

最终报告：

```text
TEMPER 0.3 BOOTSTRAP PASS

Temper repo:
Main SHA:
Reasonix baseline SHA:
Reasonix baseline tests:
Temper baseline tests:
Windows Wails build:
CI:
Next: MILESTONE 02
```

禁止第一次 Session 直接开始：

```text
Identity
CoWork
UI
Release
```

等下一个新 Session 再继续。

---

# 41. 现在开始

执行顺序：

```text
M00 Baseline
→ M01 Clean Repository
→ STOP

M02 Isolation
→ M03 Parity
→ STOP

M04 CoWork Store
→ M05 Project/Chat
→ STOP

M06 Work
→ M07 CoWork Tools
→ M08 Artifacts
→ M09 Quality
→ STOP

M10 Design System
→ M11 Home
→ M12 Chat
→ STOP

M13 Work UI
→ M14 Project/Artifacts
→ M15 Settings/Advanced
→ M16 Runtime
→ STOP

M17 Recovery/MCP/Skills/Memory
→ M18 i18n/Windows UX
→ STOP

M19 E2E
→ M20 Security
→ M21 Defender
→ STOP

M22 Release Pipeline
→ RC
→ v0.3.0
```

最终目标：

```text
一个全新的
干净的
可维护的
Reasonix-runtime-first
真正可以完整跑通工作的
具有独立 Temper UI 的
Windows CoWork 桌面应用。
```
