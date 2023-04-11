import { THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const PluginListStatusItems = [
  { label: '全部', value: '' },
  { label: '已启用', value: true },
  { label: '未启用', value: false }
]

export const PluginInstallConfigTableHeadName:THEAD_TYPE[] = [
  { title: '参数名' },
  { title: '参数值' },
  { title: '描述' }
]

export const PluginInstallConfigTableBody: EO_TBODY_TYPE[] = [
  {
    key: 'name'
  },
  {
    key: 'value',
    type: 'input',
    placeholderKey: 'placeholder'
  },
  { key: 'desc' }
]
