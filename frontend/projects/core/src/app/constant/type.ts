export interface EmptyHttpResponse{
  code:number
  data:{}
  msg:string
}

/* eslint-disable camelcase */
export interface Operator{
    user_id:number,
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

export interface ApiGroupsData{
  uuid:string
  name:string
  children:ApiGroupsData[]
  is_delete:boolean
  [key:string]:any
}

// API目录
export interface ApiGroup{
  apis:Array<{
    method:Array<string>
    name:string
    group_uuid:string
    uuid:string
  }>
  root:{
    groups:ApiGroupsData[]
    uuid:string
    name:string
  [key:string]:any
  }
}

// 获取远程类型的选项（用在服务治理-筛选条件和监控告警-选择api和上游 -api
export interface RemoteApiData{
  uuid: string,
  name: string,
  service: string,
  group: string,
  request_path: string,
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
