// webhook设置

import { SelectOption } from 'eo-ng-select'

export const methodsList: SelectOption[] = [
  { label: 'POST', value: 'POST' },
  { label: 'GET', value: 'GET' }
]

export const contentTypesList: SelectOption[] = [
  { label: 'JSON', value: 'JSON' },
  { label: 'form-data', value: 'form-data' }
]

export const noticeTypesList: SelectOption[] = [
  { label: '单次发送', value: 'single' },
  { label: '多次发送', value: 'many' }
]

export const protocolsList: SelectOption[] = [
  { label: '不设置任何协议', value: 'none' },
  { label: 'SSL协议', value: 'ssl' },
  { label: 'TLS协议', value: 'tls' }
]
