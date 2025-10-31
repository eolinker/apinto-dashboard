export interface FilterShowData{
    title?: string
    name: string
    label?: string
    values: Array<string>
    [key: string]: any
  }

export interface FilterForm {
  name: string
  title: string
  values: Array<any>
  label: string
  text: string
  allChecked?: boolean
  showAll?: boolean
  total?: number
  groupUuid?: string
  pattern: RegExp | null
  patternIsPass: boolean
  [key: string]: any
}

export interface FilterOption{
  name:string
  title:string
  type:string
  pattern:string
  options:Array<string>
  label?:string
  value?:string
  [key:string]:any
}

export type RemoteAppItem = {
  title: string
  uuid: string
  desc: string
  checked?:boolean
  [key: string]: any
}

export type RemoteApiItem = {
  uuid: string
  title: string
  service: string
  group: string
  requestPath: string
  checked?:boolean
  [key: string]: any
}

export type RemoteServiceItem = {
  uuid: string
  title: string
  scheme: string
  desc: string
  checked?:boolean
  [key: string]: any
}
export interface FilterRemoteOption{
  target: 'apis'|'services'|'applications'
  titles: Array<{title:string, field:string, [k:string]:any}>
  list:Array<RemoteApiItem|RemoteServiceItem|RemoteAppItem>
  total: number,
  key:string
}

export interface CacheStrategyData {
    name: string
    uuid?: string
    desc?: string
    priority?: number | null
    filters: Array<{
      name: string
      values: Array<string>
      type?: string
      label?: string
      title?: string
    }>
    config: {
      validTime: number
    }
    [key: string]: any
  }

export interface MaskRuleData{
      match:{
        type:'inner'|'keyword'|'regex'|'json_path'
        value:'name'|'phone'|'id-card'|'bank-card'|'date'|'amount'|string
      }
      mask:{
        type:'partial-display'|'partial-masking'|'truncation'|'replacement'|'shuffling'
        begin:number
        length:number
        replace:{
          type:'random'|'custom'
          value:string
        }
      }
      // eslint-disable-next-line camelcase
      eoKey?:string
  }

export interface DataMaskData {
    name: string
    uuid?: string
    desc?: string
    priority?: number | null
    filters: Array<{
      name: string
      values: Array<string>
      type?: string
      label?: string
      title?: string
    }>
    config: {
      rules:MaskRuleData[]
    }
  }

export interface FuseStrategyData {
    name: string
    uuid?: string
    desc?: string
    priority?: number | null
    filters: Array<{
      name: string
      values: Array<string>
      type?: string
      label?: string
      title?: string
    }>
    config: {
      metric: string
      fuseCondition: { statusCodes: Array<number>; count: number }
      fuseTime: { time: number; maxTime: number }
      recoverCondition: { statusCodes: Array<number>; count: number }
      response: {
        statusCode: number
        contentType: string
        charset: string
        header: Array<{ key: string; value: string }>
        body: string
      }
    }
    [key: string]: any
  }

export interface GreyStrategyData {
    name: string
    uuid?: string
    desc?: string
    priority?: number | null
    filters: Array<{
      name: string
      values: Array<string>
      type?: string
      label?: string
      title?: string
    }>
    config: {
      keepSession:boolean,
      nodes:Array<string>,
      distribution:string,
      percent?:number,
      match?:Array<{position:string, matchType:string, key:string, pattern?:string}>
    }
    [key: string]: any
  }

export interface TrafficStrategyData {
    name: string
    uuid?: string
    desc?: string
    priority?: number | null
    filters: Array<{
      name: string
      values: Array<string>
      type?: string
      label?: string
      title?: string
      [key: string]: any
    }>
    config: {
      metrics: Array<string>
      query: { second: number; minute: number; hour: number }
      traffic: { second: number; minute: number; hour: number }
      response: {
        statusCode: number
        contentType: string
        charset: string
        header: Array<{ key: string; value: string }>
        body: string
      }
    }
    [key: string]: any
  }

export interface StrategyListData{
    uuid: string
    name: string
    priority: number
    isStop: boolean
    isDeleted: boolean
    status: 'TOUPDATE'|'GOONLINE'|'TODELETE'|'NOTGOONLINE'
    filters: string,
    conf: string,
    operator: string,
    updateTime: string
}

export interface VisitStrategyData {
  name: string
  uuid?: string
  desc?: string
  priority?: number | null
  filters: Array<{
    name: string
    values: Array<string>
    type?: string
    label?: string
    title?: string
  }>
  config: {
    visitRule: string
    influenceSphere: Array<{
      name: string
      values: Array<string>
      type?: string
      label?: string
      title?: string
    }>,
    continue:boolean
  }
  [key: string]: any
}

export interface StrategyPublishListData{
  source:string
  strategies:Array<{
    name:string
    priority:number
    status:'TOUPDATE'|'GOONLINE'|'TODELETE'|'NOTGOONLINE'
    optTime:string
  }>
  isPublish:boolean
  versionName:string
  unpublishMsg:string
}
