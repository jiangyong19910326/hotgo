<template>
  <div class="stats-page">
    <!-- ─── Header ─── -->
    <div class="stats-header">
      <div class="stats-header-inner">
        <button class="back-btn" @click="router.back()">
          <svg viewBox="0 0 24 24" fill="none" width="14" height="14">
            <path d="M19 12H5M12 19l-7-7 7-7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ t('common.back') }}
        </button>
        <div class="page-label">{{ t('stats.pageTitle') }}</div>
      </div>
    </div>

    <n-spin :show="loading" size="large">
      <template v-if="stats">
        <div class="stats-inner">
          <!-- ─── Info bar ─── -->
          <div class="info-bar">
            <div class="info-bar-left">
              <div class="info-code">
                <span class="info-code-label">{{ t('stats.shortCode') }}</span>
                <span class="info-code-val">{{ stats.code }}</span>
              </div>
              <div class="info-details">
                <h2 class="info-title">{{ stats.title || stats.code }}</h2>
                <div class="info-url-row">
                  <span class="info-url">{{ stats.original_url }}</span>
                </div>
                <div class="info-meta">
                  <span class="lk-badge badge-ink" v-if="!isExpired">ACTIVE</span>
                  <span class="lk-badge badge-red" v-else>EXPIRED</span>
                  <span class="meta-sep">·</span>
                  <span class="meta-text">{{ stats.created_at }}</span>
                  <template v-if="stats.expire_at">
                    <span class="meta-sep">·</span>
                    <span class="meta-text">{{ t('stats.expires') }} {{ stats.expire_at }}</span>
                  </template>
                </div>
              </div>
            </div>
            <div class="info-bar-right">
              <button class="lk-btn lk-btn-primary lk-btn-sm" @click="copyUrl(shortUrl)">
                <svg viewBox="0 0 24 24" fill="none" width="13" height="13">
                  <rect x="9" y="9" width="13" height="13" rx="1" stroke="currentColor" stroke-width="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" stroke="currentColor" stroke-width="2"/>
                </svg>
                {{ t('stats.copyLink') }}
              </button>
              <div class="qr-mini">
                <qrcode-vue :value="shortUrl" :size="72" render-as="svg" foreground="#1c160d" background="#ffffff" />
              </div>
            </div>
          </div>

          <!-- ─── Stat cards ─── -->
          <div class="stat-cards">
            <div class="stat-card" v-for="s in statItems" :key="s.label">
              <div class="stat-card-num">{{ s.value.toLocaleString() }}</div>
              <div class="stat-card-label">{{ s.label }}</div>
              <div class="stat-card-en">{{ s.en }}</div>
              <div class="stat-card-accent" :style="`background: ${s.color}`"></div>
            </div>
          </div>

          <!-- ─── Trend ─── -->
          <div class="panel">
            <div class="panel-header">
              <span class="panel-title">{{ t('stats.trendTitle') }}</span>
              <span class="panel-sub">{{ t('stats.trendSub') }}</span>
            </div>
            <div v-if="stats.trend?.length" class="trend">
              <div class="trend-grid">
                <div v-for="item in stats.trend" :key="item.date" class="trend-col">
                  <div class="trend-num">{{ item.clicks }}</div>
                  <div class="trend-track">
                    <div
                      class="trend-fill"
                      :style="`height: ${trendMax ? Math.max(2, (item.clicks / trendMax) * 100) : 2}%`"
                    ></div>
                  </div>
                  <div class="trend-date">{{ item.date.slice(5) }}</div>
                </div>
              </div>
            </div>
            <div v-else class="empty-data">{{ t('common.noData') }}</div>
          </div>

          <!-- ─── Two col ─── -->
          <div class="two-col">
            <div class="panel">
              <div class="panel-header">
                <span class="panel-title">{{ t('stats.sourceTitle') }}</span>
              </div>
              <template v-if="stats.referrers?.length">
                <div class="dist-list">
                  <div v-for="(item, i) in stats.referrers" :key="item.source" class="dist-row">
                    <div class="dist-meta">
                      <span class="dist-rank">{{ String(i + 1).padStart(2, '0') }}</span>
                      <span class="dist-name">{{ item.source || t('stats.directVisit') }}</span>
                      <span class="dist-count">{{ item.count.toLocaleString() }}</span>
                      <span class="dist-pct">{{ refMax ? Math.round((item.count / refMax) * 100) : 0 }}%</span>
                    </div>
                    <div class="dist-bar-track">
                      <div
                        class="dist-bar vermillion-bar"
                        :style="`width: ${refMax ? Math.max(2, (item.count / refMax) * 100) : 0}%`"
                      ></div>
                    </div>
                  </div>
                </div>
              </template>
              <div v-else class="empty-data">{{ t('common.noData') }}</div>
            </div>

            <div class="panel">
              <div class="panel-header">
                <span class="panel-title">{{ t('stats.regionTitle') }}</span>
              </div>
              <template v-if="stats.regions?.length">
                <div class="dist-list">
                  <div v-for="(item, i) in stats.regions" :key="item.region" class="dist-row">
                    <div class="dist-meta">
                      <span class="dist-rank">{{ String(i + 1).padStart(2, '0') }}</span>
                      <span class="dist-name">{{ item.region || t('stats.unknown') }}</span>
                      <span class="dist-count">{{ item.count.toLocaleString() }}</span>
                      <span class="dist-pct">{{ regMax ? Math.round((item.count / regMax) * 100) : 0 }}%</span>
                    </div>
                    <div class="dist-bar-track">
                      <div
                        class="dist-bar green-bar"
                        :style="`width: ${regMax ? Math.max(2, (item.count / regMax) * 100) : 0}%`"
                      ></div>
                    </div>
                  </div>
                </div>
              </template>
              <div v-else class="empty-data">{{ t('common.noData') }}</div>
            </div>
          </div>
        </div>
      </template>

      <div v-else-if="!loading" class="not-found">
        <div class="nf-num">404</div>
        <p class="nf-text">{{ t('stats.notFoundText') }}</p>
        <button class="lk-btn lk-btn-primary" @click="router.push('/')">{{ t('common.goHome') }}</button>
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { NSpin, useMessage } from 'naive-ui';
import QrcodeVue from 'qrcode.vue';
import { getStats, type StatsResult } from '@/api/shortlink';

const { t }   = useI18n();
const route   = useRoute();
const router  = useRouter();
const message = useMessage();
const loading = ref(false);
const stats   = ref<StatsResult | null>(null);

// 始终用当前页面域名拼接，确保复制/二维码指向真实可访问地址
const shortUrl = computed(() =>
  stats.value ? `${window.location.origin}/r/${stats.value.code}` : ''
);

const isExpired = computed(() =>
  stats.value?.expire_at ? new Date(stats.value.expire_at) < new Date() : false
);

const statItems = computed(() => [
  { label: t('stats.statTotal'),     value: stats.value?.total_clicks ?? 0,     en: 'Total',      color: '#e05535' },
  { label: t('stats.statToday'),     value: stats.value?.today_clicks ?? 0,     en: 'Today',      color: '#2d8a5e' },
  { label: t('stats.statYesterday'), value: stats.value?.yesterday_clicks ?? 0, en: 'Yesterday',  color: '#b87a16' },
  { label: t('stats.statMonthly'),   value: stats.value?.monthly_clicks ?? 0,   en: 'This Month', color: '#1c160d' },
]);

const trendMax = computed(() =>
  stats.value?.trend?.length ? Math.max(...stats.value.trend.map((t) => t.clicks), 1) : 1
);
const refMax = computed(() =>
  stats.value?.referrers?.length ? Math.max(...stats.value.referrers.map((r) => r.count), 1) : 1
);
const regMax = computed(() =>
  stats.value?.regions?.length ? Math.max(...stats.value.regions.map((r) => r.count), 1) : 1
);

async function loadStats() {
  loading.value = true;
  try {
    const res = await getStats(route.params.code as string);
    stats.value = res.data;
  } catch (e: any) {
    message.error(e.message || t('common.loadFail'));
  } finally {
    loading.value = false;
  }
}

async function copyUrl(url: string) {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(url);
    } else {
      // HTTP 环境降级方案
      const el = document.createElement('textarea');
      el.value = url;
      el.style.cssText = 'position:fixed;opacity:0';
      document.body.appendChild(el);
      el.select();
      document.execCommand('copy');
      document.body.removeChild(el);
    }
    message.success(t('common.copied'));
  } catch {
    message.error(t('common.copyFail'));
  }
}

onMounted(loadStats);
</script>

<style scoped>
.stats-page { background: var(--paper); min-height: 100%; }

/* ─── Header ─────────────────────────────── */
.stats-header {
  border-bottom: var(--rule-heavy);
  padding: 0 32px;
}
.stats-header-inner {
  max-width: 1100px;
  margin: 0 auto;
  height: 56px;
  display: flex;
  align-items: center;
  gap: 20px;
}
.back-btn {
  display: flex;
  align-items: center;
  gap: 7px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 12.5px;
  font-weight: 500;
  letter-spacing: 0.06em;
  color: var(--ink-light);
  transition: var(--t);
  padding: 0;
}
.back-btn:hover { color: var(--ink); }
.page-label {
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--ink-faint);
  border-left: 2px solid var(--ink-faint);
  padding-left: 14px;
}

/* ─── Inner ──────────────────────────────── */
.stats-inner { max-width: 1100px; margin: 0 auto; padding: 40px 32px; }

/* ─── Info bar ───────────────────────────── */
.info-bar {
  display: flex;
  align-items: flex-start;
  gap: 24px;
  border: var(--rule-heavy);
  padding: 24px;
  background: var(--white);
  margin-bottom: 24px;
}
.info-bar-left { flex: 1; }
.info-bar-right { display: flex; flex-direction: column; align-items: flex-end; gap: 12px; flex-shrink: 0; }

.info-code {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}
.info-code-label {
  font-family: var(--font-mono);
  font-size: 9px;
  letter-spacing: 0.16em;
  color: var(--ink-faint);
}
.info-code-val {
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 500;
  color: var(--vermillion);
  border: 1px solid rgba(224, 85, 53, 0.2);
  padding: 2px 8px;
  background: var(--vermillion-bg);
}

.info-title {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 6px;
}
.info-url-row { margin-bottom: 10px; }
.info-url {
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--ink-faint);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
  max-width: 560px;
}
.info-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.meta-sep { color: var(--ink-faint); }
.meta-text { font-family: var(--font-mono); font-size: 11px; color: var(--ink-faint); }
.qr-mini { border: var(--rule-medium); padding: 6px; background: var(--white); line-height: 0; }

/* ─── Stat cards ─────────────────────────── */
.stat-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  border: var(--rule-heavy);
  background: var(--white);
  margin-bottom: 24px;
}
.stat-card {
  padding: 24px 22px 20px;
  border-right: var(--rule-medium);
  position: relative;
  overflow: hidden;
  transition: var(--t);
}
.stat-card:last-child { border-right: none; }
.stat-card:hover { background: var(--paper-warm); }

.stat-card-accent {
  position: absolute;
  top: 0; left: 0;
  width: 4px;
  height: 100%;
  opacity: 0.7;
}
.stat-card-num {
  font-family: var(--font-display);
  font-size: 38px;
  font-weight: 700;
  line-height: 1;
  color: var(--ink);
  margin-bottom: 10px;
  padding-left: 14px;
}
.stat-card-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--ink);
  padding-left: 14px;
  margin-bottom: 2px;
}
.stat-card-en {
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.1em;
  color: var(--ink-faint);
  padding-left: 14px;
}

/* ─── Panel ──────────────────────────────── */
.panel {
  background: var(--white);
  border: var(--rule-medium);
  padding: 24px;
  margin-bottom: 20px;
}
.panel-header {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 20px;
  padding-bottom: 14px;
  border-bottom: var(--rule-medium);
}
.panel-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--ink);
}
.panel-sub {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--ink-faint);
  letter-spacing: 0.06em;
}
.empty-data {
  font-size: 13px;
  color: var(--ink-faint);
  text-align: center;
  padding: 28px 0;
  letter-spacing: 0.04em;
}

/* ─── Trend chart ────────────────────────── */
.trend-grid {
  display: flex;
  align-items: flex-end;
  gap: 5px;
  height: 150px;
  padding: 0 2px;
}
.trend-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  height: 100%;
  cursor: default;
}
.trend-col:hover .trend-fill { opacity: 0.75; }
.trend-col:hover .trend-num  { color: var(--ink); }
.trend-num {
  font-family: var(--font-mono);
  font-size: 10px;
  color: transparent;
  min-height: 13px;
  line-height: 1;
  transition: color 0.15s;
}
.trend-col:hover .trend-num { color: var(--ink-light); }
.trend-track {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
  overflow: hidden;
}
.trend-fill {
  width: 100%;
  background: linear-gradient(to top, var(--vermillion), rgba(224,85,53,0.45));
  border-radius: 2px 2px 0 0;
  transition: height 0.55s var(--ease), opacity 0.15s;
  min-height: 3px;
}
.trend-date {
  font-family: var(--font-mono);
  font-size: 9px;
  color: var(--ink-faint);
  letter-spacing: 0.02em;
  white-space: nowrap;
}

/* ─── Two col ────────────────────────────── */
.two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }

/* ─── Distribution ───────────────────────── */
.dist-list { display: flex; flex-direction: column; gap: 12px; }
.dist-row { display: flex; flex-direction: column; gap: 5px; }
.dist-meta {
  display: flex;
  align-items: baseline;
  gap: 7px;
}
.dist-rank {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--ink-faint);
  width: 18px;
  flex-shrink: 0;
}
.dist-name {
  font-size: 13px;
  color: var(--ink);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.dist-count {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
  color: var(--ink);
  flex-shrink: 0;
}
.dist-pct {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--ink-faint);
  width: 36px;
  text-align: right;
  flex-shrink: 0;
}
.dist-bar-track {
  height: 6px;
  width: 100%;
  background: var(--paper-warm);
  border-radius: 3px;
  overflow: hidden;
}
.dist-bar {
  height: 100%;
  border-radius: 3px;
  transition: width 0.6s var(--ease);
}
.vermillion-bar { background: linear-gradient(to right, var(--vermillion), rgba(224,85,53,0.5)); }
.green-bar      { background: linear-gradient(to right, var(--green),      rgba(45,138,94,0.5)); }

/* ─── Not found ──────────────────────────── */
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80px 32px;
  gap: 16px;
  text-align: center;
}
.nf-num {
  font-family: var(--font-display);
  font-size: 100px;
  font-weight: 700;
  font-style: italic;
  color: var(--ink-faint);
  opacity: 0.2;
  line-height: 1;
}
.nf-text { font-size: 16px; color: var(--ink-light); }

@media (max-width: 700px) {
  .stat-cards  { grid-template-columns: 1fr 1fr; }
  .two-col     { grid-template-columns: 1fr; }
  .info-bar    { flex-direction: column; }
  .stats-inner { padding: 24px 16px; }
}
</style>
