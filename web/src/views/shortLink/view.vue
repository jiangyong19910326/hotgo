<template>
  <div>
    <n-drawer v-model:show="showModal" :width="dialogWidth">
      <n-drawer-content title="短链接详情" closable>
        <n-spin :show="loading" description="请稍候...">
          <n-descriptions label-placement="left" class="py-2" :column="1">
            <n-descriptions-item>
              <template #label>
                短码
              </template>
              {{ formValue.code }}
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                原始链接
              </template>
              <span v-html="formValue.originalUrl"></span>
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                标题
              </template>
              {{ formValue.title }}
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                总点击量
              </template>
              {{ formValue.totalClicks }}
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                今日点击（每日凌晨重置）
              </template>
              {{ formValue.todayClicks }}
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                昨日点击
              </template>
              {{ formValue.yesterdayClicks }}
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                本周点击
              </template>
              {{ formValue.weeklyClicks }}
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                本月点击
              </template>
              {{ formValue.monthlyClicks }}
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                过期时间，NULL 表示永不过期
              </template>
              {{ formValue.expireAt }}
            </n-descriptions-item>
            <n-descriptions-item label="状态：1正常 2禁用">
              <n-tag
            :type="dict.getType('sys_normal_disable', formValue.status)"
            size="small"
            class="min-left-space"
          >
                {{ dict.getLabel('sys_normal_disable', formValue.status) }}
              </n-tag>
            </n-descriptions-item>
          </n-descriptions>
        </n-spin>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>
<script lang="ts" setup>
  import { useMessage } from 'naive-ui';
  import { View } from '@/api/shortLink';
  import { State, newState } from './model';
  import { adaModalWidth } from '@/utils/hotgo';
  import { useDictStore } from '@/store/modules/dict';

  const message = useMessage();
  const dict = useDictStore();
  const loading = ref(false);
  const showModal = ref(false);
  const formValue = ref(newState(null));
  const dialogWidth = computed(() => {
    return adaModalWidth(580);
  });

  // 打开模态框
  function openModal(state: State) {
    showModal.value = true;
    loading.value = true;
    View({ id: state.id })
      .then((res) => {
        formValue.value = res;
      })
      .finally(() => {
        loading.value = false;
      });
  }

  defineExpose({
    openModal,
  });
</script>

<style lang="less" scoped></style>