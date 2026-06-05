<template>
  <div ref="screenRef" class="hmi-screen">
    <div class="warehouse-orbit orbit-a"></div>
    <div class="warehouse-orbit orbit-b"></div>

    <!-- 顶栏 -->
    <div class="hmi-top">
      <div class="top-deco-left"></div>
      <div class="top-title">
        <span class="title-kicker">DATA WAREHOUSE</span>
        <strong>设备实时监控系统</strong>
      </div>
      <div class="top-deco-right">
        <n-select
          v-model:value="mineId"
          :options="mineOptions"
          placeholder="选择矿场"
          size="small"
          style="width: 160px"
          @update:value="onMineChange"
        />
        <n-select
          v-model:value="deviceId"
          :options="deviceOptions"
          placeholder="选择设备"
          size="small"
          style="width: 200px"
          @update:value="onDeviceChange"
        />
        <span class="clock">{{ clock }}</span>
        <n-button class="fullscreen-btn" size="small" ghost @click="toggleFullscreen">
          {{ isFullscreen ? '退出全屏' : '全屏' }}
        </n-button>
      </div>
    </div>

    <div class="warehouse-strip">
      <div
        v-for="item in warehouseStats"
        :key="item.label"
        class="warehouse-chip"
        :class="item.tone"
      >
        <span class="chip-label">{{ item.label }}</span>
        <strong :key="`${item.label}-${item.value}`" class="chip-value">{{ item.value }}</strong>
        <small :key="`${item.label}-${item.desc}`" class="chip-desc">{{ item.desc }}</small>
      </div>
    </div>

    <!-- 主体三列 -->
    <div class="hmi-body">
      <!-- 左列：仪表盘 + 参数表 + 设备图 -->
      <div class="col col-left">
        <div class="block">
          <div class="block-title">仪表监测</div>
          <div class="gauges">
            <div ref="oilLevelGaugeRef" class="gauge"></div>
            <div ref="currentGaugeRef" class="gauge"></div>
            <div ref="lubeStatusGaugeRef" class="gauge gauge-wide"></div>
          </div>
        </div>

        <div class="block block-params">
          <div class="block-title">主要参数</div>
          <div class="param-table">
            <div class="param-row" v-for="row in paramRows" :key="row.label">
              <span class="param-label">{{ row.label }}</span>
              <span class="param-value" :class="row.cls">
                {{ row.value }}<small>{{ row.unit }}</small>
              </span>
            </div>
            <div v-if="paramRows.length === 0" class="param-empty">暂无参数</div>
          </div>
        </div>
      </div>

      <!-- 中列：设备示意 + 状态 -->
      <div class="col col-mid">
        <div class="block block-device">
          <div class="block-title">{{ device?.name || '设备' }} 工艺示意</div>

          <!-- 报警条 -->
          <div v-if="firstAlarm" class="alarm-banner">
            <span class="alarm-banner-icon">⚠</span>
            {{ firstAlarm.name || firstAlarm.field }} 报警，请检查
          </div>
          <div class="device-canvas">
            <div class="machine-stage">
              <img class="machine-bg" :src="deviceImageUrl" alt="设备实时监控图" />
              <div class="readout-panel">
                <div
                  v-for="item in monitorReadouts"
                  :key="item.label"
                  class="readout-card"
                  :class="item.tone"
                >
                  <span class="readout-label">{{ item.label }}</span>
                  <strong>{{ item.value }}</strong>
                  <small>{{ item.unit }}</small>
                </div>
              </div>
              <div class="live-badge">
                <span></span>
                数据实时更新
              </div>
            </div>
          </div>
        </div>

        <!-- 底部图表 -->
        <div class="bottom-charts">
          <div class="block">
            <div class="block-title">温度趋势</div>
            <div ref="tempChartRef" class="chart-area"></div>
          </div>
          <div class="block">
            <div class="block-title">电流趋势</div>
            <div ref="currentChartRef" class="chart-area"></div>
          </div>
        </div>
      </div>

      <!-- 右列：状态与报警 -->
      <div class="col col-right">
        <div class="block">
          <div class="block-title">设备状态</div>
          <div class="status-list">
            <div
              v-for="(s, idx) in statusList"
              :key="s.field"
              class="status-row"
              :class="{ 'is-active': statusCarouselIndex % statusList.length === idx }"
            >
              <span class="status-led" :class="isStatusOn(s) ? 'led-on' : 'led-off'"></span>
              <span class="status-name">{{ s.name || s.field }}</span>
              <span class="status-text" :class="isStatusOn(s) ? 'txt-on' : 'txt-off'">
                {{ statusText(s) }}
              </span>
            </div>
            <div v-if="statusList.length === 0" class="status-empty">暂无状态点位</div>
          </div>
        </div>

        <div class="block block-alarm">
          <div class="block-title">
            报警监测
            <span class="block-sub">活动 {{ activeAlarmCount }} / 共 {{ allAlarms.length }}</span>
          </div>
          <div class="alarm-list">
            <div
              v-for="(a, idx) in activeAlarms"
              :key="a.pointId"
              class="alarm-row"
              :class="{
                'alarm-on': a.active,
                'alarm-off': !a.active,
                'is-active': alarmCarouselIndex % activeAlarms.length === idx,
              }"
            >
              <span class="alarm-led"></span>
              <span class="alarm-name" :title="a.field">{{ a.name || a.field }}</span>
              <span class="alarm-tag">{{ a.active ? '报警' : '正常' }}</span>
            </div>
            <div v-if="activeAlarms.length === 0" class="alarm-empty">无</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue';
  import * as echarts from 'echarts';
  import { MineOptions, DeviceList, Overview } from '@/api/plc';
  import { http } from '@/utils/http/axios';
  import { SocketEnum } from '@/enums/socketEnum';
  import { addOnMessage, removeOnMessage, WebSocketMessage } from '@/utils/websocket';
  import coneCrusherImg from '@/assets/images/cone-crusher.jpg';
  import sandMakerImg from '@/assets/images/sand-maker.png';

  const mineId = ref<number | null>(null);
  const deviceId = ref<number | null>(null);
  const mineOptions = ref<{ label: string; value: number }[]>([]);
  const deviceOptions = ref<{ label: string; value: number }[]>([]);
  const allDevices = ref<any[]>([]);

  const device = ref<any>(null);
  const points = ref<any[]>([]);
  const allAlarms = ref<any[]>([]);
  const loading = ref(false);
  const clock = ref('');
  const screenRef = ref<HTMLElement>();
  const isFullscreen = ref(false);
  let clockTimer: any = null;
  let refreshTimer: any = null;
  let chartTimer: any = null;
  let carouselTimer: any = null;
  const autoRefresh = ref(true);
  const overviewRefreshMs = 60 * 1000;
  const chartRefreshMs = 60 * 1000;
  const realtimeStaleMs = 2 * 60 * 1000;

  // 计算
  const numericPoints = computed(() => points.value.filter((p) => p.dataType !== 'Bool'));
  const boolPoints = computed(() => points.value.filter((p) => p.dataType === 'Bool'));
  const statusList = computed(() => boolPoints.value);
  const activeAlarms = computed(() => allAlarms.value.filter((a) => a.active));
  const alarms = computed(() => allAlarms.value.filter((a) => a.active));
  const activeAlarmCount = computed(() => alarms.value.length);
  const firstAlarm = computed(() => alarms.value[0]);
  const alarmCarouselIndex = ref(0);
  const statusCarouselIndex = ref(0);

  const alarmCarouselItem = computed(() => {
    const list = allAlarms.value;
    if (list.length === 0) return null;
    return list[alarmCarouselIndex.value % list.length];
  });

  const statusCarouselItem = computed(() => {
    const list = statusList.value;
    if (list.length === 0) return null;
    return list[statusCarouselIndex.value % list.length];
  });

  const readoutValue = (keys: string[], fallback = '—') => {
    const keySet = keys.map((key) => key.toLowerCase());
    for (const p of numericPoints.value) {
      const field = String(p.field || '').toLowerCase();
      const name = String(p.name || '').toLowerCase();
      if (keySet.some((key) => field.includes(key) || name.includes(key))) {
        return {
          value: p.engValue == null ? fallback : formatValue(p.engValue),
          unit: p.unit || '',
          point: p,
        };
      }
    }
    return { value: fallback, unit: '', point: null };
  };

  const runState = computed(() => getDeviceRunState());
  const isDeviceOffline = computed(() => runState.value === '离线');

  const warehouseStats = computed(() => [
    {
      label: '采集点位',
      value: points.value.length || '--',
      desc: '采集点位',
      tone: 'tone-cyan',
    },
    {
      label: '运行状态',
      value: runState.value,
      desc: device.value?.name || '未选择设备',
      tone: runState.value === '运行' ? 'tone-green' : 'tone-amber',
    },
    {
      label: '报警通道',
      value: alarmCarouselItem.value
        ? alarmCarouselItem.value.active
          ? '报警'
          : '正常'
        : `${activeAlarmCount.value}`,
      desc: alarmCarouselItem.value
        ? alarmCarouselItem.value.name || alarmCarouselItem.value.field
        : `共 ${allAlarms.value.length} 路`,
      tone: alarmCarouselItem.value?.active ? 'tone-red' : 'tone-cyan',
    },
    {
      label: '设备状态',
      value: statusCarouselItem.value ? statusText(statusCarouselItem.value) : '暂无',
      desc: statusCarouselItem.value
        ? statusCarouselItem.value.name || statusCarouselItem.value.field
        : `兜底 ${overviewRefreshMs / 1000}s`,
      tone:
        statusCarouselItem.value && isStatusOn(statusCarouselItem.value)
          ? 'tone-green'
          : 'tone-amber',
    },
  ]);

  function pointMatches(p: any, keys: string[]) {
    const field = String(p?.field || '').toLowerCase();
    const name = String(p?.name || '').toLowerCase();
    return keys.some(
      (key) => field.includes(key.toLowerCase()) || name.includes(key.toLowerCase())
    );
  }

  function isPointActive(p: any) {
    if (typeof p?.active === 'boolean') return p.active;
    return Number(p?.engValue || 0) !== 0;
  }

  function isStatusOn(p: any) {
    return !isDeviceOffline.value && isPointActive(p);
  }

  function statusText(p: any) {
    if (isDeviceOffline.value) return '离线';
    return p?.stateText || (isPointActive(p) ? '运行' : '停止');
  }

  function hasPointValue(p: any) {
    return p?.engValue !== null && p?.engValue !== undefined;
  }

  function isFreshPoint(p: any) {
    if (!hasPointValue(p)) return false;
    if (!p.collectedAt) return false;
    const time = new Date(p.collectedAt).getTime();
    return Number.isFinite(time) && Date.now() - time <= realtimeStaleMs;
  }

  function numericValue(keys: string[]) {
    const p = numericPoints.value.find((item) => isFreshPoint(item) && pointMatches(item, keys));
    const value = Number(p?.engValue ?? 0);
    return Number.isFinite(value) ? value : 0;
  }

  function findBoolPoint(keys: string[]) {
    const candidates = boolPoints.value.filter(
      (item) =>
        isFreshPoint(item) &&
        !String(item.field || '')
          .toUpperCase()
          .startsWith('AL_')
    );
    return candidates.find((item) => pointMatches(item, keys));
  }

  function getDeviceRunState() {
    if (!points.value.some(hasPointValue)) return '无数据';
    if (!points.value.some(isFreshPoint)) return '离线';

    const explicitStop = findBoolPoint([
      'stop',
      'stopped',
      'shutdown',
      'halt',
      'fault',
      '停止',
      '停机',
      '故障',
    ]);
    if (explicitStop && isPointActive(explicitStop)) return '停止';

    const explicitRun = findBoolPoint([
      'run',
      'running',
      'motor_run',
      'crusher_run',
      'main_motor_run',
      'device_running',
      'equipment_running',
      '运行',
      '运转',
      '启动',
      '主机运行',
      '电机运行',
      '破碎机运行',
      '设备状态',
    ]);
    if (explicitRun) return isPointActive(explicitRun) ? '运行' : '停止';

    // 没有明确运行点时，用新鲜的能耗类模拟量兜底。
    const current = numericValue(['main_current', 'crusher_ampere', 'ampere', 'current', '电流']);
    const power = numericValue(['power', 'kw', '功率']);
    if (current > 0.5 || power > 0.1) return '运行';

    return '停止';
  }

  const monitorReadouts = computed(() => {
    const oilTemp = readoutValue(['oil_temp', 'lube_tank_temp', 'tank_temp', '油温', '油箱温度']);
    const gap = readoutValue(['gap', 'css', '间隙', '排矿口']);
    const current = readoutValue(['main_current', 'crusher_ampere', 'ampere', '电流']);
    const runtime = readoutValue(['run_hours', 'crusher_time_h', '运行时间']);

    return [
      {
        label: '油温',
        icon: '🌡',
        value: oilTemp.value,
        unit: oilTemp.unit || '℃',
        tone: cellClass(oilTemp.point),
      },
      {
        label: '间隙',
        icon: '⚙',
        value: gap.value,
        unit: gap.unit || 'mm',
        tone: cellClass(gap.point),
      },
      {
        label: '电流',
        icon: '⚡',
        value: current.point
          ? Math.round(Number(current.point.engValue)).toString()
          : current.value,
        unit: current.unit || 'A',
        tone: 'val-power',
      },
      {
        label: '运行时间',
        icon: '⏱',
        value: runtime.value,
        unit: runtime.unit || 'h',
        tone: '',
      },
      {
        label: '设备状态',
        icon: '●',
        value: runState.value,
        unit: '',
        tone: runState.value === '运行' ? 'val-power' : 'val-low',
      },
    ];
  });

  // 参数表
  const paramRows = computed(() => {
    const rows: any[] = [];
    const allPointList = [...numericPoints.value, ...boolPoints.value];
    const findPoint = (keys: string[]) => {
      for (const p of allPointList) {
        if (pointMatches(p, keys)) return p;
      }
      return null;
    };
    const candidates = [
      {
        label: '释放缸压力',
        keys: ['release_cyl_pressure', 'release_pressure', '释放缸压力'],
        unit: 'bar',
      },
      {
        label: '锁紧缸压力',
        keys: ['lock_pressure', 'locking_pressure', '锁紧缸压力'],
        unit: 'bar',
      },
      {
        label: '回油温度',
        keys: ['lube_return_temp', 'return_oil_temp', 'return_temp', '润滑回油温度', '回油温度'],
        unit: '℃',
        noDecimal: true,
      },
      {
        label: '油箱温度',
        keys: ['lube_tank_temp', 'tank_temp', 'oil_tank_temp', '润滑油箱温度', '油箱温度'],
        unit: '℃',
        noDecimal: true,
      },
      {
        label: '破碎机运行时间',
        keys: ['run_hours', 'crusher_time_h', 'crusher_runtime', '运行时间', '破碎机运行小时'],
        unit: 'h',
      },
      {
        label: '排料口尺寸',
        keys: [
          'discharge_opening',
          'discharge_size',
          'outlet_size',
          'css',
          'gap',
          '排料口',
          '排矿口',
          '间隙',
        ],
        unit: 'mm',
      },
      {
        label: '进料机运转状态',
        keys: [
          'feeder_run',
          'feeder_running',
          'feed_run',
          'feed_running',
          '进料机运行',
          '进料机运转',
        ],
        unit: '',
        boolText: true,
      },
    ];
    for (const c of candidates) {
      const p = findPoint(c.keys);
      if (p && p.engValue != null) {
        const isBool = c.boolText || p.dataType === 'Bool';
        rows.push({
          label: c.label,
          value: isBool
            ? statusText(p)
            : c.noDecimal
              ? Math.round(Number(p.engValue)).toString()
              : formatValue(p.engValue),
          unit: isBool ? '' : p.unit || c.unit || '',
          cls: isBool ? (isStatusOn(p) ? 'val-power' : 'val-low') : cellClass(p),
        });
      }
    }
    return rows;
  });

  function valueOf(...keys: string[]): string {
    const map: Record<string, any> = {};
    for (const p of numericPoints.value) {
      map[(p.field || '').toLowerCase()] = p;
      map[p.name || ''] = p;
    }
    for (const k of keys) {
      const p = map[k.toLowerCase()] || map[k];
      if (p && p.engValue != null) return `${formatValue(p.engValue)}${p.unit || ''}`;
    }
    return '—';
  }

  // 仪表盘
  const oilLevelGaugeRef = ref<HTMLDivElement>();
  const currentGaugeRef = ref<HTMLDivElement>();
  const lubeStatusGaugeRef = ref<HTMLDivElement>();
  let oilLevelGauge: echarts.ECharts | null = null;
  let currentGauge: echarts.ECharts | null = null;
  let lubeStatusGauge: echarts.ECharts | null = null;

  function ensureGauges() {
    if (oilLevelGaugeRef.value && !oilLevelGauge)
      oilLevelGauge = echarts.init(oilLevelGaugeRef.value);
    if (currentGaugeRef.value && !currentGauge) currentGauge = echarts.init(currentGaugeRef.value);
    if (lubeStatusGaugeRef.value && !lubeStatusGauge)
      lubeStatusGauge = echarts.init(lubeStatusGaugeRef.value);
  }

  function gaugeOption(
    name: string,
    value: number,
    max: number,
    unit: string,
    color: string,
    integer = false
  ) {
    const safeValue = Math.max(0, Math.min(max, Number(value || 0)));
    const displayVal = integer ? Math.round(safeValue) : Number(safeValue.toFixed(1));
    return {
      backgroundColor: 'transparent',
      tooltip: {
        formatter: `${name}<br/>${displayVal}${unit}`,
      },
      graphic: [
        {
          type: 'text',
          left: 'center',
          top: '38%',
          style: {
            text: '▶',
            font: 'bold 22px sans-serif',
            fill: '#f5a623',
            textAlign: 'center',
          },
          silent: true,
        },
      ],
      series: [
        {
          type: 'gauge',
          radius: '88%',
          min: 0,
          max,
          startAngle: 220,
          endAngle: -40,
          progress: {
            show: true,
            roundCap: true,
            width: 14,
            itemStyle: { color },
          },
          axisLine: {
            roundCap: true,
            lineStyle: {
              width: 14,
              color: [[1, '#e8edf2']],
            },
          },
          axisTick: { show: false },
          splitLine: { show: false },
          axisLabel: { show: false },
          pointer: { show: false },
          anchor: { show: false },
          title: {
            show: true,
            color: '#5a7a8a',
            fontSize: 12,
            fontWeight: 600,
            offsetCenter: [0, '78%'],
          },
          detail: {
            valueAnimation: true,
            formatter: `{value}${unit}`,
            color: '#1a2e3a',
            fontSize: 22,
            fontWeight: 700,
            offsetCenter: [0, '44%'],
          },
          data: [{ value: displayVal, name }],
        },
      ],
    };
  }

  function refreshGauges() {
    ensureGauges();
    const findPoint = (keys: string[]): any => {
      for (const p of [...numericPoints.value, ...boolPoints.value]) {
        const fn = (p.field || '').toLowerCase();
        if (keys.some((k) => fn.includes(k.toLowerCase()) || (p.name || '').includes(k))) {
          return p;
        }
      }
      return null;
    };
    const oilLevelPoint = findPoint(['hydraulic_oil_level', 'oil_level', '液压油油位', '油位']);
    const currentPoint = findPoint(['main_current', 'crusher_ampere', 'ampere', '电流']);
    const lubeStatusPoint = findPoint([
      'lube_status',
      'lubrication_status',
      'lubricating_oil_status',
      '润滑油状态',
      '润滑状态',
    ]);
    const oilLevel = normalizePercent(oilLevelPoint?.engValue);
    const current = Number(currentPoint?.engValue ?? 0);
    const lubeNormal = lubeStatusPoint ? Number(lubeStatusPoint.engValue || 0) !== 0 : false;
    const lubeText = isDeviceOffline.value
      ? '离线'
      : lubeStatusPoint
        ? lubeNormal
          ? '正常'
          : '异常'
        : '未知';
    const lubeColor = isDeviceOffline.value
      ? '#b0bec5'
      : lubeStatusPoint
        ? lubeNormal
          ? '#1890ff'
          : '#ff4d4f'
        : '#b0bec5';

    oilLevelGauge?.setOption(gaugeOption('液压油油位', oilLevel, 100, '%', '#f5a623'));
    currentGauge?.setOption(gaugeOption('主机电流', current, 200, 'A', '#1890ff', true));
    lubeStatusGauge?.setOption(statusGaugeOption('润滑油状态', lubeText, lubeColor));
  }

  function normalizePercent(value: any) {
    const n = Number(value ?? 0);
    if (!Number.isFinite(n)) return 0;
    return Math.max(0, Math.min(100, n <= 1 ? n * 100 : n));
  }

  function statusGaugeOption(name: string, text: string, color: string) {
    const value = text === '正常' ? 1 : text === '异常' ? 0.25 : 0.55;
    return {
      backgroundColor: 'transparent',
      tooltip: { formatter: `${name}<br/>${text}` },
      graphic: [
        {
          type: 'text',
          left: 'center',
          top: '34%',
          style: {
            text: '▶',
            font: 'bold 22px sans-serif',
            fill: '#f5a623',
            textAlign: 'center',
          },
          silent: true,
        },
      ],
      series: [
        {
          type: 'gauge',
          radius: '88%',
          min: 0,
          max: 1,
          startAngle: 220,
          endAngle: -40,
          progress: {
            show: true,
            roundCap: true,
            width: 14,
            itemStyle: { color },
          },
          axisLine: {
            roundCap: true,
            lineStyle: {
              width: 14,
              color: [[1, '#e8edf2']],
            },
          },
          axisTick: { show: false },
          splitLine: { show: false },
          axisLabel: { show: false },
          pointer: { show: false },
          anchor: { show: false },
          title: {
            show: true,
            color: '#5a7a8a',
            fontSize: 12,
            fontWeight: 600,
            offsetCenter: [0, '72%'],
          },
          detail: {
            formatter: () => text,
            color: '#1a2e3a',
            fontSize: 22,
            fontWeight: 700,
            offsetCenter: [0, '36%'],
          },
          data: [{ value, name }],
        },
      ],
    };
  }

  // 趋势图
  const tempChartRef = ref<HTMLDivElement>();
  const currentChartRef = ref<HTMLDivElement>();
  let tempChart: echarts.ECharts | null = null;
  let currentChart: echarts.ECharts | null = null;

  const deviceImageUrl = computed(() => {
    const text =
      `${device.value?.name || ''} ${device.value?.host || ''} ${device.value?.remark || ''}`.toLowerCase();
    return text.includes('制砂') || text.includes('sand') || text.includes('1263')
      ? sandMakerImg
      : coneCrusherImg;
  });

  function ensureCharts() {
    if (tempChartRef.value && !tempChart) tempChart = echarts.init(tempChartRef.value);
    if (currentChartRef.value && !currentChart) currentChart = echarts.init(currentChartRef.value);
  }

  function lineOption(legend: string[], xAxis: string[], series: any[], palette: string[]) {
    return {
      backgroundColor: 'transparent',
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(255, 255, 255, 0.96)',
        borderColor: 'rgba(24, 144, 255, 0.2)',
        textStyle: { color: '#1a2e3a' },
        extraCssText: 'box-shadow: 0 4px 16px rgba(24,144,255,0.12);',
        formatter: (params: any[]) =>
          params
            .map((p: any) => `${p.marker}${p.seriesName}：<b>${Math.round(p.value)}</b>`)
            .join('<br/>'),
      },
      legend: {
        data: legend,
        textStyle: { color: '#5a7a8a', fontSize: 11 },
        top: 4,
      },
      grid: { left: 50, right: 16, top: 32, bottom: 28 },
      xAxis: {
        type: 'category',
        data: xAxis,
        axisLine: { lineStyle: { color: 'rgba(24, 144, 255, 0.15)' } },
        axisLabel: { color: '#9ab0bc', fontSize: 9 },
        splitLine: { show: false },
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        splitLine: { lineStyle: { color: 'rgba(24, 144, 255, 0.08)', type: 'dashed' } },
        axisLabel: { color: '#9ab0bc', fontSize: 9 },
      },
      series: (series || []).map((s: any, i: number) => ({
        name: s.name,
        type: 'line',
        smooth: true,
        symbol: 'none',
        data: s.data,
        lineStyle: {
          color: palette[i % palette.length],
          width: 2.4,
          shadowBlur: 8,
          shadowColor: palette[i % palette.length],
        },
        areaStyle: { color: `${palette[i % palette.length]}22` },
      })),
    };
  }

  async function loadCharts(devId: number) {
    try {
      const tempRes: any = await http.request({
        url: '/plc/chart/temperature',
        method: 'GET',
        params: { deviceId: devId },
      });
      tempChart?.setOption(
        lineOption(tempRes?.legend || [], tempRes?.xAxis || [], tempRes?.series || [], [
          '#1890ff',
          '#f5a623',
          '#ff4d4f',
        ]),
        true
      );
    } catch {}
    try {
      const curRes: any = await http.request({
        url: '/plc/chart/current',
        method: 'GET',
        params: { deviceId: devId },
      });
      currentChart?.setOption(
        lineOption(curRes?.legend || [], curRes?.xAxis || [], curRes?.series || [], [
          '#1890ff',
          '#52c41a',
        ]),
        true
      );
    } catch {}
  }

  // 数据加载
  onMounted(async () => {
    await loadMines();
    startClock();
    startCarousel();
    document.addEventListener('fullscreenchange', onFullscreenChange);
    addOnMessage(SocketEnum.EventPlcRealtime, onPlcRealtimeMessage);
    if (autoRefresh.value) {
      refreshTimer = setInterval(() => loadOverview(false), overviewRefreshMs);
      chartTimer = setInterval(() => {
        if (deviceId.value) loadCharts(deviceId.value);
      }, chartRefreshMs);
    }
  });

  async function loadMines() {
    const res: any = await MineOptions();
    mineOptions.value = (res?.list || []).map((m: any) => ({ label: m.name, value: m.id }));
    if (mineOptions.value.length > 0 && !mineId.value) {
      mineId.value = mineOptions.value[0].value;
      await onMineChange(mineId.value);
    }
  }

  async function onMineChange(id: number | null) {
    mineId.value = id;
    deviceId.value = null;
    deviceOptions.value = [];
    allDevices.value = [];
    points.value = [];
    allAlarms.value = [];
    device.value = null;
    if (!id) return;
    const res: any = await DeviceList({ mineId: id, status: 1, page: 1, pageSize: 200 });
    allDevices.value = res?.list || [];
    deviceOptions.value = allDevices.value.map((d: any) => ({
      label: `${d.name} (${d.host})`,
      value: d.id,
    }));
    if (deviceOptions.value.length > 0) {
      deviceId.value = deviceOptions.value[0].value;
      await onDeviceChange(deviceId.value);
    }
  }

  async function onDeviceChange(id: number | null) {
    deviceId.value = id;
    points.value = [];
    allAlarms.value = [];
    device.value = null;
    if (!id) return;
    await loadOverview(true);
    await nextTick();
    ensureGauges();
    ensureCharts();
    await loadCharts(id);
    refreshGauges();
  }

  async function loadOverview(rebuild = false) {
    if (!deviceId.value) return;
    try {
      loading.value = true;
      const res: any = await Overview({ deviceId: deviceId.value, withAlarms: true } as any);
      if (!res) return;
      device.value = res.device;
      points.value = res.points || [];
      allAlarms.value = res.alarms || [];
      if (rebuild) {
        await nextTick();
        ensureGauges();
        ensureCharts();
      }
      refreshGauges();
    } finally {
      loading.value = false;
    }
  }

  function onPlcRealtimeMessage(message: WebSocketMessage) {
    const data = message?.data || {};
    if (!deviceId.value || Number(data.deviceId) !== Number(deviceId.value)) return;
    mergeRealtimePoints(data.points || []);
    refreshGauges();
  }

  function mergeRealtimePoints(realtimePoints: any[]) {
    if (!Array.isArray(realtimePoints) || realtimePoints.length === 0) return;
    const map = new Map([...points.value, ...allAlarms.value].map((p: any) => [p.field, { ...p }]));
    for (const item of realtimePoints) {
      const old = map.get(item.field) || {};
      const next = {
        ...old,
        ...item,
        dataType: old.dataType || item.dataType,
        scale: old.scale,
        offsetVal: old.offsetVal,
        sort: old.sort ?? item.sort ?? 0,
      };
      if (
        String(next.field || '')
          .toUpperCase()
          .startsWith('AL_')
      ) {
        next.active = Number(next.engValue || 0) !== 0;
        next.stateText = next.active ? '报警' : '正常';
      } else if (next.dataType === 'Bool') {
        next.active = Number(next.engValue || 0) !== 0;
        next.stateText = next.active ? '运行' : '停止';
      }
      map.set(item.field, next);
    }
    const nextPoints: any[] = [];
    const nextAlarms: any[] = [];
    for (const p of map.values()) {
      if (
        String(p.field || '')
          .toUpperCase()
          .startsWith('AL_')
      ) {
        nextAlarms.push(p);
      } else {
        nextPoints.push(p);
      }
    }
    points.value = nextPoints.sort((a, b) => Number(a.sort || 0) - Number(b.sort || 0));
    allAlarms.value = nextAlarms.sort((a, b) => Number(a.sort || 0) - Number(b.sort || 0));
  }

  function formatValue(v: any) {
    if (v == null) return '—';
    const n = Number(v);
    if (Number.isNaN(n)) return v;
    return Math.abs(n) >= 100 ? n.toFixed(1) : n.toFixed(2);
  }
  function cellClass(p: any) {
    if (!p) return '';
    if (p.alarmType === 1) return 'val-high';
    if (p.alarmType === 2) return 'val-low';
    return '';
  }

  function startClock() {
    const tick = () => {
      const d = new Date();
      clock.value = d.toLocaleString('zh-CN', { hour12: false });
    };
    tick();
    clockTimer = setInterval(tick, 1000);
  }

  function startCarousel() {
    carouselTimer = setInterval(() => {
      if (allAlarms.value.length > 0) {
        alarmCarouselIndex.value = (alarmCarouselIndex.value + 1) % allAlarms.value.length;
      }
      if (statusList.value.length > 0) {
        statusCarouselIndex.value = (statusCarouselIndex.value + 1) % statusList.value.length;
      }
      nextTick(() => {
        screenRef.value
          ?.querySelector('.status-row.is-active')
          ?.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
        screenRef.value
          ?.querySelector('.alarm-row.is-active')
          ?.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
      });
    }, 3000);
  }

  async function toggleFullscreen() {
    try {
      if (!document.fullscreenElement) {
        await screenRef.value?.requestFullscreen();
      } else {
        await document.exitFullscreen();
      }
    } catch (err) {
      console.log('[PLC Dashboard] fullscreen toggle failed:', err);
    }
  }

  function onFullscreenChange() {
    isFullscreen.value = !!document.fullscreenElement;
    setTimeout(() => {
      oilLevelGauge?.resize();
      currentGauge?.resize();
      lubeStatusGauge?.resize();
      tempChart?.resize();
      currentChart?.resize();
    }, 120);
  }

  watch(
    () => [
      tempChartRef.value,
      currentChartRef.value,
      oilLevelGaugeRef.value,
      currentGaugeRef.value,
      lubeStatusGaugeRef.value,
    ],
    () => {
      ensureGauges();
      ensureCharts();
    }
  );

  onUnmounted(() => {
    document.removeEventListener('fullscreenchange', onFullscreenChange);
    removeOnMessage(SocketEnum.EventPlcRealtime);
    if (clockTimer) clearInterval(clockTimer);
    if (refreshTimer) clearInterval(refreshTimer);
    if (chartTimer) clearInterval(chartTimer);
    if (carouselTimer) clearInterval(carouselTimer);
    oilLevelGauge?.dispose();
    currentGauge?.dispose();
    lubeStatusGauge?.dispose();
    tempChart?.dispose();
    currentChart?.dispose();
  });
</script>

<style lang="less" scoped>
  .hmi-screen {
    position: relative;
    min-height: calc(100vh - 80px);
    padding: 12px;
    box-sizing: border-box;
    background: linear-gradient(90deg, rgba(38, 111, 150, 0.08) 1px, transparent 1px),
      linear-gradient(180deg, rgba(38, 111, 150, 0.08) 1px, transparent 1px),
      radial-gradient(circle at 18% 8%, rgba(47, 112, 160, 0.16), transparent 26%),
      radial-gradient(circle at 88% 24%, rgba(19, 158, 190, 0.12), transparent 30%),
      linear-gradient(180deg, #eef3f7 0%, #dfe9f1 48%, #eef4f8 100%);
    background-size:
      26px 26px,
      26px 26px,
      auto,
      auto,
      auto;
    color: #253340;
    font-family: 'Microsoft YaHei', SimSun, sans-serif;
    overflow: hidden;
  }
  .hmi-screen:fullscreen {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    min-height: 100vh;
    padding: 12px;
    box-sizing: border-box;
    overflow: hidden;
  }
  .hmi-screen:fullscreen .hmi-top {
    flex-shrink: 0;
  }
  .hmi-screen:fullscreen .hmi-body {
    flex: 1;
    min-height: 0;
  }
  .hmi-screen:fullscreen .block-device {
    min-height: 0;
    height: 100%;
  }
  .hmi-screen:fullscreen .device-canvas,
  .hmi-screen:fullscreen .machine-stage {
    flex: 1;
    min-height: 0;
    height: 100%;
  }
  .hmi-screen::before {
    content: '';
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: radial-gradient(circle at 50% 0, rgba(255, 255, 255, 0.5), transparent 36%);
  }
  .hmi-screen::after {
    content: '';
    position: absolute;
    inset: 12px;
    pointer-events: none;
    border: 1px solid rgba(24, 144, 255, 0.1);
    border-radius: 18px;
    box-shadow: none;
  }

  // 顶栏
  .hmi-top {
    z-index: 1;
    display: flex;
    align-items: center;
    height: 58px;
    background: linear-gradient(90deg, #1890ff 0%, #096dd9 100%);
    color: #ffffff;
    border: none;
    border-radius: 16px;
    box-shadow: 0 6px 24px rgba(24, 144, 255, 0.22);
    margin-bottom: 10px;
    position: relative;
    overflow: hidden;

    .top-deco-left {
      width: 52px;
      height: 100%;
      background: rgba(255, 255, 255, 0.15);
      clip-path: polygon(0 0, 100% 0, 62% 100%, 0 100%);
      flex-shrink: 0;
    }
    .top-title {
      flex: 1;
      text-align: center;
      font-family: 'Microsoft YaHei', sans-serif;
      font-size: 20px;
      font-weight: 700;
      letter-spacing: 4px;
      color: #ffffff;
      text-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    }
    .top-deco-right {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 0 16px;
      .clock {
        background: rgba(0, 0, 0, 0.15);
        color: #ffffff;
        font-family: Menlo, Consolas, monospace;
        padding: 5px 11px;
        font-size: 13px;
        border-radius: 999px;
        border: 1px solid rgba(255, 255, 255, 0.25);
      }
    }
  }

  // 主体
  .hmi-body {
    position: relative;
    z-index: 1;
    display: grid;
    grid-template-columns: 280px minmax(680px, 1fr) 280px;
    grid-template-rows: 1fr;
    gap: 10px;
    align-items: stretch;
  }

  .col {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  // 通用模块
  .block {
    position: relative;
    background: rgba(255, 255, 255, 0.76);
    border: 1px solid rgba(86, 125, 154, 0.18);
    border-radius: 16px;
    box-shadow:
      0 14px 34px rgba(47, 70, 88, 0.1),
      inset 0 1px 0 rgba(255, 255, 255, 0.7);
    overflow: hidden;
    backdrop-filter: blur(8px);
  }
  .block::before,
  .block::after {
    content: '';
    position: absolute;
    width: 24px;
    height: 24px;
    pointer-events: none;
    border-color: rgba(24, 144, 255, 0.35);
    z-index: 2;
  }
  .block::before {
    left: 8px;
    top: 8px;
    border-left: 1px solid;
    border-top: 1px solid;
  }
  .block::after {
    right: 8px;
    bottom: 8px;
    border-right: 1px solid;
    border-bottom: 1px solid;
  }
  .block-title {
    background: linear-gradient(90deg, rgba(24, 144, 255, 0.07), transparent 70%);
    color: #1a2e3a;
    font-size: 13px;
    font-weight: 700;
    padding: 9px 12px;
    letter-spacing: 1px;
    border-bottom: 1px solid rgba(24, 144, 255, 0.1);
    display: flex;
    align-items: center;
    justify-content: space-between;
    &::before {
      content: '';
      width: 7px;
      height: 7px;
      margin-right: 8px;
      border-radius: 50%;
      background: #1890ff;
      box-shadow: 0 0 8px rgba(24, 144, 255, 0.5);
      flex-shrink: 0;
    }
    .block-sub {
      font-size: 11px;
      font-weight: 400;
      color: #5a7a8a;
    }
  }

  :deep(.n-base-selection) {
    --n-border: 1px solid rgba(24, 144, 255, 0.3) !important;
    --n-border-active: 1px solid rgba(24, 144, 255, 0.7) !important;
    --n-border-focus: 1px solid rgba(24, 144, 255, 0.7) !important;
    --n-box-shadow-active: 0 0 0 2px rgba(24, 144, 255, 0.12) !important;
    --n-box-shadow-focus: 0 0 0 2px rgba(24, 144, 255, 0.12) !important;
    --n-color: #ffffff !important;
    --n-text-color: #1a2e3a !important;
    --n-placeholder-color: rgba(90, 122, 138, 0.7) !important;
    backdrop-filter: none;
  }
  :deep(.n-base-selection-label) {
    background: #ffffff !important;
    color: #1a2e3a !important;
  }
  :deep(.n-base-selection-placeholder) {
    color: rgba(90, 122, 138, 0.7) !important;
  }

  // 仪表盘
  .gauges {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
    padding: 12px;
  }
  .gauge {
    position: relative;
    height: 154px;
    border-radius: 16px;
    background: #ffffff;
    border: 1px solid rgba(200, 218, 230, 0.6);
    box-shadow:
      0 4px 16px rgba(36, 72, 94, 0.08),
      0 1px 4px rgba(36, 72, 94, 0.04);
    overflow: hidden;
  }
  .gauge-wide {
    grid-column: 1 / -1;
    height: 128px;
  }

  .block-params {
    flex: 1;
    display: flex;
    flex-direction: column;
    .param-table {
      flex: 1;
    }
  }

  // 参数表
  .param-table {
    position: relative;
    padding: 10px;
    background: linear-gradient(180deg, rgba(239, 249, 253, 0.52), rgba(255, 255, 255, 0.22));
    .param-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      position: relative;
      margin-bottom: 6px;
      padding: 8px 10px 8px 12px;
      border-bottom: 1px solid rgba(91, 125, 148, 0.12);
      border-left: 2px solid rgba(55, 213, 255, 0.42);
      border-radius: 10px;
      font-size: 12px;
      background: rgba(255, 255, 255, 0.42);
      box-shadow: inset 0 0 0 1px rgba(55, 213, 255, 0.06);
      transition:
        transform 0.18s ease,
        box-shadow 0.18s ease,
        background 0.18s ease;
      &:last-child {
        margin-bottom: 0;
      }
      &:nth-child(even) {
        background: rgba(50, 88, 114, 0.05);
      }
      &:hover {
        transform: translateX(3px);
        background: rgba(238, 250, 255, 0.75);
        box-shadow:
          inset 0 0 0 1px rgba(55, 213, 255, 0.18),
          0 8px 20px rgba(35, 88, 120, 0.08);
      }
    }
    .param-label {
      position: relative;
      padding-left: 10px;
      color: #667786;
      font-weight: 600;
    }
    .param-label::before {
      content: '';
      position: absolute;
      left: 0;
      top: 50%;
      width: 4px;
      height: 4px;
      border-radius: 50%;
      background: #1890ff;
      box-shadow: 0 0 5px rgba(24, 144, 255, 0.5);
      transform: translateY(-50%);
    }
    .param-value {
      color: #167ba8;
      font-weight: 700;
      font-family: Menlo, Consolas, monospace;
      small {
        color: #999;
        font-weight: 400;
        margin-left: 2px;
      }
    }
    .val-high {
      color: #ff2a2a;
    }
    .val-low {
      color: #ffb000;
    }
    .param-empty {
      text-align: center;
      color: #999;
      padding: 20px 0;
      font-size: 12px;
    }
  }

  // 设备图
  .col-mid {
    gap: 10px;
  }
  .block-device {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  .alarm-banner {
    margin: 10px 12px 0;
    padding: 7px 12px;
    background: rgba(255, 209, 102, 0.9);
    color: #3d2d05;
    font-weight: 700;
    border: 1px solid rgba(190, 137, 0, 0.28);
    border-radius: 10px;
    font-size: 13px;
    display: flex;
    align-items: center;
    gap: 8px;
    animation: blink 1.5s infinite;

    .alarm-banner-icon {
      font-size: 16px;
    }
  }
  .alarm-banner-ok {
    background: rgba(226, 248, 235, 0.9);
    color: #147342;
    border-color: rgba(42, 164, 97, 0.18);
    animation: none;
  }
  @keyframes blink {
    50% {
      background: #ffe900;
    }
  }
  .device-canvas {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 0;
    background: transparent;
    overflow: hidden;
    border-radius: 0 0 16px 16px;
  }
  .machine-stage {
    position: relative;
    flex: 1;
    min-height: 460px;
    container-type: inline-size;
    overflow: hidden;
    isolation: isolate;
    background: transparent;
  }
  .machine-stage::after {
    display: none;
  }
  .machine-bg {
    position: absolute;
    inset: 0;
    display: block;
    width: 100%;
    height: 100%;
    object-fit: contain;
    object-position: center center;
    background: rgba(255, 255, 255, 0.76);
    border: 0;
    outline: 0;
    box-shadow: none;
    z-index: 1;
    filter: none;
  }
  .readout-panel {
    position: absolute;
    right: 12px;
    top: 118px;
    z-index: 3;
    width: clamp(174px, 21cqw, 214px);
    padding: 12px;
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.95);
    border: 1px solid rgba(24, 144, 255, 0.18);
    box-shadow:
      0 8px 32px rgba(24, 144, 255, 0.1),
      0 2px 8px rgba(0, 0, 0, 0.06);
    backdrop-filter: blur(8px);
  }
  .readout-panel::after {
    display: none;
  }
  .readout-panel::before {
    content: '运行数据';
    display: block;
    margin-bottom: 8px;
    color: #1890ff;
    font-size: 13px;
    font-weight: 800;
    text-align: center;
    letter-spacing: 2px;
    border-bottom: 1px solid rgba(24, 144, 255, 0.12);
    padding-bottom: 6px;
  }
  .readout-card {
    position: relative;
    display: flex;
    align-items: baseline;
    min-height: 38px;
    padding: 0 2px;
    border-radius: 0;
    color: #1a2e3a;
    border: 0;
    border-bottom: 1px solid rgba(24, 144, 255, 0.08);
    background: transparent;
    box-shadow: none;
    backdrop-filter: none;
    white-space: nowrap;
  }
  .readout-card:last-child {
    border-bottom: 0;
  }
  .readout-label {
    display: inline-block;
    flex: 1;
    min-width: 78px;
    color: #5a7a8a;
    font-size: 12px;
    font-weight: 600;
  }
  .readout-card strong {
    color: #1a2e3a;
    font-family: Menlo, Consolas, monospace;
    font-size: clamp(15px, 1.9cqw, 18px);
    font-weight: 700;
    letter-spacing: 0.2px;
    margin-right: 4px;
    text-align: right;
  }
  .readout-card small {
    color: #9ab0bc;
    font-size: 11px;
  }
  .readout-card.val-power strong {
    color: #52c41a;
  }
  .readout-card.val-high strong {
    color: #ff4d4f;
  }
  .readout-card.val-low strong {
    color: #faad14;
  }
  .live-badge {
    position: absolute;
    right: 44px;
    bottom: 22px;
    z-index: 3;
    display: flex;
    align-items: center;
    gap: 7px;
    padding: 5px 10px;
    border-radius: 999px;
    color: rgba(40, 58, 70, 0.78);
    font-size: 12px;
    font-weight: 700;
    background: rgba(255, 255, 255, 0.62);
    border: 1px solid rgba(80, 120, 150, 0.16);
  }
  .live-badge span {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #1fd17a;
    box-shadow: 0 0 12px #1fd17a;
    animation: dotPulse 1.4s ease-in-out infinite;
  }
  @keyframes dotPulse {
    0%,
    100% {
      opacity: 1;
      transform: scale(1);
    }
    50% {
      opacity: 0.52;
      transform: scale(1.35);
    }
  }

  // 底部图表
  .bottom-charts {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
  }
  .chart-area {
    height: 190px;
    width: 100%;
    background: linear-gradient(90deg, rgba(24, 144, 255, 0.04) 1px, transparent 1px),
      linear-gradient(180deg, rgba(24, 144, 255, 0.04) 1px, transparent 1px);
    background-size: 18px 18px;
  }

  // 状态列表
  .status-list {
    padding: 10px;
    max-height: 280px;
    overflow-y: auto;
  }
  .status-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 6px;
    padding: 8px 10px;
    border-bottom: 1px solid rgba(91, 125, 148, 0.12);
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.42);
    font-size: 12px;
    &:last-child {
      border-bottom: none;
    }
    .status-name {
      flex: 1;
      color: #334655;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .status-led {
      width: 12px;
      height: 12px;
      border-radius: 50%;
      flex-shrink: 0;
      border: 2px solid rgba(255, 255, 255, 0.72);
    }
    .led-on {
      background: #2dff2d;
      box-shadow: 0 0 6px #2dff2d;
    }
    .led-off {
      background: #c0c4cc;
    }
    .status-text {
      font-size: 11px;
      font-weight: 600;
      width: 40px;
      text-align: right;
    }
    .txt-on {
      color: #147342;
    }
    .txt-off {
      color: #999;
    }
  }
  .status-empty {
    text-align: center;
    color: #999;
    padding: 20px 0;
    font-size: 12px;
  }
  .status-row.is-active {
    background: rgba(45, 180, 255, 0.18);
    border-color: rgba(45, 180, 255, 0.4);
    box-shadow: 0 0 0 1px rgba(45, 180, 255, 0.25);
    transition:
      background 0.4s,
      box-shadow 0.4s;
  }

  // 报警
  .block-alarm {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  .alarm-list {
    padding: 10px;
    flex: 1;
    overflow-y: auto;
    max-height: 380px;
  }
  .alarm-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 7px;
    padding: 8px 10px;
    border-bottom: 1px solid rgba(91, 125, 148, 0.12);
    border-radius: 12px;
    font-size: 12px;
    background: rgba(255, 255, 255, 0.38);

    .alarm-led {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      flex-shrink: 0;
    }
    .alarm-name {
      flex: 1;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .alarm-tag {
      font-size: 11px;
      font-weight: 700;
      padding: 2px 6px;
      border-radius: 999px;
    }
  }
  .alarm-row.is-active {
    background: rgba(255, 200, 60, 0.18);
    border-color: rgba(255, 200, 60, 0.4);
    box-shadow: 0 0 0 1px rgba(255, 200, 60, 0.25);
    transition:
      background 0.4s,
      box-shadow 0.4s;
  }
  .alarm-on {
    background: rgba(255, 232, 232, 0.7);
    .alarm-led {
      background: #ff2a2a;
      box-shadow: 0 0 6px #ff2a2a;
      animation: pulse 1.4s infinite;
    }
    .alarm-name {
      color: #ff2a2a;
      font-weight: 600;
    }
    .alarm-tag {
      background: #ff2a2a;
      color: #fff;
    }
  }
  .alarm-off {
    .alarm-led {
      background: #c0c4cc;
    }
    .alarm-name {
      color: #888;
    }
    .alarm-tag {
      background: rgba(91, 125, 148, 0.12);
      color: #73818d;
    }
  }
  .alarm-empty {
    text-align: center;
    color: #999;
    padding: 30px 0;
    font-size: 12px;
  }

  @keyframes pulse {
    0% {
      transform: scale(1);
      opacity: 1;
    }
    50% {
      transform: scale(1.5);
      opacity: 0.4;
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }

  // 浅色主题覆盖
  .hmi-screen {
    --panel-bg: #ffffff;
    --panel-border: rgba(180, 210, 230, 0.55);
    --accent-blue: #1890ff;
    --accent-orange: #f5a623;
    --accent-green: #52c41a;
    --accent-red: #ff4d4f;
    --accent-amber: #faad14;
    --text-primary: #1a2e3a;
    --text-secondary: #5a7a8a;
    --text-muted: #9ab0bc;

    min-height: calc(100vh - 80px);
    padding: 14px;
    color: var(--text-primary);
    font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
    background: linear-gradient(90deg, rgba(24, 144, 255, 0.05) 1px, transparent 1px),
      linear-gradient(180deg, rgba(24, 144, 255, 0.05) 1px, transparent 1px),
      radial-gradient(circle at 10% 10%, rgba(24, 144, 255, 0.08), transparent 28%),
      radial-gradient(circle at 90% 6%, rgba(245, 166, 35, 0.07), transparent 26%),
      linear-gradient(160deg, #eef4fb 0%, #e6f0f8 50%, #edf3fa 100%);
    background-size:
      28px 28px,
      28px 28px,
      auto,
      auto,
      auto;
  }
  .hmi-screen::before {
    background: none;
    animation: none;
  }
  .hmi-screen::after {
    inset: 10px;
    border-color: rgba(24, 144, 255, 0.12);
    border-radius: 22px;
    box-shadow: inset 0 0 60px rgba(24, 144, 255, 0.04);
  }
  .warehouse-orbit {
    display: none;
  }

  .hmi-top {
    height: 62px;
    background: linear-gradient(90deg, #1890ff 0%, #096dd9 100%);
    border-color: rgba(24, 144, 255, 0.3);
    border-radius: 16px;
    box-shadow:
      0 6px 24px rgba(24, 144, 255, 0.22),
      0 2px 6px rgba(24, 144, 255, 0.12);
  }
  .hmi-top::before {
    display: none;
  }
  .hmi-top .top-deco-left {
    width: 60px;
    background: rgba(255, 255, 255, 0.18);
    clip-path: polygon(0 0, 100% 0, 62% 100%, 0 100%);
  }
  .hmi-top .top-title {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2px;
    letter-spacing: 4px;
  }
  .hmi-top .top-title strong {
    color: #ffffff;
    font-size: 22px;
    line-height: 1;
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.18);
  }
  .title-kicker {
    color: rgba(255, 255, 255, 0.72);
    font-size: 10px;
    letter-spacing: 5px;
  }
  .hmi-top .top-deco-right .clock {
    color: #ffffff;
    background: rgba(0, 0, 0, 0.18);
    border-color: rgba(255, 255, 255, 0.28);
    box-shadow: none;
  }
  .fullscreen-btn {
    min-width: 74px;
    color: #ffffff;
    border-color: rgba(255, 255, 255, 0.4);
    background: rgba(255, 255, 255, 0.12);
    box-shadow: none;
  }
  .fullscreen-btn:hover {
    color: #ffffff;
    border-color: rgba(255, 255, 255, 0.8);
    background: rgba(255, 255, 255, 0.22);
  }
  :deep(.n-base-selection-label) {
    background: #ffffff !important;
    color: #1a2e3a !important;
  }
  :deep(.n-base-selection-input) {
    color: #1a2e3a !important;
  }

  .warehouse-strip {
    position: relative;
    z-index: 1;
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 10px;
    margin: 0 0 10px;
  }
  .warehouse-chip {
    position: relative;
    min-height: 76px;
    padding: 12px 14px;
    overflow: hidden;
    border: 1px solid var(--panel-border);
    border-radius: 14px;
    background: #ffffff;
    box-shadow:
      0 2px 10px rgba(24, 144, 255, 0.08),
      0 1px 3px rgba(0, 0, 0, 0.04);
  }
  .warehouse-chip::after {
    content: '';
    position: absolute;
    left: 14px;
    right: 14px;
    bottom: 0;
    height: 3px;
    border-radius: 999px 999px 0 0;
    background: currentColor;
    opacity: 0.28;
  }
  .warehouse-chip::before {
    display: none;
  }
  .warehouse-chip .chip-label {
    display: block;
    color: var(--text-muted);
    font-size: 10px;
    letter-spacing: 1px;
  }
  .warehouse-chip strong {
    display: block;
    margin-top: 4px;
    color: currentColor;
    font-size: 24px;
    font-weight: 700;
    line-height: 1;
  }
  .chip-value {
    animation: chipRise 0.5s cubic-bezier(0.2, 0.86, 0.28, 1.08);
  }
  .warehouse-chip small {
    display: block;
    margin-top: 4px;
    color: var(--text-muted);
    font-size: 11px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .chip-desc {
    animation: chipRise 0.45s ease both;
  }
  .tone-cyan {
    color: var(--accent-blue);
  }
  .tone-blue {
    color: #096dd9;
  }
  .tone-green {
    color: var(--accent-green);
  }
  .tone-amber {
    color: var(--accent-amber);
  }
  .tone-red {
    color: var(--accent-red);
  }

  @keyframes chipRise {
    0% {
      opacity: 0;
      transform: translateY(14px);
      filter: blur(3px);
    }
    100% {
      opacity: 1;
      transform: translateY(0);
      filter: blur(0);
    }
  }

  .hmi-body {
    grid-template-columns: 300px minmax(720px, 1fr) 300px;
    gap: 12px;
  }
  .block {
    background: #ffffff;
    border-color: var(--panel-border);
    border-radius: 16px;
    box-shadow:
      0 2px 12px rgba(24, 144, 255, 0.07),
      0 1px 3px rgba(0, 0, 0, 0.04);
  }
  .block::before,
  .block::after {
    width: 28px;
    height: 28px;
    border-color: rgba(24, 144, 255, 0.4);
  }
  .block-title {
    min-height: 40px;
    padding: 10px 14px;
    background: linear-gradient(90deg, rgba(24, 144, 255, 0.06), transparent 70%);
    color: var(--text-primary);
    border-bottom-color: rgba(24, 144, 255, 0.1);
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 1px;
  }
  .block-title::before {
    background: var(--accent-blue);
    box-shadow: 0 0 8px rgba(24, 144, 255, 0.5);
  }

  .gauges,
  .chart-area,
  .param-table,
  .status-list,
  .alarm-list {
    background: none;
  }
  .gauge {
    background: #f7fafd;
    border-color: rgba(180, 210, 230, 0.5);
    box-shadow:
      0 2px 8px rgba(24, 144, 255, 0.06),
      inset 0 1px 0 rgba(255, 255, 255, 0.9);
  }
  .param-table .param-row,
  .status-row,
  .alarm-row {
    color: var(--text-primary);
    background: #f7fafd;
    border-color: rgba(180, 210, 230, 0.4);
    border-left-color: var(--accent-blue);
  }
  .param-table .param-row:hover,
  .status-row:hover,
  .alarm-row:hover {
    background: rgba(24, 144, 255, 0.06);
    box-shadow: inset 0 0 0 1px rgba(24, 144, 255, 0.14);
  }
  .hmi-screen .status-row.is-active {
    background: rgba(24, 144, 255, 0.1);
    border-color: rgba(24, 144, 255, 0.4);
    box-shadow: 0 0 0 1px rgba(24, 144, 255, 0.2);
  }
  .hmi-screen .alarm-row.is-active {
    background: rgba(250, 173, 20, 0.1);
    border-color: rgba(250, 173, 20, 0.4);
    box-shadow: 0 0 0 1px rgba(250, 173, 20, 0.22);
  }
  .param-table .param-label,
  .status-row .status-name,
  .alarm-row .alarm-name {
    color: var(--text-secondary);
  }
  .param-table .param-value {
    color: var(--accent-blue);
  }
  .param-table .param-value small,
  .status-empty,
  .alarm-empty {
    color: var(--text-muted);
  }

  .alarm-banner {
    background: rgba(255, 77, 79, 0.08);
    color: #cf1322;
    border-color: rgba(255, 77, 79, 0.3);
    box-shadow: none;
  }
  .alarm-banner-ok {
    background: rgba(82, 196, 26, 0.08);
    color: #389e0d;
    border-color: rgba(82, 196, 26, 0.28);
  }
  .device-canvas,
  .machine-stage {
    background: rgba(255, 255, 255, 0.76);
  }
  .machine-stage::after {
    background: linear-gradient(
      90deg,
      rgba(255, 255, 255, 0.5),
      transparent 20%,
      transparent 80%,
      rgba(255, 255, 255, 0.5)
    );
  }
  .readout-panel,
  .readout-panel::after {
    background: rgba(255, 255, 255, 0.96);
    border-color: rgba(180, 210, 230, 0.6);
  }
  .readout-panel::before {
    color: var(--accent-blue);
  }
  .readout-card {
    color: var(--text-primary);
    border-bottom-color: rgba(180, 210, 230, 0.3);
  }
  .readout-label {
    color: var(--text-muted);
  }
  .readout-card strong {
    color: var(--text-primary);
  }
  .live-badge {
    color: var(--accent-green);
    background: rgba(255, 255, 255, 0.88);
    border-color: rgba(82, 196, 26, 0.3);
  }

  .status-row .led-on {
    background: var(--accent-green);
    box-shadow: 0 0 8px rgba(82, 196, 26, 0.6);
  }
  .status-row .led-off,
  .alarm-off .alarm-led {
    background: #c8d8e2;
  }
  .status-row .txt-on {
    color: var(--accent-green);
  }
  .status-row .txt-off,
  .alarm-off .alarm-name {
    color: var(--text-muted);
  }
  .alarm-on {
    background: rgba(255, 77, 79, 0.06);
  }
  .alarm-on .alarm-led {
    background: var(--accent-red);
    box-shadow: 0 0 8px rgba(255, 77, 79, 0.5);
  }
  .alarm-on .alarm-name {
    color: #cf1322;
  }
  .alarm-on .alarm-tag {
    background: var(--accent-red);
    color: #fff;
  }
  .alarm-off .alarm-tag {
    color: var(--text-muted);
    background: rgba(180, 210, 230, 0.25);
  }

  @media screen and (max-width: 1360px) {
    .warehouse-strip {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
    .hmi-body {
      grid-template-columns: 1fr;
    }
    .col-left,
    .col-right {
      display: grid;
      grid-template-columns: 1fr 1fr;
    }
    .machine-stage {
      min-height: 460px;
    }
  }

  @media screen and (max-width: 768px) {
    .warehouse-strip {
      grid-template-columns: 1fr;
    }
    .hmi-top {
      height: auto;
      min-height: 56px;
      padding: 8px 0;
      flex-wrap: wrap;
      .top-title {
        min-width: 100%;
        order: -1;
        font-size: 18px;
      }
      .top-deco-right {
        width: 100%;
        justify-content: center;
        flex-wrap: wrap;
      }
    }
    .col-left,
    .col-right,
    .bottom-charts {
      grid-template-columns: 1fr;
    }
    .machine-stage {
      min-height: 360px;
    }
    .readout-panel {
      right: 10px;
      top: 112px;
      width: 168px;
      padding: 9px;
    }
    .readout-card {
      min-height: 34px;
    }
    .readout-label {
      min-width: 66px;
      font-size: 11px;
    }
    .readout-card strong {
      font-size: 14px;
    }
  }

  @media screen and (max-width: 520px) {
    .machine-stage {
      min-height: 420px;
    }
    .readout-card {
      min-height: 32px;
    }
    .readout-panel {
      right: 8px;
      top: 104px;
      width: 148px;
      padding: 8px;
    }
    .live-badge {
      right: 10px;
      bottom: 10px;
    }
  }
</style>
