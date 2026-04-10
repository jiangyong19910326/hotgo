<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="数据点配置" />
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
        :scroll-x="1400"
        :resizeHeightOffset="-10000"
      >
        <template #tableTitle>
          <n-button type="primary" @click="addTable" class="min-left-space">
            <template #icon><n-icon><PlusOutlined /></n-icon></template>
            添加数据点
          </n-button>
        </template>
      </BasicTable>
    </n-card>
    <Edit ref="editRef" @reloadTable="reloadTable" :deviceId="currentDeviceId" />
  </div>
</template>

<script lang="ts" setup>
  import { h, ref, reactive } from 'vue';
  import { useRoute } from 'vue-router';
  import { useDialog, useMessage, NTag } from 'naive-ui';
  import { PlusOutlined } from '@vicons/antd';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { PointList, PointDelete, PointStatus } from '@/api/plc';
  import Edit from './edit.vue';

  const route = useRoute();
  const message = useMessage();
  const dialog = useDialog();
  const editRef = ref();
  const actionRef = ref();
  const searchFormRef = ref<any>({});
  const currentDeviceId = ref(Number(route.query.deviceId) || 0);

  const [register] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 80,
    schemas: [
      {
        field: 'deviceId',
        component: 'NInputNumber',
        label: '设备ID',
        defaultValue: currentDeviceId.value || undefined,
        componentProps: { placeholder: '请输入设备ID' },
      },
      {
        field: 'name',
        component: 'NInput',
        label: '点位名称',
        componentProps: { placeholder: '请输入名称' },
      },
    ],
  });

  const columns = [
    { title: 'ID', key: 'id', width: 70 },
    { title: '设备ID', key: 'deviceId', width: 80 },
    { title: '点位名称', key: 'name', width: 140 },
    { title: '字段标识', key: 'field', width: 120 },
    { title: '存储区', key: 'area', width: 60 },
    { title: 'DB块', key: 'dbNumber', width: 60 },
    { title: '字节偏移', key: 'byteOffset', width: 80 },
    { title: '数据类型', key: 'dataType', width: 90 },
    { title: '系数', key: 'scale', width: 70 },
    { title: '偏移', key: 'offsetVal', width: 70 },
    { title: '单位', key: 'unit', width: 60 },
    {
      title: '报警下限',
      key: 'alarmMin',
      width: 90,
      render(row) { return row.alarmMin ?? '—'; },
    },
    {
      title: '报警上限',
      key: 'alarmMax',
      width: 90,
      render(row) { return row.alarmMax ?? '—'; },
    },
    { title: '排序', key: 'sort', width: 60 },
    {
      title: '状态',
      key: 'status',
      width: 80,
      render(row) {
        return h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => (row.status === 1 ? '启用' : '禁用') });
      },
    },
  ];

  const actionColumn = reactive({
    width: 180,
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
            auth: ['/plc/point/edit'],
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
            auth: ['/plc/point/delete'],
          },
        ],
      });
    },
  });

  const loadDataTable = async (res) => {
    return await PointList({ ...searchFormRef.value?.formModel, ...res });
  };

  function reloadTable() { actionRef.value?.reload(); }
  function addTable() { editRef.value.openModal(null); }

  function handleDelete(record) {
    dialog.warning({
      title: '警告',
      content: `删除点位【${record.name}】？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        PointDelete({ id: record.id }).then(() => {
          message.success('删除成功');
          reloadTable();
        });
      },
    });
  }

  function handleStatus(record) {
    PointStatus({ id: record.id, status: record.status === 1 ? 2 : 1 }).then(() => {
      message.success('状态已更新');
      reloadTable();
    });
  }
</script>
