import { PublishStatus } from '../../../constant/type'

export type APIProtocol = 'http'|'websocket'

export interface RouterEnum{
    apiId:string
    name:string
}

// 高级匹配
export interface MatchData{
    position:string
    matchType:string
    key:string
    pattern?:string
}

export interface ProxyHeaderData{
    optType:string
    key:string
    value:string
}

export type ApiMessage = {
    name:string
    uuid:string
    groupUuid:string
    desc: string
      isDisable: boolean
      scheme: string
      requestPath: string
      service:string
      method: Array<string>
      proxyPath:string
      hosts: Array<string>
      timeout: number,
      retry: number,
      match:Array<MatchData>
      proxyHeader: Array<ProxyHeaderData>
      templateUuid: string
}

// api创建表单时部分不能放在FormGroup的数据
export interface APINotFormGroupData{
    uuid:string,
    method:Array<string>,
    match:MatchData[],
    proxyHeader:ProxyHeaderData[]
}

export interface APIImportData{
    id:number
    name:string
    desc:string
}

// api批量上线检测列表
export interface APIBatchOnlineVerifyData{
    service:string
    cluster:string
    status:boolean
    statusString:string
    result:string
    solution:{params:any, name:string}
    name:string
}

export interface APIBatchPublishData{
    api:string
    cluster:string
    status:boolean
    statusString?:string
    result:string
}

export type PluginTemplateItem = {
    uuid:string
    name:string
    desc:string
    createTime:string
    updateTime:string
    operator:string
    isDelete:boolean
}

export type PluginTemplateConfigItem = {
    name:string
    config:string
    disable:boolean
    eoKey?:string
}

export type PluginTemplateData = {
    name:string
    desc:string
    plugins:PluginTemplateConfigItem[]
}

export type ApiListItem = {
    checked?:boolean
    groupUuid:string
    uuid:string
    name:string
    scheme:'http'|'websocket'
    method:string
    requestPath:string
    publish:Array<{name:string, status:PublishStatus}>
    source:string
    updateTime:string
    isDelete:boolean
    isDisable:boolean|string
    [k:string]:any
}

export type ApiData = {
    name:string
    id:string
    scheme:'http'|'websocket'
    method:string
    path:string
    service:string
    proxyPath:string
    desc:string
}

export type ApiPublishItem = {
    checked?:boolean
    name:string
    title:string
    env:string
    status:'GOONLINE'|'OFFLINE'|'NOTGOONLINE'|'TOUPDATE'
    operator:string
    updateTime:string
}
