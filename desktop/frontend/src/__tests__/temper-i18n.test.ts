// Temper i18n 验证:en / zh / zh-TW 字典 key 完全一致,且 Temper 专属
// 文案(temper.*)在三种语言中都存在。运行:tsx src/__tests__/temper-i18n.test.tsx

import { en } from "../locales/en";
import { zh } from "../locales/zh";
import { zhTW } from "../locales/zh-TW";

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

const enKeys = Object.keys(en).sort();
const zhKeys = Object.keys(zh).sort();
const zhTWKeys = Object.keys(zhTW).sort();

console.log("Temper i18n");

// zh 与 en key 完全一致(无缺无多)
{
  const missing = zhKeys.filter((k) => !enKeys.includes(k));
  const extra = enKeys.filter((k) => !zhKeys.includes(k));
  ok(missing.length === 0, `zh has no extra keys (${missing.length} extra)`);
  ok(extra.length === 0, `zh covers all en keys (${extra.length} missing)`);
  ok(zhKeys.length === enKeys.length, `zh key count ${zhKeys.length} == en ${enKeys.length}`);
}

// zh-TW 与 en key 完全一致(Record<DictKey> 编译期已保证,这里再运行时验证)
{
  ok(zhTWKeys.length === enKeys.length, `zh-TW key count ${zhTWKeys.length} == en ${enKeys.length}`);
}

// Temper 专属文案三种语言都有
{
  const temperKeys = enKeys.filter((k) => k.startsWith("temper."));
  ok(temperKeys.length >= 7, `temper.* keys present (${temperKeys.length})`);
  for (const key of temperKeys) {
    ok(zh[key] !== undefined && zh[key] !== "", `zh has ${key}`);
    ok(zhTW[key] !== undefined && zhTW[key] !== "", `zh-TW has ${key}`);
  }
}

// 值非空(不允许空字符串占位)
{
  const emptyEn = enKeys.filter((k) => en[k as keyof typeof en] === "");
  ok(emptyEn.length === 0, `en has no empty values (${emptyEn.length})`);
}

console.log(`\n${passed} passed, ${failed} failed`);
if (failed > 0) process.exit(1);
