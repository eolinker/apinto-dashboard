import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE, EO_THEAD_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const deployCertsTableHeadName: THEAD_TYPE[] = [
  { title: '证书' },
  { title: '证书有效期' },
  { title: '更新者' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const deployCertsTableBody: TBODY_TYPE[] = [
  { key: 'name' },
  { key: 'validTime' },
  { key: 'operator' },
  { key: 'updateTime' },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '修改'
      },
      {
        title: '删除'
      }
    ]
  }
]

export const deployConfigRedisTableHeadName:EO_THEAD_TYPE[] = [
  {
    title: '地址',
    titleString: '地址',
    resizeable: true,
    required: true
  },
  {
    title: '用户名',
    resizeable: true
  },
  {
    title: '密码',
    resizeable: true
  },
  {
    title: '启用',
    width: 90,
    resizeable: true
  },
  {
    title: '操作',
    width: 60,
    resizeable: false,
    right: true
  }
]

export const deployConfigRedisTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'addrs',
    type: 'input',
    placeholder: '请输入域名/IP：端口，多个以逗号分隔',
    checkMode: 'change',
    errorTip: '请输入域名/IP：端口，多个以逗号分隔'
  },
  {
    key: 'username',
    type: 'input',
    placeholder: '请输入用户名'
  },
  {
    key: 'password',
    type: 'input',
    placeholder: '请输入密码'
  },
  { key: 'enable' },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '测试'
      }
    ]
  }
]

export const deployConfigInfluxdbTableHeadName:EO_THEAD_TYPE[] = [
  {
    title: '数据源地址',
    resizeable: true,
    titleString: '数据源地址',
    required: true
  },
  {
    title: 'Organization',
    resizeable: true,
    titleString: 'Organization',
    required: true
  },
  {
    title: '鉴权token',
    resizeable: true
  },
  {
    title: '启用',
    width: 90,
    resizeable: true
  },
  {
    title: '操作',
    width: 60,
    resizeable: false,
    right: true
  }
]

export const deployConfigInfluxdbTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'addr',
    type: 'input',
    placeholder: '请输入数据源地址',
    checkMode: 'change',
    errorTip: '请输入数据源地址'
  },
  {
    key: 'org',
    type: 'input',
    placeholder: '请输入Organization'
  },
  {
    key: 'token',
    type: 'input',
    placeholder: '请输入鉴权Token'
  },
  { key: 'enable' },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '测试'
      }
    ]
  }
]

export const DeployClusterNodeThead:THEAD_TYPE[] = [
  {
    title: '名称'
  },
  {
    title: '管理地址'
  },
  {
    title: '服务地址'
  },
  {
    title: '状态'
  }
]

export const DeployClusterNodeTbody:TBODY_TYPE[] = [
  { key: 'name', ellipsis: true },
  { key: 'adminAddr', ellipsis: true },
  { key: 'serviceAddr', ellipsis: true },
  { key: 'status' }
]

export const DeployClusterOperateRecordThead:THEAD_TYPE[] = [
  { title: 'KEY', resizeable: true },
  { title: 'OLD VALUE', resizeable: true },
  { title: 'NEW VALUE', resizeable: true },
  { title: '类型', resizeable: true },
  { title: '操作时间' }
]

export const DeployClusterOperateRecordTbody:TBODY_TYPE[] = [
  { key: 'key' },
  { key: 'oldValue' },
  { key: 'newValue' },
  { key: 'optType' },
  { key: 'createTime' }
]

export const DeployClusterPublishRecordThead:THEAD_TYPE[] = [
  { width: 45 },
  { title: '版本名称', resizeable: true },
  { title: '发布者', resizeable: true },
  { title: '发布时间' }
]

export const DeployClusterPublishRecordTbody:TBODY_TYPE[] = [
  {
    key: ''
  },
  { key: 'name' },
  { key: 'operator' },
  { key: 'createTime' }
]

export const DeployClusterPublishThead:THEAD_TYPE[] = [
  { title: 'KEY', resizeable: true },
  { title: '发布的值', resizeable: true },
  { title: '未发布的值', resizeable: true },
  { title: '类型', resizeable: true },
  { title: '操作时间' }
]

export const DeployClusterPublishTbody:TBODY_TYPE[] = [
  { key: 'key' },
  { key: 'finishValue' },
  { key: 'noReleasedValue' },
  { key: 'optType' },
  { key: 'createTime' }
]

export const DeployClusterEnvConfigThead:THEAD_TYPE[] = [
  { title: 'KEY', resizeable: true },
  { title: 'VALUE', resizeable: true },
  { title: '描述', resizeable: true },
  { title: '发布状态', resizeable: true },
  { title: '更新者', resizeable: true },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const ClustersThead: THEAD_TYPE[] = [
  { title: '集群名称' },
  { title: '环境' },
  { title: '状态' },
  {
    title: '操作',
    right: true
  }
]

export const DeployClusterPluginThead:THEAD_TYPE[] = [
  { title: '插件名称' },
  { title: '发布状态' },
  { title: '状态' },
  { title: '配置' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const DeployClusterPluginStatusOptions = [
  { label: '全局启用', value: 'GLOBAL' },
  { label: '启用', value: 'ENABLE' },
  { label: '禁用', value: 'DISABLE' }
]

export const DeployClusterPluginPublishThead:THEAD_TYPE[] = [
  { title: '插件名称' },
  { title: '发布的配置' },
  { title: '未发布的配置' },
  { title: '类型' },
  { title: '操作时间' }
]

export const DeployClusterPluginPublishTbody:EO_TBODY_TYPE[] = [
  { key: 'name' },
  {
    key: 'finishValue',
    json: true
  },
  {
    key: 'noReleasedValue',
    json: true
  },
  { key: 'optType' },
  { key: 'createTime' }
]

export const DeployClusterPluginChangeHistoryThead:THEAD_TYPE[] = [
  { title: '插件名称' },
  { title: '旧配置' },
  { title: '新配置' },
  { title: '类型' },
  { title: '操作时间' }
]

export const DeployClusterPluginChangeHistoryTbody:EO_TBODY_TYPE[] = [
  { key: 'name' },
  {
    key: 'oldValue',
    json: true
  },
  {
    key: 'newValue',
    json: true
  },
  { key: 'optType' },
  { key: 'createTime' }
]
