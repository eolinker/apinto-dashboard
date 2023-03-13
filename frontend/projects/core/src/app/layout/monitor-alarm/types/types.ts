import { THEAD_TYPE } from 'eo-ng-table'

// 监控分区
export interface MonitorPartition{
  uuid:string
  name:string
  clusterNames:Array<string>
}

// 监控总览中的饼图数据
export interface SummaryPieData{
  total:number, success:number, fail:number, status4xx:number, status5xx:number
}

export interface BaseQueryData {
  clusters: Array<string>
  date: Array<string|number>
  [k:string]:any
}

export interface TotalQueryData{
  uuid:string,
  clusters:Array<string>,
  start:number,
  end:number,
  [key:string]:any
}

export interface TableQueryData{
  pageNum:number
  pageSize:number
  total:number
  keyword:string
}

export interface TableConfigEmitData{
  thead:THEAD_TYPE[],
  config:{[k:string]:boolean} | undefined
}

// 监控告警中用到的折线图封装
// 监控总览：y轴左右两侧各有标题，标题会随x轴变化而更新（时间跨度）
// 调用量统计，共六条线，默认出现四条实线
// 报文量统计，共两条线
// 调用趋势：六条线，并且会出现需要加入对比的可能
// 预留数据给标题，标题根据需求
export interface InvokeData{
  date:Array<string>
  requestTotal:Array<string>
  requestRate:Array<string>
  proxyTotal:Array<string>
  proxyRate:Array<string>
  status4xx:Array<string>
  status5xx:Array<string>
  timeInterval?:string
  [key:string]:any
}

export interface MessageData{
  date:Array<string>
  request:Array<string>
  response:Array<string>
  [key:string]:any
}

// 从总调用数据传递到调用详情的query参数
export interface QueryDetailData{
  time?:string
  startTime?:number
  endTime?:number
}

export interface QueryData {
  startTime: number
  endTime: number
  partitionId: string
  apiId?:string,
  clusters?:Array<string>,
  appId?:string,
  path?:string,
  ip?:string,
  addr?:string,
  proxyPath?:string,
  serviceName?:string
  services?: Array<string>
  apiIds?: Array<string>
  appIds?:Array<string>
}

export interface StrategyQueryData {
  strategyName: string
  warnDimension: string | Array<string>
  status: string | number
  partitionId: string
  pageNum: number
  pageSize: number
  total: number
}

export interface StrategyHistoryQueryData {
  startTime: number
  endTime: number
  partitionId: string
  pageNum: number
  pageSize: number
  total: number
  strategyName: string
}

export interface MonitorAlarmStrategyRuleConditionData {
  compare: string
  unit: string
  value: number | null
}

export interface MonitorAlarmStrategyRuleData {
  channelUuids: Array<string>
  condition: MonitorAlarmStrategyRuleConditionData[]
}

export interface MonitorAlarmStrategyData {
  uuid: string
  title: string
  desc: string
  isEnable: boolean
  dimension: string
  target: {
    rule?: string
    values?: Array<string>
  }
  quota: string
  every: number
  rule: MonitorAlarmStrategyRuleData[]
  continuity: number
  hourMax: number
  users: Array<string>
  partitionId?: string
  [key: string]: any
}

export interface MonitorAlarmStrategyListData {
  uuid: string
  strategyTitle: string
  warnDimension: string
  warnTarget: string
  warnRule: string
  warnFrequency: string
  isEnable: boolean
  operator: string
  updateTime: string
  createTime: string
}

export interface MonitorAlarmHistoryData {
  strategyTitle: string
  warnTarget: string
  warnContent: string
  createTime: string
}

export interface MonitorAlarmChannelsData {
  uuid: string
  title: string
  type: 1 | 2
}

export interface MonitorAlarmStrategyTargetValueData {
  contain: Array<string>
  // eslint-disable-next-line camelcase
  not_contain: Array<string>
  [key: string]: Array<string>
}

// 监控基础数据(表格用)
export interface MonitorData{
  proxyTotal:number,
  proxySuccess:number,
  proxyRate:number,
  statusFail:number,
  avgResp:number,
  maxResp:number,
  minResp:number,
  avgTraffic:number,
  maxTraffic:number,
  minTraffic:number,
  [key:string]:any
}

// 监控Api数据(表格用)
export interface MonitorApiData extends MonitorData{
  apiId?:string,
  apiName:string,
  serviceName:string,
  path:string,
  requestTotal:number,
  requestSuccess:number,
  requestRate:number
}
export interface MonitorApiProxyData extends MonitorData{
  apiId?:string,
  apiName:string
}

export interface MonitorAppData extends MonitorData{
  appName:string,
  appId:string,
  requestTotal:number,
  requestSuccess:number,
  requestRate:number
}

export interface MonitorUpstreamData extends MonitorData{
  serviceName:string
}

export interface MonitorIPData extends MonitorData{
  ip:string,
  requestTotal:number,
  requestSuccess:number,
  requestRate:number
}

export interface MonitorNodeData extends MonitorData{
  addr:string
}

export interface MonitorPathData extends MonitorData{
  path:string
}

export interface MonitorProxyData extends MonitorData{
  proxyPath:string
}

export interface MonitorProxyTableConfig {
  proxyTotal:boolean, proxySuccess :boolean, proxyRate :boolean, statusFail:boolean,
    avgResp :boolean, maxResp :boolean, minResp:boolean, avgTraffic :boolean, maxTraffic:boolean, minTraffic:boolean,
    [k:string]:boolean
}

export interface MonitorRequestTableConfig extends MonitorProxyTableConfig {
  requestTotal:boolean, requestSuccess:boolean, requestRate:boolean
}
