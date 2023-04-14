export type StepItem = {name:'cluster'|'upstream'|'api'|'publishApi',
title:string,
desc:Array<string>,
status:'undo'|'doing'|'done'
img?:string
toDoUrl:string
doneUrl:string}

export type TutorialItem = {
      title:string,
      content:Array<
        {text:string, url:string}
      >}

export const GuideStepList:Array<StepItem> = [
  {
    name: 'cluster',
    title: '创建集群',
    desc: ['创建 Apinto 网关集群，集群用于承载网络流量'],
    status: 'undo',
    toDoUrl: 'deploy/cluster/create',
    doneUrl: 'deploy/cluster'
  },
  {
    name: 'upstream',
    title: '添加上游服务',
    desc: ['添加上游服务器或动态服务发现，接收网关节点转发的流量'],
    status: 'undo',
    toDoUrl: 'upstream/upstream/create',
    doneUrl: 'upstream/upstream'
  },
  {
    name: 'api',
    title: '添加 API',
    desc: ['添加需要网关转发的 API'],
    status: 'undo',
    toDoUrl: 'router/api/create',
    doneUrl: 'router/api/group/list'
  },
  {
    name: 'publishApi',
    title: '发布 API',
    desc: ['发布之后就可以通过 Apinto 安全高效地访问API啦！'],
    status: 'undo',
    toDoUrl: 'router/api/group/list',
    doneUrl: 'router/api/group/list'
  }
]

export const TutorialsList:Array<TutorialItem> = [
  {
    title: '安全防护',
    content: [
      { text: '为API设置鉴权/身份认证', url: '' },
      { text: 'xxx', url: '' },
      { text: 'xxx', url: '' }
    ]
  },
  {
    title: '安全防护',
    content: [
      { text: '为API设置鉴权/身份认证', url: '' },
      { text: 'xxx', url: '' },
      { text: 'xxx', url: '' }
    ]
  },
  {
    title: '安全防护',
    content: [
      { text: '为API设置鉴权/身份认证', url: '' },
      { text: 'xxx', url: '' },
      { text: 'xxx', url: '' }
    ]
  },
  {
    title: '安全防护',
    content: [
      { text: '为API设置鉴权/身份认证', url: '' },
      { text: 'xxx', url: '' },
      { text: 'xxx', url: '' }
    ]
  }
]
