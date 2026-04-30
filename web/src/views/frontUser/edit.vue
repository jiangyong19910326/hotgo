<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      :title="formValue.id > 0 ? '编辑用户 #' + formValue.id : '添加用户'"
      :style="{ width: '520px' }"
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
          <n-form-item label="用户名" path="username">
            <n-input v-model:value="formValue.username" placeholder="登录账号" :disabled="formValue.id > 0" />
          </n-form-item>
          <n-form-item label="密码" path="password">
            <n-input
              v-model:value="formValue.password"
              type="password"
              show-password-on="click"
              :placeholder="formValue.id > 0 ? '留空则保持原密码' : '请输入初始密码'"
            />
          </n-form-item>
          <n-form-item label="昵称">
            <n-input v-model:value="formValue.nickname" placeholder="昵称" />
          </n-form-item>
          <n-form-item label="手机号">
            <n-input v-model:value="formValue.mobile" placeholder="手机号" />
          </n-form-item>
          <n-form-item label="邮箱">
            <n-input v-model:value="formValue.email" placeholder="邮箱" />
          </n-form-item>
          <n-form-item label="头像">
            <n-input v-model:value="formValue.avatar" placeholder="头像 URL" />
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
  import { ref, computed } from 'vue';
  import { useMessage } from 'naive-ui';
  import { FrontUserEdit } from '@/api/frontUser';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const formRef = ref();
  const showModal = ref(false);
  const loading = ref(false);

  const defaultForm = () => ({
    id: 0,
    username: '',
    password: '',
    nickname: '',
    mobile: '',
    email: '',
    avatar: '',
    remark: '',
    status: 1,
  });

  const formValue = ref(defaultForm());

  const rules = computed(() => ({
    username: { required: true, message: '请输入用户名', trigger: ['blur', 'input'] },
    password: {
      required: formValue.value.id === 0,
      message: '请输入初始密码',
      trigger: ['blur', 'input'],
    },
  }));

  function openModal(row: any) {
    showModal.value = true;
    if (row) {
      formValue.value = { ...defaultForm(), ...row, password: '' };
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
      await FrontUserEdit(formValue.value);
      message.success('保存成功');
      closeModal();
      emit('reloadTable');
    } finally {
      loading.value = false;
    }
  }

  defineExpose({ openModal });
</script>
