import { SelectOption } from 'eo-ng-select'
import { THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const UpstreamListTableHeadName: THEAD_TYPE[] = [
  { title: '上游名称' },
  { title: '协议类型' },
  { title: '服务类型' },
  { title: '地址' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const UpstreamListTableBody: EO_TBODY_TYPE[] = [
  {
    key: 'name',
    copy: true
  },
  {
    key: 'scheme',
    copy: true
  },
  {
    key: 'serviceType',
    copy: true
  },
  {
    key: 'config',
    copy: true
  },
  { key: 'updateTime' },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '查看'
      },
      {
        title: '删除'
      }
    ]
  }
]

export const UpstreamSchemeList:SelectOption[] = [
  { label: 'HTTP', value: 'HTTP' },
  { label: 'HTTPS', value: 'HTTPS' }
]

export const UpstreamBalanceList:SelectOption[] = [
  { label: 'round-robin', value: 'round-robin' }
]

export const ServicesTableHeadName:THEAD_TYPE[] = [
  { title: '服务名称' },
  { title: '服务类型' },
  { title: '描述' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const ServicesTablebody:EO_TBODY_TYPE[] = [
  {
    key: 'name',
    copy: true
  },
  {
    key: 'driver',
    copy: true
  },
  {
    key: 'desc'
  },
  { key: 'updateTime' },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '查看'
      },
      {
        title: '删除'
      }
    ]
  }

]
