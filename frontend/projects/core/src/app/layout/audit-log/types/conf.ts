import { SelectOption } from 'eo-ng-select'
import { THEAD_TYPE } from 'eo-ng-table'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

export const auditQueryStatusTypeList:SelectOption[] = [
  { label: '新建', value: 'create' },
  { label: '编辑', value: 'edit' },
  { label: '删除', value: 'delete' },
  { label: '发布', value: 'publish' }
]
export const auditLogsTableHeadName: THEAD_TYPE[] = [
  {
    title: '用户名',
    resizeable: true
  },
  {
    title: '操作类型',
    resizeable: true
  },
  {
    title: '操作对象',
    resizeable: true
  },
  {
    title: '操作时间',
    resizeable: true
  },
  {
    title: '操作IP',
    resizeable: true
  },
  {
    title: '操作',
    right: true
  }
]

export const auditLogsTableBody: EO_TBODY_TYPE[] = [
  {
    key: 'username',
    copy: true
  },
  {
    key: 'operateType',
    copy: true
  },
  {
    key: 'kind',
    copy: true
  },
  { key: 'time' },
  {
    key: 'ip',
    copy: true
  },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '查看'
      }
    ]
  }
]

export const auditLogDetailTableHeadName: THEAD_TYPE[] = [
  {
    title: '属性',
    resizeable: true
  },
  { title: '配置' }
]

export const auditLogDetailTableBody: EO_TBODY_TYPE[] = [
  {
    key: 'attr',
    copy: true
  },
  {
    key: 'value',
    styleFn: (item:any) => {
      if (item.attr === '请求内容') {
        return 'white-space: pre-wrap;word-wrap:break-word; word-break:break-all'
      } else {
        return 'white-space: unset;word-wrap:break-word; word-break:break-all'
      }
    },
    ellipsis: false,
    copy: true,
    json: true
  }
]
