export interface DeployCertListData{
    id:string
    name:string
    validTime:string
    operator:string
    createTime:string
    updateTime:string
}

export interface RedisData{
    addrs:string,
    username:string,
    password:string,
    enable:boolean|null,
    [key:string]:any
  }

export interface InfluxdbData{
    addr:string,
    org:string,
    token:string,
    enable:boolean|null,
    [key:string]:any
  }

export type ClusterPluginItem = {
  name:string
  publish:'UNPUBLISHED' | 'PUBLISHED' | 'DEFECT'
  changeStatus:'NONE' | 'NEW' | 'MODIFY' | 'DELETE',
  config:object
  status:'GLOBAL'|'DISABLE'|'ENABLE'
  releasedSort:number
  nowSort:number
  operator:string
  createTime:string
  updateTime:string
  isBuiltin:boolean
}

export type PluginPublishStatus = 'GLOBAL'|'DISABLE'|'ENABLE'

export type ClusterPluginConfig = {
  name:string
  status:PluginPublishStatus
  config: string
}

export type ClusterPluginPublishData = {
  name:string
  releasedConfig:{
    status:PluginPublishStatus
    config:string
  }
  noReleasedConfig:{
    status:PluginPublishStatus
    config:string
  }
  createTime:string
  optType:'NEW'|'MODIFY'|'DELETE'
  releasedSort:number
  nowSort:number,
  finishValue?:string
  noReleasedValue?:string
  finishValueTooltip?:string
  noReleasedValueTooltip?:string
}

export type ClusterEnvChangeHistoryItem = {
    key:string
    oldValue:string
    newValue:string
    createTime:string
    optType:string
}

export type ClusterPluginChangeHistoryItem = {
  name: string
  oldConfig: {
    status: PluginPublishStatus
    config: string
  },
  newConfig: {
    status: PluginPublishStatus
    config: string
  },
  createTime: string
  optType: 'NEW'|'MODIFY'|'DELETE'
  oldValue?:string
  newValue?:string
}

export type ClusterPluginPublishHistoryItem = {
  id:number
  name:string
  createTime:string
  optType:'PUBLISH'|'ROLLBACK'
  operator:string
  details:ClusterPluginChangeHistoryItem[]
  total:number
  isExpand?:boolean
}
