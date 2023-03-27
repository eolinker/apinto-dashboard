import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzCheckBoxOptionInterface } from 'ng-zorro-antd/checkbox'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const optTypeList: SelectOption[] = [
  { label: '新增或修改', value: 'ADD' },
  { label: '删除', value: 'DELETE' }
]

export const methodList:NzCheckBoxOptionInterface[] = [
  { label: 'GET', value: 'GET', checked: false },
  { label: 'POST', value: 'POST', checked: false },
  { label: 'PUT', value: 'PUT', checked: false },
  { label: 'DELETE', value: 'DELETE', checked: false },
  { label: 'PATCH', value: 'PATCH', checked: false },
  { label: 'HEAD', value: 'HEAD', checked: false }
]

export const positionList:SelectOption[] = [
  { label: 'HTTP请求头', value: 'header' },
  { label: '请求参数', value: 'query' },
  { label: 'Cookie', value: 'cookie' }
]

export const prefixMatchList:SelectOption[] = [
  { label: '全等匹配', value: 'EQUAL' },
  { label: '前缀匹配', value: 'PREFIX' },
  { label: '后缀匹配', value: 'SUFFIX' },
  { label: '子串匹配', value: 'SUBSTR' },
  { label: '非等匹配', value: 'UNEQUAL' },
  { label: '空值匹配', value: 'NULL' },
  { label: '存在匹配', value: 'EXIST' },
  { label: '不存在匹配', value: 'UNEXIST' },
  { label: '区分大小写的正则匹配', value: 'REGEXP' },
  { label: '不区分大小写的正则匹配', value: 'REGEXPG' },
  { label: '任意匹配', value: 'ANY' }
]

// 高级匹配表格参数
export const matchHeaderTableHeadName:THEAD_TYPE[] = [
  {
    title: '参数位置'
  },
  { title: '参数名' },
  { title: '匹配类型' },
  { title: '匹配值' },
  {
    title: '操作',
    right: true
  }
]

export const matchHeaderTableBody:TBODY_TYPE[] = [

  { key: 'position' },
  { key: 'key' },
  { key: 'matchType' },
  { key: 'pattern' },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '配置'
      },
      {
        title: '删除',
        action: 'delete'
      }
    ]
  }
]

// 转发上游请求头表格参数
export const proxyHeaderTableHeadName:THEAD_TYPE[] = [
  {
    title: '操作类型'
  },
  { title: '参数名' },
  { title: '参数值' },
  {
    title: '操作',
    right: true
  }
]

export const proxyHeaderTableBody:TBODY_TYPE[] = [
  { key: 'optType' },
  { key: 'key' },
  { key: 'value' },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '配置'
      },
      {
        title: '删除',
        action: 'delete'
      }
    ]
  }
]

// api列表页
export const apisTableHeadName:THEAD_TYPE[] = [
  {
    type: 'checkbox',
    resizeable: false
  },
  {
    title: 'API名称'
  },
  {
    title: '协议/方法',
    width: 140,
    resizeable: false
  },
  {
    title: '上游服务名称'
  },
  {
    title: '请求路径'
  },
  {
    title: '来源',
    filterMultiple: true,
    filterOpts: [{
      text: '自建',
      value: 'build'
    },
    {
      text: '导入',
      value: 'import'
    }
    ],
    filterFn: () => {
      return true
    }
  },
  {
    title: '更新时间'
  },
  {
    title: '操作',
    right: true
  }
]

export const apisTableBody:TBODY_TYPE[] = [
  {
    key: 'checked',
    type: 'checkbox'
  },
  {
    key: 'name'
  },
  {
    key: 'method'
  },
  {
    key: 'service'
  },
  {
    key: 'requestPath'
  },
  {
    key: 'source'
  },
  {
    key: 'updateTime'
  },
  {
    type: 'btn',
    right: true,
    btns: [{
      title: '上线管理'
    },
    {
      title: '查看'
    },
    {
      title: '删除'
    }
    ]
  }
]

// 导入api检测结果页表格参数
export const apiImportCheckResultTableHeadName:THEAD_TYPE[] = [
  {
    type: 'checkbox',
    key: 'checked',
    resizeable: false
  },
  {
    title: '序号',
    width: 58
  },
  {
    title: 'API名称'
  },
  {
    title: '协议/方法',
    width: 99
  },
  {
    title: '请求路径'
  },
  {
    title: '描述'
  },
  {
    title: '状态',
    width: 70,
    resizeable: false
  }]

export const apiImportCheckResultTableBody:TBODY_TYPE[] = [
  {
    key: 'checked',
    type: 'checkbox'
  },
  {
    key: 'id',
    styleFn: (item:any) => {
      if (item.disabled) {
        return 'color:var(--disabled-text-color)'
      }
      return ''
    }
  },
  {
    key: 'name',
    type: 'input',
    checkMode: 'change',
    errorTip: '必填项',
    placeholder: '请输入'
  },
  {
    key: 'method',
    styleFn: (item:any) => {
      if (item.disabled) {
        return 'color:var(--disabled-text-color)'
      }
      return ''
    }
  },
  {
    key: 'path',
    styleFn: (item:any) => {
      if (item.disabled) {
        return 'color:var(--disabled-text-color)'
      }
      return ''
    }
  },
  {
    key: 'desc',
    styleFn: (item:any) => {
      if (item.disabled) {
        return 'color:var(--disabled-text-color)'
      }
      return ''
    }
  },
  {
    key: 'statusString',
    styleFn: (item:any) => {
      if (item.disabled) {
        return 'color:var(--disabled-text-color)'
      }
      return ''
    }
  }]

export const apiBatchOnlineVerifyTableHeadName :THEAD_TYPE[] = [
  {
    title: '上游名称/插件模板'
  },
  { title: '集群名称' },
  { title: '状态' },
  { title: '失败原因' },
  {
    title: '操作',
    width: 94,
    right: true
  }
]

export const apiBatchOnlineVerifyTableBody :TBODY_TYPE[] = [
  {
    key: 'serviceTemplate',
    styleFn: (item:any) => {
      if (!item.status) {
        return 'color:#ff3b30'
      }
      return ''
    }
  },
  {
    key: 'cluster',
    styleFn: (item:any) => {
      if (!item.status) {
        return 'color:#ff3b30'
      }
      return ''
    }
  },
  {
    key: 'statusString',
    styleFn: (item:any) => {
      if (!item.status) {
        return 'color:#ff3b30'
      }
      return ''
    }
  },
  {
    key: 'result',
    styleFn: (item:any) => {
      if (!item.status) {
        return 'color:#ff3b30'
      }
      return ''
    }
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return item.solution.name
    },
    btns: [{
      title: '解决方案',
      click: (item:any) => {
        let routerS:string = '/' + item.data.solution.name
        const routerSArr:Array<string> = routerS.split('/')
        if (item.data.solution.params.serviceName) {
          routerSArr.splice(-1, 0, item.data.solution.params.serviceName)
        } else {
          routerSArr.push(item.data.solution.params.templateUuid)
        }
        routerS = routerSArr.join('/')
        window.open(routerS, '')
      },
      type: 'text'
    }
    ]
  },
  {
    type: 'btn',
    right: true,
    showFn: (item:any) => {
      return !item.solution.name
    }
  }
]

export const apiBatchPublishResultTableHeadName:THEAD_TYPE[] = [
  {
    title: 'API名称'
  },
  { title: '集群名称' },
  { title: '状态' },
  { title: '失败原因' }
]

export const apiBatchPublishResultTableBody:TBODY_TYPE[] = [
  {
    key: 'api',
    styleFn: (item:any) => {
      if (!item.status) {
        return 'color:#ff3b30'
      }
      return ''
    }
  },
  {
    key: 'cluster',
    styleFn: (item:any) => {
      if (!item.status) {
        return 'color:#ff3b30'
      }
      return ''
    }
  },
  {
    key: 'statusString',
    styleFn: (item:any) => {
      if (!item.status) {
        return 'color:#ff3b30'
      }
      return ''
    }
  },
  {
    key: 'result',
    styleFn: (item:any) => {
      if (!item.status) {
        return 'color:#ff3b30'
      }
      return ''
    }
  }
]

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
  { key: 'name' },
  { key: 'disable' },
  {
    key: 'config',
    json: true
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

export const PluginTemplatePublishTbody:TBODY_TYPE[] = [
  { key: 'name' },
  { key: 'env' },
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
