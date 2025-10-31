import { THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const DeployGlobalEnvDetailTableHeadName:THEAD_TYPE[] = [
  { title: '集群' },
  { title: '环境' },
  { title: 'VALUE' },
  { title: '状态' }
]

export const DeployGlobalEnvDetailTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'clusterName',
    copy: true
  },
  {
    key: 'environment'
  },
  {
    key: 'value'
  },
  { key: 'publishStatus' }
]

export const DeployGlobalEnvTableHeadName:THEAD_TYPE[] = [
  { title: 'KEY' },
  { title: '描述' },
  { title: '创建者' },
  { title: '创建时间' },
  { title: '状态' },
  {
    title: '操作',
    right: true
  }
]

export const DeployGlobalEnvTableBody:EO_TBODY_TYPE[] = [

  {
    key: 'key',
    copy: true
  },
  { key: 'description' },
  { key: 'operator' },
  { key: 'createTime' },
  { key: 'status' },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '查看'
      },
      {
        title: '删除'
      }
    ]
  }

]
