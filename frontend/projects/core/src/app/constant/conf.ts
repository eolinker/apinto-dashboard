import { THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const defaultAutoTips: Record<string, Record<string, string>> = {
  'zh-cn': {
    required: '必填项'
  },
  en: {
    required: 'Input is required'
  },
  default: {
    email: '邮箱格式不正确'
  }
}

// api上线管理列表,app上线管理列表
export const CommonPublishTableHeadName:THEAD_TYPE[] = [
  { title: '集群名称' },
  { title: '环境' },
  { title: '状态' },
  { title: '禁用状态' },
  { title: '更新者' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const PublishTableHeadName:THEAD_TYPE[] = [
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

export const CommonPublishTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'name',
    copy: true
  },
  {
    key: 'env'
  },
  { key: 'status' },
  { key: 'disable' },
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

export const PublishTableBody:EO_TBODY_TYPE[] = [

  {
    key: 'name',
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
      return item.status === 'TOUPDATE'
    },
    btns: [
      {
        title: '更新'
      },
      {
        title: '下线'
      }
    ]
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return item.status === 'GOONLINE'
    },
    btns: [
      {
        title: '下线'
      }
    ]
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return (item.status === 'OFFLINE' || item.status === 'NOTGOONLINE')
    },
    btns: [
      {
        title: '上线'
      }
    ]
  }

]
