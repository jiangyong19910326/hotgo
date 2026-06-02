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

        <div class="block">
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
              <canvas
                ref="modelCanvasRef"
                class="machine-model"
                aria-label="设备三维模型"
                @mouseenter="setModelZoomEnabled(true)"
                @mouseleave="setModelZoomEnabled(false)"
                @wheel.stop
              ></canvas>
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
              <button class="model-rotate-btn" type="button" @click="toggleModelRotate">
                {{ modelAutoRotate ? '停止旋转' : '旋转模型' }}
              </button>
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
            <div v-for="s in statusList" :key="s.field" class="status-row">
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
              v-for="a in allAlarms"
              :key="a.pointId"
              class="alarm-row"
              :class="{ 'alarm-on': a.active, 'alarm-off': !a.active }"
            >
              <span class="alarm-led"></span>
              <span class="alarm-name" :title="a.field">{{ a.name || a.field }}</span>
              <span class="alarm-tag">{{ a.active ? '报警' : '正常' }}</span>
            </div>
            <div v-if="allAlarms.length === 0" class="alarm-empty">无报警点位</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue';
  import * as echarts from 'echarts';
  import * as THREE from 'three';
  import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';
  import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader';
  import { MineOptions, DeviceList, Overview } from '@/api/plc';
  import { http } from '@/utils/http/axios';
  import { SocketEnum } from '@/enums/socketEnum';
  import { addOnMessage, removeOnMessage, WebSocketMessage } from '@/utils/websocket';
  import hp300ModelUrl from '@/assets/models/zxhp300.glb?url';
  import sandMakerModelUrl from '@/assets/models/sand-maker-1263.glb?url';

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
        value: current.value,
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
      },
      {
        label: '油箱温度',
        keys: ['lube_tank_temp', 'tank_temp', 'oil_tank_temp', '润滑油箱温度', '油箱温度'],
        unit: '℃',
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
          value: isBool ? statusText(p) : formatValue(p.engValue),
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

  function gaugeOption(name: string, value: number, max: number, unit: string, color: string) {
    const safeValue = Math.max(0, Math.min(max, Number(value || 0)));
    return {
      backgroundColor: 'transparent',
      tooltip: {
        formatter: `${name}<br/>${safeValue.toFixed(2)}${unit}`,
        backgroundColor: 'rgba(18, 30, 42, 0.92)',
        borderColor: 'rgba(63, 214, 255, 0.35)',
        textStyle: { color: '#e9f7ff' },
      },
      graphic: [
        {
          type: 'circle',
          left: 'center',
          top: 'middle',
          shape: { r: 46 },
          style: {
            fill: 'rgba(28, 52, 70, 0.35)',
            stroke: 'rgba(63, 214, 255, 0.18)',
            lineWidth: 1,
            shadowBlur: 18,
            shadowColor: 'rgba(63, 214, 255, 0.18)',
          },
          silent: true,
        },
      ],
      series: [
        {
          type: 'gauge',
          radius: '92%',
          min: 0,
          max,
          splitNumber: 4,
          startAngle: 220,
          endAngle: -40,
          progress: {
            show: true,
            roundCap: true,
            width: 10,
            itemStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 1,
                y2: 0,
                colorStops: [
                  { offset: 0, color: '#37d5ff' },
                  { offset: 0.55, color },
                  { offset: 1, color: '#4de18c' },
                ],
              },
              shadowBlur: 14,
              shadowColor: color,
            },
          },
          axisLine: {
            lineStyle: {
              width: 10,
              color: [[1, 'rgba(77, 229, 255, 0.16)']],
            },
          },
          axisTick: {
            distance: -16,
            length: 4,
            lineStyle: { color: 'rgba(215, 247, 255, 0.72)', width: 1 },
          },
          splitLine: {
            distance: -18,
            length: 10,
            lineStyle: { color: 'rgba(77, 229, 255, 0.86)', width: 1.4 },
          },
          axisLabel: { color: 'rgba(215, 247, 255, 0.78)', distance: 8, fontSize: 9 },
          pointer: {
            icon: 'path://M-3,0 L0,-62 L3,0 Z',
            length: '58%',
            width: 8,
            itemStyle: {
              color: '#e9fbff',
              shadowBlur: 10,
              shadowColor: color,
            },
          },
          anchor: {
            show: true,
            size: 15,
            itemStyle: {
              color: '#1f394d',
              borderColor: '#37d5ff',
              borderWidth: 2,
              shadowBlur: 12,
              shadowColor: 'rgba(55, 213, 255, 0.7)',
            },
          },
          title: {
            show: true,
            color: 'rgba(215, 247, 255, 0.86)',
            fontSize: 12,
            fontWeight: 700,
            offsetCenter: [0, '72%'],
          },
          detail: {
            valueAnimation: true,
            formatter: `{value}${unit}`,
            color: '#eaffff',
            fontSize: 20,
            offsetCenter: [0, '42%'],
            fontWeight: 700,
            textShadowColor: color,
            textShadowBlur: 12,
          },
          data: [{ value: Number(safeValue.toFixed(2)), name }],
        },
        {
          type: 'gauge',
          radius: '66%',
          min: 0,
          max,
          startAngle: 220,
          endAngle: -40,
          axisLine: { lineStyle: { width: 1, color: [[1, 'rgba(77, 229, 255, 0.24)']] } },
          axisTick: { show: false },
          splitLine: { show: false },
          axisLabel: { show: false },
          pointer: { show: false },
          anchor: { show: false },
          detail: { show: false },
          title: { show: false },
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
      ? '#8ea3ad'
      : lubeStatusPoint
        ? lubeNormal
          ? '#48f5a5'
          : '#ff5c7a'
        : '#ffc857';

    oilLevelGauge?.setOption(
      gaugeOption('液压油油位', oilLevel, 100, '%', oilLevelColor(oilLevel))
    );
    currentGauge?.setOption(gaugeOption('主机电流', current, 200, 'A', '#4de18c'));
    lubeStatusGauge?.setOption(statusGaugeOption('润滑油状态', lubeText, lubeColor));
  }

  function normalizePercent(value: any) {
    const n = Number(value ?? 0);
    if (!Number.isFinite(n)) return 0;
    return Math.max(0, Math.min(100, n <= 1 ? n * 100 : n));
  }

  function oilLevelColor(value: number) {
    if (value <= 20) return '#ff5c7a';
    if (value >= 80) return '#ffc857';
    return '#4de5ff';
  }

  function statusGaugeOption(name: string, text: string, color: string) {
    const value = text === '正常' ? 1 : text === '异常' ? 0.25 : 0.55;
    return {
      backgroundColor: 'transparent',
      tooltip: { formatter: `${name}<br/>${text}` },
      series: [
        {
          type: 'gauge',
          radius: '88%',
          min: 0,
          max: 1,
          startAngle: 210,
          endAngle: -30,
          progress: {
            show: true,
            width: 12,
            roundCap: true,
            itemStyle: { color, shadowBlur: 14, shadowColor: color },
          },
          axisLine: { lineStyle: { width: 12, color: [[1, 'rgba(77, 229, 255, 0.14)']] } },
          axisTick: { show: false },
          splitLine: { show: false },
          axisLabel: { show: false },
          pointer: { show: false },
          anchor: { show: false },
          title: {
            color: 'rgba(215, 247, 255, 0.86)',
            fontSize: 12,
            fontWeight: 700,
            offsetCenter: [0, '64%'],
          },
          detail: {
            formatter: () => text,
            color: '#eaffff',
            fontSize: 22,
            fontWeight: 800,
            offsetCenter: [0, '10%'],
            textShadowColor: color,
            textShadowBlur: 14,
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

  const modelCanvasRef = ref<HTMLCanvasElement>();
  const modelAutoRotate = ref(false);
  const modelZoomEnabled = ref(false);
  const gltfLoader = new GLTFLoader();
  let modelRenderer: THREE.WebGLRenderer | null = null;
  let modelScene: THREE.Scene | null = null;
  let modelCamera: THREE.PerspectiveCamera | null = null;
  let modelControls: OrbitControls | null = null;
  let activeModel: THREE.Object3D | null = null;
  let modelFrame = 0;
  let currentModelUrl = '';

  const deviceModelUrl = computed(() => {
    const text =
      `${device.value?.name || ''} ${device.value?.host || ''} ${device.value?.remark || ''}`.toLowerCase();
    return text.includes('制砂') || text.includes('sand') || text.includes('1263')
      ? sandMakerModelUrl
      : hp300ModelUrl;
  });

  function ensureCharts() {
    if (tempChartRef.value && !tempChart) tempChart = echarts.init(tempChartRef.value);
    if (currentChartRef.value && !currentChart) currentChart = echarts.init(currentChartRef.value);
  }

  function ensureModelViewer() {
    if (!modelCanvasRef.value || modelRenderer) return;
    modelScene = new THREE.Scene();
    modelScene.add(new THREE.HemisphereLight(0xdff8ff, 0x06111a, 2.4));
    const keyLight = new THREE.DirectionalLight(0x7deaff, 3.2);
    keyLight.position.set(4, 6, 5);
    modelScene.add(keyLight);
    const rimLight = new THREE.DirectionalLight(0x48f5a5, 1.6);
    rimLight.position.set(-5, 3, -4);
    modelScene.add(rimLight);

    modelCamera = new THREE.PerspectiveCamera(38, 1, 0.1, 1000);
    modelCamera.position.set(2.2, 1.7, 4.4);
    modelRenderer = new THREE.WebGLRenderer({
      canvas: modelCanvasRef.value,
      alpha: true,
      antialias: true,
    });
    modelRenderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    modelRenderer.outputColorSpace = THREE.SRGBColorSpace;
    modelControls = new OrbitControls(modelCamera, modelCanvasRef.value);
    modelControls.enableDamping = true;
    modelControls.autoRotate = modelAutoRotate.value;
    modelControls.autoRotateSpeed = 0.45;
    modelControls.enableZoom = modelZoomEnabled.value;
    modelControls.enablePan = false;
    modelControls.target.set(0, 0, 0);
    modelControls.minDistance = 1.8;
    modelControls.maxDistance = 8;
    resizeModelViewer();
    animateModel();
  }

  function toggleModelRotate() {
    modelAutoRotate.value = !modelAutoRotate.value;
    if (modelControls) modelControls.autoRotate = modelAutoRotate.value;
  }

  function setModelZoomEnabled(enabled: boolean) {
    modelZoomEnabled.value = enabled;
    if (modelControls) modelControls.enableZoom = enabled;
  }

  async function loadDeviceModel() {
    ensureModelViewer();
    if (!modelScene || !deviceModelUrl.value || currentModelUrl === deviceModelUrl.value) return;
    currentModelUrl = deviceModelUrl.value;
    if (activeModel) {
      modelScene.remove(activeModel);
      activeModel.traverse((child: any) => {
        child.geometry?.dispose?.();
        child.material?.dispose?.();
      });
      activeModel = null;
    }
    const gltf = await gltfLoader.loadAsync(deviceModelUrl.value);
    activeModel = gltf.scene;
    fitModel(activeModel);
    modelScene.add(activeModel);
  }

  function fitModel(model: THREE.Object3D) {
    const box = new THREE.Box3().setFromObject(model);
    const size = box.getSize(new THREE.Vector3());
    const center = box.getCenter(new THREE.Vector3());
    model.position.sub(center);
    const scale = 2.7 / Math.max(size.x, size.y, size.z, 1);
    model.scale.setScalar(scale);
    model.rotation.y = -Math.PI / 8;
    model.position.y = -0.08;
    model.position.x = -0.18;
    model.traverse((child: any) => {
      if (child.isMesh) {
        child.castShadow = true;
        child.receiveShadow = true;
      }
    });
  }

  function resizeModelViewer() {
    if (!modelCanvasRef.value || !modelRenderer || !modelCamera) return;
    const rect = modelCanvasRef.value.getBoundingClientRect();
    const width = Math.max(1, rect.width);
    const height = Math.max(1, rect.height);
    modelRenderer.setSize(width, height, false);
    modelCamera.aspect = width / height;
    modelCamera.updateProjectionMatrix();
  }

  function animateModel() {
    modelFrame = requestAnimationFrame(animateModel);
    if (modelControls) modelControls.autoRotate = modelAutoRotate.value;
    if (modelControls) modelControls.enableZoom = modelZoomEnabled.value;
    modelControls?.update();
    if (modelRenderer && modelScene && modelCamera) modelRenderer.render(modelScene, modelCamera);
  }

  function disposeModelViewer() {
    cancelAnimationFrame(modelFrame);
    modelControls?.dispose();
    if (activeModel) {
      activeModel.traverse((child: any) => {
        child.geometry?.dispose?.();
        child.material?.dispose?.();
      });
    }
    modelRenderer?.dispose();
  }

  function lineOption(legend: string[], xAxis: string[], series: any[], palette: string[]) {
    return {
      backgroundColor: 'transparent',
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(18, 30, 42, 0.92)',
        borderColor: 'rgba(63, 214, 255, 0.3)',
        textStyle: { color: '#e9f7ff' },
      },
      legend: {
        data: legend,
        textStyle: { color: 'rgba(215, 247, 255, 0.68)', fontSize: 11 },
        top: 4,
      },
      grid: { left: 50, right: 16, top: 32, bottom: 28 },
      xAxis: {
        type: 'category',
        data: xAxis,
        axisLine: { lineStyle: { color: 'rgba(77, 229, 255, 0.2)' } },
        axisLabel: { color: 'rgba(215, 247, 255, 0.5)', fontSize: 9 },
      },
      yAxis: {
        type: 'value',
        axisLine: { lineStyle: { color: 'rgba(77, 229, 255, 0.2)' } },
        splitLine: { lineStyle: { color: 'rgba(77, 229, 255, 0.1)' } },
        axisLabel: { color: 'rgba(215, 247, 255, 0.5)', fontSize: 9 },
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
          '#4de5ff',
          '#ffc857',
          '#ff5c7a',
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
          '#3b82f6',
          '#48f5a5',
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
    window.addEventListener('resize', resizeModelViewer);
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
    await loadDeviceModel();
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
        await loadDeviceModel();
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
      resizeModelViewer();
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
      ensureModelViewer();
    }
  );

  watch(deviceModelUrl, () => {
    loadDeviceModel();
  });

  onUnmounted(() => {
    document.removeEventListener('fullscreenchange', onFullscreenChange);
    window.removeEventListener('resize', resizeModelViewer);
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
    disposeModelViewer();
  });
</script>

<style lang="less" scoped>
  .hmi-screen {
    position: relative;
    min-height: calc(100vh - 80px);
    padding: 12px;
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
  .hmi-screen::before {
    content: '';
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: linear-gradient(
        120deg,
        transparent 0 38%,
        rgba(55, 213, 255, 0.1) 48%,
        transparent 58% 100%
      ),
      radial-gradient(circle at 50% 0, rgba(255, 255, 255, 0.62), transparent 36%);
    mix-blend-mode: screen;
  }
  .hmi-screen::after {
    content: '';
    position: absolute;
    inset: 12px;
    pointer-events: none;
    border: 1px solid rgba(34, 112, 155, 0.12);
    border-radius: 18px;
    box-shadow: inset 0 0 42px rgba(55, 213, 255, 0.08);
  }

  // 顶栏
  .hmi-top {
    z-index: 1;
    display: flex;
    align-items: center;
    height: 56px;
    background: rgba(45, 58, 69, 0.92);
    color: #f8fbff;
    border: 1px solid rgba(90, 130, 160, 0.24);
    border-radius: 16px;
    box-shadow:
      0 18px 42px rgba(38, 58, 72, 0.16),
      inset 0 1px 0 rgba(255, 255, 255, 0.08);
    backdrop-filter: blur(10px);
    margin-bottom: 10px;
    position: relative;
    overflow: hidden;

    .top-deco-left {
      width: 46px;
      height: 100%;
      background: linear-gradient(135deg, #2ed4ff 0%, #1c7da3 100%);
      clip-path: polygon(0 0, 100% 0, 62% 100%, 0 100%);
    }
    .top-title {
      flex: 1;
      text-align: center;
      font-family: 'SimSun', '宋体', serif;
      font-size: 22px;
      font-weight: 700;
      letter-spacing: 4px;
      text-shadow: 0 0 18px rgba(68, 210, 255, 0.28);
    }
    .top-deco-right {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 0 16px;
      .clock {
        background: rgba(14, 24, 32, 0.82);
        color: #41d6ff;
        font-family: Menlo, Consolas, monospace;
        padding: 5px 11px;
        font-size: 13px;
        border-radius: 999px;
        border: 1px solid rgba(65, 214, 255, 0.18);
      }
    }
  }

  // 主体
  .hmi-body {
    position: relative;
    z-index: 1;
    display: grid;
    grid-template-columns: 280px minmax(680px, 1fr) 280px;
    gap: 10px;
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
    border-color: rgba(55, 213, 255, 0.52);
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
    background: rgba(44, 59, 70, 0.94);
    color: #f4f8fb;
    font-size: 13px;
    font-weight: 700;
    padding: 8px 12px;
    letter-spacing: 1px;
    border-bottom: 1px solid rgba(43, 197, 237, 0.14);
    display: flex;
    align-items: center;
    justify-content: space-between;
    &::before {
      content: '';
      width: 7px;
      height: 7px;
      margin-right: 8px;
      border-radius: 50%;
      background: #37d5ff;
      box-shadow: 0 0 12px rgba(55, 213, 255, 0.86);
    }
    .block-sub {
      font-size: 11px;
      font-weight: 400;
      color: rgba(255, 255, 255, 0.85);
    }
  }

  :deep(.n-base-selection) {
    --n-border: 1px solid rgba(76, 177, 216, 0.28) !important;
    --n-border-active: 1px solid rgba(55, 213, 255, 0.78) !important;
    --n-border-focus: 1px solid rgba(55, 213, 255, 0.78) !important;
    --n-box-shadow-active: 0 0 0 2px rgba(55, 213, 255, 0.16) !important;
    --n-box-shadow-focus: 0 0 0 2px rgba(55, 213, 255, 0.16) !important;
    --n-color: rgba(12, 24, 34, 0.62) !important;
    --n-text-color: #eaf8ff !important;
    --n-placeholder-color: rgba(234, 248, 255, 0.72) !important;
    backdrop-filter: blur(8px);
  }
  :deep(.n-base-selection-label) {
    background: linear-gradient(180deg, rgba(31, 64, 82, 0.92), rgba(15, 30, 42, 0.9)) !important;
  }

  // 仪表盘
  .gauges {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
    padding: 12px;
    background: linear-gradient(90deg, rgba(55, 213, 255, 0.08) 1px, transparent 1px),
      linear-gradient(180deg, rgba(55, 213, 255, 0.08) 1px, transparent 1px);
    background-size: 18px 18px;
  }
  .gauge {
    position: relative;
    height: 154px;
    border-radius: 16px;
    background: radial-gradient(circle at 50% 46%, rgba(55, 213, 255, 0.18), transparent 38%),
      linear-gradient(180deg, rgba(255, 255, 255, 0.74), rgba(218, 234, 244, 0.42));
    border: 1px solid rgba(55, 145, 185, 0.18);
    box-shadow:
      inset 0 0 26px rgba(55, 213, 255, 0.08),
      0 10px 24px rgba(36, 72, 94, 0.08);
    overflow: hidden;
  }
  .gauge-wide {
    grid-column: 1 / -1;
    height: 128px;
  }
  .gauge::before,
  .gauge::after {
    content: '';
    position: absolute;
    inset: 10px;
    border-radius: 50%;
    border: 1px dashed rgba(55, 213, 255, 0.22);
    pointer-events: none;
  }
  .gauge::after {
    inset: auto 18px 12px;
    height: 1px;
    border: 0;
    border-radius: 0;
    background: linear-gradient(90deg, transparent, rgba(55, 213, 255, 0.72), transparent);
    animation: scanLine 2.4s ease-in-out infinite;
  }

  @keyframes scanLine {
    0%,
    100% {
      opacity: 0.25;
      transform: translateY(0);
    }
    50% {
      opacity: 1;
      transform: translateY(-118px);
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
      background: #37d5ff;
      box-shadow: 0 0 8px rgba(55, 213, 255, 0.9);
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
    padding: 0;
    background: #eef3f7;
    overflow: hidden;
    border-radius: 0 0 16px 16px;
  }
  .machine-stage {
    position: relative;
    min-height: 460px;
    height: 100%;
    container-type: inline-size;
    overflow: hidden;
    isolation: isolate;
    background: radial-gradient(circle at 42% 20%, rgba(255, 255, 255, 0.8), transparent 34%),
      linear-gradient(180deg, #e9eef4 0%, #f8fbff 52%, #dce7f0 100%);
  }
  .machine-stage::after {
    content: '';
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: linear-gradient(
      90deg,
      rgba(235, 243, 250, 0.18),
      transparent 18%,
      transparent 72%,
      rgba(235, 243, 250, 0.35)
    );
    z-index: 4;
  }
  .machine-model {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    z-index: 1;
    cursor: grab;
    outline: none;
  }
  .machine-model:active {
    cursor: grabbing;
  }
  .readout-panel {
    position: absolute;
    right: 12px;
    top: 118px;
    z-index: 3;
    width: clamp(174px, 21cqw, 214px);
    padding: 12px;
    border-radius: 16px;
    background: #303b45;
    border: 1px solid rgba(122, 196, 230, 0.28);
    box-shadow:
      0 20px 46px rgba(20, 34, 45, 0.34),
      inset 0 1px 0 rgba(255, 255, 255, 0.08);
  }
  .readout-panel::after {
    content: '';
    position: absolute;
    inset: -8px -10px -8px -10px;
    z-index: -1;
    border-radius: 18px;
    background: #303b45;
    box-shadow: 0 20px 46px rgba(20, 34, 45, 0.34);
  }
  .readout-panel::before {
    content: '运行数据';
    display: block;
    margin-bottom: 8px;
    color: #ff9800;
    font-size: 15px;
    font-weight: 800;
    text-align: center;
    letter-spacing: 2px;
  }
  .readout-card {
    position: relative;
    display: flex;
    align-items: baseline;
    min-height: 40px;
    padding: 0 2px;
    border-radius: 0;
    color: #f8fbff;
    border: 0;
    border-bottom: 1px solid rgba(112, 143, 160, 0.22);
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
    color: rgba(223, 233, 241, 0.72);
    font-size: 12px;
    font-weight: 700;
  }
  .readout-card strong {
    color: #f4f8ff;
    font-family: Menlo, Consolas, monospace;
    font-size: clamp(15px, 1.9cqw, 18px);
    letter-spacing: 0.2px;
    margin-right: 4px;
    text-align: right;
  }
  .readout-card small {
    color: rgba(233, 241, 247, 0.68);
    font-size: 11px;
  }
  .readout-card.val-power strong {
    color: #4de18c;
  }
  .readout-card.val-high strong {
    color: #ff6b6b;
  }
  .readout-card.val-low strong {
    color: #ffd166;
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
  .model-rotate-btn {
    position: absolute;
    left: 18px;
    bottom: 20px;
    z-index: 4;
    height: 30px;
    padding: 0 13px;
    color: #d7f7ff;
    font-size: 12px;
    font-weight: 700;
    border: 1px solid rgba(77, 229, 255, 0.3);
    border-radius: 999px;
    background: rgba(5, 15, 23, 0.72);
    box-shadow:
      inset 0 0 18px rgba(77, 229, 255, 0.08),
      0 12px 24px rgba(0, 0, 0, 0.18);
    cursor: pointer;
  }
  .model-rotate-btn:hover {
    color: #ffffff;
    border-color: rgba(77, 229, 255, 0.72);
    background: rgba(77, 229, 255, 0.16);
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
    background: linear-gradient(90deg, rgba(55, 213, 255, 0.05) 1px, transparent 1px),
      linear-gradient(180deg, rgba(55, 213, 255, 0.05) 1px, transparent 1px),
      radial-gradient(circle at 50% 20%, rgba(55, 213, 255, 0.12), transparent 48%);
    background-size:
      18px 18px,
      18px 18px,
      auto;
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

  // 科技数据仓主题覆盖
  .hmi-screen {
    --warehouse-bg: #071018;
    --warehouse-panel: rgba(9, 23, 34, 0.82);
    --warehouse-panel-strong: rgba(12, 31, 45, 0.94);
    --warehouse-border: rgba(77, 223, 255, 0.22);
    --warehouse-cyan: #4de5ff;
    --warehouse-blue: #3b82f6;
    --warehouse-green: #48f5a5;
    --warehouse-amber: #ffc857;
    --warehouse-red: #ff5c7a;

    min-height: calc(100vh - 80px);
    padding: 14px;
    color: #d7f7ff;
    font-family: 'DIN Alternate', 'Bahnschrift', 'Microsoft YaHei', sans-serif;
    background: linear-gradient(90deg, rgba(77, 229, 255, 0.08) 1px, transparent 1px),
      linear-gradient(180deg, rgba(77, 229, 255, 0.08) 1px, transparent 1px),
      radial-gradient(circle at 12% 12%, rgba(67, 190, 255, 0.26), transparent 24%),
      radial-gradient(circle at 82% 8%, rgba(72, 245, 165, 0.18), transparent 26%),
      radial-gradient(circle at 50% 110%, rgba(14, 116, 144, 0.36), transparent 42%),
      linear-gradient(135deg, #050b12 0%, #071018 42%, #0d1821 100%);
    background-size:
      34px 34px,
      34px 34px,
      auto,
      auto,
      auto,
      auto;
  }
  .hmi-screen::before {
    background: repeating-linear-gradient(
        90deg,
        transparent 0 23px,
        rgba(77, 229, 255, 0.035) 24px 25px
      ),
      linear-gradient(115deg, transparent 0 36%, rgba(77, 229, 255, 0.12) 46%, transparent 56% 100%),
      radial-gradient(circle at 50% 0, rgba(96, 239, 255, 0.16), transparent 34%);
    mix-blend-mode: screen;
    animation: warehouseSweep 9s linear infinite;
  }
  .hmi-screen::after {
    inset: 10px;
    border-color: rgba(77, 229, 255, 0.2);
    border-radius: 22px;
    box-shadow:
      inset 0 0 80px rgba(77, 229, 255, 0.08),
      0 0 0 1px rgba(5, 10, 16, 0.75);
  }
  .warehouse-orbit {
    position: absolute;
    pointer-events: none;
    border: 1px solid rgba(77, 229, 255, 0.16);
    border-radius: 50%;
    filter: drop-shadow(0 0 22px rgba(77, 229, 255, 0.15));
    z-index: 0;
  }
  .orbit-a {
    width: 420px;
    height: 420px;
    left: -150px;
    top: 80px;
    background: conic-gradient(from 45deg, transparent, rgba(77, 229, 255, 0.16), transparent 38%);
    animation: orbitRotate 24s linear infinite;
  }
  .orbit-b {
    width: 520px;
    height: 520px;
    right: -210px;
    bottom: -160px;
    background: conic-gradient(from 190deg, transparent, rgba(72, 245, 165, 0.12), transparent 42%);
    animation: orbitRotate 32s linear reverse infinite;
  }

  .hmi-top {
    height: 66px;
    background: linear-gradient(
        90deg,
        rgba(77, 229, 255, 0.16),
        transparent 18% 82%,
        rgba(72, 245, 165, 0.12)
      ),
      rgba(6, 16, 25, 0.92);
    border-color: rgba(77, 229, 255, 0.26);
    border-radius: 18px;
    box-shadow:
      0 22px 50px rgba(0, 0, 0, 0.28),
      inset 0 1px 0 rgba(151, 236, 255, 0.18),
      inset 0 -1px 0 rgba(77, 229, 255, 0.12);
  }
  .hmi-top::before {
    content: '';
    position: absolute;
    left: 72px;
    right: 72px;
    bottom: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--warehouse-cyan), transparent);
  }
  .hmi-top .top-deco-left {
    width: 70px;
    background: linear-gradient(135deg, rgba(77, 229, 255, 0.95), rgba(28, 94, 143, 0.7)),
      repeating-linear-gradient(90deg, transparent 0 7px, rgba(255, 255, 255, 0.28) 8px 9px);
  }
  .hmi-top .top-title {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2px;
    letter-spacing: 5px;
  }
  .hmi-top .top-title strong {
    color: #eaffff;
    font-size: 24px;
    line-height: 1;
    text-shadow:
      0 0 18px rgba(77, 229, 255, 0.58),
      0 0 34px rgba(72, 245, 165, 0.18);
  }
  .title-kicker {
    color: rgba(77, 229, 255, 0.76);
    font-family: 'Agency FB', 'Bahnschrift', sans-serif;
    font-size: 10px;
    letter-spacing: 6px;
  }
  .hmi-top .top-deco-right .clock {
    color: var(--warehouse-green);
    background: rgba(4, 12, 18, 0.9);
    border-color: rgba(72, 245, 165, 0.28);
    box-shadow:
      inset 0 0 18px rgba(72, 245, 165, 0.08),
      0 0 18px rgba(72, 245, 165, 0.08);
  }
  .fullscreen-btn {
    min-width: 74px;
    color: #d7f7ff;
    border-color: rgba(77, 229, 255, 0.34);
    background: rgba(5, 15, 23, 0.62);
    box-shadow: inset 0 0 18px rgba(77, 229, 255, 0.08);
  }
  .fullscreen-btn:hover {
    color: #ffffff;
    border-color: rgba(77, 229, 255, 0.78);
    background: rgba(77, 229, 255, 0.14);
  }
  :deep(.n-base-selection-label) {
    background: linear-gradient(180deg, rgba(10, 35, 50, 0.96), rgba(4, 13, 21, 0.96)) !important;
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
    border: 1px solid var(--warehouse-border);
    border-radius: 16px;
    background: linear-gradient(135deg, rgba(77, 229, 255, 0.1), transparent 36%),
      rgba(7, 20, 31, 0.78);
    box-shadow:
      inset 0 1px 0 rgba(151, 236, 255, 0.12),
      0 18px 36px rgba(0, 0, 0, 0.18);
  }
  .warehouse-chip::after {
    content: '';
    position: absolute;
    inset: auto 12px 10px 12px;
    height: 2px;
    background: linear-gradient(90deg, transparent, currentColor, transparent);
    opacity: 0.58;
  }
  .warehouse-chip::before {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(
      110deg,
      transparent 0 38%,
      rgba(255, 255, 255, 0.18) 48%,
      transparent 58% 100%
    );
    transform: translateX(-120%);
    animation: chipScan 3s ease-in-out infinite;
  }
  .warehouse-chip .chip-label {
    display: block;
    color: rgba(214, 247, 255, 0.6);
    font-size: 10px;
    letter-spacing: 2px;
  }
  .warehouse-chip strong {
    display: block;
    margin-top: 6px;
    color: currentColor;
    font-family: 'Agency FB', 'Bahnschrift', monospace;
    font-size: 26px;
    line-height: 1;
  }
  .chip-value {
    animation: chipRise 0.58s cubic-bezier(0.2, 0.86, 0.28, 1.08);
  }
  .warehouse-chip small {
    display: block;
    margin-top: 5px;
    color: rgba(214, 247, 255, 0.58);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .chip-desc {
    animation: chipRise 0.5s ease both;
  }
  .tone-cyan {
    color: var(--warehouse-cyan);
  }
  .tone-blue {
    color: #7bb5ff;
  }
  .tone-green {
    color: var(--warehouse-green);
  }
  .tone-amber {
    color: var(--warehouse-amber);
  }
  .tone-red {
    color: var(--warehouse-red);
  }

  @keyframes chipRise {
    0% {
      opacity: 0;
      transform: translateY(18px);
      filter: blur(4px);
    }
    100% {
      opacity: 1;
      transform: translateY(0);
      filter: blur(0);
    }
  }
  @keyframes chipScan {
    0% {
      transform: translateX(-120%);
    }
    45%,
    100% {
      transform: translateX(120%);
    }
  }

  .hmi-body {
    grid-template-columns: 300px minmax(720px, 1fr) 300px;
    gap: 12px;
  }
  .block {
    background: linear-gradient(145deg, rgba(77, 229, 255, 0.08), transparent 42%),
      var(--warehouse-panel);
    border-color: var(--warehouse-border);
    border-radius: 18px;
    box-shadow:
      inset 0 1px 0 rgba(151, 236, 255, 0.1),
      inset 0 -1px 0 rgba(77, 229, 255, 0.08),
      0 22px 48px rgba(0, 0, 0, 0.22);
  }
  .block::before,
  .block::after {
    width: 34px;
    height: 34px;
    border-color: rgba(77, 229, 255, 0.7);
  }
  .block-title {
    min-height: 40px;
    padding: 10px 14px;
    background: linear-gradient(90deg, rgba(77, 229, 255, 0.14), transparent 62%),
      rgba(5, 14, 22, 0.94);
    color: #e6fbff;
    border-bottom-color: rgba(77, 229, 255, 0.18);
    font-size: 13px;
    letter-spacing: 2px;
  }
  .block-title::before {
    background: var(--warehouse-cyan);
    box-shadow: 0 0 16px rgba(77, 229, 255, 0.9);
  }

  .gauges,
  .chart-area,
  .param-table,
  .status-list,
  .alarm-list {
    background: linear-gradient(90deg, rgba(77, 229, 255, 0.045) 1px, transparent 1px),
      linear-gradient(180deg, rgba(77, 229, 255, 0.045) 1px, transparent 1px);
    background-size: 22px 22px;
  }
  .gauge {
    background: radial-gradient(circle at 50% 44%, rgba(77, 229, 255, 0.2), transparent 42%),
      linear-gradient(180deg, rgba(15, 38, 54, 0.92), rgba(5, 16, 25, 0.86));
    border-color: rgba(77, 229, 255, 0.2);
    box-shadow:
      inset 0 0 42px rgba(77, 229, 255, 0.08),
      0 12px 26px rgba(0, 0, 0, 0.24);
  }
  .gauge::before {
    border-color: rgba(77, 229, 255, 0.32);
    background: radial-gradient(circle, rgba(4, 13, 21, 0.08), rgba(4, 13, 21, 0.34) 64%);
  }
  .gauge::after {
    background: linear-gradient(90deg, transparent, rgba(77, 229, 255, 0.9), transparent);
  }
  .param-table .param-row,
  .status-row,
  .alarm-row {
    color: #d7f7ff;
    background: linear-gradient(90deg, rgba(77, 229, 255, 0.08), rgba(255, 255, 255, 0.02));
    border-color: rgba(77, 229, 255, 0.12);
    border-left-color: rgba(77, 229, 255, 0.58);
  }
  .param-table .param-row:hover,
  .status-row:hover,
  .alarm-row:hover {
    background: linear-gradient(90deg, rgba(77, 229, 255, 0.16), rgba(72, 245, 165, 0.04));
    box-shadow:
      inset 0 0 0 1px rgba(77, 229, 255, 0.18),
      0 12px 24px rgba(0, 0, 0, 0.16);
  }
  .param-table .param-label,
  .status-row .status-name,
  .alarm-row .alarm-name {
    color: rgba(215, 247, 255, 0.76);
  }
  .param-table .param-value {
    color: var(--warehouse-cyan);
  }
  .param-table .param-value small,
  .status-empty,
  .alarm-empty {
    color: rgba(215, 247, 255, 0.42);
  }

  .alarm-banner {
    background: linear-gradient(90deg, rgba(255, 92, 122, 0.24), rgba(255, 200, 87, 0.16));
    color: #ffe8ed;
    border-color: rgba(255, 92, 122, 0.36);
    box-shadow: inset 0 0 20px rgba(255, 92, 122, 0.08);
  }
  .alarm-banner-ok {
    background: linear-gradient(90deg, rgba(72, 245, 165, 0.18), rgba(77, 229, 255, 0.08));
    color: #caffea;
    border-color: rgba(72, 245, 165, 0.28);
  }
  .device-canvas,
  .machine-stage {
    background: radial-gradient(circle at 45% 40%, rgba(77, 229, 255, 0.16), transparent 34%),
      linear-gradient(180deg, #07131e 0%, #0a1722 52%, #050c13 100%);
  }
  .machine-model {
    filter: drop-shadow(0 18px 38px rgba(0, 0, 0, 0.28));
  }
  .machine-stage::after {
    background: linear-gradient(
        90deg,
        rgba(77, 229, 255, 0.1),
        transparent 18%,
        transparent 72%,
        rgba(77, 229, 255, 0.12)
      ),
      radial-gradient(circle at 50% 50%, transparent 0 44%, rgba(4, 10, 16, 0.48) 88%);
  }
  .readout-panel,
  .readout-panel::after {
    background: linear-gradient(180deg, rgba(8, 22, 33, 0.96), rgba(5, 13, 20, 0.96));
    border-color: rgba(77, 229, 255, 0.28);
  }
  .readout-panel::before {
    color: var(--warehouse-amber);
  }
  .readout-card {
    color: #eaffff;
    border-bottom-color: rgba(77, 229, 255, 0.12);
  }
  .readout-label {
    color: rgba(215, 247, 255, 0.58);
  }
  .readout-card strong {
    color: #eaffff;
  }
  .live-badge {
    color: #caffea;
    background: rgba(5, 15, 23, 0.76);
    border-color: rgba(72, 245, 165, 0.26);
  }

  .status-row .led-on {
    background: var(--warehouse-green);
    box-shadow: 0 0 14px rgba(72, 245, 165, 0.9);
  }
  .status-row .led-off,
  .alarm-off .alarm-led {
    background: rgba(126, 154, 169, 0.44);
  }
  .status-row .txt-on {
    color: var(--warehouse-green);
  }
  .status-row .txt-off,
  .alarm-off .alarm-name {
    color: rgba(215, 247, 255, 0.42);
  }
  .alarm-on {
    background: linear-gradient(90deg, rgba(255, 92, 122, 0.2), rgba(255, 200, 87, 0.08));
  }
  .alarm-on .alarm-led {
    background: var(--warehouse-red);
    box-shadow: 0 0 12px rgba(255, 92, 122, 0.9);
  }
  .alarm-on .alarm-name {
    color: #ffcad4;
  }
  .alarm-on .alarm-tag {
    background: var(--warehouse-red);
  }
  .alarm-off .alarm-tag {
    color: rgba(215, 247, 255, 0.52);
    background: rgba(77, 229, 255, 0.08);
  }

  @keyframes warehouseSweep {
    0% {
      transform: translateX(-3%);
    }
    100% {
      transform: translateX(3%);
    }
  }
  @keyframes orbitRotate {
    to {
      transform: rotate(360deg);
    }
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
