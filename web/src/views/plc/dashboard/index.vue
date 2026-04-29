<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="PLC 实时监控" />
    </div>
    <n-card :bordered="false" class="proCard" style="margin-bottom: 16px">
      <n-space align="center">
        <span>选择设备：</span>
        <n-select
          v-model:value="deviceId"
          :options="deviceOptions"
          placeholder="请选择设备"
          style="width:280px"
          @update:value="onDeviceChange"
        />
        <n-button type="primary" @click="loadOverview">刷新</n-button>
        <n-button :type="autoRefresh ? 'error' : 'success'" @click="toggleAuto">
          {{ autoRefresh ? '停止自动刷新' : '开启自动刷新 (3s)' }}
        </n-button>
        <n-tag v-if="lastUpdate" type="info" size="small">最后更新：{{ lastUpdate }}</n-tag>
        <n-tag v-if="device" type="success" size="small">{{ device.name }} ({{ device.host }})</n-tag>
      </n-space>
    </n-card>

    <n-grid :cols="4" :x-gap="12" :y-gap="12" v-if="points.length > 0" responsive="screen">
      <n-gi v-for="p in points" :key="p.field">
        <n-card
          :bordered="true"
          size="small"
          :style="cardStyle(rt[p.field]?.alarmType)"
        >
          <template #header>
            <span style="font-size: 13px; color: #888">{{ p.name }}</span>
            <n-tag v-if="rt[p.field]?.alarmType === 1" type="error" size="small" style="margin-left:8px">超上限</n-tag>
            <n-tag v-else-if="rt[p.field]?.alarmType === 2" type="warning" size="small" style="margin-left:8px">超下限</n-tag>
          </template>
          <div style="font-size: 30px; font-weight: 700; text-align: center; padding: 12px 0; color: #333">
            {{ rt[p.field]?.engValue ?? '—' }}
            <span style="font-size: 14px; font-weight: 400; color: #999">{{ p.unit }}</span>
          </div>
          <div style="font-size: 11px; color: #aaa; text-align: right">{{ p.field }}</div>
        </n-card>
      </n-gi>
    </n-grid>

    <n-empty v-else-if="!loading" :description="deviceId ? '该设备暂无数据点' : '请先选择设备'" style="margin-top: 80px" />
    <div v-if="loading && points.length === 0" style="text-align:center; margin-top: 80px">
      <n-spin size="large" />
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, onMounted, onUnmounted } from 'vue';
  import { Overview, DeviceList } from '@/api/plc';

  const deviceId = ref<number | null>(null);
  const deviceOptions = ref<{ label: string; value: number }[]>([]);
  const device = ref<any>(null);
  const points = ref<any[]>([]);
  const rt = reactive<Record<string, any>>({});
  const loading = ref(false);
  const lastUpdate = ref('');
  let timer: any = null;
  const autoRefresh = ref(false);

  onMounted(async () => {
    await loadDevices();
  });

  async function loadDevices() {
    const res: any = await DeviceList({ status: 1, page: 1, pageSize: 200 });
    const list = res?.list || [];
    deviceOptions.value = list.map((d: any) => ({
      label: `${d.name} (${d.host})`,
      value: d.id,
    }));
    if (deviceOptions.value.length > 0 && !deviceId.value) {
      deviceId.value = deviceOptions.value[0].value;
      await onDeviceChange(deviceId.value);
    }
  }

  async function onDeviceChange(id: number | null) {
    deviceId.value = id;
    points.value = [];
    device.value = null;
    Object.keys(rt).forEach((k) => delete rt[k]);
    if (!id) return;
    await loadOverview(true);
  }

  // 单接口拉 device + points + 实时值
  async function loadOverview(rebuildSkeleton = false) {
    if (!deviceId.value) return;
    try {
      loading.value = true;
      const res: any = await Overview({ deviceId: deviceId.value });
      if (!res) return;
      device.value = res.device;
      // 首次或切设备时重建骨架, 否则保留 points 数组引用避免重渲染
      if (rebuildSkeleton || points.value.length === 0) {
        points.value = res.points || [];
      } else {
        // 增量更新静态字段(name/unit可能变), 跳过新增/移除导致重排
        const map = new Map(points.value.map((p: any) => [p.field, p]));
        for (const np of res.points || []) {
          const exist = map.get(np.field);
          if (exist) {
            exist.name = np.name;
            exist.unit = np.unit;
          }
        }
      }
      // 更新实时值 map
      for (const p of res.points || []) {
        rt[p.field] = { engValue: p.engValue, alarmType: p.alarmType };
      }
      lastUpdate.value = new Date().toLocaleTimeString();
    } finally {
      loading.value = false;
    }
  }

  function toggleAuto() {
    autoRefresh.value = !autoRefresh.value;
    if (autoRefresh.value) {
      loadOverview();
      timer = setInterval(() => loadOverview(false), 3000);
    } else {
      clearInterval(timer);
      timer = null;
    }
  }

  function cardStyle(alarmType?: number) {
    if (alarmType === 1) return 'border: 1px solid #d03050; background: #fff5f5;';
    if (alarmType === 2) return 'border: 1px solid #f0a020; background: #fffbe6;';
    return '';
  }

  onUnmounted(() => {
    if (timer) clearInterval(timer);
  });
</script>
