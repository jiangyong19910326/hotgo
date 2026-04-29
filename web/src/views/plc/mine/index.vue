<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="矿场管理" />
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
        :scroll-x="1000"
        :resizeHeightOffset="-10000"
      >
        <template #tableTitle>
          <n-button type="primary" @click="addTable" class="min-left-space">
            <template #icon><n-icon><PlusOutlined /></n-icon></template>
            添加矿场
          </n-button>
        </template>
      </BasicTable>
    </n-card>
    <Edit ref="editRef" @reloadTable="reloadTable" />
  </div>
</template>

<script lang="ts" setup>
  import { h, ref, reactive } from 'vue';
  import { useDialog, useMessage, NTag } from 'naive-ui';
  import { PlusOutlined } from '@vicons/antd';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { MineList, MineDelete, MineStatus } from '@/api/plc';
  import Edit from './edit.vue';

  const { } = {};
  const message = useMessage();
  const dialog = useDialog();
  const editRef = ref();
  const actionRef = ref();
  const searchFormRef = ref<any>({});

  const [register] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 80,
    schemas: [
      {
        field: 'name',
        component: 'NInput',
        label: '矿场名称',
        componentProps: { placeholder: '请输入名称' },
      },
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

  const columns = [
    { title: 'ID', key: 'id', width: 70 },
    { title: '矿场名称', key: 'name', width: 160 },
    { title: '地理位置', key: 'location', width: 200, ellipsis: { tooltip: true } },
    {
      title: '状态',
      key: 'status',
      width: 80,
      render(row) {
        return h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => (row.status === 1 ? '启用' : '禁用') });
      },
    },
    { title: '备注', key: 'remark', width: 200, ellipsis: { tooltip: true } },
    { title: '创建时间', key: 'createdAt', width: 160 },
  ];

  const actionColumn = reactive({
    width: 220,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          {
            label: '编辑',
            onClick: () => editRef.value.openModal(record),
          },
          {
            label: record.status === 1 ? '禁用' : '启用',
            type: record.status === 1 ? 'error' : 'success',
            onClick: () => handleStatus(record),
          },
          {
            label: '删除',
            type: 'error',
            onClick: () => handleDelete(record),
          },
        ],
      });
    },
  });

  const loadDataTable = async (res) => {
    return await MineList({ ...searchFormRef.value?.formModel, ...res });
  };

  function reloadTable() {
    actionRef.value?.reload();
  }

  function addTable() {
    editRef.value.openModal(null);
  }

  function handleDelete(record) {
    dialog.warning({
      title: '警告',
      content: `确认删除矿场【${record.name}】？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        MineDelete({ id: record.id }).then(() => {
          message.success('删除成功');
          reloadTable();
        });
      },
    });
  }

  function handleStatus(record) {
    MineStatus({ id: record.id, status: record.status === 1 ? 2 : 1 }).then(() => {
      message.success('状态已更新');
      reloadTable();
    });
  }
</script>
