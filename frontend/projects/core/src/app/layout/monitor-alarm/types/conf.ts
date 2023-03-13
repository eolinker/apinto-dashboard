import { EO_NG_DROPDOWN_MENU_ITEM } from 'eo-ng-dropdown'
import { RadioOption } from 'eo-ng-radio'
import { SelectOption } from 'eo-ng-select/public-api'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { TabsOptions } from 'eo-ng-tabs'
import { MonitorProxyTableConfig, MonitorRequestTableConfig, SummaryPieData } from './types'

// 监控总览饼图数据初始值
export const initPie: SummaryPieData = {
  total: 0,
  success: 0,
  fail: 0,
  status4xx: 0,
  status5xx: 0
}

export const totalTabsList: TabsOptions[] = [
  { title: 'API请求量Top10', lazyLoad: true },
  { title: '应用调用量Top10', lazyLoad: true },
  { title: '上游服务调用量Top10', lazyLoad: true }
]

export const timeButtonOptions: RadioOption[] = [
  { label: '近1小时', value: 'hour' },
  { label: '近24小时', value: 'day' },
  { label: '近3天', value: 'threeDays' },
  { label: '近7天', value: 'sevenDays' }
]

export const warnDimensionOptions: SelectOption[] | RadioOption[] = [
  { label: '按API', value: 'api' },
  { label: '按上游', value: 'service' },
  { label: '按集群', value: 'cluster' },
  { label: '按分区', value: 'partition' }
]

export const targetTypeOptions: SelectOption[] = [
  { label: '不限', value: 'unlimited' },
  { label: '包含', value: 'contain' },
  { label: '不包含', value: 'not_contain' }
]

export const isEnableOptions: SelectOption[] = [
  { label: '已启用', value: 1 },
  { label: '已停用', value: 0 }
]

export const quoteOptions: SelectOption[] = [
  { label: '请求失败状态码数', value: 'request_fail_count' },
  { label: '请求失败率', value: 'request_fail_rate' },
  { label: '请求4xx状态码数', value: 'request_status_4xx' },
  { label: '请求5xx状态码数', value: 'request_status_5xx' },
  { label: '转发失败状态码数', value: 'proxy_fail_count' },
  { label: '转发失败率', value: 'proxy_fail_rate' },
  { label: '转发4xx状态码数', value: 'proxy_status_4xx' },
  { label: '转发5xx状态码数', value: 'proxy_status_5xx' },
  { label: '请求报文量', value: 'request_message' },
  { label: '响应报文量', value: 'response_message' },
  { label: '平均响应时间', value: 'avgResp' }
]

// 告警指标为成功率、失败率、同比的单位为百分比（小数点留两位），
// 告警指标为xx数时单位为次（整数），
// 告警指标为请求报文量、响应报文量时单位为KB，
// 平均响应时间的单位为ms
export const valueUnitMap: Map<string, string> = new Map([
  ['request_fail_rate', '%'],
  ['proxy_fail_rate', '%'],
  ['request_fail_count', '次'],
  ['request_status_4xx', '次'],
  ['request_status_5xx', '次'],
  ['proxy_fail_count', '次'],
  ['proxy_status_4xx', '次'],
  ['proxy_status_5xx', '次'],
  ['request_message', 'KB'],
  ['response_message', 'KB'],
  ['avgResp', 'ms']
])

export const valueUnit: { [key: string]: string } = {
  '%': '%',
  次: 'num',
  ms: 'ms',
  KB: 'kb'
}

export const rateList: Array<string> = [
  'ring_ratio_add',
  'ring_ratio_reduce',
  'year_basis_add',
  'year_basis_reduce'
]

// 监控表格参数
export const proxyBaseTableHeadName:THEAD_TYPE[] = [
  {
    title: '转发总数',
    width: 90,
    key: 'proxyTotal',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '转发成功数',
    width: 100,
    key: 'proxySuccess',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '转发成功率',
    width: 100,
    key: 'proxyRate',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '失败状态码数',
    width: 120,
    key: 'statusFail',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '平均响应时间(ms)',
    width: 148,
    key: 'avgResp',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '最大响应时间(ms)',
    width: 148,
    key: 'maxResp',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '最小响应时间(ms)',
    width: 148,
    key: 'minResp',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '平均请求流量(KB)',
    width: 148,
    key: 'avgTraffic',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '最大请求流量(KB)',
    width: 148,
    key: 'maxTraffic',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '最小请求流量(KB)',
    width: 148,
    key: 'minTraffic',
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '操作',
    right: true
  }]

export const apiProxyTableHeadName:THEAD_TYPE[] = [
  {
    title: 'API名称',
    left: true,
    width: 86,
    key: 'apiName'
  }, ...proxyBaseTableHeadName
]

export const proxyTableHeadName:THEAD_TYPE[] = [
  {
    title: '转发路径',
    left: true,
    width: 148,
    key: 'proxyPath'
  }, ...proxyBaseTableHeadName
]

export const nodeTableHeadName:THEAD_TYPE[] = [
  {
    title: '目标节点',
    key: 'addr',
    left: true,
    width: 80
  }, ...proxyBaseTableHeadName
]

export const serviceTableHeadName:THEAD_TYPE[] = [
  {
    title: '上游服务名称',
    key: 'serviceName',
    left: true,
    width: 120
  }, ...proxyBaseTableHeadName
]

export const requestBaseTableHeadName:THEAD_TYPE[] = [
  {
    title: '请求总数',
    key: 'requestTotal',
    showSort: true,
    sortOrder: null,
    sortPriority: false,
    width: 90
  },
  {
    title: '请求成功数',
    key: 'requestSuccess',
    width: 100,
    showSort: true,
    sortOrder: null,
    sortPriority: false
  },
  {
    title: '请求成功率',
    key: 'requestRate',
    showSort: true,
    sortOrder: null,
    sortPriority: false,
    width: 100
  },
  ...proxyBaseTableHeadName
]

export const apiTableHeadName:THEAD_TYPE[] = [

  {
    title: 'API名称',
    key: 'apiName',
    width: 86,
    left: true
  },
  {
    title: '请求路径',
    key: 'path',
    width: 80
  },
  ...requestBaseTableHeadName
]

export const appTableHeadName:THEAD_TYPE[] = [
  {
    title: '应用名称',
    width: 80,
    key: 'appName',
    left: true
  },
  {
    title: '应用ID',
    key: 'appId',
    width: 80
  },
  ...requestBaseTableHeadName
]

export const pathTableHeadName:THEAD_TYPE[] = [
  {
    title: '请求路径',
    key: 'path',
    width: 148,
    left: true
  },
  ...requestBaseTableHeadName
]

export const proxyBaseTableBody:TBODY_TYPE[] =
[
  {
    key: 'proxyTotal'
  },
  {
    key: 'proxySuccess'
  },
  {
    key: 'proxyRate',
    keySuffix: '%'
  },
  {
    key: 'statusFail'
  },
  {
    key: 'avgResp'
  },
  {
    key: 'maxResp'
  },
  {
    key: 'minResp'
  },
  {
    key: 'avgTraffic'
  },
  {
    key: 'maxTraffic'
  },
  {
    key: 'minTraffic'
  },
  {
    type: 'btn',
    right: true
  }
]

export const apiProxyTableBody:TBODY_TYPE[] = [
  {
    key: 'apiName',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  }, ...proxyBaseTableBody
]

export const proxyTableBody:TBODY_TYPE[] = [
  {
    key: 'proxyPath',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  }, ...proxyBaseTableBody
]

export const nodeTableBody:TBODY_TYPE[] = [
  {
    key: 'addr',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  }, ...proxyBaseTableBody
]

export const serviceTableBody:TBODY_TYPE[] = [
  {
    key: 'serviceName',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  }, ...proxyBaseTableBody
]

export const requestBaseTableBody:TBODY_TYPE[] = [

  {
    key: 'requestTotal'
  },
  {
    key: 'requestSuccess'
  },
  {
    key: 'requestRate',
    keySuffix: '%'
  }, ...proxyBaseTableBody
]

export const apiTableBody:TBODY_TYPE[] = [
  {
    key: 'apiName',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  },
  {
    key: 'path'
  }, ...requestBaseTableBody
]

export const appTableBody:TBODY_TYPE[] = [
  {
    key: 'appName',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  },
  {
    key: 'appId'
  }, ...requestBaseTableBody
]

export const pathTableBody:TBODY_TYPE[] = [
  {
    key: 'path',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  }, ...requestBaseTableBody
]

// api转发,目标节点,转发路径,上游
export const proxyInitConfig:MonitorProxyTableConfig = {
  proxyTotal: true,
  proxySuccess: true,
  proxyRate: true,
  statusFail: true,
  avgResp: true,
  maxResp: true,
  minResp: true,
  avgTraffic: true,
  maxTraffic: true,
  minTraffic: true
}

// 请求路径
export const requestInitConfig:MonitorRequestTableConfig = {
  requestTotal: true,
  requestSuccess: true,
  requestRate: true,
  ...proxyInitConfig
}

// api
export const apiInitConfig:MonitorRequestTableConfig = {
  path: true,
  ...requestInitConfig
}

// 应用
export const appInitConfig:MonitorRequestTableConfig = {
  appId: true,
  ...requestInitConfig
}

// api转发,目标节点,转发路径,上游
export const proxyBaseInitDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [
  {
    title: '转发总数',
    key: 'proxyTotal'
  },
  {
    title: '转发成功数',
    key: 'proxySuccess'
  },
  {
    title: '转发成功率',
    key: 'proxyRate'
  },
  {
    title: '失败状态码数',
    key: 'statusFail'
  },
  {
    title: '平均响应时间(ms)',
    key: 'avgResp'
  },
  {
    title: '最大响应时间(ms)',
    key: 'maxResp'
  },
  {
    title: '最小响应时间(ms)',
    key: 'minResp'
  },
  {
    title: '平均请求流量(KB)',
    key: 'avgTraffic'
  },
  {
    title: '最大请求流量(KB)',
    key: 'maxTraffic'
  },
  {
    title: '最小请求流量(KB)',
    key: 'minTraffic'
  }
]

// 请求路径
export const requestBaseInitDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [
  {
    title: '请求总数',
    key: 'requestTotal'
  },
  {
    title: '请求成功数',
    key: 'requestSuccess'
  },
  {
    title: '请求成功率',
    key: 'requestRate'
  },
  ...proxyBaseInitDropdownMenu
]

export const apiInitDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [
  {
    title: '请求路径',
    key: 'path'
  },
  ...requestBaseInitDropdownMenu
]

export const appInitDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [
  {
    title: '应用ID',
    key: 'appId'
  },
  ...requestBaseInitDropdownMenu
]

export const compareOptions: SelectOption[] = [
  { label: '>', value: '>' },
  { label: '>=', value: '>=' },
  { label: '<', value: '<' },
  { label: '<=', value: '<=' },
  { label: '==', value: '==' },
  { label: '!=', value: '!=' },
  { label: '环比上个统计时间增率', value: 'ring_ratio_add' },
  { label: '环比上个统计时间减率', value: 'ring_ratio_reduce' },
  { label: '同比昨天同时间增率', value: 'year_basis_add' },
  { label: '同比昨天同时间减率', value: 'year_basis_reduce' }
]

export const everyOptions: SelectOption[] = [
  { label: '统计粒度1分钟', value: 1 },
  { label: '统计粒度3分钟', value: 3 },
  { label: '统计粒度5分钟', value: 5 },
  { label: '统计粒度10分钟', value: 10 },
  { label: '统计粒度30分钟', value: 30 },
  { label: '统计粒度1小时', value: 60 }
]

export function getTime (
  timeButton: string,
  datePickerValue: Date[],
  init?: boolean
): { startTime: number; endTime: number } {
  const currentSecond = new Date().getTime() // 当前毫秒数时间戳
  let currentMin = currentSecond - (currentSecond % (60 * 1000)) // 当前分钟数时间戳
  let startMin = currentMin - 60 * 60 * 1000
  if (!init && timeButton) {
    switch (timeButton) {
      case 'hour': {
        startMin = currentMin - 60 * 60 * 1000
        break
      }
      case 'day': {
        startMin = currentMin - 24 * 60 * 60 * 1000
        break
      }
      case 'threeDays': {
        startMin =
          new Date(new Date().setHours(0, 0, 0, 0)).getTime() -
          2 * 24 * 60 * 60 * 1000
        break
      }
      case 'sevenDays': {
        startMin =
          new Date(new Date().setHours(0, 0, 0, 0)).getTime() -
          6 * 24 * 60 * 60 * 1000
        break
      }
    }
  } else if (datePickerValue.length === 2) {
    startMin = new Date(
      new Date(datePickerValue[0]).setHours(0, 0, 0, 0)
    ).getTime()
    currentMin = new Date(new Date(datePickerValue[1]).setHours(23, 59, 59, 0)).getTime()
  }

  return { startTime: startMin / 1000, endTime: currentMin / 1000 }
}

export function getTimeUnit (timeInterval: string): string {
  let timeUnit = ''
  // 相差秒数
  switch (timeInterval) {
    case '1m': {
      timeUnit = '每分钟'
      break
    }
    case '5m': {
      timeUnit = '每5分钟'
      break
    }
    case '1h': {
      timeUnit = '每小时'
      break
    }
    case '1d': {
      timeUnit = '每天'
      break
    }
    case '1w': {
      timeUnit = '每周'
      break
    }
  }
  return timeUnit
}

// 当数据超过10万时，保留两个小数点，单位为万，如123212，显示12.32万；
export function changeNumberUnit (value:number):string {
  if (value > 1000000000) {
    return (value / 100000000).toFixed(2) + '亿'
  } else if (value > 1000000) {
    return (value / 10000).toFixed(0) + '万'
  } else if (value > 10000) {
    return (value / 10000).toFixed(2) + '万'
  }
  return value + '次数'
}
