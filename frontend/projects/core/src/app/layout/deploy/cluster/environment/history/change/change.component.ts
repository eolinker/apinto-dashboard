import { Component, Input, OnInit, TemplateRef } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { DeployClusterOperateRecordTbody, DeployClusterOperateRecordThead } from '../../../types/conf'

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
      [(nzPageIndex)]="operateRecordsPage.pageNum"
      [(nzPageSize)]="operateRecordsPage.pageSize"
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
  operateRecordsData:{historys:Array<{key:string, oldValue:string, newValue:string, createTime:string, optType:string}>, total:number}=
      {
        historys: [],
        total: 0
      }

  operateRecordTabelHeadName: THEAD_TYPE[] = [...DeployClusterOperateRecordThead]
  operateRecordTableBody:TBODY_TYPE[]=[...DeployClusterOperateRecordTbody]

  // 更改历史分页
  // eslint-disable-next-line camelcase
  operateRecordsPage:{pageNum:number, pageSize:number, total:number}={
    pageNum: 1,
    pageSize: 15,
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
    this.api.get('cluster/' + this.clusterName + '/variable/update-history', { pageNum: this.operateRecordsPage.pageNum, pageSize: this.operateRecordsPage.pageSize }).subscribe(resp => {
      if (resp.code === 0) {
        this.operateRecordsData = resp.data
        this.operateRecordsPage.total = resp.data.total
      }
    })
  }
}
