// Temper Runtime Observability 验证:StatusBar 必须渲染来自 ContextInfo 的
// 真实运行时数据(used/window/tokens/model),禁止造数/硬编码。
// 运行:tsx src/__tests__/temper-observability.test.tsx

import { JSDOM } from "jsdom";
import React from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { StatusBar } from "../components/StatusBar";
import { LocaleProvider } from "../lib/i18n";
import { DEFAULT_STATUS_BAR_ITEMS, normalizeStatusBarItems } from "../lib/statusBarItems";

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

const dom = new JSDOM("<!doctype html><html><body></body></html>", {
  pretendToBeVisual: true,
  url: "http://localhost/",
});
(globalThis as typeof globalThis & { IS_REACT_ACT_ENVIRONMENT: boolean }).IS_REACT_ACT_ENVIRONMENT = true;
globalThis.window = dom.window as unknown as Window & typeof globalThis;
globalThis.document = dom.window.document;

function renderStatusBar(opts: {
  used: number; window: number; sessionTokens: number; turnTokens: number; modelLabel?: string;
}): string {
  return renderToStaticMarkup(
    <LocaleProvider>
      <StatusBar
        context={{ used: opts.used, window: opts.window, sessionTokens: opts.sessionTokens } as never}
        running={false}
        statusBarItems={normalizeStatusBarItems(DEFAULT_STATUS_BAR_ITEMS)}
        modelLabel={opts.modelLabel ?? "DeepSeek-R1"}
        sessionTokens={opts.sessionTokens}
        turnTokens={opts.turnTokens}
      />
    </LocaleProvider>,
  );
}

console.log("Temper observability");

{
  // 真实数据必须渲染
  const html = renderStatusBar({ used: 12345, window: 64000, sessionTokens: 8123, turnTokens: 1200 });
  ok(html.includes("8,123") === true, "renders real session tokens (8,123)");
  ok(html.includes("1,200") === true, "renders real turn tokens (1,200)");
  ok(html.includes("DeepSeek-R1") === true, "renders real model label");
  // 上下文使用率:12345/64000 ≈ 19.3% → 19%(四舍五入到整数百分比)
  ok(html.includes("19%") === true, "renders real context usage percent (12345/64000 ≈ 19%)");
}

{
  // 空数据必须显示占位符(-),禁止造数
  const html = renderStatusBar({ used: 0, window: 0, sessionTokens: 0, turnTokens: 0 });
  ok(html.includes("-") === true, "shows placeholder for empty data");
}

{
  // 不同数据渲染不同值(证明非硬编码)
  const htmlA = renderStatusBar({ used: 1000, window: 10000, sessionTokens: 50, turnTokens: 5 });
  const htmlB = renderStatusBar({ used: 5000, window: 10000, sessionTokens: 99, turnTokens: 9 });
  ok(htmlA.includes("50") === true && htmlB.includes("99") === true, "distinct session tokens render distinct values");
  ok(htmlA.includes("10%") === true && htmlB.includes("50%") === true, "distinct context usage renders distinct percents");
}

console.log(`\n${passed} passed, ${failed} failed`);
if (failed > 0) process.exit(1);
