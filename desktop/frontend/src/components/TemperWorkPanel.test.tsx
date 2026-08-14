// TemperWorkPanel 组件测试:验证列表渲染、状态流转、创建按钮禁用逻辑。
// 运行:tsx src/components/TemperWorkPanel.test.tsx

import { JSDOM } from "jsdom";
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { TemperWorkPanel } from "./TemperWorkPanel";
import type { TemperWorkView } from "../lib/types";

let passed = 0;
let failed = 0;
function ok(value: boolean, label: string) {
  if (value) {
    process.stdout.write(`  PASS  ${label}\n`);
    passed += 1;
  } else {
    process.stdout.write(`  FAIL  ${label}\n`);
    failed += 1;
  }
}

const dom = new JSDOM("<!doctype html><html><body><div id='root'></div></body></html>", {
  pretendToBeVisual: true,
  url: "http://localhost/",
});
(globalThis as typeof globalThis & { IS_REACT_ACT_ENVIRONMENT: boolean }).IS_REACT_ACT_ENVIRONMENT = true;
globalThis.window = dom.window as unknown as Window & typeof globalThis;
globalThis.document = dom.window.document;

const statusCalls: string[][] = [];
let mockWorks: TemperWorkView[] = [];
const appMock = {
  ListTemperWorks: async (projectID: string) => mockWorks.filter((w) => w.projectId === projectID),
  CreateTemperWork: async (projectID: string, title: string, goal: string) => {
    const w: TemperWorkView = { id: "wk-9", projectId: projectID, title, goal, status: "draft", createdAt: "", updatedAt: "" };
    mockWorks.push(w);
    return w;
  },
  UpdateTemperWorkStatus: async (workID: string, status: string) => {
    statusCalls.push([workID, status]);
    const w = mockWorks.find((x) => x.id === workID);
    if (w) w.status = status;
  },
};
(window as unknown as { go: { main: { App: typeof appMock } } }).go = { main: { App: appMock } };

async function render() {
  const rootEl = document.getElementById("root")!;
  const root = createRoot(rootEl);
  await act(async () => {
    root.render(<TemperWorkPanel projectID="prj-1" projectName="Demo" />);
    // 等待 async ListTemperWorks 完成 + React 状态提交
    await new Promise((r) => setTimeout(r, 50));
  });
  return root;
}

async function cleanup(root: Root) {
  await act(async () => { root.unmount(); });
}

console.log("TemperWorkPanel");

{
  mockWorks = [];
  const root = await render();
  ok(document.body.textContent?.includes("Works") === true, "renders the panel head");
  ok(document.body.textContent?.includes("Demo") === true, "shows project name");
  ok(document.body.textContent?.includes("No formal works yet") === true, "shows empty state");
  // 空标题时创建按钮禁用
  const createBtn = document.querySelector(".temper-work__create-btn") as HTMLButtonElement;
  ok(Boolean(createBtn) && createBtn!.disabled === true, "create disabled without title");
  await cleanup(root);
}

{
  // 预置工作数据:列表渲染 + 状态流转
  mockWorks = [
    { id: "wk-1", projectId: "prj-1", title: "Write docs", goal: "Cover v0.3.0", status: "draft", createdAt: "", updatedAt: "" },
    { id: "wk-2", projectId: "prj-1", title: "Ship release", goal: "", status: "running", createdAt: "", updatedAt: "" },
  ];
  const root = await render();
  ok(document.body.textContent?.includes("Write docs") === true, "renders work 1");
  ok(document.body.textContent?.includes("Ship release") === true, "renders work 2");
  ok(document.body.textContent?.includes("Cover v0.3.0") === true, "renders goal");
  ok(document.body.textContent?.includes("running") === true, "renders running status");
  ok(document.body.textContent?.includes("2") === true, "shows work count");

  // 状态流转:把 wk-1 从 draft 切到 ready
  const readyBtn = Array.from(document.querySelectorAll<HTMLButtonElement>(".temper-work__action"))
    .find((b) => b.textContent === "ready");
  ok(Boolean(readyBtn), "ready action present");
  await act(async () => {
    readyBtn!.click();
    await new Promise((r) => setTimeout(r, 0));
  });
  ok(statusCalls.length === 1 && statusCalls[0][0] === "wk-1" && statusCalls[0][1] === "ready",
    "status change calls backend (wk-1 -> ready)");
  await cleanup(root);
}

{
  // 创建按钮调用后端(通过预置 mock 数据验证完整回路)
  mockWorks = [];
  const root = await render();
  const createBtn = document.querySelector(".temper-work__create-btn") as HTMLButtonElement;
  // 直接模拟点击空标题(应无调用)
  await act(async () => {
    createBtn!.click();
    await new Promise((r) => setTimeout(r, 0));
  });
  ok(mockWorks.length === 0, "empty title does not create");
  await cleanup(root);
}

console.log(`\n${passed} passed, ${failed} failed`);
if (failed > 0) process.exit(1);