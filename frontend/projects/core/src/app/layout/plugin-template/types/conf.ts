import { THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const PluginTemplateTableHeadName:THEAD_TYPE[] = [
  { title: '模板名称' },
  { title: '描述' },
  { title: '创建时间' },
  { title: '更新时间' },
  { title: '操作' }
]

export const PluginTemplateConfigThead:THEAD_TYPE[] = [
  { title: '插件名称' },
  { title: '状态' },
  { title: '配置' },
  { title: '操作' }
]

export const PluginTemplateConfigTbody:EO_TBODY_TYPE[] = [
  {
    key: 'name',
    copy: true
  },
  { key: 'disable' },
  {
    key: 'config',
    json: true,
    copy: true
  },
  {
    type: 'btn',
    right: true,
    btns: [{
      title: '配置'
    },
    {
      title: '删除',
      action: 'delete'
    }
    ]
  }
]

export const PluginTemplatePublishThead:THEAD_TYPE[] = [
  { title: '集群名称' },
  { title: '环境' },
  { title: '状态' },
  { title: '更新者' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const PluginTemplatePublishTbody:EO_TBODY_TYPE[] = [
  {
    key: 'title',
    copy: true
  },
  {
    key: 'env'
  },
  { key: 'status' },
  { key: 'operator' },
  { key: 'updateTime' },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return item.status === 'TOUPDATE' && !item.disable
    },
    btns: [
      {
        title: '更新'
      },
      {
        title: '下线'
      },
      {
        title: '禁用'
      }
    ]
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return item.status === 'TOUPDATE' && item.disable
    },
    btns: [
      {
        title: '更新'
      },
      {
        title: '下线'
      },
      {
        title: '启用'
      }
    ]
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return item.status === 'GOONLINE' && !item.disable
    },
    btns: [
      {
        title: '下线'
      },
      {
        title: '禁用'
      }
    ]
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return item.status === 'GOONLINE' && item.disable
    },
    btns: [
      {
        title: '下线'
      },
      {
        title: '启用'
      }
    ]
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return (item.status === 'OFFLINE' || item.status === 'NOTGOONLINE') && !item.disable
    },
    btns: [
      {
        title: '上线'
      },
      {
        title: '禁用'
      }
    ]
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return (item.status === 'OFFLINE' || item.status === 'NOTGOONLINE') && item.disable
    },
    btns: [
      {
        title: '上线'
      },
      {
        title: '启用'
      }
    ]
  }
]
