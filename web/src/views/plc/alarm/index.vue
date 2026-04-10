<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="PLC 报警记录" />
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
        :scroll-x="1300"
        :resizeHeightOffset="-10000"
      />
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { h, ref, reactive } from 'vue';
  import { useMessage, NTag } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { defRangeShortcuts } from '@/utils/dateUtil';
  import { AlarmList, AlarmResolve } from '@/api/plc';

  const message = useMessage();
  const actionRef = ref();
  const searchFormRef = ref<any>({});

  const [register] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 80,
    schemas: [
      {
        field: 'isResolved',
        component: 'NSelect',
        label: '处理状态',
        defaultValue: null,
        componentProps: {
          placeholder: '请选择',
          options: [
            { label: '未处理', value: 2 },
            { label: '已处理', value: 1 },
          ],
        },
      },
      {
        field: 'triggeredAt',
        component: 'NDatePicker',
        label: '触发时间',
        componentProps: {
          type: 'datetimerange',
          clearable: true,
          shortcuts: defRangeShortcuts(),
        },
      },
    ],
  });

  const columns = [
    { title: 'ID', key: 'id', width: 80 },
    { title: '点位', key: 'pointName', width: 140 },
    {
      title: '触发值',
      key: 'engValue',
      width: 120,
      render(row) { return `${row.engValue} ${row.unit}`; },
    },
    {
      title: '报警类型',
      key: 'alarmType',
      width: 130,
      render(row) {
        const text = row.alarmType === 1
          ? `超上限 (>${row.alarmMax})`
          : `超下限 (<${row.alarmMin})`;
        return h(NTag, { type: row.alarmType === 1 ? 'error' : 'warning', size: 'small' }, { default: () => text });
      },
    },
    {
      title: '处理状态',
      key: 'isResolved',
      width: 100,
      render(row) {
        return h(NTag, { type: row.isResolved === 1 ? 'success' : 'error', size: 'small' }, { default: () => (row.isResolved === 1 ? '已处理' : '未处理') });
      },
    },
    { title: '触发时间', key: 'triggeredAt', width: 170 },
    {
      title: '处理时间',
      key: 'resolvedAt',
      width: 170,
      render(row) { return row.resolvedAt || '—'; },
    },
    { title: '处理备注', key: 'remark', width: 180, ellipsis: { tooltip: true } },
  ];

  const actionColumn = reactive({
    width: 120,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          {
            label: '标记已处理',
            ifShow: () => record.isResolved === 2,
            onClick: () => handleResolve(record),
            auth: ['/plc/alarm/resolve'],
          },
        ],
      });
    },
  });

  const loadDataTable = async (res) => {
    return await AlarmList({ ...searchFormRef.value?.formModel, ...res });
  };

  function reloadTable() { actionRef.value?.reload(); }

  function handleResolve(record) {
    AlarmResolve({ id: record.id, remark: '手动处理' }).then(() => {
      message.success('已标记处理');
      reloadTable();
    });
  }
</script>
