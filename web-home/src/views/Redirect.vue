<template>
  <div class="redirect-page">
    <div class="redirect-box">
      <!-- Logo -->
      <div class="box-logo">
        <div class="logo-mark">
          <svg viewBox="0 0 24 24" fill="none" width="22" height="22">
            <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <span class="logo-name">LNKR</span>
      </div>

      <template v-if="!error">
        <!-- Ink loading animation -->
        <div class="loader-wrap">
          <div class="ink-bar"></div>
        </div>

        <div class="redirect-status">
          <span class="status-label">{{ t('redirect.resolving') }}</span>
        </div>

        <div class="redirect-code-display">
          <span class="code-prefix">/r/</span><span class="code-val">{{ route.params.code }}</span>
        </div>

        <p class="redirect-hint">{{ t('redirect.hintText') }}<span class="dot-anim"></span></p>
      </template>

      <template v-else>
        <div class="error-cross">✕</div>
        <div class="error-title">{{ t('redirect.errorTitle') }}</div>
        <p class="error-msg">{{ error }}</p>
        <button class="lk-btn lk-btn-primary" style="margin-top: 24px" @click="router.push('/')">
          {{ t('common.goHome') }}
        </button>
      </template>
    </div>

    <!-- Decorative ruled lines -->
    <div class="deco-lines">
      <div class="deco-line" v-for="i in 8" :key="i"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { getRedirectUrl } from '@/api/shortlink';

const { t }  = useI18n();
const route  = useRoute();
const router = useRouter();
const error  = ref('');

onMounted(async () => {
  try {
    const res = await getRedirectUrl(route.params.code as string);
    window.location.href = res.data.original_url;
  } catch (e: any) {
    error.value = e.message || t('stats.notFoundText');
  }
});
</script>

<style scoped>
.redirect-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--paper);
  position: relative;
  overflow: hidden;
}

/* ─── Decorative ruled lines ─────────────── */
.deco-lines {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  pointer-events: none;
}
.deco-line {
  width: 100%;
  height: 1px;
  background: rgba(28, 22, 13, 0.04);
}

/* ─── Box ──────────────────────────────── */
.redirect-box {
  position: relative;
  z-index: 1;
  width: 380px;
  background: var(--white);
  border: var(--rule-heavy);
  padding: 48px 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  box-shadow: 8px 8px 0 rgba(28, 22, 13, 0.06);
}

/* ─── Logo ──────────────────────────────── */
.box-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 32px;
}
.logo-mark {
  width: 38px;
  height: 38px;
  border: var(--rule-heavy);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ink);
}
.logo-name {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.12em;
  color: var(--ink);
}

/* ─── Ink loader ─────────────────────────── */
.loader-wrap {
  width: 100%;
  height: 3px;
  background: var(--paper-warm);
  margin-bottom: 28px;
  overflow: hidden;
}
.ink-bar {
  height: 100%;
  background: var(--ink);
  animation: ink-fill 1.4s ease-in-out infinite;
}
@keyframes ink-fill {
  0%   { width: 0; margin-left: 0; }
  50%  { width: 70%; margin-left: 0; }
  100% { width: 0; margin-left: 100%; }
}

/* ─── Status ─────────────────────────────── */
.redirect-status { margin-bottom: 20px; }
.status-label {
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.22em;
  color: var(--ink-faint);
  text-transform: uppercase;
  border: 1px solid var(--ink-faint);
  padding: 4px 10px;
}

.redirect-code-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  background: var(--paper-warm);
  border: var(--rule-medium);
  padding: 12px 20px;
  width: 100%;
  margin-bottom: 16px;
}
.code-prefix {
  font-family: var(--font-mono);
  font-size: 14px;
  color: var(--ink-faint);
}
.code-val {
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 500;
  color: var(--ink);
  letter-spacing: 0.06em;
}

.redirect-hint {
  font-size: 13.5px;
  color: var(--ink-light);
}
.dot-anim::after {
  content: '...';
  animation: dots 1.4s steps(4, end) infinite;
}
@keyframes dots {
  0%  { content: '';    }
  25% { content: '.';   }
  50% { content: '..';  }
  75% { content: '...'; }
}

/* ─── Error state ────────────────────────── */
.error-cross {
  font-size: 36px;
  color: var(--red);
  margin-bottom: 12px;
  font-weight: 300;
}
.error-title {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 700;
  color: var(--ink);
  margin-bottom: 8px;
}
.error-msg {
  font-size: 13.5px;
  color: var(--ink-light);
  line-height: 1.6;
}
</style>
