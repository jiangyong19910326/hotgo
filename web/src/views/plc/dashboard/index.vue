<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="PLC 实时监控" />
    </div>
    <n-card :bordered="false" class="proCard" style="margin-bottom: 16px">
      <n-space align="center">
        <span>选择设备ID：</span>
        <n-input-number v-model:value="deviceId" :min="1" placeholder="设备ID" style="width:160px" />
        <n-button type="primary" @click="loadRealtime">刷新</n-button>
        <n-button :type="autoRefresh ? 'error' : 'success'" @click="toggleAuto">
          {{ autoRefresh ? '停止自动刷新' : '开启自动刷新 (3s)' }}
        </n-button>
        <n-tag v-if="lastUpdate" type="info" size="small">最后更新：{{ lastUpdate }}</n-tag>
      </n-space>
    </n-card>

    <n-grid :cols="4" :x-gap="12" :y-gap="12" v-if="points.length > 0" responsive="screen">
      <n-gi v-for="p in points" :key="p.field">
        <n-card
          :bordered="true"
          size="small"
          :style="cardStyle(p.alarmType)"
        >
          <template #header>
            <span style="font-size: 13px; color: #888">{{ p.name }}</span>
            <n-tag v-if="p.alarmType === 1" type="error" size="small" style="margin-left:8px">超上限</n-tag>
            <n-tag v-else-if="p.alarmType === 2" type="warning" size="small" style="margin-left:8px">超下限</n-tag>
          </template>
          <div style="font-size: 30px; font-weight: 700; text-align: center; padding: 12px 0; color: #333">
            {{ p.engValue }}
            <span style="font-size: 14px; font-weight: 400; color: #999">{{ p.unit }}</span>
          </div>
          <div style="font-size: 11px; color: #aaa; text-align: right">{{ p.field }}</div>
        </n-card>
      </n-gi>
    </n-grid>

    <n-empty v-else-if="!loading" description="暂无数据，请输入设备ID后点击刷新" style="margin-top: 80px" />
    <div v-if="loading" style="text-align:center; margin-top: 80px">
      <n-spin size="large" />
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onUnmounted } from 'vue';
  import { Realtime } from '@/api/plc';

  const deviceId = ref(1);
  const points = ref<any[]>([]);
  const loading = ref(false);
  const lastUpdate = ref('');
  let timer: any = null;
  const autoRefresh = ref(false);

  async function loadRealtime() {
    try {
      loading.value = true;
      const res = await Realtime({ deviceId: deviceId.value });
      points.value = res?.data?.points || [];
      lastUpdate.value = new Date().toLocaleTimeString();
    } finally {
      loading.value = false;
    }
  }

  function toggleAuto() {
    autoRefresh.value = !autoRefresh.value;
    if (autoRefresh.value) {
      loadRealtime();
      timer = setInterval(loadRealtime, 3000);
    } else {
      clearInterval(timer);
      timer = null;
    }
  }

  function cardStyle(alarmType: number) {
    if (alarmType === 1) return 'border: 1px solid #d03050; background: #fff5f5;';
    if (alarmType === 2) return 'border: 1px solid #f0a020; background: #fffbe6;';
    return '';
  }

  onUnmounted(() => {
    if (timer) clearInterval(timer);
  });
</script>
