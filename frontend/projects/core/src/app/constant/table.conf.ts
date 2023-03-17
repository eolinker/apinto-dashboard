import { MenuOptions } from 'eo-ng-menu'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'

// 监控告警列表参数
export const monitorAlarmStrategyTableHead:THEAD_TYPE[] = [
  { title: '策略名称', key: 'strategyTitle' },
  { title: '告警维度', key: 'warnDimension' },
  { title: '告警目标', key: 'warnTarget' },
  { title: '告警规则', key: 'warnRule' },
  { title: '告警消息频率', key: 'warnFrequency' },
  { title: '告警启停', key: 'isEnable' },
  { title: '更新者', key: 'operator' },
  { title: '更新时间', key: 'updateTime', showSort: true }
]
export const monitorAlarmStrategyTableBody:TBODY_TYPE[] = [
  { key: 'strategyTitle' },
  { key: 'warnDimension' },
  { key: 'warnTarget' },
  { key: 'warnRule' },
  { key: 'warnFrequency' },
  { key: 'isEnable' },
  { key: 'operator' },
  { key: 'updateTime' }
]

export const monitorAlarmStrategyTableDropdownMenu:MenuOptions[] = [
  { title: '策略名称', key: 'strategyTitle', checked: true },
  { title: '告警维度', key: 'warnDimension', checked: true },
  { title: '告警目标', key: 'warnTarget', checked: true },
  { title: '告警规则', key: 'warnRule', checked: true },
  { title: '告警消息频率', key: 'warnFrequency', checked: true },
  { title: '告警启停', key: 'isEnable', checked: true },
  { title: '更新者', key: 'operator', checked: true },
  { title: '更新时间', key: 'updateTime', checked: true }
]

export const monitorAlarmStrategyTableConfig:{[key:string]:boolean} = {
  strategyTitle: true,
  warnDimension: true,
  warnTarget: true,
  warnRule: true,
  warnFrequency: true,
  isEnable: true,
  operator: true,
  updateTime: true
}

// 监控告警历史列表参数
export const monitorAlarmHistoryTableHead:THEAD_TYPE[] = [
  { title: '策略名称' },
  { title: '告警目标' },
  { title: '告警内容' },
  { title: '告警状态' },
  { title: '告警时间', showSort: true }
]
export const monitorAlarmHistoryTableBody:TBODY_TYPE[] = [
  { key: 'strategyTitle' },
  { key: 'warnTarget' },
  { key: 'warnContent' },
  { key: 'status' },
  { key: 'createTime' }
]

// webhook列表参数
export const webhooksTableHead:THEAD_TYPE[] = [
  { title: 'webhook名称' },
  { title: '通知url' },
  { title: '请求方式' },
  { title: '参数类型' },
  { title: '更新者' },
  { title: '更新时间' }
]
export const webhooksTableBody:TBODY_TYPE[] = [
  { key: 'title' },
  { key: 'url' },
  { key: 'method' },
  { key: 'contentType' },
  { key: 'operator' },
  { key: 'updateTime' }
]

// webhook header列表
export const responseHeaderTableBody: any[] = [
  {
    key: 'key',
    type: 'input',
    placeholder: '请输入Key'

  },
  {
    key: 'value',
    type: 'input',
    placeholder: '请输入Value'
  }
]
