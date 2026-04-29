<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      :title="formValue.id > 0 ? '编辑矿场 #' + formValue.id : '添加矿场'"
      :style="{ width: '480px' }"
    >
      <n-spin :show="loading" description="请稍候...">
        <n-form
          ref="formRef"
          :model="formValue"
          :rules="rules"
          label-placement="left"
          :label-width="90"
          class="py-4"
        >
          <n-form-item label="矿场名称" path="name">
            <n-input v-model:value="formValue.name" placeholder="请输入矿场名称" />
          </n-form-item>
          <n-form-item label="地理位置">
            <n-input v-model:value="formValue.location" placeholder="如：新疆维吾尔自治区" />
          </n-form-item>
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
  import { MineEdit } from '@/api/plc';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const formRef = ref();
  const showModal = ref(false);
  const loading = ref(false);

  const defaultForm = () => ({
    id: 0,
    name: '',
    location: '',
    remark: '',
    status: 1,
  });

  const formValue = ref(defaultForm());

  const rules = {
    name: { required: true, message: '请输入矿场名称', trigger: ['blur', 'input'] },
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
      await MineEdit(formValue.value);
      message.success('保存成功');
      closeModal();
      emit('reloadTable');
    } finally {
      loading.value = false;
    }
  }

  defineExpose({ openModal });
</script>
