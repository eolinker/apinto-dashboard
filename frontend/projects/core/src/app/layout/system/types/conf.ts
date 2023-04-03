// webhook设置

import { SelectOption } from 'eo-ng-select'
import { THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

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

export const ExternalAppListTableHeadName:THEAD_TYPE[] = [
  {
    title: '应用名称'
  },
  {
    title: '应用ID'
  },
  {
    title: '鉴权Token'
  },
  {
    title: '关联标签'
  },
  {
    title: '禁用状态',
    width: 84
  },
  {
    title: '更新者'
  },
  {
    title: '更新时间',
    showSort: true
  },
  {
    title: '操作',
    right: true
  }]

export const ExternalAppListTableBody:EO_TBODY_TYPE[] = [

  {
    key: 'name',
    copy: true
  },
  {
    key: 'id',
    copy: true
  },
  {
    key: 'token',
    copy: true
  },
  {
    key: 'tags'
  },
  {
    key: 'status'
  },
  {
    key: 'operator'
  },
  {
    key: 'updateTime'
  },
  {
    type: 'btn',
    right: true,
    btns: [{
      title: '更新鉴权Token'
    },
    {
      title: '复制Token'
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
