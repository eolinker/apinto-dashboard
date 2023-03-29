import { RadioOption } from 'eo-ng-radio'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const filterTableHeadName:THEAD_TYPE[] = [
  {
    title: '属性名称'
  },
  {
    title: '属性值'
  },
  {
    title: '操作',
    right: true
  }
]

export const filterTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'title',
    copy: true
  },
  {
    key: 'label'
  },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '配置'
      },
      {
        title: '删除'
      }
    ]
  }
]

export const nodesTableBody:TBODY_TYPE[] = [
  {
    key: 'node',
    type: 'input',
    placeholder: '请输入主机名或IP：端口',
    checkMode: 'change',
    errorTip: '请输入主机名或IP：端口'
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加',
        action: 'add'
      }
    ]
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加',
        action: 'add'
      },
      {
        title: '减少',
        action: 'delete'
      }
    ]
  }
]

export const metricsList:SelectOption[] = [
  { label: '上游服务', value: '{service}' },
  { label: 'API', value: '{api}' }
]

export const distributionOptions: RadioOption[] = [
  { label: '按百分比', value: 'percent' },
  { label: '按规则', value: 'match' }
]

export const visitRuleList:SelectOption[] = [
  { label: '允许', value: 'allow' },
  { label: '拒绝', value: 'refuse' }
]
// 策略列表表格参数
export const strategiesTableHeadName:THEAD_TYPE[] = [
  {
    title: '策略名称'
  },
  {
    title: '优先级',
    key: 'priority',
    width: 86,
    showSort: true
  },
  {
    title: '发布状态',
    width: 84
  },
  {
    title: '启停',
    width: 80
  },
  {
    title: '筛选条件'
  },
  {
    title: '限流维度'
  },
  {
    title: '更新者'
  },
  {
    title: '更新时间',
    showSort: true
  },
  {
    title: '操作',
    right: true
  }]

export const strategiesTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'name',
    copy: true
  },
  { key: 'priority' },
  { key: 'status' },
  { key: 'isStop' },
  { key: 'filters' },
  { key: 'conf' },
  { key: 'operator' },
  { key: 'updateTime' },
  {
    type: 'btn',
    showFn: (item:any) => {
      return item.isDeleted === false
    },
    right: true,
    btns: [
      {
        title: '查看'
      },
      {
        title: '删除'
      }
    ]
  },
  {
    type: 'btn',
    showFn: (item:any) => {
      return item.isDeleted === true
    },
    right: true,
    btns: [
      {
        title: '查看'
      },
      {
        title: '恢复'
      }
    ]
  }
]

export const publishTableHeadName:THEAD_TYPE[] = [
  {
    title: '策略名称'
  },
  {
    title: '优先级'
  },
  {
    title: '状态'
  },
  {
    title: '操作时间'
  }
]

export const publishTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'name',
    copy: true
  },
  {
    key: 'priority'
  },
  { key: 'status' },
  { key: 'optTime' }]

export const responseHeaderTableBody:TBODY_TYPE[] = [
  {
    key: 'key',
    type: 'input',
    placeholder: '请输入Key'
  },
  {
    key: 'value',
    type: 'input',
    placeholder: '请输入Value'
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加',
        action: 'add'
      }
    ]
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加',
        action: 'add'
      },
      {
        title: '减少',
        action: 'delete'
      }
    ]
  }
]
