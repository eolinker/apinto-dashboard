import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE } from 'eo-ng-table'

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
