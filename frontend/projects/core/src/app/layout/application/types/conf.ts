import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const algorithmList:SelectOption[] = [
  { label: 'HS256', value: 'HS256' },
  { label: 'HS384', value: 'HS384' },
  { label: 'HS512', value: 'HS512' },
  { label: 'RS256', value: 'RS256' },
  { label: 'RS384', value: 'RS384' },
  { label: 'RS512', value: 'RS512' },
  { label: 'ES256', value: 'ES256' },
  { label: 'ES384', value: 'ES384' },
  { label: 'ES512', value: 'ES512' }
]

export const positionList:SelectOption[] = [
  { label: 'Header', value: 'header' },
  { label: 'Query', value: 'query' },
  { label: 'Body', value: 'body' }
]

export const verifyList:SelectOption[] = [
  { label: 'exp', value: 'exp' },
  { label: 'nbf', value: 'nbf' }
]

export const applicationsTableHeadName: THEAD_TYPE[] = [
  { title: '应用名称' },
  { title: '应用ID' },
  { title: '描述' },
  { title: '更新者' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const applicationsTableBody: EO_TBODY_TYPE[] = [
  {
    key: 'name',
    copy: true
  },
  {
    key: 'id',
    copy: true
  },
  {
    key: 'desc'
  },
  { key: 'operator' },
  { key: 'updateTime' },
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

export const authLabelTableBody:TBODY_TYPE[] = [
  {
    key: 'key',
    type: 'input',
    placeholder: '请输入Key'
  },
  {
    key: 'value',
    type: 'input',
    placeholder: '请输入Value'
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加',
        action: 'add'
      }
    ]
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加',
        action: 'add'
      },

      {
        title: '减少',
        action: 'delete'
      }
    ]
  }
]

export const authenticationTableHeadName:THEAD_TYPE[] = [
  { title: '名称' },
  { title: '鉴权类型' },
  { title: '隐藏鉴权信息' },
  { title: '到期时间' },
  {
    title: '操作',
    right: true
  }
]

export const authenticationTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'title',
    copy: true
  },
  {
    key: 'driver'
  },
  { key: 'hideCredential' },
  {
    key: 'expireTimeString'
  },
  {
    type: 'btn',
    right: true,
    btns: [{
      title: '查看'
    },
    {
      title: '修改'
    },
    {
      title: '删除'
    }
    ]
  }]

export const customAttrTableBody: EO_TBODY_TYPE[] = [
  {
    key: 'key',
    type: 'input',
    placeholder: '请输入Key',
    checkMode: 'change',
    check: (item: any) => {
      return !item || /^[a-zA-Z_][a-zA-Z0-9_]*$/.test(item)
    },
    errorTip: '首字母必须为英文'
  },
  {
    key: 'value',
    type: 'input',
    placeholder: '请输入Value'
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加'
      }
    ]
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加'
      },
      {
        title: '减少'
      }
    ]
  }
]

export const extraHeaderTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'key',
    type: 'input',
    placeholder: '请输入Key'
  },
  {
    key: 'value',
    type: 'input',
    placeholder: '请输入Value'
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加'
      }
    ]
  },
  {
    type: 'btn',
    btns: [
      {
        title: '添加'
      },

      {
        title: '减少'
      }
    ]
  }
]

export const extraTableHeadName:THEAD_TYPE[] = [
  { title: '参数位置' },
  { title: '参数名' },
  { title: '参数值' },
  { title: '生效规则' },
  {
    title: '操作',
    right: true
  }
]

export const extraTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'position'
  },
  {
    key: 'key'
  },
  { key: 'value' },
  {
    key: 'conflict'
  },
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
  }]

export const extraConflictList:SelectOption[] = [
  {
    label: '替换原参数值',
    value: 'convert'
  },
  {
    label: '原参数值不存在时，添加额外参数',
    value: 'origin'
  },
  {
    label: '原参数值存在时返回错误',
    value: 'error'
  }
]
