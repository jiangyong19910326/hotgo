import { h } from 'vue';
import { NTag, NButton, NEllipsis, NSpace, NPopover, NText } from 'naive-ui';
import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { renderOptionTag, renderPopoverMemberSumma, MemberSumma } from '@/utils';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  public id = 0;
  public code = '';
  public originalUrl = '';
  public title = '';
  public totalClicks = 0;
  public todayClicks = 0;
  public yesterdayClicks = 0;
  public weeklyClicks = 0;
  public monthlyClicks = 0;
  public expireAt = '';
  public status = 1;
  public createdBy = 0;
  public createdBySumma?: null | MemberSumma = null;
  public updatedBy = 0;
  public updatedBySumma?: null | MemberSumma = null;
  public createdAt = '';
  public updatedAt = '';
  public deletedAt = '';

  constructor(state?: Partial<State>) {
    if (state) {
      Object.assign(this, state);
    }
  }
}

export function newState(state: State | Record<string, any> | null): State {
  if (state !== null) {
    if (state instanceof State) {
      return cloneDeep(state);
    }
    return new State(state);
  }
  return new State();
}

// 表单验证规则
export const rules = {
  code: {
    required: true,
    trigger: ['blur', 'input'],
    type: 'string',
    message: '请输入短码',
  },
  originalUrl: {
    required: true,
    trigger: ['blur', 'input'],
    type: 'string',
    message: '请输入原始链接',
  },
};

// 表格搜索表单
export const schemas = ref<FormSchema[]>([
  {
    field: 'id',
    component: 'NInputNumber',
    label: '主键',
    componentProps: {
      placeholder: '请输入主键',
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
  },
  {
    field: 'status',
    component: 'NSelect',
    label: '状态',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择状态',
      options: dict.getOption('sys_normal_disable'),
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
  },
  {
    field: 'createdAt',
    component: 'NDatePicker',
    label: '创建时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
  },
]);

// 表格列
export const columns = [
  {
    title: '短码',
    key: 'code',
    align: 'left',
    width: 120,
    render(row: State) {
      return h(
        NTag,
        {
          type: 'info',
          size: 'small',
          style: 'cursor: pointer; font-family: monospace;',
          onClick: () => {
            navigator.clipboard?.writeText(row.code).then(() => {
              window.$message?.success('短码已复制');
            });
          },
        },
        { default: () => row.code }
      );
    },
  },
  {
    title: '目标链接',
    key: 'originalUrl',
    align: 'left',
    width: 280,
    render(row: State) {
      return h(NSpace, { vertical: true, size: 2 }, {
        default: () => [
          row.title
            ? h(NText, { depth: 1, style: 'font-size: 13px;' }, { default: () => row.title })
            : null,
          h(
            'a',
            {
              href: row.originalUrl,
              target: '_blank',
              rel: 'noopener noreferrer',
              style: 'color: #18a058; font-size: 12px; text-decoration: none; word-break: break-all;',
              onMouseover: (e: MouseEvent) => {
                (e.target as HTMLElement).style.textDecoration = 'underline';
              },
              onMouseout: (e: MouseEvent) => {
                (e.target as HTMLElement).style.textDecoration = 'none';
              },
            },
            row.originalUrl.length > 50 ? row.originalUrl.slice(0, 50) + '…' : row.originalUrl
          ),
        ].filter(Boolean),
      });
    },
  },
  {
    title: '统计',
    key: 'totalClicks',
    align: 'center',
    width: 100,
    render(row: State) {
      return h(
        NPopover,
        { trigger: 'hover', placement: 'top' },
        {
          trigger: () =>
            h(NButton, { text: true, type: 'primary' }, {
              default: () =>
                h('span', { style: 'font-size: 15px; font-weight: 600;' }, String(row.totalClicks)),
            }),
          default: () =>
            h('div', { style: 'min-width: 120px; line-height: 1.8;' }, [
              h('div', null, `今日：${row.todayClicks}`),
              h('div', null, `昨日：${row.yesterdayClicks}`),
              h('div', null, `本周：${row.weeklyClicks}`),
              h('div', null, `本月：${row.monthlyClicks}`),
            ]),
        }
      );
    },
  },
  {
    title: '过期时间',
    key: 'expireAt',
    align: 'left',
    width: 160,
    render(row: State) {
      return row.expireAt
        ? h('span', null, row.expireAt)
        : h(NText, { depth: 3, style: 'font-size: 12px;' }, { default: () => '永不过期' });
    },
  },
  {
    title: '状态',
    key: 'status',
    align: 'center',
    width: 80,
    render(row: State) {
      return renderOptionTag('sys_normal_disable', row.status);
    },
  },
  {
    title: '创建者',
    key: 'createdBy',
    align: 'left',
    width: 100,
    render(row: State) {
      return renderPopoverMemberSumma(row.createdBySumma);
    },
  },
  {
    title: '创建时间',
    key: 'createdAt',
    align: 'left',
    width: 160,
  },
];

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions(['sys_normal_disable']);
}
