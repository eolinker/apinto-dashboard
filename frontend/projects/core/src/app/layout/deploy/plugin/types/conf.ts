import { THEAD_TYPE } from 'eo-ng-table'

export const PluginsTableHeadName: THEAD_TYPE[] = [
  {
    type: 'sort',
    width: 40
  },
  { title: '插件名称' },
  { title: '扩展ID' },
  { title: '描述' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]
