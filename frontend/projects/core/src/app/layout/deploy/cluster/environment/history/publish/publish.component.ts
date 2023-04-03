import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { DeployClusterOperateRecordTbody, DeployClusterOperateRecordThead, DeployClusterPublishRecordTbody, DeployClusterPublishRecordThead } from '../../../types/conf'

@Component({
  selector: 'eo-ng-deploy-cluster-environment-history-publish',
  template: `
  <div class="drawer-list-content drawer-table" style="margin-top: 0px">
    <eo-ng-apinto-table
      class="drawer-table"
      [nzTbody]="publishRecordTableBody"
      [nzThead]="publishRecordTabelHeadName"
      [(nzData)]="publishRecordsData.historys"
      [nzTrBottomTmp]="nzTrBottomTmp"
      [nzTrClick]="publishRTableClick"
      [nzNoScroll]="true"
      [nzShowPagination]="true"
      [nzFrontPagination]="false"
      [nzTotal]="publishRecordsPage.total"
      [(nzPageIndex)]="publishRecordsPage.pageNum"
      [(nzPageSize)]="publishRecordsPage.pageSize"
      [nzPageSizeOptions]="pageSizeOptions"
      (nzPageIndexChange)="getPublishRecords()"
      (nzPageSizeChange)="getPublishRecords()"
    >
    </eo-ng-apinto-table>
  </div>
  <ng-template #nzTrBottomTmp let-item="item" let-apis="apis" let-index="index">
    <tr *ngIf="item.data.isExpand">
      <td style="border-right: none; background-color: #e8e8e8"></td>
      <td
        colspan="3"
        style="
          background-color: pink;
          border: none;
          padding: 0;
          margin-top: -2px;
        "
      >
        <div style="width: 100%">
          <eo-ng-apinto-table
            class="floatR innerTable"
            [nzTbody]="publishRecordDetailsTableBody"
            [nzThead]="publishRecordDetailsTabelHeadName"
            [(nzData)]="item.data.details"
            nzTableLayout="fixed"
            [nzScrollY]="9999"
          >
          </eo-ng-apinto-table>
        </div>
      </td>
    </tr>
  </ng-template>


<ng-template #showInnerTableBtnTpl let-item="item">
<button eo-ng-row-expand-button [expand]="item.isExpand"></button>
</ng-template>
  `,
  styles: [
  ]
})
export class DeployClusterEnvironmentHistoryPublishComponent implements OnInit {
  @ViewChild('showInnerTableBtnTpl', { read: TemplateRef, static: true }) showInnerTableBtnTpl: TemplateRef<any> | undefined
  @Input() publishTypeTpl: TemplateRef<any> | undefined

  // eslint-disable-next-line camelcase
  publishRecordsData:{historys:Array<{id:number, name:string, createTime:string, optType:string, operator:string, detailShow:boolean, detail:Array<{key:string, oldValue:string, newValue:string, optType:string, createTime:string}>}>, total:number}=
      {
        historys: [],
        total: 0
      }

  publishRecordTabelHeadName: THEAD_TYPE[]= [...DeployClusterPublishRecordThead]
  publishRecordTableBody: TBODY_TYPE[]=[...DeployClusterPublishRecordTbody]

  publishRecordDetailsTabelHeadName: THEAD_TYPE[] = [...DeployClusterOperateRecordThead]
  publishRecordDetailsTableBody: TBODY_TYPE[] = [...DeployClusterOperateRecordTbody]
  pageSizeOptions:Array<number>=[15, 20, 50, 100]

  // 发布历史分页
  // eslint-disable-next-line camelcase
  publishRecordsPage:{pageNum:number, pageSize:number, total:number}={
    pageNum: 1,
    pageSize: 15,
    total: 0
  }

  clusterName:string = ''
  constructor (private message: EoNgFeedbackMessageService,
    private api:ApiService) {
  }

  ngOnInit (): void {
    this.getPublishRecords()
  }

  ngAfterViewInit () {
    this.publishRecordTableBody[0].title = this.showInnerTableBtnTpl
    this.publishRecordDetailsTableBody[3].title = this.publishTypeTpl
  }

  getPublishRecords () {
    this.api.get('cluster/' + this.clusterName + '/variable/publish-history', { pageNum: this.publishRecordsPage.pageNum, pageSize: this.publishRecordsPage.pageSize }).subscribe(resp => {
      if (resp.code === 0) {
        this.publishRecordsData = resp.data
        for (const history of resp.data.historys) {
          history.isExpand = false
        }
        this.publishRecordsPage.total = resp.data.total
      }
    })
  }

  expandChange () {
    this.publishRecordsData.historys = [...this.publishRecordsData.historys]
  }

  publishRTableClick = (item:any) => {
    item.data.isExpand = !item.data.isExpand
  }
}
