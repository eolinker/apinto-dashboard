import { THEAD_TYPE } from 'eo-ng-table'

export const IntelligentPluginDefaultThead:THEAD_TYPE[] = [
  {
    title: '名称',
    left: true
  },
  { title: 'ID' },
  { title: '描述' },
  { title: '状态' },
  { title: '更新人' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const PublishTableHeadName:THEAD_TYPE[] = [
  { type: 'checkbox' },
  { title: '集群' },
  { title: '状态' },
  { title: '更新人' },
  { title: '更新时间' }
]
