<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      :title="formValue.id > 0 ? '编辑应用 #' + formValue.id : '添加应用'"
      :style="{ width: '600px' }"
    >
      <n-spin :show="loading" description="请稍候...">
        <n-form
          ref="formRef"
          :model="formValue"
          :rules="rules"
          label-placement="left"
          :label-width="110"
          class="py-4"
        >
          <n-form-item label="AppID" path="appId">
            <n-input v-model:value="formValue.appId" placeholder="如 frontend_v1（唯一标识，调用方传此值）" />
          </n-form-item>
          <n-form-item label="AppSecret" path="appSecret">
            <n-input-group>
              <n-input v-model:value="formValue.appSecret" placeholder="HMAC-SHA256 签名密钥（建议 32+ 字节随机串）" />
              <n-button type="primary" @click="genSecret">生成</n-button>
              <n-button @click="copySecret" v-if="formValue.appSecret">复制</n-button>
            </n-input-group>
          </n-form-item>
          <n-form-item label="应用名称" path="name">
            <n-input v-model:value="formValue.name" placeholder="如 前端看板" />
          </n-form-item>
          <n-form-item label="状态">
            <n-radio-group v-model:value="formValue.status">
              <n-radio :value="1">启用</n-radio>
              <n-radio :value="2">禁用</n-radio>
            </n-radio-group>
          </n-form-item>
          <n-form-item label="备注">
            <n-input v-model:value="formValue.remark" type="textarea" :rows="2" placeholder="用途说明" />
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
  import { AppEdit, AppGenSecret } from '@/api/plc';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const formRef = ref();
  const showModal = ref(false);
  const loading = ref(false);

  const defaultForm = () => ({
    id: 0,
    appId: '',
    appSecret: '',
    name: '',
    remark: '',
    status: 1,
  });

  const formValue = ref<any>(defaultForm());

  const rules = {
    appId: { required: true, message: '请输入 AppID', trigger: ['blur', 'input'] },
    appSecret: { required: true, message: '请生成或输入 AppSecret', trigger: ['blur', 'input'] },
    name: { required: true, message: '请输入应用名称', trigger: ['blur', 'input'] },
  };

  function openModal(row: any) {
    showModal.value = true;
    if (row) {
      formValue.value = { ...row };
    } else {
      formValue.value = defaultForm();
    }
  }

  function closeModal() { showModal.value = false; }

  async function genSecret() {
    const res: any = await AppGenSecret();
    formValue.value.appSecret = res?.appSecret || '';
    message.success('已生成新的 AppSecret');
  }

  async function copySecret() {
    try {
      await navigator.clipboard.writeText(formValue.value.appSecret);
      message.success('已复制');
    } catch {
      message.error('复制失败');
    }
  }

  async function handleSubmit() {
    try {
      await formRef.value?.validate();
      loading.value = true;
      await AppEdit(formValue.value);
      message.success('保存成功');
      closeModal();
      emit('reloadTable');
    } finally {
      loading.value = false;
    }
  }

  defineExpose({ openModal });
</script>
