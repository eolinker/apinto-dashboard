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
  { label: 'Query', value: 'query' }
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
    key: 'desc',
    copy: true
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
  { title: '鉴权类型' },
  { title: '参数位置' },
  { title: '参数名' },
  { title: '参数信息' },
  { title: '过期时间' },
  { title: '透传上游' },
  { title: '更新者' },
  { title: '更新时间' },
  {
    title: '操作',
    right: true
  }
]

export const authenticationTableBody:EO_TBODY_TYPE[] = [
  {
    key: 'driver',
    copy: true
  },
  {
    key: 'paramPosition',
    copy: true
  },
  {
    key: 'paramName',
    copy: true
  },
  {
    key: 'paramInfo',
    copy: true
  },
  {
    key: 'expireTimeString',
    copy: true
  },
  { key: 'isTransparent' },
  { key: 'operator' },
  { key: 'updateTime' },
  {
    type: 'btn',
    right: true,
    btns: [{
      title: '查看'
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
