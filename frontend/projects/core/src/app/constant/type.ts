/* eslint-disable camelcase */
export interface EmptyHttpResponse{
  code:number
  data:{}
  msg:string
}

export interface Operator{
    userId:number,
    username:string,
    nickname:string,
    avatar:string
}

export interface UserListData{
  id:number
  userName:string
  nickName:string
  email:string
}

export interface UserData{
  sex?:number
  avatar?:string
  email:string
  phone:string
  user_name:string
  nick_name:string
  role_ids:Array<string>
  desc:string
  notice_user_id:string
}

export interface MonitorAlarmChannelsData {
  uuid: string
  title: string
  type: 1 | 2
}

// 创建/修改目录
export interface GroupData{
  name:string
  uuid:string
  parentUuid:string
  tagName?:string
}

export interface ApiGroupsData{
  uuid:string
  name:string
  children:ApiGroupsData[]
  isDelete:boolean
  [key:string]:any
}

export type ModuleOpenConfigDataType = {
  name:string
  value:string
  type?:string
}
export type ModuleOpenConfigData = {
  name:string
  url:string
  query:Array<ModuleOpenConfigDataType>
  header:Array<ModuleOpenConfigDataType>
  initialize:Array<ModuleOpenConfigDataType>
}

// API目录
export interface ApiGroup{
  apis:Array<{
    method:Array<string>
    name:string
    groupUuid:string
    uuid:string
  }>
  root:{
    groups:ApiGroupsData[]
    uuid:string
    name:string
  [key:string]:any
  }
}

// 集群列表
export interface ClustersData{
  env: string,
  status: 'NORMAL'|'PARTIALLY_NORMAL'|'ABNORMAL',
  desc: string,
  title:string,
  name: string,
  createTime: string,
  updateTime: string
}

export type PublishStatus = 'GOONLINE'|'OFFLINE'|'NOTGOONLINE'|'TOUPDATE'

// 简易集群列表，即将集群名与环境名组合在一起的一级目录
export type ClusterSimpleOption = {
  id:string
  name:string
  title:string
}

// 获取远程类型的选项（用在服务治理-筛选条件和监控告警-选择api和上游 -api
export interface RemoteApiData{
  uuid: string,
  name: string,
  service: string,
  group: string,
  requestPath: string,
  [key:string]:any
}

// 获取远程类型的选项（用在服务治理-筛选条件和监控告警-选择api和上游 -services
export interface RemoteServiceData{
  uuid: string,
  name: string,
  scheme: string,
  desc: string,
  [key:string]:any
}

// 获取远程类型的选项（用在服务治理-筛选条件和监控告警-选择api和上游 -applications
export interface RemoteAppData{
  uuid: string,
  name: string,
  desc: string,
  [key:string]:any
}

// 获取远程类型的选项（用在服务治理-筛选条件和监控告警-选择api和上游
export interface RemoteData{
      target: 'apis' | 'services' | 'applications',
      titles: Array<
        {
          title: string,
          field: string
        }>,
      apis: Array<RemoteApiData>,
      services:Array<RemoteServiceData>,
      applications:Array<RemoteAppData>,
      total: number,
      title?:string,
      [key:string]:any
}

// 上线管理列表
export interface PublishManagementData{
  name:string
  env:string
  status:'GOONLINE'|'OFFLINE'|'NOTGOONLINE'|'TOUPDATE'
  disable:boolean
  operator:string
  updateTime:string
}

// 集群列表接口

export type ClusterEnumData = {
  name:string, uuid:string,
  title:string
}
export interface ClusterEnum{
  clusters:Array<ClusterEnumData>
  name:string
}

export interface ArrayItemData{
  key:string
  value:string
  [k:string]:any
}

export interface RandomId{
  id:string
}
