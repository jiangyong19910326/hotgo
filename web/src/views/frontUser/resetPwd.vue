<template>
  <n-modal
    v-model:show="showModal"
    :mask-closable="false"
    :show-icon="false"
    preset="dialog"
    transform-origin="center"
    :title="`重置密码 - ${target?.username || ''}`"
    :style="{ width: '420px' }"
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
        <n-form-item label="新密码" path="password">
          <n-input
            v-model:value="formValue.password"
            type="password"
            show-password-on="click"
            placeholder="6~32 位"
          />
        </n-form-item>
      </n-form>
    </n-spin>
    <template #action>
      <n-space>
        <n-button @click="closeModal">取消</n-button>
        <n-button type="primary" :loading="loading" @click="handleSubmit">提交</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useMessage } from 'naive-ui';
  import { FrontUserResetPwd } from '@/api/frontUser';

  const message = useMessage();
  const formRef = ref();
  const showModal = ref(false);
  const loading = ref(false);
  const target = ref<any>(null);

  const formValue = ref({ password: '' });

  const rules = {
    password: {
      required: true,
      message: '请输入新密码（6~32 位）',
      trigger: ['blur', 'input'],
    },
  };

  function openModal(row: any) {
    target.value = row;
    formValue.value = { password: '' };
    showModal.value = true;
  }

  function closeModal() {
    showModal.value = false;
  }

  async function handleSubmit() {
    try {
      await formRef.value?.validate();
      if (formValue.value.password.length < 6 || formValue.value.password.length > 32) {
        message.error('密码长度需为 6~32 位');
        return;
      }
      loading.value = true;
      await FrontUserResetPwd({ id: target.value.id, password: formValue.value.password });
      message.success('密码已重置');
      closeModal();
    } finally {
      loading.value = false;
    }
  }

  defineExpose({ openModal });
</script>
