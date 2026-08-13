# PRODUCT.md — Temper 产品定义

## 产品

- 名称:Temper
- 版本:0.3.0-dev
- 标语:**Shape intent. Ship work.**
- 形态:Windows 桌面应用(React + TypeScript + Wails),最终发行 MSIX。
- 目标用户:需要把"意图"变成"交付物"的专业工作者(开发者、分析师、文档作者等)。

## 核心价值

1. **CoWork**:人与 Agent 在同一工作区协作完成正式工作(Work),而非一次性问答。
2. **正式化**:Work 有 Task Contract、Acceptance Criteria、Evidence、Decision、Artifact、Validation、Quality Gate、Completion Gate。
3. **可追溯**:每个结论有 Evidence 与 Decision;每个交付有 Artifact 与校验。
4. **可信赖**:Permission/Sandbox 管控工具动作,重启后状态恢复。

## 功能范围(非最终)

- Project(工作区注册/管理)
- Chat(Quick Chat / Project Chat,继承 Reasonix Session 全部能力)
- Formal Work(生命周期 + Task Contract 编译)
- Evidence / Decision / Artifact / Validation / Quality / Completion
- Provider / Model / Tools / MCP / Skills / Memory / Rewind / Recovery(继承 Reasonix)
- i18n:zh-CN / en-US

## 非目标

- 不做第二套 Runtime / Agent / Planner(全部继承 Reasonix)。
- 不复制 Reasonix 视觉 UI(Temper 独立 Design System)。
- v0.3.0 不要求商业代码签名证书(自签名 CN=Temper Development)。
