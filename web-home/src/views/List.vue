<template>
  <div class="list-page">
    <!-- ─── Header ─── -->
    <div class="page-header">
      <div class="page-header-inner">
        <div class="page-title-wrap">
          <span class="page-title-num">02</span>
          <div>
            <h2 class="page-title">{{ t('list.pageTitle') }}</h2>
            <p class="page-desc">{{ t('list.pageDesc') }}</p>
          </div>
        </div>
        <div class="toolbar-right">
          <div class="search-wrap">
            <svg class="search-icon" viewBox="0 0 24 24" fill="none" width="14" height="14">
              <circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2"/>
              <path d="M21 21l-4.35-4.35" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
            <input v-model="keyword" class="search-input" :placeholder="t('list.searchHint')" @input="onSearch" />
          </div>
          <button class="lk-btn lk-btn-primary lk-btn-sm" @click="router.push('/')">
            <svg viewBox="0 0 24 24" fill="none" width="13" height="13">
              <path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
            </svg>
            {{ t('list.newLink') }}
          </button>
        </div>
      </div>

      <div v-if="total" class="count-strip">
        {{ t('list.totalRecords', { n: total }) }}
      </div>
    </div>

    <!-- ─── Table ─── -->
    <div class="table-outer">
      <n-data-table
        :columns="columns"
        :data="list"
        :loading="loading"
        :pagination="pagination"
        size="medium"
        :row-key="(row: ShortLink) => row.id"
        :scroll-x="900"
        @update:page="onPageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { NDataTable, useMessage } from 'naive-ui';
import type { DataTableColumns, PaginationProps } from 'naive-ui';
import { getList, type ShortLink } from '@/api/shortlink';

const { t }   = useI18n();
const router  = useRouter();
const message = useMessage();
const loading  = ref(false);
const list     = ref<ShortLink[]>([]);
const keyword  = ref('');
const total    = ref(0);

const pagination = ref<PaginationProps>({
  page: 1, pageSize: 15,
  showSizePicker: true,
  pageSizes: [10, 15, 20, 50],
  itemCount: 0,
  onChange: (p: number) => onPageChange(p),
  onUpdatePageSize: (s: number) => { pagination.value.pageSize = s; pagination.value.page = 1; loadData(); },
});

async function doCopy(text: string) {
  try { await navigator.clipboard.writeText(text); message.success(t('common.copied')); }
  catch { message.error(t('common.copyFail')); }
}

const columns: DataTableColumns<ShortLink> = [
  {
    title: () => t('list.colIndex'),
    key: 'idx',
    width: 52,
    render: (_, i) => h('span', {
      style: 'font-family:var(--font-mono);font-size:11px;color:var(--ink-faint)',
    }, String(((pagination.value.page ?? 1) - 1) * (pagination.value.pageSize ?? 15) + i + 1).padStart(2, '0')),
  },
  {
    title: () => t('list.colCode'),
    key: 'code',
    width: 160,
    render: (row) => {
      const shortUrl = row.short_url || `${window.location.origin}/r/${row.code}`;
      return h('div', { class: 'code-wrap' }, [
        h('a', {
          class: 'code-link',
          href: shortUrl,
          target: '_blank',
          rel: 'noopener noreferrer',
          title: shortUrl,
          onClick: (e: MouseEvent) => e.stopPropagation(),
        }, row.code),
        h('button', {
          class: 'code-copy',
          title: '复制短链接',
          onClick: () => doCopy(shortUrl),
        }, h('svg', { viewBox: '0 0 24 24', fill: 'none', width: '12', height: '12' }, [
          h('rect', { x: '9', y: '9', width: '13', height: '13', rx: '2', stroke: 'currentColor', 'stroke-width': '2' }),
          h('path', { d: 'M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1', stroke: 'currentColor', 'stroke-width': '2' }),
        ])),
      ]);
    },
  },
  {
    title: () => t('list.colUrl'),
    key: 'url',
    minWidth: 220,
    render: (row) =>
      h('div', { class: 'url-cell' }, [
        row.title ? h('div', { class: 'url-title' }, row.title) : null,
        h('a', {
          class: 'url-href',
          href: row.original_url,
          target: '_blank',
          rel: 'noopener noreferrer',
          title: row.original_url,
        }, row.original_url),
      ]),
  },
  {
    title: () => t('list.colTotal'),
    key: 'total_clicks',
    width: 88,
    align: 'right',
    render: (row) =>
      h('span', { class: 'num-cell' }, row.total_clicks.toLocaleString()),
  },
  {
    title: () => t('list.colToday'),
    key: 'today_clicks',
    width: 68,
    align: 'right',
    render: (row) =>
      h('span', { class: row.today_clicks > 0 ? 'num-today-on' : 'num-today' }, row.today_clicks),
  },
  {
    title: () => t('list.colStatus'),
    key: 'status',
    width: 100,
    render: (row) => {
      if (!row.expire_at) return h('span', { class: 'lk-badge badge-green' }, 'ACTIVE');
      const expired = new Date(row.expire_at) < new Date();
      return h('span', { class: `lk-badge ${expired ? 'badge-red' : 'badge-amber'}` }, expired ? 'EXPIRED' : 'LIMITED');
    },
  },
  {
    title: () => t('list.colTime'),
    key: 'created_at',
    width: 148,
    render: (row) => h('span', { class: 'time-cell' }, row.created_at),
  },
  {
    title: '',
    key: 'action',
    width: 72,
    fixed: 'right',
    render: (row) =>
      h('button', {
        class: 'action-cell',
        onClick: () => router.push(`/stats/${row.code}`),
      }, t('list.actionStats')),
  },
];

async function loadData() {
  loading.value = true;
  try {
    const res = await getList({ page: pagination.value.page, limit: pagination.value.pageSize, keyword: keyword.value || undefined });
    list.value  = res.data.list;
    total.value = res.data.total;
    pagination.value.itemCount = res.data.total;
  } catch (e: any) {
    message.error(e.message || t('common.loadFail'));
  } finally {
    loading.value = false;
  }
}

function onPageChange(p: number) { pagination.value.page = p; loadData(); }

let timer: ReturnType<typeof setTimeout>;
function onSearch() {
  clearTimeout(timer);
  timer = setTimeout(() => { pagination.value.page = 1; loadData(); }, 380);
}

onMounted(loadData);
</script>

<style scoped>
.list-page { max-width: 1100px; margin: 0 auto; }

/* ─── Header ─────────────────────────────── */
.page-header {
  border-bottom: var(--rule-heavy);
  margin-bottom: 0;
  padding: 40px 32px 0;
}
.page-header-inner {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding-bottom: 24px;
  gap: 16px;
  flex-wrap: wrap;
}
.page-title-wrap {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}
.page-title-num {
  font-family: var(--font-display);
  font-size: 44px;
  font-weight: 700;
  font-style: italic;
  color: var(--ink-faint);
  opacity: 0.3;
  line-height: 1;
}
.page-title {
  font-family: var(--font-display);
  font-size: 26px;
  font-weight: 700;
  color: var(--ink);
  margin-bottom: 4px;
}
.page-desc {
  font-size: 13px;
  color: var(--ink-light);
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.search-wrap {
  position: relative;
  display: flex;
  align-items: center;
}
.search-icon {
  position: absolute;
  left: 10px;
  color: var(--ink-faint);
  pointer-events: none;
}
.search-input {
  background: var(--white);
  border: var(--rule-medium);
  padding: 7px 12px 7px 32px;
  font-family: var(--font-body);
  font-size: 13px;
  color: var(--ink);
  outline: none;
  transition: var(--t);
  width: 200px;
}
.search-input::placeholder { color: var(--ink-faint); }
.search-input:focus { border-color: var(--vermillion); }

.count-strip {
  font-size: 12px;
  color: var(--ink-faint);
  padding: 8px 0;
  letter-spacing: 0.02em;
}
.count-strip strong { color: var(--ink); font-weight: 600; }

/* ─── Table ──────────────────────────────── */
.table-outer {
  background: var(--white);
  border: var(--rule-medium);
  border-top: none;
}

/* ─── Cell styles ────────────────────────── */
.code-wrap {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  border: 1px solid rgba(224, 85, 53, 0.2);
  padding: 3px 6px 3px 8px;
  transition: var(--t);
}
.code-wrap:hover { background: var(--vermillion-bg); border-color: var(--vermillion); }

.code-link {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 500;
  color: var(--vermillion);
  text-decoration: none;
  letter-spacing: 0.03em;
}
.code-link:hover { text-decoration: underline; }

.code-copy {
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  color: var(--ink-faint);
  display: flex;
  align-items: center;
  transition: var(--t);
  flex-shrink: 0;
}
.code-copy:hover { color: var(--vermillion); }

.url-cell { display: flex; flex-direction: column; gap: 2px; }
.url-title { font-size: 13.5px; font-weight: 500; color: var(--ink); }
.url-href {
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--ink-faint);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 360px;
  text-decoration: none;
  display: block;
  transition: var(--t);
}
.url-href:hover { color: var(--vermillion); text-decoration: underline; }

.num-cell { font-family: var(--font-mono); font-size: 13px; font-weight: 500; color: var(--ink); }
.num-today { font-family: var(--font-mono); font-size: 13px; color: var(--ink-faint); }
.num-today-on { font-family: var(--font-mono); font-size: 13px; font-weight: 600; color: var(--green); }

.time-cell { font-family: var(--font-mono); font-size: 11.5px; color: var(--ink-faint); }

.action-cell {
  background: none;
  border: none;
  font-size: 12.5px;
  font-weight: 500;
  color: var(--ink-light);
  cursor: pointer;
  padding: 4px 0;
  transition: var(--t);
  white-space: nowrap;
  letter-spacing: 0.02em;
}
.action-cell:hover { color: var(--vermillion); }
</style>
