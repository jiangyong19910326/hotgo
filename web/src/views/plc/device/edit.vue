<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      :title="formValue.id > 0 ? '编辑设备 #' + formValue.id : '添加设备'"
      :style="{ width: '560px' }"
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
          <n-form-item label="设备名称" path="name">
            <n-input v-model:value="formValue.name" placeholder="请输入设备名称" />
          </n-form-item>
          <n-form-item label="IP 地址" path="host">
            <n-input v-model:value="formValue.host" placeholder="如 192.168.1.10" />
          </n-form-item>
          <n-grid :cols="2" :x-gap="12">
            <n-form-item-gi label="端口" path="port">
              <n-input-number v-model:value="formValue.port" :min="1" :max="65535" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="采集间隔(ms)" path="intervalMs">
              <n-input-number v-model:value="formValue.intervalMs" :min="100" :step="100" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="机架号(Rack)">
              <n-input-number v-model:value="formValue.rack" :min="0" style="width:100%" />
            </n-form-item-gi>
            <n-form-item-gi label="槽号(Slot)">
              <n-input-number v-model:value="formValue.slot" :min="0" style="width:100%" />
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
  import { DeviceEdit } from '@/api/plc';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const formRef = ref();
  const showModal = ref(false);
  const loading = ref(false);

  const defaultForm = () => ({
    id: 0,
    name: '',
    host: '',
    port: 102,
    rack: 0,
    slot: 1,
    intervalMs: 1000,
    remark: '',
    status: 1,
  });

  const formValue = ref(defaultForm());

  const rules = {
    name: { required: true, message: '请输入设备名称', trigger: ['blur', 'input'] },
    host: { required: true, message: '请输入 IP 地址', trigger: ['blur', 'input'] },
    port: { required: true, type: 'number', message: '请输入端口', trigger: ['blur'] },
    intervalMs: { required: true, type: 'number', message: '请输入采集间隔', trigger: ['blur'] },
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
      await DeviceEdit(formValue.value);
      message.success('保存成功');
      closeModal();
      emit('reloadTable');
    } finally {
      loading.value = false;
    }
  }

  defineExpose({ openModal });
</script>
