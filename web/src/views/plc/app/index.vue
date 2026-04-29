<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="应用密钥管理" />
    </div>
    <n-card :bordered="false" class="proCard">
      <BasicForm
        ref="searchFormRef"
        @register="register"
        @submit="reloadTable"
        @reset="reloadTable"
        @keyup.enter="reloadTable"
      />
      <BasicTable
        ref="actionRef"
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row) => row.id"
        :actionColumn="actionColumn"
        :scroll-x="1200"
        :resizeHeightOffset="-10000"
      >
        <template #tableTitle>
          <n-button type="primary" @click="addRow" class="min-left-space">
            <template #icon><n-icon><PlusOutlined /></n-icon></template>
            添加应用
          </n-button>
        </template>
      </BasicTable>
    </n-card>
    <Edit ref="editRef" @reloadTable="reloadTable" />
  </div>
</template>

<script lang="ts" setup>
  import { h, ref, reactive } from 'vue';
  import { useDialog, useMessage, NTag, NButton } from 'naive-ui';
  import { PlusOutlined, CopyOutlined } from '@vicons/antd';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { AppList, AppDelete, AppStatus } from '@/api/plc';
  import Edit from './edit.vue';

  const message = useMessage();
  const dialog = useDialog();
  const editRef = ref();
  const actionRef = ref();
  const searchFormRef = ref<any>({});

  const [register] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 80,
    schemas: [
      { field: 'appId', component: 'NInput', label: 'AppID', componentProps: { placeholder: '请输入 AppID' } },
      { field: 'name', component: 'NInput', label: '应用名称', componentProps: { placeholder: '请输入名称' } },
      {
        field: 'status',
        component: 'NSelect',
        label: '状态',
        defaultValue: null,
        componentProps: {
          placeholder: '请选择状态',
          options: [
            { label: '启用', value: 1 },
            { label: '禁用', value: 2 },
          ],
        },
      },
    ],
  });

  function maskSecret(s: string) {
    if (!s) return '';
    if (s.length <= 8) return s;
    return s.slice(0, 4) + '****' + s.slice(-4);
  }

  async function copyText(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      message.success('已复制到剪贴板');
    } catch {
      message.error('复制失败，请手动复制');
    }
  }

  const columns = [
    { title: 'ID', key: 'id', width: 70 },
    { title: 'AppID', key: 'appId', width: 180 },
    {
      title: 'AppSecret',
      key: 'appSecret',
      width: 280,
      render(row: any) {
        return h('div', { style: 'display:flex;align-items:center;gap:8px' }, [
          h('code', { style: 'font-size:12px;color:#666' }, maskSecret(row.appSecret)),
          h(
            NButton,
            { size: 'tiny', quaternary: true, onClick: () => copyText(row.appSecret) },
            { default: () => '复制', icon: () => h(CopyOutlined as any) }
          ),
        ]);
      },
    },
    { title: '应用名称', key: 'name', width: 160 },
    {
      title: '状态',
      key: 'status',
      width: 80,
      render(row: any) {
        return h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => (row.status === 1 ? '启用' : '禁用') });
      },
    },
    { title: '备注', key: 'remark', width: 200, ellipsis: { tooltip: true } },
    { title: '创建时间', key: 'createdAt', width: 160 },
  ];

  const actionColumn = reactive({
    width: 200,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: any) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          { label: '编辑', onClick: () => editRef.value.openModal(record), auth: ['/plc/app/edit'] },
          {
            label: record.status === 1 ? '禁用' : '启用',
            type: record.status === 1 ? 'error' : 'success',
            onClick: () => handleStatus(record),
          },
          { label: '删除', type: 'error', onClick: () => handleDelete(record), auth: ['/plc/app/delete'] },
        ],
      });
    },
  });

  const loadDataTable = async (res: any) => {
    return await AppList({ ...searchFormRef.value?.formModel, ...res });
  };

  function reloadTable() { actionRef.value?.reload(); }
  function addRow() { editRef.value.openModal(null); }

  function handleDelete(record: any) {
    dialog.warning({
      title: '警告',
      content: `确认删除应用【${record.name}】？删除后该应用的签名将立即失效。`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        AppDelete({ id: record.id }).then(() => {
          message.success('删除成功');
          reloadTable();
        });
      },
    });
  }

  function handleStatus(record: any) {
    AppStatus({ id: record.id, status: record.status === 1 ? 2 : 1 }).then(() => {
      message.success('状态已更新');
      reloadTable();
    });
  }
</script>
