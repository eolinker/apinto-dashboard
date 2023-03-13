import { Component, Input, OnInit, TemplateRef } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-deploy-cluster-environment-history-change',
  template: `
  <div class="drawer-list-content drawer-table  pb-[4px]" style="margin-top: 0px">
    <eo-ng-apinto-table
      class="drawer-table"
      [nzTbody]="operateRecordTableBody"
      [nzThead]="operateRecordTabelHeadName"
      [(nzData)]="operateRecordsData.historys"
      [nzNoScroll]="true"
      nzTableLayout="fixed"
      [nzShowPagination]="true"
      [nzFrontPagination]="false"
      [nzTotal]="operateRecordsPage.total"
      [(nzPageIndex)]="operateRecordsPage.page_num"
      [(nzPageSize)]="operateRecordsPage.page_size"
      [nzPageSizeOptions]="pageSizeOptions"
      (nzPageIndexChange)="getOperateRecords()"
      (nzPageSizeChange)="getOperateRecords()"
    >
    </eo-ng-apinto-table>
  </div>
  `,
  styles: [
  ]
})
export class DeployClusterEnvironmentHistoryChangeComponent implements OnInit {
  @Input() publishTypeTpl: TemplateRef<any> | undefined
  // eslint-disable-next-line camelcase
  operateRecordsData:{historys:Array<{key:string, old_value:string, new_value:string, create_time:string, opt_type:string}>, total:number}=
      {
        historys: [],
        total: 0
      }

  operateRecordTabelHeadName: Array<object> = [
    { title: 'KEY', resizeable: true },
    { title: 'OLD VALUE', resizeable: true },
    { title: 'NEW VALUE', resizeable: true },
    { title: '类型', resizeable: true },
    { title: '操作时间' }
  ]

  operateRecordTableBody: Array<any> =[
    { key: 'key' },
    { key: 'old_value' },
    { key: 'new_value' },
    { key: 'opt_type' },
    { key: 'create_time' }
  ]

  // 更改历史分页
  // eslint-disable-next-line camelcase
  operateRecordsPage:{page_num:number, page_size:number, total:number}={
    page_num: 1,
    page_size: 15,
    total: 0
  }

  pageSizeOptions:Array<number>=[15, 20, 50, 100]
  clusterName:string = ''
  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService) { }

  ngOnInit (): void {
    this.getOperateRecords()
  }

  ngAfterViewInit () {
    this.operateRecordTableBody[3].title = this.publishTypeTpl
  }

  getOperateRecords () {
    this.api.get('cluster/' + this.clusterName + '/variable/update-history', { page_num: this.operateRecordsPage.page_num, page_size: this.operateRecordsPage.page_size }).subscribe(resp => {
      if (resp.code === 0) {
        this.operateRecordsData = resp.data
        this.operateRecordsPage.total = resp.data.total
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }
}
