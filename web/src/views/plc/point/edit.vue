<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      :title="formValue.id > 0 ? '编辑数据点 #' + formValue.id : '添加数据点'"
      :style="{ width: '640px' }"
    >
      <n-spin :show="loading" description="请稍候...">
        <n-form
          ref="formRef"
          :model="formValue"
          :rules="rules"
          label-placement="left"
          :label-width="120"
          class="py-4"
        >
          <n-grid :cols="2" :x-gap="12">
            <n-form-item-gi label="所属设备ID" path="deviceId">
              <n-input-number v-model:value="formValue.deviceId" :min="1" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="点位名称" path="name">
              <n-input v-model:value="formValue.name" placeholder="如：润滑回油温度" />
            </n-form-item-gi>
            <n-form-item-gi label="MQTT 字段名" path="field">
              <n-input v-model:value="formValue.field" placeholder="如：LUBE_RETURN_TEMP" />
            </n-form-item-gi>
            <n-form-item-gi label="数据类型" path="dataType">
              <n-select v-model:value="formValue.dataType" :options="typeOptions" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="换算系数">
              <n-input-number v-model:value="formValue.scale" :step="0.01" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="换算偏移">
              <n-input-number v-model:value="formValue.offsetVal" :step="0.01" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="单位">
              <n-input v-model:value="formValue.unit" placeholder="如 ℃ bar rpm" />
            </n-form-item-gi>
            <n-form-item-gi label="排序">
              <n-input-number v-model:value="formValue.sort" :min="0" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="报警下限">
              <n-input-number v-model:value="formValue.alarmMin" clearable placeholder="不填则不报警" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="报警上限">
              <n-input-number v-model:value="formValue.alarmMax" clearable placeholder="不填则不报警" style="width:100%" />
            </n-form-item-gi>
          </n-grid>
          <n-form-item label="状态">
            <n-radio-group v-model:value="formValue.status">
              <n-radio :value="1">启用</n-radio>
              <n-radio :value="2">禁用</n-radio>
            </n-radio-group>
          </n-form-item>
          <n-form-item label="备注">
            <n-input v-model:value="formValue.remark" type="textarea" :rows="2" />
          </n-form-item>
        </n-form>
      </n-spin>
      <template #action>
        <n-space>
          <n-button @click="closeModal">取消</n-button>
          <n-button type="primary" :loading="loading" @click="handleSubmit">保存</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useMessage } from 'naive-ui';
  import { PointEdit } from '@/api/plc';

  const props = defineProps<{ deviceId?: number }>();
  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const formRef = ref();
  const showModal = ref(false);
  const loading = ref(false);

  const typeOptions = ['Bool', 'Real'].map((v) => ({ label: v, value: v }));

  const defaultForm = () => ({
    id: 0,
    deviceId: props.deviceId || 0,
    name: '',
    field: '',
    dataType: 'Real',
    scale: 1,
    offsetVal: 0,
    unit: '',
    alarmMin: null,
    alarmMax: null,
    remark: '',
    sort: 0,
    status: 1,
  });

  const formValue = ref(defaultForm());

  const rules = {
    deviceId: { required: true, type: 'number', message: '请输入设备ID', trigger: ['blur'] },
    name: { required: true, message: '请输入点位名称', trigger: ['blur', 'input'] },
    field: { required: true, message: '请输入 MQTT 字段名', trigger: ['blur', 'input'] },
    dataType: { required: true, message: '请选择数据类型', trigger: ['change'] },
  };

  function openModal(row: any) {
    showModal.value = true;
    if (row) {
      formValue.value = { ...row };
    } else {
      formValue.value = defaultForm();
    }
  }

  function closeModal() {
    showModal.value = false;
  }

  async function handleSubmit() {
    try {
      await formRef.value?.validate();
      loading.value = true;
      await PointEdit(formValue.value);
      message.success('保存成功');
      closeModal();
      emit('reloadTable');
    } finally {
      loading.value = false;
    }
  }

  defineExpose({ openModal });
</script>
