/**
 * i18n setup — 扩展新语言步骤：
 * 1. 新建 src/locales/<lang>.ts（参照 zh-CN.ts 结构）
 * 2. 在 messages 对象中注册
 * 3. 在 SUPPORTED_LANGS 中添加 label
 */
import { createI18n } from 'vue-i18n';
import zhCN from './zh-CN';
import enUS from './en-US';

export type LangCode = 'zh-CN' | 'en-US';

/** 已支持的语言列表 — 新增语言只需在此追加 */
export const SUPPORTED_LANGS: Record<LangCode, { label: string; short: string }> = {
  'zh-CN': { label: '简体中文', short: '中文' },
  'en-US': { label: 'English', short: 'EN'   },
};

const STORAGE_KEY = 'lnkr-lang';

function detectLang(): LangCode {
  const saved = localStorage.getItem(STORAGE_KEY) as LangCode | null;
  if (saved && saved in SUPPORTED_LANGS) return saved;
  // Browser preference
  const browser = navigator.language;
  if (browser.startsWith('zh')) return 'zh-CN';
  return 'en-US';
}

const i18n = createI18n({
  legacy: false,          // Composition API mode
  locale: detectLang(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
    // Add more locales here: 'ja-JP': jaJP, 'ko-KR': koKR ...
  },
  missingWarn:  false,
  fallbackWarn: false,
});

/** Save lang preference and update i18n locale */
export function setLang(lang: LangCode) {
  localStorage.setItem(STORAGE_KEY, lang);
  (i18n.global.locale as any).value = lang;
}

export function getCurrentLang(): LangCode {
  return (i18n.global.locale as any).value as LangCode;
}

export default i18n;
