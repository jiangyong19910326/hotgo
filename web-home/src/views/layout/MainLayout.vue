<template>
  <div class="app-shell">
    <!-- ─── Header ─── -->
    <header class="site-header">
      <div class="header-inner">
        <button class="logo-btn" @click="router.push('/')">
          <span class="logo-rule"></span>
          <span class="logo-text">LNKR</span>
        </button>

        <nav class="site-nav">
          <button
            v-for="item in navItems"
            :key="item.path"
            class="nav-item"
            :class="{ active: route.path === item.path }"
            @click="router.push(item.path)"
          >
            {{ t(item.i18nKey) }}
          </button>
        </nav>

        <div class="header-right">
          <!-- Language switcher -->
          <div class="lang-switcher">
            <button
              v-for="(meta, code) in SUPPORTED_LANGS"
              :key="code"
              class="lang-btn"
              :class="{ active: locale === code }"
              @click="switchLang(code as LangCode)"
            >
              {{ meta.short }}
            </button>
          </div>

          <span class="issue-num">Vol. 01</span>
        </div>
      </div>
    </header>

    <!-- ─── Main ─── -->
    <main class="site-main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- ─── Footer ─── -->
    <footer class="site-footer">
      <div class="footer-inner">
        <span class="footer-brand">LNKR</span>
        <div class="footer-rule"></div>
        <span class="footer-copy">{{ t('footer.service') }} · Powered by HotGo</span>
        <div class="footer-rule"></div>
        <span class="footer-copy">{{ new Date().getFullYear() }}</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { SUPPORTED_LANGS, setLang, type LangCode } from '@/locales';

const router = useRouter();
const route  = useRoute();
const { t, locale } = useI18n();

const navItems = [
  { path: '/',     i18nKey: 'nav.home'    },
  { path: '/list', i18nKey: 'nav.library' },
];

function switchLang(lang: LangCode) {
  setLang(lang);
}
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--paper);
}

/* ─── Header ─────────────────────────────── */
.site-header {
  border-bottom: var(--rule-heavy);
  background: var(--paper);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-inner {
  max-width: 1100px;
  margin: 0 auto;
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 32px;
  gap: 32px;
}

/* ─── Logo ───────────────────────────────── */
.logo-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  transition: var(--t);
  flex-shrink: 0;
}
.logo-btn:hover .logo-text { color: var(--vermillion); }

.logo-rule {
  display: block;
  width: 3px;
  height: 22px;
  background: var(--ink);
  transition: var(--t);
}
.logo-btn:hover .logo-rule { background: var(--vermillion); }

.logo-text {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 0.12em;
  color: var(--ink);
  transition: var(--t);
}

/* ─── Nav ────────────────────────────────── */
.site-nav {
  display: flex;
  gap: 2px;
  flex: 1;
}

.nav-item {
  background: none;
  border: none;
  cursor: pointer;
  font-family: var(--font-body);
  font-size: 13px;
  font-weight: 400;
  letter-spacing: 0.04em;
  color: var(--ink-light);
  padding: 8px 14px;
  position: relative;
  transition: var(--t);
}
.nav-item::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0; right: 0;
  height: 2px;
  background: var(--ink);
  transform: scaleX(0);
  transition: transform 0.2s var(--ease);
}
.nav-item:hover { color: var(--ink); }
.nav-item:hover::after { transform: scaleX(1); }
.nav-item.active { color: var(--ink); font-weight: 500; }
.nav-item.active::after { transform: scaleX(1); background: var(--vermillion); }

/* ─── Right ──────────────────────────────── */
.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 12px;
}

/* ─── Language switcher ──────────────────── */
.lang-switcher {
  display: flex;
  border: var(--rule-medium);
  overflow: hidden;
}
.lang-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.08em;
  color: var(--ink-light);
  padding: 5px 10px;
  transition: var(--t);
  border-right: var(--rule-faint);
  line-height: 1;
}
.lang-btn:last-child { border-right: none; }
.lang-btn:hover { color: var(--ink); background: var(--paper-warm); }
.lang-btn.active {
  background: var(--ink);
  color: var(--paper);
}

.issue-num {
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: 0.1em;
  color: var(--ink-faint);
  border: 1px solid var(--ink-faint);
  padding: 3px 8px;
}

/* ─── Main ───────────────────────────────── */
.site-main { flex: 1; }

/* ─── Footer ─────────────────────────────── */
.site-footer {
  border-top: var(--rule-heavy);
  padding: 20px 32px;
  background: var(--paper);
}
.footer-inner {
  max-width: 1100px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  gap: 16px;
}
.footer-brand {
  font-family: var(--font-display);
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.12em;
  color: var(--ink);
}
.footer-rule { flex: 1; height: 1px; background: var(--ink-faint); }
.footer-copy { font-size: 12px; color: var(--ink-faint); letter-spacing: 0.04em; white-space: nowrap; }
</style>
