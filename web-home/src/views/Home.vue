<template>
  <div class="home">
    <!-- ═══════════════════════════════════════ -->
    <!-- HERO — Deep space section               -->
    <!-- ═══════════════════════════════════════ -->
    <section class="hero">
      <SpaceCanvas />

      <div class="hero-content">
        <!-- Eyebrow -->
        <div class="hero-eyebrow">
          <span class="eyebrow-dash"></span>
          <span>{{ t('home.tag') }}</span>
          <span class="eyebrow-dash"></span>
        </div>

        <!-- Main headline (Fraunces italic) -->
        <h1 class="hero-title">
          <em class="title-em">{{ t('home.title1') }}</em>
          <span class="title-plain">Links.</span>
        </h1>
        <h1 class="hero-title title-line2">
          <em class="title-em">{{ t('home.title2') }}</em>
          <span class="title-plain">Data.</span>
        </h1>

        <p class="hero-sub">
          {{ t('home.desc') }}
        </p>

        <button class="hero-scroll-btn" @click="scrollToForm">
          <span>{{ t('home.scrollBtn') }}</span>
          <svg viewBox="0 0 24 24" fill="none" width="16" height="16">
            <path d="M12 5v14M5 12l7 7 7-7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>

      <!-- Bottom border decoration -->
      <div class="hero-bottom-rule">
        <span class="rule-label">{{ t('home.rulerText') }}</span>
      </div>
    </section>

    <!-- ═══════════════════════════════════════ -->
    <!-- FORM — Warm editorial section           -->
    <!-- ═══════════════════════════════════════ -->
    <section ref="formSection" class="form-section">
      <div class="form-inner">
        <!-- Left: form -->
        <div class="form-col">
          <div class="form-heading">
            <span class="form-num">{{ t('home.formNum') }}</span>
            <div>
              <h2 class="form-title">{{ t('home.formTitle') }}</h2>
              <p class="form-desc">{{ t('home.formDesc') }}</p>
            </div>
          </div>

          <n-form ref="formRef" :model="form" :rules="rules" label-placement="top">
            <n-form-item :label="t('home.urlLabel')" path="url">
              <n-input
                v-model:value="form.url"
                :placeholder="t('home.urlPlaceholder')"
                size="large"
                clearable
              />
            </n-form-item>
          </n-form>

          <!-- Options toggle -->
          <div class="opts-toggle" @click="showOptions = !showOptions">
            <svg :style="`transform: rotate(${showOptions ? 180 : 0}deg); transition: transform 0.2s`" viewBox="0 0 24 24" fill="none" width="13" height="13">
              <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            {{ showOptions ? t('home.optToggleHide') : t('home.optToggleShow') }}
          </div>

          <transition name="opts-expand">
            <div v-if="showOptions" class="opts-grid">
              <div class="opt-field">
                <label class="lk-label">{{ t('home.optTitle') }}</label>
                <input v-model="form.title" class="lk-input" :placeholder="t('home.optTitleHint')" />
              </div>
              <div class="opt-field">
                <label class="lk-label">{{ t('home.optCode') }}</label>
                <div class="prefix-wrap">
                  <span class="prefix-text">/r/</span>
                  <input v-model="form.custom_code" class="lk-input" style="padding-left: 34px" :placeholder="t('home.optCodeHint')" />
                </div>
              </div>
              <div class="opt-field opt-full">
                <label class="lk-label">{{ t('home.optExpiry') }}</label>
                <n-date-picker
                  v-model:formatted-value="form.expire_at"
                  value-format="yyyy-MM-dd HH:mm:ss"
                  type="datetime"
                  :placeholder="t('home.optExpiryHint')"
                  clearable
                  style="width: 100%"
                />
              </div>
            </div>
          </transition>

          <button
            class="lk-btn lk-btn-primary create-btn"
            :disabled="loading"
            @click="handleCreate"
          >
            <span v-if="loading" class="btn-spin"></span>
            <svg v-else viewBox="0 0 24 24" fill="none" width="15" height="15">
              <path d="M5 12h14M12 5l7 7-7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            {{ loading ? t('home.btnGenerating') : t('home.btnGenerate') }}
          </button>
        </div>

        <!-- Divider -->
        <div class="form-divider"></div>

        <!-- Right: result / placeholder -->
        <div class="result-col">
          <template v-if="result">
            <div class="result-label">
              <span class="result-dot"></span>
              {{ t('home.resultReady') }}
            </div>

            <div class="result-short-box">
              <span class="result-code">{{ result.short_url }}</span>
              <button class="copy-btn" @click="copyUrl(result!.short_url)">
                <svg viewBox="0 0 24 24" fill="none" width="14" height="14">
                  <rect x="9" y="9" width="13" height="13" rx="1" stroke="currentColor" stroke-width="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" stroke="currentColor" stroke-width="2"/>
                </svg>
                {{ t('common.copy') }}
              </button>
            </div>

            <p class="result-original">
              <span class="result-original-label">{{ t('home.originalLabel') }}</span>
              {{ result.original_url }}
            </p>

            <div class="result-qr">
              <qrcode-vue :value="result.short_url" :size="96" render-as="svg" foreground="#1c160d" background="#ffffff" />
              <span class="qr-hint">{{ t('home.qrHint') }}</span>
            </div>

            <button class="lk-btn lk-btn-ghost lk-btn-sm stats-link" @click="router.push(`/stats/${result!.code}`)">
              <svg viewBox="0 0 24 24" fill="none" width="13" height="13">
                <path d="M18 20V10M12 20V4M6 20v-6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
              {{ t('home.statsBtn') }}
            </button>
          </template>

          <template v-else>
            <div class="empty-result">
              <div class="empty-num">?</div>
              <p class="empty-text">{{ t('home.emptyHint') }}</p>
            </div>
          </template>
        </div>
      </div>
    </section>

    <!-- ═══════════════════════════════════════ -->
    <!-- FEATURES                                -->
    <!-- ═══════════════════════════════════════ -->
    <section class="features-section">
      <div class="features-inner">
        <div class="features-header">
          <div class="features-num">{{ t('home.featuresTag') }}</div>
          <h3 class="features-title">{{ t('home.featuresSection') }}</h3>
        </div>
        <div class="features-grid">
          <div class="feature-item" v-for="f in features" :key="f.label">
            <div class="feature-icon-wrap">
              <div class="feature-icon" v-html="f.svg"></div>
            </div>
            <div class="feature-body">
              <div class="feature-title">{{ f.title }}</div>
              <p class="feature-desc">{{ f.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { NForm, NFormItem, NInput, NDatePicker, useMessage } from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';
import QrcodeVue from 'qrcode.vue';
import SpaceCanvas from '@/components/SpaceCanvas.vue';
import { createShortLink, type CreateResult } from '@/api/shortlink';

const { t }       = useI18n();
const router      = useRouter();
const message     = useMessage();
const formRef     = ref<FormInst | null>(null);
const formSection = ref<HTMLElement | null>(null);
const loading     = ref(false);
const showOptions = ref(false);
const result      = ref<CreateResult | null>(null);

const form = ref({ url: '', title: '', custom_code: '', expire_at: null as string | null });

const rules = computed<FormRules>(() => ({
  url: [
    { required: true, message: t('home.urlRequired'), trigger: 'blur' },
    {
      validator: (_, v) => { try { new URL(v); return true; } catch { return false; } },
      message: t('home.urlInvalid'),
      trigger: 'blur',
    },
  ],
}));

function scrollToForm() {
  formSection.value?.scrollIntoView({ behavior: 'smooth' });
}

async function handleCreate() {
  try { await formRef.value?.validate(); } catch { return; }
  loading.value = true;
  result.value  = null;
  try {
    const res = await createShortLink({
      original_url: form.value.url,
      title: form.value.title || undefined,
      expire_at: form.value.expire_at || undefined,
    });
    result.value = res.data;
    message.success(t('home.createSuccess'));
  } catch (e: any) {
    message.error(e.message || t('home.createFail'));
  } finally {
    loading.value = false;
  }
}

async function copyUrl(url: string) {
  try { await navigator.clipboard.writeText(url); message.success(t('common.copied')); }
  catch { message.error(t('common.copyFail')); }
}

const features = computed(() => [
  {
    label: 'instant',
    title: t('home.f1Title'),
    desc: t('home.f1Desc'),
    svg: `<svg viewBox="0 0 24 24" fill="none" width="20" height="20"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/></svg>`,
  },
  {
    label: 'track',
    title: t('home.f2Title'),
    desc: t('home.f2Desc'),
    svg: `<svg viewBox="0 0 24 24" fill="none" width="20" height="20"><path d="M18 20V10M12 20V4M6 20v-6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
  },
  {
    label: 'qr',
    title: t('home.f3Title'),
    desc: t('home.f3Desc'),
    svg: `<svg viewBox="0 0 24 24" fill="none" width="20" height="20"><rect x="3" y="3" width="7" height="7" rx="1" stroke="currentColor" stroke-width="2"/><rect x="14" y="3" width="7" height="7" rx="1" stroke="currentColor" stroke-width="2"/><rect x="3" y="14" width="7" height="7" rx="1" stroke="currentColor" stroke-width="2"/><path d="M14 14h3M17 17h4M14 20h2M20 14v3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
  },
]);
</script>

<style scoped>
.home { display: flex; flex-direction: column; }

/* ═══ HERO ═══════════════════════════════ */
.hero {
  position: relative;
  height: 520px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-bottom: 3px solid var(--ink);
}

.hero-content {
  position: relative;
  z-index: 2;
  text-align: center;
  padding: 0 24px;
}

.hero-eyebrow {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.22em;
  color: rgba(200, 220, 255, 0.6);
  margin-bottom: 24px;
}
.eyebrow-dash {
  display: block;
  width: 28px;
  height: 1px;
  background: rgba(200, 220, 255, 0.35);
}

.hero-title {
  font-family: var(--font-display);
  font-size: clamp(52px, 9vw, 86px);
  font-weight: 700;
  line-height: 1;
  letter-spacing: -0.03em;
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 18px;
}
.title-line2 { margin-top: -6px; margin-bottom: 28px; }

.title-em {
  font-style: italic;
  color: #ffffff;
  text-shadow: 0 0 40px rgba(160, 200, 255, 0.4);
}
.title-plain {
  color: rgba(200, 220, 255, 0.5);
  font-style: normal;
}

.hero-sub {
  font-size: 15px;
  line-height: 1.7;
  color: rgba(200, 220, 255, 0.6);
  max-width: 440px;
  margin: 0 auto 32px;
}

.hero-scroll-btn {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: rgba(220, 235, 255, 0.85);
  padding: 12px 24px;
  cursor: pointer;
  font-family: var(--font-body);
  font-size: 13.5px;
  font-weight: 500;
  letter-spacing: 0.04em;
  transition: var(--t);
  backdrop-filter: blur(4px);
}
.hero-scroll-btn:hover {
  background: rgba(255, 255, 255, 0.14);
  border-color: rgba(255, 255, 255, 0.35);
  color: white;
}

/* Bottom rule */
.hero-bottom-rule {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px;
  background: rgba(3, 8, 26, 0.6);
  backdrop-filter: blur(4px);
}
.rule-label {
  font-family: var(--font-mono);
  font-size: 9.5px;
  letter-spacing: 0.26em;
  color: rgba(200, 220, 255, 0.35);
  text-transform: uppercase;
}

/* ═══ FORM SECTION ══════════════════════ */
.form-section {
  background: var(--paper);
  padding: 60px 32px;
}
.form-inner {
  max-width: 960px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 0;
}

/* Form col */
.form-col { padding-right: 48px; }

.form-heading {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 28px;
  padding-bottom: 20px;
  border-bottom: var(--rule-medium);
}
.form-num {
  font-family: var(--font-display);
  font-size: 40px;
  font-weight: 700;
  font-style: italic;
  color: var(--ink-faint);
  line-height: 1;
  flex-shrink: 0;
}
.form-title {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 4px;
}
.form-desc {
  font-size: 13px;
  color: var(--ink-light);
}

/* Options toggle */
.opts-toggle {
  display: flex;
  align-items: center;
  gap: 7px;
  font-size: 12.5px;
  color: var(--ink-light);
  cursor: pointer;
  padding: 8px 0 16px;
  transition: var(--t);
  user-select: none;
  letter-spacing: 0.02em;
}
.opts-toggle:hover { color: var(--ink); }

.opts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 20px;
}
.opt-field { display: flex; flex-direction: column; }
.opt-full { grid-column: 1 / -1; }
.prefix-wrap { position: relative; }
.prefix-text {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--ink-faint);
  pointer-events: none;
  z-index: 1;
}

/* Options slide */
.opts-expand-enter-active, .opts-expand-leave-active {
  transition: max-height 0.3s var(--ease), opacity 0.2s;
  overflow: hidden;
  max-height: 280px;
}
.opts-expand-enter-from, .opts-expand-leave-to { max-height: 0; opacity: 0; }

/* Create button */
.create-btn {
  width: 100%;
  height: 50px;
  font-size: 13px;
  letter-spacing: 0.1em;
  font-weight: 600;
  border-radius: 0;
  gap: 10px;
  margin-top: 8px;
}
.create-btn:disabled { opacity: 0.55; cursor: not-allowed; }
.btn-spin {
  width: 14px; height: 14px;
  border: 2px solid rgba(245, 240, 232, 0.3);
  border-top-color: var(--paper);
  border-radius: 50%;
  animation: spin 0.65s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Divider */
.form-divider {
  width: 1px;
  background: var(--ink);
  margin: 0 48px;
}

/* Result col */
.result-col { padding-left: 0; }

.result-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.16em;
  color: var(--green);
  margin-bottom: 20px;
}
.result-dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  background: var(--green);
  flex-shrink: 0;
}

.result-short-box {
  display: flex;
  align-items: center;
  gap: 10px;
  border: var(--rule-heavy);
  padding: 14px 16px;
  margin-bottom: 14px;
  background: var(--paper-warm);
}
.result-code {
  flex: 1;
  font-family: var(--font-mono);
  font-size: 15px;
  font-weight: 500;
  color: var(--ink);
  word-break: break-all;
}
.copy-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--ink);
  border: none;
  color: var(--paper);
  padding: 7px 14px;
  cursor: pointer;
  font-size: 12.5px;
  font-weight: 500;
  transition: var(--t);
  flex-shrink: 0;
}
.copy-btn:hover { background: var(--vermillion); }

.result-original {
  font-size: 12px;
  color: var(--ink-light);
  line-height: 1.5;
  margin-bottom: 20px;
  border-left: 3px solid var(--ink-faint);
  padding-left: 10px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.result-original-label {
  display: block;
  font-size: 10px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--ink-faint);
  margin-bottom: 3px;
}

.result-qr {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  border: var(--rule-medium);
  padding: 10px;
  margin-bottom: 16px;
  background: var(--white);
  line-height: 0;
}
.qr-hint {
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.1em;
  color: var(--ink-faint);
  line-height: 1;
}

.stats-link { margin-top: 4px; }

/* Empty state */
.empty-result {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 32px 0;
}
.empty-num {
  font-family: var(--font-display);
  font-size: 80px;
  font-weight: 700;
  font-style: italic;
  color: var(--ink-faint);
  line-height: 1;
  margin-bottom: 16px;
  opacity: 0.3;
}
.empty-text {
  font-size: 14px;
  color: var(--ink-faint);
  line-height: 1.7;
}

/* ═══ FEATURES ══════════════════════════ */
.features-section {
  border-top: var(--rule-heavy);
  padding: 60px 32px;
  background: var(--paper-warm);
}
.features-inner { max-width: 960px; margin: 0 auto; }

.features-header {
  display: flex;
  align-items: baseline;
  gap: 20px;
  margin-bottom: 40px;
  border-bottom: var(--rule-medium);
  padding-bottom: 20px;
}
.features-num {
  font-family: var(--font-display);
  font-size: 13px;
  font-weight: 400;
  font-style: italic;
  color: var(--ink-faint);
  letter-spacing: 0.04em;
}
.features-title {
  font-family: var(--font-display);
  font-size: 26px;
  font-weight: 700;
  color: var(--ink);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  border: var(--rule-medium);
}
.feature-item {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 28px;
  border-right: var(--rule-medium);
  transition: var(--t);
}
.feature-item:last-child { border-right: none; }
.feature-item:hover { background: var(--white); }

.feature-icon-wrap {
  width: 42px;
  height: 42px;
  border: var(--rule-medium);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ink);
  flex-shrink: 0;
  transition: var(--t);
}
.feature-item:hover .feature-icon-wrap {
  border-color: var(--vermillion);
  color: var(--vermillion);
}

.feature-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--ink);
}
.feature-desc {
  font-size: 13.5px;
  color: var(--ink-light);
  line-height: 1.65;
}

@media (max-width: 760px) {
  .form-inner { grid-template-columns: 1fr; }
  .form-divider { width: 100%; height: 1px; margin: 32px 0; }
  .form-col { padding-right: 0; }
  .features-grid { grid-template-columns: 1fr; }
  .feature-item { border-right: none; border-bottom: var(--rule-medium); }
  .feature-item:last-child { border-bottom: none; }
  .opts-grid { grid-template-columns: 1fr; }
  .hero { height: 440px; }
}
</style>
