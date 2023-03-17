/* eslint-disable dot-notation */
/* eslint-disable camelcase */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE } from '../../../constant/app.config'
import { Operator } from '../../../constant/type'
import { ApiService } from '../../../service/api.service'
import { AppConfigService } from '../../../service/app-config.service'
import { AuditLogDetailComponent } from '../detail/detail.component'

interface logsData{
  id:number,
  operator:Operator,
  operate_type: string,
  kind:string,
  time:string
  ip:string,
  [k:string]:any
}

@Component({
  selector: 'eo-ng-audit-log-list',
  templateUrl: './list.component.html',
  styles: [
    `
    .group-search-large{
      white-space:nowrap;
    }

    .group-search-large:nth-child(2){
      margin-top:var(--LAYOUT_PADDING);
      display: flex;
      align-items: center;}


     eo-ng-select.ant-select,
    eo-ng-select-top-control.ant-select-selector,
    nz-range-picker,
    .group-search-large eo-ng-input-group{
      width:254px !important;
    border-radius: var(--border-radius);
      min-height:32px;
    }

    `
  ]
})
export class AuditLogListComponent implements OnInit {
  @ViewChild('auditLogFormTpl', { read: TemplateRef, static: true }) auditLogFormTpl: TemplateRef<any> | undefined

  searchData:{keyword:string, operate_type:string, kind:string, start:Date|null, end:Date|null, page_size:number, page_num:number, total:number, [key:string]:any} = {
    keyword: '',
    operate_type: '',
    kind: '',
    start: null,
    end: null,
    page_size: 20,
    page_num: 1,
    total: 0
  }

  logsList:logsData[] = []
  logsTableHeadName: Array<object> = [
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

  logsTableBody: Array<any> =[
    { key: 'username' },
    { key: 'operate_type' },
    { key: 'kind' },
    { key: 'time' },
    { key: 'ip' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '查看',
          click: (item:any) => {
            this.getLogDetail(item.data)
          }
        }
      ]
    }
  ]

  auditLogDetail:Array<{attr:string, value:string}> = []

  date:Array<any> = [];

  listOfType:Array<{label:string, value:string}> = [
    { label: '新建', value: 'create' },
    { label: '编辑', value: 'edit' },
    { label: '删除', value: 'delete' },
    { label: '发布', value: 'publish' }
  ]

  listOfKind:Array<{label:string, value:string, [k:string]:any}> = []

  drawerRef: NzModalRef | undefined

  constructor (private message: EoNgFeedbackMessageService,
     private api:ApiService,
     private modalService: EoNgFeedbackModalService,
     private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '审计日志' }])
  }

  ngOnInit (): void {
    this.getLogList()
    this.getTargetList()
  }

  getLogList ():void {
    const { start, end, kind, operate_type, ...body } = this.searchData
    if (start) {
      body['start'] = Math.floor(start.getTime() / 1000)
    }
    if (end) {
      body['end'] = Math.floor(end.getTime() / 1000)
    } else {
      body['end'] = Math.floor(new Date().getTime() / 1000)
    }
    body['kind'] = kind || ''
    body['operate_type'] = operate_type || ''

    this.api.get('audit-logs', body).subscribe((resp:{code:number, data:{items:logsData[], total:number}, msg:string}) => {
      if (resp.code === 0) {
        for (const index in resp.data.items) {
          resp.data.items[index]['username'] = resp.data.items[index].operator.username
          resp.data.items[index].operate_type = this.changeToChinese(resp.data.items[index].operate_type)
        }
        this.logsList = resp.data.items
        this.searchData.total = resp.data.total
      } else {
        this.message.error(resp.msg || '获取列表数据失败！')
      }
    })
  }

  getTargetList ():void {
    this.api.get('audit-log/kinds').subscribe((resp:{code:number, data:{items:Array<any>}, msg:string}) => {
      if (resp.code === 0) {
        for (const index in resp.data.items) {
          resp.data.items[index]['label'] = resp.data.items[index].title
          resp.data.items[index]['value'] = resp.data.items[index].name
        }
        this.listOfKind = resp.data.items as Array<{label:string, value:string, [k:string]:any}>
      } else {
        this.message.error(resp.msg || '获取操作对象数据失败！')
      }
    })
  }

  // 接口返回成功才打开弹窗
  getLogDetail (item:any):void {
    this.api.get('audit-log', { log_id: item.id }).subscribe((resp:{code:number, data:{args:Array<{attr:string, value:string}>}, msg:string}) => {
      if (resp.code === 0) {
        this.auditLogDetail = resp.data.args
        for (const index in this.auditLogDetail) {
          if (this.auditLogDetail[index].attr === '请求内容') {
            this.auditLogDetail[index].value = JSON.stringify(JSON.parse(this.auditLogDetail[index].value), null, 4)
          }
        }
        this.openDrawer()
      } else {
        this.message.error(resp.msg || '获取日志详情失败！')
      }
    })
  }

  openDrawer ():void {
    this.drawerRef = this.modalService.create({
      nzTitle: '日志详情',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: AuditLogDetailComponent,
      nzComponentParams: {
        auditLogDetail: this.auditLogDetail
      },
      nzWrapClassName: 'audit-log-drawer',
      nzOkText: null
    })
  }

  // 重置搜索
  clearSearch ():void {
    this.searchData = {
      keyword: '',
      operate_type: '',
      kind: '',
      start: null,
      end: null,
      page_size: 20,
      page_num: 1,
      total: 0
    }
    this.date = []
  }

  logsTableClick = (item:any) => {
    this.getLogDetail(item.data)
  }

  onDateRangeChange (result: Date[]): void {
    if (result) {
      this.searchData.start = result[0]
      this.searchData.end = result[1]
    }
  }

  changeToChinese (value:string):string {
    switch (value) {
      case 'create':
        return '新建'
      case 'edit':
        return '编辑'
      case 'delete':
        return '删除'
      case 'publish':
        return '发布'
      default:
        return '编辑'
    }
  }
}
