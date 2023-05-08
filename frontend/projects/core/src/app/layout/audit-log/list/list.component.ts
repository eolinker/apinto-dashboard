/* eslint-disable dot-notation */
/* eslint-disable camelcase */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE } from '../../../constant/app.config'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'
import { AuditLogDetailComponent } from '../detail/detail.component'
import { auditLogsTableBody, auditLogsTableHeadName, auditQueryStatusTypeList } from '../types/conf'
import { AuditLogDetail, AuditLogsData } from '../types/types'

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

  logsList:AuditLogsData[] = []
  logsTableHeadName: THEAD_TYPE[] = [...auditLogsTableHeadName]
  logsTableBody: TBODY_TYPE[] =[...auditLogsTableBody]
  auditLogDetail:AuditLogDetail[]= []
  date:Array<Date> = [];
  listOfType:SelectOption[] = [...auditQueryStatusTypeList]
  listOfKind:SelectOption[] = []
  drawerRef: NzModalRef | undefined
  searchData:{keyword:string, operateType:string, kind:string, start:Date|null, end:Date|null, pageSize:number, pageNum:number, total:number, [key:string]:any} = {
    keyword: '',
    operateType: '',
    kind: '',
    start: null,
    end: null,
    pageSize: 20,
    pageNum: 1,
    total: 0
  }

  constructor (
    private message: EoNgFeedbackMessageService,
     private api:ApiService,
     private modalService: EoNgFeedbackModalService,
     private navigationService:EoNgNavigationService) {
    this.navigationService.reqFlashBreadcrumb([{ title: 'Debug日志' }])
  }

  ngOnInit (): void {
    this.logsTableBody[5].btns[0].click = (item:any) => {
      this.openDrawer(item.data.id)
    }
    this.getLogList()
    this.getTargetList()
  }

  getLogList ():void {
    const { start, end, kind, operateType, ...body } = this.searchData
    if (start) {
      body['start'] = Math.floor(start.getTime() / 1000)
    }
    if (end) {
      body['end'] = Math.floor(end.getTime() / 1000)
    } else {
      body['end'] = Math.floor(new Date().getTime() / 1000)
    }
    body['kind'] = kind || ''
    body['operateType'] = operateType || ''

    this.api.get('audit-logs', body)
      .subscribe((resp:{code:number, data:{items:AuditLogsData[], total:number}, msg:string}) => {
        if (resp.code === 0) {
          for (const index in resp.data.items) {
            resp.data.items[index]['username'] = resp.data.items[index].operator.username
            resp.data.items[index].operateType = this.changeToChinese(resp.data.items[index].operateType)
          }
          this.logsList = resp.data.items
          this.searchData.total = resp.data.total
        }
      })
  }

  getTargetList ():void {
    this.api.get('audit-log/kinds')
      .subscribe((resp:{code:number, data:{items:Array<{title:string, name:string}>}, msg:string}) => {
        if (resp.code === 0) {
          this.listOfKind = []
          for (const index in resp.data.items) {
            this.listOfKind.push({ label: resp.data.items[index].title, value: resp.data.items[index].name })
          }
        }
      })
  }

  openDrawer (auditLogId:string):void {
    this.drawerRef = this.modalService.create({
      nzTitle: '日志详情',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: AuditLogDetailComponent,
      nzComponentParams: {
        auditLogId: auditLogId
      },
      nzWrapClassName: 'audit-log-drawer',
      nzOkText: null
    })
  }

  // 重置搜索
  clearSearch ():void {
    this.searchData = {
      keyword: '',
      operateType: '',
      kind: '',
      start: null,
      end: null,
      pageSize: 20,
      pageNum: 1,
      total: 0
    }
    this.date = []
  }

  logsTableClick = (item:any) => {
    this.openDrawer(item.data.id)
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
