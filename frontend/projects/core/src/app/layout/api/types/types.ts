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

export type ApiPublishStatus = 'GOONLINE'|'OFFLINE'|'NOTGOONLINE'|'TOUPDATE'

export type ApiListItem = {
    checked?:boolean
    groupUuid:string
    uuid:string
    name:string
    scheme:'http'|'websocket'
    method:string
    requestPath:string
    publish:Array<{name:string, status:ApiPublishStatus}>
    source:string
    updateTime:string
    isDelete:boolean
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
    env:string
    status:'GOONLINE'|'OFFLINE'|'NOTGOONLINE'|'TOUPDATE'
    operator:string
    updateTime:string
}
