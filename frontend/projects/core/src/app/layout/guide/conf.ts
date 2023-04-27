export type StepItem = {name:'cluster'|'upstream'|'api'|'publishApi',
title:string,
desc:Array<string>,
status:'undo'|'doing'|'done'
img?:string
toDoUrl:string
doneUrl:string}

export type TutorialItem = {
      title:string,
      content?:Array<
        {text:string, url?:string, children?:Array<{text:string, url:string}>}
      >,

    }

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
    title: '仪表盘',
    content: [
      { text: '统计报表', url: 'https://help.eolink.com/tutorial/Apinto/c-1371' },
      { text: '监控告警', url: 'https://help.eolink.com/tutorial/Apinto/c-1372' }
    ]
  },
  {
    title: 'API',
    content: [
      { text: '添加 API', url: 'https://help.eolink.com/tutorial/Apinto/c-1373' },
      { text: '发布 API', url: 'https://help.eolink.com/tutorial/Apinto/c-1374' },
      { text: '对 API 设置额外操作', url: 'https://help.eolink.com/tutorial/Apinto/c-1375' }
    ]
  },
  {
    title: '上游',
    content: [
      { text: '添加静态上游', url: 'https://help.eolink.com/tutorial/Apinto/c-1346' },
      { text: '添加动态（服务发现）上游', url: 'https://help.eolink.com/tutorial/Apinto/c-1376' }
    ]
  },
  {
    title: '应用',
    content: [
      { text: '添加应用', url: 'https://help.eolink.com/tutorial/Apinto/c-1349' },
      { text: '发布应用', url: 'https://help.eolink.com/tutorial/Apinto/c-1350' },
      { text: '设置访问鉴权', url: 'https://help.eolink.com/tutorial/Apinto/c-1377' }
    ]
  },

  {
    title: '策略',
    content: [
      { text: '限制访问范围', url: 'https://help.eolink.com/tutorial/Apinto/c-1356' },
      { text: '熔断降级', url: 'https://help.eolink.com/tutorial/Apinto/c-1359' },
      { text: '灰度发布', url: 'https://help.eolink.com/tutorial/Apinto/c-1360' },
      { text: '数据缓存', url: 'https://help.eolink.com/tutorial/Apinto/c-1357' },
      { text: '流量限制', url: 'https://help.eolink.com/tutorial/Apinto/c-1358' }
    ]
  },
  {
    title: '基础设施',
    content: [
      { text: '集群', url: 'https://help.eolink.com/tutorial/Apinto/c-1343' },
      {
        text: '服务发现',
        children: [
          { text: 'Consul', url: 'https://help.eolink.com/tutorial/Apinto/c-1380' },
          { text: 'Eureka', url: 'https://help.eolink.com/tutorial/Apinto/c-1381' },
          { text: 'Nacos', url: 'https://help.eolink.com/tutorial/Apinto/c-1382' }
        ]
      },
      { text: '环境变量', url: 'https://help.eolink.com/tutorial/Apinto/c-1344' },
      { text: '节点插件', url: 'https://help.eolink.com/tutorial/Apinto/c-1345' }
    ]
  },
  {
    title: '开发',
    content: [
      { text: 'Open API', url: 'https://help.eolink.com/tutorial/Apinto/c-1362' },
      { text: 'Webhook', url: 'https://help.eolink.com/tutorial/Apinto/c-1386' },
      { text: 'CLI', url: 'https://help.eolink.com/tutorial/Apinto/c-1385' },
      { text: '控制台 Debug 日志', url: 'https://help.eolink.com/tutorial/Apinto/c-1361' }
    ]
  }
]
